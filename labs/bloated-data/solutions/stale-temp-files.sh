#!/usr/bin/env bash
set -euo pipefail
# Reference fix: temp files are safe to clear; the cleanup job should have done this
rm -f /var/app/data/tmp/*
