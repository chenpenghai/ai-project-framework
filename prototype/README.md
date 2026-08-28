# APF Prototype

This is the smallest experiment for the current APF idea.

It is deliberately host-neutral. It does not create or modify files in the user's project.

The prototype contains three shared guidance events:

```text
BEFORE_CHANGE  -> guidance/before-change.md
ON_PROBLEM     -> guidance/on-problem.md
BEFORE_FINISH  -> guidance/before-finish.md
```

A host adapter should do only two things:

1. detect/map the relevant host lifecycle event;
2. deliver the matching shared guidance to the coding AI.

Do not copy the guidance into each adapter.

## What this prototype is testing

We want to learn whether short, well-timed instructions improve model behavior enough to justify the framework.

Success is not measured by how many rules or hooks APF has. The sample is useful if it reliably causes the AI to:

- narrow work before modifying code;
- stop blind trial-and-error when stuck;
- verify and review the actual result before claiming completion.

Only after this works should more lifecycle stages, skills, or enforcement mechanisms be added.
