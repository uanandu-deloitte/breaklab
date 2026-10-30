FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    python3 curl procps coreutils iproute2 \
    && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /srv/www /var/log/app /workspace && \
    echo "ok" > /srv/www/healthz && \
    echo "fixture log data" > /var/log/app/app.log
WORKDIR /workspace
COPY worker.sh /opt/app/worker.sh
COPY start.sh /start.sh
RUN chmod +x /opt/app/worker.sh /start.sh
EXPOSE 80
CMD ["/start.sh"]
