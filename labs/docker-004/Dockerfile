FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    bash coreutils procps \
    && rm -rf /var/lib/apt/lists/*
ENV APP_ENV=production
ENV LOG_LEVEL=info
ENV APP_VERSION=1.4.2
RUN mkdir -p /workspace
WORKDIR /workspace
CMD ["sleep", "infinity"]
