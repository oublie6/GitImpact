# OpenCode 集成说明

## OpenCode 在本项目中的定位

OpenCode 是 GitImpact 的分析执行引擎。GitImpact 本身负责：

- 收集 Git 差异材料
- 构造提示词
- 管理任务与报告存档

OpenCode 负责：

- 根据材料输出 Markdown 报告
- 根据材料输出结构化 JSON 报告

## 当前调用模式

### CLI 模式

这是当前唯一真正投入运行的模式，对应 `internal/analyzer/cli_analyzer.go`。

执行命令基础形态：

```bash
opencode run --prompt-file analysis_prompt.md
```

结构化报告时：

```bash
opencode run --output json --prompt-file analysis_prompt_json.md
```

当配置存在时，还会附加：

- `--attach <attach_url>`
- `--model <default_model>`
- `--agent <default_agent>`

### Attach 模式

严格来说它仍然属于 CLI 模式的一部分，只是通过 `--attach` 接到外部 OpenCode 服务端点。

### Server 模式

`ServerAnalyzer` 已存在，但目前是占位实现，不会真正产生报告。

## 调用输入材料

worker 在工作目录中准备以下输入：

- `changed_files.txt`
- `diff.patch`
- `commit_log.txt`
- `repo_manifest.md`
- `analysis_prompt.md`
- `analysis_prompt_json.md`

其中提示词明确要求 OpenCode：

- Markdown 场景：输出影响分析报告
- JSON 场景：输出 `summary`、`changed_modules`、`impacted_interfaces` 等字段

## 输出报告格式

### Markdown 报告

- 直接作为字符串保存在 `analysis_reports.markdown_report`
- 下载接口默认返回 Markdown

### 结构化 JSON 报告

- 原样保存在 `analysis_reports.structured_report`
- 当前代码只提供保存和返回，不自动做 schema 校验

项目中约定的推荐字段结构见 `internal/utils/report.go`。

## 异常处理与降级逻辑

- OpenCode 命令执行失败时：
  - Markdown 阶段会让任务整体失败
  - 但 worker 仍会尽量保存原始 stdout/stderr

- 结构化 JSON 阶段失败时：
  - 当前实现不会让任务整体失败
  - 因此可能出现 Markdown 有结果、JSON 为空的情况

## 如何替换或扩展分析器实现

分析器接口：

```go
type Analyzer interface {
    RunMarkdownReport(ctx context.Context, workDir string) (stdout string, stderr string, err error)
    RunStructuredReport(ctx context.Context, workDir string) (stdout string, stderr string, err error)
}
```

扩展步骤建议：

1. 在 `internal/analyzer` 中新增实现文件。
2. 复用已有 `workDir` 和材料输入约定。
3. 在 `main.go` 中替换 `analyzer.NewCLIAnalyzer(...)`。
4. 保持 `stdout/stderr/err` 的语义一致，避免影响 worker。

## 维护建议

- 把 prompt 模板从 worker 中抽离成独立模板文件。
- 对结构化 JSON 增加校验与解析反馈。
- 为 OpenCode 执行过程增加更细的日志。
- 区分“分析器执行失败”和“输出格式不符合预期”两类错误。
