# Solution: Build a Reverse Proxy

```bash
# Confirm the backend is really up on its own port.
curl http://127.0.0.1:8090/

# Edit nginx's site config.
cat > /etc/nginx/sites-available/default <<'EOF'
server {
    listen 80 default_server;
    server_name _;

    location / {
        proxy_pass http://127.0.0.1:8090;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
EOF

# Validate, then reload -- never skip -t before reload.
nginx -t
nginx -s reload

curl http://localhost:80/
```

`proxy_pass` is the one directive that turns nginx from "serves its own
files" into "forwards to the backend" — everything else here is good
hygiene around it: `proxy_set_header Host $host` keeps the backend's
view of the request's intended host accurate, and `nginx -t` catches a
syntax mistake before it's live.
