# Users, UIDs, and login shells

Every Linux account is a line in `/etc/passwd`: a username, a numeric
UID, a home directory, and a login shell. The shell field is what
actually determines whether an account can get an interactive session —
not anything special about the account "type."

## Two different kinds of account

- A **human** account needs somewhere to put files (a home directory)
  and a real shell to log in with — typically `/bin/bash`.
- A **service** account runs a daemon or background job and should never
  accept an interactive login at all. Setting its shell to
  `/usr/sbin/nologin` (or `/sbin/nologin`, depending on the
  distribution) means any attempt to log in as that account — `su`, SSH,
  anything that execs the user's shell — is refused immediately, even if
  someone has the password.

## useradd's relevant flags

- `-m` creates the home directory if it doesn't already exist, and
  populates it from `/etc/skel`.
- `-s <path>` sets the login shell explicitly.
- `-r` creates a "system" account — conventionally used for service
  accounts, and on most distributions pulls its UID from a separate,
  lower range than ordinary human accounts.

Every UID `useradd` assigns is unique by construction — it refuses to
reuse one that's already taken — so two accounts created with plain
`useradd` calls are automatically guaranteed distinct UIDs.

## Checking what you got

`id <user>` and `getent passwd <user>` both show the UID, home
directory, and shell a new account actually ended up with — useful for
confirming the account looks the way you intended before relying on it.
