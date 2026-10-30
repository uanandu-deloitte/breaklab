#!/usr/bin/env bash
set -euo pipefail
cd /srv/app-repo

# Oops: the hotfix was committed, but onto develop instead of release --
# it's a perfectly normal, reachable commit, just on the wrong branch.
git checkout -q develop
cat > src/app.sh <<'EOF'
#!/usr/bin/env bash
echo "checkout: double-submit guarded with a lock, fix applied"
EOF
git add -A
git commit -q -m "HOTFIX: guard checkout against double-submit"
git checkout -q release

# Oops: build.sh's cache key now gets computed from README.md -- a file
# that never changes -- instead of the actual application source, so it
# will keep reporting "unchanged" even after release's real source
# changes underneath it.
sed -i 's#^HASH_SOURCE=.*#HASH_SOURCE=README.md#' /opt/pipeline/build.sh

# The pipeline "ran clean" tonight, same as always.
/opt/pipeline/build.sh
