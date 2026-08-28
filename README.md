# AI Project Framework

A general-purpose development framework for AI-assisted software projects.

The framework exists to serve the project, not the other way around. Its purpose is to help any coding model produce better code faster, with less unnecessary context and less structural drift.

## Important: this repository is the framework source

This repository is where APF itself is developed, tested, and released. Users should not clone it and build their application inside this source tree.

Consumer projects are generated separately and start with no application code, no selected language, and no prescribed architecture.

The first consumer-project command is:

```bash
go run ./cmd/apf new /path/to/my-project
```

The generated project currently contains only the minimum framework control surface:

```text
my-project/
├── .apf/
│   └── project.yaml
├── .gitignore
└── AGENTS.md
```

Framework implementation source, tests, research documents, and internal tooling are never copied into the consumer project.

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

## Current prototype

The executable prototype currently has two deliberately small surfaces:

```bash
go run ./cmd/apf new /path/to/my-project
go run ./cmd/apf scan /path/to/repository
go run ./cmd/apf scan --json /path/to/repository
```

`apf new` creates an empty consumer project without application code or framework source.

`apf scan` builds the current Structure Graph fact layer. It detects repository files, Git changed files, common project/package manifests, explicit `MODULE.md` declarations, declared local dependencies, workspace governance relationships, deterministic containment, dependency cycles, and affected projects. Unsupported languages degrade to the universal filesystem/Git baseline instead of failing.

The Go implementation is an implementation detail of the framework tool itself; consumer projects do not need to use Go.

See `docs/STRUCTURE-GRAPH.md` for the current graph semantics and `docs/FOUNDATIONS.md` for durable design constraints.

## Core direction

The implementation will grow only as evidence justifies it:

1. Repository and project discovery.
2. Structure graph: projects, modules, files, symbols, dependencies, and public boundaries.
3. Effect graph: pure / effectful / unknown code and effect boundaries.
4. Knowledge graph: authorities, modules, documentation, decisions, and tests.
5. Structural ratchet: prevent new architectural regressions without forcing cleanup of existing debt.
6. Incremental context and verification derived from the current change.

The framework may become internally powerful, but its active surface for any single task should remain small.
