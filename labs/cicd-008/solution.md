# Solution: Build Passes, Deploy Never Runs

```yaml
steps:
  - name: Build the app
    id: build
    run: |
      mkdir -p /workspace/dist
      echo "built" > /workspace/dist/app.txt
      echo "ready=true" >> "$GITHUB_OUTPUT"
```

```bash
cd /workspace/project
act push -P ubuntu-latest=-self-hosted
# [CI/build]   🏁  Job succeeded
# [CI/deploy]  🏁  Job succeeded   <- now actually runs

cat /workspace/deployed.txt
# built
```
