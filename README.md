# AI Project Framework

AI Project Framework (APF) is a guidance layer for AI-assisted software development.

The AI remains the worker. APF intervenes at useful moments, turns broad engineering behavior into small executable actions, and helps projects improve while real work is being done.

## Product direction

APF is **not an empty-project template** and does not require a project to adopt a special directory layout before it can be used.

A new project and an existing project are both valid starting points.

Installing/enabling APF should not automatically create files, reorganize directories, or perform a framework migration. Project structure should change only when the current development task provides a good reason to improve it.

The intended product shape is:

```text
shared APF core
  + thin host adapters/plugins
  + the user's existing project
```

The shared core owns development guidance. Host adapters only translate lifecycle events and feedback formats for tools such as Claude Code, Codex, Cursor, Gemini CLI, CodeBuddy, and others.

## Core behavior

APF should intervene at a small number of high-value moments rather than permanently loading a large rule manual.

The first prototype focuses on three semantic events:

```text
BEFORE_CHANGE
  -> narrow the task, identify the affected area, split the work, choose verification

ON_PROBLEM
  -> stop blind retries, inspect evidence, form a small hypothesis, choose one diagnostic next step

BEFORE_FINISH
  -> review the actual change, run sufficient affected verification, continue if important work remains
```

See `prototype/` for the current host-neutral sample.

## How projects should improve

APF favors progressive improvement instead of up-front framework conversion.

When work touches an area, the AI should prefer leaving that area easier to understand and change than before, without expanding the task into unrelated cleanup.

Important code-shape principles include:

- recursive modularity: large capabilities are composed from smaller coherent modules;
- narrow boundaries: callers should not need to understand module internals;
- pure deterministic logic where practical at lower levels;
- explicit external effects;
- one authoritative source for each fact;
- module-oriented, affected verification.

These are directions for development decisions, not requirements to rewrite an existing repository on installation.

## Control model

Natural-language guidance is useful when delivered at the right time. Mechanically observable requirements should additionally use project-native mechanisms where available, such as tests, type checks, linters, architecture checks, Git hooks, and CI.

Host hooks/plugins provide early guidance. They are not the sole correctness layer.

## Repository status

The current goal is a small working prototype of timed guidance before expanding the framework.

The existing Go scanner/graph code remains research material and is not the product core.

Read:

- `docs/FOUNDATIONS.md` — durable constraints;
- `docs/DESIGN-STATUS.md` — current design state;
- `docs/HOST-HOOKS.md` — host integration direction;
- `prototype/README.md` — first small sample.
