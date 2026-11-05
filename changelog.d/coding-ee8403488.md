### Fixed

- Bind deferred `nextTurn` custom messages to their session leaf and discard stale branch context on navigation or session replacement. Extensions receive `queued_message_discarded` with `customType` and `reason`; `PIG_DEBUG` also records each discard (SPEC-0007 D1). <!-- fix: ee8403488,cb7bf2cff -->

### Added

- `queued_message_discarded` reaches extensions when a deferred `nextTurn` message is dropped because its session leaf is no longer current, carrying `customType` and `reason` so the sender can re-queue it (SPEC-0007 D1).
