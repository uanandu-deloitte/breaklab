#!/usr/bin/env bash
set -e
# Endpoint 1: healthy — nginx listens on 8001.
nginx

# Endpoint 2 (DNS failure) and endpoint 3 (connection refused) need no
# process at all: nonexistent.internal simply never resolves, and
# nothing ever binds port 8003.

# Endpoint 4: accepts the TCP connection but never replies, so a client
# hangs until it gives up — a timeout, not an error.
nohup nc -lk -p 8004 </dev/null >/tmp/nc8004.log 2>&1 &
disown

exec sleep infinity
