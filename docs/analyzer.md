# Analyzer 说明
- 接口：`Analyzer`
  - `RunMarkdownReport(task)`
  - `RunStructuredReport(task)`
- 当前实现：
  - `CLIAnalyzer`：调用 `opencode run`。
  - `ServerAnalyzer`：**占位实现**。
- 材料准备由系统完成，不允许 OpenCode 自行拉仓库。
