FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    nginx curl procps iproute2 lsof python3 coreutils \
    && rm -rf /var/lib/apt/lists/*
RUN useradd -r -s /usr/sbin/nologin apiuser
RUN mkdir -p /opt/api /opt/db && \
    printf 'service: public-api\nenv: production\ndb_host: db.internal\ndb_port: 6000\n' > /opt/api/config.yaml && \
    chgrp apiuser /opt/api/config.yaml && \
    chmod 640 /opt/api/config.yaml
COPY api.py /opt/api/api.py
COPY db.py /opt/db/db.py
COPY nginx-default /etc/nginx/sites-available/default
CMD ["sleep", "infinity"]
