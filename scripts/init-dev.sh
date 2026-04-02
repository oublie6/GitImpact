#!/usr/bin/env bash
# 初始化本地开发环境的最小脚本：
# 1. 确保脚本本身可执行
# 2. 为后端准备 config.yaml
set -euo pipefail
cd "$(dirname "$0")/.."
chmod +x scripts/*.sh
cp -n backend/config.example.yaml backend/config.yaml || true
echo "[init-dev] config prepared at backend/config.yaml"
echo "[init-dev] vendor mode is enabled by default (GOFLAGS=-mod=vendor)."
echo "[init-dev] run './scripts/init-db.sh' before starting backend service."
