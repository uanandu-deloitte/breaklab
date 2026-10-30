FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    bash coreutils \
    && rm -rf /var/lib/apt/lists/*
RUN groupadd ops-team && \
    mkdir -p /srv/team && \
    chgrp ops-team /srv/team && \
    chmod 2775 /srv/team
RUN mkdir -p /workspace
WORKDIR /workspace
CMD ["sleep", "infinity"]
