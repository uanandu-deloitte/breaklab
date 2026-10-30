FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    git less vim-tiny ca-certificates \
    && rm -rf /var/lib/apt/lists/*
RUN git config --global user.email "ops@breaklab.local" && \
    git config --global user.name "BreakLab Ops" && \
    git config --global init.defaultBranch main
RUN git init -q --bare /srv/deploy.git
COPY post-receive /srv/deploy.git/hooks/post-receive
RUN chmod 755 /srv/deploy.git/hooks/post-receive
RUN mkdir -p /srv/app && echo "v1.0.0" > /srv/app/VERSION
RUN git clone -q /srv/deploy.git /root/app-src && cd /root/app-src && \
    echo "v2.0.0" > VERSION && git add -A && git commit -q -m "RELEASE: v2.0.0"
CMD ["sleep", "infinity"]
