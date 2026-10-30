FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    bash coreutils tar gzip \
    && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /opt/pipeline/artifacts /opt/app /workspace
COPY build.sh /opt/pipeline/build.sh
COPY deploy.sh /opt/pipeline/deploy.sh
RUN chmod +x /opt/pipeline/build.sh /opt/pipeline/deploy.sh
WORKDIR /workspace
CMD ["sleep", "infinity"]
