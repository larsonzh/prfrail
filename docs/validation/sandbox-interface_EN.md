# ProofRail SANDBOX-INTERFACE Report (P74 · layer-1 part-B introduction · implementation [SLICE] · tier M)

Date: 2026-11-15. Status: ⑧ document assembly final (①–⑦ closed). Slice id: `SANDBOX-INTERFACE` (P74 · layer-1 part-B introduction slice · implementation `[SLICE]` · tier M).

## 1. Executive summary

- Slice: `SANDBOX-INTERFACE` (P74, layer-1 part-B introduction slice, implementation `[SLICE]`, tier M); controller `DeepSeek V4.1 Flash (deepseek)`.
- Delivered: `internal/guard/sandbox.go` (the `Sandbox` check-layer interface + the three-state `Verdict` + the two sentinels + the `CheckSandbox` gate), `internal/guard/sandbox_windows.go` (form alpha), `internal/guard/sandbox_unsupported.go` (`!windows` fail-closed), `internal/console/sandbox_preflight.go` (replaceable injection point + `preflightSandbox` + `sandboxEnforcementFromEnv`), the assembly check points of both dispatch paths (`runtime_agent.go` / `runtime_supervised.go`, both ahead of the permit), the CLI trigger injection (`cli.go` / `supervised_cmd.go`), four test groups (red / green / mutation / trigger), and the `CONTRACTS{,_EN}.md` §7.2 contract section.
- Conclusion: the check layer sits ahead of the probe-permit claim; `Denied`, `Unknown` and a checker error are all fail-closed; S1 defaults to `off`, so the default control flow is sentence-for-sentence equivalent to the pre-slice build; with `PRFRAIL_SANDBOX_CHECK=enforce` both dispatch paths are refused (the S1 platform implementations always answer `Unknown`). **Part B does not unlock S-G1 / S-G8.**
- Metering: ⑤ twice (`PRE-REVIEW: PASS` → three Low fixes → `PRE-REVIEW-R2: PASS`); ⑥ once (`INDEPENDENT SCAN: PASS`, Section B no findings); ⑦ final review 1 + re-review 0 (`FINAL REVIEW: PASS`, Section B no findings); blind review 0; GLM batch assists 2 (7178 / 3016 characters).
- Counts: **R = 2 (session 1) + 1 (session 2) = 3**, **D = 0**. Session-1 rework: ① r1 → r2 falsifiability correction; ②a first attempt was CN `+35` / EN `+44` and tripped `G3(a)`, then rewritten symmetric. Session-2 rework: the first ③0 skeleton followed `design-r1.md` §2.2 and landed the platform implementation files (`sandbox_windows.go` / `sandbox_unsupported.go`), which contradicted item 2 of the continuation state file ("③0 = inert skeleton, platform implementations belong to ②b") ⇒ reshaped to the state file and the platform files were removed (same root cause, one correction, counted R+1).

### 1.1 Anomaly and deviation register

a) **Deviation P74-1 (segmented execution, §13.1.17 lenient reading, acknowledged by the user on 2026-11-15)**: this single slice was executed in two sessions because session 1 ran out of context (session 1 = step 0 + ① r2 + ②a; session 2 = ②a archival + ③0–⑩). The lenient reading ("end the session, continue the slice") was applied; this is the first use and is recorded for later implementation slices.
b) **Deviation P74-2 (the ②a archive was completed by session 2)**: session 1 never copied `design-r1.md` into the evidence pack `01-arch/` (its own state file lists this as unfinished item 1). Session 2 copied it before ③0 (normalised to BOM+LF from the original) and recomputed `SHA256SUMS.txt`.
c) **Deviation P74-3 (batch assists re-dispatched; process execution, not rework)**: session 1 recorded two GLM batch-assist conclusions but never persisted them. Session 2 re-dispatched both per §二 of the briefing and archived them as `02-batch/` payloads (7178 / 3016 characters); the session-1 conclusions are not quoted and are not evidence. **A re-dispatch re-executes the existing process per §二 of the briefing, so it is not rework and is not counted into R** (R=3 consists only of the P74-1 / P74-2 items and the ① r2 rework, itemised in f below).
d) **Deviation P74-4 (⑧ carrier deviation: controller self-execution)**: §一 of the briefing assigns ⑧ to `prfrail-documenter` + `DeepSeek V4.1 Flash`; the controller of this segment is the same model (Flash), so to avoid a second-hand restatement of a long fact chain ⑧ was self-executed. All three required points are recorded below:
   - **Deviation id**: P74-4.
   - **Independent impact assessment**: role separation is weaker (⑧ and the ①–⑦ dispatching are the same Flash controller), but ⑤/⑥/⑦ were carried by independent models and personalities and their conclusions are unaffected; deliverable consistency is checked by mechanical gates (G7-f sixstage fence, G8 frozen-pack self-consistency, G3 symmetry) rather than by human discretion; **no business semantics or fact-chain content is changed**.
   - **Must not be used as a precedent**: later slices still assign ⑧ to `prfrail-documenter` per the §一 dispatch table; this self-execution may not be cited as authority for "the controller may do ⑧".
e) **Deviation P74-5 (⑦ cross-compile evidence completed)**: ⑦ noted that `linux-build.txt` / `linux-vet.txt` were empty (the commands succeeded silently). The controller re-ran them and wrote a single self-describing evidence file `06-gates/linux-build.txt` (`build -> exit 0`, `vet -> exit 0`); the two empty files were removed. The change is an evidence readability improvement and does not alter the code state ⑦ reviewed.
f) **R = 3 trips the "high-risk slice" marker (§9.2)**: this slice has R=3 (itemised: ① r1→r2 falsifiability correction = 1; ②a asymmetric rewrite = 1; ③0 skeleton shape inconsistent with the continuation state file, one-time correction = 1; three distinct root causes, not merged) ⇒ per §9.2 "R≥2 ⇒ mark the slice high-risk and execute the §5.2 high-risk row", **the §5.2 high-risk row was executed** (① required / ② required / ③ required / ⑤ required / ⑥ required / ⑦ required / ⑧ required / ⑨ required / ⑧c required, every column delivered with evidence); R=3 stays below the "R>3 ⇒ pause for human review" threshold, so no human pause was triggered. **This is the first slice in the project to trip that marker** (P72 / P73 had R=0 and did not trip it); the fact is carried into the v58 fact pack.
g) **⑦ quota tiering change recorded (§7.5 "the row with the most re-reviews wins")**: at kick-off the ⑦ quota was set at tier M, i.e. 1 final review + 2 re-reviews; once R=3 tripped the high-risk marker, §7.5's "when several rows match, take the one with the most re-reviews" raised the ceiling to the high-risk row, i.e. 1 final review + 3 re-reviews. **Actual use = final 1 + re-review 0 + blind 0**, below both the tier-M and the high-risk ceilings ⇒ **no overrun and no extra call**. The change follows an existing row-matching rule, and **it is not the "no mid-flight tier upgrade" case of §7.5** (the ⑦ model stayed on the regular `gpt-5.3-codex` tier throughout).

## 2. Current-state check (① facts, persisted by session 1, re-checked here)

- `internal/guard/`: no interface and no `ProcessSupervisor`; `ProcessSpec` carries an `OnSpawn` callback; `ProcessIdentity{PID,StartToken}`; `ManagedProcess` keeps its fields unexported.
- `internal/chain/agent_runner_launcher.go:14`: the `AgentRunnerLauncher` interface; `:26` `AgentRunnerLaunchRequest`.
- `internal/adapters/agent_runner_pinned_cli.go:218` calls `guard.StartManaged`; `:300` calls `guard.StopProcessIdentity`.
- `internal/gates/probe_permit.go:210` `Require`: cache fast refusal → `persistClaim()` (O_EXCL) is the linearisation point and the quota is never refunded; `:253` `Complete` CAS switches to spent.
- `internal/console/runtime_agent.go`: Step 0 readiness gate (no quota) → Step 1 preconditions → permit (from `:197`, moved by this slice) → Step 2 run-root marker → … → first dispatch.
- `internal/console/runtime_supervised.go`: only a `freshStart` claims the permit.
- `docs/CONTRACTS.md`: §7 (then `:210`–`:322`) had a single `###` heading, 7.1 (`:310`); §7.2 is new in this slice.

## 3. ① design (`design-r1.md`, including the r2 correction)

- **Method set**: `Check(ctx context.Context, spec CheckSpec) (Verdict, error)`; `Verdict{Status, Reason}` with `VerdictUnknown` / `VerdictAllowed` / `VerdictDenied` (`Unknown` is the `iota` zero value, so the zero value is the fail-closed default); `CheckSpec{Reason}` carries the audit label only. No `Create/Exec/Destroy/Terminate`, and no overlap with `ProcessSpec` / `ProcessIdentity` / `ManagedProcess`.
- **Error classification**: `ErrSandboxUnavailable` (capability unavailable) and `ErrSandboxRejected` (a non-`Allowed` fail-closed refusal).
- **Check point**: option (i) — **ahead of the permit claim** and ahead of the Step 2 run-root marker. Reason: the permit quota is one-shot and never refunded, and the existing `runtime_agent.go` comment states "a configuration refusal must not burn the one-shot quota".
- **Fail-closed form**: `!windows` always answers `(Unknown, ErrSandboxUnavailable)`; Windows uses **alpha** (always `Unknown`, `nil` error) — S1 runs no check, so it must not claim certainty through `Denied`; `Unknown` is as fail-closed as `Denied` and more honest in the audit.
- **Red-test design**: the r1 `SandboxEnforced`-always-false wrapper was not falsifiable ⇒ r2 switched to an **environment-variable trigger plus call-site red tests** (`TestAgentRunnerRunSandboxRefusesBeforePermit` / `TestSupervisedRunSandboxRefusesBeforePermit`) that assert the permit is unclaimed, the Step 2 marker does not exist and the Dispatcher was called zero times.
- **Seven stop-and-escalate boundaries**: all seven judged **not triggered** (no contract section outside §7, no chain/adapters criterion change, no schema·fixtures·CI·candidate-version·R-08·R-09 change, no P-prime, no scope expansion, no S1 Exit gate change, no change to the existing refusal semantics).
- **TDD variant order**: this slice uses the red-first variant (③0 skeleton → ③ red tests → ②b implementation); the difference from the standard §5.1 order (② → ③ → ④) is an explicit arrangement of the slice definition.
- **⑦ blind-review-not-drawn confirmation** (recorded here): `Sandbox.Check` is a check-layer method, the method set carries no process lifecycle method, and the slice does not touch `ManagedProcess.Terminate` / `verifyGone` / `StopProcessIdentity` or their callers, so ownership and identity semantics are unchanged ⇒ the shutdown-semantics trigger does not apply; the post-P71 window is P72 / P73 / P74, three slices, so the proportional sampling also does not apply.

## 4. ②a contract revision (`CONTRACTS{,_EN}.md` §7.2)

- Landing: a new `### 7.2 Sandbox` section (CN: `### 7.2 Sandbox 检查层接口`) at the end of §7 (after 7.1, before §8) with eight clauses (a)–(h): check-layer position / verdict states / error classification / check-point option (i) / S1 trigger / platform implementations / boundaries / the part-B unlock disclaimer.
- Symmetry reading: `docs/CONTRACTS.md +38/−0` against `docs/CONTRACTS_EN.md +38/−0` (`G3(a)` PASS). Note the two baselines differ by design (CN 403 / EN 395 lines); the criterion is **equal added lines**, not equal file length.
- The trigger semantics follow the r2 correction: `PRFRAIL_SANDBOX_CHECK` ∈ {unset, `off`, `enforce`}; an invalid value is a CLI usage error with zero side effect, evaluated before any file read.

## 5. ③0 skeleton / ③ red tests / ②b implementation / mutation check

- **③0 skeleton**: `internal/guard/sandbox.go` (interface + three states + two sentinels + inert `NewSandbox()` / `CheckSandbox()`) + `internal/console/sandbox_preflight.go` (`sandboxForDispatch` injection point + inert `preflightSandbox`) + a new `SandboxEnforced bool` on both assembly option structs (zero value, zero behaviour change). `go build ./...` exit 0.
- **③ red tests** (verbatim output in `03-tests/red-*.txt`):
  - `TestSandboxInterfaceContract` red on `NewSandbox must never return nil`;
  - `TestAgentRunnerRunSandboxRefusesBeforePermit` red in all three sub-cases on `got <nil>`;
  - `TestSupervisedRunSandboxRefusesBeforePermit` red in all three sub-cases on `got supervised visible delivery is not accepted…` (the flow had already passed the check point).
  - `TestUnsupportedPlatformReturnsUnknown` (`//go:build !windows`) does not run on this host and is carried by the CI ubuntu leg; `GOOS=linux` build / vet proves it compiles.
- **②b implementation**: the two platform files (`!windows` / `windows` alpha), the real `NewSandbox()` / `CheckSandbox()`, the real `preflightSandbox`, `sandboxEnforcementFromEnv`, both call sites (ahead of the permit), and the single CLI read plus route-parameter injection. Every red test turned green (`03-tests/green-*.txt`) and the six trigger sub-cases pass.
- **Mutation check** (verbatim output in `03-tests/mutation-*.txt`): removing the check blocks of `runtime_agent.go` and `runtime_supervised.go` turns **both** call-site tests red (the former with `got <nil>`, the latter with `got supervised visible delivery is not accepted…`) ⇒ restoring turns them green, `gofmt` is clean and no backup file is left behind.
- **Equivalence witnesses**: no existing test was modified; `go test ./... -count=1` ⇒ 15 packages `ok`, `FAIL=0`.

## 6. ④ gate readings

Windows host, 2026-11-15 (④ batch; the ⑨ final re-run is in §9):

- `gofmt -l .` empty output; `go build ./...` exit 0; `go vet ./...` exit 0; `go test ./... -count=1` ⇒ ok packages 15 / FAIL 0.
- `node tools/gates/gate.js --all --scope=index` ⇒ `changed=25`, `TOTAL_FAIL=0`; `--scope=tree` ⇒ `changed=24`, `TOTAL_FAIL=0`; `--scope=ci` ⇒ `changed=2`, `TOTAL_FAIL=0`.
- `node tools/gates/gate.js --selftest` ⇒ `SELFTEST: PASS 737/737`.
- `GOOS=linux GOARCH=amd64 go build ./...` ⇒ exit 0; `go vet ./...` ⇒ exit 0 (`06-gates/linux-build.txt`).
- `G3(a)` (index scope): `docs/CONTRACTS.md +38/-0 vs _EN +38/-0 OK`; `docs/t027/REMAINING_SLICES.md +20/-0 vs _EN +20/-0 OK`.

## 7. ⑤ / ⑥ / ⑦ review-chain conclusions

- **⑤ pre-review r1**: `PRE-REVIEW: PASS` (3 Low + 3 Info). All three Low items were fixed: the reason enters the error chain, the `ErrSandboxUnavailable` self-wrapping is removed, and `internal/guard/sandbox_windows_test.go` pins the alpha contract. **⑤ r2 (re-run)**: `PRE-REVIEW-R2: PASS` (0 Critical/High/Medium, 1 Low + 2 Info).
  - Not fixed: the r2 Low is a forward-looking suggestion to keep the full `err` chain (both S1 platform implementations return the bare sentinel, so nothing is lost today); the two r2 Info items are not actionable (`ErrSandboxRejected` does not trim the reason; the env constant and its single read live in different files but remain unique).
- **⑥ independent scan**: `INDEPENDENT SCAN: PASS` (Section A–E complete, Section B no findings). Section D walks four hypotheses with a red/green verdict each and lists three uncovered invariants (see §11). One call (budget ≤3).
- **⑦ final review (blind review not drawn)**: `FINAL REVIEW: PASS` (Section B no findings; all four review items pass). The first round received only the slice definition, the diff, the contract sentences and the read-only constraint, with no ⑤ / ⑥ list. Metering = final 1 + re-review 0 + blind 0.
  - Two reservations from Section E: there is no independent runtime reading for the non-Windows platform (this host is Windows), and "an illegal environment value has zero side effect before any file read" rests mainly on static ordering. Both are recorded in §11.

## 8. ⑧ document assembly

- New report pair `docs/validation/sandbox-interface{,_EN}.md` (this file).
- Evidence pack `docs/validation/evidence/P74-SANDBOX-INTERFACE/`: `00-briefing.md`, `01-arch/design-r1.md`, `02-batch/` (2), `03-tests/` (8: red×3 / green×3 / mutation×2), `04-reviews/` (4: prereview-r1 / prereview-r2 / scan-r1 / final-r1), `06-gates/` (6: index / tree / ci / selftest / linux-build / full-test), `MANIFEST.md`, `SHA256SUMS.txt`, `freedom-list.md`, `input-list.md`.
- `docs/DEV_PLAN{,_EN}.md`: one added line each in the T027-B section (P74 landing record).
- `docs/t027/REMAINING_SLICES{,_EN}.md`: the P74 section flips 🔄 → ✅ plus a completion record line.

## 9. ⑨ native verification

(Section 9 readings are produced by the re-run that follows the ⑧ assembly; they are archived in `06-gates/` together with this table.)

- Three scopes + SELFTEST (re-run after the whole report pair and evidence pack were staged): `--scope=index` ⇒ `changed=47`, `TOTAL_FAIL=0`; `--scope=tree` ⇒ `changed=47`, `TOTAL_FAIL=0`; `--scope=ci` ⇒ `changed=2`, `TOTAL_FAIL=0`; `SELFTEST: PASS 737/737`.
- `go test ./... -count=1` exit 0; `gofmt -l .` empty output.
- `GOOS=linux GOARCH=amd64 go build/vet ./...` both exit 0.

## 10. ⑧c wrap-up

- `tmp/P74-SANDBOX-INTERFACE/` is cleaned item by item; five readings are read back: `Test-Path` = `False`; recursive file count under `tmp` (excluding `gates-cache`) = 1; top-level names = `gates-cache,.gitkeep`; untracked path count = 0; total recursive file count under `tmp` = 101; after the clean-up the three scopes read `TOTAL_FAIL=0` and `SELFTEST: PASS 737/737`.
- The `tmp` readings are read a second time between ⑧c and the ⑩ receipt (E-1 recurrence guard).

## 11. Known boundaries

- **Part B does not unlock S-G1 / S-G8**: the slice lands the check-layer abstraction only; the `S1 Exit` eight gates are unchanged and the real enforcement boundary stays in S2.
- **The three P73 conditions, item by item**: (a) interface and call sites land in one slice — `sandbox.go` + `sandbox_windows.go` + `sandbox_unsupported.go` + both assembly wirings are in this one main commit and no backfill commit splits them; (b) the dead-code-ban disposal — handled by ADR-044 + the B2 pointer + v1.69, and this slice adds a reachable entry point through the environment variable, so `SandboxEnforced` is not dead code; (c) the report states "part B does not unlock S-G1 / S-G8" — this item is (c).
- **The supervised S1 shape**: under `enforce` the `freshStart` branch is **always refused** (the platform implementation always answers `Unknown`). That is not an exception carved out for supervised runs; it is the honest semantics of "an explicit enforcement request is refused while enforcement is not wired in". The real wiring belongs to S2. The **re-entry path (not a `freshStart`) is neither checked nor consuming**: it dispatches nothing.
- **Uncovered invariants (flagged by both ⑥ Section D and ⑦ Section E)**: (1) the zero-side-effect behaviour of an illegal `PRFRAIL_SANDBOX_CHECK` value is covered only by a resolver unit test, without a whole-CLI I/O observation assertion; (2) the supervised re-entry "no check, no consumption" path has no standalone regression sample (the mutation covers `freshStart` only); (3) there is no local runtime reading for the non-Windows platform (the CI ubuntu leg carries it).
- **⑤ r2 Low left unfixed**: when an implementation returns a composite error that wraps the sentinel *and* supplies a reason, `CheckSandbox` drops the original error text; both S1 platform implementations return the bare sentinel, so nothing is lost today. It is kept as a note for later implementations.
- **Gate G9-e informational**: before staging, the `tmp/P74-SANDBOX-INTERFACE` reference in `REMAINING_SLICES{,_EN}.md` shows `not staged in the index`; the path is this slice's declared temporary-file location, and after ⑧c the ledger line remains as a historical record, which is not a violation.

## 12. Explicitly not executed

- `docs/DELIVERY_DIRECTIVE{,_EN}.md` is untouched (still v1.69).
- `docs/ADR_REGISTER{,_EN}.md` is untouched (ADR-044 from P73 already covers the direction and the boundary; this slice implements, it does not decide).
- `docs/PRODUCT_REQUIREMENTS{,_EN}.md` / `docs/TEST_STRATEGY{,_EN}.md` / `docs/S1_SUPPORT_MATRIX{,_EN}.md` / `docs/validation/s1-exit{,_EN}.md` / `docs/validation/remaining-rulings{,_EN}.md` / `docs/validation/strategic-planning-2026-10-31{,_EN}.md` / the v50–v57 snapshots are untouched.
- `docs/ARCHITECTURE{,_EN}.md` is untouched (judged not applicable — the document has no "execution-surface module" section, and `Sandbox` is a check-layer abstraction that changes no ownership / data flow / storage / transaction).
- `docs/SECURITY{,_EN}.md` is untouched (judged not applicable — the slice introduces no new threat, and SEC-14 "CLI Agent capability probe / independent gates·review" already covers this control surface; a `Sandbox` is a control, not a threat).
- `docs/validation/t027b-admission{,_EN}.md` is untouched (judged not applicable — the P73 report is closed, and P74 traceability is carried by the new `DEV_PLAN` T027-B line and by §13 of this report).
- `internal/chain/**` / `internal/adapters/**` criterion bodies, `schemas/**`, `testdata/**`, CI configuration, the candidate version and R-08 / R-09 are untouched.
- No P-prime; no scope expansion (no process lifecycle management, no tool/network control); the S1 Exit eight gates are unchanged; no new slice was started.
- No commit and no push (gitee included).

## 13. Fact pack for v58

- Slice: P74 `SANDBOX-INTERFACE` closed (layer-1 part-B introduction slice, implementation, tier M).
- Delivered shape: the `Sandbox` check layer (`Check(ctx, CheckSpec) (Verdict, error)`; three states `Allowed` / `Denied` / `Unknown`, the latter two fail-closed) + two assembly check points ahead of the permit + `PRFRAIL_SANDBOX_CHECK` (S1 default `off`).
- Contract: `CONTRACTS{,_EN}.md` §7.2 with eight clauses (a)–(h), CN/EN `+38/−0` each.
- Gates: three scopes `TOTAL_FAIL=0`, SELFTEST 737/737, `GOOS=linux` build/vet exit 0.
- Review: ⑤ PASS (two rounds), ⑥ PASS (one call, the mandatory layer executed), ⑦ PASS (final 1, blind not drawn).
- Counts: R = 3 (two in session 1, one in session 2), D = 0, OB-139 not repeated (D=0); **R=3 trips the §9.2 "high-risk slice" marker for the first time and the §5.2 high-risk row was executed** (P72 / P73 had R=0 and did not trip it); the batch-assist re-dispatch is process execution and is not counted into R.
- Why the ⑦ blind review was not drawn: none of the three mandatory semantics is touched (the check layer carries no lifecycle method and changes no shutdown path) and the proportional window holds three slices (P72 / P73 / P74), fewer than four.
- Layer-1 part-B status: the abstraction has landed; the real enforcement boundary stays in S2; S-G1 / S-G8 remain locked.
- Segmentation fact: P74 spans two sessions under the §13.1.17 lenient reading (acknowledged by the user on 2026-11-15), used for the first time.

## 14. Commit and CI fact chain (backfill)

- Main commit: `3300364` (`feat: add the sandbox check layer before the one-shot probe permit`; one commit carrying the whole staged set of this slice); pushed as `f834ce5..3300364 main -> main` (origin only).
- CI (main commit, run `43669748483`, head `3300364`): **both legs success** — `Go windows-latest` success (4.9 min) and `Go ubuntu-latest` success (3.35 min); the other 3 jobs are `skipped` (not triggered by this slice).
- Package-level evidence (raw `go test` `ok` lines): windows `ok internal/console 21.010s` / `ok internal/guard 3.106s`; ubuntu `ok internal/console 3.939s` / `ok internal/guard 0.992s`; ubuntu Race leg `ok internal/console 15.224s` / `ok internal/guard 2.992s`. The green ubuntu `Test` leg means the `//go:build !windows` test `TestUnsupportedPlatformReturnsUnknown` passed on a real Linux host (i.e. the §11 boundary "the non-Windows runtime reading is carried by the CI ubuntu leg" is now discharged).
- Backfill note: this section is written by a docs-only backfill commit (still pushed to origin only); **the backfill commit's own run id is an observation record only (§11.2 self-reference avoidance) and is not the authoritative CI reading of this slice**.

The following is the archived-entry list of this slice (repository-relative path + disposition, TAB separated; the CN and `_EN` fences are byte-identical).

```artifacts
docs/validation/sandbox-interface.md	repo
docs/validation/sandbox-interface_EN.md	repo
docs/validation/evidence/P74-SANDBOX-INTERFACE/00-briefing.md	repo
docs/validation/evidence/P74-SANDBOX-INTERFACE/01-arch/design-r1.md	repo
docs/validation/evidence/P74-SANDBOX-INTERFACE/02-batch/batch-01-interface-enumeration.md	repo
docs/validation/evidence/P74-SANDBOX-INTERFACE/02-batch/batch-02-contracts-structure.md	repo
docs/validation/evidence/P74-SANDBOX-INTERFACE/03-tests/red-guard-contract.txt	repo
docs/validation/evidence/P74-SANDBOX-INTERFACE/03-tests/red-agent-runner.txt	repo
docs/validation/evidence/P74-SANDBOX-INTERFACE/03-tests/red-supervised.txt	repo
docs/validation/evidence/P74-SANDBOX-INTERFACE/03-tests/green-guard-contract.txt	repo
docs/validation/evidence/P74-SANDBOX-INTERFACE/03-tests/green-agent-runner.txt	repo
docs/validation/evidence/P74-SANDBOX-INTERFACE/03-tests/green-supervised.txt	repo
docs/validation/evidence/P74-SANDBOX-INTERFACE/03-tests/mutation-agent-runner.txt	repo
docs/validation/evidence/P74-SANDBOX-INTERFACE/03-tests/mutation-supervised.txt	repo
docs/validation/evidence/P74-SANDBOX-INTERFACE/04-reviews/prereview-r1.md	repo
docs/validation/evidence/P74-SANDBOX-INTERFACE/04-reviews/prereview-r2.md	repo
docs/validation/evidence/P74-SANDBOX-INTERFACE/04-reviews/scan-r1.md	repo
docs/validation/evidence/P74-SANDBOX-INTERFACE/04-reviews/final-r1.md	repo
docs/validation/evidence/P74-SANDBOX-INTERFACE/06-gates/gate-index.txt	repo
docs/validation/evidence/P74-SANDBOX-INTERFACE/06-gates/gate-tree.txt	repo
docs/validation/evidence/P74-SANDBOX-INTERFACE/06-gates/gate-ci.txt	repo
docs/validation/evidence/P74-SANDBOX-INTERFACE/06-gates/gate-selftest.txt	repo
docs/validation/evidence/P74-SANDBOX-INTERFACE/06-gates/linux-build.txt	repo
docs/validation/evidence/P74-SANDBOX-INTERFACE/06-gates/full-test.txt	repo
docs/validation/evidence/P74-SANDBOX-INTERFACE/MANIFEST.md	repo
docs/validation/evidence/P74-SANDBOX-INTERFACE/SHA256SUMS.txt	repo
docs/validation/evidence/P74-SANDBOX-INTERFACE/freedom-list.md	repo
docs/validation/evidence/P74-SANDBOX-INTERFACE/input-list.md	repo
docs/CONTRACTS.md	repo
docs/CONTRACTS_EN.md	repo
docs/DEV_PLAN.md	repo
docs/DEV_PLAN_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
internal/guard/sandbox.go	repo
internal/guard/sandbox_windows.go	repo
internal/guard/sandbox_unsupported.go	repo
internal/guard/sandbox_interface_test.go	repo
internal/guard/sandbox_windows_test.go	repo
internal/guard/sandbox_unsupported_test.go	repo
internal/console/sandbox_preflight.go	repo
internal/console/runtime_agent.go	repo
internal/console/runtime_supervised.go	repo
internal/console/cli.go	repo
internal/console/supervised_cmd.go	repo
internal/console/runtime_agent_sandbox_test.go	repo
internal/console/runtime_supervised_sandbox_test.go	repo
```

The following is the verbatim Section A–E of the ⑥ independent scan r1 (transcribed from `docs/validation/evidence/P74-SANDBOX-INTERFACE/04-reviews/scan-r1.md`); its summary sentence reads "依据契约与主控落盘静态材料，本片满足 §7.2 的接口、时序、fail-closed 语义与可证伪性要求".

```sixstage
## Section A

INDEPENDENT SCAN: PASS

## Section B

无发现

## Section C

- 读取主控落盘的 [docs/CONTRACTS.md](docs/CONTRACTS.md) 与 [docs/CONTRACTS.md](docs/CONTRACTS.md#L321-L359)：§7.2 的八条契约彼此一致，且与实现对齐：
  - `guard.Sandbox` 只做 `Check(ctx, CheckSpec) (Verdict, error)` 判定，不提供生命周期方法，符合 (a)。
  - `Allowed / Denied(reason) / Unknown(reason)` 三态都按 fail-closed 处理，符合 (b)。
  - `ErrSandboxUnavailable` / `ErrSandboxRejected` 均不得被解释为放行，符合 (c)。
  - 检查点放在 probe permit 认领之前，且在 run-root marker 之前，符合 (d)。
  - 触发源仅是 `PRFRAIL_SANDBOX_CHECK`，默认 `off`，符合 (e)。
  - 非 Windows / Windows S1 都返回 `Unknown`，不允许谎报 `Denied`，符合 (f)。
  - 边界不碰 chain/adapters 判定本体，也不落真实 enforcement，符合 (g)。
  - B 段不解锁 S-G1/S-G8，符合 (h)。

- 读取主控落盘的 [internal/guard/sandbox.go](internal/guard/sandbox.go) 与 [internal/guard/sandbox_unsupported.go](internal/guard/sandbox_unsupported.go)：`CheckSandbox` 统一把 `nil` / 非 `Allowed` / `err != nil` 全部包装成 fail-closed 语义；`Unsupported` 平台返回 `VerdictUnknown + ErrSandboxUnavailable`，与契约一致。

- 读取主控落盘的 [internal/console/runtime_agent.go](internal/console/runtime_agent.go) 与 [internal/console/runtime_supervised.go](internal/console/runtime_supervised.go)：两条派发路径都在 permit 之前调用 `preflightSandbox`，且在 run-root marker / operator ledger 之前返回，符合“环境/配置类拒绝不得烧配额、不建 marker、不 spawn”的时序。

- 读取主控落盘的 [tmp/P74-SANDBOX-INTERFACE/03-tests/red-guard-contract.txt](tmp/P74-SANDBOX-INTERFACE/03-tests/red-guard-contract.txt)、[tmp/P74-SANDBOX-INTERFACE/03-tests/red-agent-runner.txt](tmp/P74-SANDBOX-INTERFACE/03-tests/red-agent-runner.txt)、[tmp/P74-SANDBOX-INTERFACE/03-tests/red-supervised.txt](tmp/P74-SANDBOX-INTERFACE/03-tests/red-supervised.txt) 与 [tmp/P74-SANDBOX-INTERFACE/03-tests/green-guard-contract.txt](tmp/P74-SANDBOX-INTERFACE/03-tests/green-guard-contract.txt)、[tmp/P74-SANDBOX-INTERFACE/03-tests/green-agent-runner.txt](tmp/P74-SANDBOX-INTERFACE/03-tests/green-agent-runner.txt)、[tmp/P74-SANDBOX-INTERFACE/03-tests/green-supervised.txt](tmp/P74-SANDBOX-INTERFACE/03-tests/green-supervised.txt)：红测先红、绿测后绿，说明 TDD 方向正确，且检查点确实位于 permit 之前。

- 读取主控落盘的 [tmp/P74-SANDBOX-INTERFACE/03-tests/mutation-agent-runner.txt](tmp/P74-SANDBOX-INTERFACE/03-tests/mutation-agent-runner.txt) 与 [tmp/P74-SANDBOX-INTERFACE/03-tests/mutation-supervised.txt](tmp/P74-SANDBOX-INTERFACE/03-tests/mutation-supervised.txt)：删掉调用点后，两个测试都从绿转红；这说明调用点不是“无意义装饰”，而是真实的安全断言。

- 读取主控落盘的 [tmp/P74-SANDBOX-INTERFACE/gate-index.txt](tmp/P74-SANDBOX-INTERFACE/gate-index.txt)、[tmp/P74-SANDBOX-INTERFACE/gate-tree.txt](tmp/P74-SANDBOX-INTERFACE/gate-tree.txt)、[tmp/P74-SANDBOX-INTERFACE/gate-ci.txt](tmp/P74-SANDBOX-INTERFACE/gate-ci.txt)、[tmp/P74-SANDBOX-INTERFACE/gate-selftest.txt](tmp/P74-SANDBOX-INTERFACE/gate-selftest.txt) 与 [tmp/P74-SANDBOX-INTERFACE/04-full-test.txt](tmp/P74-SANDBOX-INTERFACE/04-full-test.txt)：门禁四项均为 PASS，且全量测试输出里 `internal/guard` / `internal/console` 均为 `ok`，未见逆转证据。

## Section D

1. 删掉 [internal/console/runtime_agent.go](internal/console/runtime_agent.go) 的检查块
   - 会红：`TestAgentRunnerRunSandboxRefusesBeforePermit` 的三个子case（denied / unknown / unavailable）要求 `ErrSandboxRejected` 或 `ErrSandboxUnavailable`，而 mutation 输出中 [tmp/P74-SANDBOX-INTERFACE/03-tests/mutation-agent-runner.txt](tmp/P74-SANDBOX-INTERFACE/03-tests/mutation-agent-runner.txt) 显示 “got <nil>” 和 “the dispatch must be refused with ...”。
   - 不会红：仅看“CLI 能起”或“build 能过”的测试；这类断言不覆盖“在 permit 前拒绝并不烧配额”的安全语义。
   - 结论：删掉检查块后，安全断言会红，且这不是“假红”。

2. 删掉 [internal/console/runtime_supervised.go](internal/console/runtime_supervised.go) 的检查块
   - 会红：`TestSupervisedRunSandboxRefusesBeforePermit` 要求在 freshStart 分支中先被 sandbox 拒绝；mutation 输出 [tmp/P74-SANDBOX-INTERFACE/03-tests/mutation-supervised.txt](tmp/P74-SANDBOX-INTERFACE/03-tests/mutation-supervised.txt) 明确显示 “got supervised visible delivery is not accepted …”，说明代码已越过安全门而进入 visible delivery。
   - 不会红：“重入路径不检查、不消费配额” 的正向路径（非 freshStart）不会被这条测试覆盖到；它被设计为仅在 freshStart 上检查。
   - 结论：删掉该检查块后，fresh-start 的安全断言会红；“再入路径不检查”本身仍需单独证明。

3. 删掉 [internal/guard/sandbox_unsupported.go](internal/guard/sandbox_unsupported.go) 的 fail-closed 返回
   - 会红：`TestSandboxInterfaceContract` 明确要求 `NewSandbox()` 不为 nil，并要求 `CheckSandbox(..., nil, ...)` 以 fail-closed 拒绝；[tmp/P74-SANDBOX-INTERFACE/03-tests/red-guard-contract.txt](tmp/P74-SANDBOX-INTERFACE/03-tests/red-guard-contract.txt) 也显示 “NewSandbox must never return nil: a nil default would turn every check into a silent pass”。
   - 不会红：单纯 `go build ./...` 没有覆盖 fail-closed 安全语义；编译不证明“拒绝”的正确性。
   - 结论：删掉 fail-closed 返回，就是把默认实现从“安全拒绝”变成“静默放行/空路径”。

4. 把 `VerdictUnknown` 改成 `Allowed`
   - 会红：契约测试在 [internal/guard/sandbox_interface_test.go](internal/guard/sandbox_interface_test.go) 明确断言 “an Unknown verdict must not be allowed”；这条用例会立刻红。
   - 也会红：`TestAgentRunnerRunSandboxRefusesBeforePermit/unknown` 与 `TestSupervisedRunSandboxRefusesBeforePermit/unknown` 的语义会被破坏；它们原本依赖 “Unknown = fail-closed”，一改成 `Allowed`，运行路径就不再拒绝。
   - 不会红：只有“Happy path 允许 dispatch”的测试；它们并不会验证 fail-closed 性。
   - 结论：`Unknown` 是安全语义而非允许值，改成 `Allowed` 会直接破坏约束。

没有覆盖的不变量（我认为未被本次静态材料完全覆盖）：
- `PRFRAIL_SANDBOX_CHECK` 的非法值到真正 CLI 层的零副作用，只有 resolver 单元测试覆盖；未看到全链路 CLI 证据。
- `supervised` 的“非 freshStart 再入路径不检查、不消费配额”没有独立回归样本；现有 mutation 只覆盖 fresh-start 分支。
- 真实 Windows 主机上的实际行为仍未直接验证；代码和日志都证明了 S1 的平台语义，但没有在 Windows 真实环境上复现。

## Section E

- 我无法验证的部分：
  - 实机 Windows 行为：当前静态材料没有真实 Windows 运行证据；代码上只看到 `//go:build windows` 的 α 形态返回 `Unknown`，但没有真实宿主验证。
  - 真实探针配额/脚本层的跨进程竞争：材料证明了单元/白盒行为，但没有直接复现并发抢占后一个 claim 是否永远不被“悄悄放行”。

- 什么会推翻我的 PASS：
  - 任何一个检查点移到 permit 之后或 run-root marker 之后；
  - 任何实现把 `Unknown` / `ErrSandboxUnavailable` 当成 `Allowed`；
  - 删掉 [internal/guard/sandbox_unsupported.go](internal/guard/sandbox_unsupported.go) 的 fail-closed 返回，或让默认实现返回 nil；
  - `PRFRAIL_SANDBOX_CHECK` 的非法值被接受而不报 CLI usage error；
  - [tmp/P74-SANDBOX-INTERFACE/03-tests/mutation-agent-runner.txt](tmp/P74-SANDBOX-INTERFACE/03-tests/mutation-agent-runner.txt) 与 [tmp/P74-SANDBOX-INTERFACE/03-tests/mutation-supervised.txt](tmp/P74-SANDBOX-INTERFACE/03-tests/mutation-supervised.txt) 不再从绿转红；
  - 真实环境中出现“非 Windows 平台返回 Denied”或“Windows S1 返回 Denied”而不是 `Unknown`。

INDEPENDENT SCAN: PASS
```
