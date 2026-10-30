FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    git less vim-tiny ca-certificates \
    && rm -rf /var/lib/apt/lists/*
RUN git config --global user.email "ops@breaklab.local" && \
    git config --global user.name "BreakLab Ops" && \
    git config --global init.defaultBranch main

# --- Situation A: a staged change that was never meant to be committed ---
RUN mkdir -p /workspace/repo-staged && cd /workspace/repo-staged && git init -q && \
    printf 'mode: safe\ntimeout: 30\n' > config.yaml && \
    git add -A && \
    GIT_AUTHOR_DATE="2024-01-01T09:00:00" GIT_COMMITTER_DATE="2024-01-01T09:00:00" \
      git commit -q -m "Initial config" && \
    printf 'mode: debug\ntimeout: 30\n' > config.yaml && \
    git add config.yaml

# --- Situation B: a bad commit that only exists locally (no remote) ---
RUN mkdir -p /workspace/repo-private && cd /workspace/repo-private && git init -q && \
    printf 'stage: production\n' > settings.yaml && \
    git add -A && \
    GIT_AUTHOR_DATE="2024-01-01T09:00:00" GIT_COMMITTER_DATE="2024-01-01T09:00:00" \
      git commit -q -m "Initial settings" && \
    git rev-parse HEAD > /srv/.repo-private-baseline-sha && \
    printf 'stage: production\n# BAD-COMMIT-MARKER: disabled auth check for debugging\n' > settings.yaml && \
    GIT_AUTHOR_DATE="2024-01-02T09:00:00" GIT_COMMITTER_DATE="2024-01-02T09:00:00" \
      git commit -aq -m "DEBUG: temporarily disable auth check"

# --- Situation C: the exact same mistake, but already on a shared/pushed branch ---
RUN git init -q --bare /srv/repo-shared-origin.git && \
    git clone -q --no-hardlinks /srv/repo-shared-origin.git /tmp/seed-shared && \
    cd /tmp/seed-shared && \
    printf 'stage: production\n' > settings.yaml && \
    git add -A && \
    GIT_AUTHOR_DATE="2024-01-01T09:00:00" GIT_COMMITTER_DATE="2024-01-01T09:00:00" \
      git commit -q -m "Initial settings" && \
    git push -q origin main && \
    printf 'stage: production\n# BAD-COMMIT-MARKER: disabled auth check for debugging\n' > settings.yaml && \
    GIT_AUTHOR_DATE="2024-01-02T09:00:00" GIT_COMMITTER_DATE="2024-01-02T09:00:00" \
      git commit -aq -m "DEBUG: temporarily disable auth check" && \
    git push -q origin main && \
    git rev-parse HEAD > /srv/.repo-shared-bad-sha

RUN git clone -q --no-hardlinks /srv/repo-shared-origin.git /workspace/repo-shared

WORKDIR /workspace
CMD ["sleep", "infinity"]
