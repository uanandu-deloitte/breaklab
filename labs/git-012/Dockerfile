FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    git less vim-tiny ca-certificates \
    && rm -rf /var/lib/apt/lists/*
RUN git config --global user.email "ops@breaklab.local" && \
    git config --global user.name "BreakLab Ops" && \
    git config --global init.defaultBranch main

# Shared ancestor: handler.py (with a latent null-pointer bug) plus a
# README, tagged as the v1.0.0 release and branched off as "release".
RUN mkdir -p /workspace && cd /workspace && git init -q && \
    printf 'def handle(request):\n    return request.body.decode()\n' > handler.py && \
    printf '# App\n\nHandles incoming requests.\n' > README.md && \
    git add -A && \
    GIT_AUTHOR_DATE="2024-01-01T09:00:00" GIT_COMMITTER_DATE="2024-01-01T09:00:00" \
      git commit -q -m "Initial handler" && \
    GIT_COMMITTER_DATE="2024-01-05T09:00:00" git tag -a v1.0.0 -m "Release 1.0.0" && \
    git branch release

# main moves ahead with three commits: an unrelated feature, the hotfix
# that actually needs to ship, and an unrelated changelog update.
RUN cd /workspace && \
    printf 'RATE-LIMIT-ONLY-MARKER: 100 req/min per client\n' > ratelimit.py && \
    git add -A && \
    GIT_AUTHOR_DATE="2024-01-10T09:00:00" GIT_COMMITTER_DATE="2024-01-10T09:00:00" \
      git commit -q -m "Add rate limiting to API"

RUN cd /workspace && \
    printf 'def handle(request):\n    # HOTFIX-NULLCHECK-MARKER: guard against a null request\n    if request is None:\n        return None\n    return request.body.decode()\n' > handler.py && \
    GIT_AUTHOR_DATE="2024-01-11T09:00:00" GIT_COMMITTER_DATE="2024-01-11T09:00:00" \
      git commit -aq -m "Fix null-pointer crash in handler.py"

RUN cd /workspace && \
    printf '# Changelog\n\nCHANGELOG-ONLY-MARKER: notes for the v1.1 docs update\n' > CHANGELOG.md && \
    git add -A && \
    GIT_AUTHOR_DATE="2024-01-12T09:00:00" GIT_COMMITTER_DATE="2024-01-12T09:00:00" \
      git commit -q -m "Update changelog for v1.1"

WORKDIR /workspace
CMD ["sleep", "infinity"]
