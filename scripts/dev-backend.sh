#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/../backend"
cp -n config.example.yaml config.yaml || true
GITIMPACT_CONFIG=./config.yaml go run ./cmd/server
