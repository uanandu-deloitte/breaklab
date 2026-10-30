#!/usr/bin/env bash
set -euo pipefail
# Reference fix: the source was never wrong; reconcile the deployed state back to it
infra-apply
