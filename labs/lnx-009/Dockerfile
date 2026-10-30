FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    bash coreutils gawk \
    && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /workspace /opt/fixtures
RUN printf '%s\n' \
    '203.0.113.5 2026-01-10T09:00:01Z GET /api/orders 200 512' \
    '203.0.113.5 2026-01-10T09:00:05Z GET /api/orders 200 498' \
    '198.51.100.9 2026-01-10T09:00:10Z GET /api/cart 200 256' \
    '203.0.113.5 2026-01-10T09:00:15Z GET /api/orders 404 64' \
    '192.0.2.44 2026-01-10T09:00:20Z GET /api/orders 200 512' \
    '203.0.113.5 2026-01-10T09:00:25Z POST /api/checkout 500 128' \
    '198.51.100.9 2026-01-10T09:00:30Z GET /api/orders 200 512' \
    '203.0.113.5 2026-01-10T09:00:35Z GET /healthz 200 32' \
    '192.0.2.44 2026-01-10T09:00:40Z GET /api/cart 404 64' \
    '203.0.113.5 2026-01-10T09:00:45Z GET /api/orders 200 512' \
    '198.51.100.9 2026-01-10T09:00:50Z GET /api/orders 500 128' \
    '203.0.113.5 2026-01-10T09:00:55Z GET /api/orders 200 512' \
    '192.0.2.44 2026-01-10T09:01:00Z GET /api/login 401 64' \
    '203.0.113.5 2026-01-10T09:01:05Z GET /api/orders 200 512' \
    '198.51.100.9 2026-01-10T09:01:10Z GET /api/cart 200 256' \
    '203.0.113.5 2026-01-10T09:01:15Z GET /api/orders 200 512' \
    > /opt/fixtures/access.log.original && \
    cp /opt/fixtures/access.log.original /workspace/access.log
WORKDIR /workspace
CMD ["sleep", "infinity"]
