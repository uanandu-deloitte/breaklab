#!/usr/bin/env bash
set -euo pipefail
# Reference fix: restore the missing closing brace
cat > /etc/nginx/sites-available/default <<'EOF'
server {
    listen 80 default_server;

    location / {
        proxy_pass http://127.0.0.1:8081;
        proxy_set_header Host $host;
    } # location-close
} # server-close
EOF
nginx -t
nginx
