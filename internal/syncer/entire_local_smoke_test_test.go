package syncer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEnvBoolDefaultInvalidValueFallsBackToDefault(t *testing.T) {
	t.Setenv("GITSYNC_TEST_INVALID_BOOL", "not-a-bool")

	if got := envBoolDefault("GITSYNC_TEST_INVALID_BOOL", true); got != true {
		t.Fatalf("expected invalid value to fall back to true default, got %v", got)
	}
	if got := envBoolDefault("GITSYNC_TEST_INVALID_BOOL", false); got != false {
		t.Fatalf("expected invalid value to fall back to false default, got %v", got)
	}
}

func TestLoadEntireHostsWrappedFormat(t *testing.T) {
	configDir := t.TempDir()
	t.Setenv("ENTIRE_CONFIG_DIR", configDir)

	data := []byte(`{
  "activeHost": "127.0.0.1:8080",
  "hosts": {
    "127.0.0.1:8080": {
      "activeUser": "soph",
      "users": ["soph"]
    }
  }
}`)
	if err := os.WriteFile(filepath.Join(configDir, "hosts.json"), data, 0o600); err != nil {
		t.Fatalf("write hosts.json: %v", err)
	}

	hosts, activeHost, err := loadEntireHosts()
	if err != nil {
		t.Fatalf("load hosts: %v", err)
	}
	if activeHost != "127.0.0.1:8080" {
		t.Fatalf("unexpected active host: %q", activeHost)
	}
	if hosts["127.0.0.1:8080"].ActiveUser != "soph" {
		t.Fatalf("unexpected active user: %+v", hosts["127.0.0.1:8080"])
	}
}

func TestLoadEntireHostsLegacyMapFormat(t *testing.T) {
	configDir := t.TempDir()
	t.Setenv("ENTIRE_CONFIG_DIR", configDir)

	data := []byte(`{
  "entire.io": {
    "activeUser": "bryan",
    "users": ["bryan"]
  }
}`)
	if err := os.WriteFile(filepath.Join(configDir, "hosts.json"), data, 0o600); err != nil {
		t.Fatalf("write hosts.json: %v", err)
	}

	hosts, activeHost, err := loadEntireHosts()
	if err != nil {
		t.Fatalf("load hosts: %v", err)
	}
	if activeHost != "" {
		t.Fatalf("unexpected active host: %q", activeHost)
	}
	if hosts["entire.io"].ActiveUser != "bryan" {
		t.Fatalf("unexpected active user: %+v", hosts["entire.io"])
	}
}
