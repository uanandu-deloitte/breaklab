# A merged timeline is only as trustworthy as its clocks

Reconstructing what happened across multiple hosts almost always means
merging their logs by timestamp and reading the result as a single
timeline. That only works if every host's clock is close enough to real
time that "host-a's event at T1" and "host-b's event at T2" genuinely
mean something when compared to each other. The moment one host's clock
drifts, every timestamp it produces is wrong by the same amount, and the
merged timeline silently starts lying.

## What drift actually does

A clock offset isn't an error on one event — it's a property of the
clock itself, applied uniformly to everything that clock timestamps from
that point on. If host-a's clock is running two minutes fast, every
single event host-a logs is stamped two minutes ahead of when it really
happened, consistently, until the clock is corrected. That consistency is
exactly what makes drift dangerous: it doesn't look like noise, it looks
like a believable (if surprising) sequence of events.

## Why this breaks causality in a merged view

If host-a's clock is fast, an event it logs can appear to happen
*before* the event on host-b that actually caused it — not because
anything about the real sequence of operations was wrong, but because
the two hosts disagree about what time it currently is. Sorting by
timestamp across hosts with different offsets produces an order that
has nothing to do with what actually happened first.

## Don't trust a sync tool's exit code

A clock sync job reporting success tells you it attempted a sync against
whatever it's configured to sync against — not that the sync actually
corrected anything. A disabled sync, or one pointed at a decommissioned
or wrong server, can report a clean result on every run while the clock
keeps drifting untouched underneath it.

## How to actually check

Compare each host's own clock against an independent, trusted reference
(a known-good time source, or in this lab, the real system clock) —
never by comparing hosts only against each other. Two hosts can agree
with each other and both be wrong; the only way to catch that is
checking each one against ground truth directly.
