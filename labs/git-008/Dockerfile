FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    bash coreutils git \
    && rm -rf /var/lib/apt/lists/*

RUN git config --global user.email "dev@acme.example" && \
    git config --global user.name "Acme Dev" && \
    git config --global init.defaultBranch main

# A bare "remote" repo at /srv/origin.git, seeded with two commits and
# pushed from a throwaway clone.
RUN git init -q --bare /srv/origin.git

RUN mkdir -p /tmp/seed && cd /tmp/seed && git init -q && \
    export GIT_AUTHOR_DATE="2024-06-01T09:00:00" GIT_COMMITTER_DATE="2024-06-01T09:00:00" && \
    printf 'acme-greeter\n' > README.md && \
    git add README.md && git commit -q -m "Initial commit" && \
    export GIT_AUTHOR_DATE="2024-06-01T10:00:00" GIT_COMMITTER_DATE="2024-06-01T10:00:00" && \
    printf 'def greet(name):\n    return f"Hello, {name}!"\n' > greet.py && \
    git add greet.py && git commit -q -m "Add greeting module" && \
    git remote add origin /srv/origin.git && \
    git push -q origin main

# The local workspace clone is made here, at the two-commit state —
# before the teammate's commit below exists on origin. --no-hardlinks
# forces a real copy of the object store instead of hardlinking it,
# which a same-filesystem local clone would otherwise try and which
# fails across a Docker build layer boundary.
RUN git clone -q --no-hardlinks /srv/origin.git /workspace

# Simulate a teammate pushing a commit the workspace clone hasn't seen.
RUN mkdir -p /tmp/teammate && git clone -q --no-hardlinks /srv/origin.git /tmp/teammate && cd /tmp/teammate && \
    export GIT_AUTHOR_NAME="Teammate" GIT_AUTHOR_EMAIL="teammate@acme.example" \
           GIT_COMMITTER_NAME="Teammate" GIT_COMMITTER_EMAIL="teammate@acme.example" && \
    export GIT_AUTHOR_DATE="2024-06-02T09:00:00" GIT_COMMITTER_DATE="2024-06-02T09:00:00" && \
    printf 'import sys\n\n\ndef log(msg):\n    print(msg, file=sys.stderr)\n' > logging_helper.py && \
    git add logging_helper.py && git commit -q -m "Add logging helper" && \
    git push -q origin main

RUN rm -rf /tmp/seed /tmp/teammate

WORKDIR /workspace
CMD ["sleep", "infinity"]
