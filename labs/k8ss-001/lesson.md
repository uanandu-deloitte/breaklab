# RBAC grants exactly what its rules say, nothing you assumed

Kubernetes RBAC has four moving pieces: a subject (a user, group, or
`ServiceAccount`), a `Role` (or `ClusterRole`) listing rules, a
`RoleBinding` (or `ClusterRoleBinding`) connecting a subject to a role,
and the actual permission check the API server runs on every request.
Nothing here is implicit — a `Role`'s name ("pod-reader") is just a
label a human chose; the API server only ever looks at the rules
themselves.

## A rule is `apiGroups` + `resources` + `verbs`

Each rule in a `Role` is a cross product: it allows every verb in its
`verbs` list, on every resource in its `resources` list, in every API
group in its `apiGroups` list. There is no partial credit and no
inference — a rule that lists `["get", "delete"]` grants exactly those
two verbs, never anything in between, and never anything on a resource
type that isn't named.

This cuts both ways, and both ways matter:

- **Too narrow** is a correctness bug: a Pod reading its own Role will
  fail to `list` pods if `list` was never granted, even if `get` works
  fine — they are two separate verbs, not one "read" permission.
- **Too broad** is a security bug: a Role that still grants `delete`
  "just in case" gives any identity bound to it the ability to delete
  every pod in the namespace — an attacker who compromises that
  identity inherits that blast radius for free.

A Role is correct only when it grants precisely the verbs its subject
needs, on precisely the resources it needs — not more, not less.

## RBAC is additive, but it's still worth auditing exactly

RoleBindings only ever add permissions; there's no "deny" rule in core
RBAC. That makes the audit simple in one direction (if you don't see a
rule granting it, it isn't granted — unless some *other* binding grants
it too) but it means every `Role` bound to a subject deserves the same
scrutiny: read what it actually says, not what its name implies a
"least privilege" pass was supposed to leave behind.

## `kubectl auth can-i` asks the API server's own question

`kubectl auth can-i <verb> <resource> --as=<subject> -n <namespace>`
doesn't simulate RBAC — it runs the exact same `SubjectAccessReview`
check the API server runs for a real request from that identity. This
is the only verification that actually matters: not reading the YAML
and reasoning about it, but asking the API server itself what it would
do, for every verb your intended permission matrix depends on —
including the "should be denied" verbs, since a Role that's
accidentally too permissive looks identical to a correct one if you
only ever check the verbs it's supposed to allow.
