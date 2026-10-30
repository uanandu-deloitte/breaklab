FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    bash coreutils less vim-tiny \
    && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /var/spool/shipper/queue /var/lib/shipper /var/log
RUN for i in $(seq 1 5); do echo "baseline line $i" >> /var/log/app.log; done
RUN touch /var/log/shipped-final.log
COPY ship-logs /usr/local/bin/ship-logs
COPY drain-queue /usr/local/bin/drain-queue
RUN chmod 755 /usr/local/bin/ship-logs /usr/local/bin/drain-queue
RUN ship-logs && drain-queue
CMD ["sleep", "infinity"]
