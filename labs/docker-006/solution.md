# Solution: Nothing Listens Where You Think

```bash
# Confirm the bind address from inside the container:
ss -tlnp | grep :80
# LISTEN 0 ... 127.0.0.1:80 ...   <- bound to loopback only

sed -i 's/--bind 127.0.0.1/--bind 0.0.0.0/' /usr/local/bin/serve.sh

pkill -f 'http.server 80' || true
nohup /usr/local/bin/serve.sh >/var/log/serve.log 2>&1 &
disown
```

(`entrypoint.sh` is what originally launched `serve.sh` in the background
when the container started — killing and relaunching it directly, as
above, is enough; the container's own foreground process is unrelated
`sleep infinity`, so it keeps running throughout.)

`ss -tlnp` shows the server was listening on `127.0.0.1:80` — bound only
to the loopback interface inside the container's own network namespace.
`curl localhost` from inside the container succeeds because it also
originates from that same loopback interface. Docker's host-to-container
port mapping, however, delivers traffic to the container's regular
network interface, which a server bound to `127.0.0.1` never accepts
connections on — so the host-side `curl` gets nothing, even though the
process is healthy and the port is genuinely listening.

Editing `/usr/local/bin/serve.sh` to bind `0.0.0.0` (all interfaces)
instead, then killing the old process and starting a fresh one with that
script, makes the server accept connections from any interface —
including the one Docker's port mapping actually uses.
