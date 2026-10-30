# INC-241: Unexpanded variable in worker config

/etc/worker/worker.env is supposed to be rendered from /infra/worker.yaml,
but something in there still looks like a template placeholder instead
of a real value. Get it rendering correctly.

Get a shell: docker exec -it breaklab-unexpanded-template bash
