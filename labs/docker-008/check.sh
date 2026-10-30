#!/usr/bin/env bash
# Runs on the host (run_in: host): the final proof is whether the mapped
# host port actually reaches the backend, which can only be tested from
# outside the container.
set -uo pipefail

container="${BREAKLAB_CONTAINER:?BREAKLAB_CONTAINER not set}"
port="${BREAKLAB_PORT_80:-8080}"
fail() { echo "FAIL: $1"; exit 1; }

answers_raw=$(docker exec "$container" cat /workspace/answers.env 2>/dev/null) \
  || fail "/workspace/answers.env not found in the container"

tmpfile=$(mktemp)
printf '%s\n' "$answers_raw" > "$tmpfile"
# shellcheck disable=SC1090
source "$tmpfile"
rm -f "$tmpfile"

[ -n "${BACKEND_REAL_PORT:-}" ] || fail "BACKEND_REAL_PORT is empty"
[ -n "${ORIGINAL_PROXY_PASS_PORT:-}" ] || fail "ORIGINAL_PROXY_PASS_PORT is empty"
[ -n "${PROXY_FIXED:-}" ] || fail "PROXY_FIXED is empty"

# nginx's own config syntax check has to pass before anything else matters.
docker exec "$container" nginx -t >/dev/null 2>&1 \
  || fail "nginx -t fails — fix the config's syntax first"

# Ground truth: ask ss for the backend's real listening port, rather than
# hardcoding one here.
real_backend_port=$(docker exec "$container" sh -c \
  "ss -ltnp 2>/dev/null | grep python3 | grep -oE ':[0-9]+' | head -1 | tr -d ':'")
[ -n "$real_backend_port" ] || fail "couldn't find a live python3 listening port in the container — is the backend still running?"
[ "$BACKEND_REAL_PORT" = "$real_backend_port" ] \
  || fail "BACKEND_REAL_PORT ($BACKEND_REAL_PORT) doesn't match the backend's real listening port ($real_backend_port)"

# Ground truth: nginx's current proxy_pass must now point at that same port.
configured_port=$(docker exec "$container" sh -c \
  "grep -oE 'proxy_pass http://127.0.0.1:[0-9]+' /etc/nginx/sites-available/default | grep -oE '[0-9]+$'")
[ -n "$configured_port" ] || fail "couldn't read nginx's current proxy_pass port"
[ "$configured_port" = "$real_backend_port" ] \
  || fail "nginx's proxy_pass still points at $configured_port, not the backend's real port $real_backend_port — fix hasn't been applied/reloaded"

[ "$ORIGINAL_PROXY_PASS_PORT" != "$real_backend_port" ] \
  || fail "ORIGINAL_PROXY_PASS_PORT should be the wrong port nginx started with, not the backend's real port"

# Ground truth: the mapped host port must actually reach the backend now.
body=$(curl -fsS --max-time 3 "http://localhost:${port}/") \
  || fail "the mapped host port ${port} is still unreachable"
grep -q "BreakLab docker-008 backend OK" <<<"$body" \
  || fail "the mapped host port answered, but not with the backend's real content"

echo "$PROXY_FIXED" | grep -qi '^yes$' \
  || fail "PROXY_FIXED should be yes — the live curl above confirms the proxy now reaches the backend"

echo "PASS: nginx's proxy_pass now points at the backend's real port ($real_backend_port) and the mapped host port reaches it"
