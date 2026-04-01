#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."
chmod +x scripts/*.sh
cp -n backend/config.example.yaml backend/config.yaml || true
echo "[init-dev] config prepared at backend/config.yaml"
echo "[init-dev] vendor mode is enabled by default (GOFLAGS=-mod=vendor)."
