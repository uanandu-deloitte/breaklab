# A Dockerfile is a build recipe, not just a list of steps

Every instruction in a Dockerfile produces one layer, and three properties
of those layers explain almost every real-world image mistake:

## 1. The cache is sequential

Docker caches each layer and reuses it on the next build if that
instruction and everything before it are unchanged. The moment one layer's
cache misses — because its own inputs changed — every layer *after* it is
rebuilt too, even if that later instruction's own command never changed.
This is why instruction *order* matters as much as instruction *content*:
put the things that change often (application source) before the things
that change rarely (dependency manifests), and you force a cache miss on
every build regardless of what actually changed.

## 2. A multi-stage build only carries across what you name

`FROM ... AS builder` followed later by `COPY --from=builder <src> <dst>`
copies exactly the paths named in that instruction from the builder
stage's final filesystem state — nothing else. Anything installed into the
builder stage with `apt-get install` or similar stays behind in that
stage's own layers unless a `COPY --from` instruction explicitly reaches
back for it. This is the entire point of a slim final stage: compilers,
header files, and build tooling can be enormous, and none of it needs to
ship in production as long as the final stage only copies the build
*output*.

## 3. A layer is the filesystem *after* the whole instruction runs

Each `RUN` instruction commits one layer, and that layer's diff reflects
the filesystem only at the moment the instruction finishes — not every
intermediate state during it. A file created and deleted within the same
`RUN` instruction never appears in that layer's diff at all. But once a
layer is committed, it's committed: a later, separate `RUN` that deletes a
file only adds a "whiteout" to a newer layer that hides the file from the
merged view. The file's bytes are still sitting in the earlier layer, and
every layer that has ever existed is still part of the image — deleting
something in a different instruction than the one that created it never
shrinks the image.
