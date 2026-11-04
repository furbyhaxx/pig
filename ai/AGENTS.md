# AGENTS.md — ai

Port of Pi `packages/ai`: Transcript and message types, Providers (one file family per API), model catalogs, auth storage, Event Stream.

## Layout
- Providers: `anthropic*.go`, `openai.go` (openai-completions), `openai_responses.go` (openai-responses), `azure*.go`, `bedrock*.go`, `google*.go`, `cloudflare*.go`, ... Upstream: `.upstream/current/packages/ai/src/api/`.
- `*_generated.go` — model/provider catalogs generated from the exact published Pi package pinned by `extensions/sdk-ts`. Pin change → `make model-catalogs`. Never hand-edit.
- `event_stream*.go`, `events.go` — Event Stream (`.upstream/current/packages/ai/src/utils/event-stream.ts`). `sse.go` — the only SSE parser.
- `faux*.go`, `test_faux.go` — hermetic faux provider for tests and parity.
- Request contract checks: `internal/requestvalidity` (rules cited in `wire.go`).

## Contract
- A Provider receives a normalized Transcript. Public Message, content, Usage and event shapes match pinned Pi.
- Event Stream: preserves event order, exactly one terminal result. Setup may return an immediate error; provider, auth, runtime and cancellation failures after creation terminate the stream.
- Providers read the transcript via `TranscriptContext.borrowed` and never modify the slice, a message or anything it references. `Messages` is the copying accessor for callers that keep or change the result.

## Failure modes
| Failure | Detection | Required result |
|---|---|---|
| Invalid JSON enters trusted provider code | non-JSON value or non-finite number | reject at the public boundary |
| Result depends on event consumption | request Result without iterating past the former queue bound | Result completes |
| Iteration stops while waiting | cancel the iterator context | no waiter or goroutine remains |
| Cancellation races termination | concurrent Push, iteration, Result under `-race` | one terminal result |
| Late coroutine driver continues another turn's drain | `TestContinuationLateDriverDoesNotContinueAnotherTurnsDrain` | driver continues only its own body's drain; no body runs beside a reaction; executor never left draining without a loop |
| Invalid provider event sequence | setup, start, block, delta, end, terminal paths | Pi event order |
| Replay metadata crosses provider/model identity | replay signed content through another model | drop or transform as Pi does |
| Consumed events retained | drain a long stream | consumed queue entries released |
| Transport consumer stalls | stop reading subprocess frames | bounded backpressure + cancellation |
| Provider gets a Transcript its API rejects (orphan tool result, content between call and results, duplicate call id, refused role order) | every test; any binary built with `-tags requestvalidity` | `internal/requestvalidity` panics at the send seam. `PIG_REQUEST_VALIDITY_REPORT=<file>` records instead. Intentionally invalid test requests use `invalidRequest(t)` and state why |
| Request build copies the whole transcript | convert 100- vs 800-turn histories per API | linear cost per turn; `TestRequestBuildLeavesTheTranscriptUnchanged`, `TestRequestBuildDoesNotCopyWhatItDoesNotSend`, `TestRequestBuildAllocationsAreLinearInTheHistory` |

## Evidence
- Write failure cases before provider code. Isolated tests for normalization, JSON, payload conversion, event order, cancellation; a real Model Runtime path for `complete`, `stream`, `streamSimple`. Cover every provider family.
- Assert the serialized request shape (table-driven), not just a successful return.
- `go test -race ./ai`. Mutation-check one event discriminator and one provider content field.
- `test/pipeline-diff` sends seeded random histories through Pi's and PiG's convertToLlm → transformMessages → API converters and compares wire requests; contract rules come from provider docs and Pi's `transform-messages.ts`.

## Traps
- Don't parse SSE outside `sse.go`; don't classify errors by status-code substrings; don't map unknown/missing stop reasons to success (divergence-guard).
- Don't special-case one provider in shared code; see root "Faithful, general implementations".
