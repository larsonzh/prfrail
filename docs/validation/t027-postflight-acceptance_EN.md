# T027 · A5 · chain postflight acceptance — Validation Report

Date: 2026-09-15. Verdict: `A5 complete (all local gates green, awaiting commit authorization)`.
Slice goal: move a task's review qualification away from adapter records and into a chain-owned postflight gate — only "rebuildable frozen facts plus an independently decided pass" may write `REVIEW_PENDING`.

## 1. Scope

- NEW `internal/chain/postflight.go`: the chain-owned `AgentRunnerFrozenFacts` DTO (five digests: manifest/diff/log/usage/process-stop, mutually distinct and distinct from the request and completion digests), `PostflightOutcome` (passed/failed/uncertain), `PostflightRequest`, `PostflightDecision.Validate(facts)`, the required `PostflightPort`, the evidence-layout constants (`postflightFactsPrefix=2`, `postflightFactsCount=5`, `postflightFactsMinEvidence=7`), `agentRunnerFactsEvidence`, `frozenFactsFromEvidence`, `agentRunnerTaskFacts`, `runTaskPostflight`, `requirePostflightQualifiedReview`, and the sentinels `ErrInvalidPostflightFacts` / `ErrInvalidPostflightDecision` / `ErrPostflightUnavailable` / `ErrUnqualifiedReview`.
- NEW `internal/chain/postflight_test.go` (fact-validation matrix, evidence rebuild, decision validation, gate pass/reject/uncertain/port-failure/invalid-decision, partial-write convergence, re-entry convergence after `REVIEW_PENDING` is durable, multi-step refusal, damaged-evidence replay conflict, the three review-qualification shapes plus the re-entry refusal, and per-task attribution).
- `internal/chain/engine.go`: `Options.Postflight` is required (nil → `ErrPostflightUnavailable`); `runTask` runs the gate after every step has passed and before `Acceptance.Accept`, and only when the task is not already in `REVIEW_PENDING`; re-entry into `REVIEW_PENDING` now requires a qualification proof.
- `internal/chain/agent_runner_terminal.go` / `agent_runner_terminal_route.go`: the terminal DTO carries `Facts` (completed must, non-completed must not), the routing evidence became `[RequestHash, CompletionHash, five fact digests, ...Evidence]` (deduped, order preserved), and a completed replay compares the facts.
- adapters: `agent_runner_terminal_record.go` (intent gains a `frozenFacts` block plus validation), `agent_runner_terminal_publisher.go` (a completed outcome must carry facts before the intent is written), `agent_runner_terminal_chain.go` (one-way projection including the `toChain()` fact rebuild), `agent_runner_replay_store_terminals.go` (comment explaining why the settlement-plan comparison excludes facts), and their tests.
- `internal/console/runtime.go`: `chain.Options.Postflight = runtime`, with `localRuntime.RunPostflight` as a fail-closed stub (returns `ErrPostflightUnavailable`, so a noop-only run never triggers it).
- `internal/evidence/event.go`: the task transition table gains `STEPS_RUNNING→REPAIR_PENDING`.

## 2. Protocol-first artifacts

- `docs/CONTRACTS.md` / `_EN.md` §2.1: the task row gains `STEPS_RUNNING→REPAIR_PENDING`, noted as written only by a postflight rejection.
- New §7 "postflight acceptance (A5)" paragraph in the same section: the frozen-facts DTO and "completed must carry / non-completed must not", five mutually distinct digests distinct from the request and completion digests, positional rebuild from the evidence prefix, the port's obligation to re-derive independently (re-capture the manifest, diff it against the parent, scope/secret/file-type/side-effect checks) and reconcile field by field, the failure-routing table, idempotency and the crash windows, the **legality of a no-change execution** (the chain never refuses a fact set merely because a digest equals the parent snapshot digest; that discrimination belongs to the port), the **pre-A5 routing-evidence compatibility boundary**, and the **`REVIEW_PENDING` re-entry qualification**.
- `schemas/state-event.schema.json`: `taskTransition` gains the `REPAIR_PENDING` target.
- Contract fixtures: valid `state-event-task-postflight-rejected-valid` added; invalid `state-event-task-repair-before-steps-invalid` added (`PRECHECK→REPAIR_PENDING` must be refused, with the reason text updated to "reachable only from REVIEW_PENDING, FAILED, or a postflight rejection"); `tools/contracts/contracts.test.js` fixture count 126→128.

## 3. Fact layout and failure routing

The frozen facts appear only on completed terminals, at fixed positions:

$$\text{routeEvidence} = [\text{RequestHash},\ \text{CompletionHash},\ \text{Manifest},\ \text{Diff},\ \text{Log},\ \text{Usage},\ \text{ProcessStop},\ \ldots\text{Evidence}]$$

The gate rebuilds the facts from positions 2..6 (recoverable across restarts, without reading the replay store); the minimum is seven digests, and anything missing or ambiguous fails closed.

| Trigger | task | chain | reason | evidence |
|---|---|---|---|---|
| port returns failed | `REPAIR_PENDING` | `RUNNING→PAUSED` | `postflight-rejected` | facts + port evidence |
| port returns uncertain | `FAILED` | `PAUSED` only when RUNNING | `postflight-uncertain` / `recovery-uncertain` | port error evidence; `Run` refuses to retry |
| port returns an error | `FAILED` | `FAILED` | `postflight-failed` | — |
| invalid decision (missing facts / unknown outcome / malformed digest) | `FAILED` | `FAILED` | `postflight-invalid` | — |
| missing facts, multi-step attribution ambiguity, evidence too short | `FAILED` | `FAILED` | `postflight-facts-missing` | — |
| port passes | `REVIEW_PENDING` (evidence carries the five facts) | unchanged | `postflight-passed` | facts + port evidence |
| task without AgentRunner facts | pre-existing path unchanged | unchanged | `steps-completed` | — |
| `REVIEW_PENDING` re-entry whose qualification cannot be proven | unchanged (zero writes) | unchanged | — | returns `ErrUnqualifiedReview` |

## 4. Core invariants and focused tests

- The gate is the only qualification for review: `TestPostflightGatePassesAndBindsTheFrozenFacts` (the `REVIEW_PENDING` transition evidence carries the five fact digests, and the per-task attribution sequence is `[task-one, task-two, task-three]`); `TestPostflightGateRefusesAnUnprovenPass` (a passed decision missing any fact never reaches acceptance).
- The port is required: `TestEngineRequiresAPostflightPort`.
- Facts are provable or fail closed: `TestAgentRunnerFrozenFactsValidateFailsClosed`, `TestFrozenFactsRebuildFromRoutingEvidence` (short and ambiguous evidence refused), `TestCompletedTerminalWithoutFrozenFactsIsRefused`, `TestAgentRunnerTerminalFactsMustNotCollideWithItsBindings`.
- Attribution: `TestPostflightGateRefusesTasksWithSeveralRoutedSteps` (completed terminals from two steps for one task are refused rather than guessed, the port is never called, and the refusing record is itself a valid event chain).
- Idempotency and crash windows: `TestPostflightGateConvergesAfterReviewPendingWasPersisted` (re-entry after a durable `REVIEW_PENDING` converges and the gate does not re-run), `TestPostflightGateConvergesAfterPartialWrite` (crash after a passed gate but before the `REVIEW_PENDING` write → re-entry re-runs the gate once and converges).
- Replay safety: `TestAgentRunnerReplayWithDifferentFactsConflicts` (same request+completion with different facts → conflict with zero writes; identical facts still converge), `TestAgentRunnerReplayWithDamagedRoutingEvidenceConflicts` (recorded evidence too short to rebuild the facts → conflict).
- Compatibility and qualification boundaries: `TestRequirePostflightQualifiedReviewFailsClosed` (no review transition / transition without facts / transition with facts), `TestPostflightGateRefusesAnUnqualifiedReviewState` (a pre-A5-style `REVIEW_PENDING` re-entry → `ErrUnqualifiedReview`, zero writes, zero downstream calls, no gate re-run).
- The port's discrimination boundary: `TestPostflightGateHandsAParentDigestFactSetToThePort` (a fact set whose manifest digest equals the parent snapshot digest is not refused by the chain beforehand; the port's verdict decides the outcome).
- Adapter-side fact binding: `TestAgentRunnerTerminalIntentFactsBindingFailsClosed` (completed without facts / non-completed with facts / non-distinct facts → all refused), `TestAgentRunnerTerminalRecordValidation` (the completed-without-facts counterexample goes through the public constructor), `TestAgentRunnerTerminalPublisherRefusesACompletedOutcomeWithoutFacts` (refused before any intent and before settlement).

## 5. Review chain and remediation (closed loop: any remediated gate is re-run until it passes)

- **③ V4 Pro pre-review (round 1)**: 1 High (`runTask` wrote `REVIEW_PENDING` unconditionally, so a crash after that write made re-entry attempt a self-transition that the event validation rejects — a permanent failure) + 1 Medium (a completed replay compared only request+completion, so divergent facts converged silently) + 3 Lows (unvalidated `ErrorEvidence` digests, facts colliding with the parent digest, stale fixture reason text and count). Remediation: the `current != "REVIEW_PENDING"` guard, fact comparison on completed replays, `ErrorEvidence` digest validation, fixture text and count fixes, plus regression tests such as `TestPostflightGateConvergesAfterReviewPendingWasPersisted` and `TestAgentRunnerReplayWithDifferentFactsConflicts`.
- **③ V4 Pro re-review (round 2)**: High/Medium cleared; four Lows remained (the multi-step refusal lacked a test and a contract sentence, the `sameTerminalSettlementPlan` comment did not explain the facts exclusion, a working-doc count was stale, and damaged evidence lacked an end-to-end test). Remediation: added `TestPostflightGateRefusesTasksWithSeveralRoutedSteps` + the contract sentence, the comment, the count fix, and `TestAgentRunnerReplayWithDamagedRoutingEvidenceConflicts`.
- **④ Codex final review (round 1)**: 1 High + 2 Mediums. (1) High: `facts.DisjointFrom(parent.Hash)` was a self-added, unratified rule — `parent.Hash` is the accepted snapshot's manifest digest, and under the ratified port contract ("re-capture the manifest and diff it against the parent") a **no-change execution legitimately equals it**, so the pre-port refusal could suppress a legitimate pass → the check was withdrawn and that discrimination assigned to the port (bilingual CONTRACTS updated), with the earlier ③ Low thereby becoming a port obligation. (2) Medium: pre-A5 completed routing evidence carries no fact prefix, so replay conflicts with zero writes and cannot be resumed in place → written into the contract as a compatibility boundary. (3) Medium: `TestPostflightDecisionRejectsMalformedEvidence`'s first case was masked by the shared "missing fifth fact" cause and was not falsifying for the digest loop → changed to "all five facts present plus one extra malformed digest". Adapter-side fact counterexamples were added as well.
- **④ Codex re-review (round 2)**: a new Medium — a task already resting in `REVIEW_PENDING` skipped the gate and went straight to acceptance, so a pre-A5 review state (or a review state without a gate trace in a damaged store) could reach `PASSED` without any postflight, contradicting the new contract sentence → `requirePostflightQualifiedReview` and the `ErrUnqualifiedReview` sentinel were added: when a task with routed completed terminals is re-entered in `REVIEW_PENDING`, its review transition must prove all five fact digests are present, otherwise the run is refused with zero writes; tasks without facts keep the pre-existing path.
- **④ Codex re-review (round 3)**: `RE-REVIEW: PASS`, no Medium+ findings; it confirmed no new states or engine entry points, that adapters still cannot write policy or state, that the gate is not re-run once review is durable, and that no reference to the withdrawn parent-digest rule or to a renamed test remains.

## 6. Mutation checks (proving the new tests are falsifying)

| Mutation | Expectation | Result |
|---|---|---|
| Disable the multi-step refusal branch in `agentRunnerTaskFacts` | `TestPostflightGateRefusesTasksWithSeveralRoutedSteps` reddens | ✅ red (`facts routed by several steps must be refused, got <nil>`) |
| Remove the `REVIEW_PENDING` guard in `runTask` | `TestPostflightGateConvergesAfterReviewPendingWasPersisted` reddens | ✅ red (the self-transition is rejected by event validation) |
| Remove the `Evidence` digest loop in `PostflightDecision.Validate` | `TestPostflightDecisionRejectsMalformedEvidence` reddens | ✅ red (`got <nil>`) |
| Restore the `facts.DisjointFrom(parent.Hash)` refusal | `TestPostflightGateHandsAParentDigestFactSetToThePort` reddens | ✅ red (the port is never called) |
| Remove the fact-containment loop in `requirePostflightQualifiedReview` | `TestRequirePostflightQualifiedReviewFailsClosed/without_the_facts` and `TestPostflightGateRefusesAnUnqualifiedReviewState` redden | ✅ both red |

After each mutation the source was restored and verified by SHA-256 (the backup lived in `tmp/` and was deleted after validation).

## 7. Accepted boundaries and residual risk

- **Port implementation quality is the only trust boundary**: the chain can prove structure (positional rebuild, five mutually distinct digests, distinctness from the request/completion digests, a review state that carries its qualification) but cannot prove that the port's re-capture is real, nor tell a legitimate no-change execution from an adapter echoing the parent snapshot using digest relations alone; by contract that discrimination belongs to the port.
- The engine's existing gates (build/test/verify hook steps) still run before `Acceptance.Accept`, which differs in wording from CONTRACTS §7's "gates after freeze"; this slice does not reorder them and leaves reconciliation to the B4 final review.
- `sameTerminalSettlementPlan` deliberately does not compare facts: the persisted intent is the only recovery source for the settlement plan, and divergent facts are rejected as a replay conflict at the chain routing layer — noted in the implementation.
- Compatibility boundaries: pre-A5 completed routing records carry no fact prefix, so replay conflicts with zero writes; a pre-A5 `REVIEW_PENDING` belonging to a task with facts is refused with `ErrUnqualifiedReview` and zero writes. Both classes of run must be **re-executed as a new attempt**, never resumed in place.
- `Accept` idempotency between `REVIEW_PENDING` and `Accept` remains a pre-existing assumption; this slice does not change it.
- The review-qualification guard proves digest containment over the loaded events, so it relies on the event log staying append-only and untampered after the `New`/`Recover` validation.

## 8. Gate evidence and reproduction

- `gofmt -l internal` prints nothing; `go build ./...`, `go vet ./...` and `go test -count=1 ./...` all pass (Windows native, `GOTOOLCHAIN=local`, toolchain 1.22.12).
- Contract-fixture gate: `cd tools/contracts; npm test` → 4 tests pass, fixture count 128.
- Byte-level encoding/EOL verification of every changed file: `.md` is UTF-8 **with BOM** + LF; `.go`/`.json`/`.js` are UTF-8 **without BOM** + LF.
- Reproduction: `go test -count=1 ./internal/chain/... ./internal/adapters/... ./internal/console/... ./internal/evidence/...`; the contract gate is above.
- This offline slice used **0 paid probes** (budget 0/0), called no real CLI and accessed no network.
