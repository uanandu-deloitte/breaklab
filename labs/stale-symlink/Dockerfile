FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    bash coreutils less vim-tiny \
    && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /opt/releases/v1 /opt/releases/v2 && \
    echo "v1" > /opt/releases/v1/app.txt && \
    echo "v2" > /opt/releases/v2/app.txt && \
    ln -s /opt/releases/v2 /opt/current
CMD ["sleep", "infinity"]
