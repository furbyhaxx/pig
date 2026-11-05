# Personal PiG distribution

`furbyhaxx/pig` is Furby Haxx's personal fork of [PiG](https://github.com/MichaelKinsy/PiG), created by Michael Kinsy and originally developed at Hewlett Packard Enterprise. It is not an official Pi project and does not imply endorsement. This fork is a base for the owner's coding agent. Compatibility with Pi and upstream PiG is kept while it does not limit correctness, performance, design quality or the owner's requirements; it is not a requirement (SPEC-0001 D7). Breaking changes require a fork specification and a user migration note. Executable names, module identities, configuration paths, sessions and credentials retain their inherited names and formats until a specification decides otherwise.

## Install or switch from upstream

On Linux, macOS, or Android/Termux arm64:

```sh
curl -fsSL https://github.com/furbyhaxx/pig/releases/latest/download/install.sh | sh
```

To replace a particular installation, first run `command -v pig`, then select its directory:

```sh
curl -fsSL https://github.com/furbyhaxx/pig/releases/latest/download/install.sh | PIG_INSTALL_DIR="$HOME/.local/bin" sh
pig --version
pig update
```

Do not run a second copy earlier on `PATH`. For a package-manager installation, uninstall the old package through that manager before switching; do not overwrite files the manager owns. The installer does not delete user configuration. It writes a new standalone receipt, `update-url`, and `update-trust.pem` in the existing config root (`PIG_HOME`, then `XDG_CONFIG_HOME/pig`, then `~/.pig`). Explicit `PIG_UPDATE_URL` and `PIG_UPDATE_TRUST_ROOT` overrides still take precedence; unset old upstream overrides when switching.

The archive checksum is verified before installation. The fork public key ships in the verified archive and is also embedded in release binaries. Subsequent updates require a manifest signed by the fork key and verify the archive checksum. Old upstream receipt/key sidecars cannot safely be reused; switch once through this installer, not through upstream's `pig update`.

Windows releases contain `pig.exe`; stop PiG and replace it from the matching ZIP after checking `SHA256SUMS`. Windows standalone in-place updates remain unsupported. The upstream PowerShell installer and npm package are not fork distribution channels.

## Automatic synchronization and releases

The `Personal fork` workflow runs hourly, on pushes to `main`, and through manual dispatch. It merges `MichaelKinsy/PiG:main` into `main` (`automation/release/personal-fork-sync.sh`). A conflict stops the workflow and publishes nothing. Resolve the merge locally, preserve custom changes, and push the resolution to `main`. There is no force synchronization. The only automatic resolution is for the fork's own `AGENTS.md` files, which always keep the fork version; the merge commit lists upstream changes to them for review. Run `git config merge.ours.driver true` once in a local clone so local merges behave the same.

Upstream merges that change `.github/workflows/` can only be pushed with a token that has the workflows permission, which `GITHUB_TOKEN` never has. Store a fine-grained personal access token for `furbyhaxx/pig` with *Contents* and *Workflows* read/write as the Actions secret `PIG_SYNC_TOKEN`. Without it, such merges fail and the sync stops until you merge them locally.

After a clean merge, a new source commit produces a release. Already released source is skipped. Builds cover Linux, macOS, and Windows on amd64/arm64, plus Android arm64 using the NDK. Only Linux amd64 is executed by this workflow; other targets are cross-built, not natively verified. This lightweight personal pipeline does not claim the upstream candidate pipeline's full security, parity, SBOM, or attestation evidence.

Versions use the upstream **PiG** version, not its separate Pi pin: `0.5.3-c.1`, `0.5.3-c.2`, then `0.5.4-c.1`. The counter advances for each new released source commit within a base version, including upstream merges. Tags from failed/draft releases are reserved and never reused. Numeric identifiers sort correctly beyond 9. Although these versions are SemVer prereleases, GitHub marks them as full releases so `/releases/latest` selects the personal channel. Switching from upstream's same-base stable version requires the installer because SemVer correctly considers `-c.1` lower than the stable version.

The source branch retains upstream version pins to minimize merge conflicts. Release builds apply `make set-version VERSION=x.y.z-c.N` in the disposable build tree; plain development builds retain the upstream version. SDK dependencies and Go module paths are not renamed or republished.

`PIG_UPDATE_SIGNING_KEY` is the fork repository's Actions secret. Its public half is `automation/release/personal-update-trust.pem`. Never commit the private half. Keep a trusted backup before replacing the secret: key rotation requires an overlap release trusted by existing installations. GitHub cannot return an existing secret value.

Inherited upstream workflows are disabled in this fork until individually reviewed. The personal workflow is the supported release path; upstream package publishing, hosted-site, notification, and image automation must not run implicitly under the fork identity.

## Local checks

Run `make fork-check` before a commit (`automation/make/fork.mk`). It builds the module, lints and tests the packages a change affects, vets the changed packages for Windows, and runs the engineering ratchets (`divergence-guard`, `port-lint`, `perf-lint`, `any-trap`, `go-fix-clean`), `source-hygiene`, `docs-drift` and `changelog-check`.

Upstream's Pi-comparison, porting-ledger and publication gates (`make check`, `make verify`, `ci-drift`, `ci-contracts`, the `parity*` and `journeys*` targets, `compliance`, `release-check`, Pig Porter) remain in the tree for clean merges but are not part of the fork's checks and may fail (SPEC-0001 D10). `docs/project/test-economy.md` describes upstream's tiers.
