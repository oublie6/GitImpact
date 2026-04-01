// report.go 定义结构化报告推荐格式及其解析方法。
package utils

import "encoding/json"

// StructuredReport 为结构化报告推荐字段。
type StructuredReport struct {
	Summary                 string   `json:"summary"`
	ChangedModules          []string `json:"changed_modules"`
	ImpactedInterfaces      []string `json:"impacted_interfaces"`
	ImpactedConfigs         []string `json:"impacted_configs"`
	ImpactedScripts         []string `json:"impacted_scripts"`
	ImpactedTests           []string `json:"impacted_tests"`
	Risks                   []string `json:"risks"`
	BackwardCompatibility   string   `json:"backward_compatibility"`
	DeploymentRisks         []string `json:"deployment_risks"`
	RollbackRisks           []string `json:"rollback_risks"`
	VerificationSuggestions []string `json:"verification_suggestions"`
	Confidence              string   `json:"confidence"`
	RawNotes                string   `json:"raw_notes"`
}

// ParseStructuredJSON 解析结构化报告 JSON 字符串。
// 如果 OpenCode 输出的字段不完整，解析仍可能成功，只是缺失字段保持零值。
func ParseStructuredJSON(s string) (*StructuredReport, error) {
	r := &StructuredReport{}
	if err := json.Unmarshal([]byte(s), r); err != nil {
		return nil, err
	}
	return r, nil
}
