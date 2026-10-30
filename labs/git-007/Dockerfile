FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    bash coreutils git \
    && rm -rf /var/lib/apt/lists/*

RUN git config --global user.email "dev@acme.example" && \
    git config --global user.name "Acme Dev" && \
    git config --global init.defaultBranch main

# main and feature/raise-timeout both edit the same "timeout" line of
# config.yaml differently since their common ancestor, so merging one
# into the other is guaranteed to conflict on that line.
RUN mkdir -p /workspace && cd /workspace && git init -q && \
    export GIT_AUTHOR_DATE="2024-05-01T09:00:00" GIT_COMMITTER_DATE="2024-05-01T09:00:00" && \
    printf 'service: acme-api\nport: 8080\ntimeout: 30\n' > config.yaml && \
    git add config.yaml && git commit -q -m "Initial commit" && \
    git branch feature/raise-timeout && \
    export GIT_AUTHOR_DATE="2024-05-02T09:00:00" GIT_COMMITTER_DATE="2024-05-02T09:00:00" && \
    printf 'service: acme-api\nport: 8080\ntimeout: 45\n' > config.yaml && \
    git add config.yaml && git commit -q -m "Reduce timeout to fail fast" && \
    git checkout -q feature/raise-timeout && \
    export GIT_AUTHOR_DATE="2024-05-02T10:00:00" GIT_COMMITTER_DATE="2024-05-02T10:00:00" && \
    printf 'service: acme-api\nport: 8080\ntimeout: 60\n' > config.yaml && \
    git add config.yaml && git commit -q -m "Increase timeout for slow upstream calls" && \
    git checkout -q main

WORKDIR /workspace
CMD ["sleep", "infinity"]
