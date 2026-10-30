FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    nginx curl procps \
    && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /workspace /var/www/status-site && \
    printf '%s\n' \
      '<!doctype html>' \
      '<html>' \
      '<head><title>Service Status</title></head>' \
      '<body><h1>All systems operational</h1></body>' \
      '</html>' \
      > /var/www/status-site/index.html
COPY nginx-default /etc/nginx/sites-available/default
WORKDIR /workspace
CMD ["sleep", "infinity"]
