FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    bash coreutils jq less vim-tiny \
    && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /etc/dashboards /var/lib/metrics /var/lib/dashboards
RUN printf '{"http_request_count": 42}\n' > /var/lib/metrics/exported.json
RUN printf 'panel: requests\nmetric: http_request_count\n' > /etc/dashboards/overview.yaml
COPY render-dashboard /usr/local/bin/render-dashboard
RUN chmod 755 /usr/local/bin/render-dashboard
RUN /usr/local/bin/render-dashboard
CMD ["sleep", "infinity"]
