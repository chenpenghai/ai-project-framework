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

    def _stop(self, message: str = ""):
        return APF.dispatch(
            {
                **self.base,
                "hook_event_name": "Stop",
                "last_assistant_message": message,
            },
            ROOT,
            self.data,
        )

    def _result(self, step: str, status: str = "done", evidence: str = "checked") -> str:
        return (
            "<APF_RESULT>\n"
            f'{{"flow":"before_finish","step":"{step}","status":"{status}",'
            f'"evidence":"{evidence}"}}\n'
            "</APF_RESULT>"
        )

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

        self.assertEqual(first["hookSpecificOutput"]["permissionDecision"], "deny")
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

    def test_finish_loop_starts_with_first_small_step(self) -> None:
        output = self._stop("normal candidate final answer")
        self.assertEqual(output["decision"], "block")
        self.assertIn("check_goal", output["reason"])
        self.assertIn("<APF_RESULT>", output["reason"])

    def test_invalid_result_does_not_advance(self) -> None:
        self._stop()
        output = self._stop("I checked it but forgot the protocol")
        self.assertEqual(output["decision"], "block")
        self.assertIn("same step", output["reason"])
        self.assertIn("check_goal", output["reason"])

    def test_wrong_step_does_not_advance(self) -> None:
        self._stop()
        output = self._stop(self._result("check_structure"))
        self.assertEqual(output["decision"], "block")
        self.assertIn("check_goal", output["reason"])

    def test_needs_work_keeps_current_step(self) -> None:
        self._stop()
        output = self._stop(
            self._result("check_goal", status="needs_work", evidence="one requested behavior is missing")
        )
        self.assertEqual(output["decision"], "block")
        self.assertIn("Resolve what you found", output["reason"])
        self.assertIn("check_goal", output["reason"])

    def test_done_advances_to_next_step(self) -> None:
        self._stop()
        output = self._stop(self._result("check_goal", evidence="goal covered"))
        self.assertEqual(output["decision"], "block")
        self.assertIn("check_structure", output["reason"])

    def test_all_finish_steps_complete_before_normal_final_is_allowed(self) -> None:
        self._stop()
        sequence = [
            "check_goal",
            "check_structure",
            "check_verification",
            "check_project_knowledge",
        ]

        for index, step in enumerate(sequence):
            output = self._stop(self._result(step, evidence=f"evidence for {step}"))
            self.assertEqual(output["decision"], "block")
            if index < len(sequence) - 1:
                self.assertIn(sequence[index + 1], output["reason"])
            else:
                self.assertIn("normal final answer", output["reason"])

        final = self._stop("normal user-facing final answer")
        self.assertIsNone(final)

    def test_parser_uses_last_result_marker(self) -> None:
        message = self._result("wrong") + "\ntext\n" + self._result("check_goal", evidence="right")
        result = APF._parse_apf_result(message)
        self.assertEqual(result["step"], "check_goal")
        self.assertEqual(result["evidence"], "right")


if __name__ == "__main__":
    unittest.main()
