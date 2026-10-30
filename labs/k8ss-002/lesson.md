# A Pod with no securityContext inherits the worst defaults

Leave `securityContext` off a Pod entirely and Kubernetes doesn't pick a
safe default on your behalf — the container runs exactly as its image's
`Dockerfile` says, which for most base images means root (`uid 0`), with
the full container filesystem mounted read-write. That's not a
Kubernetes failure; it's just what happens when nothing more restrictive
was ever asked for.

## Two different securityContexts, two different scopes

`PodSecurityContext` (at `spec.securityContext`) applies to every
container in the Pod: this is where `runAsUser` and `runAsNonRoot` live.
`runAsNonRoot: true` doesn't pick a uid for you — it's a guard that makes
the kubelet refuse to start the container at all if it would otherwise
run as uid 0, so you always pair it with an explicit `runAsUser`.

A container's own `securityContext` (at
`spec.containers[].securityContext`) overrides or adds to the Pod-level
one, and some fields only exist here at all —
`readOnlyRootFilesystem` is one of them. There is no Pod-level
"read-only" switch; it's a per-container setting because different
containers in the same Pod can legitimately need different filesystem
write access.

## Read-only root filesystem means exactly that

`readOnlyRootFilesystem: true` doesn't make specific paths read-only —
it makes the entire container filesystem (the image layer) read-only,
full stop. Nothing distinguishes "a path this app happens to need" from
"a path an attacker would want" at that layer; the setting doesn't know
the difference, which is exactly why it's effective. Anything the
container legitimately needs to write — logs, cache files, a scratch
directory — has to be explicitly carved out as its own writable volume
mounted over a path on that otherwise-frozen filesystem.

`emptyDir` is the simplest such volume: a fresh, empty directory created
for the Pod's lifetime, mounted wherever the container needs write
access. It's not backed by the image layer at all, so it stays writable
even while everything else in the container is locked down.

## Why both matter together

Non-root and read-only-root address different blast radii. Running as a
non-root uid means a process that gets compromised doesn't automatically
have root-level control of the Pod's own user context. A read-only root
filesystem means even a compromised process can't persist a change to
the container's own code or binaries — any malicious write attempt
simply fails. Applied together, with a narrowly-scoped writable volume
for whatever legitimate writes still need to happen, a container goes
from "does its job, wide open if anything goes wrong" to "does its job,
with a wall around everything it doesn't need."
