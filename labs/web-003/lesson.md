# 502 means nginx is fine, its upstream isn't

A 502 Bad Gateway comes from the proxy itself, not the backend — it
means nginx successfully received the request, tried to forward it to
the address in `proxy_pass`, and got no usable response (connection
refused, reset, or an invalid HTTP response). Compare that to a 504
Gateway Timeout (upstream reachable but too slow) or a connection
error at nginx's own exposed port (nginx itself isn't listening) —
all three point at a different layer.

## Where to look

- `/var/log/nginx/error.log` records the exact address nginx tried and
  the reason it failed — usually `connect() failed (111: Connection
  refused) while connecting to upstream`.
- `ss -ltnp` shows every port actually bound on the host right now, and
  which process owns it — this is ground truth, independent of what any
  config file claims.
- `grep proxy_pass /etc/nginx/sites-available/default` shows what nginx
  believes the backend's address is.

## Fixing it correctly

The fix is always to make nginx's configured address match where the
backend actually is — never to move the backend to match a wrong
config, and never to test success by curling the backend's own port
directly (that proves the backend works, not that the proxy does).
After editing the config: `nginx -t` to catch syntax errors, then
`nginx -s reload` to apply it without dropping the listening socket.
