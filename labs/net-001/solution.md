# Solution: Identify the Network

```bash
# Which interface carries the default route?
ip route
# default via 172.17.0.1 dev eth0

# That interface's own IPv4 address:
ip addr show eth0
# inet 172.17.0.2/16 ...

# DNS resolver in use (inspection only, no answer required):
cat /etc/resolv.conf

# Every TCP port in LISTEN state, and who owns it:
ss -tlnp
# LISTEN ... *:80 ... ("nginx",...)
# LISTEN ... *:8081 ... ("python3",...)

cat <<EOF > /workspace/answers.env
INTERFACE_NAME=$(ip route show default | awk '{print $5; exit}')
CONTAINER_IP=$(ip -4 -o addr show "$(ip route show default | awk '{print $5; exit}')" | awk '{print $4}' | cut -d/ -f1)
LISTENING_PORTS=$(ss -ltn | awk 'NR>1 {print $4}' | sed 's/.*://' | sort -un | paste -sd, -)
EOF
```

`ip route`'s `default via ... dev <iface>` line names the interface that
carries traffic leaving this host — `ip addr show <iface>` then gives
that interface's own address. `ss -ltn`'s `Local Address:Port` column
lists every bound TCP port; stripping everything but the port, sorting
numerically, and deduplicating (nginx binds the same port on both IPv4
and IPv6, which would otherwise count twice) gives the comma-separated
list the check expects.
