# Runtime-cell reload and operations

This page covers atomic replacement, crash quarantine/fission, diagnostics, and the protocol boundary. Start with the [runtime-cell overview](extension-runtime-cells.md).

## Atomic reload

Each extension start receives a generation-unique socket path (`e-<index>-g<generation>.sock`). Retiring an old generation removes only its own socket, not its successor's. The host reserves space for the counters when selecting a short socket directory (SPEC-0004 D1).

`Host.Reload` treats a new extension set as a transaction:

1. resolve the requested extension configs/specs;
2. plan isolated and packed cells;
3. reuse or build cell artifacts;
4. start replacement cells beside current cells, for unchanged extensions too (upstream reload invokes every extension factory again);
5. complete per-extension protocol registration;
6. validate identities/contributions and construct the replacement registry;
7. atomically publish the new registry;
8. stop old cells only after the swap.

Every configured factory runs again, including unchanged factories in retained processes. `Host.Reload` bounds connection and registration by the global `extensionsReloadHandshakeMs` setting (default 15000 milliseconds). The clock starts at spawn or factory admission, not during compilation. Retained-process admission writes are cancellable too. Initial loading remains cancellation-only (SPEC-0004 D2).

A failed cell keeps its previous serving generations. Its staged siblings and declarations are discarded rather than retried in isolation. Newly spawned failed processes are killed and reaped; failure in a retained process does not kill its serving generations. Provider, OAuth, MCP and virtual-model declarations publish only for accepted cells. User factory side effects cannot be rolled back.

Failures populate `ReloadReport.Issues` and return an error through the CLI's `ctx.Reload()`, `/reload` and RPC paths. Healthy cells may still commit. Go callers can identify that result with `errors.Is(err, extension.ErrReloadPartiallyApplied)` and must use the returned active extension set even when an error is returned. CLI callers publish that set before reporting the error. A completely failed replacement leaves the previous runner active; a config error changes no extension generations.

## Quarantine and fission

Native packed cells retain their quarantine/fission fallback. Node cells preserve cooperating healthy extensions:

```text
attributable failure:  A + culprit + C -> (A + C) | culprit
first unknown failure: A + B + C + D -> A + B + C + D
repeated unknown:      A + B + C + D -> (A + B) | (C + D)
culprit identified:    (A + B) | C | D -> (A + B + D) | C
```

The host attributes a crash to the factory being admitted, an extension path in the bounded stderr stack tail, or the last dispatched owner. With no attribution, it restarts the whole group once. Repeated unknown failures split diagnostic groups in halves until the culprit is isolated. Healthy members then rejoin one shared Node bus. Only the culprit's recurrent singleton failure can trip its restart circuit breaker.

Recovery reruns factories, not interrupted tools or callbacks. Current invocations fail on their original dead connections. Subsequent named capability calls can use the recovered instance. Generation checks reject stale crash reports; old connection-owned UI state cannot mutate the replacement. Original owner cancellation and Host shutdown cancel and drain recovery. Every restart, diagnostic split and quarantine is reported through the existing crash-notice callback. Quarantine is host runtime state and never changes authored isolation.

## Diagnostic logs

PiG captures one temporary `pig-packed-*.log` stderr file per packed process (D20) or `pig-ext-*.log` per isolated process (D56). Each file keeps the first 1 MiB and the last 64 KiB of the output and states how many bytes it dropped, so a startup failure and a crash line both survive a flood. Shutdown and reload close the process handles and remove that file. A load or crash diagnostic that names the file retains it for troubleshooting, including after reload. A cancelled load does not retain a log. Remove a retained file manually when you no longer need it.

When a shared Go, Rust, or Python process crashes, the error for each stopped extension says that `/reload` restarts it. A Go, Rust, or Python extension that fails while it registers shows its own error in the load message, and a Rust panic shows its panic message.

A socket close and process exit claim the same failed member only once. Interrupted commands, shortcuts and events keep their error result, but the lifecycle handler owns the connection failure diagnostic when a crash handler is installed. Ordinary handler errors still reach the runner's error listener. Interactive commands use that same runner path.

PiG does not sweep logs from earlier sessions. Their names do not identify an owner or prove that its process has exited. Automatic deletion could remove a live session's log or a retained failure diagnostic.

## Reload report

`Host.LastReloadReport()` returns the last placement transaction. `/reload --explain` renders the same facts for a user.

| Fact | Examples |
|---|---|
| strategy/key/members | packed or isolated placement |
| reason | source command, shared factory, quarantine/fission |
| artifact | resolved runner path |
| cache/build | hit or cold-build duration |
| result | success/failure and wall duration |

## Protocol boundary

Every extension, including each member of a packed cell, uses the same current wire contract over its own connection. The wire has no independent version. Required behavior includes:

- registration with stable extension identity;
- request/response correlation;
- cancellation and shutdown;
- host calls and results;
- tool updates and UI notifications;
- frame-size enforcement;
- startup/register timeouts;
- error propagation without crashing the host.

Packed mode must be observationally identical to isolated mode. Fix the host or number a divergence when a conformance difference appears; do not change SDK semantics to hide it.

## Oversized results and panics

The host catches extension-call panics at the extension boundary and returns a structured internal error when a response is possible. Results above the protocol frame limit are replaced with a small `result_too_large` error so the host does not write a frame no compliant extension can read.

## Piglet component runtime

A Piglet build/release may populate neutral runtime inputs:

| Input | Runtime job |
|---|---|
| effective Piglet | provide immutable agent composition |
| component closure | expose exact prebuilt subprocess runners/assets to ordinary resolution |
| fused registrations | register build-time-linked compatible Go factories |

Stock Pig has no Piglet-specific registrations or closure. Runtime registration remains generic and does not own Piglet schema or product transport.

## Operational boundaries

| Pig owns | Product/platform owns |
|---|---|
| placement/reload facts | org policy and quotas |
| cell build/cache/start | build-environment scheduling |
| registration/cancellation/errors | vulnerability/signature policy |
| quarantine/fission | publication approval |
| local reports | external telemetry collection/storage |

## Verification

- reload success swaps all registries together;
- a failed build/start/register preserves its cell's serving generations and returns the failure;
- Node factory admission preserves interleaved native order without splitting the Node bus;
- a Node crash restarts healthy members together, with one whole-group retry and then bisection when attribution is unavailable;
- native quarantine retains its explicit fission policy;
- report getters return copies and are race-safe;
- oversized results return a small structured error;
- SDK conformance compares isolated, packed, and fused paths where applicable;
- stress tests cover reload/provider lifecycle and cancellation.

Primary sources/tests:

- `coding/extension/host/subprocess/host.go`
- `reload_cells.go`, `packed_quarantine.go`, `reload_report.go`
- `coding/extension/host/runtimecell/`
- `test/extension-conformance/`
