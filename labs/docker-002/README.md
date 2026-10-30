# Read the Dockerfile

There's no live container to break here. Instead, `/workspace/sample.Dockerfile`
holds a multi-stage Dockerfile for a fictional Node.js app, `acme-api`. It
is never built — just read it, and answer the questions below from the
text alone.

## Questions

1. **Cache busting.** Every build of this image reinstalls every npm
   dependency from scratch, even when a commit only changed application
   source code and `package.json` never changed. Which instruction in the
   builder stage causes this, and why does it also invalidate the
   `npm ci` layer that comes after it?

2. **What survives into the final image.** The builder stage installs
   `build-essential` and `python3` with `apt-get`. The final stage is
   based on `node:20-slim`. Which files or directories from the builder
   stage actually make it into the final image? Do `build-essential` and
   `python3` come with them?

3. **Deleting a file you created.** The final stage creates two 50MB
   scratch files, "Block A" and "Block B". Block A is created and removed
   within one `RUN` instruction; Block B is created in one `RUN` and
   removed by a separate, later `RUN`. Does either approach actually keep
   that 50MB out of the final image?

## Record your answers

Write these to `/workspace/answers.env`, one `KEY=value` per line:

```
CACHE_MISS_ON_SOURCE_CHANGE=
CACHE_BUSTING_INSTRUCTION=
BUILDER_TOOLS_IN_FINAL_IMAGE=
FINAL_STAGE_COPIES_FROM_BUILDER=
SAME_RUN_DELETE_AVOIDS_BLOAT=
LATER_RUN_DELETE_AVOIDS_BLOAT=
```

- `CACHE_MISS_ON_SOURCE_CHANGE`: `yes` or `no` — does changing only app
  source code (not `package.json`) still bust the dependency-install
  layer's cache in this Dockerfile?
- `CACHE_BUSTING_INSTRUCTION`: the exact instruction line (verbatim, as it
  appears in `sample.Dockerfile`) that's responsible.
- `BUILDER_TOOLS_IN_FINAL_IMAGE`: `yes` or `no` — do `build-essential` and
  `python3` end up in the final image?
- `FINAL_STAGE_COPIES_FROM_BUILDER`: a comma-separated list of the exact
  paths the final stage copies from the builder stage (the destination
  names, e.g. `dist,node_modules`).
- `SAME_RUN_DELETE_AVOIDS_BLOAT`: `yes` or `no` — does Block A's
  same-RUN create-then-delete keep that 50MB out of the image?
- `LATER_RUN_DELETE_AVOIDS_BLOAT`: `yes` or `no` — does Block B's
  later-separate-RUN delete keep that 50MB out of the image?

`CACHE_BUSTING_INSTRUCTION`'s value contains spaces — quote it
(`KEY="value with spaces"`) so sourcing the file doesn't word-split it.

Get a shell: `docker exec -it breaklab-docker-002 bash`
