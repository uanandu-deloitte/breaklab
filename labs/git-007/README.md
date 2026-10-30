# Resolve the Conflict

Companion lab to **git-006**. The `/workspace` repo has a `main` branch
and a `feature/raise-timeout` branch that both edited the same
`timeout:` line of `config.yaml` differently since their common
ancestor:

- `main` changed it to reduce the timeout (`timeout: 45`), reasoning it
  should fail fast.
- `feature/raise-timeout` changed it to increase the timeout
  (`timeout: 60`), to tolerate a slow upstream dependency.

## Decision already made

After investigating, the team decided the longer timeout is correct —
the slow upstream dependency genuinely needs it. `main`'s shorter
timeout should not survive the merge.

## Tasks

1. From `main`, run `git merge feature/raise-timeout`. This will produce
   a real conflict in `config.yaml` — that's expected.
2. Resolve the conflict by hand so `config.yaml` ends up with **exactly**
   this content:

   ```
   service: acme-api
   port: 8080
   timeout: 60
   ```

3. Remove every conflict marker (`<<<<<<<`, `=======`, `>>>>>>>`).
4. `git add config.yaml` and `git commit` to finish the merge.

## Constraints

- No conflict markers may remain anywhere under `/workspace`.
- The merge must be completed with an actual merge commit (two parents),
  not aborted and redone some other way.

Get a shell: `docker exec -it breaklab-git-007 bash`
