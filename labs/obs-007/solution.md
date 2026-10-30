# Solution: Why Clocks Have to Agree Across Hosts

```bash
real_now=$(date +%s)
log-event-a "check"
log-event-b "check"

ts_a=$(tail -n 1 /var/log/host-a.log | awk '{print $2}')
ts_b=$(tail -n 1 /var/log/host-b.log | awk '{print $2}')

skew_a=$(( ts_a - real_now )); [ "$skew_a" -lt 0 ] && skew_a=$(( -skew_a ))
skew_b=$(( ts_b - real_now )); [ "$skew_b" -lt 0 ] && skew_b=$(( -skew_b ))

if [ "$skew_a" -gt "$skew_b" ]; then
  drifted=host-a; skew=$skew_a
else
  drifted=host-b; skew=$skew_b
fi

cat <<EOF > /workspace/answers.env
DRIFTED_HOST=$drifted
SKEW_SECONDS=$skew
IS_REAL_ORDERING_ISSUE=no
EOF
```

`host-a`'s clock carries a baked-in +125 second offset, so every
timestamp it produces — including the fresh one from `log-event-a
"check"` — lands about 125 seconds ahead of the real `date +%s` value
captured just before it. `host-b`'s clock has no offset at all, so its
timestamp lines up with real time within a second or two of normal
command overhead. The ~125 second gap is far too large and far too
consistent to be ordinary noise or a real causal delay — it's `host-a`'s
clock itself that's wrong, not the sequence of events it's describing.
