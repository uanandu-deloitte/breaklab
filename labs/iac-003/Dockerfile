FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    jq coreutils less vim-tiny \
    && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /infra /workspace
COPY resources.yaml /infra/resources.yaml
COPY infra-apply /usr/local/bin/infra-apply
RUN chmod 755 /usr/local/bin/infra-apply
RUN /usr/local/bin/infra-apply
RUN sed -n 's/^value: *//p' /infra/resources.yaml > /infra/.original-value
WORKDIR /workspace
CMD ["sleep", "infinity"]
