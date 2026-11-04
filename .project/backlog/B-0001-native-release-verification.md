---
document:
  type: "backlog-item"
  id: "B-0001"
  title: "Native verification of fork release targets"
  description: "Fork releases execute only linux/amd64 before publication"
  status: "open"
  revision: 1
  created_at: "2026-11-04T12:30:29+01:00"
  updated_at: "2026-11-04T12:30:29+01:00"
  created_by: "orchestrator"
  updated_by: "orchestrator"
  changelog:
    - revision: 1
      reason: "document creation"
      changes: ["document creation"]
  kind: "feature"
  priority: "P3"
  scope: ["automation/release/personal-fork.sh", ".github/workflows/personal-fork.yml"]
---

# B-0001 Native verification of fork release targets

Current: `automation/release/personal-fork.sh` runs `pig --version` only for linux/amd64. darwin, windows, linux/arm64 and android/arm64 archives are cross-built and published unexecuted (release notes say so).

Target: each published archive passes an install + `pig --version` + `pig update` smoke on its own OS/arch before the draft is published.

Why: a broken non-Linux artifact reaches `/releases/latest` and `pig update` installs it.

Fix hint: matrix job on `macos-15`, `windows-latest`, `ubuntu-24.04-arm` that downloads the draft assets; publish only after all pass. Android stays ELF-checked (`automation/release/android-binary.sh check`).
