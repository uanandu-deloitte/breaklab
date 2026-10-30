#!/usr/bin/env bash
set -euo pipefail
# Oops: the nightly cleanup job stopped running and scratch files piled up
for i in $(seq 1 500); do
  dd if=/dev/zero of="/var/app/data/tmp/scratch-$i" bs=200K count=1 status=none
done
