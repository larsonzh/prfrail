# Slice `DIRECTIVE-REVIEW-SATURATION` validation report (directive v1.12 → v1.13)

> Status: **complete** (⑧b final; commit and push require same-turn authorization). Slice type: **medium (complex-leaning)**; size **M-L**; class `[SLICE]`; master tier High; ①⑤ Max; ⑦ Extra High.

## 1. Change summary (including the executive-freedom list)

| Item | Content |
|---|---|
| Slice | `DIRECTIVE-REVIEW-SATURATION`; directive **v1.12 → v1.13**; class `[SLICE]`; size **M-L**; type **medium (complex-leaning: pure `.md` plus `tools/gates` criteria code and fixtures)** |
| Delivered content | adds **§7.8 Review Saturation Criterion** (the sole authority for three stop objects × six situations, plus the stop-decision table, the convergence signals and the five traceability fields); adds criteria **G6-5** (OB-32 landed) and **G6-6** (OB-38 landed); closes **OB-34 / OB-35 / OB-36 / OB-41** in-slice; registers **OB-42 / OB-43 / OB-44 / OB-45**; one new §9.3 failure-mode row; the §0.3 v1.13 log row; a new appendix C **C.0i** ledger block; the §12.5 gap row turned to landed; appendix B.7.4's "complex / high risk" row synced to §7.1–§7.8 |
| Ledger closed set | `观察中` / `待办` / `待议` / `已处置` / `已裁决` / `已登记` (the EN table is isomorphic); **all 45 existing rows normalise, with no closed-set widening and no criterion relaxation** |
| Declared change set (`git diff HEAD` truth) | `docs/DELIVERY_DIRECTIVE.md` **+97/−26**, `_EN` **+97/−26** (symmetric, 1303 lines each); `tools/gates/lib/criteria/g6.js` **+311/−3**; `tools/gates/selftest.js` **+464/−3**; **8 new fixtures** (added by this slice; staged as of ⑧c) |
| Orchestration artifact | `docs/validation/evidence/directive-review-saturation-freedom-list.md` (the §1.7.3 fixed path; it sits outside G7's scan scope, so only ⑦ reviews it by hand) |
| Not touched | code, `docs/CONTRACTS`, `schemas`, `internal`, `cmd`, `testdata`, `examples`, `harnesses`, CI workflows, the role set and the model list; `g6-whitelist.txt` has zero entries and `cjk-newwords.txt` is unchanged |

### 1.1 Executive-freedom list (F-1 … F-12; entry count, numbering and §1.7.3 verdict wording consistent with the fixed-path list, checked at ⑧b)

| No. | Substantive choice | §1.7.1 basis | Affected area | §1.7.3 verdict |
|---|---|---|---|---|
| F-1 | The CN topic name adopts `评审收敛判据`; the slice ID stays the ASCII `DIRECTIVE-REVIEW-SATURATION`; new text in this slice **uses no character outside the baseline corpus** (the character once proposed for the topic name was blocked by G4b and then reworded) | contract wording completion | §7.8 and the whole slice's phrasing | all four questions no (wording changes no rule entity); confirmed by the user on 2026-09-23 |
| F-2 | The §7.8 insertion point is between §7.7 and §8 (the user named no exact position) | step order and batch splitting | §7.8 | all four questions no |
| F-3 | OB-38 takes "**add an independent G6-6 sub-check**" rather than "widen G6-1's vocabulary" | implementation-expression choice | §6.3 and `g6.js` | all four questions no; basis = G6-1's block-level semantics (a pending row and a completion marker in the same block go red) are **structurally incompatible** with §12.4's single large table — widening the table would make any pending row and any resolved row fail each other |
| F-4 | OB-34 takes "**retitle §7.2.2**" rather than "split the subsection" | implementation-expression choice | §7.2.2 | all four questions no; splitting would renumber §7.2.3 too and affect §2.4 and appendix B.7.4's references |
| F-5 | OB-41 is written into the **§9.3 failure-mode library**, not into §6.3's criterion-scope definition | compliance with the slice-scope declaration | §9.3 and §6.3 | all four questions no; §6.3's scope wording is the criterion body itself, and changing it would touch `GATES-TIGHTEN`'s area |
| F-6 | G6-5's label set includes "**line-leading bold labels**" (`^\*\*(\d+(?:\.\d+)*)\s`) | implementation within a scope the contract already fixes | `g6.js` / §6.3 | all four questions no; basis = this directive defines §1.7.1 / §1.7.2 / §1.7.3 and §3.4.1 and friends in the line-leading bold shape |
| F-7 | The closed set takes **6 words** (`观察中` / `待办` / `待议` / `已处置` / `已裁决` / `已登记`) and their EN mapping | contract wording completion | §6.3 and §12.4 | all four questions no; coverage check = **every existing ledger row (45) normalises** (see the validation report) |
| F-8 | G6-5's target regions = the §0.3 table interval plus appendix C.0–C.0i (**excluding** the C.1–C.3 v3.1 comparison table) | implementation within a scope the contract already fixes | `g6.js` / §6.3 | all four questions no; the exclusion is **structural immunity**, not a whitelist, so real dangling references are not masked |
| F-9 | Evidence capture switches to **node direct write** (`execFileSync` + `fs.writeFileSync('utf8')`), with `--suffix` so each round's reading is **stored separately** rather than overwriting the previous round's evidence | authorized in-package gate retry / orchestration artifact | `tmp/v113/**` evidence | all four questions no (the capture method changes no rule entity); the user approved its registration as OB-43 on 2026-09-23 |
| F-10 | This slice does **not** write back `DEV_PLAN`: a directive-governance slice's write-back scope = appendix C + §12.4 / §12.5 + the validation report | compliance with the slice-scope declaration | write-back scope | all four questions no; per OB-12's settled reading |
| F-11 | The remediation diff uses "`HEAD` + the previous round's full diff to **mechanically rebuild** the pre-remediation state → `git diff --no-index`" rather than a recollection-style description (**method artifact**: the `--with-remediation` branch of `tmp/v113/capture.js`; **output**: `tmp/v113/remediation.diff`) | orchestration artifact (evidence management) | `tmp/v113/remediation.diff` | all four questions no |
| F-12 | The ⑥ / ⑦ task packages carry the hard constraint "do not create / modify / delete any file; return the report as the final message; **do not claim to have run** `node` / `gate` commands" | orchestration-layer task package | ⑥ / ⑦ input | all four questions no; the user ordered on 2026-09-23 that **the master writes it at the orchestration layer and does not change the directive** |

### 1.2 User-ruling items not registered as freedom entries (7 items; item-by-item consistent with the fixed-path list's "二" section, checked at ⑧b)

1. **The registration of OB-39 / OB-40 / OB-41 and the C.0h note** — **explicit same-turn user authorization** (the quote from item 9 of the slice-start ruling: 「下一片启动第一件事：在 §12.4 追加 OB-39、OB-40（以及评估后的 OB-41+），在 C.0h / 报告记录顺延第 2 次及用户裁决。此后才进入主题工作。」) ⇒ not an exercise of freedom (⑦'s re-review Medium closed on this basis).
2. **The second pilot deferral** (including "keep executing, not as a pilot") — the user **ruled directly** on 2026-09-23, recorded independently per OB-37's meta-rule in §12.4's OB-40 and appendix C.0h's note; **not readable as a precedent**.
3. **The registration of OB-42 / OB-43 / OB-44 / OB-45** — ruled by the user item by item (scope, number and opening status word all named by the user).
4. **The literal garbled-character handling for OB-43**: the master recommended ⒝ (structural description, no literal characters registered, no whitelist widening), and the user **adopted ⒝** on 2026-09-23.
5. **OB-46 / OB-47 are not registered in this slice** (⑦'s 4 Low findings are criterion-boundary risks, not existing violations) — the user **ruled** on 2026-09-23: keep the numbers, register them formally in the next slice; that ruling is also §7.8's **first live application** of the stop criterion.
6. **D-2's "add the pointer" route** — the user stated the preference and ② judged and landed it; this is **review remediation**, not freedom exercised.
7. **OB-45's registration date stays 2026-09-23** (the ruling date, not the landing date) — the user **ruled** on 2026-09-23.

These 7 items match, item by item, the "二、不登记为自由度条目的项（用户裁定类）" section of the fixed-path list `docs/validation/evidence/directive-review-saturation-freedom-list.md` (checked at ⑧b); that list also carries a "三、⒜ 补审与额度例外的留痕" section, consistent with this report's §5 ⒜ supplementary re-review and round-2 re-review record and with §3 item 10 (checked at ⑧b).

## 2. Execution pipeline

| Step | Role / `modelId` | Status | Reworks |
|---|---|---|---|
| ① architecture | `deep-reasoner` / `deepseek-v4-pro` | `DESIGN: DONE` (1 call, cap 1) | 0 |
| ② implementation | `prfrail-implementer` / `deepseek-flash` | `IMPLEMENT: DONE` (4 calls, one DONE each: implementation 1 + ⑤ remediation D-1…D-6 1 + OB-43 / OB-44 registration 1 + OB-45 registration 1) | 1 |
| ③ testing | `prfrail-tester` / `deepseek-flash` | `TEST: DONE` (`SELFTEST: PASS 248/248`; assertions 202 → 248) | 1 |
| ④ master integration + gates | master (0 paid calls) | done | 0 |
| ⑤ pre-review | `deep-reasoner` / `deepseek-v4-pro` | `PRE-REVIEW: PASS WITH FIXES`, then after remediation `PRE-REVIEW: PASS` | 1 |
| ⑥ independent scan | `independent-reviewer` / `mai-code-1.1-flash` | `INDEPENDENT SCAN: PASS` (1 call, cap 3) | 0 |
| ⑦ final review | `independent-reviewer` / `gpt-5.3-codex` | `FINAL REVIEW: PASS` (1 final review; plus the user-ruled ⒜ supplementary re-review in the re-review slot and the round-2 re-review, see §5) | 1 |
| ⑧a / ⑧b documentation | `prfrail-documenter` / `deepseek-flash` | this report and its `_EN` mirror | 0 |
| ⑨ / ⑧c | master (0 paid calls) | done (readings in §11) | 0 |
| probes | — | **0 calls** | 0 |

## 3. Anomalies and fallbacks (11 items, one by one)

| No. | Scenario | Localisation | Disposition | Result |
|---|---|---|---|---|
| 1 | Both ① and ② called this slice's own 4 in-index registration rows "legacy leftovers" | they are in fact this slice's OB-39 / OB-40 / OB-41 and C.0h registrations, made per user-ruling item 9 | the master corrected it: they belong to this slice's own change set | not discarded, not disposed separately |
| 2 | ① used the "documentation-only degeneracy" wording in two places (the acceptance command and the "do not do" list) | this slice carries criteria code and fixtures, so it is **not documentation-only** | the master corrected it per the user ruling | ⑨ must include `node --check` and `selftest`; ②③⑧ are all carried by product-layer roles |
| 3 | D-2's cross-location inconsistency (introduced by this slice itself) | §7.8.4 claimed "§12.4's OB-11 / OB-37 disposition cells point at this section", while both cells had no pointer at the time ⇒ ⑤ reported Medium | remediation took the "**add the pointer**" route | the opening status word stays in the closed set, the in-row count of `\|` is unchanged, G6-6's `scanned 30 → 36` criterion actually strengthens, and the pointer is a reference form that copies no rule entity; recorded as instance evidence for OB-33's direction, and **this slice adds no OB** |
| 4 | The F-1 evidence-capture encoding trap | PowerShell pipeline capture of a native command's UTF-8 stdout decodes at the console code page (GBK) then re-encodes ⇒ the evidence files' Chinese becomes unrelated characters | evidence files switched to node direct write, with a per-file recheck of garbled markers | **OB-43** registered plus one §9.3 failure-mode row; after remediation the garbled-marker count is 0 (the first run measured 808 markers in `remediation.diff` and 22 in `selftest-master.txt`) |
| 5 | ⑤'s re-review wrote `tmp/v113/review-r2.md` to disk by itself | §2.2 says a review-layer role "may write files = no" | classed as a **literal deviation of the role boundary** (the file is its own report, it polluted no product text, and `tmp/` is not committed) | **OB-44** registered; this round the master wrote the hard constraint into the ⑥ / ⑦ task packages, **without changing the directive** |
| 6 | The master's freedom-list timing deviation | §1.7.3 requires the list to land **after ⑥ and before ⑦**, for ⑦ to review; in this slice it actually landed **after ⑦'s PASS** ⇒ ⑦ **did not review the list** | recorded plainly in this report and referred to the user for a ruling (use ⑦'s re-review slot, or waive it) | the user ruled on 2026-09-24: use ⑦'s re-review slot for a supplementary review (⒜) and allow the round-2 re-review as a quota exception (see item 10) |
| 7 | ③ self-caught an over-anchored assertion on its first run | `FAIL 247/248` (its own assertion text did not account for the `whitelist-source` suffix) | the assertion text was corrected, **with no criterion change and no semantic relaxation** | `PASS 248/248` afterwards |
| 8 | ② self-caught a G6-4 hit | the first draft wrote the two ledger positions in the "number + word" shape (OB-11 / OB-37), which G6-4 judged an expiring literal | rewritten to "OB-11's disposition cell / OB-37's disposition cell" | the row count is unchanged |
| 9 | The master's tool traps (disclosed plainly twice) | mechanically extracting the ⑤ / ⑦ reports first hit the master's own briefing text / reasoning text | tightened to "the report's opening sentence as prefix + a signature string" and re-extracted successfully | the corpus-external characters G4b blocked include the character once proposed for the topic name and the 4 characters in the freedom list's first draft (described structurally per OB-43's ruling, with no literal characters registered); all were **reworded**, and `cjk-newwords.txt` was **not widened**; separately G6-4 blocked the freedom list's first-draft counting shape (rewritten to "45 条") |
| 10 | the master's §1.7.3 timing deviation | the freedom list should land **after ⑥ and before ⑦** for ⑦ to review; in this slice it landed after ⑦'s first-round `PASS` ⇒ ⑦ did not review it in round 1 | the user ruled on 2026-09-24 to run a **⒜ supplementary review in the re-review slot** (`RE-REVIEW: FINDINGS`, 1 Medium + 1 Low) → re-register and add the anchor → the user's **quota exception** → round-2 re-review `RE-REVIEW: PASS` | both findings closed; candidate **OB-48** is not registered in this slice (opening status word `待办`), to land in §12.4 in the next slice together with OB-46 / OB-47 |
| 11 | one staging alignment before ⑧b | the index once held an older intermediate version (both directive files showed `MM`) | before ⑩ the master aligned the staged set with the report's declared change set by **per-file precise `git add`** (fifteen files) per §1.5 / G5 | `G7-b` turned from red to green (`32 artifact row(s) OK`) |

## 4. Environment and tiers

| Item | Value | Basis |
|---|---|---|
| Slice type / size | medium (complex-leaning) / **M-L** | §5.2 size → type mapping |
| Master | High (the medium tier) | §8.2 |
| ① / ⑤ | Max | §8.2 |
| ② ③ ⑧ | the same tier as the master (`deepseek-flash`) | §8.3 same-model roles cannot be split across tiers |
| ⑥ | the model's default / fixed tier | §8.2 |
| ⑦ | **Extra High (fixed)** | §8.2 |
| Tier-setting authority | the operator (set once before the run, never switched mid-run) | §8.1 |

## 5. Review-verdict summary and quota usage

| Step | Role / `modelId` | Verdict line | Quota | Findings |
|---|---|---|---|---|
| ⑤ pre-review (first round) | `deep-reasoner` / `deepseek-v4-pro` | `PRE-REVIEW: PASS WITH FIXES` | cap 1 | 2 Medium + 4 Low (D-1…D-6) |
| ⑤ re-review (after remediation) | `deep-reasoner` / `deepseek-v4-pro` | `PRE-REVIEW: PASS` | cap 1 (**⑤ total 2/2 spent**) | all 6 closed, no new defect, no scope creep |
| ⑥ independent scan | `independent-reviewer` / `mai-code-1.1-flash` | `INDEPENDENT SCAN: PASS` | cap 3 (**used 1/3**) | Sections A–E complete, 0 Medium+ |
| ⑦ first-round final review | `independent-reviewer` / `gpt-5.3-codex` | `FINAL REVIEW: PASS` | cap 1 final review + 1 re-review (first round used **1 + 0**) | no Critical / High / Medium; 4 Low criterion-boundary risks (not existing violations) |
| ⑦ ⒜ supplementary review (re-review slot) | `independent-reviewer` / `gpt-5.3-codex` | `RE-REVIEW: FINDINGS` | occupies the re-review slot (**1/1**) | 1 Medium + 1 Low: Medium = the freedom list's "user-ruling class" omitted **OB-39** (a listing omission, not an overreach; OB-39's registration had explicit same-turn user authorization = item 9 of the slice-start ruling); Low = **F-11 lacked a locatable evidence anchor** |
| ⑦ round-2 re-review (**user-ruled quota exception**) | `independent-reviewer` / `gpt-5.3-codex` | `RE-REVIEW: PASS` | beyond the §7.5 cap (ruled directly by the user on 2026-09-24) | both findings closed, no new issue, no Medium+; plus 1 **non-blocking** "needs a user ruling" item (whether to add an "external-evidence status declaration" to the two "user-ruling class" items that have no independent external evidence — **not adopted in this slice, left to the next slice's assessment**) |

**⒜ re-review's remediation targets**: the freedom list re-registers "user-ruling class" item 1 with the authorization quote (the original items 1–6 shift to 2–7); F-11 gains the method artifact (the `--with-remediation` branch of `tmp/v113/capture.js`) and the output (`tmp/v113/remediation.diff`) as anchors.

**Quota-exception traceability (per OB-37's meta-rule)**: ruling party = the user; date = 2026-09-24; reason = the corrected object is **the reviewed artifact itself** (the freedom list), and this is the **first timing deviation** since §1.7.3 landed; basis = the user's ruling text this round. **Stated plainly**: this entry **carries a rule-semantics change (breaching §7.5's quota cap), goes beyond OB-11's literal exception standard, and is a direct user ruling**; **it must not be read as a precedent**. This record is consistent with the freedom list's "三、⒜ 补审与额度例外的留痕" section (checked at ⑧b).

## 6. First application of the §7.8 stop criterion

⑦'s final review reported 4 **Low criterion-boundary risks**, with ⑦ itself noting "these are not existing violations in this slice". Under §7.8 this **triggered no new round**, entered no ⑦ re-review slot, and dispatched no ② re-registration. The user ruled on 2026-09-23: **not registered in this slice**, **keep numbers OB-46 (G6-5's region boundary depends on the `### C.1` heading staying stable) and OB-47 (`_EN` chooses its word table by filename suffix)**, to be registered formally in the next slice.

| Item | Content |
|---|---|
| Stop object | ⑦'s first round (final review) |
| Decision unit | this slice × a single review ring |
| Trigger situation | normal convergence (no Medium+ this round) |
| Basis sections | §4.3 + §5.3 + §7.8.3 |
| Traceability | the verdict line `FINAL REVIEW: PASS` is in this report's §5; the full ⑦ report is in appendix A |
| Disposition of the 4 Low findings | the user ruled they are not registered in this slice; numbers OB-46 / OB-47 are kept for formal registration in the next slice |
| Usable as a precedent | yes (normal flow) |

**Second application (escalation path)**: the ⒜ supplementary review reported **1 Medium + 1 Low** ⇒ after the §5.3 fixes ⑦ **must be re-run**, while the re-review slot (1 use) was already spent ⇒ the master **stopped under §7.5 and escalated for a user ruling** (**did not self-exempt**); on 2026-09-24 the user **ruled directly** to approve **round 2 of the re-review** (a quota exception, recorded independently under OB-37's meta-rule: ruling party = the user / date = 2026-09-24 / reason = the corrected object is the reviewed artifact itself and this is §1.7.3's first timing deviation / basis = the user's ruling text; **carries a rule-semantics change, goes beyond OB-11's literal standard, is a direct user ruling, and must not be read as a precedent**). Round 2's result `RE-REVIEW: PASS` (both findings closed, no new issue, no Medium+) ⇒ the ring is closed. ⇒ This record shows §7.8's **escalation path**: the correct action when the quota is exhausted = **stop + escalate + obtain a ruling**, never a self-granted pass.

**Third application (report-only changes after the quota was exhausted)**: the **4 report-only changes** after ⑦'s round-2 `RE-REVIEW: PASS` (the ⑧b finalization, the ⑧b additions to §7 / §6, the final-value numeric backfill, and the two-line self-contradiction micro-close) are ruled, under the **OB-11 exception standard**, to be **fact records / numeric backfill / wording close-out** — they carry no rule entity, change no criterion, change no verdict wording, change no criterion scope, and change no ledger-status closed set ⇒ **no ⑦ re-run is triggered**; the user **confirmed** that ruling on 2026-09-24.

⇒ §7.8 thereby gains **three field samples** in this slice: ⑦'s first-round 4 Low findings triggering no new round (**can stop**), the ⒜ Medium+ pausing to escalate (**escalation path**), and report-only changes after the quota was exhausted taking the OB-11 exception (**exception boundary**).

## 7. Review input-package summary (including the blind-review isolation proof)

- **Blind-review isolation (R2.6)**: ⑤ / ⑥ / ⑦ **attach no other party's checklist in their first round**; no round's report shows any trace of leaning on another party's conclusions.
- **What the master supplies**: the slice definition + the raw diff paths + the contract's original sentences + the read-only constraint.
- **The ⑥ / ⑦ task packages add a hard constraint**: do not create / modify / delete any file; return the report as the final message; do not claim to have run `node` / `gate` commands (the user ordered that **the master writes it at the orchestration layer and does not change the directive**).
- **Evidence landing (R6.1)**: before review the master lands the raw diff and the raw gate output under `tmp/v113/**` and hands the paths to the reviewer; no summary was substituted for the raw hunks.
- **⑤'s input**: `tmp/v113/review-input.diff` (first round), `tmp/v113/remediation.diff` and `tmp/v113/review-input-r2-r3.diff` (re-review).
- **⑦'s input**: `tmp/v113/review-input-r2-r3.diff` and `tmp/v113/review-input-r2.stat-r3.txt`.
- **⑦ ⒜ supplementary review (re-review slot)**: input package = the corrected executive-freedom list path + this slice's raw diff path (`tmp/v113/review-input-r2-r3.diff`) + the contract's original sentences + the read-only constraint + the hard constraint (do not create / modify / delete any file; return the report as the final message; **do not claim to have run** `node` / `gate` commands). Isolation proof: **no** ⑦ first-round report is attached, **no** ⑤ / ⑥ checklist is attached (the first-round isolation applies equally); ⑦ stated truthfully that it "ran no command" and marked its evidence with the master's landing paths.
- **⑦ round-2 re-review (quota exception)**: input scope identical to the previous line (the corrected list + the raw diff path + the contract's original sentences + the read-only constraint + the same hard constraint). Isolation proof: **no** ⑦ first-round report and ⒜ report, **no** ⑤ / ⑥ checklist attached.

## 8. Falsifiability (mechanical checks + mutation testing)

### 8.1 selftest T19–T28

`selftest` covers **T19–T28** for G6-5 / G6-6, including the registration and assertion chain of the 8 fixtures; `SELFTEST: PASS 248/248`, with the assertion total rising from 202 to 248.

### 8.2 Two source-level bidirectional mutations

| Mutation | Injection | Expectation | Measurement |
|---|---|---|---|
| ② | remove the line-leading bold collection from `labelsFromLines` | §1.7.1 dangles and goes red (T25) | `applied=true`; `g6.js` is byte-identical after restore |
| ③ | delete `观察中` from `LEDGER_STATUS_SETS.cn` | that row goes red (T28) | `applied=true`; `g6.js` is byte-identical after restore |

### 8.3 Constructive counterexamples (T20 / T21 / T24 / T27)

| Probe | Construction | Expectation and measurement |
|---|---|---|
| T20 | write a dangling `§7.9` reference on a line inside the §0.3 region | red with exactly one finding; the reverse mutation `§7.9` → `§7.8` turns it green (the assertion tracks the reference, not the line) |
| T21 | a row of appendix C.1's v3.1 comparison table | **structurally immune** outside the region (`refs 0`); the boundary mutation (renaming the `### C.1` heading) pulls the row back in and turns it red |
| T24 | write `§7.9` on a line outside every region | not judged (`scanned 6, refs 0`); the boundary mutation (dropping the closing heading) extends the region onto that line and turns it red |
| T27 | write `**下一片必做**` / `**已落地**` as a ledger row's opening status word | exactly two findings, with the legal control row (`已登记`) not reported; the reverse mutation (both changed to closed-set words) turns it green |

## 9. Gate results (⑩ pre-final values; raw output under `tmp/v113/**`)

**Declared change set (`git diff HEAD --numstat`, ⑩ pre-final values)**

| File | Delta |
|---|---|
| `docs/DELIVERY_DIRECTIVE.md` | `97+/26−` |
| `docs/DELIVERY_DIRECTIVE_EN.md` | `97+/26−` |
| `docs/validation/directive-review-saturation.md` | `352+/0−` |
| `docs/validation/directive-review-saturation_EN.md` | `352+/0−` |
| `docs/validation/evidence/directive-review-saturation-freedom-list.md` | `41+/0−` |
| `tools/gates/lib/criteria/g6.js` | `311+/3−` |
| `tools/gates/selftest.js` | `464+/3−` |
| fixtures `g6-5-*` / `g6-6-*` (8) | `17 / 5 / 5 / 15 / 9 / 10 / 7 / 10 +` (no deletions in either column) |

| Gate | Reading |
|---|---|
| `node tools/gates/gate.js --all --scope=tree` | `TOTAL_FAIL=0` (`tmp/v113/gate-tree-r1-r10.txt`) |
| `node tools/gates/gate.js --all --scope=index` | `TOTAL_FAIL=0` (`tmp/v113/gate-index-r1-r10.txt`) |
| `G7-b` | PASS (`32 artifact row(s) OK`) |
| `node tools/gates/gate.js --selftest` | `SELFTEST: PASS 248/248` (`tmp/v113/selftest-master-r10.txt`) |
| Bilingual mirror | `per-line mismatches = 0`, `MIRROR RESULT: PASS` (equal line counts on both sides) (`tmp/v113/mirror-r10.txt`) |
| The master's independent acceptance, `tmp/v113/verify-v113.js` | `ACCEPTED`: **45 ledger rows** in each file, all inside the closed set; 186 `§X.Y` references in-region, 0 dangling (`tmp/v113/verify-v113-r10.txt`) |
| `G6-5` | PASS `scanned 48, refs 47` (`tmp/v113/g6-5-r1-r10.txt`) |
| `G6-6` | PASS `scanned 38` (`tmp/v113/g6-6-r1-r10.txt`) |
| `G6-4` | PASS `scanned 939` |
| `G6-1` | PASS `scanned 939` |
| `G3(a)` | the two directive files `+97/-26 ↔ +97/-26`; this report's two files are equal and symmetric on both sides (the ⑧b backfill's unstaged increment over the ⑧a staged version); all five shapes true |
| `G4a` | `15 file(s) OK` |
| `G4b` | `unknown=[]` (`newCjk=704`) |
| `node --check tools/gates/lib/criteria/g6.js` | exit 0 |
| `node --check tools/gates/selftest.js` | exit 0 |

**Reading note**: all readings in this section are the ⑩ pre-final values after precise staging; the `352+/0−` for this report's two files is exactly that final value (the numeric backfill itself does not change line counts, so it has converged). Every other row is likewise drawn from the same pre-⑩ final batch.

**Second reading note**: `git diff HEAD --numstat` (the sole authoritative reading) and the gate's `G3(a)` tree-scope reading agree on `+97/-26` for the directive files; this report's two files show an equal, symmetric `G3(a)` increment from the same tree scope's unstaged reading (the ⑧b backfill against the ⑧a staged version, likewise quoted without a literal) over the `repo` artifacts in §14, which zero out once the master stages precisely at ⑩. The `+94/-27` the user named is the r2 / r3 reading from **before OB-45 was registered**.

## 10. ⑧c wrap-up cleanliness record (pointer)

The ⑧c record is in §11 (readings supplied by the master, copied here verbatim).

## 11. ⑧c wrap-up cleanliness record

⑧c is executed by the master after ⑧b; the readings below are the ⑧c re-run results, supplied by the master and copied here verbatim.

| Field | Reading |
|---|---|
| Executor | the master (0 paid calls) |
| Whole-tree diagnostics (§3.4 category table) | `gate.js --all --scope=tree` ⇒ `TOTAL_FAIL=0`; `--scope=index` ⇒ `TOTAL_FAIL=0`; `G7-b` PASS (`32 artifact row(s) OK`) (`tmp/v113/gate-tree-r1-r10.txt` / `tmp/v113/gate-index-r1-r10.txt`) |
| IDE diagnostics summary | this session has no Problems-panel reading interface; covered by `node --check` and `gate.js --selftest` (`SELFTEST: PASS 248/248`) |
| `gate.js --all --scope=tree` | `TOTAL_FAIL=0` (`tmp/v113/gate-tree-r1-r10.txt`) |
| Encoding gate (whole tree, per file) | `enc-tree` ⇒ `tracked=2923 text=2834 binary=89 violations=3`; all 3 are pre-existing violations in **OB-27's frozen evidence packs**, handled by the temporary exemption plus the self-consistency check (`tmp/v113/enc-tree-r10.txt`) |
| Frozen-evidence-pack self-consistency | `evidence-selfcheck` ⇒ `packs=4 entries=2327 missing=0 mismatch=0 not-self-checkable=1` (`tmp/v113/evidence-selfcheck-r10.txt`) |
| Temp directories / probe residue / untracked files | `git status --porcelain -uall` ⇒ the staged set equals the change set the report declares (all fifteen files staged, no `MM`, no `??` untracked residue; pre-⑩ final values) |
| Verdict line | no failing item among the readings above; the 3 encoding violations are pre-existing exemptions of OB-27's frozen evidence packs |

## 12. Cost and metering (R7.1, row by row)

| Role | `modelId` | Stage | Count | Cap |
|---|---|---|---|---|
| ① architecture | `deepseek-v4-pro` | pre-analysis | 1 | 1 |
| ② implementation | `deepseek-flash` | implementation + ⑤ remediation D-1…D-6 + OB-43 / OB-44 registration + OB-45 registration | 4 | product layer has no §7.5 cap |
| ③ testing | `deepseek-flash` | testing | 1 | product layer has no §7.5 cap |
| ④ master | — | integration + gates | 0 paid calls | — |
| ⑤ pre-review | `deepseek-v4-pro` | pre-review 1 + post-remediation re-review 1 | 2 | 1+1 (**spent**) |
| ⑥ independent scan | `mai-code-1.1-flash` | independent scan | 1 | 3 |
| ⑦ final review | `gpt-5.3-codex` | final review 1 + ⒜ supplementary review 1 + round-2 re-review 1 | 3 | 1+1 (plus 1 user-ruled exception) |
| ⑧ documentation | `deepseek-flash` | ⑧a / ⑧b | 1 | — |
| ⑨ / ⑧c | — | native validation / wrap-up cleanliness | 0 paid calls | — |
| probes | — | — | **0** | — |

## 13. Not executed, and next steps

### 13.1 Not executed

| Item | Destination |
|---|---|
| OB-46 (G6-5's region boundary depends on the `### C.1` heading staying stable) | not registered in this slice; formal registration in the next slice (landing in §12.4 together with OB-47 / OB-48) |
| OB-47 (`_EN` chooses its word table by filename suffix) | not registered in this slice; formal registration in the next slice (landing in §12.4 together with OB-46 / OB-48) |
| ⑦'s review of the freedom list | the user ruled on 2026-09-24 to use ⑦'s re-review slot for a supplementary review (⒜): `RE-REVIEW: FINDINGS` → remediation fully closed → round-2 re-review `RE-REVIEW: PASS` (quota exception, see §5) |
| `DEV_PLAN` write-back | not written back in this slice (F-10); handled by a successor slice |
| OB-30 / OB-31 / OB-33 | split out; assessed together with `GATES-G1G5` / `GATES-TIGHTEN` |
| IDE diagnostics reading | this session has no Problems-panel reading interface; ⑧c covers it with `node --check` and `gate.js --selftest` (see §11) |

### 13.2 Next-step recommendations

- ⑩ stop point: ⑦'s final review `FINAL REVIEW: PASS` (no Medium+) and ⑤'s re-review `PRE-REVIEW: PASS` ⇒ this slice may be committed; **commit and push require same-turn authorization**.
- Next-slice candidates: `GATES-G1G5`, `GATES-TIGHTEN`; OB-42 (the verdict-line closed set still has no mechanical criterion) and OB-43 / OB-44 are disposed per their registration.
- Tier recommendation: directive-governance slices keep the master at High, ①⑤ at Max and ⑦ at Extra High; ②③⑧ share `deepseek-flash` with the master, so their tiers cannot be split (§8.3).
- Candidate slice `DIRECTIVE-LEDGER-ARCHIVE`: nature = `[SLICE]` (touches §6.3's criterion scope plus a new §12 archival section); content = archive closed OBs from §12.4 to appendix D, narrow G6-6's scope to the still-open region, separate the human-reading domain from the scan domain, and guard against OB-28's "an untouched expiring row stays silent forever"; dependencies = OB-28, OB-33; queue position = after `GATES-G1G5` / `GATES-TIGHTEN` / `REWORK-THRESHOLD` (and the OB-33 slice, if it is split out), as the directive-governance closing slice (not registered as an OB in this slice).
- First thing at the next slice's start: land **OB-46 / OB-47 / OB-48** (plus the master's Sol-call registration) formally in §12.4 so they gain G6-6's mechanical coverage.

### 13.3 Deviations and items needing a master ruling

| Item | Fact | Recommendation |
|---|---|---|
| The `G3(a)` reading | the ⑧c re-run gives `+97/-26 ↔ +97/-26` for the two directive files; this report's two files read equal and symmetric on both sides (an unstaged backfill increment, quoted without an expiring literal); the summary the user named says `+94/-27` (the r2 / r3 reading from before OB-45 was registered) | take the ⑧c re-run readings (the raw output is on disk) |
| The declared change-set reading | `git diff HEAD --numstat` gives **+97/−26** (directive); this report's two files are pure additions whose counts move with the ⑧b backfill (hence no literal is quoted); the gate's `G3(a)` agrees on `+97/-26` for the directive files and takes the unstaged increment on this report | the readings differ by definition, and both are listed plainly |
| The freedom-list timing | §1.7.3 requires the list to land after ⑥ and before ⑦ for ⑦ to review; in this slice it landed after ⑦'s first-round `PASS` | the user ruled on 2026-09-24 to use ⑦'s re-review slot for a supplementary review (⒜) and to close with a quota exception: ⒜ reported `RE-REVIEW: FINDINGS` (1 Medium = the list omitted OB-39, 1 Low = F-11 lacked an anchor), and after remediation the round-2 re-review reported `RE-REVIEW: PASS`; the record is in §5 and in the freedom list's "三" section |
| Appendix A's citation marker | the ⑦ lines that collide with `G6-4`'s expiring-literal shape carry a `引文` note under that criterion's own exclusion mechanism (the rendered text is unchanged) | no new whitelist entry, no gate relaxation |
| Blocked characters | the character once proposed for the topic name and the 4 characters in the freedom list's first draft are corpus-external, so they are **described structurally with no literal characters registered** (per OB-43's ruling) | do not widen `cjk-newwords.txt` |
| G7-b's staging precondition | the `repo` disposition requires the path to be **staged and present in the worktree**; when this report landed, the 8 fixtures, the freedom list and both reports were not staged, so `G7-b` once reported "not staged" | the master aligned them before ⑩ by per-file precise `git add` (fifteen files) ⇒ `G7-b` PASSes (`32 artifact row(s) OK`) (the documenter does not stage on the master's behalf) |

## 14. Artifact list

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
tools/gates/lib/criteria/g6.js	repo
tools/gates/selftest.js	repo
tools/gates/testdata/fixtures/g6-5-c1-v31-table.txt	repo
tools/gates/testdata/fixtures/g6-5-dangling-ref.txt	repo
tools/gates/testdata/fixtures/g6-5-excluded-ref.txt	repo
tools/gates/testdata/fixtures/g6-5-legal-refs.txt	repo
tools/gates/testdata/fixtures/g6-5-outside.txt	repo
tools/gates/testdata/fixtures/g6-6-en-status.txt	repo
tools/gates/testdata/fixtures/g6-6-legacy-status.txt	repo
tools/gates/testdata/fixtures/g6-6-legal-status.txt	repo
docs/validation/evidence/directive-review-saturation-freedom-list.md	repo
docs/validation/directive-review-saturation.md	repo
docs/validation/directive-review-saturation_EN.md	repo
tmp/v113/**	local-only
```

**Post-commit archival (next batch, not part of this commit)**: the ⑤⑥⑦ reports (`review-r1.md` / `review-r2.md` / `review-r3.md`), the full diff (`review-input-r2-r10.diff`) and the final gate readings (`gate-tree-r1-r10.txt` / `selftest-master-r10.txt` / `mirror-r10.txt`) are planned for archival under `docs/validation/evidence/directive-review-saturation/`; at that point the disposition of the corresponding rows in this report's artifact list changes from `local-only` to `repo`. **Reason the archival is deferred (recorded plainly)**: those reports contain cited literals (such as line counts / counts), and entering `docs/validation/**` directly as `.md` would fall inside **G6's `.md` scan domain** and trigger **G6-4**; rewriting the reviewer's original text is unacceptable and widening the whitelist is forbidden ⇒ the next batch archives them with **an extension outside the `.md` scan domain** (or disposes of them per §6.3's terms in the next slice).

## 15. Appendix A: ⑦ final review (round 1, Codex) verbatim

> Below is the verbatim text of `tmp/v113/review-r3.md`; the one line colliding with `G6-4`'s expiring-literal shape carries a `引文` note under that criterion's own exclusion mechanism (the rendered text is unchanged), and every other line is verbatim.

未发现 Critical/High/Medium/Low 级现存违规项。以下结论全部基于主控落盘证据与只读核对，我未执行 node、gate.js、selftest 命令。

主控落盘证据：
1. [tmp/v113/gate-tree-r1-r3.txt](tmp/v113/gate-tree-r1-r3.txt#L19), [tmp/v113/gate-tree-r1-r3.txt](tmp/v113/gate-tree-r1-r3.txt#L20), [tmp/v113/gate-tree-r1-r3.txt](tmp/v113/gate-tree-r1-r3.txt#L29)
2. [tmp/v113/gate-index-r1-r3.txt](tmp/v113/gate-index-r1-r3.txt#L19), [tmp/v113/gate-index-r1-r3.txt](tmp/v113/gate-index-r1-r3.txt#L20), [tmp/v113/gate-index-r1-r3.txt](tmp/v113/gate-index-r1-r3.txt#L29)
3. [tmp/v113/selftest-master-r3.txt](tmp/v113/selftest-master-r3.txt#L253)
4. [tmp/v113/mirror-r3.txt](tmp/v113/mirror-r3.txt#L3)
5. [tmp/v113/verify-v113-r3.txt](tmp/v113/verify-v113-r3.txt#L1), [tmp/v113/verify-v113-r3.txt](tmp/v113/verify-v113-r3.txt#L2), [tmp/v113/verify-v113-r3.txt](tmp/v113/verify-v113-r3.txt#L6)
6. [tmp/v113/review-input-r2.stat-r3.txt](tmp/v113/review-input-r2.stat-r3.txt#L1-L4)

**A 安全审查**
- 无发现。未见门禁弱化、白名单扩张、或绕过路径引入。

**B 架构一致性审查**
- 无发现。单一权威收敛关系已写入并互引一致：见 [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md#L407-L408), [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md#L567), [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md#L584)；EN 同步见 [docs/DELIVERY_DIRECTIVE_EN.md](docs/DELIVERY_DIRECTIVE_EN.md#L407-L408), [docs/DELIVERY_DIRECTIVE_EN.md](docs/DELIVERY_DIRECTIVE_EN.md#L567), [docs/DELIVERY_DIRECTIVE_EN.md](docs/DELIVERY_DIRECTIVE_EN.md#L584)。
- G6-5/G6-6 正文与实现对齐：规则声明见 [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md#L491), [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md#L493)；实现锚点见 [tools/gates/lib/criteria/g6.js](tools/gates/lib/criteria/g6.js#L508), [tools/gates/lib/criteria/g6.js](tools/gates/lib/criteria/g6.js#L538), [tools/gates/lib/criteria/g6.js](tools/gates/lib/criteria/g6.js#L544), [tools/gates/lib/criteria/g6.js](tools/gates/lib/criteria/g6.js#L666)。

**C 完整性审查**
- 无发现。切片改动集严格落在 4 文件： [tmp/v113/review-input-r2.stat-r3.txt](tmp/v113/review-input-r2.stat-r3.txt#L1-L4)。
- 夹具映射与断言链覆盖到 T19-T28：注册见 [tools/gates/selftest.js](tools/gates/selftest.js#L55-L62)，关键断言见 [tools/gates/selftest.js](tools/gates/selftest.js#L1182), [tools/gates/selftest.js](tools/gates/selftest.js#L1251), [tools/gates/selftest.js](tools/gates/selftest.js#L1456), [tools/gates/selftest.js](tools/gates/selftest.js#L1519)，主控落盘 PASS 见 [tmp/v113/selftest-master-r3.txt](tmp/v113/selftest-master-r3.txt#L253)。
- C.0i 与 OB-42~45 已落文：见 [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md#L1214), [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md#L823-L826)；EN 同位见 [docs/DELIVERY_DIRECTIVE_EN.md](docs/DELIVERY_DIRECTIVE_EN.md#L1214), [docs/DELIVERY_DIRECTIVE_EN.md](docs/DELIVERY_DIRECTIVE_EN.md#L823-L826)。

**D 测试反例（按要求给最小反例；以下为判据边界风险，不是本片现存违规）**
- LOW | [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md#L491) | 本子查只覆盖“引用悬空”子类，不覆盖 C1 的“日志声称已改而正文无内容”全类 | 若要堵住该漏报面，后续片新增内容级对账机械判据。
- LOW | [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md#L491), [tools/gates/selftest.js](tools/gates/selftest.js#L1251), [tmp/v113/selftest-master-r3.txt](tmp/v113/selftest-master-r3.txt#L215) | 附录 C.1 的 v3.1 对照表与 C.2/C.3 不在目标区域 | 保持 C.1 边界标题稳定；或在实现中补边界别名防御，避免标题漂移导致扩圈误报。
- LOW | [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md#L475), [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md#L809) | 只判变更行，历史行不追溯 | 若要堵住该漏报面，后续增加周期性全量台账扫描子查。
- LOW | [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md#L493), [tools/gates/lib/criteria/g6.js](tools/gates/lib/criteria/g6.js#L666), [tools/gates/selftest.js](tools/gates/selftest.js#L1519), [tmp/v113/selftest-master-r3.txt](tmp/v113/selftest-master-r3.txt#L250) | EN 镜像使用英文闭集 | 增加命名守卫或显式语言元数据，避免仅靠文件名后缀选词表的误报面。

**三问必答**
1. §12.4 全 45 行闭集核验：通过。主控落盘显示 CN/EN 均为 45 行且 0 越界词，见 [tmp/v113/verify-v113-r3.txt](tmp/v113/verify-v113-r3.txt#L1-L2) 与 [tmp/v113/verify-v113-r3.txt](tmp/v113/verify-v113-r3.txt#L6)。我也逐行复核了 [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md#L770-L826) 与 [docs/DELIVERY_DIRECTIVE_EN.md](docs/DELIVERY_DIRECTIVE_EN.md#L770-L826)。 <!-- 引文 -->
2. G6-5 区域/排除域与实现一致性：通过。正文明确了目标域、排除域、行首加粗标签、父节前缀解析，见 [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md#L491)；实现对应于 [tools/gates/lib/criteria/g6.js](tools/gates/lib/criteria/g6.js#L508), [tools/gates/lib/criteria/g6.js](tools/gates/lib/criteria/g6.js#L538), [tools/gates/lib/criteria/g6.js](tools/gates/lib/criteria/g6.js#L544)。C.1 排除与父节/加粗标签均有断言证据，见 [tools/gates/selftest.js](tools/gates/selftest.js#L1251), [tools/gates/selftest.js](tools/gates/selftest.js#L1182), [tools/gates/selftest.js](tools/gates/selftest.js#L1363) 及主控落盘 [tmp/v113/selftest-master-r3.txt](tmp/v113/selftest-master-r3.txt#L215), [tmp/v113/selftest-master-r3.txt](tmp/v113/selftest-master-r3.txt#L230)。
3. 不做清单逐项未触碰：结论为全部未触碰。依据为改动集仅 4 文件 [tmp/v113/review-input-r2.stat-r3.txt](tmp/v113/review-input-r2.stat-r3.txt#L1-L4)。

**边界核验（逐项）**
1. 不改 G6-1~G6-4 规则实体与词表：未触碰。见改动集 [tmp/v113/review-input-r2.stat-r3.txt](tmp/v113/review-input-r2.stat-r3.txt#L1-L4)。
2. 不新增白名单条目（g6-whitelist 与 cjk-newwords）：未触碰。见改动集 [tmp/v113/review-input-r2.stat-r3.txt](tmp/v113/review-input-r2.stat-r3.txt#L1-L4)，目标文件为 [tools/gates/g6-whitelist.txt](tools/gates/g6-whitelist.txt), [tools/gates/cjk-newwords.txt](tools/gates/cjk-newwords.txt)。
3. 不改 G7 与其实现文件：未触碰。目标文件 [tools/gates/lib/criteria/g7.js](tools/gates/lib/criteria/g7.js) 不在改动集。
4. 不改 gate 主框架相关文件（gate.js、ctx.js、report.js、r25.js）：未触碰。目标文件 [tools/gates/gate.js](tools/gates/gate.js), [tools/gates/lib/ctx.js](tools/gates/lib/ctx.js), [tools/gates/lib/report.js](tools/gates/lib/report.js), [tools/gates/lib/criteria/r25.js](tools/gates/lib/criteria/r25.js) 不在改动集。
5. 不改 docs/CONTRACTS 与 docs/DEV_PLAN 与 docs/PRODUCT_REQUIREMENTS 与 docs/CODING_CONVENTIONS：未触碰。目标文件 [docs/CONTRACTS.md](docs/CONTRACTS.md), [docs/DEV_PLAN.md](docs/DEV_PLAN.md), [docs/PRODUCT_REQUIREMENTS.md](docs/PRODUCT_REQUIREMENTS.md), [docs/CODING_CONVENTIONS.md](docs/CODING_CONVENTIONS.md) 及对应 EN 镜像均不在改动集。
6. 不改 docs/validation：未触碰。目录 [docs/validation](docs/validation) 不在改动集。
7. 不改 internal/cmd/schemas/testdata/examples/harnesses：未触碰。目录 [internal](internal), [cmd](cmd), [schemas](schemas), [testdata](testdata), [examples](examples), [harnesses](harnesses) 均不在改动集。
8. 不新增或修改 CI workflow：未触碰。目录 [.github/workflows](.github/workflows) 不在改动集。
9. 不改角色与模型名单定义：未触碰。目录 [.github/agents](.github/agents) 不在改动集。
10. 不扩散到其他文档体系：未触碰。改动仅在 [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md), [docs/DELIVERY_DIRECTIVE_EN.md](docs/DELIVERY_DIRECTIVE_EN.md)。
11. 仅实现目标代码范围（G6 与 selftest）：满足。改动代码仅 [tools/gates/lib/criteria/g6.js](tools/gates/lib/criteria/g6.js), [tools/gates/selftest.js](tools/gates/selftest.js)。

**无法验证的部分**
1. 我未执行任何命令，无法独立重放 gate 与 selftest，只能引用主控落盘。
2. 无法证明落盘文件生成后是否被外部改写；本结论默认这些证据文件可信且对应同一工作树快照。

**什么会推翻我的结论**
1. 若对同一提交重放得到与 [tmp/v113/gate-tree-r1-r3.txt](tmp/v113/gate-tree-r1-r3.txt) 或 [tmp/v113/selftest-master-r3.txt](tmp/v113/selftest-master-r3.txt) 不一致结果。
2. 若出现额外改动文件超出 [tmp/v113/review-input-r2.stat-r3.txt](tmp/v113/review-input-r2.stat-r3.txt#L1-L4)。
3. 若 §12.4 台账在我审阅后被继续改写而未同步复核。

**需用户裁决项**
1. 无（本轮未触发 §1.7.2 升级边界）。

FINAL REVIEW: PASS

## 16. Appendix B: ⑤ / ⑥ verdict lines and key passages

| Step | Report (local-only) | Verdict line | Key passage |
|---|---|---|---|
| ⑤ pre-review (first round) | `tmp/v113/review-r1.md` | `PRE-REVIEW: PASS WITH FIXES` | 2 Medium (D-1 / D-2) + 4 Low (D-3…D-6); D-1 is §7.8.6 mis-citing §12.1's companion ②; D-2 is §7.8.4's pointer claim not matching the OB-11 / OB-37 cells |
| ⑤ re-review | `tmp/v113/review-r2.md` | `PRE-REVIEW: PASS` | D-1…D-6 all closed; 0 new defects; no scope creep; the true increment = 9 insertions / 8 deletions per file |
| ⑥ independent scan | returned as the final message (no standalone file on disk) | `INDEPENDENT SCAN: PASS` | Sections A–E complete; 0 Medium+ |
