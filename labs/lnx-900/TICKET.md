# INC-900: Customer portal unavailable

Customers report that the customer portal (http://localhost:8080) is
failing with server errors since a change window earlier today. Restore
service.

Architecture: `client -> nginx -> application -> local data`.

## Constraints

- do not replace the application image
- do not delete retained data

Get a shell: `docker exec -it breaklab-lnx-900 bash`
