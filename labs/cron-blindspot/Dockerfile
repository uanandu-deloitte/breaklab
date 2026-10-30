FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    bash coreutils tar cron less vim-tiny passwd \
    && rm -rf /var/lib/apt/lists/*
RUN useradd -m -s /bin/bash opsuser

COPY backup.sh /usr/local/bin/backup.sh
RUN chmod 755 /usr/local/bin/backup.sh

RUN mkdir -p /var/app/data /var/backups && \
    echo v1 > /var/app/data/marker.txt && \
    chown -R opsuser:opsuser /var/app/data /var/backups
RUN touch /var/log/backup.log && chown opsuser:opsuser /var/log/backup.log

# A fresh baseline backup, taken the same way a real cron run would: as
# opsuser, with $HOME as the working directory.
RUN su - opsuser -c '/usr/local/bin/backup.sh'

RUN echo "0 2 * * * /usr/local/bin/backup.sh >> /var/log/backup.log 2>&1" | crontab -u opsuser -

CMD ["sleep", "infinity"]
