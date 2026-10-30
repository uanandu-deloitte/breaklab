FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    bash coreutils diffutils \
    && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /workspace && \
    printf '%s\n' \
      'database_host: db-old.internal' \
      'log_level: debug' \
      'timeout: 30' \
      'max_connections: 100' \
      'retry_count: 3' \
      > /workspace/config.old && \
    printf '%s\n' \
      'database_host: db-new.internal' \
      'log_level: info' \
      'timeout: 60' \
      'max_connections: 100' \
      'retry_count: 3' \
      > /workspace/config.new
WORKDIR /workspace
CMD ["sleep", "infinity"]
