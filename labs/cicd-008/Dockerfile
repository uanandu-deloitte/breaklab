FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    curl ca-certificates tar git \
    && rm -rf /var/lib/apt/lists/*
RUN arch="$(dpkg --print-architecture)"; \
    curl -fsSL -o /tmp/act.tar.gz "https://github.com/nektos/act/releases/download/v0.2.89/act_Linux_${arch}.tar.gz" && \
    tar xzf /tmp/act.tar.gz -C /usr/local/bin act && \
    rm /tmp/act.tar.gz && \
    chmod +x /usr/local/bin/act
RUN mkdir -p /workspace/project/.github/workflows
COPY ci.yml /workspace/project/.github/workflows/ci.yml
WORKDIR /workspace/project
RUN git init -q . && \
    git config user.email "ci@example.com" && \
    git config user.name "ci" && \
    git add -A && \
    git commit -q -m "add CI workflow"
CMD ["sleep", "infinity"]
