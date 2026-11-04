# AGENTS.md — mcp

Port of Pi `packages/mcp`: standalone Model Context Protocol client (transport-neutral core, stdio, Streamable HTTP, SSE, OAuth). Must not import an MCP SDK or other PiG packages.

- Upstream: `.upstream/current/packages/mcp/`. `mcptest/` — in-memory transport; `oauth/` — auth flows.
- Transports own framing and I/O; `Client` owns correlation, init, timeouts, cancellation, server requests. Keep that split.
- Every unawaited upstream Promise is a goroutine owned by the client or transport and drained on `Close`.
- Official conformance suite: `make test-mcp-conformance` (`test/mcp-conformance/README.md`).
- Coding-agent integration lives in `coding/mcpext/` and `coding/extension/mcp_servers*.go`, not here.
