FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    bash coreutils \
    && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /srv/app /opt/reference /workspace
COPY requirements.lock /srv/app/requirements.lock
COPY cache-key /usr/local/bin/cache-key
COPY cache-key /opt/reference/cache-key.orig
RUN chmod 755 /usr/local/bin/cache-key && chmod 444 /opt/reference/cache-key.orig
WORKDIR /workspace
CMD ["sleep", "infinity"]
