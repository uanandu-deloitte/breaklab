# Declaring a version and resolving one are two different steps

Plenty of config declares, somewhere, which version of shared code it
expects to run — a module reference in a `.tf` block, a dependency
line in a manifest, a `module:` field like the one in this lab's
`resources.yaml`. That declaration is intent. It is not, by itself,
what actually runs.

Something else — a pin file, a lockfile, a registry lookup — is the
thing an apply tool actually consults at execution time to resolve
"which implementation do I run right now." In this lab, that's
`/infra/MODULE_VERSION`. `infra-apply` never even reads the `module:`
field in `resources.yaml` — it reads the pin file, full stop. If the
pin is stale (still pointing at an old version after the declaration
moved on), missing, or simply never updated to match, the tool resolves
to whatever the pin says — not whatever the declaration says — and
nothing about that disagreement produces an error.

## A clean exit doesn't mean the right version ran

`infra-apply` in this lab exits 0 whether the pin points at the
version `resources.yaml` declares or a completely different one. A
successful exit only proves the tool ran *some* module without
crashing — it says nothing about whether that module was the one
anyone intended. The only way to know is to check both files directly:
what does the config declare, what does the pin actually say, and do
they agree?

## Why this matters

This exact pattern shows up anywhere code gets resolved indirectly
instead of referenced directly: a Terraform module source pinned to an
old ref, a `package-lock.json` out of sync with `package.json`, a
container image tag that moved out from under a deployment that pinned
a digest months ago. The config can be completely correct and the
pipeline can be completely green, and the thing that actually executed
can still be the wrong version — because resolution and declaration
are two separate steps, and only one of them is usually checked.
