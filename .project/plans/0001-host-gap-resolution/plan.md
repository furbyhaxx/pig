---
document:
  type: "plan"
  id: "PLAN-0001"
  title: "Resolve the Workbench host gaps"
  description: "Slices that implement SPEC-0002..SPEC-0010 after re-steering the fork (SPEC-0001 D7-D10), with tests, wire/SDK parity and docs per slice"
  status: "draft"
  revision: 1
  created_at: "2026-11-04T15:05:00+01:00"
  updated_at: "2026-11-04T15:05:00+01:00"
  created_by: "orchestrator"
  updated_by: "orchestrator"
  changelog:
    - revision: 1
      reason: "Owner asked for a plan that resolves the verified host gaps, linked to the specs, including docs"
      changes: ["document creation"]
---

# PLAN-0001 Resolve the Workbench host gaps

Input: `../../analyses/2026-11-04-host-gap-verification/README.md`. Policy: `../../specs/0001-personal-distribution.md` D7–D10.

## Progress

| Slice | Goal | Specs | Status |
|---|---|---|---|
| S0 | Steering, gates, public fork statement | 0001 | now |
| S1 | Reliability fixes and skill roots | 0004 D1–D3, 0007 D1–D3, 0009 | next |
| S2 | Model request pipeline | 0003 | later |
| S3 | System prompt composition | 0002 | later |
| S4 | Headless output and dialogs | 0005 | later |
| S5 | Session listing, tools, credentials | 0007 D4–D6, 0008 | later |
| S6 | UI surfaces and renderers | 0006, 0004 D4 | later |
| S7 | Packed build inputs | 0010 | later |
| S8 | Retained processes | 0004 D5 | later |

No work: U-09 and U-26 (wrong findings), U-20 (watch; regression test when it reproduces).

## Rules for every slice

- Per task: red test first for bugs; implementation; wire + all SDKs + `test/extension-conformance/` for extension API changes; docs in both trees (`docs/site/docs/`, `internal/pigdocs/content/`); `make generate` when settings, flags or exported API change; `changelog.d/` fragment for user-facing changes.
- Upstream tests and parity scenarios that pin replaced behavior are rewritten or deleted in the same change (SPEC-0001 D10); the commit body cites the spec.
- Gate per task: `make test-changed lint-changed`, `go vet` on touched packages, `-race` on edited packages. Before closing a slice: `make ci-drift ci-contracts`, `make check`.
- A spec decision that turns out wrong is revised in the spec in the same change.

## S0 Steering and gates (now)

| # | Task | Ref | Status |
|---|---|---|---|
| 0.1 | Rewrite root `AGENTS.md` for the fork goal | 0001 D7–D10 | done |
| 0.2 | SPEC-0001 revision 2 (D7–D10) | 0001 | done |
| 0.3 | Child `AGENTS.md` files: scope parity-mandating rules to the parity machinery or replace them with fork rules: `docs/AGENTS.md` (D<N> rule), `docs/parity/AGENTS.md`, `test/parity/AGENTS.md`, `coding/extension/AGENTS.md`, `ai/AGENTS.md`, `agent/AGENTS.md`, `extensions/AGENTS.md`, `internal/AGENTS.md`, `tui/AGENTS.md`, `.github/AGENTS.md`, `automation/AGENTS.md`, `coding/AGENTS.md`, `piglets/AGENTS.md`, `mcp/AGENTS.md`; skill `agents-md` | 0001 D7, D10 | next |
| 0.4 | `divergence-guard`: accept `// fork (SPEC-NNNN): <reason>` when `.project/specs/NNNN-*.md` exists; guard test; then update root `AGENTS.md` Code and lint and SPEC-0001 D10 | 0001 D10 | next |
| 0.5 | `automation/ci/check-public-claims.py`: make sure it accepts the fork statement of 0.6; adjust rules that require Pi-faithful wording | 0001 D7 | next |
| 0.6 | Docs: `docs/personal-fork.md` replaces "Runtime behavior … remain upstream-compatible" with the D7 policy; `README.md` replaces "If PiG behaves differently from Pi, that is either a bug or a documented divergence" with the fork statement and a link to `docs/personal-fork.md`; changelog fragment | 0001 D7 | next |
| 0.7 | Outside this repo: pi-accordion resolves a Node executable instead of `process.execPath` (stray billed sessions) | analysis §3 | owner |
| 0.8 | `TestFamilyCoverageReportMatchesCurrentStaticDashboard` (`test/parity/closure/report_projection_test.go:744`) reads the family dashboard from root `AGENTS.md`, which `f0db6902b` removed; it fails on `main`. Read the dashboard from the generated coverage report instead | 0001 D6 | next |

Done when: `make check` passes, no child `AGENTS.md` demands Pi faithfulness for fork behavior, public docs state the fork policy.

## S1 Reliability fixes and skill roots (next)

| # | Task | Ref | Tests | Docs |
|---|---|---|---|---|
| 1.1 | Packed Go reload hang: looping red test reproducing P29; confirm cause; generation-unique sockets | 0004 D1 | `TestPackedGoReloadLoopNeverHangs` (fails before fix) | `docs/runtime-cell-reload.md`, bundle `runtime-cells.md` |
| 1.2 | Bounded reload handshake, atomic reload, error to `ctx.Reload()`, `/reload`, RPC; setting `extensions.reloadHandshakeMs` | 0004 D2 | stuck-peer test (process never registers) keeps old generation serving | `settings.md`, `extensions.md` reload section |
| 1.3 | Event-bus cancellation (`events.cancel`) host + all SDKs | 0004 D3 | conformance test: caller abort ends listener context in Go/Rust/Python/Node | `extensions.md` event bus |
| 1.4 | `nextTurn` bound to leaf; discard + `queued_message_discarded` on navigation/new/fork/resume/switch | 0007 D1 | red test from probe P20; poisoned history case | `extensions.md`, `events.md` |
| 1.5 | Strict `deliverAs`; SDK constants | 0007 D3 | wire test with unknown value → error | `extensions.md`, `rpc.md` |
| 1.6 | `session_start` `resume` for continued sessions; `initial` flag | 0007 D2 | `--continue`, `--resume`, `--session`, fresh start | `events.md`, changelog `### Upgrading` |
| 1.7 | `skills.roots[]`: settings object form, legacy migration fix, tokens, precedence, trust; `make generate-schemas` | 0009 | migration keeps `roots`/`paths`; collision order; untrusted project ignores project roots; `HOME`-independent isolation | `skills.md`, `settings.md`, `settings-examples.md`, root `AGENTS.md` test-isolation bullet |

Done when: all S1 tests red-proven where they are bug fixes, `make check` passes, docs and changelog fragments merged.

## S2 Model request pipeline (later)

Goal: SPEC-0003 D1–D6. Key steps: move hooks into `ModelRuntime` (`coding/services.go`), `purpose`/`origin`, re-entrancy guard, fail-closed table, `block` result, phases, `extensions.hookTimeoutMs`, wire + SDKs, benchmark of the no-handler path. Docs: `extensions.md`, `extension-api.md`, `events.md`, `security.md`, `compaction.md`, `settings.md`.

## S3 System prompt composition (later)

Goal: SPEC-0002 D1–D10. Key steps: gonja dependency + `nix/vendor-hash.sh update`; `internal/codingagent/prompts/compose` package; `builtin:default` byte-identical (existing prompt tests unchanged); catalog with owners; section/layout API in all SDKs; layout selection and trust; forced prompt recording; transcript-truth invariant test; `pig prompt render`, `/prompt`; benchmarks. Docs: new `system-prompt.md` (both trees, `docs.json`, bundle index), `settings.md`, `extensions.md`, `extension-api.md`, `trust.md`, `prompt-templates.md`, `cli.md`, `slash-commands.md`.

## S4 Headless output and dialogs (later)

Goal: SPEC-0005 D1–D4. Unblocks Workbench `--list-tools`. Docs: `extensions.md`, `extension-api.md`, `rpc-extension-ui.md`, `json.md`, `settings.md`; changelog `### Upgrading` for Go SDK dialog signatures.

## S5 Session listing, tools, credentials (later)

Goal: SPEC-0007 D4–D6, SPEC-0008 D1–D2 (D2 starts with the hermetic repro). Docs: `sessions.md`, `extensions.md`, `rpc-commands.md`, `providers.md`, `custom-provider.md`, `cli.md` (`--tools` hides MCP and codemode tools; `-ne` disables built-in extensions).

## S6 UI surfaces and renderers (later)

Goal: SPEC-0006 D1–D2, SPEC-0004 D4. Docs: `tui.md`, `extensions.md`, `codemode.md`.

## S7 Packed build inputs (later)

Goal: SPEC-0010 D1–D2 with before/after build benchmarks. Docs: `runtime-cells.md`, `docs/runtime-cell-build-cache.md`, `docs/pig-package-spec.md`, `packages.md`, `piglets.md`.

## S8 Retained processes (later)

Goal: SPEC-0004 D5. Docs: `extensions.md`, `extension-api.md`, `settings.md`.
