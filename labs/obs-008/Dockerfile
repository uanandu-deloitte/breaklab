FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    curl ca-certificates tar \
    && rm -rf /var/lib/apt/lists/*
RUN arch="$(dpkg --print-architecture)"; \
    curl -fsSL -o /tmp/prometheus.tar.gz "https://github.com/prometheus/prometheus/releases/download/v3.15.1/prometheus-3.15.1.linux-${arch}.tar.gz" && \
    tar xzf /tmp/prometheus.tar.gz -C /tmp && \
    cp /tmp/prometheus-3.15.1.linux-${arch}/promtool /usr/local/bin/promtool && \
    rm -rf /tmp/prometheus.tar.gz /tmp/prometheus-3.15.1.linux-${arch}
RUN mkdir -p /infra /workspace
COPY rules.yml /infra/rules.yml
COPY rules_test.yml /infra/rules_test.yml
WORKDIR /workspace
CMD ["sleep", "infinity"]
