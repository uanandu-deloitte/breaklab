# Solution: Create Operations Accounts

```bash
useradd -m -s /bin/bash developer
useradd -r -s /usr/sbin/nologin -g operator operator
```

`-m` creates `/home/developer` and populates it from `/etc/skel`; `-s
/bin/bash` gives `developer` a real interactive shell. `-r` marks
`operator` as a system account and `-s /usr/sbin/nologin` means any
attempt to log in as `operator` is refused immediately. `useradd`
assigns each account the next free UID, so the two end up with distinct
UIDs automatically.

Ubuntu's base system already reserves a group named `operator` (a
standard low-numbered system group); `useradd` would otherwise try to
create a new group of the same name for the account and refuse because
that name is already taken. `-g operator` tells it to use the existing
group as this account's primary group instead of creating a new one.
