# T027 A1 pure admission Prerequisite Validation Report

[中文](t027-pure-admission.md)

Date: 2026-09-14. Conclusion: `PREREQUISITE COMPLETE / T027 BLOCKED / NOT IMPLEMENTED / AT-23 NOT PASS`.

## Validation scope

This report covers only the T027 offline slice A1 pure admission and does not claim a production AgentRunner implementation. A1 delivered:

- `AgentRunnerCompositeAdmission` dropped its `RequestIndex` field and its in-memory replay consumption: admission is now pure preflight and immutable binding with no replay-write side effects and no launch.
- The one-shot replay blocking sentinel `ErrAgentRunnerRequestReplayBlocked` plus the in-memory `AgentRunnerRequestIndex` type, constructor, and `Record` method were removed because that admission gate was their only consumer. Same-key idempotent replay, same-key conflicting digests, and first-dispatch semantics remain owned by the A0 durable replay store on the dispatch path; `ErrAgentRunnerRequestConflict` is retained and still exercised by replay-store tests.
- The responsibility boundary is frozen in text: `docs/CONTRACTS.md` / `CONTRACTS_EN.md` now state that production admission is a pure preflight and immutable binding that must not consume request replay identities, publish R/C or write to the replay store, or launch processes, and that first-dispatch eligibility plus replay identity consumption belong to the dispatch path, decided after R is published, based on the replay store; the admission verdict is a point-in-time conclusion, not a lock, so after publishing R and before allowing launch the dispatch path must reconfirm that the authorization is not revoked and the budget reservation is still outstanding; the `AgentRunnerAdmission` interface docs in `internal/chain/router.go` and the struct/method docs in `agent_runner_admission.go` declare the same boundary.
- Focused tests: `TestAgentRunnerCompositeAdmissionIsStatelessPreflight` (three consecutive passes on one admission; a same-requestId different-digest conflict is decided by the dispatch/store boundary, not admission); `...FailureLeavesNoConsumedState` (after an invalid record and a missing clock block admission, repairing the condition lets the same instance pass); `...ConcurrentRepeatsPass` (8 concurrent repeated admissions all pass); `...LeavesInputsUnchanged` (deep-equality assertions over every admission input on both the passing and blocked paths, promoting "no write side effects" from a review conclusion into a regression assertion); `...RejectsNilReceiverCancelledContextAndZeroClock` (the nil-receiver, cancelled-context, and zero-time-clock fail-closed branches); the memory-index-only tests were deleted and stale `BeforeReplay` naming was cleaned up.

## Execution results

### Native Windows execution

| Command | Result |
|---|---|
| `gofmt -l internal` | Clean (no output) |
| `go build ./...` | Pass |
| `go vet ./...` | Pass |
| `go test -count=1 ./...` | All pass |
| `go test -count=1 -run 'AgentRunnerCompositeAdmission' ./internal/adapters` | Pass |

Note: this Windows machine has no gcc (cgo unavailable), so `-race` was not run locally; the Linux CI `-race` step has actually executed and passed in run 34809403487 (below). Concurrent correctness rests on the regular execution of the new concurrency case plus the internal `RWMutex` read lock of `CostLedger.RequireOutstandingReservation` and the fact that admission is read-only throughout.

### CI workflow change (frozen contract synced)

- New Linux CI step `Race`: `CGO_ENABLED=1 go test -race -count=1 ./internal/adapters/... ./internal/chain/...`, executed only when the runner is Linux.
- The frozen workflow contract in `internal/release` was updated in step: the `test` job step list gains `Race`, with its script digest and condition written into the frozen table; the local `go test -count=1 ./...` suite (including the release workflow contract and mutation-rejection suites) passes.
- After the push, CI run [34809403487](https://github.com/larsonzh/prfrail/actions/runs/34809403487) succeeded: `Go ubuntu-latest` passed every step (Build/Vet/Test/Race/Contract fixtures) and `Go windows-latest` passed; `Race` was the first real execution of the race-regression gate.

### Review conclusion

- V4 Pro pre-analysis (run retroactively): A1 itself PASS (the implemented design matches the ideal rulings item by item); it raised R1 (a contract amendment gating A2) plus R2–R4 follow-ups, all remediated: R1 freezes in the bilingual CONTRACTS that the admission verdict is not a lock and that dispatch must, after publishing R and before launch, reconfirm a non-revoked authorization and an outstanding budget reservation, failing closed otherwise; R2 adds the three fail-closed branch tests; R3 documents the conflict sentinel ownership; R4 documents the invariance helper criteria; the Codex remediation re-review passed.
- V4 Pro pre-review (run retroactively): PASS (no Critical/High/Medium). All four Low suggestions adopted: L1 a Linux CI `-race` gate (with the frozen workflow contract synced); L2 the input-invariance regression test; L3 the A2 handoff note (admission is a preflight, not a lock; dispatch must decide based on the replay store after publishing R); L4 contract wording tightened to "publish R/C or write to the replay store".

- Codex independent final review: PASS (no Critical/High/Medium). Every Low suggestion was adopted: stale `BeforeReplay` naming cleaned up, plus the early-failure no-side-effect and concurrent admission counterexamples added; the no-launch observation is covered by the existing router test proving admission blocking prevents dispatch. Follow-up edits were test naming and added cases only with no production change, so no mandatory re-review was triggered.

## Known boundaries

- Admission's "no launch" property is structural: `AgentRunnerCompositeAdmission` holds no port, process, or replay dependency, and the router layer already tests that a blocked admission never calls `AgentRunnerPort`.
- Replay identity consumption and first-dispatch eligibility are not yet implemented; they belong to the A2 replay-aware dispatcher. Until then no production dispatch path exists, so this slice's removal of the admission block creates no real redispatch exposure.
- This Windows machine has no cgo/gcc, so `-race` did not run locally; the Linux CI `-race` step has been proven by GitHub Actions (run 34809403487).
- Admission is a preflight, not a lock: a verdict can expire before dispatch (grant revocation, budget settlement); A2 dispatch must decide based on the replay store after publishing R.

## Explicitly not performed

- No real dispatch was introduced; `AgentRunnerPort`/Engine wiring is untouched;
- No real model CLI was invoked, no network access, no SecretStore reads;
- No commit, push, or publish happened while this report was written; both followed under same-turn user authorization (`452d173`), with CI evidence above;
- T027 remains `BLOCKED / NOT IMPLEMENTED` and AT-23 has not passed.
