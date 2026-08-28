# AGENTS.md

This file is a small entry point, not the project knowledge base.

Before making a substantial change, read `docs/FOUNDATIONS.md` and only the additional documents relevant to the current target.

Keep changes local and minimal. Prefer existing module boundaries. Prefer deterministic pure logic for business decisions and keep I/O, time, randomness, persistence, network access, and UI/runtime effects at explicit edges.

Do not add framework ceremony unless it measurably improves development speed, code quality, or maintainability.

Do not introduce dependencies on a specific programming language, operating system, coding assistant, model, or second-model review workflow.

When implementation begins, use the framework's own incremental checks rather than relying on these instructions alone.
