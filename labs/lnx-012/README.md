# Create Operations Accounts

Two new operations roles need Linux accounts on this host:

- `developer` — a human engineer who needs to log in interactively.
  Give them a real home directory under `/home/developer` and the
  `/bin/bash` shell.
- `operator` — a service account that should never get an interactive
  login shell at all.

## Tasks

- create the `developer` user, with shell `/bin/bash` and a home
  directory created at `/home/developer`
- create the `operator` user, with shell `/usr/sbin/nologin`

## Constraints

- Both accounts must end up with distinct UIDs — from each other and
  from every account that already exists on this host.

Get a shell: `docker exec -it breaklab-lnx-012 bash`
