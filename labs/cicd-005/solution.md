# Solution: Lockfiles and Crash Safety

```bash
cat /var/run/deploy/deploy.lock
# PID=51234
# STARTED=<~30 days ago>
# HOST=ci-runner-03

ps -p 51234               # no such process
date -u                   # current date, weeks after STARTED

cat <<EOF > /workspace/answers.env
LOCK_PID=51234
PID_RUNNING=no
LOCK_STATUS=stale
EOF
```

The lock records PID `51234` as the process that created it. `ps -p
51234` finds nothing — this container's PID namespace has no process
with that number, so whatever created the lock is long gone. The
`STARTED` timestamp is about 30 days old, far longer than any real
deploy step runs. Both signals point the same way: this lock isn't
protecting anything in progress, it's debris from a run that crashed
(or was killed) before it could remove its own lock on the way out.
That makes it safe to clear by hand once you've actually confirmed —
not assumed — that nothing is still using it.
