---
document:
  type: "specification"
  id: "SPEC-0001"
  title: "Personal PiG distribution"
  description: "Fork identity, upstream synchronization, versioning, release, update trust and agent-steering ownership of furbyhaxx/pig"
  status: "published"
  revision: 1
  created_at: "2026-11-04T12:30:29+01:00"
  updated_at: "2026-11-04T12:30:29+01:00"
  created_by: "orchestrator"
  updated_by: "orchestrator"
  changelog:
    - revision: 1
      reason: "document creation"
      changes: ["record the decisions behind the fork setup and the AGENTS.md rewrite"]
---

# SPEC-0001 Personal PiG distribution

## Context

`furbyhaxx/pig` is a personal distribution of `MichaelKinsy/PiG`. It does not contribute upstream. It must track upstream continuously, build and publish its own releases, update installed binaries from its own releases, and stay a drop-in replacement that reuses the user's existing PiG configuration.

## Decisions

### D1 Identity

- Executable `pig`, config root (`PIG_HOME` → `XDG_CONFIG_HOME/pig` → `~/.pig`), session and credential formats, Go module path `github.com/MichaelKinsy/PiG`: unchanged.
- Reason: renaming the module or paths would conflict with every upstream merge and break config reuse.
- User-facing install and update docs point to `furbyhaxx/pig` releases (`README.md`, `docs/personal-fork.md`, `internal/pigdocs/content/install.md`, `docs/site/public/install.sh`).

### D2 Upstream synchronization

- `.github/workflows/personal-fork.yml` job `sync` merges `MichaelKinsy/PiG:main` into `main` hourly, on push and on dispatch.
- Real conflict → job fails, nothing is published. A human resolves locally and pushes. No force, no automatic conflict resolution.
- Fork-owned files resolve to the fork side automatically through `.gitattributes` `merge=ours` (driver `merge.ours.driver true`, configured by the sync job): every `AGENTS.md` except fixture files under `testdata/`. The sync job lists upstream changes to those files in the merge commit message and the job summary, because they may carry new upstream rules to port into the fork's files.

### D3 Versioning

- Release version = upstream `PigVersion` + `-c.N`, e.g. `0.5.3-c.1`. N starts at 1 per upstream version and increases for every released source commit (upstream merges and fork changes alike).
- Tags of drafts and failed publications reserve their N; numbers are never reused.
- `main` keeps the upstream `PigVersion`; `make set-version VERSION=x.y.z-c.N` runs only in the disposable release build tree. Reason: committing the custom version would conflict with every upstream version bump.
- SemVer orders `x.y.z-c.N` below `x.y.z`. Switching from an upstream install of the same base therefore goes through the fork installer, not `pig update`.

### D4 Release and update trust

- `automation/release/personal-fork.sh` builds linux/darwin/windows × amd64/arm64 and android/arm64 (NDK), signs `update.json` with Ed25519, and publishes a GitHub release through a draft. Already-released source commits are skipped.
- Signing key: Actions secret `PIG_UPDATE_SIGNING_KEY`; public key committed at `automation/release/personal-update-trust.pem`, embedded in binaries and shipped in archives. The installer writes it to `<config-root>/update-trust.pem` together with `update-url`.
- GitHub marks `-c.N` releases as full releases so `/releases/latest` serves the fork channel.
- Only linux/amd64 is executed in the pipeline; other targets are cross-built. The pipeline does not reproduce upstream's SBOM, attestation and parity evidence.

### D5 Inherited workflows

- All upstream workflows except `personal-fork.yml` are disabled in the fork and re-disabled by the sync job when a merge brings new ones. Reason: they publish to upstream channels (npm, site, Discord, images) or spend runner minutes on upstream's governance gates.

### D6 Agent steering

- All `AGENTS.md` files follow skill `agents-md`; artifacts follow skill `dot-project` under `.project/`.
- Root `AGENTS.md` carries no generated status: the parity coverage block was removed and `automation/ci/check-public-claims.py` reads `test/parity/coverage.md` instead.
- No test or gate reads `AGENTS.md` content. Upstream couplings were removed: `TestGoVersionPolicy` no longer pins a sentence in `AGENTS.md`; `TestAgentsDocRendersTheTierTable` was deleted with the upstream CI tier table it checked. Reason: with `merge=ours`, such a test fails as soon as upstream moves the data it compares against.
- `changelog.d/` has no `AGENTS.md`: the fragment checker accepts only fragments there; fragment rules live in the root file.
- Upstream-owned artifact directories (`docs/plan/`, `docs/specs/`, `docs/findings/`, `plans/`, `research/`, `handoffs/`) stay in place to keep merges clean.

## Consequences

- Upstream edits to `AGENTS.md` never block the sync but are silently superseded; the merge commit message is the review trigger.
- Upstream edits to other files the fork changed (installer, release tooling, README, docs bundle install page, public-claims gate) can still conflict and stop the sync.
