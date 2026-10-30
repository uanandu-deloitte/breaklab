FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    nginx curl less vim-tiny gettext-base procps \
    && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /infra/templates
COPY site.yaml /infra/site.yaml
COPY index.html.tmpl /infra/templates/index.html.tmpl
COPY infra-apply /usr/local/bin/infra-apply
RUN chmod 755 /usr/local/bin/infra-apply
RUN /usr/local/bin/infra-apply
CMD ["sleep", "infinity"]
