FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    bash coreutils findutils grep \
    && rm -rf /var/lib/apt/lists/*

# The live, active configuration — this is the one the running service
# actually reads, and it's the only one that's supposed to change.
RUN mkdir -p /opt/app/config && \
    printf '%s\n' \
      'service: acme-backend' \
      'database_host: database-old.internal' \
      'port: 5432' \
      > /opt/app/config/app.conf

# An archived copy of a previous release's own config — a snapshot of
# what v1 shipped with, kept for audit/rollback reference. It happens to
# mention the same old hostname, but it is not live and must not change.
RUN mkdir -p /opt/app/releases/v1/config && \
    printf '%s\n' \
      'service: acme-backend' \
      'database_host: database-old.internal' \
      'port: 5432' \
      > /opt/app/releases/v1/config/app.conf

# A changelog documenting the migration historically — it's supposed to
# say the old hostname was replaced, as a historical record of the fact.
RUN printf '%s\n' \
      '# Changelog' \
      '' \
      '## v2.0' \
      '- migrated production database host from database-old.internal to database.internal' \
      '' \
      '## v1.0' \
      '- initial release, backend connects to database-old.internal' \
      > /opt/app/CHANGELOG.md

# Migration notes, same idea — historical record, not live config.
RUN mkdir -p /opt/app/docs && \
    printf '%s\n' \
      '# Migration Notes' \
      '' \
      'Historical record: the v1 release used database-old.internal as its' \
      'database host. This document intentionally preserves that reference' \
      'for audit purposes and must not be edited.' \
      > /opt/app/docs/migration-notes.md

# A tiny script the learner can run themselves to self-check the active
# config, without needing to know this lab's own check.sh logic.
RUN mkdir -p /opt/app/bin && \
    printf '%s\n' \
      '#!/usr/bin/env bash' \
      'if grep -qF "database.internal" /opt/app/config/app.conf && ! grep -qF "database-old.internal" /opt/app/config/app.conf; then' \
      '  echo "OK: active config points at database.internal"' \
      'else' \
      '  echo "STALE: active config still references the old backend"' \
      'fi' \
      > /opt/app/bin/check-config.sh && \
    chmod +x /opt/app/bin/check-config.sh

RUN mkdir -p /workspace
WORKDIR /workspace
CMD ["sleep", "infinity"]
