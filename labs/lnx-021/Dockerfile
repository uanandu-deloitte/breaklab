FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    bash coreutils findutils \
    && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /data/app /data/cache /data/logs /workspace && \
    head -c 1000000  /dev/zero > /data/app/data1.bin && \
    head -c 500000   /dev/zero > /data/app/data2.bin && \
    head -c 2000000  /dev/zero > /data/logs/app.log && \
    head -c 3000000  /dev/zero > /data/logs/app.log.1 && \
    head -c 2000000  /dev/zero > /data/cache/other.bin && \
    head -c 15000000 /dev/zero > /data/cache/big-dump.bin
WORKDIR /workspace
CMD ["sleep", "infinity"]
