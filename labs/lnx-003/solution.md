# Solution: Explore the Filesystem

```bash
cd /etc && pwd >> /workspace/transcript.txt
cd /var/log && pwd >> /workspace/transcript.txt
cd /opt && pwd >> /workspace/transcript.txt
cd ~ && pwd >> /workspace/transcript.txt
cd /workspace
```

`pwd` after each `cd` records exactly where the shell actually ended up,
rather than trusting the path you typed was spelled correctly.
