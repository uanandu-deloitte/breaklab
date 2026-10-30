FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    curl ca-certificates xz-utils git rsync \
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
RUN npm install -g gitlab-ci-local@4.77.0
RUN mkdir -p /workspace/project
COPY .gitlab-ci.yml /workspace/project/.gitlab-ci.yml
WORKDIR /workspace/project
RUN git init -q . && \
    git config user.email "ci@example.com" && \
    git config user.name "ci" && \
    git add -A && \
    git commit -q -m "add GitLab CI pipeline"
CMD ["sleep", "infinity"]
