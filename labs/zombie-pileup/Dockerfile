FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    python3 procps bash coreutils vim-tiny less \
    && rm -rf /var/lib/apt/lists/*
COPY supervisor.py /usr/local/bin/supervisor.py
COPY worker_pool.py /usr/local/bin/worker_pool.py
COPY restart-supervisor /usr/local/bin/restart-supervisor
COPY restart-worker-pool /usr/local/bin/restart-worker-pool
COPY restart-workers /usr/local/bin/restart-workers
RUN chmod 755 /usr/local/bin/supervisor.py /usr/local/bin/worker_pool.py \
    /usr/local/bin/restart-supervisor /usr/local/bin/restart-worker-pool \
    /usr/local/bin/restart-workers
CMD ["sleep", "infinity"]
