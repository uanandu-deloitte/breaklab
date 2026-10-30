FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    bash coreutils git \
    && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /workspace
WORKDIR /workspace

# A repo with only main and a single commit. original-main-tip tags
# that commit so check.sh can confirm main's tip never moves, computed
# dynamically rather than hardcoding the SHA.
RUN git init -q -b main && \
    git config user.email "ops@acme.example" && \
    git config user.name "Acme Ops" && \
    printf '# Acme Service\n\nProduction service for Acme.\n' > README.md && \
    git add README.md && \
    GIT_AUTHOR_DATE="2024-04-01T09:00:00" GIT_COMMITTER_DATE="2024-04-01T09:00:00" \
      git commit -q -m "Initial commit" && \
    git tag original-main-tip

WORKDIR /workspace
CMD ["sleep", "infinity"]
