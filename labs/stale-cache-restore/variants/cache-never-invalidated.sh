#!/usr/bin/env bash
set -euo pipefail
# Oops: a refactor stripped the cache-key check out entirely and replaced
# ci-build with a version that always restores whatever is cached
cat > /usr/local/bin/ci-build <<'EOF'
#!/usr/bin/env bash
set -euo pipefail
CACHE_DIR=/var/cache/ci/deps
VENDOR_DIR=/srv/app/vendor

rm -rf "$VENDOR_DIR"
mkdir -p "$VENDOR_DIR"
cp "$CACHE_DIR/requirements.lock" "$VENDOR_DIR/requirements.lock"
echo "ci-build: build complete"
EOF
chmod 755 /usr/local/bin/ci-build

# Real dependency bump: libwidget needs 3.0, but the cache dir below still
# holds whatever was resolved for the last build (2.0)
sed -i 's/^libwidget==.*/libwidget==3.0/' /srv/app/requirements.lock

ci-build
