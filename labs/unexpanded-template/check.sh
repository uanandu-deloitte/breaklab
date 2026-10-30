#!/usr/bin/env bash
set -uo pipefail
f=/etc/worker/worker.env
if [ ! -f "$f" ]; then
  echo "FAIL: $f missing"
  exit 1
fi
if grep -q '\${' "$f"; then
  echo "FAIL: an unexpanded template variable leaked into $f"
  exit 1
fi
if grep -q "^REPLICAS=3$" "$f" && grep -q "^QUEUE=orders$" "$f"; then
  echo "PASS: worker config rendered correctly"
else
  echo "FAIL: worker config values don't match /infra/worker.yaml"
  exit 1
fi
