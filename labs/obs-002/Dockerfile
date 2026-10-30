FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    bash coreutils jq \
    && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /workspace /var/log
COPY legacy.log /var/log/legacy.log
COPY app.jsonl /var/log/app.jsonl
WORKDIR /workspace
CMD ["sleep", "infinity"]
