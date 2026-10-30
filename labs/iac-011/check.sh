#!/usr/bin/env bash
set -uo pipefail

fail() { echo "FAIL: $1"; exit 1; }

cd /infra || fail "/infra is missing"
[ -f playbook.yml ] || fail "/infra/playbook.yml is missing"
[ -f app.conf.pristine ] || fail "lab setup problem: app.conf.pristine missing"

grep -qiE '^\s*(shell|command|raw):' playbook.yml \
  && fail "playbook.yml still uses a shell/command/raw task -- use a real idempotent module instead (e.g. ansible.builtin.lineinfile)"

# Reset to a known starting point every time this check runs, regardless
# of what's already been tried manually.
cp app.conf.pristine /etc/app/app.conf

run1=$(ansible-playbook -i inventory playbook.yml 2>&1)
echo "$run1" | grep -qE '\bfailed=0\b' || fail "first ansible-playbook run failed:
$run1"

occurrences=$(grep -c '^max_connections=200$' /etc/app/app.conf)
[ "$occurrences" -eq 1 ] || fail "/etc/app/app.conf has $occurrences line(s) matching max_connections=200 after one run, expected exactly 1"

run2=$(ansible-playbook -i inventory playbook.yml 2>&1)
echo "$run2" | grep -qE 'changed=0 .*failed=0' \
  || fail "second ansible-playbook run wasn't idempotent (expected changed=0):
$run2"

occurrences2=$(grep -c '^max_connections=200$' /etc/app/app.conf)
[ "$occurrences2" -eq 1 ] || fail "/etc/app/app.conf has $occurrences2 line(s) matching max_connections=200 after a second run, expected exactly 1"

echo "PASS: playbook.yml idempotently manages max_connections -- a second run makes no changes"
