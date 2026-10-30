#!/usr/bin/env bash
set -uo pipefail

fail() { echo "FAIL: $1"; exit 1; }

ANSWERS=/workspace/answers.env
[ -f "$ANSWERS" ] || fail "/workspace/answers.env not found"

# shellcheck disable=SC1090
source "$ANSWERS"

check_one() {
  local file="$1" varname="$2"
  local recorded actual
  recorded="${!varname:-}"
  [ -n "$recorded" ] || fail "$varname is empty"
  case "$recorded" in
    yes|no) ;;
    *) fail "$varname must be exactly 'yes' or 'no', got '$recorded'" ;;
  esac

  if su -s /bin/bash -c "test -r '$file'" tester 2>/dev/null; then
    actual=yes
  else
    actual=no
  fi

  [ "$recorded" = "$actual" ] || fail "$varname should be '$actual' (tester's real access to $file), got '$recorded'"
}

check_one /workspace/playground/owner-only.txt OWNER_ONLY_READABLE_BY_TESTER
check_one /workspace/playground/group-readable.txt GROUP_READABLE_READABLE_BY_TESTER
check_one /workspace/playground/world-readable.txt WORLD_READABLE_READABLE_BY_TESTER

echo "PASS: all three predictions match tester's actual access"
