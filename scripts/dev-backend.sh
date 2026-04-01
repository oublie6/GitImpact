#!/usr/bin/env bash
# 启动后端开发服务；如缺少配置文件则先复制示例配置。
set -euo pipefail
cd "$(dirname "$0")/../backend"
cp -n config.example.yaml config.yaml || true
GITIMPACT_CONFIG=./config.yaml go run ./cmd/server
