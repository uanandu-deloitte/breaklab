#!/usr/bin/env bash
set -e
azurite-blob --location /data --blobHost 0.0.0.0 >/var/log/azurite.log 2>&1 &
exec sleep infinity
