---
document:
  type: "specification"
  id: "SPEC-0004"
  title: "Extension runtime lifetime: bounded reloads, cancellable bus, live renderer resolution, retained processes"
  description: "Generation-unique sockets and bounded handshakes for reload, cancellation through event-bus listeners, renderer re-resolution without reload, and host-owned processes that survive extension reloads"
  status: "draft"
  revision: 5
  created_at: "2026-11-04T14:40:00+01:00"
  updated_at: "2026-11-04T23:10:42+01:00"
  created_by: "orchestrator"
  updated_by: "coding-agent"
  changelog:
    - revision: 1
      reason: "Resolve U-01, U-14, U-22 and U-23 of the host-gap verification"
      changes: ["document creation"]
    - revision: 2
      reason: "Implement generation isolation with deterministic teardown evidence; distinguish it from the unreproduced P29 hang"
      changes: ["D1 per-Host start counter covers isolated and packed starts", "D1 teardown regression and retained-process stress evidence", "record P29 reproduction gap"]
    - revision: 3
      reason: "Use the approved flat reload setting without changing the extensions path array"
      changes: ["D2 global extensionsReloadHandshakeMs", "clarify retained-process admission and per-cell rollback", "use consistent flat D5 setting names"]
    - revision: 4
      reason: "Complete Host and CLI deadline classification and cancellation ownership checks"
      changes: ["D2 cancellation is not runtime incompatibility", "D2 retained probes do not hold the registry mutex", "record generic SDK ownership as an unresolved implementation requirement"]
    - revision: 5
      reason: "Decide the SDK loader contract: per-cell publication over one persistent Host"
      changes: ["D2 SDK reloads reuse one Host and publish per cell", "D2 runtime facade: fresh snapshot per reload, caller-supplied Runtime keeps identity", "D2 factory failures retain the serving extension record"]
---

# SPEC-0004 Extension runtime lifetime

Resolves U-01, U-14, U-22, U-23 of `../analyses/2026-11-04-host-gap-verification/README.md`.

## D1 Generation-unique sockets (U-01)

- Every extension start gets its own socket path: `e-<index>-g<generation>.sock` from `sockPathFor` in `coding/extension/host/subprocess/host.go`. A stable per-name index and a monotonic per-Host start counter cover packed members, isolated extensions and park acknowledgements; cell-local counters do not cover all of those callers. Socket-directory selection reserves the full counter widths.
- The old generation's listener and socket file are removed only by the old generation's teardown; the new generation never removes a path it did not create.
- Socket directory sweep at host start removes files of dead hosts only (existing ownership check).
- `TestSocketGenerationTeardownDoesNotUnlinkSuccessor` fails before the fix: tearing down the previous packed member removes its successor's bound socket, and a connection to the successor fails. This proves the path-ownership defect, not the complete P29 hang.
- `TestPackedGoReloadLoopNeverHangs` exercises 200 reloads with factory-time host calls, bus subscription and tool execution; it also checks retained-process identity and shutdown cleanup. The original loop passed before the socket fix; this is stress coverage, not a red-proven hang reproduction. It does not drive tree navigation or an agent turn. The P29 sequence remains a reproduction obligation before claiming the original hang resolved.

## D2 Bounded reload handshake (U-01)

- Connect and register of a reloaded extension are bounded by the global setting `extensionsReloadHandshakeMs` (default 15000; positive integer milliseconds representable as a Go duration). The existing `extensions` path array is unchanged; project settings cannot override this budget. The timer starts at spawn or factory admission in a retained process; build time and waiting for a configured admission turn are excluded.
- On timeout or failure: a newly spawned process is killed and reaped; a failed admission in a shared retained process releases only the new generation, never the serving generations' process references. The old generation keeps serving, and the reload reports an error to the `ctx.Reload()` caller, `/reload` and RPC. Reload is atomic per planned cell: a failed member prevents its staged siblings from replacing their predecessors. Other successful cells may commit. User factory side effects cannot be rolled back.
- Initial loads keep their current cancellation-only behavior (startup shows progress; a slow first build is legitimate).
- Cancellation and reload deadline expiry are not runtime incompatibility: they do not trigger TypeScript runtime fallback, a Node-only record or SDK-rebuild advice. Retained-process control probes are cancellable without holding the shared retention registry mutex; their writer is owned and joined at process teardown.
- The SDK caller follows the same rule: `DefaultResourceLoader` reloads run `Host.Reload` on the Host it already owns instead of starting a new one, so extensions publish per cell and a failed factory or file extension keeps its serving generation. A trust bootstrap pass still starts a fresh Host.
- Runtime facade: a reload returns a fresh `ExtensionRuntime` view over the same live provider, MCP and virtual-model registration state, so the previous runner's factory APIs go stale while the replacement runner keeps the bound registries. A caller that supplied `DefaultResourceLoaderOptions.Runtime` keeps that runtime's identity instead. An extension a reload removed has its factory API retired. The SDK interface carries the same partial-success error as the CLI; a caller must handle it rather than assume a clean replacement.

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
- Session shutdown, session replacement and host exit terminate retained processes (graceful signal, then kill after `extensionsProcessKillGraceMs`, default 5000).
- Limit `extensionsMaxRetainedProcesses` (default 32) per extension; exceeding it fails `Spawn`.
- Processes run with the extension's isolation and spawn hooks; no new privilege.

## Wire and SDKs

New calls `processes.spawn|attach|list|signal|write`, notification `processes.output|exit`, `events.cancel`. All SDKs in the same change; conformance tests in `test/extension-conformance/`.

## Docs

- `extensions.md` / `extension-api.md`: reload atomicity and failure reporting, bus cancellation, renderer invalidation, process API.
- `runtime-cells.md` (bundle) / `docs/runtime-cell-reload.md`: generation sockets, handshake bound.
- `settings.md`: `extensionsReloadHandshakeMs`, `extensionsProcessKillGraceMs`, `extensionsMaxRetainedProcesses`; `make generate-schemas`.
