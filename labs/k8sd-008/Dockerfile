FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    curl ca-certificates tar \
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
RUN arch="$(dpkg --print-architecture)"; \
    curl -fsSL -o /tmp/helm.tar.gz "https://get.helm.sh/helm-v4.3.1-linux-${arch}.tar.gz" && \
    tar xzf /tmp/helm.tar.gz -C /tmp && \
    mv /tmp/linux-${arch}/helm /usr/local/bin/helm && \
    chmod +x /usr/local/bin/helm && \
    rm -rf /tmp/helm.tar.gz /tmp/linux-${arch}
ENV KUBECONFIG=/etc/rancher/k3s/k3s.yaml
RUN mkdir -p /chart /workspace
COPY chart /chart
COPY wait-for-k8s.sh /usr/local/bin/wait-for-k8s
RUN chmod 755 /usr/local/bin/wait-for-k8s
CMD ["k3s", "server", "--disable", "traefik", "--disable", "servicelb", "--disable", "metrics-server", "--snapshotter=native"]
