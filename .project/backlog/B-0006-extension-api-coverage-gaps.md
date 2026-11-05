---
document:
  type: "backlog-item"
  id: "B-0006"
  title: "Extension API coverage gaps found while landing SPEC-0004 D3 and SPEC-0007 D1-D3"
  description: "Three behaviours shipped on contract-level evidence only: no cross-SDK conformance row for strict deliverAs, no row for an emit inside an aborted tool call, and no end-to-end test for the interactive --resume picker"
  status: "done"
  revision: 2
  created_at: "2026-11-05T14:21:09+01:00"
  updated_at: "2026-11-05T23:55:00+01:00"
  created_by: "coding-agent"
  updated_by: "coding-agent"
  kind: "debt"
  priority: "P2"
  scope: ["test/extension-conformance/**", "coding/cli/**", "coding/**"]
  changelog:
    - revision: 1
      reason: "PLAN-0001 S1.3-S1.6 landed with three behaviours proven below the row each contract names"
      changes: ["record the three gaps with the evidence that exists today", "name the row each one needs"]
    - revision: 2
      reason: "All three rows implemented and green"
      changes: ["record the per-SDK strict deliverAs row and Node's void API", "record the aborted-tool emit row", "correct the picker row to resume initial=false per SPEC-0007 D2"]
---

# Extension API coverage gaps from PLAN-0001 S1

Every behaviour below is implemented and has evidence at one layer, but not the layer its spec decision or the slice's own "tests" column names. Each is a missing row, not a defect: nothing here reports broken behaviour today.

## 1. Strict `deliverAs` has no cross-SDK row (SPEC-0007 D3)

- Current: rejection is proven once, at the shared Session boundary every binding uses (`extension.ValidateDeliverAs`, called from `SendMessage`, `SendMessageStaged`, `sendCustomMessage`, `SendUserMessage` and `SendExtensionUserMessageStaged`). Tests: `TestSendMessageRejectsUnknownDeliverAs`, `TestSendUserMessageRejectsUnknownDeliverAs`, `TestSendUserMessageRejectsNextTurn`, `TestKnownDeliverAsValuesAreAccepted` in `./coding`.
- Gap: no conformance row drives an unknown value through a Node, Rust or Python process, so an SDK that silently substituted a default for the value would still pass every gate.
- Target: one `test/extension-conformance/` row per native SDK that calls `send_message` with `deliverAs: "bogus"` and requires the call to fail with the host's error, with a value distinguishable from each SDK's fallback.

## 2. Cancellation of an emit inside an aborted tool call has no row (SPEC-0004 D3)

- Current: `TestNativeEventBusCallerCancellationEndsListener` cancels an emitting **command** request and checks the Go, Python and Rust listener's reason, later-listener suppression and snapshot cleanup in both placements. Separate Go wire tests cover cancellation before `Run` and before ready.
- Gap: an extension that emits from inside a tool call the user aborted reaches the same host path with a different parent (`tool_call` rather than `command`), and no test drives that parent. Recorded in PLAN-0001 under S1.3 evidence.
- Target: a row that registers a tool which emits, runs it, cancels the turn mid-call, and requires the listener's context to end.

## 3. The interactive `--resume` picker is untested end to end (SPEC-0007 D2)

- Current: `--continue`, `--resume <file>`, `--session <existing file>`, a fresh start and a `--fork` start are covered by `TestSessionStartReportsResumeAndInitialAtProcessStart` (print and RPC) and `TestInteractiveSessionStartReportsResumeAndInitial` (PTY). The picker path is not: it selects an existing file through the same `CreateAgentSessionRuntime` call and should report `resume`, but nothing drives it. Recorded as UNVERIFIED in PLAN-0001 S1.6 and S1.4 evidence.
- Target: a PTY test that opens `/resume`, selects an existing session from the picker, and asserts `session_start reason=resume initial=false` (SPEC-0007 D2: a replacement is not the Runtime's first `session_start`).

## Why P2

No user-visible behaviour is wrong. The risk is that the same class of SDK-level divergence D3 was written to prevent reappears without a failing gate, which is exactly what the stale-bundle incident (B-0004) showed: the check existed and nothing ran it.

## Implementation (05.11.2026)

1. `TestConformance_ExtensionAPIStrictDeliverAs` adds `send_message_bogus_deliver_as` to the Go, Node, Rust and Python fixtures. Go, Rust and Python return the host's `unsupported deliverAs "bogus"` error. Node's `pi.sendMessage` is declared void (`loader.js:305-308`), so its row requires the real host rejection through the absence of the Session's custom message after an ordered `active_tools_get` call. An SDK fallback would append `bogus-deliver-as` and fail. Subtests: subprocess-go, subprocess-node, subprocess-node-packed, subprocess-rust, subprocess-python, fused-go.
2. `TestNativeEventBusToolCallerCancellationEndsListener` registers an `emitwait` tool in the Go bus fixture. It runs the tool with a cancellable context while a native listener waits in `awaitcancel`, cancels the call, and requires the tool to fail and the listener to log `context canceled`. Both `isolated` and `shared-ok` pass.
3. `TestInteractiveSessionStartReportsResumeAndInitial/picker` seeds an existing Session, starts interactive mode, runs `/resume`, selects the first picker entry and asserts `session_start reason=resume initial=false`. SPEC-0007 D2 revision 3 makes `initial` the first event of a Runtime; a picker replacement omits it. The earlier `initial=true` target contradicted D2 and was corrected here.

Evidence: `go test ./coding/cli -run '^TestInteractiveSessionStartReportsResumeAndInitial$/^picker$'` and `go test ./test/extension-conformance -run '^(TestConformance_ExtensionAPIStrictDeliverAs|TestNativeEventBusToolCallerCancellationEndsListener)$'` pass under pinned Node 24.22.1.
