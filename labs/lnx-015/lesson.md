# Ownership, permissions, and least privilege

A process can only read a file if its own user identity satisfies the
file's permission bits — owner, group, or "other." When a service account
can't read its own configuration, the fix is almost always one of: wrong
owner, wrong group, or a missing read bit, not a reason to loosen
permissions wider than necessary.

## Reading the bits

`ls -l` shows permissions as three triplets — owner, group, other — each
with read/write/execute. `stat -c '%a'` gives the same information as
three octal digits, which is easier to compare against policy
(`640` = owner read/write, group read, no access for everyone else).

## Checking as another user

`su -s /bin/bash -c "<command>" <user>` runs a single command as another
user without giving them an interactive login shell — useful for proving
a fix actually works from the service account's own point of view, not
just as root (which can read anything regardless of permissions).

## Why "least privilege," not "just make it work"

`chmod 777` or making a file world-writable always "fixes" a permission
error — by removing the protection the permissions existed for in the
first place. The correct fix grants the specific identity that needs
access exactly the access it needs, usually through group membership
plus a narrow permission change, and nothing more.
