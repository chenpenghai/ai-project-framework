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

Do not automatically generate framework files, reorganize modules, migrate tests, or rewrite architecture.

When real work touches an area, APF should guide the AI toward a cleaner local result without unnecessarily widening task scope.

## Core control model now established

Hooks only find coarse lifecycle moments. They are not precise enough by themselves to encode a complex developer workflow.

Inside a coarse event, APF should use a small deterministic loop:

```text
coarse host event
  -> ordered small rule/step
  -> AI works
  -> structured result protocol
  -> program validates flow/step/status/evidence
  -> state advances or stays put
  -> next small rule/step
  -> ...
  -> completion
```

The program owns sequence and state. The model owns semantic judgment and execution.

Do not parse arbitrary natural-language keywords when a small explicit protocol can be required instead.

A model declaration such as `done` is not proof by itself. Mechanically verifiable evidence should still be checked mechanically when possible.

## First prototype lifecycle

The first prototype still uses three semantic events:

```text
BEFORE_CHANGE
ON_PROBLEM
BEFORE_FINISH
```

`BEFORE_CHANGE` and `ON_PROBLEM` remain simple one-shot guidance in the current Codex prototype.

`BEFORE_FINISH` is now the first real APF state-machine loop because Codex `Stop` exposes `last_assistant_message`, allowing the adapter to observe the model's structured response to the previous APF instruction.

The current loop is defined in:

```text
prototype/flows/before-finish.json
```

It contains:

```text
check_goal
  -> check_structure
  -> check_verification
  -> check_project_knowledge
```

The model reports each step using a structured `APF_RESULT` containing:

```text
flow
step
status
evidence
```

Invalid/missing protocol output does not advance. `needs_work` keeps the model on the current step. After the last step succeeds, APF requests the normal user-facing final answer and then allows Stop.

A safety limit prevents an endless malformed-protocol loop.

## Important host constraint

Do not assume every lifecycle hook can run the same state machine.

A loop needs a reliable return channel from the model to APF. Codex `Stop` currently provides this through `last_assistant_message`.

For stages such as `BEFORE_CHANGE`, a future loop requires either:

- an equivalent host event that exposes the model's response; or
- an explicit APF reporting tool/protocol endpoint the model can call.

Do not fake state progress from an event that cannot actually observe the model's result.

## Codex adapter

`prototype/` is a self-contained Codex plugin prototype.

It maps:

```text
UserPromptSubmit -> small preparation fallback
PreToolUse       -> one-shot BEFORE_CHANGE
PostToolUse      -> one-shot ON_PROBLEM for explicit failures
Stop             -> multi-step BEFORE_FINISH loop
```

Turn state lives in Codex `PLUGIN_DATA`; the user's repository is not framework state storage.

The hook script intentionally fails open if its own code crashes.

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
- GitHub Issues as duplicate internal state;
- advancing workflow state by guessing from arbitrary model prose.

## Current design frontier

Do not add another host yet.

First validate the APF Loop itself in real Codex work:

1. Does the model reliably follow the structured return protocol?
2. Are the small review steps useful, or merely repetitive?
3. Does `needs_work` actually cause the model to fix omissions before advancing?
4. Does the extra loop improve final quality enough to justify the added turns?

Only after those answers are known should we generalize the loop to other lifecycle stages or hosts.
