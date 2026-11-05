# AGENTS.md — test

Cross-package suites. Package-local tests stay next to their package.

| Path | What | Run |
|---|---|---|
| `parity/` ● | Pi vs PiG scenario parity, ledgers, coverage | `make parity*`, `make coverage` |
| `extension-conformance/` ● | cross-SDK, cross-transport extension behavior | `make test-sdk-conformance` |
| `integration/` | real pig + faux provider, tmux, no network | `make test-integration` |
| `journeys/` | PTY journeys vs Pi: requests, sessions, frames | `make journeys`, `make journeys-changed` |
| `pipeline-diff/` | wire-request differential fuzz vs Pi | `make test-pipeline-fuzz` |
| `perfgate/` | whole-process count gates | `go test ./test/perfgate` |
| `docs-drift/` | docs bundle vs code, public claims | `make docs-drift` |
| `upstream-parity/` | declaration-driven API parity | `go test ./test/upstream-parity` |
| `gomodule/` | public Go module contract | `go test ./test/gomodule` |
| `ci-images/` | CI image pins, installer tests (`install_sh*_test.go`) | `go test ./test/ci-images` |
| `evals/`, `longsession/`, `mcp-conformance/`, `wiring/`, `generate/` | evals, long-session fixtures, MCP suite, SDK wiring, generator help | see each README/doc.go |

- A Pi golden names its regeneration script in `testdata-provenance.toml`. Fork-owned goldens need no Pi provenance; prefer an invariant assertion over a value snapshot (SPEC-0001 D10).
- `docs/site/public/install.sh` change → `test/ci-images/install_sh_test.go`.
