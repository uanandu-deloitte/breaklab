FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    bash coreutils \
    && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /workspace /opt/pipeline
WORKDIR /workspace

COPY scripts/build.sh /opt/pipeline/build.sh
COPY scripts/test.sh /opt/pipeline/test.sh
COPY scripts/deploy.sh /opt/pipeline/deploy.sh
COPY scripts/run.sh /opt/pipeline/run.sh
RUN chmod +x /opt/pipeline/build.sh /opt/pipeline/test.sh /opt/pipeline/deploy.sh /opt/pipeline/run.sh

CMD ["sleep", "infinity"]
