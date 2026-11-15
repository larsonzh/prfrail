# ProofRail R3-REAL-CLI-ASSESSMENT report (P75 · residual assessment of the layer-1 T027-B main-line milestone · assessment [SLICE] · M tier)

Date: 2026-11-15. Status: step-2 write-up (step 1 already ruled at its stop point; steps 2a-2c land in sequence). Slice: `R3-REAL-CLI-ASSESSMENT`.

R-3 ruling: **(A) cannot start for now** (user ruling in the same round, 2026-11-15; attached rulings A1-A5 in §4).

## 1. Executive summary

- This slice answers the three R-3 residual questions of the T027-B main-line milestone: whether the milestone can be declared complete, whether S-G8 can advance, and whether the closed AT-23 stage-1 state is maintained.
- Conclusion: R-3 = (A) cannot start for now; T027-B is declared "complete (partly not applicable)"; S-G8 stays blocked until S2 enforcement lands; the closed AT-23 stage-1 state is registered as a permanent fact.
- Evidence chain: B1 found no acceptable candidate inside the evaluated scope + B2 external enforcement is not wired in + P74 landed only the check layer (S1 defaults to off) + P61 / F2-FIX already closed stage 1.
- Downstream: real CLI end-to-end is unreachable inside S1; the S1 Exit closure point does not move earlier because of this ruling; `p57` and `DEV_PLAN` carry the ruling through additive pointers and a section update.
- This slice spends no probe, changes no gate-table decision, and changes no contract or product code.

## 2. Current-state check

### 2.1 T027-B: C1-C7 / N1-N4 / R-1 / R-2

- C1-C7: P39 C6 / P40 C1 / P41 C2 / P42 C3 / P43 C4 and C5 / P44 C7 have all landed (`docs/DEV_PLAN.md:180-185`).
- N1-N4: definitions in `docs/validation/t027b-tactical-blueprint.md:24`; landed by P41 / P42 / P43 / P44.
- R-1: closed by P69 `T027B-E2E-WIRING` (`docs/validation/p57-closeout-assessment.md:72`).
- R-2: ruled by P73 = option b (`docs/validation/p57-closeout-assessment.md:74`), landed by P74 as the `Sandbox` check layer (`docs/validation/sandbox-interface.md:8-9`).

### 2.2 AT-23 stage-1 conditional closure

- Current: stage 1 conditional closure reached (`docs/validation/s1-exit.md:43`; `docs/validation/at23-restricted-candidate.md:499`).
- Premise 1: the mechanism leg is green, judged by step 7 of the F2 fix slice (`docs/validation/at23-restricted-candidate.md:423`).
- Premise 2: the real-candidate leg recorded unsupported at least three times, satisfied at the P61 point (`docs/validation/at23-restricted-candidate.md:19-20`).
- Boundary: R-08.5 states that conditional closure is not a proof of real-candidate compatibility (`docs/validation/r08-r09-ruling.md:38-40`).

### 2.3 S-G1 / S-G8 and the eight-gate decisions

- S-G1: unreachable / not passed (`docs/validation/p57-closeout-assessment.md:88`).
- S-G8: not passed, depends on T027-B (`docs/validation/p57-closeout-assessment.md:95`).
- Other gates: S-G2 wording is stale; S-G3 and S-G4 not passed; S-G5 unknown (manual gate); S-G6 and S-G7 unreachable (same file `:89-94`).
- P74 boundary: the `Sandbox` check layer does not unlock S-G1 / S-G8 (`docs/validation/sandbox-interface.md:133`).

### 2.4 B1 / B2 and R-08 / R-09

- B1: no acceptable candidate inside the evaluated scope (`docs/validation/t027-b1-candidate-live-discovery.md:4`, `:111`).
- B2: external OS enforcement, with the in-boundary and out-of-boundary split at `docs/validation/B2-EXTERNAL-ENFORCEMENT.md:12`, `:18-19`.
- B2 limit: the candidate session fails inside a zero-capability container (`docs/validation/B2-EXTERNAL-ENFORCEMENT.md:162-164`).
- Matrix wording: the pinned candidate is 10 verified / 2 unsupported / 0 unknown, incompatible as a whole (`docs/validation/s1-exit.md:3`).
- R-08 elements: R-08.3 only P; R-08.4 no VM and no proxy; R-08.5 conditional closure; R-08.6 the BYOK parallel path (`docs/validation/r08-r09-ruling.md:28-57`).

### 2.5 R-3 origin and definition

- Origin: `docs/validation/p57-closeout-assessment.md:70` (layer-1 item 8, `AT23-RESTRICTED-CANDIDATE`).
- Position: an assessment item, not a re-run (stated by the user ruling of 2026-11-15).

## 3. R-3 three-way comparison

- (A) cannot start for now: keep the current state and register it; zero cost; S-G8 stays not passed; no existing closure is disturbed; rollback = no action.
- (A) does not do: no re-run of the real-candidate leg, no P-prime, no candidate swap, no gate-decision change.
- (B) start in stages: reuse the F2-FIX evidence for the mechanism leg and re-run the real-candidate leg at least three times; cost = per-run probe declarations plus one M-tier implementation slice; the outcome is at most a re-confirmation of conditional closure.
- (B) does not do: no P-prime, no candidate swap, no criterion change, no S-G8 advance.
- (C) introduce P-prime or swap the candidate: needs a new candidate assessment slice or a device premise; highest cost and argument burden; the only path that could actually advance S-G8.
- (C) does not do: keep the pinned candidate as the only channel; land anything inside this slice (a new slice is required).

| Axis | (A) | (B) | (C) |
|---|---|---|---|
| Does | keep and register | reuse the mechanism leg + re-run the real leg | new candidate assessment or a P-prime device |
| Does not | no re-run / no swap | no swap / no S-G8 advance | not landed in this slice |
| Cost | 0 | probe + one M slice | high (several slices + argument) |
| T027-B completion | wording can land | wording can land | unaffected |
| S-G8 advance | no advance | no advance | only if the swap succeeds |
| AT-23 closure | maintained | maintained | upgrades on success |
| S1 Exit timing | no move | no move | earliest |
| Rollback | no action | slice rollback | device needs paired discipline |

- Step 1 recommended (B); the user ruled (A) on 2026-11-15 and noted that "consistent with the R-08/R-09 baseline" does not hold - R-08/R-09 are the targets of AT-23 stage 1, which is already reached. A recommendation does not replace a ruling.

## 4. R-3 ruling registration

| Element | Content |
|---|---|
| Ruled by | user |
| Date | 2026-11-15 |
| Reason | stage-1 closure is already reached, so a re-run is unnecessary; the root cause of S-G1 / S-G8 is "no real enforcement"; the P74 default of off would make a re-run report unsupported again; (B) does not advance S-G8; (A) costs nothing and disturbs no existing closure |
| Basis | the user's same-round ruling text of 2026-11-15 (including the five attached rulings A1-A5 and the O-7 = option-2 ruling); archived at `docs/validation/evidence/P75-R3-REAL-CLI-ASSESSMENT/02-batch/ruling-2026-11-15-r3.md` |
| Ruling | R-3 = (A) cannot start for now |

### 4.1 Attached rulings A1-A5

| Id | Content |
|---|---|
| A1 | T027-B declaration = "complete (partly not applicable)": the main line has landed, and R-3 is partly not applicable (real CLI end-to-end needs the S2 enforcement premise) |
| A2 | S-G8 stays blocked until S2 enforcement lands; the `s1-exit` gate table is unchanged; the S-G8 row of `p57` §4.2 gets an additive pointer |
| A3 | the closed AT-23 stage-1 state is registered as a permanent fact; stage 2 is not applicable; the S-G1 row of `p57` §4.2 is unchanged |
| A4 | three landing points for R-3: the T027-B section of `DEV_PLAN`, the P75 section of `REMAINING_SLICES`, and the R-3 row of `p57` §3.3 |
| A5 | add ADR-045 (a decision record; not bound by "must not be used as a precedent") |

### 4.2 Open-question ruling

- O-1 to O-6 and O-8: covered by this round's ruling and A1-A5; they are in-slice rulings.
- O-7: ruled = option 2 (no separate slice); the obligation is handed to the S1 exit closure slice (see §11).
- ⑥ skipped: this slice produces documents only (assessment, ruling registration, report assembly); the §7.2.1 three conditions judge it skippable, pending step-7 confirmation; if step 7 does not confirm, step 6 runs and this line is replaced by a `sixstage` fence.

## 5. T027-B completion decision

- Criterion: the milestone composition (the November refactor milestone list in `docs/DEV_PLAN.md`) plus the P57 §3.3 residue prongs (`docs/validation/p57-closeout-assessment.md:64-76`).
- Main line: C1-C7 and N1-N4 landed by P39-P44; R-1 closed by P69; R-2 ruled by P73 and landed as the check layer by P74.
- R-3 classification: **partly not applicable** - its carrier is "an end-to-end run on a real CLI with a real candidate", whose premise is S2 enforcement, absent inside S1.
- Conclusion (A1): T027-B is declared "**complete (partly not applicable)**"; this slice writes neither "not done" nor "S1 Exit passed".
- Boundary: this decision changes no `s1-exit` gate table row, no AT-23 verdict and no `unproven` marker.

## 6. S-G8 advance condition

- Criterion verbatim: the supervised path runs end-to-end on a real CLI and a real candidate (`docs/validation/p57-closeout-assessment.md:95`).
- Current: not passed, depends on P5 residue R-3; S-G8 unlocks together with S-G1.
- A2 ruling: S-G8 **stays blocked until S2 enforcement lands**; the `s1-exit{,_EN}.md` gate table is unchanged; the S-G8 row of `p57` §4.2 gets an additive pointer (landing point: `docs/validation/p57-closeout-assessment.md:99`).
- Three-way cross-check: under (A) / (B) / (C) alike this slice does not advance S-G8 - (B) at most re-confirms conditional closure, and (C) can only advance it if the swap succeeds and the real end-to-end run passes.
- Boundary: this slice changes neither the criterion nor the verdict of S-G8; a later slice that wants to change the criterion must stop and escalate for a user ruling.

## 7. Keeping the AT-23 stage-1 closure

- **Permanent fact registration** (A3): the AT-23 stage-1 conditional closure is a historical fact; once registered it is not revoked by the R-3 conclusion and not rewritten by later slices.
- P61 point (2026-11-08): premise 2 (the real-candidate leg ran start / execute / stop and recorded unsupported) was satisfied; premise 1 (mechanism leg) was not, due to the F2 inherited drift (`docs/validation/at23-restricted-candidate.md:19-20`).
- F2-FIX point (2026-11-09): step-7 final review `FINAL: PASS WITH FIXES` ⇒ premise 1 satisfied ⇒ **the AT-23 stage-1 conditional closure is reached** (`docs/validation/at23-restricted-candidate.md:423`, `:499`).
- Stage 2: not applicable (user ruling A1 of 2026-11-09); this slice does not change that ruling.
- Boundary: R-08.5 states that conditional closure is not a proof of real-candidate compatibility; the **gate-level verdict** of S-G1 (unreachable / not passed) and the **verdict state** of AT-23 are two different dimensions, and this slice does not change the S-G1 verdict (A3).

## 8. Effect on the S1 Exit closure point

- Eight-gate state (P57 point plus the P74 boundary): S-G1 and S-G8 not passed; S-G3 and S-G4 not passed; S-G5 a manual gate (unknown); S-G6 and S-G7 unreachable; S-G2 wording stale.
- Effect of R-3 = (A): **the closure point does not move earlier**; real progress on S1 Exit depends on S2 enforcement landing, not on a repeated assessment inside S1.
- Admission constraint (new): when the S1 exit closure slice starts, if the residual stale statement of §11.2 (item 3) is not yet handled, that becomes one of the slice's **admission conditions** - the closure slice must not declare S1 Exit passed while that residual item remains unhandled.
- Boundary: this slice gives no pass date for S1 Exit and changes no gate verdict.

## 9. Reachable path for real CLI end-to-end

- Fact chain: the S-G8 criterion demands a real end-to-end run → a compatible candidate is needed → B1 found no acceptable candidate inside the evaluated scope (`docs/validation/t027-b1-candidate-live-discovery.md:4`, `:111`) → the B2 external boundary is not wired into the product (`docs/validation/B2-EXTERNAL-ENFORCEMENT.md:12`, `:18-19`) → the candidate session fails inside a zero-capability container (same file `:162-164`).
- Reachability inside S1: **unreachable** - the P74 platform implementation is permanently `Unknown` and the trigger defaults to `off`; even `enforce` only refuses explicitly and is not enforcement.
- Reachable paths: (1) S2 enforcement landing plus a compatible candidate (main path); (2) a user ruling for (C) to swap the candidate family or introduce P-prime (which needs a separate ruling on the R-08.3 / R-08.4 exemption premises).
- Conclusion: this slice spends no probe and starts no end-to-end run; the S1-side "do not do it" is fixed by ADR-045.

## 10. Evidence gaps

- Available: the green mechanism-leg evidence, the unsupported records of the P61 real-candidate leg, the B1 / B2 conclusions, the P69 wiring evidence, the P74 check-layer evidence, and the AT-23 stage-1 closure decision.
- Still missing: real end-to-end evidence for the supervised path (including billing), compatible-candidate evidence, S2 real enforcement, and the evidence that unlocks S-G1 / S-G8.
- Conclusion: every gap closes only in S2 or through the user-ruled (C) path; S1 cannot fill them.

## 11. Open questions and follow-ups

### 11.1 O-1 to O-8 disposition

| Id | Disposition | Note |
|---|---|---|
| O-1 | in-slice ruling | R-3 = (A), see §4 |
| O-2 | in-slice ruling | lapses with (A) (the probe channel and budget wording does not apply) |
| O-3 | in-slice ruling | lapses with (A) (the candidate family and P-prime argument do not apply), kept as an S2 entry |
| O-4 | in-slice ruling | wording = "complete (partly not applicable)", see §5 |
| O-5 | in-slice ruling | S-G8 does not advance under any of the three rulings, see §6 |
| O-6 | in-slice ruling | ADR-045 is added at `docs/ADR_REGISTER.md:56` |
| O-7 | ruled = option 2 (no separate slice) | obligation hand-off in §11.2 and §11.3 |
| O-8 | in-slice ruling | the escalation boundary was not triggered, see §6 and §7 |

### 11.2 O-7 obligation hand-off (three items, not summarised)

p57 §4.3 is a point-in-time register of 2026-11-05; items 1 and 2 below were corrected and confirmed by the 2026-11-15 disk re-check, and item 3 is the residual.

| P75 §11.2 item | P57 §4.3 original item | Disposition |
|---|---|---|
| item 1 (the AT-24 row at `s1-exit.md:44`) | original item 1 | corrected by P64 ⇒ closed |
| item 3 (the core task-chain row at `s1-exit.md:52`) | original item 2 | the only residual ⇒ handed to the S1 exit closure slice |
| item 2 (row G-D of `s1-reach`) | original item 3 | corrected by s1-reach ⇒ closed |


1. Item 1 (the AT-24 row at `docs/validation/s1-exit.md:44`): **disk re-check on 2026-11-15 = already corrected** - the row now reads "conditionally closed (2026-11-11, P64 acceptance)", consistent with `COMPLETED` at `docs/DEV_PLAN.md:119` ⇒ closed, no hand-off.
2. Item 2 (row G-D of `docs/validation/s1-reach-reachability.md` §2): **disk re-check on 2026-11-15 = already corrected** - the row now cites `DEV_PLAN.md:119` and notes the s1-exit correction ⇒ closed, no hand-off.
3. Item 3 (the core task-chain row at `docs/validation/s1-exit.md:52`): **residual, the only hand-off** - the first half still says "T027 not done, AT-23 not passed" while the same cell later says "the AT-23 stage-1 conditional closure is reached" ⇒ inconsistent wording inside one row; fix timing = when the S1 exit closure slice restates the gate conclusions.

### 11.3 Admission condition and follow-ups

- Admission condition: when the S1 exit closure slice starts, if item 3 of §11.2 (the inconsistent wording inside the `s1-exit` core task-chain row) is not yet handled, that slice must not declare S1 Exit passed; the closure slice's briefing must cite the three re-check results of §11.2.
- Follow-ups: (1) the S2 enforcement wiring slice (a precondition of the real end-to-end run); (2) the S1 exit closure slice (carrying the §11.2 cleanup and the gate-conclusion restatement); (3) this slice starts no separate cleanup slice.
- Identified contract revision point (registered only): the wording mapping between R-08.5 "conditional closure" and the S-G8 criterion "real end-to-end run" is not yet written.

## 12. Review chain conclusion

- ① `DESIGN: DONE` (`deep-reasoner` + `DeepSeek V4 Pro`; includes a controller line-number verification note).
- ②a / ②b / ②c landed step by step; each step ran `gate.js --all --scope=tree` plus the CN/EN six-shape check green.
- ③ consistency check: controller self-executed (carrier deviation, see §14).
- ④ gates: the three scopes and SELFTEST all green (readings in §13).
- ⑤ pre-review: r1 `PASS WITH FIXES` (7 Medium + 2 Low) → fixes → r2 `PASS WITH FIXES` (3 Low) → fixes; the 1+1 quota is used up.
- ⑥ independent scan: skipped under the three §7.2.1 conditions (documents only / no machine-structure semantics / none of the four semantic classes), confirmed by ⑦ as `SKIP CONFIRMED`.
- ⑦ blind review: `BLIND REVIEW: PASS` (0 findings; the §7.7 proportional draw hit this slice and ran before the final review).
- ⑦ final review: `FINAL REVIEW: PASS` (0 findings); no Medium+ ⇒ no re-review.
- Blind and final agree (both PASS with 0 findings) ⇒ no root-cause arbitration needed.
- Metering: ⑦ = blind 1 + final 1 = 2, within the M-tier 1+2 plus the blind +1 total of 4; batch assistance 2 times (floor 2).

## 13. Gate readings and ruling landing points

- Gates: `--scope=index` / `--scope=tree` / `--scope=ci` all `TOTAL_FAIL=0`; `--selftest` = `SELFTEST: PASS 737/737`; base `go test ./...` = `ok=15 FAIL=0`.
- A1 landing: `docs/DEV_PLAN.md:203` (`_EN:204`).
- A2 landing: `docs/validation/p57-closeout-assessment.md:99` (additive S-G8 pointer; the `s1-exit` gate table is unchanged).
- A3 landing: §7 of this report (permanent fact registration); the S-G1 row of `p57` §4.2 is unchanged.
- A4 landings: `docs/DEV_PLAN.md:203`, the P75 section of `docs/t027/REMAINING_SLICES.md` (⑧ completion record), `docs/validation/p57-closeout-assessment.md:76`.
- A5 landing: `docs/ADR_REGISTER.md:56` (`_EN:57`).
- Metering: R = 1 (fix and re-run after the ⑤ r1 Medium+ report); D = 0 (no temp-file violation).

## 14. Known boundaries

- The four captures in `06-gates/` are point-in-time snapshots; the pack self-consistency lines reflect the declaration state at capture time, and the canonical readings are §13 and the ⑩ receipt.
- The ruling archive is an **excerpt**: the ruling body (the three-way ruling and the five attached rulings) is verbatim; the process-instruction sections are not transcribed and the archive header states the scope.
- A batch-assistance subagent created and immediately deleted a temporary file while counting characters (outside this slice's temp directory); zero residue was confirmed and it is not a waiver for other slices.
- The ③ consistency check was controller self-executed (carrier deviation); the independent impact assessment is in the report's exception section.

## 15. Explicitly not done

- No real CLI end-to-end probe was started (option (A) spends no probe).
- `docs/validation/s1-exit{,_EN}.md` and `docs/validation/s1-reach-reachability{,_EN}.md` were not modified.
- `CONTRACTS` / schema / fixtures / CI / product code were not modified.
- No ruling on the user's behalf; no commit / push; no new slice started.

### 15.1 Exceptions and deviation register

- P75-1 (fixes after the ⑤ r1 Medium+ report): the three O-7 items were carried over from p57 §4.3 without re-checking, two dangling references, a missing "does not do" axis, stale gate captures, and an unarchived ruling text ⇒ all fixed, with ⑤ re-run in place (r2). Counts as R = 1.
- P75-2 (gate-compliance fixes during authoring, not counted as R): the G8-b metadata row counted as payload, the G1(e2) completion row, the G7-f new-report declaration, the G6-4 expiring literal, the missing G5(d) fence row, and the G8-b self-referential hash; all are mechanical fixes inside the same authoring pass and were re-run green item by item.
- P75-3 (carrier deviation): the ③ consistency check was controller self-executed (§13.3.18); independent impact assessment = the action is fact comparison, and its conclusions were independently re-checked by ⑤ r2 and by both ⑦ readings; not to be treated as a precedent.
- P75-4 (timing deviation): the freedom-list was finalised later than the "after ⑥, before ⑦" window of §1.7.3 (finalised at ⑧); impact assessment = the list items are unchanged since step 0 and both ⑦ readings checked scope against them; not to be treated as a precedent.
- P75-5 (subagent temp file): a batch-assistance subagent created and immediately deleted a temp file while counting characters; zero residue confirmed; not a waiver for other slices.

## 16. Commit and CI fact chain (backfilled)

- Commit: `ed2e26e` (`docs: rule R-3 as not startable and declare T027-B complete in part`), 27 files `+2287 / -0`.
- Push: `origin/main` `c5c3a78..ed2e26e` (origin only; no gitee push).
- CI: run `43732384194`; `Go ubuntu-latest` success (3m19s), `Go windows-latest` success (4m41s); both legs `GATE REPORT scope=ci changed=27`, `TOTAL_FAIL=0`, `SELFTEST: PASS 737/737`.
- Note: this section was added by the docs-only commit after CI turned green (§11.2); the CI run of that backfill commit itself is an observation record only.

```artifacts
docs/ADR_REGISTER.md	repo
docs/ADR_REGISTER_EN.md	repo
docs/DEV_PLAN.md	repo
docs/DEV_PLAN_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/evidence/P75-R3-REAL-CLI-ASSESSMENT/00-briefing.md	repo
docs/validation/evidence/P75-R3-REAL-CLI-ASSESSMENT/01-arch/design-r1.md	repo
docs/validation/evidence/P75-R3-REAL-CLI-ASSESSMENT/02-batch/batch-01-t027b-at23-b1b2-facts.md	repo
docs/validation/evidence/P75-R3-REAL-CLI-ASSESSMENT/02-batch/batch-02-r08-r09-sandbox-facts.md	repo
docs/validation/evidence/P75-R3-REAL-CLI-ASSESSMENT/02-batch/ruling-2026-11-15-r3.md	repo
docs/validation/evidence/P75-R3-REAL-CLI-ASSESSMENT/04-reviews/blind-r1.md	repo
docs/validation/evidence/P75-R3-REAL-CLI-ASSESSMENT/04-reviews/final-r1.md	repo
docs/validation/evidence/P75-R3-REAL-CLI-ASSESSMENT/04-reviews/prereview-r1.md	repo
docs/validation/evidence/P75-R3-REAL-CLI-ASSESSMENT/04-reviews/prereview-r2.md	repo
docs/validation/evidence/P75-R3-REAL-CLI-ASSESSMENT/06-gates/gate-ci.txt	repo
docs/validation/evidence/P75-R3-REAL-CLI-ASSESSMENT/06-gates/gate-index.txt	repo
docs/validation/evidence/P75-R3-REAL-CLI-ASSESSMENT/06-gates/gate-selftest.txt	repo
docs/validation/evidence/P75-R3-REAL-CLI-ASSESSMENT/06-gates/gate-tree.txt	repo
docs/validation/evidence/P75-R3-REAL-CLI-ASSESSMENT/MANIFEST.md	repo
docs/validation/evidence/P75-R3-REAL-CLI-ASSESSMENT/SHA256SUMS.txt	repo
docs/validation/evidence/P75-R3-REAL-CLI-ASSESSMENT/freedom-list.md	repo
docs/validation/evidence/P75-R3-REAL-CLI-ASSESSMENT/input-list.md	repo
docs/validation/p57-closeout-assessment.md	repo
docs/validation/p57-closeout-assessment_EN.md	repo
docs/validation/r3-real-cli-assessment.md	repo
docs/validation/r3-real-cli-assessment_EN.md	repo
```
