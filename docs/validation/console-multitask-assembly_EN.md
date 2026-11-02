# CONSOLE-MULTITASK-ASSEMBLY Verification Report

- Positioning: v39 layer 2 item 3 (P50); console multi-task production assembly expansion.
- Directive: v1.60; start anchor `8bd8224`; date 2026-11-03.
- Implementation class: `[SLICE]` (implementation); size tier L; ledger row `P50`.

## 1. Slice Summary

- **Goal**: lift the console assembly layer's single-task limit so several tasks of one chain are assembled and driven serially under the existing A6 exclusivity contract.
- **Scope**: the production assembly layer (`internal/console/runtime_agent.go`, `runtime_supervised.go`), a new adapters registry (`agent_runner_intent.go`), the offline harness and test rewrites; `internal/chain` plus the launcher / mirror / slot implementations are untouched.
- **Ruling**: Plan A (user ruling 2026-11-03) - CLI wiring deferred; this slice is limited to the production assembly layer plus the offline harness plus test rewrites; stage 2 was dispatched in small serial steps 2a-2d.
- **Risk rating**: tier-L core assembly refactor; never downgraded to `[ROUTINE]`; no gitee push.
- **Untouched**: the existing A2/A4/A5 semantics, historical validation reports, `.github/workflows` and gate scripts.

## 2. Protocol First

- Stage 1 verdict: no wire / Schema / fixture change (the mirror, slot and starting marker are store-local and take no part in the R/C recordHash); the CONTRACTS A6 paragraph must change to state that "one mirror per runRoot" is **instant-level** exclusion, superseding the P45 D9 single-task reading.
- Stage 2a landing: CN `docs/CONTRACTS.md` gains one paragraph at the end of A6 (+1/-0) with the EN `docs/CONTRACTS_EN.md` synced (+1/-0); both files keep their BOM and LF.
- Contract text (EN, verbatim):

> **Multi-task serial assembly (P50)**: the runRoot-level identity mirror and the launch slot are the exclusion point for "**at most one running launch at any instant**", not a one-time uniqueness over the run lifetime. One chain under the same `runId` may contain multiple tasks (one isolated-workspace code step per task); the console assembly layer binds each `(taskId, stepId, attempt)` to its own already-persisted `agent-runner-request` (distinct `requestId` values, with R/C/terminals/artifacts isolated by requestId), and the chain engine schedules them **serially**: only after the previous task's launch has been provably stopped (slot released, the process named by the mirror is dead, or the mirror already atomically rotated by the new launch) may the next task publish R and launch-intent; mirror rotation means atomically overwriting the old mirror with the new launch identity, consistent with the existing semantics that "a process named by the mirror that is dead is stale and may be replaced". Concurrent running launches within one runRoot are forbidden; a loser of a concurrent Start still must not create log artifacts or overwrite the identity mirror. The event log remains the single total-order source `events/state-events.jsonl` for that runRoot; each task's six frozen facts attach via the positional evidence of that task's own terminal routing event (A5 `evidence[2..7]`) and must not be reconciled across tasks.

- Implementation mapping: `agentRunnerAssemblyPlan` / `supervisedAssemblyPlan` validate "one step per task"; `agentRunnerAssemblyRecords` enforces distinct `requestId` values with attempt fixed at 1; `agentRunnerBindTasks` builds an independent workspace, replay store and terminal publisher per task; the driver calls `proveTaskStopped` before the next `engine.Run`; the adapters slot / mirror are untouched, which already satisfies "a losing concurrent Start creates no artifacts and does not overwrite the mirror".

## 3. Red-to-Green Flip

- Stage 3 red baseline (production code untouched): all three cases in `internal/console/runtime_multitask_red_test.go` PASS - their assertion shape is "the current implementation refuses multi-task"; `gofmt -l` empty, `go vet ./internal/console/` clean, whole package `ok console 8.495s`. Evidence `03-tests/logs/03-tdd-red-run.txt`, `03-tests/red-baseline.md`.
- Stage 2 flip: the three cases were rewritten into accept form and run green - `TestAgentRunnerAssemblyAcceptsMultiTaskDefinitionShape`, `TestSupervisedAssemblyAcceptsMultiTaskDefinitionShape`, `TestMultiTaskDefinitionShapeIsOtherwiseValid`.
- Single-task gate wording zeroed (mechanical reading): `requires exactly one task` appears **0** times in either production file; the new wording is `requires at least one task` (`runtime_agent.go`, `runtime_supervised.go`), and the step-shape gate stays by design. Evidence `04-impl/logs/04-refusal-strings.txt`.
- **Regression leg**: the accept-form cases permanently assert that a 2-task definition is accepted, so any regression back to "exactly one" goes red immediately; `TestMultiTaskDefinitionShapeIsOtherwiseValid` keeps the per-task shape leg.
- Production-entry evidence: `TestAgentRunnerDriverTwoTaskSerialAssembly` calls `ExecuteAgentRunnerRun` directly; `TestAgentRunnerDriverStopsTheLaunchOfTheTaskThatJustFailed` pins H1; `TestAgentRunnerAssemblyRefusesNonInitialAttempt` pins M1.

## 4. Implementation and Remediation

| Step | Carrier | Landing point | Result |
|---|---|---|---|
| 2a | `prfrail-implementer` | `docs/CONTRACTS{,_EN}.md` | +1/-0 each; the CN paragraph matches the design draft verbatim |
| 2b | `prfrail-implementer` | `internal/console/runtime_agent.go` | assembly plan, record registry, per-task binding, serial driver; `return fail(` grows from 7 to 8 sites |
| 2c | `prfrail-implementer` | `internal/console/runtime_supervised.go` | `supervisedAssemblyPlan`, per-task state map, parked and return phases |
| 2d | `prfrail-implementer` | `internal/adapters/agent_runner_intent.go` and tests | `AgentRunnerStepIntentRegistry` (+90/-0), 4 new test files, 2 existing test rewrites |

- The four stage-5 r1 remediations and their closure:
  - **H1 (High)**: the non-awaiting branch bound the previous task's `requestId`, so an orphan launch of a non-first task was never stopped ⇒ switched to observation-first resolution (`agentRunnerDrivingRequestID` plus `currentRequestID`).
  - **H2 (High)**: the production multi-task driver had zero behavioural coverage ⇒ added `runtime_multitask_driver_test.go`, three fixtures calling the production entry, plus all three stop-proof states (stopped / still alive / report missing).
  - **M1 (Medium)**: attempt was not validated ⇒ `agentRunnerAssemblyRecords` now rejects `Attempt != 1`, before the marker and the engine.
  - **M2 (Medium)**: the supervised five-fact positional order was not pinned ⇒ recomputed from the production source functions and compared position by position.
- New seams (all console-local and unexported): `testStub`, `agentRunnerLaunchReceiptReader` / `agentRunnerLaunchStopper`, `agentRunnerDispatchWaiter`, plus the dispatch / admission / intent / postflight muxes; production carriers reuse existing adapters types and `internal/chain` is untouched.
- Deviation recorded: `AgentRunnerStepIntentRegistry` is implemented and unit-tested, yet the production path still builds a preparer per record (7 Info 1; a later slice adopts it or deletes one of the two).

## 5. Review Chain

| Step | Carrier | Verdict | Budget |
|---|---|---|---|
| 5 r1 | `deep-reasoner` + `DeepSeek V4 Pro (deepseek)` | `PRE-REVIEW: PASS WITH FIXES` (2 High + 2 Medium) | 1/1 |
| 5 r2 | same as above | `PRE-REVIEW r2: PASS` (`SEVERITY=H0/M0/L0/INFO2`) | 1+1 used up |
| 6 | `independent-reviewer` + `MAI-Code-1.1-Flash (copilot)` | `INDEPENDENT SCAN: PASS` (Section B no findings) | 1/3 |
| 7 | `independent-reviewer` + `GPT-5.3-Codex (copilot)` | `FINAL REVIEW: PASS` (`SEVERITY=INFO2`) | 1/1 |

- Step 6 self-run: `go test ./internal/console -run 'MultiTask|AgentRunnerAssembly|SupervisedAssembly|MultiTaskDefinitionShape|Assembly' -count=1` returned `ok console 3.541s`.
- The two step-7 Info items are registered: (1) the intent registry is not wired into production; (2) local Go 1.27 differs from CI Go 1.22 ⇒ step 9 takes the CI Go 1.22 legs as the final release evidence.
- Two step-7 release conditions (non-blocking): step 9 must capture both CI legs; step 8 must land the ARCHITECTURE / DEV_PLAN write-back (see section 10).

- The step-6 verdict is transcribed verbatim below:

```sixstage
## Section A
This review was a read-only forensic check of the live P50 multi-task serial assembly in [docs/CONTRACTS.md](docs/CONTRACTS.md#L255-L289), [internal/console/runtime_agent.go](internal/console/runtime_agent.go#L86-L200), [internal/console/runtime_supervised.go](internal/console/runtime_supervised.go#L332-L524), and [internal/adapters/agent_runner_intent.go](internal/adapters/agent_runner_intent.go#L18-L83). I did not write or modify any files; the review package was taken directly from the on-disk working tree and checked against the contract and focused tests.

## Section B
(no findings)

The effective contract for P50 is preserved in the current tree: one `runRoot` remains the serial launch barrier, each task/step/attempt is bound by its own request record, and the next task may launch only after the previous task is proven stopped. I did not find a material violation of the A6 semantics in the live implementation.

## Section C
Evidence supporting the pass result:

- The contract explicitly says the same `runId` may contain multiple tasks, but the run-root identity mirror and launch slot remain the exclusivity point: “at most one running launch at any instant” and “the next launch must be proven stopped before the next R/launch-intent can be published.” See [docs/CONTRACTS.md](docs/CONTRACTS.md#L279-L289).
- The registry enforces per-task binding and rejects cross-task drift, duplicate step bindings, duplicate request IDs, and cross-run mixing. See [internal/adapters/agent_runner_intent.go](internal/adapters/agent_runner_intent.go#L27-L83).
- The serial engine barrier is explicit in the driver: `proveTaskStopped` calls `managedProcessesGone`, which fails closed on missing pid reports or still-alive PIDs, so the next task cannot start while the previous one is not proven stopped. See [internal/console/runtime_agent.go](internal/console/runtime_agent.go#L912-L926) and [internal/console/runtime_agent.go](internal/console/runtime_agent.go#L853-L910).
- The acceptance tests cover the exact P50 contract shape: two-task definitions are accepted, retired single-task refusal wording is rejected, and the serial harness is enforced in the real driver path. See [internal/console/runtime_multitask_red_test.go](internal/console/runtime_multitask_red_test.go#L124-L209) and [internal/console/runtime_multitask_driver_test.go](internal/console/runtime_multitask_driver_test.go#L283-L332).
- Fresh verification: I ran:
  `cd D:\LZProjects\prfrail; go test ./internal/console -run 'MultiTask|AgentRunnerAssembly|SupervisedAssembly|MultiTaskDefinitionShape|Assembly' -count=1`
  and it returned:
  `ok      github.com/larsonzh/prfrail/internal/console    3.541s`

## Section D
The claim is falsifiable by the exact failure modes the code is designed to reject:

- Duplicate task/step or duplicate request ID would fail in `NewAgentRunnerStepIntentRegistry`, which is the direct guard for task aliasing and request leakage. See [internal/adapters/agent_runner_intent.go](internal/adapters/agent_runner_intent.go#L27-L83).
- A still-alive PID or a missing PID report would fail closed in `managedProcessesGone` and therefore also in `proveTaskStopped`. See [internal/console/runtime_agent.go](internal/console/runtime_agent.go#L853-L926).
- A non-serial or re-used run-root launch would be visible in the driver order and in the multi-task harness tests. See [internal/console/runtime_multitask_driver_test.go](internal/console/runtime_multitask_driver_test.go#L283-L332).

In other words, if the contract were being violated in the current tree, the nearest observable failure points are the registry binding, the stop-proof gate, or the multi-task serial-harness assertions; none of them fail under the live code.

## Section E
This is a current-state review of the checked tree, not a guarantee for every future path. The implementation comments themselves call out the constraint that the serial/offline assumption is what makes the `dispatched` observation sound and non-resetting. A future resume/concurrent driver that reuses stale state would require a fresh re-review; under the current serial/offline design, I did not find a contract gap.

I also checked the live worktree via `git status --short` and `git diff --numstat -- .`; no review-only file writes were introduced by this scan.

INDEPENDENT SCAN: PASS
```

## 6. Gate Readings

- Local run (Go 1.27.0 / windows): `gofmt -l` empty, `go build ./...` exit 0, `go vet ./...` exit 0, `go test -count=1 ./...` **15 packages all ok, FAIL=0** (console 18.188s / adapters 51.380s); raw output `04-impl/logs/04-full-suite-green.txt`.
- Step-5 r2 re-run after the fix round: also 15 packages all ok (console 22.6s / adapters 51.6s / chain 6.6s).
- Step-9 measurement: all three scopes tree / index / ci are `TOTAL_FAIL=0` (the tree/index measured counts and the `G6-1 scanned` count live in `08-gates/gate-readings.txt`); `SELFTEST: PASS 697/697`; the base four legs (gofmt/build/vet/test) are all green (`ok=15`). The final criterion is the CI Go 1.22 legs (Ubuntu plus Windows) at the `Test` and `Gates` steps.
- Local `-race` is out of scope; the Race leg of the six core packages is carried by the P49 CI coverage.

## 7. Accounting and Environment Facts

- **R = 1**: the step-5 r1 remediation round (H1 / H2 / M1 / M2 and I5 fixed, then step 5 re-run).
- **D = 1**: stage-2 dispatch-channel network failures, an environment fact rather than a code defect; disposition and evidence below.
- The stage-2 dispatch channel failed 5 times: `UND_ERR_CONNECT_TIMEOUT` x4 and `ECONNRESET` x1; failures 4 and 5 returned an error only after partial backend execution, leaving **two partial artifacts** (CONTRACTS not matching the draft, supervised side failing to compile).
- Disposition: `git checkout` reverted exactly the 6 affected files (keeping step 0, step 3 and the evidence pack), the partial diff was preserved as `tmp/CONSOLE-MULTITASK-ASSEMBLY/02-partial-attempt-diff.txt` and the blockage logged as `02-impl-dispatch-blocked.txt`; after recovery `go build ./...` was green and stage 2 switched to the small serial steps 2a-2d.
- Environment facts: local Go 1.27.0 (windows/amd64) while the repository gates and CI use Go 1.22; the version gap is covered by the CI legs.
- Probe budget: none (zero real external model candidate calls in this slice).

## 8. Open Items

- **CLI wiring deferred**: the multi-record supply contract for `cmd/prfrail` and the console flags is undecided and left to a later slice; this slice stops at the `RequestRecords` registry boundary.
- **Intent registry not wired into production** (7 Info 1): `agentRunnerIntentMux` still builds a preparer per record, behaviourally equivalent.
- **Windows R/C write path not run locally**: on Windows `PublishDurability` is always `unproven`, so the green tests only assert per-requestId path and digest distinctness; the write path is carried by the Ubuntu CI leg.
- **First CI run observation**: both legs' `Test` / `Gates` results must be captured after the commit; a first-run red is registered under the DR rules.
- **Existing exported sentinel retained**: `ErrSupervisedStepRequired` is unused and undeleted (stage 2c deviation 6), keeping the change surface narrow.

## 9. Artifacts Fence

- Process files: `normalize-step0.js`, `verify-step0.js`, `rehash.js`, `02-impl-dispatch-blocked.txt`, `02-partial-attempt-diff.txt`, `02a-gate-index.txt`, all disposed as `local-only` together with `tmp/CONSOLE-MULTITASK-ASSEMBLY/` (cleared at 8c).

```artifacts
docs/ARCHITECTURE.md	repo
docs/ARCHITECTURE_EN.md	repo
docs/CONTRACTS.md	repo
docs/CONTRACTS_EN.md	repo
docs/DEV_PLAN.md	repo
docs/DEV_PLAN_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/console-multitask-assembly.md	repo
docs/validation/console-multitask-assembly_EN.md	repo
docs/validation/evidence/CONSOLE-MULTITASK-ASSEMBLY/00-briefing.md	repo
docs/validation/evidence/CONSOLE-MULTITASK-ASSEMBLY/01-arch/design.md	repo
docs/validation/evidence/CONSOLE-MULTITASK-ASSEMBLY/03-tests/logs/03-tdd-red-run.txt	repo
docs/validation/evidence/CONSOLE-MULTITASK-ASSEMBLY/03-tests/red-baseline.md	repo
docs/validation/evidence/CONSOLE-MULTITASK-ASSEMBLY/04-impl/implementer-reports.md	repo
docs/validation/evidence/CONSOLE-MULTITASK-ASSEMBLY/04-impl/logs/04-02a-gate-index.txt	repo
docs/validation/evidence/CONSOLE-MULTITASK-ASSEMBLY/04-impl/logs/04-full-suite-green.txt	repo
docs/validation/evidence/CONSOLE-MULTITASK-ASSEMBLY/04-impl/logs/04-refusal-strings.txt	repo
docs/validation/evidence/CONSOLE-MULTITASK-ASSEMBLY/05-prereview/review-r2.md	repo
docs/validation/evidence/CONSOLE-MULTITASK-ASSEMBLY/05-prereview/review.md	repo
docs/validation/evidence/CONSOLE-MULTITASK-ASSEMBLY/06-scan/review.md	repo
docs/validation/evidence/CONSOLE-MULTITASK-ASSEMBLY/07-final/review.md	repo
docs/validation/evidence/CONSOLE-MULTITASK-ASSEMBLY/08-gates/gate-readings.txt	repo
docs/validation/evidence/CONSOLE-MULTITASK-ASSEMBLY/09-glm/glm-1-reference-sweep.md	repo
docs/validation/evidence/CONSOLE-MULTITASK-ASSEMBLY/09-glm/glm-2-test-coverage.md	repo
docs/validation/evidence/CONSOLE-MULTITASK-ASSEMBLY/MANIFEST.md	repo
docs/validation/evidence/CONSOLE-MULTITASK-ASSEMBLY/SHA256SUMS.txt	repo
docs/validation/evidence/CONSOLE-MULTITASK-ASSEMBLY/freedom-list.md	repo
docs/validation/evidence/CONSOLE-MULTITASK-ASSEMBLY/input-list.md	repo
internal/adapters/agent_runner_intent.go	repo
internal/adapters/agent_runner_intent_registry_test.go	repo
internal/console/runtime_agent.go	repo
internal/console/runtime_agent_e2e_test.go	repo
internal/console/runtime_multitask_agent_test.go	repo
internal/console/runtime_multitask_driver_test.go	repo
internal/console/runtime_multitask_red_test.go	repo
internal/console/runtime_multitask_supervised_test.go	repo
internal/console/runtime_supervised.go	repo
internal/console/runtime_supervised_test.go	repo
tmp/CONSOLE-MULTITASK-ASSEMBLY/	local-only
```

## 10. Documentation Sync Status

- `docs/ARCHITECTURE{,_EN}.md`: three paragraphs appended to section 3 (multi-task serial assembly / A6 mirror sharing semantics / scheduling mechanism), +6 lines each (3 body paragraphs plus 3 blank lines), with equal CN / EN line counts.
- `docs/DEV_PLAN{,_EN}.md`: the "November Refactor Milestone (T027-B)" list grows from 6 entries to the full P37-P50 set (status and report paths), with equal added and deleted line counts on the CN / EN sides.
- `docs/t027/REMAINING_SLICES{,_EN}.md`: the P50 evidence cell in the status table is written back and the P50 section gains a completion-record row; the status glyph was converted to `✅ Completed` at 8b; the commit hash and CI evidence are recorded after the authorised commit.
- Synced: the `08-gates/` payload of `MANIFEST.md` was appended at stage 9; `08-ci/` is appended after the authorised commit.

This slice stops at ⑩ awaiting authorization.
