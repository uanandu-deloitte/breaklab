FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    bash coreutils tar less vim-tiny \
    && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /srv/artifacts /srv/app /etc/deploy /tmp/build-prod /tmp/build-staging
RUN echo "prod build 42: checkout-v3 (stable)" > /tmp/build-prod/CONTENT && \
    echo "staging build 42: checkout-v4-experimental (unreleased)" > /tmp/build-staging/CONTENT && \
    tar -cf /srv/artifacts/app-prod-42.tar -C /tmp/build-prod CONTENT && \
    tar -cf /srv/artifacts/app-staging-42.tar -C /tmp/build-staging CONTENT && \
    echo "app-prod-42.tar" > /srv/artifacts/PROD_LATEST && \
    echo "PROD_ARTIFACT=app-prod-42.tar" > /etc/deploy/promote.conf
COPY promote /usr/local/bin/promote
RUN chmod 755 /usr/local/bin/promote
RUN /usr/local/bin/promote
CMD ["sleep", "infinity"]
