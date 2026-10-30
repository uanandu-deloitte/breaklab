#!/usr/bin/env bash
set -uo pipefail

ANSWERS=/workspace/answers.env
LOCKFILE=/srv/app/requirements.lock
CACHE_KEY_SCRIPT=/usr/local/bin/cache-key
ORIG_REF=/opt/reference/cache-key.orig
fail() { echo "FAIL: $1"; exit 1; }

[ -f "$ANSWERS" ] || fail "/workspace/answers.env not found"
[ -x "$CACHE_KEY_SCRIPT" ] || fail "$CACHE_KEY_SCRIPT is missing or not executable"
[ -f "$LOCKFILE" ] || fail "$LOCKFILE is missing — fixture was tampered with"
[ -f "$ORIG_REF" ] || fail "$ORIG_REF is missing — fixture was tampered with"

# shellcheck disable=SC1090
source "$ANSWERS"

# ORIGINAL_KEY_SOURCE: a diagnosis of the untouched reference copy, not
# whatever the learner has since edited cache-key into.
[ -n "${ORIGINAL_KEY_SOURCE:-}" ] || fail "ORIGINAL_KEY_SOURCE is empty"
echo "$ORIGINAL_KEY_SOURCE" | grep -qiE 'static|hardcod|constant|literal' \
  || fail "ORIGINAL_KEY_SOURCE doesn't name a static/hardcoded source — read $ORIG_REF again"
echo "$ORIGINAL_KEY_SOURCE" | grep -qiE 'lockfile|requirements|hash|sha' \
  && fail "ORIGINAL_KEY_SOURCE claims the original script was derived from the lockfile — it wasn't; that's the bug"

# Ground truth: the frozen reference copy really never referenced the
# lockfile at all.
grep -q "requirements.lock" "$ORIG_REF" && fail "fixture problem: $ORIG_REF already references the lockfile"

# CACHE_KEY_NOW_TRACKS_LOCKFILE: the learner's own conclusion after
# fixing (or not fixing) cache-key.
[ -n "${CACHE_KEY_NOW_TRACKS_LOCKFILE:-}" ] || fail "CACHE_KEY_NOW_TRACKS_LOCKFILE is empty"
echo "$CACHE_KEY_NOW_TRACKS_LOCKFILE" | grep -qi '^yes$' \
  || fail "CACHE_KEY_NOW_TRACKS_LOCKFILE should be 'yes' once cache-key is fixed"

# Dynamic check: does editing the lockfile's content actually change
# cache-key's output now? Probe with a temporary edit, then restore.
key_before="$("$CACHE_KEY_SCRIPT")"
cp "$LOCKFILE" "${LOCKFILE}.bak"
echo "# check-probe-$$" >> "$LOCKFILE"
key_after="$("$CACHE_KEY_SCRIPT")"
mv "${LOCKFILE}.bak" "$LOCKFILE"

[ "$key_before" != "$key_after" ] \
  || fail "cache-key prints the same key ('$key_before') before and after editing $LOCKFILE — it's still not derived from the lockfile"

echo "PASS: cache-key now derives its key from requirements.lock's content ($key_before -> $key_after), not a static string"
