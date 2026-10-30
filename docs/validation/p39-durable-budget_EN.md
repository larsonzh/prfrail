# ProofRail P39 report: C6 durable budget (directive v1.60)

Slice `P39-T027B-DURABLE-BUDGET` (implementation type `[IMPL]`); run on 2026-10-31; directive v1.60; R = 4, D = 0.

## 1 Executive summary
- **Landed**: a durable budget ledger with a cross-process shared lock; `Reserve`, `Settle` and `RequireOutstandingReservation` now decide only after re-reading under the lock.
- **Fixed scope**: C6 (the budget was never persisted) - the shared cap cannot be bypassed across processes, a failed commit issues no acknowledgement, and an `unknown` charge still counts after restart.
- **Verification**: `go build ./...` and `go vet ./...` pass; `go test ./...` is green; all six counter-example legs pass.
- **Core judgement**: design, then red tests, then green implementation, then the first pre-review FAIL closing five BLOCKERs, then a PASS re-review, then a CLEAR scan and a PASS final review.

## 2 Preflight and implementation
### 2.1 Pre-review (two rounds)
- **First round FAIL**: five BLOCKERs - a read could return an unconfirmed tail, the three commit states were indistinguishable, a read path created the lock file, a missing ledger re-issued the initial allocation, and several matrix legs never really ran.
- **Fix and re-review**: the implementer fixed every item and added a staged fault-injection surface plus a cancellable lock wait; the re-review returns PASS with four non-blocking notes.
### 2.2 Implementation and tests
- **Tests first**: the tester wrote red tests against undefined interfaces, then the implementer turned them green.
- **Matrix coverage**: two-process cap race, lock identity and aliasing, history-drift refusal, read-cannot-mint, staged commit faults, and missing-ledger refusal to re-issue.
- **Platform notes**: the `TestDurableHelperReserve` and symlink skips are a subprocess guard and a platform capability limit, each registered.

## 3 Independent scan (6)
```sixstage
## Section A Scope and method
- Scope: the implementation files and the test file; method: read-only static scan against the contract.

## Section B Findings
| # | SEVERITY | Issue | Location |
|---|---|---|---|

## Section C Coverage and gaps
- Coverage: lock ordering, atomic commit, error classification, test credibility, contract consistency.
- Gaps: none.

## Section D Difference from the previous stage
- Consistent with the pre-review re-review; no new disagreement.

## Section E Recommendation
- Raise the Windows directory-barrier grade in the durability contract slice.

INDEPENDENT SCAN: PASS
```

## 4 Final review (7) and open items
- **Final verdict**: the final review returns PASS - C6 is judged resolved, with no deadlock, starvation, double write, forged acknowledgement or platform semantic error.
- **Open item**: the Windows directory barrier stays `unproven` per contract; whether to raise it is deferred to the `CONTRACTS-S7-DURABILITY` slice.
- **Observation candidates**: the four non-blocking pre-review notes are registered here for `CONTRACTS-S7-DURABILITY` or a later refactor - one, no test leg for lock release after a holder crash; two, no test leg for the visible branch of a `rename` that reports failure; three, the `confirmed` field naming in `newDurableAck` mismatches its semantics; four, the pre-confirmation cancellation leg is not landed.

## 5 Process declarations
- **Carriers and models**: the architect used `deep-reasoner` with `GPT-6.1 Sol (copilot)` under user authorisation; the implementer, tester and documenter used `DeepSeek V4.1 Flash (deepseek)`; the pre-review used `DeepSeek V4 Pro (deepseek)`; the scan used `independent-reviewer` with `MAI-Code-1.1-Flash (copilot)`; the final review used `independent-reviewer` with `GPT-5.3-Codex (copilot)`.
- **Deviation record**: the first implementer dispatch wrongly used `DeepSeek V4 Pro (deepseek)`; all four artifacts were voided with hashes kept (section 7).

## 6 Artifact registration
The change set holds the implementation files, one test file, the evidence pack, the ledger pair and the v28 checklist.

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/t027/CHECKLIST_v28.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/evidence/P39-T027B-DURABLE-BUDGET/01-arch-design/review-batch-b.md	repo
docs/validation/evidence/P39-T027B-DURABLE-BUDGET/01-arch-design/review.md	repo
docs/validation/evidence/P39-T027B-DURABLE-BUDGET/MANIFEST.md	repo
docs/validation/evidence/P39-T027B-DURABLE-BUDGET/SHA256SUMS.txt	repo
docs/validation/evidence/P39-T027B-DURABLE-BUDGET/freedom-list.md	repo
docs/validation/p39-durable-budget.md	repo
docs/validation/p39-durable-budget_EN.md	repo
internal/tickets/durable_cost_fault.go	repo
internal/tickets/durable_cost_ledger.go	repo
internal/tickets/durable_cost_ledger_test.go	repo
internal/tickets/durable_cost_lock_shared.go	repo
internal/tickets/durable_cost_lock_unix.go	repo
internal/tickets/durable_cost_lock_windows.go	repo
docs/validation/evidence/P39-T027B-DURABLE-BUDGET/02-implement/voided-v4pro-waste.md	repo
```

## 7 Exception handling / fact pack
1. **Model binding deviation**: the first implementer dispatch wrongly used `DeepSeek V4 Pro (deepseek)`, breaching the role-model binding; the stop-loss was to halt immediately, void all four artifacts with sha256 recorded, and redo tests-first with `DeepSeek V4.1 Flash (deepseek)`.
2. **First pre-review FAIL**: a normal in-process loop; the five BLOCKERs are closed and this is not counted as a deviation.
3. **Platform limit**: Windows has no `chmod` semantics, so the read-only directory leg was replaced by a fault-injection leg.
4. **Nothing else**: no out-of-scope call, no unregistered temporary artifact, no contract change.
