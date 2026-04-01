# 部署指南

## 推荐部署模式

当前推荐且已落地的方案是：

- 前端在联网环境构建为 `dist`
- 将 `dist` 同步到 `backend/web/dist`
- 由后端统一托管前端页面、静态资源和 SPA history fallback

这样离线环境只需要启动 GitImpact 后端服务即可。

## 开发环境部署

- 使用仓库自带脚本启动即可。
- 建议保留 `mixed` 认证模式，方便使用 `config_users` 验证。
- OpenCode CLI 与 Git 直接安装在宿主机。
- 前端开发模式下仍通过 Vite 代理 `/api` 到后端。

## 测试环境部署

- 使用独立数据库。
- 明确 `workdir` 持久化路径。
- 把 `jwt_secret`、数据库凭证、OpenCode 配置替换为测试环境专用值。

## 生产部署注意事项

- 更换所有默认密码和 `jwt_secret`。
- 不建议继续使用示例 `config_users`。
- 控制 `local_cache_dir` 和 `workdir` 的磁盘权限。
- 为数据库和运行目录设计备份策略。

## Docker 构建与运行

构建：

```bash
docker build -t gitimpact/all-in-one:test .
```

运行：

```bash
docker run --rm -p 8080:8080 gitimpact/all-in-one:test
```

Dockerfile 特点：

- 包含前端构建阶段和后端构建阶段
- 前端构建阶段使用 `npm ci` 与 `vite build`
- 构建阶段基于 `golang:1.22`
- 后端仍使用 `GOFLAGS=-mod=vendor`
- 运行时镜像已包含 `web/dist`
- 使用 `GOFLAGS=-mod=vendor`
- 不在构建期联网拉取 Go 模块
- 需要联网的部分只发生在前端 npm 依赖安装阶段和镜像拉取阶段，而不是离线部署阶段

## 二进制离线部署

在线构建：

```powershell
powershell -ExecutionPolicy Bypass -File scripts/package-offline.ps1
```

离线部署：

1. 解压离线包
2. 复制 `config.example.yaml` 为 `config.yaml`
3. 调整数据库与 OpenCode 配置
4. 运行 `run-offline.ps1` 或直接启动后端二进制

## 前端静态资源部署说明

- 后端默认读取 `frontend.enabled` 和 `frontend.dist_dir`
- 默认 `frontend.dist_dir` 为 `./web/dist`
- 部署包和 Docker 镜像都会把前端 dist 放到这个目录
- 浏览器访问 `/`、`/tasks`、`/repositories` 等路径时，由后端回退到 `index.html`

## 离线/vendor 构建说明

Makefile 默认：

```makefile
GOFLAGS ?= -mod=vendor
```

Dockerfile 默认：

```dockerfile
ENV GOFLAGS=-mod=vendor
```

因此常规构建链路依赖仓库中已经提交的 `backend/vendor/`。如果更新了依赖但未同步 vendor，构建和测试可能失败。

## 交叉编译

Linux AMD64：

```bash
make build-linux-amd64
```

产物默认输出到：

```text
bin/gitimpact-backend-linux-amd64
```

## 环境变量说明

- `GITIMPACT_CONFIG`：指定后端配置文件路径

目前项目没有实现更多环境变量覆盖机制，主要依赖 YAML 文件。

## 持久化目录说明

建议在部署时持久化以下目录：

- `workdir.root`
- `workdir.repo_cache`
- `workdir.artifacts`
- 如自行扩展了文件型报告落盘，也应持久化 `workdir.reports`

前端静态资源 `web/dist` 一般随发布包或镜像分发，不需要单独持久化。

## 安全建议

- JWT 密钥不要使用示例值
- 不要把包含真实凭证的 `config.yaml` 提交到仓库
- OpenCode attach 服务如果暴露到网络，应加访问控制
- 仓库缓存目录不要使用过于宽泛的系统路径
