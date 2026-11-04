# AGENTS.md — chord

Port of Pi `@earendil-works/chord`: strict JSON copies, shared Context, delta operations (`delta/`), replicated state, watch, host, SQLite and telemetry adapters.

- Upstream: `.upstream/current/packages/chord/`. Users: `durable/`, `env/`, extension replicated state.
- `delta.JsonObject` keeps JavaScript property order and `delta/` Ops are tuples; don't replace them with Go maps where order is observable.
- JSON helpers shared with other packages: `internal/chord/chordjson`.
