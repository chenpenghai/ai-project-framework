# Foundations

This document records durable design constraints for AI Project Framework. Keep it small enough to remain authoritative. Implementation details belong in focused documents.

## 1. The framework serves the project

The framework exists to make AI-assisted development faster, higher quality, and more reliable.

Do not add ceremony, metadata, agents, rules, or machinery merely because they are theoretically complete. Every mechanism must justify its cognitive, runtime, and maintenance cost.

## 2. The AI is the worker; the framework is the guide

The coding model still performs the development work.

The framework should not assume that the model already knows how to perform a complex engineering activity reliably. It should decompose difficult behavior into small actions that a weaker but capable coding model can understand and execute.

For example, do not rely on a vague instruction such as "review the code" when it can be decomposed into concrete checks such as:

- look for duplicate authoritative facts,
- inspect module-boundary violations,
- identify business logic mixed with external effects,
- verify affected module behavior,
- check whether authoritative documentation changed.

The quality target is that a merely competent coding model can follow the framework successfully; the framework must not depend on exceptional model intelligence.

## 3. Recursive modularity is the primary code shape

A project should be decomposable recursively:

```text
project
└── module
    ├── submodule
    │   └── smaller module
    └── submodule
```

A module is a cohesive unit that can be understood, changed, and verified with a clear responsibility and boundary. A module may contain smaller modules using the same principle.

The goal is composability: internally complex units should present simple external surfaces, like building blocks.

Do not define modules by arbitrary line counts or file counts. Split by responsibility, boundary, and independent reason to change.

Respect native language/package conventions; recursive modularity is a semantic rule, not one universal directory layout.

## 4. Reduce the reasoning surface

The framework should continually reduce how much a coding model must understand at once.

Primary techniques:

- recursive modules reduce structural scope,
- narrow public boundaries reduce dependency scope,
- pure functions reduce state/effect scope,
- single sources of truth reduce ambiguity,
- on-demand documentation reduces knowledge scope,
- affected verification reduces verification scope,
- focused changes reduce modification scope.

The normal reasoning path should narrow from repository -> module -> submodule -> relevant implementation units.

## 5. Pure functions are the preferred leaf shape

At the lower levels of the module tree, deterministic business logic should prefer simple pure functions where practical.

External effects such as persistence, network access, filesystem access, current time, randomness, shared mutable state, UI/runtime interaction, and process execution should be explicit and kept near boundaries.

The framework does not require every function to be pure. It should make deterministic logic easier for models to understand, compose, modify, and test.

## 6. Timed guidance beats permanent instruction load

A giant `AGENTS.md` is not the control system.

Rules and guidance should be delivered when they become relevant. Where a coding host supports lifecycle hooks, the framework should use key events such as:

- session/task start,
- user prompt submission,
- before a tool or first code-changing action,
- after code-changing actions,
- before the agent stops,
- before context compaction when useful.

At a key event, inject only the small instruction set needed for that stage.

A critical gate may block progression and return a precise reason to the model when required evidence is missing.

Host-specific lifecycle names are adapters. The framework behavior must remain host-neutral.

## 7. Instruction, evidence, gate

For important workflow rules, prefer this shape:

```text
event
  -> small instruction
  -> observable evidence
  -> gate
```

Example:

```text
before first code change
  -> identify target module, plan the change, choose verification
  -> required preparation evidence exists
  -> allow or block the change
```

This is stronger than placing the same instruction in a long prompt and hoping the model remembers it.

## 8. Mechanisms over model obedience

If a rule can be checked mechanically with high confidence, eventually express it as a project-owned mechanism rather than relying only on natural-language compliance.

Useful mechanisms include:

- compiler and type-system constraints,
- module/package boundaries,
- tests,
- linters,
- architecture tests/checks,
- schema/contract validation,
- Git hooks where appropriate,
- GitHub Actions.

Host hooks provide early intervention; project-native checks and CI provide durable protection independent of a particular coding assistant.

Do not hard-block on low-confidence semantic guesses.

## 9. Documentation is a primary system

Coding models depend on repository knowledge, so documentation is part of the operating structure of the project.

Durable constraints:

- `AGENTS.md` stays small,
- detailed knowledge is loaded on demand,
- documentation follows the project/module structure where practical,
- one fact has one authoritative source,
- facts derivable from code or tooling are not manually duplicated,
- current state belongs in current repository sources,
- history belongs in Git,
- models should update authoritative documentation when their work changes authoritative project knowledge.

A document map may be useful, but it must not become a large hand-maintained task-to-file routing database. The exact document-map and module-README contracts remain under design.

## 10. Sufficient affected verification

The objective is not minimum testing and not universal full-repository testing.

Use enough verification to cover the changed behavior and its reasonably affected dependencies.

A module is the natural organizational unit for tests, while individual pure functions may still have focused tests inside that module's test suite.

As systems grow, verification may include module tests, integration tests, architecture checks, and critical end-to-end paths.

Do not declare success while relevant verification is failing.

## 11. Large capability, small active surface

The framework may know many workflows, skills, review checks, and host integrations, but a normal task should activate only the relevant subset.

Skills, agents, and specialist guidance are on-demand capabilities, not permanent context.

Framework scale and per-task overhead must remain decoupled.

## 12. Universal by default

The core design must not require a particular:

- programming language,
- application framework,
- operating system,
- coding assistant,
- AI model,
- second model,
- package manager,
- build system,
- test framework.

Where hosts expose equivalent lifecycle events under different names, use thin adapters.

When a host lacks lifecycle hooks, degrade to the strongest portable combination available: small repository instructions plus project-native checks, Git, and CI.

## 13. Git and GitHub are first-class mechanisms

Git is the authoritative history of project evolution.

GitHub should be used where it provides durable value, especially Actions for model-independent automated verification and pull requests for change boundaries when appropriate.

Do not duplicate architecture facts or routine internal state into GitHub Issues. Issues are for genuine tracked problems, feedback, or collaboration.

## 14. Inference over duplicated configuration

Users and models should not maintain information the repository can reliably derive.

Avoid hand-maintained dependency registries, file inventories, context-routing databases, or duplicate ownership tables when code structure, manifests, Git, tests, or native tooling already express the same fact.

Explicit declarations are reserved for semantic facts machines cannot infer reliably.

## 15. Deterministic gates, heuristic guidance

Use blocking gates only for high-confidence, mechanically observable conditions.

Use non-blocking guidance for uncertain judgments such as possible module extraction, suspicious coupling, or possible semantic duplication.

The framework should help the model make semantic decisions, not pretend uncertain inference is mechanical truth.

## 16. Framework source and consumer project are separate

This repository is the framework source and research project. It is not the application project a user builds inside.

The canonical user-facing artifact lives in `empty-project/` and starts with:

- zero application code,
- zero prescribed language,
- zero prescribed application architecture,
- no framework executable,
- no framework runtime.

The static artifact may contain host configuration and control files as the design matures, but framework implementation source must not leak into consumer projects.

## 17. Graphs are research models, not the product core

Structure Graph, Effect Graph, and Knowledge Graph remain useful research concepts for reasoning about dependencies, effects, ownership, and impact.

They are not the defining product architecture and do not imply a graph database, scanner runtime, daemon, or required framework executable inside user projects.

The framework should adopt graph-derived mechanisms only when they demonstrably improve the core goals above without slowing normal development.
