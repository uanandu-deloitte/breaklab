#!/usr/bin/env bash
set -uo pipefail

fail() { echo "FAIL: $1"; exit 1; }

BASE=/opt/acme-app

[ -f "$BASE/bin/acme-app" ] || fail "$BASE/bin/acme-app is missing"
[ -x "$BASE/bin/acme-app" ] || fail "$BASE/bin/acme-app exists but isn't executable"
grep -q "acme-app running" "$BASE/bin/acme-app" || fail "$BASE/bin/acme-app content doesn't match the original executable"

[ -f "$BASE/config/app.conf" ] || fail "$BASE/config/app.conf is missing"
grep -q "service: acme-app" "$BASE/config/app.conf" || fail "$BASE/config/app.conf content doesn't match the original config"

[ -f "$BASE/data/data.db" ] || fail "$BASE/data/data.db is missing"
grep -q "ACME-DATA-V1" "$BASE/data/data.db" || fail "$BASE/data/data.db content doesn't match the original data file"

# The files should have been moved, not copied — the wrong locations
# shouldn't still hold a duplicate.
[ -e "$BASE/config/acme-app" ] && fail "$BASE/config/acme-app still exists — move it, don't leave a copy behind"
[ -e "$BASE/data/app.conf" ] && fail "$BASE/data/app.conf still exists — move it, don't leave a copy behind"
[ -e "$BASE/bin/data.db" ] && fail "$BASE/bin/data.db still exists — move it, don't leave a copy behind"

echo "PASS: every file is in its correct directory with its original content intact"
