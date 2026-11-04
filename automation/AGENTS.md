# AGENTS.md — automation

Every non-Go script. The root `Makefile` and `make/*.mk` are the interface; folder map in `README.md`. Add a make target when a script is called by hand twice.

- `release/personal-fork.sh`, `release/personal-fork-sync.sh`, `release/personal-update-trust.pem` — fork-owned release and sync (`.project/specs/0001-personal-distribution.md`). Other scripts are inherited; weigh merge cost against design quality under the root fork policy.
- Ratchet baselines are shrink-only; a fix deletes its entries; never add entries for new code: `ci/divguard/baseline.toml`, `ci/portlint/baseline.toml`, `ci/anytrap` baseline, `perf/ratchet.toml`, `ci/changelogfrag/legacy.txt`.
- `changelog.d/` accepts only fragments; `ci/changelogfrag` rejects any other file there.
- `ci/autofix-paths.txt` / `ci/autofix-gates.txt` — what the upstream autofix bot may regenerate; a new generated file needs a reviewed entry.
- `ci/check-public-claims.py` reads coverage figures from `test/parity/coverage.md`.
- Scripts find the repo root from their own location; keep it that way. Python tests: `python3 -m unittest <file>` (no pytest dependency).

## Cross-updates
- Release asset names, archive layout or update manifest change → `docs/site/public/install.sh`, `test/ci-images/install_sh_test.go`, `internal/codingagent` self-update tests, `docs/personal-fork.md`, `internal/pigdocs/content/install.md`.
- `release/personal-update-trust.pem` rotation → Actions secret `PIG_UPDATE_SIGNING_KEY` in the same change; installed clients only trust the old key until they reinstall.
