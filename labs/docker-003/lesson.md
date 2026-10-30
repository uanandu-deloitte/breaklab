# Shell form vs. exec form, and why PID 1 cares

A Dockerfile's `CMD` (and `ENTRYPOINT`) can be written two ways:

```dockerfile
CMD ["sleep", "infinity"]   # exec form — a JSON array
CMD sleep infinity          # shell form — a plain string
```

Exec form runs the named program directly as PID 1. Shell form is
different: Docker takes the whole string and runs it as
`/bin/sh -c "sleep infinity"`. That means the process Docker actually
starts is `/bin/sh`, not `sleep` — the shell parses the string, then forks
and execs `sleep` as its own child process. PID 1 inside the container is
the shell; the program you actually meant to run ends up as PID 1's child,
one layer removed.

## Why PID 1 is special

The kernel treats PID 1 differently from every other process: it's the
one process a namespace can't do without, and it's the one `docker stop`
targets directly with `SIGTERM`. Whether that signal reaches the program
you care about depends entirely on what's sitting at PID 1:

- If PID 1 **is** the real program (exec form, or a shell-form command
  simple enough for the shell to replace itself via exec), the signal goes
  straight to it.
- If PID 1 **is a shell wrapping the real program** (shell form, with the
  shell staying resident as a parent process), the signal goes to the
  shell. Whether the real program ever finds out depends on whether that
  shell has signal handling set up to forward it — and a plain,
  non-interactive `/bin/sh` generally does not.

## How to check, instead of assuming

Never trust what the Dockerfile's `CMD` line *looks like* it will run as
PID 1 — verify it on the live container:

- `cat /proc/1/cmdline | tr '\0' ' '` prints PID 1's actual argv. If it
  starts with `/bin/sh -c ...`, the shell is PID 1, not your program.
- `ps -p 1 -o comm=` prints just PID 1's command name — `sh` versus the
  name of the program you expected.

This is purely observational: nothing about inspecting PID 1 changes
anything about the container. It's the first thing worth checking whenever
graceful shutdown isn't behaving the way you'd expect.
