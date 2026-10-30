FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    bash coreutils less vim-tiny passwd \
    && rm -rf /var/lib/apt/lists/*
RUN groupadd appdata && \
    useradd -r -g appdata -s /usr/sbin/nologin producer && \
    useradd -r -g appdata -s /usr/sbin/nologin consumer && \
    mkdir -p /var/lib/appdata && \
    chown producer:appdata /var/lib/appdata && \
    chmod 2775 /var/lib/appdata
RUN mkdir -p /etc/profile.d
RUN printf 'umask 022\n' > /etc/default/producer
COPY run-producer.sh /usr/local/bin/run-producer.sh
RUN chmod 755 /usr/local/bin/run-producer.sh
CMD ["sleep", "infinity"]
