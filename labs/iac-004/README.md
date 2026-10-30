# Same Template, Different Environment

This container's config was templated out during setup, and whoever
set it up applied it against the **dev** environment's variables to
test it. Ops has since asked you to roll this out to **production** —
the real resource currently shown at `/etc/app/config.yaml` is still
dev-shaped, and that needs to change.

Look at `/infra/resources.yaml.tmpl` to see what's templated, and at
`/infra/vars/` to see what environments exist and what each one
declares. Figure out the right way to apply this template against the
production environment's variables, then do it.

When you're done, `/etc/app/config.yaml` should reflect production's
values, not dev's.

Get a shell: `docker exec -it breaklab-iac-004 bash`
