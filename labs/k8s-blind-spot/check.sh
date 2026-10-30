#!/usr/bin/env bash
set -uo pipefail

# --- K8s half: is the checkout Service actually routing to a Ready
# backend right now? Measured directly, never read from a cached file. ---
k8s_ok=0
real_healthy=0
ip=$(kubectl get svc checkout -o jsonpath='{.spec.clusterIP}' 2>/dev/null)
port=$(kubectl get svc checkout -o jsonpath='{.spec.ports[0].port}' 2>/dev/null)
if [ -n "$ip" ] && [ -n "$port" ] && curl -fsS --max-time 3 "http://${ip}:${port}/" >/dev/null 2>&1; then
  real_healthy=1
  k8s_ok=1
fi

# --- Observability half, part 1: force a refresh right now and confirm it
# actually produced a current, accurate reading. This catches a frozen or
# disconnected metrics pipeline independent of whether the outage itself
# is fixed. ---
refresh-metrics >/dev/null 2>&1

fresh=0
accurate=0
if [ -f /var/lib/metrics/current.json ]; then
  observed_value=$(jq -r '.checkout_healthy // empty' /var/lib/metrics/current.json 2>/dev/null)
  observed_at=$(jq -r '.observed_at // empty' /var/lib/metrics/current.json 2>/dev/null)
  now=$(date +%s)
  if [[ "$observed_at" =~ ^[0-9]+$ ]]; then
    age=$(( now - observed_at ))
    if [ "$age" -ge 0 ] && [ "$age" -le 120 ]; then
      fresh=1
    fi
  fi
  if [ "$observed_value" = "$real_healthy" ]; then
    accurate=1
  fi
fi

# --- Observability half, part 2: does the alert rule itself correctly
# tell a real breach (0) apart from healthy (1)? Tested against synthetic
# values so a rule that merely happens not to have fired yet can't pass
# this by coincidence of whatever the live metric is right now. ---
rule_ok=0
if [ -f /etc/alerts/rules.yaml ]; then
  op=$(sed -n 's/^op: *//p' /etc/alerts/rules.yaml | tr -d '"')
  threshold=$(sed -n 's/^threshold: *//p' /etc/alerts/rules.yaml)
  fires() {
    local v="$1" fire=0
    case "$op" in
      ">")  [ "$v" -gt "$threshold" ] 2>/dev/null && fire=1 ;;
      ">=") [ "$v" -ge "$threshold" ] 2>/dev/null && fire=1 ;;
      "<")  [ "$v" -lt "$threshold" ] 2>/dev/null && fire=1 ;;
      "<=") [ "$v" -le "$threshold" ] 2>/dev/null && fire=1 ;;
    esac
    echo "$fire"
  }
  if [ "$(fires 0)" = "1" ] && [ "$(fires 1)" = "0" ]; then
    rule_ok=1
  fi
fi

observability_ok=0
if [ "$fresh" = "1" ] && [ "$accurate" = "1" ] && [ "$rule_ok" = "1" ]; then
  observability_ok=1
fi

if [ "$k8s_ok" = "1" ] && [ "$observability_ok" = "1" ]; then
  echo "PASS: checkout is actually routing traffic, and monitoring correctly reflects that"
  exit 0
fi

msg="FAIL:"
if [ "$k8s_ok" = "0" ]; then
  msg="$msg checkout Service still has no healthy backend to route to."
fi
if [ "$observability_ok" = "0" ]; then
  msg="$msg monitoring does not currently reflect reality (stale/inaccurate metric, or the alert rule can't tell a breach from healthy)."
fi
echo "$msg"
exit 1
