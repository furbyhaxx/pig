# AGENTS.md — coding/extension

Pi extension API port: registration, events, tools, commands, UI surfaces, providers, MCP, and the host that runs extensions as subprocesses.

## Layout
- Root: API types and dispatch (`api.go`, `events.go`, `eventbus.go`, `tool*.go`, `provider*.go`, `ui*.go`). Upstream: `.upstream/current/packages/coding-agent/src/core/extensions/types.ts`, `runner.ts`.
- `host/subprocess/protocol.go` — the one wire contract all SDKs target. No independent wire version, no negotiation.
- `host/subprocess/cell_plan.go`, `host/runtimecell/` — packed runtime cells (D20); `host/fusepack/` — fused members (`registry_generated.go` generated).
- `host/subprocess/runtime-node/` — Node runtime; `shims/*` are vendored Pi package trees, keep bytes.
- `builtin/`, `factoryload/`, `installresolver/`, `pigsdk/`, `extensiontest/`.
- Docs: `docs/extension-authoring.md` (read before any host feature), `docs/extension-api-parity.md`, `docs/extension-runtime-cells.md`, `docs/additive-features.md` D19/D20.

## Contract
- Every SDK exposes one extension contract. Isolated, packed and fused realizations of one factory pass the same conformance behavior; packed mode is observationally identical to isolated from the extension's view.
- Registration owns identity and capabilities. Reload and shutdown remove owned handlers and work.
- Each contained extension keeps its own socket and `register` handshake.
- Two source forms: conventional factory (one standard symbol per language) and exact standalone (always isolated). Pig derives wire, SDK, runner, placement, registration.

## Adding a host feature
1. Upstream Pi behavior.
2. Protocol shape in `host/subprocess/protocol.go`.
3. Host wiring in `host/subprocess/`.
4. SDK ergonomics in every `extensions/sdk*/` (`extensions/AGENTS.md`).
5. Conformance row in `test/extension-conformance/`.
6. Parity row in `docs/extension-api-parity.md`, incl. its async contract; exercise isolated and packed modes.
7. Numbered divergence only if needed.

## Never without an approved spec
- WASM or embedded-JS extension runtime; `plugin.Open` loader; linked/in-process production runtime — except the compiled-in `extension.ExtensionFactory` passed by value from the binary's own `main` (`docs/specs/extension-factory-trust.md`). Settings and `-e builtin:<name>` may select/disable a registered built-in; nothing user-, project-, Package-, Piglet- or wire-supplied adds factory code. Piglet members (incl. fused) never use it.
- Multi-register (`RegisterPayload.Extensions`, `RequestPayload.TargetExtension`).
- SDK semantics diverging from the in-process Go reference without a numbered divergence and a failing conformance test.
- An authored extension manifest or another global extension configuration plane.
- Changing the SDK API to hide a packed/isolated difference; fix the host or number a divergence.

## Failure modes
| Failure | Detection | Required result |
|---|---|---|
| An SDK returns a fallback instead of host state | value distinct from every SDK default | match the in-process reference |
| One realization has a weaker API | same row isolated, packed, fused | same interface and behavior |
| One extension crash stops another topology | crash isolated, call fused and packed tools | healthy extensions stay available |
| Reload duplicates a handler | subscribe, reload, dispatch | each live handler once |
| Unsubscribe mutates active dispatch | remove during dispatch | only the next snapshot changes |
| Large extension set stalls startup | representative mixed set | bounded startup/registration |
| Frames exhaust memory | oversized, malformed, partial, stalled frames | bounds enforced, accurate cancel |
| Shutdown leaks a process or socket | cancel during host calls and streaming | resources drained and closed |

## Evidence
- Define the cross-language row before implementation; exercise Go, Rust, Python, Node through real bridges and isolated/packed/fused/mixed realizations. External components through their own protocol, not the SDK.
- Measure startup, reload, dispatch, memory, fds, processes, shutdown. Missing toolchain = environment failure, not a skip.
- No blocking extension IPC on the TUI input/render loop.
