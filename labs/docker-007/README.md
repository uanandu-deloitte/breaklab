# Where docker logs Actually Comes From

This container runs a small app that writes a heartbeat line once a
second — to its own stdout, and to a separate file on disk, at the same
time. Figure out which of those two `docker logs` actually captures.

## Tasks

1. From your host shell, run `docker logs breaklab-docker-007 --tail 20`
   and note what it shows.
2. Get a shell inside the container
   (`docker exec -it breaklab-docker-007 bash`) and find the on-disk log
   file the same app is writing to. It is real, current, and growing —
   confirm that for yourself before answering.
3. Decide, from what you actually observed in steps 1 and 2 — not from
   assumption — whether `docker logs` ever shows the file-only content.

Record your findings in `/workspace/answers.env`, one `KEY=value` per
line:

```
STDOUT_CAPTURED=
FILE_LOG_CAPTURED=
LOG_FILE_PATH=
```

- `STDOUT_CAPTURED`: `yes` or `no` — does `docker logs` show the
  heartbeat lines this app writes to its own stdout?
- `FILE_LOG_CAPTURED`: `yes` or `no` — does `docker logs` ever show the
  separate file this same app writes to disk?
- `LOG_FILE_PATH`: the absolute path of the on-disk log file you found.

Get a shell: `docker exec -it breaklab-docker-007 bash`
