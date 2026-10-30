FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    bash coreutils jq less vim-tiny \
    && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /etc/app /var/lib/app
RUN printf '{"service":"checkout","version":"2.4.1"}\n' > /etc/app/release.json
COPY deploy-sync /usr/local/bin/deploy-sync
RUN chmod 755 /usr/local/bin/deploy-sync
RUN deploy-sync
CMD ["sleep", "infinity"]
