FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    python3 python3-pip ca-certificates \
    && rm -rf /var/lib/apt/lists/*
RUN pip3 install --no-cache-dir --break-system-packages "ansible-core==2.17.*"
RUN mkdir -p /infra /etc/app /workspace
RUN printf 'app_name=acme-checkout\nlog_level=info\n' > /infra/app.conf.pristine
COPY inventory /infra/inventory
COPY playbook.yml /infra/playbook.yml
WORKDIR /infra
RUN cp app.conf.pristine /etc/app/app.conf
WORKDIR /workspace
CMD ["sleep", "infinity"]
