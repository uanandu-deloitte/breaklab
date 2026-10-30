#!/usr/bin/env bash
set -uo pipefail

ANSWERS=/workspace/answers.env
fail() { echo "FAIL: $1"; exit 1; }

[ -f "$ANSWERS" ] || fail "/workspace/answers.env not found"

# shellcheck disable=SC1090
source "$ANSWERS"

# Read the real baked-in values dynamically from this running container's
# own environment, never hardcoded, so this check stays correct even if
# the Dockerfile's values are tweaked later.
real_app_env=$(printenv APP_ENV || true)
real_log_level=$(printenv LOG_LEVEL || true)
real_app_version=$(printenv APP_VERSION || true)

[ -n "${APP_ENV:-}" ] || fail "APP_ENV is empty"
[ "$APP_ENV" = "$real_app_env" ] || fail "APP_ENV ($APP_ENV) doesn't match this container's real environment ($real_app_env)"

[ -n "${LOG_LEVEL:-}" ] || fail "LOG_LEVEL is empty"
[ "$LOG_LEVEL" = "$real_log_level" ] || fail "LOG_LEVEL ($LOG_LEVEL) doesn't match this container's real environment ($real_log_level)"

[ -n "${APP_VERSION:-}" ] || fail "APP_VERSION is empty"
[ "$APP_VERSION" = "$real_app_version" ] || fail "APP_VERSION ($APP_VERSION) doesn't match this container's real environment ($real_app_version)"

# Also verify straight from /proc/1/environ, not just the shell's printenv
# view, since those are two different sources that should agree.
proc_env=$(tr '\0' '\n' < /proc/1/environ)
grep -qx "APP_ENV=$real_app_env" <<<"$proc_env" || fail "/proc/1/environ doesn't show the same APP_ENV printenv reports"
grep -qx "LOG_LEVEL=$real_log_level" <<<"$proc_env" || fail "/proc/1/environ doesn't show the same LOG_LEVEL printenv reports"
grep -qx "APP_VERSION=$real_app_version" <<<"$proc_env" || fail "/proc/1/environ doesn't show the same APP_VERSION printenv reports"

[ -n "${OVERRIDE_CANDIDATE:-}" ] || fail "OVERRIDE_CANDIDATE is empty"
case "$OVERRIDE_CANDIDATE" in
  LOG_LEVEL|APP_ENV) ;;
  APP_VERSION) fail "APP_VERSION identifies the build itself — that's a reasonable thing to bake in, not the override candidate" ;;
  *) fail "OVERRIDE_CANDIDATE ($OVERRIDE_CANDIDATE) isn't one of the three baked-in keys" ;;
esac

echo "PASS: answers.env correctly reports the baked-in environment and a defensible override candidate"
