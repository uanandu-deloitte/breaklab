# Correlating logs across services

A single failed request in a multi-service system almost never leaves
its evidence in one place. The service that returns the error to the
user is often just the messenger — the actual cause happened one or two
hops further down the call chain, in a different log file entirely.

## The request ID is the thread

Timestamps alone aren't enough to line up related log lines: multiple
unrelated requests are being handled concurrently, so "what happened
around 10:22" matches dozens of irrelevant lines. A request ID (or trace
ID, correlation ID) that's logged consistently by every service touched
by the same request is what actually ties the story together —
`grep <request-id>` across every log file collapses the noise down to
just the lines that matter.

## Reading the chain in order

Once you've isolated a request ID, read its lines in timestamp order
across all the files, not just within one file. The order usually tells
the story directly: a downstream service logs a slow or failed call
first, the service that called it logs the resulting error a few seconds
later, and the user-facing service logs the final failure last. The
*first* ERROR in that timestamp-ordered chain — not the one closest to
the user — is the root cause.

## Don't stop at the first error you see

`grep ERROR` across all three logs will surface multiple lines for the
same incident (the proximate failure the frontend saw, and the
underlying failure that caused it). Resist reporting the first one you
find; trace it back to where the chain actually started.
