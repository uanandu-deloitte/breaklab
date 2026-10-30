FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    bash coreutils procps \
    && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /workspace
COPY worker.sh /usr/local/bin/worker.sh
COPY entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/worker.sh /usr/local/bin/entrypoint.sh
WORKDIR /workspace
CMD ["/usr/local/bin/entrypoint.sh"]
