# 排障指南

## 启动失败

排查顺序：

1. 确认 `backend/config.yaml` 是否存在
2. 确认 `GITIMPACT_CONFIG` 指向的路径是否正确
3. 查看启动日志里的 `load config failed` 或 `db init failed`
4. 确认端口是否被占用
5. 如果页面能打开 API 不能用，继续检查数据库初始化与后端日志

## 配置错误

常见现象：

- YAML 缩进错误
- `database.type` 写错
- MySQL 缺少 `database.dsn`
- 达梦 `port` 不是整数

排查方法：

- 用 YAML 校验器检查配置
- 对照 [配置参考](./config-reference.md)

## 数据库连接失败

排查：

- MySQL/达梦服务是否可访问
- 用户名、密码、库名或 schema 是否正确
- 初始化 SQL 是否执行过
- 网络和防火墙是否允许连接

## 启动报错：database schema is not initialized

现象：

- 日志提示缺失核心表（如 `users`、`analysis_tasks`）
- 同时提示需要执行 `sql/mysql/init.sql` 或 `sql/dameng/init.sql`

原因：

- 项目默认禁用 AutoMigrate，不会在启动时自动建表

处理：

1. 根据 `database.type` 执行初始化脚本
   - MySQL：`sql/mysql/init.sql`
   - 达梦：`sql/dameng/init.sql`
2. 再次启动服务
3. 可先执行 `make init-db` 查看推荐脚本路径，避免选错

## OpenCode 调用失败

常见原因：

- `opencode` 不在 PATH 中
- `binary_path` 配置错误
- `attach_url` 不可达
- OpenCode 命令执行超时

排查：

- 在宿主机手动执行 `opencode --help`
- 检查 worker 生成的 prompt 文件是否存在
- 查看 `analysis_reports.raw_stderr`

## 任务卡住

当前实现没有队列监控，“卡住”通常表现为长期停留在 `running`。

排查：

- 查看 `task_logs`
- 检查 OpenCode 是否长时间无返回
- 检查 Git clone/fetch/checkout 是否被卡住
- 确认服务进程是否仍存活

## 报告为空

排查：

- 任务是否真的成功
- `generate_markdown` / `generate_structured` 是否启用
- `raw_stdout` / `raw_stderr` 是否有内容
- OpenCode 输出是否符合预期格式

## 登录失败

排查：

- `auth.mode` 是否符合预期
- `config_users` 密码是否按明文配置
- 数据库用户的 `password_hash` 是否正确
- 用户状态是否为 `disabled`

## 跨平台构建问题

当前 Makefile 与脚本偏向 Unix 风格。

排查：

- Windows 下是否使用 Git Bash、WSL 或兼容的 `make`
- `mkdir -p`、`rm -rf`、`chmod` 等命令是否可用
- Go 交叉编译时是否显式设置 `CGO_ENABLED=0`

## 前端离线部署后页面打不开

排查：

- `backend/web/dist/index.html` 是否存在
- `frontend.enabled` 是否为 `true`
- `frontend.dist_dir` 是否指向正确目录
- 是否启动了新版本后端

## 前端刷新二级路由 404

排查：

- 确认当前运行的后端是否包含前端 fallback 逻辑
- 确认访问路径不是 `/api/*`
- 确认 `index.html` 可被后端读取

## 前端接口请求仍指向 localhost

排查：

- 检查是否使用了旧版 dist
- 检查 built dist 中是否仍残留 `http://127.0.0.1:8080`
- 重新执行 `scripts/build-frontend.ps1`

## 离线环境为什么不需要 npm install

因为离线环境只消费已经构建完成的 `dist`，不会重新执行 Vite 构建。

如果你要在离线环境重新构建前端，那是“离线重建”问题，不属于当前默认部署链路。

## vendor / 依赖问题

现象：

- `inconsistent vendoring`
- 构建或测试提示模块与 `vendor/` 不一致

排查：

- 检查 `backend/go.mod`、`backend/go.sum` 是否被修改
- 在允许联网的开发环境执行 `go mod vendor`
- 重新提交同步后的 `backend/vendor/`

## 常见 HTTP 错误

- `400`：请求体不合法、业务校验失败
- `401`：缺少 Bearer Token、JWT 无效、登录失败
- `404`：任务/仓库/报告不存在
- `409`：注册用户名冲突
- `500`：Git 操作失败、数据库异常或其他未处理问题
