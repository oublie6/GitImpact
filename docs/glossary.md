# 术语表

- `repository`：被纳入 GitImpact 管理的代码仓库记录，包含远端地址、本地缓存路径等信息。
- `task`：一次影响分析请求，描述 old/new 仓库与引用、输出类型和关注点。
- `report`：任务执行后的 Markdown 或结构化 JSON 结果。
- `analyzer`：负责调用 OpenCode 或其他分析引擎的抽象层。
- `worker`：后台异步执行任务的组件。
- `artifact`：任务执行过程中生成的中间材料文件。
- `attach`：OpenCode CLI 的连接参数，表示把命令附着到外部服务端点。
- `config user`：定义在 `config.yaml` 中的静态登录用户，不落库。
- `db user`：存储在 `users` 表中的用户。
- `mixed mode`：认证模式之一，登录时先查配置用户，再查数据库用户。
- `repo cache`：本地 Git 仓库缓存目录。
- `workdir`：运行时目录根，用于保存仓库缓存、任务材料和其他文件。
