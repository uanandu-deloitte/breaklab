# A Policy Engine That Isn't Enforcing Anything

This cluster runs a real Kyverno (the actual admission-controller,
background-controller, reports-controller, and cleanup-controller — a
real policy engine, not a simulation). A `ClusterPolicy` called
`disallow-privileged-containers` exists. A privileged Pod, `bad-pod`,
is running right now in the `default` namespace.

That should not be possible.

## Tasks

1. Look at the real policy and its real effect:

   ```
   kubectl get clusterpolicy disallow-privileged-containers -o yaml
   kubectl get pod bad-pod
   ```

2. Figure out why a Pod that clearly violates this policy was admitted
   anyway, fix the policy for real, and clean up the Pod that
   shouldn't be there.

3. Prove it: attempting to create a privileged Pod again should now be
   rejected by Kyverno's real admission webhook, not just theoretically
   blocked.

Get a shell: `docker exec -it breaklab-k8ss-004 bash`
