#!/usr/bin/env bash
set -uo pipefail

fail() { echo "FAIL: $1"; exit 1; }

HOSTS=/workspace/hosts.txt
OUT=/workspace/normalized.txt

[ -f "$HOSTS" ] || fail "/workspace/hosts.txt not found"
[ -f "$OUT" ] || fail "/workspace/normalized.txt not found"

expected=$(cat <<'EOF'
app-01.example.com
cache-01.example.com
db-01.example.com
lb-01.example.com
web-01.example.com
web-02.example.com
EOF
)

actual=$(cat "$OUT")

if [ "$actual" != "$expected" ]; then
  fail "normalized.txt doesn't match the expected canonical list:
--- expected ---
$expected
--- got ---
$actual"
fi

echo "PASS: normalized.txt is lowercase, sorted, deduplicated, and matches the canonical host list"
