#!/usr/bin/env bash
# 启动前端开发服务。
# 仅在 node_modules 缺失时补装依赖，避免每次开发启动都联网拉包。
set -euo pipefail
cd "$(dirname "$0")/../frontend"
if [ ! -d node_modules ]; then
  npm install
fi
npm run dev
