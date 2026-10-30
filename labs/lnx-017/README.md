# Find the Runaway Process

This container is running three background processes. Two of them are
healthy, idle workers. One of them is quietly pinning a full CPU core and
needs to be stopped.

## Tasks

1. Use `ps` or `top` to find the offending process. Don't guess from the
   name alone — confirm it by its actual CPU usage.
2. Record its identity into `/workspace/answers.env`, one `KEY=value` per
   line:

   ```
   OFFENDING_PID=
   OFFENDING_PPID=
   OFFENDING_USER=
   OFFENDING_CMD=
   ```

3. Terminate it **gracefully** — send it `SIGTERM` (plain `kill <pid>`),
   not `SIGKILL` (`kill -9`).
4. Confirm the two healthy background processes are still running.

Get a shell: `docker exec -it breaklab-lnx-017 bash`
