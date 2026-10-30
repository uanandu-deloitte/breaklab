#!/usr/bin/env bash
set -uo pipefail

fail() { echo "FAIL: $1"; exit 1; }

run() {
  local desc="$1"
  shift
  if ! "$@" >/dev/null 2>&1; then
    fail "$desc"
  fi
}

run "curl is not installed or not working (curl --version)" curl --version
run "iproute2's ip is not installed or not working (ip -V)" ip -V
run "lsof is not installed or not working (lsof -v)" lsof -v
run "jq is not installed or not working (jq --version)" jq --version

echo "PASS: curl, iproute2 (ip), lsof, and jq are all installed and working"
