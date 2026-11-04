---
document:
  type: "backlog-item"
  id: "B-0002"
  title: "Community files still name the upstream maintainer"
  description: "CODEOWNERS, SECURITY.md and CONTRIBUTING.md address upstream contacts and processes"
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
  kind: "chore"
  priority: "P3"
  scope: [".github/**"]
---

# B-0002 Community files still name the upstream maintainer

Current: `.github/CODEOWNERS` assigns paths to `@MichaelKinsy`; `.github/SECURITY.md` routes reports to the upstream maintainer's email; `.github/CONTRIBUTING.md` requires DCO and signed commits for upstream contribution.

Target: decide per file whether the fork keeps the upstream text (merge-friendly) or states fork contacts and process.

Why: GitHub requests reviews and shows security contacts from these files on the fork.

Fix hint: decision first (user); fork-owned versions can join the `merge=ours` set in `.gitattributes` (SPEC-0001 D2).
