# Foundations

This document records durable constraints for AI Project Framework (APF).

## 1. The AI is the worker; APF is the guide

APF should not try to replace the coding AI. It should make difficult engineering behavior easier to execute by delivering small, concrete instructions at useful moments.

Prefer "check these three things now" over a permanent manual the model must remember.

## 2. Existing projects are first-class

APF must work on an existing repository without requiring a migration first.

Enabling APF must not automatically create framework files, reorganize directories, or rewrite project structure.

A project should improve progressively as normal development touches real areas of the codebase.

## 3. Improve locally, not ceremonially

When a task touches an area, prefer leaving that area easier to understand, change, and verify than before.

Do not expand a feature or bug-fix task into unrelated repository-wide cleanup merely to satisfy framework ideals.

## 4. Recursive modularity is the preferred code shape

Projects should be decomposable into coherent building blocks:

```text
project
└── module
    ├── submodule
    └── submodule
```

A module may contain smaller modules using the same principle. Internally complex units should expose simple boundaries.

Do not define modules by arbitrary line or file counts. Split by responsibility, boundary, and independent reason to change.

## 5. Reduce the AI's reasoning surface

The framework should help the AI narrow from repository -> affected module -> smaller unit -> relevant implementation.

Useful techniques include narrow interfaces, local documentation, pure functions, explicit effects, single sources of truth, and affected verification.

## 6. Prefer pure deterministic leaf logic

Where practical, lower-level business logic should be expressed as simple deterministic functions.

Filesystem, network, persistence, time, randomness, UI/runtime interaction, process execution, and shared mutable state should remain explicit near boundaries.

This is a preference, not a requirement that every function be pure.

## 7. Timed guidance beats giant prompts

APF should activate guidance at key lifecycle moments rather than loading every rule permanently.

The first prototype uses only three semantic moments:

```text
BEFORE_CHANGE
ON_PROBLEM
BEFORE_FINISH
```

More stages should be added only when a real failure pattern proves they are useful.

## 8. Decompose expert behavior

Do not assume the model can reliably execute vague commands such as "review the code".

Break complex behavior into small checks, for example:

- inspect the affected boundary;
- check for duplicated authoritative facts;
- separate deterministic logic from external effects where useful;
- verify affected behavior;
- update authoritative documentation only when project knowledge changed.

The target is that a merely competent coding model can execute the workflow successfully.

## 9. Guidance first; hard gates only when justified

Natural-language guidance is appropriate for semantic judgment.

Use blocking gates only when the condition is high-confidence and observable. Project-native tests, type checks, linters, architecture checks, Git hooks, and CI are stronger than relying on model obedience for mechanical facts.

Do not pretend uncertain architectural judgments are deterministic truth.

## 10. One shared core, thin host adapters

Framework behavior must have one source of truth.

Host-specific plugins/adapters may translate events, payloads, blocking semantics, and feedback formats, but must not duplicate the development rules.

The same APF guidance should be reusable across coding assistants as far as their extension systems allow.

## 11. No mandatory project footprint

APF should not require `.apf/`, an `AGENTS.md`, generated scaffolding, an APF executable, daemon, scanner, or runtime inside every user project merely to function.

Project-owned files may be created later only when they provide value to that specific project, not because installation demands them.

## 12. Documentation stays close to authority

Coding models depend on repository knowledge.

Keep documentation small, local, and authoritative. Avoid duplicating facts derivable from code or tooling. Current state belongs in current repository sources; history belongs in Git.

Module-local documentation remains a useful direction when the module actually needs persistent knowledge.

## 13. Sufficient affected verification

Verification should cover changed behavior and reasonably affected dependencies without defaulting to unrelated full-repository work.

Tests should naturally follow module boundaries where practical.

Do not claim completion while relevant verification is failing.

## 14. Large capability, small active surface

APF may eventually contain many skills, checks, workflows, and host integrations, but each task should activate only what is relevant now.

Framework capability must not translate into permanent context or workflow overhead.

## 15. Git and CI remain durable mechanisms

Host hooks provide early intervention. Git and CI provide model-independent history and enforcement.

GitHub Issues should represent genuine tracked problems or collaboration, not duplicate routine internal architecture or task state.

## 16. Research is not the product

Existing scanner and graph work may remain useful research, but it is not the defining product architecture.

Do not revive a consumer runtime, graph database, daemon, mandatory scanner, or static empty-project product unless future evidence clearly justifies it.
