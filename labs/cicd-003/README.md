# Artifacts Between Stages

A pipeline splits work into stages, and each stage hands its output to the
next one as an *artifact* — a file whose name and location the next stage
has to agree on. When one stage changes what it produces without the next
stage changing what it expects, the handoff breaks.

## The pipeline here

- `/opt/pipeline/build.sh` packages the application into a tarball under
  `/opt/pipeline/artifacts/`.
- `/opt/pipeline/deploy.sh` reads an artifact from that same directory,
  extracts it to `/opt/app/`, and records the version it finds inside at
  `/workspace/deployed-version.txt`.

## Tasks

1. Run `/opt/pipeline/build.sh` and see what artifact it actually produces.
2. Run `/opt/pipeline/deploy.sh` and watch it fail to find what it's
   looking for.
3. Fix `deploy.sh` (or how it locates its artifact) so it correctly picks
   up whatever `build.sh` actually produced — not a hardcoded filename
   left over from a prior release.
4. Re-run both scripts and confirm the right version ends up in
   `/opt/app/` and `/workspace/deployed-version.txt`.

## Constraints

- Don't just rename the artifact file to satisfy deploy.sh's old
  expectation — the real fix is in how deploy.sh locates its input, so the
  pipeline keeps working the next time build.sh's naming changes again.

Get a shell: `docker exec -it breaklab-cicd-003 bash`
