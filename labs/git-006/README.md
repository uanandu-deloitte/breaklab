# Merge It In

The `/workspace` repo has a `main` branch and a `feature/signup-form`
branch that both gained commits after they diverged from a shared
starting point.

## Tasks

1. Run `git log --graph --oneline --all` and look at how `main` and
   `feature/signup-form` relate to each other.
2. Decide: will merging `feature/signup-form` into `main` fast-forward,
   or will it need a true merge commit? (Hint: has `main` moved since the
   branch point?)
3. From `main`, merge `feature/signup-form` in.

## Requirement

The merge must produce a real merge commit — a commit on `main` with two
parents. A fast-forward does not satisfy this lab.

Get a shell: `docker exec -it breaklab-git-006 bash`
