FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    nginx python3 curl procps iproute2 lsof \
    && rm -rf /var/lib/apt/lists/*
RUN echo "BreakLab net-002 OK" > /var/www/html/index.html
COPY start.sh /start.sh
RUN chmod +x /start.sh
RUN mkdir -p /workspace
WORKDIR /workspace
CMD ["/start.sh"]
