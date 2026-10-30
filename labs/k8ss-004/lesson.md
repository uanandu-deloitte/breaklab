# Audit reports violations. Only Enforce blocks them.

Kyverno's `validationFailureAction` has exactly two values that matter
here: `Audit` evaluates every matching resource against the policy and
records a violation in a PolicyReport, but **always admits the
resource anyway**. `Enforce` evaluates the same rule and actually
**rejects** the request at admission time if it violates the policy. A
policy left on `Audit` isn't broken — it's working exactly as
configured, and what it's configured to do is watch, not stop.

## Why this matters

`Audit` is a legitimate, common first phase when rolling out a new
policy (so you can see what it *would* block before anything is
actually at risk of being rejected) — but it's also a trap: a policy
that was deliberately set to `Audit` during rollout and never flipped
to `Enforce` provides zero actual protection while still looking,
operationally, like a real guardrail. Nothing about `kubectl get
clusterpolicy` screams "this isn't blocking anything" unless you read
`validationFailureAction` specifically.

## Where to look

- `kubectl get clusterpolicy disallow-privileged-containers -o yaml` —
  read `spec.validationFailureAction` directly.
- `kubectl get pod bad-pod` — a real, currently-running privileged Pod
  is the proof the policy never actually stopped anything.

## The fix

```yaml
spec:
  validationFailureAction: Enforce
```

Switching the mode doesn't retroactively remove Pods that already
slipped through — `bad-pod` has to be deleted separately. Only new
admission attempts are evaluated against the updated policy.
