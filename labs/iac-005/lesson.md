# A lock is only as good as its cleanup

Two copies of the same apply tool running against the same resource at
the same time is a real hazard — they can interleave writes and leave
the resource in a state neither of them actually intended. A lock file
is the usual fix: before doing any real work, a script writes a marker
(often including its own PID and a timestamp); before starting, it
checks whether that marker already exists, and refuses to run if so.

This only works if the script reliably removes its own lock when it's
done — including when it's done because it crashed. A lock file has no
idea whether the process that created it is still alive. It's just a
file. If the process holding it dies without running its cleanup step
(a `trap ... EXIT`, a `finally` block, whatever the language calls it),
the lock stays on disk forever, and every future run refuses to start —
not because anything is actually in progress, but because the lock
can't tell the difference between "still running" and "never cleaned up
after itself."

## Stale vs. live: how to actually tell

A lock file refusing a run looks identical whether it's protecting a
real in-progress apply or just rotting there from a crash. The fix is
never to guess — it's to check:

- **Is the recorded PID actually alive?** A well-behaved lock records
  the PID of the process that created it. `ps -p <pid>` either finds a
  real, running process or it doesn't. If it doesn't, nothing is
  holding this lock — it's just a leftover file.
- **How old is it?** `stat` shows a lock's modification time. A lock
  that's been sitting for far longer than any real apply should ever
  take is further evidence it's stale, even before you check the PID.

Only once you've confirmed the PID is dead is it safe to remove the
lock. Removing a lock that's genuinely still held by a live process
re-introduces the exact race condition the lock exists to prevent.

## Why this matters

Lock files show up anywhere two things must never run at once: package
manager installs, database migrations, deploy scripts, this lab's tiny
`infra-apply`. Every one of them has the same failure mode — a crash
that skips cleanup leaves a stale lock that blocks every future run
until a human steps in and clears it. Knowing how to tell "stale" from
"live" before touching anything is what makes clearing a lock safe
instead of reckless.
