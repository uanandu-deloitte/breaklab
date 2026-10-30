# Groups and the setgid bit

A group is just a way to let more than one user share permissions over
the same files, without making anything world-accessible. Adding users
to a shared group is only half the job, though — the directory itself
also has to be configured to keep files in that group once they're
created.

## Groups

`groupadd <name>` creates a new group. `usermod -aG <group> <user>`
appends an existing user to a group without touching their existing
group memberships (`-G` alone, without `-a`, would overwrite the whole
list — always use `-aG` when adding to one group among several).

## Why files don't just inherit the directory's group by default

Normally, a file a user creates gets that user's own **primary** group,
not the directory's group — even if the directory belongs to a shared
group. Two users with different primary groups writing into the same
shared directory would produce files owned by two different groups,
defeating the point of a shared group.

## The setgid bit

Setting the setgid bit on a directory (`chmod g+s <dir>`, which shows up
as an `s` in the group-execute position of `ls -l`) changes that: every
new file or subdirectory created inside inherits the **directory's**
group, not the creating user's primary group. Combined with making the
directory group-writable (`chmod g+w`) and owned by the shared group
(`chgrp <group> <dir>`), every member of that group can create files
there and they all land in the same group automatically.

## Proving it

The real test isn't reading the permission bits — it's creating a file
as one of the group's members and checking which group the new file
actually landed in. `su -s /bin/bash -c "<command>" <user>` runs a
single command as that user without opening an interactive session.
