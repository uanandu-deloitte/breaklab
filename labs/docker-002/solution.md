# Solution: Read the Dockerfile

```bash
cat <<'EOF' > /workspace/answers.env
CACHE_MISS_ON_SOURCE_CHANGE=yes
CACHE_BUSTING_INSTRUCTION="COPY . ."
BUILDER_TOOLS_IN_FINAL_IMAGE=no
FINAL_STAGE_COPIES_FROM_BUILDER=dist,node_modules
SAME_RUN_DELETE_AVOIDS_BLOAT=yes
LATER_RUN_DELETE_AVOIDS_BLOAT=no
EOF
```

1. `COPY . .` runs before `npm ci` in the builder stage, and it copies
   every file in the build context — including application source that
   changes on nearly every commit. Docker's build cache is sequential: as
   soon as one layer's cache misses (because the files `COPY . .` copied
   changed), every layer defined after it is rebuilt too, regardless of
   whether that later instruction's own inputs changed. So `npm ci`
   reinstalls from scratch on every build that touches any source file,
   even though `package.json` itself never changed. The fix would be to
   `COPY package.json ./` and run `npm ci` first, then `COPY . .` for the
   rest — so the dependency layer only invalidates when `package.json`
   actually changes.

2. Only `/app/dist` and `/app/node_modules` are pulled from the builder
   stage, because those are the only two paths a `COPY --from=builder`
   instruction names. `build-essential` and `python3` were installed with
   `apt-get` directly into the builder stage's own layers; nothing copies
   those layers (or the packages in them) into the final `node:20-slim`
   stage. The final image never has a compiler or Python toolchain in it.

3. Block A's create-and-delete happens inside a single `RUN` instruction,
   and a layer is a snapshot of the filesystem only after the whole
   instruction finishes — so the 50MB file never exists at the moment
   that layer's diff is captured, and it adds nothing to the image. Block
   B creates the file in one `RUN`, which commits a layer containing the
   full 50MB; the next `RUN` deletes the file, but that only adds a
   "whiteout" marker to a *later* layer hiding the file from the merged
   view — the earlier layer, and its 50MB, are still part of the image
   and still get pulled/stored. Deleting a file in a later, separate
   instruction never shrinks the image; only deleting it in the same
   instruction that created it does.
