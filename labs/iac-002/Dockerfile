FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    jq coreutils \
    && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /infra /workspace
COPY resources.yaml /infra/resources.yaml
COPY infra-apply /usr/local/bin/infra-apply
RUN chmod 755 /usr/local/bin/infra-apply
# Simulate a resource that was already applied once under an older
# declared value, before resources.yaml was edited to its current
# content without ever re-running infra-apply. Freeze that pre-change
# value as ground truth for check.sh — written before infra-apply ever
# runs in this image.
RUN printf '%s' "Patch 2.2.0: initial release" > /infra/managed-file.txt && \
    cp /infra/managed-file.txt /infra/.original-value
WORKDIR /workspace
CMD ["sleep", "infinity"]
