#!/usr/bin/env bash
set -euo pipefail
DATA_DIR=/var/app/data
OUT=/var/backups/latest.tar.gz
tar czf "$OUT" -C "$DATA_DIR" .
