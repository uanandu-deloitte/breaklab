# Fronting a Container With a Reverse Proxy

This container runs nginx in front of a small backend server. Host port
`8080` is mapped to the container's port `80`, where nginx listens. Right
now, `curl http://localhost:8080/` returns a `502 Bad Gateway`.

## Tasks

1. Confirm the symptom: `curl -I http://localhost:8080/` from your host
   shell.
2. Get a shell (`docker exec -it breaklab-docker-008 bash`) and run
   `nginx -t` — rule out a syntax error before suspecting anything else.
3. Run `ss -ltnp` to see what's actually listening, and on which port the
   backend really is.
4. Compare that against `grep proxy_pass /etc/nginx/sites-available/default`
   — find the mismatch between what nginx is configured to forward to and
   where the backend actually listens.
5. Fix the config so `proxy_pass` points at the backend's real port, then
   reload nginx (`nginx -s reload`) to apply it.
6. Confirm the fix: `curl http://localhost:8080/` should now return the
   backend's real content.

Record your findings in `/workspace/answers.env`, one `KEY=value` per
line:

```
BACKEND_REAL_PORT=
ORIGINAL_PROXY_PASS_PORT=
PROXY_FIXED=
```

- `BACKEND_REAL_PORT`: the port the backend process is actually listening
  on, as shown by `ss -ltnp` — not a guess, not what the config used to
  say.
- `ORIGINAL_PROXY_PASS_PORT`: the wrong port nginx's `proxy_pass` pointed
  at before you fixed it.
- `PROXY_FIXED`: `yes` once the mapped host port reaches the backend
  through the proxy.

Get a shell: `docker exec -it breaklab-docker-008 bash`
