# Solution: Why Did the Script Fail?

```bash
cat <<'EOF' > /opt/app/start.sh
#!/usr/bin/env bash
set -euo pipefail

if [ -z "${APP_CONFIG_PATH:-}" ]; then
  echo "ERROR: APP_CONFIG_PATH is not set; cannot start" >&2
  exit 1
fi

if [ ! -f "$APP_CONFIG_PATH" ]; then
  echo "ERROR: configuration file not found" >&2
  exit 1
fi

if ! grep -q "service:" "$APP_CONFIG_PATH"; then
  echo "ERROR: configuration file is invalid" >&2
  exit 1
fi

echo "Service started successfully"
exit 0
EOF
chmod +x /opt/app/start.sh
```

Two independent fixes:

1. **Missing precondition.** An explicit `[ -z "${APP_CONFIG_PATH:-}" ]`
   check runs before the variable is ever used, failing loudly with a
   clear message instead of letting an unset variable flow silently into
   `cat`.
2. **Masked exit status.** The original piped `cat ... | grep -q ...`
   and never looked at the result. The fix runs `grep -q` directly
   against the file (no pipe) inside an `if` that checks its real exit
   status, and `exit 1` on failure instead of always falling through to
   `exit 0`.
