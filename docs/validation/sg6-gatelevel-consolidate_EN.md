# P80 Verification Report: S-G6 Gate-Level Consolidation and the S1 Exit Next-Step List

Date: 2026-11-18. Status: complete. Slice: P80 `SG6-GATELEVEL-CONSOLIDATE` (`[SLICE]`, tier-S documentation type).

## 1. Slice Identity and Authorization

- Adjudicated object: consolidate the gate-level verdicts and evidence landing spots of S-G1 through S-G5, decide whether the S-G6 gate-level side holds, state the side-by-side wording that separates it from the T018 task level, and produce the S1 Exit next-step list.
- Authorization: the 2026-11-18 start message is the same-turn authorization; the S-G5 wording ruling was made by the user on the same day and recorded in v62 §6.
- Boundary: this slice carries adjudicated facts, it does not modify criteria; the p57 and `s1-exit` gate tables are untouched, and the v62 snapshot is untouched.
- Budget: probes 0; (7) is 1 final review + 1 re-review (tier S, §7.5); one blind review follows the §7.7 proportional spot-check and is counted separately.

## 2. Step-0 Anchor Measurements

| Anchor | Expected | Measured |
|---|---|---|
| `HEAD` / `origin/main` | `b5ff8ee` | `b5ff8ee` / `b5ff8ee` |
| Directive version | v1.69 | v1.69 |
| `SELFTEST` | 737/737 | 737/737 |
| Three scopes | `TOTAL_FAIL=0` | index / tree / ci all 0 |
| `tmp` top level | `gates-cache` + `.gitkeep`; 1 file after excluding the cache | same (the slice directory was created after the check) |
| `git status --porcelain -uall` | 0 entries | 0 entries |
| Rolling window v60 / v61 / v62 | 242 / 214 / 217 | 242 / 214 / 217 (CN and EN symmetric) |

## 3. Design Phase Conclusions

The full design is archived at `01-arch/design-r1.md` (`DESIGN: DONE`). Item by item:

- Q1: the five verdicts and their landing spots are tabulated; criteria are quoted verbatim from p57 §4.2, verdicts verbatim from v62 §4, and every gate carries at least two usable coordinates.
- Q2: **the S-G6 gate-level side holds**; the second half of the criterion (the four remaining gates closed) is satisfied, while the first half (dependency fields satisfied) stays on the task level.
- Q3: the side-by-side wording is "the gate level holds is not the criterion passing and not T018 being complete"; the T018 task level stays bound by T027 and the product-owner approval, and the `DEV_PLAN` status field is unchanged.
- Q4: the four next-step items (S-G5 step two, re-evaluation prerequisites, the P79 security review, the S2 scheduling advice) each carry an action, a precondition, an evidence shape and an owner.
- Q5: three first-layer items must change; every second-layer item was read and judged not to change (including `DEV_PLAN{,_EN}.md`); the third layer definitely does not change.
- Q6: the three searches ran independently with zero new hits; no upgrade or downgrade action was needed.

## 4. Q1: Gate-Level Verdicts and Landing Spots for S-G1 through S-G5

| Gate | Criterion (quoted from p57 §4.2) | Current verdict (v62 §4) | Evidence landing spot | Source slice |
|---|---|---|---|---|
| S-G1 | "a CLI can really launch an admitted AI candidate and bill it" | not passed (Q2=(a), criterion literal) | `p57:88`; `v62:115`; P76 report §4 | P76 |
| S-G2 | "`supervised-black-box` implemented and the AT-24 verdict holds" | passed (Q2=(a), criterion literal) | `p57:89`; `v62:116`; `s1-exit:52`; `DEV_PLAN:121` | P76 |
| S-G3 | "the final ZIP is bound to the release record and re-verified" | passed (P78 closed) | P78 report §5 / §7; `p57:90`; `v62:117`; `s1-exit:66` | P78 |
| S-G4 | "ADR-010 approved by the product owner" | passed (ADR-010 flipped to `APPROVED`) | `ADR_REGISTER:20`; P77 report §5; `p57:91`; `v62:118` | P77 |
| S-G5 | "the product owner makes the approval" | release part approved (step one) | `v62:119`; `v62:196`; ADR-046 (`ADR_REGISTER:57`); P76 report §12 | user ruling |

Reading note: the coordinates above come from the 2026-11-18 on-disk read; both v62 and p57 are cited read-only, and no gate-table word was changed.

## 5. Q2: The S-G6 Gate-Level Verdict

- Criterion verbatim (`p57:93`): "the dependency fields are satisfied and the four remaining gates are closed"; the same row states the "the gate level and the task level are not the same source" constraint.

| Criterion component | Verdict | Basis |
|---|---|---|
| S-G2 passed at gate level (P76) | satisfied | P76 report §4 (`s1-exit-closeout-assessment.md` §4); `v62:116`; `s1-exit:52` |
| S-G3 passed | satisfied | P78 report §7 (`release-zip-binding.md:69`-`:73`); `v62:117` |
| S-G4 passed | satisfied | P77 report §5 (`adr010-approval.md:90`-`:97`); `ADR_REGISTER:20` |
| S-G5 release part approved (step one) | satisfied (the user ruled on 2026-11-18 that this suffices for the gate level) | `v62:196`; ADR-046 (`ADR_REGISTER:57`) |
| Dependency fields | task-level component, not adjudicated at the gate level | `DEV_PLAN.md:94` (dependency field verbatim); `DEV_PLAN.md:96` (status `IN_PROGRESS / BLOCKED`) |

Exact wording of the gate-level verdict (this section is the on-disk text; the ledger completion record cites the same):

> **The S-G6 gate-level side holds** (never "S-G6 passed"): the second half of the criterion, "the four remaining gates are closed", is satisfied - S-G2 at gate level (P76 Q2=(a)), S-G3 (P78), S-G4 (P77), and S-G5 recorded as "release part approved, step one" (user ruling of 2026-11-18, v62 §6). The first half, "the dependency fields are satisfied", belongs to the task level (`DEV_PLAN.md:94` lists T027 and T028 among the dependency fields), so the task level does not hold and stays side by side (T018 is not checked, and the `DEV_PLAN.md:96` status field is unchanged). The split follows P76 Q8=(c). This report carries the verdict; the p57 and `s1-exit` gate tables are untouched.

Falsification conditions (any single observation overturns the verdict and forces a stop-and-rejudge):

1. Any of S-G2 / S-G3 / S-G4 / S-G5 has its gate-level verdict reverted by later evidence or a user ruling.
2. The user rules later that the step-one approval is insufficient to close the gate level.
3. The P76 Q8=(c) side-by-side wording is replaced by a new ADR or user ruling.
4. The ADR-046 rollback clause triggers (later evidence overturns the Q2 criterion-literal reading, which requires a new slice and a new ADR).

## 6. Q3: Side-by-Side Wording for the T018 Task Level

Exact side-by-side text (this section is the on-disk text):

> **The S-G6 gate-level side holds is not the criterion passing and not T018 complete.** The T018 task level remains bound by its dependency fields: T027 is not implemented (`BLOCKED / NOT IMPLEMENTED`), AT-23 as a whole is not passed (`s1-exit:52`: "T027 is incomplete and AT-23 as a whole is not passed"), and the product-owner S1 exit approval is missing; therefore T018 stays unchecked and the `DEV_PLAN.md:96` status field keeps `IN_PROGRESS / BLOCKED`; the P78 note at `DEV_PLAN.md:98` ("the T018 check state is unchanged") stays in force. The gate-level verdict records closure progress of the four remaining gates only and converts into no T018 progress.

P76 Q8=(c) verbatim (`s1-exit-closeout-assessment.md` §10):

- "Measured fields: the T018 dependency at `DEV_PLAN.md:94` reads T017, T025, T026, T027, T028 and the whois input authorization; the status at `:96` is `IN_PROGRESS / BLOCKED`."
- "Wording (user ruling c): state the gate level and the task level side by side without a single conclusion; this slice changes nothing."

The third sentence, "S-G6 stays unreachable at gate level until S2", is a 2026-11-16 statement from when the four remaining gates were still open. P77 and P78 closed S-G3 and S-G4, the user ruled the S-G5 wording on 2026-11-18, and v62 §4 already records S-G6 as "gate level closable, pending the P80 consolidation"; item 5 of the P76 §12 roadmap is exactly "advance the gate level side by side and keep the task level unchecked", and this slice occupies that position.

The two `DEV_PLAN` spots: `:96` is a dated historical status sentence with a reading note and `:98` is the additive P78 note; both are historical records, and this slice neither rewrites them nor flips the check state.

## 7. Q4: S1 Exit Next-Step List

| Item | Action | Precondition | Evidence shape | Owner |
|---|---|---|---|---|
| S-G5 step two | after S2 unlocks, approve the S1 Exit as a whole and record the four elements | S-G1 / S-G8 re-evaluated and passed; T027 complete; AT-23 passed as a whole; the independent security review completed | the ruling text (ruling person / date / reason / basis) | product owner |
| S-G1 / S-G8 prerequisites | re-run the real E2E after S2 enforcement lands and re-evaluate | R-08, R-09, the T027 assembly and a compatible candidate are all in place, and enforcement is wired | a re-evaluation report and a real E2E evidence chain | later slice, after S2 |
| S-G7 | no re-evaluation; stays unreachable inside S1 | depends on S-G1 through S-G6 all green | this list is the explaining carrier | later slice, after S2 |
| P79 independent security review | complete one independent review and record the conclusion before S2 starts | before S2 starts; scope = the S0-to-S1 security gate in `PRODUCT_REQUIREMENTS` §7 | the review report; an updated `SECURITY.md` note; an updated ADR-047 status | later slice P79 |

S2 scheduling advice (**for the user reference only, not a ruling**): run the P79 independent security review first, then start S2 and the enforcement wiring, then re-evaluate S-G1 / S-G8, and finally take S-G5 step two. The ADR-047 due date is the earlier of the S2 start date and 2027-03-31 (`ADR_REGISTER:58`); if the review is not completed by then, the related S2 measurements must not start.

Two objects must not be mixed: the P79 independent security review is the historical gap under `PRODUCT_REQUIREMENTS` §7 (never executed); the "independent review" row in the `s1-exit` gate table (`:59`) refers to the T017 baseline and Node 24 replacement slice reviews. Their evidence is not interchangeable.

## 8. Q5: Documentation Sync Decision (Three Layers)

First layer (must change in this slice):

| Document | Action |
|---|---|
| `docs/validation/sg6-gatelevel-consolidate{,_EN}.md` | create the report pair (this document) |
| `docs/t027/REMAINING_SLICES{,_EN}.md` | P80 registration (step 0c) and the ⑧b completion record |
| `docs/validation/evidence/P80-SG6-GATELEVEL-CONSOLIDATE/` | the evidence pack (created in step 0, appended round by round) |

Second layer (all read and judged not to change):

- `DEV_PLAN{,_EN}.md`: not changed. This report carries the verdict; `:96` and `:98` are dated historical records, and `DEV_PLAN` carries no criterion meaning that needs updating. Had the design judged otherwise, this slice would have stopped and escalated under §1.7.2; it does not.
- The other documents hit by the three Q6 searches are all known historical statements or carriers of this slice, each judged not to change (details in §9).

Third layer (definitely unchanged): the directive, `CONTRACTS`, `SECURITY`, `ARCHITECTURE`, `TEST_STRATEGY`, `BUSINESS_WORKFLOWS`, `CODING_CONVENTIONS`, `DOCUMENTATION_PLAN`, `RFC-*`, `CHECKLIST_v62{,_EN}`, every snapshot from v50 to v62, the p57 gate table, the `s1-exit` gate table, product code, CI, and the other ADRs.

## 9. Q6: Combined Search Sweep

The three searches ran independently (never combined); the baseline is HEAD's tracked files and the delta is the files this slice adds:

| Pattern | HEAD baseline hits | Carriers added after the baseline | New hits elsewhere |
|---|---|---|---|
| `S-G6` | 75 | this report pair, the ledger P80 segment, the evidence pack | 0 |
| `SG6` | 8 | this report pair, the ledger P80 segment, the evidence pack | 0 |
| `gate level` | 6 | this report pair, the ledger P80 segment, the evidence pack | 0 |

Attribution: every HEAD hit is a known carrier (the p57 / `s1-exit` / P76 / P75 / P73 / P61 reports, the v59 to v62 snapshots, the ledger, and the frozen evidence tree); both sides are decided by comparing file sets rather than hit counts, so nothing new appears outside this slice's carriers, no upgrade or downgrade action is needed, and the second-layer "no change" judgement agrees with the search results.

## 10. Phase-2 Execution and Per-File Changes

- Report pair: `docs/validation/sg6-gatelevel-consolidate{,_EN}.md` created, CN and EN paired line by line with all six shapes equal.
- Ledger: the P80 segment was registered in step 0c (CN and EN gained the same number of rows), and ⑧b appends the completion record and flips the status.
- Evidence pack: the three initial payloads were created in step 0; the design, the batch enumeration, the review outputs and the raw gate outputs were appended round by round, and ⑧c recomputes the checksums.

| File pair | Added | Deleted |
|---|---|---|
| `docs/validation/sg6-gatelevel-consolidate{,_EN}.md` | see the ⑨ measurement | see the ⑨ measurement |
| `docs/t027/REMAINING_SLICES{,_EN}.md` | see the ⑨ measurement | see the ⑨ measurement |

Counting rule: the per-pair equality of additions and deletions comes straight from the raw `G3(a)` output; this table copies no numbers, so it cannot drift from the measured readings.

## 11. Phase-3 Consistency Check

- The S-G6 gate-level verdict agrees with the `p57` §4.2 criterion text: the criterion asks only for the four remaining gates closed and the dependency fields satisfied, and this slice rules on the first and keeps the second side by side.
- No conflict with the `s1-exit` gate table: rows and verdict words are untouched; the `:52` wording "AT-23 as a whole is not passed" points the same way as this slice stating T027 is not implemented.
- No conflict with v62 §4: v62 records "gate level closable, pending the P80 consolidation", and this slice supplies that consolidation verdict.
- No conflict with ADR-046: the two-step structure of Q10=(C) is preserved in full, step one is recorded as "release part approved", and step two still waits until after S2.
- No conflict with P76 Q8=(c): the side-by-side wording is kept verbatim and no single conclusion is given.
- CN and EN align sentence by sentence: the report pair is paired line by line with isomer terms and verdict words, and the ledger segment has the same row count on both sides.

## 12. Phase-4 Gate Readings

| Criterion | Reading |
|---|---|
| `--scope=index` | `TOTAL_FAIL=0` (`changed=19`) |
| `--scope=tree` | `TOTAL_FAIL=0` (`changed=19`) |
| `--scope=ci` | `TOTAL_FAIL=0` (local degraded reading, `changed=2`) |
| `--selftest` | `SELFTEST: PASS 737/737` |

Reading scope: the four raw outputs under `06-gates` come from this round's snapshot (`changed=` appears in each file header) and were produced in the same batch as the `artifacts` fence row count and the `G7-b` count; afterwards only review payloads and the ledger completion record are appended, the gate verdict does not change, and ⑨ re-runs on the final snapshot and states the final `changed=` in the ⑩ receipt.

`--scope=ci` degrades to the previous commit range when no CI environment variables are present and only covers `HEAD^1..HEAD`, so its `changed` count does not reflect this slice's staging set; the real coverage of this slice comes from `--scope=tree` and `--scope=index`.

Reading rule: this is a documentation-only slice with zero product-code change and adds no build step; the raw outputs of the three scopes and the self-test are archived under `06-gates/`.

## 13. Review Chain

- (5) Pre-review: `DeepSeek V4 Pro` (`deep-reasoner`, read-only); rounds and verdicts are recorded in §13.1.
- ⑥ skipped: all three §7.2.1 conditions hold (changes land only on documents and ledgers; no script or code semantics change; no state machine, gate criterion, evidence model, role or authorization contract semantics), confirmed by (7).
- (7) Review chain: a blind review (§7.7 proportional spot-check, separate budget) plus the final review and the re-review; completed rounds are listed in §13.1 one by one and unfinished rounds are not described in advance.

### 13.1 Review Round Record

| Stage | Round | Verdict line | Findings | Disposition |
|---|---|---|---|---|
| (5) pre-review | r1 | `PRE-REVIEW: PASS WITH FIXES` | High 0 / Medium 0 / Low 3 / Info 3 | the three Low items were fixed minimally: search counts restated against the HEAD baseline, the S-G2 gate-level component added to §5, and the fence rebuilt from the staging truth; the three Info items are recorded one by one |
| (7) blind review | r1 | `BLIND REVIEW: FINDINGS` | High 0 / Medium 1 / Low 1 / Info 0 | both items are point-in-time gaps (the round record and the ledger state were not landed in the same batch) and touch no verdict; this row plus the ⑧b ledger flip are their closure evidence |
| (7) final review | r1 | `FINAL REVIEW: FINDINGS` | High 0 / Medium 1 / Low 0 / Info 0 | the Medium is the gate archive not being in the same batch as the snapshot; `06-gates` was re-run and refilled from one snapshot, for r2 to re-check |
| (7) re-review | r2 | `RE-REVIEW: FINDINGS` | High 0 / Medium 1 / Low 0 / Info 0 | the Medium points at the missing local-degradation note for `--scope=ci` in §12; that note is now in §12, but the tier-S budget (1 final review + 1 re-review) is exhausted, so the slice stops and escalates to the user instead of self-approving |
| (7) narrowed read | r1 (authorized exception) | `RE-REVIEW: PASS` | High 0 / Medium 0 / Low 0 / Info 0 | the single r2 Medium is judged `CLOSED`; the exceptions chapter passed the truthfulness check; the authorization applies to this slice only and is not a precedent (§14.3) |


## 14. Exceptions (R counting / risk marking / D clarification)

### 14.1 Per-round R counting (§9.2 is the single source of truth)

| Id | Stage | Event | R | Disposition |
|---|---|---|---|---|
| R-1 | (7) final review r1 | one Medium (`06-gates` archive not in the same batch as the snapshot) | +1 | fixed in place: rerun and refill `06-gates` from one batch; the rerun counts toward R per §5.3 |
| R-2 | (7) re-review r2 | one Medium (§12 missing the `--scope=ci` degradation note) | +1 | fixed in place: three reading-scope sentences added to §12; the rerun counts toward R per §5.3 |
| Total | - | - | **R = 2** | reaches the §9.2 high-risk threshold |

Two rounds do not enter R: (5) pre-review r1 (3 Low / 3 Info, no Medium+) and the (7) blind review r1 (1 Medium + 1 Low, both point-in-time and already closed) - the blind review is a spot-check read under §7.7 whose findings are recorded separately in §13.1 and stay outside the R sequence.

### 14.2 High-risk slice marking (§9.2: R>=2)

- Marking: **this slice is a high-risk slice** (per the §5.2 high-risk row and §9.2 clause 3).
- Budget: the **tier-S budget is not exceeded** (the start message preset 1 final review + 1 re-review, with the blind review counted separately). Although the §5.2 high-risk row makes a (5) re-check mandatory, this slice **does not self-upgrade to 1+3**; the disposition is to **stop and escalate to the user**.
- User ruling (same turn, 2026-11-18): **option A - one narrowed read is authorized** (the §7.8.3 "user-authorized exception" row): budget = one (7) Codex regular-tier call, valid for this slice only and not a precedent (same reading as the GATES-EXT precedent). The authorization text and the input-narrowing declaration are in §14.3, and the narrowed read verdict is in §13.1.

### 14.3 Authorization text and input-narrowing declaration

- Authorization text (verbatim excerpt of the 2026-11-18 user message): the user ruled option A, authorizing one narrowed read under the §7.8.3 user-authorized exception row with a budget of one regular-tier Codex (7) call, valid for this slice only and not usable as a precedent (the GATES-EXT precedent reading).
- Input narrowing: the read receives (1) the three reading-scope sentences of §12; (2) the same-batch readings under `06-gates/`; (3) `impl.numstat`; (4) this chapter. It does **not** receive the slice definition, the full change set, the (5) list, or the blind-review and final-review finding lists.
- Task narrowing: verify only whether the single r2 Medium (the `--scope=ci` degradation note in §12) is CLOSED and whether the patches in this chapter are truthful.
- Validity boundary: the authorization **does not generalize** and applies to this slice only; the verdict line comes from the §4.1 closed set and lands in `04-reviews/narrow-r1.md`.

### 14.4 D classification clarification

- Fact: the r1 remediation script `fix-final1.js` was written but **never executed** before r2 was dispatched, so r2 reviewed text that had not been remediated.
- Classification: the slip is an **orchestration-sequence error** (it violated the "after remediation, redo the paired action in place" rule between (5) and (7)) and **not** a temp-file location violation - every `tmp/` artifact sits under this slice directory, so OB-139 does not trigger.
- Counting: this slice **keeps a self-imposed D = 1**, recorded in this section rather than in the temp-file violation sequence, and marked "**self-imposed reading, for user ratification**". If the user ratifies D = 0 under the strict reading, this section stands as the substitute note for that exception (the only difference between the two readings is whether an orchestration error counts as D).

## 15. Known Limits and Residuals

- The `p57` and `s1-exit` gate tables keep their P76-time verdict words and sit side by side with this slice stating the S-G6 gate-level side holds; leaving the tables untouched is P76's explicit order, and this slice flips no row.
- `DEV_PLAN.md:96` and `:98` stay as they are; the live carriers of the S-G6 gate-level verdict are this report and the ledger P80 completion record.
- Six wording differences between the start message and the on-disk facts were judged one by one: zero conclude-level impact, four are wording-level items, and each is recorded at the matching place in §4 / §5 / §6 using the measured coordinates.
- The concrete differences: `ADR_REGISTER` holds ADR-009 at `:19` and ADR-010 at `:20`; the v62 §1 HEAD is a generation-time anchor; the 2026-11-17 and 2026-11-18 rulings have different objects; T028 is already complete while the live task-level constraints are T027 and two further causes; the P76 ruling c is sourced from the §10 body, not the §15 table.
- When the pack was created in step 0 its metadata was not yet in the index, so `G8-a` is a timing red before the ⑧b staging; it is not a defect.
- The historical double-BOM residue (t025-related files) is a HEAD fact and outside this slice.
- This slice stopped at the user-ruling point after the (7) re-review (r1 and r2 each reported one Medium and the tier-S budget was exhausted); on 2026-11-18 the user ruled option A - one narrowed read, see §14.2 and §14.3.

## 16. Accounting

| Stage | Count |
|---|---|
| (1) architecture | 1 (`deep-reasoner` + `DeepSeek V4 Pro`) |
| Batch assistance | 1 (`GLM-5.3-Flash`) |
| (5) pre-review | 1 call (`deep-reasoner` + `DeepSeek V4 Pro`) |
| (6) independent scan | 0 calls (skipped under §7.2.1, confirmed by (7) `SKIP CONFIRMED`) |
| (7) reviews | 4 calls: blind review 1 (separate §7.7 budget) + final review 1 + re-review 1 + narrowed read 1 (user-authorized exception) |
| Probes | 0 |
| R / D | see §14 exceptions (R = 2, D = 1 under the self-imposed reading) |

## 17. Commit and CI

The mapping between the main commit and the CI run is filled in by the §11.2 backfill commit; the backfill commit's own run is an observation record only, and no further commit is made for it.

```artifacts
docs/t027/REMAINING_SLICES_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/validation/evidence/P80-SG6-GATELEVEL-CONSOLIDATE/00-briefing.md	repo
docs/validation/evidence/P80-SG6-GATELEVEL-CONSOLIDATE/01-arch/design-r1.md	repo
docs/validation/evidence/P80-SG6-GATELEVEL-CONSOLIDATE/02-batch/batch-01-sg6-sources.md	repo
docs/validation/evidence/P80-SG6-GATELEVEL-CONSOLIDATE/04-reviews/blind-r1.md	repo
docs/validation/evidence/P80-SG6-GATELEVEL-CONSOLIDATE/04-reviews/final-r1.md	repo
docs/validation/evidence/P80-SG6-GATELEVEL-CONSOLIDATE/04-reviews/final-r2.md	repo
docs/validation/evidence/P80-SG6-GATELEVEL-CONSOLIDATE/04-reviews/narrow-r1.md	repo
docs/validation/evidence/P80-SG6-GATELEVEL-CONSOLIDATE/04-reviews/prereview-r1.md	repo
docs/validation/evidence/P80-SG6-GATELEVEL-CONSOLIDATE/06-gates/gate-ci.txt	repo
docs/validation/evidence/P80-SG6-GATELEVEL-CONSOLIDATE/06-gates/gate-index.txt	repo
docs/validation/evidence/P80-SG6-GATELEVEL-CONSOLIDATE/06-gates/gate-selftest.txt	repo
docs/validation/evidence/P80-SG6-GATELEVEL-CONSOLIDATE/06-gates/gate-tree.txt	repo
docs/validation/evidence/P80-SG6-GATELEVEL-CONSOLIDATE/freedom-list.md	repo
docs/validation/evidence/P80-SG6-GATELEVEL-CONSOLIDATE/input-list.md	repo
docs/validation/evidence/P80-SG6-GATELEVEL-CONSOLIDATE/MANIFEST.md	repo
docs/validation/evidence/P80-SG6-GATELEVEL-CONSOLIDATE/SHA256SUMS.txt	repo
docs/validation/sg6-gatelevel-consolidate_EN.md	repo
docs/validation/sg6-gatelevel-consolidate.md	repo
```
