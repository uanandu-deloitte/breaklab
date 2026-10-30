# Solution: Repair the Broken Backup Script

The seeded script has four defects:

1. `cp $f $STAGE/` is unquoted — a filename containing a space gets
   word-split into two arguments, so `cp` fails to find either half.
2. `DEST=/opt/backup-destination` is the wrong destination: that
   directory is never created and was never where backups are expected —
   `/opt/backups` already exists and is the real destination.
3. Neither the `cp` loop nor the `tar` call checks whether it actually
   succeeded.
4. `echo "backup succeeded: ..."` runs unconditionally, even after a
   failure.

```bash
cat <<'SCRIPT' > /opt/app/backup.sh
#!/usr/bin/env bash
set -u
SRC=/data/to-backup
DEST=/opt/backups
STAMP=$(date +%Y%m%d%H%M%S)
STAGE="/tmp/backup-stage-$STAMP"
ARCHIVE="$DEST/backup-$STAMP.tar.gz"

mkdir -p "$STAGE" "$DEST"

for f in "$SRC"/*; do
  cp "$f" "$STAGE"/ || { echo "backup failed: could not copy $f" >&2; exit 1; }
done

if ! tar czf "$ARCHIVE" -C "$STAGE" .; then
  echo "backup failed: tar did not complete successfully" >&2
  exit 1
fi

echo "backup succeeded: $ARCHIVE"
SCRIPT
chmod +x /opt/app/backup.sh
```

Quoting every variable expansion fixes the space-in-filename bug, the
destination now matches where backups are actually expected to land,
both external commands are checked before anything is reported, and the
"succeeded" message only prints on the line that's reachable only after
both of those checks already passed.
