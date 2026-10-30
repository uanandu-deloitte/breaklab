# INC-388: Pipeline green, but prod is still running the old dependency

We bumped `libwidget` in /srv/app/requirements.lock for a security fix.
The CI build (`ci-build`) ran clean and reported success, same as always.
Prod is still hitting the bug the bump was supposed to fix. Figure out
why the deployed vendor copy doesn't match requirements.lock and get it
actually in sync.

Get a shell: docker exec -it breaklab-stale-cache-restore bash
