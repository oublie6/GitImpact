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

### config_users 认证规则
- 配置用户现在仅使用明文字段 `password` 进行登录校验（不再使用 `password_hash` / `password_plain` / `allow_dev_plain`）。
- 数据库用户登录逻辑保持不变，仍使用 bcrypt 对 `password_hash` 做比对。
- `auth.mode=mixed` 时，登录顺序为先 `config_users`，再 DB。

## 数据库初始化说明
- MySQL: `mysql -uroot -proot < sql/mysql/init.sql`
- 达梦: 使用达梦客户端在目标 schema 执行 `sql/dameng/init.sql`

### 数据库驱动与连接方式
- `database.type=mysql`：使用 `gorm.io/driver/mysql`，连接串来自 `database.dsn`。
- `database.type=dameng`：使用 `github.com/godoes/gorm-dameng`。
  - 若 `database.dsn` 非空，直接使用该达梦 DSN（如 `dm://SYSDBA:SYSDBA@127.0.0.1:5236?schema=SYSDBA`）。
  - 若 `database.dsn` 为空，则按 `database.dameng.host/port/user/password/dbname` 自动拼接，
    其中 `dbname` 会映射为达梦 DSN 的 `schema` 参数。

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
- `make build-linux-amd64`：交叉编译 Linux amd64 后端二进制到 `bin/`
- `make clean`：清理构建产物
- `make docker-build` / `make docker-run` / `make docker-push` / `make docker-build-run` / `make deploy`：容器与部署入口

## Docker 构建
- 直接执行：`docker build -t gitimpact/backend:test .`
- 或使用 Makefile：`make docker-build`
- Dockerfile 在 builder 阶段使用 `backend/vendor` + `GOFLAGS=-mod=vendor` 离线构建，不拉取远端依赖。

## 当前限制说明
- 达梦已切换为专用 GORM 驱动接入（`github.com/godoes/gorm-dameng`）。
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
在仓库根目录执行：`docker build -t gitimpact/backend:test .`。
