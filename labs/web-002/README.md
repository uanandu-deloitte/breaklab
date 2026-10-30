# Build a Reverse Proxy

nginx is running on port 80 and a small backend server is already
running on port 8090 — but nginx has no idea the backend exists. Right
now, requests to nginx get nginx's own default page, never the backend.

## Tasks

- confirm the backend is really listening on its own port, independent
  of nginx
- edit nginx's site config so requests to `/` are reverse-proxied to the
  backend instead of being served by nginx itself
- forward a sensible `Host` header so the backend sees the original
  request
- validate the config and reload nginx so the change takes effect

## Constraints

- the backend's own files and port must not change
- the fix goes through nginx's config, not by moving the backend to
  port 80

Get a shell: `docker exec -it breaklab-web-002 bash`

Check it from outside the container: `curl http://localhost:8080/`
