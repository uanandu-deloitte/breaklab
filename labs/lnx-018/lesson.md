# Trapped signals vs. unconditional termination

A process can choose how it reacts to most signals by installing a
**trap** — a handler that runs instead of (or before) the signal's
default action. `SIGTERM` is the standard "please shut down" signal, and
well-behaved long-running processes trap it to flush buffers, close
connections, or write a final status before exiting on their own terms.

## Why SIGKILL is different

`SIGKILL` (signal 9) cannot be trapped, blocked, or ignored. The kernel
terminates the process immediately, without ever handing control back to
it — no trap runs, no cleanup code executes, nothing gets written. This
is exactly what makes `SIGKILL` dangerous for anything that hasn't already
reached a safe state to stop in: whatever the process was in the middle
of just stops, mid-step.

## What this looks like in practice

- `kill <pid>` sends `SIGTERM` by default. A process with a `SIGTERM`
  trap gets to run it, then typically calls `exit` itself.
- `kill -9 <pid>` (or `kill -KILL <pid>`) sends `SIGKILL`. The same
  process, given this instead, vanishes with no trap run at all — any
  marker file, log line, or cleanup step its trap would have produced
  simply never happens.

## Why this matters operationally

`SIGKILL` is sometimes necessary — a process that's hung and ignoring
`SIGTERM` entirely eventually needs it. But reaching for `kill -9` as a
first resort on every stuck or slow-to-stop process throws away the one
chance a well-written service had to shut down cleanly. The standard
sequence is: `SIGTERM` first, wait a reasonable grace period, `SIGKILL`
only if it's still running afterward.
