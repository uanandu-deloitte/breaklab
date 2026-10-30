#!/usr/bin/env bash
set -euo pipefail
kubectl patch service orders -p '{"spec":{"ports":[{"port":80,"targetPort":8080}]}}'
