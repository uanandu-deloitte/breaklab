# Find the Failure in the Logs

Customers reported a single failed checkout this morning. Three services
logged something at the time it happened: `/var/log/frontend.log`,
`/var/log/api.log`, and `/var/log/payments.log`. Each of them is noisy —
most of their traffic that morning was unrelated and succeeded fine.

## Tasks

- correlate timestamps and request IDs across all three logs to find the
  one request that actually failed end to end
- identify which downstream dependency's own failure was the true root
  cause (not just which service happened to return the error to the
  customer)
- write your findings to `/workspace/answers.env`:

```
FAILING_DEPENDENCY=<the service whose own call failed>
REQUEST_ID=<the request ID that ties all three logs together>
```

Get a shell: `docker exec -it breaklab-lnx-023 bash`
