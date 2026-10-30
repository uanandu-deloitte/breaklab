# The Playbook That's Never Actually Done

`/infra/playbook.yml` is supposed to make sure `/etc/app/app.conf` has
`max_connections=200` set. Run it once — fine. Run it again — it's
still "changing" something, every single time, forever.

A playbook that reports `changed` on every run isn't safe to run on a
schedule or in CI: nobody can tell a real change from this one repeating
itself. Fix the task so running it twice in a row is a true no-op the
second time, with the config still correct.

Get a shell: `docker exec -it breaklab-iac-011 bash`
