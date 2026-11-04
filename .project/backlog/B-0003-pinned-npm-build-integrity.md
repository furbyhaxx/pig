---
document:
  type: "backlog-item"
  id: "B-0003"
  title: "Pinned Pi build fails npm integrity verification"
  description: "Pinned source tarballs differ from the locked npm integrities, blocking host-gap gates"
  status: "done"
  revision: 2
  created_at: "2026-11-04T15:28:56+01:00"
  updated_at: "2026-11-04T17:13:00+01:00"
  created_by: "coding-agent"
  updated_by: "coding-agent"
  kind: "bug"
  priority: "P1"
  scope: ["automation/pin/**", "extensions/sdk-ts/**"]
  changelog:
    - revision: 1
      reason: "Found while verifying PLAN-0001 S0"
      changes: ["record prerequisite failure without overriding integrity checks"]
    - revision: 2
      reason: "Confirmed runtime-dependent compression; rebuilt with the repository-pinned Node"
      changes: ["record seven matching frozen package hashes", "add rejecting-runtime regression and preflight", "close integrity blocker"]
---

# Pinned Pi build integrity failure

## Observed

On 04.11.2026, `make upstream-mirror` and `make test-changed` fail in `parity-deps`: the pinned source build for commit `6c384aa30f619d7b28963f2bab24d8dacbe73cb6` packs Pi 1.4.0 packages, then `npm ci --ignore-scripts --no-audit --no-fund` reports `EINTEGRITY`.

For pi-coding-agent, locked SHA-512 is `htEa1twE2jJnjWV8Yru5xnLL/44mGFkdeBCF2C96V6vJxTvH1yVi35CEphEquBZyDOII7V93xleEOhU41xuetw==`; rebuilt SHA-512 is `4xlXry7boWMDh2S73wyU9bW2hMXt9gMTcPN/YJNamltIKKRv48qY6DcszbE8EJvpdz2+QEfPFecsE0vL7+j1RQ==`. Other local Pi packages also report integrity mismatches.

The direct mirror script extracts the source but cannot hydrate generated provider data without the matching built package. It does not establish `.upstream/current`. `make ci-contracts` then fails correspondence extraction, and `make check` fails Node bundle regeneration because `extensions/sdk-ts/node_modules` is missing.

## Required outcome

Reproduce the source-pack/lock mismatch at the pin builder, identify why packed bytes differ, and restore reproducible verified installation. Do not disable integrity verification or silently accept newly generated hashes. Then rerun PLAN-0001 S0 gates before S1.

## Resolution

The host used Node 26.13.0; `.node-version` pins 24.22.1. With the same built source and npm 12.1.0, packing on Node 24.22.1 reproduces all seven frozen workspace package SHA-512 values exactly; Node 26.13.0 does not. Changing npm alone did not repair the mismatch. Archive compression is runtime-dependent, not a changed dependency closure.

`pack()` now checks the Node version against `.node-version` before creating or removing output. `ClosureTests.test_pack_rejects_unpinned_node_before_running_npm` is red-proven; all 64 pin unit tests pass. The developer guide explains the runtime requirement. No frozen hashes, source pin or dependency versions changed.

The pinned runtime was selected on the verification command PATH; the system Node installation was not replaced. `make upstream-mirror` succeeds after refreshing the partial mirror left by the failed setup. `make ci-drift ci-contracts` and `make lint-changed` pass. Remaining full-check failures are tracked in B-0004.
