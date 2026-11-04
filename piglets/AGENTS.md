# AGENTS.md — piglets

Authored Piglets, separate from Stock PiG. Bare `pig` never imports or activates them.

| Path | What |
|---|---|
| `standard/` | PiG Standard: branded presentation, defaults, selected extensions, onboarding, games (own module) |
| `porter/` | Pig Porter: parity workbench; procedure `porter/skills/pig-porter/SKILL.md`; run `make porter`, `make porter-task` |
| `aot/` | experimental aot TypeScript runtime Piglet |

- Product behavior goes here, not into `internal/`, `cmd/` or `coding/`. Bootstrap mechanisms (parser, resolver, verifier, SDK bridge, isolation, build) stay in Stock.
- Standard changes: `make standard-check`. Porter: `make porter-check`, `make porter-smoke`.
- Porter campaigns (`make porter-campaign`) are read-only `inventory`/`verify` jobs on disposable snapshots; their reports are proposals tied to one commit, upstream version and family. Production/evidence edits stay serialized in the canonical workspace.
