FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    git less vim-tiny ca-certificates \
    && rm -rf /var/lib/apt/lists/*
RUN git config --global user.email "ops@breaklab.local" && \
    git config --global user.name "BreakLab Ops" && \
    git config --global init.defaultBranch main

# 1. Create the bare "origin" repo and seed it with an initial commit.
RUN git init -q --bare /srv/origin.git

RUN git clone -q --no-hardlinks /srv/origin.git /tmp/seed && \
    cd /tmp/seed && \
    printf 'def main():\n    print("service v1")\n' > app.py && \
    printf '# Changelog\n\n## v1.0.0\n- initial release\n' > CHANGELOG.md && \
    git add -A && \
    GIT_AUTHOR_DATE="2024-01-01T09:00:00" GIT_COMMITTER_DATE="2024-01-01T09:00:00" \
      git commit -q -m "Initial release" && \
    git push -q origin main

# 2. Clone to /workspace -- this is the clone the user inherits, made
#    *before* the teammate's push below.
RUN git clone -q --no-hardlinks /srv/origin.git /workspace

# 3. A teammate clones separately, commits, and pushes to origin -- after
#    /workspace's clone was made, so /workspace doesn't know about it yet.
RUN git clone -q --no-hardlinks /srv/origin.git /tmp/teammate && \
    cd /tmp/teammate && \
    printf '## v1.1.0\n- TEAMMATE-COMMIT-MARKER: fix login timeout\n' >> CHANGELOG.md && \
    GIT_AUTHOR_DATE="2024-01-02T10:00:00" GIT_COMMITTER_DATE="2024-01-02T10:00:00" \
      git commit -aq -m "Fix login timeout" && \
    git push -q origin main && \
    git rev-parse HEAD > /srv/.teammate-sha

# 4. Back in /workspace (which still only knows about the initial
#    commit), make a local commit of our own -- and never push it.
RUN cd /workspace && \
    printf 'def main():\n    print("service v1")\n\n# LOCAL-COMMIT-MARKER: added retry logic\n' > app.py && \
    GIT_AUTHOR_DATE="2024-01-02T11:00:00" GIT_COMMITTER_DATE="2024-01-02T11:00:00" \
      git commit -aq -m "Add retry logic to service calls" && \
    git rev-parse HEAD > /srv/.local-sha

WORKDIR /workspace
CMD ["sleep", "infinity"]
