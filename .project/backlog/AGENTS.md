Open work items. Type `backlog-item`, id `B-NNNN`, file `B-NNNN-<slug>.md`; next free number, never reused.

- Frontmatter adds `kind` (bug|drift|debt|feature|chore), `priority` (P1|P2|P3), `scope` (path globs). Status: open|done|wontfix.
- Body: current behavior (file:line), target behavior, why it matters, fix hint. Evidence → `B-NNNN-<slug>/evidence/` if needed.
- Every open item has exactly one pointer line in the `BACKLOG.md` of the deepest directory covering its `scope` (repo root if cross-cutting): `- B-NNNN <kind> <prio> — <title> → .project/backlog/B-NNNN-<slug>.md`.
- Create: item file + pointer line in one commit. Fix: delete pointer line + `status: "done"` + changelog entry in the fixing commit.
- Claims/in-progress live in project-scope `tasks`, never in these files.
