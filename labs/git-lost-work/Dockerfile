FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    git less vim-tiny ca-certificates \
    && rm -rf /var/lib/apt/lists/*
RUN git config --global user.email "ops@breaklab.local" && \
    git config --global user.name "BreakLab Ops" && \
    git config --global init.defaultBranch main
RUN mkdir -p /srv/repo && cd /srv/repo && git init -q && \
    echo "# App" > README.md && \
    echo "v1.1.0" > RELEASE.md && \
    git add -A && git commit -q -m "Initial release v1.1.0"
CMD ["sleep", "infinity"]
