#!/usr/bin/env bash
set -euo pipefail
sed -i '/SITE_MESSAGE/d' /etc/bash.bashrc
infra-apply
