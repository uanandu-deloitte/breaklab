FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    curl ca-certificates jq \
    && rm -rf /var/lib/apt/lists/*
RUN arch="$(dpkg --print-architecture)"; \
    case "$arch" in \
      amd64) bin="k3s" ;; \
      arm64) bin="k3s-arm64" ;; \
      *) echo "unsupported arch $arch" >&2; exit 1 ;; \
    esac; \
    curl -fsSL -o /usr/local/bin/k3s "https://github.com/k3s-io/k3s/releases/download/v1.37.1%2Bk3s1/${bin}" && \
    chmod +x /usr/local/bin/k3s && \
    ln -s /usr/local/bin/k3s /usr/local/bin/kubectl
RUN mkdir -p /manifests /etc/alerts /var/lib/metrics
COPY manifests/checkout.yaml /manifests/
COPY wait-for-k8s.sh /usr/local/bin/wait-for-k8s
RUN chmod 755 /usr/local/bin/wait-for-k8s
COPY metrics-agent /usr/local/bin/metrics-agent
COPY refresh-metrics /usr/local/bin/refresh-metrics
COPY alert-evaluator /usr/local/bin/alert-evaluator
RUN chmod 755 /usr/local/bin/metrics-agent /usr/local/bin/refresh-metrics /usr/local/bin/alert-evaluator
RUN printf 'metric: checkout_healthy\nop: "<"\nthreshold: 1\n' > /etc/alerts/rules.yaml
RUN printf '{"checkout_healthy": 1, "observed_at": 1700000000}\n' > /var/lib/metrics/current.json
RUN touch /var/log/alerts.log
CMD ["k3s", "server", "--disable", "traefik", "--disable", "servicelb", "--disable", "metrics-server", "--snapshotter=native"]
