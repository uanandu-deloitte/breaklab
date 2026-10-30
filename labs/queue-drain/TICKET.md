# INC-318: Job queue backlog not clearing

Ops says the processing queue hasn't moved all afternoon. Five jobs are
sitting in /var/spool/jobs/ waiting. Get them processed.

Run `process-jobs` to drain the queue once you think it's fixed.

Get a shell: docker exec -it breaklab-queue-drain bash
