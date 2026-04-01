# 开发指南

## 本地开发准备

需要的软件：

- Go 1.22+
- Node.js 18+
- npm 9+
- Git
- MySQL 8+ 或达梦
- OpenCode CLI
- Bash 环境，用于执行 `scripts/*.sh`

## 初始化

```bash
./scripts/init-dev.sh
```

这个脚本会：

- 为脚本赋执行权限
- 复制 `backend/config.example.yaml` 为 `backend/config.yaml`

## 数据库初始化

MySQL：

```bash
mysql -uroot -proot < sql/mysql/init.sql
```

达梦：

- 在目标 schema 中执行 `sql/dameng/init.sql`

## 后端启动

```bash
./scripts/dev-backend.sh
```

或手动启动：

```bash
cd backend
GITIMPACT_CONFIG=./config.yaml go run ./cmd/server
```

## 前端启动

```bash
./scripts/dev-frontend.sh
```

说明：当前脚本每次都会先执行 `npm install`，如果你只想快速启动，可以在 `frontend/` 目录手动执行 `npm run dev`。

## 常用脚本说明

- `scripts/init-dev.sh`：初始化本地开发环境
- `scripts/dev-backend.sh`：启动后端
- `scripts/dev-frontend.sh`：安装依赖并启动前端

## Makefile 常用命令说明

- `make build`：构建后端二进制到 `bin/gitimpact-backend`
- `make test`：执行后端测试
- `make build-linux-amd64`：交叉编译 Linux AMD64 二进制
- `make clean`：清理构建目录
- `make docker-build`：构建后端镜像

默认都带：

```makefile
GOFLAGS=-mod=vendor
```

## 从零到跑起来

1. 初始化配置文件。
2. 初始化数据库。
3. 修改 `backend/config.yaml`。
4. 启动后端。
5. 使用 `config_users` 或数据库用户登录。
6. 创建仓库记录。
7. 执行仓库抓取。
8. 创建分析任务。
9. 查看任务详情、日志和报告。
10. 下载报告。

## 最小演示路径

建议使用仓库自带配置用户快速验证：

1. 保持 `auth.mode: mixed`
2. 使用 `config_admin / Admin@123456` 登录
3. 创建一个仓库记录
4. 调用仓库抓取
5. 用 [examples/sample-task.json](../examples/sample-task.json) 的结构创建任务
6. 查看任务结果

## curl 示例

登录：

```bash
curl -X POST http://127.0.0.1:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"config_admin","password":"Admin@123456"}'
```

创建仓库：

```bash
curl -X POST http://127.0.0.1:8080/api/repositories \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name":"gitimpact-self",
    "repo_url":"https://github.com/example/gitimpact.git",
    "default_branch":"main",
    "local_cache_dir":"./runtime/repos/gitimpact-self",
    "auth_note":"readonly"
  }'
```

创建任务：

```bash
curl -X POST http://127.0.0.1:8080/api/tasks \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d @examples/sample-task.json
```

## 如何新增接口

1. 在 `internal/handler` 新增 handler 方法。
2. 在 `internal/service` 增加业务逻辑。
3. 需要数据库访问时，在 `internal/repository` 增加方法。
4. 在 `internal/router/router.go` 注册路由。
5. 为 service 或 handler 增加测试。
6. 更新 [API 参考](./api-reference.md)。

## 如何新增业务字段

1. 在 `internal/model/models.go` 增加字段。
2. 确认 SQL 初始化脚本同步更新。
3. 确认 API 请求/响应需要哪些字段。
4. 如有前端页面，补充表单或展示逻辑。
5. 更新文档。

## 如何新增分析器

1. 实现 `Analyzer` 接口。
2. 在 `main.go` 中装配新的 analyzer。
3. 确保 worker 的 `workDir` 约定仍然成立。
4. 增加测试或最小验证样例。

## 如何调试任务链路

- 查看任务状态：`GET /api/tasks/:id`
- 查看任务日志：`GET /api/tasks/:id/logs`
- 查看报告：`GET /api/tasks/:id/report`
- 查看中间产物目录：`workdir.artifacts/task_<id>_<ts>`

重点排查点：

- 仓库缓存目录是否存在 `.git`
- `old_ref/new_ref` 是否能被 checkout
- OpenCode CLI 是否在 PATH 中
- 结构化 JSON 是否为空

## 如何查看日志和报告

- 任务状态和错误信息：数据库 `analysis_tasks`
- 简要过程日志：数据库 `task_logs`
- 报告正文：数据库 `analysis_reports`
- 中间文件：`runtime/artifacts/`

## 前端页面操作路径

- `/login`：登录
- `/register`：注册
- `/`：首页
- `/repositories`：仓库列表
- `/tasks/new`：新建任务
- `/tasks`：任务列表
- `/tasks/:id`：任务详情
- `/settings`：系统设置
