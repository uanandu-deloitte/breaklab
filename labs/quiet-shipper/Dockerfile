FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    bash coreutils jq less vim-tiny \
    && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /var/log
RUN for i in $(seq 1 10); do \
      printf '{"ts":"2026-10-18T00:00:%02dZ","level":"info","msg":"request handled","id":%d}\n' "$i" "$i" >> /var/log/app.log; \
    done
RUN printf 'DEST=/var/log/shipped.log\n' > /etc/log-shipper.conf
COPY ship-logs /usr/local/bin/ship-logs
RUN chmod 755 /usr/local/bin/ship-logs
CMD ["sleep", "infinity"]
