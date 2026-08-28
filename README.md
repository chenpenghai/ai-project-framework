# AI Project Framework

A general-purpose development framework for AI-assisted software projects.

The framework exists to serve the project, not the other way around. Its purpose is to help any coding model produce better code faster, with less unnecessary context and less structural drift.

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

## Core direction

The first implementation will focus on a small kernel:

1. Repository and project discovery.
2. Structure graph: projects, modules, files, symbols, dependencies, and public boundaries.
3. Effect graph: pure / effectful / unknown code and effect boundaries.
4. Knowledge graph: authorities, modules, documentation, decisions, and tests.
5. Structural ratchet: prevent new architectural regressions without forcing cleanup of existing debt.
6. Incremental context and verification derived from the current change.

The framework may become internally powerful, but its active surface for any single task should remain small.
