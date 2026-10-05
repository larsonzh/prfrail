# Slice validation report · `DIRECTIVE-GLM-LANDING` (directive v1.26)

> Slice id: `DIRECTIVE-GLM-LANDING`; directive v1.26; date 2026-10-06. Tier `[SLICE]`; size M; probes 0; `[ROUTINE]` not applicable.
> **Status**: **⑧b finalized** (the local part of the ⑨ native verification has run, with readings in §7 and §7.7; the ⑧c mechanical verification and close-out cleanup are recorded in §11, with its reading sequence deviation truthfully registered). **Authorization boundary**: commit no / push no / gitee no / probes 0; commit and push **require same-round authorization**.

## 1. Goal and scope

- Slice nature: a directive-governance / criterion-body change slice, tier `[SLICE]`, size M, probes 0, `[ROUTINE]` not applicable.
- Trigger: the user's 2026-10-06 kickoff declaration; this is a directive-governance slice, so the **pilot is deferred for the 5th consecutive time** — registered **independently as OB-85** under §1.7.4 / the OB-37 meta-rule (four elements + truthful annotation + "must not be taken as precedent").
- Eight landing goals: ① §2.2 gains `GLM-5.3 (glm)` (⑤ pre-review **first choice**) and `GLM-5.3-Flash (glm)` (**batch-auxiliary first choice**) as two rows; ② §2.2's existing pre-reviewer row is re-noted **standby**; ③ §5.2 gains the ⑤ note paragraph; ④ §6.3's G5-b whitelist closed set gains `glm-5.3` / `glm-5.3-flash`; ⑤ §12.4's OB-83 disposition column flips to `Resolved`; ⑥ batch-auxiliary is **limited to four classes** and **excludes sorting and Top-N selection**; ⑦ positions ① / ②③ unchanged; ⑧ cost comparison is not in this slice.
- Out of scope (boundaries stated): §2.4 blacklist body, §5.2 matrix body, §7.5 numbers, §1.5.1, §1.7, §7.7, the role set and carriers, `.github/agents/**`, `internal/**`, `docs/CONTRACTS{,_EN}.md`, `schemas/**`, contract fixtures, `.github/workflows/**`, `go.mod`.

## 2. Change summary

- Change set = the **13** repo artifacts of `git diff --cached` (numstat) (9 implementation / text + the report pair + the ledger pair); every artifact is **staged file by file** (**13** items in the index); this slice is **not committed and not pushed** (commit and push **require same-round authorization**).

| Artifact | numstat |
|---|---|
| `docs/DELIVERY_DIRECTIVE.md` | +37/−9 |
| `docs/DELIVERY_DIRECTIVE_EN.md` | +37/−9 |
| `tools/gates/lib/criteria/g5b.js` | +2/−0 |
| `tools/gates/selftest.js` | +69/−1 |
| `tools/gates/testdata/fixtures/g5b-closed-set-legal-glm.txt` | +6 (new) |
| `tools/gates/testdata/fixtures/g5b-closed-set-legal-glm-flash.txt` | +6 (new) |
| `tools/gates/testdata/fixtures/g5b-illegal-glm-case.txt` | +6 (new) |
| `tools/gates/testdata/fixtures/g5b-illegal-glm-qualified.txt` | +6 (new) |
| `docs/validation/evidence/DIRECTIVE-GLM-LANDING-freedom-list.md` | +23/−0 (new) |
| `docs/validation/directive-glm-landing.md` | +180/−0 (new) |
| `docs/validation/directive-glm-landing_EN.md` | +180/−0 (new) |
| `docs/t027/REMAINING_SLICES.md` | +1/−0 |
| `docs/t027/REMAINING_SLICES_EN.md` | +1/−0 |

- The nine touches of directive v1.26 (enumerated truthfully by the `§0.3` v1.26 log row): §2.2 / the §5.2 note paragraph / §6.3's G5-b / §7.1 / §7.2.3 / §12.4 / appendix B.5 / appendix C.0o / the §5.1 pipeline diagram.
- Landing point (models entering the table): `GLM-5.3 (glm)` enters §2.2 as the **⑤ pre-review first choice**, and the original `DeepSeek V4 Pro (deepseek)` ⑤ position is re-noted **standby**; `GLM-5.3-Flash (glm)` is the **batch-auxiliary first choice**, **limited to four classes** — enumeration / statistics, long-document excerpting with numeric cross-checking, bilingual-mirror language-side checks, MANIFEST draft lists — and **explicitly excludes** sorting and Top-N selection (basis = the BATCH-a counterexample in `docs/validation/directive-glm-eval{,_EN}.md` §4.2).
- Landing point (criterion and ledger): the G5-b closed set goes 6 → **8** members and is **element-wise equal** to §2.2's deduplicated `modelId` set (the ⑦ first reading mechanically verified `SET_EQUAL=YES`); §12.4's OB-83 disposition column flips from `Registered` to `Resolved` (parenthetical "landed (standing): temporary authorization ended, escalation trigger closed"); appendix C.0o is added (at ⑧a its status read `⬜ in progress` with no premature completion claim; ⑧b finalization flips it to `✅ completed` as planned).

## 3. Pipeline and environment

- Pipeline: ① plan → ②a rule text → ②b implementation → ③ tests → ④ gates → ⑤ first pre-review → ⑥ independent scan → ⑦ first reading (final review) → ⑦ second reading (narrowed-input re-review) → ⑧a this report → ⑨ native verification (local re-run + the §7.7 A12 dedicated independent experiment) → ⑧b finalization.
- Roles and models: ① plan = `DeepSeek V4 Pro (deepseek)`; ②a / ②b / ③ / ⑧ are carried by `DeepSeek V4.1 Flash (deepseek)`; ⑤ first pre-review = `DeepSeek V4 Pro (deepseek)`; ⑥ = `MAI-Code-1.1-Flash (copilot)`; ⑦ = `GPT-5.3-Codex (copilot)`, **not downgraded**.
- Environment: Windows host; repository `d:\LZProjects\prfrail` (the `prfrail` folder of a multi-root workspace); gate entry `node tools/gates/gate.js`.

## 4. Anomalies and fallbacks

- **A-1 §1.7.3 sequence deviation (driver orchestration error)**: the freedom list did **not** enter the ⑦ first reading's input pack, while §1.7.3 requires it to land after ⑥ / before ⑦ and to be part of the ⑦ input pack. Remedy = the driver landed it late plus the ⑦ second reading (narrowed input) performing the out-of-bounds judgement. **Must not be taken as precedent.** This item is consistent with the last paragraph of `docs/validation/evidence/DIRECTIVE-GLM-LANDING-freedom-list.md`.
- **A-2 the ⑥ coverage gap**: ⑥'s Section C / Section E did **not cover** four of the driver's must-verify items — "the mechanical enumeration of §2.2 against the closed set", "independent mutation testing of the assertions", "the enumeration consistency between §0.3 and appendix C.0o", "the historical-vs-present description boundary" — and did **not** list them under "cannot verify"; that gap is carried by ⑦ under §7.8.2. The driver did **not** add a ⑥ call (cap ≤ 3, this slice used 1).
- **A-3 remediation preceded review**: the ⑤ first reading's Low item 4 (appendix C.0o's deliverables row omitting §5.1) and the driver's own same-family omission (appendix C.0o's goal row not naming "four-class limit / sorting and Top-N exclusion") were both remediated, and the remediation happened **before** ⑥ and ⑦ ⇒ ⑥ / ⑦ reviewed the **remediated** state.
- **A-4 two empirical corrections to ①**: (a) the §5.2 claim quoted in the user's kickoff message ("⑦'s three readings are not downgraded in directive-governance slices") has **no such wording** in §5.2 (what exists is the matrix's "⑦ final review = required" + the note paragraph "the ⑦ independent final review is indispensable" + §5.2's ⑦ three-model tier table); (b) §12.4's OB-83 disposition column originally read **`Registered`** (not "temporary authorization").
- **A-5 two empirical "goal-missed" findings**: (a) §2.4's parenthetical is **pointer-style** ⇒ a new whitelist member takes effect automatically and §2.4 has **zero diff** (so the v1.26 log row **must not** and **does not** claim §2.4 was changed); (b) §5.2 originally had **no** ⑤ note paragraph ⇒ goal ③ is in fact a **new paragraph** (pointer only, carrying no `glm` string).
- **A-6 line loss in the ④ phase pipeline**: the driver found that a PowerShell pipe **loses lines / garbles bytes** (the same command produced inconsistent hunk counts across two runs) ⇒ switched to `cmd /c` byte-level redirection plus a node re-check, yielding `hunks=13 / adds=37 / dels=9` (identical for CN and `_EN`).
- **A-7 ⑧b / ⑧c sequence deviation (driver orchestration, truthfully registered)**: the ⑧c mechanical-verification readings were taken before ⑧b was finalized (to write the ⑧c readings into the report and avoid editing the report after finalization), which contradicts the literal order "⑧c after ⑧b" in §5.1; after ⑧b was finalized the driver additionally re-ran `--all --scope=tree` (readings in §11). **Must not be taken as precedent.**

## 5. Review conclusion summary (⑤⑥⑦)

- **⑤ first reading**: `PRE-REVIEW: PASS WITH FIXES` (**no Medium+**), 4 Low + 2 Info: ① Low — the v1.26 log row carries one ⑤ item more than the ① plan's C5 draft (§5.1 / §7.1 / §7.2.3 label sync in three places), **a direction more faithful to the user's enumeration constraint**, kept; ② Low — the qualified-name fixture value was changed from the draft's `GLM-5.3 (glm)` to `glm-5.3 (glm)`, **decoupling** the case boundary from the R2.2 boundary, kept; ③ Low — appendix C.0o's status row was changed from the draft's `✅ completed` to `⬜ in progress`, a truthfulness direction, kept; ④ Low — **appendix C.0o's deliverables row omits §5.1** ⇒ **remediated** (1 place each in CN / `_EN`; after remediation G3-a / G4a / G6-5 / G6-6 re-ran all green); ⑤ Info — T63 / T65 are weakly falsifiable (mutation-restoration completeness), adequately disclosed, passed to ⑦ to confirm; ⑥ Info — the remaining "V4 Pro pre-review" wording elsewhere in the repository is **historical execution record** (`docs/DEV_PLAN{,_EN}.md`'s 2026-09-14 / 15 entries, the pipeline line of `docs/validation/directive-evidence-scope{,_EN}.md`, `docs/t027/FLASH_OPERATING_DIRECTIVE*`), not retroactively rewritten under §11.1.
- **⑥ independent scan**: `INDEPENDENT SCAN: PASS` (0 file:line-level findings; Sections A–E present); one coverage gap truthfully registered (see A-2); the driver did **not** add a ⑥ call.
- **⑦ first reading (final review)**: `FINAL REVIEW: PASS` (**0 Medium+**); 1 Info = §2.2's first-choice / standby annotation and the batch-auxiliary four-class limit have **no mechanical criterion coverage** (consistent with the user's R-2 ruling); all 8 check surfaces "verified", including boundary zero diff, the closed set's `SET_EQUAL=YES`, the criterion adding set elements only, independent in-memory mutation (`DROP_GLM_RESULT=FAIL` / `DROP_GLM_FLASH_RESULT=FAIL`), zero out-of-bounds hits, and the gate re-run in three scopes.
- **⑦ second reading (narrowed-input re-review)**: `RE-REVIEW: PASS` (**0 Medium+**); **F-1…F-11 each judged in-bounds**; the sequence-deviation registration judged **sufficient** (no new observation item needed; ⑧b writes it back truthfully); the 8 files of the first reading show no drift (numstat identical + current hash / index fingerprint + "write-file tool hits after the first reading = 0"); its Info = the first reading's input pack **did not attach the then-current hash baseline**, so numstat identity + fingerprint + 0 hits served as the compensating evidence chain.

## 6. Review input pack summary (with blind-review isolation proof)

- Input packs: ⑦ first reading = the 8 changed artifacts (**without** the freedom list, see A-1); ⑦ second reading = the **narrowed input** (the freedom list added, for the out-of-bounds judgement).
- **Blind-review isolation proof**: this slice **adds no blind review** (it does not touch §7.7's mandatory-sampling three semantics); ⑦ = 1 final review + 1 re-review = **2 readings**, within the cap (1 + 1), with no resend and no voided round.
- Sanitization: the review input contains only repository paths and text, no credentials and no external connection material.

## 7. Falsifiability and gate results

| Item | Reading |
|---|---|
| `node tools/gates/gate.js --all --scope=tree` | `changed=13`, `TOTAL_FAIL=0` (final-state reading; the ⑨ time point read `changed=11`) |
| `node tools/gates/gate.js --all --scope=index` | `changed=13`, `TOTAL_FAIL=0` (final-state reading; the ⑨ time point read `changed=11`) |
| `node tools/gates/gate.js --selftest` | `SELFTEST: PASS 335/335` (this slice's baseline 328/328, plus 7 assertions T61–T67) |
| `G1-b` / `G4a` / `G7-a` / `G7-b` / `G8-a` | `header=v1.26 last=v1.26 2026-10-06/2026-10-06 rows=27` (OK for both CN and `_EN`); `11 file(s) OK (1 excluded: evidence tree)`; `2 report(s) OK`; `40 artifact row(s) OK`; `1 path(s) OK` |
| `G8-b` / `G8-c` | `INFO` (this slice does not touch the frozen evidence packs) |
| `G6-4` / `G3(a) symmetry + mirror parity` | scan reading `392`; CN `+37/−9` vs `_EN` `+37/−9`, the five shapes (lines / headings / bold / pipes / blanks) in sync |
| Independent frozen-region re-check | §2.4 blacklist body, the §5.2 **matrix body**, the §7.5 numbers and appendix **B.1** are **byte-untouched** in both CN and `_EN` (the driver compared the first 12 hex digits of sha256 with `tmp/glm-landing/check-boundaries.js`, six items `identical=true`; the ⑦ first reading independently re-checked `IDENTICAL` / `MATRIX_TABLE_IDENTICAL=YES`) |

- The readings above are **local mechanical readings**, already measured by ④ and ⑨ (see this table, §5 and §7.7); the `40` in `G7-b` is the ⑨-time reading; after the report is finalized the count is measured by G7-b against the finalized §16 rows = **48** (the finalization-time reading).
- **⑨ executed (local part)**: the local re-run readings are in this table and §7.7; the **CI leg** is observed after commit and push (same-round authorization), and this report **does not claim** CI has run.

### §7.7 A12 dedicated independent experiment

- Because ⑥ reported `PASS`, ⑨ must independently prove the key invariant with a criterion unrelated to ⑥: `tmp/glm-landing/verify-a12.js` ⇒ **all 26 items OK, `A12_EXPERIMENT_FAILURES=0`**, in three parts.
- (a) **live criterion measurement** (a synthetic in-memory ctx, bypassing git and writing nothing to the repository): `glm-5.3` PASS, `glm-5.3-flash` PASS, non-whitelist `gpt-5.4-terra` FAIL, `GLM-5.3` (case) FAIL, `glm-5.3 (glm)` (qualified-name form) FAIL, `'glm-5.3'` (quoted) FAIL, missing `model` key PASS (§2.6), empty value `model: ""` PASS (leaving room for G5-a ②).
- (b) **document ↔ code invariant, equal in both directions**: the §2.2 table has `11` rows in CN and in `_EN`; the §2.2 `modelId` deduplicated set == `g5b.js`'s `CLOSED_SET` (8 members) with **no extra member in the reverse direction**; `GLM-5.3 (glm)` → `glm-5.3` and `GLM-5.3-Flash (glm)` → `glm-5.3-flash`; ⑤ first choice = GLM and ⑤ backup = `deepseek-v4-pro`; the architect (①) and the ②③⑧ positions are unchanged.
- (c) **boundary invariants**: §2.4's body / the whole §5.2 section / the §7.5 numbers carry **no `glm` token** in CN or `_EN`; appendix B.1 is still `DeepSeek V4 Pro (deepseek)` and appendix B.5 is already `GLM-5.3 (glm)` (same shape on both sides).

## 8. Known boundaries and residual risks

- **C-1 observation-item candidate (not landed in this slice)**: §2.2's first-choice / standby annotation and the batch-auxiliary four-class limit have **no mechanical criterion coverage** (independently flagged by both the ⑦ first reading and ⑥); the user ruled it out for this slice (§12.4's goal enumeration contains OB-83 only) ⇒ kept for the next slice's kickoff to register in the same batch.
- **C-2 known boundary + residual risk**: §7.5's "cost-gradient hint (measured pricing)" lists only Flash / V4 Pro and is **not the same source** as §2.2 after the GLM switch ⇒ **§7.5 untouched** (boundary stated), registered as a known boundary.
- **A-2 coverage gap**: the four must-verify items not covered by ⑥'s Section C / E are carried by ⑦ under §7.8.2; this must not be read as ⑥ having covered them.
- **Historical wording not retroactive**: the remaining "V4 Pro pre-review" wording is historical execution record, not retroactively rewritten under §11.1 (see §5's ⑤ Info item 6).
- **Boundaries stated**: this slice does not touch §2.4's blacklist body, §5.2's matrix body or §7.5's numbers; §2.4 is pointer-style, so a new member takes effect automatically.

## 9. Explicitly not executed

- **⑨ executed (local part)**: the local re-run and the §7.7 A12 dedicated independent experiment are both complete, with readings in §7; the **CI leg** has not happened (it can only be observed after commit and push), and this report **does not claim** CI has run.
- **commit / push not executed** (`git add` staging only); **gitee not executed**; commit and push **require same-round authorization**.
- **C-1 not landed in this slice**; **C-2 leaves §7.5 untouched**; **cost comparison is not in this slice**.
- Probes 0; no downgrade ADR.

## 10. Deviations and degrees of freedom

- Same source as the freedom list, item by item in `docs/validation/evidence/DIRECTIVE-GLM-LANDING-freedom-list.md` (F-1 to F-11); every judgement is **in-bounds / an implementation choice**.
- F-1: "⑤ first choice / standby" is expressed with §2.2's **row labels**, **no new table column** (in-bounds).
- F-2: the `V4 Pro` labels in the §5.1 pipeline diagram / §7.1 three-layer review table / §7.2.3 three-level fallback are synced (OB-69 consistency duty, **not scope expansion**) (in-bounds).
- F-3: §0.3's v1.26 log row is rewritten as a whole line (the step-0 trailing clause became void on landing) (in-bounds).
- F-4: OB-83's disposition column **opens with `Resolved`**, with "landed (standing)" only as a parenthetical (G6-6 closed set + the user's R-1) (in-bounds).
- F-5: the qualified-name fixture value is `glm-5.3 (glm)` (decoupling the case boundary from the R2.2 boundary) (in-bounds).
- F-6: appendix C.0o's status row reads `⬜ in progress`, with **no** premature completion (in-bounds).
- F-7: appendix C.0o's goal row gains "four-class limit, sorting and Top-N exclusion" and the deliverables row gains §5.1 (in-bounds).
- F-8: the ⑥ / ⑦ call counts and other **cost vectors enter the validation report only** (the G6-2 caliber, F-9 precedent) (implementation choice).
- F-9: 4 new fixtures and the seven self-test assertions T61–T67 (in-bounds).
- F-10: this slice **does not land** the C-1 observation-item candidate (the user's R-2) (in-bounds).
- F-11: **§7.5 untouched** (the user's R-3) (in-bounds).
- **Sequence deviation (A-1, truthfully registered)**: the freedom list did **not** enter the ⑦ first reading's input pack, a **driver orchestration error**; remedy = late landing + the ⑦ second reading performing the out-of-bounds judgement. **This deviation must not be taken as precedent.**

## 11. ⑧c close-out cleanup record

- The ⑧c mechanical verification is run by the driver; **the sequence is truthfully registered** — its readings were taken **before** this report was finalized (⑧b) (purpose = to avoid editing the report after finalization; after ⑧b was finalized the driver additionally re-ran `--all --scope=tree`, see the last bullet of this section), a **truthful deviation registration** against §5.1's "⑧c after ⑧b" sequence; the process artifacts under `tmp/glm-landing/` are one-off process products of §1.5.1 **class ③**, removed immediately after use and not committed.
- Encoding check (`tmp/glm-landing/wrap-enc.js`, all of `git ls-files`): the **non-frozen tree is `tracked=3414 text=3101 binary=313 violations=0`**; the frozen evidence tree separately reports `violations=9`, all under already-registered / already-confirmed calibers.
- Frozen-tree caliber (stated truthfully): the CRLF in `b2-2026-09-17/measuring/probe10-dir-listing.raw.txt` (registered as OB-27), the BOM in `b3b-2026-09-21/{session,variant}/pin.json` (a pre-existing legacy), and six items under `docs/validation/evidence/glm-eval/archive/**` (including `hash-bom.txt`, which is the **BOM compatibility fixture** exception); those six are in the **expected shape** under v1.24's "archive payloads are normalized by their archived type" caliber (`.ps1` keeps the BOM, the rest do not). This slice's script judges the frozen tree with the generic `.md` rule, hence the 9 reported items; the verdict rests on the **non-frozen tree's 0 violations**.
- Residue: untracked residue is **0** (`tmp/` is excluded by `.gitignore` and its content is class ③ process material under §1.5.1); the 10 process artifacts under `tmp/glm-landing/` are listed in §16.
- IDE diagnostics: the 7 changed files (directive CN / `_EN`, report CN / `_EN`, `g5b.js`, `selftest.js`, the freedom list) show **0 errors**.
- `node tools/gates/gate.js --all --scope=tree` re-run at the ⑧c time point ⇒ `TOTAL_FAIL=0`.
- The frozen regions (§2.4 / §5.2 matrix / §7.5 / appendix B.1) take no part in the rewrite; this slice has zero diff against them.

## 12. Cost and accounting

- R7.1 line by line (field caliber = `{role, model, stage, startedAt, reason, retryOf}`): ① `DeepSeek V4 Pro (deepseek)` 1 time (plan); ⑤ `DeepSeek V4 Pro (deepseek)` 1 time (first pre-review, **no re-review**, reason = the first round had **no Medium+**, the user's caliber being "re-review only if the first round reports Medium+"); ⑥ `MAI-Code-1.1-Flash (copilot)` 1 time (effective call, cap ≤ 3); ⑦ `GPT-5.3-Codex (copilot)` **2 times** (final review + narrowed-input re-review; cap = 1 + 1; **no downgrade, no added blind review**; no resend, no voided round).
- Bearers (all `DeepSeek V4.1 Flash (deepseek)`): ②a 1 time / ②b 1 time / ③ 1 time / ⑧a 1 time / ⑧b 1 time / ⑧b′ 1 time / ⑧b″ 1 time ⇒ **12** calls in total (① 1 + ⑤ 1 + ⑥ 1 + ⑦ 2 + Flash 7).
- **The cost vector enters this section only, not the directive (§0.3 / appendix C.0o)** (G6-2).

## 13. Candidate observation-item assessment (not landed in this slice)

- **C-1**: §2.2's first-choice / standby annotation and the batch-auxiliary four-class limit have **no mechanical criterion coverage** — independently flagged by both the ⑦ first reading and ⑥; the user ruled it out for this slice, kept for the next slice's kickoff to register in the same batch.
- **C-2**: §7.5's "cost-gradient hint (measured pricing)" is **not the same source** as §2.2 after the GLM switch — §7.5 untouched, registered as a known boundary + residual risk, handled in the same batch as C-1.

## 14. Follow-up recommendations (including tier advice)

- Follow-up recommendations (recommendations only, not registered): C-1 and C-2 are kept for the next slice's kickoff to register in the same batch; the pilot-deferral count uses §12.4 as its single source of truth.
- Tier advice: for a follow-up criterion-body slice, ⑦ still uses `GPT-5.3-Codex (copilot)`, **not downgraded**.

## 15. ⑩ stopping-point report (including the write-back check line)

- Executed through **⑧b finalization**; the local part of ⑨ has run, and the ⑧c mechanical-verification readings were taken before ⑧b was finalized (**the sequence deviation is truthfully registered**, see §11); **the ⑩ stopping-point receipt has been issued**; this slice **stops before commit authorization** (no commit, no push, no gitee push).
- **Write-back check line** (§11.2 caliber): validation report = **written back** (this file pair); `DEV_PLAN` = **not applicable** (a directive-governance slice writes back appendix C, not `DEV_PLAN`); ledger = **written back** (the `docs/t027/REMAINING_SLICES{,_EN}.md` completion row, report path only); ADR = **not applicable** (no model downgrade / no budget exception).
- Authorization boundary: **commit no / push no / gitee no / probes 0**; commit and push **require same-round authorization**.

## 16. Artifacts list

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/validation/directive-glm-landing.md	repo
docs/validation/directive-glm-landing_EN.md	repo
docs/validation/evidence/DIRECTIVE-GLM-LANDING-freedom-list.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
tools/gates/lib/criteria/g5b.js	repo
tools/gates/selftest.js	repo
tools/gates/testdata/fixtures/g5b-closed-set-legal-glm-flash.txt	repo
tools/gates/testdata/fixtures/g5b-closed-set-legal-glm.txt	repo
tools/gates/testdata/fixtures/g5b-illegal-glm-case.txt	repo
tools/gates/testdata/fixtures/g5b-illegal-glm-qualified.txt	repo
tmp/glm-landing/PLAN-v1.md	local-only
tmp/glm-landing/cn.diff	local-only
tmp/glm-landing/en.diff	local-only
tmp/glm-landing/check-boundaries.js	local-only
tmp/glm-landing/norm-fixtures.js	local-only
tmp/glm-landing/mutate-g5b.js	local-only
tmp/glm-landing/selftest-t3.txt	local-only
tmp/glm-landing/gate-tree-t3.txt	local-only
tmp/glm-landing/verify-a12.js	local-only
tmp/glm-landing/wrap-enc.js	local-only
tmp/glm-landing/	local-only
```

Note: the `repo` rows are this slice's staged artifacts (including this report pair); the `local-only` rows are one-off process artifacts under `tmp/glm-landing/`, cleared at ⑧c; the former slice's `tmp/glm-eval*` leftovers are **not included**.
