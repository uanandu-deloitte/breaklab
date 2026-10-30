# INC-205: Release commit missing from main

Whoever shipped tonight's release swears they committed
"RELEASE: v1.2.0 going out tonight" in the repo at /srv/repo, but
`git log main` doesn't show it and RELEASE.md on main doesn't mention
v1.2.0. Get that commit onto main without losing anything else that's
there now.

Get a shell: docker exec -it breaklab-git-lost-work bash
