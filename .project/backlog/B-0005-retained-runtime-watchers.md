---
document:
  type: "backlog-item"
  id: "B-0005"
  title: "Retained runtime watchers grow with each reload"
  description: "One live OS runtime accumulates a packed process state and watcher for every replacement generation until shutdown"
  status: "done"
  revision: 2
  created_at: "2026-11-04T19:55:21+01:00"
  updated_at: "2026-11-05T23:25:00+01:00"
  created_by: "coding-agent"
  updated_by: "coding-agent"
  kind: "bug"
  priority: "P2"
  scope: ["coding/extension/host/subprocess/**"]
  changelog:
    - revision: 1
      reason: "Resource-lifetime inspection during PLAN-0001 S1.2"
      changes: ["record measured retained watcher growth and ownership constraint"]
    - revision: 2
      reason: "One watcher per retained process; crash attribution follows the current generation"
      changes: ["record the red 201-state / green 1-state reproduction", "record the node-recovery follow-on found by the full package", "record the lifetime and race evidence"]
---

# Retained runtime watchers grow with each reload

`startGoPackedCellOnce` in `coding/extension/host/subprocess/packed_go.go` creates a new `packedProcessState` for a retained admission, inserts it into `Host.watchedPacked`, and starts `watchPackedProcess`. A retained process does not exit when a generation retires, so its earlier watchers and states remain until the process ends.

Observed with a temporary log immediately before shutdown in `TestPackedGoReloadLoopNeverHangs`: one retained OS process and 201 watched process states after 200 reloads. Command: `automation/ci/with-isolated-pig-home.sh go test -trimpath -v ./coding/extension/host/subprocess -run '^TestPackedGoReloadLoopNeverHangs$' -timeout 90s`. The test passed, including process exit and socket-directory removal after shutdown. The temporary log was removed. This is a live retained-state count, not a before/after memory or CPU claim.

Target: watcher/state ownership bounded by live OS processes, not historical reload count. Preserve immutable generation IDs, cancellation, crash attribution, old-generation command drain, and retained processes shared across successive Hosts. A watcher must dispatch a crash to the currently published state of its OS process; watching only its original immutable generation would discard later crashes.

Add a red resource-lifetime assertion and test ordinary reload churn, failed reload rollback, cross-Host retention, crash recovery and final shutdown. Profile live heap and goroutines before claiming the growth fixed. Do not resolve this by recycling generation IDs or suppressing crash reports.

## Implementation (05.11.2026)

- `processShare` owns `current` (newest admitted generation) and the process's single `watcherDone`. `startGoPackedCellOnce` starts a watcher only for the spawn state; a retained admission publishes its state as `current` and advances `packedCellActiveGeneration` under the host lock.
- `watchPackedProcess` reads `current` for crash attribution, checks the current generation's cancellation, and treats the retiring spawn state as a clean stop only when its share is actually released. The spawn generation's `stopping` flag alone no longer suppresses a crash in a live retained process.
- The full package exposed a node-recovery regression in the first attempt: recovery rejected the new generation as stale. Publishing the active generation fixed it; `TestNodeCrashAfterRetainedReloadRecoversWithoutReplay` passes again.
- Red proof: `TestPackedGoReloadLoopNeverHangs` against `HEAD` failed with `retained reloads left 201 watched process states, want 1`; with the fix it passes at 1.
- Evidence: full `./coding/extension/host/subprocess` package 274.5 s, churn/crash/goroutine/retention group under `-race` 7.9 s, `make port-lint divergence-guard` clean.
- No memory or CPU improvement is claimed; this is the live retained-state count the item measured.
