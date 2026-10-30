# INC-560: Today's logs never made it to the aggregator

/var/log/app.log has ten lines from today, but the aggregator's inbox
at /var/log/shipped.log is missing them. Get them shipped.

Run `ship-logs` once you think it's fixed.

Get a shell: docker exec -it breaklab-quiet-shipper bash
