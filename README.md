# GitImpact（Git 变更影响分析与报告平台）

## 快速启动
1. 初始化：`./scripts/init-dev.sh`
2. 准备数据库：MySQL 执行 `sql/mysql/init.sql`（或达梦执行 `sql/dameng/init.sql`）
3. 后端启动：`./scripts/dev-backend.sh`
4. 前端启动：`./scripts/dev-frontend.sh`

## 配置说明
后端读取 `backend/config.yaml`，可参考 `backend/config.example.yaml`，核心字段包括：
- `auth.mode`（config/db/mixed）
- `auth.jwt_secret`
- `auth.token_expire_minutes`
- `auth.enable_register`
- `auth.init_admin_enabled`
- `auth.config_users`
- `database.type`（mysql/dameng）
- `database.dsn`
- `database.mysql.*`
- `database.dameng.*`
- `opencode.*`
- `workdir.*`

## 数据库初始化说明
- MySQL: `mysql -uroot -proot < sql/mysql/init.sql`
- 达梦: 使用达梦客户端在目标 schema 执行 `sql/dameng/init.sql`

## 默认管理员说明
- 当 `auth.init_admin_enabled=true` 且模式非 config 时，后端启动会尝试初始化 admin 用户（若不存在）。
- SQL 也提供了默认 admin 初始化语句。
- 默认管理员用户名 `admin`，默认密码请通过配置或 SQL 初始化后立即重置。

## OpenCode 集成说明
- 默认 `CLIAnalyzer` 调用 `opencode run`。
- 若配置 `opencode.attach_url`，将优先使用 `opencode run --attach <url>`。
- 系统先生成分析材料：`changed_files.txt`、`diff.patch`、`commit_log.txt`、`repo_manifest.md`、`analysis_prompt.md`。
- 若结构化 JSON 生成失败，依旧保存 Markdown 与原始 stdout/stderr。
- `ServerAnalyzer` 当前为**占位实现**。

## vendor 使用说明
- 后端构建/测试默认通过 `GOFLAGS=-mod=vendor` 执行。
- 已提供 `backend/vendor/` 目录，Makefile 中默认开启 vendor 模式。

## Makefile 使用说明
- `make build`：构建后端二进制到 `bin/`
- `make test`：执行后端测试（vendor 模式）
- `make clean`：清理构建产物
- `make docker-build` / `make docker-run` / `make docker-push` / `make docker-build-run` / `make deploy`：容器与部署入口

## 当前限制说明
- 达梦当前采用 MySQL 兼容策略接入（可运行基线），生产环境建议替换为专用驱动与方言增强。
- `ServerAnalyzer` 为占位实现。


## Windows PowerShell 交叉编译 Linux x64
```powershell
$env:CGO_ENABLED="0"
$env:GOOS="linux"
$env:GOARCH="amd64"
$env:GOFLAGS="-mod=vendor"
cd backend
go build -trimpath -ldflags "-s -w" -o ../bin/gitimpact-backend-linux-amd64 ./cmd/server
```

也可以直接在仓库根目录执行：`make build-linux-amd64`。

## Docker 构建验证
在仓库根目录执行：`make docker-build`。
