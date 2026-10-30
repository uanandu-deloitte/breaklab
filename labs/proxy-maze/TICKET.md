# INC-340: Storefront proxy is misbehaving

Since a cleanup pass touched the edge nginx config, the storefront is either
unreachable or returning errors through its public port. Get it serving
normally again.

Get a shell: docker exec -it breaklab-proxy-maze bash
