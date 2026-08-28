# APF Prototype

This directory is the smallest experiment for the current APF idea and is also a self-contained Codex plugin prototype.

It does not create or modify framework files in the user's project.

## Current experiment

Hooks find coarse lifecycle moments. APF then breaks a complex activity into small steps and drives those steps with an explicit protocol and state machine.

```text
host event
  -> APF flow
  -> one small instruction
  -> AI works
  -> structured APF_RESULT
  -> state advances
  -> next small instruction
  -> ...
  -> normal user-facing completion
```

## Prototype structure

```text
prototype/
├── .codex-plugin/plugin.json
├── flows/
│   └── before-finish.json
├── guidance/
│   ├── before-change.md
│   └── on-problem.md
├── hooks/
│   ├── hooks.json
│   └── apf_hook.py
└── tests/
```

`BEFORE_CHANGE` and `ON_PROBLEM` are still simple one-shot guidance.

`BEFORE_FINISH` is the first real APF Loop because Codex `Stop` exposes the model's last assistant message, giving the adapter a reliable return channel for the protocol.

## BEFORE_FINISH loop

`flows/before-finish.json` currently contains four small steps:

```text
check_goal
  -> check_structure
  -> check_verification
  -> check_project_knowledge
```

For each step, the model must return:

```text
<APF_RESULT>
{"flow":"before_finish","step":"check_goal","status":"done","evidence":"brief concrete evidence"}
</APF_RESULT>
```

Allowed status values:

- `done` — the current step is complete;
- `needs_work` — the step found unresolved work, so APF keeps the model on the same step until it is resolved.

Wrong flow IDs, wrong step IDs, missing evidence, malformed JSON, or a missing marker do not advance the state machine.

After the last review step succeeds, APF asks the model once more for the normal user-facing final answer and then allows the turn to stop. Internal APF protocol markers should therefore not become the final answer shown to the user.

A safety limit prevents a malformed protocol from trapping the host in an endless Stop loop.

Turn state is stored under Codex `PLUGIN_DATA`, not inside the user's repository.

## Why only BEFORE_FINISH loops today

Do not pretend every hook can run the same protocol.

A structured loop requires a reliable way for the adapter to observe the model's response to the previous APF instruction. Codex `Stop` provides `last_assistant_message`, so `BEFORE_FINISH` can do this directly.

`PreToolUse` currently does not give this prototype an equivalent model-response channel. A future `BEFORE_CHANGE` loop therefore needs either a comparable host capability or an explicit APF report tool/protocol endpoint.

## Install for a real Codex test

The repository includes `.agents/plugins/marketplace.json`, so it can be added directly as a development marketplace:

```text
codex plugin marketplace add chenpenghai/ai-project-framework
codex plugin add apf-prototype@apf-development
```

Restart Codex after installation. Plugin hooks still require Codex's normal hook trust review before they run.

## Run tests

```text
python -m unittest discover prototype/tests
```

CI runs both the existing Go tests and these prototype tests.

## What this prototype is testing

We want to learn whether a deterministic program can make a coding model perform complex developer behavior more reliably by sequencing small semantic checks instead of giving one large instruction.

Do not add more lifecycle stages until this loop has been exercised in real coding sessions.
