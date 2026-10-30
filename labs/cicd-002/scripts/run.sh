#!/usr/bin/env bash
echo "=== build ==="; /opt/pipeline/build.sh; echo "=== test ==="; /opt/pipeline/test.sh; echo "=== deploy ==="; /opt/pipeline/deploy.sh
