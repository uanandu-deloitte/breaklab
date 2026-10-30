# Agent instructions

## Code discipline
- Write the minimum code that correctly solves the task. No speculative features, unused params, wrapper functions that add nothing, single-implementation interfaces, or config nobody asked for.
- No logging, comments, validation, or error handling beyond what the task needs. Handle real errors (returned `error` values, I/O, user input); don't guard against impossible states.
- Reuse existing helpers and patterns in this codebase before writing new ones. Never duplicate logic that already exists.
- Delete dead code you touch. No commented-out code, no TODOs.
- Prefer Go stdlib. A new module dependency needs a one-line justification in the PR/commit message.
- Default to efficient algorithms and data structures: no repeated work, no needless allocations, no N+1 calls (e.g. file reads, subprocess spawns, HTTP calls in a loop).

## Context discipline
- Treat the files as currently shown (via Read, not memory) as source of truth. If a file hasn't been read in its current form this session, read it before editing it.
- Don't reintroduce code, names, or patterns that were removed or replaced.
- Keep changes scoped to the task. Don't refactor, rename, or reformat unrelated code.
- If a request conflicts with existing code or an earlier decision, flag it in one line before proceeding.

## Destructive actions
- Never delete or overwrite state files (e.g. `progress.json`, lab state) without explicit confirmation — this repo has no safety net for that data (see memory: no-rm-on-real-state-files).

## Output format
- Code only: a diff or the changed functions. Full files only when new or explicitly requested.
- No intro, recap, or summary. State assumptions and non-obvious decisions in one line each, nothing more.
