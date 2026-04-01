# 配置参考

GitImpact 通过 `backend/config.yaml` 读取配置，默认示例文件为 `backend/config.example.yaml`。后端启动时会优先读取环境变量 `GITIMPACT_CONFIG` 指定的路径，否则使用当前目录下的 `./config.yaml`。

## 顶层结构

```yaml
server:
auth:
database:
opencode:
workdir:
frontend:
```

## server

### `server.port`

- 类型：string
- 默认值：`"8080"`
- 必填：否
- 作用：后端 HTTP 服务监听端口
- 示例：

```yaml
server:
  port: "8080"
```

## auth

### `auth.mode`

- 类型：string
- 默认值：无
- 必填：是
- 可选值：`config`、`db`、`mixed`
- 作用：
  - `config`：只允许配置用户登录
  - `db`：只允许数据库用户登录
  - `mixed`：先配置用户、后数据库用户

### `auth.jwt_secret`

- 类型：string
- 默认值：无
- 必填：强烈建议
- 作用：JWT 签名密钥

### `auth.token_expire_minutes`

- 类型：int
- 默认值：`120`
- 必填：否
- 作用：JWT 过期时间，单位分钟

### `auth.enable_register`

- 类型：bool
- 默认值：`false`
- 必填：否
- 作用：是否允许注册数据库用户

### `auth.init_admin_enabled`

- 类型：bool
- 默认值：`false`
- 必填：否
- 作用：启动时是否自动创建数据库默认管理员

### `auth.config_users`

- 类型：数组
- 默认值：空
- 必填：仅 `config` / `mixed` 模式下建议配置
- 作用：定义静态登录用户

每个元素字段：

- `username`：登录名
- `password`：明文密码，当前代码直接比对
- `display_name`：展示名
- `email`：邮箱
- `role`：角色，空值时后端会回退为 `user`
- `status`：`active` 或 `disabled`

重要说明：

- 当前代码不会对 `config_users.password` 做哈希。
- `config_users` 用户不会写入数据库。
- 配置用户登录后生成的 JWT 中 `user_id` 为 `0`。

## database

### `database.type`

- 类型：string
- 默认值：空字符串按 MySQL 处理
- 必填：建议填写
- 可选值：`mysql`、`dameng`

### `database.dsn`

- 类型：string
- 默认值：空
- 必填：MySQL 必填；达梦可选
- 作用：
  - MySQL 时直接作为 GORM DSN
  - 达梦时若填写，直接作为最终 DSN

### `database.mysql.*`

- 类型：map[string]string
- 默认值：空
- 作用：仅作为配置展示和补充说明，当前 MySQL 实际连接只用 `database.dsn`

字段示例：

- `host`
- `port`
- `dbname`
- `user`
- `password`

### `database.dameng.*`

- 类型：map[string]string
- 默认值：空
- 作用：当 `database.type=dameng` 且 `database.dsn` 为空时用于拼接达梦 DSN

必需字段：

- `host`
- `port`
- `user`
- `password`

可选字段：

- `dbname`：实际会映射为达梦 DSN 的 `schema`

## opencode

### `opencode.binary_path`

- 类型：string
- 默认值：`opencode`
- 必填：否
- 作用：OpenCode CLI 可执行文件路径

### `opencode.attach_url`

- 类型：string
- 默认值：空
- 必填：否
- 作用：非空时，CLIAnalyzer 会追加 `--attach <url>`

### `opencode.default_model`

- 类型：string
- 默认值：空
- 必填：否
- 作用：CLIAnalyzer 会追加 `--model`

### `opencode.default_agent`

- 类型：string
- 默认值：空
- 必填：否
- 作用：CLIAnalyzer 会追加 `--agent`

### `opencode.timeout_seconds`

- 类型：int
- 默认值：`600`
- 必填：否
- 作用：OpenCode 命令最大执行时间

## workdir

### `workdir.root`

- 类型：string
- 默认值：无
- 必填：建议填写
- 作用：运行时根目录

### `workdir.repo_cache`

- 类型：string
- 默认值：无
- 必填：建议填写
- 作用：本地仓库缓存目录

### `workdir.artifacts`

- 类型：string
- 默认值：无
- 必填：建议填写
- 作用：任务材料目录

### `workdir.reports`

- 类型：string
- 默认值：无
- 必填：建议填写
- 作用：报告目录

说明：当前代码会创建这些目录，但真正落盘的任务中间文件主要在 `workdir.artifacts` 下；报告正文同时保存到数据库。

## frontend

### `frontend.enabled`

- 类型：bool
- 默认值：`false` 的零值会关闭托管；示例配置里建议设为 `true`
- 必填：建议填写
- 作用：是否启用后端静态托管前端

### `frontend.dist_dir`

- 类型：string
- 默认值：`./web/dist`
- 必填：否
- 作用：后端读取前端 dist 的目录

推荐值：

```yaml
frontend:
  enabled: true
  dist_dir: "./web/dist"
```

说明：

- 开发/部署推荐都把前端构建产物同步到这个目录
- 开启后，后端会同时处理静态资源和 SPA history fallback

## 不同环境配置建议

### 开发环境

- `auth.mode: mixed`
- `enable_register: true`
- 保留少量 `config_users`
- `opencode.binary_path: opencode`
- `workdir.*` 使用项目内相对路径

### 测试环境

- 使用独立数据库
- 更换 `jwt_secret`
- 配置固定仓库缓存目录
- 明确 `attach_url` 是否启用

### 生产环境

- 不建议继续使用默认 `config_users` 示例账号
- 必须更换 `jwt_secret`
- 建议把 `workdir` 放到持久化磁盘
- 为 `local_cache_dir` 和 `workdir` 设置可控权限

## 常见配置错误及排查

- `load config failed`：配置文件路径错误或 YAML 语法错误
- `database.dsn is required for mysql`：MySQL 场景缺少 `database.dsn`
- `unsupported database.type`：配置值不是 `mysql` / `dameng`
- 达梦连接失败：检查 `database.dsn` 或 `database.dameng.*` 是否完整，特别是 `port` 是否可转整数
- OpenCode 调用失败：确认 `binary_path` 可执行、`attach_url` 可访问
