# Solution: Prepare a Troubleshooting Host

```bash
apt-get update
apt-get install -y curl iproute2 lsof jq

curl --version
ip -V
lsof -v
jq --version
```

`apt-get update` refreshes the package index first, so `apt-get install`
can actually find the four requested packages. Each verification command
confirms its tool is genuinely installed, on `PATH`, and executable —
not just that `apt-get` reported success.
