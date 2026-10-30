FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    nginx curl procps iproute2 lsof vim-tiny less python3 \
    && rm -rf /var/lib/apt/lists/*
RUN echo "BreakLab OK" > /var/www/html/index.html
CMD ["sleep", "infinity"]
