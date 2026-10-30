#!/usr/bin/env bash
set -e
# Two independent local services on two different ports: nginx on 80
# (its package default) and a plain Python HTTP server on 8081.
nginx
nohup python3 -m http.server 8081 --directory /workspace >/tmp/http8081.log 2>&1 &
disown
exec sleep infinity
