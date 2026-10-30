FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    coreutils tar gzip findutils \
    && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /data/to-backup /opt/backups && \
    echo "note contents" > "/data/to-backup/notes.txt" && \
    echo "draft contents" > "/data/to-backup/report draft.txt" && \
    echo "service: acme" > "/data/to-backup/config.yml"
COPY backup.sh /opt/app/backup.sh
RUN chmod +x /opt/app/backup.sh
CMD ["sleep", "infinity"]
