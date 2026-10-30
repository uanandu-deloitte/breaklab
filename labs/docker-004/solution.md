# Solution: Build Time vs Run Time

```bash
cat <<EOF > /workspace/answers.env
APP_ENV=$(printenv APP_ENV)
LOG_LEVEL=$(printenv LOG_LEVEL)
APP_VERSION=$(printenv APP_VERSION)
OVERRIDE_CANDIDATE=LOG_LEVEL
EOF
```

`APP_ENV`, `LOG_LEVEL`, and `APP_VERSION` were all set with `ENV` in the
Dockerfile, so they show up identically in `env`, `printenv`, and
`/proc/1/environ` — there's only one running environment here, not three
different views of it.

`LOG_LEVEL` is the strongest candidate to override per deployment: how
verbose to log is an operational decision that differs between
environments (quiet in production, noisy while debugging), not something
tied to what got built. `APP_VERSION` identifies the actual build and
should change only when the image itself changes — baking it in is
correct, not a problem. `APP_ENV` is defensible either way in real
systems, but `LOG_LEVEL` is the clearer answer.
