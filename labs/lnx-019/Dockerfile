FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    bash coreutils \
    && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /opt/app /workspace && \
    printf 'service: acme\nenv: production\n' > /opt/app/config.yaml
COPY start.sh /opt/app/start.sh
RUN chmod +x /opt/app/start.sh
WORKDIR /workspace
CMD ["sleep", "infinity"]
