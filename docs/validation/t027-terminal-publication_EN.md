# T027 A3 terminal publication and settlement prerequisite validation report

[中文](t027-terminal-publication.md)

Date: 2026-09-14. Verdict: `PREREQUISITE COMPLETE / T027 BLOCKED / NOT IMPLEMENTED / AT-23 NOT PASS`.

## Validation scope

This report covers only the T027 offline slice A3 terminal publication and settlement and makes no claim about the AgentRunner production implementation. A3 delivered:

- Contract freeze (before implementation): `docs/CONTRACTS.md` / `CONTRACTS_EN.md` gained the terminal-publication paragraph — the four-phase chain (terminal-intent → ledger `Settle` → publish C → terminal-closure), the two-file store-local protocol (never on the wire, never part of R/C recordHash), ledger deduplication by `idempotencyKey`, settlement evidence carrying the judged C and R digests, pre-settlement rejection of any new publication on unproven-durability platforms with write-free replay of complete chains only, orphan C fail-closed without adoption, no C when the reservation is already settled under a different key, unknown settlements keeping the hold, recovery from store-local records only (the ledger is an enforcer), and takeover remaining with a later operator slice.
- `AgentRunnerTerminalPublisher` (adapters, value semantics): the four-phase `PublishTerminal` plus `convergeTerminal` crash recovery; preflight reuses `ValidateAgentRunnerCompletionBinding` (including the resume session binding, eliminating the "settled then rejected by C" deadlock); the settlement plan is derived deterministically from (R, C, usage) (idempotency key, evidence composition and dedupe, unknown normalization); it depends on `*tickets.CostLedger` plus test-only single-shot hooks (afterIntentWrite/afterSettle).
- Store extensions: `terminals/terminal-intent.<requestId>.jsonl` and `terminals/terminal-closure.<requestId>.jsonl` records, accessors, and binding validation; write paths follow the existing discipline (read-only replay resolution → durability gate → binding checks → no-replace with bounded reread convergence); `terminals/` joins bootstrap, the production constructor, the locked runtime path-safety re-check, and the ownership-publish preflight enumeration.
- Idempotency key: `settle-<20hex>` = digest over domain `proofrail:agent-runner-settlement-key:1` of (requestId, C.recordHash); `EntryID` reuses the value and `CompletionID` is derived; "a repeated receipt never settles twice" is enforced by the ledger's `settlementByKey`/`settlementByReservation`.
- Focused tests: ordering coupling (hooks assert that the intent is visible before settlement and C before its settlement, with both hooks asserted to fire), idempotent replay, 8-way concurrency single winner (including different-clock barrier concurrency and deterministic sequential convergence), crash windows i/ii/iii recovery (exactly one settlement, closure completed), orphan fail-closed, foreign-settlement rejection, unknown retention (`UnknownHoldReservations()==1`, `OutstandingReservations()==0`, summary 1000), the five-case settlement matrix, charged-over-reserved, different-outcome-same-request conflict, unproven pre-settlement rejection (fresh and converge variants), write-free replay of a complete chain on unproven, closure key-mismatch and closure-without-completion corruption, recovery requiring the persisted request, resume positive/negative cases, evidence prevalidation/dedupe/negative rejection, 24-way cross-store single winner, and Unix construction-time symlink rejection, runtime swap, dangling-slot convergence, sync failure without rollback, and the proven full chain.

## Execution results

### Windows native run (native evidence)

| Command | Result |
|---|---|
| `gofmt -l internal` | clean (no output) |
| `go build ./...` | pass |
| `go vet ./...` | pass |
| `go test -count=1 ./...` | all pass |
| `go test -count=1 -run 'AgentRunnerTerminal\|AgentRunnerReplayStoreTerminal' ./internal/adapters` | pass |

Notes: the A3 concurrency counterexamples (8-way, different-clock, 24-way cross-store) ran natively on Windows; Windows has no gcc so `-race` was not run locally. The Unix-specific tests (symlinks, swap, convergence, sync failure) are carried by GitHub Actions Ubuntu, with CI evidence written back after the push.

### Review outcomes

- V4 Pro pre-analysis (mandatory for this slice): the four-phase protocol, the crash-recovery algorithm, the rejection matrix, and all Q1–Q10 rulings were adopted (including "the intent embeds the full C", "no new ledger query path", and "the durability precheck moves before Settle").
- V4 Pro implementation pre-review: conditional pass (no High). H1 (a missing resume session binding could settle and then fail C publication) was remediated through the shared binding validator; M1 (different-clock concurrency falsely reporting a conflict) was remediated with semantic slot matching plus clock-skew tests; M2 (caller evidence could only be rejected by the ledger after the intent was durable) was remediated with pre-write evidence validation, deterministic dedupe, and exact membership matching; M3 (no converge+unproven test), L1 (closure without completion corruption), and L5 (recovery requires the persisted request) were all remediated.
- Codex independent final review (including a per-test counterexample audit): conditional pass (no High). Both Medium findings were remediated: unproven complete-chain replay is now strictly write-free (no ledger call), and converge reuses the full closure binding check (including settlement keys) while `RecordTerminalClosure` gained an intent runId re-check. The three Lows were remediated too: the different-clock concurrency test gained a start barrier plus a deterministic sequential case, and the happy path asserts both hooks actually fire. Three missing counterexamples were added: write-free replay, closure key mismatch corruption, and recovery requiring the persisted request record.
- Codex remediation re-review (2026-09-15, reviewing the remediated code; per the directive that medium-or-higher findings must be resubmitted for re-review and must never be self-certified): **PASS — no Critical/High/Medium remain and the slice is signable from the review standpoint**. The re-review also recorded three low-risk test-rigor gaps (the write-free counterexample's probe branch is unreachable, the key-mismatch case does not isolate its variable because C is not published first, and the `RecordTerminalClosure` intent runId re-check lacks a dedicated regression test), all below the medium threshold and therefore not blocking sign-off under the directive.
- Tail-item remediation and extra confirmation (2026-09-15, authorized by the user in the same round): all three low-risk test-rigor gaps were fixed — the publisher gained an unexported test-only `beforeSettle` probe at the single settlement choke point `settleTerminalIntent` (never set on production paths), the write-free counterexample now asserts under both fatal probes (`beforeSettle`/`afterSettle`), the closure key-mismatch case publishes C first to isolate its variable, and a dedicated `RecordTerminalClosure` intent-runId regression test was added. Codex extra confirmation (beyond the 1+1 budget, explicitly authorized): **PASS — all three are revert-sensitive genuine counterexamples, no new defects, no Critical/High/Medium remain, and the slice is signable**. Mechanical confirmation on this host: `gofmt`/`go build`/`go vet`/`go test -count=1 ./...` all green.

## Known boundaries

- Plan irreparability: the terminal-intent is a no-replace single winner; correcting a failed plan for the same request belongs to the operator or a later slice (fail-closed by design; the slot is never rewritten).
- Recovery semantics: on unproven platforms only a complete chain may be replayed write-free; incomplete chains (including settled-but-unpublished C) are refused and must be recovered on a proven instance; recovery is an idempotent completion of a slot we already own, not takeover (W1/W2 unchanged).
- Unknown settlement is terminal; this slice defines no upgrade semantics (deferred to A5).
- The ledger exposes no by-reservation query path: recovery relies on store-local records only, with the ledger acting as an idempotency/conflict enforcer.
- Items guaranteed by A6 port implementations (for example non-zero `StartedAt`) are unchanged.
- Terminal publication success says nothing about task PASS; task-state advancement belongs to A4.

## Explicitly not executed

- No Engine wiring, no real candidate run, no chain terminal routing or postflight acceptance (A4/A5);
- no real model CLI invoked, no network access, no SecretStore reads;
- at the time of writing this report no commit, push, or publish was performed; CI evidence will be written back after same-turn authorization under the discipline;
- T027 remains `BLOCKED / NOT IMPLEMENTED` and AT-23 has not passed.
