# k8ss-004: A Policy Engine That Isn't Enforcing Anything

Real k3s, real Kyverno (admission-controller, background-controller,
reports-controller, cleanup-controller). `disallow-privileged-containers`
is set to `validationFailureAction: Audit`, which only ever reports
violations -- it never blocks them. `bad-pod`, a real privileged Pod,
is running as proof.
