FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    bash coreutils passwd \
    && rm -rf /var/lib/apt/lists/*
RUN groupadd opsgroup && \
    useradd -m -s /bin/bash tester
RUN mkdir -p /workspace/playground && \
    echo "owner-only secret" > /workspace/playground/owner-only.txt && \
    chown root:root /workspace/playground/owner-only.txt && \
    chmod 600 /workspace/playground/owner-only.txt && \
    echo "group readable" > /workspace/playground/group-readable.txt && \
    chown root:opsgroup /workspace/playground/group-readable.txt && \
    chmod 640 /workspace/playground/group-readable.txt && \
    echo "world readable" > /workspace/playground/world-readable.txt && \
    chown root:root /workspace/playground/world-readable.txt && \
    chmod 644 /workspace/playground/world-readable.txt
WORKDIR /workspace
CMD ["sleep", "infinity"]
