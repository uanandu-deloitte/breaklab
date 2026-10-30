# Why Clocks Have to Agree Across Hosts

Two simulated hosts, `host-a` and `host-b`, each log events with their
own timestamp. One of them is telling the truth about what time it is.
The other isn't. Find out which, and by how much.

## Tasks

1. Get the real current time directly: `date +%s`.
2. Produce a fresh timestamp from each host: `log-event-a "check"` and
   `log-event-b "check"`, then read the last line each wrote to
   `/var/log/host-a.log` / `/var/log/host-b.log` (the second
   whitespace-separated field is that host's unix timestamp).
3. Compare each host's timestamp against the real current time you
   captured in step 1 to compute that host's skew.
4. Decide which host is drifted and by roughly how many seconds, and
   confirm this is clock drift, not a real event-ordering problem.

Record your findings in `/workspace/answers.env`, one `KEY=value` per line:

```
DRIFTED_HOST=
SKEW_SECONDS=
IS_REAL_ORDERING_ISSUE=
```

- `DRIFTED_HOST`: `host-a` or `host-b` — whichever one's clock disagrees
  with real time.
- `SKEW_SECONDS`: roughly how many seconds that host's clock is off by
  (a positive number — direction doesn't matter here, just magnitude).
- `IS_REAL_ORDERING_ISSUE`: `yes` or `no` — is this actually a genuine
  causality/ordering problem, or just one clock disagreeing with
  reality?

Get a shell: `docker exec -it breaklab-obs-007 bash`
