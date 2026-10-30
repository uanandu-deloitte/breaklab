FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    python3 curl coreutils \
    && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /workspace /opt/app
COPY app.py /opt/app/app.py
COPY start.sh /start.sh
RUN chmod +x /start.sh
WORKDIR /workspace
CMD ["/start.sh"]
