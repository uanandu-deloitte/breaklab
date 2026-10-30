# HTTP Failure Classification

Five endpoints are reachable from inside this container, each failing (or
succeeding) in a different way:

| Label      | Address                              |
|------------|---------------------------------------|
| ENDPOINT_1 | `http://nosuchhost.invalid:9001/`     |
| ENDPOINT_2 | `http://127.0.0.1:9002/`              |
| ENDPOINT_3 | `http://127.0.0.1:9003/`              |
| ENDPOINT_4 | `http://127.0.0.1:9004/anything`      |
| ENDPOINT_5 | `http://127.0.0.1:9005/`              |

## Tasks

Probe each endpoint with `curl` and classify it. Write your answers to
`/workspace/answers.env`, one `KEY=value` per line, using exactly one of
these values for each: `dns-failure`, `refused`, `timeout`, `404`, `500`.

```
ENDPOINT_1=
ENDPOINT_2=
ENDPOINT_3=
ENDPOINT_4=
ENDPOINT_5=
```

## Useful commands

```
curl -v --max-time 3 <url>
```

Use `--max-time` on every probe — one of these endpoints accepts the
connection and then never responds, and without a timeout you'll just be
waiting on it.

Get a shell: `docker exec -it breaklab-net-007 bash`
