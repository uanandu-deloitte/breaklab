FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    nginx python3 curl procps iproute2 lsof vim-tiny less \
    && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /srv/backend && \
    echo "Storefront v2 -- all systems normal" > /srv/backend/index.html
COPY nginx-default /etc/nginx/sites-available/default
CMD ["sleep", "infinity"]
