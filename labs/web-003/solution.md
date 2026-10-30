# Solution: 502 Bad Gateway

```bash
# See the exact address nginx is failing to reach.
tail -n 5 /var/log/nginx/error.log

# Find out where the backend is really listening.
ss -ltnp | grep python3

# Compare that against nginx's configured upstream.
grep proxy_pass /etc/nginx/sites-available/default

# Fix the mismatch.
sed -i 's/127.0.0.1:8098/127.0.0.1:8095/' /etc/nginx/sites-available/default

# Validate, then reload through nginx -- never skip -t.
nginx -t
nginx -s reload

# Confirm the fix through nginx's own exposed port, not the backend directly.
curl http://localhost:80/
```

The error log already names the broken address before you touch
anything; `ss -ltnp` gives the backend's real port independently of any
config file's claims. The fix is a one-line correction to
`proxy_pass`, reloaded through nginx — the backend itself was never the
problem.
