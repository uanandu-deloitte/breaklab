#!/usr/bin/env bash
set -euo pipefail
# Reference fix: restore the service's selector to match the pods' real label
kubectl patch service storefront --type merge -p '{"spec":{"selector":{"app":"storefront"}}}'
