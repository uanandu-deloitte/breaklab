# INC-414: Orders Service times out, Pod looks fine

The orders Pod is Running and Ready, but every request through the
orders Service times out. Find why the Service isn't actually routing
to it and get it serving again.

Get a shell: docker exec -it breaklab-selector-mismatch bash
