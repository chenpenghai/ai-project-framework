# Foundations

This document records the framework's durable design constraints. It should stay small. Implementation details belong elsewhere.

## 1. The framework serves the project

The framework is successful only when it helps a project move faster while improving code quality and maintainability.

Do not add a mechanism merely because it is theoretically complete or architecturally interesting. Every default-path feature must justify its runtime, cognitive, and maintenance cost.

## 2. Large capability, small active surface

The framework may be internally sophisticated, but a normal change should activate only the smallest useful subset of its capabilities.

Framework scale and per-task overhead must remain decoupled.

## 3. Zero user ceremony

Users should not need to understand framework internals, select agents or skills, maintain routing tables, or repeatedly configure concepts the repository can infer.

Prefer:

1. inference,
2. conventions already expressed by the project,
3. explicit overrides only when inference is insufficient.

Derived facts must not be manually duplicated.

## 4. Universal by default

The core must not require a particular:

- programming language,
- application architecture,
- operating system,
- IDE or coding assistant,
- AI model,
- second model,
- hosted platform,
- build system,
- test framework.

Language-, tool-, and platform-specific capabilities are progressive enhancements, not prerequisites.

## 5. Reduce the reasoning surface

The framework should reduce how much of the project a coding model must understand for one task.

Primary techniques:

- modular boundaries reduce structural scope,
- pure functions reduce state and effect scope,
- single sources of truth reduce ambiguity,
- context routing reduces knowledge scope,
- affected verification reduces test scope,
- minimal changes reduce modification scope.

## 6. Modular code is the default shape

Prefer cohesive capability-oriented modules with explicit boundaries and narrow public surfaces.

Do not impose one physical directory pattern across languages. Respect native package/module systems where they exist.

The framework should prevent new structural regressions such as dependency cycles, boundary leaks, unnecessary public-surface growth, and ownerless shared business logic when these can be determined reliably.

## 7. Pure core, effects at the edge

Business decisions and deterministic transformations should be pure whenever practical.

External effects such as persistence, network access, filesystem access, environment access, current time, randomness, shared mutable state, UI/runtime interaction, and process execution should remain explicit and concentrated near boundaries.

The framework does not require every function to be pure. It should make pure logic easier to create, understand, test, and preserve.

Purity analysis must be conservative:

- PURE: sufficient evidence of determinism and no known effects,
- EFFECTFUL: confirmed direct or transitive effects,
- UNKNOWN: insufficient evidence.

UNKNOWN must never be silently treated as PURE.

## 8. Single source of truth

A fact should have one authority.

Facts derivable from code or project tooling should be generated, not manually restated. Semantic business ownership may be declared explicitly only where machines cannot infer it reliably.

Documentation should describe intent and authoritative semantics rather than duplicate generated code facts.

## 9. Deterministic gates, heuristic guidance

Use hard gates only for high-confidence, mechanically observable violations.

Use non-blocking guidance for heuristic observations such as possible extraction opportunities, suspicious shared modules, or possible semantic duplication.

Never block development based on low-confidence inference.

## 10. Structural ratchet

Existing projects may contain architectural debt. Adoption must not require rewriting the project first.

The default rule is:

- existing known debt may remain,
- new structural regressions should be rejected when deterministically detectable,
- improvements pass silently.

For projects created from this framework, the clean initial baseline should be preserved as the project grows.

## 11. Incremental by design

Normal work must operate on changed files and affected graph regions rather than repeatedly rescanning and retesting the entire repository.

Expensive analysis belongs on cached, warm, or cold paths unless the current change requires it.

Framework overhead itself is a performance budget that must be measured.

## 12. Local reasoning and local verification

Good code should be locally understandable, locally testable, and locally verifiable.

The framework should favor code that is:

- local,
- explicit,
- deterministic,
- composable.

This principle is more important than enforcing any named architecture style.

## 13. Current model only

The framework must not require a second model or cross-model review process.

The coding model receives deterministic evidence from the framework and corrects its own work. Optional external integrations may exist later, but they can never be required for the core workflow.

## 14. Core graph model

The current architectural hypothesis is that most useful behavior can be derived from three incremental graphs:

### Structure Graph

Projects, logical modules, files, symbols, dependencies, public surfaces, and structural boundaries.

### Effect Graph

Calls, direct effects, transitive effects, and PURE / EFFECTFUL / UNKNOWN classification.

### Knowledge Graph

Semantic authorities, modules, documentation, decisions, and tests.

A change set queries these graphs to derive minimal context and minimal verification.

This hypothesis should be tested before additional framework systems are added.

## 15. Framework source and consumer projects are separate products

This repository is the framework's implementation source. It is not the starter project users should build applications inside.

Framework implementation complexity must never leak into consumer repositories. A newly created consumer project starts with:

- zero application code,
- zero prescribed language,
- zero prescribed architecture,
- zero framework source code,
- only the minimum control surface required to activate the framework.

The framework may dogfood itself as a normal software project, but that does not make its source tree part of a user's project.

The consumer artifact is a static directory. It must not require a framework executable or framework runtime. The current coding model applies the framework through static project guidance and project-native mechanisms.

The canonical consumer artifact lives directly in `empty-project/`; do not maintain a second generated copy of the same static files.
