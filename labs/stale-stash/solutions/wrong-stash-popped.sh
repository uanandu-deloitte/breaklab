#!/usr/bin/env bash
set -euo pipefail
cd /srv/repo
# Reference fix: find the stash that actually has the hotfix (not necessarily
# stash@{0}) and pop that specific one
ref=$(git stash list | grep "hotfix: bump timeout_seconds" | head -n1 | cut -d: -f1)
git stash pop -q "$ref"
