# A resource's address is an identity, not a description

Terraform tracks resources in state by their **address**
(`local_file.legal_notice`) — not by what they manage. If you rename
that address in your config without telling state about it, Terraform
sees two unrelated things: an address in state with no matching config
(scheduled for **destroy**) and an address in config with no matching
state (scheduled for **create**). It has no idea they're the same
real-world resource.

## Why this matters

For a `local_file` this is mostly harmless. For almost anything else —
a database, a load balancer, a DNS record — a destroy-then-create cycle
from a pure rename can mean real downtime or, worse, actually losing
the real-world resource (a new cloud resource rarely keeps the old one's
data). Renaming things in `.tf` files is extremely common during
refactors; silently taking the destroy/create path is not acceptable.

## The fix: `terraform state mv`

`terraform state mv <old_address> <new_address>` tells Terraform "this
is the same real resource, just relabeled" — it rewrites state in place,
with no destroy, no create, no real-world effect at all.

```
terraform state mv local_file.legal_notice local_file.privacy_notice
terraform plan -detailed-exitcode   # exit 0 -- nothing left to reconcile
```
