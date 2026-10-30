# Solution: Where Did the Space Go?

```bash
du -sb /data/*/ | sort -rn | head -1
# e.g. 17000000  /data/cache/

find /data -type f -printf '%s %p\n' | sort -rn | head -1
# e.g. 15000000 /data/cache/big-dump.bin

cat <<EOF > /workspace/answers.env
LARGEST_DIR=/data/cache
LARGEST_FILE=/data/cache/big-dump.bin
EOF
```

`du -sb /data/*/` sums every top-level directory under `/data` in bytes;
`sort -rn` puts the largest first. `find ... -printf '%s %p\n'` prints
every file's size next to its path across the whole tree, so the single
largest file can be ranked regardless of which directory it's actually
in.
