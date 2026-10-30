# Exit Codes and Pipeline Stages

This container has three scripts representing pipeline stages —
`/opt/pipeline/build.sh`, `/opt/pipeline/test.sh`, and
`/opt/pipeline/deploy.sh` — plus a driver, `/opt/pipeline/run.sh`, that's
supposed to run them in order and stop if one of them fails.

`test.sh` genuinely fails every time it runs — that's fixture data for
this lab, not something to fix. And yet, running `run.sh` deploys anyway:
it creates a marker file at `/workspace/bad-deploy-happened`, proving
`deploy.sh` ran even though the stage before it failed.

## Tasks

1. Run `/opt/pipeline/run.sh` and watch its output. Confirm `test.sh`
   really does fail, and that `deploy.sh` still runs right after it.
2. Read `/opt/pipeline/run.sh`. Find exactly what separates its three
   script invocations, and explain why that separator lets a failure in
   `test.sh` have no effect on whether `deploy.sh` runs.
3. Fix `run.sh` so a failing `test.sh` genuinely stops the pipeline before
   `deploy.sh` ever runs — chain the three invocations with `&&` instead
   of `;` (or add `set -e` near the top of the script).
4. Delete the marker file if it exists (`rm -f /workspace/bad-deploy-happened`),
   then re-run `/opt/pipeline/run.sh` from scratch and confirm the marker
   is never recreated.

Get a shell: `docker exec -it breaklab-cicd-002 bash`
