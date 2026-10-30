#!/usr/bin/env bash
set -uo pipefail

fail() { echo "FAIL: $1"; exit 1; }

[ -f /infra/main.bicep ] || fail "/infra/main.bicep is missing"

out=$(az bicep build --file /infra/main.bicep --stdout 2>/tmp/bicep-err.log)
build_exit=$?
[ "$build_exit" -eq 0 ] || fail "az bicep build failed:
$(cat /tmp/bicep-err.log)"

type=$(echo "$out" | jq -r '.parameters.retentionDays.type')
[ "$type" = "int" ] || fail "retentionDays compiled as type \"$type\", expected \"int\" -- don't just change its declared type to string"

default=$(echo "$out" | jq -r '.parameters.retentionDays.defaultValue')
[ "$default" = "30" ] || fail "retentionDays' compiled defaultValue is \"$default\", expected 30"

echo "PASS: main.bicep compiles and retentionDays is a real int parameter"
