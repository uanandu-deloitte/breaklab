FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    bash coreutils git \
    && rm -rf /var/lib/apt/lists/*

RUN git config --global user.email "dev@acme.example" && \
    git config --global user.name "Acme Dev" && \
    git config --global init.defaultBranch main

# main and feature/signup-form both gain commits after their common
# ancestor, touching different files, so merging feature into main
# cannot fast-forward: it must produce a real merge commit.
RUN mkdir -p /workspace && cd /workspace && git init -q && \
    export GIT_AUTHOR_DATE="2024-04-01T09:00:00" GIT_COMMITTER_DATE="2024-04-01T09:00:00" && \
    printf 'acme-newsletter\n\nSignup app.\n' > README.md && \
    git add README.md && git commit -q -m "Initial commit" && \
    git branch feature/signup-form && \
    export GIT_AUTHOR_DATE="2024-04-02T09:00:00" GIT_COMMITTER_DATE="2024-04-02T09:00:00" && \
    printf 'acme-newsletter\n\nSignup app.\n\nSetup instructions: pip install -r requirements.txt\n' > README.md && \
    git add README.md && git commit -q -m "Update README with setup instructions" && \
    git checkout -q feature/signup-form && \
    export GIT_AUTHOR_DATE="2024-04-02T10:00:00" GIT_COMMITTER_DATE="2024-04-02T10:00:00" && \
    printf '<form id="signup"></form>\n' > signup.html && \
    git add signup.html && git commit -q -m "Add signup form template" && \
    export GIT_AUTHOR_DATE="2024-04-02T11:00:00" GIT_COMMITTER_DATE="2024-04-02T11:00:00" && \
    printf 'function validate(form) { return true; }\n' > validate.js && \
    git add validate.js && git commit -q -m "Add signup form validation" && \
    git checkout -q main

WORKDIR /workspace
CMD ["sleep", "infinity"]
