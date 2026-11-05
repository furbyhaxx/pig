---
document:
  type: "specification"
  id: "SPEC-0009"
  title: "Skill roots setting: skills.roots[]"
  description: "Configurable list of directories searched for skills, replacing the fixed auto-discovery roots; tokens keep the default roots addressable"
  status: "draft"
  revision: 2
  created_at: "2026-11-04T14:58:00+01:00"
  updated_at: "2026-11-05T15:20:00+01:00"
  created_by: "orchestrator"
  updated_by: "coding-agent"
  changelog:
    - revision: 1
      reason: "Owner asked for a skills.roots[] setting (U-24)"
      changes: ["document creation"]
    - revision: 2
      reason: "Implemented and accepted through the built CLI"
      changes: ["record the canonical object form and malformed-entry handling", "record configured-root precedence and the shared discovery mechanism", "record implementation evidence and test coverage"]
---

# SPEC-0009 Skill roots

Resolves U-24 of `../analyses/2026-11-04-host-gap-verification/README.md`.

## Context

- Skill auto-discovery searches fixed roots: project `.pig/skills`, ancestor `.agents/skills` up to the git root (trusted projects only), `<agentDir>/skills`, `~/.agents/skills` (`internal/codingagent/resource_paths.go`). `~/.agents/skills` follows `HOME`, so `PIG_HOME` does not isolate it.
- Setting `skills` is a list of extra skill paths and `+`/`-`/`!` override patterns (`internal/codingagent/settings.go:357`).
- An object-valued `skills` is treated as Pi's deprecated `{enableSkillCommands, customDirectories}` form and migrated: an object without `customDirectories` is deleted (`settings.go:1194-1205`, `3575-3588`). A new object form must change that migration.

## D1 Shape

```json
{
  "skills": {
    "roots": ["@project", "@project-agents", "~/work/skills", "@agent"],
    "paths": ["~/one-off/SKILL.md", "-noisy-skill"]
  }
}
```

- `skills` accepts an object `{roots?, paths?}` or, as today, a list (shorthand for `{paths: [...]}`).
- `paths` keeps the current meaning of the list form.
- Legacy migration maps `customDirectories` → `paths` and `enableSkillCommands` → top-level `enableSkillCommands`; it never deletes `roots` or `paths`.

## D2 Roots

- `roots` replaces the auto-discovery root list. Absent → default list.
- Entries: absolute path, `~/…`, path relative to the settings file's directory, or a token:

| Token | Directory |
|---|---|
| `@project` | `<cwd>/.pig/skills` (project config dir) |
| `@project-agents` | every `.agents/skills` from cwd up to the git root |
| `@agent` | `<agentDir>/skills` |
| `@user-agents` | `~/.agents/skills` |

- Default list: `["@project", "@project-agents", "@agent", "@user-agents"]` (current order and behavior).
- Order is precedence: on a skill name collision the earlier root wins and the existing collision diagnostic names both files.
- Duplicate roots (same real path) are searched once.
- A missing root is skipped silently; a root that is not a directory is a settings diagnostic.

## D3 Scopes and trust

- Project `skills.roots` overrides global `skills.roots` (normal settings precedence); tokens re-include defaults.
- Project-scope roots and the `@project*` tokens apply only in trusted projects, as today.
- CLI `--skill` adds paths; `--no-skills` disables root discovery and keeps explicit paths, as today.

## D4 Isolation

`skills.roots` set without `@user-agents` stops `~/.agents/skills` from loading, so test setups isolate with settings instead of `HOME`. Root `AGENTS.md` test-isolation line is updated when this lands.

## Docs

- `skills.md` (both trees): roots, tokens, precedence, trust.
- `settings.md` (both trees) and `settings-examples.md`: `skills` object form; `make generate-schemas` (`internal/codingagent/configschema/settings.go` gains the object form).
- Root `AGENTS.md` test-isolation bullet.

## Implementation (05.11.2026)

- The object form is canonical when written back; the list form stays `["…"]`, so the settings-file byte oracle is unchanged for existing files. `customDirectories` migrates to `paths` only when `paths` is absent.
- A malformed `roots` or `paths` entry is dropped individually with a diagnostic, both when decoding and when migrating the stored object. One bad entry cannot erase the other member on the next setter write.
- Duplicate roots are compared by real path, so a symlink to the same directory is searched once. Roots configured in settings sort before other auto-discovered resources while keeping their list order; settings `paths` entries keep their existing higher precedence.
- `AmbientSkillResourcesIn` is the one discovery mechanism. `coding/cli` config and status inventory consume it rather than performing a second fixed-root scan.
- Evidence: `TestSkillRootsSettingsSurviveMigrationAndRoundTrip`, `TestSkillRootsReplaceDefaultDiscoveryInOrder`, `TestProjectSkillRootsOverrideGlobalAndRequireTrust`, `TestSkillRootsPersistAcrossPathSettersAndLegacyFileMigration`, `TestSkillRootTokensRelativePathsDuplicatesAndDiagnostics`, `TestSkillRootsWrongTypedEntryDoesNotEraseTheObject`, `TestSkillRootsDeduplicateASymlinkedRoot`, `TestStatusHonorsConfiguredSkillRoots`. The built binary listed only the configured project skill, then the home skill too after `@user-agents` was added.
