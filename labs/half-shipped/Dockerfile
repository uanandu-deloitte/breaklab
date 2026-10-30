FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    bash coreutils git less vim-tiny ca-certificates \
    && rm -rf /var/lib/apt/lists/*
RUN git config --global user.email "ops@breaklab.local" && \
    git config --global user.name "BreakLab Ops" && \
    git config --global init.defaultBranch main
RUN mkdir -p /infra/modules /etc/site
COPY site.yaml /infra/site.yaml
COPY MODULE_VERSION /infra/MODULE_VERSION
COPY modules/render-v1.sh /infra/modules/render-v1.sh
COPY modules/render-v2.sh /infra/modules/render-v2.sh
COPY infra-apply /usr/local/bin/infra-apply
RUN chmod 755 /usr/local/bin/infra-apply
RUN /usr/local/bin/infra-apply
RUN mkdir -p /srv/release-repo && cd /srv/release-repo && git init -q && \
    echo "Welcome to the store" > RELEASE_MESSAGE.md && \
    git add -A && git commit -q -m "Initial release message"
CMD ["sleep", "infinity"]
