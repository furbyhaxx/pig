# AGENTS.md — telemetry

Port of Pi `packages/telemetry`: backend-neutral spans, propagation, links, metrics, in-memory recording, rendering.

- Upstream: `.upstream/current/packages/telemetry/`. `telemetrytest/` — conformance fixtures shared with adapters (e.g. `chord/telemetrysqlite`).
- New adapter capability → conformance case in the shared catalogue (`conformance_*_test.go`), not an adapter-local test only.
