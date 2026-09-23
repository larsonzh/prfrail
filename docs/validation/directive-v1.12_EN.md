# T027 · DIRECTIVE-V1.12 directive-revision validation report (three rounds of third-party independent review disposition + user-ruling landing)

> Status: **complete (commit and push require same-turn authorization)**. Slice type: **medium** (documentation only, multiple files / rule semantics involved); tier: `high`; ⑦ always Extra High.

## 1. Change summary (including the executive-freedom list)

| Item | Content |
|---|---|
| Changed files | `docs/DELIVERY_DIRECTIVE.md` and its strict positional mirror `docs/DELIVERY_DIRECTIVE_EN.md` (version v1.11 → **v1.12**); this report and its mirror are new; the freedom list at the §1.7.3 fixed path is new |
| Disposition scale | **64 findings** accumulated across three rounds of third-party independent review (**1 Critical + 11 High + 35 Medium + 17 Low**) disposed item by item, plus the A3 and H4 user rulings landed (64 = round 2's 28 + round 3's 14 + the pre-commit 22, all three already inside the 64); additionally the **post-disposition re-verification round**'s **12 findings**, the **fourth re-review (R6)**'s **10 findings** and the **fifth re-review (R7)**'s **9 findings** were disposed (see §3.5 to §3.9) |
| New rule content | §0.1 authority table plus 4 rows and a terminology disambiguation; §0.2 MAJOR/MINOR ruling note; §0.3 new "affected in-flight slices" column; §1.2 three new fields; §1.6 production-code boundary plus counting unit plus multi-trigger hardening; §1.7 "§1.6 boundary outranks" plus the three-branch test; §2.2 two new columns (modelId / UI display name); §4.1 closed-set verdict table; §4.2 and §5.1 ②a/②b (contract first); §5.1 ⑧c frozen-evidence-pack exemption; §6.1 focused-gate definition; §6.3 G1 stopped-state exemption; §12.1 fixed denominator and `[SLICE]` kickoff handling; §12.2 re-evaluation clause and the DR-FIX decision record |
| New ledger entries | 附录 C **C.0h**; §12.4 new observation items **OB-28 – OB-33** (with OB-32 promoted to mandatory for the next slice); §12.5 gap row turned to "landed" |
| Orchestration artifact | `docs/validation/evidence/directive-v1.12-freedom-list.md` (the §1.7.3 fixed path; it sits outside G7's scan scope, so only ⑦ reviews it by hand) |
| Not touched | the existing judgement logic (only two exemption clauses are added: A10's stopped-state exemption and H5's ⑧c frozen-evidence-pack exemption); the **structure** of the §5.2 role-enablement matrix (only M1's blind-review wording changed); the §7.5 cost ceiling values; the role set and carriers; code / CONTRACTS / schema / fixtures / CI |

### 1.1 Executive-freedom list (§1.7.3 four questions: (1) does it change axioms or contract semantics; (2) is it reversible; (3) does it need a user ruling; (4) is it inside the slice scope)

| No. | Substantive choice | §1.7.1 basis | Affected area | §1.7.3 verdict |
|---|---|---|---|---|
| F-1 | ④ implementation landed in four batches (附录 B → §12 pilot closure → text sync → smaller new items) rather than one | step order and batch splitting | ④ | all four questions no ⇒ executive freedom (batch order changes no rule entity) |
| F-2 | The 19 C1 items follow the same four-layer batch order and land in the same batches | step order and batch splitting | ④ | all four questions no |
| F-3 | New observation items numbered consecutively OB-28 to OB-33, no gaps and no reuse of closed numbers | ledger numbering completion | §12.4 | all four questions no (number allocation carries no semantics) |
| F-4 | The 12 historical rows of the §0.3 "affected in-flight slices" column all read `历史缺失`, with no per-row archaeology | wording-level unification of a factual reading | §0.3 | all four questions no (historical reconstruction is outside this slice) |
| F-5 | The C.0h "field-table applicability note" points at the three fields in one sentence instead of inventing one row per field | contract wording completion | 附录 C.0h | all four questions no; and **manufacturing field values was deliberately avoided** (it would have produced values contrary to fact) |
| F-6 | M4 disposition: no new `tools/gates/g6-whitelist.txt` entry, only a measured finding recorded in this report | implementation-expression choice | gate whitelist | all four questions no; grounded in a **reproducible measurement** (the shipped G6-4 criteria score 0 hits on the "NN 条" / "NN 项" shapes) |
| F-7 | M5 recorded OB-32 as "mandatory for the next slice" only; G6-5 is not implemented here | compliance with the slice-scope declaration | §12.4 and successor slices | all four questions no; implementing G6-5 needs `gate.js` plus `selftest` plus tests, beyond this slice's authorization |
| F-8 | The §12.1 reading note sits after the baseline table rather than becoming a table column | contract wording placement | §12.1 | all four questions no (placement does not affect the reading) |
| F-9 | The "three rounds" reading was unified in one pass across the §0.3 v1.12 log row and the C.0h heading | wording-level unification of a factual reading | §0.3 and 附录 C.0h | all four questions no (the two statements had to agree anyway) |
| F-10 | User-ruling items (A3 / H4 / M1 / M2a / M2b / M3 / the M5 priority promotion) are not listed here | §1.7.3 records only master discretion | whole slice | these are **user rulings**, not freedom exercised; they are named in §5 instead, so they are not confused with the "(1) recommendation + master choice" pairs of F-6 and F-7 |
| F-11 | L1 to L6 are not fixed in this slice per user ruling and are recorded in the "not executed" section of this report | compliance with the slice-scope declaration | whole slice | landing a user ruling, not freedom exercised |
| F-12 | ⑨ degenerates for a documentation slice to "encoding gate + G1 to G7 self-check + link and reference resolvability check" | the §5.2 degeneracy reading | ⑨ | all four questions no (the degenerate form is fixed by §5.2; there is no second option) |
| F-13 | ⑧ is carried out by the master, with no sub-agent dispatched | the startup declaration's budget line lists only ①⑤⑥⑦ | ⑧a and ⑧b | all four questions no; a wrong call here changes the **execution role**, not code or contract semantics |
| F-14 | This report follows the v1.11 report's nine-section order (change summary → pipeline → item-by-item disposition → gates → ⑧c → not executed → cost → next steps → artifact list) and adds "C1 event" and "user-ruling provenance" | structural choice for an orchestration artifact | this report (eleven sections) | all four questions no (report structure is not a contract entity) |
| F-15 | R6's N7, N8 and N9 (the user ruled they are "left to the next slice") are not traced by report number but registered in this directive's observation ledger as OB-34 to OB-36 | ledger numbering completion (same reading as F-3) | §12.4 and the successor slices | all four questions no (number allocation carries no semantics) |
| F-16 | The "not touched" row of this report's §1 adopts the N6-corrected wording of Appendix C.0h rather than changing C.0h alone | wording-level unification of a factual reading (same reading as F-9) | this report §1 and Appendix C.0h | all four questions no (the two statements had to agree anyway) |
| F-17 | P6 is landed with the fuller reading "64 findings plus 22 findings from the two later re-reviews" (rather than either of the two options the user offered), so that the block covers both R5 and R6 | wording-level completion of a factual reading (same reading as F-9) | Appendix C.0h's goal row | all four questions no (only the numeric range changes, no rule entity) |
| F-18 | P1's timing resolution is written as "**the re-evaluation point is not cancelled, only merged into the judgement at the new pilot's wrap-up**" (the re-evaluation point the user ruled on earlier is not revoked) | contract wording completion (§1.7.1) | §12.2's dual-track priority rule | all four questions no; a wrong call here changes **wording**, not a rule entity (both readings satisfy P1's safety requirement) |

## 2. Execution pipeline

| Step | Role and model | Result |
|---|---|---|
| ① pre-analysis | `deep-reasoner` / `deepseek-v4-pro` | plan book: the 42-item disposition table, the A3 and H4 landing points, the boundary self-check, drafts OB-28 – OB-33, and the ⑤⑥⑦ input-package points |
| ② implementation and ④ integration | master (⑧ not delegated in this slice) | four batches: 附录 B → §12 pilot closure → text sync → smaller new items; then the 19 C1 items completed per user ruling |
| C1 mechanical re-check | master (0 paid calls) | reconciliation of the 34 claimed changes: `items=34 missing=0 result=C1 CLOSED` |
| M4 mechanical re-check | master (0 paid calls) | the shipped G6-4 criteria score 0 hits on "NN 条" / "NN 项", the positive control fires, the real C.0h block scores 0 ⇒ not reproducible |
| Pre-commit review (v1.12) | third-party independent review | `DIRECTIVE REVIEW: FINDINGS`: 1 Critical (C1) + 4 High + 13 Medium + 4 Low |
| Post-disposition re-verification round | third-party independent review | C1 re-verified 19/19 in place; this round 1 High + 5 Medium + 6 Low, disposed per user ruling (H1 and M1 – M5 landed, L1 – L6 deferred to the next slice) |
| ⑤ pre-review | `deep-reasoner` / `deepseek-v4-pro` | `PRE-REVIEW: PASS WITH FIXES`: 2 Medium + 5 Low + 2 Info, **all fixed** (§6.2's `.tmp/`, §12.4's OB-27 disposition column, §9.2's section-number precision, three wordings and two references inside this report, and the freedom list's F-14 typo and wording alignment) |
| ⑥ independent scan | `independent-reviewer` / `mai-code-1.1-flash` | `INDEPENDENT SCAN: PASS`, Sections A–E complete, **0 findings** |
| Fourth re-review (R6) | third-party independent review | `DIRECTIVE REVIEW: FINDINGS`: 1 High + 4 Medium + 5 Low; N1 and N2 were ruled directly by the user, N3, N4, N6 and N10 are ledger- and wording-class fixes that landed, N5 was measured as not reproducible, and N7, N8 and N9 are registered as OB-34 to OB-36 |
| Fifth re-review (R7) | third-party independent review | `DIRECTIVE REVIEW: FINDINGS`: 1 High + 4 Medium + 4 Low; P1, P2 and P3 landed per user ruling (carrying rule-semantics changes), P4 and P5 are registered as OB-37 and OB-38, and P6 to P9 are ledger- and wording-class fixes that landed |
| ⑦'s third round (narrowed input, user-authorized) | `independent-reviewer` / `gpt-5.3-codex` | `RE-REVIEW: PASS`, **no Medium+ this round**; all five groups closed (P1's dual-track priority rule / P2's deferral companions / P3's judge and anchor / P4 and P5's registration truthfulness / the report and list's truthfulness) |

## 3. Review summary and item-by-item disposition

### 3.0 Reading and reconciliation

| Round | Report (local-only) | Critical | High | Medium | Low | Subtotal |
|---|---|---|---|---|---|---|
| Round 2 | `tmp/DELIVERY_DIRECTIVE_v111_REVIEW.md` | 0 | 5 | 13 | 10 | 28 |
| Round 3 | `tmp/DELIVERY_DIRECTIVE_v111_REVIEW_R3.md` | 0 | 2 | 9 | 3 | 14 |
| Pre-commit | `tmp/DELIVERY_DIRECTIVE_v112_REVIEW.md` | 1 | 4 | 13 | 4 | 22 |
| **Total (the C.0h 64-item target)** | — | **1** | **11** | **35** | **17** | **64** |
| Post-disposition re-verification | `tmp/DELIVERY_DIRECTIVE_v112_REVIEW_R5.md` | 0 | 1 | 5 | 6 | 12 |
| Fourth re-review (R6) | `tmp/DELIVERY_DIRECTIVE_v112_REVIEW_R6.md` | 0 | 1 | 4 | 5 | 10 |
| Fifth re-review (R7) | `tmp/DELIVERY_DIRECTIVE_v112_REVIEW_R7.md` | 0 | 1 | 4 | 4 | 9 |

**Counting reading**: accumulated round by round from the reports, with items that recur across rounds counted again (this reading is written into C.0h). **The de-duplicated item-by-item disposition is the set below**: series A (14) + series H (5) + series M (13) + series L (10) + series V (22) = **64**, where series A expands round 3, series H/M/L expand round 2, and series V expands the pre-commit review. **The 12 findings of the re-verification round** are not part of that 64 (its H1 corrected C.0h's own numbers about the first three rounds); their disposition is in §3.6; the **fourth re-review (R6)'s 10 findings** and the **fifth re-review (R7)'s 9 findings** likewise do not count into that 64 and are disposed in §3.8 and §3.9.

### 3.1 Series A (round 3: 2 High + 9 Medium + 3 Low, expanded item by item)

| ID | Disposition | Finding summary | Landing point and fix |
|---|---|---|---|
| **A1** | fixed | CONTRACTS' gate and this directive's G1–G7 are the same word for different things | §0.1 gained a disambiguation row: "gate criteria" here means CONTRACTS' **runtime gate**, while the §6.3 G1–G7 are **execution discipline** |
| **A2** | fixed | The B.7.1 kickoff template's first line is necessarily judged `[ROUTINE]`; the counting unit is undefined | first line became `切片：[切片 ID] [描述]`; §1.6 gained the 20-code-point counting unit and the `[id]` character set |
| **A3** | fixed (user ruling) | the daily tier / production-code boundary was undefined | §1.6 gained "production-code (`.go`) semantic changes require `[SLICE]`"; §1.7.2 gained "the §1.6 boundary outranks §1.7.1 freedom" |
| **A4** | fixed | v1.7 / v1.9 / v1.11 should have been MAJOR under the current definition | §0.2 gained a ruling note: the fact is admitted, historical version numbers are not rewritten, strict enforcement from v2.0 |
| **A5** | fixed | the third of §0.2's three hard requirements had no landing point | §0.3 gained the "affected in-flight slices" column; §0.2's item ③ gained the filling rule; historical rows read `历史缺失` |
| **A6** | fixed | §12.1 and §1.6's "may suggest, never self-escalate" were not joined up | §12.1 gained the `[SLICE]` kickoff handling (including the duty to record "no pilot-metric declaration") |
| **A7** | fixed | hard rule R2.1 was self-exempted by a note | the §2.2 note became a known-limitation entry and states **this row does not exempt R2.1** |
| **A8** | fixed | `[ROUTINE]` can change the repo yet had no tier rule | the §8.2 table gained a `[ROUTINE]` row; §1.6 keeps a single reference |
| **A9** | fixed | the §1.7.3 two-branch test had gaps | the test became three-branch (wrong ⇒ fix code means freedom; wrong ⇒ fix contract / direction / scope / budget / authorization means escalation boundary), and anything uncovered stops on ambiguity |
| **A10** | fixed (the only criteria exemption the boundary allows) | the stopped state is by design a three-way inconsistency | the §6.3 G1 row gained the stopped-state exemption, with completeness owned by G7-d |
| **A11** | fixed (merges M13 and L7) | the tier / probe budget / freedom list had no home; the list's producer, timing and path contradicted each other | the §1.2 table gained three fields; the §1.7.3 traceability paragraph became a fixed-path orchestration artifact |
| **A12** | partly fixed | the multi-trigger row left a misreading path | §1.6 hardened to "stop and require the user to resend", and **must not continue as `[ROUTINE]`** |
| **A13** | fixed (merged with M10) | the slice ledger / ADR / gates / agents had no authority row | see the M10 entry |
| **A14** | fixed | §1.6's `[ROUTINE]` and §1.1's "not applicable" overlapped | §1.6 gained a decidable dividing line: **whether anything lands on disk** |

### 3.2 Series H (round 2 High)

| ID | Disposition | Finding summary | Landing point and fix |
|---|---|---|---|
| **H1** | fixed (the G8 proposal ⇒ OB-29) | the call string had no single authority and the UI display name was easily written as one | the §2.2 header gained modelId and UI display name columns; seven carriers in 附录 B became modelId; a note bans display names from task packages and templates |
| **H2** | fixed | verdict lines were not a closed set, `RE_` was a truncated token, and ⑥/⑦ were indistinguishable without prefixes | §4.1 gained the closed-set table (including `DIRECTIVE REVIEW`), `RE_` was deleted, bare tokens now count as format deviation; B.6 now references §4.1 |
| **H3** | fixed | CONTRACTS / schema / fixtures counted as "non-code documents" ⇒ contract-first had no legal executor | the §4.2 prohibition list excludes the three; §5.1 gained the ②a contract / ②b code split in the implementer's hands |
| **H4** | fixed (user ruling) | the pilot denominator drifted and was not recomputable; no re-evaluation clause and no decision record | §12.1 fixed the denominator (our findings + reviewer findings + user findings, per item, same root cause merged into 1); §12.2 gained the re-evaluation clause and the DR-FIX decision record |
| **H5** | fixed (the ⑧c exemption the boundary allows) | ⑧c's whole-tree per-file encoding check is unsatisfiable for frozen evidence packs | §5.1 gained the frozen-evidence-pack exemption (self-consistency check instead; new packs still pass the encoding gate); §5.3 echoes it; OB-27 turned to landed |

### 3.3 Series M (round 2 Medium)

| ID | Disposition | Finding summary | Landing point and fix |
|---|---|---|---|
| **M1** | fixed | §5.2 and §7.7 had no arbiter for blind-review triggering | the §5.2 high-risk row dropped "+blind review"; the note below the matrix says triggering is judged by §7.7, so the two no longer double-count |
| **M2** | fixed | B.0 and §5.2 still said G1–G5 | both now say G1–G7; historical rows are kept per §11.1 |
| **M3** | fixed | the acceptance evidence cited a criterion that does not exist | C.0g was rewritten to the G3(a)+(b) reading with a correction note; C.0d gained a historical-reading note (L9) |
| **M4** | fixed + OB-28 added | stale ledger rows; G6 only judges changed rows so it cannot see them | OB-4 and OB-2/3/5/6 were closed row by row with the confirming round named; the §12.5 gap row turned to landed; OB-28 registered |
| **M5** | stopped to ask the user (recommending option B) | the master's ownership of governance-slice validation reports was not closed | the user adopted option B ⇒ §3.2 gained the orchestration-artifact exception (it does not exempt ⑤⑥⑦ review) |
| **M6** | fixed (a definitional addition) | the focused gate had no command and no selection rule | §6.1 gained the focused-gate definition, requiring the "steps" field of a slice definition to name the command |
| **M7** | partly fixed (the threshold system ⇒ OB-31) | rework counting and thresholds formed parallel ladders | §9.2's third item gained the counting reading and shares one counter with §5.3; the threshold rebuild is registered as OB-31 |
| **M8** | fixed | the fallback "stop after 2" and ⑥'s ceiling of 3 pointed opposite ways | §9.1 gained the ⑥-layer exception: only reaching the §7.2 ceiling starts the failure streak |
| **M9** | observation item OB-30 | the frontmatter five-key criterion is a new criterion; the 附录 D template is not built here | both are registered as OB-30 for the criteria-revision slice |
| **M10** | fixed (merged with A13) | the ledger paths referenced by G1 and G5 were undefined | the §0.1 table gained the ledger / ADR / gates / agents rows and the full requirements-document path |
| **M11** | fixed | ⑧b could introduce un-reviewed document diffs after ⑦ | the §5.3 ⑧a/⑧b row gained the ⑦ re-review clause bound to the existing quota |
| **M12** | fixed | the minimum reading set missed the blacklist and the anchoring isolation sections | B.7.4's "all" row gained §2.4 and §3.1–§3.3; the "simple" row gained §7.1 and §7.3 |
| **M13** | fixed (merged into A11) | the freedom list's producer, step and path disagreed | see the A11 entry |

### 3.4 Series L (round 2 Low)

| ID | Disposition | Finding summary | Landing point and fix |
|---|---|---|---|
| **L1** | fixed | ⑦'s first-round exemption list said ③ where it should say ⑤ | §7.1 now reads "first round attaches no ⑤/⑥ checklist" |
| **L2** | fixed | `.tmp/` has no gitignore guarantee | §6.2 now points at the repo-root `tmp/` and at §1.5 |
| **L3** | fixed | §7.2 numbering skipped | §7.2 gained the 7.2.1 trigger and 7.2.2 output-contract headings |
| **L4** | fixed | the language note described a timing that had passed | the header note now says the EN mirror is updated once the CN side freezes |
| **L5** | fixed | the focused gate's executor was unclear | §5.3 now reads "④ reruns the focused gate (§6.1)" |
| **L6** | fixed | size, type and high-risk had no mapping | the §5.2 table gained the size → type mapping and states that high-risk is an independent dimension |
| **L7** | fixed (merged into A11) | the probe-budget field had no home in acceptance | see the A11 entry |
| **L8** | fixed | role versioning pointed at §2.6 where it is §2.5 | R2.4 now points at §2.5 |
| **L9** | fixed (note added) | C.0d's acceptance evidence used the old G3 reading | C.0d gained a historical-reading note; its body is unchanged per §11.1 |
| **L10** | fixed | the blacklist's neighbouring-version boundary was unstated | the §2.4 blacklist row now covers the whole GPT-5.x line (only `gpt-5.3-codex` is current) |

### 3.5 Series V (pre-commit review: 1 Critical + 4 High + 13 Medium + 4 Low)

| ID | Disposition | Finding summary | Landing point and fix |
|---|---|---|---|
| **C1** | **disposed** (completed per user ruling) | the §0.3 v1.12 row claimed 34 changes of which 19 were absent from the body (56%) | the 19 were completed rather than the log rewritten, per user ruling; the mechanical reconciliation shows 0 missing of 34 (see §4 and §4.1) |
| **V-H1** | fixed | §4.1's closed set directly conflicted with the B.6 template (bare tokens, no final/blind prefixes, open-ended wording) | B.6 now references the §4.1 closed set and bans bare tokens (landed with the C1 completion) |
| **V-H2** | fixed | §1.6 referenced a "§8.2 `[ROUTINE]` row" that did not exist | §8.2 gained the row; §1.6 keeps a single reference (landed with the C1 completion) |
| **V-H3** | fixed | the call string was half-migrated: §2.2 done, seven 附录 B templates not | seven carriers in 附录 B.1 – B.6 became modelId; R2.2 gained the display-name constraint (landed with the C1 completion) |
| **V-H4** | fixed | the pilot closure was still open while the log claimed it fixed | §12.1 fixed denominator plus reading note; §12.2 re-evaluation clause and DR-FIX decision record (landed with the C1 completion) |
| **V-M1** | fixed | §6.1 still had no focused-gate definition | §6.1 gained the definition row |
| **V-M2** | fixed | the §6.3 G1 stopped-state exemption was missing | the G1 row gained the exemption (same as A10) |
| **V-M3** | fixed | stale OB rows, OB-28 – 31 unregistered, the §12.5 gap row not turned | rows closed one by one, OB-28 – OB-33 registered, §12.5 turned to landed |
| **V-M4** | fixed | B.0's wrap-up list still said G1–G5 | it now says G1–G7 |
| **V-M5** | fixed | the ③ legacy identifier and the §7.2 numbering skip | §7.1 now says ⑤; §7.2 gained 7.2.1 and 7.2.2 |
| **V-M6** | fixed | the §9.1 ⑥ exception and the §9.2 counting object were missing | §9.1 and §9.2's third item each gained one row (wording level; the threshold rebuild is OB-31) |
| **V-M7** | fixed | R2.4 still pointed at §2.6 | it now points at §2.5 |
| **V-M8** | fixed | the field table has three new fields yet two places still said "nine fields" | both now say "field table" (no number, so it cannot go stale again) |
| **V-M9** | fixed | the B.7.1 kickoff template's first line fails the trigger-word grammar | first line became `切片：[切片 ID] [描述]` (same as A2) |
| **V-M10** | fixed | the B.7.4 minimum reading set was not updated | same as the M12 entry |
| **V-M11** | fixed | C.0d and C.0g acceptance evidence used the old G3 reading | C.0g rewritten with a correction note; C.0d noted (same as M3 and L9) |
| **V-M12** | fixed | ⑦'s enum contains `PARTIAL` with no disposition defined | §4.1 gained the `PARTIAL` disposition: a subclass of `FINDINGS`, handled by "eliminate, then re-review" |
| **V-M13** | fixed | the §0.3 v1.12 log row's last column said "none", contradicting §0.2's reading | it now names the successor slices and the impact nature (needs write-back) |
| **V-L1** | fixed | size, type and tier had no mapping | the §5.2 mapping sentence was added (same as L6) |
| **V-L2** | fixed | the freedom list sits outside G7's scan scope | §1.7.3 gained the coverage-boundary sentence, plainly stating there is no mechanical protection |
| **V-L3** | fixed | §4.1's twelve combinations did not say which are legal | the closed-set table names the legal values per prefix; `PARTIAL` follows V-M12 |
| **V-L4** | fixed | UI display names were hard-coded in the table | the §2.2 header and note now mark them "example, follow the UI" |

### 3.6 R5 re-verification round (post-disposition re-verification: 1 High + 5 Medium + 6 Low)

| ID | Disposition | Finding summary | Landing point and fix |
|---|---|---|---|
| **R5-H1** | fixed (user ruled: change the numbers) | C.0h's "three rounds running" disagreed with its own severity breakdown, and its "42 items" covered only the first two rounds | C.0h now says **64 items** with the counting reading and the breakdown; the §0.3 v1.12 log row and the C.0h heading both moved from "two rounds" to "three rounds" |
| **R5-M1** | fixed (reading note added) | the §12.1 1/4 baseline and the fixed denominator are not the same reading | §12.1 gained a note after the baseline table: the two cannot be compared directly, and the target and the measured value are compared on the same reading |
| **R5-M2** | fixed (user ruling) | the re-evaluation point was vague and the pilot trigger chain had a never-executes gap | the re-evaluation point became whichever of `GATES-G1G5` / `GATES-TIGHTEN` completes first; §12.1 gained "the next slice must restart as `[PILOT]`, otherwise the pilot is never executed" |
| **R5-M3** | fixed (applicability note added) | C.0h's tier / probe budget / list path had no landing point | C.0h gained the field-table applicability note (one sentence, no manufactured fields) |
| **R5-M4** | **measured as not reproducible** | C.0h's "42 items" / "19 items" were suspected to be G6-4 aggregate-count shapes needing whitelist entries | the shipped G6-4 criteria's four patterns do not include the "NN 条" / "NN 项" shapes: negative set 0 hits, positive control fires, the real C.0h block 0 hits ⇒ no whitelist entry was added (method in §4.2) |
| **R5-M5** | disposed (priority raised) | G6-5 (the §X.Y reference set must be a subset of the real heading set) was not implemented | OB-32 moved from "watching" to "mandatory for the next slice" and batches with `GATES-G1G5` and `GATES-TIGHTEN`; G6-5 is not implemented here; the three-layer consistency architecture is registered as OB-33 |
| **R5-L1 – L6** | not fixed in this slice | six Low items (legal-combination notes, `_EN` wording, list verifiability, and similar) | deferred to the next slice per user ruling; recorded in §8 |

### 3.7 Results of this slice's review chain (⑤, ⑥ and ⑦)

| Step | Verdict line | Disposition |
|---|---|---|
| ⑤ pre-review | `PRE-REVIEW: PASS WITH FIXES` | 2 Medium (§6.2's `.tmp/` not landed and §12.4's OB-27 still marked to-be-discussed) + 5 Low + 2 Info: **all fixed** (directive CN and EN in step, this report and the freedom list in step); both gate scopes and the mirror re-check stay green afterwards |
| ⑥ independent scan | `INDEPENDENT SCAN: PASS` | 0 findings; Sections A–E complete (Section D gives the falsifiability audit, Section E gives the invalidation conditions) |
| ⑦ final review (round 1, no checklist) | `FINAL REVIEW: PASS WITH FIXES` | 1 Medium: this report's §3.0 said `处置见 §3.7` where it should say §3.6 (a CN-side section-number error; the `_EN` side was already right); **fixed** (CN now says §3.6); no Critical / High / Low; all of F-1 to F-14 ruled **`in-bounds`** |
| ⑦ re-review | `RE-REVIEW: PASS` | Medium-1 closure confirmed; four counterexample attempts against the new change were all ruled out; **no Medium+ this round** |
| ⑦'s third round (narrowed input) | `RE-REVIEW: PASS` | the user-authorized exception round: it judged only P1 / P2 / P3's rule changes and P4 / P5's registration; **all five groups closed (P1 / P2 / P3 / P4 and P5 / report and list truthfulness), with no Medium+ this round** (P1 removes the dual-track coexistence without revoking the re-evaluation object or point; P2's closure is complete; P3's judge and anchor are executable and non-subjective; OB-37 and OB-38 are registered truthfully) |

### 3.8 Disposition of the fourth re-review (R6), N1 – N10

> **Reading**: C.0h's **64 findings** are the first three rounds' accumulation; the R5 re-verification round's **12 findings** and the R6 fourth re-review's **10 findings** are recorded separately and **do not count into the 64 again** (the user made this explicit on 2026-09-23). R5's item-by-item disposition is in §3.6; this table is R6's **N1 – N10**.

| ID | Disposition | Finding summary | Landing point and fix |
|---|---|---|---|
| **N1** | fixed (user ruling) | the re-evaluation clause defined only "when" to re-evaluate, never "what"; the two readings have opposite consequences (reading A forces a suspension of this directive at the re-evaluation point) | §12.2 gained a "**re-evaluation object**" clause: the object is **DR-FIX itself**, recomputed by "**how many findings the extended criteria could have caught**" (`G6` / `G7`, had they existed while DR-FIX ran, could have caught the 3 findings missed then ⇒ those count as **findings by us**); the §12.2 ruling record and OB-14's disposition cell in §12.4 follow suit |
| **N2** | fixed (user ruling) | "the first slice is the pilot" directly conflicted with "that slice does not count towards the pilot evaluation"; the subject / object / manner of the "restart" were undefined | §12.1 now reads "**the first slice should be the pilot**; if it starts as `[SLICE]`, the **deferral clause** applies"; the clause names subject = the user, object = the pilot-metric declaration, manner = the next slice starts as `[PILOT]`, and requires **"pilot deferred" to be registered in the ledger**; C.0h's note registers this slice's deferral |
| **N3** | fixed (ledger-metadata class) | C.0h's "64 findings" was labelled with two mutually exclusive readings ("with recurrences" and "de-duplicated"); the number covered only three rounds | C.0h's acceptance-evidence row now reads "the **per-item disposition table** (item count per this report)"; C.0h's "why it is needed" row gained the disposition accounts of R5's 12 findings and R6's 10 findings |
| **N4** | fixed (ledger-metadata class) | the §0.3 v1.12 log row did not record the later re-reviews' disposition (the second in-slice empirical case of the OB-32 defect class) | the §0.3 row tail gained "**Follow-up note (fourth round, R5 re-review)**" and "**Follow-up note (fifth round, R6 re-review)**"; OB-32's "mandatory for the next slice" status is stated in §5 |
| **N5** | **measured as not reproducible** | it claimed C.0h's "64 findings", "1 Critical + 11 High + 35 Medium + 17 Low", "19 items" and "34 items" fall in G6-4's aggregate-count shapes and need whitelist entries | measured `gate.js --check=G6-4 --scope=tree` = `PASS` (`scanned 945`, 0 hits); the shipped probe reproduces the same verdict ⇒ **no whitelist entry added** (the third OB-6 precedent); its long-term advice (do not write bare counts that expire) is absorbed by the criterion-style wording of N3 |
| **N6** | fixed (wording class) | C.0h's boundary sentence contradicted itself ("no change to the G1–G7 criteria body" immediately followed by two criteria exemptions) | C.0h's boundary row and this report's §1 "not touched" row now read "no change to the existing judgement logic (only two exemption clauses are added)" (see F-16 in §1.1) |
| **N7** | left to the next slice | §7.2.2's heading overflows its content (cost ceiling and failure handling hang under the output contract) | registered as **OB-34** (disposed in the next slice) |
| **N8** | left to the next slice | §1.6 references "§7.2 ③" while the three conditions now live in §7.2.1 - the reference is imprecise | registered as **OB-35** (disposed in the next slice) |
| **N9** | left to the next slice | OB-14's disposition cell substitutes a pointer for the four elements required by §12.2's re-evaluation clause ① | registered as **OB-36** (disposed in the next slice) |
| **N10** | fixed (wording class) | §1.6's boundary ② code item is already a subset of the production-code boundary, so listing both reads as two independent boundaries | §1.6's boundary ② now notes "**the item below already covers this, and it is kept for contrast only**" |

**The second remediation after ⑧b and its relation to ⑦ (user ruling 2026-09-23, landed verbatim)**: every remediation here comes from an **explicit user ruling** (N1 and N2) or from **ledger-metadata / wording-class** fixes (N3, N4, N5, N6 and N10), with **no self-initiated rule-semantics change**, so under the **OB-11 exception ⑦ is not re-run**; the package stands as **⑦'s historical input** (folded in at the next revision). **The same ruling also covers** the ⑧b factual record already explained in §3.7 (the `⬜` → `✅` flip, the ⑨ measurements and the ⑦ verdict lines): under the OB-11 exception standard (ledger-metadata class, no rule entity, closure verifiable literally, no new semantics) it **does not trigger a ⑦ re-run**; ⑦'s quota (one final review plus one re-review) is not exceeded.

### 3.9 Disposition of the fifth re-review (R7), P1 – P9

> **Reading**: C.0h's **64 findings** are the first three rounds' accumulation; the R5 **12 findings**, the R6 **10 findings** and the R7 **9 findings** are each recorded separately and do not count into the 64 again (R5 is in §3.6, R6 in §3.8; this table is R7's **P1 – P9**).

| ID | Disposition | Finding summary | Landing point and fix |
|---|---|---|---|
| **P1** | fixed (user ruling, carries rule semantics) | the deferred pilot and the DR-FIX re-evaluation form two independent tracks with no priority rule ⇒ "the new pilot passed" and "this directive is forcibly suspended" can hold at once | §12.2 gained the "**dual-track priority rule**": the **new pilot**'s result **outranks** DR-FIX's historical re-evaluation (new pilot meets every metric ⇒ **counts as having passed the pilot** and DR-FIX's re-evaluation **becomes a record item**; new pilot **misses** ⇒ the rollback path runs and the re-evaluation **merges into that same judgement**); a re-evaluation point earlier than the new pilot's wrap-up is **merged** into that judgement (the point is not cancelled) |
| **P2** | fixed (user ruling, carries rule semantics) | the deferral clause had no qualification filter and said nothing about a second deferral | §12.1 gained the "**two companion clauses to the deferral**": ① the deferral target must **meet the pilot-selection criteria**, otherwise the deferral **continues** with the **deferral count** registered; ② **two or more deferrals** ⇒ the master stops per §1.7.2 and **escalates for a user ruling** |
| **P3** | fixed (user ruling, carries rule semantics) | the re-evaluation test was not mechanically checkable, the judge was undefined, and the "3 findings" count was wrong | §12.2 gained the "**judge and recomputable anchor**": the judge is **⑦** (deciding inside the re-evaluation input package and recorded in the report); the anchor is running `gate.js --scope=range:<DR-FIX commit>^..<DR-FIX commit>` over G6 / G7 for DR-FIX's commit, where **a red verdict counts as "catchable"**; the count became the **classifiable items (OB-14 records 2)** |
| **P4** | fixed (registration; this round's meta-level finding) | "no ⑦ re-run under the OB-11 exception" was misapplied: N1 and N2 carry **rule entities and new semantics** and fail the OB-11 exception standard | registered as **OB-37** (with the four elements: who ruled / date / reason / basis) plus the **generalizable meta-rule**: a direct user ruling may break the OB-11 exception, but **must be recorded independently** and **labelled truthfully**; §0.3's ⑪ gained a **correction note**; **later remediation in this slice must not use the OB-11 exception again** |
| **P5** | fixed (registration) | G6-1's marker vocabulary and §12.4's actual status vocabulary barely intersect | registered as **OB-38** (batched with `GATES-TIGHTEN` / `GATES-G1G5`); this slice **does not widen G6-1** (that is a new criterion, beyond the boundary) |
| **P6** | fixed (ledger-metadata class) | C.0h's goal row did not cover the later re-reviews | C.0h's goal row now reads "**64 findings** plus **22 findings from the two later re-reviews** (R5's 12 and R6's 10; the item count is per this report's disposition table)" |
| **P7** | fixed (wording class) | OB-14's disposition cell said "this slice (DR-FIX)", easily misread as the current slice | it now reads "**DR-FIX (the slice OB-14 belongs to)**" with the judge / anchor / count added |
| **P8** | fixed (wording class) | §1.6's boundary ② kept a code item whose strength differed from the production-code boundary | it now reads "**the production-code boundary governs this item, which is listed for contrast only**" |
| **P9** | fixed (wording class) | §12.5's "criterion self-reference risk" row sat among the "landed" rows | its timing cell now reads "the same root as OB-28; **candidate slice = `GATES-TIGHTEN`** (user ruling)" |

**The OB-11 exception lapses in this slice (user ruling 2026-09-23, landed verbatim)**: P1, P2 and P3 **all carry rule-semantics changes ⇒ they fail the OB-11 exception**, so this slice authorized **⑦'s third round** (narrowed input: only P1 / P2 / P3's rule changes and whether P4 / P5's registration is truthful); **later remediation in this slice must not use the OB-11 exception again**. ⑦'s quota (one final review plus one re-review) is already spent ⇒ that third round is a **user-ruled** exception (same precedent as `GATES-EXT`) and **does not generalise** to other slices. That third round returned `RE-REVIEW: PASS` (**no Medium+**): all five groups closed (P1 / P2 / P3 / P4 and P5 / report and list truthfulness).

**Pilot evaluation attribution (R7's closing judgement, adopted)**: this slice **has deferred** per §12.1 ⇒ the pilot evaluation moves to the next slice and **this slice's review findings count towards no pilot metric**; however, this slice's defect self-capture rate is still recorded as an **instance sample under §12.1's fixed denominator reading**, for future evaluation.

## 4. The C1 event (Critical) and two mechanical re-checks

### 4.1 C1: 34 items claimed, 19 not landed

| Field | Content |
|---|---|
| Fact | The §0.3 v1.12 log row claimed 34 changes, of which **19 items** did not exist in the body (the reviewer reconciled them one by one) |
| Mechanism | ⑦'s review scope is given by the log ⇒ **areas claimed but never landed are skipped systematically**; the reviewer's words: the directive itself supplies the miss mechanism |
| User ruling | The 19 **must be completed**; rewriting the log to match reality was refused (it is cover-up and would concede that the audit record is untrustworthy) |
| Completion | the 19 were landed in four layers: ① 附录 B (B.0 / B.6 / B.7.1 / B.7.4); ② §12 pilot closure (§12.1 / §12.2 / §12.4 / §12.5); ③ text sync (R2.2 / R2.4 / §6.1 / §6.3 G1 / §7.1 / §7.2 numbering / §9.1 / §9.2's third item / the §8.2 `[ROUTINE]` row); ④ remaining items (the C.0d and C.0g correction notes, and the new C.0h) |
| Mechanical re-check | `node tmp/v112/c1-recheck.js` → `items=34 missing=0 result=C1 CLOSED` (each item probed on both the CN and the EN side) |
| Re-verification | the third-party re-verification round confirmed C1 item by item: 19/19 in place and of acceptable quality ⇒ C1 dropped from Critical to disposed, **no acknowledgment needed and no change made** |
| Institutional meaning (empirical evidence for G6-5) | once the log becomes the review's source of truth, a "claim" is automatically exempted; this slice therefore promotes **OB-32** (the §X.Y reference set must be a subset of the real heading set) from an observation item to **mandatory for the next slice** |
| Positive outcome | this slice's review verdict line `DIRECTIVE REVIEW: FINDINGS` is **legal for the first time** - the `DIRECTIVE REVIEW` role added in v1.12's §4.1 finally gives a directive review a self-consistent verdict line |

### 4.2 The R5-M4 measurement (verdict: not reproducible, no whitelist entry added)

| Step | Result |
|---|---|
| Criteria taken | the shipped module `tools/gates/lib/criteria/g6.js` G6-4 sub-check was loaded directly (not re-implemented) |
| Negative set (the accused shapes: `64 条` / `19 项` / the severity-breakdown literals) | 6 input rows fed in → `PASS`, `no expiring literal on added lines (scanned 6)` |
| Positive control (**引文** shapes: `14 个文件` / `42 行` / `164/80 行` / `共 42 个文件`) | 4 input rows fed in → `FAIL`, all four hits listed ⇒ the probe itself works |
| The real C.0h block (line by line from the live file) | 17 input rows fed in → `PASS`, 0 hits |
| Root cause | G6-4's four patterns are "ratio plus lines", "count plus lines", "count plus files", "total plus count"; **the "NN 条" and "NN 项" shapes are not among them** |
| Verdict | no `tools/gates/g6-whitelist.txt` entry added (adding one is a no-op and would mask a shape the criterion does not ban); recorded per the OB-6 precedent, "measurement overturns the reviewer's fix direction" |
| Probe location | `tmp/v112/m4-g6-4-probe.js` (`local-only`) |
| R6 re-run (N5) | `gate.js --check=G6-4 --scope=tree` = `PASS` (`scanned 1003`, 0 hits); the shipped probe reproduces the same verdict ⇒ the same conclusion holds for the third time, and **no whitelist entry is added** |
| R7 re-run (N5's evidence request) | `gate.js --check=G6-4 --scope=tree` = `PASS` (`scanned 1079`, 0 hits); the probe reproduces the same verdict ⇒ the same conclusion holds for the fourth time; **the evidence is on disk** (`tmp/v112/g6-4-check.txt`, `local-only`) |

## 5. User-ruling provenance

| Ruling | Content | Landing point | Class |
|---|---|---|---|
| A3 | the daily tier may not change `.go` semantics; the §1.6 boundary outranks §1.7.1 freedom | §1.6 and §1.7.2 | user ruling |
| H4 | DR-FIX fell short of the pilot target ⇒ continue executing; add the re-evaluation clause and the fixed-denominator definition | §12.1 and §12.2 | user ruling |
| M5 (round 2) | option B adopted: governance-slice validation reports are produced directly by the master (orchestration-artifact exception, no ⑤⑥⑦ exemption) | §3.2 | user ruling |
| R5-M1 | severity numbers corrected to 64 items with the four-band breakdown; "two rounds" became "three rounds" | §0.3 and 附录 C.0h | user ruling (the only honest option) |
| R5-M2 | the re-evaluation point became "whichever completes first"; the `[SLICE]` kickoff handling and "the next slice must restart as `[PILOT]`" were added | §12.1 and §12.2 | user ruling |
| R5-M3 | the field-table applicability note (no manufactured fields) | 附录 C.0h | user ruling |
| R5-M5 | **landed as a (1) recommendation plus user ruling**: OB-32 promoted to "mandatory for the next slice" (batching with `GATES-G1G5` and `GATES-TIGHTEN`), with G6-5 not implemented here | §12.4 | (1) recommendation + user ruling |
| R5-M4 | no whitelist change; the measured finding is recorded in this report | this report §4.2 | (1) recommendation + master choice (**not** a user ruling) |
| R5-L1 – L6 | not fixed in this slice, deferred | §8 | user ruling |
| The ② and ⑧ execution roles | carried by the master (the budget line lists only ①⑤⑥⑦) | this report §2 | master freedom (see F-13) |
| R6-N1 | the re-evaluation object is DR-FIX itself, recomputed by "how many findings the extended criteria could have caught" (removing the buried mine of "a mandatory suspension at the re-evaluation point") | §12.2's re-evaluation-object clause and ruling record, §12.4's OB-14 | user ruling |
| R6-N2 | the pilot evaluation is the next slice that starts as `[PILOT]`; the first slice defers per the deferral clause and registers "pilot deferred" | §12.1's deferral clause, Appendix C.0h's note | user ruling |
| R6-N3 to N10 | N3, N4, N5, N6 and N10 are ledger-metadata / wording-class fixes; N7, N8 and N9 are left to the next slice; the whole package is closed under the **OB-11 exception with no ⑦ re-run** | §0.3, §12.4 and Appendix C.0h; this report §3.8 | user ruling (whole-package closure reading) |
| R7-P1 to P3 | the dual-track priority rule, the two deferral companions, and the re-evaluation judge / anchor / count reading (**carrying rule-semantics changes**) | §12.1 and §12.2 | user ruling (the OB-11 exception **does not apply**, hence the authorized ⑦ third round) |
| R7-P4 and P5 | OB-37 registered (with its generalizable meta-rule) and OB-38 registered; **later remediation in this slice must not use the OB-11 exception again** | §12.4 and the correction note in §0.3's ⑪ | user ruling |
| R7-P6 to P9 | ledger-metadata / wording-class fixes (C.0h's goal row, OB-14's wording, §1.6 boundary ②'s strength, §12.5's candidate slice) | Appendix C.0h, §12.4, §1.6, §12.5 | user ruling |
| Pilot evaluation attribution | this slice has deferred ⇒ its review findings **count towards no pilot metric**; its self-capture rate is recorded as an **instance sample** | §12.1 and this report §3.9 | R7's judgement adopted (confirmed by the user) |

**Provenance note**: the table above separates "user ruling" from "① recommendation + master choice" item by item, so that master discretion is never written up as a user ruling (or the reverse). **R5-M5 landed as a ① recommendation plus user ruling**, and its evidence chain is this report together with the slice ledger. **R6's N1 and N2 are user rulings too** (N3 to N10 are its closure reading), and the table above names them item by item.

## 6. Gates and native validation

| Item | Result |
|---|---|
| `node tools/gates/gate.js --selftest` | `SELFTEST: PASS 203/203` |
| `--all --scope=tree` | `TOTAL_FAIL=0` |
| The four §6.1 base gates | `gofmt` prints nothing, `go build` and `go vet` exit 0, `go test ./...` reports `ok=14 FAIL=0` |
| Encoding gate (change set) | both directive files and both report files are UTF-8 with BOM plus LF |
| Bilingual positional mirror | the directive pair has 0 per-line mismatches (CN and EN both have 1232 lines); changed-line counts are symmetric (numstat 175/84 against 175/84) |
| Link and reference resolvability | every section reference added or changed (§0.1 / §0.3 / §1.2 / §1.6 / §1.7 / §2.2 / §2.4 / §2.5 / §3.2 / §4.1 / §4.2 / §5.1 / §5.2 / §5.3 / §6.1 / §6.2 / §6.3 / §7.1 / §7.2 / §7.5 / §7.6 / §7.7 / §8.2 / §9.1 / §9.2 / §12.1 / §12.2 / §12.4 / §12.5 plus 附录 B and 附录 C) exists in the directive |
| The degenerate form for a documentation slice | per §5.2, ⑨ degenerates to "encoding gate + G1 to G7 self-check + link and reference resolvability check"; the measurement here: 1097 §X.Y references inside the changed-line scope, **0 unresolved** (`tmp/v112/refs-check.txt`; the whole file carries 18 further external citations such as v3.1, all on inherited rows and excluded) |

## 7. ⑧c wrap-up cleanliness record

| Field | Result |
|---|---|
| Time and executor | 2026-09-23; the master (0 paid calls) |
| Whole-tree diagnostics (by §3.4.1 category) | compiled: `go build ./...` and `go vet ./...` (through the `gate.js` base gates); scripting: `node --check` had nothing to cover (no `.js` change here); marked data: the four files' encoding gate and `gate.js` parsing; plain text others: none; uncovered types: none |
| IDE diagnostics summary | this session has no Problems-panel reading interface, recorded as uncovered (see §8) |
| `gate.js --all --scope=tree` | `TOTAL_FAIL=0` |
| Encoding gate (whole tree, per file) | whole-tree check: 2912 tracked, of which 2823 text and 89 binary, with **3 violations**, all inside frozen evidence packs; handled by the OB-27 temporary exemption and the self-consistency check below. Reading: this slice's 3 new files raise each count by 3 against the v1.11 report |
| Frozen evidence pack self-consistency (exemption condition two) | the `SHA256SUMS.txt` of b2 (41), b3a (128) and b3b (2158) agrees entry by entry with the real files, 0 missing, 0 mismatched; `b1-20260917` has no `SHA256SUMS.txt` ⇒ **not self-checkable** (stated plainly) |
| Temp directories, probe residue and untracked files | `git status --porcelain -uall` shows only this slice's changes and no untracked residue; `tmp/v112/**` and this round's four review reports sit outside version control |
| Verdict line | `WRAP-UP CLEANUP: FINDINGS` (the whole-tree per-file check still finds 3 pre-existing violations **inside frozen evidence packs**; fixing them would break pack integrity ⇒ a §1.7.2 escalation boundary, handled by the OB-27 temporary exemption and the pack self-consistency check, same source and same verdict as the v1.11 report) |

**Milestones**: ① **C1 was closed by completing the work, not by rewriting the log, for the first time** - `items=34 missing=0` made the log's claim true without sacrificing the audit record's trustworthiness; ② **the R5-M4 measurement overturned the reviewer's fix direction** (the second appearance of the OB-6 precedent) - implementability is decided by measurement rather than by a reviewer's impression; ③ **the `DIRECTIVE REVIEW` verdict line is legal for the first time** - the directive-review role added in v1.12 is self-consistent in the first report that uses it.

## 8. Not executed and not covered

| Item | Reason |
|---|---|
| R5's L1 – L6 | the user ruled they are **not fixed in this slice** and are deferred to the next one |
| G6-5 (the §X.Y reference set must be a subset of the real heading set) | not implemented here (it touches `gate.js` plus `selftest` plus tests, a separate capability); promoted to OB-32, "mandatory for the next slice" |
| The three-layer cross-location consistency review architecture | registered as OB-33, to be assessed together with OB-32 |
| A mechanical criterion for the freedom list | that path sits in `docs/validation/evidence/**` (outside G7's scan scope) ⇒ no mechanical protection, only hand review by ⑦ (stated plainly in §1.7.3 and in the list itself) |
| IDE diagnostics reading | this session has no Problems-panel reading interface; recorded as uncovered, with no claim of having been checked |
| The 3 pre-existing encoding violations inside frozen evidence packs | the OB-27 temporary exemption (limited to archived frozen evidence packs carrying MANIFEST / SHA256SUMS); new packs still pass the encoding gate before commit |
| The ③ test role | a documentation slice, so no product-layer sub-agent is dispatched per §5.2 |
| OB-34, OB-35 and OB-36 | the fourth re-review's (R6) N7, N8 and N9, left to the next slice per user ruling (each registered in §3.8) |
| OB-37 and OB-38 | the fifth re-review's (R7) P4 and P5; OB-37's meta-rule **takes effect from this slice on**, and OB-38 is batched with `GATES-TIGHTEN` / `GATES-G1G5` (each registered in §3.9) |
| ⑦'s third round (narrowed input) | authorized by the user on 2026-09-23 (an exception to the spent quota); it judges only P1 / P2 / P3's rule changes and whether P4 / P5's registration is truthful; **a further Medium+ ⇒ stop and escalate, with no automatic round 4** |
| The ⑦ exception for later remediation here | **the OB-11 exception lapses in this slice** (P4): later remediation must not skip ⑦ on OB-11 grounds again; the basis for skipping it during the R6 round (the fourth round of remediation) is corrected by OB-37 to a **direct user ruling** (the old row read "not re-run per user ruling", whose reason is exactly this correction - see §3.8 and §3.9) |

## 9. Cost and metering

| Item | Budget | Actual |
|---|---|---|
| ① pre-analysis | 1 | 1 (`deepseek-v4-pro`) |
| ② implementation, ④ and ⑧ | master | the master (0 paid calls) |
| Pre-commit third-party review | — | 4 rounds (the C1 round, the R5 post-disposition re-verification round, the R6 fourth re-review, and the R7 fifth re-review) |
| ⑤ and ⑥ and ⑦ | see §7.1 to §7.3 and §7.5 | ⑤ once (`deepseek-v4-pro`), ⑥ once (`mai-code-1.1-flash`), ⑦ three times (`gpt-5.3-codex`: one final review, one re-review and the user-authorized third round); the final review and the re-review sit inside the §7.5 quota, and the third round is a user-ruled exception (registered in §3.9) |
| Probe budget | none (documentation slice) | 0 external calls |

## 10. Next steps

- ⑩ stop point: the first ⑦ final review `PASS WITH FIXES` (1 Medium, fixed), the re-review `RE-REVIEW: PASS` and the **user-authorized third round `RE-REVIEW: PASS` (no Medium+)** have all passed; the fourth re-review's (R6) N1 to N10 and the fifth re-review's (R7) P1 to P9 landed per user ruling (§3.8 and §3.9); ⑧c is green ⇒ this slice may be committed; **commit and push require same-turn authorization** (nothing was committed or pushed during this slice).
- Candidate next slices (listed by the user): `GATES-G1G5` and `GATES-TIGHTEN`; the **mandatory item** is OB-32 (the approximate criterion: the §X.Y reference set appearing in §0.3 and the 附录 C ledger blocks must be a subset of the file's real heading set); OB-30, OB-31 and OB-33 are assessed in the same batch.
- Carried over from this slice: R5's L1 – L6; the priority ruling between `*.raw.txt`'s `-text` and ⑧c (pending with OB-28).

## 11. Artifact list

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/validation/directive-v1.12.md	repo
docs/validation/directive-v1.12_EN.md	repo
docs/validation/evidence/directive-v1.12-freedom-list.md	repo
tmp/v112/**	local-only
tmp/DELIVERY_DIRECTIVE_v111_REVIEW.md	local-only
tmp/DELIVERY_DIRECTIVE_v111_REVIEW_R3.md	local-only
tmp/DELIVERY_DIRECTIVE_v112_REVIEW.md	local-only
tmp/DELIVERY_DIRECTIVE_v112_REVIEW_R5.md	local-only
```

## 12. Freeze ruling and CI result (v1.12 frozen)

> **This section is a post-freeze record** (user ruling 2026-09-23, item 4: the directive text is not changed again, and the report may only gain a "freeze ruling / CI result" subsection). It **supersedes** the state statement in §10's stop receipt that "nothing was committed or pushed during this slice" - that statement no longer holds once the actions recorded here are done; the rest of §10 is unaffected.

### 12.1 The freeze ruling (a direct user ruling, recorded independently under OB-37's meta-rule)

| Element | Content |
|---|---|
| Who ruled | the user |
| Date | 2026-09-23 |
| Reason | R7's disposition and ⑦'s third-round verdict have been accepted; the directive still has no written criterion for when a third-party re-review may stop ⇒ further iteration on the same revision buys less than its regression risk |
| Basis | the user's ruling text this round: accept R7's disposition and ⑦'s third-round verdict, send no further third-party re-review, no R8, no ⑦ round 4 |
| Carries rule-semantics changes | **No**: this section only records the ruling and the evidence, and changes none of the directive's clauses, criteria or vocabulary |
| Precedent force | **not a precedent**: the stop condition still has no written criterion, so this stop is a direct user ruling; the formal criterion is defined by the next slice (see §12.3) |

**Landing reading**: ① accept R7's disposition and ⑦'s third-round verdict (it returned `RE-REVIEW: PASS` with no Medium+), with **no further third-party re-review**; ② v1.12 is **frozen**: P1 to P9 landed and ⑦'s third round completed is what freezes it, whatever the verdict; ③ findings arising after R7 all become observation items, numbered from **OB-39** (see §12.3); ④ the next slice is `DIRECTIVE-REVIEW-SATURATION` (the review-convergence criterion: when a third-party re-review and ⑦ round N may stop), and `GATES-G1G5` and `GATES-TIGHTEN` are not folded into it.

### 12.2 Commit and CI result

| Item | Result |
|---|---|
| Commit | `5c8a6a4` (`docs: freeze DELIVERY_DIRECTIVE v1.12 after R7 disposition`, 5 files) |
| Push | `origin` succeeded (`4514afb..5c8a6a4`); **gitee was not pushed** |
| CI run | run 35856548089 (`CI`, push, sha `5c8a6a4`) |
| Verdict | **double green**: `Go windows-latest` and `Go ubuntu-latest` both `success`; the `Gates` step ran `completed/success` on both legs (not skipped); the ubuntu leg's `Race` and `Contract fixtures` passed; three further jobs (candidate build, bootstrap and candidate probe) are `skipped` by workflow condition |
| Write-back reading | this section records only the **first** run after the freeze; to avoid a write-push-run loop, later run results are not written back here and are reported in the session only |

### 12.3 Observation items after R7 (numbered from OB-39, left to the next slice)

| Number | Source | Content | Disposition |
|---|---|---|---|
| **OB-39** | freeze re-check (master) | this report's §11 artifact list names 4 review reports such as `tmp/DELIVERY_DIRECTIVE_v111_REVIEW.md` but omits the fourth and fifth re-review reports (`_REVIEW_R6.md` and `_REVIEW_R7.md`), although both exist on disk | the next slice reconciles and completes the list reading; ledger-truthfulness class, no change to this slice's rule text |

**Reading note**: §4.2's `scanned 1079` and §6's `1097` are readings of the **whole change set of commit `5c8a6a4`**; this section is a record added after that commit, so it neither rewrites those readings nor recomputes them.
