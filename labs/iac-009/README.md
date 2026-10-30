# The Config Says One Thing, The State Remembers Another

Whoever's on call handed you a one-liner: "today's cart-fix release
updated the banner, can you confirm it's live?" `/etc/motd-app/banner.txt`
is the file every shell session on this box prints on login, and it's
supposed to reflect whatever `/infra/main.tf` currently declares.

Check whether it actually does. If it doesn't, figure out why — and fix
it the way this stack is supposed to be fixed, not by hand-editing the
file.

Get a shell: `docker exec -it breaklab-iac-009 bash`
