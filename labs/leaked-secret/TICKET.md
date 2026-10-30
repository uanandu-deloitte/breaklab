# INC-330: Secret found in repo history

Security scanned /srv/repo before it goes into a shared archive and
found a live-looking API key somewhere in its history. It's not in the
current files, but `git log` never really forgets. Get it out of every
commit before this repo leaves the building.

Get a shell: docker exec -it breaklab-leaked-secret bash
