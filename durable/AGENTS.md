# AGENTS.md — durable

Port of Pi `@earendil-works/pi-durable`: durable conversation, task and document runtime, ID kinds, errors, entry kinds, harness, storage, sessions, tools. Go mapping decisions: `doc.go` (read before changing a type).

- Upstream: `.upstream/current/packages/durable/`. Depends on `chord/`, `env/`, `telemetry/`.
- `core/` — scheduling core compiled natively and to WASM (Go and TinyGo); own `core/Makefile`; design: `docs/plan/durable-core/` (ADR, ABI). Keep `core/` TinyGo-compatible; performance rules: `docs/performance/go-practices.md`.
- `contract/` — Pi golden traces; `make durable-contract-trace` compares call/commit interleaving with Pi.
- Branded IDs are distinct named int64 types; never convert between kinds to make code compile.
- Overloaded TS methods → erased tokens + generic helpers (`TxDoc`, `TxEntry`, `CreateTask`, ...); extend that pattern, don't add parallel APIs.
