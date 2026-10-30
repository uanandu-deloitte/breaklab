FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    bash coreutils procps \
    && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /workspace
WORKDIR /workspace

# Deliberate shell-form CMD: a plain command string, not a JSON array, so
# Docker wraps it as `/bin/sh -c "sleep infinity"` — PID 1 inside the
# container ends up being the shell, not `sleep` itself.
CMD sleep infinity
