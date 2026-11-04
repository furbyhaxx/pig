# Install and use PiG

Use this page to get from a fresh binary to a working session. Return to the
[docs router](README.md) for focused tasks.

## 1. Build or install the binary

This is the `furbyhaxx/pig` personal distribution. On macOS, Linux, or Android with Termux arm64, install from this fork's GitHub Releases:

```sh
curl -fsSL https://github.com/furbyhaxx/pig/releases/latest/download/install.sh | sh
pig update
```

To replace an existing standalone installation, set `PIG_INSTALL_DIR` on the `sh` side of the pipe to the directory containing the current `pig`. The installer preserves settings, sessions, and credentials. It replaces the installation receipt, update URL, and public trust-key sidecars so subsequent updates use signed personal releases. Explicit update environment overrides still take precedence. Uninstall a package-manager installation through its manager before switching; do not overwrite package-owned files.

On Windows, stop PiG and replace `pig.exe` with the matching ZIP from [this fork's releases](https://github.com/furbyhaxx/pig/releases), after checking `SHA256SUMS`. Windows standalone in-place updates remain unsupported.

The upstream npm package, upstream Go module installation, hosted website installer, and PowerShell installer install upstream, not this fork. This fork does not republish those packages. From this fork's checkout, build with the toolchain in `.go-version`:

```sh
go build -trimpath -o bin/pig ./cmd/pig
./bin/pig --version
```

From a source checkout, build it with Go 1.27.3. Plain source builds keep the upstream version pin; personal release builds apply `x.y.z-c.N` in a disposable build tree. Go 1.27 release binaries for macOS require macOS 13 or later. Nix users can build `github:furbyhaxx/pig` and update that flake input; `pig update` does not replace a Nix installation.

PiG stores user state under `~/.pig`. Set `PIG_HOME` only for an isolated config
root, tests, or CI.

### Change what the install script does

`install.sh` reads these environment variables. Set them on the `sh` side of the pipe, for example `curl -fsSL https://github.com/furbyhaxx/pig/releases/latest/download/install.sh | PIG_VERSION=0.5.3-c.1 sh`.

| Variable | Meaning |
|---|---|
| `PIG_VERSION` | Install this release (for example `0.5.0`) instead of the latest. |
| `PIG_INSTALL_DIR` | Directory that receives `pig`. The default is `$HOME/.local/bin`, or `$PREFIX/bin` in Termux. |
| `PIG_HOME` | PiG's settings directory, where the install receipt goes. The default is `~/.pig`. |
| `PIG_API_BASE` | Optional API that names the latest release. Unset by default; the installer uses GitHub directly. |
| `PIG_LATEST_RELEASE_URL` | The page that redirects to the latest release, used when the API does not answer. |
| `PIG_DOWNLOAD_BASE` | The release download root. The default is the GitHub Releases download URL of this repository. |
| `PIG_UPDATE_URL` | The update manifest that `pig update` uses. The default is the latest release's `update.json`. |

The script verifies the archive against the release's `SHA256SUMS` before it writes anything to the install directory, whichever mirror these variables select.

### macOS Gatekeeper after a manual download

A browser download can add the local `com.apple.quarantine` attribute to a PiG
archive. First verify the published SHA-256 checksum. For a verified PiG binary,
run:

```bash
xattr -p com.apple.quarantine /path/to/pig
xattr -d com.apple.quarantine /path/to/pig
/path/to/pig --version
```

Do not recursively remove quarantine attributes from a broad directory.

## 2. Authenticate a model provider

```bash
pig login github-copilot
```

API-key providers may use their documented environment variable:

```bash
export OPENAI_API_KEY=...
pig --model openai/gpt-4.1
```

See [Providers](providers.md). Product extensions may contribute additional
login targets, such as a platform host, through the same `pig login` surface.
Use `pig login --help`/`pig --help` for the targets available in your binary.

## 3. Install packages

Raw PiG accepts local, git, npm, and HTTP package sources:

```bash
pig install ./my-package
pig install git:https://github.com/example/my-package.git
pig install npm:@example/package
```

An installed product extension may add another resolver scheme, such as a
private marketplace. Its authentication and catalog commands are product
features, not assumptions in raw PiG.

## 4. Run PiG

```bash
pig                              # interactive TUI
pig -p "summarize this repo"     # one-shot prompt
pig --model github-copilot/gpt-5-mini
```

## 5. Install and run a piglet

Install and validate a registerable Piglet, then invoke it by name:

```bash
pig piglet add ./piglets/my-agent.yaml
pig piglet validate my-agent
pig --piglet my-agent
```

A Piglet with relative `local:` extension or skill origins is source-bound.
Run it directly instead of registering it:

```bash
pig piglet validate ./piglets/my-agent.yaml
pig --piglet ./piglets/my-agent.yaml
```

See [Piglets](piglets.md) for schema and resolution.

## 6. Build an explicit Piglet entry point or Binary

A source script is a thin Host Piglet entry point. It is not managed by PiG and
creates no PiG state or record:

```bash
pig piglet build my-agent --format script --out ~/.local/bin/my-agent
my-agent
```

On Windows the script is a cmd.exe batch file, so give it a `.cmd` name, for
example `--out my-agent.cmd` (D69).

A Piglet Binary is a target-native executable for one immutable Piglet
composition:

```bash
pig piglet build my-agent --format binary --out ./pig-my-agent
./pig-my-agent
```

The Piglet artifact lifecycle lives under `pig piglet`.

A Piglet Binary may fuse compatible Go components and host exact subprocess
components through PiG's ordinary subprocess extension host. Its registered
Piglet release records where each component comes from: binary, release
closure, required `agentEnv`, or external mapping. Subprocess use does not create
a different binary kind.

Piglet Image is the OCI carrier when the Piglet needs an environment-complete
runtime. Artifact startup never silently installs, pulls, logs in, or substitutes
an origin.

## 7. Set up an extension toolchain only when authoring/building

You need a language toolchain to author a source extension or build its release
closure. A consumer needs no compiler for a prebuilt native component; interpreted
subprocess components still need the exact runtime recorded by the Piglet
release or supplied by `agentEnv`.

Check before installing anything:

```bash
go version
cargo --version
python3 --version
node --version
```

Prefer a version manager already present on the machine. Otherwise use the OS
package manager or the language's official per-user installer, with the user's
consent. Do not remove or replace an existing system toolchain without asking.

For a Go extension, the required version is the `go` directive in the staged SDK
module. `pig extension init` stages an offline-resolvable project and
`pig install <path> --validate-only --json` validates it.

## Troubleshooting

| Symptom | Likely cause | Action |
|---|---|---|
| `pig` not found | binary directory absent from `PATH` | add the install directory |
| no selectable models | no provider credential | run `pig login <provider>` or set its API key env var |
| contributed source scheme is unknown | product resolver is not installed/loaded | install or run the product distribution that owns it |
| piglet origin does not resolve | missing local source/package/product resolver | run `pig piglet validate <name>` and fix the reported origin |
| source extension fails on first use | missing language toolchain | install it with consent or use a prebuilt Piglet release/Image |
| Piglet artifact cannot start a component | registered release/environment/external requirement missing | rebuild with a supported format or follow the reported mapping remediation; do not substitute PATH content |

## Offline docs

```bash
pig docs list
pig docs show install
pig docs path
```
