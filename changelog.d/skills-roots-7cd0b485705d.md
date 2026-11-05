### Added

- `skills.roots` chooses the directories searched for skills, in precedence order. The object form accepts `roots` and `paths`; tokens `@project`, `@project-agents`, `@agent` and `@user-agents` keep the default locations addressable, and leaving out `@user-agents` isolates a setup from `~/.agents/skills`. Absolute, `~/` and settings-relative paths work, duplicate real roots are searched once, and a project's roots apply only after you trust it. `pig status --json` and `pig config` use the same discovery as startup.

### Changed

- `skills` accepts an object `{roots, paths}` as well as the existing path list. `customDirectories` migrates to `paths`; a malformed root or path entry is ignored with a diagnostic instead of discarding the whole object.
