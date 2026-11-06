# AGENTS.md — extensions

Extension SDKs bridging Pi's extension API through PiG's subprocess host (D19).

| Path | What | Test |
|---|---|---|
| `sdk/` | Go SDK (own module); `upgrade/` rewrite rules for `pig extension upgrade` | `make test-sdk-go`, `make sdk-apidiff` |
| `sdk-rs/` | Rust SDK | `make test-sdk-rs` |
| `sdk-py/` | Python SDK | `make test-sdk-py` |
| `sdk-lisp/` | Common Lisp SDK (SBCL; `automation/ci/install-sbcl.sh`) | `make test-sdk-lisp` |
| `sdk-ts/` | declaration-only adapter: re-exports pinned Pi declarations + PiG-only types; no runtime | `make test-sdk-ts` |
| `RUBRIC.md` | Go vs Rust choice for new extensions | |

## Parity across SDKs
- All SDKs target the same wire contract (`coding/extension/host/subprocess/protocol.go`) and the same registration/host-call/tool-result shapes.
- Interface and behavior are identical across SDKs (modulo naming convention). Implementation may differ (host call vs replicated state, sync vs async, errors as `Result`/`(T, error)`/exceptions).
- Authority: `test/extension-conformance/`. A capability is done only when it lands in every SDK with a conformance row that fails when any one SDK is broken. Bind the row to a value distinguishable from every SDK's fallback.
- Don't return a constant where another SDK returns host data; that is a behavior divergence, worse than a missing method. If the host call is async and the API sync, replicate the value into state.
- An extension ports between SDKs by rewriting, never by rearchitecting or dropping a feature.

## Cross-updates
- SDK capability change → all four SDK dirs + conformance row + `docs/extension-api-parity.md` + `make sdk-surface-drift`.
- Breaking Go SDK change → upgrade rule + fixture under `sdk/upgrade/testdata/fixtures` + release note (`make sdk-apidiff` fails otherwise).
- `sdk-ts` Pi dependency must equal `coding.UpstreamVersion`; pin bumps update it.
- READMEs here are product surface; update them with behavior changes.
