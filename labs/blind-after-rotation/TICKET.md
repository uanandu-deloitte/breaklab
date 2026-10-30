# INC-560: Metrics frozen since last night's log rotation

The dashboard still shows yesterday's request count. The app is
definitely still logging. Get the metrics agent reporting current
numbers again.

Run `collect-metrics` once you think it's fixed.

Get a shell: docker exec -it breaklab-blind-after-rotation bash
