# Application Cannot Read Configuration

The application runs as `appuser` and reads `/srv/app/config.yaml` on
every health check. Right now `/healthz` is failing.

## Tasks

- inspect the application process's identity
- inspect the ownership and permissions on `/srv/app/config.yaml` (and the
  directories above it)
- repair access using least privilege
- confirm the health endpoint succeeds

## Constraints

- no `chmod 777`
- no world-writable configuration
- the application must keep running as a non-root user
- the configuration's content must not change

Get a shell: `docker exec -it breaklab-lnx-015 bash`

Health check: `curl http://localhost:8080/healthz`
