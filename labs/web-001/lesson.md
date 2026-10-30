# Serving a specific file at a specific route

nginx doesn't know about a file just because it exists on disk — every
route it serves has to be named in its configuration, mapping a URL path
to a location on the filesystem. A file with no matching `location` block
is invisible to nginx no matter how correct its content is.

## root vs. alias

A `location` block's `root` directive appends the request path onto a
base directory: `location /status/ { root /var/www; }` would look for
`/var/www/status/index.html`. `alias` instead *replaces* the location's
path prefix with the given directory: `location /status/ { alias
/var/www/status-site/; }` maps `/status/` straight onto
`/var/www/status-site/`, regardless of the directory name. When the
route's prefix doesn't match the directory's own name, `alias` is
usually what you want.

## Validate before you reload

`nginx -t` checks a configuration file's syntax without starting or
restarting anything. It catches a surprising amount — missing semicolons,
unmatched braces, directives in the wrong context — before they become a
downtime instead of a one-line diff. Making this check a reflex, every
time, before touching the running process, costs nothing and prevents a
real outage.

## Validating isn't starting

A config that passes `nginx -t` still hasn't done anything by itself —
nginx has to actually be running, and either started fresh or told to
reload, before any change takes effect. On a host with no service
manager, that's as direct as running `nginx` (to start it for the first
time) or `nginx -s reload` (to pick up a config change on an
already-running instance).
