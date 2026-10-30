FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    bash coreutils passwd \
    && rm -rf /var/lib/apt/lists/*
RUN useradd -r -s /usr/sbin/nologin deployer
RUN mkdir -p /var/run/deploy /srv/app /usr/local/libexec && \
    chown deployer:deployer /var/run/deploy
RUN echo "v4.0.0" > /srv/app/VERSION && chown deployer:deployer /srv/app/VERSION
COPY do-deploy.sh /usr/local/libexec/do-deploy
COPY deploy.sh /usr/local/bin/deploy
RUN chmod 755 /usr/local/libexec/do-deploy /usr/local/bin/deploy
CMD ["sleep", "infinity"]
