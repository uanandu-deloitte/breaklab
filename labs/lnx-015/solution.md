# Solution: Application Cannot Read Configuration

```bash
# Confirm the diagnosis: appuser can't read the file as root-owned 600.
ls -ld /srv/app /srv/app/config.yaml
su -s /bin/bash -c "cat /srv/app/config.yaml" appuser   # Permission denied

# Least-privilege fix: give appuser read access via group ownership,
# not by loosening world permissions.
chgrp appuser /srv/app/config.yaml
chmod 640 /srv/app/config.yaml

# Confirm the fix from appuser's own point of view.
su -s /bin/bash -c "cat /srv/app/config.yaml" appuser   # now succeeds

curl http://localhost:80/healthz
```

`chgrp appuser` plus `chmod 640` grants the owning group (now `appuser`)
read access while leaving "other" with no access at all — narrower than
the original 600-root-only state was ever meant to be loosened, and far
narrower than `777`. The application never needed to run as root, and
its content was never touched — only who's allowed to read it changed.
