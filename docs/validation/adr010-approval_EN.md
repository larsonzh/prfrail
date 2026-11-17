# P77 Validation Report: ADR-010 approval and the EOL notification policy sync

Date: 2026-11-17. Status: `[SLICE]`, documentation type, tier S. Slice id: `ADR010-APPROVAL` (P77).

## 1. Slice identity and authorization

- Goal: after the product owner approves ADR-010 (the EOL notification policy), land the documentation - flip ADR-010 to `APPROVED` and sync the now-effective policy to every related document.
- Four elements of the approval: ruling party = the user (product owner); date = 2026-11-16; reason = the EOL notification policy must be frozen before v0.1.0 enters the release flow so S-G4 can close; basis = the user's same-turn P77 kickoff ruling of 2026-11-16.
- Baseline: `HEAD` = `origin/main` = `fb90a20`; directive v1.69; `SELFTEST` 737/737; the three scopes report `TOTAL_FAIL=0`.
- R / D baseline: R = 0, D = 0 (step-0 registration, §13.1.1).

## 2. Boundary and no-change confirmation

- The S1 Exit eight gate criteria are unchanged; the `p57-closeout-assessment{,_EN}` and `s1-exit{,_EN}` gate tables are untouched.
- Item 1 of the `s1-exit{,_EN}` residual-risk section is narrative (not the gate table) and is synced per the P76 report §8 "post-approval change points" list; the dated 2026-09-09 conclusion paragraph in the same file stays untouched per §11.1 (see §8).
- The EOL notification policy itself is unchanged: columns 2 and 3 of ADR-010 keep their exact bytes.
- The `S1_SUPPORT_MATRIX` platform list and support tiers are unchanged; the INSTALLATION install model is unchanged.
- No external GitHub Release notes action (that belongs to P78 `RELEASE-ZIP-BINDING`).
- No other ADR is changed; no ADR or OB is added.
- The directive is not bumped (stays at v1.69); `CHECKLIST_v59/v60{,_EN}` and the v50-v58 snapshots are untouched.
- No change to `.go` / `.json` / `.yml` / schema / fixtures / CI / `go.mod`; no gitee push.

## 3. Per-file change evidence

| File | Location | Nature | CN delta | EN delta |
|---|---|---|---|---|
| `docs/ADR_REGISTER.md` / `_EN` | ADR-010 status cell | replace | `+3/-3` | `+3/-3` |
| `docs/ADR_REGISTER.md` / `_EN` | support-policy paragraph tail | replace | counted above | counted above |
| `docs/ADR_REGISTER.md` / `_EN` | decision-table platform/support cell | closing note | counted above | counted above |
| `docs/S1_SUPPORT_MATRIX.md` / `_EN` | tail policy section, first and last item | replace | `+2/-2` | `+2/-2` |
| `docs/INSTALLATION_PLAN.md` / `_EN` | EOL "still open" clause | replace | `+1/-1` | `+1/-1` |
| `docs/PRODUCT_REQUIREMENTS.md` / `_EN` | §7 open-details sentence | additive pointer | `+1/-1` | `+1/-1` |
| `docs/validation/s1-exit.md` / `_EN` | residual-risk item 1 | replace | `+1/-1` | `+1/-1` |
| `docs/t027/REMAINING_SLICES.md` / `_EN` | P77 ledger section | append | `+20/-0` | `+20/-0` |

### 3.1 The ADR-010 status cell, before and after

Before (CN):

> `` `PROPOSED / NOT APPROVED`；90 天适合当前早期原型，通知渠道与演进规则须在首版发布前获产品所有者明确批准 ``

After (CN):

> `` `APPROVED`（裁决人 = 用户；日期 = 2026-11-16；理由 = v0.1.0 进入 Release 流程前冻结 EOL 通知方式；依据 = 用户 2026-11-16 P77 起跑消息同轮裁定原文）；90 天适合当前早期原型，通知渠道与演进规则须在首版发布前获产品所有者明确批准（**已于 2026-11-16 落实；P77**） ``

EN side mirrors the CN structure: the status word is `APPROVED`, the four elements use the same English sentence pattern, the historical completion sentence keeps its exact bytes, and the closing note mirrors CN.

### 3.2 Verdicts for the six candidate documents and two in-file statements

| Document | Verdict | Basis |
|---|---|---|
| `S1_SUPPORT_MATRIX{,_EN}` | tier 1, changed (two spots) | carried the "pending" and "not approved" wording |
| `INSTALLATION_PLAN{,_EN}` | tier 1, changed | the sentence listed EOL notice alongside two other "still open" items |
| `PRODUCT_REQUIREMENTS{,_EN}` | tier 2, triggered, additive pointer | the open-details list contained the EOL notice item |
| `S1_RELEASE_NOTES{,_EN}` | not touched (downgraded from tier 2) | the line is a binding obligation, not an "unapproved" claim; approval is not binding completion |
| `INSTALLATION{,_EN}` | not touched | no EOL statement anywhere |
| `OPERATIONS{,_EN}` | not touched | no EOL statement anywhere |
| `ADR_REGISTER` support-policy paragraph tail | tier 1, changed | a standing claim that directly contradicted the approval |
| `ADR_REGISTER` decision-table platform/support cell | tier 1, closing note added | the cell is a "current constraint" status column |
| `DEV_PLAN{,_EN}` | not changed (boundary exclusion, disclosed) | all three are dated historical records and the file is explicitly out of bounds |

### 3.3 Combined search sweep

- Five patterns were searched separately; the union scope is everything outside the evidence tree (hits inside `docs/validation/evidence/**` belong to this slice's own transcripts and payloads and are not counted as reach); the union covers twenty-three files.
- Four buckets: changed here, registered here (ledger), history frozen, and boundary-excluded with disclosure. The per-file mapping is in the ① design, section Q4.
- Point-in-time statements in the frozen files (`CHECKLIST_v59/v60{,_EN}`, `p57-closeout-assessment{,_EN}`, `s1-exit-closeout-assessment{,_EN}`, `s1-reach-reachability{,_EN}`) are preserved and trace back to the ADR-010 status cell.

### 3.4 Differences between the ② landed text and the ① target text

| Location | ① design target text | ② landed text | Basis for the choice |
|---|---|---|---|
| ADR-010 status cell | bold status word plus a "Four elements:" prefix plus nature / citable / effective-date / rollback tail notes | the §4 ②a reference format of the kickoff message: backticked status word plus parenthesized four elements plus the historical sentence plus a bold closing note | the §4 ②a reference format is authoritative for ②; nature-style tail notes stay out of the status cell, the S-G4 verdict is carried by §5 of this report and the four elements and date are in §1, avoiding duplication with ADR-045 / ADR-046 / ADR-047 |
| the `S1_SUPPORT_MATRIX` first and last items | the first item reads "Approved ADR-010 (2026-11-16):"; the last item keeps "this page does not replace the final notice" | the first item becomes a bold "ADR-010 (approved 2026-11-16)"; the last item becomes "this page does not replace the formal notice in the Release notes" | §4 ②a items 1 and 2 of the kickoff message state these forms; the general limitation is kept and only the "not approved" layer is removed |
| the §7 open-details sentence in `PRODUCT_REQUIREMENTS` | a minimal in-sentence split | an additive pointer at the sentence tail | the §4 ②b item 4 asks for an additive pointer, and an addition never rewrites an existing sentence (§11.1) |
| the dated 2026-09-09 conclusion paragraph in `s1-exit` | listed as a change point | not changed | that paragraph is a dated historical conclusion and §11.1 forbids rewriting it; the minimal list in P76 report §8 names only residual-risk item 1 |

- Note: the "status cell landed per the Q1 target text" wording in §Q3 of the ① design was forward-looking when written; ② landed the texts in the table above, and this section supersedes that wording.

## 4. Acceptance evidence and gates

| Check | Result |
|---|---|
| The ADR-010 status cell starts with `APPROVED` | yes (true on both sides) |
| All four elements present | yes (ruling party / date / reason / basis) |
| The historical completion sentence is kept byte-for-byte with a closing note | yes (one assertion pair per side) |
| Zero hits for the "not approved yet" phrase in `S1_SUPPORT_MATRIX` | yes (zero on both sides) |
| CN and EN changed-line symmetry | yes (both sides agree per file pair) |
| The three scopes and `SELFTEST` | see §7.2 |

## 5. S-G4 closeout

Criterion (unchanged): the S-G4 row in §4.2 of `p57-closeout-assessment{,_EN}` reads "EOL notification method | ADR-010 approved by the product owner"; this slice makes zero change to that row and its gate table.

- Evidence: the status cell at `docs/ADR_REGISTER.md:20` now carries the after text in §3.1 (EN side `docs/ADR_REGISTER_EN.md:20`); this report pair carries the conclusion; the four approval elements are in §1.
- Conclusion: **S-G4 moves from not passed to passed**; this report carries the verdict, while the gate tables in `p57-closeout-assessment{,_EN}` and `s1-exit{,_EN}` stay untouched per the P76 ruling.
- Downstream: P78 `RELEASE-ZIP-BINDING` is unblocked (S-G4 is one of its dependencies); S-G3 is judged after P78 lands.
- Unchanged: this approval must not be read as S1 Exit passing; S-G5 remains a human gate; the core task-chain caliber and the S-G2 clause in `s1-exit` stay with a later slice.

## 6. Downstream unblock and residue

- Consistent with ADR-046 Q10 = (C) two-step approval: this slice completes the ADR-010 part of step one; the v0.1.0 Release part of step one belongs to P78; step two (S1 Exit as a whole) comes after S2.
- The other three bullets in the `S1_SUPPORT_MATRIX` tail policy section (three-part versioning, the ninety-day window and review, and the dated-ADR deviation rule) keep their exact bytes and agree with ADR-010 columns 2 and 3.

## 7. Metering and review chain

| Stage | Count |
|---|---|
| ① architecture | one run (`deep-reasoner` + `DeepSeek V4 Pro`) |
| Batch assistance | two runs (`GLM-5.3-Flash`) |
| Probes | 0 |

### 7.1 Review-chain results

- ① architecture: `DESIGN: DONE`; the design is `01-arch/design-r1.md` in the evidence pack.
- ⑥ skipped: the three §7.2.1 conditions hold (documentation and ledger files only; no code, schema or fixture change; no state-machine, gate-criterion or authorization-contract semantics), confirmed by ⑦.
- ⑤ pre-review: two rounds - r1 `PRE-REVIEW: PASS WITH FIXES` (one medium, three low); r2 `PRE-REVIEW: PASS WITH FIXES` (no medium, three low); budget 2/2 exhausted, and the three low-risk fixes from r2 are confirmed by ⑦.
- ⑦ final review: one round - `FINAL REVIEW: PASS WITH FIXES` (no medium or above, one low; `SKIP CONFIRMED` included); budget used 1 of 2.
- ⑥ skip confirmation (made by ⑦): covered = the change-surface boundary and file types, bilingual symmetry and report shape, mechanical gate convergence, and the S-G4 carrier chain; omitted = the post-push CI run state, a byte-level comparison with the user's original kickoff message, and runtime behavior validation (this slice touches no code execution semantics).

### 7.2 Gate readings

| Check | Reading |
|---|---|
| `--scope=index` | `TOTAL_FAIL=0` |
| `--scope=tree` | `TOTAL_FAIL=0` |
| `--scope=ci` | `TOTAL_FAIL=0` (the local reading covers the previous commit range; the real post-push reading is in the ⑩ receipt) |
| `--selftest` | `SELFTEST: PASS 737/737` |

## 8. Known limits and residual differences

- The "state and evidence" cell of S-G4 in §4.2 of `p57-closeout-assessment{,_EN}` still quotes the old status text; the gate table stays untouched per the P76 ruling, so the `ADR_REGISTER` status cell is authoritative.
- The dated 2026-09-09 conclusion paragraph in `s1-exit{,_EN}` still says the EOL notification policy is not frozen; it is a dated historical conclusion and stays untouched per §11.1.
- Three dated records in `DEV_PLAN{,_EN}` still list the EOL notification policy among the T018 remaining gates; the file is explicitly out of bounds and all three are dated progress notes.
- Point-in-time "not passed / not approved" statements in the frozen snapshots (`CHECKLIST_v59/v60{,_EN}`, `s1-reach-reachability{,_EN}`) are preserved under the §3.3 rule and are not rewritten.

## 9. Receipt elements

- The final-state readings, commit hash and CI run id are in the ⑩ receipt (this report records no commit hash or run id, per §11.2).
- Anchor for the next slice: P78 `RELEASE-ZIP-BINDING` takes the post-closeout `HEAD` and CI run id as its step-0 anchor.

## artifacts

```artifacts
docs/ADR_REGISTER.md	repo
docs/ADR_REGISTER_EN.md	repo
docs/S1_SUPPORT_MATRIX.md	repo
docs/S1_SUPPORT_MATRIX_EN.md	repo
docs/INSTALLATION_PLAN.md	repo
docs/INSTALLATION_PLAN_EN.md	repo
docs/PRODUCT_REQUIREMENTS.md	repo
docs/PRODUCT_REQUIREMENTS_EN.md	repo
docs/validation/s1-exit.md	repo
docs/validation/s1-exit_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/adr010-approval.md	repo
docs/validation/adr010-approval_EN.md	repo
docs/validation/evidence/P77-ADR010-APPROVAL/00-briefing.md	repo
docs/validation/evidence/P77-ADR010-APPROVAL/01-arch/design-r1.md	repo
docs/validation/evidence/P77-ADR010-APPROVAL/02-batch/batch-01-adr010-register-format.md	repo
docs/validation/evidence/P77-ADR010-APPROVAL/02-batch/batch-02-eol-survey.md	repo
docs/validation/evidence/P77-ADR010-APPROVAL/04-reviews/final-r1.md	repo
docs/validation/evidence/P77-ADR010-APPROVAL/04-reviews/prereview-r1.md	repo
docs/validation/evidence/P77-ADR010-APPROVAL/04-reviews/prereview-r2.md	repo
docs/validation/evidence/P77-ADR010-APPROVAL/06-gates/gate-ci.txt	repo
docs/validation/evidence/P77-ADR010-APPROVAL/06-gates/gate-index.txt	repo
docs/validation/evidence/P77-ADR010-APPROVAL/06-gates/gate-selftest.txt	repo
docs/validation/evidence/P77-ADR010-APPROVAL/06-gates/gate-tree.txt	repo
docs/validation/evidence/P77-ADR010-APPROVAL/MANIFEST.md	repo
docs/validation/evidence/P77-ADR010-APPROVAL/SHA256SUMS.txt	repo
docs/validation/evidence/P77-ADR010-APPROVAL/freedom-list.md	repo
docs/validation/evidence/P77-ADR010-APPROVAL/input-list.md	repo
```
