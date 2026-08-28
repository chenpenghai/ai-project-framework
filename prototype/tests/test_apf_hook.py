from __future__ import annotations

import importlib.util
import tempfile
import unittest
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]
SPEC = importlib.util.spec_from_file_location("apf_hook", ROOT / "hooks" / "apf_hook.py")
assert SPEC is not None and SPEC.loader is not None
APF = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(APF)


class APFHookTests(unittest.TestCase):
    def setUp(self) -> None:
        self.temp = tempfile.TemporaryDirectory()
        self.data = Path(self.temp.name)
        self.base = {
            "session_id": "session-1",
            "turn_id": "turn-1",
        }

    def tearDown(self) -> None:
        self.temp.cleanup()

    def test_user_prompt_adds_small_preparation_context(self) -> None:
        output = APF.dispatch(
            {**self.base, "hook_event_name": "UserPromptSubmit"}, ROOT, self.data
        )
        self.assertEqual(
            output["hookSpecificOutput"]["hookEventName"], "UserPromptSubmit"
        )
        self.assertIn("smallest affected area", output["hookSpecificOutput"]["additionalContext"])

    def test_first_edit_is_paused_once_per_turn(self) -> None:
        payload = {**self.base, "hook_event_name": "PreToolUse"}
        first = APF.dispatch(payload, ROOT, self.data)
        second = APF.dispatch(payload, ROOT, self.data)

        self.assertEqual(
            first["hookSpecificOutput"]["permissionDecision"], "deny"
        )
        self.assertIn("BEFORE_CHANGE", first["hookSpecificOutput"]["permissionDecisionReason"])
        self.assertIsNone(second)

    def test_explicit_failure_adds_problem_guidance_once(self) -> None:
        payload = {
            **self.base,
            "hook_event_name": "PostToolUse",
            "tool_response": {"exit_code": 1},
        }
        first = APF.dispatch(payload, ROOT, self.data)
        second = APF.dispatch(payload, ROOT, self.data)

        self.assertEqual(first["hookSpecificOutput"]["hookEventName"], "PostToolUse")
        self.assertIn("ON_PROBLEM", first["hookSpecificOutput"]["additionalContext"])
        self.assertIsNone(second)

    def test_successful_tool_does_not_trigger_problem_guidance(self) -> None:
        output = APF.dispatch(
            {
                **self.base,
                "hook_event_name": "PostToolUse",
                "tool_response": {"exit_code": 0},
            },
            ROOT,
            self.data,
        )
        self.assertIsNone(output)

    def test_stop_blocks_once_then_allows_continuation(self) -> None:
        first = APF.dispatch(
            {**self.base, "hook_event_name": "Stop", "stop_hook_active": False},
            ROOT,
            self.data,
        )
        second = APF.dispatch(
            {**self.base, "hook_event_name": "Stop", "stop_hook_active": True},
            ROOT,
            self.data,
        )

        self.assertEqual(first["decision"], "block")
        self.assertIn("BEFORE_FINISH", first["reason"])
        self.assertIsNone(second)


if __name__ == "__main__":
    unittest.main()
