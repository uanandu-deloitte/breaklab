FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    bash coreutils less vim-tiny passwd \
    && rm -rf /var/lib/apt/lists/*
RUN useradd -r -s /usr/sbin/nologin jobsrunner
RUN mkdir -p /var/spool/jobs/done /etc/jobs /usr/local/libexec && \
    chown -R jobsrunner:jobsrunner /var/spool/jobs
RUN for i in 1 2 3 4 5; do echo "{\"job\":$i}" > /var/spool/jobs/job-$i.json; done && \
    chown jobsrunner:jobsrunner /var/spool/jobs/job-*.json
COPY worker.sh /usr/local/libexec/worker.sh
COPY process-jobs.sh /usr/local/bin/process-jobs
RUN chmod 755 /usr/local/libexec/worker.sh /usr/local/bin/process-jobs
RUN printf 'WORKER_BIN=/usr/local/libexec/worker.sh\n' > /etc/jobs/worker.conf
CMD ["sleep", "infinity"]
