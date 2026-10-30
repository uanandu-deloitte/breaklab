FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    curl ca-certificates jq python3 gnupg lsb-release libicu74 \
    && rm -rf /var/lib/apt/lists/*
RUN curl -fsSL https://aka.ms/InstallAzureCLIDeb -o /tmp/install-az.sh && \
    bash /tmp/install-az.sh && \
    rm /tmp/install-az.sh
RUN az bicep install
RUN mkdir -p /infra /workspace
COPY main.bicep /infra/main.bicep
WORKDIR /workspace
CMD ["sleep", "infinity"]
