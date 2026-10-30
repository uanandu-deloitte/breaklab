# INC-2214: Disk is nearly full, but we can't find anything big

/var/log/app/service.log doesn't exist anymore — `ls -la /var/log/app`
and `du -sh /var/log/app` both come back clean. But something on this
box is still holding onto the space, and deleting more stuff doesn't
bring it back down.

Get a shell: docker exec -it breaklab-deleted-file-still-open bash
