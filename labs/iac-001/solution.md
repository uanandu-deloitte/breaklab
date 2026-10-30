# Solution: Declarative vs Imperative

```bash
before=$(sed -n 's/^content: *//p' /infra/resources.yaml)

sed -i 's/^content: .*/content: Build server banner v2/' /infra/resources.yaml

infra-apply

after=$(sed -n 's/^content: *//p' /infra/resources.yaml)

cat <<EOF > /workspace/answers.env
DESIRED_VALUE_BEFORE="$before"
DESIRED_VALUE_AFTER="$after"
SOURCE_OF_TRUTH_FILE=/infra/resources.yaml
EOF

cat /etc/banner.txt   # now prints "Build server banner v2"
```

`/infra/resources.yaml` is the declaration; `/etc/banner.txt` is just
whatever `infra-apply` last wrote there. Editing `/etc/banner.txt`
directly would "work" for about as long as it takes someone to run
`infra-apply` again — at which point the manual edit is silently
overwritten back to whatever the declaration says. Changing the
declaration first, then re-applying, is the only edit that survives.
