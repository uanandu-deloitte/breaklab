# A Secret is a different resource type, not a different encoding of the same one

`ConfigMap` and `Secret` have nearly identical shapes — both are just
key/value data a Pod can consume as environment variables, command-line
arguments, or mounted files. It's tempting to treat the choice between
them as cosmetic. It isn't: they're different API resource types with
different handling guarantees, and only one of them is meant to hold
anything sensitive.

## What actually differs

A `ConfigMap`'s `data` is stored and returned by the API server exactly
as plain text, with no special handling. A `Secret`'s `data` is
base64-encoded in the API representation and, depending on the cluster's
configuration, may be encrypted at rest in etcd, have tighter default
RBAC expectations, and is excluded from things like `kubectl describe`'s
normal output (which shows `<set>` rather than the value for Secrets, but
prints a ConfigMap's data plainly). None of that happens for a ConfigMap,
no matter how sensitive the value you put in it — the resource type is
what triggers the special handling, not anything about the value itself.

## Base64 is encoding, not encryption

This is the detail that trips people up: a `Secret`'s `data` field is
base64, and base64 is **reversible by design** — it exists to represent
arbitrary bytes as ASCII text in a JSON/YAML-safe way, nothing more.
Anyone who can `kubectl get secret -o yaml` (or `-o jsonpath`) a Secret
they have read access to can decode it in one command:
`kubectl get secret <name> -o jsonpath='{.data.<key>}' | base64 -d`. There
is no cryptographic step here at all.

This means a Secret's actual security property has nothing to do with
its encoding and everything to do with **who is allowed to read the
Secret object in the first place.** RBAC `get`/`list` permission on
`secrets` is the real access boundary — restricting that (and, in a
production cluster, enabling encryption-at-rest for etcd) is what
actually protects a Secret's contents. Switching a value from a
ConfigMap to a Secret without also thinking about who can read Secret
objects in that namespace closes none of the actual exposure; it just
moves the plaintext one `base64 -d` further away from a casual
`kubectl get -o yaml`.

## Why migrate at all, then

Even though the encoding itself isn't protection, moving sensitive
values into Secrets still matters: it's what makes the RBAC boundary
possible to draw in the first place. You can scope a Role to allow
reading ConfigMaps but not Secrets in a namespace — that distinction is
meaningless if a token is sitting in a ConfigMap's `data` next to
ordinary log-level config. Putting sensitive values in the resource type
built for them is what lets the rest of your access control actually
mean something.
