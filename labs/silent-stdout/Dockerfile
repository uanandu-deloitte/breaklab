FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    python3 bash coreutils procps less vim-tiny \
    && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /etc/app /var/log/app
RUN printf 'exec python3 /usr/local/bin/app.py\n' > /etc/app/launch-cmd
RUN printf 'log_target=stdout\n' > /etc/app/logging.conf
COPY app.py /usr/local/bin/app.py
COPY supervisor.sh /usr/local/bin/supervisor.sh
RUN chmod 755 /usr/local/bin/app.py /usr/local/bin/supervisor.sh
CMD ["/usr/local/bin/supervisor.sh"]
