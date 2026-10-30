FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    bash coreutils git \
    && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /workspace
WORKDIR /workspace

# One committed file (config.yaml) gets an uncommitted edit, and an
# unrelated new file (scratch-notes.txt) is left sitting untracked —
# the learner must stage and commit only the tracked edit.
RUN git init -q -b main && \
    git config user.email "ops@acme.example" && \
    git config user.name "Acme Ops" && \
    printf 'service: acme-api\ntimeout: 10\nretries: 3\n' > config.yaml && \
    printf '# Acme API\n\nService configuration lives in config.yaml.\n' > README.md && \
    git add config.yaml README.md && \
    GIT_AUTHOR_DATE="2024-02-01T09:00:00" GIT_COMMITTER_DATE="2024-02-01T09:00:00" \
      git commit -q -m "Initial commit" && \
    sed -i 's/timeout: 10/timeout: 30/' config.yaml && \
    printf 'local scratch notes - do not commit\n' > scratch-notes.txt

WORKDIR /workspace
CMD ["sleep", "infinity"]
