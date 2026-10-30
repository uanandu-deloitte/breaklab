#!/usr/bin/env bash
# Starts the status page server. The lab's bug lives in the --bind address
# below: a server only accepting connections on its own loopback interface
# looks completely healthy from inside the container, but is unreachable
# through Docker's port mapping from anywhere else.
exec python3 -m http.server 80 --bind 127.0.0.1 --directory /var/www/html
