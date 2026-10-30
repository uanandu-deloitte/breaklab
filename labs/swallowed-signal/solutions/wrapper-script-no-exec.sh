#!/usr/bin/env bash
set -euo pipefail
# Reference fix: the wrapper needs exec too, not just launch-cmd.
cat > /usr/local/bin/run-app.sh <<'EOF'
#!/usr/bin/env bash
echo "$(date -u +%FT%TZ) launching app" >> /var/log/app/startup.log
exec python3 /usr/local/bin/app.py
EOF
chmod 755 /usr/local/bin/run-app.sh
kill -HUP 1
sleep 0.5
