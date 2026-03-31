# API 列表
- 健康检查：`GET /healthz`
- 认证：`POST /api/auth/register`、`POST /api/auth/login`、`POST /api/auth/logout`、`GET /api/auth/me`
- 仓库：`POST/PUT/GET /api/repositories`，`POST /api/repositories/:id/fetch`
- 任务：`POST /api/tasks`、`GET /api/tasks`、`GET /api/tasks/:id`、`GET /api/tasks/:id/logs`、`GET /api/tasks/:id/report`、`GET /api/tasks/:id/report/download?format=md|json`
- 系统设置：`GET /api/settings`、`POST /api/settings`

统一返回：`{code,message,data}`。
