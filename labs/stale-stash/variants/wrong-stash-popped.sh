#!/usr/bin/env bash
set -euo pipefail
cd /srv/repo
# The real hotfix gets stashed first...
sed -i 's/^timeout_seconds:.*/timeout_seconds: 60/' CONFIG.md
git stash push -q -m "hotfix: bump timeout_seconds to 60"
# ...then someone stops to jot down unrelated notes and stashes that too,
# burying the real fix at stash@{1} instead of stash@{0}
echo "TODO: investigate flaky test on CI runner 3" >> NOTES.md
git stash push -q -m "wip: flaky test notes, ignore"
