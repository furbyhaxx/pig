package codingagent

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/MichaelKinsy/PiG/internal/testenv"
)

func TestSkillRootsSettingsSurviveMigrationAndRoundTrip(t *testing.T) {
	var settings Settings
	if err := json.Unmarshal([]byte(`{"skills":{"roots":["@project","./shared"],"paths":["./one-off/SKILL.md"],"customDirectories":["./legacy"],"enableSkillCommands":false}}`), &settings); err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(settings.SkillRoots, []string{"@project", "./shared"}) {
		t.Fatalf("SkillRoots = %v", settings.SkillRoots)
	}
	if !slices.Equal(settings.Skills, []string{"./one-off/SKILL.md"}) {
		t.Fatalf("Skills = %v; paths must win over legacy customDirectories", settings.Skills)
	}
	if settings.EnableSkillCommands == nil || *settings.EnableSkillCommands {
		t.Fatal("enableSkillCommands was not migrated")
	}
	data, err := json.Marshal(settings)
	if err != nil {
		t.Fatal(err)
	}
	var wire struct {
		Skills struct {
			Roots               []string `json:"roots"`
			Paths               []string `json:"paths"`
			CustomDirectories   []string `json:"customDirectories"`
			EnableSkillCommands *bool    `json:"enableSkillCommands"`
		} `json:"skills"`
	}
	if err := json.Unmarshal(data, &wire); err != nil {
		t.Fatalf("round trip produced %s: %v", data, err)
	}
	if !slices.Equal(wire.Skills.Roots, settings.SkillRoots) || !slices.Equal(wire.Skills.Paths, settings.Skills) {
		t.Fatalf("round trip = %s", data)
	}
	if wire.Skills.CustomDirectories != nil || wire.Skills.EnableSkillCommands != nil {
		t.Fatalf("legacy fields leaked into the object form: %s", data)
	}
}

func TestSkillRootsReplaceDefaultDiscoveryInOrder(t *testing.T) {
	cwd, agentDir, home := t.TempDir(), t.TempDir(), t.TempDir()
	t.Setenv("HOME", home)
	makeRootSkill := func(root, name string) string {
		t.Helper()
		path := filepath.Join(root, name, "SKILL.md")
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		body := "---\nname: shared\ndescription: " + name + "\n---\n" + name
		if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
			t.Fatal(err)
		}
		return root
	}
	first := makeRootSkill(filepath.Join(t.TempDir(), "unused"), "first")
	second := makeRootSkill(filepath.Join(t.TempDir(), "unused"), "second")
	makeRootSkill(filepath.Join(home, ".agents", "skills"), "home-must-not-load")
	if err := os.MkdirAll(filepath.Join(cwd, CONFIG_DIR_NAME), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(agentDir, 0o755); err != nil {
		t.Fatal(err)
	}
	data := []byte(`{"skills":{"roots":[` + strconvQuote(first) + `,` + strconvQuote(second) + `],"paths":[]}}`)
	if err := os.WriteFile(filepath.Join(agentDir, "settings.json"), data, 0o600); err != nil {
		t.Fatal(err)
	}
	sm := NewSettingsManager(cwd, agentDir)
	resources := AmbientSkillResources(cwd, agentDir, sm, true, true)
	if len(resources) != 2 {
		t.Fatalf("resources = %#v; default and home roots must be replaced", resources)
	}
	want := []string{filepath.Join(first, "first"), filepath.Join(second, "second")}
	if resources[0].Path != want[0] || resources[1].Path != want[1] {
		t.Fatalf("root precedence = [%s %s], want %v", resources[0].Path, resources[1].Path, want)
	}
	if diagnostic := resources[1].Metadata.Scope; diagnostic != "user" {
		t.Fatalf("root scope = %q", diagnostic)
	}
}

func TestProjectSkillRootsOverrideGlobalAndRequireTrust(t *testing.T) {
	for _, trusted := range []bool{true, false} {
		t.Run(map[bool]string{true: "trusted", false: "untrusted"}[trusted], func(t *testing.T) {
			base := t.TempDir()
			cwd, agentDir := filepath.Join(base, "project"), filepath.Join(base, "agent")
			for _, dir := range []string{filepath.Join(cwd, CONFIG_DIR_NAME), agentDir} {
				if err := os.MkdirAll(dir, 0o755); err != nil {
					t.Fatal(err)
				}
			}
			globalRoot := filepath.Join(base, "global-skills")
			projectRoot := filepath.Join(base, "project-skills")
			for _, root := range []string{globalRoot, projectRoot} {
				if err := os.MkdirAll(root, 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(root, "SKILL.md"), []byte("---\nname: root\ndescription: root\n---\nroot"), 0o600); err != nil {
					t.Fatal(err)
				}
			}
			if err := os.WriteFile(filepath.Join(agentDir, "settings.json"), []byte(`{"skills":{"roots":[`+strconvQuote(globalRoot)+`]}}`), 0o600); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(cwd, CONFIG_DIR_NAME, "settings.json"), []byte(`{"skills":{"roots":[`+strconvQuote(projectRoot)+`]}}`), 0o600); err != nil {
				t.Fatal(err)
			}
			sm := NewSettingsManagerWithProjectTrust(cwd, agentDir, trusted)
			resources := AmbientSkillResources(cwd, agentDir, sm, trusted, true)
			want := globalRoot
			if trusted {
				want = projectRoot
			}
			if len(resources) != 1 || resources[0].Path != want {
				t.Fatalf("trusted=%v resources=%#v, want only %s", trusted, resources, want)
			}
		})
	}
}

func TestSkillRootsPersistAcrossPathSettersAndLegacyFileMigration(t *testing.T) {
	agentDir := t.TempDir()
	path := filepath.Join(agentDir, "settings.json")
	if err := os.WriteFile(path, []byte(`{"skills":{"roots":["@agent"],"customDirectories":["./legacy"],"enableSkillCommands":false}}`), 0o600); err != nil {
		t.Fatal(err)
	}
	sm := NewSettingsManager(t.TempDir(), agentDir)
	if err := sm.SetSkillPaths([]string{"./explicit"}); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var file struct {
		Skills struct {
			Roots             []string `json:"roots"`
			Paths             []string `json:"paths"`
			CustomDirectories []string `json:"customDirectories"`
		} `json:"skills"`
		EnableSkillCommands *bool `json:"enableSkillCommands"`
	}
	if err := json.Unmarshal(data, &file); err != nil {
		t.Fatalf("settings file %s: %v", data, err)
	}
	if !slices.Equal(file.Skills.Roots, []string{"@agent"}) || !slices.Equal(file.Skills.Paths, []string{"./explicit"}) || file.Skills.CustomDirectories != nil {
		t.Fatalf("settings file lost the skills object: %s", data)
	}
	if file.EnableSkillCommands == nil || *file.EnableSkillCommands {
		t.Fatalf("legacy enableSkillCommands was not moved: %s", data)
	}
}

func TestSkillRootTokensRelativePathsDuplicatesAndDiagnostics(t *testing.T) {
	base := t.TempDir()
	repo := filepath.Join(base, "repo")
	cwd := filepath.Join(repo, "nested")
	agentDir := filepath.Join(base, "agent")
	home := filepath.Join(base, "home")
	t.Setenv("HOME", home)
	for _, dir := range []string{filepath.Join(repo, ".git"), filepath.Join(cwd, CONFIG_DIR_NAME), agentDir} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	skill := func(root, name string) string {
		t.Helper()
		path := filepath.Join(root, name, "SKILL.md")
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte("---\nname: "+name+"\ndescription: "+name+"\n---\n"), 0o600); err != nil {
			t.Fatal(err)
		}
		return filepath.Dir(path)
	}
	wantProject := skill(filepath.Join(cwd, CONFIG_DIR_NAME, "skills"), "project")
	wantNested := skill(filepath.Join(cwd, ".agents", "skills"), "nested-agents")
	wantRepo := skill(filepath.Join(repo, ".agents", "skills"), "repo-agents")
	wantRelative := skill(filepath.Join(cwd, CONFIG_DIR_NAME, "relative"), "relative")
	wantAgent := skill(filepath.Join(agentDir, "skills"), "agent")
	wantUser := skill(filepath.Join(home, ".agents", "skills"), "user-agents")
	if err := os.WriteFile(filepath.Join(cwd, CONFIG_DIR_NAME, "not-a-dir"), []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	projectSettings := `{"skills":{"roots":["@agent","@project","@project-agents","./relative","./relative/../relative","./missing","./not-a-dir","@user-agents"]}}`
	if err := os.WriteFile(filepath.Join(cwd, CONFIG_DIR_NAME, "settings.json"), []byte(projectSettings), 0o600); err != nil {
		t.Fatal(err)
	}
	sm := NewSettingsManagerWithProjectTrust(cwd, agentDir, true)
	errs := sm.DrainErrors()
	if len(errs) != 1 || errs[0].Scope != "project" || !strings.Contains(errs[0].Error.Error(), "not a directory") || !strings.Contains(errs[0].Error.Error(), `"./not-a-dir"`) {
		t.Fatalf("settings diagnostics = %v", errs)
	}
	resources := AmbientSkillResources(cwd, agentDir, sm, true, true)
	got := make([]string, len(resources))
	for i, resource := range resources {
		got[i] = resource.Path
	}
	want := []string{wantAgent, wantProject, wantNested, wantRepo, wantRelative, wantUser}
	if !slices.Equal(got, want) {
		t.Fatalf("resources = %v, want %v", got, want)
	}
	ordered := OrderResolvedResources(nil, resources)
	got = got[:0]
	for _, resource := range ordered {
		got = append(got, resource.Path)
	}
	if !slices.Equal(got, want) {
		t.Fatalf("resolver order = %v, want root order %v", got, want)
	}
}

func TestSkillRootsWrongTypedEntryDoesNotEraseTheObject(t *testing.T) {
	agentDir := t.TempDir()
	path := filepath.Join(agentDir, "settings.json")
	if err := os.WriteFile(path, []byte(`{"skills":{"roots":["@agent",true],"paths":["./keep",7]}}`), 0o600); err != nil {
		t.Fatal(err)
	}
	sm := NewSettingsManager(t.TempDir(), agentDir)
	errs := sm.DrainErrors()
	if len(errs) != 1 || !strings.Contains(errs[0].Error.Error(), "skills.roots") || !strings.Contains(errs[0].Error.Error(), "skills.paths") {
		t.Fatalf("diagnostics = %v", errs)
	}
	if settings := sm.GetGlobalSettings(); !settings.skillRootsSet || !slices.Equal(settings.SkillRoots, []string{"@agent"}) || !slices.Equal(settings.Skills, []string{"./keep"}) {
		t.Fatalf("settings after wrong-typed entries = %+v", settings)
	}
	if err := sm.SetSkillPaths([]string{"./new"}); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var file struct {
		Skills struct {
			Roots []string `json:"roots"`
			Paths []string `json:"paths"`
		} `json:"skills"`
	}
	if err := json.Unmarshal(data, &file); err != nil {
		t.Fatalf("settings file %s: %v", data, err)
	}
	if !slices.Equal(file.Skills.Roots, []string{"@agent"}) || !slices.Equal(file.Skills.Paths, []string{"./new"}) {
		t.Fatalf("settings write lost the roots: %s", data)
	}
}

func TestSkillRootsDeduplicateASymlinkedRoot(t *testing.T) {
	base := t.TempDir()
	agentDir := filepath.Join(base, "agent")
	target := filepath.Join(base, "skills")
	alias := filepath.Join(base, "skills-alias")
	for _, name := range []string{"one", "two"} {
		path := filepath.Join(target, name, "SKILL.md")
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte("---\nname: "+name+"\ndescription: "+name+"\n---\n"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.MkdirAll(agentDir, 0o755); err != nil {
		t.Fatal(err)
	}
	testenv.Symlink(t, target, alias)
	data := []byte(`{"skills":{"roots":[` + strconvQuote(target) + `,` + strconvQuote(alias) + `]}}`)
	if err := os.WriteFile(filepath.Join(agentDir, "settings.json"), data, 0o600); err != nil {
		t.Fatal(err)
	}
	resources := AmbientSkillResources(t.TempDir(), agentDir, NewSettingsManager(t.TempDir(), agentDir), true, true)
	if len(resources) != 2 {
		t.Fatalf("resources = %#v; a symlinked root must be searched once", resources)
	}
}

func strconvQuote(value string) string {
	encoded, _ := json.Marshal(value)
	return string(encoded)
}
