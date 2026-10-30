# Solution: Artifacts Between Stages

Run the build step and see what it actually names its output:

```bash
/opt/pipeline/build.sh
ls /opt/pipeline/artifacts/
```

`deploy.sh` is hardcoded to look for `app-bundle.tar.gz`, but `build.sh`
now produces `app-bundle-v2.tar.gz`. Fix `deploy.sh` to discover whatever
artifact actually exists instead of hardcoding last release's filename:

```bash
cat <<'EOF' > /opt/pipeline/deploy.sh
#!/usr/bin/env bash
set -euo pipefail

ARTIFACTS_DIR=/opt/pipeline/artifacts
APP_DIR=/opt/app

ARTIFACT="$(ls "$ARTIFACTS_DIR"/app-bundle-*.tar.gz 2>/dev/null | head -n1)"

if [ -z "$ARTIFACT" ] || [ ! -f "$ARTIFACT" ]; then
  echo "ERROR: no build artifact found in $ARTIFACTS_DIR" >&2
  exit 1
fi

rm -rf "$APP_DIR"
mkdir -p "$APP_DIR"
tar -xzf "$ARTIFACT" -C "$APP_DIR"

version="$(cat "$APP_DIR/VERSION")"
echo "$version" > /workspace/deployed-version.txt
echo "deployed $version"
EOF
chmod +x /opt/pipeline/deploy.sh

/opt/pipeline/deploy.sh
```

`deploy.sh` now globs for whatever versioned artifact actually exists in
`/opt/pipeline/artifacts/` instead of hunting for one hardcoded filename —
so the pipeline keeps working the next time `build.sh`'s naming changes
again, not just this once.
