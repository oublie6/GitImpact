#!/usr/bin/env bash
# 启动前端开发服务。当前脚本会先执行 npm install，适合首次启动或依赖变化后使用。
set -euo pipefail
cd "$(dirname "$0")/../frontend"
npm install
npm run dev
