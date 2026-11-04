---
document:
  type: "specification"
  id: "SPEC-0010"
  title: "Packed build inputs: source snapshots and build tags"
  description: "Packed cells compile a snapshot of their source set, so concurrent edits never fail a build, and extensions declare Go build tags in their manifest"
  status: "draft"
  revision: 1
  created_at: "2026-11-04T15:01:00+01:00"
  updated_at: "2026-11-04T15:01:00+01:00"
  created_by: "orchestrator"
  updated_by: "orchestrator"
  changelog:
    - revision: 1
      reason: "Resolve U-25 of the host-gap verification"
      changes: ["document creation"]
---

# SPEC-0010 Packed build inputs

Resolves U-25 of `../analyses/2026-11-04-host-gap-verification/README.md`. The report's claim that `.project/` is packed is wrong: Go cell identity covers only the Go source set (`coding/extension/host/runtimecell/hash_go_module.go`).

## D1 Source snapshots

- Before compiling a Go packed cell, the host copies the cell's source set (the files `hashGoModule` hashes: `go.mod`, `go.sum`, non-test sources, embedded files) into a staging directory, hashes the copy and compiles the copy.
- Copy uses reflinks where the filesystem supports them, plain copies otherwise; never hardlinks (an in-place edit would change the snapshot).
- The artifact is addressed by the snapshot's digest. An edit during the build no longer fails it (`cell_inputs.go` "changed while it was being built" path is removed for snapshot builds); the next load sees a new digest and rebuilds.
- Benchmark before/after on a cold build of the largest example cell and of a 50k-file module; snapshot overhead reported.
- Rust and Node packed cells adopt the same mechanism in a later step once their source sets are defined.

## D2 Build tags

- Extension manifest field `build.go.tags: ["tag", …]` (Package/Piglet manifest, `docs/pig-package-spec.md`).
- Tags join the cell digest; extensions with different tag sets never share a cell.
- Validation: tags match `^[A-Za-z0-9_.]+$`; invalid → manifest error.

## Docs

- `runtime-cells.md` / `docs/runtime-cell-build-cache.md`: snapshot builds.
- `docs/pig-package-spec.md`, `packages.md`, `piglets.md`: `build.go.tags`.
