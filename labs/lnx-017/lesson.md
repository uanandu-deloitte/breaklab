# Finding a runaway process, and stopping it politely

A "runaway" process isn't identified by its name — plenty of legitimate
services have unremarkable names, and nothing stops a misbehaving one from
having a perfectly innocent-sounding command. It's identified by what it's
actually doing to the machine: usually pinning a CPU core, or consuming
memory without bound.

## Finding it

- `ps aux --sort=-%cpu | head` sorts every process by CPU usage, highest
  first — a process stuck in a tight loop stands out immediately in the
  `%CPU` column, regardless of what it's called.
- `top` shows the same thing live, refreshing continuously, which is
  useful for confirming the CPU usage is sustained and not a brief spike.
- `ps -eo pid,ppid,user,cmd` prints exactly the four fields you need to
  record about any process you've identified, in one line per process.

## Signals: ask first, force second

`kill <pid>` sends `SIGTERM` by default — a request that the process
clean up and exit on its own terms (closing files, flushing buffers,
releasing locks). Most processes honor it immediately. `kill -9 <pid>`
sends `SIGKILL`, which the kernel enforces unconditionally and the target
process never even sees — it cannot clean up, trap the signal, or object.

Reaching for `SIGKILL` first on a process that would have exited cleanly
on `SIGTERM` is how corrupted state and half-written files happen. Send
`SIGTERM` first, and only escalate to `SIGKILL` if the process is still
running after a few seconds.

## Recording before you act

Capturing a process's PID, parent PID, user, and command *before*
terminating it — not from memory afterward — is what makes an incident
record trustworthy. Once you've sent the signal, the process (and the
easy way to answer "what was that, exactly?") is gone.
