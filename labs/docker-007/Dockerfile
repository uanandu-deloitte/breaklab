FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    bash coreutils procps less \
    && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /workspace /var/log/app
COPY app.sh /usr/local/bin/app.sh
RUN chmod 755 /usr/local/bin/app.sh
WORKDIR /workspace
CMD ["/usr/local/bin/app.sh"]
