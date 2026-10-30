FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    python3 bash coreutils procps less vim-tiny \
    && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /var/spool/jobs/done
COPY handler.py /usr/local/bin/handler.py
COPY metrics-collector.sh /usr/local/bin/metrics-collector.sh
COPY flood.sh /usr/local/bin/flood.sh
RUN chmod 755 /usr/local/bin/metrics-collector.sh /usr/local/bin/flood.sh
CMD ["bash", "-c", "nohup python3 /usr/local/bin/handler.py >/var/log/handler.log 2>&1 & exec bash -c 'while true; do wait -n 2>/dev/null || sleep 0.2; done'"]
