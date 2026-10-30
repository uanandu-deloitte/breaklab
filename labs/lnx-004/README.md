# Find the Production Log

A support ticket references `transaction-id=TX-8842`, but nobody on the
ticket knows which log file recorded it. The host has several log
locations and several log files — only one of them contains this
transaction ID.

## Tasks

- Search the filesystem for the log file that contains the line with
  `transaction-id=TX-8842`.
- Write that file's absolute path to `/workspace/answer.txt` (just the
  path, nothing else).

Get a shell: `docker exec -it breaklab-lnx-004 bash`
