# Graceful Shutdown

A worker process called `worker1` is already running in this container.
It's built to write a completion marker file when it shuts down cleanly —
but only if it gets the chance.

## Tasks

1. Find `worker1`'s PID.
2. Send it a **graceful** terminate signal: `kill <pid>` (this sends
   `SIGTERM`, not `SIGKILL`). Give it a second, then confirm
   `/workspace/worker-graceful.marker` now exists.
3. Start a second instance of the same worker yourself, using a
   **different** marker path so it doesn't collide with worker1's:

   ```
   /usr/local/bin/worker.sh /workspace/worker-forced.marker &
   ```

4. Find its PID and this time force it: `kill -9 <pid>` (`SIGKILL`).
   Confirm `/workspace/worker-forced.marker` was **not** created.
5. Record what you observed in `/workspace/answers.env`:

   ```
   GRACEFUL_MARKER_EXISTS=yes
   FORCED_MARKER_EXISTS=no
   ```

Get a shell: `docker exec -it breaklab-lnx-018 bash`
