# AGENTS.md

This repository develops AI Project Framework (APF).

Before substantial design or implementation work, read:

1. `docs/FOUNDATIONS.md` — durable constraints.
2. `docs/DESIGN-STATUS.md` — current direction and next experiment.
3. The focused document relevant to the task, especially `docs/HOST-HOOKS.md` for host integration.

Current product direction:

- APF is a plugin-oriented guidance layer, not an empty-project template.
- Existing projects are first-class; activation must not require migration or automatic restructuring.
- The shared core owns development guidance; host adapters only translate host events and feedback semantics.
- Start with the small three-event prototype in `prototype/` before adding lifecycle complexity.
- Recursive modularity, narrow boundaries, pure deterministic leaf logic, explicit effects, authoritative documentation, and affected verification remain preferred development directions.

Do not revive the consumer APF runtime, mandatory `.apf/` structure, giant permanent rule prompts, graph/scanner product core, or separate duplicated rule sets per coding assistant.
