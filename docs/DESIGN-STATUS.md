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

The first prototype uses three semantic events:

```text
BEFORE_CHANGE
ON_PROBLEM
BEFORE_FINISH
```

### BEFORE_CHANGE

Before meaningful code modification, help the AI narrow the goal, identify the smallest affected area, split the work into small steps, and choose verification. Do not require a long formal plan.

### ON_PROBLEM

When there is explicit failure evidence, stop blind retries, state the observed failure, inspect the smallest useful evidence, form one narrow hypothesis, and choose one diagnostic action.

### BEFORE_FINISH

Before declaring completion, compare the actual result with the user's goal, review the changed area, run sufficient affected verification, update persistent project knowledge only when needed, and continue working if an important problem remains.

## Codex adapter now exists

`prototype/` is now a self-contained Codex plugin prototype.

It maps:

```text
UserPromptSubmit -> small preparation fallback
PreToolUse       -> BEFORE_CHANGE
PostToolUse      -> ON_PROBLEM when explicit failure is observable
Stop             -> BEFORE_FINISH
```

The adapter reads the shared guidance files rather than copying their content.

Turn state lives in Codex `PLUGIN_DATA`; the user's repository is not used as framework state storage.

The prototype intentionally fails open if its own hook script errors.

## Known prototype limits

- Hook coverage is a host capability, not a complete security/enforcement boundary.
- `ON_PROBLEM` currently reacts only to explicit structured failure signals; it does not guess from arbitrary output text.
- The first version uses Python for the command hook and therefore currently assumes a Python 3 command is available.
- Windows Codex hook command handling has had recent path/quoting bugs; real Windows testing is required before treating the adapter as portable.

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

Do not add more framework concepts yet.

Next, run the Codex prototype in real coding sessions and observe only three things:

1. Does the before-change interruption make the model narrow and plan better without becoming annoying?
2. Does problem guidance reduce blind retries?
3. Does the finish interruption catch real omissions often enough to justify its cost?

Use those observations to change the three guidance steps before adding more hosts or lifecycle stages.
