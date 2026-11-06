---
document:
  type: "specification"
  id: "SPEC-0005"
  title: "Headless output and dialogs: extension output channel, typed dialog results, secret input"
  description: "Same stdio routing for every extension runtime, an explicit output API for print/JSON/RPC, dialog results with a status, and a masked secret input in every mode that has a UI"
  status: "published"
  revision: 2
  created_at: "2026-11-04T14:44:00+01:00"
  updated_at: "2026-11-06T11:37:12+01:00"
  created_by: "orchestrator"
  updated_by: "coding-agent"
  changelog:
    - revision: 1
      reason: "Resolve U-03, U-12 and U-13 of the host-gap verification"
      changes: ["document creation"]
    - revision: 2
      reason: "Owner accepted the decisions as the fork contract"
      changes: ["status draft -> published"]
---

# SPEC-0005 Headless output and dialogs

Resolves U-03, U-12, U-13 of `../analyses/2026-11-04-host-gap-verification/README.md`.

## D1 One stdio routing for all runtimes

- Go, Rust and Python extension processes get the routing Node processes have (`coding/extension/host/subprocess/host.go:5087` `nodeExtensionOutput`): stdout to the terminal in TUI mode, stdout and stderr to the host's stderr in print/JSON/RPC mode, stderr also to the extension log.
- Applies to isolated and packed processes (`host.go:3212,3244`, `packed_go.go:643`).

## D2 Output API

`ctx.Output(text, OutputOptions{Stream})`, `Stream`: `result` (default) or `diagnostic`.

| Mode | `result` | `diagnostic` |
|---|---|---|
| print | host stdout, serialized with the host's own output, newline-terminated | host stderr |
| json, rpc | event `{"type":"extension_output","extension":…,"stream":…,"text":…}` | same event |
| tui | info notification | info notification |

- Works from commands, flags and `session_start` before a prompt, so `pig --list-tools` prints and exits.
- `Notify` in print mode writes `[<extension>] <level>: <message>` to stderr instead of being dropped; JSON/RPC emit the existing notify event.

## D3 Dialog results with a status

- Every dialog (`Select`, `Confirm`, `Input`, `Editor`, custom) returns a status: `answered`, `cancelled`, `timeout`, `unavailable` (no UI in this mode).
- Wire: the host's dialog reply carries `status`. RPC: `extension_ui_response` with `cancelled: true` → `cancelled`; host-side timer → `timeout`.
- Go SDK, typed: `Select(…) (string, DialogStatus, error)`, `Confirm(…) (bool, DialogStatus, error)`, `Input(…) (string, DialogStatus, error)`, `Editor(…) (string, DialogStatus, error)`. The `any`-typed variants are removed. A non-string RPC answer to a string dialog is an error, not a value. Rust/Python/TS declarations follow the same shape.
- Setting `extensions.dialogTimeoutMs` (default 0 = none) applies when the extension passes no timeout; RPC no longer blocks forever when it is set.

## D4 Secret input

- `Input` takes `InputOptions{Secret: true}`.
- TUI: the masked input of `/login` (`tui/login_dialog.go`) with a reveal toggle; the value never enters history, logs, the transcript or debug dumps.
- RPC: `extension_ui_request` carries `secret: true`; client docs require masking and no logging.
- Print/JSON: `unavailable`.

## Docs

- `extensions.md` / `extension-api.md`: output API, notify in print mode, dialog statuses, secret input, removed `any` variants (upgrade note).
- `rpc-extension-ui.md`: `status`, `secret`, `extension_output` event; `json.md`: `extension_output` event.
- `settings.md`: `extensions.dialogTimeoutMs`; `make generate-schemas`.
- Changelog `### Upgrading`: Go SDK dialog signatures.
