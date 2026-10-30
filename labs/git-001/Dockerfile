FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    bash coreutils git \
    && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /workspace
WORKDIR /workspace

# A ready-made repo with deterministic history: three commits on main,
# a second branch (feature-login) that exists but isn't checked out,
# and a working tree left with one modified tracked file plus one
# untracked file so the learner has to tell the two apart.
RUN git init -q -b main && \
    git config user.email "ops@acme.example" && \
    git config user.name "Acme Ops" && \
    printf '# Acme Notes\n\nInternal notes repo.\n' > README.md && \
    git add README.md && \
    GIT_AUTHOR_DATE="2024-01-01T09:00:00" GIT_COMMITTER_DATE="2024-01-01T09:00:00" \
      git commit -q -m "Initial commit" && \
    printf 'TODO: write onboarding doc\n' > TODO.md && \
    git add TODO.md && \
    GIT_AUTHOR_DATE="2024-01-02T09:00:00" GIT_COMMITTER_DATE="2024-01-02T09:00:00" \
      git commit -q -m "Add TODO list" && \
    printf 'v1 release notes\n' > NOTES.md && \
    git add NOTES.md && \
    GIT_AUTHOR_DATE="2024-01-03T09:00:00" GIT_COMMITTER_DATE="2024-01-03T09:00:00" \
      git commit -q -m "Add release notes" && \
    git branch feature-login && \
    git checkout -q feature-login && \
    printf 'package main\n\nfunc Login() {}\n' > login.go && \
    git add login.go && \
    GIT_AUTHOR_DATE="2024-01-04T09:00:00" GIT_COMMITTER_DATE="2024-01-04T09:00:00" \
      git commit -q -m "Start login feature" && \
    git checkout -q main && \
    printf '# Acme Notes\n\nInternal notes repo.\n\nUpdated by on-call.\n' > README.md && \
    printf 'scratch notes, not for commit\n' > scratch.txt

WORKDIR /workspace
CMD ["sleep", "infinity"]
