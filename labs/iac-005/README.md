# Lockfiles and Crash Safety

This container has a small infra-as-code setup, same shape as before:

- `/infra/resources.yaml` declares the desired content of a file.
- `/etc/release-notes.txt` is the real file `infra-apply` writes to match it.
- `/usr/local/bin/infra-apply` reconciles the real file to match the
  declaration — but first it checks for a lock at
  `/var/run/infra/apply.lock`. If that lock already exists, it refuses
  to run at all, to stop two applies from racing against each other.

Right now, `infra-apply` is refusing to run. There's a lock in place.

## Tasks

1. Run `infra-apply` and read exactly what it says.
2. `cat /var/run/infra/apply.lock` — see which PID the lock claims is
   holding it, and when it was started.
3. Check whether that PID is actually a real, running process right now.
4. `stat /var/run/infra/apply.lock` — look at how old the lock is.
5. Decide: is this a live run you should leave alone, or a stale lock
   left behind by a crash that never got to clean up after itself?
6. If it's stale, remove it, then run `infra-apply` again and confirm it
   finally succeeds.

Record your findings in `/workspace/answers.env`, one `KEY=value` per
line:

```
LOCK_PID=
LOCK_PROCESS_ALIVE=
APPLY_SUCCEEDED=
```

- `LOCK_PID` — the PID recorded inside the original lock file.
- `LOCK_PROCESS_ALIVE` — `yes` or `no`: is that PID actually a running
  process?
- `APPLY_SUCCEEDED` — `yes` or `no`: after you dealt with the lock, did
  `infra-apply` run to completion and update `/etc/release-notes.txt`?

Get a shell: `docker exec -it breaklab-iac-005 bash`
