# AGENTS.md — .github

GitHub config. In the fork only `workflows/personal-fork.yml` runs; every inherited upstream workflow is disabled and re-disabled after each sync.

## Fork workflow
- `personal-fork.yml`: job `sync` → `automation/release/personal-fork-sync.sh`; job `release` → `automation/release/personal-fork.sh`. Triggers: push to `main`, hourly, dispatch.
- Secrets: `PIG_UPDATE_SIGNING_KEY` (Ed25519 update signing), `PIG_SYNC_TOKEN` (fine-grained PAT, Contents + Workflows write; needed when an upstream merge changes `.github/workflows/`).
- Changing it → run locally with `act` first (`act workflow_dispatch -W .github/workflows/personal-fork.yml -j <job>`).
- Pin actions by full commit SHA, like upstream.

## Inherited upstream files
- Don't edit upstream workflows, `CODEOWNERS`, `CONTRIBUTING.md`, `SECURITY.md` for fork needs; every edit becomes a merge conflict. Fork behavior goes into `personal-fork.yml` or new files.
- Upstream CI design (required checks, tiers in `automation/ci/ci-tiers.txt`, merge queue, autofix bot, CI images, Nix vendor hash): `git show upstream/main:.github/AGENTS.md`. Its local gates still work: `make ci-bookkeeping ci-preflight`.

## Release reliability
Contract: each artifact identifies one source commit, target, toolchain and composition; published bytes equal built bytes; Stock artifacts never activate PiG Standard.

| Failure | Detection | Required result |
|---|---|---|
| Artifact differs from built bytes | recompute digests before publish | stop publication |
| Partial release visible | assets staged on a draft | publish only after all uploads |
| Version reused | tags + drafts reserve `-c.N` | never reuse a number |
| Archive escapes extraction root | traversal, links, devices, overwrite | installer rejects it |
| Update replaces a binary with bad data | truncated/corrupt download | verify checksum, size, signature first |
| Public evidence leaks private context | scan logs, paths, credentials | publish nothing private |
| Target advertised without native proof | native install + smoke | release notes say which targets ran |

Don't publish, tag or change visibility outside `personal-fork.yml` without explicit user approval.
