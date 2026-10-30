#!/usr/bin/env bash
set -uo pipefail

fail() { echo "FAIL: $1"; exit 1; }

BUILD=/opt/pipeline/build.sh
DEPLOY=/opt/pipeline/deploy.sh
DEPLOYED=/workspace/deployed-version.txt

[ -f "$BUILD" ] || fail "$BUILD is missing"
[ -f "$DEPLOY" ] || fail "$DEPLOY is missing"

# Read the version build.sh actually produces out of build.sh's own
# source, rather than hardcoding a version string here — the check
# should track whatever build.sh really does, not a snapshot of it.
VERSION="$(grep -oE '^VERSION="[^"]+"' "$BUILD" | head -n1 | sed -E 's/^VERSION="([^"]+)"$/\1/')"
[ -n "$VERSION" ] || fail "couldn't read VERSION out of $BUILD — has its build-output naming changed shape?"

ARTIFACT="/opt/pipeline/artifacts/app-bundle-${VERSION}.tar.gz"
[ -f "$ARTIFACT" ] || fail "expected build artifact $ARTIFACT doesn't exist — run build.sh"

[ -f "$DEPLOYED" ] || fail "$DEPLOYED doesn't exist — deploy.sh hasn't successfully deployed anything yet"
deployed_version="$(cat "$DEPLOYED")"
[ "$deployed_version" = "$VERSION" ] || fail "$DEPLOYED says '$deployed_version' but build.sh actually produced '$VERSION'"

[ -x /opt/app/app.sh ] || fail "/opt/app/app.sh is missing or not executable — the bundle wasn't extracted correctly"
grep -q "acme-app ${VERSION} running" /opt/app/app.sh || fail "/opt/app/app.sh doesn't contain the expected bundle content for version $VERSION"

[ -f /opt/app/VERSION ] || fail "/opt/app/VERSION is missing — the bundle wasn't extracted correctly"
[ "$(cat /opt/app/VERSION)" = "$VERSION" ] || fail "/opt/app/VERSION doesn't match the version build.sh produced"

echo "PASS: deploy.sh correctly found and deployed the artifact build.sh actually produced ($VERSION)"
