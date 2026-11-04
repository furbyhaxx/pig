### Changed

- Bound `Host.Reload` and CLI extension reload admission and registration with the global `extensionsReloadHandshakeMs` setting (default 15000 milliseconds), excluding builds. Failed cells retain their serving generations and declarations; successful cells can still publish while the caller receives the failure (SPEC-0004 D2).

- `DefaultResourceLoader` reloads now run through one persistent Extension Host, so they inherit the reload budget, per-cell rollback and partial-success error. Successful factories and file extensions publish while a failed sibling's serving generation stays loaded; a caller-supplied `Runtime` keeps its ownership, and the aggregate runtime is a fresh snapshot per reload (SPEC-0004 D2).

- Extension starts now use generation-unique sockets so an older generation's teardown cannot remove its successor's socket (SPEC-0004 D1).

- Pi oracle source packing now rejects Node versions other than the repository pin with an actionable error, before npm reports a frozen archive integrity mismatch.

- Clarify the personal fork policy: Pi and upstream PiG compatibility is retained while it does not limit correctness, performance, design quality or the owner's requirements; breaking changes require a specification and migration notes (SPEC-0001 D7).

### Upgrading

- Handle CLI reload errors from `ctx.Reload()` and RPC, and SDK reload errors from `DefaultResourceLoader.Reload`/`ReloadWith`. A failed replacement now remains loaded instead of being silently dropped; other successful cells may already be active (SPEC-0004 D2).

- An SDK reload returns a fresh `ExtensionRuntime` view unless you supplied `DefaultResourceLoaderOptions.Runtime`. Compare registration state with `SharesRegistrations` instead of comparing runtime pointers across reloads.
