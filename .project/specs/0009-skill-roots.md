---
document:
  type: "specification"
  id: "SPEC-0009"
  title: "Skill roots setting: skills.roots[]"
  description: "Configurable list of directories searched for skills, replacing the fixed auto-discovery roots; tokens keep the default roots addressable"
  status: "draft"
  revision: 1
  created_at: "2026-11-04T14:58:00+01:00"
  updated_at: "2026-11-04T14:58:00+01:00"
  created_by: "orchestrator"
  updated_by: "orchestrator"
  changelog:
    - revision: 1
      reason: "Owner asked for a skills.roots[] setting (U-24)"
      changes: ["document creation"]
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
