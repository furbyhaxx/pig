### Fixed

- Cancel active Go, Python and Rust event-bus listener contexts when the emitting request is cancelled, including Go factory-time dispatches, and stop dispatching subsequent listeners. Registration and unregistration remain unconditional. JavaScript listeners retain their `(data)` signature and cannot be preempted. See SPEC-0004 D3. <!-- fix: f01e57c2e -->
