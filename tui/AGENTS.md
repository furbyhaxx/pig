# AGENTS.md — tui

Port of Pi `packages/tui`: Main Screen, differential renderer, Editor, components, keybindings, alt screen. Interactive-mode components ported from `packages/coding-agent` also live here.

- Upstream: `.upstream/current/packages/tui/src/tui-main-screen.ts`, `components/editor.ts`. `main_screen_upstream_parity_test.go` documents a Pi 0.84 model; confirm it against the pinned source before relying on it.
- `internal/hljs/` — snapshot of Pi's highlight.js grammars (`make highlight-grammars`); `internal/tui/termsim` — terminal simulator for tests.
- `**/testdata/**` terminal goldens keep blank rows and mixed indentation (`.gitattributes`); don't let editors normalize them.
- Width/units: terminal cells, not bytes or runes; reuse `internal/jsstring`, `internal/wordsegmenter`.

## Contract
- Main Screen preserves scrollback. Editor changes only the required visible rows. Viewport behavior follows the documented contract; Pi comparisons cover unchanged inherited behavior. Frame work never scales with unchanged Session history.

## Failure modes
| Failure | Detection | Required result |
|---|---|---|
| Ordinary editor shrink repaints the viewport | delete one newline or wrap row | local differential update |
| Viewport snaps backward | origin before/after shrink | preserve the viewport contract |
| A frame replays transcript history | short vs large Sessions | bounded frame work |
| Stale or duplicate row | apply bytes to the terminal simulator | expected final screen |
| Scrollback cleared | detect `CSI 2J`, `CSI 3J` | neither on an ordinary edit |
| Cursor homes | detect `CSI H` | relative movement |
| Wrapped input enters history early | Up within wrapped input | expected cursor and history state |
| Resize corrupts layout | width/height change with ANSI, wide text, images | expected final screen |
| Hidden-row change can't repair locally | change content above viewport | repair the visible screen without losing scrollback |

## Evidence
- Define the terminal scenario before renderer code. Compare escaped Pi and PiG output and simulator state; keep inputs, dimensions, bytes, screen state, mutation failure.
- Isolated tests for width, row, cursor, viewport math. Profile editor shrink, streaming, overlays, large transcripts.
- `go test -race ./tui` + owning parity scenarios (`interactive-rendering`, `tui-components`, `fullscreen`).
- Don't fix a rendering bug in the renderer when the cause is upstream-source logic elsewhere.
