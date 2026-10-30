# Deploy Is Never Even Attempted

`/workspace/project/.gitlab-ci.yml` is a real pipeline with two stages:
`build`, then `deploy`. `build` goes green on every push. `deploy`
never runs at all — no failure, no error, it's just absent from every
pipeline.

Run it yourself with `gitlab-ci-local` (the same real pipeline file
GitLab CI would run, executed locally):

```
cd /workspace/project
gitlab-ci-local --force-shell-executor
```

Figure out why `deploy`'s `rules:` never matches on a normal push, and
fix the real pipeline YAML so it does.

Get a shell: `docker exec -it breaklab-cicd-009 bash`
