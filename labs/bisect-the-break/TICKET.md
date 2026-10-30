# INC-301: healthcheck.sh is failing on main

Something in the last several commits on main broke healthcheck.sh — it
used to pass. Find the commit that broke it and revert just that
change, without rewriting history anyone else may have already pulled.

Get a shell: docker exec -it breaklab-bisect-the-break bash
