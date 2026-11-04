---
document:
  type: "specification"
  id: "SPEC-0002"
  title: "System prompt composition: Jinja layouts, section catalog, extension sections"
  description: "Replaces last-writer-wins prompt mutation with a section catalog, declarative extension sections and one Jinja layout render per run whose output is recorded exactly as sent"
  status: "draft"
  revision: 1
  created_at: "2026-11-04T14:30:00+01:00"
  updated_at: "2026-11-04T14:30:00+01:00"
  created_by: "orchestrator"
  updated_by: "orchestrator"
  changelog:
    - revision: 1
      reason: "Owner asked for a Jinja system prompt with sections that extensions can contribute to and take over"
      changes: ["document creation"]
---

# SPEC-0002 System prompt composition

Resolves U-05, U-06, U-07, U-15 (prompt part) and the context-file suppression item of `../analyses/2026-11-04-host-gap-verification/README.md`.

## Context

Today a run's prompt is built in three steps:

1. Base sections from options: `internal/codingagent/prompts/coding.go` `BuildSystemPromptSections` (`preamble`, `tools`, `rules`, `docs`, `addendum`, `project_context`, `skills`, `cwd`).
2. `before_agent_start` handlers, in load order, edit `systemPromptOptions`, add `sections`, or force a whole prompt.
3. `coding/session_transcript.go` `effectiveSystemSections` applies the custom sections; the result is diffed against the transcript and sent. A forced prompt is projected onto the request and not recorded.

Problems:
- Last writer wins. A handler sees only what earlier handlers did; built-in extensions (MCP) run last, so a forcing extension loses `mcp_servers` (U-05).
- The transcript can differ from the request (U-06).
- No control over the preamble, deletion or order (U-07).

## Decisions

### D1 Sections are the output

- Composition produces ordered named sections (`ai.OrderedSections`), never one opaque string.
- The transcript records exactly what is sent. The section diff and the mid-conversation patch path (`internal/codingagent/prompts/sections.go` `DiffSystemPromptSections`) stay unchanged.
- Invariant, tested: for every request, the system text rendered from the transcript equals the system text of the provider payload before `before_provider_request` handlers run.

### D2 Pipeline per run

| Phase | What happens | Who |
|---|---|---|
| 1 Catalog | Built-in producers fill the section catalog: name, body, owner, structured data | host |
| 2 Handlers | `before_agent_start` handlers run unchanged; their `sections` enter the catalog with the extension as owner; MCP's `mcp_servers` arrives here | extensions, built-ins |
| 3 Contributions | Declarative sections (D6) render and apply add/replace/remove to the catalog | extensions |
| 4 Layout | The selected layout template (D4) renders the catalog into ordered sections | host |
| 5 Record | Diff against the transcript, record, send | host |

- Phases 3 and 4 run after every handler, built-ins included. Handler load order no longer decides what the layout sees.
- Seam: phase 3–4 run inside `effectiveSystemSections` (one call into a new `internal/codingagent/prompts/compose` package). The rest is new files.
- Composition runs once per run start and again only when inputs change mid-run (tool loadout change). Never per request.

### D3 Engine

- `github.com/nikolalohinski/gonja/v2` (pure Go Jinja2, MIT; v2.9.1 at the time of writing). Version pinned in `go.mod`; `nix/vendor-hash.sh update`.
- Autoescape off. Templates are parsed once per load and cached by content hash.
- Template loader: built-in layouts (embedded) and the directory of the selected template file only. `include`/`import`/`extends` resolve inside that set. No other file, environment or process access from templates.

### D4 Layout templates

Syntax added on top of Jinja:

| Construct | Meaning |
|---|---|
| `{% section "name" %}…{% endsection %}` | emit section `name` with the rendered body |
| `{% section "name", wrap=false %}` | emit without the `<name>…</name>` tag (default for `preamble`, `wrap=true` otherwise) |
| `{{ section("name") }}` | emit the catalog section `name` as it is (body + owner); nothing when absent |
| `{{ rest() }}` | emit every catalog section not emitted yet, in catalog order |
| `has("name")` | test whether the catalog has `name` |

Rules:
- Non-whitespace text outside a section is a render error.
- Emitting a name twice is a render error.
- Section names match `^[a-z][a-z0-9_-]*$`.
- A section whose rendered body is empty is omitted.
- A catalog section the layout neither emits nor passes to `rest()` is dropped. That is how a layout deletes a section.

Variables (read-only data; values are never parsed as templates, so `{{` inside `AGENTS.md` stays text):

| Variable | Content |
|---|---|
| `sections` | ordered catalog: `name`, `body`, `owner` (`builtin` or extension name) |
| `tools` | active tools: `name`, `snippet`, `guidelines`, `hidden`, `source` |
| `skills` | `name`, `description`, `location`, `disable_model_invocation` |
| `context_files` | `path`, `content` |
| `model` | `provider`, `id`, `api`, `reasoning`, `context_window`, `max_output`, `mid_conversation_system` |
| `thinking`, `mode`, `cwd`, `os`, `date` | thinking level; `tui`/`print`/`json`/`rpc`; working directory; OS; local date (day only, so the prompt stays cache-stable within a day) |
| `custom_prompt` | `SYSTEM.md` content or empty |
| `append.file`, `append.flags` | `APPEND_SYSTEM.md` content; `--append-system-prompt` values (the built-in default keeps today's rule: flags replace the file) |
| `vars` | owner-defined values from `systemPrompt.vars` |

The built-in layout `builtin:default` reproduces today's output byte for byte. Existing prompt tests (`internal/codingagent/prompts/*_test.go`, `coding/session_prompt_sections_test.go`, `coding/session_system_prompt_updates_upstream_test.go`) run against it unchanged.

Example owner layout:

```jinja
{% section "preamble" %}You are my coding agent. Be terse.{% endsection %}
{{ section("rules") }}
{% section "environment" %}OS: {{ os }} · model: {{ model.provider }}/{{ model.id }} · {{ date }}{% endsection %}
{{ section("tools") }}
{{ section("project_context") }}
{{ rest() }}
```

### D5 Selecting the layout

Precedence, first match wins:

1. CLI `--system-layout <path|builtin:name>`.
2. Project `.pig/SYSTEM.md.j2`, only in a trusted project (added to the trust-gated files in `internal/codingagent/trust_manager.go`).
3. Setting `systemPrompt.layout` (path relative to the settings file, or `builtin:<name>`), project scope over global.
4. Global `<agentDir>/SYSTEM.md.j2`.
5. The layout an extension registered (D6); two extensions registering layouts without an owner choice is a load error naming both.
6. `builtin:default`.

`SYSTEM.md` and `APPEND_SYSTEM.md` keep their discovery and become template inputs (`custom_prompt`, `append`).

### D6 Extension API

Registered at factory time, applied per run in phase 3. Same API in every SDK (Go, Rust, Python, TypeScript declarations) and over the wire.

| Call | Effect |
|---|---|
| `RegisterPromptSection{Name, Mode, Before/After, Wrap, Render}` | `Mode`: `add` (new name; error when it exists), `replace` (take over an existing name, `preamble` included), `remove`. `Render(PromptContext) (string, error)` gets the variables of D4; empty result omits the section this run |
| `RegisterPromptLayout{Name, Source}` | offer a layout template; used per D5 |

- One owner per name for `replace`/`remove`. A second claim fails registration with an error naming the current owner, unless the owner setting (below) names the claimant.
- Placement: `Before`/`After` an existing name; none → after the last built-in section. Unknown anchor → appended plus a diagnostic.
- `Render` errors fail the run before any request, with the extension named. No silent fallback.

Owner settings (`settings.json`, project over global):

```json
{
  "systemPrompt": {
    "layout": "builtin:default",
    "vars": { "persona": "terse" },
    "sections": {
      "docs": { "enabled": false },
      "rules": { "owner": "workbench" }
    }
  }
}
```

`enabled: false` removes a section whatever extensions do. `owner` resolves a takeover conflict.

### D7 Forced prompts

- `before_agent_start` returning `systemPrompt` (or setting `forceSystemPrompt`) still forces the whole prompt. Phases 3–4 are skipped and a diagnostic names the forcing extension.
- A forced prompt is recorded: the transcript gets an opaque system message (content, no sections) whenever the sent prompt differs from the recorded one. Compaction entries store it. U-06 is closed for both paths.
- Sections set by handlers after a force are dropped with a diagnostic, never recorded.

### D8 Order changes mid-session

The section patch carries values and deletions (`null`), not order. An order change takes effect at the next full render: every request for models without mid-conversation system messages, the next session or compaction for models with them. Documented for layout authors.

### D9 Inspection

- `pig prompt render [--json]` prints the composed sections for the current cwd, settings and extensions: name, owner, size, body.
- `/prompt` in the TUI shows the same for the live session.
- Template parse errors report file:line at load and reload.

### D10 Performance

- No configured layout and no contributions → phase 3–4 cost is one catalog copy; benchmark against the current path, allocation counts reported.
- Layout render target: under 1 ms for the default layout with 50 skills and 10 context files (benchmark `BenchmarkComposeDefault`).

## Consequences

- Workbench drops its `force` delivery and the payload splice in `before_provider_request`; it registers sections and, optionally, a layout.
- Extensions written against `before_agent_start` keep working.
- "System prompt layout" is distinct from Pi's prompt templates (`/name` slash-command templates); docs use that term only.

## Docs

- New `docs/site/docs/system-prompt.md` and `internal/pigdocs/content/system-prompt.md`: catalog, layout syntax, variables, selection, owner settings, order rule (D8), inspection; add to `docs/site/docs/docs.json` and the bundle index.
- `settings.md` (both trees): `systemPrompt.*`; `make generate-schemas`.
- `extensions.md` / `extension-api.md`: section and layout API, forced-prompt recording; SDK docs for Rust, Python, TS declarations.
- `trust.md`: `.pig/SYSTEM.md.j2` is trust-gated.
- `prompt-templates.md`: one line pointing to `system-prompt.md` to separate the two features.
- `cli.md`: `--system-layout`, `pig prompt render`; `slash-commands.md`: `/prompt`.
