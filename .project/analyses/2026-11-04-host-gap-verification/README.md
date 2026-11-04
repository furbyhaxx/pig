---
document:
  type: "analysis"
  id: "ANA-2026-11-04-host-gap-verification"
  title: "Verification of the Workbench host-gap report against the fork and pinned Pi"
  description: "Per-finding check of U-01..U-26 from the pig-workbench report: what the fork code does, what pinned Pi does, which findings are wrong, and which spec resolves each"
  status: "published"
  revision: 1
  created_at: "2026-11-04T14:23:15+01:00"
  updated_at: "2026-11-04T14:23:15+01:00"
  created_by: "orchestrator"
  updated_by: "orchestrator"
  changelog:
    - revision: 1
      reason: "Owner asked to check each finding of the Workbench report before designing the fixes"
      changes: ["Source verification of U-01..U-26 and the minor items", "Spec mapping"]
---

# 1. Scope

Input: `host:/mnt/Projects/furbyhaxx/pig-workbench/.project/analyses/2026-11-04-pig-upstream-candidates/README.md` (findings U-01..U-35, probes P1–P30, S1–S6, PR1–PR8).

Checked against:
- Fork `main` at `87e8a2f79` (`0.5.3`, Pi pin `1.4.0+6c384aa3`); `upstream/main` had no newer commits on 04.11.2026.
- Pi source at the pin (`badlogic/pi-mono@6c384aa3`, `packages/coding-agent/src`), fetched read-only. `make upstream-mirror` failed locally (`npm ci` EINTEGRITY), so `.upstream/current/` was not used.

Method: source reading plus the probe evidence in the Workbench report. Behavioral probes were not re-run. "Pi" below names upstream behavior only to explain where a behavior comes from; under SPEC-0001 D7 it is not a constraint.

# 2. Findings

| ID | Verified fact | Origin of behavior | Verdict | Spec |
|---|---|---|---|---|
| U-01 | Reload blocks in `awaitConnection` (`coding/extension/host/subprocess/accept_client.go:64`), no deadline. Socket leaf per extension name is stable across generations (`host.go:2596` `sockPathFor` → `e-N.sock`); the new generation removes and re-listens on the same path (`packed_go.go` `startGoPackedCellOnce`). Probe: new instance's dial refused on `e-0.sock` (`reload-hunt/probe-9.jsonl`). Root cause UNVERIFIED: path reuse across generations is the leading hypothesis | PiG | bug | 0004 |
| U-02 | `before_provider_request` handler error is recorded and the loop continues with the unchanged payload (`coding/extension/host/inproc/runner.go:2120-2152`) | Pi `runner.ts:1383` | change | 0003 |
| U-03 | Node extension console output reaches the terminal (stderr outside TUI, `host.go:5087` `nodeExtensionOutput`); Go/Rust/Python processes write stderr to a log only (`host.go:3212,3244`, `packed_go.go:643`). Print mode `Notify` is dropped (no-op UI context) | PiG drift + Pi no-op UI | gap | 0005 |
| U-04 | Compaction uses its own completer (`internal/codingagent/compaction/compaction.go:664`, `provider_completer.go`); payload hooks are installed on the agent only | Pi `compaction.ts:636`, `sdk.ts:376-424` | gap | 0003 |
| U-05 | Built-in extensions load after configured ones; MCP sets `mcp_servers` in its own `before_agent_start` | Pi `package-manager.ts:986-1002` | gap | 0002 |
| U-06 | Forced prompt is projected onto the request and not recorded; sections set after a force are recorded but not sent | Pi `agent-session.ts:1748-1820` (stated design) | change | 0002 |
| U-07 | Custom sections can replace a non-preamble built-in section in place (`internal/codingagent/prompts/custom_sections.go:12`); `preamble` is rejected, empty values are skipped, new names append | Pi `system-prompt.ts:121-186` | gap | 0002 |
| U-08 | `nextTurn` messages live in `agent.pendingNextTurn` (`agent/agent.go:894-899`) and are not cleared on `NavigateTree` | Pi `agent-session.ts:399,1960,2371` | bug | 0007 |
| U-09 | **Wrong.** The block reason is recorded (`nestedCalls[].error = "WB-BARRIER: …"`, `runs/codemode/rpc.jsonl:25`). `{"reason":{}}` is `JSON.stringify` of a JS `Error` in the probe script | — | none | — |
| U-10 | Interactive mode skips nested `tool_execution_*` events; codemode draws its children | Pi `interactive-mode.ts:3565` | gap | 0006 |
| U-11 | `viewEvent` routes `ui.view.event` to overlays only (`extensions/sdk/view.go:311`) | PiG kit | gap | 0006 |
| U-12 | No secret dialog; a masked input exists for `/login` (`tui/login_dialog.go:47-171`) | Pi has none | gap | 0005 |
| U-13 | `Select`/`Input` return `(any, ok, error)`; cancel, timeout and no-UI all give `ok=false` (`extensions/sdk/context.go:478-570`) | Pi undefined result | change | 0005 |
| U-14 | No process retention across reload | — | gap | 0004 |
| U-15 | Handlers run in load order; no phases | Pi | gap | 0002 (prompt), 0003 (phases) |
| U-16 | Stored `!command` keys are deliberately not resolved on read (Pi `auth-storage.ts:264`; fork `ai/auth_store.go:602`); resolution happens in the provider composer (Pi `provider-composer.ts:472`). Native provider `Resolve` therefore sees the raw value. pig-web observation not reproduced here | Pi | inconsistency | 0008 |
| U-17 | Initial session reports `startup` also for `--continue` (`coding/session.go:901`) | Pi `agent-session.ts:509`, `main.ts:740` | change | 0007 |
| U-18 | Unknown `deliverAs` takes the steer/append path; only TS types restrict it | Pi | bug | 0007 |
| U-19 | Built-in tools prefer the session cwd over the factory argument (`protocol_builtin_tool.go:120`) | Pi `ctx?.cwd \|\| cwd` (`read.ts:124` and siblings) | change | 0008 |
| U-20 | One observation, not reproduced | — | watch | — |
| U-21 | `listAll` includes symlinked bucket dirs without dedupe | Pi `session-manager.ts` `listAll` | bug | 0007 |
| U-22 | Renderer resolution cached per tool until `/reload` | PiG | change | 0004 |
| U-23 | Caller abort cancels its `events.emit`; listener runs on (`event_bus.go:481`) | PiG wire | bug | 0004 |
| U-24 | `~/.agents/skills` derives from `HOME` (`internal/codingagent/resource_paths.go:41,156`); `skills` setting is a path list (`internal/codingagent/settings.go:357`) | Pi `package-manager.ts:2489` | feature | 0009 |
| U-25 | Go cell identity hashes only the Go source set; dot/underscore dirs, `testdata`, `vendor`, nested modules are excluded (`coding/extension/host/runtimecell/hash_go_module.go:30-200`), so `.project/` is **not** packed. "Inputs changed while built" is a deliberate guard (`cell_inputs.go:15`). No build-tag option exists | PiG | partly wrong; gap | 0010 |
| U-26 | **Wrong.** `github.com/MichaelKinsy/PiG/extensions/sdk@v0.5.3` resolves from the Go proxy; the `replace` only pins the host's staged SDK | — | none | — |

Minor items (Workbench report §5–6):

| Item | Verdict | Spec |
|---|---|---|
| `GetBranch` omits custom entry data | gap | 0007 |
| Context-file suppression needs per-run option edits | resolved by section control | 0002 |
| `SessionManagerList/ListAll` always include `allMessagesText` | change | 0007 |
| "Nothing to compact" only via `OnError`/`compaction_end` | change | 0007 |
| `--append-system-prompt` replaces `APPEND_SYSTEM.md` | exposed to templates | 0002 |
| `--tools` hides MCP/codemode, `-ne` disables built-ins | document only | — |

# 3. Outside this repo

- pi-accordion spawns `spawn(process.execPath, [runnerPath])` (`host:$HOME/.pi/agent/git/github.com/furbyhaxx/pi-accordion/extension/accordion.ts:761`). In an ahead-of-time compiled harness `process.execPath` is the harness, so each subagent start opened a billed agent session. Fix belongs in pi-accordion: resolve a Node executable explicitly.
