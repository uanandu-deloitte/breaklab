#!/usr/bin/env bash
set -euo pipefail

VERSION="v2"
STAGE_DIR="$(mktemp -d)"
ARTIFACTS_DIR=/opt/pipeline/artifacts
OUT="$ARTIFACTS_DIR/app-bundle-${VERSION}.tar.gz"

mkdir -p "$ARTIFACTS_DIR"

echo "$VERSION" > "$STAGE_DIR/VERSION"
cat > "$STAGE_DIR/app.sh" <<EOF
#!/usr/bin/env bash
echo "acme-app ${VERSION} running"
EOF
chmod +x "$STAGE_DIR/app.sh"

tar -czf "$OUT" -C "$STAGE_DIR" .
rm -rf "$STAGE_DIR"

echo "build complete: $OUT"
