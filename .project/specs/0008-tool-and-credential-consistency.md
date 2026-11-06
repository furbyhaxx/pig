---
document:
  type: "specification"
  id: "SPEC-0008"
  title: "Tool and credential consistency: explicit tool cwd, one !command resolver"
  description: "Built-in tool factories honour an explicit cwd, and command-valued credentials resolve in one place for every provider registration kind"
  status: "published"
  revision: 2
  created_at: "2026-11-04T14:55:00+01:00"
  updated_at: "2026-11-06T11:37:12+01:00"
  created_by: "orchestrator"
  updated_by: "coding-agent"
  changelog:
    - revision: 1
      reason: "Resolve U-16 and U-19 of the host-gap verification"
      changes: ["document creation"]
    - revision: 2
      reason: "Owner accepted the decisions as the fork contract"
      changes: ["status draft -> published"]
---

# SPEC-0008 Tool and credential consistency

Resolves U-16, U-19 of `../analyses/2026-11-04-host-gap-verification/README.md`.

## D1 Explicit tool cwd (U-19)

- `Create<Tool>Tool(cwd, …)` with a non-empty `cwd` pins the tool to that directory; the session cwd applies only when `cwd` is empty.
- Host-registered built-in tools are created with an empty cwd and keep following the session.
- Wire: `protocol_builtin_tool.go` resolves `Cwd` before the session context; the SDK doc comments state the rule.

## D2 One `!command` resolver (U-16)

- Precondition: reproduce the pig-web observation (native provider `Resolve` receives `!printf …`; `GetApiKeyForProvider` returns it raw) in a hermetic test. If it does not reproduce, this decision is revised.
- Command-valued credentials (`!cmd`) from `auth.json`, `models.json`, `ProviderConfig` and native provider registrations resolve in one function (`internal/configvalue`) at use time, before a provider's `Resolve` runs and before `GetApiKeyForProvider` returns.
- Raw command strings never reach provider code, extension code or logs. Resolution failures name the provider and the source, not the command output.
- Existing caching semantics of command values stay.

## Docs

- `extensions.md` / SDK docs: tool factory cwd rule.
- `providers.md` / `custom-provider.md`: `!command` values resolve for every registration kind.
