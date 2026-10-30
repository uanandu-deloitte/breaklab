#!/usr/bin/env bash
set -uo pipefail

TMPL=/infra/resources.yaml.tmpl
DEST=/etc/app/config.yaml

fail() { echo "FAIL: $1"; exit 1; }

[ -f "$TMPL" ] || fail "/infra/resources.yaml.tmpl is missing"
[ -f "$DEST" ] || fail "$DEST does not exist — has infra-apply ever run?"

render() {
  local vars_file="$1"
  (
    # shellcheck disable=SC1090
    set -a
    source "$vars_file"
    set +a
    sed -e "s/{{ *ENVIRONMENT *}}/${ENVIRONMENT}/g" \
        -e "s/{{ *REPLICA_COUNT *}}/${REPLICA_COUNT}/g" \
        "$TMPL"
  )
}

[ -f /infra/vars/prod.env ] || fail "/infra/vars/prod.env is missing"
[ -f /infra/vars/dev.env ] || fail "/infra/vars/dev.env is missing"

expected_prod=$(render /infra/vars/prod.env)
expected_dev=$(render /infra/vars/dev.env)
actual=$(cat "$DEST")

[ "$actual" = "$expected_prod" ] \
  || fail "$DEST doesn't match what prod's variables would render — wrong (or no) environment applied"

[ "$actual" != "$expected_dev" ] \
  || fail "$DEST still matches dev's rendering — production was never actually applied"

echo "PASS: $DEST reflects production's variables, not dev's"
