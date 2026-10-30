# Serve a Static Site

A one-page status site already exists on disk at
`/var/www/status-site/index.html`. It needs to be reachable at the route
`/status/` on port 80 — but nginx isn't configured to serve it from there
yet, and nginx itself isn't running.

## Tasks

- add a `location` block to `/etc/nginx/sites-available/default` that
  serves `/var/www/status-site/index.html` at the route `/status/`
- validate the configuration with `nginx -t`
- start nginx
- confirm `curl http://localhost:8080/status/` returns the status page

Get a shell: `docker exec -it breaklab-web-001 bash`
