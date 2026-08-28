# APF Prototype

This directory is the smallest experiment for the current APF idea and is also a self-contained Codex plugin prototype.

It does not create or modify framework files in the user's project.

## Shared guidance

```text
BEFORE_CHANGE  -> guidance/before-change.md
ON_PROBLEM     -> guidance/on-problem.md
BEFORE_FINISH  -> guidance/before-finish.md
```

The Codex adapter does not copy this guidance. It loads these files at runtime.

## Codex prototype

```text
prototype/
├── .codex-plugin/plugin.json
├── guidance/
├── hooks/
│   ├── hooks.json
│   └── apf_hook.py
└── tests/
```

Behavior:

1. `UserPromptSubmit` adds a very small preparation fallback.
2. The first `apply_patch` / `Edit` / `Write` attempt in a turn is paused once and receives `BEFORE_CHANGE`.
3. `PostToolUse` adds `ON_PROBLEM` once when the tool response contains an explicit structured failure signal.
4. The first `Stop` is continued with `BEFORE_FINISH`; a continuation Stop (`stop_hook_active=true`) is allowed through.

Turn state is stored under Codex `PLUGIN_DATA`, not inside the user's repository.

## Install for a real Codex test

The repository includes `.agents/plugins/marketplace.json`, so it can be added directly as a development marketplace:

```text
codex plugin marketplace add chenpenghai/ai-project-framework
codex plugin add apf-prototype@apf-development
```

Restart Codex after installation. Plugin hooks still require Codex's normal hook trust review before they run.

## Run adapter tests

```text
python -m unittest discover prototype/tests
```

## What this prototype is testing

We want to learn whether short, well-timed instructions improve model behavior enough to justify the framework.

Do not add more lifecycle stages until this small loop has been exercised in real coding sessions.
