# Memory limits are enforced by the kernel, not the application

A container's CPU limit and memory limit are enforced completely
differently, and conflating them leads to the wrong mental model.
**CPU** is throttled: a container that tries to use more CPU than its
limit allows just gets slowed down — it keeps running, only slower.
**Memory** is not throttled at all. The instant a container's memory
usage crosses its configured limit, the kernel's OOM (out-of-memory)
killer sends it `SIGKILL` immediately. No warning, no graceful
shutdown, no chance for the process to free memory and recover. One
moment it's running; the next, it's dead.

Because `restartPolicy: Always` (the default for most workloads)
restarts a killed container automatically, what this looks like from
the outside is deceptively mundane: a restart count that keeps
climbing. Nothing about `kubectl get pod` alone tells you *why* it
keeps restarting — a plain application crash and a kernel OOM kill
both just show up as "RESTARTS: 4" and climbing.

## Finding the real reason

`kubectl describe pod <name>` is the direct evidence. Its Containers
section includes a "Last State" block with the previous termination's
exact `Reason` and `Exit Code`, right next to the container's
currently configured memory limit:

```
Last State:   Terminated
  Reason:     OOMKilled
  Exit Code:  137
State:        Running
...
Limits:
  memory: 6Mi
```

`OOMKilled` in the reason field is the kernel's own OOM killer,
unambiguously. The exit code alone is **not** enough evidence:
`137` is just `128 + 9`, the generic exit code for any process killed
by `SIGKILL` — a liveness probe failure, `kubectl delete --force`, or
an operator's own script could produce the identical `137` with a
`Reason` of plain `Error` instead. Confirm the reason field before
calling anything an OOM kill.

## Two different root causes, same symptom

An OOMKilled container can mean the limit was set too low for a
program whose own memory use is perfectly reasonable (a limit copied
from a smaller sibling service, say), or it can mean the program's own
memory use grew past a limit that used to be fine (a real leak, or an
unbounded in-memory structure). Both produce the exact same
`OOMKilled` / `137` signature — the limit number itself and knowledge
of what the program is supposed to need are what tell them apart, not
anything in the termination event itself.
