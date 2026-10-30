# Prepare a Troubleshooting Host

This host was built from a minimal Ubuntu base image. None of the usual
troubleshooting tools are installed on it yet.

## Task

Using `apt-get`, install exactly these four packages — no more, no less:

- `curl`
- `iproute2`
- `lsof`
- `jq`

Then prove each is actually available and working by running its
version (or help) flag:

```bash
curl --version
ip -V
lsof -v
jq --version
```

Get a shell: `docker exec -it breaklab-lnx-020 bash`
