# Solution: Log Rotation Breaks Naive File Tailing

```bash
pid=$(cat /var/run/reader.pid)

fd_target=""
for fd in /proc/"$pid"/fd/*; do
  tgt=$(readlink -f "$fd")
  case "$tgt" in
    /var/log/app.log*) fd_target="$tgt" ;;
  esac
done

current_inode=$(stat -c %i /var/log/app.log)
fd_inode=$(stat -c %i "$fd_target")

inodes_match=no
[ "$current_inode" = "$fd_inode" ] && inodes_match=yes

cat <<EOF > /workspace/answers.env
READER_STALE=yes
READER_FD_PATH=$fd_target
INODES_MATCH=$inodes_match
EOF
```

`/proc/<pid>/fd/` lists the reader's open descriptors; resolving each
with `readlink -f` shows that its log-reading descriptor points at
`/var/log/app.log.1` — the file that `app.log` got renamed to during
rotation — not at the current `/var/log/app.log`. `stat -c %i` on both
paths confirms they're genuinely different files now (different inode
numbers): the rotation created a brand-new `app.log`, while the reader
is still attached to the old one under its new name. Since nothing
writes to `app.log.1` anymore, `reader-output.log` stops growing the
moment rotation happens, even though `app.log` keeps gaining new lines —
exactly the silent, errorless gap this lesson is about.
