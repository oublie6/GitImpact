# FAQ

## 1. 为什么我登录成功后 `user_id` 是 `0`？

因为你登录的是 `config_users` 中的静态用户，这类用户不写数据库，JWT 中的 `user_id` 会保持为 `0`。

## 2. 为什么注册接口可用，但有时又提示关闭？

注册是否可用同时取决于：

- `auth.enable_register`
- `auth.mode` 不能是 `config`

## 3. 任务创建成功为什么没有立刻看到报告？

`POST /api/tasks` 只表示任务入库成功，真正执行发生在后台 goroutine 中。

## 4. 为什么前端任务创建页默认仓库 ID 是 1？

这是当前原型页面的最小演示写法，不代表系统只支持仓库 1。

## 5. 当前支持哪些数据库？

运行时支持 MySQL 和达梦。SQLite 只在测试里使用。

## 6. 当前支持哪些 OpenCode 集成方式？

真正可用的是 CLIAnalyzer。ServerAnalyzer 还是占位实现。

## 7. Docker 构建时会联网拉 Go 依赖吗？

按当前 Dockerfile，不会。构建走 `GOFLAGS=-mod=vendor`。

## 8. 可以并发执行多个任务吗？

代码层面可以同时起多个 goroutine，但由于同一仓库缓存会被 `git checkout`，并发任务存在互相影响的风险。

## 9. 为什么登出后仍然能用旧 token？

当前没有服务端 token 吊销或黑名单机制，登出只清理前端本地状态。

## 10. 为什么现在默认禁用 AutoMigrate？

因为在达梦环境中，AutoMigrate 可能触发兼容性问题（例如 MySQL 类型透传、重复加列等），会降低部署稳定性。

当前策略是：

- 默认关闭 `database.auto_migrate`
- 统一由手工 SQL 初始化表结构
  - MySQL：`sql/mysql/init.sql`
  - 达梦：`sql/dameng/init.sql`

仅建议在本地临时调试时显式开启 AutoMigrate，且不推荐在达梦环境开启。
