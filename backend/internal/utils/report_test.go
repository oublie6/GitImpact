// report_test.go 验证结构化报告 JSON 的解析成功与失败分支。
package utils

import "testing"

func TestParseStructuredJSON(t *testing.T) {
	ok := `{"summary":"s","changed_modules":["a"]}`
	if _, err := ParseStructuredJSON(ok); err != nil {
		t.Fatal(err)
	}
	if _, err := ParseStructuredJSON("not-json"); err == nil {
		t.Fatal("expected parse error")
	}
}
