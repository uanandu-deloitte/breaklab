FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    jq coreutils \
    && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /infra /workspace
COPY resources.yaml /infra/resources.yaml
COPY state.json /infra/state.json
COPY infra-apply /usr/local/bin/infra-apply
RUN chmod 755 /usr/local/bin/infra-apply
# Freeze the original declared value as ground truth for check.sh, before
# infra-apply ever runs.
RUN sed -n 's/^content: *//p' /infra/resources.yaml > /infra/.original-content
RUN /usr/local/bin/infra-apply
WORKDIR /workspace
CMD ["sleep", "infinity"]
