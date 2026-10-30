# Lockfiles and Crash Safety

This container has a deploy lock sitting at `/var/run/deploy/deploy.lock`.
Nothing in this container is actively deploying right now — your job is
to figure out, from the lock's own content and the live state of the
container, whether this lock is still meaningful or just leftover
debris from a crash.

## Tasks

1. Read the lock file: `cat /var/run/deploy/deploy.lock`. It records a
   `PID` and a `STARTED` timestamp — figure out what each one is
   supposed to mean.
2. Check whether the PID it names is actually a running process in this
   container right now (`ps -p <pid>`, or `ps aux | grep <pid>`).
3. Check the `STARTED` timestamp against the current date (`date -u`).
   How old is this lock?
4. Decide: is this lock still protecting an in-progress deploy, or is it
   stale — left behind by a process that's gone?

Record your findings in `/workspace/answers.env`, one `KEY=value` per
line:

```
LOCK_PID=
PID_RUNNING=
LOCK_STATUS=
```

- `LOCK_PID`: the PID number recorded in `/var/run/deploy/deploy.lock`.
- `PID_RUNNING`: `yes` or `no` — is that PID actually a running process
  in this container right now?
- `LOCK_STATUS`: `stale` or `active` — your overall conclusion about
  whether this lock still means anything.

Get a shell: `docker exec -it breaklab-cicd-005 bash`
