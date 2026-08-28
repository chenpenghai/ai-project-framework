# Host Integration and Hooks

Research snapshot: 2026-08-28.

Vendor APIs change quickly. Re-check official documentation before implementing each host adapter.

## Purpose

APF needs lifecycle access so it can give the coding AI a small instruction at the moment that instruction matters.

The host integration is not the framework logic. It is only the bridge between a coding assistant and the shared APF core.

```text
host event
  -> thin adapter/plugin
  -> shared APF semantic event
  -> shared guidance
```

## Product direction

Prefer plugin/extension integration over placing many vendor-specific configuration directories in every user repository.

The desired model is:

```text
APF core
├── shared guidance
└── semantic events

host adapters
├── Claude adapter
├── Codex adapter
├── Cursor adapter
├── Gemini adapter
└── ...
```

Each adapter should be as thin as possible.

It may translate:

- event names;
- event payloads;
- tool/change classification;
- blocking/continue semantics;
- feedback returned to the model.

It must not contain a copied version of the development rules.

## No mandatory project installation

APF activation must not require every project to contain `.claude/`, `.codex/`, `.cursor/`, `.gemini/`, `.apf/`, or similar framework-generated configuration.

If a host requires project-local configuration and offers no better integration mechanism, that may be supported as a fallback, not treated as the universal architecture.

Hosts with plugin/extension support should prefer an install-once adapter that can operate on whichever repository the coding assistant currently opens.

## First prototype lifecycle

Do not start with every lifecycle event exposed by every vendor.

The first APF prototype uses three framework-level events:

```text
BEFORE_CHANGE
ON_PROBLEM
BEFORE_FINISH
```

### BEFORE_CHANGE

Goal: intervene just before meaningful modification begins.

Possible host signals include a pre-tool event for write/edit operations or an equivalent code-changing action.

APF should inject the shared `before-change` guidance.

### ON_PROBLEM

Goal: intervene when the agent is getting stuck or evidence invalidates its current approach.

This event may not map to one universal native hook. Adapters may infer it from host-visible signals such as:

- tool/command failure;
- test failure;
- repeated failed attempts;
- explicit model report that an assumption is wrong.

Because inference confidence varies, the first implementation may treat this as guidance rather than a hard gate.

### BEFORE_FINISH

Goal: intervene before the agent declares the task complete.

Possible host signals include `Stop`, task-complete, after-agent, or equivalent completion lifecycle events.

Where the host permits it, APF should be able to return feedback that causes the model to continue when important work or verification remains.

## Common host vocabulary

Mainstream coding tools currently expose many comparable lifecycle concepts, often around names similar to:

```text
SessionStart
UserPromptSubmit
PreToolUse
PostToolUse
Stop
PreCompact
```

Other hosts use different names such as `BeforeTool`, `AfterTool`, `AfterAgent`, or plugin-specific event identifiers.

APF should map semantic equivalents rather than standardizing on one vendor's spelling.

## Compatibility strategy

Use three levels:

```text
1. Plugin/extension adapter
   best integration, no normal project clutter

2. Thin project-local adapter
   only when the host requires it

3. No lifecycle adapter
   APF capabilities degrade; project-native tests/Git/CI still work
```

A host does not need to implement every APF enhancement to be useful.

## Guidance versus gates

The first prototype should optimize for useful intervention, not maximum blocking.

Use a hard gate only when:

1. the host reliably supports blocking/continuation;
2. the condition is observable with high confidence;
3. blocking improves development rather than adding ceremony.

Semantic architectural judgments should normally be guidance.

## Shared source of truth

The prototype guidance text lives in `prototype/guidance/`.

Adapters should point to or package that shared source rather than recreate its content.

The experiment succeeds if the same guidance can be delivered through different hosts with only thin translation code.

## Current research conclusion

There is not yet one plugin binary/package format that runs unchanged in every coding assistant.

However, lifecycle concepts are similar enough that one APF core plus multiple thin adapters is practical.

Open plugin standards may reduce adapter differences over time, but APF must not assume full cross-host Hook compatibility until vendors actually provide it.
