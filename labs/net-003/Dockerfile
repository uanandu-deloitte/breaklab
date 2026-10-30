FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    nginx curl iproute2 netcat-openbsd procps \
    && rm -rf /var/lib/apt/lists/*
RUN echo "BreakLab net-003 healthy endpoint" > /var/www/html/index.html
COPY nginx-8001.conf /etc/nginx/sites-available/default
COPY start.sh /start.sh
RUN chmod +x /start.sh
RUN mkdir -p /workspace
WORKDIR /workspace
CMD ["/start.sh"]
