# T027-B End-to-End Wiring (P69 `T027B-E2E-WIRING`) Validation Report

Date: 2026-11-12. Status: ⑧b final (after the ⑨ native validation went all green). Slice: `T027B-E2E-WIRING` (P69 · layer-2 governance side · implementation class [SLICE] · M tier).

## §1 Executive summary

- Goal: wire the P44 readiness / permit decision into the production assembly path and the CLI entry, closing R-1 of P57 §3.3.
- Result: `DESIGN: DONE` (①); the wiring lands 5 production files under `internal/console/**` (4 modified, 1 added) plus 5 test files (4 modified, 1 added); all five acceptance counterexamples carry measured proof (③ landed 7 tests).
- Review chain: ⑤ `PRE-REVIEW: PASS` (1 Low fixed at comment level, ④ rerun in place all green); ⑥ `INDEPENDENT SCAN: PASS` (0 findings); ⑦ `FINAL REVIEW: PASS` (0 findings; metering = 1 final review + 0 re-review + 0 blind review).
- Why ⑦ did not sample a blind review: this slice touches no ownership / stop / identity semantics (§7.7 mandatory sampling does not apply); the proportional sampling window starts at v1.32, and by the §12.4 tail block (counting only formally closed slices) P67 = 1st, P68 = 2nd, P69 = 3rd ⇒ fewer than 4 slices, so no sample.

## §2 Wiring points

Content from design-r1 §2.1 (`docs/validation/evidence/P69-T027B-E2E-WIRING/01-arch/design-r1.md`); the added/deleted counts are `git diff --cached --numstat` readings (the kickoff 6 + 6 file figures disagree with the change set, so the measured 5 + 5 is recorded).

| # | Path | Add/Modify/Delete | Content | Impact |
|---|---|---|---|---|
| 1 | `internal/console/readiness_gate.go` | added | `requireProbePermit` triple (Evaluate(GoalProbe) → Require → release/Complete) as the single source | new file (+29) |
| 2 | `internal/console/cli.go` | modified | `--probe-permit` flag + pre-routing presence decision (decide only, never consume) + `runUsage` + permit forwarding at two call sites + gates import | +32 / −4 |
| 3 | `internal/console/runtime_agent.go` | modified | options gain the `ProbePermit` field; Step 0 decision (after stub, before Step 1); Require (after all of Step 1, before the Step 2 marker) + `defer release` | +23 / −0 |
| 4 | `internal/console/runtime_supervised.go` | modified | options gain the same field; fresh gate (`SupervisedInteractionsLedgerPath` absent = first start only, then decide+Require+defer; re-entry safe) | +23 / −0 |
| 5 | `internal/console/supervised_cmd.go` | modified | `executeSupervisedRun` signature gains the permit and assigns `options.ProbePermit`; `ExecuteSupervisedReturn` untouched | +4 / −1 |
| 6 | `internal/console/readiness_gate_test.go` | added | 1 positive + 3 counterexample groups, 7 tests in total | new file (+225) |
| 7 | 4 existing test files | modified | injection points only, zero assertion changes (fixture `newAssemblyFixture`, harness `newSupervisedHarness`, 10 bare literals in `runtime_multitask_red_test.go`, `p51SevenAgentRunnerArgs` gains the flag) | +32 / −16 |

## §3 Call order and failure paths

- AgentRunner: ctx → Step 0 decision (missing permit ⇒ `ErrProbePermitMissing`, zero traces) → Step 1 preconditions → `requireProbePermit` (Require) → `defer release` → Step 2 marker → Steps 3–8.
- Supervised: ctx → plan → RunRoot/VisibleSender validation → fresh gate (first start: decide+Require+defer; re-entry: skip) → marker → engine.
- Three failure shapes: ① decision failure ⇒ CLI/entry refuses, zero consumption and zero traces; ② permit already consumed ⇒ the CLI decision still passes (presence-only is P44's deliberate design), the entry Require refuses with `ErrProbePermitConsumed`, and no double start happens (O_EXCL cross-process backstop); ③ decision passes but the permit stays unconsumed ⇒ only the Step 1 refusal between decision and Require — fail-closed, no spawn, quota intact; spawn ⇒ Require already succeeded ⇒ the deferred release always Completes (no consumption gap).
- Relation to existing mechanics: Require precedes the event-log marker (a refused assembly writes nothing into runRoot); it is unrelated to the inter-task stop proof barrier / R publication (the gate sits only at the assembly entry).

## §4 Counterexamples and mutation checks

7 tests (③ measured, `internal/console/readiness_gate_test.go`):

| # | Test | One-line assertion |
|---|---|---|
| 1 | `TestProbePermitConsumedAfterSuccessfulAssembly` | positive: the marker lands on disk, a disk reload Require yields `ErrProbePermitConsumed`, and a second run on the same object yields Consumed rather than InUse |
| 2 | `TestAgentRunnerAssemblyRefusesNilProbePermit` | nil permit ⇒ `ErrProbePermitMissing`, no run-root marker |
| 3 | `TestSupervisedAssemblyRefusesNilProbePermit` | supervised first start without a permit ⇒ same sentinel, no run-root marker |
| 4 | `TestAgentRunnerAssemblyRefusesSpentProbePermit` | spent permit ⇒ `ErrProbePermitConsumed`, no marker, zero spawn |
| 5 | `TestSupervisedAssemblyRefusesSpentProbePermit` | supervised twin: Consumed, no marker |
| 6 | `TestRunRefusesMissingProbePermitFlag` | missing `--probe-permit` ⇒ non-zero exit and stderr names the flag |
| 7 | `TestRunRefusesIncompleteRecordsBeforeBurningPermit` | missing `--agent-cli` ⇒ the noop-only wording is kept and no permit marker appears |

Mutation checks (M1–M5, from `docs/validation/evidence/P69-T027B-E2E-WIRING/06-gates/mutation-evidence.txt`):

| Mutation | Object | Focused test | Red | Restore sha256 equal | Green after restore |
|---|---|---|---|---|---|
| M1 | `runtime_agent.go`: the Step 0 decision block plus the `requireProbePermit`/defer block (removed together) | `TestAgentRunnerAssemblyRefusesNilProbePermit` | FAIL | equal (1100E2DE…C3F336) | PASS |
| M2 | `runtime_agent.go`: only the `requireProbePermit` plus defer block | `TestAgentRunnerAssemblyRefusesSpentProbePermit` | FAIL | equal (1100E2DE…C3F336; first attempt hit a file lock, the retry succeeded) | PASS |
| M3 | `cli.go`: the whole CLI presence gate block (kept compilable with `_ = *probePermit`) | `TestRunRefusesMissingProbePermitFlag` | FAIL | equal (5C8D6B15…620EE2B) | PASS |
| M4 | `cli.go`: the `!inputs.complete()` refusal block | `TestRunRefusesIncompleteRecordsBeforeBurningPermit` | FAIL | equal (5C8D6B15…620EE2B) | PASS |
| M5 | `runtime_agent.go`: `defer releaseProbePermit()` → `_ = releaseProbePermit` | `TestProbePermitConsumedAfterSuccessfulAssembly` | FAIL | equal (1100E2DE…C3F336) | PASS |

Counterexample 4 (the existing P44 gates decision tests stay unchanged): `internal/gates/**` zero diff + `go test -count=1 ./internal/gates/...` ok.

## §5 Review chain verdicts

- ① `DESIGN: DONE` (design-r1); batch assistance GLM-5.3-Flash twice (5546 / 5288 characters).
- ⑤ r1 `PRE-REVIEW: PASS` (1 Low: the release closure swallowed a Complete error, making it unobservable — the minimal fix argues at comment level that the failure is unreachable, and ④ reruns the full suite green afterwards; 7 Info recorded). ⑤ quota 1/2 (no re-review used).
- ⑥ r1 `INDEPENDENT SCAN: PASS` (Section B has no findings; Section E residuals = live two-process contention evidence, live platform durability path). Quota 1/3.
- ⑦ r1 `FINAL REVIEW: PASS` (0 findings; all six mandatory items pass; residual risk same as ⑥ Section E). Metering = 1 final review + 0 re-review + 0 blind review.
- Five-element trace (§7.8.5): object = ⑦ round N; round = r1; basis = §4.3 + §5.3 normal convergence; verbatim conclusion line = `FINAL REVIEW: PASS`.

## §6 Gate readings (⑨ final measurement, 2026-11-12, p9 batch; the intermediate ④ batch readings live in this slice's temp artefacts and are cleaned per §5.1 at ⑧c)

- `gofmt -l .` empty; `go build ./...` exit 0; `go vet ./...` exit 0; `go test -count=1 ./...` 15 packages ok / FAIL=0.
- `GOOS=linux go build ./...` and `GOOS=linux go vet ./...` produce empty output (local evidence that Unix is unaffected; the CI ubuntu leg is the backstop).
- `node tools/gates/gate.js --selftest` ⇒ `SELFTEST: PASS 737/737`.
- Three scopes: `--scope=index` / `--scope=tree` both `TOTAL_FAIL=0` (`changed=38`; G4b `newCjk=405 unknown=[]`; G8-a/b OK; G5(f) index = worktree); `--scope=ci` locally degrades to the `HEAD^1..HEAD` range (changed=2, `TOTAL_FAIL=0`); this slice's change-set coverage comes from index/tree, and the real ci reading is produced by CI after the push (existing repo convention).
- Archiving: the gate-reading captures were not archived (they can be regenerated by re-running the gates); the mutation evidence is archived.

## §7 Known limits

- The noop path and the `ai check` path are not wired in this slice (the permit lifecycle of `ai check` needs its own design and belongs to a later slice).
- The positive test is platform aware: on Windows `agentRunnerPublicationDurability` is unproven and the production dispatcher refuses before R ⇒ the local assertion lands on "the quota is consumed exactly once + a second run is Consumed"; the strict success branch on durability-proven platforms runs on the CI ubuntu leg.
- A live two-process contention over the same marker is not covered (⑥ Section E and the ⑦ residual risk point at the same gap; O_EXCL linearization is carried by the P44 decision itself and the existing gates tests, and this slice does not change the decision body).
- `ExecuteSupervisedReturn` is not gated (return/resume is not a new start; the quota was already spent at the first start).
- The permit injection shape deviates from the wording at P44 report line 156 (F-1) and is recorded as is.
- Temporary artefacts D = 1 (the ⑧a1 doc carrier deleted its one-shot helper scripts inside the slice directory before the pre-deletion count was recorded - an ordering violation, escalated to the user at ⑩; the final zero-residue reading passed); R = 0. From the next slice on, the ⑧a doc carrier must not delete its helper scripts inside `tmp/<slice-id>/`; the ⑧c wrap-up removes them.

## §8 Explicitly not executed

No change to the directive body (`DELIVERY_DIRECTIVE{,_EN}.md` zero diff); no change to the decision bodies of `internal/gates/**`, `internal/chain/**`, `internal/adapters/**`, `internal/tickets/**` or `internal/guard/**`; no change to schemas / fixtures / `.github/workflows/**` / `go.mod`; no change to CONTRACTS semantics (the tense update of the §4.2 last sentence is item 10 of group [five] and lands in ⑧a2); no change to the v50–v55 snapshots; no commit / push; no new slice started; no gitee push.

artifacts fence: the archived entry list of this slice (path + disposition, TAB separated, ASCII lexicographic; 38 entries).

```artifacts
docs/ADR_REGISTER.md	repo
docs/ADR_REGISTER_EN.md	repo
docs/ARCHITECTURE.md	repo
docs/ARCHITECTURE_EN.md	repo
docs/CONTRACTS.md	repo
docs/CONTRACTS_EN.md	repo
docs/DEV_PLAN.md	repo
docs/DEV_PLAN_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/evidence/P69-T027B-E2E-WIRING/00-briefing.md	repo
docs/validation/evidence/P69-T027B-E2E-WIRING/01-arch/design-r1.md	repo
docs/validation/evidence/P69-T027B-E2E-WIRING/02-batch/batch-01-seam-enumeration.md	repo
docs/validation/evidence/P69-T027B-E2E-WIRING/02-batch/batch-02-run-callchain.md	repo
docs/validation/evidence/P69-T027B-E2E-WIRING/04-reviews/final-r1.md	repo
docs/validation/evidence/P69-T027B-E2E-WIRING/04-reviews/prereview-r1.md	repo
docs/validation/evidence/P69-T027B-E2E-WIRING/04-reviews/scan-r1.md	repo
docs/validation/evidence/P69-T027B-E2E-WIRING/06-gates/mutation-evidence.txt	repo
docs/validation/evidence/P69-T027B-E2E-WIRING/MANIFEST.md	repo
docs/validation/evidence/P69-T027B-E2E-WIRING/SHA256SUMS.txt	repo
docs/validation/evidence/P69-T027B-E2E-WIRING/freedom-list.md	repo
docs/validation/evidence/P69-T027B-E2E-WIRING/input-list.md	repo
docs/validation/p44-gate-readiness.md	repo
docs/validation/p44-gate-readiness_EN.md	repo
docs/validation/p57-closeout-assessment.md	repo
docs/validation/p57-closeout-assessment_EN.md	repo
docs/validation/t027b-e2e-wiring.md	repo
docs/validation/t027b-e2e-wiring_EN.md	repo
internal/console/cli.go	repo
internal/console/cli_multitask_wiring_test.go	repo
internal/console/readiness_gate.go	repo
internal/console/readiness_gate_test.go	repo
internal/console/runtime_agent.go	repo
internal/console/runtime_agent_e2e_test.go	repo
internal/console/runtime_multitask_red_test.go	repo
internal/console/runtime_supervised.go	repo
internal/console/runtime_supervised_test.go	repo
internal/console/supervised_cmd.go	repo
```

sixstage fence: the verbatim Section A–E transcript of ⑥ independent scan r1.

> The ⑥ conclusion-line transcription convention is defined by G7-f in §6.3 (v1.48).

```sixstage
## Section A

INDEPENDENT SCAN: PASS

## Section B

(no findings)

## Section C

- The contract is matched exactly by the level-1 gate semantics in [internal/gates/readiness.go](internal/gates/readiness.go#L141-L170): `EvaluateT027BReadiness` is a pure function, and for `GoalProbe` it checks presence only, while the actual quota consumption is delegated to `ProbePermit.Require`. That matches the contract sentence: “EvaluateT027BReadiness is the only construction entry and is a pure function; zero I/O; no permit consumption.”
- The one-shot lifecycle is enforced in [internal/gates/probe_permit.go](internal/gates/probe_permit.go#L1-L95) and [internal/gates/probe_permit.go](internal/gates/probe_permit.go#L176-L244): the state machine is `Unused → Claimed → Spent`, the linearization point is the `O_CREATE|O_EXCL` marker file persisted with `Sync`, and `Complete` is idempotent only after a valid claim. This is the direct contract anchor for the “fail-closed” rule.
- The production wiring sits in the allowed scope only: [internal/console/readiness_gate.go](internal/console/readiness_gate.go#L1-L24), [internal/console/cli.go](internal/console/cli.go#L282-L395), [internal/console/runtime_agent.go](internal/console/runtime_agent.go#L120-L190), and [internal/console/runtime_supervised.go](internal/console/runtime_supervised.go#L259-L320). The flow is consistent with the design: CLI resolves and validates the permit before routing, AgentRunner evaluates before preconditions and consumes only after preconditions pass, and supervised starts enforce a fresh-start ledger gate before reusing the permit.
- The negative-path test bundle in [internal/console/readiness_gate_test.go](internal/console/readiness_gate_test.go#L34-L205) is the real behavioral proof for the five acceptance conditions: nil permit denied, spent permit denied, missing CLI flag denied, incomplete record set kept fail-closed, and successful assembly consumes the permit exactly once.
- The mutation evidence in [docs/validation/evidence/P69-T027B-E2E-WIRING/06-gates/mutation-evidence.txt](docs/validation/evidence/P69-T027B-E2E-WIRING/06-gates/mutation-evidence.txt) is especially important: M1–M5 are all red when the guard is removed or weakened, then restored green after the original state is restored. That is exactly the correct “can this be falsified?” check.
- The recorded gate readouts are all green in the evidence bundle: [tmp/P69-T027B-E2E-WIRING/06-gates/p4b-07-selftest.txt](tmp/P69-T027B-E2E-WIRING/06-gates/p4b-07-selftest.txt), [tmp/P69-T027B-E2E-WIRING/06-gates/p4b-08-index.txt](tmp/P69-T027B-E2E-WIRING/06-gates/p4b-08-index.txt), [tmp/P69-T027B-E2E-WIRING/06-gates/p4b-09-tree.txt](tmp/P69-T027B-E2E-WIRING/06-gates/p4b-09-tree.txt), and [tmp/P69-T027B-E2E-WIRING/06-gates/p4b-10-ci.txt](tmp/P69-T027B-E2E-WIRING/06-gates/p4b-10-ci.txt). The index/tree/ci reports all say `TOTAL_FAIL=0`, and the selftest report is `737/737`.

## Section D

- If the CLI presence gate in [internal/console/cli.go](internal/console/cli.go#L295-L395) is removed:
  - The check in [internal/console/readiness_gate_test.go](internal/console/readiness_gate_test.go#L167-L205) goes red: `TestRunRefusesMissingProbePermitFlag` fails because the missing `--probe-permit` path becomes a later record-read or assembly error instead of a fail-closed refusal.
  - The “which check is not red?” answer is that the existing `--request-record` and “local runtime only supports noop steps” validations still remain red; the unique thing that becomes green wrongly is the missing-permit gate itself.

- If the `Require` claim is removed from the shared helper in [internal/console/readiness_gate.go](internal/console/readiness_gate.go#L8-L24) or from the assembly entry in [internal/console/runtime_agent.go](internal/console/runtime_agent.go#L146-L180):
  - `TestAgentRunnerAssemblyRefusesSpentProbePermit` in [internal/console/readiness_gate_test.go](internal/console/readiness_gate_test.go#L105-L166) goes red because presence-only evaluation would allow an already-spent permit to proceed.
  - The “not red” part is the initial `EvaluateT027BReadiness` presence check; that one is intentionally not enough to decide quota consumption, so it is not the guard that catches a spent permit.

- If the deferred `Complete` release is removed:
  - The mutation case M5 in [docs/validation/evidence/P69-T027B-E2E-WIRING/06-gates/mutation-evidence.txt](docs/validation/evidence/P69-T027B-E2E-WIRING/06-gates/mutation-evidence.txt) goes red; the second run would report `ErrProbePermitInUse` instead of `ErrProbePermitConsumed`.
  - This is the exact “complete must be recorded” invariant; no other check catches it.

- If the supervised fresh-start gate in [internal/console/runtime_supervised.go](internal/console/runtime_supervised.go#L269-L320) is removed:
  - A second supervised run on an existing interaction ledger can wrongly pass the permit through again; the re-entry behavior becomes fail-open.
  - The current evidence does not show a dedicated second-park mutation on a real ledger; the invariant is still logically covered by the ledger check and the design, but it is the one place where an OS-level concurrency test would be the strongest proof.

- If the stronger “reject before marker creation” ordering in [internal/console/runtime_agent.go](internal/console/runtime_agent.go#L146-L180) is removed:
  - The nil/spent permit tests in [internal/console/readiness_gate_test.go](internal/console/readiness_gate_test.go#L66-L166) are the ones that go red, because the run-root marker could be created before the permit failout.
  - The design is intentionally ordered to prevent a marker from appearing on a fail-closed refusal.

- The no-coverage invariant I would call out is the real cross-process race itself: the code relies on `O_CREATE|O_EXCL` plus `Sync`, but the artifacts do not contain a true concurrent two-process reproduction of the race. That is not a logic bug in this slice; it is the one remaining externality that would require a live two-process OS test to eliminate all doubt.

## Section E

- The things that are not directly observable from the read-only evidence:
  - A real two-process race on the same permit marker across separate OS processes.
  - A live Windows-specific agent-runner publication-durability path that depends on the actual pinned CLI and replay store behavior.
  - A full end-to-end run of the real CLI in a production environment outside the recorded artifacts.
  - The precise behavior of the Linux-side `GOOS=linux` path beyond the recorded pass states; the output files are empty, but the higher-level gate reports still confirm overall pass.

- What would overturn this PASS:
  - Any proof that a `nil` or already-spent permit is accepted into a real run.
  - Any proof that a second supervised entry with an existing ledger bypasses the fresh-start gate.
  - Any proof that the permit is consumed without `Require`, or that `Complete` is skipped on the success path.
  - Any diff outside the allowed scope that modifies the gate logic itself under [internal/gates](internal/gates) or alters the chain/adapters semantics.
  - Any real runtime failure where the mutation tests in [docs/validation/evidence/P69-T027B-E2E-WIRING/06-gates/mutation-evidence.txt](docs/validation/evidence/P69-T027B-E2E-WIRING/06-gates/mutation-evidence.txt) no longer hold.

- The current evidence is strong enough to say the slice is fail-closed in the intended sense and stays within the declared contract. The remaining uncertainty is only the untested OS-level race behavior, not a logical defect in the gate ordering or the permit lifecycle.

INDEPENDENT SCAN: PASS
```
