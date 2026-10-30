FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    bash coreutils less vim-tiny \
    && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /var/log/app /etc/agent /var/lib/agent
RUN for i in 1 2 3 4 5; do echo "old request $i" >> /var/log/app/app.log.1; done
RUN for i in 1 2 3 4 5 6 7 8; do echo "request $i" >> /var/log/app/app.log; done
RUN printf 'LOGFILE=/var/log/app/app.log\n' > /etc/agent/agent.conf
COPY collect-metrics /usr/local/bin/collect-metrics
RUN chmod 755 /usr/local/bin/collect-metrics
RUN /usr/local/bin/collect-metrics
CMD ["sleep", "infinity"]
