---
document:
  type: "backlog-item"
  id: "B-0004"
  title: "Resolve remaining host-gap full-check failures"
  description: "Failures exposed after repairing the pinned Pi oracle setup"
  status: "open"
  revision: 7
  created_at: "2026-11-04T17:13:00+01:00"
  updated_at: "2026-11-05T10:05:00+01:00"
  created_by: "coding-agent"
  updated_by: "coding-agent"
  kind: "bug"
  priority: "P1"
  scope: ["ai/**", "coding/extension/**", "coding/mcpext/**", "env/**", "internal/codingagent/**", "internal/nativeplatform/**", "internal/perfgate/**", "internal/testenv/**", "internal/yaml12/**", "test/extension-conformance/**", "test/parity/cmd/sdksurface/**", "test/parity/interface-closure/autobind/**", "test/upstream-parity/**"]
  changelog:
    - revision: 1
      reason: "Full S0 gate reaches tests after B-0003 resolution"
      changes: ["record failing tests without suppressing or retrying them"]
    - revision: 2
      reason: "Continue S0 scanner and process-fixture repairs without blocking all S1 work"
      changes: ["record scanner and process-fixture red-green results", "record timezone and X11 reproduction", "keep remaining full-check failures open"]
    - revision: 3
      reason: "S1 changed-package gate exposes additional failures and timeouts"
      changes: ["record MCP configuration data race and provider cancellation mismatch", "record four package timeouts and Rust cancellation wait", "distinguish repaired fork-contract failures from remaining gate failures"]
    - revision: 4
      reason: "Separate pre-existing failures from the SDK reload change with a pristine-HEAD baseline"
      changes: ["record startup allocation, RPC usage golden, raw-pig scan and regex oracle failures", "record their reproduction on unmodified HEAD"]
    - revision: 5
      reason: "Removed the stale foreign worktree and ran ci-contracts on the committed tree"
      changes: ["record the raw-pig scan passing after the worktree removal", "record the host core.hooksPath failure in the Porter commit-failure test"]
    - revision: 6
      reason: "Record the isolation wrapper hiding rustup"
      changes: ["add the RUSTUP_HOME and CARGO_HOME requirement for Rust rows under with-isolated-pig-home.sh"]
    - revision: 7
      reason: "The rustup isolation gap is fixed"
      changes: ["mark the isolation-wrapper rustup section resolved by ddcf654d8"]
---

# Remaining full-check failures

## Reproduction

Run `make upstream-mirror` and `make check` with Node 24.22.1 and npm 12.1.0 selected on PATH. Isolate agent credentials through `automation/ci/with-isolated-pig-home.sh`; preserve compiler/package caches separately from HOME, including Rust toolchain discovery.

Observed full-run failures on 04.11.2026:

| Test | Failure |
|---|---|
| `coding/extension.TestExecCommandStartsTheTimeoutAfterTheSpawnAsPiDoes` | Pi oracle spawn returns ELOOP |
| `env.TestSSHBootstrap/starts_the_daemon_through_the_login_shell_only_when_asked` | daemon exits 127 before ready |
| `env.TestOpenTerminalRunsTheUsersShellAsALoginShellByDefault` | Pig and Pi both expose HOME from the login shell instead of the test's expected unset value |
| `internal/nativeplatform.TestX11DisplaySpellingsConnectAsPiDoes/unix:0.-1` | Pig returns a root window; xcb_connect reports error |
| `internal/perfgate.TestDescendantFindsAGrandchildByName` | descendant lookup returns 0 |
| `internal/testenv.TestTimingGatesAreListed` | scanner enters `.worktrees/` and treats another checkout's tests as this checkout's |
| `internal/yaml12.TestDateStringMatchesNode` | Eire historical timezone name differs from Node |
| `test/upstream-parity.TestKeyReleaseChokePoint`, `TestInputDispatchFiltersKeyReleases` | scanners enter `.worktrees/` and report another checkout's paths |

The format-version inventory had the same nested-worktree issue and was fixed with a red-proven regression in S0. Other scanner paths above remain unchanged. Root causes of the non-scanner failures are not established; do not label them flakes or weaken their assertions.

`make test-changed lint-changed` exceeded the 600 s command deadline in its changed-package race group. Several edited-package suites completed, but the target did not finish. A separate `make lint-changed` passed.

## Follow-up (04.11.2026)

- `TestTimingGatesAreListed`, `TestKeyReleaseChokePoint` and `TestInputDispatchFiltersKeyReleases` reproduced their nested-worktree failures, then passed after their walkers excluded `.worktrees/`.
- `TestDescendantFindsAGrandchildByName` reproduced alone. Its shell fixture could replace the named shell with its last `sleep` command. Starting the sleep in the background and using the shell's `wait` builtin keeps the named process alive; the test passes with no production change or retry.
- `TestExecCommandStartsTheTimeoutAfterTheSpawnAsPiDoes` passed in its first narrow follow-up run without a fix; the prior full-run ELOOP is not resolved or dismissed as a flake.
- The X11 failure reproduces with installed libxcb 1.17.0, while the source commentary cites 1.15. The timezone failure reproduces on the pinned Node: `Eire` is accepted as `Europe/Dublin` by Intl but Date prints a fixed Greenwich Mean Time name for it; regeneration differs only for `Eire`. Neither the oracle nor generated data was changed to hide these differences.

## S1 changed-package run (04.11.2026, 21:36)

`make test-changed` exceeded the 1200 s outer command limit. Its main race group reported eight failing packages; the subsequent bounded group reported additional failures. Do not treat increasing a timeout as a repair.

| Test or package | Observed failure |
|---|---|
| `coding/extension.TestExecCommandStartsTheTimeoutAfterTheSpawnAsPiDoes` | `spawn sleep ELOOP` reproduced in the broad run |
| `coding/mcpext.TestMcpManagerEnablesAServerInTheBackground` | race detector reports `Extension.saveConfig` writes overlapping command/TUI configuration reads |
| `internal/nativeplatform.TestX11DisplaySpellingsConnectAsPiDoes/unix:0.-1` | prior X11 disagreement repeats |
| `ai.TestMistralCancelMatchesPi/pending-body` | run 2 emits `[start]`, expected `[start error]` |
| `ai`, `internal/codingagent`, `test/parity/cmd/sdksurface`, `test/parity/interface-closure/autobind` | each hits its inner 600 s package timeout; autobind stack is in `go/packages` parsing |
| `test/extension-conformance.TestCancelReasonReachesEveryNativeSDK/subprocess-rust` | wait condition not satisfied within 120 s; one narrow race run later passed without a fix, so the broad failure remains unresolved |

Two additional failures were within the current fork change and repaired: the docs gate read `SPEC-0004 D1` as upstream divergence D1; the packed Python register-shape test still expected a nil reload error. Their focused race regressions pass while retaining the upstream-ID and unknown-field/rebuild assertions.

The complete `coding` race package passed in the broad run; focused Host, CLI, Session, API registry and interactive reload selections also pass. These narrower results do not close this item.

## Pre-existing failures confirmed against unmodified HEAD (04.11.2026, 23:05)

A throwaway `git worktree` at `5972948ba` (no local changes; removed afterwards) was tested with the same isolated home and toolchain. These fail identically there, so they are not caused by the S1.2 work:

| Test | Baseline result |
|---|---|
| `coding/cli.TestStartupAllocationBudget` | warm first prompt 4,574,080 bytes in 26,174 objects against a 2,589,000-byte budget |
| `coding/cli.TestRPCUsageGoldenMatchesPi` | live stream differs from `rpc-usage-pi.json` and `rpc-usage-extension-pi.json` |
| `coding/cli.TestRawPigContainsNoPlatformAPIPaths` | reports `.worktrees/docs/agents-md-rewrite/bench/k8s/cmd/k8sbench/kube.go`, the nested-checkout walker class above |
| `internal/codingagent.TestSessionSearchRegexMatchesPi` | 5 of 3000 seeded patterns differ from Pi; duplicate capture-group names and quantifier-error precedence |
| `coding/cli.TestRPCInputEndResponsePrecedesShutdownNotificationComparedWithPi` | passes alone in 6 s and completed in the second full run; the first full run hit the package timeout under load |

A full `coding` run (all packages in that directory) and a focused `-race` reload selection passed after the change.

## Worktree removal and ci-contracts (04.11.2026, 23:40)

- `.worktrees/docs/agents-md-rewrite` was removed by owner decision: its branch had no commits outside `main` and the checkout had no local changes. `coding/cli.TestRawPigContainsNoPlatformAPIPaths` now passes, so that row is cleared by the removal rather than by a code change.
- `make ci-drift` passes on the committed tree. `make ci-contracts` fails at `porter-check` in `automation/porter/test_task_runner.py::test_a_failed_commit_fails_the_task_and_the_runner_and_is_not_marked_done`, which installs a failing `.git/hooks/pre-commit` in a scratch repository and expects the commit to fail.
- Cause, reproduced outside the repository: the host sets a global `core.hooksPath` (`/home/furbyhaxx/.config/git/hooks`), so a repository-local `.git/hooks/pre-commit` is never consulted and the commit succeeds. `git -c user.email=t@t -c user.name=t commit` in a fresh repository with a failing `pre-commit` returns 0 with that configuration. This is a host configuration effect on a test, not a defect in the Porter runner or in this change, and it will fail any test that relies on a repository-local hook.
- Not repaired here: neutralizing the global setting for CI would change the host's git configuration, which is outside this repository. Verified: `make porter-check` passes all 26 tests with `GIT_CONFIG_GLOBAL=/dev/null` and `GIT_AUTHOR_*`/`GIT_COMMITTER_*` set, so the runner's behavior is correct and only the host's hook discovery differs.

## Isolation wrapper hides rustup (05.11.2026)

- `automation/ci/with-isolated-pig-home.sh` points `HOME` at a scratch directory, so rustup's default `~/.rustup` and `~/.cargo` disappear and every Rust extension row fails with "rustup could not choose a version of cargo to run". Seen in `TestJSONModeWritesNothingSessionStartEmittedBeforeTheWriterSubscribesInEverySDK/rust` and the RPC sibling.
- Resolved by `ddcf654d8` (05.11.2026): the wrapper pins `RUSTUP_HOME` and `CARGO_HOME` to the caller's directories before it moves `HOME`, as it already did for the Go caches. A caller's explicit value wins; a missing directory stays unset. With neither variable exported, both Rust rows failed before the change and pass after it. The remaining sections of this item stay open.

## Acceptance
Reproduce each failure narrowly, fix its owner with regression evidence, then rerun `make check` and the complete changed-package target without retries, timing overrides or skipped tests. PLAN-0001 S0 remains open until its gates pass.
