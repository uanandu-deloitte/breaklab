# Solution: Identify the Machine

```bash
cat <<EOF > /workspace/answers.env
DISTRIBUTION=$(. /etc/os-release; echo "$ID")
KERNEL=$(uname -r)
HOSTNAME=$(hostname)
USER=$(whoami)
SHELL=$SHELL
WORKING_DIRECTORY=$(pwd)
ARCHITECTURE=$(uname -m)
EOF
```

Each line pulls its value directly from a live command rather than typing in
a guess, so the file stays correct even if run on a different host.
