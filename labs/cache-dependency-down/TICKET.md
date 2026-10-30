# INC-556: Checkout is slow and sometimes empty

The frontend Pod has been Running the whole time and never restarted,
but customers report checkout is failing intermittently. Find what
frontend actually depends on and get the user journey working again.

Get a shell: docker exec -it breaklab-cache-dependency-down bash
