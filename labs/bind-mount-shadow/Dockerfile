FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    bash coreutils procps passwd less vim-tiny \
    && rm -rf /var/lib/apt/lists/*
RUN useradd -r -u 5000 -s /usr/sbin/nologin appuser
RUN mkdir -p /data && chown appuser:appuser /data
COPY writer.sh /usr/local/bin/writer.sh
RUN chmod 755 /usr/local/bin/writer.sh
CMD ["su", "-s", "/bin/bash", "-c", "/usr/local/bin/writer.sh", "appuser"]
