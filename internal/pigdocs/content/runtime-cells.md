# Runtime cells

Runtime cells are PiG's internal process-packing optimization. Authors select a
conventional factory or an exact standalone. Authors do not configure cells.

PiG can pack Go, Rust, Python and Node factories by language. Standalones run in isolated processes. Each packed extension keeps one socket
and one registration handshake. Packed behavior must match isolated behavior.

Each extension start receives a generation-unique socket path. Teardown removes only that generation's socket, so retiring an old generation cannot remove its successor's socket (SPEC-0004 D1).

`Host.Reload` and CLI reloads bound connection and registration by the global `extensionsReloadHandshakeMs` setting (default 15000 milliseconds), starting at spawn or factory admission and excluding build time. Initial loading stays cancellation-only. A failed cell keeps its serving generations and discards its staged declarations and siblings. Newly spawned failed processes are killed and reaped; a retained process stays alive for its serving generations. Other cells may commit, and the CLI updates its active set before returning the failure to `ctx.Reload()`, `/reload` or RPC. User factory side effects cannot be rolled back (SPEC-0004 D2).

PiG derives a cell cache key from source content, SDK content, generated runner
content, composition, language runtime, and toolchain. `/reload` reuses an exact
cache hit and rebuilds changed input.

A shared-process exit quarantines that cell. One stalled handler or one logical
member failure does not quarantine healthy siblings. The next reload fissions a
quarantined composition into isolated processes.
