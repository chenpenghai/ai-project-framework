# BEFORE_FINISH

Before declaring the task complete:

1. Compare the actual result with the user's requested goal.
2. Review the changed area for obvious boundary leaks, duplicated authoritative logic/data, and unnecessary mixing of deterministic logic with external effects.
3. Run sufficient verification for the changed behavior and reasonably affected dependencies.
4. Update persistent project documentation only if authoritative project knowledge changed.
5. If an important problem or failed relevant check remains, continue working instead of claiming completion.

Keep the review focused on the work actually affected by this task.
