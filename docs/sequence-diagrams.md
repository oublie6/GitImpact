# 时序图

## 用户登录时序图

```mermaid
sequenceDiagram
    participant User as 用户
    participant Frontend as 前端
    participant Handler as AuthHandler
    participant Service as AuthService
    participant Repo as UserRepository
    participant JWT as JWT 工具

    User->>Frontend: 输入用户名和密码
    Frontend->>Handler: POST /api/auth/login
    Handler->>Service: ValidateLogin
    Service->>Service: 先校验 config_users
    Service->>Repo: mixed/db 模式下查数据库用户
    Service->>JWT: GenerateToken
    JWT-->>Service: token
    Service-->>Handler: user + token
    Handler-->>Frontend: 返回 token 和 user
```

说明：

- `mixed` 模式下会先查 `config_users`，再查数据库。
- 配置用户登录成功时，`user_id` 为 `0`。

## 创建仓库记录时序图

```mermaid
sequenceDiagram
    participant User as 用户
    participant Frontend as 前端
    participant Handler as RepoHandler
    participant Service as RepositoryService
    participant Repo as RepoRepository
    participant DB as 数据库

    User->>Frontend: 填写仓库信息
    Frontend->>Handler: POST /api/repositories
    Handler->>Service: Create
    Service->>Repo: Create
    Repo->>DB: INSERT repositories
    DB-->>Repo: 成功
    Repo-->>Service: 成功
    Service-->>Handler: 成功
    Handler-->>Frontend: 返回仓库对象
```

说明：创建仓库记录只写数据库，不会自动执行 `git clone`。

## 创建分析任务时序图

```mermaid
sequenceDiagram
    participant User as 用户
    participant Frontend as 前端
    participant Handler as TaskHandler
    participant Service as TaskService
    participant Repo as TaskRepository
    participant Worker as TaskWorker

    User->>Frontend: 提交任务
    Frontend->>Handler: POST /api/tasks
    Handler->>Service: CreateTask
    Service->>Repo: INSERT analysis_tasks(status=pending)
    Repo-->>Service: 成功
    Service-->>Handler: 任务对象
    Handler->>Worker: Enqueue(taskID)
    Handler-->>Frontend: 返回 pending 任务
```

说明：HTTP 返回成功后，任务才开始进入异步执行阶段。

## Worker 执行任务时序图

```mermaid
sequenceDiagram
    participant Worker as TaskWorker
    participant TaskRepo as TaskRepository
    participant RepoRepo as RepoRepository
    participant Git as Git 命令
    participant Analyzer as CLIAnalyzer
    participant ReportRepo as ReportRepository

    Worker->>TaskRepo: GetByID
    Worker->>TaskRepo: Update(status=running)
    Worker->>TaskRepo: AddLog(task started)
    Worker->>RepoRepo: Get old/new repository
    Worker->>Git: clone/fetch/checkout
    Worker->>Worker: writeMaterials
    Worker->>Analyzer: RunMarkdownReport
    Worker->>Analyzer: RunStructuredReport
    Worker->>ReportRepo: Upsert(report)
    Worker->>TaskRepo: Update(status=success)
    Worker->>TaskRepo: AddLog(task completed)
```

说明：如果 Markdown 生成失败，任务会直接失败；结构化 JSON 失败不会让任务整体失败。

## OpenCode 分析调用时序图

```mermaid
sequenceDiagram
    participant Worker as TaskWorker
    participant CLI as CLIAnalyzer
    participant OpenCode as opencode run

    Worker->>CLI: RunMarkdownReport(workDir)
    CLI->>CLI: buildArgs
    CLI->>OpenCode: exec.CommandContext
    OpenCode-->>CLI: stdout/stderr
    CLI-->>Worker: markdown result
    Worker->>CLI: RunStructuredReport(workDir)
    CLI->>OpenCode: opencode run --output json
    OpenCode-->>CLI: stdout/stderr
    CLI-->>Worker: json result
```

说明：工作目录内必须已有 prompt 和分析材料，CLIAnalyzer 不会自行生成这些文件。

## 下载报告时序图

```mermaid
sequenceDiagram
    participant User as 用户
    participant Frontend as 前端
    participant Handler as TaskHandler
    participant Service as TaskService
    participant Repo as ReportRepository

    User->>Frontend: 点击下载报告
    Frontend->>Handler: GET /api/tasks/:id/report/download
    Handler->>Service: GetReport
    Service->>Repo: GetByTaskID
    Repo-->>Service: report
    Service-->>Handler: report
    Handler-->>Frontend: markdown 或 json 附件流
```

说明：下载接口直接从数据库读取报告正文，不依赖磁盘文件。
