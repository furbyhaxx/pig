# AGENTS.md — codemode

Port of Pi `packages/codemode`: runs model-written JavaScript against host tools in QuickJS compiled to WebAssembly, hosted by wazero. Spec: `docs/specs/builtin-codemode-tool-search.md`.

- Upstream: `.upstream/current/packages/codemode/`. `assets/` — embedded QuickJS build; `cloudflare/` — Cloudflare executor.
- The script reaches the host only through the single bridge function. Don't add a second host channel.
