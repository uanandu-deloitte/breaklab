FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    nginx python3 curl procps iproute2 lsof \
    && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /srv/backend && \
    echo "BreakLab backend reached directly -- marker web003-backend-ok" > /srv/backend/index.html
COPY nginx-default /etc/nginx/sites-available/default
COPY start.sh /start.sh
RUN chmod +x /start.sh
CMD ["/start.sh"]
