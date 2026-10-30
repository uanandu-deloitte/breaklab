# Solution: Serve a Static Site

Add a location block to `/etc/nginx/sites-available/default`, inside the
existing `server { ... }` block:

```nginx
location /status/ {
    alias /var/www/status-site/;
    index index.html;
}
```

Then validate and start nginx:

```bash
nginx -t
nginx
```

Confirm the fix:

```bash
curl -i http://localhost:80/status/
```

```
HTTP/1.1 200 OK
Content-Type: text/html

<!doctype html>
<html>
<head><title>Service Status</title></head>
<body><h1>All systems operational</h1></body>
</html>
```

`alias` maps the `/status/` route directly onto
`/var/www/status-site/`, rather than appending `/status/` onto a root
directory (which would have looked for a nonexistent
`.../status/status-site/...` path). `nginx -t` confirms the edit is
syntactically valid before anything is (re)started, and because nginx
wasn't running yet in this container, starting it — not just reloading —
is what actually brings the site up.
