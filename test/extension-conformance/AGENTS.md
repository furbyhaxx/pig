# AGENTS.md — test/extension-conformance

Authority for extension SDK behavior: one scenario table run against every SDK (Go, Rust, Python, Node) through the real host, isolated and packed. The in-process Go reference defines expected behavior.

- Run: `make test-sdk-conformance`; aot host: `make test-aot-conformance`.
- `testdata/*-go/` etc. — per-SDK fixture extensions; `sdkscenarios/` — the shared table.
- A row binds a value distinguishable from every SDK's fallback, so an SDK that never makes the call fails.
- New SDK capability → one row that fails when any single SDK is broken; add fixtures for all SDKs in the same change.
- Missing toolchain is an environment failure, never a skip.
