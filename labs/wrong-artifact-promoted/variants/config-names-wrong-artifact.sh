#!/usr/bin/env bash
set -euo pipefail
# New release cycle: both staging and prod get rebuilt and retagged as build 43
mkdir -p /tmp/build-prod-43 /tmp/build-staging-43
echo "prod build 43: checkout-v3.1 (security patch)" > /tmp/build-prod-43/CONTENT
echo "staging build 43: checkout-v5-experimental (unreleased)" > /tmp/build-staging-43/CONTENT
tar -cf /srv/artifacts/app-prod-43.tar -C /tmp/build-prod-43 CONTENT
tar -cf /srv/artifacts/app-staging-43.tar -C /tmp/build-staging-43 CONTENT
echo "app-prod-43.tar" > /srv/artifacts/PROD_LATEST

# Oops: whoever updated the promote config for this release copy-pasted the
# staging artifact's filename instead of the prod one
sed -i 's/^PROD_ARTIFACT=.*/PROD_ARTIFACT=app-staging-43.tar/' /etc/deploy/promote.conf

promote
