FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    coreutils procps \
    && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /infra /workspace /var/run/infra
COPY resources.yaml /infra/resources.yaml
COPY infra-apply /usr/local/bin/infra-apply
RUN chmod 755 /usr/local/bin/infra-apply
# Simulate a run that crashed mid-apply: it wrote its lock but never got
# to the trap that removes it, and the real resource was never written
# this time around either. The PID below belongs to a process that
# finished and exited long before this image was even built, so it can
# never be alive inside this (or any) container.
RUN printf 'PID=99999\nSTARTED=2026-10-21T03:14:00Z\n' > /var/run/infra/apply.lock && \
    cp /var/run/infra/apply.lock /infra/.original-lock && \
    echo "Build server release notes v2" > /etc/release-notes.txt
WORKDIR /workspace
CMD ["sleep", "infinity"]
