#!/usr/bin/env bash
set -euo pipefail

ARTIFACTS_DIR=/opt/pipeline/artifacts
# Last release's artifact was a plain app-bundle.tar.gz. build.sh has
# since moved to versioned filenames, but this was never updated to
# match.
ARTIFACT="$ARTIFACTS_DIR/app-bundle.tar.gz"
APP_DIR=/opt/app

if [ ! -f "$ARTIFACT" ]; then
  echo "ERROR: expected artifact not found: $ARTIFACT" >&2
  exit 1
fi

rm -rf "$APP_DIR"
mkdir -p "$APP_DIR"
tar -xzf "$ARTIFACT" -C "$APP_DIR"

version="$(cat "$APP_DIR/VERSION")"
echo "$version" > /workspace/deployed-version.txt
echo "deployed $version"
