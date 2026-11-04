package codingagent

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

func TestExtensionsReloadHandshakeSetting(t *testing.T) {
	cwd, agentDir := t.TempDir(), t.TempDir()
	sm := NewSettingsManager(cwd, agentDir)
	if got, err := sm.GetExtensionsReloadHandshakeTimeout(); err != nil || got != 15*time.Second {
		t.Fatalf("default = %v, %v", got, err)
	}
	path := filepath.Join(agentDir, "settings.json")
	if err := os.WriteFile(path, []byte(`{"extensions":["./extension"],"extensionsReloadHandshakeMs":321}`), 0o600); err != nil {
		t.Fatal(err)
	}
	sm.Reload()
	if got, err := sm.GetExtensionsReloadHandshakeTimeout(); err != nil || got != 321*time.Millisecond {
		t.Fatalf("configured = %v, %v", got, err)
	}
	if err := sm.UpdateGlobal(func(s *Settings) { s.Theme = "light" }); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var written Settings
	if err := json.Unmarshal(data, &written); err != nil {
		t.Fatal(err)
	}
	if written.ExtensionsReloadHandshakeMs == nil || *written.ExtensionsReloadHandshakeMs != 321 || len(written.Extensions) != 1 || written.Extensions[0] != "./extension" {
		t.Fatalf("round trip lost the timeout or extensions array: %s", data)
	}
	copy := sm.GetGlobalSettings()
	*copy.ExtensionsReloadHandshakeMs = 999
	if got, err := sm.GetExtensionsReloadHandshakeTimeout(); err != nil || got != 321*time.Millisecond {
		t.Fatalf("caller mutated manager settings: %v, %v", got, err)
	}
}

func TestExtensionsReloadHandshakeIgnoresProjectSetting(t *testing.T) {
	cwd, agentDir := t.TempDir(), t.TempDir()
	if err := os.MkdirAll(filepath.Join(cwd, ".pig"), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(cwd, ".pig", "settings.json"), []byte(`{"extensionsReloadHandshakeMs":1}`), 0o600); err != nil {
		t.Fatal(err)
	}
	sm := NewSettingsManager(cwd, agentDir)
	sm.SetProjectTrusted(true)
	if got, err := sm.GetExtensionsReloadHandshakeTimeout(); err != nil || got != 15*time.Second {
		t.Fatalf("project overrode the global budget: %v, %v", got, err)
	}
}

func TestExtensionsReloadHandshakeRejectsInvalidDuration(t *testing.T) {
	values := []int{0, -1}
	if strconv.IntSize == 64 {
		values = append(values, int(^uint(0)>>1))
	}
	for _, value := range values {
		t.Run(strconv.Itoa(value), func(t *testing.T) {
			sm := NewSettingsManager(t.TempDir(), t.TempDir())
			if err := sm.UpdateGlobal(func(s *Settings) { s.ExtensionsReloadHandshakeMs = &value }); err != nil {
				t.Fatal(err)
			}
			if _, err := sm.GetExtensionsReloadHandshakeTimeout(); err == nil {
				t.Fatal("invalid duration accepted")
			}
		})
	}
}
