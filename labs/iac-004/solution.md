# Solution: Same Template, Different Environment

```bash
infra-apply prod
```

`infra-apply` takes the environment name as its one argument, loads
`/infra/vars/prod.env` (`ENVIRONMENT=prod`, `REPLICA_COUNT=5`), and
substitutes those values into `/infra/resources.yaml.tmpl`'s
placeholders before writing the rendered result to
`/etc/app/config.yaml` — overwriting the dev-shaped content that was
applied there before.
