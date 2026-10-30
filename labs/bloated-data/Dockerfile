FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    bash coreutils less vim-tiny \
    && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /var/app/data/tmp
RUN echo "app data" > /var/app/data/app.log
CMD ["sleep", "infinity"]
