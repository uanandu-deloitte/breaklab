#!/usr/bin/env bash
set -euo pipefail
# Oops: an old internal tool from a previous migration got symlinked to
# the name "jq" for some other team's convenience, and it was never
# cleaned up. It takes similar-looking flags but doesn't behave the same.
mkdir -p /opt/legacy-tools
cat > /opt/legacy-tools/json-query <<'EOF'
#!/usr/bin/env bash
# long-unmaintained internal tool, kept around for a different script;
# mishandles the -r flag and the filter argument entirely
echo "legacy-unset"
EOF
chmod 755 /opt/legacy-tools/json-query
ln -sf /opt/legacy-tools/json-query /usr/local/bin/jq
