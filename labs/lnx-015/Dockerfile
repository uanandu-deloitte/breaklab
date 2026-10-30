FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    python3 curl procps coreutils \
    && rm -rf /var/lib/apt/lists/*
RUN useradd -r -s /usr/sbin/nologin appuser
RUN mkdir -p /srv/app && \
    printf 'service: acme\nenv: production\n' > /srv/app/config.yaml && \
    chown root:root /srv/app/config.yaml && \
    chmod 600 /srv/app/config.yaml
COPY app.py /opt/app.py
COPY start.sh /start.sh
RUN chmod +x /start.sh
CMD ["/start.sh"]
