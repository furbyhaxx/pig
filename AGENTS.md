# AGENTS.md — pig (furbyhaxx fork)

Go port of Pi (`earendil-works/pi`): coding-agent CLI, TUI, SDK and providers. This repo is a personal distribution forked from `MichaelKinsy/PiG`. Target: Pi's observable behavior, except numbered divergences in `docs/parity/DIVERGENCES.md`. Go module path stays `github.com/MichaelKinsy/PiG`.

## Fork

- Decisions: `.project/specs/0001-personal-distribution.md`. User guide: `docs/personal-fork.md`.
- Hourly sync merges `MichaelKinsy/PiG:main` (`automation/release/personal-fork-sync.sh`). Each new source commit on `main` releases `x.y.z-c.N` (`automation/release/personal-fork.sh`).
- Every fork edit to an upstream file is a future merge conflict. Prefer new files. Keep fork logic in `automation/release/personal-fork*`, `.github/workflows/personal-fork.yml`, `docs/personal-fork.md`, `.project/`.
- All `AGENTS.md` outside `testdata/` are fork-owned; merges keep ours (`.gitattributes` `merge=ours`). The sync lists upstream changes to them in the merge commit → port applicable rules here; read upstream's with `git show upstream/main:<path>`. Local merges: `git config merge.ours.driver true` once.
- Upstream tests or gates that read `AGENTS.md` content break after syncs (ours is kept, their data moves on). Decouple them from `AGENTS.md`; don't mirror upstream text into ours.
- Don't commit `-c.N` versions or run `make set-version` on `main`; the release script sets the version in its build tree. Upstream owns `PigVersion` bumps.
- Don't edit inherited upstream workflows; GitHub runs only `personal-fork.yml`.
- Upstream-owned artifact dirs: `docs/plan/`, `docs/specs/`, `docs/findings/`, `plans/`, `research/`, `handoffs/`. Read; never add fork artifacts there.
- `.project/` — development artifacts; never in `docs/`; committed. Load skill `dot-project` before writing one. Every Markdown artifact needs `document:` frontmatter (`.project/AGENTS.md`). Open work: `BACKLOG.md` pointers + `.project/backlog/`; live claims: project-scope `tasks`. Never put drift/TODOs in AGENTS.md.

## Sources of truth

No prose claim overrides these. Never hand-edit generated files. After changing exported Go API, CLI flags, settings, scenarios or docs mirrors: `make generate`, commit the result.

| What | Where |
|---|---|
| Upstream source (read-only, generated) | `.upstream/current/`; missing → `make upstream-mirror` |
| Upstream pin | `.upstream/pin`; move only with `make pin-bump PI_UPSTREAM_REF=<tag\|commit>` (`docs/parity/upstream-pin.md`); generates `internal/coding/pigversion/pigversion.go` |
| Upstream file → Go file | `docs/parity/PORT_MAP.md` |
| Observed behavior | `test/parity/scenarios/<family>/*.toml` |
| Verification report (generated) | `test/parity/coverage.md` (`make coverage`) |
| Intentional differences | `docs/parity/DIVERGENCES.md`, generated from `docs/parity/divergences/D<N>.md` |
| Additive (non-Pi) capabilities | `docs/additive-features.md` |
| Public interface denominator, mapping | `test/parity/interfaces/upstream-v<ver>.json`, `mapping-v<ver>.json` |
| Version leap accounting | `test/parity/upstream-sync/v<ver>.toml` |
| Promise/async contracts | `test/parity/async-contracts.toml` |
| Glossary (Session, Provider, Model Runtime, Main Screen, Piglet, ...) | `docs/project/CONTEXT.md`; use its terms, never redefine |
| TS → Go semantics | `docs/typescript-to-go-porting.md` |
| Which doc to read for a task | `docs/README.md` |

## Layout

Own `AGENTS.md` in each dir marked ●. Read it before working there.

| Path | What |
|---|---|
| `agent/` ● | agent loop, message types (Pi `packages/agent`) |
| `ai/` ● | providers, model catalogs, Event Stream (Pi `packages/ai`) |
| `coding/` ● | coding-agent SDK: Session, Model Runtime, modes, CLI (`coding/cli/`) |
| `coding/extension/` ● | extension API, host, subprocess wire, runtime cells |
| `tui/` ● | terminal rendering, editor, components (Pi `packages/tui`) |
| `internal/` ● | shared helpers incl. JS/Node semantics; interactive app in `internal/codingagent/` |
| `chord/` `durable/` `env/` `telemetry/` `codemode/` `mcp/` ● | ports of the other Pi monorepo packages |
| `extensions/` ● | extension SDKs: Go, Rust, Python, TypeScript declarations |
| `piglets/` ● | PiG Standard, Porter, aot Piglets |
| `examples/` ● | example extensions and SDK programs (product surface) |
| `test/` ● | cross-package suites; `test/parity/` ●, `test/extension-conformance/` ● |
| `docs/` ● | maintainer refs + user docs; `docs/parity/` ● |
| `automation/` ● | every non-Go script; `Makefile` is the interface |
| `changelog.d/` | changelog fragments; every file there must be a fragment (no AGENTS.md); rules under Git |
| `.github/` ● | fork release workflow, inherited upstream CI |
| `cmd/pig/` | binary entry; calls `coding/cli.Main`. Other `cmd/*`: generators |
| `schemas/` | generated JSON Schemas (`make generate-schemas`) |
| `third_party/` | vendored forks of external modules; each has its own notes (e.g. `FORK.md`) |
| `bench/` | benchmarks and eval harnesses |
| `nix/` | flake package; `vendorHash` in `nix/package.nix` |

Cross-updates:
- `go.mod`, `go.sum`, `go.work*` or a new module import in `cmd/pig` → `nix/vendor-hash.sh update` (or set the hash it prints).
- New package/module/area → its `AGENTS.md` in the same change (skill `agents-md`).

## Stock PiG vs product

Stock PiG = upstream-parity behavior + inert, product-neutral mechanisms to discover, validate, select, build, verify, load, isolate and run Packages, Resources, extensions, Piglets, Piglet Binaries. Classify every additive capability before implementing:

1. Required substrate → Stock.
2. Inert capability that an extension or Piglet must select → Stock.
3. Composition, presentation, policy, branding → PiG Standard (`piglets/standard/`) or another Piglet/Package.
4. Caller-free or unneeded → delete.
5. Overlaps Pi but differs observably → fix, or numbered divergence after explicit user approval.

- Don't activate Standard behavior in Stock. Don't put product UI, Piglet-specific recommendations, workflow logic or linked extension products under `internal/` or `cmd/` beyond the minimum generic host capability.
- A Piglet cannot supply its own parser, resolver, verifier, SDK bridge, isolation boundary or build machinery; keep those bootstrap mechanisms in Stock.
- Additive records: `docs/additive-features.md`. Observable differences: `docs/parity/DIVERGENCES.md`. One short typed source marker at each production decision point; tests don't count as markers.
- `explicit unsafe opt-in`: off by default, visible at configuration boundaries, never selected by PiG Standard. `product-neutral diagnostics`: generic reports about Stock mechanisms.

## Porting rules

- Port semantics, not syntax: omitted/undefined/null, closed unions, insertion order, UTF-16/code-point/terminal-cell units, numeric conversion, exceptions, callback ownership, Node path/URL/regex/JSON behavior, every state or wire effect callers observe. Reuse `internal/js*` and `internal/node*` (`internal/AGENTS.md`).
- Mirror upstream names and structure 1:1 (camelCase → exported PascalCase, same literal union values, same field names). Rename only when Go forces it (no `Symbol.for`, no field-bearing interfaces, discriminated union → sealed interface); record each case. A needless rename inflates the next sync.
- TS `interface` ≠ Go interface: struct for data; small consumer-owned interface only for a real behavioral contract. Constructors return concrete types.
- New port: production comment `// Ports packages/<pkg>/src/<file>.ts` + PORT_MAP row in the same change. Status `partial` until behavioral closure is reviewed; the drift gate rejects port claims left not-started, deferred or n/a.
- A production file is `ported` only when every upstream test covering it is ported with the same inputs and expectations, or designed out with a numbered reason (`test/parity/AGENTS.md`).
- Choose the option that makes the next upstream sync smaller. No pig-specific behavior in mirrored code unless it is a numbered divergence.
- Don't change a public extension or TUI contract for presumed speed; profile first, prove allocations and byte-faithful output before and after.
- Unsure → read upstream source and tests, probe upstream, match it. Still unsure → ask the user.

### Promise/async

| Upstream | Go |
|---|---|
| awaited Promise | blocking call returning value/error; no `go func()` just because TS says `async` |
| `Promise.all`/`allSettled` | owned, joined concurrency with deterministic result/error order |
| intentionally unawaited Promise | owned background task with cancellation, error reporting, shutdown drain; never a naked goroutine |
| rejection | error on the same operation unless upstream handles it |
| `AbortSignal`, request/session lifetime | `context.Context` or equivalent cancellation |
| async TUI/event handler | slow work off the UI loop; UI mutation via the loop's run-on-main path |

- Probe call sites, not declared types: an `async` function without `await` still adds a microtask boundary; a Promise-returning callback may be awaited by its host.
- Evidence must prove completion order, cancellation and error propagation; happy-path output is not async proof.

## Faithful, general implementations

- Port Pi's design, not one case's output. Use Pi's shared data (model catalog, `defaultModelPerProvider`, auth metadata) and one shared function where Pi has one.
- Never hard-code a provider or model in shared code; when Pi itself special-cases one, cite the Pi line.
- Test shared paths with several provider shapes, never only Copilot or `test-faux`. Fix a special case's siblings together.

Shared provider/model/auth/tool regression matrix (explain inapplicable rows from the upstream contract, not from available credentials; `--list-models` is not stream/payload proof):

| Axis | Minimum cases |
|---|---|
| provider shape | OAuth, API key, OpenAI-compatible with custom base URL, no default model — as the path allows |
| API path | each API kind reached; both `openai-completions` (`ai/openai.go`) and `openai-responses` (`ai/openai_responses.go`) when both apply; `agent/transform.go` before conversion |
| mode | pure function/converter unit test + nearest caller (print, interactive or RPC) |
| history | fresh + poisoned persisted history when relevant |
| prior assistant turn | error, aborted, empty text, tool-call/orphaned tool-call when relevant |
| auth/network | hermetic success + hermetic failure when touching auth refresh/error display |

## Divergences

- Allowed only when user-visible or interop-relevant and it cannot or should not be faithful now. Required: record `docs/parity/divergences/D<N>.md` (`docs/parity/AGENTS.md`), `SCRUTINIZED:approved`, remove-when condition, `// pig divergence (D<N>): ...` at every call site, parity coverage or explicit allowance.
- Not divergences: TS → Go mechanics, harness differences, per-SDK implementation details. An SDK missing a capability or answering differently is drift; record and close it.
- Where a user can observe a difference, the docs state it at the point of use and cite `D<N>`.
- A prior divergence is evidence, not permanent approval; re-probe at every upstream leap.
- `make divergence-guard` (in `make check`) flags in `agent`, `ai`, `coding`, `internal/codingagent`, `tui`, `cmd/pig`: `magic-literal`, `cancel-drop-send`, `default-drop-send`, `stream-parse-continue`, `stop-reason-success`, `hand-rolled-sse` (SSE only in `ai/sse.go`), `discarded-io-error`, `hook-fail-open`, `error-status-substring`, `agent-run-outside-session` (agent `Send`/`Continue` only in `coding/session*.go`); `bare-recover` everywhere. A hit passes only with `// pig divergence (D<N>): ...` or `// upstream: <file>:<symbol>` on its line or the line above. The cited file must exist in `.upstream/current` or a verified dependency snapshot (`test/parity/dependency-sources.json`, `docs/parity/dependency-source-proof.md`); a literal's value must appear there.
- `automation/ci/divguard/baseline.toml` is shrink-only: a fix deletes its entries in the same change; never add an entry for new code.

## Verification

- Bug fix → regression guard in the same change that fails on the bug, or name the existing failing test. Reproduce first; a test written after the fix is "not red-proven" — say so.
- Core parity surfaces (agent loop, message normalization, provider payload conversion, session persistence, auth, tool execution, interactive rendering): red → green → refactor. Encode the upstream rule in the test name/comment. Fix at the lowest shared layer that matches upstream; no UI/provider band-aids.
- Cross-layer bug → unit test at the pure layer + regression at the nearest caller (e.g. `NormalizeMessages` + OpenAI converter; session loop + print/interactive error display).
- Persisted-session bugs → poisoned history: errored/aborted assistant turns, empty text blocks, nil content, orphaned tool calls, model/provider switches.
- No accepted flake. Don't hide failures with retries, sleeps, longer budgets, env overrides, skips or normalization. Timeouts/retries only for faults outside our control, and surface them. Unavailable external dependency = blocker, not pass.
- A new production API used to fix behavior needs a production call site or a test that drives the production path.
- Exact-count assertions derive the expectation from an independent denominator or a reviewed generated snapshot plus its regeneration command.
- Mutation evidence must compile, unless compilation is the contract.
- Critical paths (extension, TUI, model/provider, Session/RPC, Package/Resource, startup/reload): representative benchmark/profile + resource-lifetime check. Test empty, ordinary, boundary and stress inputs; verify ordering, cancellation, replacement, cleanup, retained state.
- Performance claims need a representative benchmark with allocations and a CPU/allocation profile. Interface dispatch is not a defect without evidence.
- Known bug found in scope → fix with red/mutation evidence, or report it as a user-owned blocker and file a backlog item.
- Green is necessary, not sufficient. Stop and report when everything passed first try or a check passed before and after the supposed fix.
- Core parity fix report: failing symptom, regression test names, mirrored upstream rule, commands run.

## Test isolation

- `PIG_CODING_AGENT_DIR` and `PI_CODING_AGENT_DIR` beat `PIG_HOME`, and every child process inherits them, including tmux sessions started from your shell. Any manual probe or test that runs pig sets or clears all of `PIG_CODING_AGENT_DIR`, `PI_CODING_AGENT_DIR`, `PIG_HOME`, `PI_HOME`, `HOME`. Otherwise `/logout` and login rewrite the real credentials.
- Plain `go test` from an agent shell: point those four `*_DIR`/`*_HOME` variables at scratch dirs first, unless the package's `TestMain` calls `testenv.ScopeTempDir` or `testenv.RunScoped`. `make test` fails a group that changes them (`automation/ci/agent-dir-guard.sh`).
- tmux: unique session names; only `tmux kill-session -t <name>`; never `tmux kill-server`.

## Code and lint

- Go toolchain: `.go-version` is the only pin; scripts read it. Moving it updates its recorded copies (workflows, images, devcontainer, `toolchain` lines in `go.mod`/`go.work`) in one change; `make compliance` and `TestGoVersionPolicy` enforce it. Module `go` directives stay at the 1.26 floor unless a language need is approved.
- Use modern stdlib: `slices`, `maps`, `cmp`, `strings.Cut`/`CutPrefix`/`SplitSeq`, `for i := range n`, `wg.Go`, `min`/`max`. `go fix -diff ./...` stays empty (`make go-fix-clean`).
- Production-grade Go: concrete validated data, minimal consumer-owned interfaces, deterministic order, bounded resources, explicit lifetime and error ownership, one mechanism per concept.
- Comments state current behavior and non-obvious protocol, concurrency, security or compatibility invariants. The source-hygiene gate rejects roadmap phases, workstream labels, future promises, placeholder plans and pointers to PORT_MAP/status ledgers. Nothing is grandfathered: revalidate old comments, placeholders and compatibility paths when touched.
- Observable behavior change → update the doc comment in the same commit.
- Lint: fix findings at the source; upstream behavior beats style. Narrowest suppression with a reason (`nolintlint` rejects stale or unexplained pragmas); security suppressions name the CLI threat-model reason. Don't suppress `iface`; delete the unneeded interface. Don't enable or apply `ireturn`/`interfacebloat` mechanically.
- `make port-lint` (`automation/ci/portlint/README.md`): fix at the source or `//portlint:allow <check> <reason>` for a reviewed false positive; never add baseline entries for new code.
- Pig-owned formats (source, manifest, record, graph/plan, builder, Marketplace, publication, subprocess): one strict current unversioned shape. No format discriminator, negotiation, compatibility reader, migration, redirect or dual write unless upstream versions the equivalent schema or the user approves.
- Self-update: one native ownership tier per operation (standalone download, package manager, immutable Piglet Binary, OCI/Image, read-only); never substitute one for another after work starts. Downloads need full checksum and overflow-safe size checks. Startup checks stay bounded and best-effort; explicit update failures stay actionable.
- Release version: only via `make set-version VERSION=x.y.z` (fork: see Fork).

## Commands

- Iterate: `go test ./<pkg> -run '<Tests>'`. `-race` on edited packages only. No `-count=1` unless the test reads inputs through a child process. Bound fuzzing with `-parallel` and `-fuzztime`. Shared host: build with `-trimpath`.
- Before a commit: `make test-changed lint-changed`, `go vet` on touched packages.
- Before reporting ready: `make ci-drift ci-contracts` + affected `make parity-family FAMILY=<f>`.
- Full: `make check` (deterministic gates + `parity-fast`), `make verify` (end-of-loop, regenerates coverage), `make test`, `make parity`, `make lint`. Tiers: `docs/project/test-economy.md`.
- `make test` requires Go, Node, Python, Rust; fixture build failures are fatal. Prereqs: `make doctor`; setup: skill `setup-pig`.
- Lint config change: `go tool golangci-lint config verify`.
- All targets: `make help`, `make help-parity`. Don't `go install ./cmd/pig`; use `make pig` (→ `bin/pig`) or `make install`.

## Public writing

Applies to `README.md`, `docs/`, `changelog.d/`, user-visible Go strings; `automation/ci/check-public-claims.py` enforces part of it.
- Describe PiG as a faithful Go implementation of Pi; name and link Pi as the reference. Credit Michael Kinsy as creator and state PiG was originally developed at Hewlett Packard Enterprise, separate from legal notices.
- Don't call PiG an official Pi project or imply endorsement. Don't claim HPE sponsorship, OSPO approval, community governance, certification or trademark approval. No HPE logos/wordmarks/styling. Stock PiG branding stays product-neutral.
- Direct technical prose; verifiable facts (pinned Pi version, observed parity, named divergences, published evidence). No slogans, superlatives, competitive framing or future promises.
- Coverage figures in prose must match `test/parity/coverage.md`.
- Don't hard-wrap Markdown prose or comments.

## Git

- User-facing change → fragment `changelog.d/<lane>-<sha>.md`, `<sha>` = a commit of the range reachable from `HEAD`. Never edit `CHANGELOG.md`; `make changelog-fold VERSION=x.y.z` folds fragments at release.
  - Sections in `CHANGELOG.md` style, bullets only: `### Added`, `### Changed`, `### Fixed`, `### Upgrading` (what a previous-release user must change), `### Deprecated`, `### Removed`, `### Security`, `### Known issues`.
  - Every `### Fixed` bullet ends with `<!-- fix: <sha>[,<sha>] -->` naming fixing commits reachable from `HEAD` (or a squash commit's `Fix-Commit:` trailer); the bug ledger reads it.
  - Rebase/squash/amend replaced the SHA → rename the fragment and update markers. Check: `make changelog-check` (full history); preview: `make changelog-unreleased`. `automation/ci/changelogfrag/legacy.txt` only shrinks.
- Conventional Commits. Upstream's DCO sign-off and signed-commit rules govern upstream contributions only.
- Testdata dirs accumulate temp files; inspect before staging.
- Forbidden: `git stash`, `git commit --no-verify`, `tmux kill-server`.
- Read local files yourself; don't dispatch subagents to read them.

## Windows

First-class target. Fix Windows-only failures on a Windows host, not by CI round trips.
- Toolchain: `.go-version`, `.node-version`, Python 3.12, Rust per `RUST_VERSION` in `.github/workflows/ci.yml`, Git for Windows (run `make`/`automation/` from Git Bash), `fd`/`rg` via `automation/ci/install-windows-tools.ps1`. Mirror: `MSYS=winsymlinks:nativestrict ./automation/gen/mirror-upstream.sh`.
- Reproduce with the `windows` job's command from `.github/workflows/ci.yml`, `go test -v` and a `-timeout` below the job limit.
- Set `PIG_HOME` and `PIG_CODING_AGENT_DIR` to a path containing a space. Most Windows bugs are path bugs (separators, drive letters, 8.3 names); compare with Pi on Windows before changing a path rule.
