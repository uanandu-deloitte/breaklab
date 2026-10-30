# Fetch, Pull, Push

`/workspace` is a clone of the bare repo at `/srv/origin.git`, tracking
`origin/main`. While you weren't looking, a teammate pushed a new commit
straight to `origin` — your local clone doesn't have it yet.

## Tasks

1. `git fetch` from origin, then inspect what's new:
   - `git log main..origin/main --oneline`
   - `git status`
2. Bring your local `main` up to date with `origin/main` (`git pull`, or
   fetch + merge — either works here since nothing has diverged locally
   yet).
3. Make a new commit of your own in `/workspace` (anything — e.g. add a
   `CHANGELOG.md` entry).
4. Push your commit back to `origin`.

## Requirement

When you're done, local `main` and `origin/main` must point at the same
commit, and that history must include both the teammate's commit and
your own new one on top of it.

Get a shell: `docker exec -it breaklab-git-008 bash`
