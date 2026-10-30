# Shared Operations Access

`developer` and `operator` already exist on this host. Both need to
share files in `/srv/team`, which currently exists but is only writable
by root.

## Tasks

- create a group named `ops-team`
- add both `developer` and `operator` to it
- configure `/srv/team` so:
  - its group is `ops-team`
  - new files created inside it automatically inherit the `ops-team`
    group, regardless of the creating user's own primary group
  - the directory is writable by the group

## Constraints

- Don't make `/srv/team` world-writable — only `ops-team` members
  should be able to create files there.

Get a shell: `docker exec -it breaklab-lnx-013 bash`
