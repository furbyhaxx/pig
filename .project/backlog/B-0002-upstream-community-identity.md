---
document:
  type: "backlog-item"
  id: "B-0002"
  title: "Community files still name the upstream maintainer"
  description: "CODEOWNERS, SECURITY.md and CONTRIBUTING.md address upstream contacts and processes"
  status: "done"
  revision: 2
  created_at: "2026-11-04T12:30:29+01:00"
  updated_at: "2026-11-05T22:55:00+01:00"
  created_by: "orchestrator"
  updated_by: "coding-agent"
  changelog:
    - revision: 1
      reason: "document creation"
      changes: ["document creation"]
    - revision: 2
      reason: "Fork-owned community, conduct, support and security files"
      changes: ["replace upstream owners, email and contribution process", "delete the upstream parity issue form", "protect the files with merge=ours and record SPEC-0001 D11"]
  kind: "chore"
  priority: "P3"
  scope: [".github/**"]
---

# B-0002 Community files still name the upstream maintainer

Current: `.github/CODEOWNERS` assigns paths to `@MichaelKinsy`; `.github/SECURITY.md` routes reports to the upstream maintainer's email; `.github/CONTRIBUTING.md` requires DCO and signed commits for upstream contribution.

Target: decide per file whether the fork keeps the upstream text (merge-friendly) or states fork contacts and process.

Why: GitHub requests reviews and shows security contacts from these files on the fork.

Fix hint: decision first (user); fork-owned versions can join the `merge=ours` set in `.gitattributes` (SPEC-0001 D2).

## Implementation (05.11.2026)

- Fork-owned: `CODEOWNERS` (`* @furbyhaxx`), `SECURITY.md` (GitHub private reporting, no email), `CONTRIBUTING.md` (personal fork, no external contributions or upstream DCO/signature), `SUPPORT.md`, `CODE_OF_CONDUCT.md` (owner enforcement), `security-insights.yml` (fork URLs and owner), `pull_request_template.md`, `ISSUE_TEMPLATE/bug.yml` and `ISSUE_TEMPLATE/feature.yml`.
- Removed `ISSUE_TEMPLATE/parity.yml`; the fork is not bound to upstream parity policy. `ISSUE_TEMPLATE/config.yml` stays as the generic blank-issue configuration.
- All owned paths are in the `merge=ours` block in `.gitattributes`; SPEC-0001 D11 records the decision.
- Evidence: `make compliance` passes all seven checks, `security-insights.yml` parses, and the owned files contain no upstream maintainer, email or DCO process (the remaining `MichaelKinsy/PiG` strings are attribution links).
