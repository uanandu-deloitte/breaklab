FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    bash coreutils less vim-tiny \
    && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /srv/app /var/cache/ci/deps
RUN printf 'libwidget==2.0\n' > /srv/app/requirements.lock && \
    printf 'libwidget==2.0\n' > /var/cache/ci/deps/requirements.lock
COPY ci-build /usr/local/bin/ci-build
RUN chmod 755 /usr/local/bin/ci-build
RUN /usr/local/bin/ci-build
CMD ["sleep", "infinity"]
