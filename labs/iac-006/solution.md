# Solution: Config Precedence: File vs Environment vs Flag

```bash
file_value=$(sed -n 's/^value: *//p' /infra/resources.yaml)
env_value="$INFRA_VALUE"

infra-apply
# applied (source: env): maintenance-mode-on

cat /etc/app/mode.conf
# maintenance-mode-on

infra-apply --value=maintenance-mode-forced
# applied (source: flag): maintenance-mode-forced

cat <<EOF > /workspace/answers.env
FILE_VALUE=$file_value
ENV_VALUE=$env_value
ACTUAL_SOURCE=env
RESOLVED_VALUE=$env_value
EOF
```

`resources.yaml` declares `maintenance-mode-off`, and `$INFRA_VALUE` is
set to `maintenance-mode-on` — baked into this container's own
environment, not hidden in any rc file or wrapper. `infra-apply` checks
a CLI flag first, then the environment variable, and only falls back
to the file if neither is set. With no flag given, the environment
variable wins, so `/etc/app/mode.conf` ends up holding the env value,
not the file's — even though the file was never wrong. Passing
`--value=...` demonstrates the top of the order: a flag beats both the
environment and the file.
