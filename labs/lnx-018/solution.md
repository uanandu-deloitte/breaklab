# Solution: Graceful Shutdown

```bash
# Terminate worker1 gracefully.
pid1=$(pgrep -f worker1)
kill "$pid1"
sleep 1
cat /workspace/worker-graceful.marker   # exists — the trap ran

# Start a second instance with a different marker path, then force-kill it.
/usr/local/bin/worker.sh /workspace/worker-forced.marker &
pid2=$!
kill -9 "$pid2"
sleep 1
ls /workspace/worker-forced.marker      # No such file — SIGKILL gave it no chance

cat <<EOF > /workspace/answers.env
GRACEFUL_MARKER_EXISTS=yes
FORCED_MARKER_EXISTS=no
EOF
```

`kill` (no flag) sends `SIGTERM`, which `worker.sh`'s trap catches to
write its marker before exiting on its own terms. `kill -9` sends
`SIGKILL`, which the kernel enforces immediately and unconditionally —
the process never regains control, so its trap never runs and no marker
is ever written.
