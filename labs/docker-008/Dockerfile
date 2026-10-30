FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    nginx python3 curl procps iproute2 less vim-tiny \
    && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /srv/backend /workspace && \
    echo "BreakLab docker-008 backend OK" > /srv/backend/index.html
COPY nginx-default /etc/nginx/sites-available/default
COPY entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod 755 /usr/local/bin/entrypoint.sh
WORKDIR /workspace
CMD ["/usr/local/bin/entrypoint.sh"]
