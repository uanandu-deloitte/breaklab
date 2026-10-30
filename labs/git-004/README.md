# Branch and Switch

`/workspace` is a git repo with a single branch, `main`, and one
commit.

## Tasks

1. Create a new branch named `feature-logging` and switch to it.
2. On `feature-logging`, add a file `logging.md` with this content:

   ```
   # Logging plan

   Add structured logging to the service.
   ```

3. Commit it (any message).
4. Switch back to `main`.
5. Confirm `main` still has only its original commit — `logging.md`
   should not exist there.

Get a shell: `docker exec -it breaklab-git-004 bash`
