### Fixed

- Bind deferred `nextTurn` custom messages to their session leaf and discard stale branch context on navigation or session replacement. Extensions receive `queued_message_discarded` with `customType` and `reason`; `PIG_DEBUG` also records each discard (SPEC-0007 D1). <!-- fix: ee8403488,cb7bf2cff -->
