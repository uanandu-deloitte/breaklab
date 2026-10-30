FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    python3 curl iproute2 procps vim-tiny less \
    && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /var/www/html /workspace && \
    echo "BreakLab docker-006 OK" > /var/www/html/index.html
COPY serve.sh /usr/local/bin/serve.sh
RUN chmod +x /usr/local/bin/serve.sh
COPY entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh
WORKDIR /workspace
CMD ["/usr/local/bin/entrypoint.sh"]
