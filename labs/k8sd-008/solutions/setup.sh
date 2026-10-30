#!/usr/bin/env bash
set -euo pipefail
export KUBECONFIG=/etc/rancher/k3s/k3s.yaml
sed -i 's/targetPort: 8080/targetPort: 80/' /chart/templates/service.yaml
helm upgrade web /chart
