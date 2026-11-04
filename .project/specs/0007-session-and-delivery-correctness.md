---
document:
  type: "specification"
  id: "SPEC-0007"
  title: "Session and delivery correctness"
  description: "Queued nextTurn messages bound to their branch, resume reason for continued sessions, strict deliverAs, deduplicated session listing, lighter listing payloads, custom data in branch projections and an explicit no-op compaction result"
  status: "draft"
  revision: 1
  created_at: "2026-11-04T14:52:00+01:00"
  updated_at: "2026-11-04T14:52:00+01:00"
  created_by: "orchestrator"
  updated_by: "orchestrator"
  changelog:
    - revision: 1
      reason: "Resolve U-08, U-17, U-18, U-21 and the session minor items of the host-gap verification"
      changes: ["document creation"]
---

# SPEC-0007 Session and delivery correctness

Resolves U-08, U-17, U-18, U-21 and the `GetBranch`, `allMessagesText` and "Nothing to compact" items of `../analyses/2026-11-04-host-gap-verification/README.md`.

## D1 `nextTurn` bound to its branch (U-08)

- A queued `nextTurn` message records the session leaf id at queue time (`agent/agent.go` `pendingNextTurn`).
- Tree navigation, new session, fork, resume and session switch discard queued messages whose leaf is no longer the current leaf.
- Each discard emits `queued_message_discarded {customType, reason}` to extensions and a debug log line; the sender can re-queue.

## D2 Resume reason (U-17)

- Process start with an existing session (`--continue`, `--resume`, `--session <file>`) reports `session_start.reason = "resume"`; a fresh session at process start reports `startup`.
- New field `initial: true` on the first `session_start` of the process, for extensions that need "process just started".

## D3 Strict `deliverAs` (U-18)

- The host rejects an unknown `deliverAs` with an error on the `sendMessage`/`sendUserMessage` call.
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
