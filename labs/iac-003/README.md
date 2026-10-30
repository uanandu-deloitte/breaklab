# Nothing Changes When Nothing Changed

This container has a tiny apply tool, `infra-apply`, and a declared
source of truth at `/infra/resources.yaml`. The real resource it
manages already matches that declaration — it was applied once when
the container was built.

Do the following, in order:

1. Run `infra-apply` once, with nothing edited yet. Read its message —
   it should tell you plainly whether it changed anything.
2. Edit the `value:` line in `/infra/resources.yaml` to any new value.
3. Run `infra-apply` again. Read its message again — this time it
   should report a real change, and the resource file it manages
   should now hold your new value.

Record what you observed in `/workspace/answers.env`, one `KEY=value`
per line:

```
FIRST_RUN_RESULT=
SECOND_RUN_RESULT=
FINAL_VALUE=
```

`FIRST_RUN_RESULT` and `SECOND_RUN_RESULT` should each describe, in a
word or two, whether that run changed anything (for example `no-change`
and `changed`). `FINAL_VALUE` should be the value `/infra/resources.yaml`
declares after your edit.

Get a shell: `docker exec -it breaklab-iac-003 bash`
