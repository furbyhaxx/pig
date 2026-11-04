---
document:
  type: "specification"
  id: "SPEC-0003"
  title: "Model request pipeline: every model call hooked, fail-closed transforms, hook phases"
  description: "One request pipeline in the Model Runtime for agent turns, compaction, summaries and extension calls; security transforms fail closed; early/normal/late phases instead of priorities"
  status: "draft"
  revision: 1
  created_at: "2026-11-04T14:36:00+01:00"
  updated_at: "2026-11-04T14:36:00+01:00"
  created_by: "orchestrator"
  updated_by: "orchestrator"
  changelog:
    - revision: 1
      reason: "Resolve U-02, U-04 and U-15 of the host-gap verification"
      changes: ["document creation"]
---

# SPEC-0003 Model request pipeline

Resolves U-02, U-04, U-15 (non-prompt part) of `../analyses/2026-11-04-host-gap-verification/README.md`.

## Context

- Payload and header hooks are installed on the agent's stream options only. Compaction (`internal/codingagent/compaction/compaction.go:664`), branch summaries, cache warming (`coding/session_cache_warming.go`), virtual-model routes (`coding/virtual_models.go`) and extension calls through `ModelRegistry`/`ModelRuntime` (`coding/services.go:489-512`) reach providers without `before_provider_request`, `before_provider_headers`, `after_provider_response` or `provider_stream_event`.
- A `before_provider_request` handler error or panic is reported and the unchanged payload is sent (`coding/extension/host/inproc/runner.go:2120-2152`). A redactor that crashes leaks.
- Handlers run in extension load order; an extension cannot run after all others.

## Decisions

### D1 One pipeline

- `ModelRuntime` owns the request pipeline. `Stream`, `StreamSimple` and `CompleteSimple` in `coding/services.go` run the provider hooks for every request, whatever the caller.
- The session-level hook installation is removed; the agent path goes through the same pipeline.
- No handlers registered for an event → no work beyond one registry check (current zero-cost behavior kept, benchmarked).

### D2 Request purpose

Every provider event carries `purpose` and `origin`:

| `purpose` | Caller |
|---|---|
| `turn` | agent loop |
| `compaction` | compaction and its retries |
| `branch_summary` | tree navigation summaries |
| `cache_warm` | cache warming |
| `extension` | an extension's own model call (`origin` = extension name) |
| `internal` | other host calls (title generation, classifiers) |

Handlers filter on `purpose`. Existing handlers that ignore it see more requests than before; docs state this.

### D3 Re-entrancy

A model call made by an extension while one of its own provider-hook handlers runs is dispatched to every other extension's hooks, never back into the handler that issued it. Depth beyond 2 fails the request with an error naming the cycle.

### D4 Fail-closed transforms

| Event | On handler error, panic, timeout or crashed extension |
|---|---|
| `before_provider_request`, `before_provider_headers` | request withheld; the operation fails with an error naming extension and event |
| `context`, `context_with_system` | run fails before the request |
| `tool_call` (block hook) | tool call blocked with the error as reason |
| observers (`after_provider_response`, `provider_stream_event`, `message_*`, `turn_*`) | reported, run continues |

- Explicit block result: `before_provider_request` may return `{block: true, reason}`; the request fails with that reason.
- `ctx.Abort()` keeps its meaning.
- Handler timeout: setting `extensions.hookTimeoutMs` (default 30000) for the fail-closed events; 0 disables.

### D5 Phases

- Subscriptions take an optional phase: `early`, `normal` (default), `late`.
- Dispatch order: all `early` handlers in load order, then `normal`, then `late`.
- Phases apply to transform events: `context`, `context_with_system`, `before_provider_request`, `before_provider_headers`, `tool_call`, `tool_result`, `message_end`, `before_agent_start`.
- Intended use: redactors and policy checks register `late` so they see the final payload; normalizers register `early`. No numeric priorities.

### D6 Wire and SDKs

- Event payloads gain `purpose` and `origin`; subscribe call gains `phase`; handler result gains `block`/`reason` for `before_provider_request`.
- All SDKs (Go, Rust, Python, TS declarations) in the same change. Go SDK: typed `Purpose` and `Phase` constants.

## Consequences

- A secrets extension enforces redaction for every outbound request with one `late` `before_provider_request` handler.
- Compaction and summaries become visible to observers (usage accounting sees them).
- Extensions that relied on fail-open behavior now stop the request; the error tells the user which one.

## Docs

- `extensions.md` / `extension-api.md` / `events.md`: purposes, re-entrancy rule, fail-closed table, phases, block result.
- `security.md`: redaction guarantee and its scope (all model requests of the process).
- `compaction.md`: provider hooks fire for compaction with `purpose: compaction`.
- `settings.md`: `extensions.hookTimeoutMs`; `make generate-schemas`.
