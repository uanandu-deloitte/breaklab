#!/usr/bin/env bash
set -euo pipefail
REPO=/srv/app-repo
SRC=src/app.sh
HASH_SOURCE=src/app.sh   # the file the cache key is derived from
ARTIFACT_DIR=/opt/pipeline/artifacts
HASH_FILE=/opt/pipeline/.last-built-hash
DEPLOYED=/srv/deployed/current-build.txt
REF_FILE=/opt/pipeline/BUILD_REF

build_ref=$(cat "$REF_FILE" 2>/dev/null || echo "release")

cd "$REPO"

if ! git rev-parse -q --verify "${build_ref}" >/dev/null 2>&1; then
  echo "build: ref '$build_ref' does not exist, nothing to build" >&2
  exit 1
fi

content=$(git show "${build_ref}:${SRC}" 2>/dev/null) || {
  echo "build: $SRC not found at ref $build_ref" >&2
  exit 1
}
hash_input=$(git show "${build_ref}:${HASH_SOURCE}" 2>/dev/null || echo "")

# Cache key: derived from whatever HASH_SOURCE points at. Skipping the
# rebuild when this key hasn't changed saves repackaging an identical
# artifact -- but only if HASH_SOURCE is actually the application source.
current_hash=$(printf '%s' "$hash_input" | sha256sum | cut -d' ' -f1)
stored_hash=$(cat "$HASH_FILE" 2>/dev/null || echo "")

if [ "$current_hash" = "$stored_hash" ] && [ -f "$DEPLOYED" ]; then
  echo "build: cache key unchanged, reusing cached artifact (built from $build_ref)"
  exit 0
fi

mkdir -p "$ARTIFACT_DIR"
version="${build_ref}-$(printf '%s' "$content" | sha256sum | cut -c1-8)"
workdir=$(mktemp -d)
printf '%s\n' "$content" > "$workdir/app.sh"
tar -C "$workdir" -czf "$ARTIFACT_DIR/app-${version}.tar.gz" app.sh
rm -rf "$workdir"

mkdir -p "$(dirname "$DEPLOYED")"
printf '%s\n' "$content" > "$DEPLOYED"
echo "$current_hash" > "$HASH_FILE"
echo "build: rebuilt artifact app-${version}.tar.gz from $build_ref (cache key from $HASH_SOURCE)"
