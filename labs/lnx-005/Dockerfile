FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    bash coreutils \
    && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /workspace /opt/fixtures && \
    printf '#!/usr/bin/env bash\necho "acme-app running"\n' > /opt/fixtures/acme-app && \
    chmod +x /opt/fixtures/acme-app
WORKDIR /workspace
CMD ["sleep", "infinity"]
