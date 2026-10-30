# Where Did the Space Go?

`/data` has grown unexpectedly. Find out which top-level directory under
it, and which single file anywhere in its tree, is responsible.

## Tasks

1. Identify the **largest top-level directory** directly under `/data`
   (by total size of everything inside it).
2. Identify the **single largest file** anywhere under `/data`, no matter
   how deep it is nested.
3. Record both as absolute paths in `/workspace/answers.env`:

   ```
   LARGEST_DIR=
   LARGEST_FILE=
   ```

Compute both from the live filesystem — don't guess.

Get a shell: `docker exec -it breaklab-lnx-021 bash`
