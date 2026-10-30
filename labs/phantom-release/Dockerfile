FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    bash coreutils git tar gzip ca-certificates less vim-tiny \
    && rm -rf /var/lib/apt/lists/*
RUN git config --global user.email "ops@breaklab.local" && \
    git config --global user.name "BreakLab Ops" && \
    git config --global init.defaultBranch release

RUN mkdir -p /srv/app-repo/src /opt/pipeline/artifacts /srv/deployed
COPY src/app.sh /srv/app-repo/src/app.sh
COPY src/README.md /srv/app-repo/README.md
COPY build.sh /opt/pipeline/build.sh
RUN chmod 755 /opt/pipeline/build.sh /srv/app-repo/src/app.sh

RUN cd /srv/app-repo && git init -q -b release && \
    git add -A && git commit -q -m "Initial checkout flow" && \
    git tag v2.0.0 && \
    git branch develop

RUN echo "release" > /opt/pipeline/BUILD_REF
RUN /opt/pipeline/build.sh

CMD ["sleep", "infinity"]
