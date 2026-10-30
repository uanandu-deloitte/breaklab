#!/usr/bin/env bash
set -euo pipefail
# New release cycle: both staging and prod get rebuilt and retagged as build 43
mkdir -p /tmp/build-prod-43 /tmp/build-staging-43
echo "prod build 43: checkout-v3.1 (security patch)" > /tmp/build-prod-43/CONTENT
echo "staging build 43: checkout-v5-experimental (unreleased)" > /tmp/build-staging-43/CONTENT
tar -cf /srv/artifacts/app-prod-43.tar -C /tmp/build-prod-43 CONTENT
tar -cf /srv/artifacts/app-staging-43.tar -C /tmp/build-staging-43 CONTENT
echo "app-prod-43.tar" > /srv/artifacts/PROD_LATEST

# Oops: someone replaced the explicit, config-driven promotion with a "just
# deploy whatever was built most recently" script using a glob and a sort
cat > /usr/local/bin/promote <<'EOF'
#!/usr/bin/env bash
set -euo pipefail
latest=$(ls /srv/artifacts/app-*.tar | sort | tail -n1)
rm -rf /srv/app/live
mkdir -p /srv/app/live
tar -xf "$latest" -C /srv/app/live
echo "promote: deployed $(basename "$latest") to /srv/app/live"
EOF
chmod 755 /usr/local/bin/promote

promote
