FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    tcpdump python3 curl iproute2 procps \
    && rm -rf /var/lib/apt/lists/*
COPY start.sh /start.sh
RUN chmod +x /start.sh
RUN mkdir -p /workspace
WORKDIR /workspace
CMD ["/start.sh"]
