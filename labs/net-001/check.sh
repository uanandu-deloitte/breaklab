#!/usr/bin/env bash
set -uo pipefail

ANSWERS=/workspace/answers.env
fail() { echo "FAIL: $1"; exit 1; }

[ -f "$ANSWERS" ] || fail "/workspace/answers.env not found"

# shellcheck disable=SC1090
source "$ANSWERS"

live_iface=$(ip route show default | awk '{print $5; exit}')
[ -n "$live_iface" ] || fail "could not determine this host's default-route interface (test setup problem)"

[ -n "${INTERFACE_NAME:-}" ] || fail "INTERFACE_NAME is empty"
[ "$INTERFACE_NAME" = "$live_iface" ] || fail "INTERFACE_NAME ($INTERFACE_NAME) doesn't match the interface carrying the default route ($live_iface)"

live_ip=$(ip -4 -o addr show "$live_iface" | awk '{print $4}' | cut -d/ -f1 | head -1)
[ -n "$live_ip" ] || fail "could not determine this host's live IPv4 address (test setup problem)"

[ -n "${CONTAINER_IP:-}" ] || fail "CONTAINER_IP is empty"
[ "$CONTAINER_IP" = "$live_ip" ] || fail "CONTAINER_IP ($CONTAINER_IP) doesn't match this container's live address on $live_iface ($live_ip)"

live_ports=$(ss -ltn | awk 'NR>1 {print $4}' | sed 's/.*://' | sort -un | paste -sd, -)
[ -n "$live_ports" ] || fail "could not determine this host's live listening ports (test setup problem)"

[ -n "${LISTENING_PORTS:-}" ] || fail "LISTENING_PORTS is empty"
[ "$LISTENING_PORTS" = "$live_ports" ] || fail "LISTENING_PORTS ($LISTENING_PORTS) doesn't match the live listening ports ($live_ports)"

echo "PASS: answers.env matches this host's live network identity"
