# Build Passes, Deploy Never Runs

`/workspace/project/.github/workflows/ci.yml` is a real workflow with two
jobs: `build`, then `deploy`. Every push shows `build` going green — and
`deploy` never runs at all, not even as a failure. It's just skipped,
every single time.

Run it yourself with `act` (the same real workflow file GitHub Actions
would run, executed locally):

```
cd /workspace/project
act push -P ubuntu-latest=-self-hosted
```

Figure out why `deploy`'s `if:` condition never comes out true, and fix
the real workflow YAML so it does.

Get a shell: `docker exec -it breaklab-cicd-008 bash`
