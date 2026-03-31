package analyzer

import "context"

// ServerAnalyzer 占位实现：用于后续接入 OpenCode Server/SDK。
type ServerAnalyzer struct{}

func (s *ServerAnalyzer) RunMarkdownReport(ctx context.Context, workDir string) (string, string, error) {
	return "", "", nil
}
func (s *ServerAnalyzer) RunStructuredReport(ctx context.Context, workDir string) (string, string, error) {
	return "", "", nil
}
