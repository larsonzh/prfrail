# P76 S1 Exit Closeout Assessment (`S1-EXIT-CLOSEOUT-ASSESSMENT`)
Date: 2026-11-16. Slice: `S1-EXIT-CLOSEOUT-ASSESSMENT` (P76, `[SLICE]` assessment, tier M).
[中文](s1-exit-closeout-assessment.md)

**Line-number caliber**: every `file:line` below was read from disk on 2026-11-16; when a line number drifts, the file name plus the section number is the primary anchor.

---

## 1. Change summary

- What this slice did: on the post-P75 state of the eight S1 Exit gates it produced a closeout roadmap and landed, per the user ruling of 2026-11-16, six rulings (Q2 / Q3 / Q7 / Q9 / Q10 plus the closeout order), two ADRs, three additive notes and one in-cell caliber correction.
- What this slice did not do: no real CLI end-to-end run, no AT-23 / AT-24 re-run, no change to `s1-reach-reachability{,_EN}`, no directive version bump (v1.69 stands), and no change to any v50–v59 snapshot.
- One-line conclusion: the four gates that can progress inside S1 now have executable closeout paths; S1 Exit as a whole cannot pass inside S1 (the S-G1 criterion is unmet, S-G8 stays blocked until S2, and S-G7 is unreachable).
- Execution-freedom list: `docs/validation/evidence/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/freedom-list.md`.

## 2. Fact package

| Item | Value |
|---|---|
| Slice / tier | `S1-EXIT-CLOSEOUT-ASSESSMENT` / `[SLICE]` assessment, tier M |
| Controller | `DeepSeek V4.1 Flash` |
| Step 0 baseline | `HEAD` = `origin/main` = `b90de21`; directive v1.69; `SELFTEST: PASS 737/737`; `TOTAL_FAIL=0` in the index / tree / ci scopes |
| Rolling window | v57 / v58 / v59 count 160 / 148 / 180 LF lines, with CN and EN equal on each side |
| User ruling | six rulings plus the additional item plus the X authorization on 2026-11-16 (see section 15) |
| Evidence pack | `docs/validation/evidence/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/` |
| Mechanical readings | max ADR = 47 (this slice adds ADR-046 / ADR-047); max OB = 139; max DR = 14; pilot deferral round 25 |

## 3. Q1: eight-gate and ten-gate mapping

| S1 Exit gate | s1-reach gate | Mapping basis | Verdict difference |
|---|---|---|---|
| S-G1 | G-A + G-B + G-C | `p57:88` against `s1-reach:22`-`:24` | one-to-three; G-B is a standalone mechanism leg with no counterpart in the eight gates |
| S-G2 | G-D | `p57:89` against `s1-reach:25` | consistent (G-D folds the real CLI end-to-end path into its unlock condition) |
| S-G3 | G-E | `p57:90` against `s1-reach:26` | consistent |
| S-G4 | G-F | `p57:91` against `s1-reach:27` | consistent |
| S-G5 | G-G | `p57:92` against `s1-reach:28` | different state words: `unknown` and `BLOCKED` carry the same meaning |
| S-G6 | G-H | `p57:93` against `s1-reach:29` | consistent; both declare that gate level and task level are not the same source |
| S-G7 | G-I | `p57:94` against `s1-reach:30` | consistent; the G-I criterion is stated more completely |
| S-G8 | no counterpart | s1-reach has no separate row | the semantics fold into the G-A / G-C unlock conditions and G-D; a coverage gap, not a contradiction |
| (none) | OB-21 | cross-slice observation item | not part of the eight gates; must not be folded into any S-Gx verdict |

The verdict sets differ: the eight-gate side mixes `PASS` / `PARTIAL` / `BLOCKED` / unreachable / not passed / unknown, while the ten-gate set is unreachable / not passed / executed but not independently established / unknown; mapping must not translate one verdict word into another across sets.

## 4. Q2: S-G1 / S-G2 gate-level verdicts

- Caliber = `(a) criterion-literal` (user ruling): each verdict follows the literal criterion in `p57` section 4.2, and an AT-level achievement is never converted automatically into a gate pass.
- S-G1 = not passed: the real-candidate leg is unreachable (the section 9 fact chain of `r3-real-cli-assessment.md`; `s1-exit.md:43` states that the real-candidate leg is not a real-candidate compatibility proof).
- S-G2 = passed: the criterion "the `supervised-black-box` mode is implemented and the AT-24 verdict holds" is met (T028 `COMPLETED` at `DEV_PLAN.md:119`; AT-24 conditionally closed at `s1-exit.md:44`).
- Consistency with the ten gates: S-G1 points the same way as G-A "unreachable"; S-G2 does not conflict with G-D being "executed but not independently established" (G-D folds the real CLI end-to-end path into its unlock condition, while this caliber assigns it to S-G8).
- Registered residual: the "the S-G2 gate-level verdict is pending the v54 snapshot" clause inside `s1-exit.md:52` is a historical statement outside this slice's boundary and is left untouched; synchronizing it is proposed for the follow-up slice `S1-EXIT-SG2-SYNC` (a working name; starting it is the user's call) and does not affect this slice's gate-level basis.

## 5. Q3: the `s1-exit.md:52` correction

- Read-back result: the status cell is `BLOCKED`; the first half says "T027 is incomplete and AT-23 is not passed" while the second half says "AT-23 stage-1 condition closure achieved", and the two read as competing claims inside one cell.
- Chosen form: option A (the default under the user X authorization) - the qualifier "as a whole" separates the overall verdict from the stage-1 verdict; option B was not used because option A already removes the ambiguity (the dependency-chain parenthesis states twice that stage 2 does not apply).
- Before (CN): `T027 未完成、AT-23 未通过；`
- After (CN): `T027 未完成、AT-23 整体未通过；`
- Before (EN): `T027 is incomplete and AT-23 is not passed, while`
- After (EN): `T027 is incomplete and AT-23 as a whole is not passed, while`
- Option texts under the X authorization (option B is drafted in this slice along the Q3 direction rather than taken from the user text, and serves only as a comparison record): option A = `AT-23 整体未通过` (the qualifier 整体 added); option B = `AT-23 整体未通过（阶段 1 条件闭合已达成，详见本行依赖链说明）`. This slice uses option A.
- Three grounds: against the three `PARTIAL` local gates at `:56`-`:58`, the core task chain carries an incomplete task and an unpassed AT and is genuinely blocked; against the residual item defined in `p57` section 4.3, the fix aligns wording rather than re-judging the state; and it matches the G-I basis sentence "core task chain `BLOCKED`" at `s1-reach:30`.
- Boundaries: the status cell is unchanged, the gate-table row names and row count are unchanged, `:56`-`:58` and `:60` are unchanged, and every other row is untouched.

## 6. Q4: the S-G7 basis row

- The verdict is unchanged: S-G7 stays unreachable (it depends on S-G1 through S-G6, where S-G1 is not passed, S-G8 is blocked and S-G5 is unapproved).
- The `p57` basis row needs no change: it cites only the state word, which this correction leaves intact.
- Synchronization actually landed: one additive pointer line after the `p57` section 4.2 table, pointing at this report and ADR-046.

## 7. Q5: the S-G3 unlock path

| Step | Action | Precondition | Evidence form |
|---|---|---|---|
| One | the product owner approves the formal Release and its version | ADR-010 approved, or approved inside the same act | ADR status update plus a release authorization record |
| Two | build the final ZIP candidate artifact | step one done; CI candidate-build available | the artifact plus `SHA256SUMS`, SBOM and license manifest |
| Three | bind the release record (version / commit / run / file name) | step two done | the three bindings required by `INSTALLATION.md:206` are written down |
| Four | rerun the install lifecycle against the final ZIP | step three done | install / upgrade / rollback / uninstall rerun records |

S-G3 depends on S-G4: the version lives inside ADR-010, so the S-G4 approval must precede step two or run in parallel with it, while binding and reruns wait for the version to freeze.

## 8. Q6: the S-G4 unlock path

- Approving party: the product owner (the `ADR_REGISTER` header requires the ruling party and date for every approval and forbids model sign-off).
- Post-approval edit points: the `ADR_REGISTER.md:20` status cell, the EOL item at `INSTALLATION_PLAN.md:32`, the EOL wording at `s1-exit.md:66`, the EOL item at `PRODUCT_REQUIREMENTS.md:107`, and the policy paragraph at the end of the support matrix.
- Boundary: approving ADR-010 is not the same as S-G5; they authorize different objects (release policy versus stage exit).

## 9. Q7: S-G5 and the independent security review

- Historical-status finding: never executed. The basis is `SECURITY.md:5` ("without independent security sign-off"), the reviewer still outstanding at `PRODUCT_REQUIREMENTS.md:107`, the later-gate wording at `DEV_PLAN.md:27`, and the absence of any independent security review report under `docs/validation/`.
- Effect on S-G5: the S-G5 criterion is only "the product owner approves", so the gate is not a direct precondition; it is, however, an outstanding S0-to-S1 gate in `PRODUCT_REQUIREMENTS` section 7 and must be dispositioned explicitly.
- Disposition (user ruling A): waived and deferred to S2, recorded as ADR-047; the uncovered scope is the threat-model surface, the hard-security-gate surface, and the secret and credential surface; the due date is the earlier of the S2 start date and 2027-03-31.
- Does it block the S1 Exit declaration: it is not a direct block under the eight-gate criteria; the S-G5 stepwise approval material must nevertheless carry this disposition.

## 10. Q8: S-G6 gate level versus task level

- Measured fields: the T018 dependency at `DEV_PLAN.md:94` reads "T017, T025, T026, T027, T028 and the whois input authorization"; the status at `:96` is `IN_PROGRESS / BLOCKED`.
- Caliber (user ruling c): state gate level and task level side by side without a single verdict; this slice changes nothing.
- Consequence: closing S-G3 / S-G4 / S-G5 moves only the gate level; T018 stays unchecked because T027 is unimplemented; the S-G6 gate level stays unreachable until S2.

## 11. Q9: support-matrix caliber correction

| Location | Corrected wording |
|---|---|
| `:13` | `R-3 = (A), cannot start for now; S-G8 stays blocked until S2 enforcement lands (ADR-045)` |
| `:14` | `AT-24 conditionally closed; the real CLI end-to-end path depends on S2 enforcement (R-3 = (A), ADR-045)` |

`:5` is untouched (the security-design baseline status has nothing to do with P75); the support levels are unchanged (Windows stays a production-support candidate and Linux stays core CI) and only the status wording changes.

## 12. Q10: when the product owner approves

- Caliber (user ruling C, two steps): the first step, inside S1, approves ADR-010 and the v0.1.0 Release so that S-G3 / S-G4 close; the second step, after S2, approves S1 Exit as a whole.
- How S-G5 is booked after the first step: "release portion approved, whole-stage approval pending S2", never as passed.
- Hard boundary: no approval may be read as S1 Exit passing; no declaration before `s1-exit.md:52` is corrected; S-G8 stays blocked.

## 13. Q11: REQ stage-ownership boundaries

| Item | Ownership | Basis |
|---|---|---|
| REQ-007 black-box sub-item | deferred to S2 | `PRODUCT_REQUIREMENTS.md:69`; the P71 note at `s1-exit.md:62` |
| REQ-009 full unified TUI | not required in S1 | `PRODUCT_REQUIREMENTS.md:71`; `s1-exit.md:62` |
| REQ-020 three-language runs | S1 keeps only the C plus Go contract | `PRODUCT_REQUIREMENTS.md:82` |
| AT-13 self-hosting | still in S1 and T017 is complete | `s1-exit.md:62`; the T017 row in `DEV_PLAN` |

Constraint: none of the above may ever serve as a "not passed" reason for any S-Gx; the blocked basis for S-G1 and S-G8 may only come from R-3 = (A), the unfinished T027 assembly and candidate compatibility.

## 14. Closeout roadmap

1. Correct the `s1-exit{,_EN}.md:52` caliber (landed in this slice).
2. Settle the S-G2 gate level (recorded in this slice, Q2 = a).
3. S-G4: the product owner approves ADR-010.
4. S-G3: formal Release, final ZIP build, the three bindings, and the lifecycle rerun.
5. S-G6: advance the gate level while explaining the task level; the task stays unchecked.
6. S-G5: two-step approval per Q10.
7. S-G1 / S-G8 / S-G7: deferred to S2, and the whole S1 Exit declaration waits for the post-S2 final approval.

Parallel window: the S-G4 approval can run alongside the S-G3 ZIP build preparation, while binding and reruns wait for the version to freeze. Blockers: the product-owner approval and the S2 enforcement landing, the latter outside S1 control.

## 15. User ruling record

| Item | Ruling |
|---|---|
| Q2 | `(a) criterion-literal` |
| Q3 | `A additive (BLOCKED preserved)` |
| Q7 | `A waive and defer to S2 plus ADR` |
| Q9 | `included in ②b` |
| Q10 | `(C) two-step approval` |
| Closeout order | `the section 12 sequence is accepted` |
| Additional | `p57 unchanged apart from a pointer; ⑥ may be skipped` |
| X | `for Q3 option A comes first and B applies only if A proves insufficient on read-back; if the line is already unambiguous, stop and escalate` |

The full ruling text is archived at `02-ruling/ruling-2026-11-16.md` in the evidence pack; it takes effect on 2026-11-16 and the ruling party is the product owner.

## 16. CI fact chain

This is an assessment slice, so the mapping between the landing commit and the CI run is added by a backfill commit (§11.2); the backfill commit run itself is an observation record only.

- Landing commit `77f23d2` (`docs: assess the S1 exit closeout and land the P76 rulings`), pushed to `origin/main` over the range `b90de21..77f23d2` (no gitee push).
- CI run `43785603244` was green on both legs at the first attempt: ubuntu `3m24s`, windows `5m12s`; both legs reported `SELFTEST: PASS 737/737`, `GATE REPORT scope=ci changed=33` and `TOTAL_FAIL=0`.

## 17. Known boundaries

- Real enforcement stays in S2; S-G8 stays blocked; the AT-23 stage-1 condition closure is a permanent fact; no AT is re-run; the probe budget is zero.
- The independent security review is waived and deferred to S2; the uncovered scope is listed in section 9 and ADR-047.
- The S-G2 clause at `s1-exit.md:52` is not synchronized (this slice only edits the AT-23 fragment) and is registered as a residual; the proposed carrier slice is `S1-EXIT-SG2-SYNC` (pending the user's decision to start it).

## 18. Next steps

- An ADR-010 approval slice (a product-owner action).
- A Release and ZIP binding plus rerun slice.
- An independent-security-review slice ahead of the S2 start.

## 19. Cost and review chain

| Stage | Count |
|---|---|
| ① Architecture | 1 (`deep-reasoner` plus `DeepSeek V4 Pro`) |
| Batch assistance | 3 (`GLM-5.3-Flash`, fact enumeration) |
| ⑤ Pre-review | 2 rounds: r1 `FINDINGS` (1 high + 5 low) then r2 `PASS WITH FIXES` (none at medium or above, 2 low) |
| ⑥ Independent scan | skipped, confirmed by ⑦ |
| ⑦ Final review | 1 round: `PASS WITH FIXES` (none at medium or above, 1 low; includes `SKIP CONFIRMED`) |
| Probes | 0 |
| R / D | R = 1, D = 0 |

- ⑥ skipped: the three §7.2.1 conditions hold (the change set is documentation and ledger material only; it carries no `.go` / `.json` / schema / fixtures / script semantic change; and it touches no state machine, gate criterion, evidence model or role/authorization contract semantics), confirmed by ⑦ with the reason recorded.

## artifacts

```artifacts
docs/ADR_REGISTER.md	repo
docs/ADR_REGISTER_EN.md	repo
docs/PRODUCT_REQUIREMENTS.md	repo
docs/PRODUCT_REQUIREMENTS_EN.md	repo
docs/S1_SUPPORT_MATRIX.md	repo
docs/S1_SUPPORT_MATRIX_EN.md	repo
docs/SECURITY.md	repo
docs/SECURITY_EN.md	repo
docs/validation/p57-closeout-assessment.md	repo
docs/validation/p57-closeout-assessment_EN.md	repo
docs/validation/s1-exit.md	repo
docs/validation/s1-exit_EN.md	repo
docs/validation/s1-exit-closeout-assessment.md	repo
docs/validation/s1-exit-closeout-assessment_EN.md	repo
docs/validation/evidence/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/00-briefing.md	repo
docs/validation/evidence/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/input-list.md	repo
docs/validation/evidence/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/freedom-list.md	repo
docs/validation/evidence/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/MANIFEST.md	repo
docs/validation/evidence/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/SHA256SUMS.txt	repo
docs/validation/evidence/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/01-arch/design-r1.md	repo
docs/validation/evidence/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/02-batch/batch-01-s1-exit-p57.md	repo
docs/validation/evidence/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/02-batch/batch-02-s1-reach.md	repo
docs/validation/evidence/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/02-batch/batch-03-lagging-docs.md	repo
docs/validation/evidence/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/02-ruling/ruling-2026-11-16.md	repo
docs/validation/evidence/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/04-reviews/prereview-r1.md	repo
docs/validation/evidence/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/04-reviews/prereview-r2.md	repo
docs/validation/evidence/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/04-reviews/final-r1.md	repo
docs/validation/evidence/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/06-gates/gate-index.txt	repo
docs/validation/evidence/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/06-gates/gate-tree.txt	repo
docs/validation/evidence/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/06-gates/gate-ci.txt	repo
docs/validation/evidence/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/06-gates/gate-selftest.txt	repo
```
