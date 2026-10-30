FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    jq coreutils less vim-tiny \
    && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /infra/vars /workspace
COPY resources.yaml.tmpl /infra/resources.yaml.tmpl
COPY vars/dev.env /infra/vars/dev.env
COPY vars/prod.env /infra/vars/prod.env
COPY infra-apply /usr/local/bin/infra-apply
RUN chmod 755 /usr/local/bin/infra-apply
RUN /usr/local/bin/infra-apply dev
WORKDIR /workspace
CMD ["sleep", "infinity"]
