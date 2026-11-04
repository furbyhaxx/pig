# AGENTS.md — piglets

Authored Piglets, separate from Stock PiG. Bare `pig` never imports or activates them.

| Path | What |
|---|---|
| `standard/` | PiG Standard: branded presentation, defaults, selected extensions, onboarding, games (own module) |
| `porter/` | Pig Porter: parity workbench; procedure `porter/skills/pig-porter/SKILL.md`; run `make porter`, `make porter-task` |
| `aot/` | experimental aot TypeScript runtime Piglet |

- Opinionated composition, presentation and workflow go here; generic mechanisms (parser, resolver, verifier, SDK bridge, isolation, build) stay in core. Core defaults may be decided in a fork spec.
- Standard changes: `make standard-check`. Porter: `make porter-check`, `make porter-smoke`.
- Porter campaigns (`make porter-campaign`) are read-only `inventory`/`verify` jobs on disposable snapshots; their reports are proposals tied to one commit, upstream version and family. Production/evidence edits stay serialized in the canonical workspace.
