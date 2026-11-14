# P72 `REQ020-S2-RESTRICTION` — registering the S2 commitment narrowing for REQ-020 (change report)

Slice: `REQ020-S2-RESTRICTION` (P72 · layer-2 governance · documentation `[SLICE]` · S tier). Controller `DeepSeek V4.1 Flash (deepseek)`; the user's same-turn explicit authorisation of 2026-11-14 is carried by the kick-off message (§1.7.2). Baseline `70eccee` (directive v1.68).

## 1. Executive summary

- **Goal**: change the REQ-020 stage column from `S0-S2` to "Not required by S1 (C+Go contract only); S2 commitment narrowed (Python/Node official harness + three-language product tasks + the release-support commitment frozen; existing tool scripts / open language IDs / generic capability are not frozen)", and sync AT-01 in `TEST_STRATEGY` plus the generic/C/Go row of `S1_SUPPORT_MATRIX`.
- **Conclusion**: closed. The REQ entry count stays 28 with contiguous numbering; the REQ-020 "testable requirement" column is untouched; the MVA composite stays unchanged; the S1 Exit eight gates themselves are unchanged; `S1_SUPPORT_MATRIX` keeps 6 rows (no standalone Python/Node row added).
- **Change surface**: 6 document pairs (REQ / AT / MATRIX / ledger) plus 2 additive pointer lines (`remaining-rulings` §5.4, `s1-scope-annotate` §7 item 3).
- **Stop-and-escalate boundaries**: all 8 checks are not triggered; (h) adopts the narrow reading — REQ-011's "extend languages in S2" belongs to "open language IDs / generic capability", which the narrowing text explicitly does not freeze.

## 2. Difference list

| Document | Current state vs the Q3 freeze requirement | Nature |
|---|---|---|
| `PRODUCT_REQUIREMENTS{,_EN}.md` | REQ-020 stage column is `S0-S2` with no narrowing note | modify (column 3 only) |
| `TEST_STRATEGY{,_EN}.md` AT-01 | scenario column `C+Go/三语言` (EN `dual/triple-language`) carries no stage ownership | modify (that fragment only) |
| `TEST_STRATEGY{,_EN}.md` AT-11 | `C+Go 同 task 整体恢复` is S1 recovery semantics | unchanged |
| `S1_SUPPORT_MATRIX{,_EN}.md` row 12 | constraint column has no three-language commitment note | additive note |
| `t027/REMAINING_SLICES{,_EN}.md` P72 section | status and completion record | modify (current: 🔄 already landed at step 0; planned: ✅ plus the completion record when ⑧ completes) |
| `remaining-rulings{,_EN}.md` §5.4 | no landing pointer for the Q3 freeze | additive 1 line |
| `s1-scope-annotate{,_EN}.md` §7 item 3 | the follow-up item is not yet discharged | additive 1 line |

## 3. Per-document write-back plan

- **REQ-020 stage column**: EN = `Not required by S1 (C+Go contract only); S2 commitment narrowed (Python/Node official harness + three-language product tasks + the release-support commitment frozen; existing tool scripts / open language IDs / generic capability are not frozen)`; CN = `S1 不要求（仅 C+Go 契约）；S2 承诺收窄（Python/Node 官方 harness + 三语言产品任务 + 正式支持承诺冻结；不冻既有工具脚本 / 开放语言 ID / generic 能力）` (terms taken from the existing translation in the ledger's `REMAINING_SLICES_EN.md` goal field).
- **AT-01**: `dual/triple-language` → `C+Go contract (three-language execution frozen to S2, see REQ-020)`; CN `C+Go/三语言` → `C+Go 契约（三语言实跑冻结至 S2，见 REQ-020）`; all other cells byte-identical.
- **AT-11**: judged unchanged — `C+Go whole-task recovery` is an S1 recovery assertion, not a support-commitment stage ownership.
- **S1_SUPPORT_MATRIX row 12**: append `Python/Node official harness + three-language product tasks + the release-support commitment frozen to S2 (REQ-020 S2 commitment narrowing); existing tool scripts / open language IDs / generic capability are not frozen` to the constraint column; the "S1 candidate status" column stays "Configuration and noop orchestration".
- **Five assessment items**: (1) `remaining-rulings` §5 = applicable (additive pointer, following the P71 precedent); (2) `S1_RELEASE_NOTES` = not applicable (the support section names neither Python/Node official support nor three-language execution); (3) `s1-exit` AT-01 row = not applicable (a `PASS` summary row, not a stage-ownership statement); (4) `p57-closeout-assessment` §4.2 = not applicable (gates S-G1..S-G8 contain no language-support gate); (5) `s1-scope-annotate` §7 item 3 = applicable (that follow-up is this very slice; discharged by an additive pointer).

## 4. Counter-examples and verification

| Counter-example | Mechanical detection | Measured |
|---|---|---|
| changing REQ-020 column 4 | split-column comparison of field 4 | original 103-char value unchanged |
| adding or removing REQ rows / renumbering | `^\| REQ-\d{3} ` count | CN=28 / EN=28 |
| adding a standalone Python/Node row | `^\| (Python\|Node)` count | 0 / 0; in-table rows 7 (header + 6) |
| changing the harness row status column | field 2 comparison | `Configuration and noop orchestration` unchanged |
| widening the freeze (dropping the "not frozen" clause) | byte comparison with the target row | the clause is present |

## 5. Review-chain conclusion

- ① `DESIGN: DONE` (`design-r1`; includes the (h) narrow-reading ruling and the eight stop-boundary verdicts).
- ⑤ pre-review: `PRE-REVIEW: PASS WITH FIXES` (2 Low + 4 Info, all addressed; see §7 and `04-reviews/prereview-r1.md`).
- ⑥ skipped: documentation-only write-back with no code / schema / fixtures / CI change, so the three §7.2.1 conditions hold; confirmed by ⑦ (`INDEPENDENCE GATE (⑥): SKIP CONFIRMED`).
- ⑦ final review: `FINAL REVIEW: PASS WITH FIXES` (1 Low + 1 Info, both fixed; see `04-reviews/final-r1.md`); blind review 0 (P71 already consumed the slot, so the window restarted and P72 is the first of four, hence not drawn).

## 6. Gate readings

```text
gate.js --selftest         : SELFTEST: PASS 737/737
gate.js --all --scope=index: TOTAL_FAIL=0
gate.js --all --scope=tree : TOTAL_FAIL=0
gate.js --all --scope=ci   : TOTAL_FAIL=0
```

## 7. Known boundaries

- The REQ-020 stage column no longer spells out `S0`; the "S0 dual/triple-language configs" fact stays in column 4 (the zero-change invariant).
- `S1_SUPPORT_MATRIX` has no standalone C / Go row; the harness row (row 12) is the only language row, so the note lands in that row's constraint column.
- The `C+Go` wording in AT-11 is preserved (recovery semantics) and is not rewritten as stage ownership.
- ⑤ fix items: (1) the section reference in report §3 was corrected from "§5" to "§7 item 3" (the same typo in `design-r1` §6 stays untouched as evidence and is only recorded here); (2) §1 / §5 review verdicts are not pre-written and are finalised at ⑧ from the actual results; (3) the §2 ledger row now reads "🔄 already landed at step 0".
- The "not frozen" enumeration in the narrowing text follows the slice card's goal wording with three items (existing tool scripts / open language IDs / generic capability), one term short of the P68 freeze text ("golden configs"); it is a non-exhaustive enumeration and does not widen the freeze, so the slice card governs.
- No "slice sync (P72 ...)" tail note was added to `S1_SUPPORT_MATRIX`: the in-row note already states the same position, and design did not authorise that new line, so the change stays minimal.
- The evidence pack `docs/validation/evidence/P72-REQ020-S2-RESTRICTION/` is part of this slice's audit surface (deliverable 2 of the kick-off message's section 5; the kick-off message itself is an artifact under §13.1.10) and is not a product-surface change; this sentence fixes the boundary wording raised as item 1 of the ⑦ r1 review.

## 8. Explicitly not executed

- No change to the directive body / criterion bodies / `CONTRACTS` semantics / schema / fixtures / CI / the v50-v57 snapshots.
- No change to the REQ table structure (no entry added or removed, no renumbering, the testable-requirement column untouched), no change to the MVA composite, no change to the S1 Exit eight gates themselves.
- No standalone Python/Node row added to `S1_SUPPORT_MATRIX`; no ruling made on the user's behalf; no commit / push; no new slice started.

- Carrier deviation (§1.7.4 record): the write-back and the consistency check in ② / ③ / ⑧ were executed by the controller (no prfrail-documenter dispatch), because they are whole-line replacements plus mechanical checks whose dispatch cost outweighs the benefit; the deviation changes neither the change surface nor the criteria, is recorded in the ⑩ receipt, and must not be treated as a precedent.
```artifacts
docs/PRODUCT_REQUIREMENTS.md	repo
docs/PRODUCT_REQUIREMENTS_EN.md	repo
docs/TEST_STRATEGY.md	repo
docs/TEST_STRATEGY_EN.md	repo
docs/S1_SUPPORT_MATRIX.md	repo
docs/S1_SUPPORT_MATRIX_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/remaining-rulings.md	repo
docs/validation/remaining-rulings_EN.md	repo
docs/validation/s1-scope-annotate.md	repo
docs/validation/s1-scope-annotate_EN.md	repo
docs/validation/req020-s2-restriction.md	repo
docs/validation/req020-s2-restriction_EN.md	repo
docs/validation/evidence/P72-REQ020-S2-RESTRICTION/00-briefing.md	repo
docs/validation/evidence/P72-REQ020-S2-RESTRICTION/input-list.md	repo
docs/validation/evidence/P72-REQ020-S2-RESTRICTION/freedom-list.md	repo
docs/validation/evidence/P72-REQ020-S2-RESTRICTION/MANIFEST.md	repo
docs/validation/evidence/P72-REQ020-S2-RESTRICTION/SHA256SUMS.txt	repo
docs/validation/evidence/P72-REQ020-S2-RESTRICTION/01-arch/design-r1.md	repo
docs/validation/evidence/P72-REQ020-S2-RESTRICTION/02-batch/batch-01-req-enumeration.md	repo
docs/validation/evidence/P72-REQ020-S2-RESTRICTION/02-batch/batch-02-at-matrix-enumeration.md	repo
docs/validation/evidence/P72-REQ020-S2-RESTRICTION/04-reviews/prereview-r1.md	repo
docs/validation/evidence/P72-REQ020-S2-RESTRICTION/06-gates/gates-r1.txt	repo
docs/validation/evidence/P72-REQ020-S2-RESTRICTION/06-gates/selftest-r1.txt	repo
docs/validation/evidence/P72-REQ020-S2-RESTRICTION/04-reviews/final-r1.md	repo
```
