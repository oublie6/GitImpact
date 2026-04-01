# 前端构建与打包指南

## 版本要求

- Node.js 18+
- npm 9+

## 构建命令

### 只构建前端并同步到后端静态目录

```powershell
powershell -ExecutionPolicy Bypass -File scripts/build-frontend.ps1
```

### 联合构建发布目录

```powershell
powershell -ExecutionPolicy Bypass -File scripts/build-release.ps1
```

### 生成离线部署包

```powershell
powershell -ExecutionPolicy Bypass -File scripts/package-offline.ps1
```

### 本地验证离线部署链路

```powershell
powershell -ExecutionPolicy Bypass -File scripts/verify-offline.ps1
```

## 产物目录说明

- `frontend/dist/`：Vite 输出目录
- `backend/web/dist/`：后端托管目录
- `artifacts/release/`：联合构建产物
- `artifacts/offline/`：压缩后的离线包

## 打包方式

默认脚本会生成：

- 后端二进制
- `config.example.yaml`
- `web/dist/`
- `run-offline.ps1`

这些文件会被压缩为 zip 包，适合内网传输和投产。

## API 访问方式

前端统一通过 `VITE_API_BASE_URL` 访问 API。

当前约定：

- 开发：`/api`，由 Vite 代理到后端
- 生产：`/api`，由浏览器直接请求同源后端

## 注意事项

- 离线部署和离线重新构建是两回事
- 当前项目已经实现“离线部署”
- 如果要在完全离线环境重新构建前端，仍需要额外准备 npm 离线源或依赖缓存

## 不建议的做法

- 在离线环境执行 `npm install`
- 让生产环境继续依赖 Vite dev server
- 在前端代码中硬编码后端主机名和端口
