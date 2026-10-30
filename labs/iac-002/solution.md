# Solution: Preview Before You Apply

```bash
before=$(cat /infra/managed-file.txt)

plan_output=$(infra-apply --plan)
echo "$plan_output"

after_plan=$(cat /infra/managed-file.txt)
[ "$before" = "$after_plan" ] && echo "confirmed: --plan made no changes"

cat <<EOF > /workspace/answers.env
PLANNED_CHANGE="$plan_output"
EOF

infra-apply

cat /infra/managed-file.txt   # now matches resources.yaml's declared content
```

`infra-apply --plan` reads `/infra/resources.yaml` and the real file's
current content and compares them, but only *prints* the difference — it
never writes. That's why `/infra/managed-file.txt` is still exactly what
it was before, immediately after running `--plan`: a plan step describes
a change, it doesn't make one.

Only running `infra-apply` with no flag actually writes
`resources.yaml`'s declared content to the real file, bringing it in
line with what the plan already told you to expect — no surprises,
because you read the diff first.
