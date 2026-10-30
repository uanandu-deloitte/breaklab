FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    git git-filter-repo less vim-tiny ca-certificates \
    && rm -rf /var/lib/apt/lists/*
RUN git config --global user.email "ops@breaklab.local" && \
    git config --global user.name "BreakLab Ops" && \
    git config --global init.defaultBranch main
RUN mkdir -p /srv/repo && cd /srv/repo && git init -q && \
    echo "# App" > README.md && echo "notes" > NOTES.md && \
    git add -A && git commit -q -m "Initial commit"
CMD ["sleep", "infinity"]
