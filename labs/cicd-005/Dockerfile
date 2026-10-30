FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    bash coreutils procps \
    && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /var/run/deploy /workspace

# A lock file left behind by a deploy that crashed before it could clean
# up after itself: it names a PID that is not, and never again will be,
# a running process in this container, and a STARTED timestamp computed
# relative to build time so it's always well in the past by the time
# this lab is used.
RUN STALE_TS="$(date -u -d '-30 days' +%Y-%m-%dT%H:%M:%SZ)" && \
    printf 'PID=51234\nSTARTED=%s\nHOST=ci-runner-03\n' "$STALE_TS" > /var/run/deploy/deploy.lock

WORKDIR /workspace
CMD ["sleep", "infinity"]
