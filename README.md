# GitImpact（Git 变更影响分析与报告平台）

## 快速启动
1. 初始化：`./scripts/init-dev.sh`
2. 准备数据库：MySQL 执行 `sql/mysql/init.sql`（或达梦执行 `sql/dameng/init.sql`）
3. 后端启动：`./scripts/dev-backend.sh`
4. 前端启动：`./scripts/dev-frontend.sh`

## 配置说明
后端读取 `backend/config.yaml`，可参考 `backend/config.example.yaml`，核心字段包括：
- `auth.mode`（config/db/mixed）
- `database.type`（mysql/dameng）
- `database.dsn`
- `opencode.*`
- `workdir.*`

## 数据库初始化说明
- MySQL: `mysql -uroot -proot < sql/mysql/init.sql`
- 达梦: 使用达梦客户端在目标 schema 执行 `sql/dameng/init.sql`

## 默认管理员说明
- 当 `auth.init_admin_enabled=true` 且模式非 config 时，后端启动会尝试初始化 admin 用户（若不存在）。
- SQL 也提供了默认 admin 初始化语句。
- 默认密码对应 hash 为示例值，建议上线前替换。

## OpenCode 集成说明
- 默认 `CLIAnalyzer` 调用 `opencode run`。
- 若配置 `opencode.attach_url`，将优先使用 `opencode run --attach <url>`。
- 系统先生成分析材料：`changed_files.txt`、`diff.patch`、`commit_log.txt`、`repo_manifest.md`、`analysis_prompt.md`。
- 若结构化 JSON 解析失败，依旧保存 Markdown 与原始 stdout/stderr。

## 当前限制说明
- `ServerAnalyzer` 目前为**占位实现**。
- 达梦通过 MySQL 兼容 DSN/语法进行适配验证，实际生产请按驱动与方言深化。
