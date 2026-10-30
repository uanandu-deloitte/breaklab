FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    git less vim-tiny ca-certificates \
    && rm -rf /var/lib/apt/lists/*
RUN useradd -m -s /bin/bash deployer

RUN git init -q --bare /srv/deploy.git && chown -R deployer:deployer /srv/deploy.git
COPY post-receive /srv/deploy.git/hooks/post-receive
RUN chmod 755 /srv/deploy.git/hooks/post-receive && \
    chown deployer:deployer /srv/deploy.git/hooks/post-receive

RUN mkdir -p /etc/deploy
COPY deploy.conf /etc/deploy/deploy.conf

RUN mkdir -p /srv/releases && chown deployer:deployer /srv/releases
RUN mkdir -p /srv/releases/v1 /srv/app && echo "v1.0.0" > /srv/releases/v1/VERSION && \
    ln -s /srv/releases/v1 /srv/app/current && \
    chown -R deployer:deployer /srv/releases/v1 /srv/app

RUN su - deployer -c "git config --global user.email deployer@breaklab.local && \
    git config --global user.name 'BreakLab Deployer' && \
    git config --global init.defaultBranch main && \
    git clone -q /srv/deploy.git /home/deployer/app-src && \
    cd /home/deployer/app-src && git checkout -q -b main && \
    echo v2.0.0 > VERSION && \
    git add -A && git commit -q -m 'RELEASE: v2.0.0'"

CMD ["sleep", "infinity"]
