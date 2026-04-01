// Package analyzer 定义影响分析器抽象与具体实现。
package analyzer

import "context"

// Analyzer 定义影响分析接口，支持 Markdown 与 Structured 报告。
type Analyzer interface {
	RunMarkdownReport(ctx context.Context, workDir string) (stdout string, stderr string, err error)
	RunStructuredReport(ctx context.Context, workDir string) (stdout string, stderr string, err error)
}
