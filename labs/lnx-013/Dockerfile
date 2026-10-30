FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    bash coreutils passwd \
    && rm -rf /var/lib/apt/lists/*
RUN useradd -m -s /bin/bash developer && \
    useradd -r -s /usr/sbin/nologin -g operator operator
RUN mkdir -p /srv/team && chown root:root /srv/team && chmod 755 /srv/team
RUN mkdir -p /workspace
WORKDIR /workspace
CMD ["sleep", "infinity"]
