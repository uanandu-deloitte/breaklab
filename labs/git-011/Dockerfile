FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    git less vim-tiny ca-certificates \
    && rm -rf /var/lib/apt/lists/*
RUN git config --global user.email "ops@breaklab.local" && \
    git config --global user.name "BreakLab Ops" && \
    git config --global init.defaultBranch main

RUN mkdir -p /workspace && cd /workspace && git init -q && \
    printf 'def main():\n    print("service v1")\n' > app.py && \
    git add -A && \
    GIT_AUTHOR_DATE="2024-01-01T09:00:00" GIT_COMMITTER_DATE="2024-01-01T09:00:00" \
      git commit -q -m "Initial version" && \
    printf 'def main():\n    print("service v1")\n    # TRACKED-CHANGE-MARKER: added debug line\n' > app.py && \
    printf 'UNTRACKED-FILE-MARKER: scratch notes, not committed\n' > notes.txt

WORKDIR /workspace
CMD ["sleep", "infinity"]
