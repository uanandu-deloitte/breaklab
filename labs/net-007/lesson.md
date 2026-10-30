# Five ways an HTTP request can fail

"It doesn't work" collapses several genuinely different failures into one
sentence. Each one happens at a different stage of making a request, and
each one calls for a different response.

## Before any connection: DNS failure

If a hostname doesn't resolve, nothing is ever attempted on the network —
`curl` fails immediately with something like `Could not resolve host`.
This is purely a naming problem; the server (if it exists at all) is
never involved.

## Connection refused

The name resolved, but when curl tried to open a TCP connection to that
address and port, the operating system on the other end answered
immediately with a reset: nothing is listening there. This fails fast,
just like DNS failure — the two are easy to confuse from the outside, but
`curl -v` labels them differently (`Could not resolve host` vs
`Connection refused`).

## Timeout

The connection succeeded — something accepted it — but then nothing ever
came back. Curl has no way to know whether a response is still coming or
never will be, so it just waits, by default for a very long time. A short
`--max-time` turns an indefinite hang into a fast, deliberate failure you
can actually use.

## 404 — the request reached a real server, which said "not here"

A connection succeeded, a real HTTP response came back, and that response
says the specific thing being asked for doesn't exist at that path. This
is the server working correctly and telling you about a problem with the
request, not with itself.

## 500 — the request reached a real server, which failed to handle it

Also a connection succeeded and a real HTTP response came back — but this
time the status code says the server itself hit a problem trying to
produce an answer. Same successful connection as 404, opposite meaning:
the request was probably fine, the server wasn't.

## Why the distinction matters

DNS and refusal point you at configuration (a hostname, a port, whether
a process is even running). A timeout points you at a process that's
alive but stuck. A 404 points you at the request itself. A 500 points you
at the server's own logic. Treating all four as "the website is down"
sends you investigating the wrong layer at least three times out of four.
