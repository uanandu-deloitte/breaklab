FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    bash coreutils \
    && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /workspace
RUN printf '%b\n' \
    'WEB-01.Example.com, web-02.example.com' \
    'DB-01.example.COM' \
    'cache-01.example.com,  CACHE-01.Example.Com' \
    '  lb-01.example.com  ' \
    'APP-01.EXAMPLE.COM' \
    'web-01.example.com' \
    'WEB-02.Example.Com,\tdb-01.example.com' \
    'LB-01.example.com' \
    'app-01.example.com, web-01.example.com' \
    'cache-01.example.com' \
    > /workspace/hosts.txt
WORKDIR /workspace
CMD ["sleep", "infinity"]
