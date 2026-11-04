# AGENTS.md — env

Port of Pi `packages/env`: remote execution environment for Durable (tools run on another machine over SSH; worker, storage, credentials stay local). Mapping notes: package comment in `connection.go`.

- Upstream: `.upstream/current/packages/env/`. `daemon/` — remote daemon; `cmd/` — its binary; cross-build: `make build-env-daemons`.
- Calls start in call order (frame write order is part of the daemon contract) and are awaited with `Wait`; don't reorder or parallelize frame writes.
- Windows/FIFO/descriptor behavior is platform-split (`*_unix.go`, `*_other.go`); change both sides together.
