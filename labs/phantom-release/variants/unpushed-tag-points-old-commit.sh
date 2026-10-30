#!/usr/bin/env bash
set -euo pipefail
cd /srv/app-repo

# The hotfix lands directly on release as a normal, later commit...
git checkout -q release
cat > src/app.sh <<'EOF'
#!/usr/bin/env bash
echo "checkout: double-submit guarded with a lock, fix applied"
EOF
git add -A
git commit -q -m "HOTFIX: guard checkout against double-submit"

# ...but the v2.0.0 tag was cut earlier tonight (it already points at the
# initial commit from image build) and was never moved to the fix.

# The pipeline builds and ships from tags, not branch heads -- a real
# pipeline convention, not a bug by itself.
echo "v2.0.0" > /opt/pipeline/BUILD_REF

# The pipeline "ran clean" tonight: it built exactly what v2.0.0 pointed
# to, which is the old, pre-hotfix source.
/opt/pipeline/build.sh
