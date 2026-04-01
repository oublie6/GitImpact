// cli_analyzer_test.go 验证 OpenCode CLI 参数拼装是否符合配置预期。
package analyzer

import (
	"strings"
	"testing"

	"gitimpact/backend/internal/config"
)

func TestCLIArgsBuild(t *testing.T) {
	c := NewCLIAnalyzer(config.OpenCodeConfig{AttachURL: "http://127.0.0.1:4096", DefaultModel: "m", DefaultAgent: "a"})
	args := c.buildArgs("analysis_prompt.md", true)
	line := strings.Join(args, " ")
	for _, token := range []string{"--attach", "http://127.0.0.1:4096", "--model", "m", "--agent", "a", "--output", "json"} {
		if !strings.Contains(line, token) {
			t.Fatalf("missing %s", token)
		}
	}
}
