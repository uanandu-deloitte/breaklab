#!/usr/bin/env bash
set -euo pipefail
# Reference fix: the internal registry isn't reachable from this
# cluster; roll back to the last known-good revision rather than
# standing up a registry.
kubectl rollout undo deployment/catalog
kubectl rollout status deployment/catalog --timeout=30s
