# Solution: Debug an API with curl

```bash
curl -i -X POST \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer test-token-fixture' \
  -d '{"action": "sync"}' \
  http://127.0.0.1:8080/
```

Expected response:

```
HTTP/1.0 200 OK
Content-Type: application/json

{"result": "accepted"}
```

Each piece maps to one of the server's four requirements: `-X POST` sets
the method, the two `-H` flags supply the `Content-Type` and
`Authorization` headers exactly as required, and `-d` sends the JSON body
containing `{"action": "sync"}`. Dropping any one of them (wrong method,
missing header, wrong body) gets a `400` or `401` instead — and because
the server only remembers the *last* request it saw, the final call you
make is the one the check actually grades.
