FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    curl ca-certificates unzip \
    && rm -rf /var/lib/apt/lists/*
RUN arch="$(dpkg --print-architecture)"; \
    case "$arch" in \
      amd64) tf_arch="amd64" ;; \
      arm64) tf_arch="arm64" ;; \
      *) echo "unsupported arch $arch" >&2; exit 1 ;; \
    esac; \
    curl -fsSL -o /tmp/terraform.zip "https://releases.hashicorp.com/terraform/1.16.4/terraform_1.16.4_linux_${tf_arch}.zip" && \
    unzip -q /tmp/terraform.zip -d /usr/local/bin && \
    rm /tmp/terraform.zip && \
    chmod +x /usr/local/bin/terraform
RUN mkdir -p /infra /workspace
COPY main.tf /infra/main.tf
WORKDIR /infra
RUN terraform init -input=false
RUN terraform apply -auto-approve -input=false
# Someone updated the declared banner text in main.tf for today's
# release, but the pipeline that should run `terraform apply` for this
# change never ran -- the real file on disk and Terraform's own state
# both still hold yesterday's text.
COPY main.tf.updated /infra/main.tf
WORKDIR /workspace
CMD ["sleep", "infinity"]
