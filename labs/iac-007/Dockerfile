FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    coreutils \
    && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /infra/modules /etc/app /workspace
COPY resources.yaml /infra/resources.yaml
COPY MODULE_VERSION /infra/MODULE_VERSION
COPY modules/render-v1.sh /infra/modules/render-v1.sh
COPY modules/render-v2.sh /infra/modules/render-v2.sh
COPY infra-apply /usr/local/bin/infra-apply
RUN chmod 755 /usr/local/bin/infra-apply
# Freeze the pin's original (stale) value as ground truth for check.sh,
# and apply once under that stale pin — so the resource is already
# rendered by the wrong module before the learner touches anything.
RUN cp /infra/MODULE_VERSION /infra/.original-pin && \
    /usr/local/bin/infra-apply
WORKDIR /workspace
CMD ["sleep", "infinity"]
