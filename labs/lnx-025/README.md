# Repair the Broken Backup Script

`/opt/app/backup.sh` is supposed to back up `/data/to-backup` into a
`tar.gz` archive. It runs without any obvious error, but the backups it
produces aren't trustworthy — it has four separate defects.

## Tasks

Fix `/opt/app/backup.sh` so that:

- it correctly backs up every file in `/data/to-backup`, including the
  one with a space in its name
- the resulting archive lands at the correct destination — `/opt/backups`
  already exists and is where the backup is expected to end up
- it checks whether the backup actually succeeded, instead of assuming
  it did
- it only reports success when the backup genuinely succeeded, and
  reports failure (with a non-zero exit code) otherwise

Get a shell: `docker exec -it breaklab-lnx-025 bash`
