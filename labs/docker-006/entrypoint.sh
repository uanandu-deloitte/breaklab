#!/usr/bin/env bash
set -e
# Runs the status-page server in the background so the container keeps
# running even after the learner kills and restarts it while fixing the
# bind address.
nohup /usr/local/bin/serve.sh >/var/log/serve.log 2>&1 &
disown
exec sleep infinity
