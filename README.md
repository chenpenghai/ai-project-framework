# AI Project Framework

A general-purpose static project framework for AI-assisted software development.

The framework exists to serve the user's project, not the other way around. Its purpose is to help any coding model produce better code faster, with less unnecessary context and less structural drift.

## User-facing product

The user-facing artifact is the `empty-project/` directory in this repository.

Copy that directory, rename it for the new project, open it with a coding model, and start describing the product to build.

The initial project is intentionally static and empty:

```text
empty-project/
├── AGENTS.md
└── .apf/
    └── project.yaml
```

It contains no application code, no prescribed programming language, no prescribed architecture, no executable, and no framework runtime.

The coding model uses the static project guidance together with the project's own language, build system, tests, linters, compiler, Git, and CI as those mechanisms emerge.

## Framework source repository

Everything outside `empty-project/` exists to design, test, and maintain the framework itself. Framework development source, experiments, tests, and documentation are not part of a user's project.

The current Go scanner and Structure Graph code are framework-development prototypes used to test architectural ideas. They are not required by projects copied from `empty-project/`.

## Goals

- Keep normal development fast; framework overhead must stay small.
- Work with any programming language, platform, coding assistant, and model.
- Require no second model or external AI reviewer.
- Prefer modular code with clear boundaries.
- Prefer pure, deterministic business logic with side effects pushed to the edges.
- Reduce coupling and duplicate sources of truth.
- Give the model only the context relevant to the current change.
- Verify only what the current change can affect whenever possible.
- Automate framework behavior so users do not need to learn its internals.

See `docs/FOUNDATIONS.md` for durable design constraints.
