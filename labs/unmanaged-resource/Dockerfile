FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    jq less vim-tiny \
    && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /infra
COPY resources.yaml /infra/resources.yaml
COPY state.json /infra/state.json
COPY infra-apply /usr/local/bin/infra-apply
RUN chmod 755 /usr/local/bin/infra-apply
RUN /usr/local/bin/infra-apply
CMD ["sleep", "infinity"]
