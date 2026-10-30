# Name resolution comes before connection

Every time something connects by hostname instead of by raw IP, there's a
step before any packet goes anywhere: turning that name into an address.
If that step fails, the connection never even gets attempted — the
symptom shows up as a resolution error, not a timeout or a refusal.

## Where names come from

On a typical Linux host (and inside a container), a name lookup checks, in
order: an in-memory cache if one exists, then `/etc/hosts` for a literal
match, then whatever's configured in `/etc/resolv.conf` (a real DNS
server). `/etc/hosts` is the simplest and most direct of these — one line
per name, in the form:

```
<ip-address>  <hostname>
```

`getent hosts <name>` asks the system to resolve a name exactly the way an
application would, using that same lookup order — it's the fastest way to
confirm whether a name resolves at all, independent of whatever tool
you're about to point at it.

## Why this is a different failure than "down"

"Connection refused" means a name resolved to an address, and something
answered at that address saying no one's listening. A DNS/hosts failure
never gets that far — there's no address to connect to yet. Conflating the
two wastes time: you'll get nowhere restarting a backend that was never
reachable in the first place because the name in front of it was never
wired up.

## Fixing it in the right layer

The durable fix is always in the resolution layer itself — add the
missing `/etc/hosts` entry (or fix DNS) — not in the code or command that
does the lookup. Hard-coding a raw IP into an application's request
"fixes" today's symptom while quietly removing the indirection that makes
the name useful in the first place: the day that address changes, the
hard-coded version breaks silently while the name-based one would have
kept working.
