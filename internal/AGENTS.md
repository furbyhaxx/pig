# AGENTS.md — internal

Shared non-public helpers. Search here before writing any JS/Node-semantics, JSON, path, process or file-system helper; duplicating one is drift.

| Group | Packages |
|---|---|
| JavaScript semantics | `ecmascript`, `jsarray`, `jsclock`, `jsdiff`, `jsjson`, `jsloop` (event-loop ordering), `jsmath`, `jsnumber`, `jsonparse`, `jsonstringify`, `jsparity`, `jsstring` (UTF-16/WTF-8), `orderedjson`, `scalarjson`, `jsonvalid` |
| Node semantics | `nodebase64`, `nodeerrno`, `nodefs`, `nodehome`, `nodepath`, `nodesemver`, `nodespawn`, `nodeurl`, `crossspawn`, `realpath`, `resolvepath` |
| Ported npm libraries | `minimatch`, `typebox`, `yaml12`, `jsdiff`, `wordsegmenter`, `latex`, `mermaid` |
| Files, processes | `filelock`, `pilock`, `fspublish`, `fsretry`, `execfile`, `ownerfile`, `childwait`, `jobgroup`, `linkerexec`, `termuxenv`, `shellconfig`, `nativeplatform` |
| Config, identity | `configroot` (config root resolution), `configvalue`, `coding/pigversion` (generated from `.upstream/pin`), `coding/pigidentity`, `contentdirs`, `pigdocs` (embedded docs bundle) |
| Extension host plumbing | `boundarywire`, `eventwire`, `changeobservers`, `handleryield`, `compiledin`, `frontendpack`, `pigsdklock`, `pigstrip` (`ids_generated.go` generated) |
| Checks, test support | `testenv`, `testbudget`, `testscaling`, `bodyreadhook`, `lockorder`, `pioracle` (runs compiled Pi modules as oracle), `perfgate`, `requestvalidity`, `upstreammirror`, `tui/termsim` |

- `codingagent/` — interactive app (Pi `modes/interactive`), self-update, changelog display. Divergence-guard scope.
- `pigdocs/content/` — agent/offline docs bundle shipped in the binary; product surface (`docs/AGENTS.md`).
- New helper → row in this table.
