package codingagent

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/MichaelKinsy/PiG/coding/packagecontent"
	"github.com/MichaelKinsy/PiG/internal/resolvepath"
)

type configuredSkillRoot struct {
	dir       string
	baseDir   string
	scope     string
	overrides []string
	agents    bool
}

func resolveConfiguredSkillRoots(scan *AgentsSkillScan, cwd, agentDir string, sm *SettingsManager, project, user bool) ([]configuredSkillRoot, bool) {
	global := sm.GetGlobalSettings()
	settings, baseDir, scope := global, agentDir, "user"
	configured := global.skillRootsSet
	if project {
		if candidate := sm.GetProjectSettings(); candidate.skillRootsSet {
			settings, baseDir, scope, configured = candidate, ProjectConfigDir(cwd), "project", true
		}
	}
	if !configured {
		return nil, false
	}
	roots := make([]configuredSkillRoot, 0, len(settings.SkillRoots))
	seen := make([]string, 0, len(settings.SkillRoots))
	add := func(dir, base, rootScope string, overrides []string, agents bool) {
		enabled := rootScope == "project" && project || rootScope == "user" && user
		if !enabled {
			return
		}
		resolved, err := ResolvePath(dir, base)
		if err != nil {
			return
		}
		key := resolved
		if real, err := filepath.EvalSymlinks(resolved); err == nil {
			key = real
		}
		for _, prior := range seen {
			if resolvepath.PathsEqual(prior, key) {
				return
			}
		}
		seen = append(seen, key)
		roots = append(roots, configuredSkillRoot{dir: resolved, baseDir: base, scope: rootScope, overrides: overrides, agents: agents})
	}
	for _, raw := range settings.SkillRoots {
		switch raw {
		case "@project":
			add(filepath.Join(ProjectConfigDir(cwd), "skills"), ProjectConfigDir(cwd), "project", sm.GetProjectSettings().Skills, false)
		case "@project-agents":
			for _, dir := range scan.Dirs() {
				add(dir.Dir, filepath.Dir(dir.Dir), "project", sm.GetProjectSettings().Skills, true)
			}
		case "@agent":
			add(filepath.Join(agentDir, "skills"), agentDir, "user", global.Skills, false)
		case "@user-agents":
			if home, err := PackageManagerHomeDir(); err == nil {
				dir := filepath.Join(home, ".agents", "skills")
				add(dir, filepath.Dir(dir), "user", global.Skills, true)
			}
		default:
			add(raw, baseDir, scope, settings.Skills, false)
		}
	}
	return roots, true
}

func (sm *SettingsManager) noteInvalidSkillRoots(scope, path, baseDir string, loaded Settings) {
	if !loaded.skillRootsSet {
		return
	}
	invalid := make([]string, 0)
	for _, raw := range loaded.SkillRoots {
		if strings.HasPrefix(raw, "@") {
			continue
		}
		resolved, err := ResolvePath(raw, baseDir)
		if err != nil {
			invalid = append(invalid, fmt.Sprintf("%q (%v)", raw, err)) //perf:ok one diagnostic per invalid configured root at settings load
			continue
		}
		if info, err := os.Stat(resolved); err == nil && !info.IsDir() { //perf:ok one stat per configured root at settings load
			invalid = append(invalid, fmt.Sprintf("%q (not a directory)", raw)) //perf:ok one diagnostic per invalid configured root at settings load
		}
	}
	if len(invalid) > 0 {
		sm.errors = append(sm.errors, SettingsError{Scope: scope, Path: path, Error: fmt.Errorf("invalid skills.roots: %s", strings.Join(invalid, ", "))})
	}
}

func discoverConfiguredSkillRoot(root configuredSkillRoot) []string {
	if root.agents {
		return packagecontent.DiscoverAgentSkillDirs(root.dir)
	}
	return packagecontent.DiscoverSkillDirs(root.dir)
}
