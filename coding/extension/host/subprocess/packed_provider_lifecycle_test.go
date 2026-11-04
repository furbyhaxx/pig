package subprocess

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/MichaelKinsy/PiG/coding/extension"
)

// writePackedFactoryModuleWithProvider scaffolds a factory-style Go module
// that registers exactly one provider in addition to a no-op tool. Used for
// packed-cell provider lifecycle tests.
func writePackedFactoryModuleWithProvider(t *testing.T, modulePath, extName, toolName, providerName string) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), fmt.Appendf(nil, "module %s\n\ngo 1.26\n\nrequire github.com/MichaelKinsy/PiG/extensions/sdk v0.0.0\n", modulePath), 0o644); err != nil {
		t.Fatal(err)
	}
	src := fmt.Sprintf(`package ext

import "github.com/MichaelKinsy/PiG/extensions/sdk"

func Extension() *sdk.Extension {
	e := sdk.New(%q)
	e.Tool(%q, "noop", sdk.Schema{"type": "object"}, func(ctx sdk.Context, params map[string]any) (any, error) {
		return map[string]any{"content": "ok"}, nil
	})
	e.RegisterProvider(%q, sdk.ProviderConfig{"baseUrl": "http://localhost"})
	return e
}
`, extName, toolName, providerName)
	if err := os.WriteFile(filepath.Join(dir, "ext.go"), []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	return dir
}

func TestReloadPackedCellFailureKeepsServingGeneration(t *testing.T) {
	marker := filepath.Join(t.TempDir(), "stall")
	rootA := writePackedFactoryModuleWithProvider(t, "example.com/rollback/a", "rollback-a", "tool_a", "rollback-provider-a")
	rootB := writePackedFactoryModuleWithProvider(t, "example.com/rollback/b", "rollback-b", "tool_b", "rollback-provider-b")
	apiPath := filepath.Join(rootA, "ext.go")
	apiSource, err := os.ReadFile(apiPath)
	if err != nil {
		t.Fatal(err)
	}
	apiText := strings.Replace(string(apiSource), `import "github.com/MichaelKinsy/PiG/extensions/sdk"`, `import ("os"; "github.com/MichaelKinsy/PiG/extensions/sdk")`, 1)
	apiText = strings.Replace(apiText, `e := sdk.New("rollback-a")`, fmt.Sprintf(`e := sdk.New("rollback-a")
	url, name := "https://mcp.invalid", "Route"
	if _, err := os.Stat(%q); err == nil { url, name = "https://replacement.invalid", "Replacement" }
	if err := e.RegisterMcpServer("rollback-mcp", sdk.McpServerConfig{Type: "http", URL: url}); err != nil { panic(err) }
	if err := e.RegisterVirtualModel(sdk.VirtualModel{Provider: "rollback", ID: "route", Name: name, Route: func(sdk.Context, sdk.ModelRouteRequest) (sdk.ModelRoute, error) { return sdk.ModelRoute{}, nil }}); err != nil { panic(err) }`, marker), 1)
	if err := os.WriteFile(apiPath, []byte(apiText), 0o600); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(rootB, "ext.go")
	source, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	text := strings.Replace(string(source), `import "github.com/MichaelKinsy/PiG/extensions/sdk"`, `import ("os"; "time"; "github.com/MichaelKinsy/PiG/extensions/sdk")`, 1)
	text = strings.Replace(text, `e := sdk.New("rollback-b")`, fmt.Sprintf(`if _, err := os.Stat(%q); err == nil { time.Sleep(time.Hour) }; e := sdk.New("rollback-b")`, marker), 1)
	if err := os.WriteFile(path, []byte(text), 0o600); err != nil {
		t.Fatal(err)
	}
	configs := []ExtConfig{
		packedFactoryConfig("rollback-a", rootA, "example.com/rollback/a", "a"),
		packedFactoryConfig("rollback-b", rootB, "example.com/rollback/b", "b"),
	}
	h := NewHostWithConfigRoot(t.TempDir(), t.TempDir())
	t.Cleanup(func() { h.Shutdown("test done") })
	h.SetConfigLoader(func() ([]ExtConfig, error) { return configs, nil })
	h.SetReloadHandshakeTimeoutFunc(func() (time.Duration, error) { return 250 * time.Millisecond, nil })
	var mu sync.Mutex
	providers := map[string]bool{}
	h.SetProviderCallbacks(func(name string, _ extension.ProviderConfig) error {
		mu.Lock()
		defer mu.Unlock()
		providers[name] = true
		return nil
	}, func(name string) {
		mu.Lock()
		defer mu.Unlock()
		delete(providers, name)
	})
	virtualModels := map[string]extension.VirtualModelDefinition{}
	h.Runtime().BindProviderActions(extension.ProviderActions{
		RegisterVirtualModel: func(model extension.VirtualModelDefinition) error {
			mu.Lock()
			defer mu.Unlock()
			virtualModels[model.ID] = model
			return nil
		},
		UnregisterVirtualModel: func(_ string, id string) {
			mu.Lock()
			defer mu.Unlock()
			delete(virtualModels, id)
		},
	}, nil)
	ctx, cancel := buildContext(t)
	defer cancel()
	if _, errs := h.LoadAll(ctx, configs); len(errs) != 0 {
		t.Fatalf("initial load: %v", errs)
	}
	applyQueuedRegistrations(h)
	h.mu.Lock()
	oldA, oldB := h.exts["rollback-a"], h.exts["rollback-b"]
	h.mu.Unlock()
	if err := os.WriteFile(marker, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := h.Reload(ctx); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("reload error = %v, want handshake timeout", err)
	}
	h.mu.Lock()
	gotA, gotB := h.exts["rollback-a"], h.exts["rollback-b"]
	h.mu.Unlock()
	if gotA != oldA || gotB != oldB {
		t.Fatal("a failed member replaced a serving generation of its cell")
	}
	for _, old := range []*managedExt{oldA, oldB} {
		for _, tool := range old.ext.Tools {
			if _, err := tool.Definition.Execute(ctx, "after-rollback", json.RawMessage(`{}`), nil); err != nil {
				t.Fatalf("old generation stopped serving: %v", err)
			}
		}
	}
	applyQueuedRegistrations(h)
	mu.Lock()
	if len(providers) != 2 || !providers["rollback-provider-a"] || !providers["rollback-provider-b"] {
		t.Errorf("rollback changed the serving provider set: %v", providers)
	}
	if server, ok := h.Runtime().McpServers().Get("rollback-mcp"); !ok || server.Config.URL != "https://mcp.invalid" {
		t.Errorf("rollback changed the serving MCP server: %#v", server)
	}
	if model, ok := virtualModels["route"]; !ok || model.Name != "Route" {
		t.Errorf("rollback changed the serving virtual model: %#v", model)
	}
	mu.Unlock()
	h.quarantinePackedCellGeneration(oldA.packedCellKey, oldA.packedProcess.generation, "serving process crashed")
	if len(h.Extensions()) != 0 {
		t.Fatal("failed reload made the serving process's crash report stale")
	}
}

// TestPackedCell_ProviderLifecycle ensures that packed-cell members can
// register providers, that reload preserves them when packed artifacts are
// reused, and that quarantine unregisters every packed-member provider.
//
// pig-specific. Production-grade reload safety: globally-visible provider
// registrations from packed members must not leak after quarantine.
func TestPackedCell_ProviderLifecycle(t *testing.T) {
	rootA := writePackedFactoryModuleWithProvider(t, "example.com/provpack/a", "provpack-a", "tool_a", "prov-a")
	rootB := writePackedFactoryModuleWithProvider(t, "example.com/provpack/b", "provpack-b", "tool_b", "prov-b")
	configs := []ExtConfig{
		packedFactoryConfig("provpack-a", rootA, "example.com/provpack/a", "ha"),
		packedFactoryConfig("provpack-b", rootB, "example.com/provpack/b", "hb"),
	}

	var (
		mu         sync.Mutex
		registered = map[string]extension.ProviderConfig{}
	)
	h := NewHost(t.TempDir())
	h.SetProviderCallbacks(
		func(name string, cfg extension.ProviderConfig) error {
			mu.Lock()
			defer mu.Unlock()
			registered[name] = cfg
			return nil
		},
		func(name string) {
			mu.Lock()
			defer mu.Unlock()
			delete(registered, name)
		},
	)
	h.SetConfigLoader(func() ([]ExtConfig, error) { return configs, nil })
	t.Cleanup(func() { h.Shutdown("test done") })

	ctx, cancel := buildContext(t)
	defer cancel()
	if _, err := h.Reload(ctx); err != nil {
		t.Fatalf("first reload: %v", err)
	}

	check := func(want []string) {
		t.Helper()
		applyQueuedRegistrations(h)
		mu.Lock()
		got := make([]string, 0, len(registered))
		for k := range registered {
			got = append(got, k)
		}
		mu.Unlock()
		sort.Strings(got)
		sort.Strings(want)
		if len(got) != len(want) {
			t.Fatalf("provider set = %v, want %v", got, want)
		}
		for i, name := range want {
			if got[i] != name {
				t.Fatalf("provider set = %v, want %v", got, want)
			}
		}
	}
	check([]string{"prov-a", "prov-b"})

	// Reload with identical config: packed artifact is reused; provider set must
	// remain stable (no double-register, no temporary unregister).
	if _, err := h.Reload(ctx); err != nil {
		t.Fatalf("second reload: %v", err)
	}
	check([]string{"prov-a", "prov-b"})

	// Quarantine the packed cell → all packed-member providers must be
	// unregistered when the host stops the quarantined extensions.
	h.mu.Lock()
	packedKey := h.exts["provpack-a"].packedCellKey
	h.mu.Unlock()
	h.quarantinePackedCell(packedKey, "boom")
	check(nil)
}
