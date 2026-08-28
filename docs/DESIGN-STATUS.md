# Design Status

Updated: 2026-08-28.

This file is the short handoff for continuing framework design in a new session. Durable principles live in `FOUNDATIONS.md`; host-hook research lives in `HOST-HOOKS.md`.

## Current product definition

AI Project Framework is a generic empty project used as the starting environment for future software projects.

The coding AI is the worker. The framework should make a weaker but capable AI behave more like a disciplined developer by giving it small, concrete instructions at the right moments and by using mechanisms that do not depend on model memory or obedience.

The user-facing project remains static and must not require an APF executable/runtime.

## Core idea now established

### 1. Recursive modularity is the code core

The project is a tree of building blocks:

```text
project
└── module
    ├── submodule
    │   └── smaller module
    └── submodule
```

Large modules contain smaller modules; lower levels should become simple enough for an AI to understand and modify locally. Narrow boundaries and pure deterministic leaf logic are preferred.

### 2. The framework must guide at key moments

A giant permanent rule prompt is not sufficient.

Main control idea:

```text
key event
  -> small instruction
  -> evidence
  -> gate
```

Example: if planning is required before coding, do not merely write "plan first" in `AGENTS.md`. Intercept the first code-changing action, require the preparation evidence, block if it is missing, and tell the AI exactly what to do next.

### 3. Complex developer behaviors should be decomposed

Do not tell a weaker AI only "review the code". Break review into simple actions such as:

- check duplicate authoritative facts,
- check module-boundary leaks,
- check pure logic versus external effects,
- check affected tests,
- check authoritative documentation.

The framework should progressively encode expert development behavior as small executable steps.

### 4. Mainstream coding hosts are converging on hooks

Research currently supports a common lifecycle around concepts equivalent to:

```text
SessionStart
UserPromptSubmit
PreToolUse
PostToolUse
Stop
PreCompact
```

Host-specific adapters should translate these into one framework lifecycle. See `HOST-HOOKS.md`.

## Important corrections from earlier design

- Do not make Structure/Effect/Knowledge Graphs the product core. They remain research models only.
- Do not ship an APF runtime, executable, daemon, or scanner into consumer projects.
- Do not use manually maintained context-routing tables as the default solution.
- Do not use giant `AGENTS.md` / rule manuals as the main control mechanism.
- Do not require multi-model review.
- Do not optimize for "minimal verification"; optimize for sufficient verification of affected behavior.
- Do not use GitHub Issues as a duplicate internal architecture/task database.

## Documentation direction already agreed

The documentation system remains a major framework capability:

- `AGENTS.md` should stay small,
- detailed knowledge should load on demand,
- module-local knowledge should live near its module,
- one fact should have one authoritative source,
- documentation should stay aligned with code,
- models should leave useful repository context as work progresses/completes.

A lightweight document map may still be useful, but the exact design is not decided. Do not recreate a large hand-maintained task-to-document router.

The exact module `README.md` contract is also not finalized, though module-local README/documentation remains a strong direction.

## Verification direction already agreed

Tests should be organized around modules rather than one unrelated global test pile. Small functions may have focused tests inside the module's test suite.

Mechanically observable rules should become project-native checks/tests/CI where practical. Hooks provide timing and early feedback; they are not the sole enforcement layer.

## Current design frontier

Do not jump back into scanner/graph implementation yet.

The next design work should continue from the lifecycle insight:

1. decide the smallest useful framework lifecycle stages,
2. decide what instruction/evidence/gate belongs at each stage,
3. then connect documentation, module structure, planning, review, testing, skills/agents, and GitHub mechanisms to those stages,
4. only after that decide the consumer static-file layout and host adapter files.

Do not prematurely freeze `.apf`, `AGENTS.md`, hook configuration, or module README schemas before the above behavior is clear.
