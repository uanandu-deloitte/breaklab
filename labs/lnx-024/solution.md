# Solution: Write a Health Script

```bash
cat <<'SCRIPT' > /workspace/health.sh
#!/usr/bin/env bash
set -uo pipefail

PROCESS_PATTERN="worker.sh"
PORT=80
HEALTH_URL="http://localhost:${PORT}/healthz"
DISK_PATH="/var/log/app"
DISK_THRESHOLD_MB=100

if ! pgrep -f "$PROCESS_PATTERN" >/dev/null 2>&1; then
  echo "UNHEALTHY: required process ($PROCESS_PATTERN) is not running"
  exit 1
fi

if ! (exec 3<>"/dev/tcp/127.0.0.1/${PORT}") 2>/dev/null; then
  echo "UNHEALTHY: nothing is listening on port ${PORT}"
  exit 1
fi
exec 3>&- 2>/dev/null || true

code=$(curl -s -o /dev/null -w '%{http_code}' --max-time 3 "$HEALTH_URL")
if [ "$code" != "200" ]; then
  echo "UNHEALTHY: $HEALTH_URL returned $code"
  exit 1
fi

disk_mb=$(du -sm "$DISK_PATH" 2>/dev/null | cut -f1)
disk_mb=${disk_mb:-0}
if [ "$disk_mb" -ge "$DISK_THRESHOLD_MB" ]; then
  echo "UNHEALTHY: $DISK_PATH is using ${disk_mb}MB (threshold ${DISK_THRESHOLD_MB}MB)"
  exit 1
fi

echo "HEALTHY: process up, port ${PORT} listening, ${HEALTH_URL} returned 200, disk ${disk_mb}MB/${DISK_THRESHOLD_MB}MB"
exit 0
SCRIPT
chmod +x /workspace/health.sh
```

Each check is independent and short-circuits the moment it fails, so the
script's exit code always matches its single line of output: `0` and
`HEALTHY: ...` together, or non-zero and exactly one `UNHEALTHY: ...`
line naming which specific check failed.
