# Structure Graph v0.1

The Structure Graph is the framework's observed structural model of a repository. It is deliberately smaller than an architecture model: it records what can be established reliably before the framework starts making suggestions.

## Levels

The graph is designed to grow across these levels without assuming that directories are modules:

`Repository -> Project/Package -> Logical Module -> File -> Symbol`

The first scanner slice implements repository, project/package, explicit logical module, and file nodes. Symbol and dependency adapters come later.

## Facts and hypotheses

Observed facts and heuristic hypotheses must remain separate.

Facts may come from sources such as:

- Git and the filesystem,
- native project manifests,
- explicit `MODULE.md` declarations,
- language adapters that can resolve imports or visibility deterministically.

Hypotheses may later come from signals such as directory cohesion, naming, change coupling, or clustering. Hypotheses are guidance only and must never silently become hard module boundaries.

## Stable IDs

Current node IDs are path-derived and deterministic:

- `repository:.`
- `project:<relative-path>`
- `module:<relative-path>`
- `file:<relative-path>`

This is intentionally simple. IDs should only become more complex if real repository shapes prove that path identity is insufficient.

## Derived facts are not authored

Dependencies, file membership, exports, test lists, and similar facts should be derived from project tooling or source code whenever possible. `MODULE.md` is reserved for semantic intent that machines cannot reliably infer.

A minimal declaration can be as small as:

```yaml
---
module: order
purpose: Own order lifecycle and business rules.
owns:
  - order.lifecycle
pure:
  - core/**
---
```

Only `module` is currently read by the scanner. The remaining fields are placeholders for later milestones and must not be implemented until their behavior is justified.

## Scanner behavior

`apf scan [repository]` currently:

1. resolves the Git root when Git is available;
2. uses Git's tracked/untracked non-ignored file set when possible;
3. falls back to a filesystem walk when Git is unavailable;
4. records changed files without treating a dirty repository as an error;
5. detects common project/package manifests;
6. detects explicit `MODULE.md` boundaries;
7. derives nearest-parent containment edges;
8. emits either a compact summary or a deterministic JSON snapshot.

Unsupported languages and unknown repository shapes are valid states. The scanner should report what it knows and leave the rest unknown.

## Non-goals for this slice

Do not add yet:

- automatic module splitting,
- heuristic blocking rules,
- purity/effect analysis,
- semantic duplicate detection,
- agent orchestration,
- architecture scoring,
- a large configuration language.

The next step is to validate this fact layer against real repositories, then add import/project adapters only where they materially improve structural accuracy.
