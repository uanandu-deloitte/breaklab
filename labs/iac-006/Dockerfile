FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    coreutils \
    && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /infra /etc/app /workspace
COPY resources.yaml /infra/resources.yaml
COPY infra-apply /usr/local/bin/infra-apply
RUN chmod 755 /usr/local/bin/infra-apply
# This environment variable is baked into the image itself, so it's set
# in every shell this container ever runs, including a fresh
# `docker exec` session — there is no rc file or wrapper script hiding
# it; it's simply part of this container's own environment.
ENV INFRA_VALUE=maintenance-mode-on
WORKDIR /workspace
CMD ["sleep", "infinity"]
