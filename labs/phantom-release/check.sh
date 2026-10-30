#!/usr/bin/env bash
set -uo pipefail
cd /srv/app-repo 2>/dev/null || { echo "FAIL: app repo missing"; exit 1; }

# The hotfix's real content, found live by what it actually changed,
# wherever it's reachable in history -- never a hardcoded string, so this
# works no matter which ref or path actually fixed it.
hotfix_sha=$(git log --all --oneline --grep='^HOTFIX' 2>/dev/null | awk '{print $1}' | head -n1)
if [ -z "$hotfix_sha" ]; then
  echo "FAIL: no commit matching the hotfix is reachable anywhere in this repo"
  exit 1
fi
hotfix_content=$(git show "${hotfix_sha}:src/app.sh" 2>/dev/null)

build_ref=$(cat /opt/pipeline/BUILD_REF 2>/dev/null || echo "release")

ref_ok=0
ref_content=$(git show "${build_ref}:src/app.sh" 2>/dev/null)
if [ -n "$hotfix_content" ] && [ "$ref_content" = "$hotfix_content" ]; then
  ref_ok=1
fi

build_ok=0
/opt/pipeline/build.sh >/tmp/phantom-release-build.log 2>&1
deployed_content=$(cat /srv/deployed/current-build.txt 2>/dev/null)
if [ -n "$hotfix_content" ] && [ "$deployed_content" = "$hotfix_content" ]; then
  build_ok=1
fi

if [ "$ref_ok" = "1" ] && [ "$build_ok" = "1" ]; then
  echo "PASS: '$build_ref' contains the hotfix and the deployed build matches it"
  exit 0
fi

msg="FAIL:"
if [ "$ref_ok" = "0" ]; then
  msg="$msg '$build_ref' (what the pipeline actually builds from) still doesn't contain the hotfix's real change."
fi
if [ "$build_ok" = "0" ]; then
  msg="$msg the deployed build artifact still doesn't match the hotfix content."
fi
echo "$msg"
exit 1
