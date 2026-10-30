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
# Record the real file's mtime now, while it's still the original
# resource -- this is how check.sh proves later whether the file was
# ever actually destroyed and recreated (a rename-only config change
# should never touch it, so its mtime should never move) versus
# correctly relabeled in place.
RUN stat -c %Y /var/www/legal/notice.txt > /infra/.original-mtime
# A config cleanup renamed the resource's label (legal_notice ->
# privacy_notice) without telling Terraform's state about the rename.
COPY main.tf.updated /infra/main.tf
WORKDIR /workspace
CMD ["sleep", "infinity"]
