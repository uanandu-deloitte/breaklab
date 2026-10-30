# Solution: Fronting a Container With a Reverse Proxy

```bash
nginx -t   # confirms the config is syntactically valid already

ss -ltnp   # shows python3 listening on 5001, nginx on 80 -- nothing on 5002

sed -i 's/5002/5001/' /etc/nginx/sites-available/default
nginx -t
nginx -s reload

curl http://localhost:8080/   # now returns the backend's real content

cat <<EOF > /workspace/answers.env
BACKEND_REAL_PORT=5001
ORIGINAL_PROXY_PASS_PORT=5002
PROXY_FIXED=yes
EOF
```

`nginx -t` passing from the start proves this was never a syntax problem.
`ss -ltnp` shows the backend genuinely listening on `5001`, while the
config's `proxy_pass` pointed at `5002` — a port nothing was bound to.
Editing the file alone does nothing until `nginx -s reload` tells the
running master process to pick up the change. Once the upstream address
matches the backend's real port, the mapped host port `8080` reaches it
end-to-end.
