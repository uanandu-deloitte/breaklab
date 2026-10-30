FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    bash coreutils less vim-tiny \
    && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /var/lib/clock /etc/clocksync /var/log
RUN echo 0 > /var/lib/clock/offset_seconds
RUN printf '%s' 'true' > /etc/clocksync/enabled
RUN printf '%s' 'ntp.internal.example' > /etc/clocksync/server
RUN printf '2026-10-18T09:58:00Z 1760780280 request sent to host-a\n' > /var/log/host-b.log
COPY log-event /usr/local/bin/log-event
COPY sync-clock /usr/local/bin/sync-clock
RUN chmod 755 /usr/local/bin/log-event /usr/local/bin/sync-clock
RUN log-event "baseline check-in"
CMD ["sleep", "infinity"]
