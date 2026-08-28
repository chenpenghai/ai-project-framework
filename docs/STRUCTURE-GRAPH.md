# Structure Graph v0.1

The Structure Graph is the framework's structural model of a repository. It is deliberately smaller than an architecture model: it records what can be established reliably before the framework starts making suggestions.

## Levels

The graph is designed to grow across these levels without assuming that directories are modules:

`Repository -> Project/Package -> Logical Module -> File -> Symbol`

The current scanner slice implements repository, project/package, explicit logical module, and file nodes. Symbol adapters come later.

## Confidence is part of the data

The framework must preserve both evidence and how strongly that evidence supports a structural claim.

Current confidence classes are:

- `observed` — direct repository facts such as files;
- `declared` — explicit native/project declarations such as `package.json`, `go.mod`, or `MODULE.md`;
- `derived` — deterministic relationships computed from declared/observed facts, such as nearest-parent containment;
- `inferred_high` — strong but not authoritative indicators, such as a directory containing only `requirements.txt` as its Python project signal;
- `candidate` — reserved for later heuristic module hypotheses.

Low-confidence inference must never silently become a blocking architecture rule. Hard enforcement should use only evidence classes appropriate to the specific invariant.

Multiple manifests at the same path are aggregated rather than discarded. A polyglot workspace can therefore retain evidence such as both `package.json` and `go.work` without creating duplicate path nodes.

## Facts and hypotheses

Observed/declared facts and heuristic hypotheses must remain separate.

Reliable evidence may come from sources such as:

- Git and the filesystem,
- native project manifests,
- explicit `MODULE.md` declarations,
- native dependency declarations,
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

Only `module` is currently read by the scanner, and only from the opening frontmatter block. The remaining fields are placeholders for later milestones and must not be implemented until their behavior is justified.

## Declared dependency edges

The scanner now derives repository-local `depends_on` edges from native manifests when both ends can be resolved with high confidence.

Current adapters:

- Go: `go.mod` `require` entries are matched to other repository-local Go module declarations.
- Node: `package.json` dependency sections are matched to other repository-local package names.

External dependencies remain outside the repository Structure Graph for now. A manifest declaration is evidence of a project-level dependency; source import analysis will later refine relationships where useful.

Project dependency cycles are computed from these local edges. At this milestone they are diagnostic only, not enforcement.

## Affected projects

Changed files are mapped to their nearest owning project through containment edges. The framework then walks reverse local dependency edges to derive the project-level affected closure.

This is the first concrete use of the Structure Graph to reduce work: a change should not imply that every project in the repository must be verified.

The current changed-file source is the local Git working tree/index. Base-commit/PR-range analysis will be added later before this becomes a CI verification primitive.

## Scanner behavior

`apf scan [repository]` currently:

1. resolves the Git root when Git is available;
2. uses Git's tracked/untracked non-ignored file set when possible;
3. falls back to a filesystem walk when Git is unavailable;
4. records changed files without treating a dirty repository as an error;
5. detects common project/package/workspace manifests and preserves multiple manifests at one path;
6. uses high-confidence fallback indicators for ecosystems whose projects do not always have one canonical manifest;
7. detects explicit `MODULE.md` boundaries;
8. derives nearest-parent containment edges;
9. derives high-confidence repository-local project dependencies from native manifests;
10. reports project dependency cycles and affected projects;
11. emits either a compact summary or a deterministic JSON snapshot.

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

The next step is to validate the graph against real repositories, then add source import adapters only where they materially improve module/dependency accuracy.
