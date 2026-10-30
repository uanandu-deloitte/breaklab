#!/usr/bin/env bash
set -uo pipefail

deploy_ok=0
expected=$(sed -n 's/^message: *//p' /infra/site.yaml)
actual=$(cat /etc/site/deployed.conf 2>/dev/null)
[ "$expected" = "$actual" ] && deploy_ok=1

release_ok=0
if cd /srv/release-repo 2>/dev/null \
  && git log main --oneline 2>/dev/null | grep -q "RELEASE: v1.2.0" \
  && grep -q "v1.2.0" RELEASE_MESSAGE.md 2>/dev/null; then
  release_ok=1
fi

if [ "$deploy_ok" = "1" ] && [ "$release_ok" = "1" ]; then
  echo "PASS: deploy pipeline healthy and this release's commit is on main"
  exit 0
fi

msg="FAIL:"
if [ "$deploy_ok" = "0" ]; then
  msg="$msg deployed config still doesn't match /infra/site.yaml."
fi
if [ "$release_ok" = "0" ]; then
  msg="$msg this release's commit is still missing from main."
fi
echo "$msg"
exit 1
