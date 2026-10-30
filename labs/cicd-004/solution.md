# Solution: Secrets and Config in CI

Replace the hardcoded comparison in `/opt/pipeline/deploy.sh` with a
runtime read of `$DEPLOY_TOKEN`, checked against a hash — never the
plaintext — stored in the image:

```bash
cat <<'EOF' > /opt/pipeline/deploy.sh
#!/usr/bin/env bash
set -euo pipefail

EXPECTED_HASH_FILE=/opt/pipeline/.deploy_token.sha256

if [ -z "${DEPLOY_TOKEN:-}" ]; then
  echo "DEPLOY DENIED: DEPLOY_TOKEN not set" >&2
  exit 1
fi

actual_hash="$(printf '%s' "$DEPLOY_TOKEN" | sha256sum | awk '{print $1}')"
expected_hash="$(cat "$EXPECTED_HASH_FILE")"

if [ "$actual_hash" = "$expected_hash" ]; then
  echo "DEPLOYED" > /workspace/deploy-status.txt
  echo "deploy authorized, shipping release"
else
  echo "DEPLOY DENIED: invalid token" >&2
  exit 1
fi
EOF
chmod +x /opt/pipeline/deploy.sh

export DEPLOY_TOKEN=brk-ci-7f2a9d41c8b3
/opt/pipeline/deploy.sh
```

`deploy.sh` no longer carries any secret literal of its own — it reads
`$DEPLOY_TOKEN` from the environment at run time and compares its hash
against a hash stored in the image, so the real token value is never
written to disk anywhere in the pipeline's own files, only supplied by
whoever runs the deploy.
