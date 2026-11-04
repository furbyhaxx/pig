---
document:
  type: "specification"
  id: "SPEC-0004"
  title: "Extension runtime lifetime: bounded reloads, cancellable bus, live renderer resolution, retained processes"
  description: "Generation-unique sockets and bounded handshakes for reload, cancellation through event-bus listeners, renderer re-resolution without reload, and host-owned processes that survive extension reloads"
  status: "draft"
  revision: 1
  created_at: "2026-11-04T14:40:00+01:00"
  updated_at: "2026-11-04T14:40:00+01:00"
  created_by: "orchestrator"
  updated_by: "orchestrator"
  changelog:
    - revision: 1
      reason: "Resolve U-01, U-14, U-22 and U-23 of the host-gap verification"
      changes: ["document creation"]
---

# SPEC-0004 Extension runtime lifetime

Resolves U-01, U-14, U-22, U-23 of `../analyses/2026-11-04-host-gap-verification/README.md`.

## D1 Generation-unique sockets (U-01)

- Every extension start gets its own socket path: `e-<index>-g<generation>.sock` from `sockPathFor` (`coding/extension/host/subprocess/host.go:2596`), with the generation counter the packed cells already keep.
- The old generation's listener and socket file are removed only by the old generation's teardown; the new generation never removes a path it did not create.
- Socket directory sweep at host start removes files of dead hosts only (existing ownership check).
- Before the fix: a looping red test that reproduces the hang (`TestPackedGoReloadLoopNeverHangs`, 200 reloads, reload after tree navigation and an agent run as in probe P29). The root-cause hypothesis is confirmed or replaced by that test, and this decision is revised if the cause differs.

## D2 Bounded reload handshake (U-01)

- Connect and register of a reloaded extension are bounded by `extensions.reloadHandshakeMs` (default 15000). The timer starts when the process is spawned; build time is excluded.
- On timeout or failure: the new process is killed, the old generation keeps serving, and the reload reports an error to the `ctx.Reload()` caller, `/reload` and RPC. Reload is atomic per cell.
- Initial loads keep their current cancellation-only behavior (startup shows progress; a slow first build is legitimate).

## D3 Cancellable event bus (U-23)

- `events.emit` carries the caller's request id. When the caller's context ends, the host sends `events.cancel` to every listener process still running that emit.
- SDKs cancel the listener handler's context. Late replies stay discarded.
- Same for in-process listeners: the handler context derives from the caller's.

## D4 Live renderer resolution (U-22)

- The renderer resolution cache is invalidated when a tool renderer resolver registers or unregisters, or the extension set changes.
- Live rows and rows in the visible viewport re-resolve immediately; older history re-resolves on its next render.
- Historical rows render from the persisted tool call arguments and details; the U-20 empty-args observation is covered by a regression test once reproduced.

## D5 Retained processes (U-14)

Host-owned processes that outlive an extension reload:

| Call | Effect |
|---|---|
| `ctx.Processes().Spawn(ProcessSpec{Key, Command, Args, Cwd, Env, Retain})` | start a process owned by the host; `Retain`: `reload` (survives extension reloads, ends with the session) or `none` (ends with the extension generation) |
| `ctx.Processes().Attach(key)` | handle of a running process of this extension: stdout/stderr stream from a bounded ring buffer (default 1 MiB) plus live output, exit status, `Signal`, `Kill`, `Write` stdin |
| `ctx.Processes().List()` | this extension's processes with key, pid, state, start time |

- Keys are scoped per extension name; a new generation of the same extension attaches to its predecessor's retained processes.
- Session shutdown, session replacement and host exit terminate retained processes (graceful signal, then kill after `extensions.processKillGraceMs`, default 5000).
- Limit `extensions.maxRetainedProcesses` (default 32) per extension; exceeding it fails `Spawn`.
- Processes run with the extension's isolation and spawn hooks; no new privilege.

## Wire and SDKs

New calls `processes.spawn|attach|list|signal|write`, notification `processes.output|exit`, `events.cancel`. All SDKs in the same change; conformance tests in `test/extension-conformance/`.

## Docs

- `extensions.md` / `extension-api.md`: reload atomicity and failure reporting, bus cancellation, renderer invalidation, process API.
- `runtime-cells.md` (bundle) / `docs/runtime-cell-reload.md`: generation sockets, handshake bound.
- `settings.md`: `extensions.reloadHandshakeMs`, `extensions.processKillGraceMs`, `extensions.maxRetainedProcesses`; `make generate-schemas`.
