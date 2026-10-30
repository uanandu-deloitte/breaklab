# A host's network identity, in four questions

Before you can debug *why* two machines can't talk to each other, you need
to establish the basic network facts about the machine you're sitting on:
which interface it talks through, what address that interface holds,
where it sends traffic that isn't local, and which of its own services
are actually reachable.

## Interfaces and addresses: `ip addr`

Every network-capable process on a Linux host sends and receives through
a network **interface** — a named object (`eth0`, `lo`, `wlan0`, ...)
that the kernel attaches one or more IP addresses to. `ip addr` (or `ip
addr show`) lists every interface this host has, each with:

- its name and state (`UP`/`DOWN`)
- its MAC address
- its IPv4/IPv6 addresses, each written as `address/prefix-length`
  (e.g. `172.17.0.2/16` — the `/16` is the subnet mask in CIDR notation)

`lo`, the loopback interface (`127.0.0.1`), always exists and only ever
talks to the host itself — it's never the interface carrying traffic to
anywhere else.

## Where traffic goes when it isn't local: `ip route`

A **route** tells the kernel which interface (and which next-hop gateway)
to use for traffic to a given destination network. `ip route` lists every
route this host knows; the line starting `default via <gateway> dev
<interface>` is the **default route** — where everything goes that isn't
covered by a more specific route. The `dev` field on that line is the
interface that actually carries traffic leaving this host.

## Turning names into addresses: `/etc/resolv.conf`

IP routing only understands addresses, not hostnames — something has to
translate `example.com` into an address before a connection can be routed
anywhere. `/etc/resolv.conf` names the DNS resolver(s) this host asks to
do that translation, as one or more `nameserver <address>` lines. If this
file points at a resolver that's unreachable or misconfigured, every
hostname lookup fails even though routing itself is completely healthy —
a classic way "the network is down" turns out to mean "DNS is down."

## What's actually listening: `ss -tlnp`

A process has to explicitly **bind** to a port and put it into `LISTEN`
state before anything can connect to it. `ss -tlnp` lists every TCP
socket in that state:

- `-t` — TCP sockets only
- `-l` — listening sockets only (not established connections)
- `-n` — show numeric ports/addresses instead of resolving service names
- `-p` — show the process (name and PID) that owns each socket

The `Local Address:Port` column is what matters most: `0.0.0.0:80` means
a socket is listening on port 80 on every interface, while `127.0.0.1:80`
would mean it's only reachable from the host itself (not from outside,
even across a real network).

## Why this comes before any break/fix work

Every networking incident downstream of this lab boils down to one of
these four facts being wrong: the wrong interface, an unreachable
gateway, a dead DNS resolver, or a process not actually listening where
you expect. Knowing how to read all four quickly, on any unfamiliar host,
is the foundation every other networking lab in this stage builds on.
