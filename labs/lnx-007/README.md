# Inspect Configuration

Two versions of the same application config shipped:
`/workspace/config.old` and `/workspace/config.new`. Exactly 3 keys
changed value between them; the rest are identical.

## Tasks

- Compare the two files and identify the 3 keys whose value changed.
- Record each one's old and new value in `/workspace/answers.env`:

```
DATABASE_HOST_OLD=
DATABASE_HOST_NEW=
LOG_LEVEL_OLD=
LOG_LEVEL_NEW=
TIMEOUT_OLD=
TIMEOUT_NEW=
```

Get a shell: `docker exec -it breaklab-lnx-007 bash`
