FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    bash coreutils man-db \
    && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /workspace /opt/fixtures && \
    printf 'acme-app release notes\nversion: 1.4.2\n' > /opt/fixtures/readme.txt
WORKDIR /workspace
CMD ["sleep", "infinity"]
