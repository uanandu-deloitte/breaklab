# Solution: PID 1 and the Command

```bash
pid1_command=$(tr '\0' ' ' < /proc/1/cmdline | sed -e 's/[[:space:]]*$//')

cat <<EOF > /workspace/answers.env
PID1_COMMAND="$pid1_command"
IS_SHELL_FORM=yes
EOF
```

This lab's Dockerfile writes `CMD sleep infinity` — shell form, a plain
string rather than a JSON array (`CMD ["sleep", "infinity"]` would be exec
form). Docker runs shell-form `CMD` as `/bin/sh -c "sleep infinity"`, so
`/proc/1/cmdline` shows `/bin/sh -c sleep infinity`: PID 1 really is the
shell. `ps -p 1 -o comm=` confirms it reports `sh`, and `ps -eo
pid,ppid,comm` shows a separate `sleep` process whose `ppid` is `1` — the
shell forked it as a child instead of becoming it. If a signal like
`SIGTERM` were sent to this container, it would reach the shell at PID 1
first, and whether `sleep` ever finds out depends entirely on whether that
shell forwards it — which a plain, non-interactive `/bin/sh` generally
does not do on its own.
