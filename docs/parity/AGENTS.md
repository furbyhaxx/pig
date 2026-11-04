# AGENTS.md — docs/parity

Parity maintenance authorities.

| File | What | Edit |
|---|---|---|
| `PORT_MAP.md` | upstream file → Go file, status | by hand; `make port-map-drift` validates |
| `divergences/D<N>.md` | one divergence record per file | by hand |
| `DIVERGENCES.md` | assembled ledger | generated: `make divergences` |
| `DIVERGENCE-IDS.txt` | ID allocation notes | append when allocating |
| `KNOWN-GAPS-0.3.x.md` | known-gap block | generated: `make known-gaps` |
| `upstream-pin.md`, `dependency-source-proof.md`, `pending-tests-batches.md` | procedures | |
| other `*.md`, `gap-closure/` | upstream-owned closure notes | read only |

## Divergence records
- Active: file holds one `## D<N> <title>` section with `SCRUTINIZED:approved`, remove-when condition, parity coverage or allowance. Retired: one `- D<N> — ...` line.
- IDs are global across `DIVERGENCES.md` and additive ledgers (`docs/additive-features.md`); duplicates fail.
- Edit only `D<N>.md`, then `make divergences`. `make divergence-consistency` matches every record to its `// pig divergence (D<N>)` call sites; `make divergence-quality` validates contracts.
- Re-probe records at every upstream leap and release review. Delete or reclassify entries that are language mechanics, additive behavior, test limitations, unreachable/equal behavior, stale scaffolds or caller-free code.
- User-facing summary: `internal/pigdocs/content/divergences.md`.
