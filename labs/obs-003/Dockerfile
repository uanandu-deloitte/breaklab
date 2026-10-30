FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    bash coreutils jq less vim-tiny \
    && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /etc/alerts /var/lib/metrics
RUN printf '{"cpu_pct": 42, "disk_used_pct": 55, "error_rate": 7}\n' > /var/lib/metrics/current.json
RUN touch /etc/alerts/rules.yaml
RUN touch /var/log/alerts.log
COPY alert-evaluator /usr/local/bin/alert-evaluator
RUN chmod 755 /usr/local/bin/alert-evaluator
CMD ["sleep", "infinity"]
