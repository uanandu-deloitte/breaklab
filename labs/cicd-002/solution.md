# Solution: Exit Codes and Pipeline Stages

```bash
cat <<'EOF' > /opt/pipeline/run.sh
#!/usr/bin/env bash
echo "=== build ===" && /opt/pipeline/build.sh && \
echo "=== test ===" && /opt/pipeline/test.sh && \
echo "=== deploy ===" && /opt/pipeline/deploy.sh
EOF
chmod +x /opt/pipeline/run.sh

rm -f /workspace/bad-deploy-happened
/opt/pipeline/run.sh
echo "exit code: $?"
ls /workspace/bad-deploy-happened 2>&1   # No such file or directory
```

The original `run.sh` joined the three stage invocations with `;`, which
is just a command separator — it carries no information about exit
codes. `build.sh; test.sh; deploy.sh` runs all three in order no matter
what any of them returned, so a failing `test.sh` had zero effect on
whether `deploy.sh` ran next.

Replacing `;` with `&&` between the invocations means each stage only
runs if the one before it exited `0`. `build.sh` succeeds, so `test.sh`
runs; `test.sh` fails (by design, as fixture data), so the `&&` chain
stops right there — `deploy.sh` is never reached, and
`/workspace/bad-deploy-happened` never gets created. `run.sh` itself now
exits non-zero too, correctly reporting the pipeline as failed instead of
silently succeeding.

(Adding `set -e` near the top of the script instead, with the three
invocations on separate lines, would have the same effect: the script
would exit immediately at `test.sh`'s failure, never reaching
`deploy.sh`.)
