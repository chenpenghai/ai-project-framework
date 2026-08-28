# Design Status

Updated: 2026-08-28.

## Current product definition

AI Project Framework (APF) is a plugin-oriented guidance layer for AI-assisted development.

It is not an empty-project template. It should work on both new and existing projects without requiring a migration or mandatory project structure.

The AI is the worker. APF intervenes at useful moments, narrows the current problem, and breaks expert development behavior into small executable actions.

## Current architecture direction

```text
shared APF core
        ↓
thin host adapters/plugins
        ↓
existing coding assistant
        ↓
user's current project
```

The shared core owns development guidance. Host adapters only translate host lifecycle events and feedback semantics.

Do not maintain separate rule sets for Claude, Codex, Cursor, Gemini, CodeBuddy, or other hosts.

## Existing-project principle

Enabling APF should do nothing destructive or ceremonial to the repository.

Do not automatically:

- generate `.apf/`;
- create `AGENTS.md`;
- reorganize modules;
- create README files everywhere;
- migrate tests;
- rewrite architecture.

Instead, when real work touches an area, APF should guide the AI toward a cleaner local result without unnecessarily widening task scope.

## First prototype

Stop expanding the design before testing the central idea.

The first host-neutral sample uses three semantic events:

```text
BEFORE_CHANGE
ON_PROBLEM
BEFORE_FINISH
```

### BEFORE_CHANGE

Before meaningful code modification, help the AI:

1. narrow the user's goal;
2. identify the smallest affected area/module;
3. split the work into small steps;
4. decide how the result will be verified.

Do not require a long formal plan.

### ON_PROBLEM

When there is evidence of failure, repeated unsuccessful attempts, or a broken assumption:

1. stop blind retries;
2. state the observed failure;
3. inspect the smallest useful evidence;
4. form a narrow hypothesis;
5. choose one next diagnostic action.

The goal is to prevent uncontrolled scope expansion.

### BEFORE_FINISH

Before declaring completion:

1. compare the actual result with the user's goal;
2. review the changed area for obvious boundary, duplication, and effect problems;
3. run sufficient affected verification;
4. update persistent project knowledge only if the work changed it;
5. continue working if an important unresolved problem remains.

The current sample lives in `prototype/`.

## Code direction retained

Recursive modularity remains the preferred code shape: large capabilities should be composed from smaller coherent modules with narrow boundaries.

At lower levels, deterministic logic should prefer simple pure functions where practical and keep external effects explicit.

These are progressive development directions, not installation-time migration requirements.

## Important rejected directions

Do not return to these by default:

- static `empty-project/` as the product;
- mandatory APF project files;
- consumer APF executable/runtime/daemon;
- graph/scanner as product core;
- giant permanent `AGENTS.md` or rule manual;
- one separate framework rule set per coding assistant;
- mandatory multi-model review;
- framework-wide cleanup during unrelated tasks;
- GitHub Issues as duplicate internal state.

## Current design frontier

Build and test the small three-event prototype before adding more lifecycle stages or framework machinery.

The next implementation question is narrow:

> Which host should receive the first thin adapter, and can it reliably deliver the three shared guidance events without copying the guidance itself?

Vendor hook/plugin APIs must be checked against current official documentation when that adapter is implemented.
