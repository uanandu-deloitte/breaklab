# PID 1 and the Command

This container's Dockerfile has a `CMD` that looks like it just runs
`sleep infinity`. Find out what's actually running as PID 1 — don't trust
the Dockerfile's appearance, check the live process.

## Tasks

1. Inspect PID 1's real command line with
   `cat /proc/1/cmdline | tr '\0' ' '`.
2. Cross-check it with `ps -p 1 -o comm=`.
3. Decide whether PID 1 is the real program (`sleep`) itself, or a shell
   wrapping it.

Record your findings in `/workspace/answers.env`, one `KEY=value` per
line:

```
PID1_COMMAND=
IS_SHELL_FORM=
```

- `PID1_COMMAND`: PID 1's actual command line, exactly as
  `/proc/1/cmdline` shows it (with its null separators turned into
  single spaces).
- `IS_SHELL_FORM`: `yes` or `no` — is PID 1 a shell wrapping the real
  command, rather than the real command running directly as PID 1?

`PID1_COMMAND`'s value contains spaces — quote it
(`KEY="value with spaces"`) so sourcing the file doesn't word-split it.

Get a shell: `docker exec -it breaklab-docker-003 bash`
