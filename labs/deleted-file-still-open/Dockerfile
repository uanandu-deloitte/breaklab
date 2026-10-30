FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    bash coreutils procps lsof less vim-tiny \
    && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /var/log/app /var/run
# Pre-grow the "live" log so it looks like what a long-running service
# would have accumulated before this incident started.
RUN dd if=/dev/zero of=/var/log/app/service.log bs=1M count=20 >/dev/null 2>&1
COPY logger-daemon.sh /usr/local/bin/logger-daemon.sh
COPY disk-usage-report /usr/local/bin/disk-usage-report
RUN chmod 755 /usr/local/bin/logger-daemon.sh /usr/local/bin/disk-usage-report
CMD ["bash", "-c", "nohup /usr/local/bin/logger-daemon.sh >/var/log/app/daemon.out 2>&1 & exec bash -c 'while true; do wait -n 2>/dev/null || sleep 0.2; done'"]
