# Which permission triplet actually applies to you

Every file has three permission triplets — owner, group, and other — but
any one user is only ever judged against exactly one of them for a given
access.

## Picking the right triplet

For a given user and file, the kernel checks in this order:

1. Is the user the file's owner? If so, the **owner** triplet decides
   access — full stop, even if the group or other triplets would have
   allowed more.
2. Otherwise, is the user a member of the file's owning group (primary
   or supplementary)? If so, the **group** triplet decides.
3. Otherwise, the **other** triplet decides.

A user who isn't the owner and isn't in the group never falls back to a
more permissive triplet just because one exists — if the group triplet
denies them and they're in that group, that's the final answer for them,
regardless of what "other" allows.

## Reading the bits

`ls -l` shows a 10-character string like `-rw-r-----`: the first
character is the file type, then three groups of three (owner, group,
other), each `r`/`w`/`x` or `-`. `640` in octal means owner
read+write, group read-only, other nothing at all — exactly
`-rw-r-----`.

## Don't test as root

Root bypasses permission checks entirely, so testing access by just
reading the file yourself as root (or inside the container's default
shell, which may already be root) tells you nothing about what a
specific unprivileged user can do. `su -s /bin/bash -c "test -r <file>"
<user>` runs a single command as that user and reports success or
failure through its exit code — a real test of their actual access.
