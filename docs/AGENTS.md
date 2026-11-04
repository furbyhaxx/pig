# AGENTS.md — docs

Maintainer references and the public user documentation. Product surface: keep current with every behavior change. Task router: `README.md` — read the one row you need, not everything.

| Path | What |
|---|---|
| `site/docs/` | public user docs (rendered by the separate site repo); `site/public/install.sh` is the fork installer |
| `parity/` ● | PORT_MAP, divergences, parity ledgers |
| `project/` | CONTEXT glossary, governance, RELEASING, PROVENANCE, quality, test economy |
| `*.md` (root) | architecture and extension/Piglet/Package specs |
| `personal-fork.md` | fork install, update, sync and release guide |
| `plan/`, `specs/`, `findings/` | upstream-owned artifacts; read only, never add (fork artifacts → `.project/`) |
| `bug-ledger.json`, `knowledge-graph/`, `performance/scoreboard.md` | generated (`make generate`) |

- The agent/offline docs bundle is `internal/pigdocs/content/`, separate from `site/docs/`; a user-visible change usually updates both.
- Language: simple direct technical English, one instruction per sentence, present tense, active voice, one term per thing. Don't apply this mechanically to essays, interviews or release narratives.
- Observable PiG ≠ Pi difference → state it at the point of use and cite `D<N>`.
- Pin-dependent statements must name the pinned Pi version; earlier versions only in files listed in `automation/ci/version-records.toml`.
