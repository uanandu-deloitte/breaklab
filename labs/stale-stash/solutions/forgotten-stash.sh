#!/usr/bin/env bash
set -euo pipefail
cd /srv/repo
# Reference fix: there's exactly one stash, pop it
git stash pop -q
