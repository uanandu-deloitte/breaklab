# Solution: Where docker logs Actually Comes From

```bash
cat <<EOF > /workspace/answers.env
STDOUT_CAPTURED=yes
FILE_LOG_CAPTURED=no
LOG_FILE_PATH=/var/log/app/app.log
EOF
```

`docker logs breaklab-docker-007` shows a growing stream of `heartbeat N`
lines — that's the app's stdout, which is the only thing Docker's logging
driver ever attaches to. `/var/log/app/app.log`, inspected from a shell
inside the container, is just as real: it exists, it's current, and it's
growing at the same rate, one `file-only N <timestamp>` line per second.
`docker logs` never shows a single line of it, and never will, because it
was never attached to that file — only to PID 1's stdout/stderr file
descriptors.
