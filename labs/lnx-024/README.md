# Write a Health Script

This host runs a small application: a background worker process, and an
HTTP server exposing a `/healthz` path on port 80. Write a health-check
script other tooling can rely on.

## Tasks

Write `/workspace/health.sh`, executable, that checks:

1. the worker process is actually running (by process, not by assuming)
2. something is listening on TCP port 80
3. `http://localhost:80/healthz` returns a successful HTTP response
4. disk usage under `/var/log/app` stays below a reasonable threshold
5. exits `0` when everything above is healthy, and non-zero with a
   clear one-line message identifying which check failed otherwise
6. output stays concise — one line per outcome, not verbose debug spew

Get a shell: `docker exec -it breaklab-lnx-024 bash`

Manual check: `curl http://localhost:8080/healthz`
