FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    curl ca-certificates jq gnupg lsb-release libicu74 xz-utils \
    && rm -rf /var/lib/apt/lists/*
RUN arch="$(dpkg --print-architecture)"; \
    case "$arch" in \
      amd64) node_arch="x64" ;; \
      arm64) node_arch="arm64" ;; \
      *) echo "unsupported arch $arch" >&2; exit 1 ;; \
    esac; \
    curl -fsSL -o /tmp/node.tar.xz "https://nodejs.org/dist/v24.22.1/node-v24.22.1-linux-${node_arch}.tar.xz" && \
    tar -xJf /tmp/node.tar.xz -C /usr/local --strip-components=1 && \
    rm /tmp/node.tar.xz
RUN npm install -g azurite@3.37.0
RUN curl -fsSL https://aka.ms/InstallAzureCLIDeb -o /tmp/install-az.sh && \
    bash /tmp/install-az.sh && \
    rm /tmp/install-az.sh

ENV AZURE_STORAGE_CONNECTION_STRING="DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;BlobEndpoint=http://127.0.0.1:10000/devstoreaccount1;"

RUN mkdir -p /data /workspace
COPY notice.txt /tmp/notice.txt
# Seed one real container and blob by actually running Azurite once at
# build time -- public-access is deliberately left "off" here, the
# real misconfiguration this lab exists to fix.
RUN azurite-blob --location /data --blobHost 0.0.0.0 >/tmp/azurite-build.log 2>&1 & \
    sleep 2 && \
    az storage container create --name public-assets --public-access off >/dev/null && \
    az storage blob upload --container-name public-assets --name notice.txt --file /tmp/notice.txt --overwrite >/dev/null && \
    pkill -f azurite-blob && sleep 1

COPY entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod 755 /usr/local/bin/entrypoint.sh

WORKDIR /workspace
CMD ["/usr/local/bin/entrypoint.sh"]
