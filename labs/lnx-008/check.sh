#!/usr/bin/env bash
set -uo pipefail

fail() { echo "FAIL: $1"; exit 1; }

ACTIVE=/opt/app/config/app.conf
ARCHIVED=/opt/app/releases/v1/config/app.conf
CHANGELOG=/opt/app/CHANGELOG.md
NOTES=/opt/app/docs/migration-notes.md

[ -f "$ACTIVE" ] || fail "$ACTIVE is missing"
[ -f "$ARCHIVED" ] || fail "$ARCHIVED is missing — it must stay untouched, not deleted"
[ -f "$CHANGELOG" ] || fail "$CHANGELOG is missing — it must stay untouched, not deleted"
[ -f "$NOTES" ] || fail "$NOTES is missing — it must stay untouched, not deleted"

grep -qF "database.internal" "$ACTIVE" || fail "the active config at $ACTIVE doesn't reference database.internal yet"
grep -qF "database-old.internal" "$ACTIVE" && fail "the active config at $ACTIVE still references the old hostname"

# The important assertion: historical/archived files must still mention
# the old hostname unchanged. A learner who ran a blind find-and-replace
# across the whole /opt/app tree fails here even though the active
# config above now looks correct.
grep -qF "database-old.internal" "$ARCHIVED" || fail "$ARCHIVED no longer mentions database-old.internal — it's an archived snapshot and must not be edited"
grep -qF "database-old.internal" "$CHANGELOG" || fail "$CHANGELOG no longer mentions database-old.internal — it's a historical record and must not be edited"
grep -qF "database-old.internal" "$NOTES" || fail "$NOTES no longer mentions database-old.internal — it's a historical record and must not be edited"

echo "PASS: only the active config was updated; archived and historical references are untouched"
