# HTTP 500 Investigation

`/workspace/access.log` is a web access log with columns:

```
<client-ip> <timestamp> <method> <path> <status> <size>
```

Operations needs to know which client IP is responsible for the most HTTP
500 responses.

Write an executable script at `/workspace/top-500-ip.sh` that reads
`/workspace/access.log` and prints exactly one line to stdout:

```
<count> <client-ip>
```

where `<count>` is that IP's number of 500 responses.

## Constraints

- Don't manually count lines — build a pipeline that computes the answer.
- Don't edit `access.log`.
- Your script's answer must stay correct even if the log's lines are
  reordered — build it to compute the answer, not to echo one you worked
  out by eye.

Get a shell: `docker exec -it breaklab-lnx-010 bash`
