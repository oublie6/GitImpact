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

func NewCLIAnalyzer(cfg config.OpenCodeConfig) *CLIAnalyzer { return &CLIAnalyzer{cfg: cfg} }

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

func (c *CLIAnalyzer) RunMarkdownReport(ctx context.Context, workDir string) (string, string, error) {
	return c.run(ctx, workDir, "analysis_prompt.md", false)
}
func (c *CLIAnalyzer) RunStructuredReport(ctx context.Context, workDir string) (string, string, error) {
	return c.run(ctx, workDir, "analysis_prompt_json.md", true)
}

func BuildCommandPreview(bin string, args []string) string {
	return strings.TrimSpace(bin + " " + strings.Join(args, " "))
}
