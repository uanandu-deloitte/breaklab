# INC-3301: Deploy marker never matches what we shipped

/etc/app/release.json says the current release is 2.4.1, but
/var/lib/app/active-version.txt — written by deploy-sync on every run —
keeps disagreeing with it. Re-running deploy-sync by hand doesn't fix it.

Get a shell: docker exec -it breaklab-path-shadow bash
