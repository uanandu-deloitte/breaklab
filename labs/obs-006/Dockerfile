FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    bash coreutils jq less vim-tiny \
    && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /etc/dashboards /var/lib/metrics /var/lib/dashboards /workspace
RUN printf 'title: Overview\nmetric: disk_free_pct\n' > /etc/dashboards/overview.yaml
RUN printf '{"disk_free_percent": 23, "cpu_pct": 41, "error_rate": 2}\n' > /var/lib/metrics/exported.json
COPY render-dashboard /usr/local/bin/render-dashboard
RUN chmod 755 /usr/local/bin/render-dashboard
RUN render-dashboard
WORKDIR /workspace
CMD ["sleep", "infinity"]
