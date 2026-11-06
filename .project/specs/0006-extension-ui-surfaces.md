---
document:
  type: "specification"
  id: "SPEC-0006"
  title: "Extension UI surfaces: view events for widgets, header and footer; nested tool-call renderers"
  description: "Routes kit view events to every extension view, not only overlays, and lets extension tool renderers draw the child calls that codemode and other tools make"
  status: "published"
  revision: 2
  created_at: "2026-11-04T14:48:00+01:00"
  updated_at: "2026-11-06T11:37:12+01:00"
  created_by: "orchestrator"
  updated_by: "coding-agent"
  changelog:
    - revision: 1
      reason: "Resolve U-10 and U-11 of the host-gap verification"
      changes: ["document creation"]
    - revision: 2
      reason: "Owner accepted the decisions as the fork contract"
      changes: ["status draft -> published"]
---

# SPEC-0006 Extension UI surfaces

Resolves U-10, U-11 of `../analyses/2026-11-04-host-gap-verification/README.md`.

## D1 View events for every view (U-11)

- The host delivers `ui.view.event` for widget, header and footer kit views, not only overlays (`extensions/sdk/view.go:311` routes to overlays only today).
- Delivered: mouse events inside the view's `MouseRegion`s (press, release, wheel, hover when the terminal reports motion). Keyboard focus stays with the editor; non-overlay views receive no key events.
- SDK: kit views of every kind accept the same event handler as overlays; one routing table keyed by view id.
- Mouse capture follows the existing fullscreen/mouse setting; without mouse reporting the views receive nothing and must offer a keyboard path.

## D2 Nested tool-call renderers (U-10)

- A tool that runs child calls (`ctx.ExecuteTool`, codemode) renders each child row through the tool renderer registry.
- The resolver gets the child's tool name, arguments, status, result/details and `nested: true` with the parent call id. A claiming renderer draws the row inside the parent's block at the parent's inner width; otherwise the parent's own row format is used.
- The codemode renderer keeps its summary header and ordering; only child rows are delegated.

## Docs

- `tui.md` / `extensions.md`: view events per view kind, focus rule, keyboard fallback.
- `codemode.md` and tool-renderer section of `extensions.md`: nested rendering, `nested` flag, width.
