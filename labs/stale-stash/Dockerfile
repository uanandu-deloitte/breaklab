FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    git less vim-tiny ca-certificates \
    && rm -rf /var/lib/apt/lists/*
RUN git config --global user.email "ops@breaklab.local" && \
    git config --global user.name "BreakLab Ops" && \
    git config --global init.defaultBranch main
RUN mkdir -p /srv/repo && cd /srv/repo && git init -q && \
    printf '# App Config\ntimeout_seconds: 30\n' > CONFIG.md && \
    printf '# Notes\nnothing to see here\n' > NOTES.md && \
    git add -A && git commit -q -m "Initial config"
CMD ["sleep", "infinity"]
