---
document:
  type: "plan"
  id: "PLAN-0001"
  title: "Resolve the Workbench host gaps"
  description: "Slices that implement SPEC-0002..SPEC-0010 after re-steering the fork (SPEC-0001 D7-D10), with tests, wire/SDK parity and docs per slice"
  status: "draft"
  revision: 13
  created_at: "2026-11-04T15:05:00+01:00"
  updated_at: "2026-11-05T15:40:00+01:00"
  created_by: "orchestrator"
  updated_by: "coding-agent"
  changelog:
    - revision: 1
      reason: "Owner asked for a plan that resolves the verified host gaps, linked to the specs, including docs"
      changes: ["document creation"]
    - revision: 2
      reason: "S0 implementation and targeted verification; full gates blocked by pinned npm build"
      changes: ["record S0 tasks and red-green evidence", "record B-0003 prerequisite blocker", "leave S1-S8 pending"]
    - revision: 3
      reason: "Owner requested fixing the integrity blocker"
      changes: ["close B-0003 with pinned runtime and preflight", "fix nested-worktree format scans", "record passing drift/contracts and remaining B-0004 failures"]
    - revision: 4
      reason: "Owner requested continuing the plan; carry known S0 failures while starting S1"
      changes: ["record additional S0 red-green repairs", "implement S1.1 socket ownership and stress coverage", "record P29 and S1.2 setting-shape questions"]
    - revision: 5
      reason: "Implement the approved flat handshake setting and Host/CLI transaction path without claiming the SDK loader or full slice complete"
      changes: ["record S1.2 timeout and rollback evidence", "record caller freshness and partial-publication checks", "identify SDK loader ownership gap", "file retained watcher growth as B-0005"]
    - revision: 6
      reason: "Verify cancellation classification, retained writer ownership and cross-process progress; record broader gate failures honestly"
      changes: ["record additional red-green regressions and focused race coverage", "refresh generated schema and interface evidence", "record passing drift/contracts components and install proof", "extend B-0004 with changed-suite failures"]
    - revision: 7
      reason: "Close the current verification pass and identify the remaining SDK transaction policy decision"
      changes: ["record final ci-drift and ci-contracts passes", "make SDK partial-publication versus whole-attempt rollback choice explicit"]
    - revision: 8
      reason: "Owner approved per-cell publication for the SDK loader; implement it and record its evidence"
      changes: ["route DefaultResourceLoader reloads through one persistent Host", "record red-green, mutation and baseline comparison evidence", "record the remaining pre-existing gate failures in B-0004", "mark S1.1 and S1.2 implemented in the progress and task tables"]
    - revision: 9
      reason: "Commit the S0/S1.1/S1.2 work, remove the stale foreign worktree and land SPEC-0007 D3"
      changes: ["record the two commits that now carry the previously uncommitted work", "record the removed .worktrees/docs/agents-md-rewrite checkout and the scan failure it cleared", "implement strict deliverAs with red-green evidence and SDK constants", "record S1.5 evidence and the still-open S1.3, S1.4, S1.6 and S1.7"]
    - revision: 10
      reason: "Prepare S1.6 without landing it, and record the one decision it needs first"
      changes: ["record the single implementation anchor for the resume reason", "record the three candidate anchors for the initial flag with a recommendation", "correct the earlier commit count in the S1.5 evidence"]
    - revision: 11
      reason: "S1.6 implemented after the owner chose the Runtime-scoped initial flag"
      changes: ["replace the S1.6 preparation with evidence", "record that the prepared resume anchor did not hold and where resume is decided instead", "record the two delivery defects found and fixed", "mark the resume picker path UNVERIFIED"]
    - revision: 12
      reason: "Fix the test wrapper that hid rustup"
      changes: ["record ddcf654d8 and drop the manual RUSTUP_HOME/CARGO_HOME workaround from the S1.6 evidence"]
    - revision: 13
      reason: "S1.3 and S1.4 landed through delegated workers and were integrated"
      changes: ["replace the 1.3 and 1.4 rows with their real evidence", "record the S1.3 and S1.4 evidence sections", "record the two blockers found during integration: Node has no handler context, and the SDK bundles were stale"]
---

# PLAN-0001 Resolve the Workbench host gaps

Input: `../../analyses/2026-11-04-host-gap-verification/README.md`. Policy: `../../specs/0001-personal-distribution.md` D7–D10.

## Progress

| Slice | Goal | Specs | Status |
|---|---|---|---|
| S0 | Steering, gates, public fork statement | 0001 | implemented; full gate blocked (B-0004) |
| S1 | Reliability fixes and skill roots | 0004 D1–D3, 0007 D1–D3, 0009 | now; 1.1–1.6 implemented, P29 agent-turn reproduction open, 1.7 next |
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
- Gate per task: `make test-changed lint-changed`, `go vet` on touched packages, `-race` on edited packages. Any change to an SDK source (`extensions/sdk*`) also runs that SDK's own suite (`make test-sdk-go`, `test-sdk-py`, `test-sdk-rs`) and `go generate ./...` in that bundle directory, because a stale `.gz` mirror otherwise ships silently (B-0004). Before closing a slice: `make ci-drift ci-contracts`, `make check`.
- A spec decision that turns out wrong is revised in the spec in the same change.

## S0 Steering and gates (now)

| # | Task | Ref | Status |
|---|---|---|---|
| 0.1 | Rewrite root `AGENTS.md` for the fork goal | 0001 D7–D10 | done |
| 0.2 | SPEC-0001 revision 2 (D7–D10) | 0001 | done |
| 0.3 | Child `AGENTS.md` files: scope parity-mandating rules to the parity machinery or replace them with fork rules: `docs/AGENTS.md` (D<N> rule), `docs/parity/AGENTS.md`, `test/parity/AGENTS.md`, `coding/extension/AGENTS.md`, `ai/AGENTS.md`, `agent/AGENTS.md`, `extensions/AGENTS.md`, `internal/AGENTS.md`, `tui/AGENTS.md`, `.github/AGENTS.md`, `automation/AGENTS.md`, `coding/AGENTS.md`, `piglets/AGENTS.md`, `mcp/AGENTS.md`; skill `agents-md` | 0001 D7, D10 | implemented |
| 0.4 | `divergence-guard`: accept `// fork (SPEC-NNNN): <reason>` when `.project/specs/NNNN-*.md` exists; guard test; then update root `AGENTS.md` Code and lint and SPEC-0001 D10 | 0001 D10 | implemented |
| 0.5 | `automation/ci/check-public-claims.py`: make sure it accepts the fork statement of 0.6; adjust rules that require Pi-faithful wording | 0001 D7 | implemented |
| 0.6 | Docs: `docs/personal-fork.md` replaces "Runtime behavior … remain upstream-compatible" with the D7 policy; `README.md` replaces "If PiG behaves differently from Pi, that is either a bug or a documented divergence" with the fork statement and a link to `docs/personal-fork.md`; changelog fragment | 0001 D7 | implemented |
| 0.7 | Outside this repo: pi-accordion resolves a Node executable instead of `process.execPath` (stray billed sessions) | analysis §3 | owner |
| 0.8 | `TestFamilyCoverageReportMatchesCurrentStaticDashboard` (`test/parity/closure/report_projection_test.go:744`) reads the family dashboard from root `AGENTS.md`, which `f0db6902b` removed; it fails on `main`. Read the dashboard from the generated coverage report instead | 0001 D6 | implemented |

Done when: `make check` passes, no child `AGENTS.md` demands Pi faithfulness for fork behavior, public docs state the fork policy.

### S0 verification (04.11.2026)

- Audited all child steering files listed in 0.3; changed only conflicting rules. `extensions/AGENTS.md`, `internal/AGENTS.md`, `.github/AGENTS.md` and `mcp/AGENTS.md` needed no edits.
- Red → green: `TestFamilyCoverageReportMatchesCurrentStaticDashboard` failed on the missing AGENTS dashboard, then passed against the generated family table. The coverage generator now shares its family-table renderer between report and condensed block.
- Red → green: `TestForkMarkerRequiresExistingSpecAndReason`; recorded spec accepted, missing spec and missing reason rejected. Complete divguard unit suite passed.
- Public-claims tests and checker passed without changing checker policy; added fork wording to the supported-statements fixture.
- Targeted `-race` runs passed for fork markers, family projection, coverage report/status/evidence and public-claims tests. `go vet` passed on all four touched Go packages. `make lint-changed` passed with 0 issues. `git diff --check` passed.
- Full divguard + closure race run exceeded the 120 s command deadline; targeted regression race runs completed. Full coverage-package test run fails its CLI test because `.upstream/current` is absent.
- `make test-changed` and `make upstream-mirror` fail pinned npm tarball integrity checks; `make ci-drift ci-contracts` stops at correspondence extraction without `.upstream/current`; `make check` stops at missing Node SDK dependencies. See B-0003. No gate bypass or integrity override was applied.
- Coverage regenerated using `go run -trimpath ./test/parity/cmd/coverage -upstream .upstream/v1.4.0+6c384aa3 -out test/parity/coverage.md`. The extracted pinned source has the manifests this generator consumes; full mirror hydration remains blocked. Generated diff adds only the family table.
- 0.7 remains owner work outside this repository. S0 is not closed; S1–S8 have not started.

### S0 follow-up verification (04.11.2026, 17:13)

- B-0003 resolved: identical built source packed with Node 24.22.1 and npm 12.1.0 matches all seven frozen workspace hashes. Host Node 26.13.0 does not. No integrity values or dependency versions changed; the system Node installation was not replaced.
- Added `pack()` runtime preflight and red-proven `ClosureTests.test_pack_rejects_unpinned_node_before_running_npm`; all 64 pin tests pass. Developer docs explain selecting `.node-version` before oracle builds.
- `make upstream-mirror` succeeds after refreshing the earlier partial mirror. Use a persistent build cache separately from isolated HOME so the installed SDK links outlive the validation command. Preserve installed Rust toolchain discovery separately as well.
- Red → green: `TestScansExcludeNestedWorktrees`; both format-version scanners now exclude `.worktrees/`. Full formatversions suite passes with `-race`.
- `make ci-drift ci-contracts` passes. `make lint-changed` passes with 0 issues.
- `make check` now reaches the full Go tests and fails outside the edited runtime/guard packages: shell/SSH, X11, process discovery, historical timezone names and other nested-worktree scanners. Exact tests and symptoms: B-0004.
- `make test-changed lint-changed` exceeded the 600 s command deadline in the changed-package race group; not a pass. Separate lint passed. S0 remains open; S1–S8 remain pending.

## S1 Reliability fixes and skill roots (now)

| # | Task | Ref | Tests | Docs |
|---|---|---|---|---|
| 1.1 | Generation-unique sockets implemented; original P29 tree-navigation and agent-turn reproduction remains open | 0004 D1 | red-proven `TestSocketGenerationTeardownDoesNotUnlinkSuccessor`; 200-reload stress loop is not a P29 reproduction | `docs/runtime-cell-reload.md`, bundle `runtime-cells.md` |
| 1.2 | Bounded reload handshake, atomic reload, error to `ctx.Reload()`, `/reload`, RPC and the SDK loader; global setting `extensionsReloadHandshakeMs` | 0004 D2 | stuck-peer test (process never registers) keeps old generation serving; SDK per-cell publication test | `settings.md`, `extensions.md` reload section |
| 1.3 | Event-bus cancellation reuses the existing `cancel` frame; Go, Python and Rust handler contexts end, Node keeps its `(data)` signature | 0004 D3 | `TestNativeEventBusCallerCancellationEndsListener` (Go/Python/Rust, isolated and shared-ok), `TestEventBusCancelledDispatchDoesNotStartListeners`, `TestEventBusCancelledCallerStillRegistersAndUnregisters` | `extensions.md` event bus, `extension-api.md`, `extension-api-parity.md` |
| 1.4 | `nextTurn` bound to session and leaf; discard + `queued_message_discarded` on navigation/new/fork/resume/switch, including in-place slash forks | 0007 D1 | `TestNextTurnLeafBinding` (ordinary, errored, aborted history; keep and discard cases), interactive send guard | `events.md` both trees |
| 1.5 | Strict `deliverAs`; SDK constants | 0007 D3 | `TestSendMessageRejectsUnknownDeliverAs`, `TestSendUserMessageRejectsUnknownDeliverAs`, `TestSendUserMessageRejectsNextTurn`; unknown value fails the call in every SDK through the shared Session action | `extensions.md`, `extension-api.md` |
| 1.6 | `session_start` `resume` for continued sessions; `initial` flag | 0007 D2 | `--continue`, `--resume`, `--session`, fresh start | `events.md`, changelog `### Upgrading` |
| 1.7 | `skills.roots[]`: settings object form, legacy migration fix, tokens, precedence, trust; `make generate-schemas` | 0009 | migration keeps `roots`/`paths`; collision order; untrusted project ignores project roots; `HOME`-independent isolation | `skills.md`, `settings.md`, `settings-examples.md`, root `AGENTS.md` test-isolation bullet |

Done when: all S1 tests red-proven where they are bug fixes, `make check` passes, docs and changelog fragments merged.

### S1.1 evidence and next decision (04.11.2026)

- S0 remains open under B-0004. Its remaining nested-worktree scanners and named-process fixture are now red → green; the timezone and X11 failures reproduce, and the earlier ELOOP did not reproduce in the first narrow run. No full gate is claimed green.
- Implemented generation-unique sockets across packed, isolated and park-ack starts. Removed pre-listen deletion of stable paths and reserved full-width counters in the Unix socket budget. Both docs trees and the changelog fragment are updated; SPEC-0004 D1 revision 2 records the implementation and evidence limits.
- Red → green: `TestSocketGenerationTeardownDoesNotUnlinkSuccessor` demonstrates old teardown unlinking the successor's bound socket before the fix. `TestPackedGoReloadLoopNeverHangs` ran 200 reloads with factory-time parser calls, bus subscription and a tool call each time. Retained-process identity and shutdown cleanup pass. The initial loop also passed before the fix, so it does not reproduce P29; tree navigation plus a real agent turn remains unverified.
- The loop's per-admission timeout explicitly retains the session context as `runtimeParent`: cancelling the operation context must not terminate the retained process being measured. The added retained-identity assertion initially exposed that fixture lifetime mistake; it was corrected without removing the assertion.
- Focused socket/reload tests and the broader reload/retention selection pass with `-race`; `go vet`, `make lint-changed` and `make ci-drift ci-contracts` pass. The teardown regression uses the production `ListenExtension` boundary (avoiding a false-positive TempDir socket lint warning). Warm-cache reload benchmark: 3 iterations, 4,547,431 ns/op, 3,601,690 B/op, 5,657 allocs/op; CPU and allocation profiles captured locally. No before/after performance claim.
- The owner approved flat `extensionsReloadHandshakeMs`. SPEC-0004 D2 and the settings wire, clone, serialization, global accessor, CLI wiring and schema now use it; the `extensions` path array is unchanged. D5 setting names are flat in the draft spec too.

### S1.2 progress and remaining work (04.11.2026)

- Implemented `Host.Reload` admission/connect/register budgets (default 15000 ms), excluding builds and admission-order waits. Isolated and packed paths preserve process lifetime independently of the operation deadline. Initial loads remain cancellation-only. A retained-process probe/write uses the same cancellation mechanism; its one owned writer is drained when the process ends and an expired write prevents that process being claimed again.
- Red → green: `TestHost_Reload_FailedReplacementKeepsOldExtension`; `TestReloadPackedCellFailureKeepsServingGeneration` (old tools, provider set, MCP URL, virtual-model metadata and crash ownership); `TestReloadReportsUnresolvedConfigsAndKeepsServingGenerations` (source-based unresolved identities). The connect/register deadline test was added after the first timeout implementation and mutation-checked: removing the deadline made both cases wait for the outer deadline. `TestReloadDeadlineIncludesStalledRetainedAdmission` red-proved an unbounded stdin write before the handshake wait.
- Staged providers, OAuth, native providers, MCP servers and virtual models publish only for accepted cells. Failed cell siblings are discarded without factory replay; crash allocation IDs remain monotonic while rollback restores the live generation's crash ownership. Ready activation precedes registry publication. Partial success carries `extension.ErrReloadPartiallyApplied`; loader/Session/headless/interactive adapters must publish the actual set and still return the error.
- `TestSessionReloadKeepsTheOldRuntimeActiveOnFailure` was red → green. API registry cleanup now occurs after successful loading and keeps registrations made during the load, including source IDs. Partial reloads do not reset the serving registry. Real print/RPC tests pass for a never-registering Node factory, an invalid configured budget, an ordinary resource-resolution failure, and mixed healthy Node/failed Go cells. The mixed test checks new healthy tools, retained failed-cell tools and removal of replaced tools before the failure reaches its caller.
- Existing failure-dropping/stale-instance tests are being rewritten to SPEC-0004 D2. The first broader Host reload selection exposed only old-contract expectations and one synthetic generation fixture; the first interactive selection exposed four old-contract expectations. Follow-up gate status is recorded below only after observing it.
- `make generate` initially stopped because the schema tests were renamed and active unit-evidence locators still named them. Those locators now name the fork-contract tests. The settings schema is checked against its pinned origin plus the exact fork property, rather than falsely claiming byte equality for the changed settings schema.
- SDK path implemented after the owner approved per-cell publication: an ordinary `DefaultResourceLoader` reload now runs `Host.Reload` on the Host the loader already owns instead of creating a new Host, so it inherits the budget, per-cell rollback and `ErrReloadPartiallyApplied`. A trust bootstrap pass still starts a fresh Host, and a caller-supplied `DefaultResourceLoaderOptions.Runtime` keeps its identity. The aggregate runtime of a reload is a fresh `ReloadSnapshot` view over the same live provider registration state, so factory APIs from the previous generation go stale while MCP, provider and virtual-model registries stay bound. Compiled-in and inline factory failures keep their serving extension record through `extensionset.Loader.Previous`, and factory APIs of extensions a reload removed are retired.
- Red → green: `TestDefaultResourceLoaderPublishesHealthyFactoriesAndKeepsFailedGenerations` (both `Runtime` ownership cases) failed first with a nil reload error, then with the serving command dropped when `Previous` retention was removed; both mutations were reverted. The test also asserts the Host survives the reload, the supplied runtime keeps identity and the omitted one gets a new snapshot.
- Baseline comparison for the `coding` and `coding/cli` suites used a throwaway `HEAD` worktree, since these failures are outside the reload change: `TestStartupAllocationBudget`, `TestRPCUsageGoldenMatchesPi`, `TestRawPigContainsNoPlatformAPIPaths` and `internal/codingagent.TestSessionSearchRegexMatchesPi` fail identically on unmodified `5972948ba`. None is claimed fixed or caused here; all are recorded in B-0004.
- Resource-lifetime probe: one retained OS process has 201 `watchedPacked` states after 200 reloads, though final shutdown drains it. Filed B-0005 with the measured reproduction and required crash/retention constraints; no bounded steady-state resource claim. P29's original tree-navigation plus agent-turn reproduction and B-0004 full-gate failures also remain open.
- Further red → green: `TestReloadDeadlineDoesNotSelectAnotherRuntime` exposed cancellation/deadline errors being classified as runtime startup incompatibility; deadline tests exposed incorrect SDK-rebuild advice. Both are corrected. `TestRetentionProbeDoesNotBlockAnotherProcessClaim` timed out with an unrelated claim waiting on the retention mutex behind blocked stdin I/O; probes now reserve an owned reference and perform I/O outside the mutex.
- Admission writes now use one `sync.Once`-started, `WaitGroup`-owned worker per process. `TestProcessShareCancellationReturnsBeforeTheOwnerDrainsItsWriter` proves caller cancellation returns before blocked I/O drains, while owner release waits for that drain. No lint suppression was added.
- `TestDivergenceCitationsExcludeSpecDecisions` red-proved the documentation gate confusing qualified fork decisions with upstream D<N> entries. Scoped single decisions and ranges are excluded; adjacent upstream citations remain checked. The packed Python shape-skew test now requires the caller error and still checks the exact unknown-field/rebuild diagnostics in the report.
- Focused reload/retention/publication/setting selections pass with `-race` in Host, CLI, Session, API registry and interactive packages. `make generate`, `make generate-schemas`, `make interface-gaps-update`, `make lint-changed` and touched-package `go vet` passed. Generated mapping changes are limited to this work’s API-registry/reload evidence and production references.
- Final `make ci-drift`, `make ci-contracts`, `make lint-changed` and touched-package `go vet` passed. Drift used an alternate index containing the new files because the source-zip gate omits untracked files; the real index remains unchanged. An earlier ordinary-index invocation failed for those missing files, not because the complete source failed to compile. No complete slice gate is claimed.
- The latest `make test-changed` exceeded its 1200 s command limit and exposed ELOOP, an MCP configuration data race, X11 disagreement, a provider cancellation mismatch, four inner 600 s package timeouts and a Rust cancellation wait failure. Two in-scope failures (scoped spec citations and shape-skew reload error expectations) were fixed and rerun. The Rust cancellation test passed once narrowly without a fix; the broad failure remains open. Details: B-0004.
- Post-change warm-cache reload benchmark: 3 iterations, 4,988,982 ns/op, 3,590,306 B/op, 5,692 allocs/op, with CPU/allocation profiles captured. No comparative performance or bounded steady-state resource claim; B-0005 remains open.
- Remaining 1.2 caveat: SDK reload coverage exercises inline factories, built-ins and file extensions through the loader. Compiled-in `extensionFactories` still go through the same inline path, but a real agent-turn reload against a compiled-in factory has not been driven end to end.
- S1.3–S1.7 and S2–S8 have not started. No complete slice gate or commit is claimed by these implementation notes.

### S1.5 evidence (04.11.2026)

- The previously uncommitted S0/S1.1/S1.2 work is committed as `37eb840a5`; SPEC-0007 D3 is `e14430cfa` plus its release-note commit. The tree is clean after both.
- `.worktrees/docs/agents-md-rewrite` was removed by owner decision. Its branch had no commits outside `main` and the checkout had no local changes, so nothing was lost. `coding/cli.TestRawPigContainsNoPlatformAPIPaths` passes after the removal; it had been scanning that foreign checkout.
- Red → green: `TestSendMessageRejectsUnknownDeliverAs`, `TestSendUserMessageRejectsUnknownDeliverAs` and `TestSendUserMessageRejectsNextTurn` failed first (each returned nil and delivered the message), then passed. `TestKnownDeliverAsValuesAreAccepted` keeps the empty value and the three defined modes working.
- One validation sits at the shared boundary (`extension.ValidateDeliverAs` called from `SendMessage`, `SendMessageStaged`, `sendCustomMessage`, `SendUserMessage` and `SendExtensionUserMessageStaged`), so the in-process, subprocess and session-graph bindings all reject through the same path. An unknown value fails before any queueing or dispatch work.
- Go, Rust and Python SDKs export `DeliverSteer`/`DeliverFollowUp`/`DeliverNextTurn` (`DELIVER_*` in Rust and Python). The TypeScript surface already restricts the literal union through the pinned Pi declarations and exports no runtime values, so it gained no constant.
- RPC-mode user prompts are not validated by this change: `PromptOptions.StreamingBehavior` still comes from CLI flags, and SPEC-0007 D3 names only the two extension calls. A `rpc.md` entry therefore was not added; `extensions.md` and `extension-api.md` carry the contract.
- Verified: `go build ./...`, `go vet ./coding ./coding/extension`, the deliverAs and queue/custom-message/prompt selections in `./coding`, `cargo check --all-targets` in `extensions/sdk-rs` (pre-existing warnings only), Python parse of `pig_sdk/__init__.py`, `make lint-changed` (0 issues), `make changelog-check`, `make ci-drift`.
- `make ci-contracts` passes its correspondence and Porter packages, then fails the Porter unit test that installs a failing repository-local `pre-commit` hook, because this host sets a global `core.hooksPath`. `make porter-check` passes all 26 tests with `GIT_CONFIG_GLOBAL=/dev/null`. Recorded in B-0004 as a host-configuration effect, not a defect of this change.
- Not verified in this pass: no cross-SDK conformance row drives an unknown `deliverAs` through Node, Rust and Python processes. The rejection is proven at the Session boundary that every binding shares, not through each runtime. Record it as a conformance gap for S1.5 follow-up rather than a covered row.
- Still open in S1: 1.7 (`skills.roots[]`) and the P29 agent-turn reproduction. 1.3, 1.4 and 1.6 landed (below). S2–S8 have not started.
- Anchors for the next items: `events.cancel` (1.3) is the bus in `coding/extension/host/subprocess/event_bus.go`; the `nextTurn` leaf binding (1.4) is `agent.pendingNextTurn`.

### S1.3 evidence: cancellable event dispatch (05.11.2026)

Commits `f01e57c2e` (fix), `b9abc0da2` (changelog), `95e40703e` (docs). Implemented by a delegated worker after one refusal as "mis-sized"; the resume carried the lead decision that the host change and its proof are one unit.

- No new wire method. `handleEventBusCall` no longer shields `events.emit` from the caller's context (`context.WithoutCancel` now applies only to `events.on`/`events.off`, so registration stays atomic for a cancelled caller), `dispatchEvent` starts no further listener once the caller's context is done, and `Conn.request` already sends the existing `MsgCancel` frame for each outstanding dispatch. Node listeners additionally had their emission waiters released on cancellation.
- Go factory-pump and pre-ready dispatches now arm a request scope instead of `context.Background()`, so they cancel like registered handlers.
- Red → green, mutation-checked by me on the merged tree: restoring `context.WithoutCancel(ctx)` for every method turns `TestNativeEventBusCallerCancellationEndsListener/go/{isolated,shared-ok}` red after 240 s; the unmutated selection is green.
- Node/TypeScript is a documented limit, not a gap: Pi's handler signature is `(data)`, so there is no context to end and JavaScript cannot preempt a running handler. SPEC-0004 D3 records it, the changelog says it, and the docs say it.
- Docs: `extensions.md`, `extension-api.md` and the `pi.events` row of `extension-api-parity.md` now state the contract. `docs-drift` caught that `events.md` (added by 1.4) had no provenance classification and no navigation entry; `docs.json` gained both.
- Not covered: the cross-SDK row cancels an emitting *command* request. Cancellation of an emit made inside a tool call that the user aborted is the same host path but has no row of its own.

### S1.4 evidence: branch-bound `nextTurn` (05.11.2026)

Commits `ee8403488` (fix), `cb7bf2cff` (in-place slash forks), `8694de46c` (changelog and bug ledger). Also delegated, also after one "mis-sized" refusal.

- A queued message records its session and leaf at queue time; navigation and every replacement path drop messages whose leaf is no longer current, including the in-place slash forks (`internal/codingagent`), and each drop emits `queued_message_discarded {customType, reason}` plus a `PIG_DEBUG` line.
- Red → green, mutation-checked by me on the merged tree: removing the discard emission turns `TestNextTurnLeafBinding` red for the ordinary, errored and aborted history cases; the keep case (same leaf) still passes, which is what makes the discard meaningful rather than a blanket clear.
- The extension event is typed (`extension.QueuedMessageDiscardedEvent`) and reached every SDK surface without a wire change, since it travels as an ordinary extension event.
- Not covered: the interactive `--resume` picker path still has no end-to-end test (carried from S1.6).

### Integration notes (05.11.2026)

- `feat/s1-4-nextturn-leaf` was merged with `--no-ff` as `fb244b67b`; S1.3 was committed directly on `main`. No conflicts.
- The embedded SDK bundles (`*.go.gz`, `*.py.gz`, `*.rs.gz`) were stale before this work: the Python `__init__.py.gz` still lacked the S1.5 `DELIVER_*` constants. S1.3 regenerated them with `go generate` in each SDK; after the merge I regenerated again so the mirror matches the source (`29dfffd16`).
- Gates on the merged tree: `make generate`, `make docs-drift`, `make ci-drift`, `make changelog-check`, `go vet` on the touched packages, `-race` on `agent`, `coding` and `coding/extension/host/subprocess`, the SDK and host selections, and `TestNativeEventBusCallerCancellationEndsListener|QueuedMessageDiscarded|ApiOnEvents` in `test/extension-conformance` (green, 88 s). `make ci-contracts` still fails only the known `porter-check` hook test (B-0004).
- Delegation note for the next slice: both workers refused the first brief as mis-sized within seconds. The resume-with-lead-decision countermeasure from `selecting-worker-model` worked again, and cost less than splitting.

### S1.6 evidence: `session_start` `resume` and `initial` (05.11.2026)

Owner decision 04.11.2026: `initial` is Runtime-scoped; SPEC-0007 D2 revision 2 records the alternatives and why they lost. Commits `453651743` (fix), `f076665fd` (changelog), `3cc9957f2` (context), `314a3a615` (perf annotation).

```
launch                         reason   initial
pig "prompt"                   startup  true
pig --session <missing file>   startup  true
pig --session <existing file>  resume   true
pig --continue                 resume   true
  then /new                    new      absent
```

- The prepared anchor did not hold. `SessionOptions.resumed` is set only on the `Runtime.startSession` path, which no CLI mode uses: all four modes build a `SessionManager` first and pass it to `CreateAgentSessionRuntime`. The CLI now decides `resume` in `resolveStartupSessionSelectionIn` by checking that the selected file exists, before anything can create it, and passes the event to `CreateAgentSessionRuntime`. `NewAgentSessionRuntime` marks the first event `initial`.
- Two delivery defects surfaced. Print, JSON and RPC re-emitted through `EmitSessionStartTransition(reason, previous)`, which dropped every other field: `Session.EmitSessionStartEvent(ctx, event)` forwards the whole event and now uses the mode's context instead of `context.Background()`. Interactive mode bound its first Session with a fixed `startup` event; `InteractiveOptions.SessionStartEvent` now carries the Runtime's.
- Red then green, each mutation-checked: `TestRuntimeMarksOnlyItsFirstSessionStartInitial` (`./coding`), `TestSessionStartReportsResumeAndInitialAtProcessStart` (print and RPC binaries), `TestInteractiveSessionStartReportsResumeAndInitial` (PTY, fresh and `--continue`). Removing the `resumesExisting` assignment, the `Initial = true` line or the interactive option each turns its test red.
- Cross-SDK: a second `session_start` payload row with `initial: true` passes for Go, Node, Rust and Python in isolated, packed and fused realizations. No SDK type changed: handlers receive the event as JSON, and the SDK `SessionStartEvent` types are child-session inputs, where `initial` is host-computed and never accepted.
- Tests that pinned the old first event were rewritten to D2: `TestAgentSessionRuntimeKeepsTheSessionStartEventUpstream`, `TestRuntimeOriginalNewResumeEvents`, `TestRPCRunEmitsSessionStartOnce`.
- Gates: `make generate`, `make lint-changed` (0 issues), `go vet`, `make changelog-check`, `make ci-drift` (port-lint baseline down by two `ctxbackground` entries, one cold `os.Stat` annotated for perf-lint). `make ci-contracts` fails only the known `porter-check` hook test (B-0004).
- Environment note: `with-isolated-pig-home.sh` hid rustup, so Rust rows needed `RUSTUP_HOME` and `CARGO_HOME` exported by hand. Fixed in `ddcf654d8` (05.11.2026); the wrapper now runs Rust rows without a manual setup (B-0004).
- Open: the `pig --resume` picker path was not driven end to end. It reaches the same `CreateAgentSessionRuntime` call with a selected existing file, so it should report `resume`, but no test exercises it. UNVERIFIED.

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
