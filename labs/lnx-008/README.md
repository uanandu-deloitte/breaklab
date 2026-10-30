# Find the Stale Backend

Production is supposed to talk to `database.internal`, but the backend
is still failing to connect to the right host. The old hostname,
`database-old.internal`, shows up in more than one place under
`/opt/app` — only one of them is the live configuration the running
service actually reads.

## Tasks

- Find every occurrence of `database-old.internal` under `/opt/app`.
- Identify which one is the active, live configuration versus archived
  or historical references to the old hostname.
- Update only the active configuration to use `database.internal`.
- You can self-check your work with the provided script:
  `/opt/app/bin/check-config.sh`.

## Constraints

- Archived release snapshots and historical documentation must keep
  mentioning the old hostname exactly as they do now — they're a record
  of what used to be true, not something to "fix."
- Do not blindly find-and-replace the string across the whole
  `/opt/app` tree.

Get a shell: `docker exec -it breaklab-lnx-008 bash`
