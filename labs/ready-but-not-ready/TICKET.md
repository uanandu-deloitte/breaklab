# INC-308: api Service has no healthy backends

`kubectl get pods` shows the api Pod as Running, but requests through
the api Service never get a response. Find out why it isn't eligible
for traffic and get it serving again.

Get a shell: docker exec -it breaklab-ready-but-not-ready bash
