# INC-230: Worker keeps restarting since today's deploy

`kubectl get pods` shows the worker Pod's restart count climbing every
few seconds. Find out why it keeps exiting and get it staying up.

Get a shell: docker exec -it breaklab-crash-loop-worker bash
