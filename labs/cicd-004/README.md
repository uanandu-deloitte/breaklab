# Secrets and Config in CI

`/opt/pipeline/deploy.sh` "authenticates" a deploy by comparing a token
against an expected value — except both values are hardcoded literally in
the script itself, so the check always passes for anyone who can read the
file. That's not authentication, it's security theater: the secret is
sitting in plaintext in a script that ships inside this very container.

## Ops gave you this one-time deploy token

```
DEPLOY_TOKEN=brk-ci-7f2a9d41c8b3
```

This is the real, current token. It should never be hardcoded anywhere in
the pipeline's own files — only supplied at the moment deploy.sh actually
runs.

## Tasks

- Refactor `/opt/pipeline/deploy.sh` so it reads its token from the
  `DEPLOY_TOKEN` environment variable instead of a literal baked into the
  script.
- Don't just replace the old hardcoded literal with the real token value
  hardcoded in its place — the real value must only ever be supplied at
  run time, e.g.:

  ```bash
  export DEPLOY_TOKEN=brk-ci-7f2a9d41c8b3
  /opt/pipeline/deploy.sh
  ```

- Confirm deploy.sh still refuses to deploy when `DEPLOY_TOKEN` is unset
  or wrong.

Get a shell: `docker exec -it breaklab-cicd-004 bash`
