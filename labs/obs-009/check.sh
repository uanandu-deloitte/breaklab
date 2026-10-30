#!/usr/bin/env bash
set -uo pipefail

fail() { echo "FAIL: $1"; exit 1; }

for i in $(seq 1 60); do
  n=$(curl -s "http://localhost:9090/api/v1/query?query=up" 2>/dev/null | jq -r '.data.result | length' 2>/dev/null)
  [ "${n:-0}" != "0" ] && break
  sleep 1
done
[ "${n:-0}" != "0" ] || fail "Prometheus never scraped anything -- it isn't up yet"

for i in $(seq 1 60); do
  curl -sf http://localhost:3000/api/health >/dev/null 2>&1 && break
  sleep 1
done
curl -sf http://localhost:3000/api/health >/dev/null 2>&1 || fail "Grafana isn't up yet"

dash_json=$(curl -s "http://localhost:3000/api/dashboards/uid/checkout-overview")
expr=$(echo "$dash_json" | jq -r '.dashboard.panels[0].targets[0].expr')
[ -n "$expr" ] && [ "$expr" != "null" ] || fail "couldn't read the Checkout overview dashboard's panel query from Grafana's API"

result=$(curl -s "http://localhost:9090/api/v1/query" --data-urlencode "query=$expr")
count=$(echo "$result" | jq -r '.data.result | length')
[ "$count" != "0" ] || fail "the dashboard panel's query ($expr) returns no data from Prometheus -- it's still querying a metric name nothing exports"

value=$(echo "$result" | jq -r '.data.result[0].value[1]')
[ "$value" = "4321" ] || fail "the dashboard panel's query returned an unexpected value: $value"

echo "PASS: the Checkout overview dashboard's panel query matches the real metric Prometheus is scraping"
