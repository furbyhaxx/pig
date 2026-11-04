package configschema_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/MichaelKinsy/PiG/internal/coding/pigversion"
	"github.com/MichaelKinsy/PiG/internal/codingagent"
	"github.com/MichaelKinsy/PiG/internal/codingagent/configschema"
)

func assertSchemaContract(t *testing.T, path, rendered string, origin []byte) {
	t.Helper()
	if path != "schemas/settings.schema.json" {
		if rendered != string(origin) {
			t.Errorf("%s differs from its origin", path)
		}
		return
	}
	var got, want map[string]any
	if err := json.Unmarshal([]byte(rendered), &got); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(origin, &want); err != nil {
		t.Fatal(err)
	}
	want["properties"].(map[string]any)["extensionsReloadHandshakeMs"] = map[string]any{
		"type": "integer", "minimum": float64(1), "maximum": float64(9223372036854),
		"default": float64(15000), "description": "Global reload-only connect/register budget in milliseconds; build time is excluded.",
	}
	if !reflect.DeepEqual(got, want) {
		t.Error("settings schema differs from its origin plus SPEC-0004 D2")
	}
}

func TestRenderedSchemasMatchPinnedOriginAndForkContract(t *testing.T) {
	pi := filepath.Join(repoRoot(t), "extensions", "sdk-ts", "node_modules", "@earendil-works", "pi-coding-agent")
	var manifest struct{ Version string }
	data, err := os.ReadFile(filepath.Join(pi, "package.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(data, &manifest); err != nil {
		t.Fatal(err)
	}
	if manifest.Version != pigversion.UpstreamPackageVersion {
		t.Fatalf("Pi oracle version %s, want %s", manifest.Version, pigversion.UpstreamPackageVersion)
	}
	rendered, err := configschema.Render(codingagent.KeybindingsSchemaArtifact())
	if err != nil {
		t.Fatal(err)
	}
	for path, got := range rendered {
		want, err := os.ReadFile(filepath.Join(pi, path))
		if err != nil {
			t.Fatal(err)
		}
		assertSchemaContract(t, path, got, want)
	}
}

func TestRenderedSchemasMatchUpstreamMirrorAndForkContract(t *testing.T) {
	pi := filepath.Join(repoRoot(t), ".upstream", "current", "packages", "coding-agent")
	rendered, err := configschema.Render(codingagent.KeybindingsSchemaArtifact())
	if err != nil {
		t.Fatal(err)
	}
	for path, got := range rendered {
		want, err := os.ReadFile(filepath.Join(pi, path))
		if err != nil {
			t.Fatal(err)
		}
		assertSchemaContract(t, path, got, want)
	}
}
