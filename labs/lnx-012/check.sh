#!/usr/bin/env bash
set -uo pipefail

fail() { echo "FAIL: $1"; exit 1; }

id developer >/dev/null 2>&1 || fail "developer user not found"
id operator >/dev/null 2>&1 || fail "operator user not found"

dev_uid=$(id -u developer)
op_uid=$(id -u operator)

[ "$dev_uid" != "$op_uid" ] || fail "developer and operator must have distinct UIDs"
[ "$dev_uid" != "0" ] || fail "developer must not be UID 0"
[ "$op_uid" != "0" ] || fail "operator must not be UID 0"

dev_shell=$(getent passwd developer | cut -d: -f7)
[ "$dev_shell" = "/bin/bash" ] || fail "developer's shell should be /bin/bash, got $dev_shell"

[ -d /home/developer ] || fail "/home/developer does not exist"

op_shell=$(getent passwd operator | cut -d: -f7)
case "$op_shell" in
  /usr/sbin/nologin|/sbin/nologin) ;;
  *) fail "operator's shell should be a nologin shell, got $op_shell" ;;
esac

echo "PASS: developer and operator exist with the expected identities"
