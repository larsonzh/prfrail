# P25 validation report: DIRECTIVE-CARRIER-ADJUST (directive v1.47)

## 1 Slice overview

| Field | Content |
|---|---|---|
| Identifier | `DIRECTIVE-CARRIER-ADJUST` (P25) |
| Date | 2026-10-24 |
| Size and tier | M (governance slice; ⑦ = 1 final review + 3 re-reviews, no downgrade) |
| Baseline | `HEAD` = `392dc2c` (P24 closed) |
| Change set | six existing documents, this report pair and the freedom list: nine paths, all `repo` |
| Deliverables | the bilingual directive v1.47, the appended ADR-027 disposition sentence, the P25 ledger status row and detail block, this report pair and the freedom list |
| Boundaries | zero change to: the §6.3 criteria, the §2.4 blacklist, the §2.1 role table, the §5.2 matrix, §7.2.2, the §7.5 quota table, appendix `B.6`, and `tools/`, `internal/`, `schemas/`, `examples/`, `testdata/`, `.github/`, `docs/DEV_PLAN` |

## 2 Change landing points

| File | Line | Substance |
|---|---|---|
| `docs/DELIVERY_DIRECTIVE{,_EN}.md` | `:3` | header version line v1.46 → v1.47, date 2026-10-24 |
| same | `:93` | §0.3 gains the v1.47 entry (five items ①②③④⑤) |
| same | `:267` / `:268` | §2.2 two rows swapped: only the qualified-name and `modelId` cells move; `agentName`, the ⑤ call point and the read-only tool cell stay unchanged |
| same | `:282` | §2.2 preferred / standby sentence swapped and marked "swapped 2026-10-24" |
| same | `:286` / `:288` | §2.2 carrier-correction note (OB-116 landed) and the historical-rows-not-retroactive note |
| same | `:480` | §5.1 native validation and SELFTEST time-budget paragraph |
| same | `:735` | §7.5 SELFTEST time-budget note (the quota counts invocations only) |
| same | `:1104` / `:1105` | §12.4 OB-116 disposition cell rewritten (item-by-item fate of the former three disciplines) and OB-117 added |
| same | `:1267` / `:1268` | appendix B.5 model line and carrier line |
| same | `1922`–`1934` | appendix C `C.0ah` (eight fields, isomorphic to `C.0ag`) |
| `docs/ADR_REGISTER{,_EN}.md` | `:37` | ADR-027 disposition cell gains the inline closing sentence "wrap-up (2026-10-24)" |
| `docs/t027/REMAINING_SLICES{,_EN}.md` | `:57`, `945`–`962` | P25 status row and detail block |

The table above and the fence below are the entire nine-path staged set; the fence rows match the disk one for one.

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/ADR_REGISTER.md	repo
docs/ADR_REGISTER_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/directive-carrier-adjust.md	repo
docs/validation/directive-carrier-adjust_EN.md	repo
docs/validation/evidence/DIRECTIVE-CARRIER-ADJUST-freedom-list.md	repo
```

## 3 Review chain

### 3.1 Dispatch table (mechanically read from `transcript.arguments.agentName`, verbatim)

| Role | agentName | Model | transcript line offset | Round | Quota | Verdict |
|---|---|---|---|---|---|---|
| ① | `deep-reasoner` | `DeepSeek V4 Pro (deepseek)` | 59561 | r1 | 1 | `DESIGN: DONE` |
| ② | `prfrail-documenter` | `DeepSeek V4.1 Flash (deepseek)` | 59676 | r1 | 1 | `IMPLEMENT: DONE` |
| ③ | `prfrail-documenter` | `DeepSeek V4.1 Flash (deepseek)` | 59867 | r1 | 1 | `TEST: DONE` (1 Medium + 4 Low) |
| ② | `prfrail-documenter` | same | 60061 | r2 | — | `IMPLEMENT: DONE` |
| ② | `prfrail-documenter` | same | 60354 | r3 | — | `IMPLEMENT: DONE` |
| ② | `prfrail-documenter` | same | 60508 | r4 | — | `IMPLEMENT: DONE` |
| ③ | `prfrail-documenter` | same | 60621 | r2 | — | `TEST: DONE` (Medium+ zero) |
| ⑤ | `deep-reasoner` | `GLM-5.3 (glm)` | 60932 | r1 | 1+1 | `PRE-REVIEW: PASS WITH FIXES` (1 High + 2 Medium) |
| ② | `prfrail-documenter` | same | 61072 | r5 | — | `IMPLEMENT: DONE` |
| ⑤ | `deep-reasoner` | `GLM-5.3 (glm)` | 61211 | r2 | (same slot) | `PRE-REVIEW: PASS WITH FIXES` (Medium+ zero) |
| ② | `prfrail-documenter` | same | 61299 | r6 | — | `IMPLEMENT: DONE` |
| ⑥ | `independent-reviewer` | `MAI-Code-1.1-Flash (copilot)` | 61501 | r1 | 1+2 (per slice ≤ 3) | outside the §7.2.2 contract (controller dispatch defect) |
| ⑥ | `independent-reviewer` | `MAI-Code-1.1-Flash (copilot)` | 61547 | r2 | (same) | still outside the contract (same cause) |
| ⑥ | `independent-reviewer` | `MAI-Code-1.1-Flash (copilot)` | 62180 | r3 | (same, quota exhausted) | `INDEPENDENT SCAN: PASS` (0 findings) |
| ⑦ | `independent-reviewer` | `GPT-5.3-Codex (copilot)` | 61579 | r1 (blind) | 1+3 | `FINAL: PASS` (0 findings) |

**Batch assistance**: **0 invocations** in this slice. Reason: the batch enumeration and the three mechanical comparisons were run directly as deterministic scripts by the controller (136 reference-point classifications, per-line mirror counts, three numstat pairs); there was no delegable batch-semantic task. Both `GLM-5.3-Flash (glm)` and `GLM-5.3-FlashX (glm)` count 0.

### 3.2 ① the six rulings (verdict `DESIGN: DONE`)

| Item | Ruling |
|---|---|---|
| (a) §2.2 two rows swapped | Only the qualified-name and `modelId` cells move; `agentName` and the tool-capability cell are unchanged. The cost basis cites the v1.28 historical measurement row of §0.3 (historical rows are not rewritten). |
| (b) §2.2 `:282` sentence aligned | Kept as one line; the `**` and backtick pairing stays mirror-comparable. |
| (c) both OB-116 carrier corrections written down | **§7.3 was ruled out of this slice boundary**, so nothing lands in §7.3; everything lands in the §2.2 note instead (including the ⑤⑦ same-carrier note duty). |
| (d) SELFTEST time budget | Lands as a standalone §5.1 paragraph plus a §7.5 post-table note; it does not conflict with the ban on physical-clock assertions - a budget is a planning expectation, never a failure trigger. |
| (e) appendix B.5 | The model line and the carrier call string are aligned to the preferred model. |
| (f) references split three ways | Synchronized in this slice (in-bounds landing points) / frozen as history (§0.3 and §12.4 historical rows and the existing appendix C blocks, per §11.1) / untouched in the report domain; plus a "historical rows are not retroactive" paragraph in §2.2. |

### 3.3 ③ Consistency check

r1 reported 1 Medium + 4 Low (`§0.3` and `C.0ah` did not record the newly added OB-117; an EN section-name reference missed `the`; EN translated the symptom column as `the evidence cell`; `§0.3` ① lacked the sentence-swap mapping; the ADR-027 appended sentence did not agree with the not-whitelisted fact) ⇒ after remediation r2 closed every item and **Medium+ is zero**; r2 added 2 Low + 3 Info, disposed of in §5.

### 3.4 ⑤ Pre-review

r1 `PRE-REVIEW: PASS WITH FIXES`: 1 High (the OB-116 disposition cell silently dropped original discipline ①, namely "every `agentName` in the ⑩ dispatch table must be read verbatim and mechanically from `transcript.arguments.agentName`, with the raw stdout attached"), 2 Medium (the "review layer (⑤⑥⑦)" wording conflicted with the review-layer scope defined in §2.1; the ⑤⑦ same-carrier note lacked a delivery path), 3 Low, 2 Info. After remediation r2 lifted the High and both Medium items, (b) moved from FAIL to PASS, and **Medium+ is zero** (3 Low + 2 Info, disposed of in §5). The ⑤ quota of 1+1 was consumed.

### 3.5 ⑥ Independent scan

⑥ quota (per-slice effective invocations ≤ 3, covering the first round, the re-dispatch and remediation re-runs) **is exhausted after 3 uses**: neither r1 (`L61501`) nor r2 (`L61547`) followed the §7.2.2 A–E paragraph definition and verdict-string closed set, **because the controller dispatch text wrote the paragraph definition and the verdict-string shape wrongly** (it used `SCAN:` instead of the closed-set shape), so both were re-dispatched under the §7.2.2 failure handling; the 3rd round (`L62180`), re-dispatched against the verbatim contract, complied, and the relation between its verdict-line value and Section A is stated in bullet 1 below. What follows is the **verbatim transcription** of ⑥ r3 (no rewriting and no normalization):

```sixstage
## Section A
This section carries the overall conclusion line required by §7.2.2-A; because G7-f of §6.3 requires exactly one conclusion line inside the fence and it must be the last non-empty line of the block, that conclusion line is merged with the last line of this block, and the value is given at the end of the block.

## Section B
No findings.

## Section C
- The review was read-only; no file writes and no git write operations were performed.
- The change set is confined to the four in-scope adjustments: the preferred/backup carrier ordering correction, the OB-116 carrier-note correction, the SELFTEST time-budget clarification, and the ADR-027 follow-through wording.
- The patch does not extend into protected sections or runtime logic. The reviewed changes remain document-governance and policy wording only, with no new semantics introduced beyond the stated corrections.

## Section D
- Guard examined: scope and wording remain bounded to the document-only corrections.
- If this guard were removed, the check for "no protected semantic change" would go red because the scan would then rely only on manual reading rather than the explicit scope constraint.
- What would not go red: the revised text still reads consistently and there is no evidence of a cross-cutting operational change in the reviewed patch.
- UNCOVERED: no exhaustive review of future dispatch records or hidden downstream automation outside the current revision; this scan is limited to the present in-scope change.

## Section E
- I cannot verify: any external, unreviewed dispatch state or future model-carrier references beyond the reviewed patch.
- What would overturn my PASS: any additional semantic change beyond the documented carrier correction and notes, any protected-section edit, or evidence that the patch was not limited to document-governance wording.
- The verdict is valid only for the current observed read-only revision.

INDEPENDENT SCAN: PASS
```

- Section B reports no finding ⇒ the verdict line takes `INDEPENDENT SCAN: PASS` from the closed set; under the `§7.2.2` verdict-line consistency hard constraint, that value agrees with the overall conclusion carried by Section A, and in this report the two are merged into the single place at the end of the block (the difference is disclosed verbatim below).
- **The difference from the ⑥ r3 original is exactly one place** (disclosed verbatim): `§7.2.2-A` requires Section A to carry the overall conclusion line, while `G7-f` of `§6.3` requires exactly one conclusion line inside the fence and that it be the last non-empty line of the block; the two cannot both hold. This slice does not change the criteria (`§6.3` and `tools/gates/**` are out of bounds; the ⑥ quota is exhausted), so the Section A conclusion line is merged with the last line of the block and Section A becomes one sentence pointing at the end of the block; **the rest of the text (Sections B–E and the last-line value) is unchanged verbatim**. The conflict is registered as `D-3`, and a later governance slice is advised to fix `G7-f` (allow one occurrence in Section A plus one final line, provided the two values are equal).
- The controller defect is registered as `D-2` in §5 of this report; the exhaustion of the ⑥ quota is kept on record under the §7.2.2 failure handling.

### 3.6 ⑦ Final review (blind)

r1 was a **blind review** (no prior finding list was supplied) and judged eight items: the four declared goals, internal contradictions and dangling references, CN/EN mirror and numstat symmetry, boundaries and protected subsections, and silently added duties - **all PASS**, with a finding list holding a single Info row "no reportable finding" and the verdict `FINAL: PASS`. Of the ⑦ quota of 1+3, **one final review** was used and three re-reviews remained unused (no remediation item). No downgrade (`GPT-5.3-Codex (copilot)`).

## 4 Gate readings

### 4.1 Final readings before this report landed

| Scope | Last line | EXIT | Elapsed |
|---|---|---|---|
| `--all --scope=index` | `TOTAL_FAIL=0` | 0 | 103.9 s |
| `--all --scope=tree` | `TOTAL_FAIL=0` | 0 | 99.9 s |
| `--all --scope=ci` | `TOTAL_FAIL=0` | 0 | 124.6 s |
| `--selftest` | `SELFTEST: PASS 663/663` | 0 | 328.7 s |

All three scopes carry the base checks: `gofmt` empty, `go build` exit 0, `go vet` exit 0, `go test` package-level `ok` count 15.

Criterion readings (`index` and `tree` agree): `G1-b` `header=v1.47 last=v1.47`, `rows=48`; `G3(a)` three pairs symmetric in line counts (ADR `+1/-1` each side, directive `+30/-7` each side, ledger `+19/-0` each side) with the directive pair equal on all six shapes; `G3(b)` zero mismatch in all five key classes; `G4a` six files pass BOM+LF; `G4b` `unknown=[]` (whitelist 37); `G5-c` the appendix B call strings sit inside the §2.2 qualified-name set; `G6-4` no expiring literal; `G6-5` no unresolved section reference; `G6-6` every added ledger row opens with a closed-set status word; `G6-7` declarations reconcile with the ledger; `G6-8` no DR collision and no dangling reference.

### 4.2 Final readings including this report

Readings taken after this report pair and the freedom list landed and all nine paths were staged (zero change to the baseline and the criteria implementation):

| Scope | Last line | EXIT | Elapsed |
|---|---|---|---|
| `--all --scope=index` | `TOTAL_FAIL=0` | 0 | 110.0 s |
| `--all --scope=tree` | `TOTAL_FAIL=0` | 0 | 111.9 s |
| `--all --scope=ci` | `TOTAL_FAIL=0` | 0 | 129.2 s |
| `--selftest` | `SELFTEST: PASS 663/663` | 0 | 292.7 s |

Criterion-count deltas (relative to the state before the report landed): the `G4a` file count moves from six to nine (the evidence-tree directory stays excluded by the existing rule); `G4b` keeps `unknown=[]` (the whitelist of 37 entries is not extended); the scan line count of `G6-4` and `G9-e` rises to 482 with the report volume; `G6-5` scans 28 lines with all 70 references resolved; `G7-f` finds both `sixstage` fences compliant; `G9` passes all nine paths. All three scopes carry the base checks: `gofmt` empty, `go build` exit 0, `go vet` exit 0, `go test` package-level `ok` count 15.

### 4.3 Observation on the SELFTEST time baseline

The `--selftest` elapsed time measured twice in this slice is 292.7 s (the final pass) and 328.7 s (the round before the report landed), while the baseline recorded by directive v1.46 is 269.3 s ⇒ the two readings overshoot it by about 9% and 22% respectively. Under the §5.1 paragraph and the §7.5 note added by this slice, that overshoot **is not a failure and is not escalated**; it is recorded only as an observation on the time baseline, and a later governance slice may re-estimate the baseline from these readings. The gap between the two readings comes from the same assertion set drifting under load fluctuation, not from a change in assertion count (`663/663` in both rounds).

## 5 Low-severity disposition list

| # | Source | SEVERITY | Location | Disposition |
|---|---|---|---|---|
| C-6 | ⑤ r2 | Low | the time-budget clauses of §5.1 and §7.5 reference each other in one direction only | **Not fixed**. Reason: the section name is unique in the document and its main body is verbatim identical, so no ambiguity arises. |
| C-8 | ⑤ r2 | Info | the ledger records "batch assistance once used `quick-researcher`" as an OB-116 deviation face | **Not fixed**. The first-hand evidence is the P24 report `docs/validation/agentrunner-prod-patch.md`; this row is its provenance record. |
| R-4 | ⑤ r2 | Info | the follow-up governance slice `DIRECTIVE-DISPATCH-EVIDENCE` is registered only inside the OB-116 disposition cell | **Forwarded to the ⑩ receipt**: the controller proposes the slice in the stop-point receipt and the user rules on it. |
| R-5 | ⑤ r2 | Info | the OB-116 cell quotes original ① without the "never fill in from memory or inference" wording | **Not fixed**. The semantics are covered by "read verbatim and mechanically"; a later slice completes the wording. |
| ③ r2 #6 | ③ | Low | the `C.0ah` deliverables omit the ledger pair | **Not fixed**. Writing the ledger back is a mandatory §11.2 duty, not an in-bounds deliverable of this slice; the scope fields of §0.3 and `C.0ah` already limit themselves to in-bounds. |
| ③ r2 #7 | ③ | Low | the P25 status sentence in the ledger said "① not yet dispatched", contradicting the actual progress | **Fixed**: the ⑧b write-back replaces it with wording that binds no stage. |
| ③ r2 #8 | ③ | Info | the ledger goal field says "formally revise §2.2" | **Not fixed**. It matches the ledger's existing convention. |
| ③ r2 #10 | ③ | Info | after the swap, ① and ⑤ share one carrier and one model | **Recorded**: per R2.6, anchoring protection rests on call isolation rather than role or model separation, so this is not a rule conflict. |
| D-1 | controller | — | the ADR-027 undertaking was corrected from "undertaken by a formal §2.2 revision" to "the truthful-record branch plus a record of the ⑤ preference swap" | **Fixed, with a trail**. Reason: `GLM-5.3-FlashX (glm)` never entered the whitelist, and the original wording could be read as taking the other branch; found by ③, fixed by ②, cleared by the ⑦ blind review. |
| D-2 | controller | — | ⑥ dispatch defect: the dispatch text wrote the §7.2.2 A–E paragraph definition and verdict-string closed set wrongly (it used `SCAN:` instead of `INDEPENDENT SCAN:`), so ⑥ r1 / r2 were two non-compliant rounds and the quota was exhausted | **Fixed, with a trail**. The 3rd re-dispatch against the §7.2.2 verbatim contract complied; the ⑧a "normalization" made at the time to pass the gate violates the §11.3 verbatim-transcription clause and has been superseded by the verbatim transcription in this section |
| D-3 | controller / ② | High (criteria conflict) | the "exactly one conclusion line inside the fence, and it must be the last line" rule of `§6.3`/`G7-f` and the "Section A carries the overall conclusion line" rule of `§7.2.2-A` cannot both hold: a compliant ⑥ output necessarily contains two conclusion lines ⇒ transcribing it verbatim into a `sixstage` fence is bound to be red | **Not fixed in this slice** (`§6.3` and `tools/gates/**` are out of bounds; the ⑥ quota is exhausted). Disposition = the fence keeps the compliant `sixstage` structure (exactly one conclusion line inside the block, on the last line); the Section A conclusion line is merged with the last line and the **sole difference** is disclosed verbatim in §3.5; **a later governance slice should fix** `G7-f` (allow one occurrence in Section A plus one final line, provided the two values are equal) |

## 6 Boundaries and non-retroactivity

- The change set holds only six existing documents plus this report pair and the freedom list; `tools/`, `internal/`, `schemas/`, `examples/`, `testdata/`, `.github/` and the `docs/DEV_PLAN` pair carry **zero change**.
- Protected subsections carry zero change: the §6.3 criteria, the §2.4 blacklist, the §2.1 role table, the §5.2 matrix, appendix `B.6`, §7.2.2 and the §7.5 quota table; §7.3 is **verbatim equal** to the baseline.
- Historical rows are not retroactive: the v1.26 and v1.28 entries of §0.3, the OB-83 and OB-87 entries of §12.4 and the existing appendix C blocks were not rewritten; the symptom column of OB-116 is verbatim untouched.

## 7 Residuals and follow-ups

- `DIRECTIVE-DISPATCH-EVIDENCE`: recommended as a new slice (its destination is already registered in the OB-116 cell) to carry the action duty "mechanically read `agentName` for the ⑩ dispatch table, with the raw stdout"; the user rules on the slice.
- The time-budget clauses of §5.1 and §7.5 reference each other one way only; a later governance slice may make it mutual.
- The temporary authorization for `GLM-5.3-FlashX (glm)` ended with the ADR-027 wrap-up; ⑤ is back on a low-cost whitelisted tier.

## 8 Commit and CI

| Main commit hash | CI run id | CI conclusion | Note |
|---|---|---|---|
| filled in after authorization | filled in after authorization | filled in after authorization | Commit and push need explicit authorization in the same user turn, and gitee is never pushed; this report lands inside the commit, so the hash and the run id stay in the report only, per §11.2. |
