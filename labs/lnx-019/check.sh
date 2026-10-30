#!/usr/bin/env bash
set -uo pipefail

SCRIPT=/opt/app/start.sh
fail() { echo "FAIL: $1"; exit 1; }

[ -f "$SCRIPT" ] || fail "$SCRIPT not found"

out_unset=$(env -u APP_CONFIG_PATH bash "$SCRIPT" 2>&1)
code_unset=$?

if [ "$code_unset" -eq 0 ]; then
  fail "start.sh exited 0 with APP_CONFIG_PATH unset — it must fail loudly, not silently succeed"
fi
echo "$out_unset" | grep -qi "error" || \
  fail "start.sh's failure (APP_CONFIG_PATH unset) didn't print a clear error message, got: $out_unset"

out_set=$(APP_CONFIG_PATH=/opt/app/config.yaml bash "$SCRIPT" 2>&1)
code_set=$?

[ "$code_set" -eq 0 ] || fail "start.sh exited $code_set with a valid APP_CONFIG_PATH — it should succeed, got: $out_set"
echo "$out_set" | grep -qi "started successfully" || \
  fail "start.sh didn't print its success marker with a valid config, got: $out_set"

echo "PASS: start.sh fails loudly with no config and succeeds truthfully with a valid one"
