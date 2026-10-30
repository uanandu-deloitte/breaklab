# Solution: Find It Yourself

```bash
# Discover the options (any of these work):
ls --help | grep -i hidden
mkdir --help | grep -i parent
cp --help | grep -i prompt

cat <<EOF > /workspace/answers.env
HIDDEN_FILES_OPTION=-a
NESTED_MKDIR_OPTION=-p
INTERACTIVE_COPY_OPTION=-i
EOF

touch /workspace/.hidden-check
mkdir -p /workspace/deep/nested/tree
cp -i /opt/fixtures/readme.txt /workspace/fixture-copy.txt
```

`ls -a` lists entries whose name starts with `.`, including `.` and
`..`. `mkdir -p` creates every missing directory on the way to the final
path instead of refusing because a parent doesn't exist yet. `cp -i`
prompts before silently overwriting a destination that already exists —
here there's nothing to overwrite yet, but the flag is what demonstrates
you found and used the right option.
