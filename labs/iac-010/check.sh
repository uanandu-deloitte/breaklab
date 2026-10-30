#!/usr/bin/env bash
set -uo pipefail

fail() { echo "FAIL: $1"; exit 1; }

cd /infra || fail "/infra is missing"

state_addrs=$(terraform state list 2>/dev/null)
echo "$state_addrs" | grep -qx "local_file.privacy_notice" \
  || fail "local_file.privacy_notice is not in Terraform state -- the renamed resource was never reconciled"
echo "$state_addrs" | grep -qx "local_file.legal_notice" \
  && fail "local_file.legal_notice is still in state -- it should have been renamed to privacy_notice, not left behind as a duplicate"

FILE=/var/www/legal/notice.txt
[ -f "$FILE" ] || fail "$FILE does not exist"

content=$(cat "$FILE")
expected=$(printf 'All rights reserved, Acme Corp.\n')
[ "$content" = "$expected" ] || fail "$FILE's content changed -- it should be untouched by a rename-only config change"

[ -f /infra/.original-mtime ] || fail "/infra/.original-mtime is missing (lab setup problem)"
original_mtime=$(cat /infra/.original-mtime)
current_mtime=$(stat -c %Y "$FILE")
[ "$current_mtime" = "$original_mtime" ] \
  || fail "$FILE's mtime changed ($current_mtime vs $original_mtime) -- it was destroyed and recreated instead of being renamed in Terraform's state"

plan_out=$(terraform plan -input=false -detailed-exitcode 2>&1)
plan_exit=$?
[ "$plan_exit" -eq 0 ] || fail "terraform plan still reports pending changes (exit $plan_exit): $plan_out"

echo "PASS: the resource was renamed to local_file.privacy_notice in state without ever destroying the real file"
