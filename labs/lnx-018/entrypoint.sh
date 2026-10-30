#!/usr/bin/env bash
# worker1: the instance the learner will terminate gracefully.
(exec -a worker1 bash /usr/local/bin/worker.sh /workspace/worker-graceful.marker) &
exec sleep infinity
