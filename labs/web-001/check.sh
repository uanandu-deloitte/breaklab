#!/usr/bin/env bash
set -uo pipefail
port="${BREAKLAB_PORT_80:-8080}"
fail() { echo "FAIL: $1"; exit 1; }

body_file=$(mktemp)
headers=$(curl -fsS --max-time 3 -D - -o "$body_file" "http://localhost:${port}/status/") || {
  rm -f "$body_file"; fail "GET /status/ failed or returned an HTTP error"
}
body=$(cat "$body_file")
rm -f "$body_file"

echo "$headers" | grep -qi '^Content-Type: *text/html' || \
  fail "response is missing a Content-Type: text/html header"

expected='<!doctype html>
<html>
<head><title>Service Status</title></head>
<body><h1>All systems operational</h1></body>
</html>'

[ "$body" = "$expected" ] || fail "response body doesn't match the expected status page exactly"

echo "PASS: /status/ serves the expected status page"
