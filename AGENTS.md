# AGENTS.md — pig (furbyhaxx fork)

Personal coding-agent base: a Go CLI, TUI, SDK and providers, derived from PiG (`MichaelKinsy/PiG`, a Go port of Pi `earendil-works/pi`). Goal: fast, performant, well-engineered, built to the owner's needs. Go module path stays `github.com/MichaelKinsy/PiG`.

## Goal and compatibility

Decisions: `.project/specs/0001-personal-distribution.md` (D7–D10).

- Pi and upstream PiG compatibility (behavior, extension API, wire, config, sessions, CLI) is kept while it costs nothing. It is not a requirement.
- Don't keep a defect, a weak design or a slow path because Pi or upstream PiG does it. Don't justify a change with "faithful to Pi"; justify it with correctness, performance, design quality or the owner's requirement.
- Breaking compatibility (config shape, session format, extension API, wire protocol) → decided in a `.project/specs/` spec first, with a migration note in user docs.
- Upstream remains a source of fixes and features: read upstream code and tests to learn intent, then choose the best design for this fork.
- Child `AGENTS.md` files under `docs/parity/`, `test/parity/` and other ported areas describe upstream's parity machinery. Their parity rules apply only when maintaining that machinery, never to fork behavior decisions. This file wins on conflict.

## Fork

- User guide: `docs/personal-fork.md`.
- Hourly sync merges `MichaelKinsy/PiG:main` (`automation/release/personal-fork-sync.sh`). Each new source commit on `main` releases `x.y.z-c.N` (`automation/release/personal-fork.sh`).
- Merge cost is weighed, not a veto. Same design quality → new files and a narrow seam in the upstream file. Upstream structure blocks a good design → change it.
- All `AGENTS.md` outside `testdata/` are fork-owned; merges keep ours (`.gitattributes` `merge=ours`). The sync lists upstream changes to them in the merge commit → port the rules that serve the fork; read upstream's with `git show upstream/main:<path>`. Local merges: `git config merge.ours.driver true` once.
- Upstream tests or gates that read `AGENTS.md` content break after syncs (ours is kept, their data moves on). Decouple them from `AGENTS.md`; don't mirror upstream text into ours.
- Don't commit `-c.N` versions or run `make set-version` on `main`; the release script sets the version in its build tree. Upstream owns `PigVersion` bumps.
- Don't edit inherited upstream workflows; GitHub runs only `personal-fork.yml`.
- Upstream-owned artifact dirs: `docs/plan/`, `docs/specs/`, `docs/findings/`, `plans/`, `research/`, `handoffs/`. Read; never add fork artifacts there.
- `.project/` — development artifacts; never in `docs/`; committed. Load skill `dot-project` before writing one. Every Markdown artifact needs `document:` frontmatter (`.project/AGENTS.md`). Open work: `BACKLOG.md` pointers + `.project/backlog/`; live claims: project-scope `tasks`. Never put drift/TODOs in AGENTS.md.

## Behavior changes

- Every deliberate change of behavior inherited from upstream: decision in a `.project/specs/` spec, user docs updated in the same change, `changelog.d/` fragment.
- Don't record fork changes as `docs/parity/divergences/D<N>.md`, `// pig divergence (D<N>)` markers or `docs/additive-features.md` entries; those are upstream's ledgers. Cite the spec in the commit body.
- Go tests that pin the replaced behavior or a Pi golden → rewrite them to the new behavior or delete them, in the same change. Parity scenarios under `test/parity/scenarios/` are not maintained by the fork (SPEC-0001 D10); don't edit them for fork changes. Never weaken, skip or normalize a test to keep unrelated behavior green.
- Existing `D<N>` records and `docs/additive-features.md` stay as upstream maintains them; read them to learn why upstream PiG differs from Pi.

## Sources of truth

No prose claim overrides these. Never hand-edit generated files. After changing exported Go API, CLI flags, settings or docs mirrors: `make generate`, commit the result.

| What | Where |
|---|---|
| Fork decisions and designs | `.project/specs/`; work breakdown `.project/plans/` |
| Upstream Pi source (read-only, generated) | `.upstream/current/`; missing → `make upstream-mirror` |
| Upstream pin | `.upstream/pin`; move only with `make pin-bump PI_UPSTREAM_REF=<tag\|commit>` (`docs/parity/upstream-pin.md`); generates `internal/coding/pigversion/pigversion.go` |
| Upstream Pi file → Go file | `docs/parity/PORT_MAP.md` |
| Pi-comparison scenarios (upstream's parity suite) | `test/parity/scenarios/<family>/*.toml`; report `test/parity/coverage.md` (`make coverage`) |
| Upstream PiG's differences from Pi | `docs/parity/DIVERGENCES.md`, `docs/additive-features.md` |
| Glossary (Session, Provider, Model Runtime, Main Screen, Piglet, ...) | `docs/project/CONTEXT.md`; use its terms, never redefine |
| TS → Go semantics of ported code | `docs/typescript-to-go-porting.md` |
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
- Extension API change (host call, event, field, wire message) → every SDK (`extensions/sdk` Go, `sdk-rs`, `sdk-py`, `sdk-ts` declarations), `test/extension-conformance/`, `docs/site/docs/extensions.md` + `internal/pigdocs/content/extension-api.md` in the same change. An SDK left behind is drift: backlog item, never silent.
- Setting added or changed → settings type in `internal/codingagent/settings.go`, `make generate-schemas`, `docs/site/docs/settings.md` + `internal/pigdocs/content/settings.md`.

## Mechanism vs policy

Generic mechanisms (discover, validate, load, isolate, run, render, compose) live in the core packages. Opinionated composition, presentation and personal workflow live in extensions, Piglets (`piglets/`) or config, and select those mechanisms.

- A core mechanism has one owner, a typed contract and documented defaults. Defaults may be opinionated when the owner decided them in a spec.
- A Piglet cannot supply its own parser, resolver, verifier, SDK bridge, isolation boundary or build machinery; those stay in core.
- Caller-free or unneeded code → delete.

## Engineering

- Read the code and tests you change and their callers first. Upstream source and tests explain intent; they don't bind the design.
- Fix at the lowest shared layer. No UI or provider band-aids for a core defect.
- Never hard-code a provider or model in shared code. Use shared catalog data (model catalog, `defaultModelPerProvider`, auth metadata) and one shared function per concept.
- Test shared paths with several provider shapes, never only Copilot or `test-faux`. Fix a special case's siblings together.
- Changing ported code: preserve the JS/Node semantics callers rely on (omitted/undefined/null, insertion order, UTF-16/code-point/terminal-cell units, Node path/URL/regex/JSON behavior) unless the spec changes them. Reuse `internal/js*` and `internal/node*` (`internal/AGENTS.md`).
- Go structure: struct for data; small consumer-owned interface only for a real behavioral contract. Constructors return concrete types.
- Don't change a public extension or TUI contract for presumed speed; profile first, prove allocations and byte-faithful output before and after.

Shared provider/model/auth/tool regression matrix (explain an inapplicable row from the contract, not from available credentials; `--list-models` is not stream/payload proof):

| Axis | Minimum cases |
|---|---|
| provider shape | OAuth, API key, OpenAI-compatible with custom base URL, no default model — as the path allows |
| API path | each API kind reached; both `openai-completions` (`ai/openai.go`) and `openai-responses` (`ai/openai_responses.go`) when both apply; `agent/transform.go` before conversion |
| mode | pure function/converter unit test + nearest caller (print, interactive or RPC) |
| history | fresh + poisoned persisted history when relevant |
| prior assistant turn | error, aborted, empty text, tool-call/orphaned tool-call when relevant |
| auth/network | hermetic success + hermetic failure when touching auth refresh/error display |

### Concurrency and lifetime

| Situation | Go |
|---|---|
| blocking operation | blocking call returning value/error; no `go func()` without a reason |
| parallel work | owned, joined concurrency (`errgroup`, `wg.Go`) with deterministic result/error order |
| background work | owned task with cancellation, error reporting, shutdown drain; never a naked goroutine |
| failure | error on the same operation; never swallowed |
| request/session lifetime | `context.Context`; cancellation propagates to every callee, remote listeners included |
| slow work triggered by TUI/events | off the UI loop; UI mutation via the loop's run-on-main path |
| waits on external processes or sockets | bounded by a deadline or cancellation; a stuck peer is reported, never waited on forever |

- Evidence must prove completion order, cancellation and error propagation; happy-path output is not async proof.
- Security-relevant hooks (redaction, policy, blocking) fail closed: a handler error or panic withholds the action.

## Verification

- Bug fix → regression guard in the same change that fails on the bug, or name the existing failing test. Reproduce first; a test written after the fix is "not red-proven" — say so.
- Core surfaces (agent loop, message normalization, provider payload conversion, session persistence, auth, tool execution, interactive rendering): red → green → refactor.
- Cross-layer bug → unit test at the pure layer + regression at the nearest caller (e.g. `NormalizeMessages` + OpenAI converter; session loop + print/interactive error display).
- Persisted-session bugs → poisoned history: errored/aborted assistant turns, empty text blocks, nil content, orphaned tool calls, model/provider switches.
- No accepted flake. Don't hide failures with retries, sleeps, longer budgets, env overrides, skips or normalization. Timeouts/retries only for faults outside our control, and surface them. Unavailable external dependency = blocker, not pass.
- A new production API used to fix behavior needs a production call site or a test that drives the production path.
- Exact-count assertions derive the expectation from an independent denominator or a reviewed generated snapshot plus its regeneration command.
- Mutation evidence must compile, unless compilation is the contract.
- Critical paths (extension, TUI, model/provider, Session/RPC, Package/Resource, startup/reload): representative benchmark/profile + resource-lifetime check. Test empty, ordinary, boundary and stress inputs; verify ordering, cancellation, replacement, cleanup, retained state.
- Performance claims need a representative benchmark with allocations and a CPU/allocation profile. Interface dispatch is not a defect without evidence.
- Known bug found in scope → fix with red/mutation evidence, or report it and file a backlog item.
- Green is necessary, not sufficient. Stop and report when everything passed first try or a check passed before and after the supposed fix.
- Fix report: failing symptom, regression test names, spec cited, commands run.

## Test isolation

- `PIG_CODING_AGENT_DIR` and `PI_CODING_AGENT_DIR` beat `PIG_HOME`, and every child process inherits them, including tmux sessions started from your shell. Any manual probe or test that runs pig sets or clears all of `PIG_CODING_AGENT_DIR`, `PI_CODING_AGENT_DIR`, `PIG_HOME`, `PI_HOME`, `HOME`. Otherwise `/logout` and login rewrite the real credentials.
- Plain `go test` from an agent shell: point those four `*_DIR`/`*_HOME` variables at scratch dirs first, unless the package's `TestMain` calls `testenv.ScopeTempDir` or `testenv.RunScoped`. `make test` fails a group that changes them (`automation/ci/agent-dir-guard.sh`).
- Skills under `~/.agents/skills` follow `HOME`, not `PIG_HOME`: set `HOME` too. A settings `skills.roots` without `@user-agents` keeps them out; use both for a hermetic probe.
- tmux: unique session names; only `tmux kill-session -t <name>`; never `tmux kill-server`.

## Code and lint

- Go toolchain: `.go-version` is the only pin; scripts read it. Moving it updates its recorded copies (workflows, images, devcontainer, `toolchain` lines in `go.mod`/`go.work`) in one change; `make compliance` and `TestGoVersionPolicy` enforce it. Module `go` directives stay at the 1.26 floor unless a language need is approved.
- Use modern stdlib: `slices`, `maps`, `cmp`, `strings.Cut`/`CutPrefix`/`SplitSeq`, `for i := range n`, `wg.Go`, `min`/`max`. `go fix -diff ./...` stays empty (`make go-fix-clean`).
- Production-grade Go: concrete validated data, minimal consumer-owned interfaces, deterministic order, bounded resources, explicit lifetime and error ownership, one mechanism per concept.
- Comments state current behavior and non-obvious protocol, concurrency, security or compatibility invariants. The source-hygiene gate rejects roadmap phases, workstream labels, future promises, placeholder plans and pointers to PORT_MAP/status ledgers. Nothing is grandfathered: revalidate old comments, placeholders and compatibility paths when touched. Don't leave "Ports packages/..." or "upstream:" claims on code the fork redesigned.
- Observable behavior change → update the doc comment in the same commit.
- Lint: fix findings at the source. Narrowest suppression with a reason (`nolintlint` rejects stale or unexplained pragmas); security suppressions name the CLI threat-model reason. Don't suppress `iface`; delete the unneeded interface. Don't enable or apply `ireturn`/`interfacebloat` mechanically.
- `make divergence-guard` (in `make fork-check`) flags engineering smells in `agent`, `ai`, `coding`, `internal/codingagent`, `tui`, `cmd/pig`: `magic-literal`, `cancel-drop-send`, `default-drop-send`, `stream-parse-continue`, `stop-reason-success`, `hand-rolled-sse` (SSE only in `ai/sse.go`), `discarded-io-error`, `hook-fail-open`, `error-status-substring`, `agent-run-outside-session` (agent `Send`/`Continue` only in `coding/session*.go`); `bare-recover` everywhere. Fix the smell. The guard accepts `// fork (SPEC-NNNN): <reason>` when `.project/specs/NNNN-*.md` exists, plus inherited upstream markers (`// pig divergence (D<N>): ...`, `// upstream: <file>:<symbol>`); don't invent `D<N>` numbers for fork code. `automation/ci/divguard/baseline.toml` is shrink-only.
- `make port-lint` (`automation/ci/portlint/README.md`): fix at the source or `//portlint:allow <check> <reason>` for a reviewed false positive; never add baseline entries for new code.
- Pig-owned formats (source, manifest, record, graph/plan, builder, Marketplace, publication, subprocess): one strict current unversioned shape. No format discriminator, negotiation, compatibility reader, migration, redirect or dual write unless a spec decides it.
- Self-update: one native ownership tier per operation (standalone download, package manager, immutable Piglet Binary, OCI/Image, read-only); never substitute one for another after work starts. Downloads need full checksum and overflow-safe size checks. Startup checks stay bounded and best-effort; explicit update failures stay actionable.
- Release version: only via `make set-version VERSION=x.y.z` (fork: see Fork).

## Commands

- Iterate: `go test ./<pkg> -run '<Tests>'`. `-race` on edited packages only. No `-count=1` unless the test reads inputs through a child process. Bound fuzzing with `-parallel` and `-fuzztime`. Shared host: build with `-trimpath`.
- Before a commit and before reporting ready: `make fork-check` (`automation/make/fork.mk`: build, `lint-changed`, Windows vet of changed packages, `test-changed`, `go-fix-clean`, `divergence-guard`, `port-lint`, `perf-lint`, `any-trap`, `source-hygiene`, `docs-drift`, `changelog-check`). SDK source change → also that SDK's suite (`make test-sdk-go|test-sdk-py|test-sdk-rs|test-sdk-ts`) and `make sdk-bundles-check`.
- Don't run upstream's Pi-compatibility, porting-ledger or publication gates unless the owner asks: `make check`, `verify`, `ci-drift`, `ci-contracts`, `parity*`, `journeys*`, `portmap-check`, `port-map-drift`, `coverage-drift`, `divergence-consistency|quality`, `test-provenance`, `test-expectation-lint`, `sdk-surface-drift`, `sdk-apidiff`, `bug-ledger-check`, `compliance`, `module-publication`, `module-size`, `release-check`, `porter-*` (SPEC-0001 D10). Their failures are not fork defects.
- `make test` requires Go, Node, Python, Rust; fixture build failures are fatal. Prereqs: `make doctor`; setup: skill `setup-pig`.
- Lint config change: `go tool golangci-lint config verify`.
- All targets: `make help`, `make help-parity`; fork gates in `automation/make/fork.mk`. Don't `go install ./cmd/pig`; use `make pig` (→ `bin/pig`) or `make install`.

## Public writing

Applies to `README.md`, `docs/`, `changelog.d/`, user-visible Go strings; `automation/ci/check-public-claims.py` enforces part of it.
- Describe this repo as a personal fork of PiG; link PiG and Pi as its origins. Credit Michael Kinsy as PiG's creator and state PiG was originally developed at Hewlett Packard Enterprise, separate from legal notices.
- Don't call it an official Pi or PiG project or imply endorsement. Don't claim HPE sponsorship, OSPO approval, community governance, certification or trademark approval. No HPE logos/wordmarks/styling.
- Don't claim Pi parity for behavior the fork changed. Docs describe the fork's behavior at the point of use and name the spec when it differs from Pi or upstream PiG.
- Direct technical prose; verifiable facts. No slogans, superlatives, competitive framing or future promises.
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
- Set `PIG_HOME` and `PIG_CODING_AGENT_DIR` to a path containing a space. Most Windows bugs are path bugs (separators, drive letters, 8.3 names).
