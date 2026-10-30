# Solution: What Even Is a Container

```bash
process_count=$(ps aux --no-headers | wc -l)
live_hostname=$(hostname)

cat <<EOF > /workspace/answers.env
PROCESS_COUNT_SEEN=$process_count
KERNEL_SHARED=yes
HOSTNAME_SOURCE=container-runtime
EOF
```

`ps aux` inside the container can only see processes in this container's
own PID namespace, so the count is a small handful, never the hundreds a
real host would show. `uname -r` reports the host's kernel version — a
container image never bundles a kernel, so this is always the host's,
regardless of which image is running. The hostname is a 12-character hex
string matching `/etc/hostname`: the container runtime (Docker) assigns it
as the container's own short container ID the moment the container starts
— nothing in this lab's Dockerfile ever sets a hostname.
