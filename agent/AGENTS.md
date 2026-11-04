# AGENTS.md — agent

Port of Pi `packages/agent`: agent loop, agent state, message types, prompt queueing, context-token accounting.

- Upstream: `.upstream/current/packages/agent/src/`. Users: `coding/` (Session drives it), `ai/` (Provider calls), extension host (events).
- `transform.go` ports upstream transform-messages: normalizes history before provider conversion. Message-shape bugs are fixed here or in `ai/` normalization, not in a provider converter or renderer.
- Agent `Send`/`Continue` run only from `coding/session*.go` (divergence-guard `agent-run-outside-session`).
- Event order, queue clearing and abort semantics are observable; tests compare with Pi oracles (`*_pi_oracle_test.go`, `*_upstream_test.go`). Extend those before adding new test files.
