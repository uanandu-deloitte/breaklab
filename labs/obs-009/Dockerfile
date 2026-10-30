FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    curl ca-certificates tar python3 jq \
    && rm -rf /var/lib/apt/lists/*
RUN arch="$(dpkg --print-architecture)"; \
    curl -fsSL -o /tmp/prometheus.tar.gz "https://github.com/prometheus/prometheus/releases/download/v3.15.1/prometheus-3.15.1.linux-${arch}.tar.gz" && \
    tar xzf /tmp/prometheus.tar.gz -C /tmp && \
    cp /tmp/prometheus-3.15.1.linux-${arch}/prometheus /usr/local/bin/prometheus && \
    cp /tmp/prometheus-3.15.1.linux-${arch}/promtool /usr/local/bin/promtool && \
    rm -rf /tmp/prometheus.tar.gz /tmp/prometheus-3.15.1.linux-${arch}
RUN arch="$(dpkg --print-architecture)"; \
    curl -fsSL -o /tmp/grafana.tar.gz "https://dl.grafana.com/oss/release/grafana-13.2.4.linux-${arch}.tar.gz" && \
    tar xzf /tmp/grafana.tar.gz -C /opt && \
    mv /opt/grafana-* /opt/grafana && \
    rm /tmp/grafana.tar.gz

RUN mkdir -p /infra /opt /etc/grafana/provisioning/datasources /etc/grafana/provisioning/dashboards/json /var/lib/prometheus
COPY prometheus.yml /infra/prometheus.yml
COPY metrics_server.py /opt/metrics_server.py
COPY grafana.ini /etc/grafana/grafana.ini
COPY provisioning/datasources/prometheus.yaml /etc/grafana/provisioning/datasources/prometheus.yaml
COPY provisioning/dashboards/dashboards.yaml /etc/grafana/provisioning/dashboards/dashboards.yaml
COPY provisioning/dashboards/json/checkout.json /etc/grafana/provisioning/dashboards/json/checkout.json
COPY entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod 755 /usr/local/bin/entrypoint.sh

WORKDIR /workspace
CMD ["/usr/local/bin/entrypoint.sh"]
