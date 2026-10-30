#!/usr/bin/env bash
set -e
git daemon --export-all --base-path=/repos --reuseaddr &
exec k3s server --disable traefik --disable servicelb --disable metrics-server --snapshotter=native
