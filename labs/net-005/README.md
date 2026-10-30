# Broken Name Resolution

A backend API is already running inside this container, listening on
`127.0.0.1:8080`. A frontend request is supposed to reach it as
`http://api.internal:8080/`, but that call fails before it even gets to
the network layer.

## Tasks

- determine exactly why `curl http://api.internal:8080/` fails
- fix name resolution for `api.internal` so it resolves to the right
  address
- do **not** fix this by hard-coding an IP address into the curl call (or
  anywhere else) — the fix belongs in the resolver configuration
- confirm `curl http://api.internal:8080/` succeeds and returns the
  backend's JSON marker

## Useful commands

```
getent hosts api.internal
cat /etc/hosts
curl http://api.internal:8080/
curl http://127.0.0.1:8080/        # the backend itself, to confirm it's healthy
```

Get a shell: `docker exec -it breaklab-net-005 bash`
