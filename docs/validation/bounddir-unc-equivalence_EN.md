# P86 `BOUNDDIR-UNC-EQUIVALENCE` Slice Report

[简体中文](bounddir-unc-equivalence.md)

**Verdict: `CLOSED`** (both N13-UNC and W10 closed; host and gate readings all green, the CI two-leg readings are in section 11.2). This slice closes only the last two items of F1; it does not touch the other already-closed F1 items and does not touch F3 / F6 / F8 (the P85 scope).

## 1. Slice identity and authorization

| Item | Content |
|---|---|
| Number / identifier | P86 / `BOUNDDIR-UNC-EQUIVALENCE` |
| Class | `[SLICE]` - implementation - S tier |
| Authorization | The user named the start on 2026-11-22 (P86 before P85); the verbatim start message is `evidence/P86-BOUNDDIR-UNC-EQUIVALENCE/00-briefing.md` |
| Basis | P84 report section 8 deferral list (N13-UNC plus the evidence-attribution Low) and the v68 snapshot section 4 queue item |
| Boundary | Fix N13-UNC and W10 only; no change to the gate table / the directive / CI / the threat model / dependencies |

## 2. Step 0 anchor readings

| # | Item | Expected | Measured |
|---|---|---|---|
| 1 | HEAD | `c7ebf3d` | `c7ebf3d` |
| 2 | origin/main | `c7ebf3d` | `c7ebf3d` |
| 3 | Directive version | v1.71 | v1.71 |
| 4 | SELFTEST | 740/740 | 740/740 |
| 5 | Three scopes | all `TOTAL_FAIL=0` | index / tree / ci all 0 |
| 6 | v66 line count (LF) | 216 | 216 (CN / EN symmetric) |
| 7 | v67 line count (LF) | 216 | 216 (CN / EN symmetric) |
| add | v68 line count (LF) | 218 | 218 (CN / EN symmetric) |
| add | Mechanical readings | ADR 47 / OB 139 / DR 17 / pilot deferral 25 | all match |

**Registration (not counted as drift)**: the P86 ledger section pre-registered "(7) `GPT-6.2 Sol` deep tier + 1 final review + 3 re-reviews + mandatory blind review", which conflicts with the user's 2026-11-22 start message ("`GPT-5.3-Codex` regular tier + 1 final + 1 re-review + no budget exception"). Per section 0.1 / section 13.2.11 discipline F the start message governs, and the ledger was corrected in place. The R/D baseline is 0 / 0.

## 3. Goal and scope

- **N13-UNC** (F1 residual, Medium): the two canonical spellings of a UNC share root (`\\host\share` and `\\host\share\`) were falsely refused by the `boundRelativeSuffix` equality judgement - a deterministic in-bound false rejection in the fail-closed direction, an availability defect.
- **W10** (evidence-attribution Low): the mutation annotation in `p84_bounddir_residual_test.go` called the double-separator prefix `D:\\` "the round-11-before shape", contradicting (7) re-review r3 (r3 records the single-separator prefix `D:\`).

## 4. Implementation and closure status

| Item | Owner | Change | Status |
|---|---|---|---|
| N13-UNC | F1 | `internal/evidence/bounddir_windows.go` `boundRelativeSuffix`: a new equality branch (after step 3, before the strict-descendant prefix) plus the helpers `uncShareRootName` / `uncShareRootEqual` / `lowerASCII` | **closed** |
| W10 | evidence attribution | the `Mutation check:` comment block of `internal/evidence/p84_bounddir_residual_test.go` (comment only, zero code change) | **closed** |
| Tests | - | new `internal/evidence/p86_bounddir_unc_test.go` (242 lines; 5 positives / 12 refusals / 10 classifier negatives / 1 production-entry case) | all green (the production entry SKIPs on this host, see section 11) |

**N13-UNC decision order (after the fix)**: (1) non-canonical-spelling guard, (2) `IsAbs`, (3) canonical-form identity, (4) **new: both are UNC share roots and their volume names are ASCII-fold equal => admit**, (5) empty trimmed base refuses, (6) strict-descendant Ordinal prefix plus the component guard. **The only new admission point is step (4)**.

**Zero-relaxation argument**: a different host / a different share name => the folded volume names differ; host-only (`\\host`, `\\host\`) => not two non-empty components => the classifier answers `""`; `\\?\` / `\\.\` => the classifier answers `""` explicitly; non-canonical spellings (`/`, repeated separators, a non-root trailing separator, `.`/`..`) => step (1) refuses first; an out-of-bound descendant => the step (6) prefix does not match; the drive form => the volume name has no `\\` prefix; case variants of descendant components => step (6) stays Ordinal and refuses.

**W10 correction**: M1 (revert the equality only, prefix still built from the trimmed base => `D:\`) and M2 (additionally build the prefix from the untrimmed form => `D:\\`) are annotated separately; the wrong "`D:\\` = the round-11-before shape" wording is gone and replaced by "NOT the round-11-before arrangement"; the text is verbatim-compatible with r3's "the target equal to that prefix".

## 5. Contract impact judgement

| Item | Verdict | Reason |
|---|---|---|
| N13-UNC | **DOES-NOT-TOUCH / IN-BOUNDARY** | `docs/CONTRACTS{,_EN}.md` carries no bounddir / reparse / UNC-specific sentence; the nearest sentence (`CONTRACTS.md:238`) belongs to the replay-store writer-root/protected-root discipline and does not constrain this file. The change keeps the "Ordinal byte comparison" discipline and only adds one share-root equality shape |

No new schema, no new CLI, no new dependency (only the existing `strings` / `filepath`). **The contract is untouched => no protocol-first step, and `CONTRACTS{,_EN}.md` / schema / fixtures / `go.mod` are unchanged**.

## 6. Domain escalation and user rulings

| Item | Ruling | Record |
|---|---|---|
| Main control | stays `DeepSeek V4.1 Flash (deepseek)` High (P84's V4 Pro is not inherited) | start message section 0 (the 2026-10-31 long-standing default) |
| (7) budget / tier | 1 final review + 1 re-review; `GPT-5.3-Codex (copilot)` regular tier; no budget exception | user ruling a of 2026-11-22 |
| Blind review | mandatory once (counted separately as +1, executed before the final review) | (1) judgement (bounddir touches identity/ownership semantics); recorded in plan section 7 |
| **Case-variant conflict** | **executed per start message section 6**: ASCII case differences of a UNC share-root volume name are judged equal | Four elements: decider = user (same-turn start message); date = 2026-11-22; reason = the start message is newer and more specific and (7)'s "keep Ordinal" is a P84 leftover suggestion acting as a lower bound; basis = start message section 6 plus directive section 0.1; **may not be cited as a precedent for other slices** |
| Scheduling | P86 before P85 | user ruling b of 2026-11-22 |

**Safety argument for the ASCII fold**: the UNC volume namespace is defined by SMB, and host and share names are case-insensitive there, so one host cannot carry two shares differing only by case => the fold **does not widen** the admitted object set, it only removes false refusals; the fold is **volume-name only and ASCII only** (`A`-`Z`); descendant components, the drive-form root and `isUnderDirectory` stay Ordinal (W1's case-variant-twin refusal is unchanged).

## 7. Review chain

| Stage | Carrier | Model | Budget | Verdict |
|---|---|---|---|---|
| (1) plan | `deep-reasoner` | `DeepSeek V4 Pro (deepseek)` | 1 | contract judgement + design + conflict ruling + blind-review decision |
| (2) implementation | `prfrail-implementer` | `DeepSeek V4.1 Flash (deepseek)` | 1 | two changes + full regression |
| (3) tests | `prfrail-tester` | `DeepSeek V4.1 Flash (deepseek)` | 1 | new tests + mutation run |
| (5) pre-review | `independent-reviewer` | `DeepSeek V4 Pro (deepseek)` | 1 + 1 | r1 `PASS WITH FIXES` (1 Medium + 4 Low) -> r2 `PASS` |
| (6) scan | `independent-reviewer` | `MAI Code 2 Flash (copilot)` | <= 3 (used 2) | r1 `FINDINGS` (2 Low) -> r2 `PASS` |
| (7) blind | `independent-reviewer` | `GPT-5.3-Codex (copilot)` | +1 separate | `PASS WITH FIXES` (the only Low = no non-SKIP production-entry evidence) |
| (7) final | `independent-reviewer` | `GPT-5.3-Codex (copilot)` | 1 (first of 1+1) | `PASS WITH FIXES` (no Medium+; Low-L1 = the same production-entry evidence) |

The verbatim transcripts are archived in `04-reviews/` (pre-review r1/r2, scan r1/r2, blind, final r1).

## 8. Residuals and deferrals

- No new deferral in this slice. **F1 is closed after this slice** (the T4 residual N13-UNC and the evidence-attribution W10 are both closed here).
- The P85 scope (F3 / F6 / F8: T2 / T3 / T8a) is untouched and needs the user to name its start.
- Pending user ruling: whether `bounddir-residual{,_EN}.md` section 8 gets a "N13-UNC closed by P86" timestamp note (start message section 7 item 7 lists it as a check item; this slice does not add it unilaterally).

## 9. Documentation sync

| File | Action |
|---|---|
| `docs/validation/bounddir-unc-equivalence{,_EN}.md` | this report (new) |
| `docs/t027/REMAINING_SLICES{,_EN}.md` | P86 section: status -> `closed` plus the completion record; the (7) budget line corrected per the start message; the P84 section's "deferred to P86" wording replaced with "closed by P86" |
| `docs/SECURITY{,_EN}.md:5` | remediation note becomes "F1 fully closed (the T4 residual N13-UNC closed by P86)"; the F2/F5/F7 fixed sentence and the F3/F6/F8-to-P85 sentence stay |
| `docs/DEV_PLAN{,_EN}.md` | a P86 paragraph appended to the T027 section (conclusion and report path only; the T027 status field is not rewritten) |
| `docs/t027/CHECKLIST_v68{,_EN}.md` | **untouched** (v69 absorbs it) |
| Directive / ADR / OB / gate table / CI | **untouched** (no version bump, no new entries) |

## 10. Gate results

| Gate | Result |
|---|---|
| `go build ./...` | EXIT=0 |
| `go vet ./...` | EXIT=0 |
| `gofmt -l internal/evidence` | empty |
| `go test ./... -count=1` | 15 packages ok, 0 FAIL (`internal/evidence` ok) |
| `gate.js --selftest` | `SELFTEST: PASS 740/740` |
| `gate.js --all --scope=tree` / `--scope=ci` | `TOTAL_FAIL=0` |
| `gate.js --all --scope=index` | 1 before the commit (G5(f) index hygiene on a dirty tree, expected); 0 after the (8)b commit (see section 11.2) |

Logs: `03-verification/build.txt`, `vet.txt`, `gate-selftest.txt`, `gate-{index,tree,ci}.txt`, `test-all-windows.txt`, `test-p86-verbose.txt`.

## 11. (9) verification readings and host-side non-executables

- **Windows host**: `go test ./internal/evidence -run TestP86 -count=1 -v` -> `TestP86N13UncShareRootEquivalence` (5 subcases) PASS, `TestP86N13UncShareRootRefusals` (12 subcases + 10 classifier negatives) PASS, `TestP86UncShareRootProductionEntry` **SKIP** (this host is not elevated: `net share ... exit status 2 (system error 5, access denied)`; the skip reason is spelled out in `t.Skipf`).
- **CI guard (the mechanical closure of the (7) blind/final Low)**: when the `CI` environment variable is set, the production-entry case turns three precondition failures (share creation, opening the UNC root handle, the trailing-separator handle spelling) from a skip into a **hard failure**, so **a green CI proves the production entry really ran and passed**; on this host (`CI` unset) it still skips with the reason. This does not change the CI workflow (inside the boundary).
- **Host flake observation (not introduced here)**: 1 of 3 local `go test ./... -count=1` runs failed `TestP39DurablePrecommitFaultsRejectAck/dirClose` in `internal/tickets` with `rename ... Access is denied` (the same environmental IO/AV jitter family as DR-17); the other 2 runs were green and neither CI leg reproduced it, so it is recorded as an observation-level flake and no DR is opened.
- **Mutation check (run and recorded)**: `mutation-record.txt` records `BEFORE_SHA256=35df2761...dba62` / `ANCHOR_COUNT=1` / `MUTATED_SHA256=f6ab720f...64db7` / `MUTATED_TEST_EXIT=1` / `RESTORED_SHA256=35df2761...dba62` / `RESTORED_IDENTICAL=True` / `RESTORED_TEST_EXIT=0`; `mutation-branch-removed.txt` shows **exactly** 3 spelling-difference subcases FAIL (`trailing-base-bare-target`, `bare-base-trailing-target`, `ascii-case-folded-trailing-target`) while `identical-trailing-spellings`, `strict-descendant-regression`, every refusal and `TestP84W10...` stay PASS.
- **No `-race`**: this host has no gcc (same caliber as P84); the Race step is carried by the CI ubuntu leg.
- **The Ubuntu VM is not started** (start message section 3); the Linux half is untouched and is regressed by the CI ubuntu leg.

## 11.2 CI readings (written after the (8)b push)

- Main commit `f1ae821`; first run **`44962802936`**, **success on both legs** (`Go windows-latest` / `Go ubuntu-latest`); windows leg: `ok github.com/larsonzh/prfrail/internal/evidence`, `SELFTEST: PASS 740/740`, `TOTAL_FAIL=0`; log payloads `03-verification/ci/run-44962802936-windows.log` / `-ubuntu.log`.
- The first run predates the CI guard (no `-v`, so no PASS/SKIP lines); the guard round (precondition failures become hard failures under CI) was rerun with the fix commit.
- Guard round first run **`44965434402`** (commit `a5ef122`): the **windows leg went red** (`TestP86UncShareRootProductionEntry` precondition: the CI host derives the share root **without** the trailing separator, `\\localhost\<share>`, while the fixture assumed the separator form), the ubuntu leg succeeded => registered as **DR-18**.
- Fix: the fixture now takes the OTHER canonical spelling as the stored target based on the observed handle spelling (trailing handle => bare target; bare handle => trailing target), so both host shapes enter the N13 equality branch.
- Fix round run **`44967047189`** (commit `1eb63d9`): **success on both legs** (`Go windows-latest` / `Go ubuntu-latest`); windows leg `ok github.com/larsonzh/prfrail/internal/evidence`, `SELFTEST: PASS 740/740`, `base go test ./... :: ok=15 FAIL=0`, `TOTAL_FAIL=0`; ubuntu leg `ok internal/evidence`, `SELFTEST 740/740`, `TOTAL_FAIL=0`, Race step success.
- **Production entry closed**: under CI the guard turns the three precondition failures into hard failures, so a green windows leg proves `TestP86UncShareRootProductionEntry` **really ran and passed** (the N13 equality branch admitted the target through the production entry and the bytes were read back unchanged).

## 12. Environment and tiers

- Host: Windows (local), Go toolchain 1.22+ (measured go1.27.0 locally; CI follows the pinned workflow version).
- Main control / (2)(3): `DeepSeek V4.1 Flash` High (same tier; role separation is guaranteed by separate dispatches); (1)(5): `DeepSeek V4 Pro`; (6): `MAI Code 2 Flash`; (7): `GPT-5.3-Codex` regular tier.
- Network: no proxy (SSH push / `gh` direct); no gitee.

## 13. Review cost and metering

- **R = 0** (remediation re-runs inside the review loop are not counted; the (5)/(6) remediations were same-loop fixes; (7) had no remediation round).
- **D = 0** (all temp files stayed under `tmp/P86-BOUNDDIR-UNC-EQUIVALENCE/` and are removed at (8)c; no in-slice violation).
- Probes = 0 (no external real call); the GLM batch assistant was not used (this is an S slice with no enumeration/mirroring batch need; the planned >= 2 versus measured 0 difference is registered in section 15).

## 14. Explicitly not executed

No gitee; no directive version bump; no new ADR / OB; no CI or gate-table change; no `go.mod` / schema / CLI change; no `-race` (no local gcc); no Linux-only case run on a real Linux host (carried by the CI ubuntu leg); no `bounddir-residual` section 8 timestamp note (pending the user's ruling).

## 15. (8)a self-check

1. Is the start message landed verbatim and registered in MANIFEST / SHA256SUMS (section 13.1.10)? Yes.
2. Are the reviewer transcripts archived verbatim (section 13.1.13)? Yes (6 files).
3. Is the payload three-way equal (disk / MANIFEST / SHA256SUMS)? Yes (rebuilt at (8)a and checked by G8-b).
4. Does the change surface match the boundary (zero touch on F3/F6/F8, the gate table, the directive, CI, dependencies)? Yes.
5. Are the meter and budgets honest (R = 0, D = 0, (6) used 2 of <= 3, (7) used 1 final + 1 blind, the re-review budget unused)? Yes.
6. Are the non-executed items disclosed one by one (section 14)? Yes.
7. **Deviation register**: (a) the plan's "GLM batch assistant >= 2" was not used (no fitting task; none was invented to pad the count); (b) the ledger's (7) pre-registration conflicted with the start message and was corrected (section 2).

## 16. Six-stage self-check (sixstage)

```sixstage
## Section A Overall verdict
CLOSED (both N13-UNC and W10 closed; F1 is closed after this slice)

## Section B Finding list

| SEVERITY | Location | Disposition |
|---|---|---|
| Medium | N13-UNC (the boundRelativeSuffix equality branch) | closed: a UNC share-root equality branch plus helpers; the mutation run proves the branch is necessary |
| Low | W10 (the p84 test mutation comment) | closed: M1/M2 annotated separately, the "D:\\ = the round-11-before shape" wording removed |

## Section C Per-item evidence and reasoning
- Decision order: non-canonical guard -> IsAbs -> canonical-form identity -> **UNC share-root equality (the only new admission point)** -> empty trimmed base refuses -> strict-descendant Ordinal prefix plus the component guard.
- 5 positives (both trailing-separator directions / ASCII fold / identity / descendant regression) and 12 refusals (other host / other share / descendant case variant / repeated separator / non-root trailing separator / forward slash / `..` / out of bound / host-only x2 / trailing blank) are pinned by `p86_bounddir_unc_test.go`; the 10 classifier negatives include `\\?\` and `\\.\`.
- Mutation: removing the branch turns exactly the 3 spelling-difference subcases red and the restore is byte-identical in SHA256 (`03-verification/mutation-record.txt`).

## Section D Falsifiability audit
- Remove the equality branch => the 3 positives go red while the refusals and the P84 W10/W1 cases stay green (run and recorded).
- Widen the volume-name fold to descendants => `case-variant-descendant-stays-ordinal` turns red.
- Relax the canonical-form guard => `repeated-separator` / `parent-component` / `non-root-trailing-separator` turn red.

## Section E What could not be verified
- The host production entry SKIPs (not elevated, `net share` access denied); the CI windows leg reading is in section 11.2.
- No `-race` was run (no local gcc); the Linux half is untouched and is regressed by the CI ubuntu leg.

INDEPENDENT SCAN: FINDINGS
```

## 17. Artifacts fence

```artifacts
internal/evidence/bounddir_windows.go	repo
internal/evidence/p84_bounddir_residual_test.go	repo
internal/evidence/p86_bounddir_unc_test.go	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/SECURITY.md	repo
docs/SECURITY_EN.md	repo
docs/DEV_PLAN.md	repo
docs/DEV_PLAN_EN.md	repo
docs/validation/bounddir-unc-equivalence.md	repo
docs/validation/bounddir-unc-equivalence_EN.md	repo
docs/validation/evidence/P86-BOUNDDIR-UNC-EQUIVALENCE/MANIFEST.md	repo
docs/validation/evidence/P86-BOUNDDIR-UNC-EQUIVALENCE/SHA256SUMS.txt	repo
docs/validation/evidence/P86-BOUNDDIR-UNC-EQUIVALENCE/00-briefing.md	repo
docs/validation/evidence/P86-BOUNDDIR-UNC-EQUIVALENCE/input-list.md	repo
docs/validation/evidence/P86-BOUNDDIR-UNC-EQUIVALENCE/freedom-list.md	repo
docs/validation/evidence/P86-BOUNDDIR-UNC-EQUIVALENCE/01-arch/01-arch-design.md	repo
docs/validation/evidence/P86-BOUNDDIR-UNC-EQUIVALENCE/03-verification/build.txt	repo
docs/validation/evidence/P86-BOUNDDIR-UNC-EQUIVALENCE/03-verification/vet.txt	repo
docs/validation/evidence/P86-BOUNDDIR-UNC-EQUIVALENCE/03-verification/gate-selftest.txt	repo
docs/validation/evidence/P86-BOUNDDIR-UNC-EQUIVALENCE/03-verification/gate-index.txt	repo
docs/validation/evidence/P86-BOUNDDIR-UNC-EQUIVALENCE/03-verification/gate-tree.txt	repo
docs/validation/evidence/P86-BOUNDDIR-UNC-EQUIVALENCE/03-verification/gate-ci.txt	repo
docs/validation/evidence/P86-BOUNDDIR-UNC-EQUIVALENCE/03-verification/test-all-windows.txt	repo
docs/validation/evidence/P86-BOUNDDIR-UNC-EQUIVALENCE/03-verification/test-p86-verbose.txt	repo
docs/validation/evidence/P86-BOUNDDIR-UNC-EQUIVALENCE/03-verification/probe-clean-windows.txt	repo
docs/validation/evidence/P86-BOUNDDIR-UNC-EQUIVALENCE/03-verification/mutation-record.txt	repo
docs/validation/evidence/P86-BOUNDDIR-UNC-EQUIVALENCE/03-verification/mutation-branch-removed.txt	repo
docs/validation/evidence/P86-BOUNDDIR-UNC-EQUIVALENCE/03-verification/test-post-restore.txt	repo
docs/validation/evidence/P86-BOUNDDIR-UNC-EQUIVALENCE/03-verification/local-flake-tickets-p39.txt	repo
docs/validation/evidence/P86-BOUNDDIR-UNC-EQUIVALENCE/03-verification/ci/run-44962802936-windows.log	repo
docs/validation/evidence/P86-BOUNDDIR-UNC-EQUIVALENCE/03-verification/ci/run-44962802936-ubuntu.log	repo
docs/validation/evidence/P86-BOUNDDIR-UNC-EQUIVALENCE/03-verification/ci/run-44965434402-windows.log	repo
docs/validation/evidence/P86-BOUNDDIR-UNC-EQUIVALENCE/03-verification/ci/run-44967047189-windows.log	repo
docs/validation/evidence/P86-BOUNDDIR-UNC-EQUIVALENCE/03-verification/ci/run-44967047189-ubuntu.log	repo
docs/validation/evidence/P86-BOUNDDIR-UNC-EQUIVALENCE/04-reviews/05-pre-review-r1.md	repo
docs/validation/evidence/P86-BOUNDDIR-UNC-EQUIVALENCE/04-reviews/05-pre-review-r2.md	repo
docs/validation/evidence/P86-BOUNDDIR-UNC-EQUIVALENCE/04-reviews/06-scan-r1.md	repo
docs/validation/evidence/P86-BOUNDDIR-UNC-EQUIVALENCE/04-reviews/06-scan-r2.md	repo
docs/validation/evidence/P86-BOUNDDIR-UNC-EQUIVALENCE/04-reviews/07-blind-review.md	repo
docs/validation/evidence/P86-BOUNDDIR-UNC-EQUIVALENCE/04-reviews/07-final-review-r1.md	repo
```

Untracked files: the process artifacts under `tmp/P86-BOUNDDIR-UNC-EQUIVALENCE/` (removed at (8)c).
