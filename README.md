# AI Project Framework

A general-purpose project framework for AI-assisted software development.

The framework exists to help a coding AI build software faster and more reliably by reducing the amount of judgment, memory, and project-wide reasoning required for each step.

The coding AI is still the worker. The framework's job is to guide that worker at the right moments, break difficult engineering behavior into small executable actions, and use project-native mechanisms to catch mistakes that should not depend on model obedience.

## User-facing product

The user-facing artifact is the `empty-project/` directory in this repository.

Copy that directory, rename it for the new project, open it with a coding assistant, and start describing the product to build.

The initial project is intentionally static and empty:

```text
empty-project/
├── AGENTS.md
└── .apf/
    └── project.yaml
```

It contains no application code, no prescribed programming language, no prescribed application architecture, no executable, and no framework runtime.

## Core direction

### Recursive modularity

The project should grow as a tree of composable modules:

```text
project
└── module
    ├── submodule
    │   └── smaller module
    └── submodule
```

A module may contain smaller modules. Each module should remain understandable as one unit, expose a narrow boundary, and hide internal details from callers. At the bottom of the tree, deterministic logic should prefer simple pure functions where practical.

The goal is not a particular folder layout. The goal is to let an AI narrow its reasoning from project -> module -> submodule -> small implementation unit instead of repeatedly understanding the whole repository.

### Timed guidance instead of giant prompts

Do not put every rule into `AGENTS.md` and hope the model remembers it.

Where the coding host supports lifecycle hooks, the framework should give the AI small instructions at the moment they become relevant, for example:

```text
Task starts
  -> clarify and decompose the work
Before the first code change
  -> confirm module, plan, boundaries, and verification
After changes
  -> run focused checks and collect evidence
Before the AI stops
  -> perform structured review and required verification
```

Complex instructions such as "review the code" should be decomposed into smaller checks the model can execute reliably, such as checking duplicate sources of truth, module-boundary leaks, mixed side effects, missing module tests, and stale documentation.

### Mechanisms over model obedience

Rules describe desired behavior. Mechanisms should enforce everything that can be observed reliably.

Examples include project-native tests, compiler/type checks, linters, architecture tests, Git hooks where useful, and GitHub Actions. Host hooks provide earlier feedback but must not be the only protection because hook support differs between coding assistants.

## Documentation direction

Documentation is a primary part of the framework because coding models depend on repository context.

- `AGENTS.md` must remain small.
- Detailed knowledge should be loaded on demand rather than placed in one permanent prompt.
- Module-local documentation should travel with the module it explains.
- One fact should have one authoritative source.
- Facts derivable from code or tooling should not be manually duplicated in documents.
- Current project state belongs in repository files; history belongs in Git.

The exact documentation map and module README contract are still under design and should not be prematurely frozen.

## Verification direction

Verification should be organized around modules and affected behavior.

The objective is not "minimal verification". It is **sufficient affected verification**: enough checks to cover the changed behavior and its reasonably affected dependencies, without running unrelated work by default.

## GitHub

Git and GitHub are first-class mechanisms:

- Git stores project history and change truth.
- GitHub Actions provides model-independent automated gates.
- Pull requests provide a natural change-review boundary when used.
- Issues are for genuine tracked problems or collaboration, not as a duplicate architecture or internal state database.

## Framework source repository

Everything outside `empty-project/` exists to design, test, and maintain the framework itself. Framework development source, experiments, tests, and documentation are not part of a user's project.

The current Go scanner and graph code are framework-development research tools. They are not the product architecture and are not required by projects copied from `empty-project/`.

See `docs/FOUNDATIONS.md` for durable constraints and `docs/HOST-HOOKS.md` for the current lifecycle-hook research snapshot.
