# Identify the Network

You have been given shell access to an unfamiliar Linux host that talks to
other services over the network. Before debugging any connectivity
problem, establish this host's own network identity.

Record these values in `/workspace/answers.env`, one `KEY=value` per line:

```
INTERFACE_NAME=
CONTAINER_IP=
LISTENING_PORTS=
```

- `INTERFACE_NAME` — the name of the network interface that carries this
  host's default route (e.g. `eth0`).
- `CONTAINER_IP` — this container's own IPv4 address on that interface.
- `LISTENING_PORTS` — every TCP port this host is listening on, as a
  comma-separated, numerically sorted list with no spaces (e.g. `80,8081`).

Also inspect (no answer required, but you'll need it to understand what
you're looking at):

- the default route: `ip route`
- the DNS resolver configuration: `cat /etc/resolv.conf`

Collect every value from the running environment itself — don't guess or
hard-code anything.

Get a shell: `docker exec -it breaklab-net-001 bash`
