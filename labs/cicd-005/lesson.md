# A lock only means something while its owner is alive

A lock file is a cheap, common way to stop two deploys from racing each
other: before doing anything real, a script checks whether the lock
exists, creates it if not, does its work, then removes it. As long as
every run reliably removes its own lock when it's done, this works fine.

The failure mode is just as simple: whatever creates the lock has to be
the thing that removes it, and if that process is killed, crashes, or
loses power between "create" and "remove," the lock survives it. From
that point on, every future run sees the lock, assumes a deploy is still
in progress, and refuses to start — forever, or until a human notices
and clears it by hand.

## The lock file by itself can't tell you anything

A bare empty lock file (just its existence) gives you no way to tell
"a deploy is genuinely running right now" apart from "a deploy crashed
three weeks ago and never cleaned up." Both look identical: a file that
exists. That ambiguity is exactly why some locks carry metadata — a PID,
a timestamp, a hostname — recording *who* created the lock and *when*.
That metadata turns "does this file exist" into a question you can
actually answer: "is the process that created this still alive?"

## Reading the metadata is the whole job

Given a lock with a recorded PID:

- Check whether that PID is a real, running process (`ps -p <pid>`, or
  `kill -0 <pid>` which succeeds only if a process with that PID exists
  and you have permission to signal it). If it isn't running, nothing
  is actually using this lock.
- Check the recorded timestamp against now. A lock that's hours old
  might still be a legitimately slow deploy. A lock that's weeks old is
  not — no real deploy step runs that long.

Either signal alone is already fairly convincing; together, they turn
"I think this lock is stale" into "I confirmed this lock is stale."
Only once you've actually checked — not assumed — is it safe to clear
it and let deploys resume.

## Why this matters

Treating every lock as sacred, untouchable state leads to the opposite
failure: a genuinely stale lock blocks legitimate work indefinitely,
and nobody wants to be the one who deletes a lock without being sure.
Knowing exactly what a lock's own metadata tells you — and what it
doesn't — is what lets you clear a stale one with confidence instead of
guessing.
