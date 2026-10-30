FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    bash coreutils jq \
    && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /workspace /var/log /var/lib/metrics
COPY app.jsonl /var/log/app.jsonl
COPY current.json /var/lib/metrics/current.json
WORKDIR /workspace
CMD ["sleep", "infinity"]
