# A listening port isn't automatically a reachable one

"The process is running" and "the port is reachable from the host" are
two different facts, and a healthy-looking container can satisfy the
first while failing the second. The missing piece is almost always which
*address* a listening socket is bound to, not whether it's listening at
all.

## One port, several possible addresses

A process that opens a listening socket chooses which address (or
addresses) to accept connections on:

- `0.0.0.0` (all interfaces) accepts a connection arriving on any network
  interface the host — or container — has.
- `127.0.0.1` (loopback only) accepts a connection only if it originates
  from inside the same network namespace as the listening process itself.

Both show up identically in `ps aux` (the process is running) and both
show up in `ss -tlnp` as a `LISTEN` line on the right port — the
difference is only visible in the *address* column of that same line.

## Why curl from inside the container lies to you

`curl localhost` run inside the container originates from that
container's own loopback interface — so it succeeds against a server
bound to `127.0.0.1` just as easily as one bound to `0.0.0.0`. This makes
"it works when I exec in and curl it" misleading: it proves the process
is alive and the port is open, but it proves nothing about whether
traffic arriving from *outside* that network namespace — which is exactly
what Docker's host-to-container port mapping delivers — can reach it.

## Narrowing it down

Running the same `curl` from two different vantage points localizes the
problem immediately:

- `curl` from the **host**, against the mapped port, exercises the real
  path a real client would use.
- `curl localhost` from **inside the container** exercises only the
  loopback path.

A pass from inside and a failure from the host, with the process
confirmed alive and the port confirmed listening, points at exactly one
thing: the bind address. The fix is to make the server listen on
`0.0.0.0` so it accepts connections arriving on the container's real
network interface, not only its loopback one.
