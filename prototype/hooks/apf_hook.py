#!/usr/bin/env python3
from __future__ import annotations

import json
import os
import re
import sys
import tempfile
from pathlib import Path
from typing import Any


RESULT_PATTERN = re.compile(r"<APF_RESULT>\s*(\{.*?\})\s*</APF_RESULT>", re.DOTALL)
MAX_FINISH_BLOCKS = 12


def _plugin_root() -> Path:
    value = os.environ.get("PLUGIN_ROOT")
    if value:
        return Path(value)
    return Path(__file__).resolve().parents[1]


def _plugin_data() -> Path:
    value = os.environ.get("PLUGIN_DATA")
    if value:
        return Path(value)
    return Path(tempfile.gettempdir()) / "apf-prototype"


def _guidance(root: Path, name: str) -> str:
    path = root / "guidance" / name
    return path.read_text(encoding="utf-8").strip()


def _flow(root: Path, name: str) -> dict[str, Any]:
    path = root / "flows" / f"{name}.json"
    return json.loads(path.read_text(encoding="utf-8"))


def _safe_id(value: Any) -> str:
    text = str(value or "unknown")
    return re.sub(r"[^A-Za-z0-9_.-]+", "_", text)[:120]


def _state_path(data_root: Path, payload: dict[str, Any]) -> Path:
    session_id = _safe_id(payload.get("session_id"))
    turn_id = _safe_id(payload.get("turn_id"))
    return data_root / "turn-state" / f"{session_id}--{turn_id}.json"


def _load_state(path: Path) -> dict[str, Any]:
    try:
        return json.loads(path.read_text(encoding="utf-8"))
    except (FileNotFoundError, json.JSONDecodeError, OSError):
        return {}


def _save_state(path: Path, state: dict[str, Any]) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    tmp = path.with_suffix(".tmp")
    tmp.write_text(json.dumps(state, ensure_ascii=False), encoding="utf-8")
    tmp.replace(path)


def _explicit_failure(value: Any) -> bool:
    """Return True only for explicit structured failure signals."""
    if isinstance(value, dict):
        for key, item in value.items():
            normalized = str(key).replace("-", "_").lower()
            if normalized in {"is_error", "iserror", "failed"} and item is True:
                return True
            if normalized in {"success", "ok"} and item is False:
                return True
            if normalized in {"exit_code", "exitcode", "return_code", "returncode"}:
                if isinstance(item, int) and item != 0:
                    return True
            if normalized in {"status", "state"} and isinstance(item, str):
                if item.strip().lower() in {"error", "failed", "failure"}:
                    return True
            if _explicit_failure(item):
                return True
        return False
    if isinstance(value, list):
        return any(_explicit_failure(item) for item in value)
    return False


def _parse_apf_result(message: Any) -> dict[str, Any] | None:
    if not isinstance(message, str):
        return None
    matches = RESULT_PATTERN.findall(message)
    if not matches:
        return None
    try:
        value = json.loads(matches[-1])
    except json.JSONDecodeError:
        return None
    return value if isinstance(value, dict) else None


def _result_is_valid(result: dict[str, Any] | None, flow_id: str, step_id: str) -> bool:
    if not result:
        return False
    if result.get("flow") != flow_id or result.get("step") != step_id:
        return False
    if result.get("status") not in {"done", "needs_work"}:
        return False
    evidence = result.get("evidence")
    return isinstance(evidence, str) and bool(evidence.strip())


def _protocol(flow_id: str, step_id: str) -> str:
    return (
        "When this step is complete, end your response with exactly one structured result:\n"
        "<APF_RESULT>\n"
        f'{{"flow":"{flow_id}","step":"{step_id}","status":"done","evidence":"brief concrete evidence"}}\n'
        "</APF_RESULT>\n"
        "Use status \"needs_work\" instead of \"done\" if this check finds unresolved work."
    )


def _step_prompt(flow: dict[str, Any], step: dict[str, Any], prefix: str = "") -> str:
    parts = []
    if prefix:
        parts.append(prefix.strip())
    parts.extend(
        [
            f"APF {flow['event']} step: {step['id']}",
            step["instruction"],
            "Required evidence: " + step["evidence"],
            _protocol(flow["id"], step["id"]),
        ]
    )
    return "\n\n".join(parts)


def _stop_loop(payload: dict[str, Any], root: Path, data_root: Path) -> dict[str, Any] | None:
    flow = _flow(root, "before-finish")
    steps = flow["steps"]
    state_path = _state_path(data_root, payload)
    state = _load_state(state_path)
    finish = state.setdefault("before_finish", {})

    if finish.get("ready_for_final"):
        finish["completed"] = True
        _save_state(state_path, state)
        return None

    step_index = int(finish.get("step_index", 0))
    blocks = int(finish.get("blocks", 0))

    if blocks >= MAX_FINISH_BLOCKS:
        finish["ready_for_final"] = True
        finish["safety_limit_reached"] = True
        _save_state(state_path, state)
        return {
            "decision": "block",
            "reason": (
                "APF loop safety limit reached. Stop the APF review loop and now give the user "
                "a normal final answer. Do not include an APF_RESULT marker."
            ),
        }

    if not finish.get("started"):
        finish.update({"started": True, "step_index": 0, "blocks": blocks + 1})
        _save_state(state_path, state)
        return {"decision": "block", "reason": _step_prompt(flow, steps[0])}

    current = steps[step_index]
    result = _parse_apf_result(payload.get("last_assistant_message"))

    if not _result_is_valid(result, flow["id"], current["id"]):
        finish["blocks"] = blocks + 1
        _save_state(state_path, state)
        return {
            "decision": "block",
            "reason": _step_prompt(
                flow,
                current,
                "The previous response did not contain a valid APF_RESULT for the current step. "
                "Complete the same step and return the required structured result.",
            ),
        }

    if result["status"] == "needs_work":
        finish["blocks"] = blocks + 1
        _save_state(state_path, state)
        return {
            "decision": "block",
            "reason": _step_prompt(
                flow,
                current,
                "This step reported unresolved work. Resolve what you found before marking this step done. "
                f"Reported evidence: {result['evidence']}",
            ),
        }

    next_index = step_index + 1
    finish["step_index"] = next_index
    finish["blocks"] = blocks + 1
    finish.setdefault("evidence", {})[current["id"]] = result["evidence"]

    if next_index >= len(steps):
        finish["ready_for_final"] = True
        _save_state(state_path, state)
        return {
            "decision": "block",
            "reason": (
                "APF BEFORE_FINISH review is complete. Now give the user your normal final answer. "
                "Do not include any APF_RESULT marker or internal APF protocol details."
            ),
        }

    _save_state(state_path, state)
    return {"decision": "block", "reason": _step_prompt(flow, steps[next_index])}


def dispatch(payload: dict[str, Any], root: Path, data_root: Path) -> dict[str, Any] | None:
    event = payload.get("hook_event_name")

    if event == "UserPromptSubmit":
        context = (
            "APF is active. Before the first meaningful code edit in this turn, "
            "identify the smallest affected area, split the work into a few small steps, "
            "and choose how the changed behavior will be verified. Keep this preparation short."
        )
        return {
            "hookSpecificOutput": {
                "hookEventName": "UserPromptSubmit",
                "additionalContext": context,
            }
        }

    if event == "PreToolUse":
        state_path = _state_path(data_root, payload)
        state = _load_state(state_path)
        if state.get("before_change_checkpoint_seen"):
            return None

        state["before_change_checkpoint_seen"] = True
        _save_state(state_path, state)
        reason = (
            "APF BEFORE_CHANGE checkpoint. Before retrying this edit, briefly complete this preparation:\n\n"
            + _guidance(root, "before-change.md")
        )
        return {
            "hookSpecificOutput": {
                "hookEventName": "PreToolUse",
                "permissionDecision": "deny",
                "permissionDecisionReason": reason,
            }
        }

    if event == "PostToolUse":
        if not _explicit_failure(payload.get("tool_response")):
            return None

        state_path = _state_path(data_root, payload)
        state = _load_state(state_path)
        if state.get("on_problem_guidance_seen"):
            return None

        state["on_problem_guidance_seen"] = True
        _save_state(state_path, state)
        return {
            "hookSpecificOutput": {
                "hookEventName": "PostToolUse",
                "additionalContext": _guidance(root, "on-problem.md"),
            }
        }

    if event == "Stop":
        return _stop_loop(payload, root, data_root)

    return None


def main() -> int:
    try:
        payload = json.load(sys.stdin)
        output = dispatch(payload, _plugin_root(), _plugin_data())
        if output is not None:
            json.dump(output, sys.stdout, ensure_ascii=False)
            sys.stdout.write("\n")
        return 0
    except Exception as exc:  # Prototype should fail open rather than break the host.
        print(f"APF prototype hook error: {exc}", file=sys.stderr)
        return 0


if __name__ == "__main__":
    raise SystemExit(main())
