# ProofRail Delivery Execution Directive (DELIVERY_DIRECTIVE)

Version: v1.13. Date: 2026-09-23. Status: **CN body frozen** ( `_EN` mirror in `docs/DELIVERY_DIRECTIVE_EN.md`).

> **Multilingual note**: this file is a strict positional mirror of the CN body, maintained in step with it (updated after CN is frozen, to avoid bilingual rework during review); `docs/DELIVERY_DIRECTIVE.md` is authoritative and wins on any conflict (see §11.1).
> **Encoding**: this document and its `_EN` mirror must stay UTF-8 **with BOM** + **LF**.
> **The origin and inheritance of this directive** are described in [Appendix A](#appendix-a-historical-versions-and-sources-of-experience); see [Appendix C](#appendix-c-migration-mapping) for the mapping from the legacy T027 directive.

---

## 0. Status, Versions, and Changes of the Directive

### 0.1 Authority relationships

| Topic | Authoritative document |
|---|---|
| Product requirements, acceptance criteria | `docs/PRODUCT_REQUIREMENTS{,_EN}.md`; slice-level acceptance ledger in the "slice ledger" row |
| Protocol and semantics (wire / state machine / gate semantics) | `docs/CONTRACTS.md` (+ `_EN`), `schemas/`, contract fixtures |
| Implementation plan and milestones | `docs/DEV_PLAN.md` (+ `_EN`) |
| Evidence and conclusions | `docs/validation/`, `docs/validation/evidence/` |
| **Delivery execution discipline (this directive)** | **`docs/DELIVERY_DIRECTIVE.md` (+ `_EN`)** |
| Coding / line-ending / Git discipline | `docs/CODING_CONVENTIONS.md`, `.github/copilot-instructions.md` |
| Slice ledger and remaining slices | `docs/t027/REMAINING_SLICES{,_EN}.md` |
| Architecture decision records (ADR) | `docs/ADR_REGISTER{,_EN}.md` |
| Gate-criteria implementation and source of truth | `tools/gates/**` (criterion text in §6.3) |
| Role set (carrier definitions) | `.github/agents/*.agent.md` (rules in §2.5 / §2.6) |

**Conflict handling**: execution-discipline conflicts are decided by this directive; protocol-semantics conflicts are decided by CONTRACTS. When the two conflict, **revise the documents first, then execute**; do not pick one over the other mid-execution.

**Category test**: anything involving wire format, state transitions, gate criteria, or evidence fields = **protocol semantics** (CONTRACTS prevails); anything involving process, roles, authorization, document form, or report structure = **execution discipline** (this directive prevails). If a case still cannot be decided, **CONTRACTS prevails** (conservative principle).

**Terminology disambiguation**: the "gate criteria" in this section means the **runtime gates** defined by CONTRACTS (in-step gate / postflight / freeze); the G1–G7 **delivery self-check criteria** in §6.3 belong to **execution discipline** and are governed by this directive.

### 0.2 Versions and changes

- Version numbers are `vMAJOR.MINOR`; `MAJOR` = a change to the role set or pipeline structure; `MINOR` = an added rule.
- Every change must: ① update the version number and date; ② append a change-log entry in §0.3 (including the "Slices in flight (affected)" column); ③ fill that column with: slice ID + effect (blocked / needs write-back / no action), marking historical gaps as "历史缺失" (historically missing).
- **This directive is itself subject to validation** (see §12): it must not claim "verified effective" until the pilot in §12.2 passes.

**Ruling note (2026-09-23, v1.12)**: v1.7 (role-file `model` deletion and identifier disambiguation), v1.9 (§1.6 three-tier startup) and v1.11 (the §5.1 ⑧c stage) would be MAJOR under this section's current definition (role-set / pipeline-structure changes) but were released as MINOR. This note acknowledges that fact and does not retroactively rewrite historical version numbers; **from v2.0 the MAJOR/MINOR definition above is enforced strictly**.

### 0.3 Change log

| Version | Date | Change | Slices in flight (affected) |
|---|---|---|---|
| v1.0 | 2026-09-21 | First version: used the "Flash Master–Subagent Architecture v1.0" document as the skeleton, folding in the accumulated experience of T027 directive v3.1 and BRIEFING v3.1; established the three-layer role model, the serial pipeline, the evidence-persistence obligation, the wrap-up self-check gate, the non-separable-effort-tier constraint, and fallback thresholds. Supersedes `docs/t027/FLASH_OPERATING_DIRECTIVE*.md` (the original text is kept as a historical version with a supersede pointer). **Same-day review: two rounds of same-family review (including 13 P0–P3 findings) + one independent blind review by Codex (first round with no checklist; 3 Critical/High + 5 Medium + 5 coverage gaps), all remediated**; added R1.1/R2.6/R7.1/R9.1, §5.0 startup gate, §6.3 G5 structured input, §7.7 proportional spot checks, and Appendix C.0/C.0b/C.3. | historically missing |
| v1.1 | 2026-09-21 | SW-1/SW-2 execution-phase remediation (all **measured findings**, not presumptions): ① §6.3 G2 adds "excluded-rule quotation lines", G3 becomes "changed-line-count symmetry", and G4b becomes "the difference set ⊖ `tools/gates/cjk-newwords.txt` allow-list must be empty" — all three original criteria were **self-referential defects** (G2 self-scanned 4 false positives; G3 was always false for existing bilingual pairs; G4b was unsatisfiable for any new document — measured on a 51 KB new document whose 7-character difference set was verified character by character to be entirely legitimate new words); ② §9.3 adds the "write-file tools emit CRLF" failure mode (measured by the SW-2 probe); ③ Appendix C adds `SW-3` (role-set commit + encoding normalization, pending authorization) + C.0/C.0b status lines + C.2 execution results; ④ §12.4 registers OB-2/OB-3/OB-4/OB-5. | historically missing |
| v1.2 | 2026-09-21 | ⑦ independent final review (Codex, first round with no checklist) outcome `FINDINGS`: 3 High + 5 Medium, **all verified item by item, all valid and all remediated** — F1 §2.2 still stated the carrier status as "not yet established" while C.0b was already ✅ (mutually exclusive) ⇒ changed to "all seven carriers are ready"; F2 SW-1's trigger condition conflicted with its completion state ⇒ stated the execution-phase rule explicitly (executable once CN is frozen, `_EN` mirror to follow); F3 role files were inconsistent with §2.6 ⇒ **measurement overturns ⑦'s fix direction**: §2.6 now says the `model` key must be **omitted** (writing an empty value makes the carrier fail to load, proven by isolation verification); F4 SW-1's acceptance assertion "every hit" conflicted with "historical drafts only get a header" ⇒ changed to file-level + section-level dual assertions + new C.0a allow-list; F5 §6.3 did not state the scripting status ⇒ added a "scripting status" paragraph (only the allow-list has been persisted); F6 the probe covered only 1 carrier ⇒ changed to a **three-carrier full-coverage probe matrix**; F7 reference assertions searched file names only ⇒ added section-number keyword search; F8 G3 gained "key-field alignment". | historically missing |
| v1.3 | 2026-09-21 | Remediation from ⑦ round-2 re-review + ⑥ independent scan: ① **split the concatenated v1.1/v1.2 change-log table rows** (structural breakage; independently and overlappingly found by ⑥ and ⑦); ② C.0b "does not currently exist" changed to "did not exist **when SW-2 was chartered**" (tense closure); ③ **tightened ⑥'s skip criterion in §7.2** and added a **counter-example** (this slice revised gate criteria ⇒ ③ does not hold ⇒ **re-run ⑥**); ④ **supplied the missing definition of `Section A–E`** (previously referenced in three places — §4.2/§7.2/B.6 — with **no definition**, a broken reference found in self-check); ⑤ added **G5-a**, a frontmatter hard gate for role files (the `model` key must not appear); ⑥ registered OB-7 (independent overlap between ⑥/⑦). | historically missing |
| v1.4 | 2026-09-21 | Remediation from ⑥ round 2 (rescan after remediation): ① **split the concatenated OB-6/OB-7 table rows in §12.4** (**the same class** of defect as the change-log concatenation fixed in v1.3 ⇒ added the **G1-a table-structure integrity** mechanical criterion, and recorded the measured lesson that "not stripping escaped `\|` produces 4 false positives"); ② §7.2 now supplies the **sole authoritative definition** of `Section A–E` and its inheritance source; ③ the header version was raised to v1.4 (the old header v1.3 and the log v1.4 were **out of sync**) + added **G1-b (version-metadata consistency)** (⑥ Section D pointed out that this invariant had no mechanical coverage). | historically missing |
| v1.5 | 2026-09-21 | **SW-3 executed (user-authorized option A)**: ① **11 existing role files normalized to BOM+LF and all committed** (R2.5 wording rewritten per the user's formulation: "the generator toolchain is not committed; `.agent.md` role files, including existing ones, are always committed", removing ambiguity); ② **uniform header note for the 14 role files** ("非生成物（sol-orchestrator 已按 §2.4 禁用）：手工维护；不得由生成器覆盖") + a mechanical criterion for the header note; ③ registered **OB-8** (the default `model` of 8 existing files is blacklisted; not fixed in this slice). | historically missing |
| v1.6 | 2026-09-21 | **Remediation from ⑦ round-3 re-review (user authorized a one-time breach of the §7.5 cap)**: 7 of R1–R8 closed, R7 partially closed, and 4 new findings (3 Medium + 1 Low) all valid and all remediated: ① **blacklist count corrected 9→8** with a family breakdown (Terra 3 / Luna 2 / Gemini 1 / GPT-5.4 1 / GPT-5.6-Sol 1) + registered the residual risk that "a newly added `.agent.md` can bring in a blacklisted default `model`"; ② **G5-a hardened into two criteria** (the original criterion matched only a fixed form and covered only `prfrail-*` ⇒ blank variants and other files escaped detection); ③ **header-note criterion changed to position anchoring** (the first non-empty line after the frontmatter must equal the header-note text; "merely appearing in the document" is no longer accepted); ④ the **header-note blocks of the 3 `prfrail-*` files were moved above the H1**, unifying the position across all 14 files (measured 14/14). **Correction (2026-09-22, during v1.9)**: the actual bytes for the "header-note blocks of the 3 `prfrail-*` files were moved above the H1" listed in this row landed in **the v1.5 commit `e0b08ec`** (15 seconds apart from the v1.6 commit); this note governs the attribution accuracy. | historically missing |
| v1.7 | 2026-09-21 | **OB-8 fixed at the root + identifier disambiguation + user rulings recorded**: ① **OB-8 fixed at the root per the user's ruling** ("delete defaults + hard gate", **do not delete files**) — first a **field-deletion probe** on the inactive file `standard-builder` (after deleting `model` the carrier loads and works normally), and once it passed the `model` field was **deleted from all 8 blacklisted files** (measured blacklist hits = 0), while the `model` of the 3 non-blacklisted files is **retained**; added **G5-b** (the `model` value of `*.agent.md` added/modified in the change set must not be a blacklisted value) to prevent regression; ② **identifier disambiguation**: the legacy **`③.5` is renamed throughout to `⑥` (independent scan layer)** (23 occurrences), and the use of `④` as a "final-review alias" is eliminated (`独立终审员(④)`→`(⑦)`, `⑦④终审`→`⑦ 终审`, `B3c ④ 复审`→`⑦ 复审`, `④ 无效轮次`→`⑦ 无效轮次`, `（⑥ ③.5 / ⑦ ④）`→`（⑥ 独立扫描 / ⑦ 独立终审）`), while the pipeline step `④集成` and the list marker `④` stay unchanged; C.1 adds an old/new identifier mapping; ③ registered **OB-9** (the independent-verification gap for T3/T4 changes in the same area as this round, and the focus required of the next ⑦). | historically missing |
| v1.8 | 2026-09-21 | **Added Appendix B.7 "Slice Startup Card" (user adjudication: option A)**: the old briefing's two **functions with no landing place** — **session seeding (launch package structure)** and **minimum context feeding** — are folded into the directive as a **thin card** of "entry point + slots + section-number references": it contains the user kickoff message template, the master opening receipt (`SLICE: STARTED` / `SLICE: BLOCKED`), the stop-point receipt (⑩), and the **minimum reading set** by slice type; **hard constraint: this card contains no rule substance** (a rule change edits the body only and the card edits only the section-number references, structurally eliminating the old disease of "the same wording written in two places ⇒ drift"; the old briefing's "where it conflicts with the directive, the directive prevails" is exactly this kind of drift permit, not inherited); §5.0 gains one line of guidance; A.1 notes that this inheritance has landed; C.3 gains the row "briefing slice-level hard gates are not inherited"; Appendix C gains `SW-4`. | historically missing |
| v1.9 | 2026-09-22 | **Added §1.6 work tiers and startup authorization (slice `SW-5`; user rulings Q1–Q6)**: ① three-tier startup — `[ROUTINE]` (default, no trigger word) / `[SLICE]` / `[PILOT]`, with the trigger-word syntax **defined in that one place in the whole document**; ② the routine tier's **non-bypassable boundaries** (① is **a reference to §7.2 ③**; ②–⑥ are incremental boundaries: CONTRACTS / schema / fixtures, `.agent.md`, state machine / crash / ownership / concurrency code, `.github/workflows/**`, dependency additions and removals); ③ the master **determines the tier first and declares it on the receipt's first line**, **may recommend but never upgrades unilaterally**, **one-way tier movement with no mid-flight downgrade** (a user request to downgrade ⇒ stop the slice and restart under `[ROUTINE]`; a stopped state must be marked "stopped — incomplete" ＋ an unfinished list and is **not treated as a closed loop**), and **is not exempt from §10 because of the routine tier**; ④ routine-tier record-keeping = diff + §6.1 + §6.3, **0 paid calls**; ⑤ §1.1 repointed, §5.0 gains a "startup precondition" sentence and **`SW-3` is added to the step 0 upfront hard gate**, Appendix B.7 expanded for the three tiers (heading renamed from "Slice Startup Card" to "**Startup Card**"), Appendix C.0e records `SW-5` (including **the first recorded role-trimming downgrade judgement**: the letter lands "Medium", substantively "Simple", ①⑤ exempted; **cost**: budget ⑥×1 + ⑦×1, **actual ⑥×1 + ⑦×2 = 3 paid calls** (see Appendix C.0e). | historically missing |
| v1.10 | 2026-09-22 | **§6.3 extension (G6 / G7) and `tools/gates/` scripting (slice `GATES-EXT`)**: ① two new criterion classes - **G6 ledger metadata consistency** (status tense / cost clause / run-id shape / expiring literals) and **G7 review-package completeness** (artifacts fence block / resolvable path and disposition / finding disposition cell / stopped-state marker); both judge **changed lines only** (same scope as G1-a / G3 / G4); ② the criteria run from an **offline in-repo script** (single entry `gate.js`, exit codes 0 / 1 / 2, no external dependency) that migrates and replaces the throwaway `tmp/gate.js` (kept verbatim for comparison), while base **G1 / G5 stay unscripted** and print `SKIP` truthfully (see **OB-15**); ③ CI gains one **Gates** step (both legs, no `continue-on-error`, no new job or Action) and `Checkout source` gains `fetch-depth: 0` so `--scope=ci` resolves the change set on full history (**the 3rd frozen-contract edit, user-authorised**); ④ **scope guards**: an unresolvable range endpoint, a failing `git diff` / `git status`, and an inconclusive shallow probe all become **usage / environment errors** (exit 2) and **never** degrade into "empty change set ⇒ every criterion vacuous pass"; ⑤ §12.4 registers **OB-15 / OB-16 / OB-17 / OB-18 / OB-19 / OB-20**, §12.5 registers this slice's gap row, and `docs/t027/REMAINING_SLICES` registers **DR-6**; ⑥ the criterion text is aligned with the implementation line by line (G6-1's completion marker takes the **union** the directive enumerates). | historically missing |
| v1.11 | 2026-09-23 | **Added §1.7 execution freedom and escalation boundaries, §3.4 post-edit diagnostics scan, and §5.1 ⑧c wrap-up cleanup (slice `DIRECTIVE-V1.11`)**: ① §1.7 - execution freedom (single-file implementation choices / implementation approaches within the scope the contract has already defined / gate retries within an already-authorized package / refactors that do not affect contract semantics / completion of contract wording) and escalation boundaries (any change touching contract semantics / a side fix beyond the slice scope / budget quota exhausted / ambiguity in deciding "is this scope expansion or execution freedom" / the non-bypassable boundaries of §1.6), the principle "how to do it may be autonomous, what to do belongs to the user", with the record-keeping requirement re-checked by ⑦ final review; ② §3.4 - post-edit diagnostics scan (tools chosen by language / file category, no extension enumeration; any newly added warning must not be left to the next step; an "uncovered type" is not "exempt"; a diagnostics aggregator is to be implemented in `tools/gates/` with priority); ③ §5.1 gains **⑧c wrap-up cleanup** (after ⑧b and before ⑩, executed by the master, a fixed stage parallel to ⑨ native validation); ④ the §5.2 matrix gains a ⑧c column, §5.3 gains a ⑧c row, §7.1 gains a ⑧c note, and the §11.3 report structure gains a ⑧c record; ⑤ Appendix B.7.4's minimum reading set gains the new section numbers, and Appendix C gains this slice's ledger block C.0g. OB-22/23/24 were already registered at B4 wrap-up; this slice only references them and does not re-register them. | historically missing |
| v1.12 | 2026-09-23 | **DIRECTIVE-V1.12 (remediation of three rounds of third-party independent review + landing of user rulings)**: ① §0.1 authority table gains the ledger / ADR / gates / agents rows + a "gate criteria" terminology disambiguation (A1/A13/M10); §0.2 gains the MAJOR/MINOR ruling note (A4); §0.3 gains the "Slices in flight (affected)" column (A5); ② §1.2 gains the tier / probe quota / freedom-list fields (A11/L7); §1.6 gains the production-code semantic boundary and the §5.2 scope sentence, the `[id]` character set and the 20-code-point counting unit, multi-trigger hardening, and the §1.1 dividing line (A2/A3/A12/A14); §1.7 gains "the §1.6 boundary outranks §1.7.1" and the three-branch test (A3/A9); ③ §2.2 gains the modelId / UI display-name columns and Appendix B's seven call strings move to modelId (H1); R2.2/R2.4 corrected (H1/L8); §2.4 gains the blacklist-boundary note (L10); ④ §4.1 builds the closed-set verdict-line table, drops RE_, and makes bare tokens a format deviation (H2); §4.2's prohibition list excludes CONTRACTS/schema/fixtures and §5.1 gains ②a/②b (H3); ⑤ §5.1 ⑧c gains the frozen-evidence-pack exemption and §5.3 ⑧c echoes it (H5); the §5.2 high-risk row drops +blind (M1); §5.3 ⑧b gains the ⑦ re-review clause (M11); ⑥ §6.1 gains the focused-gate definition (M6); §6.3 G1 gains the stopped-state exemption (A10); ⑦ §7.1 list corrected (L1), §7.2 numbering completed (L3), §9.1 ⑥ exception (M8), the §9.2 counting rule (M7 in part); ⑧ §12.1 fixed denominator + [SLICE] startup handling (H4/A6); §12.2 gains the re-evaluation clause + the DR-FIX ruling record (H4); §12.4 closes OB-2/3/4/5/6, registers OB-28–OB-33, and §12.5's gap row flips to landed (M4); ⑨ Appendix B.0/B.6/B.7.1/B.7.4 synchronised (M2/H2/A2/M12); C.0d/C.0g gain correction notes (L9/M3); Appendix C gains the C.0h ledger block; ⑩ **disposition of the pre-commit third-party review (v1.12)**: its C1 (§0.3 log ↔ body consistency; 19 items claimed as changed but not landed) was **completed** per the user's ruling (R2.2/R2.4, B.0/B.6/B.7.1/B.7.4, §6.1, §6.3 G1, §7.1, §7.2 numbering, §9.1/§9.2, the §8.2 [ROUTINE] row, §12.1/§12.2/§12.4/§12.5, the C.0d/C.0g correction notes, the new C.0h), and its M8/M12/M13 plus L1–L4 landed in step (the §1.2 field table, the §4.1 PARTIAL handling, the §0.3 last-column reading, the §2.2 display-name example note, the §1.7.3 coverage boundary, the §5.2 size mapping; M12 and L3 share the §4.1 PARTIAL handling as their landing point); G6-5 and the three-layer consistency architecture are registered as OB-32 / OB-33; ⑪ **Follow-up note (fourth re-review, R6)**: of the 10 findings of the fourth re-review received after this row was finalised (1 High + 4 Medium + 5 Low, `tmp/DELIVERY_DIRECTIVE_v112_REVIEW_R6.md`), N1 (re-evaluation object) and N2 (pilot deferral) landed per user ruling (§12.1 and §12.2), N3, N4, N6 and N10 are ledger- and wording-class fixes that landed, N5 was measured as not reproducible (G6-4 needs no whitelist entry), and N7, N8 and N9 are registered as OB-34 – OB-36 for the next slice; every remediation in this round came from an explicit user ruling or from ledger-metadata / wording-class fixes, so under the OB-11 exception ⑦ is **not re-run** (**Correction note (fifth round, R7)**: the actual basis of that exemption is a **direct user ruling**, which goes beyond the literal standard of the OB-11 exception - see OB-37); ⑫ **Follow-up note (fifth round, R7)**: of the fifth re-review's 9 findings (1 High + 4 Medium + 4 Low, `tmp/DELIVERY_DIRECTIVE_v112_REVIEW_R7.md`), P1 (dual-track priority rule), P2 (deferral companions) and P3 (re-evaluation judge / anchor / count reading) landed per user ruling (§12.1 and §12.2), P4 and P5 are registered as OB-37 and OB-38, and P6, P7, P8 and P9 are ledger- and wording-class fixes that landed; **this round's remediation carries rule-semantics changes, so the OB-11 exception does not apply**, and ⑦'s third round is authorized in the same batch (narrowed input: only P1 / P2 / P3's rule changes and whether P4 / P5's registration is truthful). | no blocked slices; follow-on slices `GATES-G1G5` / `GATES-TIGHTEN` need write-back |
| v1.13 | 2026-09-23 | **DIRECTIVE-REVIEW-SATURATION (review saturation criterion)**: ① adds **§7.8 Review Saturation Criterion** - the **sole authority** for the three stop objects (third-party independent re-review / ⑦ round N / re-review after remediation) and the five situations (quota / user-authorized exception / OB-11 / OB-37 / pilot deferral); ② §6.3 lands two mechanical sub-checks, **G6-5** (section-reference resolution, the approximate criterion for OB-32) and **G6-6** (ledger status closed set, OB-38), and backfills the first status word of every §12.4 ledger row against the closed set; ③ closes **OB-34** (§7.2.2 heading), **OB-35** (§1.6 citing §7.2.1 ③), **OB-36** (the four elements in OB-14) and **OB-41** (§9.3 failure-mode row) within this slice, and splits out **OB-30 / OB-31 / OB-33**; ④ §12.5's "the §5.3 vs §7.5 conflict decision standard" flips to **Landed**; ⑤ Appendix C gains this slice's ledger block **C.0i**. The clauses are in §7.8, §6.3, §9.3, §12.4 and §12.5. It also registers OB-42 (no mechanical criterion for the verdict-line closed set), OB-43 (mojibake in evidence capture) , OB-44 (the review layer writing to disk) and OB-45 (the review layer's report of its data source). | no blocked slices; follow-on slices `GATES-G1G5` / `GATES-TIGHTEN` need write-back |
---

## 1. Scope and Slice Model

### 1.1 Scope

- Applies to: **all engineering activity in the ProofRail repository that produces a deliverable** (code, contracts, schemas, fixtures, scripts, documents, evidence bundles).
- Does not apply to: purely consultative Q&A; ad-hoc troubleshooting with a single command (these need no slice and no review, but must not produce any repository change).
- Work tiers (`[ROUTINE]` / `[SLICE]` / `[PILOT]`) and each tier's flow, review layer, record-keeping and cost are in **§1.6**; the "does not apply" items of this section **produce no repository change** and are **not the same class** as the `[ROUTINE]` tier of §1.6 (for the difference see the §1.6 「`[ROUTINE]` Applies」 paragraph).

### 1.2 Slice

**A slice = the smallest reviewable delivery unit**. Every slice must have a **slice definition** containing:

| Field | Requirement |
|---|---|
| Identifier | A stable ID of the form `A0`/`B3c`/`DR-2` |
| Goal | One sentence, truth-decidable |
| Dependencies | Preceding slices or external authorization |
| Size | S / M / M-L / L |
| Steps | A checkable checklist |
| Current gap | Why not doing it now would cause problems |
| Deliverables | File-level list |
| Acceptance evidence | Machine-verifiable assertions |
| Boundary | What is **explicitly not** done (scope-creep prevention) |
| Tier | [ROUTINE] / [SLICE] / [PILOT] (see §1.6; no trigger word ⇒ [ROUTINE]) |
| Probe quota | the cap declared for real external calls (§7.6); write "none" if not needed |
| Freedom-list path | see the record-keeping requirement in §1.7.3 (fixed path) |

### 1.3 One slice at a time

At any moment **only one slice may be "in progress"**.

**Parallelism policy (hard rule R1.1)**: parallelism **is allowed only at the master's own tool-call layer** (for example, issuing several read-only commands in one message); **multiple subagents must not be invoked in parallel** (`runSubagent` is serial and blocking, see R5.1). Cross-slice parallelism is allowed only with explicit user authorization, and **read-only**; **do not** run real-machine experiments that interfere with each other (processes / Job Objects / ports / temp directories) in parallel, see the original §四.11 text (Appendix A).

### 1.4 Contracts first (hard rule)

For any semantic change the order is fixed: **CONTRACTS in both languages → schema / fixtures → code → focused tests → document write-back**. Code must not go first.

### 1.5 Workspace hygiene

- The workspace must be clean before work starts; uncommitted historical changes must first be disposed of by the user (commit / stash / discard).
- Temporary artifacts go only to the repository root `tmp/` (gitignored), and must be **removed immediately after use**, never left behind across slices.
- Do not run interfering real-machine experiments (processes / Job Objects / ports / temp directories) in parallel while a slice is running.

### 1.6 Work tiers and startup authorization

**This directive defines the trigger-word syntax in this one place only; every other subsection merely references §1.6.**

**Three tiers** (with no trigger word the default is `[ROUTINE]`):

| Tier | Trigger | Flow | Review layer | Cost |
|---|---|---|---|---|
| Routine `[ROUTINE]` | no trigger word (default) | the master **personally** executes + §6.1 base gates + §6.3 self-check; **does not go through the §5.1 pipeline and dispatches no subagent** | ⑤⑥⑦ **all exempt** | **0 paid calls** |
| Slice `[SLICE]` | see the trigger-word syntax below | goes through the pipeline per §5.1; roles trimmed per the §5.2 matrix | per §5.2 | per §7.5 |
| Pilot `[PILOT]` | see the trigger-word syntax below | same as slice, **plus** the §12.1 quantitative metrics declaration | same as slice (per §5.2) | same as slice (per §7.5) |

Where `[ROUTINE]` applies: conversational collaboration, small changes, producing proposals; it **may produce repository changes** (writing to disk and authorization follow §10, unlike the "does not apply" items of §1.1). **The decidable dividing line from §1.1's "does not apply" = whether anything is written to disk**: producing repository changes (writing to disk) ⇒ `[ROUTINE]`; producing no repository change at all (pure consultation / ad-hoc troubleshooting) ⇒ §1.1 "does not apply".

**Trigger-word syntax**:

- Form: `「切片：<id> <描述>」` / `SLICE:<id> <description>` / `「试点：<id> <描述>」` / `PILOT:<id> <description>`.
- Position: it takes effect only if it appears **within the first 20 characters of the message's first line** (to prevent a trigger word mentioned in the body from being misjudged).
- **Counting unit**: "the first 20 characters" is counted in **Unicode code points**; leading whitespace and `>` quote markers are **skipped** before counting; the **first character** of the trigger word must fall inside that 20-code-point window.
- **`[id]` character set**: letters, digits, `-`, `_`, `.` (1–32 characters, no whitespace and no full-width punctuation); a historical ID containing whitespace (e.g. `B4 · AT-23`) is rewritten as `-` or `.` inside a trigger word, while its ledger display form is unchanged.
- Full-width / half-width colons are both accepted; the English form is **case-insensitive**.
- **The routine tier has no trigger word.**
- If a trigger word appears **beyond the first 20 characters of the first line** (for example, mentioned in the body), it **does not take effect**; the master judges it as `[ROUTINE]`, but **must** explicitly note in the receipt "**a suspected trigger word was detected but its position is invalid**" and recommend that the user resend it.
- If **more than one trigger word matches within the first 20 characters of the first line** (including the same form appearing twice or two different IDs), it **does not take effect**; the master **must stop and ask the user to resend**, must not choose one itself, **and must not continue under `[ROUTINE]`**.

**The routine tier's non-bypassable boundaries** (touching one means an upgrade must be recommended; **execution must not continue**):

- ① **all the semantic categories listed in §7.2.1 ③** (state machine / gate criteria / evidence model / roles / authorization contract). **This item is a reference, not a redefinition.**
- ②–⑥ are this tier's **incremental** boundaries: the semantics of CONTRACTS / schema / fixtures; any `.agent.md`; code touching state machines, crash paths, ownership, or concurrency (the production-code boundary **governs this item, which is listed for contrast only**); `.github/workflows/**`; **adding / removing / upgrading / downgrading dependencies (including lock-file changes)**.
- **Also (production-code semantic boundary, user ruling 2026-09-23)**: **any semantic change** to production code (`.go`) — additions / modifications / refactors / error handling / branches / return values — **must go through `[SLICE]`**; changes that **do not alter behaviour**, such as pure comments, pure formatting or docstrings, may still use `[ROUTINE]`.

> **Scope-relationship ruling (2026-09-23)**: "⑦ cannot be omitted" in §5.2 and R5.3 constrains only slices inside the §5.1 pipeline; "⑤⑥⑦ all exempt" in §1.6 applies only to the `[ROUTINE]` tier (which does not go through the §5.1 pipeline). The two scopes are mutually exclusive and do not conflict — once `[SLICE]` / `[PILOT]` is entered, ⑤⑥⑦ are enabled per §5.2 and ⑦ cannot be omitted.

> **Handling an over-reaching instruction**: if the user explicitly says "complete it in the routine tier", the master **must refuse** and give the reason; if the user insists, the master must **record the "user over-reaching instruction" in writing**, then execute the following downgrade: **do only the minimal part that does not touch a boundary**, and explicitly put the remainder in the "**unfinished list**", stating that it must be re-run under `[SLICE]`.

**The master's determination and recommendation duty**:

- On receiving a message, **determine the tier first**, and declare `[ROUTINE]` / `[SLICE]` / `[PILOT]` ＋ a one-sentence reason on **the receipt's first line**.
- **May recommend, never upgrades unilaterally**: if the work turns out larger than expected ⇒ reply with an upgrade recommendation and **stop**, waiting for the user to say so explicitly; if the user does not upgrade ⇒ the master may only switch to the part outside the boundary or stop entirely, and **must not** keep touching boundary content under the name of "a small scope".
- **One-way tier movement**: after entering `[SLICE]` / `[PILOT]`, **no mid-flight downgrade**; if the user explicitly requests a downgrade, **stop the current slice** and **restart** under `[ROUTINE]`, not switching inside the original slice; if `[ROUTINE]` turns out larger ⇒ recommend an upgrade.
  - **Artifact marking in the stopped state**: artifacts already produced **stay where they are** but **must not be committed as a completed slice**; their report header must carry `**Status: stopped — incomplete**` and list the "**unfinished list**".
  - **A stopped state is not treated as a closed loop**: it must not be handled as "complete" in the ledger / `DEV_PLAN`; that row of `REMAINING_SLICES` stays `🔄` or reverts to `⬜` (**the user decides**). The write-back duty of §11.2 targets "slice wrap-up" and **does not apply** to a stopped state.
  - **Restart path**: restarting under `[ROUTINE]` **counts as new work**; the stop record **stays in the original slice report** as context for "why it was downgraded to the routine tier".
- **Not exempt from §10 because of the routine tier**: `[ROUTINE]` changes **may accumulate** into one commit, but committing / pushing **still requires the same-round explicit authorization of §10**.

**Record-keeping and cost**:

- `[ROUTINE]`: minimum record-keeping = diff ＋ §6.1 gate results ＋ §6.3 self-check results; **the §6.3 self-check applies in full** (no routine-tier exemption); **0 paid calls**.
- `[SLICE]` / `[PILOT]`: record-keeping per §5.3 and §11, cost per §7.5.

> **Disambiguation (two different dimensions)**: the work tier (`[ROUTINE]` / `[SLICE]` / `[PILOT]`) and **the §8.2 effort tier (low / high / max)** are **two different dimensions** and must not be mixed — the former decides the flow and the review layer, the latter is set once by the operator before startup; a tier declaration does not replace the §8.2 effort-tier setting, and the effort-tier setting does not replace this section's tier determination. Also: the effort tier for the `[ROUTINE]` tier is defined by the `[ROUTINE]` row in §8.2.

### 1.7 Execution freedom and escalation boundaries

**1.7.1 Execution freedom (autonomous)**: the master or a product-layer role **may decide the following on its own** without stopping to ask the user:

- single-file implementation choices;
- implementation approaches within the scope the contract has already defined;
- gate retries within an already-authorized package;
- refactors that do not affect contract semantics;
- completion of contract wording.

**1.7.2 Escalation boundaries (must stop and ask the user)**: touching any of the following **must stop** and confirm with the user; it must not continue in the name of execution freedom:

- any change touching contract semantics;
- a side fix beyond the slice scope;
- budget quota exhausted;
- ambiguity in deciding "is this scope expansion or execution freedom";
- the non-bypassable boundaries of §1.6.

**The non-bypassable boundaries of §1.6 outrank the execution freedom of §1.7.1** - anything touching a §1.6 boundary is always handled as an escalation boundary and must not invoke §1.7.1.

**1.7.3 The distinction principle**: "how to do it" may be autonomous; "what to do" always belongs to the user. The test (three branches): **if this decision were wrong, what would need to change?** - **code** ⇒ execution freedom; **contract / direction / scope / budget / authorization** (i.e. the items listed in §1.7.2) ⇒ an escalation boundary. Anything the list does not cover is handled per the last item of §1.7.2 (stop when in doubt).

**Record-keeping requirement**: the master **does not report line by line** when exercising execution freedom, but the "change summary" section of the slice validation report **must** list the substantive choices made under execution freedom. **Production and landing place**: that list is written to disk by **the master** as an **orchestration artifact** between ⑥ and ⑦ at the fixed path `docs/validation/evidence/[slice-id]-freedom-list.md`; ⑧b finalisation then writes the list back into the validation report and the two must agree; ⑦ final review re-checks **the same list inside its review input package** (the list is produced before ⑦, which is compatible with the §5.1 order) and rules on whether the boundary was crossed. **Coverage boundary**: that path lives under `docs/validation/evidence/**` (outside G7's scan scope) ⇒ the list **has no mechanical-criterion protection** and rests on ⑦'s manual re-check (same exclusion reading as §6.3 G7).

---

## 2. Roles and Model Allow-List

### 2.1 Role model (three layers)

| Layer | Role | One-line responsibility | May write files |
|---|---|---|---|
| **Orchestration layer** | Master | Command, scheduling, gates, evidence, fallback, reporting | Only orchestration artifacts (§3.2) |
| **Product layer** | Implementer / Tester / Documenter | The **sole** producer of product artifacts in their own domain | Yes |
| **Review layer** | Architect / Pre-reviewer / Independent scanner (⑥) / Independent final reviewer (⑦) | Read-only reasoning and gates | **No** |

### 2.2 Role–carrier–model triple (**`model` must be passed explicitly**)

| Role | Call point | `agentName` (carrier) | `modelId` (call string, sole authority) | UI display name (for manual cross-check only) | Same model as the master? | Tool capability |
|---|---|---|---|---|---|---|
| Architect | ① | `deep-reasoner` | `deepseek-v4-pro` | `DeepSeek V4 Pro (deepseek)` | No | Read-only |
| Pre-reviewer | ⑤ | `deep-reasoner` | `deepseek-v4-pro` | `DeepSeek V4 Pro (deepseek)` | No | Read-only |
| Implementer | ② | `prfrail-implementer` | `deepseek-flash` | `DeepSeek V4.1 Flash (deepseek)` | **Yes (transitional state)** | Read/write + execute |
| Tester | ③ | `prfrail-tester` | `deepseek-flash` | `DeepSeek V4.1 Flash (deepseek)` | **Yes (transitional state)** | Read/write + execute |
| Documenter | ⑧ | `prfrail-documenter` | `deepseek-flash` | `DeepSeek V4.1 Flash (deepseek)` | **Yes (transitional state)** | Read/write + execute |
| Independent scanner | ⑥ | `independent-reviewer` | `mai-code-1.1-flash` | `MAI-Code-1.1-Flash (copilot)` | No | Read-only |
| Independent final reviewer | ⑦ | `independent-reviewer` | `gpt-5.3-codex` | `GPT-5.3-Codex (copilot)` | No | Read-only |

**`modelId` is the sole authority for the call string** (R2.2): task packages and Appendix B templates always write `modelId`. A UI display name is for manual cross-check only and **must not be written into task packages or Appendix B templates**; the display name follows whatever the operator UI currently shows (its version number, e.g. `V4.1`, is not part of the ID) and **is not treated as a stable value**; the "UI display name" column of this table carries **example values, subject to what the operator UI currently shows**.

> **Known limitation (transitional state)**: the three product-layer roles share `deepseek-flash` with the master ⇒ **not only the effort tier but the model capability itself is inseparable** (see §8.3). **This does not exempt R2.1** - every invocation must still pass `modelId` explicitly; this limitation lifts automatically once the long-term fix of §8.3 (a different model per role) takes effect.

**Carrier readiness status (as of 2026-09-23, actually checked)**: `deep-reasoner` (architecture/pre-review) and `independent-reviewer` (⑥/⑦) are **ready**; `prfrail-implementer` / `prfrail-tester` / `prfrail-documenter` have been **created by slice SW-2, and all three carriers passed real-machine probes** (each completed "create → read back → delete" with 0 residue; evidence in `docs/validation/sw-directive-switchover.md` §5) ⇒ **all seven carriers listed in §2.2 are ready**. The 11 existing generated role files have been **normalized to BOM+LF and fully committed by slice SW-3** ⇒ **all 14 role files (3 `prfrail-*` plus 11 existing) are ready** (SW-3 status row and acceptance evidence in appendix C.0c), **with nothing left outstanding**, so both the §5.0 step 0 upfront hard gate and the step 2 role-readiness pre-check pass normally.

**Hard rule R2.6 (call isolation, anti-anchoring)**: ① the architect and the ⑤ pre-reviewer are **two independent invocations of the same carrier**, with **no shared context** — ⑤ may receive only "the solution document + the original implementation diff (**paths**) + verbatim contract sentences", and **must not** receive ①'s conversation history, the master's implementation-phase reasoning, or any pre-implementation assumption. ⇒ **Anti-anchoring relies on "call isolation", not on "role or model separation"**; ⑤ must not be cancelled on the grounds of "same carrier", nor may ①'s reasoning be passed to ⑤.

**Hard rule R2.1**: `runSubagent` **silently falls back to the session's main model when `model` is not passed**. Therefore "not passing model" is forbidden; every invocation must state the model string in the task package.
**Hard rule R2.2**: Write only currently-active model IDs; **a UI display name must not be used as a call string in task packages or templates** (see the modelId note in §2.2). For Flash write only `deepseek-flash` (`deepseek-v4-flash` and `deepseek-v4-flash-vision-exp` are already listed in `LEGACY_FLASH_MODEL_IDS` and are forbidden).
**Hard rule R2.3**: **Subagent depth ≤ 1, call chains forbidden** (no master → subagent → another model). Therefore orchestration-type agents carrying the `agent` tool **must not enter the role set** (see §2.4).
**Hard rule R2.4**: Role files are the **reproducible source of the role set**. Adding or editing a role must go through version control (see §2.5).

### 2.3 Model self-reporting is untrustworthy

Subagents **cannot reliably self-report** their model or effort tier (measured: a Flash subagent self-reported `GitHub Copilot（系统提示声明为 DeepSeek V4.1 Flash）| vendor=unknown`). ⇒ **Do not** use "the subagent says it is model X" as acceptance evidence; the model and effort tier can only be set by the operator on the UI side and **recorded** in the task package.

### 2.4 Prohibited list (blacklist)

The following models **must not be invoked**, under any circumstances: Luna, the Gemini family, Terra, GPT-5.4, **Sol (strategic master)**, the Kimi family, and "any other model doing drafting / review / formatting / checklisting / validation / coordination". **The blacklist covers the whole GPT-5.x family** (such as GPT-5.4 and GPT-5.6-Sol) **unless §2.2 lists them explicitly** (currently only `gpt-5.3-codex`).

**Orchestration-type agent prohibited**: `sol-orchestrator` (`model: gpt-5.6-sol`, with `agent` among its `tools` and declaring `agents: [...]`) —
- It carries call-chain capability, violating R2.3;
- Its model is on the blacklist, and not passing model explicitly falls back to the master itself ⇒ producing "fake orchestration";
- Its built-in `routing` / `Kimi K3` / `ultra tier` rules are **inconsistent** with this directive's allow-list and cost ceiling;
- Per its own configuration `strategy.manual_override.detection: 非 Sol 主控会话即视为手动模式，跳过全部子代理自动编排` — **when the current master is not Sol, its automatic orchestration is off by design**, so there is no need to enable it in the first place.

**Sole exception**: layer ⑥ is cleared for `mai-code-1.1-flash`; Haiku is cleared as a ⑥ supplement on a time-limited basis only when all enabling preconditions in §7.2.3 are satisfied.

### 2.5 Versioning of the role set (D1-a selection)

**Rule R2.5**: `.github/agents/*.agent.md` is **always** put under version control (**including the existing generated role files**); the **generator toolchain is not committed**.

- Version-control scope: `.github/agents/*.agent.md` (role definitions, **including the 11 existing** files originally generated by `sol-orchestrator`).
- Still ignored: `.github/agents/sol-orchestrator/` (the generator toolchain with its `__pycache__`, `_discovered_models*.json`, pricing tables), `.github/agents/*.managed.json`.
- **Header-note hard rule**: **all 14 role files** (3 `prfrail-*` + 11 existing) must carry their **header note in a fixed position**: the **first non-empty line** after the end of the frontmatter must be **verbatim equal** to the following line:
  `> **非生成物（sol-orchestrator 已按 DELIVERY_DIRECTIVE §2.4 禁用）：手工维护；不得由生成器覆盖。**`
  **Mechanical criterion**: ① parse the end position of the frontmatter, take `the first non-empty line after that point` and compare it with the header-note text; it **must be equal** (**not** accepting "merely appearing in the document" — otherwise moving the header note into the body would also pass, a Low finding by ⑦ round 3); ② supporting evidence: `grep -L "非生成物（sol-orchestrator 已按" .github/agents/*.agent.md` produces **no output**.
- **If the generator is enabled in the future**: `_config.yaml` and the generation scripts must first be brought under version control, the generation re-run, and the output compared **verbatim** against the hand-maintained files; differences must be explained before switching.
- **Known latent defect (**OB-8, already fixed at the root on 2026-09-21**)**: an actual check of the 11 existing role files found **8** of them (not 9; corrected by ⑦ round 3) whose default `model` is on the §2.4 blacklist — family breakdown: **Terra 3** (`expert-arbiter`/`heavy-builder`/`planning-specialist`), **Luna 2** (`fast-implementer`/`quick-researcher`), **Gemini 1** (`context-researcher`), **GPT-5.4 1** (`standard-builder`), **GPT-5.6-Sol 1** (`sol-orchestrator`); the 3 non-blacklisted ones: `deep-reasoner` (`deepseek-v4-pro`), `independent-reviewer` (`gpt-5.3-codex`), `quick-verifier` (`deepseek-flash`). ⇒ Invoking directly by `agentName` without passing `model` explicitly can hit the blacklist while bypassing R2.1. **Disposition (user ruling: "delete default values + hard gate", no file deletion)**: ① a **field-deletion probe** was completed on the inactive file `standard-builder` (after deleting the `model` line the carrier **loads and works normally**); ② the probe passed ⇒ the `model` field was **deleted** from **all 8 blacklisted files** (measured blacklist hits = 0); ③ the `model` of the 3 non-blacklisted files is **retained** (as legitimate defaults); ④ **G5-b** was added (§6.3) to prevent regression.
- **Additional residual risk (noted by ⑦ round 3, now closed by G5-b)**: `.agent.md` files are now all under version control ⇒ future `.agent.md` files added by a generator re-run may **directly bring in a blacklisted default `model`**. ⇒ **G5-b** has been added (§6.3): the `model` value of added/modified `*.agent.md` must not be a blacklisted value (mechanical criterion: blacklist hits in the change set must be 0). R2.5's "prerequisites for re-running the generator" still stand.
- **Ignore-rule alignment (accepted by SW-2)**: the directory-level rule `.github/agents/` in `.gitignore` must be **narrowed** to ignore only the generator toolchain and its outputs — keep ignoring `.github/agents/sol-orchestrator/`, `.github/agents/*.managed.json`, `**/__pycache__/`, and **allow** `.github/agents/*.agent.md`. Acceptance assertion: `git check-ignore -v .github/agents/prfrail-implementer.agent.md` produces **no output** (not ignored), while `git check-ignore -v .github/agents/sol-orchestrator/_config.yaml` **still matches**.

### 2.6 Minimum contract for role files

Every role file must contain: `name` / `description` / `tools` / `user-invocable` / `target`, plus the four body sections **responsibilities, inputs, output contract, prohibitions** (template in Appendix B). **The `model` key must be omitted** — **do not** write an empty value, nor a concrete model. Reasons (all measured on 2026-09-21): ① writing `model: ""` makes VS Code **fail to load that carrier** (`runSubagent` reports `Requested agent '...' not found`; isolation verification: adding the empty-value key ⇒ failure, deleting the key ⇒ recovery); ② writing a concrete model value conflicts with R2.1 "must be passed explicitly by the caller" and invites silent fallback. ⇒ **Omitting the key** is the only form that satisfies both "loadable" and "no default model"; the body must explicitly declare in one line that "`model` is left empty and is passed explicitly by the master in `runSubagent` (R2.1)".

**Transitional clause**: `prfrail-*` role files **newly created** after this directive takes effect must satisfy the four-section contract **immediately**; existing role files produced by the `sol-orchestrator` generator (`deep-reasoner`, `independent-reviewer`, etc.) are **exempt** from the body-section requirement — they are used only as `model`/`tools` carriers, and their responsibilities/inputs/outputs/prohibitions **are governed by §4.2 and Appendix B of this directive**.

---

## 3. Master Responsibilities and Boundaries

### 3.1 What the master does

1. **Task decomposition and dispatch**: break the slice into atomic work packages, stating for each: input materials, expected output, acceptance conditions, and the **model string**.
2. **Minimal context feeding**: give subagents only the materials needed to finish the task; for any reviewer, always give **evidence paths** instead of paraphrases (§6.2).
3. **Gate execution**: personally run build / vet / focused tests / encoding gate / wrap-up self-check at key points (§6).
4. **Review scheduling**: schedule in the order ⑤→⑥→⑦ and decide pass or return based on the conclusions; enforce closure discipline (§5.3).
5. **Evidence management**: persist diffs, gate output, and probe results, and deliver their paths (§6.2).
6. **Fallback**: see §9.
7. **Stop points**: §10.4.

### 3.2 What the master does not do (hard boundaries)

- **Does not originate product artifacts**: no production code, test code, CONTRACTS/schema/fixtures, or validation-report body text. After locating a problem, **assign the corresponding subagent** to fix it.
- **Does not draw conclusions on behalf of reviewers**: must not decide by itself that "pre-review / scan / final review passed".
- **Exception (orchestration artifacts, allowed and required from the master)**: task packages, status tables, evidence files (persisted diffs/logs), exception and fallback records, report **assembly** (aggregating subagent outputs and gate results into a report without rewriting their conclusions), ledger and status-marker write-back.
- **Mechanical integration exemption**: applying a subagent-produced patch to files counts as **mechanical integration**, not "writing code", but it must pass three checks: ① the diff line count is identical before and after persistence; ② for `.md` files the encoding gate passes; ③ focused gates are run immediately after persistence.
- **Exemption extension (directive-governance slices, user ruling 2026-09-23)**: the validation report of an **Appendix C entry** is produced directly by the master, as an **orchestration-artifact exemption**; ① it is **limited to** Appendix C entries (other slices' validation reports still follow this section's original rule); ② it **does not exempt** ⑤⑥⑦ review; ③ it **does not change** the rules for other slices.

### 3.3 The master's self-check obligation

The master must run the wrap-up self-check of §6.3 on **every round of its own changes**; on finding an inconsistency, **fix it on the spot** or assign a subagent to fix it, and never leave it for the user or a reviewer to find.

### 3.4 Post-edit diagnostics scan (hard step)

**After any edit and before moving to the next step**, the master must run a diagnostics scan over this edit following the categories below; this step stands alongside the self-check of §3.3 and **cannot be skipped**.

**3.4.1 Tool selection (by language / file category, no extension enumeration)**:

- compiled languages (Go and alike): `go build ./...` + `go vet ./...`;
- scripting languages: that language's syntax checker / static analyzer;
- markup and data (`.md` / `.json` / schema): the encoding gate (§6.1 encoding row) ＋ a structure parser;
- plain text and others: the encoding gate;
- uncovered types: **fallback** - the encoding gate ＋ a manual readability check.

**3.4.2 Hard constraints**: any newly added warning **must not be left to the next step** - it must be fixed in place or assigned to a subagent, and the scan re-run until there is no newly added warning. An "uncovered type" is **not "exempt"**: it must still pass the encoding gate ＋ a manual readability check.

**3.4.3 Tools first**: implement a "diagnostics aggregator" in `tools/gates/` with priority (a single entry point that aggregates each language's diagnostics output per the category table of §3.4.1).

**3.4.4 Relation to the encoding layer**: §6.3 G4 (encoding + anomalous characters) and the §6.1 encoding row cover the **encoding layer**; this section covers the **semantic layer** (generating, aggregating, and zeroing diagnostics warnings); the two complement each other and do not replace one another.

**3.4.5 Basis**: the user's observation of 2026-09-23 - Flash often leaves warnings behind after editing code; the root cause is the missing fixed step of "post-edit diagnostics", **not the effort tier**.

---

## 4. Subagent Responsibilities and Output Contracts

### 4.1 Common Output Contract (All Roles)

1. **Conclusion line**: the last line must be a form from the **closed-set enumeration** below; **a bare `PASS` / `FINDINGS` (no prefix) or any form outside the table counts as a format deviation** and is handled per §7.2.

| Role / stage | Prefix | Legal enumeration values |
|---|---|---|
| ① Architect | `DESIGN` | `DONE` / `NEEDS_CLARIFICATION` |
| ② Implementer | `IMPLEMENT` | `DONE` / `BLOCKED` |
| ③ Test Engineer | `TEST` | `DONE` / `BLOCKED` |
| ⑤ Pre-reviewer | `PRE-REVIEW` | `PASS` / `PASS WITH FIXES` / `FINDINGS` |
| ⑥ Independent scanner | `INDEPENDENT SCAN` | `PASS` / `FINDINGS` |
| ⑦ Independent final reviewer (final / re-review / blind) | `FINAL REVIEW` / `RE-REVIEW` / `BLIND REVIEW` | `PASS` / `PASS WITH FIXES` / `PARTIAL` / `FINDINGS` |
| ⑧ Documenter | `DOCS` | `DONE` / `BLOCKED` |
| Master stages (④ / ⑨ / ⑧c / ⑩) | `SLICE` | `STARTED` / `READY_FOR_REVIEW` / `BLOCKED` |
| Directive document review (third-party independent review, directive-governance slices only) | `DIRECTIVE REVIEW` | `PASS` / `FINDINGS` |

The residual `RE_` placeholder form is abolished (不得残留).

**Handling of `PARTIAL`**: `PARTIAL` is treated as a subclass of `FINDINGS` (partially closed) and is handled by "eliminate, then re-run the ⑦ re-review" (§5.3).
2. **Locatable**: every finding must carry `file:line` (or `file:?` + a verbatim quote); prose alone is not acceptable.
3. **Severity**: four levels `Critical / High / Medium / Low`, defined in §4.3.
4. **Proposed fix**: give text or a patch that can be written directly; when bilingual, give both the CN and EN sentence.
5. **No scope expansion**: must not propose unfreezing conclusions, skipping review, or lowering gates.
6. **Read-only roles must not modify files**; if a change is needed, hand it to the master for dispatch.

### 4.2 Per-Role Responsibilities

| Role | Input | Required output | Prohibited |
|---|---|---|---|
| Architect | slice definition, relevant contracts, historical rulings | Architecture Design Proposal: conclusion, change-point list (file level), test-point list (executable), risks and mitigations | write code; **carry pre-implementation reasoning context into ⑤** (anti-anchoring relies on call isolation, see R2.6 — ①'s conversation history or implementation-time assumptions must not be passed to ⑤) |
| Implementer | proposal, contracts, exact context of the files to change | code diff + implementation notes (which change points were covered, which were not) | make architecture decisions; write tests; modify non-code documents (except **CONTRACTS / schema / fixtures** - these three are **contract artifacts** and the implementer owns them, see §1.4 and ②a of §5.1) |
| Test Engineer | proposal test-point list, implementation diff | test diff + coverage matrix (positive/negative/boundary/concurrency/crash) | modify production code; modify test infrastructure |
| Documenter | proposal, implementation diff, test diff, current documents | document diff + change summary | modify code; add subjective judgment |
| Pre-reviewer | proposal, implementation diff, verbatim contract text | Pre-review Report: overall conclusion + deviation items + counterexample analysis (concurrency/crash/boundary) | modify files; look at test code (keep an independent perspective) |
| Independent scanner (⑥) | diff, verbatim contract text | Section A–E report + conclusion line | modify files; reference others' lists (stay independent) |
| Independent final reviewer (⑦) | final diff + verbatim contract text (**no list attached in the first round**) | four sections: security / architecture consistency / completeness / test counterexamples + conclusion line | modify files; in the first round reference any intermediate conclusion |

### 4.3 Severity Definitions

| Level | Meaning | Handling |
|---|---|---|
| Critical | data loss, security boundary failure, silently wrong conclusion | must fix and re-run that review stage |
| High | functional failure, crash, contract conflict | same as above |
| Medium | latent defect, coverage gap, bilingual mismatch | same as above ("Medium+" is blocking) |
| Low | style, naming, comments, readability | recorded on the ledger, may be deferred |

---

## 5. Workflow (Serial Pipeline)

### 5.0 Startup (Slice Entry)

0. **Upfront hard gate (directive switchover period)**: slices **SW-1** (reference rewriting), **SW-2** (role carriers ready) and **SW-3** (role-set check-in and encoding normalization) must all have passed acceptance; otherwise this directive's pipeline **must not start** (prevent two directives running in parallel). For their definitions and acceptance commands see appendices C.0 / C.0b / C.0c.
1. The user provides the **slice definition** (the field table of §1.2) or indicates its location in the ledger.
2. Master verifies: whether dependencies are satisfied, whether the workspace is clean (§1.5), whether budget/authorization is still needed (§7.6); and performs a **role readiness pre-check** — check one by one that the role files required by §2.2 exist, that **the `model` key criterion of a role file is governed by §2.6 and by G5-a / G5-b in §6.3** (this step restates no rule substance; failure mode in §9.3), that **the `tools` of file-writing roles include `edit`** (passing a model without choosing the right carrier yields no write permission), and that **`modelId` must be passed explicitly at call time** (R2.1); any mismatch ⇒ mark "role blocked" and stop.
3. Master determines the slice type and **sets the effort tier** (§8.2; once set, it is not switched during continuous runs).
4. Generate the **task package** (each package contains: input materials, expected output, acceptance conditions, the `modelId` string), then proceed to ①.
5. If any item is not satisfied ⇒ **mark "dependency blocked" and stop**; must not infer on its own or degrade execution.

> **Startup precondition**: this pipeline **starts only when a trigger word appears** (trigger-word syntax see §1.6 — that section is the sole definition).
>
> **Entry template**: the user kickoff message, the master opening receipt, the stop-point receipt, and the "minimum reading set" are in **Appendix B.7 (Startup Card)** (its hard constraint is in the first paragraph of B.7).

### 5.1 Pipeline

```
①architecture (if needed) → ②implementation → ③testing → ④master integration+gates → ⑤pre-review (V4 Pro)
   → ⑥ independent scan (MAI) → ⑦ final review (Codex) → ⑧a document draft → ⑨ native validation → ⑧b document finalization → ⑧c wrap-up cleanup → ⑩ stop point (awaiting authorization)
```

**Note: ⑧ is split into two steps (prevent document lag)** — **⑧a draft** is done immediately after ⑦ (does not block ⑨; get the settled conclusions onto disk first); **⑧b finalization** is done after ⑨, absorbing native validation results and any rollback fixes. ⇒ Avoids "a ⑨ failure forcing a document rewrite".

**Note: ② is split into two steps (contract first)** — **②a contract**: change CONTRACTS / schema / fixtures first (order per §1.4); **②b code**: change code afterwards. Both steps are executed by the implementer and each leaves its own artifact on disk (R5.2).

**Note: ⑧c wrap-up cleanup (fixed stage)** - ⑧c runs after ⑧b and before ⑩, executed **by the master**, following the category table of §3.4: whole-repo diagnostics, IDE diagnostics summary, `gate.js --all --scope=tree` all green, the encoding gate (§6.1 encoding row / §6.3 G4, **file-by-file for all types, not only the change set**), and a checklist of the temp directory / probe leftovers / untracked files. **Responsibility assignment**: ⑧c is executed by the master, but **not as "model self-awareness"** - it is a **fixed stage** of the pipeline (**not subject to §5.2 trimming**), **parallel** to ⑨ native validation. **Long-term evolution**: if ⑧c repeatedly finds the same class of problems, consider promoting it to an **independent review stage** (§12 revision to assess). **Exemption (2026-09-23)**: archived **frozen evidence packs** carrying a MANIFEST/SHA256SUMS do not take part in the file-by-file encoding rewrite - ⑧c checks evidence packs by **self-consistency** instead (the MANIFEST list and SHA256SUMS must verify); a **new** evidence pack must still pass the encoding gate before it is committed.

**Hard rule R5.1 (serial)**: `runSubagent` is **serial and blocking** (agents are not asynchronous, not backgrounded). "Parallelism" is allowed only at the master's own tool-call layer (e.g. sending several read-only commands in one message); it **must not** be promised for subagents in task packages or flowcharts.
**Hard rule R5.2 (every stage leaves an artifact)**: every stage must leave a verifiable artifact (diff / report / gate output / evidence file). An "already passed" without an artifact does not hold.
**Hard rule R5.3 (tailor by scale)**: role enablement is tailored to slice scale (§5.2), but ⑦ independent final review **cannot be omitted**.

### 5.2 Role Enablement Matrix

**Size → type mapping (2026-09-23)**: §1.2's size S ≈ simple, M ≈ medium, M-L ≈ medium (leaning complex), L ≈ complex; "high risk" is a **separate dimension** (triggered by rework / contradictory reviews / repeated native failures) and does not follow size automatically.

| Slice type | ①architecture | ②implementation | ③testing | ⑤pre-review | ⑥scan | ⑦final review | ⑧docs a/b | ⑨native | ⑧c wrap-up |
|---|---|---|---|---|---|---|---|---|---|
| Simple (**definition same as the three skippable conditions in §7.2**, must not be defined separately in two places) | — | required | required | optional | decided per §7.2 | **required** | required | required | required |
| Medium (behavior change, multiple files, test points) | required | required | required | required | **required** | **required** | required | required | required |
| Complex (concurrency/transaction/crash window/cross-platform) | **required** | required | required | **required** | **required** | **required** | required | **required** | required |
| High risk (rework ≥2, or contradictory reviews, or repeated native failures) | **required** | required | required | **required** | **required** | **required** | required | **required** | required |

**Note (document-type slices, removing the conflict with the "simple" definition in §7.2)**: when a slice satisfies the three conditions of §7.2 (pure `.md`/comments, no code semantics), this matrix's **②implementation/③testing degenerate into the documenter's "document implementation" and "document consistency check"** (no code produced, no test engineer needed); **⑨native validation degenerates into encoding gate + G1–G7 self-check + link/reference resolvability check**. ⇒ In that form, the three columns ②/③/⑨ **are not executed under the code-slice criteria**, but **⑦ independent final review cannot be omitted**. **Blind review is triggered per §7.7** (mandatory for ownership / shutdown / identity semantics; plus 1 in 4 slices), **not by this matrix**, to avoid double counting.

### 5.3 Closure Discipline (Hard Rules)

When any stage fails or needs remediation, **the stage itself must be re-executed after the fix**; skipping steps is forbidden:

| Stage | Action after failure |
|---|---|
| ① | after revising the proposal/contract, **re-run ①** |
| ②③ | after fixing code/tests, **re-run ④ focused gates** |
| ⑤ | after eliminating deviations/counterexamples, **re-run ⑤** to confirm zero |
| ⑥ | after fixing, **re-run ⑥** until no Medium+ |
| ⑦ | after fixing, **re-run ⑦ re-review** until no Medium+ (**must not self-declare pass**); **exceptions and stop decisions are in §7.8** |
| ⑧a/⑧b | any document inconsistency returns to ⑧a/⑧b for rewrite; if the ⑧b rewrite's diff touches **semantics / criteria / conclusion wording**, ⑦ re-review must be re-run - that re-review **uses the §7.5 re-review quota**, and when the quota is exhausted ⇒ escalate to the user per §7.5; **exceptions and stop decisions are in §7.8** |
| ⑨ | fix in place and re-run; **two consecutive failures escalate to ⑤ re-check**; fall back to ① only when an architecture assumption is falsified by native evidence |
| ⑧c | a wrap-up cleanup finding ⇒ fix in place and **re-run ⑧c** until all green (frozen evidence packs run the self-consistency check per the ⑧c exemption entry in §5.1) |

---

## 6. Gates and Evidence

### 6.1 Base Gates (after every change)

| Gate | Command / criterion |
|---|---|
| Format | `gofmt -l .` output is empty |
| Build | `go build ./...` exit 0 |
| Static check | `go vet ./...` exit 0 |
| Unit tests | `go test ./...` exit 0 (record the count of ok packages and the FAIL count) |
| Contract fixtures | when contracts are involved, additionally run the Node contract check (fixture counts must stay in sync) |
| Encoding | `.md`/`.ps1` = UTF-8 **with BOM** + LF; `.json`/`.go` etc. = **without BOM** + LF |

**Focused-gate definition (2026-09-23)**: `focused gate = go test ./[affected package]/... -run [test-point prefix]` (or the equivalent contract-fixture subset); the slice definition's "steps" field **must state that command**.

### 6.2 Evidence Persistence Obligation (Hard Rule R6.1)

**All review-type subagents lack the `execute` capability** (measured: `deep-reasoner`, `independent-reviewer` have only `read`+`search`). Therefore:

- Before review, the master must persist the **raw diff** and the **raw gate output** to disk (under the repo-root `tmp/`, cleaned up after use; the path rule is §1.5), and hand the **paths** to the reviewer.
- It is **forbidden** to paste only a diff summary/paraphrase into the prompt and treat that as "evidence provided".
- Precedent: the B3c ⑦ re-review was voided for an entire round because "the raw hunk was not provided"; after switching to persisting `tmp/b3c-review.diff` it passed in one round.

### 6.3 Closing Self-check Gate (seven mechanical checks, **mandatory**)

**Location**: `tools/gates/` (under version control). **Trigger**: before closing any slice, before any docs-only commit.
**Principle**: replace "human memory" with "mechanical assertions" — this directive does not assume that any role is incapable of forgetting.

| # | Check | Criterion |
|---|---|---|
| G1 | **Status-marker consistency** | subsection heading status symbol ↔ status table row ↔ completion record row, all three must agree (e.g. B3c once had a heading still `⬜` while the status table was already ✅). **Exemption (2026-09-23)**: a stopped-state slice (§1.6) does not take part in this criterion's three-place consistency judgement; it follows the stopped-state marking rules of §1.6 (the report header reads "stopped — incomplete", the ledger keeps 🔄 or reverts to ⬜), and G7-d covers its completeness. |
| G2 | **Placeholder residue** | the full text must not retain placeholders such as "待…回填" "尚未提交" "待同轮授权" "pending the push" (after the switchover these must already have been replaced by facts). **The scan scope must exclude this line's own rule quotation** — otherwise the forbidden strings quoted on this line would be judged by the line itself as violations (**self-reference defect**; measured 2026-09-21: §6.3 self-scan yielded 4 false positives); implement it as "forbidden string inside quotes + the line contains 『不得残留』 ⇒ skip" |
| G3 | **Bilingual symmetry** | (a) **Changed-line-count symmetry**: for every changed CN/EN document pair, the added/deleted line counts from `git diff --numstat` must be equal (measured 2026-09-21: all 6 pairs in this slice PASS). It is **not required** that full-text line counts or bold-marker counts be equal — existing bilingual pairs are historically asymmetric (`REMAINING_SLICES` 501/488, `B2-EXTERNAL-ENFORCEMENT` 270/315), and "full-text equality" is unsatisfiable for them (**self-reference defect**, same as G4b). (b) **Key-field alignment (scope = changed lines, not the full text)**: among the **added lines** of `git diff -U0`, the counts of **commit hash / run number / status enumeration / conclusion enumeration** on the CN and EN sides must be equal — line-count symmetry **cannot** prove that semantics did not drift (⑦ first-round Medium finding). **Do not** compare the full text: full-text differences may come from historical legacy (measured: `B2-EXTERNAL-ENFORCEMENT` full text 10/9, while changed lines 0/0). |
| G4 | **Encoding + anomalous characters** | (a) BOM+LF verified file by file; **scope = the change set pending commit this time** (≠ the whole dirty tree). (b) **Anomalous-character scan**: the set of CJK characters in the added content ⊖ the HEAD repository-wide `.md` corpus = the **difference set**; the difference set then ⊖ `tools/gates/cjk-newwords.txt` (the legitimate-new-word allow-list, each character must come with a reason and a context excerpt) **must be empty**. ⇒ The original wording "the difference set must be empty" is unsatisfiable for any new document (measured 2026-09-21: a 51 KB new document had a 7-character difference set, all confirmed character by character as legitimate new words), so this is changed to an **allow-list** that fixes "human judgment" into a mechanically re-runnable assertion. (Origin of this rule: 「拒绝」 was once miswritten as 「拒绍」 in 5 places) |
| G5 | **Change-set contract** | the three must agree: **staged list** ↔ **change set declared in the report** ↔ **ledger evidence references**; `git add` must be done file by file (`-A` is forbidden). **Structured input** (the script reads a fixed JSON, avoiding reliance on humans reading reports): `{stagedFiles[], declaredFiles[], evidenceRefs[], addMode}`; if `addMode` is not `explicit`, or the three sets are not equal ⇒ **exit code 2** |
| G6 | **Ledger metadata consistency** | The ledger metadata of changed lines must be self-consistent: status tense (G6-1), cost clause (G6-2), run-id format (G6-3), expiring literal (G6-4), section-reference resolution (G6-5), status-word closed set (G6-6). Scope = changed lines (the added lines of `git diff -U0`); **historical rows are not retroactive** (same as G1-a / G3 / G4). |
| G7 | **Review-bundle completeness** | A changed `docs/validation/*.md` must carry a fenced block (info string `artifacts`) of `<repo-relative path><TAB><disposition>` rows, the dispositions being `repo` (may be omitted), `local-only` and `missing-historical`; every finding row needs a disposition cell and a stopped state needs `已停止——未完成`. **Structure only**: artifact content, run reachability and number truth are NOT judged. |

**Sub-rule G1-b (version metadata consistency, added 2026-09-21)**: **the version number declared in the file header must equal the version number on the latest line of the §0.3 change log**, and the header date must match the latest line's date. Mechanical criterion: extract `版本：vX.Y` from the file header, **extract `vX.Y` of the last row only from within the §0.3 table range** (**must not scan the full text** — other tables (e.g. the v3.1 comparison table in appendix C.1) also contain the `| vN.N |` form, and a full-text scan yields 5 false hits, measured), **the two must be equal** (not equal ⇒ G1 fails ⇒ stop the commit).
> **Origin**: an instance of "header v1.3 out of sync with log v1.4" appeared in this slice (⑥ round 2); ⑥'s Section D explicitly pointed out that "version metadata consistency" was previously **not within the mechanical coverage of G1–G5** ⇒ this sub-rule is exactly that missing coverage.

**Hard rule R6.2**: if any of G1–G7 fails ⇒ **stop the commit**; after fixing, re-run all gates.

**Scripting status (stated truthfully, updated 2026-09-22)**: `tools/gates/gate.js` is the single entry point (`--all` / `--check=<ID>` / `--list` / `--selftest`, `--scope=tree|index|range:<A>..<B>|ci`; exit code 0 = every executed criterion PASS, 1 = at least one FAIL, 2 = usage / environment error). **Scripted**: G1-a, G1-b, G2, G3(a.1) / G3(a.2) / G3(b) (the sub-checks `--check=G3-a.1` / `--check=G3-a.2` are individually selectable), G4a, G4b, G5-a①②, G5-b, R2.5, G6-1–G6-6, G7-a–G7-d, plus the four section 6.1 base gates (gofmt / build / vet / test). **Still not scripted** (the script prints `SKIP`, which must not be read as PASS): base **G1** (heading ↔ status table ↔ completion record: judging three human-authored places against each other) and base **G5** (the change-set contract, which needs a structured JSON input produced outside the repository). CI appends one step at the end of the `test` job that runs this script (`--selftest` first, then `--all --scope=ci`); that step **must not** carry `continue-on-error`, and a red first run must be fixed to green. This line must not be read as "fully automated".

**G5 sub-rule G5-a (role-file frontmatter hard gate, added 2026-09-21; hardened the same day per ⑦ round 3)**: the frontmatter of `prfrail-*.agent.md` **must not contain the `model` key** (writing an empty value makes the carrier **fail to load entirely**; see §2.6 / §9.3). Mechanical criteria (**both must be satisfied**): ① `grep -l "^model:" .github/agents/prfrail-*.agent.md` **must produce no output**; ② the frontmatter of **all** `.agent.md` files **must not contain an empty `model`** — the three forms `^\s*model\s*:\s*$`, `model: ""`, `model: ''` all count as empty, and all must have **0 hits** (coverage includes non-`prfrail-*` files and whitespace variants). Violation ⇒ G5 fails ⇒ **stop the commit**. It is a mandatory check item once scripted.

> **Origin**: ⑦ round 2 pointed out that "the F3 defense rests only on manual pre-check" ⇒ G5-a was added; ⑦ round 3 further pointed out that the original criteria **matched only a fixed form and covered only `prfrail-*`** (whitespace variants and other role files were missed) ⇒ this clause is hardened accordingly.

**G5 sub-rule G5-b (role-file `model` value gate, added 2026-09-21, landing the user's ruling)**: for **newly added or modified** `.github/agents/*.agent.md`, the `model` field (**if present**) **must not be a §2.4 blacklist value** (Luna / Terra / Gemini / GPT-5.4 / Sol / Kimi, case-insensitive). Mechanical criterion: for `*.agent.md` **added/modified in the change set**, extract the frontmatter `model` value; **blacklist hits must be 0**. Violation ⇒ G5 fails ⇒ **stop the commit**.
> **Origin (root fix for OB-8)**: measurement on 2026-09-21 found that the default `model` of 8 existing role files belonged to the blacklist (when `model` is not passed explicitly, the black is reachable directly via `agentName`). The user's ruling adopted the combination of "**delete the default value + hard gate**" (**do not delete the files**): first perform a **field-deletion probe** on one inactive file (delete the `model` line → invoke that carrier: **it loads and works normally**); after the probe passes, delete the `model` field from all 8 blacklist files; this G5-b clause then prevents its re-introduction.

**Sub-rule G1-a (table structural integrity, added 2026-09-21)**: **any line starting with `|` must occupy its own line**; two table records must not be concatenated into one line (affects: change log, status table, observation item ledger, and the tables in appendix C). Mechanical criteria (two independent commands):

1. `grep -n "^|.*||" <file>` against **the lines added/modified in this change set** **must produce no output** (**scope = changed lines, not the full text**, same as G3/G4; full-text comparison would report historical legacy defects as this slice's defects);
2. compare **column counts** table block by table block (before counting, **first strip escaped `\|`** — otherwise a `\|` inside a code span causes false positives (measured 2026-09-21: 4 false positives without stripping, 0 after stripping).

> **Origin**: the same "table-row concatenation" defect **appeared twice** within this slice (§0.3 change log, §12.4 observation item ledger), and both were **found independently** by ⑥ and ⑦ ⇒ the defect **cannot be reliably caught by human eyeballing** and must be mechanized (⑦ round 2 R3 and ⑥'s two rounds of conclusions).
>
> **Known pre-existing exception (not fixed in this slice)**: the M6/M7 rows at `docs/t027/A6_PILOT_LAUNCH.md:269` have the `||` concatenation **already in HEAD** (verified with `git show HEAD:...` as **not introduced by this slice**); because the SW-1 boundary requires historical drafts to **only have a header pointer added, with the body unchanged**, it is **not fixed**, merely **registered as a pre-existing defect**, and does not block this slice's commit.

**G6 sub-rule (ledger metadata consistency, added 2026-09-22)**: the four G6 sub-checks judge **changed lines only** - the **added lines** of `git diff -U0` (in tree scope an untracked file counts in full, a staged-only file is taken with `--cached`), **the same scope as G1-a / G3 / G4**; **historical rows are never retroactive** (a committed old line never condemns this slice). Target files: every `.md` in the change set.

**Self-reference guard**: lines carrying any exclusion marker (`不得残留` / `示例` / `引文` / `冲突` / `矛盾`) are skipped by every G6 sub-check - isomorphic to G2's own guard. ⇒ A sentence in this rule that quotes a forbidden shape or historical text must carry one of those markers, otherwise the gate judges the directive itself.

**G6-1 (status tense, catches the OB-11 class)**: a changed line carries a pending marker (`待` plus one of ①–⑩ / 终审 / 独立扫描 / 复审 / 预审 / 盲审 / 提交授权), and **the same block** (between two adjacent `^#{2,4} ` headings) contains a **same-label** completion marker (the 实得 parenthesised form / bold 实得 / 已闭合 / ✅ 已完成) ⇒ FAIL. The completion marker must share the label and fall inside its 40-character window, otherwise a block-level "✅ 已完成" would kill any pending line. **示例行 (example line)**: 待 ⑥ 独立扫描 coexisting with ⑥ 独立扫描（实得）PASS in one block ⇒ FAIL.
> **Origin**: OB-11's C.0e status-line tense conflict (待 ⑥ 独立扫描 vs ⑥ 独立扫描实得 PASS below); SW-5 missed 4 more of the same class before closing ⇒ human eyeballing is unreliable. **Measured evidence (引文 / quotation)**: `fb12ff9` reads 待 ⑦ 终审 plus ⑥ 实得 below ⇒ **no false positive** (no same-label completion marker).

**G6-2 (cost clause, catches the A1 class)**: the **realized clause** of a changed line (`⑥×n` / `⑦×n` tokens, with `×` / `x` / `X` all accepted as the multiplication sign, or a `= n 次` total folded in only when the line already carries a qualifying ⑥ / ⑦ token - otherwise the line is not compared) and the **same-key** realized clause found **anywhere in the file** must have identical token vectors; otherwise ⇒ FAIL (both positions and both vectors are printed). Clause qualification uses the **nearest qualifier**: the **closest** `预算` within the 60 characters before the token ⇒ budget clause (never compared); the nearest `实得` or **bare `成本`** ⇒ realized clause. ⇒ Do not degrade to a "预算 appeared somewhere in the window" containment test (that lets a preceding budget keyword shadow a following realized clause and makes the comparison vacuous). Slice keys: `SW-n` / `DR-n` / `OB-n` / `GATES-EXT` / uppercase-letter-plus-digits on the line, otherwise the same-shaped id in the enclosing block heading.
> **Basis (A1)**: the DR-FIX §0.3 v1.9 row and the appendix C.0e cost row conflict on the **same key** (`SW-5`) - 引文 (quotation): 成本 ⑥×1 + ⑦×1 versus 实得 ⑥×1 + ⑦×2 = 3 次. **Measured evidence**: `--scope=range:b738e6b^..b738e6b` must be red; `--scope=range:9f68a52^..9f68a52` (after the fix) must be green, matching the fix DR-FIX actually applied.

**G6-3 (run-id format, format only)**: the legal shape is `run` plus a space plus **9–13 digits** (this repository's run ids are 11 digits; the constant can be adjusted by a later slice and verified by the T-series tests). No reachability, no network (user ruling). **A digit run is judged only when it is NOT immediately followed by a word** - an ASCII letter or a CJK ideograph; a run followed by a word is a prose count (such as `run 1234 times` / `连续 run 12 次`), not a run id. A digit run that is **backticked, or at end of line / before punctuation** (`run 111`, `run 1234。`, `` `run 1234` ``) is still judged **only when it is also NOT followed by a word**; when it IS followed by a word it is not judged (backticks alone do not judge it: `` `run 1234 times` ``). **Every G6-3 rule (the placeholder rules included) additionally carries a leading boundary**: `run` immediately preceded by a word character is not judged (`rerun 1234`, `prerun 12`, `在run 1234`). The placeholder forms (angle brackets / word list / CJK wording) carry no following-word restriction and are checked in **ANY** context (the leading boundary still applies). **示例行 (example line; any of these shapes ⇒ FAIL - an all-zero run id fails at any length)**: `run 0` / `run 000000` / `run <RUN_ID>` / `run TBD` / `run N/A` / `run TODO` / `run 待定` / `run 占位…` / `run 7` / `run 12345678` / `run 123456789012345`.
> **Known residual false positive (registered, not fixed - G1-a precedent)**: the angle-bracket placeholder rule also catches CLI usage examples - two instances exist: `prfrail run <chain-file>` at `docs/RFC-proofrail-unattended-ai-engineering-product.md:737` and `:1085`; because only changed lines are judged this **does not affect the gate today**, but either line would produce a spurious red as soon as it is touched. **No whitelist entry was added** (a whitelist entry would mask this shape); registered as **OB-17**. This line carries the `示例` marker so that the quoted shape does not trip the rule (see OB-18).

**G6-4 (expiring literal)**: a changed line carrying a line-count / file-count / aggregate-count literal (such as `NN/NN 行`, `NN 行`, `NN 个文件`, `共 NN`) ⇒ FAIL; the whitelist may exempt them one by one (reason and context excerpt mandatory). **示例行 (example line)**: 995/995 行、共 128 个文件。
> **Origin (引文 / quotation)**: the OB-11 lesson - a ledger must not carry expiring literals (line counts / file counts), only criteria such as "verified consistent by G3/G4" (v1.9 removed 995/995 accordingly).

**G6-5 (section-reference resolution, added 2026-09-23, landing OB-32)**: on the added lines of the change set that sit inside the §0.3 change-log table interval or the appendix C.0–C.0i ledger interval, every §X.Y reference must **exist exactly** in the heading set of **this file and its `_EN` mirror** (`^#{2,4} ` headings + leading bold labels such as §1.7.1 / §3.4.1), or resolve via a **parent-section prefix** (§7.2 resolves through §7.2.1). **Approximate-criterion statement**: this sub-check covers only the "dangling reference" subclass and **does not cover** C1's full class of "the log claims a change while the body has no such content" (content-level reconciliation stays with ⑤ ⑥ ⑦ and G6-1 / G6-2, see OB-32 in §12.4); the v3.1 comparison table in Appendix C.1 and C.2 / C.3 are **outside the target region**; an external document's section number is exempted by an exclusion-marked line or by `g6-whitelist.txt`. **示例行 (example line, any of these shapes ⇒ FAIL)**: `| v1.13 | … | 新增 §7.9 评审收敛判据 | … |`.

**G6-6 (ledger status closed set, added 2026-09-23, landing OB-38)**: the disposition cell's **first status word** of a changed row inside the §12.4 table interval (a `| OB-N |` line) must be ∈ the closed set `观察中 / 待办 / 待议 / 已处置 / 已裁决 / 已登记` (the `_EN` mirror uses `Observing / To do / To be discussed / Resolved / Ruling / Registered`; **word** may be followed by a parenthetical note). **示例行 (example line, any of these shapes ⇒ FAIL)**: `| OB-32 | … | **下一片必做**（…）|`. **由来 (引文 / quotation)**: P5 measured that §12.4's status vocabulary and G6-1's marker vocabulary barely intersect ⇒ ledger status conflicts previously had no mechanical coverage; this sub-check is independent of G6-1 (block-level semantics are incompatible, **no widening** of G6-1, argued in §7.8 and this subsection).

**G7 sub-rule (review-bundle completeness, added 2026-09-22)**: target files = `docs/validation/*.md` in the change set (**excluding** `docs/validation/evidence/**`); the rule text lives inside the directive, outside the scan domain ⇒ the self-reference trap is **structurally immune**. **Structure only**: artifact content, run reachability, number truth and report-conclusion truth are NOT judged.
**Machine contract**: a changed validation report must carry a fenced block (triple-backtick fence, info string `artifacts`) of `<repo-relative path><TAB><disposition>` rows; the dispositions are `repo` (may be omitted, the default), `local-only` and `missing-historical`. **Self-compliance**: this slice's own validation report (`docs/validation/gates-ext{,_EN}.md`) must carry that block.
**Sub-checks**: G7-a a missing fence ⇒ FAIL; G7-b an unresolvable disposition / path (`repo` or omitted, the same meaning in every scope: `git cat-file -e <endpoint>:<path>` for range / ci; for tree / index the path must be staged in the index (`git ls-files --error-unmatch`) AND present on disk ⇒ a tree-scope run is a faithful pre-flight only after `git add` of the declared artifact) or an unknown disposition token ⇒ FAIL, while `local-only` and `missing-historical` are accepted as declared; G7-c a changed finding row (line starting with `|` plus a bold Critical/High/Medium/Low) whose **last cell** (the disposition cell) is empty or carries no disposition token ⇒ FAIL (structural only: the token must sit in the row's **last cell**; a disposition word inside the description cell does not count - `| **High** | ledger says “Fixed” but the disposition cell is empty | |` FAILs; the accepted token set is the CN / EN union - `已修` / `已整改` / `已闭合` / `已处置` / `接受` / `拒收` / `记录在案` / `待办` / `观察中` / `待议` / `已停止`, plus `Fixed` / `Resolved` / `Closed` / `Accepted` / `Rejected` / `Deferred` / `Recorded` / `To do` / `Observed` / `To be discussed` / `Not a problem` / `Documented`); G7-d a changed status line on an added line (`状态` or `Status` followed by a separator) declaring a stopped state (the token set is the CN / EN union, with ASCII tokens matched case-insensitively: `已停止` / `停止态` / `Stopped` / `Halted`) without an incomplete mark (`已停止——未完成` / `Stopped — incomplete`) on the same line or in the same block ⇒ FAIL (added lines only, never the whole report; a prose line that merely names the rule is not a status line, avoiding a self-reference false positive).
> **Basis**: C1 of DR-FIX (a report conclusion with no artifact path ⇒ broken evidence chain) is caught by G7-a / G7-b; OB-10 ① (the stopped-state report header) **evaluation: adopted** as G7-d; OB-10 ② (the pilot-metrics slot of the kickoff message) **not adopted** (it exists only in session text and cannot be checked mechanically ⇒ OB-10 stays observing).

> **The four reports that predate this contract are not retroactively failed**: `docs/validation/dr-fix-enforcement-proxy{,_EN}.md` and `docs/validation/sw5-work-tiers{,_EN}.md` were committed before this contract existed ⇒ the criterion only looks at the change set, so it is **not retroactive** and **does not backfill** (they are history, see §11.1); the fence block is required only when they are changed next.
**Fixture convention**: `tools/gates/testdata/fixtures/*.txt` are mapped to virtual `docs/**/*.md` paths by `selftest.js` ⇒ the criteria still select `.md` targets, yet the tool **never** trips on its own deliberate violations; fixtures are deliberately not stored with a `.md` suffix and an assertion in `selftest.js` keeps that convention.
**Two whitelists**: `tools/gates/g6-whitelist.txt` (three mandatory fields: match string / file / reason plus context excerpt; one entry exempts one finding, a crippled entry ⇒ exit code 2) and `tools/gates/cjk-newwords.txt` (used by G4b; its role is unchanged and every character is still registered one by one). **Loader-source rule (GATES-EXT F1)**: `tree` / `index` read the working tree; `range` / `ci` read the scope's endpoint revision (`ci` means HEAD) - both whitelists share one loader, and each criterion prints the source it used in its `detail`.

### 6.4 CI Observation

After pushing, GitHub Actions must be observed and **reported back as double green**. A failed first run requires: ① obtain the raw log of the failing step; ② `gh run rerun --failed` to re-run and collect evidence; ③ determine whether it is flaky; ④ **register it as a defect item** (DR-N) stating the run number, failing step, failing test, error message, and re-run result. It is **forbidden** to cover up a red first run with "the re-run went green".

---

## 7. Review system

### 7.1 Three-layer review

| Layer | Role | Positioning | Key requirement |
|---|---|---|---|
| ⑤ | Pre-reviewer (V4 Pro) | whether the implementation deviates from the plan | covers the plan item by item + the three classes of counter-examples |
| ⑥ | Independent scan (MAI) | low-cost supplementary scan | **Does not reference anyone else's checklist**; Section A–E |
| ⑦ | Independent final review (Codex) | the sole final-review gate | **The first round attaches no ⑤/⑥ checklist**; four-item review |

> **Note (⑧c is not one of the three review layers)**: the ⑧c wrap-up cleanup is executed **by the master** and **does not belong to** the ⑤⑥⑦ review layers - it is a fixed stage of §5.1 (see §5.1 / §3.4), parallel to ⑨ native validation.

### 7.2 ⑥ (independent scan layer)

### 7.2.1 Trigger

**Trigger**:
- **Behavior-changing slice** (contract / state machine / crash path / process governance / ownership semantics changed) → **enabled**;
- **Skippable** only when the following three conditions are **simultaneously** met: ① the change **lands only in documentation and configuration-wording files** (`.md`/comments/wording; non-machine-structural files such as `.gitignore`, `.txt`, `.html` count as this class); ② it contains no semantic change to `.go`/`.json`/`.yml`/schema/fixtures/**script logic**; ③ it **does not touch the semantics of the state machine, gate criteria, evidence model, or role/authorization contracts**. **Skipping must be confirmed by ⑦ or the user** (the master must not unilaterally declare it "skippable"), and the reason must be recorded in the report.
    - **Counter-example (measured 2026-09-21, must not be repeated)**: SW-1/SW-2 look like a "documentation slice" on the surface, but it **revised the §6.3 gate criteria and the evidence model** ⇒ **③ does not hold, ⑥ must not be skipped**. The master's initial skip reason **does not hold**, and after ⑦ pointed this out in round 2 **⑥ was re-run to make up for it** (result in validation report §7.2b). ⇒ Judging "skippable" requires **checking ①②③ item by item**, not inferring from "it looks like documentation".

### 7.2.2 Output contract, cost cap and failure disposition of the ⑥ layer

**⑥ Output contract (Section A–E, isomorphic to ⑦'s contract)**:

- **Section A**: overall verdict line —— `INDEPENDENT SCAN: PASS` or `INDEPENDENT SCAN: FINDINGS`.
- **Section B**: findings list, each line in the fixed format `SEVERITY | file:line | violated sentence (verbatim quote) | minimal fix`.
- **Section C**: item-by-item evidence and reasoning.
- **Section D**: **falsifiability audit (mandatory)** —— for each item, state "if guard X were deleted, which check/assertion would go red and which would not", and list the invariants it considers **not covered**.
- **Section E**: **"the parts I cannot verify" + "what would overturn my PASS"** (anti-rubber-stamp).

> **Inherited source and sole authority**: the A–E definitions come from v3.1 §2.2 (original text in `docs/t027/FLASH_OPERATING_DIRECTIVE_v3.1.md`). The three references to `Section A–E` in this directive §4.2 / §7.2 / appendix B.6 **take this definition as the sole authority** (completed 2026-09-21: previously those three places referenced it with **no definition**, which was a **reference break**).

**Cost cap**: effective MAI calls for one slice (including first round, resend, re-run after remediation) ≤ **3**; resends are limited to "truncation/anomaly" and count toward the cap; Haiku is counted separately, ≤1.
**Failure handling**: MAI reports Medium+ ⇒ **re-run MAI** after fixing; if MAI's output **deviates in format** (not following Section A–E / no verdict line) ⇒ one resend is allowed; **if it deviates again or the cap is reached ⇒ this layer must not declare a pass**: it must ① truthfully leave a trace and record "⑥ produced no actionable findings"; ② **enter ⑦ independent final review and close the loop per ⑦'s verdict**; ③ escalate to user adjudication when necessary; afterwards ⑨ may still add a "dedicated independent corroboration experiment" to take on the A12 role. **This layer's verdict does not constitute proof of review sufficiency.**

### 7.2.3 Prerequisites for enabling Haiku

It must be **confirmed before it is enabled**: ① `Claude Haiku 4.5` exists in the Copilot model list; ② the access channel (proxy/network egress) is normal and a first real call completes the availability pre-check. If either is unmet ⇒ **handle it by a three-level fallback**: ① fall back to ⑤ and ask V4 Pro for one extra scan round (one extra pre-review round); ② if ⑤ has already converged, then **skip ⑥ and record the reason in the report**, with ⑦ taking over the independent-scan duty; ③ if ⑦ also cannot cover it, then **escalate to user adjudication**; **do not retry repeatedly**, and do not count its absence as a ⑥ failure.

### 7.3 Input sanitization and landing of conclusions

- **Input sanitization**: material handed to the reviewing party must filter out IPs, usernames, hostnames, SSH paths, credentials, and tokens.
- **Verdict landing check**: before landing, every `file:line` given by the reviewing party must be **re-checked item by item** for whether it really exists and matches the criterion; anything that cannot be located is treated as a "lead" and must not be taken at face value.
- **Hard anchoring isolation**: in the first round, ⑦ is given only the slice definition, the diff, the verbatim contract sentences, and read-only constraints; only the **re-review round** may attach a checklist, marked "read only after completing the independent scan".

### 7.4 Getting it right the first time (fewer ⑦ rejections)

Before implementation the following are mandatory: ① the plan's change points are precise down to the file; ② the test points are executable and include counter-examples; ③ contract-first has been completed; ④ the gates are fully green locally. **The cost of a rejection is far higher than the cost of upfront clarification.**

### 7.5 Cost caps (hard metrics)

| Role | Per-slice cap | Remarks |
|---|---|---|
| Architect | ≤ 1 | complex/high-risk slices |
| Pre-reviewer | ≤ 1 (re-review after remediation +1) | excludes the ⑦ blind-review add-on |
| ⑥ MAI | ≤ 3 (effective calls) | see §7.2 |
| ⑦ Codex | 1 (final review) + 1 (re-review) | re-review after remediation is an obligation; the "mandatory-sampled blind review" is counted separately as +1 and may be merged with the latter but **must be judged and its reason recorded before start-up**; **exceptions and stop decisions are in §7.8** |

Exceeding a cap ⇒ escalate to user adjudication. **Cost-gradient note** (measured pricing, USD per million tokens, offPeak input/output): Flash `0.15/0.6`, V4 Pro `0.66/1.98` ⇒ **prefer pushing coding/testing/documentation/gates down to Flash**.

**Hard rule R7.1 (source of truth for metering and de-duplication)**: the counts in the table above may only come from a **single source of truth** —— the **session call receipts/logs** (not "I guess I called it a few times"). Every call must be identifiable as `{role, model, stage, startedAt, reason, retryOf}`; **de-duplication rules**: a **resend** of the same `stage` (truncation/anomaly) **counts toward the cap**; a **re-review after remediation** is counted on its own line; **blind review** and final review are listed separately. The report's "cost and metering" must list items by these fields; verbal counts from outside the session logs **do not count**.

### 7.6 Pre-authorized probe quota

Slices that need real external calls must declare a quota cap in advance; when the quota is exhausted, stop, and **must not** add more within the slice. Probe usage must be accounted for (count, number of paid calls, conclusion).

### 7.7 Residual governance (A11 / A12)

- **A11 symmetric spot-check**: the master **randomly samples 1–2 items from each of ⑥'s and ⑦'s "delete guard X ⇒ test Y goes red" claims and actually executes them**; if a claim does not match reality ⇒ that audit's credibility is downgraded and logged as a negative signal. A spot-check only provides **existence corroboration** and does not constitute proof of coverage.
- **A12 independent corroboration**: when ⑥ reports PASS, ⑨ must have at least one **dedicated experiment** that independently proves the key invariant (a criterion unrelated to ⑥).
- **Mandatory-sampled blind review**: when a slice touches **ownership / shutdown / identity** semantics, **1 blind review** is added before ⑦ (with no checklist attached). Its call count and the convention for "whether it is merged with the §7.2 fallback" are **uniformly specified by §7.5**; this section does not accumulate separately, avoiding double counting in two places.
    - **Proportional spot-check (inheriting v3.1 §3.11)**: besides the "mandatory sampling", **1 slice out of every 4** undergoes one ⑦ blind review (choose the one with the broadest reach among the most recent 4 slices); the sampling record (slice number / number of blind-review findings / whether ⑥'s PASS was overturned) is written into the report.

### 7.8 Review Saturation Criterion (when reviewing may stop)

> **Positioning**: this section is the **sole authority** for the three stop objects (third-party independent re-review / ⑦ round N / re-review after remediation); §5.3 / §7.5 / §12.1 / §12.2 / §12.4 only cite it and never define it separately (one authority, the same hard constraint as Appendix B.7).

**7.8.1 Stop objects and the unit of judgement**: three objects - a **third-party independent re-review** (the DIRECTIVE REVIEW round of a directive-governance slice), **⑦ round N** (the re-review chain after the final review), and a **re-review after remediation** (re-judging a fix diff). The **unit of judgement** = one slice × one review loop; **stopping** = that object is not re-run this round and its conclusion must be labelled truthfully (**not run ≠ PASS**; §7.2's reading "that layer must not declare a pass" is unchanged). ⑥'s cap of three rounds already has a failure disposition in §7.2 and is only cited here.

**7.8.2 Convergence signals (judgement input; each with its test)**:

- **Severity convergence**: no Medium+ this round (§4.3) - mechanically checkable; the verdict line ∈ the closed-set enum of §4.1 is covered by the §4.1 contract and human verification (not G7-c); G7-c checks only the disposition column of finding rows.
- **Finding-rate convergence**: this round's new findings are fewer than the previous round's (per item, same root cause merged, denominator reading in §12.1) - declared by ⑦ in the re-review conclusion.
- **Repetition-rate convergence**: most new findings are recurrences of a registered OB ⇒ recorded as "recurrence" rather than new findings - declared by ⑦.
- **Coverage-gap convergence**: the "parts I cannot verify" of Section E are covered by a mechanical criterion or a registered OB - declared by ⑦.
- **Quota and exceptions**: see the decision table in 7.8.3.

**7.8.3 Stop decision table**: the four stop situations and the pilot deferral are all decided by the table below.

| Situation | Stop objects it applies to | May stop | Basis section | Record-keeping | May serve as a precedent |
|---|---|---|---|---|---|
| Normal convergence (no Medium+ this round, legal verdict line) | re-review after remediation / ⑦ round N / third-party independent re-review | **May stop** (that object closes) | §4.3 + §5.3 | verdict line verbatim into the report; the re-review report into the review input package (G7 fence block) | yes (the normal flow) |
| Quota exhausted (the §7.5 caps; ⑥'s three rounds in §7.2) | ⑦ round N / third-party independent re-review | **May stop** (no self-funded extra round) | §7.5 + §1.7.2 + §7.2 | unclosed findings registered in §12.4 (`待办：并入下次 ⑦ 输入`); R7.1 metered line by line | no (counted per slice) |
| User-authorized exception (pre-authorized extra round / narrowed input) | ⑦ round N | **May continue** (only within the authorized scope and rounds) | the sentence after the table in §7.5 ("Exceeding a cap ⇒ escalate to user adjudication") + the authorization text | the authorization text + the input-narrowing statement into §12.4 / Appendix C | this slice only (`GATES-EXT` writes explicitly "the authorization does not generalize") |
| OB-11 exception (the ledger-metadata class with all five elements) | re-review after remediation | **May stop** (folded into the next ⑦ input) | §7.8.4 (sole authority) + OB-11 in §12.4 | the OB disposition cell is marked `待办：并入下次 ⑦ 输入` | may serve as a precedent (the same nature only) |
| OB-37 meta-rule (a direct user ruling breaks the literal OB-11 standard) | re-review after remediation | **May stop** (direct user ruling) | §7.8.4 + OB-37 in §12.4 | an independent OB or ADR entry (who ruled / date / reason / basis) + a truthful label | **must not serve as a precedent** |
| Pilot deferral (the two companions in §12.1 + the second pre-ruling in OB-40) | the pilot evaluation (not one of the three review objects) | **Deferral** (not a stop); a second deferral ⇒ escalate for a user ruling | §12.1 + §12.2 + §7.8.6 | "pilot deferred" + the deferral count in the ledger | the second one is a direct user ruling (OB-40); from the third on it must not continue as a precedent |

**7.8.4 Exceptions and meta-rules (sole authority)**:

- **OB-11 exception standard (five elements, restated here as the sole authoritative text)**: a ledger-metadata-class remediation must satisfy all of - **no rule entity**, **no gate impact**, **closure literally verifiable**, **no new semantics**, **ledger-metadata class**; OB-11's disposition cell in §12.4 points here.
- **OB-37 meta-rule (restated here)**: a direct user ruling **may** break the OB-11 exception, but it must ① **be recorded independently** (as an OB or an ADR, never only in a note under §0.3 / Appendix C); ② **be labelled truthfully** as "beyond the literal standard of the OB-11 exception, by direct user ruling"; ③ **not serve as a precedent**. OB-37's disposition cell in §12.4 points here.
- **Stated explicitly**: the OB-11 exception **must not** be used for a remediation that "carries rule-semantics changes" (OB-37 has demonstrated its lapse).

**7.8.5 Record-keeping**: every stop decision must land five fields - {object, round, basis section, verdict line verbatim, quota count (R7.1 field by field)}. **Where it is registered**: the quota class ⇒ §12.4 (`待办：并入下次 ⑦ 输入`); the user-ruling class ⇒ an independent OB or ADR entry. **Mechanical backstops**: the verdict line ∈ the closed-set enum of §4.1 is covered by the §4.1 contract and human verification (not G7-c); G7-c checks only the disposition column of finding rows; the §12.4 status word ∈ the closed set (G6-6); a log / ledger § reference is resolvable (G6-5); the realized cost agrees (G6-2).

**7.8.6 Interface with the pilot deferral**: the deferral clause and its two companions are in §12.1; the dual-track priority rule is in §12.2; the second deferral's pre-ruling is in OB-40 of §12.4. This section only judges "whether the deferral may continue": fewer than two deferrals follows companion ① of §12.1; reaching two follows companion ② and must escalate for a user ruling; the second one has been ruled directly by the user (OB-40) and **must not serve as a precedent for the third or later**.

---

## 8. Effort tier (thinking mode) rules

### 8.1 Basic facts (measured)

- The effort tier is an **API-layer parameter**; a model **cannot perceive, confirm, or change** its own tier, and a subagent's self-reported tier is also untrustworthy (measured `effort=unknown`).
- The tier is **set by the operator before the run starts**; **it cannot be changed dynamically during continuous operation**.
- `deepseek-flash` (used by both the master and the product layer) supports `low / high / max`, default `high`.

### 8.2 Tier setting (by slice type, **set once before start-up**)

| Slice type | Master | Architect/Pre-reviewer | Product layer (implementation/testing/documentation) | ⑥ / ⑦ |
|---|---|---|---|---|
| Simple | High | — | same tier as the master (see §8.3) | each model's own default/fixed tier |
| Medium | High | Max | same as above | same as above |
| Complex | High | Max | same as above | same as above |
| High risk | **Max** | Max | same as above | same as above |
| Routine [ROUTINE] | inherits the current session tier (not set from this table) | — | — | — |
| ⑦ Codex | — | — | — | **Extra High (fixed)** |

### 8.3 Tiers of same-model roles cannot be separated (**known limitation**, D2=A)

The master and the product layer are **both `deepseek-flash`** ⇒ their tiers are **the same setting**, and cannot be set separately to "master High / implementer Low" during continuous operation.

- **Accept the same tier**: by default everything takes **High**.
- **When Max is needed**: ① configure the Flash tier to Max **before the task starts**; or ② during the task, **pause → change tier → continue** (the gate must be re-run before and after the pause to confirm state).
- **Long-term solution**: configure a **different model** for the product layer (model per role); the limitation is then lifted.

### 8.4 Recording obligation for tier changes

Every tier change must be recorded in the report's "environment and tier" table: time of setting, tier value, whether a pause was involved, and the gate results before and after the pause.

---

## 9. Master fallback and exception handling

### 9.1 Fallback trigger scenarios

| Scenario | Master action |
|---|---|
| Subagent output is empty / clearly unreasonable | **retry once** after re-trimming the context; if it still fails ⇒ log an exception, switch to a backup carrier or downgrade the process step |
| Subagent fails twice in a row | pause that loop, output a diagnostic report, **mark as requiring human intervention**. **Exception**: resends of the ⑥ independent-scan layer follow the §7.2 cost cap (≤ 3 effective calls), and count towards this row's consecutive-failure count only once the cap is reached |
| ⑤/⑥/⑦ verdicts contradict each other | the master does **root-cause arbitration**: reproduce the minimal counter-example → decide who is wrong → decide the fix path → write it into the report |
| Gate anomaly that is not a code problem | check environment/dependencies/configuration; rebuild the environment if necessary (must not modify production code to mask it) |
| Ambiguous boundary / unmet dependency | trace back up to the slice definition, **mark "dependency blocked" and stop**; do not infer on your own |
| Native validation fails with an unknown cause | peel off layer by layer to locate the root cause → assign the corresponding subagent to fix it |

### 9.2 Fallback principles

1. **Fallback does not mean writing on their behalf**: after locating the problem, still assign the product-layer role to fix it; the master only does mechanical integration and gates.
2. **Every fallback must be recorded in the report's "exception handling" section** (scenario, localization process, assignment path, result).
3. **Threshold merging**: fallback intervention in the same slice **> 3 times** ⇒ automatically mark it a "high-risk slice" and **pause for human review**; it **shares the same counter** as §5.3's "⑨ failing twice in a row escalates to ⑤", and must not form two parallel ladders. **Counting rule (2026-09-23)**: the "fallback count" = each **individual action** the master takes to localise / fix / re-run; a fix-and-rerun of the same root cause in the same loop counts once.

### 9.3 Known anomaly patterns (prior checklist, for fast localization)

| Pattern | Symptom | Handling |
|---|---|---|
| Stale SCM display | the editor shows "untracked" but git is clean | trust the four commands such as `git status --porcelain -uall`; refresh/restart the editor |
| Phantom terminal interrupt | the command echoes `^C` but it has actually executed | verify with `git log`/file state; do not resend destructive commands |
| Character-level mis-write | Chinese characters are written as look-alike characters | run the G4 anomalous-character scan |
| Loose tool matching | the replace tool still reports success for a look-alike oldString | always run mechanical assertions after the change; do not trust the "success" echo |
| Encoding defaults | a newly created `.md` has no BOM / has CRLF | run the encoding gate immediately after writing to disk and correct it |
| **The file-writing tool emits CRLF** | a newly created file written to disk by the file-writing tool has line endings `0d 0a`, violating the hard LF rule | measured by the SW-2 probe (31 B including CRLF). **A new file has no existing line endings to infer from ⇒ it is always CRLF**; existing files usually keep their original line endings. ⇒ after any "new file" is written to disk, its line endings must be normalized with an explicit `UTF8Encoding($false)` (no BOM) or `UTF8Encoding($true)` (`.md` with BOM) and read back for verification. **Supplementary test (2026-09-21, probes 4/4)**: the file-writing tool is **equally unfaithful** for `.md` —— it produces **no BOM + CRLF** (54 B, `first3=23 20 50`, `0D×2`); after byte-level correction to BOM+LF, 55 B passes all items. ⇒ **after a `.md` is written to disk you must explicitly add the BOM and remove the CRs**, and must not rely on the tool or the editor to handle it automatically |

| **`model: ""` makes the carrier fail to load** | the role file contains an empty `model:` value ⇒ `runSubagent` reports `Requested agent '...' not found` | measured 2026-09-21 (isolated verification: adding the empty-value key ⇒ failure; deleting the key ⇒ recovery). ⇒ role files **omit** the `model` key (§2.6) |
| **A post-freeze append drifts the criteria's scope** | appending a "freeze ruling / CI result" record to an already-committed revision leaves the index without the original slice's commit ⇒ ⑨ and G6's changed-line scope **degenerates from the whole change set to this append**; with an empty change set there is not even a line to judge | a record's readings **must be bound to the commit they were taken on** (keep the readings file and the commit id together under `tmp/` and state the reading in the report); if the body is already frozen, **do not change** it and store this round's readings separately - landing OB-41, see G6's scope in §6.3 |
| **A PowerShell pipe capturing a native command's UTF-8 stdout ⇒ mojibake in evidence files** | decoding through the console code page (GBK) and re-encoding as UTF-8 ⇒ Chinese text and ⑥/⑦-class symbols turn into unrelated characters (GBK decode traces) | always write evidence files straight from node (`execFileSync` + `fs.writeFileSync('utf8')`), or use a bare `>` redirect (a native command's stdout passes through byte for byte); re-check every file for mojibake markers = 0 after writing |

**Rule R9.1 (the failure-mode library keeps growing)**: §9.3 is a **continuously growing library of known failure modes** —— for every new, reproducible anomaly pattern met in the pilot or in later slices, **one line must be appended** (symptom / handling) when the slice is closed out; otherwise the documentation write-back is deemed incomplete (detectable via G2/G5 in §6.3).

---

## 10. Authorization, Git, and Release Discipline

### 10.1 Authorization

- `git commit` / `git push` **require** explicit authorization in the same turn; **never** commit just because "it is convenient at wrap-up".
- Authorization scope follows the user's wording strictly: if only commit is authorized, do not push; if only origin is authorized, do not touch gitee.

### 10.2 Push

- **Push only to `origin`**; **gitee is a mirror only — without an explicit request in the same turn, any gitee push is forbidden**.
- After authorization, push and **observe CI** (§6.4).

### 10.3 Staging

- Stage precisely with `git add <specific files>`; **`git add -A` is forbidden**.
- Commit message style `<type>: <summary>` (chore/feat/fix/docs/test/build); bilingual commit messages use English.

### 10.4 Stop Point

After the pipeline finishes ⑨ it **stops at ⑩**: output the slice report and the status table, and wait for user authorization to commit. **Do not** commit/push/publish on your own.

---

## 11. Documentation and Ledger

### 11.1 Domain Authority and Bilingual Mirrors

- Domain authority is defined in §0.1; Chinese is authoritative and `_EN` is the mirror; the mirror must be **aligned sentence by sentence** (not merely keyword by keyword).
- Historical versions (e.g. the v3.1 directive, frozen validation reports) are **not rewritten**; only a pointer or a correction note is added.

### 11.2 Write-Back Standard

At slice wrap-up the following must be written back: the validation report (created or appended), the `DEV_PLAN{,_EN}` section, the slice ledger (status line / checkboxes / completion record), and the ADR row (if a ruling is involved). Write-back content must be **taken from the slice's own evidence**, never from memory; missing fields are marked "historically missing".

### 11.3 Report Structure

Change summary → execution pipeline (step/role/status/duration/rework count) → anomaly and fallback record → environment and effort tier → review conclusion summary (⑤⑥⑦) → **review input package summary (including the blind-review isolation proof: whether a checklist was attached in the first round, whether the input was already sanitized)** → falsifiability (mechanical checks + mutation testing) → gate results → ⑧c wrap-up cleanup record (§5.1) → cost and metering (call counts by type, probe quota) → explicitly unexecuted items → next-step recommendations (including effort-tier recommendations for each role).

---

## 12. Validation and Evolution of This Directive

### 12.1 Pilot Requirements

**Pilot slice selection criteria**: prefer a **medium-complexity, behavior-changing** slice (real gates, real review, controllable cost); **avoid** making the very first pilot a high-risk slice or a documentation-only slice. Candidates for T027: the **DR-2 + DR-3 fix slice** (recommended: medium, behavior-changing, complete review surface); if **B4 · AT-23** (complex) is chosen instead, the budget must be doubled and more fallbacks must be expected.

After a new directive or a major revision, **the first slice should be the pilot**, and quantitative metrics are declared in advance; if the first slice starts as `[SLICE]`, this section's **deferral clause** applies (the pilot evaluation defers to the next slice, and "pilot deferred" is registered in the ledger):

**Baseline snapshot = B3c (2026-09-21, documentation-only slice)**; on user request it may be replaced before the pilot by the snapshot of another slice of the same scale.

| Metric | v3.1 baseline snapshot (measured on B3c) | Pilot target |
|---|---|---|
| Defect self-capture rate (found by us ÷ total found) | **1/4** (3 of the 4 omission classes were found by the user/reviewer) | **≥ 3/4** (with master self-check + subagent review as the primary discovery channels) |
| Rework count (reruns per loop) | ③×3, ④×3 (v3.1 numbering, of which **1 round was voided for "raw diff not provided"**) | ③ ≤ 2, ④ ≤ 2 |
| ⑦ invalid rounds (voided for evidence/format problems) | 1 | **0** |
| Master fallback count (**i.e. the fallback count of §9.2; this item uses the same definition as §9.2**) | Not counted separately (did not occur on B3c) | ≤ 3 (exceeding it triggers the §9.2 high-risk flag) |
| Omissions found by the user | **3** | **0** (caught by the §6.3 self-check) |

**Denominator definition (user ruling 2026-09-23, fixed)**: `total found = found by us + found by reviewers + found by the user`, counted **per item**; **the same root cause merges into 1 item** (the same problem in CN and `_EN` counts as 1). `Found by us` = problems the master's self-check / the gates / the ①②③④⑧ roles reported and fixed before the commit; `found by reviewers` = findings by ⑤⑥⑦ and by third-party independent review; `found by the user` = problems the user pointed out that we did not report first. **Missing any single target makes the pilot a failure.**

**Reading note (2026-09-23)**: the 1/4 baseline above uses the "omission-class" reading, which is **not the same denominator** as this section's fixed per-item reading and therefore cannot be compared directly - it is a historical reference only; the ≥ 3/4 target and the measured 2/5 are compared on the **same denominator**.

**Handling when the slice starts as `[SLICE]` (2026-09-23)**: if the first slice of the pilot period starts as `[SLICE]` (rather than `[PILOT]`), the master must flag "this slice should be a pilot per §12.1" in the opening receipt and stop per §1.6 to await the user's confirmation; if the user insists on `[SLICE]`, the validation report must record in writing "no pilot metrics were declared" with the reason, and that slice does not count towards the pilot evaluation; **the pilot evaluation defers to the next slice** (user ruling 2026-09-23): that next slice must start as a `[PILOT]` and **declare the pilot metrics itself** (subject = the user; object = the pilot-metric declaration; manner = the next slice starts as `[PILOT]`), and "pilot deferred" must be registered in that slice's ledger; otherwise the pilot is never executed. **Two companion clauses to the deferral (user ruling 2026-09-23)**: ① the deferral target must **meet this section's pilot-selection criteria** (it must not land on a high-risk slice or a documentation-only slice); if it does not, the deferral **continues**, with the **deferral count** registered in the ledger; ② **two or more deferrals** ⇒ the master must stop per §1.7.2 and **escalate for a user ruling**, and must not keep deferring. **Stop / deferral / exception decisions are unified in §7.8.**

### 12.2 Rollback Clause

**v1.0 has no project-level previous version to roll back to** (v3.1 covers only the T027 slice; it is not a repository-level directive). Rollback path when the pilot misses its targets: ① **suspend this directive**; ② T027-related slices **resume execution under `FLASH_OPERATING_DIRECTIVE_v3.1`**; ③ all other repository activity is **suspended**, pending user adjudication. After a revision the pilot must be run again. **Until the pilot passes, this directive must not be claimed to be "validated".**

**Re-evaluation clause (user ruling 2026-09-23)**: when the pilot misses its targets, the path for a user ruling of "continue" is: ① register the ruling in §12.4 (who ruled / date / reason / re-evaluation point); ② extend the criteria or fix the root cause as registered; ③ recompute the metric with the fixed denominator at the **re-evaluation point** - if it still misses, the rollback path of this section is enforced with no second exemption. **Re-evaluation object (user ruling 2026-09-23)**: the object is **DR-FIX itself**, but its metric is recomputed by "**how many findings the extended criteria could have caught**" - the test being whether they could catch the **classifiable** findings among those missed at the time (OB-14 records **2 items**; the rest fall outside `G6` / `G7`'s coverage), and those count into **findings by us** before the self-capture rate is recomputed. **Judge and recomputable anchor (user ruling 2026-09-23)**: the judge is **⑦** (deciding inside the re-evaluation input package and recorded in the report); the anchor is running `gate.js --scope=range:<DR-FIX commit>^..<DR-FIX commit>` over G6 / G7 for DR-FIX's commit, where **a red verdict counts as "catchable"** (so the call is never subjective). ⇒ This section **does not create** the reading "the directive is necessarily suspended at the re-evaluation point"; only a recomputation that still misses triggers the rollback path of this section.

**Dual-track priority rule (user ruling 2026-09-23, P1)**: the **new pilot**'s result after a deferral **outranks** DR-FIX's historical re-evaluation. ① If the new pilot **meets every** §12.1 metric at its wrap-up ⇒ this directive **counts as having passed the pilot** and DR-FIX's re-evaluation **becomes a record item** that no longer triggers this section's rollback; ② if the new pilot **misses** ⇒ this section's rollback path runs and DR-FIX's re-evaluation **merges into that same judgement**; ③ if the re-evaluation point falls before the deferred new pilot wraps up ⇒ the re-evaluation is **merged into the single judgement at that wrap-up** (the point is not cancelled, only merged), removing the timing inversion in which "the re-evaluation arrives first and the deferral loses its point".

**This ruling record (DR-FIX pilot)**: the DR-FIX pilot's self-capture rate was **2/5** (missed, see OB-14); the user ruled on **2026-09-23** to **continue** - reason: the root cause is clear (no mechanical criteria for the ledger-metadata and review-bundle-completeness classes), and it has been addressed by `GATES-EXT`'s G6 / G7; registered per ① in the **OB-14** disposition cell of §12.4; re-evaluation point = **whichever of `GATES-G1G5` / `GATES-TIGHTEN` completes first**; **the re-evaluation object is DR-FIX itself** (recomputed by "how many findings the extended criteria could have caught", see the clause above). **Stop / deferral / exception decisions are unified in §7.8.**

### 12.3 Observation Item Registration

Observation items (phenomena not yet sufficient to change the directive) are registered in §12.4 and evaluated before the next slice starts.

### 12.4 Observation Item Ledger

| # | Registered on | Phenomenon | Disposition status |
|---|---|---|---|
| OB-1 | 2026-09-21 | **Effectiveness of ⑥ on documentation-only slices is pending evaluation**: MAI outputs normally on code slices (A6), yet on the documentation-only `.md` slice (B3c) it deviated in format twice in a row with no actionable finding (n=1, insufficient to change the directive). The next **code slice** serves as the control (**currently expected to be the DR-2 + DR-3 fix slice; if B4 · AT-23 comes first, B4 takes precedence**; trigger condition: observe as soon as that slice enters ⑤ pre-review, and register the conclusion at its wrap-up): if normal ⇒ attribute it to "documentation-type tasks"; if it still deviates ⇒ attribute it to the model. If 2–3 more slices behave the same way, revisit adding an exception to §7.2 (contract revisions that are pure `.md` and involve no executable semantics may skip ⑥, with ⑦ confirming). **2026-09-21 update (n=2)**: when ⑥ was rerun for SW-1/SW-2, MAI's output was **fully compliant in format** (Sections A–E complete, conclusion line correct) and produced 1 real finding ⇒ the "documentation-type task characteristic" attribution is **weakened**; it is more likely an incidental/prompt factor specific to B3c at the time; keep observing and do not change §7.2 for now. | Observing |
| OB-2 | 2026-09-21 | **G2 self-reference defect**: a rule line itself quoted the forbidden placeholder string ⇒ a full-text scan inevitably produces false positives (4 measured). G2 was amended in place to "exclude rule quotation lines". | **Resolved** (v1.1); the ⑦ re-review confirmation **was not run separately** (v1.6 round 3, C.0f's ⑦×3 and B4's ⑦×2 all left this row uncovered), and it **has been folded into V1.12's ⑦ input** (same reading as OB-9) |
| OB-3 | 2026-09-21 | **G4b criterion unsatisfiable**: the original wording "the difference set must be empty" can never hold for any new document (a new 51 KB document measured a 7-character difference set, and character-by-character re-check found all of them to be legitimate new words: 飞/咨/肉/忘/百/摸/纠). Changed to "the difference set ⊖ the `tools/gates/cjk-newwords.txt` whitelist must be empty", fixing human judgment as a mechanically rerunnable assertion. | **Resolved** (v1.1); the ⑦ re-review confirmation **was not run separately** (v1.6 round 3, C.0f's ⑦×3 and B4's ⑦×2 all left this row uncovered), and it **has been folded into V1.12's ⑦ input** (same reading as OB-9) |
| OB-4 | 2026-09-21 | **Existing role files have non-compliant encoding**: the 11 generated `*.agent.md` files under `.github/agents/` are no-BOM+CRLF, conflicting with the `.md` hard rule; because of a directory-level ignore they had **never been covered by G4**. ⇒ Immediately exposed after `.gitignore` was narrowed. | **Resolved** (`SW-3` done, 2026-09-21, user ruling: option A, see Appendix C.0c; the former "observing, pending user adjudication" was a stale status and was closed on 2026-09-23; the first status word is checked against the G6-6 closed set in §6.3) |
| OB-5 | 2026-09-21 | **The G3 criterion is unsatisfiable for existing bilingual pairs**: "equal CN/EN line counts and equal bold-marker counts" is permanently false for `REMAINING_SLICES` (501/488) and `B2-EXTERNAL-ENFORCEMENT` (270/315). Changed to "**symmetric changed-line count**" (compared via `git diff --numstat`); all 6 pairs of this slice PASS. | **Resolved** (v1.1); the ⑦ re-review confirmation **was not run separately** (v1.6 round 3, C.0f's ⑦×3 and B4's ⑦×2 all left this row uncovered), and it **has been folded into V1.12's ⑦ input** (same reading as OB-9) |
| OB-6 | 2026-09-21 | **The direction of ⑦'s first-round F3 fix was overturned by measurement**: ⑦ suggested "adding `model: \"\"` to comply with §2.6", but measurement showed that empty-value key makes the agent **completely unloadable** (isolated verification: adding the key fails / removing it restores). ⇒ The remediation was **executed in reverse**: fix the directive (§2.6 mandates **omitting** the key) rather than the files. Lesson: the reviewer's **proposed fix** likewise needs measurement verification, and must not be written down just because the "source is authoritative". | **Resolved** (v1.2); the ⑦ re-review confirmation **was not run separately** (same reading as OB-2/3/5), and it **has been folded into V1.12's ⑦ input** |
| OB-7 | 2026-09-21 | **⑥ and ⑦ overlap independently (positive signal)**: under the condition of **no checklist attached in the first round**, ⑥ (MAI) independently reported **the same** defect as ⑦'s 2nd round (a changelog table row joined together), and its Section A–E output was **fully compliant in format** with substantive Section D/E content. ⇒ Supports "⑥ has independent value and should not be skipped lightly". | Observing |
| OB-8 | 2026-09-21 | **The default `model` of existing role files is on the blacklist**: among the 11 existing files, **8** (Terra 3 / Luna 2 / Gemini 1 / GPT-5.4 1 / GPT-5.6-Sol 1) had a default `model` on the §2.4 blacklist (most of them not in the §2.2 role table) ⇒ calling directly by `agentName` without passing `model` explicitly can bypass R2.1 and hit blacklisted models. | **Resolved** (v1.7, user adjudication "delete default values + hard gate", files not deleted): the delete-field probe passed ⇒ 8 files had their `model` field removed, 3 non-blacklisted defaults were retained, and **G5-b** was added to prevent regression |
| OB-9 | 2026-09-21 | **Independent-verification gap for T3/T4 and v1.7/v1.8 (user adjudication)**: the T1–T4 remediation of ⑦'s 3rd round **went through self-verification + ⑥ independent scan only**, without running ⑦ again (the 4th round is not automatically authorized under the user's rules); **v1.7** (OB-8 / identifier disambiguation), the **`_EN` mirror**, and **v1.8** (Appendix B.7) likewise **did not pass independent final review**. | **To do**: **the input of the next ⑦ must contain** — the diff of the T3/T4 criterion changes, the change diffs of v1.7 and v1.8, and the `_EN` mirror (**without running a separate 4th round**) |
| OB-10 | 2026-09-22 | **⑦ final review (v1.9) found that the new §1.6 clauses lack mechanical criteria** — ① a stopped-state report header must read "stopped — incomplete", ② an empty "pilot metrics declaration" slot in a `[PILOT]` kickoff message — both currently have **wording but no G1–G5 criterion**; ⑦ suggested adding them to §6.3. **Not adopted in this slice** because of the SW-5 boundary (**does not change the G1–G5 criteria**). **Empirical note (2026-09-22)**: while closing out SW-5, **4 more temporal inconsistencies of the same class were missed before the commit** (found by Flash self-review) — this class of defect occurred **at least 3 times** within this slice, and happened right after two same-class findings ⇒ showing that **manual visual checking is itself unreliable**, which is the basis for **promoting this item to a mandatory item** the next time §6.3 is opened. | **Observing**: a candidate for a later §6.3 revision (planned to be promoted to mandatory) |
| OB-11 | 2026-09-22 | **The C.0e ledger temporal fix was remediated without re-running the ⑦ re-review (user-ruled exception, option A)**: the C.0e status row said "awaiting the ⑥ independent scan", conflicting with the "⑥ independent scan (**actual**) `PASS`" recorded below (⑦ re-review Medium); after the remediation **⑦ was not re-run** — §5.3 requires "until no Medium+", while the §7.5 budget (1 final review + 1 re-review) was exhausted. **The decision standard for this exception**: the remediation belongs to the **ledger-metadata class** (no rule substance, no gate impact, **closure literally verifiable**, no new semantics), which differs in nature from the **criterion-semantics class** remediation authorized in SW-3's round 3 (changing the G5-a decision standard, containing rule substance, affecting gates ⇒ semantics verifiable ⇒ independent re-review mandatory); the former may be folded into the next revision's ⑦ input, the latter must be independently re-reviewed. **Related lesson**: a ledger **must not carry expiring literals** (line counts / file counts); write only criteria such as "verified consistent by G3/G4" (this slice removed `995/995` accordingly). | **To do**: fold into the next revision's ⑦ input (the exception standard is in §7.8.4) |
| OB-12 | 2026-09-22 | **§11.2's write-back scope does not distinguish document domains (to be discussed)**: the `DEV_PLAN` format is "a T027 prerequisite status update / prfrail project implementation plan" and **does not cover "authoring the delivery directive" itself** ⇒ the authoritative ledger for directive-governance slices (`SW-*`) is the **directive's Appendix C.0–C.0e**, so they **do not write back to `DEV_PLAN`** (user ruling 2026-09-22; measured: `DEV_PLAN` contains no `SW-*` section at all, and `SW-1`'s write-back was the §1 pointer to the new directive, not a slice section). | **To be discussed**: whether to add a sentence to §11.2 ("the write-back scope of directive-governance slices: Appendix C rather than `DEV_PLAN`"); evaluate at the next revision |
| OB-13 | 2026-09-22 | **The acceptance evidence of a "flaky-fix slice" must be split into two stages (`DR-FIX` pilot finding)**: ① **before the commit** — variation evidence (injection class 甲-i: the regression would be caught / injection class 甲-ii: the fix body itself is falsified) + `-count=10`, proving **the fix mechanism works**; ② **after the commit** (after the ⑩ authorized push) — **N consecutive pushes that are green on the first attempt**, proving **the fix works on the real CI**. A failure means going back to ① per §5.3 and separately registering a "fix failure" under **boundary ③** per the DR disposition standard. Root cause: the order in §5.1 is ⑦→⑨→⑩, whereas the core acceptance evidence of such slices (multiple consecutive first-attempt greens) **can only be obtained after a push** ⇒ neither ⑦ nor ⑨ **can verify** it before the commit. | **To be discussed**: whether to revise the position of ⑨ in §5.1 or introduce the concept of "closing evidence after ⑩"; evaluate at the next §12 revision. **Executed in this slice**: ⑨ is split into these two stages accordingly, and the 3 consecutive pushes after the commit are written into the §5 "after the commit" subsection of the report |
| OB-14 | 2026-09-22 | **The §12.1 "defect self-capture rate" was NOT met in the `DR-FIX` pilot (2/5, target ≥ 3/4)**: **root cause** — 2 of the 3 reviewer-side findings belong to the **ledger-metadata class** (in-document self-contradiction, status-line tense) and to the **review-bundle completeness class** (missing evidence chain), while §6.3's **G1–G5 have no mechanical criterion for either class** ⇒ the **same root as OB-10** (the same class of gap demonstrated again within one slice). **Suggested action**: promote "**ledger-metadata consistency** (a row must not outlive the fact it records)" and "**review-bundle completeness** (every claim must carry its raw evidence)" to the next §6.3 extension candidates (proposed as **G6 / G7**); **if such defects then become self-captured, the "self-capture rate" metric can be re-evaluated**. **Next slice candidate (same domain, mergeable)**: `[SLICE]` — the §6.3 extension (G6 / G7) plus scripting `tools/gates/` (move `tmp/gate.js` in and add CI integration, path references and a scripting-status declaration; **not a one-line move**). | **Ruling** (who ruled = the user; date = 2026-09-23; reason = the root cause is clear - no mechanical criteria for the ledger-metadata / review-bundle-completeness classes, and `GATES-EXT` has already extended G6 / G7; re-evaluation point = whichever of `GATES-G1G5` / `GATES-TIGHTEN` completes first): continue; the ruling and the re-evaluation point are in the re-evaluation clause of §12.2; G6 / G7 have landed via `GATES-EXT`; follow-on slices `GATES-G1G5` → `GATES-TIGHTEN` (both after `B4`, user ruling 2026-09-22); **the re-evaluation object is DR-FIX (the slice OB-14 belongs to)**, recomputed by "how many findings the extended criteria could have caught" (judge = ⑦, anchor = running G6 / G7, 2 classifiable items), see the re-evaluation-object clause of §12.2 |
| OB-15 | 2026-09-22 | **Three coverage gaps of `tools/gates/` before the migration (user-dictated, registered truthfully with this slice)**: ① a hardcoded repository root (`D:/LZProjects/prfrail/`) made the old throwaway script unusable on Linux CI; ② a constant exit code 0 plus silent misses on **untracked / staged-only** files meant that **every "gate green" claim before 2026-09-22 covered G2 / G3(b) / G4b / G1-a incompletely**; ③ base G1 and G5 were never implemented. **Handling in this slice**: ① repository-root discovery via `git rev-parse --show-toplevel` (the ② miss surface is fixed in the same pass), ③ printed honestly as `SKIP` (no false claim of coverage). **Historical "gate green" claims are read as covering what the tool actually covered at the time and are not retroactively rewritten.** | **Resolved** (this slice `GATES-EXT`, batches 2a / 2b): ① no hardcoded path, ② `--scope` covers untracked and staged-only files and returns a non-zero exit code, ③ coverage is carried jointly by the script's `--all` list and this honest registration |
| OB-16 | 2026-09-22 | **When a slice should register the section 12.5 gap row it produces (to be discussed)**: this slice's gap row was registered before the slice started (with an honest receipt annotation), while section 1.5 wants a clean workspace before starting ⇒ the two pull against each other in time. **To be discussed**: register such a gap row **at slice start** (current practice: yes, with an honest receipt annotation) or **after the commit**. **Proposal**: evaluate it together in the next section 12 revision. **Numbering note**: OB-15 / OB-16 are the numbers this slice took; the user may renumber. | **To be discussed**: evaluate at the next section 12 revision |
| OB-17 | 2026-09-22 | **Known residual false positive of the G6-3 angle-bracket rule (measured during `GATES-EXT`; pre-existing)**: CLI usage examples such as `prfrail run <chain-file>` are taken for run ids by the angle-bracket placeholder rule - two instances exist (`docs/RFC-proofrail-unattended-ai-engineering-product.md:737` / `:1085`). Because only changed lines are judged this **does not affect the gate today**, but **either line would produce a spurious red as soon as it is touched** (measured: the guard-removal comparison shows NEW-only=0 ⇒ not introduced by this slice). **Not fixed here** (narrowing the rule to "id-like placeholders" is a criterion-semantics change beyond this slice's authorization), and **no whitelist entry was added** (it would mask this shape). This row carries the `示例` marker for the same reason as the note above (see OB-18). | **To do**: evaluate at the next section 12 revision or in a later `[SLICE]` - narrow the angle-bracket rule, or move that example line onto an exclusion-marked line |
| OB-18 | 2026-09-22 | **G6's exclusion markers are Chinese tokens, so an `_EN` mirror must carry one of them to be exempt** (measured by `GATES-EXT` self-reference): the marker set is `不得残留` / `示例` / `引文` / `冲突` / `矛盾`; if an `_EN` mirror quotes a banned shape verbatim **without any marker**, G6 trips on it (**measured in this slice**: `docs/DELIVERY_DIRECTIVE_EN.md:394` and `:638` were flagged by G6-3 for lacking `示例`). **What this slice did**: keep the `示例` marker inside the English lines (the established convention of the §6.3 English example line). | **To be discussed**: whether to extend the marker set to both languages (e.g. `example` / `quote` / `conflict`) - the exemption surface of common English words needs evaluation; evaluate at the next revision |
| OB-19 | 2026-09-22 | **Four measured verification mechanisms of DELIVERY_DIRECTIVE, proposed as design input for S2 (user-confirmed)**: ① **mechanical assertions instead of human memory** (a pluggable criterion library built into the product); ② **tests cannot self-certify** (the A11 / A12 mutation checks as a differentiating capability); ③ **the cost gradient is an architectural constraint** (a cost gradient built into the model strategy instead of defaulting to the most expensive model); ④ **the closure iron rule** (AgentRunner's "fix ⇒ re-run the current step" state machine). **Out of scope**: concrete rules, naming schemes, moving the directive file. **Timing**: before the S2 design starts, together with the Sol strategy review question set. **Boundary note**: this row is unrelated to this slice's subject (§6.3 extension G6 / G7); the user registered it during this slice as a temporary observation item, and it is **not** a `GATES-EXT` deliverable. | **To be discussed**: before the S2 design starts, together with the Sol strategy review question set |

| OB-20 | 2026-09-22 | **`g5a.js` error-routing gap (`GATES-EXT` ⑦ round-3 Low, registered not fixed)**: after G5-a① calls `git grep`, a falsy `r.ok` turns the output into an empty string ⇒ indistinguishable from "zero hits", so the sub-check still reports PASS (`grep` exit 1 = no match, any other non-zero = execution failure; the two share one path today). **The same defect class was dug out by ⑦ on three consecutive rounds of this slice** (empty change set vacuous pass, failing `git status`, inconclusive shallow probe); this is the last remnant of that class. | **Ruling (user, 2026-09-22)**: handled in the next **§6.3 criterion-narrowing slice** (`GATES-TIGHTEN`); **not fixed in this slice** - rationale: ① ⑦'s quota is exhausted (3 rounds used by this slice) and fixing it would need a round 4, which falls inside red line 3's semantic zone; ② round 3 of ⑦ explicitly returned `safe to commit and push as-is`, and a Low does not block; ③ the next slice, "scripting base G1 / G5" (`GATES-G1G5`), adds the **unimplemented** G1 / G5 main criteria and is **not** about fixing the already-implemented `g5a.js`, so OB-20 does not belong to it. The fix (treat only exit 1 as "no match", raise a usage error for other non-zero) is executed with that slice. |

| OB-21 | 2026-09-22 | **B4's premise does not match the current state (measured by the B4 ①)**: ① the `prfrail` CLI has **no AgentRunner wiring** (`internal/console/runtime.go` only has `ExecuteNoopRun`, and `RunPostflight` is a fail-closed stub); ② **no legal, startable candidate exists** (T026's verdict "the current candidate is not compatible" + B2 §8.3 the candidate failing inside the zero-capability container / §8.4 the endpoint returning WAF 403 + DR-1 not fixed). ⇒ **the "real AI candidate" leg of AT-23 is unreachable today**. **Impact**: it affects not only `B4` but **also the S1-exit reading**. **Timing**: the next S1-exit assessment plus any candidate / assembly progress. **Side note**: whether the remaining T018 gates share this dependency is registered as a T018 open question in `docs/DEV_PLAN`. | **To be discussed**: the next S1-exit assessment plus any candidate / assembly progress |

| OB-22 | 2026-09-23 | **Internal ordering inconsistency inside `CONTRACTS` §7 (found by the B4 ⑦ final review)**: L205/L227 say "freeze the candidate, then run the gates", L229 says "postflight before `REVIEW_PENDING`"; the **measured behaviour and the implementation both agree with L229** (evidence: `internal/chain/engine.go:186/197/364` and `tmp/b4/runs/f5-pos/events/state-events.jsonl:13/14/15`) ⇒ **the implementation is correct and the contract text must be unified**. | **Resolved** (2026-09-23, user ruling ⒜): the `CONTRACTS` §7 wording is unified to the measured order "in-step gate → postflight → freeze / review / promotion", **wording only, no code change**, CN and `_EN` kept in sync; **no ADR and no `DEV_PLAN` write-back** (user ruling ③: unifying wording is not in the ADR-required class, and B4 is already carried by `DEV_PLAN`; the first status word is checked against the G6-6 closed set in §6.3) |

| OB-23 | 2026-09-23 | **B4 ⑤/⑦ disagreement on "sufficiency of the alternative evidence for skipped faces" (to be discussed)**: within one slice ⑤ judged it **sufficient** under the honest-fail-closed reading while ⑦ judged **all five insufficient** under the equivalent-coverage-of-production-semantics reading. ⇒ the criterion **lacks one agreed reading**; both readings are reasonable but not equivalent. **Suggested**: adopt ⑦'s equivalence reading and state explicitly that a "**non-equivalent substitute**" is a legal but **must-be-labelled** status. | **To be discussed**: evaluated at the next §12 revision (whether to define skipped-face evidence sufficiency in §7 or §12) |

| OB-24 | 2026-09-23 | **B4 ⑦ re-review (round 2) reported 2 Medium findings (the same issue in CN and `_EN`, prose scope)**: the report's §9.7 sentence "this slice changed the contract wording and no code" could be read as a whole-slice no-code claim while the staged set contains `tools/b4-e2e/*.go`. **The decision is recorded as measured**: ① it was made autonomously by Flash **while the user was offline**; ② it rests on **the exception standard OB-11 already established** (prose only, no rule entity, no gate impact, closure literally verifiable, no new semantics); ③ **the reason no round 3 of ⑦ was run** is that the §7.5 quota was exhausted (2/2 for B4); ④ **if the user later judges it semantic rather than prose, running round 3 is all that is needed - that case is a judgement disagreement, not a rule violation**. Closure: constructive elimination (the sentence scoped to the R-1 sub-change) plus a deterministic criterion (`git diff --cached --numstat -- internal/` empty) plus a green ④ re-run. | **Ruling** (user, 2026-09-23: no round 3 of ⑦; see §7.8) |

| OB-25 | 2026-09-23 | **Whether §3.4 applies to the `[ROUTINE]` tier is not stated (V1.11 ⑤ pre-review Low, registered to be discussed)**: §3.4 requires a diagnostics scan "after any edit and before moving to the next step" but does not say whether it applies to the routine tier (the §1.6 routine-tier flow does not reference §3.4) ⇒ an interaction ambiguity about its scope. **Not ruled on here** (an "expansion or execution freedom" ambiguity, so it was stopped and referred to the user, who ruled that it be registered as to-be-discussed). | **To be discussed**: `DIRECTIVE-V1.12` to assess whether the scope is stated in §3.4 or §1.6 |

| OB-26 | 2026-09-23 | **§3.4.2's "newly added warning" lacks a baseline definition (V1.11 ⑤ pre-review Low, registered to be discussed)**: §3.4.2 requires that "any newly added warning must not be left to the next step" but never defines "newly added" against which baseline (existing warnings or the previous run) ⇒ that hard constraint is **partly not mechanically checkable**. **Not ruled on here**. | **To be discussed**: `DIRECTIVE-V1.12` to assess it together with the §3.4.3 diagnostics aggregator and define the baseline snapshot |

| OB-27 | 2026-09-23 | **⑧c's "all types, file by file" encoding check meets historical leftovers in frozen evidence packs (found by the V1.11 ⑧c first run)**: the whole-tree file-by-file sweep reported 3 pre-existing violations (`b2-2026-09-17/measuring/probe10-dir-listing.raw.txt` is CRLF; `b3b-2026-09-21/{session,variant}/pin.json` carry a BOM); these 3 files **do not violate G4a** (whose scope is explicitly the change set); what they violate is **the "not only the change set" wording ⑧c adds**. This slice disposed of it under the user's ruling, option A: a temporary exemption (strictly limited to "archived packs that carry a MANIFEST/SHA256SUMS"; a new evidence pack must still pass the encoding gate before it is committed) plus ⑧c checking evidence packs by self-consistency. | **Resolved** (this slice, v1.12, written into §5.1 and §5.3); it must **coordinate G4a and ⑧c and the B2 fidelity ruling** - **the CRLF of `probe10-dir-listing.raw.txt` directly conflicts with `.gitattributes`'s `*.raw.txt -text`** (the B2 fidelity ruling) |
| OB-28 | 2026-09-23 | **G6 judges only changed lines, so it structurally cannot catch a ledger row that is stale without having been touched (measured in V1.12)**: adding a column to §0.3 pulled 12 historical rows into the criteria's scope, and one file-count phrasing in the v1.6 row was immediately flagged by G6-4 as an expiring literal ⇒ **a historical row is re-judged as soon as it is touched, while an untouched stale row stays silent forever** (the stale statuses of OB-2/3/5/6 share this root). **Related**: the `*.raw.txt` `-text` conflict (the B2 fidelity ruling) left over from OB-27 is discussed together with this item. | **Observing**: the next §6.3 revision evaluates a "whole-ledger tense re-check" criterion or a periodic full sweep, plus the priority ruling for `-text` versus ⑧c |
| OB-29 | 2026-09-23 | **The G8 mechanical criterion suggested by H1** (the set of `model=` values in Appendix B ⊆ the `modelId` column of §2.2): a **new criterion**, outside V1.12's boundary (this slice does not change any G1–G7 criterion body). | **Observing**: candidate for the next §6.3 criterion revision (evaluated together with `GATES-TIGHTEN` / `GATES-G1G5`) |
| OB-30 | 2026-09-23 | **M9's five-key frontmatter criterion for role files plus Appendix D "role file template"**: the former is a new criterion (beyond V1.12's boundary); the latter is not done here, avoiding the half-finished state of "a template without a criterion". | **Observing**: evaluated together with the next §6.3 criterion revision or a role-file-template slice |
| OB-31 | 2026-09-23 | **M7's threshold-system rebuild**: a unified rework count R (R≥2 ⇒ high risk, R>3 ⇒ pause) touches the §5.2 matrix's trigger semantics; the §7.5 architect row "≤1 (+1 after a revised proposal)" touches a cost-cap number. V1.12 fixes only the terminology definition (the §9.2 counting clause); the threshold rebuild is parked as a package. | **To be discussed**: user adjudication or the next §12 revision |
| OB-32 | 2026-09-23 | **No mechanical criterion for §0.3 change-log ↔ body consistency (the C1 evidence of V1.12)**: the v1.12 §0.3 row claimed 34 changes, of which **19 do not exist in the body** (the third-party v1.12 review's section 2 reconciled them item by item and called it Critical: the miss mechanism is provided by the directive itself - ⑦ scopes its review by the log and therefore skips the areas claimed as disposed). The review suggested a new criterion (G6-5 / folded into G1-b): **the section numbers / IDs / wording a log row mentions must really exist in the same change set**. | **To do** (next slice, see §7.8 and G6-5 in §6.3): **the approximate criterion G6-5 has landed in this slice** (the set of §X.Y references appearing in §0.3 and in Appendix C ledger blocks ⊆ the file's actual heading set); **full-class coverage (C1's content-level reconciliation) stays an observation for `GATES-TIGHTEN`** |
| OB-33 | 2026-09-23 | **Cross-position consistency review is missing (raised by the user on 2026-09-23)**: the three-layer architecture of "script layer (`tools/gates/`) + slice-level MAI + periodic whole-corpus layer" is absent - the same fact is written in several places (body / templates / ledger / validation reports), any single edit can miss a sync, and today's criteria cover only local consistency on changed lines (G6's scope = changed lines) ⇒ a "cross-position consistency" review layer is needed. **Relation to OB-32**: OB-32 is one concrete gap of that direction, on "log ↔ body". | **Observing**: the next §6.3 revision or a dedicated slice evaluates the three layers (script-layer criteria / slice-level independent review / periodic full sweep) |
| OB-34 | 2026-09-23 | **§7.2.2's heading overflows its content (R6's N7)**: the subsection is titled "⑥ output contract (Section A–E)", yet "cost ceiling" and "failure handling" sit under it, which is not output-contract matter ⇒ the title misnames the content. | **Resolved** (v1.13; retitled "output contract, cost cap and failure disposition of the ⑥ layer", see §7.2.2) |
| OB-35 | 2026-09-23 | **§1.6's old section reference was not updated when §7.2 was split (R6's N8)**: §1.6 says "all the semantic categories listed in §7.2 ③", while the three conditions now live in §7.2.1; the reference is imprecise (not dangling, but a reader lands nowhere). | **Resolved** (v1.13; boundary ① of §1.6 now reads "§7.2.1 ③"; historical rows such as `C.0e` stay unchanged per §11.1) |
| OB-36 | 2026-09-23 | **OB-14's disposition cell substitutes a pointer for the four elements (R6's N9)**: §12.2's re-evaluation clause ① requires "register the ruling in §12.4 (who ruled / date / reason / re-evaluation point)", while OB-14's disposition cell currently points at §12.2 ⇒ the four elements are not in that cell. | **Resolved** (v1.13; the four elements are written directly into OB-14's disposition cell, with the pointer kept only as a supplement) |
| OB-37 | 2026-09-23 | **The applicability boundary of "no ⑦ re-run under the OB-11 exception" was broken by this slice (R7's P4)**: the OB-11 exception standard is "ledger-metadata class: **no rule entity**, no gate impact, closure verifiable literally, **no new semantics**", yet the N1 (re-evaluation-object definition) and N2 (deferral clause) landed in this slice's R6 round **both carry rule entities and new semantics** ⇒ **the OB-11 condition does not hold**; that exemption was in fact a **direct user ruling** (who = the user; date = 2026-09-23; reason = the two fixes accompany the re-evaluation-object reading and had to be settled together with the pilot deferral; basis = the user's ruling text this round). **Generalizable meta-rule (this item must be read as a precedent)**: a direct user ruling **may** break the OB-11 exception, but **must** ① **be recorded independently** (as an OB or an ADR, never only in a note under §0.3 / Appendix C); ② **be labelled truthfully** as "this item carries rule-semantics changes and goes beyond the literal standard of the OB-11 exception, by direct user ruling". ⇒ later slices **must not** cite this item as a precedent for normalising "rule-semantics changes may skip ⑦". | **Registered** (this slice, v1.12); and its **applicability is narrowed immediately in this slice**: later remediation here **must not** use the OB-11 exception again (user ruling in the same batch, 2026-09-23) (the meta-rule is restated in §7.8.4) |
| OB-38 | 2026-09-23 | **G6-1's marker vocabulary and §12.4's actual status vocabulary barely intersect (R7's P5)**: G6-1 takes its pending markers as `待` plus one of ①–⑩ / 终审 / 独立扫描 / 复审 / 预审 / 盲审 / 提交授权, and its completion markers as the realised-parenthesis form / bold realised / 已闭合 / ✅ 已完成; yet the disposition statuses actually used across §12.4's rows contain about 28 forms (观察中 / 待办 / 待议 / 已处置 / 已认可 / 已落地 / 已按用户裁决继续执行 / 已裁决 / 下一片必做 / 下一片处置 and others) ⇒ **ledger status conflicts have no mechanical coverage inside §12.4** (G6's scope covers changed lines). | **To do** (next slice, see §7.8 and G6-6 in §6.3): **the closed status enum and the backfill have landed in this slice** (a six-word bilingual closed set); **widening G6-1 instead is not done** (block-level semantics are incompatible, argued in G6-6 of §6.3) |
| OB-39 | 2026-09-23 | **The v1.12 validation report's artifact list omits the fourth and fifth re-review reports (found in the v1.12 freeze re-check)**: its §11 names 4 review reports such as `tmp/DELIVERY_DIRECTIVE_v111_REVIEW.md` but not `tmp/DELIVERY_DIRECTIVE_v112_REVIEW_R6.md` or `_R7.md`, although both exist on disk ⇒ the reading "list every review report round by round" had nowhere to be fixed. | **Registered** (this slice `DIRECTIVE-REVIEW-SATURATION`): the v1.12 report is frozen ⇒ **its §11 is not changed**; instead this slice fixes the reading as "every review report a slice cites is listed round by round in its artifact list" and records the difference and its reason truthfully in this slice's report |
| OB-40 | 2026-09-23 | **§12.1's second deferral companion meets a documentation-only governance slice (found at the v1.12 freeze)**: this slice, `DIRECTIVE-REVIEW-SATURATION`, is still a directive-governance / review-saturation slice ⇒ it fails §12.1's pilot-selection criteria (the target must not fall on a high-risk or documentation-only slice) ⇒ the **deferral count reaches 2** ⇒ companion ② fires: "two or more deferrals ⇒ the master must stop per §1.7.2 and escalate for a user ruling". **Direct user ruling** (who ruled = the user; date = 2026-09-23; reason = this slice is directive governance and the review-saturation criterion, not a pilotable behaviour-change slice; basis = the user's ruling text this round): **continue executing, not as a pilot**, registering "pilot deferred for the second time" in this slice's kickoff receipt and ledger. **Recorded independently per OB-37's meta-rule**: this is a direct user ruling and **must not be read as a precedent** (a third or later deferral must escalate again). | **Registered** (this slice); deferral count = 2 is recorded in this slice's kickoff receipt and report; the convergence of this class ("quota / deferral / exception") is now defined by §7.8 |
| OB-41 | 2026-09-23 | **A post-freeze append drifts the criteria's scope (found in the v1.12 freeze re-check; registered by the user)**: when a "freeze ruling / CI result" record is appended to an already-committed revision, the index no longer holds the original slice's commit ⇒ ⑨ and G6-4's "changed-line scope" **degenerates from the whole change set to this append** (measured: 1097 refs / scanned 1079 ⇒ 14 refs / scanned 70, and 0 on an empty change set), which **contradicts the readings already written into the report**; this time it was handled as "store this round's readings separately as `tmp/v112/*-sec12.txt`, leave the frozen readings untouched, and state the reading in the report". | **Resolved** (v1.13; written into §9.3's known-anomaly list as "a post-freeze append drifts the criteria's scope", with an executable reading; §6.3's scope definition stays with `GATES-TIGHTEN`) |
| OB-42 | 2026-09-23 | **The verdict line ∈ the §4.1 closed-set enum has no mechanical criterion (human verification covers it today)**: this slice previously recorded "the verdict line ∈ the §4.1 closed-set enum" under G7-c, while G7-c only judges the disposition cell of a finding row; §7.8.2 and §7.8.5 now state this truthfully per OB-15. | **To do** (a later criterion-revision slice assesses whether to add a mechanical criterion; **this slice adds no criterion**) |
| OB-43 | 2026-09-23 | **A PowerShell pipe capturing a native command's UTF-8 stdout ⇒ mojibake in evidence files**: when PS 5.1 writes a native command's stdout to disk through a pipe (`Out-File` and the like), it decodes the UTF-8 bytes through the console code page (GBK) and then writes them as UTF-8 ⇒ Chinese text and ⑥/⑦-class symbols all turn into unrelated characters; measured in this slice: 808 mojibake markers in `remediation.diff` and 22 in `selftest-master.txt`, while a bare `git diff HEAD >` redirect stays clean. | **To do**: fixed immediately in this slice - evidence files are written straight from node (`execFileSync` + `fs.writeFileSync(..., 'utf8')`) and every file is re-checked so that mojibake markers = 0; §9.3 now carries this failure-mode row; later slices keep this capture method |
| OB-44 | 2026-09-23 | **A review-layer role wrote to disk itself (a literal deviation at the role boundary)**: in this slice the ⑤ re-review wrote its report to `tmp/v113/review-r2.md`, whereas §2.2 states that a review-layer role's "may write files = **No**"; that role file's `tools` list holds only `read` / `search`. Classification: **a literal deviation, not a substantive violation** (the content is its report, no product text was polluted, `tmp/` is not committed, the master has confirmed it). | **To be discussed**: ① how a review-layer report is delivered (transcribed by the master vs allowed to write under `tmp/`); ② a note - G6-5's rule text says "appendix C.0–C.0i" while the implementation reads "C.0 up to but excluding C.1", **equivalent today**, and any future C block must be described in step. **This slice neither amends §2.2 nor re-runs ⑤** (its conclusions hold and its content was not polluted); this round's ⑥/⑦ carry a hard constraint written into the task package by the master - "create, modify or delete no file; return the report as your final message" - **without amending the directive**. Both items are left for the next slice to assess |
| OB-45 | 2026-09-23 | **A review-layer report misstates its "data source"**: this slice's ⑥ independent scan described "**reading** the raw output the master wrote to disk" as "**running** `node tools/gates/gate.js` / `selftest`", whereas that role file's `tools` list holds only `read` / `search` and it cannot run commands; the quoted values match the on-disk evidence string for string, yet that phrasing gives the impression of "independent execution", while §7.7's A11 / A12 depend precisely on **traceable provenance**. | **To be discussed**: ① whether the ⑥ / ⑦ task package should explicitly write "read-only paths + **make no claim of running** `node` / `gate` commands, citing only the paths the master wrote to disk"; ② whether the wording template for the "data source" statement in ⑥ / ⑦ reports needs to be unified. **This round does not change the ⑥ / ⑦ prompt template** (already covered as a hard constraint inside this slice's task package); the rule layer is left for the next slice to assess |

> **This slice's (`GATES-EXT`) ⑦ round exception (user pre-authorization, 2026-09-22)**: the ⑦ re-review (round 2) reported Medium+, which forces a re-review under §5.3 and, under the pre-authorization, entitled a **round 3** (`GPT-5.3-Codex`, **valid for this slice only**); round 3 returned `PASS WITH FIXES` with no Medium+ (one Low only, i.e. OB-20). **The authorization does not generalize**: other slices stay on §7.5 (one final review + one re-review), and a round 4 is never automatic.

### 12.5 Gaps and Future Extensions

| Gap | Suggested timing |
|---|---|
| Integration coordinator (cross-slice merge conflicts, dependency upgrades) | When multiple slices run in parallel |
| Performance reviewer | When performance-sensitive slices increase |
| Incident response process (rollback / hotfix) | After the first production incident |
| Retrospective analyst (multi-slice data review) | After execution data accumulates |
| Infrastructure maintenance (CI/CD, test infrastructure, model version upgrades) | Continuous, independent of slices |
| **The decision standard for the §5.3 vs §7.5 conflict** (after the budget is exhausted, which class of remediation may skip the ⑦ re-run and which must break through) | **Landed** (`DIRECTIVE-REVIEW-SATURATION`, v1.13, see §7.8) |
| **The §6.3 criteria are not scripted (G1–G5 are executed by hand by the master) + no mechanical criteria for ledger-metadata consistency / review-bundle completeness** | **Landed** (`GATES-EXT`, 2026-09-22, see Appendix C.0f); **base G1 / G5 are still unscripted** (they print `SKIP` truthfully) |
| **Criterion self-reference risk**: a G6 / G7 false positive is amplified by CI into a repository-wide block; and G6 judges only changed lines, so it cannot catch a stale ledger row that was not touched | the same root as OB-28; **candidate slice = `GATES-TIGHTEN`** (user ruling 2026-09-23) |

---

## Appendix A: Historical Versions and Sources of Experience

### A.1 Sources of This Directive

| Source | Absorbed content |
|---|---|
| 《Flash Master-Subagent Architecture v1.0》 (`tmp/主控子代理架构_v1.0.md`) | Three-layer role division, fallback mechanism, effort tier table, prompt template skeleton |
| 《Flash Master Operating Directive v3.1》 (`docs/t027/FLASH_OPERATING_DIRECTIVE_v3.1.md`) | Model whitelist/blacklist, ⑥'s three trigger conditions and cost cap, ⑦'s anchoring-based hard isolation, A11/A12, input sanitization, conclusion-landing check, protocol-first, cost awareness, Git discipline, report structure |
| 《T027 Briefing v3.1》 | Launch package structure and minimum context feeding (**landed in v1.8 as Appendix B.7 "Slice Startup Card"**; its "slice-level hard gates" belong to the slice definition by nature, see C.3) |
| B3c slice practice (2026-09-21) | Evidence landing obligation (§6.2), G1–G5 self-check gates (§6.3), character-level miswriting lessons (§9.3), SCM stale display, causes of ⑦ invalid rounds |

### A.2 Disposition of v3.1

`docs/t027/FLASH_OPERATING_DIRECTIVE{,_EN}.md`, `..._v3.1{,_EN}.md`, `FLASH_T027_BRIEFING*.md` are **retained in their original form** (historical trace), with one line added at the head of the file:

> This file has been superseded by `docs/DELIVERY_DIRECTIVE.md` v1.0 (2026-09-21); it is retained only as a historical version and source of experience.

### A.3 Section Number Cross-Reference

See [Appendix C](#附录-c迁移映射).

---

## Appendix B: Role Prompt Templates

> How to use: the master picks each template up by role, **trims the context**, and feeds it in; the `<...>` in the templates are mandatory slots.
> All templates share the §4.1 general output contract.

### B.0 Master (orchestration layer, the session itself)

The master does not dispatch itself via `runSubagent`; this template is the **master's own working contract**, executed as each slice starts (consistent with §3, §5, §6, §9).

```
Role: master (development team lead). Model: deepseek-flash (the session itself).
Sole responsibility: command, dispatch, gates, evidence, fallback, reporting. **Do not originate product artifacts** (§3.2).

Startup checklist (§5.0):
  [ ] receive the slice definition and verify dependencies / workspace hygiene / budget
  [ ] determine the slice type → set the effort tier (§8.2; once set, do not switch mid-way)
  [ ] generate task packets (each packet: input material / expected output / acceptance criteria / modelId string)

Execution checklist:
  [ ] **pass modelId explicitly** when dispatching (R2.1); do not give write permission to read-only roles
  [ ] before review, **persist the raw diff and gate output to disk**, and hand the **paths** to the reviewer (R6.1)
  [ ] keep an artifact at every ring; on failure, rerun that ring (§5.3)
  [ ] ⑦ first round **attaches no checklist of any kind** (§7.3); ① and ⑤ are **call-isolated** (R2.6)
  [ ] record each fallback action; **> 3 times ⇒ mark high risk and stop** (§9.2)

Closing checklist:
  [ ] §6.1 baseline gates all green
  [ ] §6.3 all seven categories of mechanical checks G1–G7 green
  [ ] report assembly (§11.3; do not rewrite subagent conclusions)
  [ ] stop at ⑩ and await authorization; **must not** commit / push yourself (§10)
```

The output must end with one line: `SLICE: READY_FOR_REVIEW` (or `SLICE: BLOCKED` + the blocking items).

### B.1 Architect (①)

```
Role: architect. Model: deepseek-v4-pro. Read-only; must not modify any file.
Vehicle: runSubagent(agentName='deep-reasoner', model='deepseek-v4-pro')

Input: ① slice definition (goal / boundary / dependencies / scale) ② the relevant contract text ③ historical rulings and known constraints
Output: the architecture design document, containing:
  1) analysis conclusions (feasibility / complexity / effort estimate)
  2) change-point list (file level: path / add-modify-delete / description / impact scope)
  3) test-point list (ID / description / positive-negative-boundary-concurrency-crash / priority)
  4) risks and mitigations
  5) special cautions (cross-platform, compatibility, performance)
Forbidden: writing code (including pseudocode); **carrying pre-implementation reasoning context into ⑤** (anti-anchoring relies on call isolation, see R2.6).
The output must end with one line: `DESIGN: DONE` (or `DESIGN: NEEDS_CLARIFICATION` + the items needing clarification).
```

### B.2 Implementer (②)

```
Role: implementer. Model: deepseek-flash. May read and write files and run commands.
Vehicle: runSubagent(agentName='prfrail-implementer', model='deepseek-flash')

Input: ① the change-point list from the design document ② the contract text ③ the exact context of the files to modify
Requirements: implement strictly per the change points; follow protocol-first (change CONTRACTS/schema/fixtures first); gofmt; introduce no unauthorized dependencies.
Output: code diff + implementation notes (which change points are covered; whether any are not covered, and why)
Self-check: after the change you must run the focused tests and gofmt yourself, and paste the raw output back.
Forbidden: architectural decisions; writing tests; editing non-code documents. If the design document is ambiguous ⇒ write "needs clarification"; do not guess.
The output must end with one line: `IMPLEMENT: DONE` (or `IMPLEMENT: BLOCKED` + the blocking items).
```

### B.3 Test Engineer (③)

```
Role: test engineer. Model: deepseek-flash. May read and write files and run commands.
Vehicle: runSubagent(agentName='prfrail-tester', model='deepseek-flash')

Input: ① the test-point list from the design document ② the implementation diff ③ the contract text (to verify assertion expectations)
Requirements: cover positive / negative / boundary / concurrency / crash cases; every test point must be falsifiable by "delete the guard ⇒ it turns red".
Output: test diff + coverage matrix (ID / test name / type / corresponding test point)
Self-check: run `go test` yourself and paste the raw output back.
Forbidden: modifying production code; modifying test infrastructure. Untestable items must be explicitly marked "untestable".
The output must end with one line: `TEST: DONE` (or `TEST: BLOCKED` + the blocking items).
```

### B.4 Documenter (⑧)

```
Role: documenter. Model: deepseek-flash. May read and write files and run commands.
Vehicle: runSubagent(agentName='prfrail-documenter', model='deepseek-flash')

Input: ① the design document ② the implementation diff ③ the test diff ④ the list of documents to update plus their current content
Requirements: sentence-by-sentence bilingual alignment (CN authoritative + _EN mirror); consistent terminology; `.md` keeps UTF-8 with BOM + LF;
      fields written back must come from actual evidence (commit ID / run ID / counts); where missing, mark "historical gap"; do not fabricate.
Output: documentation diff + change summary + list of unsynced items
Self-check: after writing to disk, run the encoding gate and the G3 bilingual symmetry count, and paste the results back.
Forbidden: changing code; adding subjective commentary.
The output must end with one line: `DOCS: DONE` (or `DOCS: BLOCKED` + the blocking items).
```

### B.5 Pre-reviewer (⑤)

```
Role: pre-reviewer. Model: deepseek-v4-pro. Read-only; must not modify any file.
Vehicle: runSubagent(agentName='deep-reasoner', model='deepseek-v4-pro')

Input: ① the design document ② the implementation diff (raw hunks, paths already persisted) ③ the contract text
Task: check item by item whether the implementation deviates from the design document; focus on concurrency counterexamples, crash windows, resource leaks, missing error handling, and boundary omissions.
Output the pre-review report: overall verdict (pass / conditional pass / fail) + deviation table (description / corresponding clause / severity / fix)
      + counterexample analysis (concurrency / crash / boundary; write "none" if none) + other concerns
Requirements: every finding carries `file:line`; severity per §4.3.
Forbidden: editing files; inspecting test code or documentation changes (preserve an independent perspective).
The output must end with one line: `PRE-REVIEW: PASS` / `PRE-REVIEW: PASS WITH FIXES` / `PRE-REVIEW: FINDINGS`.
```

### B.6 Independent Review (⑥ independent scan / ⑦ independent final review)

```
Role: independent reviewer. Read-only; must not modify, create, or delete any file.
⑥ Vehicle: runSubagent(agentName='independent-reviewer', model='mai-code-1.1-flash')
⑦ Vehicle: runSubagent(agentName='independent-reviewer', model='gpt-5.3-codex')   ← most expensive; capped per §7.5

Input: ① slice definition ② the raw diff (**paths**, not paraphrase) ③ the original contract sentences ④ the read-only constraint
     !! ⑦ first round **attaches no** ⑤/⑥ checklist; only a re-review round may attach one, and it must be marked "read only after completing the independent scan".
Task: ⑦ four items — security / architectural consistency / completeness / test counterexample review; ⑥ per Section A–E.
Output: each finding = severity + `file:line` + verbatim quote + rationale + a fix that can be written down directly;
      points that "look like problems but are in fact correct" must explicitly state "checked, not a problem".
Forbidden: editing files; recommending lifting a frozen conclusion / loosening a gate / expanding the slice scope.
The output must end with one line whose form comes from the **closed-set enumeration table** of §4.1 (⑦ final = `FINAL REVIEW: …`; re-review = `RE-REVIEW: …`; blind = `BLIND REVIEW: …`; ⑥ = `INDEPENDENT SCAN: …`); a bare `PASS` / `FINDINGS` **must not** be used.
```

---

### B.7 Startup Card (three-tier entry templates; **this card contains no rule substance**)

> **Hard constraint**: this card **holds only the entry point, slots, and section-number references**; it **must not** copy the substance of any rule — the sole authority for a rule is always its **corresponding section in the body**. A rule change **edits the body only**; this card edits only the section-number references. ⇒ This structurally eliminates the old disease of "the same wording written in two places ⇒ inevitable drift" (the old briefing's file-header line "where it conflicts with the directive, the directive prevails" is exactly this kind of drift permit, and is **not inherited**).

**B.7.1 User kickoff message** (per the three tiers of §1.6; square brackets are slots to be filled)

- `[ROUTINE]`: **no template needed** — just write the requirement (tier determination and declaration see §1.6).
- `[SLICE]` / `[PILOT]`: use the templates below; `[PILOT]` additionally needs to fill the "pilot metrics declaration" slot (definition see §12.1).

```
SLICE:[slice ID] [description] (for the `[PILOT]` tier use the trigger word `试点：` instead)
Execute per docs/DELIVERY_DIRECTIVE.md v[version].
Slice definition: [ledger file §subsection ｜ or give the field table of §1.2 inline]
Effort tier: [low / high / max] (the operator has already set it on the UI side; this line is for the record only)
Probe budget: [n times ｜ none]
Authorization boundary: [commit pre-authorized?] / [push pre-authorized?] / [paid-call cap]
[pilot period only] Pilot metrics declaration: [cite the metric line in §12.1]
Notes: [other constraints]
```

**B.7.2 Master opening receipt** (per the three tiers of §1.6; **read back first, then start work**)

`[ROUTINE]` (thin form):

```
[ROUTINE] [one-line reason]
Boundary self-check (§1.6): [not touched ｜ suspected touch + recommend upgrading to [SLICE]]
Record-keeping (§1.6): [diff + §6.1 gates + §6.3 self-check]
```

`[SLICE]` / `[PILOT]` (full form):

```
SLICE: STARTED | SLICE: BLOCKED
Upfront hard gate (§5.0 step 0): SW-1 [✅/⬜]  SW-2 [✅/⬜]  SW-3 [✅/⬜]
Role readiness pre-check (§5.0 step 2): [per role: file exists / the model-key criterion in §5.0 step 2 and G5-a / G5-b in §6.3 / tools include edit / the call will pass model explicitly]
Slice type (§5.2) and effort tier (§8.2): [type] / [low|high|max] (once set, it is not switched during continuous runs)
Budget and authorization boundary: probes [n]; paid-call cap [list each §7.5 item by category]; commit [yes/no]; push [yes/no]
Workspace hygiene (§1.5): [clean ｜ not clean + handling plan]
Read-only this round (pick section numbers per B.7.4): [list of section numbers]
Blocking items (if any): [item by item + the section number each rests on]
```

**B.7.3 Stop-point receipt** (⑩, `[SLICE]` / `[PILOT]` only)

```
SLICE: READY_FOR_REVIEW | SLICE: BLOCKED
Artifacts: [item by item + path / hash]
Gates (§6.1 + §6.3): [result item by item]
Review summary (⑤ ⑥ ⑦): [conclusion line]
Cost and metering (R7.1): [line by line, fields per §11.3]
Items not executed: [item by item]
Awaiting authorization: [commit / push / probe top-up / user adjudication items]
```

**B.7.4 Minimum reading set** (pick section numbers by slice type; the "**All**" row must be read for every type)

| Slice type | Required section numbers |
|---|---|
| **All** | §0.1, §0.2, §1.2, §1.5, §1.6, §1.7, §2.2, §2.4, §2.6, §3.1–§3.4, §5.0–§5.3, §6.1–§6.4, §7.5, §7.6, §10.1–§10.4, §11.1–§11.3, Appendix B.0 + the role templates enabled this round |
| **Document-type / simple** | the table above + §5.2's "Note (document-type slices)", §7.1–§7.3 (including the ⑥ skip criterion in §7.2), §7.7, §9.3, Appendix B.4 |
| **Medium** | the table above + §4.1–§4.3, §7.1–§7.8, §8, §9.1–§9.3, Appendix B.1–B.6 |
| **Complex / high risk** | the table above + §4.1–§4.3, §7.1–§7.8, §8, all of §9, Appendix B.1–B.6 |

> **How to use**: the master **must not** re-read the whole body for every slice; read the section numbers given in the table above, and treat the rest as §-number references. If a needed rule is not in the read set, **read the corresponding section and record it in the opening receipt** — do not infer rules from memory.

---

## Appendix C: Migration Mapping

### C.0 Switchover Slice `SW-1` (reference rewrite)

The switchover itself is a delivery unit, so it must have a slice definition:

| Field | Content |
|---|---|
| ID | `SW-1` |
| Status | ✅ **Completed** (2026-09-21; evidence in `docs/validation/sw-directive-switchover.md`) |
| Goal | Repoint every reference in the repository to the old execution directive at this directive |
| Trigger | As soon as this directive's **CN body is frozen, execution may begin**; the **`_EN` mirror is scheduled as a follow-up slice** (2026-09-21 execution-period wording revision: the original text read "after the CN is frozen and `_EN` is generated", which conflicts with the actual cadence of "freeze CN first, add the mirror later"; the latter is now applied uniformly, and the wording is stated explicitly here to remove ambiguity) |
| Scale | S (documentation only) |
| Executing role | Documenter (⑧a/⑧b) |
| Dependency | None; **must be completed before the §12.1 pilot** (the pilot must run on a tree whose references are already correct) |
| Deliverable | Reference updates for all files listed in table C.2 below + the supersede pointer in the v3.1-series file headers |
| Acceptance evidence | ① **File-level assertion**: **every file hit** by `git grep -l "FLASH_OPERATING_DIRECTIVE\|FLASH_T027_BRIEFING" -- docs` must either (a) belong to the "historically retained file allow-list" (C.0a) and **carry a supersede pointer in its file header**, or (b) already point at this directive; rewriting the in-line references inside the bodies of historical drafts is **not required** (historical drafts only get a file-header pointer). ② **Section-number-level assertion**: every hit from `git grep -n "§3\.10\|§3\.11\|第五章\|Section 5"` must be classified as "historically retained" or "already repointed". ③ §6.3 G1–G5 all green |
| Boundary | **Does not modify the v3.1 source text** (only adds one file-header line); does not change code; does not widen the rewrite scope to historical validation reports |

### C.0a Historically Retained File Allow-List (basis for executing SW-1)

The following files **only get a supersede pointer added to the file header and are not rewritten in the body**; the old-directive references inside their bodies are **not rewritten line by line** (otherwise it would amount to rewriting historical drafts):
`docs/t027/FLASH_OPERATING_DIRECTIVE.md`, `_EN.md`, `_v3.1.md`, `_v3.1_EN.md`, `FLASH_T027_BRIEFING.md`, `FLASH_T027_BRIEFING_v3.1.md`, `FLASH_OPERATING_DIRECTIVE_v3.1.html`, `A6_PILOT_LAUNCH.md`.

### C.0b Role-Carrier Readiness Slice `SW-2` (unblocking an execution deadlock)

**Why it is necessary**: the three product-layer role carriers of §2.2 in this directive (`prfrail-implementer` / `prfrail-tester` / `prfrail-documenter`) **did not exist when SW-2 was chartered** (they were built on 2026-09-21 and passed a real-machine probe; see the "Status" row of this table); meanwhile the models that could write carriers at the time (Luna / Terra / GPT-5.4) were **all on the blocklist**. ⇒ Without building the carriers first, pipeline steps ②③⑧ **cannot be executed**.

| Field | Content |
|---|---|
| ID | `SW-2` |
| Status | ✅ **Completed** (2026-09-21; real-machine probe passed: carriers writable/readable/deletable, 0 residue; evidence in `docs/validation/sw-directive-switchover.md`) |
| Goal | Make the three product-layer roles loadable with correct permissions |
| Trigger | After the CN body is frozen (can run in parallel with SW-1, but the drafts must be kept in sync) |
| Scale | S |
| Executing role | Documenter (⑧a) + master runs the verification commands |
| Deliverable | 3 role files (the header note declares "not a generated artifact, do not let `_sync.py` overwrite it"; the body contains the four sections of responsibilities / inputs / output contract / prohibitions; `tools: ['read','search','edit','execute']`); finely narrowed `.gitignore` (same as §2.5); the three tiers and the model are all left blank for the caller to pass explicitly |
| Acceptance evidence | ① `git check-ignore -v .github/agents/prfrail-implementer.agent.md` **produces no output**; ② all three files are hit by version control (`git status --short` shows them as newly added); ③ **real-machine probe matrix**: **each of the three carriers** runs "create→read back→delete" once, byte-identical with 0 residue (⑦ first-round Medium finding: testing only one carrier would mask permission anomalies in the other two; hence changed from "any carrier" to **full coverage**); ④ the frontmatter of all three files **contains no `model` key** (§2.6; writing an empty value would make carrier loading fail), and the bodies all explicitly declare "the caller passes `model` explicitly"; ⑤ G1–G5 all green |
| Boundary | Does not change any existing role file (including the `sol-orchestrator` toolchain); does not change code |

### C.0c Role-Set Check-In and Encoding Normalization Slice `SW-3`

**Why it is necessary**: after `SW-2` narrowed `.gitignore`, **testing exposed** that the **11 existing generated role files under `.github/agents/` had never been checked in**, and were **no-BOM + CRLF** — conflicting with the repository hard rule "`.md` = UTF-8 with BOM + LF" (G4 flags red); meanwhile R2.5 requires all `*.agent.md` to be brought under version control. ⇒ Both can only be satisfied at once by "**normalize first, then check in**".

| Field | Content |
|---|---|
| ID | `SW-3` |
| Status | ✅ **Completed** (2026-09-21, **user-authorized option A**, see §0.3 v1.5) |
| Goal | Make the existing role files satisfy the encoding hard rule, carry the unified header note, and be brought under version control |
| Trigger | User authorization (already triggered); the originally suggested timing was before the §12.1 pilot |
| Scale | S |
| Executing role | Documenter (⑧a) + master runs the gates |
| Deliverable | The 11 existing `.agent.md` files normalized to BOM+LF and then checked in, and **all 14 role files carry the unified header note** (note text and hard rule in §2.5) — **apart from that header-note line and the encoding/line endings, no semantic word is changed** |
| Acceptance evidence | ① G4a all green (14 files BOM=true / CRLF=false); ② **the only semantic change is the added header-note line** (no semantic change once the header note and line-ending differences are excluded; verified file by file); ③ header-note criterion `grep -L "非生成物（sol-orchestrator 已按" .github/agents/*.agent.md` **produces no output**; ④ `git check-ignore -v` still produces no output for `prfrail-implementer.agent.md`; ⑤ G1–G5 all green |
| Boundary | Does not change `model`/`tools` values, does not change the responsibilities/constraints body; does not touch the `sol-orchestrator/` toolchain |
| **Alternative (rejected)** | A more conservative option: narrow the `.gitignore` exception to `prfrail-*.agent.md`, versioning only hand-maintained files. **The user decided on 2026-09-21 to adopt option A (check in everything)**, on the grounds that `deep-reasoner` / `independent-reviewer` / `quick-verifier` are directly cited by §2.2 as active carriers, and if they are not checked in then **the mechanical gates cannot cover them** (hidden technical debt) |
| **Regression risk** | If the `sol-orchestrator` generator is re-run later, the artifacts will revert to no-BOM+CRLF ⇒ G4 flags red again (R2.5 already requires "check in the generator source and compare verbatim before re-running") |

### C.0d Slice Startup Card Slice `SW-4`

**Why it is necessary**: the old briefing's two functions — `launch package structure` and `minimum context feeding` — have **no landing place** in the new directive (the body already exceeds 800 lines, so feeding it in full for every slice is impossible; §5.0 only says "the user provides the slice definition", with no entry template). But **a second rule copy that would drift must not be made** ⇒ instead it is folded into Appendix B.7 as a **thin card**.

| Field | Content |
|---|---|
| ID | `SW-4` |
| Status | ✅ **Completed** (2026-09-21, **user adjudication: option A**; see §0.3 v1.8) |
| Goal | Supply session seeding and the minimum reading set, **without introducing a second copy of the rules** |
| Trigger | User authorization (already triggered) |
| Scale | S (documentation only) |
| Executing role | Documenter (⑧a) + master gates |
| Deliverable | The directive's **Appendix B.7 Slice Startup Card** + one line of guidance in §5.0 + traceability updates in A.1/C.3 + `_EN` mirror sync |
| Acceptance evidence | ① Appendix B.7 **contains no rule substance** (only templates, slots, and section-number references); ② the §5.0 guidance exists; ③ CN / `_EN` **line counts / bold markers / section-number sequences all equal** (**Note**: historical reading; the current G3 criteria are in §6.3 - changed-line-count symmetry + key-field alignment); ④ G1–G5 all green |
| Boundary | **No existing rule is changed**; no second document is added; code, contracts, and role files are not touched |
| ⑥ independent scan | MAI-Code-1.1-Flash, **1** effective call (within the §7.5 budget of 3): `INDEPENDENT SCAN: PASS`, Section A–E complete with substantive Section D/E; **independently judged** that "B.7 copies no rule substance" and "no CN/EN semantic drift" |
| **Did not pass independent final review** | Per OB-9, this slice's changes go together with v1.7 into **the input of the next ⑦** (without running a separate round) |

### C.0e Work-Tier and Startup-Authorization Slice `SW-5`

**Why it is necessary**: §1.1 says only "does not apply to: purely consultative Q&A; ad-hoc troubleshooting with a single command" and **does not define the most common case, the "routine change"** ⇒ for any repository change the master can only infer "slice", producing two bad paths: **over-engineering** (even a one-line copy edit dispatches paid ⑥/⑦ calls) and **implicit over-reach when there is no formal channel** (boundary content mixed into a small change ⇒ gates and review are bypassed). Measured basis: several small changes in v1.7/v1.8 (identifier disambiguation, ledger-row revision, encoding normalization) were in fact routine-tier work.

| Field | Content |
|---|---|
| ID | `SW-5` (v1.9) |
| Status | ✅ **Completed** (v1.9; commit `fb12ff9`, CI run `35632485437` green on both legs on the first attempt) |
| Goal | Give "non-slice work" a formal channel: three-tier startup (`[ROUTINE]` / `[SLICE]` / `[PILOT]`) ＋ trigger-word syntax ＋ tier determination and recommendation duty ＋ the routine tier's non-bypassable boundaries ＋ handling of over-reach / stopped states; **does not lower** the slice tier's review intensity |
| Dependency | v1.8 (Appendix B.7) committed as `30cf799`; user rulings Q1–Q6 |
| Scale | S–M (documentation only) |
| Executing role (**first role-trimming judgement**) | **Documenter (⑧a) ＋ master gates**; **① architecture and ⑤ pre-review exempted**. **Recorded downgrade judgement**: per the letter of §5.2, §1.6 touches authorization-contract semantics ⇒ §7.2 ③ does not hold ⇒ lands "Medium" (the Medium row requires test points); **substantively lands "Simple"**, three reasons: ① only the "Complex" row touches the technical-architecture dimension, and SW-5 has none; ② the "Medium" row requires test points, and SW-5 has none; ③ the authorization-contract semantics of §1.6 are **rule semantics**, not technical architecture. |
| Deliverable | §1.6 (new) ＋ §1.1 pointer ＋ one §5.0 startup-precondition sentence ＋ **the §5.0 step 0 upfront hard gate gains `SW-3`** (user ruling 2026-09-21) ＋ Appendix B.7 three-tier expansion ＋ §0.3 v1.9 row ＋ version number ＋ this ledger block ＋ **G4b allow-list sync** (`tools/gates/cjk-newwords.txt` gains the character 「竖」) ＋ strict `_EN` positional mirroring |
| Acceptance evidence (**target**) | ① the trigger-word syntax is **defined in §1.6 alone in the whole document**; ② boundary ① of §1.6 is **a reference to §7.2 ③** (not a redefinition); ③ the B.7 three-tier templates **contain no rule substance**; ④ `[ROUTINE]` cost = **0 paid calls**; ⑤ CN / `_EN` line counts / bold markers / pipes / heading line numbers / ID counts all equal; ⑥ G1–G5 all green; ⑦ ⑥ `PASS`; ⑧ ⑦ `PASS` (folded into the OB-9 open-item input) |
| Cost | **⑥ ×1 (MAI) ＋ ⑦ ×2 (Codex: 1 final review + 1 re-review) = 3 paid calls** (all within the §7.5 budget; the original budget statement was ⑥×1 + ⑦×1) |
| Boundary | **does not change** the G1–G5 criteria (the G4b allow-list merely gains 1 character through that criterion's own registration procedure, which is **data, not a criterion**); the §5.0 step 0 hard gate **only gains `SW-3`** (user ruling 2026-09-21, see §0.3 v1.9); does not change the §7.5 cap values; does not change the §5.1 pipeline sequence; does not change the role set / carriers; adds no document; does not touch code / CONTRACTS / schema / fixtures / CI |
| Note | the Appendix B.7 heading is renamed from "Slice Startup Card" to "**Startup Card**" (from v1.9 it also covers `[ROUTINE]`); the §5.0 guidance line is renamed accordingly |
| ⑥ independent scan (**actual**) | MAI-Code-1.1-Flash, **1** effective call (within the §7.5 budget of 3): `INDEPENDENT SCAN: PASS`, Section A–E complete; **0 findings** (Medium+ = 0) |
| ⑦ independent final review (**actual**) | Codex **×2** (1 final review + 1 re-review): the final review's `FINDINGS` **3 High + 4 Medium**, the re-review's `FINDINGS` **1 Medium**; **all 7 findings from the first round are "closed"**; the re-review Medium and the missing pieces in the review bundle have both been remediated; for the ledger-metadata-class exception see **OB-11** |
| **⑦ closed** | one final review plus one re-review, with no unclosed findings; the post-re-review ledger metadata fix is folded, per **OB-11**, into the next revision's ⑦ input |

### C.0f §6.3 extension (G6 / G7) and `tools/gates/` scripting slice `GATES-EXT`

**Why it was needed (the root cause behind OB-14)**: the DR-FIX pilot's defect self-capture rate was 2/5 (target ≥ 3/4), and the two missed items both belonged to the **ledger-metadata class** and the **review-package-completeness class**; G1–G5 have **no mechanical criterion** for either class, so they rest on human eyesight - and defects of exactly those classes were still being missed by measurement inside one and the same slice.

| Field | Content |
|---|---|
| Id | `GATES-EXT` (v1.10) |
| Status | ✅ **completed** (commit `8a2841f`, CI run `35718884321` first run **green on both legs**) |
| Objective | write the G6 / G7 criterion classes into §6.3 and run them from an **offline in-repo script** behind a CI hard gate; do not implement base G1 / G5 and do not change the semantics of G1–G5 |
| Dependencies | v1.9 committed as `fb12ff9`; the DR-FIX pilot conclusions (OB-13 / OB-14); the user's rulings on option F1 and the unattended authorization package |
| Size | M (documentation + code) |
| Roles | ① architecture → ② implementation → ③ testing → ④ gates → ⑤ pre-review → ⑥ scan → ⑦ final review → ⑧a / ⑧b documentation |
| Deliverables | the §6.3 text and the scripting-status paragraph, `tools/gates/**`, the CI Gates step and `fetch-depth`, the three `workflow_test.go` syncs, this ledger block, the validation report, the §12.4 / §12.5 registrations, a strict `_EN` positional mirror |
| Acceptance evidence | ① staged step 1 / 2 / 3 acceptance; ② the historical real defects reproduce (`range:b738e6b^..b738e6b` turns G6-2 red, `range:1b01115^..1b01115` turns G7-a red) while the fixed side is green; ③ two-way mutation checks (revert the criterion ⇒ the assertion goes red; restore ⇒ byte-identical); ④ **scope guards falsifiable one by one**; ⑤ the ⑥ and ⑦ (including round 3) conclusions |
| Cost | budget **①×1 + ⑤×1 + ⑥×1 + ⑦×2**; **realized ①×1 + ⑤×1 + ⑥×1 + ⑦×3** (⑦ round 3 is a user-pre-authorized exception; see the round-exception note in §12.4) |
| Boundaries | **not** changing the semantics of the G1–G5 criteria; **not** implementing base G1 / G5; no second criteria document; no change to the §5.1 sequence or the §7.5 caps; no change to the role set / carriers; the G4b whitelist gains this slice's legitimate new words only through the criterion's own registration flow (data, not criteria) |
| Notes | ① ⑦ round 3 reported **one Low** (`g5a.js` error routing) which was **left unfixed** because this slice's pre-authorized red lines include "fixing G1/G5 along the way"; it is **folded into the next §6.3 criterion-narrowing slice** (`GATES-TIGHTEN`; user ruling 2026-09-22); ② OB-19's boundary annotation (a user's temporary registration unrelated to this slice) is in the report; ③ the `internal/gates` flake is registered separately as **DR-6** per DR disposition boundary ① |
| ⑥ independent scan (**实得**) | MAI-Code-1.1-Flash **1 call**: `FINDINGS` **2 items** - one confirmed and fixed (G6-1 accepted only the bold shape of the completion marker, contradicting the criterion's enumeration and its own example line), one **rejected** (G7-b already implements the written rule verbatim) |
| ⑦ independent final review (**实得**) | Codex **3 calls**: round 1 `FINDINGS` 1 High + 1 Medium; round 2 1 Medium + 1 Low and judged the round-1 High as `PARTIAL`; round 3 `PASS WITH FIXES` with no Medium+. **Two-way coverage**: the ②f leading boundary and the NF-1 wording were both named and reviewed |

### C.0g Execution freedom, post-edit diagnostics and ⑧c wrap-up cleanup slice `DIRECTIVE-V1.11`

**Why it is needed**: ① the directive long had only an implicit convention, no written clause, for "what the master may decide on its own and what must stop and ask the user", so every boundary call depended on an ad-hoc reading; ② the fixed step of "post-edit diagnostics" was missing ⇒ warnings left behind after an edit could only be found by eye or by a later gate (the user's observation of 2026-09-23); ③ the wrap-up stage (after ⑧b and before ⑩) had no fixed whole-repo cleanup, so the temp directory / probe leftovers / untracked files depended on the master's self-awareness alone.

| Field | Content |
|---|---|
| Id | `DIRECTIVE-V1.11` |
| Status | ✅ **Complete** (2026-09-23; change and review evidence in `docs/validation/directive-v1.11.md`) |
| Goal | Write execution freedom and escalation boundaries plus the post-edit diagnostics scan into the directive, and fix the wrap-up cleanup as a pipeline stage |
| Trigger | User startup declaration of 2026-09-23 (`[SLICE]`) |
| Scale | Medium (documentation only; behavior change / multiple files / touches rule semantics) |
| Executing role | Master lands the text + ① plan + ⑤⑥⑦ review (⑧ not separately assigned) |
| Deliverable | §1.7, §3.4, §5.1 ⑧c, the §0.3 v1.11 row and the header version of `docs/DELIVERY_DIRECTIVE{,_EN}.md`, plus the §5.2 / §5.3 / §7.1 / §11.3 / Appendix B.7.4 companions; this slice's validation report |
| Acceptance evidence | CN and `_EN` positional: **what G3(a)+(b) check = changed-line-count symmetry + equal key-field counts**; full-text equality of line counts / bold markers / pipes / blank lines / heading positions is a **manual check item** (G3 does not cover it, see §6.3). **Correction note (2026-09-23)**: the original row's "checked by G3" was inaccurate and is corrected to the current §6.3 G3 criteria. `gate.js --all` all green; the first ⑧c run record (whole-repo diagnostics / IDE diagnostics summary / `gate.js --all --scope=tree` / the all-types encoding gate / the temp-directory and probe-leftover checklist) |
| Boundary | **No change to** the G1–G7 criteria themselves; no change to the existing §5.1 step sequence (⑧c is an addition); no change to the §7.5 cap; no change to the role set / carriers; code, CONTRACTS, schema, fixtures and CI are not touched |
| Note | The substantive choices made under execution freedom (e.g. rewriting a dead reference, adding section numbers to B.7.4, landing this ledger block) are listed one by one in the "change summary" section of this slice's validation report for ⑦ final review to judge whether any boundary was crossed |

### C.0h Three rounds of third-party independent review remediation and user-ruling landing slice `DIRECTIVE-V1.12`

**Why it is needed**: after v1.11 was frozen, three consecutive rounds of external independent review raised **64 findings** in total against the directive body (counting rule: accumulated round by round as reported, with cross-round recurrences counted more than once; severity breakdown **1 Critical + 11 High + 35 Medium + 17 Low**); all of them fall into the two classes of "following the letter yet silently going wrong" and "text out of sync with criteria". Two user rulings (the routine tier's production-code boundary, and the DR-FIX pilot re-evaluation) also had to be landed; the **fourth re-review (R6) raised 10 further findings** (1 High + 4 Medium + 5 Low, `tmp/DELIVERY_DIRECTIVE_v112_REVIEW_R6.md`) - N1 (re-evaluation object) and N2 (pilot deferral) were ruled by the user and landed by this slice (§12.1 and §12.2), N3, N4, N6 and N10 are ledger- and wording-class fixes that landed, N5 was measured as not reproducible, and N7, N8 and N9 are registered as OB-34 – OB-36 for the next slice.

| Field | Content |
|---|---|
| Id | `DIRECTIVE-V1.12` |
| Status | ✅ **Completed** (2026-09-23; the change and review evidence is in `docs/validation/directive-v1.12.md`) |
| Goal | Dispose of the **64 findings** accrued across the review rounds plus **22 findings from the two later re-reviews** (R5's 12 and R6's 10; **the item count is per the validation report's disposition table**), land the A3 / H4 user rulings and raise the directive to v1.12, registering the C1 event of the v1.12 review truthfully |
| Trigger | User startup declaration of 2026-09-23 (`[SLICE]`) |
| Scale | Medium (documentation only, multiple files, touching rule / authorization-contract semantics) |
| Executing role | Master lands the text (document-type slice: ②③ degenerate per the §5.2 note; the validation report under the Appendix C exception of §3.2) + ① plan + ⑤⑥⑦ review (⑧ not separately assigned, same reading as C.0g) |
| Deliverable | `docs/DELIVERY_DIRECTIVE{,_EN}.md` v1.12 + `docs/validation/directive-v1.12{,_EN}.md` |
| Acceptance evidence | CN / `_EN` positional checking by G3(a)+(b) (changed-line-count symmetry + equal key-field counts); `gate.js --all` all green; the ⑧c record; the **per-item disposition table** (item count per the validation report) in the validation report |
| Boundary | No change to the existing judgement logic (only two exemption clauses are added: A10's G1 stopped-state exemption and H5's ⑧c evidence-pack exemption); no change to the §5.2 matrix structure (blind-review wording only); no change to §7.5 numbers; no change to the role set / carriers; code, CONTRACTS, schema, fixtures and CI are not touched |
| Note | new observation items OB-28–OB-38; M5 lands per ①'s recommendation plus the user's ruling (§3.2 orchestration-artifact exception). **Fourth re-review (R6) disposition**: N1 and N2 landed per user ruling, N3, N4, N6 and N10 are ledger- and wording-class fixes that landed, N5 was measured as not reproducible, and N7, N8 and N9 are registered as OB-34 – OB-36; **every remediation of this fourth round came from an explicit user ruling or from ledger-metadata / wording-class fixes, with no self-initiated rule-semantics change, so under the OB-11 exception ⑦ is not re-run** (kept as ⑦'s historical input). **Pilot deferral**: this slice started as `[SLICE]` ⇒ it defers per §12.1, and the pilot metrics are declared by the next slice when it starts as `[PILOT]`. **Deferral count = 2**: the first is this slice; the second is `DIRECTIVE-REVIEW-SATURATION` (still a directive-governance / review-saturation slice, so it fails the pilot-selection criteria), which the user **ruled directly** on 2026-09-23 to continue executing without acting as a pilot (who ruled = the user; date = 2026-09-23; reason = not a pilotable behaviour-change slice; basis = the user's ruling text this round), recorded independently per OB-37's meta-rule and **not readable as a precedent** (see OB-40 in §12.4). **Fifth re-review (R7) disposition**: P1, P2 and P3 landed per user ruling (§12.1's deferral companions, §12.2's dual-track priority rule and the re-evaluation judge / anchor / count reading), P4 and P5 are registered as OB-37 and OB-38, and P6 to P9 are ledger- and wording-class fixes that landed; **later remediation in this slice must not use the OB-11 exception again** (its standard lapsed here because of P4, user ruling 2026-09-23), and that ⑦ third round is authorized in the same batch. **The C1 event**: the v1.12 review judged it Critical before the commit - the §0.3 v1.12 row claimed 34 changes while 19 had not landed; this slice **completes them** per the user's ruling rather than rewriting the log (see the same-titled subsection of the validation report). **Field-table applicability note**: this is a document-type governance slice; tier = [SLICE] (see the v1.12 log row in §0.3); probe quota = none; freedom-list path = `docs/validation/evidence/directive-v1.12-freedom-list.md`; the dependency / steps / current-gap fields are in the "execution pipeline" section of the validation report. |

### C.0i Review-saturation-criterion slice `DIRECTIVE-REVIEW-SATURATION`

**Why it is needed**: across v1.11 → v1.12 five rounds of third-party independent review **every round produced new findings**, and whether to stop could previously only be explained ad hoc; OB-32 (log ↔ body consistency) and OB-38 (the ledger status vocabulary and G6-1's vocabulary barely intersect) have each been demonstrated once ⇒ "when may this stop" and "the status-word closed set" had neither a mechanical criterion nor a sole authority.

| Field | Content |
|---|---|
| Id | `DIRECTIVE-REVIEW-SATURATION` |
| Status | 🔄 **In progress** (2026-09-23; the change and review evidence is in `docs/validation/directive-review-saturation.md`) |
| Goal | Define when a third-party independent re-review / ⑦ round N / re-review after remediation may stop (the five situations of quota / user-authorized exception / OB-11 / OB-37 / pilot deferral unified), land the two mechanical criteria **G6-5** (section-reference resolution) and **G6-6** (ledger status closed set), and close OB-34 / OB-35 / OB-36 / OB-41 within this slice |
| Dependency | v1.12 is committed; OB-40 (the user's pre-ruling on the second deferral); the `tools/gates/**` landed by `GATES-EXT` |
| Trigger | the user's startup declaration of 2026-09-23 (`[SLICE]`; the **second pilot deferral**, see OB-40 in §12.4) |
| Scale | M-L (pure `.md` plus `tools/gates` JS; multiple files, rule semantics included, two new criteria) |
| Steps | ① plan → ② implementation → ③ testing → ④ gates → ⑤ pre-review → ⑥ scan → ⑦ final review and re-review → ⑧a/⑧b → ⑨ native validation → ⑧c → ⑩ |
| Current gap | No stop decision has a sole authority (§5.3 / §7.5 / §12.1 / §12.2 each define a piece); the ledger status words have no closed set ⇒ the OB-32 and OB-38 classes have no mechanical coverage inside §12.4 |
| Deliverable | `docs/DELIVERY_DIRECTIVE{,_EN}.md` v1.13; G6-5 / G6-6 in `tools/gates/lib/criteria/g6.js`; test fixtures; the validation report; the freedom list |
| Acceptance evidence | `gate.js --all --scope=tree` fully green; `gate.js --selftest` fully green; CN / `_EN` positional checking by G3(a); positive and negative fixtures plus mutation assertions for both new sub-checks; the ⑨ native-validation readings |
| Boundary | no change to §7.5 numbers, the role set, the §5.2 matrix structure or G6-1's vocabulary; `internal/**`, `docs/CONTRACTS{,_EN}.md`, `schemas/**` and `.github/workflows/**` are not touched; OB-30 / OB-31 / OB-33 are split out |
| Tier | `[SLICE]` |
| Probe quota | none |
| Freedom-list path | `docs/validation/evidence/directive-review-saturation-freedom-list.md` |
| Note | the **second pilot deferral** (a direct user ruling on 2026-09-23, see OB-40 in §12.4); this slice does not use the OB-11 exception (OB-37 has narrowed its applicability); it also registers OB-42 (no mechanical criterion for the verdict-line closed set), OB-43 (mojibake in evidence capture) , OB-44 (the review layer writing to disk) and OB-45 (the review layer's report of its data source). |

### C.1 v3.1 section numbers → this directive

| Original v3.1 section | Content | Location in this directive |
|---|---|---|
| §2.0 allow-list overview | Three roles | §2.2 (expanded to 7 roles as triples) |
| §2.1 V4 Pro | Pre-analysis + pre-review | §2.2 (①/⑤) + appendix B.1/B.5 |
| §2.2 MAI / §2.2.1 Haiku | ⑥ and enabling Haiku | §7.2 / §7.2.3 |
| §2.3 Codex | ④ final review | §7.1 / §7.5 |
| §2.4 prohibited list | Blocklist | §2.4 (adds `sol-orchestrator`) |
| §3.1 do it yourself by default | Master boundary | §3.1 / §3.2 (changed to "produces only orchestration artifacts") |
| §3.2 get it right the first time | Fewer rejections | §7.4 |
| §3.3 expected iterations | Don't get discouraged | §9.1 (merged into fallback) |
| §3.4 protocol first | Fixed order | §1.4 |
| §3.5 cost awareness | Cost discipline | §7.5 |
| §3.6 discipline common to all tiers | One change per round / no scope expansion | §3.2 / §3.3 |
| §3.7 thinking-mode tiers | Tier rules | §8 (including the same-model inseparability limit) |
| §3.8 continuous execution and stop points | Continuous segments / stop points | §5.1 / §10.4 |
| §3.9 pre-authorized probe budget | probe budget | §7.6 |
| §3.10 ⑥ scanning layer | Three triggering conditions / cost ceiling | §7.2 |
| §3.11 conclusion handling / sanitization / residual | A11/A12 / blind review | §7.3 / §7.7 |
| §四 work discipline | Collaboration discipline | §3 / §9 |
| §四.11 workspace hygiene | Concurrent interference prohibited | §1.5 |
| §五 closed-loop discipline and the ⑤ failure state machine | Re-run that loop after a fix | §5.3 / §9.2 |
| §六 Git discipline | Authorization / push | §10 |
| §7.1 / §7.4 report structure and next steps | Report format | §11.3 / §11.3 (last item) |
| Appendix A / B cost ceilings and observation clauses | Budget and observation items | §7.5 / §12.4 |

> **Identifier mapping (2026-09-21 disambiguation)**: the layer identifier **`③.5`** that the old directive used is **uniformly renamed `⑥` (the independent scanning layer)** in this directive, to avoid confusion with pipeline step numbers; likewise, the old directive called the final review **`④`**, whereas this directive **uniformly calls it `⑦`**. The left column of the table above keeps the old identifier (because that is what is being mapped).

### C.2 References to rewrite after the switchover (executed by the switchover slice)

| File | Existing content | Action |
|---|---|---|
| `docs/t027/REMAINING_SLICES{,_EN}.md` | `见 docs/t027/FLASH_OPERATING_DIRECTIVE.md 第五章`, `按 §3.10 三条件`, the `§3.10` of the B3c observation item | Repoint to this directive §5.3 / §7.2 / §7.2 |
| `docs/DEV_PLAN{,_EN}.md` | If it contains `FLASH_OPERATING_DIRECTIVE` references | Same as above |
| `docs/DOCUMENTATION_PLAN{,_EN}.md` | Authoritative domain table | Add one row: delivery execution discipline = `DELIVERY_DIRECTIVE` |
| `.github/copilot-instructions.md` | Locating and authority section | Add a pointer (**pending user confirmation on whether to change it as well**) |
| `docs/t027/FLASH_OPERATING_DIRECTIVE*`, `FLASH_T027_BRIEFING*` | Full text | Keep the original text, add a supersede pointer (§A.2) |
| `docs/t027/B3A_RIG_DESIGN.md` and other historical working drafts | If they contain references | Historical drafts are not changed; covered by this mapping table |

**C.2 execution results (2026-09-21, `SW-1`)**: file-by-file measured results; any differences from the table above are flagged.

| File | Measured hits | Result |
|---|---|---|
| `docs/t027/REMAINING_SLICES.md` / `_EN.md` | 2 each | Changed: `第五章` → this directive §5.3; observation item `§3.10` → §7.2 (CN/EN symmetric) |
| `docs/DEV_PLAN.md` / `_EN.md` | **0** | Added a new-directive pointer in the "working method and constraints" section per §11.2 (the original table said "if it contains references"; measured as none) |
| `docs/DOCUMENTATION_PLAN.md` / `_EN.md` | — | Added the authoritative row "delivery execution discipline" |
| `.github/copilot-instructions.md` | — | Pointer added (executed within the same round's authorization scope) |
| `docs/t027/FLASH_OPERATING_DIRECTIVE.md`, `_EN.md`, `_v3.1.md`, `_v3.1_EN.md`, `FLASH_T027_BRIEFING.md`, `FLASH_T027_BRIEFING_v3.1.md` | ≥1 each | Supersede pointer added (file header only, body untouched) |
| `docs/t027/FLASH_OPERATING_DIRECTIVE_v3.1.html` | 1 | Added a historical-version banner (**beyond the original table**: `.html` was not listed) |
| `docs/t027/A6_PILOT_LAUNCH.md` | **4** | Supersede pointer added (**beyond the original table**: the original table classified it as "historical working draft, not changed"; the pointer was added because the `SW-1` acceptance assertion requires "every hit carries a superseded-by note", file header only) |
| `docs/validation/B2-EXTERNAL-ENFORCEMENT.md` / `_EN.md` | 1 each | Changed the "related documents" line (**beyond the original table**, same reason as above) |

**Bootstrap exemption (2026-09-21)**: `SW-1`/`SW-2` list the documenter (⑧a) as the executing role, but **the product-layer carriers are exactly the objects `SW-2` is meant to create** ⇒ there is a **bootstrap dependency**. These two slices are therefore **produced directly by the master** (§3.2's "do not originate product artifacts" does not apply here), **but independent review is not exempted**: once both are done they must still go through ⑦ independent final review (Codex), and **no checklist is attached in the first round**. This exemption **applies only during the directive switchover period**; from the pilot slice onward, the §5.0 role-readiness precheck and dispatch rules resume.

### C.3 List of items not inherited (present in v3.1, not inherited verbatim by this directive)

| Original v3.1 clause | Disposition | Reason |
|---|---|---|
| §3.11 "blind review is mandatory for selected slices" | **Inherited** (§7.7) | Residual governance must not be cut |
| §3.11 "sample proportionally (1 in every 4 slices)" | **Inherited** (new in §7.7) | Same as above |
| §3.11 blind-review budget "theoretically up to 4 times per slice" | **Not inherited** | This project-level directive instead enumerates per-category ceilings in §7.5 and no longer uses a "theoretical maximum" phrasing (to avoid it being treated as a quota) |
| v3.1 §12 "trial period (first 2 slices)" rollback conditions | **Rewritten** as the §12.1 quantitative metrics + the §12.2 rollback path | A project-level directive needs quantifiable targets, not a "trial period" description |
| §2.2.1 Haiku "time-limited clearance" details | **Inherited and tightened** into the three-level fallback of §7.2.3 | Keeps the capability, makes the failure path explicit |
| §四.11 concurrent-experiment interference prohibition | **Inherited** (§1.3 R1.1) | — |
| Appendix A cost-ceiling table | **Inherited and rewritten** as §7.5 | Paired with the measurement ground-truth source of R7.1 |
| The "hard gates (not bypassable)" section of Briefing v3.1 (B1→B2 / B3 / A7 etc.) | **Not inherited into the general directive** | They are **slice-level** preconditions and should be written in the slice definition and the ledger; writing them into the general document gives you "the slice changed while the document did not" |
