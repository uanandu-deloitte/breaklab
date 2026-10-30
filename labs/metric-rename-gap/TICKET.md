# INC-985: Requests panel shows no data, not an error

The overview dashboard's requests panel has been showing "no data" since
yesterday. Traffic hasn't stopped — the exporter at
/var/lib/metrics/exported.json still shows real numbers. Get the panel
showing them again.

Get a shell: docker exec -it breaklab-metric-rename-gap bash
