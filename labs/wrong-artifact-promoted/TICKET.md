# INC-614: Prod users are seeing a staging checkout flow

Build 43 shipped tonight. The pipeline log says "promoted prod build"
and nobody saw an error. Support is now getting reports of an
unreleased checkout flow in production — the one that was only ever
supposed to be on staging. Find out which artifact is actually live
under /srv/app/live and get the real prod build promoted.

Get a shell: docker exec -it breaklab-wrong-artifact-promoted bash
