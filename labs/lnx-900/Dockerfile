FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    nginx curl procps iproute2 lsof python3 coreutils \
    && rm -rf /var/lib/apt/lists/*
RUN useradd -r -s /usr/sbin/nologin appuser
RUN mkdir -p /opt/app /var/log/app-data && \
    printf 'service: customer-portal\nenv: production\n' > /opt/app/config.yaml && \
    chgrp appuser /opt/app/config.yaml && \
    chmod 640 /opt/app/config.yaml && \
    printf 'RETAINED_AUDIT_RECORD\n' > /var/log/app-data/retained-audit.log && \
    chown -R appuser:appuser /var/log/app-data
COPY app.py /opt/app/app.py
COPY nginx-default /etc/nginx/sites-available/default
CMD ["sleep", "infinity"]
