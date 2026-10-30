#!/usr/bin/env bash
set -uo pipefail

fail() { echo "FAIL: $1"; exit 1; }

timeout 3 getent hosts api.internal >/dev/null 2>&1 || \
  fail "api.internal does not resolve — check /etc/hosts"

body=$(curl -fsS --max-time 3 http://api.internal:8080/) || \
  fail "curl http://api.internal:8080/ failed"

grep -q '"service": "api-internal"' <<<"$body" || \
  fail "api.internal resolved but the response body is wrong: $body"

echo "PASS: api.internal resolves and the backend responds correctly"
