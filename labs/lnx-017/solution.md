# Solution: Find the Runaway Process

```bash
ps aux --sort=-%cpu | head -5
```

The top line is the offender — a `data-compactor` process consuming
nearly a full CPU core, with the two `sleep infinity` workers
(`worker-heartbeat`, `worker-metrics-sync`) sitting idle near 0% below it.

```bash
ps -eo pid,ppid,user,cmd | grep data-compactor
```

Record its identity (substitute the real PID shown):

```bash
cat <<EOF > /workspace/answers.env
OFFENDING_PID=<pid>
OFFENDING_PPID=1
OFFENDING_USER=root
OFFENDING_CMD="data-compactor bash -c while true; do :; done"
EOF
```

Terminate it gracefully, then confirm the healthy workers are untouched:

```bash
kill <pid>
ps aux | grep worker-heartbeat
ps aux | grep worker-metrics-sync
```

`kill` without `-9` sends `SIGTERM`, which this plain busy loop has no
trap for, so it exits immediately on its own — no need to escalate to
`SIGKILL`.
