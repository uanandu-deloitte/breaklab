FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    coreutils grep gawk \
    && rm -rf /var/lib/apt/lists/* && mkdir -p /workspace
COPY frontend.log /var/log/frontend.log
COPY api.log /var/log/api.log
COPY payments.log /var/log/payments.log
CMD ["sleep", "infinity"]
