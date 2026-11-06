---
document:
  type: "backlog-item"
  id: "B-0001"
  title: "Native verification of fork release targets"
  description: "Fork releases execute only linux/amd64 before publication"
  status: "done"
  revision: 3
  created_at: "2026-11-04T12:30:29+01:00"
  updated_at: "2026-11-06T13:50:00+01:00"
  created_by: "orchestrator"
  updated_by: "coding-agent"
  changelog:
    - revision: 1
      reason: "document creation"
      changes: ["document creation"]
    - revision: 2
      reason: "Implement the native verify matrix; acceptance needs one fork release run"
      changes: ["record the workflow split and the three native targets", "record the local syntax/act evidence and the remaining external acceptance"]
    - revision: 3
      reason: "Owner halted macOS and Windows; first fork release on 0.7.0 accepted"
      changes: ["verify matrix is linux-arm64 only", "record the accepting run"]
  kind: "feature"
  priority: "P3"
  scope: ["automation/release/personal-fork.sh", ".github/workflows/personal-fork.yml"]
---

# B-0001 Native verification of fork release targets

Current: `automation/release/personal-fork.sh` runs `pig --version` only for linux/amd64. darwin, windows, linux/arm64 and android/arm64 archives are cross-built and published unexecuted (release notes say so).

Target: each published archive passes an install + `pig --version` + `pig update` smoke on its own OS/arch before the draft is published.

Why: a broken non-Linux artifact reaches `/releases/latest` and `pig update` installs it.

Fix hint: matrix job on `macos-15`, `windows-latest`, `ubuntu-24.04-arm` that downloads the draft assets; publish only after all pass. Android stays ELF-checked (`automation/release/android-binary.sh check`).

## Implementation (05.11.2026)

- `automation/release/personal-fork.sh` uploads the draft and writes `version` to `GITHUB_OUTPUT`; it no longer publishes.
- `.github/workflows/personal-fork.yml`: `release` builds the draft, `verify` downloads the matching draft asset on `ubuntu-24.04-arm`, `macos-15` and `windows-latest` and runs `pig --version` and `pig update --self --no-approve`, and `publish` edits the release only when every verify leg succeeded. An already-released source emits no version, so verify and publish skip.
- Local evidence: `bash -n automation/release/personal-fork.sh`, YAML parse, and `act -l -W .github/workflows/personal-fork.yml` lists `sync`, `release`, `verify`, `publish` in dependency order.
- 06.11.2026: macOS and Windows halted (SPEC-0001 D4 revision 7); the verify matrix is linux-arm64 only and darwin/windows archives are no longer built.
- Accepted 06.11.2026: run 42218594059 (`v0.7.0-c.1`) ran sync, release, the linux-arm64 verify leg and `publish` last, all green.
