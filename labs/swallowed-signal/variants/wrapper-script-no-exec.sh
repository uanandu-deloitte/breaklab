#!/usr/bin/env bash
set -euo pipefail
# Oops: a small wrapper was added to log each startup. launch-cmd itself
# still correctly uses exec -- but the wrapper it hands off to doesn't.
printf 'exec /usr/local/bin/run-app.sh\n' > /etc/app/launch-cmd
cat > /usr/local/bin/run-app.sh <<'EOF'
#!/usr/bin/env bash
echo "$(date -u +%FT%TZ) launching app" >> /var/log/app/startup.log
python3 /usr/local/bin/app.py
EOF
chmod 755 /usr/local/bin/run-app.sh
kill -HUP 1
sleep 0.5
