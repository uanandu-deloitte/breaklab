# Solution: What Actually Happens When a Container Hits Its Memory Limit

```bash
reason=$(kubectl get pod mem-demo -o jsonpath='{.status.containerStatuses[0].lastState.terminated.reason}')
exit_code=$(kubectl get pod mem-demo -o jsonpath='{.status.containerStatuses[0].lastState.terminated.exitCode}')
limit=$(kubectl get pod mem-demo -o jsonpath='{.spec.containers[0].resources.limits.memory}')

cat <<EOF > /workspace/answers.env
TERMINATION_REASON=$reason
EXIT_CODE=$exit_code
MEMORY_LIMIT=$limit
IS_OOM=yes
EOF
```

`kubectl describe pod mem-demo`'s Last State shows `Reason:
OOMKilled` and `Exit Code: 137` together — the kernel's OOM killer,
unambiguously. `137` by itself would not be enough: it's the generic
exit code for any `SIGKILL`, which a liveness probe failure or a
forced delete could also produce with a `Reason` of plain `Error`
instead. The `Limits` block alongside it shows `memory: 6Mi` — the
exact number the kernel enforced the instant the container's own
allocations (ten 1MiB chunks, comfortably over that limit) crossed it.
Because this is a kernel-level kill rather than an application
exception, there's no stack trace, no log line from the program
itself explaining what happened — the Pod's own last-state fields are
the only record.
