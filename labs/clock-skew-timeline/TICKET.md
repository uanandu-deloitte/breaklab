# INC-1002: host-a's events show up before host-b's requests

Postmortem review shows host-a logging responses to requests before
host-b's log says the request was even sent. Get host-a's clock back in
sync with reality.

Get a shell: docker exec -it breaklab-clock-skew-timeline bash
