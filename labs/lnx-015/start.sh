#!/usr/bin/env bash
set -euo pipefail
exec su -s /bin/bash -c "python3 /opt/app.py" appuser
