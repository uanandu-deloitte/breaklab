# Solution: Cache Keys Must Track What Changed

```bash
cat /usr/local/bin/cache-key
# echo "deps-v1"   <- a static string, never changed by anything

echo "# touch" >> /srv/app/requirements.lock
/usr/local/bin/cache-key    # still prints "deps-v1"

cat <<'EOF' > /usr/local/bin/cache-key
#!/usr/bin/env bash
echo "deps-$(sha256sum /srv/app/requirements.lock | cut -c1-12)"
EOF
chmod +x /usr/local/bin/cache-key

/usr/local/bin/cache-key    # now reflects the lockfile's real content

cat <<EOF > /workspace/answers.env
ORIGINAL_KEY_SOURCE=static
CACHE_KEY_NOW_TRACKS_LOCKFILE=yes
EOF
```

The original `cache-key` just echoed the literal string `"deps-v1"` —
nothing about it ever referenced `requirements.lock`, so editing the
lockfile had zero effect on the key a build system would use to decide
whether to reuse its cache. It would report "cache hit" forever,
regardless of how many times the real dependencies changed underneath
it.

Deriving the key from `sha256sum` of `requirements.lock` fixes this
properly: the hash changes if and only if the lockfile's content
changes, which is exactly the property a cache key needs. Any edit to
the lockfile now produces a different key; leaving it untouched
produces the same key every time — the cache is finally watching the
right thing.
