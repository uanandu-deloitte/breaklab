FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    bash coreutils \
    && rm -rf /var/lib/apt/lists/*

# A failed deployment scattered these three files into the wrong
# subdirectories under /opt/acme-app. Each file's own content and name
# make clear where it actually belongs:
#   - acme-app (an executable script)  belongs in bin/
#   - app.conf (a config file)         belongs in config/
#   - data.db  (an application data file) belongs in data/
# Instead, right now:
#   bin/    has data.db   (should be in data/)
#   config/ has acme-app  (should be in bin/, and needs +x)
#   data/   has app.conf  (should be in config/)
#   logs/   is empty, and already correct
RUN mkdir -p /opt/acme-app/bin /opt/acme-app/config /opt/acme-app/data /opt/acme-app/logs && \
    printf '#!/usr/bin/env bash\necho "acme-app running"\n' > /opt/acme-app/config/acme-app && \
    chmod +x /opt/acme-app/config/acme-app && \
    printf 'service: acme-app\nenv: production\n' > /opt/acme-app/data/app.conf && \
    printf 'ACME-DATA-V1\nrecord-count: 42\n' > /opt/acme-app/bin/data.db

RUN mkdir -p /workspace
WORKDIR /workspace
CMD ["sleep", "infinity"]
