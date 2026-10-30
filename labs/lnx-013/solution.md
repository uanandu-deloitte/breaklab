# Solution: Shared Operations Access

```bash
groupadd ops-team
usermod -aG ops-team developer
usermod -aG ops-team operator

chgrp ops-team /srv/team
chmod 2770 /srv/team
```

`groupadd` creates `ops-team`; `usermod -aG` appends each user to it
without disturbing their existing group memberships. `chgrp` makes
`ops-team` the directory's own group. `chmod 2770` sets read/write/exec
for the owner and the group, nothing for everyone else, and the leading
`2` is the setgid bit — so every file created inside `/srv/team` from
now on inherits the `ops-team` group automatically, regardless of the
creating user's own primary group.
