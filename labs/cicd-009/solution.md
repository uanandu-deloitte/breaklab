# Solution: Deploy Is Never Even Attempted

```yaml
deploy:
  stage: deploy
  needs: ["build"]
  rules:
    - if: '$CI_COMMIT_BRANCH'
  script:
    - cp /workspace/dist/app.txt /workspace/deployed.txt
    - echo "deployed"
```

```bash
cd /workspace/project
gitlab-ci-local --force-shell-executor
#  PASS  build
#  PASS  deploy   <- now actually runs

cat /workspace/deployed.txt
# built
```
