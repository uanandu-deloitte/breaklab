FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    git less vim-tiny ca-certificates \
    && rm -rf /var/lib/apt/lists/*
RUN git config --global user.email "ops@breaklab.local" && \
    git config --global user.name "BreakLab Ops" && \
    git config --global init.defaultBranch main
RUN mkdir -p /srv/repo && cd /srv/repo && git init -q && \
    printf '#!/usr/bin/env bash\nset -euo pipefail\nTHRESHOLD=90\nusage=50\nif [ "$usage" -lt "$THRESHOLD" ]; then\n  echo ok\n  exit 0\nelse\n  echo "usage too high"\n  exit 1\nfi\n' > healthcheck.sh && \
    chmod +x healthcheck.sh && \
    echo "# App" > README.md && echo "notes" > NOTES.md && \
    git add -A && git commit -q -m "Initial commit"
CMD ["sleep", "infinity"]
