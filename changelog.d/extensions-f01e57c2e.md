### Fixed

- Cancel active Go, Python and Rust event-bus listener contexts when the emitting request is cancelled, including Go factory-time dispatches, and stop dispatching subsequent listeners. Registration and unregistration remain unconditional. JavaScript listeners retain their `(data)` signature and cannot be preempted. See SPEC-0004 D3. <!-- fix: f01e57c2e -->

### Upgrading

- A Go, Python or Rust `pi.events` listener now runs with the emitting request as its parent, so its context ends when that request is cancelled. A handler that keeps work running in the background should stop when its context ends instead of assuming it always runs to completion. Node and TypeScript handlers are unchanged.
