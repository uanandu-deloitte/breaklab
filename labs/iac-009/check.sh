#!/usr/bin/env bash
set -uo pipefail

fail() { echo "FAIL: $1"; exit 1; }

cd /infra || fail "/infra is missing"
[ -f main.tf ] || fail "/infra/main.tf is missing"

declared_raw=$(sed -n 's/^\s*content\s*=\s*"\(.*\)"$/\1/p' main.tf | head -1)
[ -n "$declared_raw" ] || fail "couldn't read the declared content attribute from main.tf"
declared=$(printf '%b' "${declared_raw}")

[ -f /etc/motd-app/banner.txt ] || fail "/etc/motd-app/banner.txt does not exist"
actual=$(cat /etc/motd-app/banner.txt)

[ "$actual" = "$declared" ] || fail "/etc/motd-app/banner.txt doesn't match what main.tf declares -- it's still drifted"

# Prove this was reconciled through Terraform, not a hand-edit of the
# real file: plan must report zero pending changes once state, config,
# and the real world all agree.
plan_out=$(terraform plan -input=false -detailed-exitcode 2>&1)
plan_exit=$?
[ "$plan_exit" -eq 0 ] || fail "terraform plan still reports pending changes (exit $plan_exit): $plan_out"

echo "PASS: /etc/motd-app/banner.txt matches main.tf and terraform plan reports no pending changes"
