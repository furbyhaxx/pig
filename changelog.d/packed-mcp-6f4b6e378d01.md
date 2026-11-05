### Fixed

- A packed Node cell exposes an earlier member's MCP registrations to a later member's factory, so `pi.getMcpServers()` sees them and a duplicate `pi.registerMcpServer()` throws to that factory as it does in Pi. <!-- fix: 6f4b6e378d01 -->
