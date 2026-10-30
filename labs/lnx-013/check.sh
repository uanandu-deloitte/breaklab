#!/usr/bin/env bash
set -uo pipefail

fail() { echo "FAIL: $1"; exit 1; }

getent group ops-team >/dev/null 2>&1 || fail "ops-team group not found"

groups developer 2>/dev/null | grep -qw ops-team || fail "developer is not a member of ops-team"
groups operator 2>/dev/null | grep -qw ops-team || fail "operator is not a member of ops-team"

[ -d /srv/team ] || fail "/srv/team does not exist"

group=$(stat -c '%G' /srv/team)
[ "$group" = "ops-team" ] || fail "/srv/team's group should be ops-team, got $group"

perms=$(stat -c '%A' /srv/team)
group_triplet=${perms:4:3}
[ "$group_triplet" = "rws" ] || fail "/srv/team needs group read+write+exec plus setgid (group bits: $group_triplet)"

rm -f /srv/team/testfile 2>/dev/null
if ! su -s /bin/bash -c "touch /srv/team/testfile" developer 2>/dev/null; then
  fail "developer could not create a file in /srv/team"
fi

file_group=$(stat -c '%G' /srv/team/testfile)
[ "$file_group" = "ops-team" ] || fail "file created by developer in /srv/team should belong to ops-team, got $file_group"

echo "PASS: ops-team is shared correctly and /srv/team inherits its group"
