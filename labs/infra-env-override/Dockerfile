FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    bash coreutils less vim-tiny \
    && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /infra /etc/site
COPY site.yaml /infra/site.yaml
COPY infra-apply /usr/local/bin/infra-apply
COPY apply-site /usr/local/bin/apply-site
RUN chmod 755 /usr/local/bin/infra-apply /usr/local/bin/apply-site
# Legacy cruft: someone added this during a one-off debugging session a
# long time ago and never removed it.
RUN echo 'export SITE_MESSAGE="Legacy notice v1"' >> /etc/bash.bashrc
RUN /usr/local/bin/infra-apply
CMD ["sleep", "infinity"]
