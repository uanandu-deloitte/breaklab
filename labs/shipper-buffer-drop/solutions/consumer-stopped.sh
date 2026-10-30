#!/usr/bin/env bash
set -euo pipefail
# Reference fix: drain the stuck backlog so the queue has room again
drain-queue
