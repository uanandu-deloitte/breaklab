# INC-247: Timeout hotfix isn't live, but nothing looks broken

Someone on the team swears they already fixed the `timeout_seconds`
bug in the repo at /srv/repo before tonight's check ran. `git status`
is clean, `git log` looks completely normal, there's no stray branch
or weird commit — but the deploy check still reads the old value.
Find out where the fix actually went and get it into the working tree.

Get a shell: docker exec -it breaklab-stale-stash bash
