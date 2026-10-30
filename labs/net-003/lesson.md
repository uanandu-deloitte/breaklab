# Four ways a connection can fail, and how to tell them apart

"It doesn't work" covers at least four completely different failures,
each pointing at a different layer and a different fix. Telling them
apart — instead of treating every failed `curl` the same — is most of
what network troubleshooting actually is.

## Healthy: everything worked

A 200 (or any valid HTTP response) means DNS resolved, the TCP handshake
completed, the server accepted the connection, and it sent a response
back in reasonable time. Nothing to diagnose — but it's worth knowing
what "fine" looks like before you go looking for what isn't.

## DNS failure: the name never became an address

`curl`'s very first job is turning a hostname into an IP address before
it can even attempt a connection. If that lookup fails, curl reports
something like `Could not resolve host` — and note that this happens
*before* any TCP connection is even attempted; the destination's port
state is irrelevant because nothing ever reached it. `cat
/etc/resolv.conf` tells you where this host is even asking; a name that
isn't registered anywhere (or a resolver that's unreachable) both show
up this way.

## Connection refused: a reachable host, but nobody home

DNS resolved (or you used a bare IP), the host is reachable, and the
kernel got an actual answer back fast: "nothing is listening on that
port." This produces an immediate, explicit `Connection refused` — not a
hang. `ss -ltn` run against the target itself would show no process has
that port in `LISTEN` state at all. `nc -zv <host> <port>` reproduces
exactly this signal without needing an application protocol: it just
tries to open the TCP connection and reports whether it succeeded.

## Timeout: something answered the connection, then said nothing

This is the trickiest of the four to recognize, because there's no error
message at all — just silence, until your client gives up waiting. A
process accepted the TCP connection (so it isn't "refused"), but then
never sent any application-layer response. This is why giving `curl` a
short `--max-time` matters when probing somewhere unfamiliar: without
one, you can wait indefinitely for a reply that was never coming.
`ss -ltn` against this target *would* show something listening — the
difference from "healthy" only shows up once you actually wait for a
response and never get one.

## Reading the symptom instead of guessing

The shape of the failure — an immediate resolver error, an immediate
refusal, or silence until a timeout — tells you which layer broke before
you've touched the target at all. Guessing based on "the request failed"
alone throws that information away.
