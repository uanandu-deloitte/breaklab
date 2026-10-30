# What Actually Happens When a Container Hits Its Memory Limit

The `mem-demo` Pod's restart count keeps climbing. Before assuming
anything, confirm exactly what's killing it and what limit it's
hitting.

## Tasks

1. `kubectl get pod mem-demo` — confirm RESTARTS is climbing.
2. `kubectl describe pod mem-demo` — read the Containers section's
   Last State: the exact termination `Reason` and `Exit Code`, and
   the container's configured memory `Limits`.
3. `kubectl get pod mem-demo -o jsonpath='{.status.containerStatuses[0].lastState.terminated.reason}'`
   and the matching `.exitCode` — confirm the reason specifically,
   not just the exit code (137 alone doesn't prove an OOM kill).
4. `kubectl get pod mem-demo -o jsonpath='{.spec.containers[0].resources.limits.memory}'`
   — read the exact memory limit configured on this container.

Record your findings in `/workspace/answers.env`, one `KEY=value` per
line:

```
TERMINATION_REASON=
EXIT_CODE=
MEMORY_LIMIT=
IS_OOM=
```

- `TERMINATION_REASON`: the exact value of the last termination's
  `reason` field.
- `EXIT_CODE`: the exact value of the last termination's `exitCode`
  field.
- `MEMORY_LIMIT`: the container's configured memory limit, exactly as
  `kubectl` prints it (for example `6Mi`).
- `IS_OOM`: `yes` or `no` — is this genuinely a kernel OOM kill (based
  on the reason field, not the exit code alone)?

Collect every value from the running cluster itself — don't hard-code
anything from the lesson.

Get a shell: `docker exec -it breaklab-k8sd-005 bash`
