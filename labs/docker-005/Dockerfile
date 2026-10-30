FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    bash coreutils \
    && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /data && echo 'Welcome to the shared data volume.' > /data/welcome.txt
RUN mkdir -p /workspace
WORKDIR /workspace
CMD ["sleep", "infinity"]
