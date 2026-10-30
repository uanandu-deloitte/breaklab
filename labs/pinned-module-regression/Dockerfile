FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    bash coreutils less vim-tiny \
    && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /infra/modules /etc/site
COPY site.yaml /infra/site.yaml
COPY MODULE_VERSION /infra/MODULE_VERSION
COPY modules/render-v1.sh /infra/modules/render-v1.sh
COPY modules/render-v2.sh /infra/modules/render-v2.sh
COPY infra-apply /usr/local/bin/infra-apply
RUN chmod 755 /usr/local/bin/infra-apply
RUN /usr/local/bin/infra-apply
CMD ["sleep", "infinity"]
