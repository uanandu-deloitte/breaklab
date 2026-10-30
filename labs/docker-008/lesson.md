# A reverse proxy needs three things true at once

Fronting a container with a reverse proxy (nginx, in this lab) only works
end-to-end when three independent facts all hold:

1. **The config is syntactically valid.** A malformed config can mean
   nginx refuses to start, or refuses to reload an edit while an older
   instance keeps running on stale settings.
2. **The proxy is actually listening on the public port.** The process
   has to be up and bound to the address clients will actually connect
   to.
3. **The upstream target is correct.** The `proxy_pass` address has to
   point at a host:port where the real backend is actually listening —
   not where it used to listen, not where a renumbering pass assumed it
   would be.

Any one of these can break independently while the other two stay
perfectly fine. "nginx is running" only confirms #2. "The config exists
and nginx -t passes" only confirms #1. Neither tells you anything about
#3.

## Validate syntax before anything else

`nginx -t` checks that the config file is well-formed and every directive
it uses is one nginx recognizes. It is fast, free, and catches an entire
category of mistakes (a missing closing brace, a typo'd directive) before
you waste time chasing a symptom that was really just a syntax error.

What `nginx -t` **cannot** tell you: whether the address a `proxy_pass`
line names is one anything is actually listening on. That's a runtime
connectivity fact, not a syntax one — a config can be perfectly
well-formed and still point nowhere useful.

## Find out what's actually listening, don't trust the config

`ss -ltnp` lists every socket actually in a `LISTEN` state right now, the
port it's bound to, and the process holding it. This is ground truth,
independent of what any config file *says* should be happening. Comparing
this output against the `proxy_pass` address in the nginx config is the
fastest way to spot a mismatch: the proxy configured to forward
to one port, while the real backend sits on another.

## Applying a fix

Editing the config file on disk changes nothing on its own — nginx reads
it once at startup and otherwise keeps running on the version it already
loaded. `nginx -s reload` tells the running master process to re-read the
config and swap in the new settings without dropping the listening
socket, which is the normal way to apply a fix without restarting the
whole container.
