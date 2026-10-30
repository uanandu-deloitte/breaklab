FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    bash coreutils \
    && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /srv/artifacts /workspace

# Three build artifacts, intentionally built oldest-to-newest by name AND
# by mtime (app-2.2.0.tar is both the highest version number and the most
# recently built) -- but the "current" symlink deliberately points at an
# older one, app-2.1.0.tar, simulating a pinned/rolled-back release. The
# promotion mechanism is the symlink itself, never "whichever file looks
# newest."
RUN touch -d "$(date -u -d '-20 days')" /srv/artifacts/app-2.1.0.tar && \
    touch -d "$(date -u -d '-10 days')" /srv/artifacts/app-2.1.5.tar && \
    touch -d "$(date -u -d '-2 days')"  /srv/artifacts/app-2.2.0.tar && \
    ln -s app-2.1.0.tar /srv/artifacts/current

WORKDIR /workspace
CMD ["sleep", "infinity"]
