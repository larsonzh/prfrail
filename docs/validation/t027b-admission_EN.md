# P73 `T027B-ADMISSION` — OB-55 admission assessment and ruling landing (change report)

Slice: `T027B-ADMISSION` (P73 · layer-2 governance · tier-M assessment `[SLICE]`). Controller `DeepSeek V4.1 Flash (deepseek)`; the user same-turn explicit authorisation of 2026-11-14 is carried by the kick-off message. Baseline `0cac174`; directive v1.69 (this slice bumped it before the main commit per §0.2, as a metadata completion).

## 1. Executive summary

- **OB-55 ruling = option b** (land the platform-neutral `Sandbox` interface early), the same direction as R-06 (2026-11-12), recorded as confirmation and refinement; all three conditions are accepted.
- **Landing surface**: a new P73 admission entry in the `DEV_PLAN` T027-B section; an additive pointer after item 8 of §8 in the B2 report; the OB-55 disposition cell flipped to "disposed"; a new ADR-044; three additive landing pointers; the P73 ledger section closed.
- **Boundary**: no production code, no change to `CONTRACTS` / schema / fixtures / CI, no change to the S1 Exit eight gates themselves; the implementing slice `P74-SANDBOX-INTERFACE` is not started by this slice.

## 2. Current-state check

| Item | Fact |
|---|---|
| C1–C7 and N1–N4 | the P38 working drawing gives the verdicts (C1 / C4 / C6 = BLOCKED, C2 / C3 / C5 / C7 = PARTIAL, plus N1–N4); landing is carried by P39–P44 and the P38 text is not rewritten to the landed state. |
| P39–P44 | land C6 / C1 / C2+N1 / C3+N2 / C4+C5+N3 / C7+N4 in order; P44 limits itself to "must not claim end-to-end closure". |
| P69 | wires the P44 readiness / permit gates into the production assembly and the CLI entry; the noop path and the `ai check` path are not wired. |
| OB-55 | the disposition cell read "open item", with candidate scopes = standalone slice / land the `Sandbox` interface early / keep S2. |
| item 8 of §8 in the B2 report | states that the `Sandbox` interface is not built in S1, because landing it early would create dead code with no caller. |
| R-06 | on 2026-11-12 the user ruled option b and required the P44 end-to-end wiring R-1 to be closed first. |
| S1 Exit eight gates | S-G1 / S-G6 / S-G7 are unreachable or failing, S-G3 / S-G4 / S-G8 fail, S-G5 is a human gate, and the S-G2 wording already follows the AT-24 conditional closure. |
| `DEV_PLAN:153` | still reads "R-2 OB-55 ownership (pending the R-06 ruling)", a P57-time statement that predates the R-06 ruling. |

Key finding: R-06 already fixed the direction at option b and its sequencing condition (close R-1 first) was closed by P69; what actually remained for ruling was the implementation boundary, the ban disposal, the implementing slice tier and the record form.

## 3. Three-way comparison for OB-55

| Option | What it does | Cost | Main risk | Rollback |
|---|---|---|---|---|
| 1 standalone slice for the full part B | in-container tool execution path plus the assembly enforcement check plus policy wiring | tier M/L over months | the candidate cannot complete a session inside a zero-capability container, so there may be no real call chain to verify; full regression of the refusal path | revertible, but a long chain |
| 2 land the `Sandbox` interface early | contract first plus a platform-neutral interface plus the Windows wrapper plus a fail-closed empty implementation plus a same-slice call site | tier S/M, zero probes | literal conflict with the dead-code ban in item 8 of §8, which must be disposed first | purely additive, removable as a whole |
| 3 keep S2 | only record the conclusion in the OB-55 disposition cell | 0 | the conflict persists to S2 and the interface would be built from scratch then | nothing to roll back |

Relation to R-06: the R-06 option b text ("land the platform-neutral `Sandbox` interface early") is literally the same as the second candidate scope of OB-55; what was ruled already = direction and sequencing, what remained = implementation boundary, ban disposal, implementing slice tier and record form.

## 4. User ruling registration

- **Ruling party**: the user. **Date**: 2026-11-14. **Reason**: the same direction as R-06 plus the eight-step assessment in ① plus the three conditions. **Basis**: this round's ruling text; the registration is in `02-ruling/ob55-ruling.md` and the ruling message itself remains the source text.
- **Ruling 1**: OB-55 = option b; all three conditions are accepted (same-slice landing, ban disposed first, the report states the gates are not unlocked).
- **Ruling 2**: item 8 of §8 in the B2 report is disposed of by an additive ruling note; the original line stays untouched (§11.1).
- **Ruling 3**: the implementing slice is the tier-M standalone slice `P74-SANDBOX-INTERFACE`, not started by P73.
- **Ruling 4**: the call sites are the AgentRunner assembly plus the supervised fail-closed empty implementation, on the same layer as the P69 permit.
- **Ruling 5**: "part B does not unlock S-G1 / S-G8" goes into this report and `DEV_PLAN`; `s1-exit` stays untouched.
- **Ruling 6**: dual record = the OB-55 disposition cell flipped to "disposed" plus ADR-044; it contains a rule-semantics change, is a direct user ruling and must not be treated as a precedent.
- **Ruling 7**: the AT-23 re-review (R-3) and part B have no mandatory ordering and proceed independently.

## 5. Part-B implementation boundary

- **To change**: a new Sandbox section in `CONTRACTS` §7 (contract first); `internal/guard/sandbox.go`, `sandbox_windows.go`, `sandbox_unsupported.go`; the assembly check point (same layer as the P69 permit); the matching tests.
- **Not changed**: the criterion bodies of `internal/chain` and `internal/adapters`, schema, fixtures, CI, the candidate version, the R-08 / R-09 rulings and the refusal semantics.
- **Call-site scope**: the AgentRunner assembly plus the supervised path (fail-closed empty implementation), on the same layer as the P69 permit; the noop path and the `ai check` path are not wired.
- **Interface freeze**: the `internal/guard` location and the three file names are frozen by item 8 of §8 in the B2 report; the gate-layer order is frozen by what P42 / P43 / P44 / P69 already landed.
- **Position**: **part B does not unlock S-G1 / S-G8** - the gate-level verdicts do not change with the interface landing (S-G1 rests on the AT-23 stage-1 conditional closure, and S-G8 shares that origin).

## 6. Candidate version and Windows / P

- **Candidate version**: stays at GitHub Copilot CLI **1.0.83** (frozen by R-08.1); the `Sandbox` abstraction depends on neither and changes neither.
- **Windows / P**: keeps the R-08.3 "P only" premise (process-level crash, no power loss, local Win11 / NTFS); **P-prime is not introduced**; P-prime belongs to the B3b power-loss rig premise and any reference must carry the "P-prime covering a physical target deployment is unproven" note.

## 7. AT-23 and refusal-path mapping

- **Execution path**: assembly entry -> the Step 0 permit verdict (wired by P69) -> the Sandbox check (empty implementation, fail-closed) -> refusal with zero traces; no new "pass" branch is added.
- **AT-23 coverage**: part B adds no AT-23 evidence and weakens no existing conclusion; the AT-23 stage-1 conditional closure is unaffected.
- **Refusal path**: semantically equivalent; the reason wording moves from "no enforcement record" to "Sandbox unsupported" while the fail-closed behaviour is unchanged.

## 8. Evidence gaps

- **Available**: the B2 boundary measurements, the P44 verdict implementation, the P69 wiring evidence, the P39-P43 landing evidence, the AT-23 stage-1 conditional closure and the AT-24 conditional closure.
- **Missing**: the Sandbox contract section, the call-site wiring evidence, the dead-code-ban disposal record (landed by this slice) and the implementing slice ⑦ review record; S-G1 / S-G6 / S-G7 gain no unlocking evidence from part B.

## 9. Resource and probe authorisation needs

- **This slice**: tier-M assessment, zero probes and zero production code; ⑤ pre-review 1+1; ⑦ 1+2 (no blind review).
- **The P74 implementing slice**: tier M; **probe demand = 0** (the interface and the empty implementation call no real candidate); the R-08.2 stage-1 "declare per use" probe allowance applies only to the R-3 re-review and is handled separately.

## 10. Open questions and follow-ups

- The AT-23 re-review (R-3) and part B have **no mandatory ordering**: R-3 does not depend on the `Sandbox` interface, so part B neither enables nor blocks it.
- The `P74-SANDBOX-INTERFACE` kick-off must state the call-site scope and the same-slice landing constraint up front.
- The Sandbox section in `CONTRACTS` §7 is a contract-semantics change and must follow that contract's own flow (RFC / schema / samples / tests) inside P74; this slice does not land it.
- ADR-044 records the dead-code-ban disposal and the implementation boundary; any later boundary change needs a new ADR or OB entry rather than an edit of ADR-044.
- The stale wording at `DEV_PLAN:153` ("pending the R-06 ruling") is superseded in meaning by the new T027-B entry; the original line stays untouched per §11.1.
- The `s1-exit` gate table is untouched (ruling 5); any later change in how the S-G1 position is adopted is carried by the gate-level reports.
- This slice does not start `P74` and makes no ruling on later slice content.

## 11. Review-chain conclusion

- ① `DESIGN: DONE` (`design-r1`; eight-step assessment plus the three-way comparison plus recommendation 2 plus the six assessment items plus the per-item stop-boundary verdicts).
- ⑤ pre-review: `PRE-REVIEW: PASS` (r2 re-check; original in `04-reviews/prereview-r1.md`).
- ⑥ skipped: documentation-only assessment and ruling landing with no code / schema / fixtures / CI change, so the three §7.2.1 conditions hold, confirmed by ⑦.
- ⑦ final review: `FINAL REVIEW: PASS` (0 findings plus 1 Info synced; blind review 0; ⑥ `SKIP CONFIRMED`).

## 12. Gate readings

```text
gate.js --selftest         : SELFTEST: PASS 737/737
gate.js --all --scope=index: TOTAL_FAIL=0
gate.js --all --scope=tree : TOTAL_FAIL=0
gate.js --all --scope=ci   : TOTAL_FAIL=0
```

## 13. Known boundaries

- The "not built in S1" wording in item 8 of §8 stays as-is; an additive pointer records that R-06 and P73 updated it.
- The time-point verdicts in the P38 working drawing and the P57 closeout assessment are not rewritten to the landed state (historical reports stay untouched, §11.1).
- The stale wording at `DEV_PLAN:153` is kept and superseded in meaning by the new entry.
- The change to `DELIVERY_DIRECTIVE` is limited to the OB-55 disposition cell in §12.4, with the version bumped to **v1.69** per §0.2 and the historical precedents (header version line plus a §0.3 log row plus the appendix C.0av record block); this is a metadata completion and changes no rule semantics.
- The version was bumped to v1.69 before the main commit (header version line plus a §0.3 log row plus the appendix C.0av record block) as a metadata completion; `G1-b` / `G3(a.2)` / `G6-5` / `G6-7` all green.

## 14. Explicitly not executed

- No change to production code / `CONTRACTS` / schema / fixtures / CI; no change to the S1 Exit eight gates themselves; no change to the `s1-exit` gate table.
- No change to the v50-v57 snapshots; no change to the body of the P38 / P57 and other historical reports (additive pointers only).
- `P74-SANDBOX-INTERFACE` is not started; no ruling is made on later slice content.
- No commit / push (stopped at ⑩ pending authorisation); no new slice started.


## 15. Commit and CI fact chain

- Main commit: `8c59712` (assessment plus ruling landing plus the directive ledger cell and the v1.69 metadata; English subject `feat: rule the OB-55 admission as option b and add the sandbox interface to S1`).
- CI: run `43622770750` (ubuntu-latest 3m15s / windows-latest 4m44s, both legs `success`).
- Both legs: `SELFTEST: PASS 737/737`; `TOTAL_FAIL=0`.
- This section is written by the backfill commit and records facts only; no business content changes.

```artifacts
docs/DEV_PLAN.md	repo
docs/DEV_PLAN_EN.md	repo
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/ADR_REGISTER.md	repo
docs/ADR_REGISTER_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/B2-EXTERNAL-ENFORCEMENT.md	repo
docs/validation/B2-EXTERNAL-ENFORCEMENT_EN.md	repo
docs/validation/p57-closeout-assessment.md	repo
docs/validation/p57-closeout-assessment_EN.md	repo
docs/validation/remaining-rulings.md	repo
docs/validation/remaining-rulings_EN.md	repo
docs/validation/strategic-planning-2026-10-31.md	repo
docs/validation/strategic-planning-2026-10-31_EN.md	repo
docs/validation/t027b-admission.md	repo
docs/validation/t027b-admission_EN.md	repo
docs/validation/evidence/P73-T027B-ADMISSION/00-briefing.md	repo
docs/validation/evidence/P73-T027B-ADMISSION/input-list.md	repo
docs/validation/evidence/P73-T027B-ADMISSION/freedom-list.md	repo
docs/validation/evidence/P73-T027B-ADMISSION/MANIFEST.md	repo
docs/validation/evidence/P73-T027B-ADMISSION/SHA256SUMS.txt	repo
docs/validation/evidence/P73-T027B-ADMISSION/01-arch/design-r1.md	repo
docs/validation/evidence/P73-T027B-ADMISSION/02-batch/batch-01-p36-p69-facts.md	repo
docs/validation/evidence/P73-T027B-ADMISSION/02-batch/batch-02-rulings-gates.md	repo
docs/validation/evidence/P73-T027B-ADMISSION/02-ruling/ob55-ruling.md	repo
docs/validation/evidence/P73-T027B-ADMISSION/04-reviews/prereview-r1.md	repo
docs/validation/evidence/P73-T027B-ADMISSION/04-reviews/final-r1.md	repo
docs/validation/evidence/P73-T027B-ADMISSION/06-gates/gates-r1.txt	repo
docs/validation/evidence/P73-T027B-ADMISSION/06-gates/selftest-r1.txt	repo
```
