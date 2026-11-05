### Fixed

- A launch with `--continue`, `--resume`, or `--session` naming an existing file reports `session_start` reason `resume` instead of `startup`, in print, JSON, RPC and interactive mode (SPEC-0007 D2). <!-- fix: 453651743 -->

### Added

- `session_start` carries `initial: true` on the first Session of a run, so an extension can run once-per-launch work without tracking state itself. Replacements and reloads omit it; an SDK program gets it on the first `session_start` of each `AgentSessionRuntime` (SPEC-0007 D2).

### Upgrading

- An extension that treated `reason: "startup"` as "the process just started" now also sees `reason: "resume"` for a continued session at launch. Check `initial` instead; use `reason` only to tell a continued conversation from a fresh one.
