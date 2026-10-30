FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    jq coreutils \
    && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /infra/resources /etc/app /etc/cache /workspace
COPY resources /infra/resources
COPY state.json /infra/state.json
COPY infra-apply /usr/local/bin/infra-apply
RUN chmod 755 /usr/local/bin/infra-apply
# /etc/cache/cache.conf predates this IaC setup entirely — it was
# created by hand, long before anyone wrote resources/cache.yaml or
# ran infra-apply for the first time, and it was never imported into
# state.json's managed list.
RUN echo "cache config v0-manual" > /etc/cache/cache.conf && \
    cp /infra/state.json /infra/.original-state.json && \
    /usr/local/bin/infra-apply
WORKDIR /workspace
CMD ["sleep", "infinity"]
