#!/usr/bin/env bash
set -e
python3 /opt/metrics_server.py &
/usr/local/bin/prometheus --config.file=/infra/prometheus.yml --storage.tsdb.path=/var/lib/prometheus --log.level=warn &
/opt/grafana/bin/grafana server --homepath /opt/grafana --config /etc/grafana/grafana.ini >/var/log/grafana.log 2>&1 &
exec sleep infinity
