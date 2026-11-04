package coding

import (
	"context"
	"encoding/json"
	"errors"
	"path/filepath"
	"testing"

	"github.com/MichaelKinsy/PiG/ai"
	"github.com/MichaelKinsy/PiG/coding/extension"
)

func TestSessionReloadPublishesPartialCellSuccessAndReturnsTheFailure(t *testing.T) {
	services := newTestServices(t)
	retained := registryTool("retained", "Retained", "Serving failed cell", "")
	healthyOld := registryTool("healthy_old", "Old", "Previous healthy cell", "")
	healthyNew := registryTool("healthy_new", "New", "Replacement healthy cell", "")
	partial := false
	cause := errors.New("another cell failed to register")
	runtime := extension.CreateExtensionRuntime()
	loader := NewDefaultResourceLoader(DefaultResourceLoaderOptions{
		CWD: services.CWD(), AgentDir: filepath.Join(t.TempDir(), "agent"), NoSkills: true, NoPromptTemplates: true, NoThemes: true, NoContextFiles: true,
		LoadExtensions: func(context.Context, ExtensionLoadRequest) (LoadExtensionsResult, error) {
			healthy, loadErr := healthyOld, error(nil)
			if partial {
				healthy = healthyNew
				loadErr = errors.Join(extension.ErrReloadPartiallyApplied, cause)
			}
			return LoadExtensionsResult{Runtime: runtime, Extensions: []extension.Extension{
				{Path: "<retained>", Tools: map[string]extension.RegisteredTool{retained.Name: {Definition: retained}}},
				{Path: "<healthy>", Tools: map[string]extension.RegisteredTool{healthy.Name: {Definition: healthy}}},
			}}, loadErr
		},
	})
	if err := loader.ReloadWith(t.Context()); err != nil {
		t.Fatal(err)
	}
	session, err := NewSession(services, SessionOptions{NoSession: true, ResourceLoader: loader})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = session.Close() })
	partial = true
	if err := session.Reload(t.Context()); !errors.Is(err, cause) {
		t.Fatalf("partial reload error = %v", err)
	}
	if _, found := session.GetToolDefinition(healthyOld.Name); found {
		t.Fatal("the replaced healthy cell's old tool stayed published")
	}
	for _, name := range []string{healthyNew.Name, retained.Name} {
		tool, found := session.GetToolDefinition(name)
		if !found {
			t.Fatalf("serving tool %s was not published", name)
		}
		if result, err := tool.Execute(t.Context(), "partial-probe", json.RawMessage(`{}`), nil); err != nil || result.Text() != "ok" {
			t.Fatalf("serving tool %s failed: %v, %v", name, result, err)
		}
	}
}

type reloadContextKey struct{}

// contextReloadLoader records the context its ReloadWith ran under.
type contextReloadLoader struct {
	ResourceLoader
	got context.Context
}

func (l *contextReloadLoader) ReloadWith(ctx context.Context, _ ...ResourceLoaderReloadOptions) error {
	l.got = ctx
	return nil
}

// Session.Reload reloads a loader that can take a context under the context of the reload call, so the caller's request values reach the resource reload (agent-session.ts reload awaits resourceLoader.reload for its caller).
func TestSessionReloadPassesItsContextToTheResourceLoader(t *testing.T) {
	t.Setenv("PIG_HOME", t.TempDir())
	services, err := CreateAgentSessionServices(CreateAgentSessionServicesOptions{CWD: t.TempDir()})
	if err != nil {
		t.Fatal(err)
	}
	loader := &contextReloadLoader{ResourceLoader: NoResources}
	session, err := NewSession(services, SessionOptions{ResourceLoader: loader})
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = session.Close() }()
	ctx := context.WithValue(context.Background(), reloadContextKey{}, "caller")
	if err := session.Reload(ctx); err != nil {
		t.Fatal(err)
	}
	if loader.got == nil || loader.got.Value(reloadContextKey{}) != "caller" {
		t.Fatalf("ReloadWith ran under %v, want the reload call's context", loader.got)
	}
}

// runtimeReloadLoader reports, after each reload, the extension runtime its extensions registered against, or fails the reload.
type runtimeReloadLoader struct {
	ResourceLoader
	runtime *extension.ExtensionRuntime
	err     error
}

func (l *runtimeReloadLoader) Reload(...ResourceLoaderReloadOptions) error { return l.err }

func (l *runtimeReloadLoader) GetExtensions() LoadExtensionsResult {
	return LoadExtensionsResult{Runtime: l.runtime}
}

func TestSessionReloadKeepsTheOldRuntimeActiveOnFailure(t *testing.T) {
	for _, tc := range []struct {
		name      string
		fresh     bool
		err       error
		wantStale bool
	}{
		{name: "kept"},
		{name: "replaced", fresh: true, wantStale: true},
		{name: "failed", err: errors.New("resource reload failed")},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("PIG_HOME", t.TempDir())
			services, err := CreateAgentSessionServices(CreateAgentSessionServicesOptions{CWD: t.TempDir()})
			if err != nil {
				t.Fatal(err)
			}
			old := extension.CreateExtensionRuntime()
			loader := &runtimeReloadLoader{ResourceLoader: NoResources, runtime: old}
			session, err := NewSession(services, SessionOptions{ResourceLoader: loader})
			if err != nil {
				t.Fatal(err)
			}
			defer func() { _ = session.Close() }()
			if _, err := session.ReloadExtensions(nil, old); err != nil {
				t.Fatal(err)
			}
			ai.RegisterAPIProvider(ai.APIProvider{API: "reload-serving-api"}, "reload-serving-source")
			t.Cleanup(ai.ResetAPIProviders)
			servingAPI := ai.GetAPIProvider("reload-serving-api")
			loader.err = tc.err
			if tc.fresh {
				loader.runtime = extension.CreateExtensionRuntime()
			}
			if err := session.Reload(t.Context()); !errors.Is(err, tc.err) {
				t.Fatalf("Reload error = %v, want %v", err, tc.err)
			}
			if stale := old.AssertActive() != nil; stale != tc.wantStale {
				t.Fatalf("old runtime stale = %v, want %v", stale, tc.wantStale)
			}
			if tc.err != nil && ai.GetAPIProvider("reload-serving-api") != servingAPI {
				t.Fatal("failed reload replaced the serving API registration")
			}
			if session.ExtensionRunner().Runtime().AssertActive() != nil || session.ExtensionRunner().IsStale() {
				t.Fatal("the serving runner or its runtime is stale")
			}
		})
	}
}
