#!/usr/bin/env bash
set -euo pipefail
kubectl patch service orders -p '{"spec":{"selector":{"app":"orders"}}}'
