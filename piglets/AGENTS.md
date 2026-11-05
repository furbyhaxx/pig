# AGENTS.md — piglets

Authored Piglets, separate from Stock PiG. Bare `pig` never imports or activates them.

| Path | What |
|---|---|
| `standard/` | PiG Standard: branded presentation, defaults, selected extensions, onboarding, games (own module) |
| `porter/` | Pig Porter: upstream's parity workbench; not used by the fork (SPEC-0001 D10); don't run `porter*` targets or fix their failures |
| `aot/` | experimental aot TypeScript runtime Piglet |

- Opinionated composition, presentation and workflow go here; generic mechanisms (parser, resolver, verifier, SDK bridge, isolation, build) stay in core. Core defaults may be decided in a fork spec.
- Standard changes: `make standard-check`.
