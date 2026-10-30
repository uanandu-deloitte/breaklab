#!/usr/bin/env bash
set -euo pipefail
# Oops: a stub someone wrote to test deploy-sync offline (so it wouldn't
# need real release data) never got removed, and /usr/local/bin comes
# before /usr/bin on $PATH.
cat > /usr/local/bin/jq <<'EOF'
#!/usr/bin/env bash
# leftover test stub — always "works", never reads its input
echo "0.0.0-stub"
EOF
chmod 755 /usr/local/bin/jq
