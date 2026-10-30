#!/usr/bin/env bash
set -euo pipefail
# Reference fix: the bad tag was never a real image; roll back to the
# last known-good revision rather than guessing at a correct tag.
kubectl rollout undo deployment/catalog
kubectl rollout status deployment/catalog --timeout=30s
