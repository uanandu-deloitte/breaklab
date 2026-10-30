# Fix Unsafe Defaults

`/srv/team` is a shared workspace for the `ops-team` group. This shell's
current umask leaves every newly created file too open: `other` (everyone
who isn't the owner or in the group) ends up with access it should never
have.

## Task

1. Check the umask currently in effect for your session (`umask`).
2. Set a umask for your session such that any **new** file ends up:
   - readable **and** writable by its owning group
   - completely inaccessible to "other" (no read, no write, no execute)
3. Prove it: with that umask active, create a literal file at
   `/srv/team/newfile.txt` (e.g. `touch /srv/team/newfile.txt`).
4. Confirm its resulting permissions with `ls -l /srv/team/newfile.txt` or
   `stat -c '%a' /srv/team/newfile.txt`.

Get a shell: `docker exec -it breaklab-lnx-016 bash`
