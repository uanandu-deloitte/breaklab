FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    bash coreutils git \
    && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /workspace
WORKDIR /workspace

# fixture-base tags the last clean commit. On top of it, two unrelated
# changes sit uncommitted at once, across two files: a real off-by-one
# bugfix in app.py, and an unrelated stray debug print left in
# utils.py. Both are already applied to the working tree — the
# learner's job is to split them into two separate commits, not to
# write the changes themselves.
RUN git init -q -b main && \
    git config user.email "ops@acme.example" && \
    git config user.name "Acme Ops" && \
    printf 'def get_last_item(items):\n    return items[len(items)]\n\ndef main():\n    sample = [1, 2, 3]\n    print(get_last_item(sample))\n\nif __name__ == "__main__":\n    main()\n' > app.py && \
    printf 'def add(a, b):\n    return a + b\n' > utils.py && \
    git add app.py utils.py && \
    GIT_AUTHOR_DATE="2024-03-01T09:00:00" GIT_COMMITTER_DATE="2024-03-01T09:00:00" \
      git commit -q -m "Initial commit" && \
    git tag fixture-base && \
    printf 'def get_last_item(items):\n    return items[len(items) - 1]\n\ndef main():\n    sample = [1, 2, 3]\n    print(get_last_item(sample))\n\nif __name__ == "__main__":\n    main()\n' > app.py && \
    printf 'print("DEBUG: utils.py loaded")\n\ndef add(a, b):\n    return a + b\n' > utils.py

WORKDIR /workspace
CMD ["sleep", "infinity"]
