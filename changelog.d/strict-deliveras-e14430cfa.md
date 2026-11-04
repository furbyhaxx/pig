### Fixed

- `sendMessage` and `sendUserMessage` reject an unknown `deliverAs` with an error on the call instead of silently taking the steering path or ignoring the value. An omitted `deliverAs` keeps the call's default; `sendUserMessage` also rejects `"nextTurn"`, which only `sendMessage` defines (SPEC-0007 D3). <!-- fix: e14430cfa -->

### Changed

- The Go, Rust and Python SDKs export typed delivery modes (`DeliverSteer`, `DeliverFollowUp`, `DeliverNextTurn`, and `DELIVER_STEER`, `DELIVER_FOLLOW_UP`, `DELIVER_NEXT_TURN` in Rust and Python) beside the existing string options.

### Upgrading

- An extension that passed a misspelt or unsupported `deliverAs` used to queue the message as a steering message. It now fails the call and the message is not delivered; use one of the exported constants.
