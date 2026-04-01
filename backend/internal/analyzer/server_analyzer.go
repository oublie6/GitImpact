// server_analyzer.go 预留给未来的 OpenCode Server/SDK 集成。
package analyzer

import "context"

// ServerAnalyzer 占位实现：用于后续接入 OpenCode Server/SDK。
type ServerAnalyzer struct{}

// RunMarkdownReport 当前尚未实现，返回空结果。
func (s *ServerAnalyzer) RunMarkdownReport(ctx context.Context, workDir string) (string, string, error) {
	return "", "", nil
}

// RunStructuredReport 当前尚未实现，返回空结果。
func (s *ServerAnalyzer) RunStructuredReport(ctx context.Context, workDir string) (string, string, error) {
	return "", "", nil
}
