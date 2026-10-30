#!/usr/bin/env bash
# Two normal, lightweight background workers...
(exec -a worker-heartbeat sleep infinity) &
(exec -a worker-metrics-sync sleep infinity) &
# ...and one that looks like a plausible service name but is actually
# a tight busy loop burning a full CPU core.
(exec -a data-compactor bash -c 'while true; do :; done') &
exec sleep infinity
