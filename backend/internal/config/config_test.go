// config_test.go 验证配置加载与默认值填充行为。
package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	tmp := t.TempDir() + "/cfg.yaml"
	_ = os.WriteFile(tmp, []byte("server:\n  port: '9000'\nauth:\n  mode: db\n  jwt_secret: x\ndatabase:\n  dsn: root:root@tcp(127.0.0.1:3306)/gitimpact\n"), 0o644)
	cfg, err := Load(tmp)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Server.Port != "9000" || cfg.Auth.Mode != "db" {
		t.Fatal("unexpected config")
	}
}
