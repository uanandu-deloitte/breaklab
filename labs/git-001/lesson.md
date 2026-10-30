# Reading a repo's state before you touch it

A Git repository carries more state than whatever's in front of you in
an editor: a checked-out branch, a commit history, other branches that
exist but aren't checked out, and a working tree that may or may not
match the last commit. Before you commit, branch, or reset anything,
read all of that first.

## Where to look

- `git status` is the single most useful command for orientation. Its
  first line names the checked-out branch. Below that, it separates
  **staged changes**, **unstaged changes to tracked files**, and
  **untracked files** into their own sections.
- `git branch` (no arguments) lists every local branch and marks the
  checked-out one with a `*`. A branch existing is not the same as a
  branch being checked out — only one branch can be checked out at a
  time, but a repo can hold many.
- `git log --oneline` prints commit history newest-first, one line per
  commit: a short SHA followed by its message. The first line is
  always the commit `HEAD` currently points at.

## Tracked, modified, and untracked — not the same thing

- A **tracked** file is one Git already knows about (it's in the last
  commit, or staged to be in the next one).
- A tracked file becomes **modified** when its content on disk differs
  from what's recorded in the last commit — `git status` lists it under
  "Changes not staged for commit" (or "Changes to be committed" if
  you've already staged it).
- An **untracked** file is one Git has never committed and isn't
  staged — brand new to the repo as far as Git knows. `git status`
  lists these separately, under "Untracked files".

Confusing "untracked" with "modified" is how an unrelated scratch file
ends up accidentally committed, or a real edit gets ignored because it
looked unfamiliar in the status output.
