FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    bash coreutils less vim-tiny gettext-base \
    && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /infra/templates
COPY worker.yaml /infra/worker.yaml
COPY worker.env.tmpl /infra/templates/worker.env.tmpl
COPY infra-apply-worker /usr/local/bin/infra-apply-worker
RUN chmod 755 /usr/local/bin/infra-apply-worker
RUN /usr/local/bin/infra-apply-worker
CMD ["sleep", "infinity"]
