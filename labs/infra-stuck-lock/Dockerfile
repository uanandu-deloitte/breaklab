FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    nginx curl less vim-tiny gettext-base procps passwd \
    && rm -rf /var/lib/apt/lists/*
RUN useradd -r -s /usr/sbin/nologin infraops
RUN mkdir -p /infra/templates /var/run/infra && \
    chown infraops:infraops /var/run/infra && chmod 755 /var/run/infra
COPY site.yaml /infra/site.yaml
COPY index.html.tmpl /infra/templates/index.html.tmpl
COPY infra-apply /usr/local/bin/infra-apply
RUN chmod 755 /usr/local/bin/infra-apply
RUN /usr/local/bin/infra-apply
CMD ["sleep", "infinity"]
