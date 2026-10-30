# INC-941: Catalog rollout stuck

A routine image update to the catalog service never finished rolling
out. `kubectl get pods` shows no Ready replicas. Get catalog serving
again.

Get a shell: docker exec -it breaklab-image-pull-backoff bash
