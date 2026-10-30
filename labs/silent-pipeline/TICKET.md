# INC-412: Last push didn't ship

The on-call engineer pushed tonight's release to /srv/deploy.git from
/root/app-src and git reported success, but /srv/app still shows the
old version. Get v2.0.0 actually deployed.

Get a shell: docker exec -it breaklab-silent-pipeline bash
