# The staging area is a deliberate step, not a formality

Git never commits your working tree directly — everything you commit
passes through the **staging area** (sometimes called the index)
first. That extra step exists so you can choose exactly what goes into
the next commit, even when your working tree has several unrelated
changes sitting in it at once.

## Two different diffs

- `git diff` compares your working tree against the staging area — it
  shows changes you haven't staged yet.
- `git diff --staged` (equivalently `git diff --cached`) compares the
  staging area against the last commit — it shows exactly what would
  be committed right now.
- `git status` summarizes both at once: a "Changes not staged for
  commit" section (what `git diff` shows) and a "Changes to be
  committed" section (what `git diff --staged` shows), plus a separate
  "Untracked files" section for anything Git has never seen before.

## Staging on purpose

`git add <path>` moves exactly that file's current changes into the
staging area. Nothing else in the working tree is affected — an
unrelated modified file stays unstaged, and a new untracked file stays
untracked, until you explicitly `git add` it too.

This is what makes it possible to have two unrelated changes sitting in
the same working tree — one you're ready to commit, one you're not —
and commit only the one you mean to. Reaching for `git add .` or
`git commit -a` out of habit stages everything indiscriminately,
which is exactly how an unrelated scratch file or half-finished edit
ends up in a commit it has nothing to do with.
