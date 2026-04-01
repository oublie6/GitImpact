// cli_analyzer.go 使用 OpenCode CLI 执行影响分析。
//
// 它读取 worker 预先写入的 prompt 文件与任务材料，不直接访问数据库。
package analyzer

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"gitimpact/backend/internal/config"
)

// CLIAnalyzer 使用 opencode run 命令执行分析。
type CLIAnalyzer struct{ cfg config.OpenCodeConfig }

// NewCLIAnalyzer 创建基于 CLI 的分析器实现。
func NewCLIAnalyzer(cfg config.OpenCodeConfig) *CLIAnalyzer { return &CLIAnalyzer{cfg: cfg} }

// buildArgs 根据配置拼装 opencode run 参数。
// jsonMode 为 true 时会额外要求 CLI 输出 JSON。
func (c *CLIAnalyzer) buildArgs(promptFile string, jsonMode bool) []string {
	args := []string{"run"}
	if c.cfg.AttachURL != "" {
		args = append(args, "--attach", c.cfg.AttachURL)
	}
	if c.cfg.DefaultModel != "" {
		args = append(args, "--model", c.cfg.DefaultModel)
	}
	if c.cfg.DefaultAgent != "" {
		args = append(args, "--agent", c.cfg.DefaultAgent)
	}
	if jsonMode {
		args = append(args, "--output", "json")
	}
	args = append(args, "--prompt-file", promptFile)
	return args
}

// run 在指定工作目录下执行 OpenCode CLI。
// 返回值会把 stdout/stderr 分开暴露，便于 worker 原样存档。
func (c *CLIAnalyzer) run(ctx context.Context, workDir, promptFile string, jsonMode bool) (string, string, error) {
	bin := c.cfg.BinaryPath
	if bin == "" {
		bin = "opencode"
	}
	t, cancel := context.WithTimeout(ctx, time.Duration(c.cfg.TimeoutSeconds)*time.Second)
	defer cancel()
	cmd := exec.CommandContext(t, bin, c.buildArgs(promptFile, jsonMode)...)
	cmd.Dir = workDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", string(out), fmt.Errorf("opencode run failed: %w", err)
	}
	return string(out), "", nil
}

// RunMarkdownReport 执行 Markdown 报告生成。
func (c *CLIAnalyzer) RunMarkdownReport(ctx context.Context, workDir string) (string, string, error) {
	return c.run(ctx, workDir, "analysis_prompt.md", false)
}

// RunStructuredReport 执行结构化 JSON 报告生成。
func (c *CLIAnalyzer) RunStructuredReport(ctx context.Context, workDir string) (string, string, error) {
	return c.run(ctx, workDir, "analysis_prompt_json.md", true)
}

// BuildCommandPreview 用于测试或日志中展示即将执行的命令。
func BuildCommandPreview(bin string, args []string) string {
	return strings.TrimSpace(bin + " " + strings.Join(args, " "))
}
