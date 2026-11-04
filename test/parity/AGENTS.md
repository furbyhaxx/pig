# AGENTS.md — test/parity

For unchanged inherited behavior, pinned Pi is the comparison oracle. Fork behavior follows `.project/specs/`; rewrite or delete scenarios that pin intentionally replaced behavior, rather than masking differences. Scenarios (`scenarios/<family>/*.toml`, families in `families.toml`) run Pi and PiG through `runner/` and compare. Read `README.md` before editing a scenario.

## Contract
- A scenario proves only its asserted path. Generated inventories define review obligations, not acceptance.
- A comparator exposes every observed difference unless a numbered divergence owns it.
- Probe upstream first; record pin, environment, inputs, raw output, comparator, expected result before production changes. At least one compiling mutation per new behavioral scenario.

## Scenario rules
- `covers = [...]`: exact upstream paths from PORT_MAP that the asserted output exercises. Startup, registration or import doesn't count. Keep it narrow.
- Comparator strength: `wait_contains` → `output_normalized_equal` → `output_equal` → `escaped_output_equal`. Tighten, never weaken to keep a bug green. Raise `runs` where flake risk matters; `runs = 1` is wrong for non-trivial behavior.
- Weak tags: `boot-only`, `registration-only` (catalog/auth wiring only, never stream/payload files), `smoke-only`, deferred. They don't count as behavioral verification.
- Comment every `normalize_replace`, `lint-known-gap`, `skip`, `serial`, group override or weak tag. Never normalize ANSI/whitespace to hide a real difference.
- Output differs → fix a defect, cite an inherited `D<N>`, or rewrite the affected scenario for a spec-decided fork change; a comment alone is not enough.
- Writable paths via `{{TEMP}}`; never write into checked-in fixture roots; agent/home dirs are ephemeral copies. `pre_clear_paths` is deprecated.
- Default scheduler groups unless `group`/`serial` is justified.
- Shared behavior fix → re-probe affected scenarios in other families.
- `make lint-scenarios` enforces part of this.

## Done criteria for a surface
1. Canonical scenario under `scenarios/<family>/`, derived from upstream first, behavioral quality.
2. `covers` exact and exercised; legacy duplicate/weak claims removed.
3. Drift fixed or numbered.
4. Unit + caller-path tests for cross-layer bugs; `go test ./...` and the family gate pass.
5. Upstream leap: `make upstream-delta` (every changed tracked file has a disposition + evidence).
6. `make interface-inventory-drift` and `make interface-mapping-strict` pass.
7. `make async-contracts` passes.
8. Interaction/state-machine behavior not expressible as declarations → `behavior-contracts.toml` (hashed upstream range, keybinding/event inventory, Pig targets, failing unit regression, behavioral scenario). Boundary transitions (wraparound, empty/singleton, cancellation, ordering, persistence) are separate obligations.

## Ledgers
- `interfaces/behavior-inputs-v<ver>.json` — upstream input/render surfaces; reviewed mapping per family. `make behavior-contracts` (drift), `make behavior-contracts-strict` (no pending/partial). Generated facets never prove semantics.
- `interfaces/test-mapping-v<ver>.json` — upstream test → Go test locations and hashes. `partial` lists missing cases.
- `interfaces/test-porting-policy-v<ver>.json` — `hot-path` tags and ported baseline; `make ci-contracts` rejects closing pending/partial hot paths or lowering the baseline. Approved `deferred-0.3.x` rows: `README.md#approved-03x-test-porting-gaps`. Never remove a hot-path tag or lower the baseline to pass.
- `unit-evidence/*.json` — reviewed mutation-proven unit tests counted as behavioral evidence.
- Test-file batches: `docs/parity/pending-tests-batches.md`.
- `coverage.md` is generated (`make coverage`); fix `cmd/coverage/` if the math is wrong.

## Failure modes
| Failure | Detection | Required result |
|---|---|---|
| Scenario passes without reaching the feature | mutate the owning branch | scenario fails behaviorally |
| Normalization hides a defect | inspect raw vs normalized | fix PiG or cite approved divergence |
| Crop excludes the difference | compare the complete owning surface | difference stays visible |
| Timing changes queue order | barriers and acknowledgements | exact order, no retries |
| Generated count treated as proof | trace to reviewed mappings and execution | denominator and evidence separate |
| Oracle version drifts | binary version, tag, commit, package lock | fail before comparison |
| Evidence not repeatable | replay command in isolation | same contract reproduced |

Keep failure and final-pass artifacts outside the source tree. Don't add retries, widen timeouts or rewrite expected output to hide a diagnosed defect.

## Loop discipline
Value order: (1) find and fix a Pi/PiG difference at the source or number it; (2) tighten a comparator; (3) promote silent drift to a numbered divergence; (4) honest new file coverage. Scenario counts and green checks are not value.

Stop and report instead of claiming success when: everything passed first try; a scenario passed before and after the fix; the loop found zero bugs, tightenings and divergences; a scenario asserts only banner/prompt/ready text; `family-gaps` is green because one broad scenario claims many files; a comment explains a rendering difference without a `D<N>`.

Family-loop report (required by the Porter skill):
```md
### Loop self-report
Family: <name>
Bugs fixed at the source (N): <cause> → <fix> @ <file:line>
Comparators tightened (N): <scenario>: <before> → <after>
Divergences numbered (N): D<N> <one line>
Lint suppressions added/changed (N): <rule> @ <file>: <why>
New PORT_MAP ✅ coverage (N): <upstream path>
Band-aids (N, should be 0): <site>: <symptom>: <why deferred>: <backlog item>
Scope expansion (N): <file>: <reason>: <call sites audited>
Cross-family scenarios re-probed (N): <family>/<scenario>: <change>
Smells investigated (N): <smell> → <finding>
Smells deferred to user (N, 0 unless blocked): <smell>: <blocker>
```
All-zero counts = failed loop unless the user accepts it.
