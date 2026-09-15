# T027 · A4 · chain terminal routing and resume — Validation Report

Date: 2026-09-15. Verdict: `A4 implementation and local gates PASS`; commit/push and CI evidence are written back only after explicit same-turn user authorization.
Slice goal: converge terminal semantics into chain core routing and recovery, so an adapter receipt can never drive a task conclusion or skip the downstream acceptance flow.

## 1. Scope delivered

- New chain-owned terminal fact `AgentRunnerTerminal` (five statuses: completed/failed/cancelled/operator-action-required/uncertain) binding RequestID/RequestHash/CompletionHash/Run/Task/Step/Attempt/SessionID/Prior* plus deduplicated ordered evidence hashes; no settlement digests, no wire details.
- New step state `TERMINAL_PENDING`: `RUNNING→TERMINAL_PENDING` (evidence = the dispatch result hashes) means "dispatched, awaiting terminal"; `TERMINAL_PENDING→WAITING_FOR_OPERATOR|PASSED|FAILED|CANCELLED`. It is not a schedulable resting state.
- New single external write entry `Engine.SubmitAgentRunnerTerminal`; runStep for AgentRunner steps now dispatches, parks as `TERMINAL_PENDING`, and returns `ErrAwaitingAgentRunnerTerminal` instead of taking the step-passed path.
- New one-way adapters translation `ToChainAgentRunnerTerminal` / `LoadChainAgentRunnerTerminal` (adapters→chain); chain reads no replay store.
- Files added/changed: `internal/chain/agent_runner_terminal.go`, `internal/chain/agent_runner_terminal_route.go`, `internal/chain/agent_runner_terminal_test.go`, `internal/adapters/agent_runner_terminal_chain.go`, `internal/adapters/agent_runner_terminal_chain_test.go`; `internal/chain/engine.go`, `internal/chain/operator_interaction.go`, `internal/chain/recover.go`, `internal/evidence/event.go`, `internal/chain/engine_test.go`, `schemas/state-event.schema.json`, `testdata/contracts/{valid,invalid}/runtime-records.json`, `tools/contracts/contracts.test.js`, `docs/CONTRACTS{,_EN}.md`, `docs/t027/REMAINING_SLICES{,_EN}.md`.

## 2. Protocol-first artifacts

- `docs/CONTRACTS.md` / `_EN.md` §2.1: the step row becomes `RUNNING→TERMINAL_PENDING|WAITING_FOR_OPERATOR|PASSED|FAILED|CANCELLED` with a new `TERMINAL_PENDING→WAITING_FOR_OPERATOR|PASSED|FAILED|CANCELLED`, noting that `TERMINAL_PENDING` is not a resting state and that a restart without terminal evidence is an uncertain pause.
- New §7 paragraph "Chain terminal routing (A4)" in the same chapter: the route table, the operator route pausing the chain, the additive open-side acceptance plus the ledger-exclusivity boundary, the uncertain-pause boundary, idempotency/conflict rules, and the resume-continuity criterion with the U-9 boundary.
- `schemas/state-event.schema.json`: `stepTransition` gains the `TERMINAL_PENDING` branch and the extended `RUNNING` target set.
- Contract fixtures: new valid `state-event-step-terminal-pending-valid` (RUNNING→TERMINAL_PENDING), valid `state-event-step-terminal-pending-passed-valid` (TERMINAL_PENDING→PASSED), and invalid `state-event-step-terminal-pending-resume-invalid` (TERMINAL_PENDING→RUNNING must be rejected); `tools/contracts/contracts.test.js` fixture count 123→126, aligned line-by-line with the Go transition table in `internal/evidence/event.go`.

## 3. Route table and core invariants

| Terminal | step | task | chain | Notes |
|---|---|---|---|---|
| completed | `TERMINAL_PENDING→PASSED` | unchanged (`STEPS_RUNNING`) | unchanged (`RUNNING`) | the task still walks freeze/gates/review/completed promotion on the next `Run` |
| failed | `TERMINAL_PENDING→FAILED` | `FAILED` | `FAILED` | no auto retry/relaunch |
| uncertain | `TERMINAL_PENDING→FAILED` | `FAILED` | `RUNNING→PAUSED` (reason `recovery-uncertain`) | stays paused and never retries (`Run` returns `ErrRecoveryUncertain`); an already-paused chain keeps its pause reason |
| cancelled | `TERMINAL_PENDING→CANCELLED` | `CANCELLED` | `CANCELLED` | archives stop evidence first (`Stopper.Stop`); on failure the step stays parked and the chain pauses |
| operator-action-required | `TERMINAL_PENDING→WAITING_FOR_OPERATOR` | `STEPS_RUNNING→WAITING_FOR_OPERATOR` | `RUNNING→PAUSED` | routing itself pauses the run and hands over to the T025 machinery; evidence carries the C digest |

Invariant evidence (focused tests):

- No direct pass: `TestAgentRunnerCompletedTerminalAloneDoesNotCompleteTask` (after Submit the task is still `STEPS_RUNNING` with zero downstream calls; the next `Run` completes it); the routing matrix asserts zero downstream calls at routing time for all five statuses.
- Dispatch only parks: `TestAgentRunnerDispatchParksStepInsteadOfPassingIt` (re-entry neither redispatches nor writes events).
- Binding and conflict: `TestAgentRunnerTerminalRefusesUnboundAndConflictingReceipts` (unparked step, non-AgentRunner step, RequestHash mismatch, and same request with a different completion are all refused with zero writes; a repeated terminal converges with `Converged=true` and zero writes).
- Crash convergence: `TestAgentRunnerTerminalConvergesAfterPartialWrite` (step written, task not) and `TestAgentRunnerTerminalConvergesAfterTaskWriteBeforeChainWrite` (step+task written, chain not; the replay writes exactly the one missing transition).
- Resume continuity: `TestAgentRunnerTerminalResumeContinuityProof` (`PriorCompletionHash` outside the step's event evidence history is refused with zero writes; once proven it routes; a terminal replayed after the operator answered never pushes the task back into waiting).
- Operator handover: `TestAgentRunnerOperatorRouteSurvivesRecoveryAndInteractionOpen` (route→restart→`Recover`→`OpenOperatorInteraction`→`ResumeOperatorInteraction`→`Run` closes the loop with `RecoveryUncertain` still false).
- Stop evidence: `TestAgentRunnerTerminalCancelledArchivesStopEvidence` (the CANCELLED transitions carry the stop digest and the C digest) and `TestAgentRunnerTerminalCancelledArchiveFailuresStayParked` (a stop failure keeps the step parked, pauses the chain, and returns `ErrRecoveryUncertain`).
- Restart uncertainty: `TestAgentRunnerRecoverTreatsParkedStepAsUncertain` (`TERMINAL_PENDING` joins the recover uncertainty list; a proven terminal can still be routed later, but the run never resumes silently).
- Translation layer: `TestToChainAgentRunnerTerminalMapsThePublishedChain` (identity/status/session/evidence mapping), `TestToChainAgentRunnerTerminalFailsClosed` (zero values, intent binding, settlement-slot mismatch, resume session inequality), `TestLoadChainAgentRunnerTerminalRequiresASettledClosure` (unsettled, unclosed, and missing-store cases are refused; the translation is strictly read-only — a full path→size snapshot of the store is identical before and after).

## 4. Review chain (protocol steps ③④, closed loop)

1. **V4 Pro pre-review (③)**: `PRE-REVIEW: FINDINGS` — 1 Medium + 2 Low.
   - Medium: the operator route wrote only the task/step waiting states and left the chain `RUNNING`; a crash before the interaction record existed made `Recover` write a `recovery-uncertain` pause that neither `OpenOperatorInteraction` nor `Run` could leave, so the route could never converge.
   - Low 1: the open side does not verify the interaction-record hash for an agent waiting state (the control ledger owns exclusivity) — document it; Low 2: an uncertain terminal while the chain is already paused never reaches `PAUSED` — document the boundary.
2. **Remediation**: the route now pauses the chain (`RUNNING→PAUSED`, evidence carrying the C digest); new additive helper `operatorWaitingPauseMatches` (classic binding first); both boundaries recorded in the bilingual contract.
3. **V4 Pro re-review**: `RE-REVIEW: FINDINGS` — no Medium+, but a new Low that is a real gap: replaying an operator terminal after the operator machine had returned control pushed the task/chain back into waiting.
4. **Remediation**: `convergeAgentTerminalRoute` treats a route whose step is no longer `WAITING_FOR_OPERATOR` as already converged (zero writes); the superseded-replay assertion was added to `TestAgentRunnerTerminalResumeContinuityProof`.
5. **Codex independent final review (④)**: `FINAL REVIEW: FINDINGS` — 1 Medium + 1 Low + 3 coverage suggestions.
   - Medium: the translation bound only closure `requestId`/`intentRecordHash`/`completionHash`, not the settlement slot (`settlementEntryId`/`settlementIdempotencyKey`), so a tampered closure with a recomputed recordHash could pass the translation.
   - Low: it reported four `.md` files missing a BOM (see §5 boundary 5 — byte-level measurement shows this is a false positive).
   - Coverage: add a "task written, chain not written" convergence counterexample, and finish the route→restart→Recover→Open→Resume test with a final `Run`.
6. **Remediation**: the translation now cross-binds the settlement slot (same rule as `AgentRunnerReplayStore.validateTerminalClosureBinding`) with two genuine counterexamples built from internally valid records (recomputed recordHash); new `TestAgentRunnerTerminalConvergesAfterTaskWriteBeforeChainWrite`; the operator recovery test ends with a final `Run`.
7. **Codex re-review (④ closure)**: `RE-REVIEW: PASS`, no findings; it confirmed the Medium is closed, the two settlement counterexamples are falsifying, the crash-window test is genuine, and no new problem was introduced. On R4 (BOM) Codex stated its tooling could not obtain byte-level evidence and declined to judge the claim unverified (recorded as-is).

Conclusion: every Medium+ finding was fixed and re-reviewed by the independent reviewer (Codex); nothing was self-certified.

## 5. Known boundaries (deliberate, recorded in the contract or explained here)

1. **First attempt only**: the chain models attempt 1 (`taskEntity`/`stepEntity` are fixed), so any terminal claiming another attempt is refused; a new attempt is built by the CLI from the durable envelope (U-9) and is owned by A6/later slices.
2. **Open-side binding strength**: the waiting transition written by the agent route binds the terminal receipt (no interaction record existed yet), so Open does not verify the interaction-record hash; interaction exclusivity remains the control ledger's (`OperatorInteractionLedger`) responsibility. Recorded in CONTRACTS §7.
3. **Uncertain with an existing pause**: the transition table forbids `PAUSED→PAUSED`, so an uncertain terminal only pauses a chain that is `RUNNING`; when the chain is already paused for another reason, resuming stops at the task `FAILED` boundary (`runTask` fast-fail) and never retries the step. Recorded in CONTRACTS §7.
4. **Run during a wait**: calling `Run` while an interaction is pending first resumes the chain and then re-pauses it with `repair-required`, which no longer matches the interaction binding — the same behaviour the classic T025 path already had; callers must not run during a wait.
5. **BOM false-positive clarification**: the Codex Low claimed four `.md` files lack a BOM. Byte-level measurement in the repository shows `EF BB BF` as the first three bytes (via `[System.IO.File]::ReadAllBytes`) with no CRLF, matching `docs/CODING_CONVENTIONS.md`; the Low is judged a false positive caused by a text-level read that stripped the BOM, and no change was made.
6. **Publisher/chain decoupling**: the A3 publisher does not depend on chain; A4 only hands an already published and closed terminal fact to chain through the translation, so a receipt itself has zero state side effects.

## 6. Gate evidence (local, Windows + Go 1.22 toolchain)

| Command | Result |
|---|---|
| `gofmt -l internal` | empty output |
| `go build ./...` | pass |
| `go vet ./...` | pass |
| `go test -count=1 ./...` | every package passes (including `internal/chain`, `internal/adapters`, and the frozen `internal/release` workflow contract test) |
| `npm test` (`tools/contracts`) | 4/4 pass (126 contract fixtures including the 3 new ones) |
| byte checks (`.md`) | `docs/CONTRACTS{,_EN}.md`, `docs/t027/REMAINING_SLICES{,_EN}.md`, this report: BOM=True, CRLF=False |

Linux CI evidence (with the `Race` step and the `Contract fixtures` step) is written back after an authorized commit/push.

## 7. How to reproduce

```text
go build ./... && go vet ./... && go test -count=1 ./...
# focused: go test -count=1 -run 'AgentRunnerTerminal|ChainAgentRunner' ./internal/chain/... ./internal/adapters/...
cd tools/contracts && npm test
```

## 8. Not done / follow-ups

- A5 postflight acceptance integration; A6 pinned CLI adapter plus the dispatcher/admission resume lift (U-7 keeps that out of A4); the B-series real-candidate and durability proofs.
- T027 remains `BLOCKED / NOT IMPLEMENTED`; AT-23 has not passed.
