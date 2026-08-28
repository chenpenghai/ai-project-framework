#!/usr/bin/env python3
from __future__ import annotations

import json
import os
import re
import sys
import tempfile
from pathlib import Path
from typing import Any


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
        if payload.get("stop_hook_active") is True:
            return None
        return {
            "decision": "block",
            "reason": (
                "APF BEFORE_FINISH checkpoint. Complete this focused review before ending the turn:\n\n"
                + _guidance(root, "before-finish.md")
            ),
        }

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
