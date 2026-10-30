FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    bash coreutils gawk \
    && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /workspace /opt/fixtures
RUN printf '%s\n' \
    '203.0.113.5 2026-10-19T10:00:01Z GET /api/orders 200 512' \
    '198.51.100.9 2026-10-19T10:00:02Z GET /api/orders 500 128' \
    '203.0.113.5 2026-10-19T10:00:03Z GET /api/cart 200 256' \
    '198.51.100.9 2026-10-19T10:00:04Z POST /api/checkout 500 128' \
    '192.0.2.44 2026-10-19T10:00:05Z GET /api/orders 404 64' \
    '198.51.100.9 2026-10-19T10:00:06Z GET /api/cart 500 128' \
    '203.0.113.5 2026-10-19T10:00:07Z GET /healthz 200 32' \
    '192.0.2.44 2026-10-19T10:00:08Z GET /api/orders 500 128' \
    '198.51.100.9 2026-10-19T10:00:09Z GET /api/orders 200 512' \
    '192.0.2.44 2026-10-19T10:00:10Z GET /api/cart 200 256' \
    '198.51.100.9 2026-10-19T10:00:11Z GET /api/orders 500 128' \
    '203.0.113.5 2026-10-19T10:00:12Z GET /api/orders 200 512' \
    > /opt/fixtures/access.log.original && \
    cp /opt/fixtures/access.log.original /workspace/access.log
WORKDIR /workspace
CMD ["sleep", "infinity"]
