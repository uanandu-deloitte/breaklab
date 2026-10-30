FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    bash coreutils \
    && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /opt/pipeline /workspace

# The real deploy token's plaintext is never baked into this image —
# only a sha256 hash of it, computed outside this Dockerfile. deploy.sh
# (once fixed) compares a hash of whatever $DEPLOY_TOKEN it's given at
# run time against this file; the plaintext itself only ever exists in
# this lab's README, handed to the learner the way ops would hand over
# a real one-time deploy credential.
RUN echo "7faed5537a001553dcb08a6a3760bfdbd441d6288896110666e20dc55133c367" > /opt/pipeline/.deploy_token.sha256

COPY deploy.sh /opt/pipeline/deploy.sh
RUN chmod +x /opt/pipeline/deploy.sh

WORKDIR /workspace
CMD ["sleep", "infinity"]
