---
document:
  type: "specification"
  id: "SPEC-0007"
  title: "Session and delivery correctness"
  description: "Queued nextTurn messages bound to their branch, resume reason for continued sessions, strict deliverAs, deduplicated session listing, lighter listing payloads, custom data in branch projections and an explicit no-op compaction result"
  status: "draft"
  revision: 4
  created_at: "2026-11-04T14:52:00+01:00"
  updated_at: "2026-11-05T12:32:01+01:00"
  created_by: "orchestrator"
  updated_by: "coding-agent"
  changelog:
    - revision: 1
      reason: "Resolve U-08, U-17, U-18, U-21 and the session minor items of the host-gap verification"
      changes: ["document creation"]
    - revision: 2
      reason: "Owner decided the initial flag is Runtime-scoped; implement D2 and D3"
      changes: ["D2 initial marks the first session_start of a Runtime, not of the OS process", "D2 record where resume is decided and how every mode receives it", "D3 record the shared validation boundary"]
    - revision: 3
      reason: "Implement D1 branch-bound nextTurn delivery"
      changes: ["Record queue-time session and leaf identity, lifecycle discard boundaries, typed discard event and regression evidence"]
    - revision: 4
      reason: "Cover the interactive in-place fork path"
      changes: ["Record the direct slash fork boundary and its red-proven regression; share discard dispatch between modes"]
---

# SPEC-0007 Session and delivery correctness

Resolves U-08, U-17, U-18, U-21 and the `GetBranch`, `allMessagesText` and "Nothing to compact" items of `../analyses/2026-11-04-host-gap-verification/README.md`.

## D1 `nextTurn` bound to its branch (U-08)

- A queued `nextTurn` message records the session identity and leaf id at queue time in `agent/agent.go` `pendingNextTurnMessage`. `Agent.BindNextTurnSession` binds the leaf reader to the SessionManager accessor `(*internal/codingagent.Session).GetLeafID()` (`coding.Session.SessionManager()` returns that underlying Session). A nil leaf identifies the position before the first entry; forked sessions sharing entry ids remain distinct by session identity.
- Tree navigation and in-place fork discard messages whose stored identity differs from the destination. Selecting a user message can resolve to its parent leaf: if that is the queued leaf, the message stays queued. Ordinary transcript appends do not invalidate the queue; messages queued during a running turn still join the next prompt.
- In-place new/resume/switch replacements compare the destination session and leaf before rebinding the reader. A Runtime replacement retires the old Agent and discards its entire queue before extension shutdown; it does not transfer deferred messages between Agent instances. Cancelled navigation or replacement does not discard.
- Each discard emits the typed `extension.QueuedMessageDiscardedEvent` as `{"type":"queued_message_discarded","customType":"…","reason":"…"}` through the normal extension dispatch and writes a `PIG_DEBUG` line. Reasons are `tree`, `new`, `fork`, `resume` (including switch to an existing session file), and `switch` for direct `ReplaceInner`. Events retain queue order; removing a message before dispatch prevents a second discard of it. Extensions can choose to queue fresh context for the destination instead of retaining stale branch state.
- Production boundaries: `Session.navigateTree`, `Session.Fork`, `Session.replaceInner`, `Runtime.teardownCurrent`, and interactive `buildSlashContext.ForkAtEntry` (which calls `forkAndRebuild` directly). All three prompt injection paths drain the same Agent queue, including interactive delivery.
- Evidence: `TestNextTurnLeafBinding` failed before the fix with injected stale context and zero discard events, for ordinary, errored and aborted prior assistant turns; it now checks exactly one discard and same-leaf retention after navigation to a child user entry. `TestNextTurnLeafStampAndDiscard` checks queue-time identity, ordered filtering, nil root and session identity. `TestInteractiveForkDiscardsNextTurn` independently failed with a retained queue and no event before the slash fork boundary was covered; it now verifies the interactive sendMessage path and the production fork callback. `TestNextTurnRuntimeReplacement` covers new/fork/resume; `TestNextTurnSessionReplacement` covers direct SDK boundaries and same-session/same-leaf resume retention. `TestConformance_EventPayloadsReachEverySDKUnchanged` includes the new event for Go, Node, Rust and Python in isolated, packed and fused realizations.

## D2 Resume reason (U-17)

- Process start with an existing session (`--continue`, `--resume`, `--session <file>` naming an existing file) reports `session_start.reason = "resume"`; a fresh session at process start reports `startup`. `--fork` at launch creates a new log and stays `startup`.
- New field `initial: true` on the first `session_start` of a `Runtime`; replacements and reloads omit it. In the CLI one process owns one Runtime, so this is "process just started". An SDK program that creates several Runtimes sees `initial` once per Runtime. Owner decision 04.11.2026: a package-level once-flag would match "process" literally but is a hidden global that makes event tests order-dependent; a CLI-only marker would give embedders a different contract for the same event.
- The CLI decides `resume` when it selects the startup Session, before anything can create the file, and passes the event to `CreateAgentSessionRuntime`; `NewAgentSessionRuntime` marks it `initial`. Print, JSON and RPC forward the whole event (`Session.EmitSessionStartEvent`); interactive mode receives it through `InteractiveOptions.SessionStartEvent`.
- Evidence: `TestRuntimeMarksOnlyItsFirstSessionStartInitial`, `TestSessionStartReportsResumeAndInitialAtProcessStart` (print and RPC), `TestInteractiveSessionStartReportsResumeAndInitial`, each red before the change; the payload row in `TestConformance_EventPayloadsReachEverySDKUnchanged` carries `initial` to Go, Node, Rust and Python, isolated, packed and fused.

## D3 Strict `deliverAs` (U-18)

- The host rejects an unknown `deliverAs` with an error on the `sendMessage`/`sendUserMessage` call. One validation at the Session actions (`extension.ValidateDeliverAs`) serves every binding; `sendUserMessage` also rejects `nextTurn`.
- SDKs expose typed constants (`DeliverSteer`, `DeliverFollowUp`, `DeliverNextTurn`).

## D4 Session listing (U-21, `allMessagesText`)

- `ListAll` and the all-projects `/resume` view resolve bucket directories and session files with `EvalSymlinks` and list each real file once (newest mtime wins on equal ids).
- `SessionManagerList/ListAll` take `ListOptions{IncludeText bool}`; default false. Search features that need text request it.

## D5 Branch projection with custom data

`GetBranch` returns custom entries with their `data`, as `GetEntries` does, so extensions restore state without walking parent chains.

## D6 Explicit no-op compaction

`Compact()` (command, SDK, RPC) returns a result `{compacted: false, reason: "nothing_to_compact"}` instead of reporting only through `OnError`/`compaction_end`.

## Docs

- `extensions.md` / `events.md`: `queued_message_discarded`, `session_start.reason` and `initial`, `deliverAs` errors, `GetBranch` data, `Compact` result.
- `sessions.md`: symlinked buckets listed once.
- `rpc.md` / `rpc-commands.md`: compaction result, `deliverAs` error.
- SDK docs: `ListOptions`, constants; changelog `### Upgrading` for `ListAll` text default and `startup`→`resume`.
