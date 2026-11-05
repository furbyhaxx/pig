# Events and hooks

PiG-native extensions react to **events**. A separately selected extension can
translate another harness's hook format into those events. Stock Pig does not
load a hook bridge.

## Register PiG-native events in an extension

A Go factory extension registers event handlers on the SDK extension object:

```go
func Extension() *sdk.Extension {
    ext := sdk.New("demo")

    ext.OnSessionStart(func(ctx sdk.Context, data map[string]any) (any, error) {
        ctx.SetStatus("demo", "demo extension loaded")
        return nil, nil
    })

    ext.OnToolResult(func(ctx sdk.Context, data map[string]any) (any, error) {
        ctx.SetStatus("demo", "a tool finished")
        return nil, nil
    })

    return ext
}
```

Validate the extension after changing event code:

```bash
pig install ./my-ext --validate-only --json
```

## Event names are the extension API surface

The Go SDK exposes helpers for common events:

| SDK helper | Event name |
| --- | --- |
| `OnSessionStart` | `session_start` |
| `OnSessionShutdown` | `session_shutdown` |
| `OnToolResult` | `tool_result` |
| `OnEvent("session_info_changed", handler)` | `session_info_changed` |
| `OnEvent("session_before_compact", handler)` | `session_before_compact` |
| `OnEvent("session_compact", handler)` | `session_compact` |
| `OnEvent("session_compact_failed", handler)` | `session_compact_failed` |
| `OnEvent(name, handler)` | any named event |

`session_start` carries `reason` (`startup`, `resume`, `new`, `fork`, `reload`) and `initial: true` on the first Session of a run. A launch with `--continue`, `--resume`, or `--session` naming an existing file reports `resume`; a new Session at launch reports `startup`. Replacements and reloads omit `initial`. An SDK program gets `initial` on the first `session_start` of each `AgentSessionRuntime` (SPEC-0007 D2).

The lower-level event bridge also emits agent, turn, message, tool-execution,
input, model-selection, and session lifecycle events. The
`session_info_changed` event contains the effective Session name. Its `name`
field is absent when an extension clears the name. The
`session_before_compact` event contains the preparation, branch entries,
reason, retry state, custom instructions, and cancellation context. A handler
can cancel compaction or provide the compaction result. PiG persists a provided
result before it emits `session_compact`. The `session_compact` event contains
the persisted entry and identifies an extension-provided result. Use the SDK
helper when one exists. Use `OnEvent` only when you need a named event without
a helper.

## Branch-bound `nextTurn` messages

In this personal fork, `deliverAs: "nextTurn"` binds a custom message to the session and leaf where it is queued (SPEC-0007 D1). Tree navigation, in-place fork and session replacement remove stale messages before another prompt can consume them. Navigating to the same leaf keeps the queue, including selecting a user entry whose parent is that leaf. Normal transcript growth during a turn does not discard deferred context.

### `queued_message_discarded`

Each removed message emits `{"type":"queued_message_discarded","customType":"branch-note","reason":"tree"}` to extensions, in queue order. The host also writes a line to the debug log when `PIG_DEBUG` is set. Register with `pi.on("queued_message_discarded", handler)` or the Go SDK’s `OnEvent(sdk.EventQueuedMessageDiscarded, handler)`.

`reason` is `tree`, `new`, `fork`, `resume` (including switching to an existing session file), or `switch` (direct SDK `ReplaceInner`). A Runtime replacement discards the outgoing Agent’s queue before extension shutdown; it does not transfer queued context to the replacement. Cancelled operations retain the queue. Extensions can prepare fresh context for the destination instead of relying on a deferred message from the old branch.

## Hooks are external Resources

Some harnesses use declarative command triggers called hooks. A Package can
distribute hook metadata, and a Piglet can select an extension that translates
that metadata into PiG events. Stock Pig does not interpret cross-harness hook
files by itself.

Use SDK events when you write a PiG extension and need direct access to the SDK
context, host calls, and cancellation. Use a selected bridge extension when a
composition requires another harness's hook format.
