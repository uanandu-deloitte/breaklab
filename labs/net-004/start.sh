#!/usr/bin/env bash
set -e
nohup python3 -m http.server 8005 --directory /workspace >/tmp/http8005.log 2>&1 &
disown
exec sleep infinity
