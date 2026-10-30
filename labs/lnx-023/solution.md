# Solution: Find the Failure in the Logs

```bash
# Find the user-visible failure first.
grep ERROR /var/log/frontend.log
# 2026-01-15T10:22:06Z ERROR frontend checkout failed: 500 from api req=req-8841

# Follow its request ID across all three logs, in order.
grep req-8841 /var/log/frontend.log /var/log/api.log /var/log/payments.log
```

Reading the matched lines in timestamp order:

```
api.log:      10:21:30 handling req=req-8841
api.log:      10:21:31 calling payments service req=req-8841
payments.log: 10:21:31 received authorization request req=req-8841
payments.log: 10:21:33 WARN ledger-service response slow req=req-8841
payments.log: 10:22:01 ERROR upstream timeout contacting ledger-service req=req-8841
api.log:      10:22:05 ERROR payment authorization failed: upstream timeout req=req-8841
api.log:      10:22:06 ERROR returning 500 req=req-8841
frontend.log: 10:22:06 ERROR checkout failed: 500 from api req=req-8841
```

`api.log`'s 500 is the error the customer saw, but it isn't the root
cause — it's `payments.log`'s own upstream timeout at 10:22:01, calling
out to `ledger-service`, that started the chain. `payments` is the
failing dependency.

```bash
cat <<EOF > /workspace/answers.env
FAILING_DEPENDENCY=payments
REQUEST_ID=req-8841
EOF
```
