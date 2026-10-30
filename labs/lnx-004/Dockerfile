FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    bash coreutils findutils grep \
    && rm -rf /var/lib/apt/lists/*

# Decoy logs — plausible locations, plausible content, none contain the
# transaction ID the support ticket references.
RUN mkdir -p /var/log/app /opt/acme-app/logs /tmp/old-logs && \
    printf '%s\n' \
      '2026-10-10T09:00:01Z level=info transaction-id=TX-1001 status=ok' \
      '2026-10-10T09:00:02Z level=info transaction-id=TX-1002 status=ok' \
      '2026-10-10T09:00:03Z level=warn transaction-id=TX-1003 status=retry' \
      > /var/log/app/app.log && \
    printf '%s\n' \
      '2026-10-11T12:30:00Z checkout transaction-id=TX-2041 status=ok' \
      '2026-10-11T12:30:05Z checkout transaction-id=TX-2042 status=ok' \
      > /opt/acme-app/logs/current.log && \
    printf '%s\n' \
      '2025-06-01T00:00:00Z legacy transaction-id=TX-0099 status=ok' \
      > /tmp/old-logs/archive.log

# The real file: buried a few directories deep, under a name that
# doesn't scream "log", mixed in with other diagnostic dumps.
RUN mkdir -p /var/log/app/archive/2024-01/host-03 && \
    printf '%s\n' \
      '2024-01-14T03:11:58Z transaction-id=TX-8839 status=ok' \
      '2024-01-14T03:12:00Z transaction-id=TX-8842 status=timeout upstream=payments' \
      '2024-01-14T03:12:03Z transaction-id=TX-8845 status=ok' \
      > /var/log/app/archive/2024-01/host-03/trace.log

RUN mkdir -p /workspace
WORKDIR /workspace
CMD ["sleep", "infinity"]
