# A cache is only as trustworthy as its key

A build cache exists to skip redoing expensive work — installing
dependencies, compiling, whatever — when nothing that affects the
result has changed since last time. The entire mechanism rests on one
promise: the cache key changes exactly when the thing it represents
changes, and stays the same otherwise. Break that promise in either
direction and the cache stops doing its job, usually without raising
any error at all.

## The dangerous direction is the quiet one

A key that changes too often just means wasted rebuilds — annoying,
but safe; you notice it because builds take longer than they should.
The dangerous failure is a key that doesn't change often enough: it
keeps reporting "cache hit, nothing to do" even after the real input
has changed. The build reports success. Nothing looks wrong. The
artifact it produces is just quietly built from yesterday's
dependencies instead of today's.

This happens whenever the key is derived from something *incidental*
to the real input rather than the input itself — a static string that
never changes at all, a build counter that increments for unrelated
reasons, a directory's own timestamp rather than the content inside
it. Each of these changes (or doesn't change) according to its own
rules, which have nothing to do with whether the dependency lockfile
actually changed.

## The fix is always the same shape

Derive the key from the actual content of the thing that should
invalidate the cache — typically a hash of the lockfile itself
(`sha256sum requirements.lock`, or similar). A hash changes if and only
if the file's content changes, which is exactly the property a cache
key needs. It doesn't matter whether the lockfile gained a line, lost
a line, or just bumped one version number — any content change
produces a different hash, and any non-change produces the same one.

## Why this matters

"The pipeline says build complete" is not evidence that the right
dependencies shipped — it's only evidence the script ran without
throwing an error. A cache keyed on the wrong thing can report success
forever while serving content that stopped matching the real input
days or weeks ago. The only way to trust a cache is to know, concretely,
what its key is actually computed from.
