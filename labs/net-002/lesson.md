# A port can only have one owner

A TCP port is a scarce resource on a given IP address: only one process
can hold a listening socket on `0.0.0.0:80` at a time. Whichever process
calls `bind()` on that address/port pair *first* wins it — every process
that tries afterward gets `EADDRINUSE` ("Address already in use") and
fails to start, no matter how correctly it's configured otherwise.

## Why "the web server won't start" doesn't mean nginx is broken

When nginx fails to bind its configured port, its own config can be
perfect — the problem is almost always that something else got there
first. This is a different failure mode than a bad config (which `nginx
-t` would catch) or a missing file (which a 403/404 would show): the
server process itself never even reaches the point of serving a
response, because it couldn't start.

## Finding the real owner of a port

`ss -ltnp` lists every TCP socket in `LISTEN` state along with the
process that holds it (name and PID, in the `users:(("name",pid=...))`
field). This is the fastest way to answer "what, specifically, is
sitting on this port right now" — don't guess from the symptom alone.
`ps aux` (or `ps -p <pid> -o cmd=`) then shows that process's full
command line, which usually makes it obvious whether it's the service
you expect or something else entirely (a leftover debugging process, a
second instance, a completely unrelated tool).

## The fix is precise, not broad

Once you know the exact PID holding the port, stop *that* process
specifically (`kill <pid>`, or `pkill -f` against a command-line pattern
specific enough to match only it) — not every process that might be
vaguely related. Killing the wrong thing can take down a process that
was actually fine. Only after the port is free does starting the real
service make sense; starting it while the decoy still holds the port
just produces the same bind failure again.
