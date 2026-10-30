# Read the Access Log

`/workspace/access.log` is a web access log with columns:

```
<client-ip> <timestamp> <method> <path> <status> <size>
```

Using `head`, `tail`, `wc`, `cut`, `sort`, and `uniq`, determine:

- the total number of requests
- the timestamp of the first and last request
- how many responses were 2xx, 4xx, and 5xx
- the single most frequent request path

Write all six answers to `/workspace/answers.env`, one `KEY=value` per line:

```
TOTAL_REQUESTS=
FIRST_TIMESTAMP=
LAST_TIMESTAMP=
COUNT_2XX=
COUNT_4XX=
COUNT_5XX=
MOST_FREQUENT_PATH=
```

## Constraints

- Don't edit `access.log` — read it, don't change it.
- Derive every value with a command against the live file; don't eyeball
  the log and type in a guess.

Get a shell: `docker exec -it breaklab-lnx-009 bash`
