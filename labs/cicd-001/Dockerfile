FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    bash coreutils \
    && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /workspace
WORKDIR /workspace
COPY pipeline.yml /workspace/pipeline.yml
CMD ["sleep", "infinity"]
