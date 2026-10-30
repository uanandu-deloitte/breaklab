FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    bash coreutils less vim-tiny \
    && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /var/lib/clock /var/log /workspace
RUN echo 125 > /var/lib/clock/host-a-offset_seconds
RUN echo 0 > /var/lib/clock/host-b-offset_seconds
COPY log-event-a /usr/local/bin/log-event-a
COPY log-event-b /usr/local/bin/log-event-b
RUN chmod 755 /usr/local/bin/log-event-a /usr/local/bin/log-event-b
RUN log-event-a "startup" && log-event-b "startup"
WORKDIR /workspace
CMD ["sleep", "infinity"]
