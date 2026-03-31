#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."
chmod +x scripts/*.sh
cd backend && go mod tidy
