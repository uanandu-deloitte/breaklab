FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    bash coreutils git \
    && rm -rf /var/lib/apt/lists/*

RUN git config --global user.email "dev@acme.example" && \
    git config --global user.name "Acme Dev" && \
    git config --global init.defaultBranch main

# Known-good history: two commits on main, tagged known-good-base. Then
# two more commits are made directly on main that should have gone on a
# feature branch instead. original-main-tip freezes a reference to that
# mistaken tip so it can still be identified after the fixture is solved
# and main has moved back.
RUN mkdir -p /workspace && cd /workspace && git init -q && \
    export GIT_AUTHOR_DATE="2024-03-01T09:00:00" GIT_COMMITTER_DATE="2024-03-01T09:00:00" && \
    printf 'acme-cli\n\nCommand line export tool.\n' > README.md && \
    git add README.md && git commit -q -m "Initial commit" && \
    export GIT_AUTHOR_DATE="2024-03-02T09:00:00" GIT_COMMITTER_DATE="2024-03-02T09:00:00" && \
    printf 'def export_json(rows):\n    return rows\n' > exporter.py && \
    git add exporter.py && git commit -q -m "Add exporter module" && \
    git tag known-good-base && \
    export GIT_AUTHOR_DATE="2024-03-05T14:00:00" GIT_COMMITTER_DATE="2024-03-05T14:00:00" && \
    printf 'def export_json(rows):\n    return rows\n\n\ndef export_csv(rows):\n    return ",".join(rows)\n' > exporter.py && \
    git add exporter.py && git commit -q -m "Add CSV export handler" && \
    export GIT_AUTHOR_DATE="2024-03-05T15:00:00" GIT_COMMITTER_DATE="2024-03-05T15:00:00" && \
    printf 'def test_export_csv():\n    assert export_csv(["a", "b"]) == "a,b"\n' > test_exporter.py && \
    git add test_exporter.py && git commit -q -m "Add tests for CSV export" && \
    git tag original-main-tip

WORKDIR /workspace
CMD ["sleep", "infinity"]
