# API 参考

统一响应结构：

```json
{
  "code": 0,
  "message": "ok",
  "data": {}
}
```

失败时 `code` 等于 HTTP 状态码，`message` 为错误信息。

## 健康检查

### `GET /healthz`

- 是否需要鉴权：否
- 作用：确认服务存活

响应示例：

```json
{
  "status": "ok"
}
```

## 认证接口

### `POST /api/auth/register`

- 是否需要鉴权：否
- 作用：注册数据库用户
- 请求体：

```json
{
  "username": "alice",
  "password": "Secret@123",
  "email": "alice@example.com"
}
```

- 成功响应：

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "username": "alice"
  }
}
```

- 典型错误场景：
  - `400`：参数不合法、注册关闭、`auth.mode=config`
  - `409`：用户名已存在

### `POST /api/auth/login`

- 是否需要鉴权：否
- 作用：登录并获取 JWT
- 请求体：

```json
{
  "username": "config_admin",
  "password": "Admin@123456"
}
```

- 成功响应：

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "token": "jwt-token",
    "user": {
      "id": 0,
      "username": "config_admin",
      "display_name": "Config Admin",
      "email": "config_admin@example.com",
      "role": "admin",
      "status": "active",
      "source": "config"
    }
  }
}
```

说明：

- 配置用户登录成功时，`user.id` 会是 `0`。
- `mixed` 模式登录顺序是先 `config_users`，再数据库用户。

- 典型错误场景：
  - `401`：用户名或密码错误、用户停用

### `POST /api/auth/logout`

- 是否需要鉴权：是
- 作用：返回退出确认
- 请求头：`Authorization: Bearer <token>`

成功响应：

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "message": "logout success"
  }
}
```

说明：当前未实现服务端 token 黑名单。

### `GET /api/auth/me`

- 是否需要鉴权：是
- 作用：返回 JWT 中的当前用户信息

成功响应：

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "username": "config_admin",
    "user_id": 0,
    "role": "admin"
  }
}
```

## 仓库接口

### `POST /api/repositories`

- 是否需要鉴权：是
- 作用：创建仓库记录
- 请求体字段：
  - `name`
  - `repo_url`
  - `default_branch`
  - `local_cache_dir`
  - `auth_note`

请求示例：

```json
{
  "name": "gitimpact-self",
  "repo_url": "https://github.com/example/gitimpact.git",
  "default_branch": "main",
  "local_cache_dir": "./runtime/repos/gitimpact-self",
  "auth_note": "使用只读部署凭证"
}
```

- 典型错误场景：
  - `400`：参数绑定失败、数据库唯一键冲突

### `PUT /api/repositories/:id`

- 是否需要鉴权：是
- 作用：更新仓库记录
- 路径参数：`id`
- 请求体：与创建接口相同

### `GET /api/repositories`

- 是否需要鉴权：是
- 作用：查询仓库列表

### `GET /api/repositories/:id`

- 是否需要鉴权：是
- 作用：查询单个仓库详情
- 典型错误场景：
  - `404`：仓库不存在

### `POST /api/repositories/:id/fetch`

- 是否需要鉴权：是
- 作用：执行本地缓存同步
- 实际行为：
  - 若 `local_cache_dir/.git` 不存在，则先 `git clone`
  - 然后执行 `git fetch --all --prune`

- 典型错误场景：
  - `404`：仓库不存在
  - `500`：git 命令执行失败

## 任务接口

### `POST /api/tasks`

- 是否需要鉴权：是
- 作用：创建分析任务并异步入队
- 请求体字段：

```json
{
  "task_name": "核心接口变更影响分析",
  "mode": "same_repo_commits",
  "old_repo_id": 1,
  "old_ref": "main~3",
  "new_repo_id": 1,
  "new_ref": "main",
  "generate_markdown": true,
  "generate_structured": true,
  "custom_focus": "重点关注登录鉴权和部署脚本",
  "remark": "发布前评估"
}
```

- 成功响应：返回创建后的任务对象
- 实际行为：
  - 服务端会把 `status` 强制设为 `pending`
  - `created_by` 来自 JWT 中的 `user_id`
  - 返回成功不代表任务执行完成

- 典型错误场景：
  - `400`：参数绑定失败或入库失败

### `GET /api/tasks`

- 是否需要鉴权：是
- 作用：返回任务列表

### `GET /api/tasks/:id`

- 是否需要鉴权：是
- 作用：返回任务详情
- 典型错误场景：
  - `404`：任务不存在

### `GET /api/tasks/:id/logs`

- 是否需要鉴权：是
- 作用：返回任务日志列表

日志对象字段：

- `id`
- `task_id`
- `level`
- `message`
- `created_at`

### `GET /api/tasks/:id/report`

- 是否需要鉴权：是
- 作用：返回任务报告

报告对象字段：

- `task_id`
- `markdown_report`
- `structured_report`
- `raw_stdout`
- `raw_stderr`

### `GET /api/tasks/:id/report/download`

- 是否需要鉴权：是
- 作用：下载报告附件
- 查询参数：
  - `format=md`，默认值
  - `format=json`

实际行为：

- `format=json` 返回 `application/json`
- 其他值按 Markdown 返回 `text/markdown`

## 系统设置接口

### `GET /api/settings`

- 是否需要鉴权：是
- 作用：返回系统设置列表

### `POST /api/settings`

- 是否需要鉴权：是
- 作用：保存单个系统设置
- 请求体：

```json
{
  "key": "opencode.default_model",
  "value": "gpt-4.1"
}
```

- 典型错误场景：
  - `400`：参数绑定失败或保存失败

## curl 最小调用示例

```bash
TOKEN=$(curl -s -X POST http://127.0.0.1:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"config_admin","password":"Admin@123456"}' | jq -r '.data.token')

curl -X GET http://127.0.0.1:8080/api/tasks \
  -H "Authorization: Bearer $TOKEN"
```
