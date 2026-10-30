# Solution: Fix Unsafe Defaults

```bash
umask 007
touch /srv/team/newfile.txt
stat -c '%a' /srv/team/newfile.txt   # 660
```

`umask 007` removes nothing from owner or group but strips every bit from
other, so a plain file's default request of `666` comes out as `660`
(`rw-rw----`): the owner and the group can both read and write it, and
"other" has no access at all.
