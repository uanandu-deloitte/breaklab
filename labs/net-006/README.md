# Debug an API with curl

A fixture API is running inside this container on `127.0.0.1:8080`. It
only accepts a request that is, all at once:

- method `POST`
- header `Content-Type: application/json`
- header `Authorization: Bearer test-token-fixture`
- a JSON body containing `{"action": "sync"}`

Anything else gets a `400` or `401`, and the server records the outcome
of your most recent request.

## Tasks

- use `curl` to send a request that satisfies all four requirements at
  once
- confirm the server's response is a success (`200`)

## Useful commands

```
curl -i -X POST http://127.0.0.1:8080/
```

Start by sending something incomplete and reading exactly what comes
back — the response body names which requirement failed.

Get a shell: `docker exec -it breaklab-net-006 bash`
