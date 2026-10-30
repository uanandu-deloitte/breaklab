#!/usr/bin/env bash
set -euo pipefail
# Reference fix: add the missing hosts entry. The api resolves
# db.internal fresh on every request, so no restart is needed.
echo "127.0.0.1 db.internal" >> /etc/hosts
curl -fsS http://127.0.0.1:5000/ >/dev/null
