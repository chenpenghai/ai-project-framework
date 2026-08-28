# Coding Host Hook Research

Research snapshot: 2026-08-28.

This document records the current host-integration hypothesis. It is not a promise that vendor APIs will remain unchanged. Verify official documentation again before implementing an adapter.

## Why hooks matter

The framework should not depend on a model remembering a large permanent rule set.

The preferred control pattern is:

```text
key lifecycle event
  -> inject a small stage-specific instruction
  -> require observable evidence
  -> allow or block progression
```

Examples:

- task begins -> decompose the requested work,
- first code-changing action -> require preparation/plan evidence,
- after changes -> run relevant fast checks,
- agent attempts to stop -> perform structured completion review and required verification.

## Emerging lifecycle vocabulary

Across mainstream coding agents, the following Claude-style event names are increasingly common or easy to map to equivalent events:

```text
SessionStart
UserPromptSubmit
PreToolUse
PostToolUse
Stop
PreCompact
```

The framework should use these names as its host-neutral vocabulary unless later research shows a better common denominator.

These names describe semantics, not a dependency on Claude Code.

## Current compatibility picture

### Near-direct compatibility

The following hosts currently expose the same or very similar lifecycle concepts and, in several cases, the same event names:

- Claude Code
- OpenAI Codex
- CodeBuddy IDE
- VS Code Agent / GitHub Copilot integrations
- Kiro
- Qwen Code
- Augment Code
- Cursor (similar names, often lower camel case)
- Cline (tool hooks align closely; completion event naming differs)

Important behavior for the framework is not exact spelling. It is whether the host can:

1. observe the event,
2. run project-controlled logic,
3. block a code/tool action when required,
4. return a reason/context to the coding model,
5. on completion gates, cause the model to continue when work is incomplete.

### Semantic adapters required

Some hosts expose equivalent lifecycle points under different names or plugin APIs.

Examples observed during research:

- Gemini CLI: concepts such as `BeforeTool`, `AfterTool`, `AfterAgent`, plus model-level lifecycle hooks.
- OpenCode: plugin events such as session creation and tool execution before/after.
- Windsurf: more operation-specific pre/post hooks rather than the common generic vocabulary.

These should be thin adapters to the framework lifecycle, not separate framework logic.

### Partial/no full lifecycle

Some tools may not expose a complete blocking lifecycle. Aider is a representative example: it provides strong automatic lint/test feedback loops but not the same general hook surface.

For these hosts the framework must degrade gracefully to:

```text
small repository instructions
+ project-native tests/checks
+ Git hooks where useful
+ GitHub Actions
```

The project must remain usable even without host hooks.

## Proposed framework lifecycle

Do not freeze more lifecycle stages than needed yet. The currently useful semantic stages are:

```text
TASK_START
BEFORE_CHANGE
AFTER_CHANGE
BEFORE_FINISH
```

These are framework concepts. They can be implemented using host events such as:

```text
TASK_START
  <- UserPromptSubmit / session-task equivalent

BEFORE_CHANGE
  <- PreToolUse on the first code-changing action

AFTER_CHANGE
  <- PostToolUse or change batching logic

BEFORE_FINISH
  <- Stop / agent-completion equivalent
```

`SessionStart` and `PreCompact` are supporting lifecycle events rather than necessarily separate workflow stages.

## Critical design rule

Do not make hooks merely print reminders.

Where a requirement matters, prefer:

```text
instruction
+ evidence
+ gate
```

For example, if a task requires a plan before coding, the reliable point is not "hope the model plans after reading AGENTS.md". The reliable point is the first attempted code-changing tool call:

```text
PreToolUse(first change)
  -> preparation evidence missing
  -> block
  -> tell the model exactly what preparation is missing
```

Likewise, a generic request such as "review the code" should be decomposed at the completion gate into small review actions a weaker model can execute one by one.

## Implementation constraints

- Framework workflow logic must have one source of truth.
- Host adapters only translate event names, payloads, and blocking/feedback formats.
- Do not duplicate the rule set for each coding assistant.
- Do not require every host to support every enhancement.
- Do not require a consumer-side APF executable merely to implement hooks.
- Host hooks are early-control mechanisms; project-native checks and CI remain the durable enforcement layer.

## Sources to re-check before implementation

Official documentation examined during the 2026-08-28 research included:

- Claude Code hooks documentation
- OpenAI Codex hooks documentation
- CodeBuddy IDE hooks documentation
- Cursor hooks documentation
- VS Code / GitHub Copilot hooks documentation
- Gemini CLI hooks documentation
- Kiro hooks documentation
- Qwen Code hooks documentation
- Augment Code hooks documentation
- OpenCode plugin/hook documentation
- Cline hook documentation
- Windsurf hooks documentation
- Aider lint/test automation documentation

Vendor capabilities change quickly. Adapter implementation should always be based on the current official specification, while this document preserves the framework-level conclusion: mainstream coding agents are converging on interceptable lifecycle events, and the framework should exploit that convergence.
