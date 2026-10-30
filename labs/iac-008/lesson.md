# Declared isn't the same as managed

It's tempting to assume that once a resource is written down in your
source of truth, an apply tool will eventually get around to applying
it. Most real apply tools don't work that way. They maintain their own
separate record — a state file — of exactly which resources they're
allowed to touch, and they only ever act on resources that appear in
that record. A resource can be completely, correctly declared and
still never get applied, because the tool was never told it owns that
resource.

This isn't a bug. It's a deliberate safety mechanism: without it, an
apply tool pointed at the wrong directory, or handed a declaration file
with a typo'd path, could start overwriting files it has no business
touching. Refusing to act on anything outside its own tracked list is
what makes "declare it and the tool does the rest" safe to use at all.
The cost is that a resource genuinely needs two things to be true
before it ever gets written — a declaration, and a matching entry in
state — and having only one of them produces no error, just silence.

## What this looks like from the outside

In this lab, `infra-apply` runs clean every time, exit code 0, no
error text thrown at you. But its own printed output tells the real
story if you read it: one resource says `applied: ...`, the other says
`skipping ...: not in state, run with --import first`. The tool is
telling you exactly what's wrong, in its own normal output — the only
mistake is not reading past "it didn't crash."

## Why this matters

This is the same shape as a Terraform resource that exists for real in
a cloud account but was never run through `terraform import` — perfectly
real, perfectly declared in a `.tf` file somewhere, and permanently
invisible to every `terraform apply` until someone explicitly imports
it into state. The fix is never to "try applying harder." It's to
compare what's declared against what the tool's own state file
actually tracks, find the gap, and close it explicitly.
