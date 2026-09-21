# SW-5 Work tiers and startup authorization (directive v1.9)

Date: 2026-09-22. Status: `COMPLETE`. Commit: `fb12ff9` (v1.9 body) + `b738e6b` (ledger closeout).
CI: run `35632485437`, `35633761736` (**both green on both legs on the first attempt**). Push: `origin/main` (**gitee not pushed**).

## 0. Honesty rules

1. This slice does **exactly** three things: the directive adds **§1.6 Work tiers and startup authorization** and wires it in (§1.1 / §5.0 / Appendix B.7 / §0.3 / version number / Appendix C.0e); it adds OB-10 / OB-11 and the §12.5 gap row; and, per G4b's **own** whitelist registration procedure, it adds 1 character (`竖`). It changes **no** G1–G5 criterion, no §7.5 cap value, no role set / carrier, and does not touch code / CONTRACTS / schema / fixtures / CI.
2. Every commit id, run id, gate result and count in this report is **taken from actual evidence**; anything not obtained is explicitly marked "**historically missing**", with **no inference and no fabrication**. This slice is a **documentation-only slice** (executed in degraded form per the §5.2 "Note (document-type slices)"), so it contains no code-class mutation testing.

## 1. Change summary

1. **§1.6 (new) Work tiers and startup authorization**: three-tier startup — `[ROUTINE]` (default, no trigger word) / `[SLICE]` / `[PILOT]`; **the trigger-word syntax is defined in this one place only in the whole document** (form, position, full-width/half-width, case, handling of multiple trigger words).
2. **The routine tier's non-bypassable boundaries**: ① is **a reference to §7.2 ③** (state machine / gate criteria / evidence model / roles / authorization contract) and **is not redefined**; ②–⑥ are this tier's **incremental boundaries** — the semantics of CONTRACTS / schema / fixtures, any `.agent.md`, code touching state machines / crash paths / ownership / concurrency, `.github/workflows/**`, and **adding / removing / upgrading / downgrading dependencies (including lock-file changes)**.
3. **The master's determination and recommendation duty**: **determine the tier first** and declare it on **the receipt's first line**; **may recommend, never upgrades unilaterally**; **one-way tier movement**, with no mid-flight downgrade — if the user asks for a downgrade ⇒ **stop the slice** and restart under `[ROUTINE]`; **three rules for a stopped state**: artifacts stay where they are but must not be committed as a completed slice, the report header must carry `**Status: stopped — incomplete**` and list the unfinished items, and the ledger row stays `🔄` or reverts to `⬜` (the user decides). **Not exempt from §10 because of the routine tier**; routine-tier record-keeping = diff + §6.1 + §6.3, **0 paid calls**.
4. **Wiring**: §1.1 gains a pointer; §5.0 gains a "startup precondition" sentence (the pipeline starts only when a trigger word appears) ＋ **the step 0 upfront hard gate gains `SW-3`**; Appendix B.7 is expanded per the three tiers (its heading changes from "Slice Startup Card" to "**Startup Card**", and B.7.1 kickoff message / B.7.2 opening receipt / B.7.3 stop-point receipt are all split by tier); §0.3 gains a v1.9 row; the header version number goes v1.8 → v1.9; Appendix C.0e records this slice (including the **first role-trimming downgrade judgement record**).
5. **Observation items and gate data**: it adds **OB-10** (the new §1.6 clauses lack G1–G5 mechanical criteria) and **OB-11** (the decision standard for skipping the ⑦ re-run after a ledger-metadata-class remediation); §12.5 gains a "decision standard for the §5.3 vs §7.5 conflict" gap row; **the G4b whitelist gains 1 character** (`竖`, which belongs to that criterion's own **data-registration procedure** — every character must carry a reason and a context excerpt, so this is **not a criterion change**).

## 2. Pipeline execution

| Stage | Role / carrier | Status | Note |
|---|---|---|---|
| ① architecture / ⑤ pre-review | — (V4 Pro) | **Exempt** | user ruling: this slice has no technical-architecture dimension, no test points, and touches no code |
| ② implementation | documenter (⑧a) | Done | degenerates into **document implementation** per the §5.2 "Note (document-type slices)" |
| ③ testing | documenter (⑧a) | Done | degenerates into a **document consistency check** |
| ④ master integration + gates | master | **All green** | §6.1 base gates + §6.3 G1–G5 |
| ⑥ independent scan | MAI-Code-1.1-Flash | **×1 `PASS`** | Sections A–E complete, **0 findings** |
| ⑦ independent final review | Codex (GPT-5.3-Codex) | **×2** (1 final review + 1 re-review) | 7 findings in round 1 → 1 in the re-review, all closed |
| ⑧a / ⑧b documentation write-back | documenter (⑧a) | Done | CN authority + `_EN` mirror strictly positional |
| ⑨ native verification | master | Degraded | encoding gate + G1–G5 + a **reference resolvability** check |
| ⑩ stop point | — | **Awaiting authorization** | commit / push request same-round authorization per §10 |

**Rework count**: after the ⑦ final review, remediation took **1** round (3 High + 4 Medium); after the re-review, remediation took **1** round (1 Medium); **⑦ invalid rounds = 0**.

**Environment and effort tier**: Windows master + PowerShell (BOM / EOL / gate-verification script); Go 1.22 toolchain running `gofmt` / `build` / `vet` / `test` natively; paid review carriers ⑥ `MAI-Code-1.1-Flash (copilot)`, ⑦ `GPT-5.3-Codex (copilot)`. **Effort tier (§8.2): historically missing** — no operator-side effort-tier record could be obtained; per §1.6 a tier declaration does **not replace** the effort-tier setting, so no inference is made.

**Anomaly and fallback record**: this slice dispatched no ①⑤, ran no native experiment, and had 0 ⑦ invalid rounds; **no fallback action count record was obtained ⇒ historically missing** (§9.2 definition). Anomalies that did occur and were remediated (not fallbacks): 7 ⑦ final-review findings, 1 ⑦ re-review finding, 3 ledger-class defects of ours within the same slice (§4.4), and the missing pieces in the ⑦ review bundle (M4).

## 3. Role-trimming downgrade judgement record (a **first** for this slice)

- Per **the letter of §5.2**: §1.6 touches **authorization-contract semantics** ⇒ **§7.2 ③ does not hold** ⇒ it lands in the "**Medium**" row, but the Medium row requires **test points** — this slice has none.
- **Substantively it lands in "Simple"**, for three reasons: ① only the "Complex" row touches the **technical-architecture dimension**, and SW-5 has none; ② the "Medium" row requires **test points**, and SW-5 has none; ③ the "authorization-contract semantics" of §1.6 are **rule semantics**, **not technical architecture**. ⇒ Conclusion: **① architecture and ⑤ pre-review are exempt**; ⑥, because §7.2 ③ does not hold, **may not be skipped**, and it was in fact run. This judgement and its reasons were also recorded in Appendix C.0e.

## 4. Review summary

**⑥ (MAI ×1)**: `INDEPENDENT SCAN: PASS`; Sections A–E complete; **0 findings** (Medium+ = 0); **no checklist was attached** in round 1 (isolation proof in the last paragraph of this section).

**⑦ final review (Codex ×1)**: `FINDINGS` = **3 High + 4 Medium**, each one closed:

| # | Sev | Finding (summary) | Disposition |
|---|---|---|---|
| H1 | High | the routine-tier boundary covered only "adding or removing dependencies", not upgrading / downgrading / lock files | changed to "adding / **removing / upgrading / downgrading** dependencies (including lock-file changes)"; the user ruled it was a "**completion of intent**" |
| H2 | High | C.0e claimed "does not change the §5.0 step 0 hard gate" while in fact changing it | the ledger's **three rows were rewritten truthfully** |
| H3 | High | inside C.0e, "awaiting ⑥ / ⑦ not passed" conflicted with "⑦ `PASS`" | split into "acceptance evidence (**target**)" and "⑥ independent scan (**actual**)" |
| M1 | Medium | no resolution rule for multiple trigger words on the first line | added "**does not take effect + stop and ask the user to resend**, must not choose one itself" |
| M2 | Medium | B.7 contained sentences with a **rule-like tone** | three places were **stripped of rule wording** |
| M3 | Medium | the deliverables did not declare the **3rd** changed file | added the G4b whitelist item |
| M4 | Medium | the T3/T4 review bundle **lacked the role-file hunk** | measured: the header-note-move bytes are in **`e0b08ec` (v1.5)** rather than `52c28af` (docs-only); `T4-note-move-e0b08ec.diff` was added |
| R1 | Rejected | suggested adding §6.3 mechanical criteria for the new §1.6 clauses | **beyond the SW-5 boundary, not adopted**; registered as **OB-10** |

**⑦ re-review (Codex ×1)**: `FINDINGS` = **1 Medium** (the C.0e status row "awaiting ⑥" **conflicted in tense** with the "⑥ actual `PASS`" below) ⇒ fixed (and the expiring `995/995` literal was removed). **Its independent re-check confirmed all 7 first-round findings as "closed"** (giving CN / `_EN` **positional line-number** evidence); **no semantic drift** was found between CN / `_EN`.

**Negative findings (recorded truthfully)**: ① **our ledger-class defects occurred 3 times within the same slice** — H3, the re-review Medium, and the pre-closeout self-check **missed 4 more places**; ② **the ⑦ round-1 review bundle we delivered had missing pieces** (M4); ③ **1 of the reviewer's proposed fixes was out of bounds** (OB-10).

**Review input package summary and isolation proof**: ⑥ and ⑦ **both had no checklist attached in round 1** (§7.1 / §7.2; §7.3 anchors hard isolation, and only the **re-review round** attaches a checklist, marked "read only after completing the independent scan"); input = the raw diff + a read-only repo, sanitized per §7.3 (**no item-by-item record obtained ⇒ historically missing**). Per **OB-9**, this round's ⑦ input **folded in together** the four classes of historical debt (the T3/T4 criterion-change diff, the v1.7 change diff, the v1.8 change diff, the `_EN` mirror); the review bundle's landing path is **historically missing**. Blind review (§7.7): this slice **did not touch** ownership / shutdown / identity semantics ⇒ mandatory sampling was not triggered; **the proportional-sampling record ⇒ historically missing**.

## 5. Gates and evidence

| Gate | Result |
|---|---|
| G1-a table-structure integrity | PASS ✅ |
| G1-b version-metadata consistency | PASS (header v1.9 = last §0.3 row v1.9) ✅ |
| G2 placeholder residue | PASS ✅ |
| G3 bilingual symmetry (a.1) / (a.2) / (b) | PASS ✅ |
| G4a encoding (BOM + LF) | PASS ✅ |
| G4b anomalous characters (difference set ⊖ whitelist = empty) | PASS ✅ |
| G5-a ① / ② | PASS ✅ |
| G5-b change-set `model` value gate | PASS ✅ |
| R2.5 header-note position anchoring | PASS ✅ |
| Base gates | `gofmt -l .` empty; `go build ./...` / `go vet ./...` / `go test ./...` = `ok=14 FAIL=0` ✅ |
| **TOTAL_FAIL** | **0** ✅ |

- CN / `_EN` **strictly positional**: **1000/1000** at commit time; **1001/1001** after closeout; **heading line-number sequence drift = 0**; **bold / pipe / blank-line / ID counts all equal item by item**; both files are UTF-8 **BOM + LF** (G4a checked file by file).
- How the gates were run: `tools/gates/` is **not yet scripted** (only the whitelist is on disk), so this slice's gates were run by the master through an **ad-hoc script**, **deleted once used** (see §8).

### 5.1 Falsifiability

This slice is a documentation-only slice, with **no code-class mutation testing** (**N/A**); the counter-inspection of the mechanical guards is as follows (**stated at the design level; not run item by item in this slice**) — if the G4b whitelist entry `竖` were removed ⇒ the difference set is non-empty ⇒ **G4b goes red**; if only CN or only `_EN` were changed ⇒ the changed-line counts are asymmetric ⇒ **G3 goes red**; if the header version number ≠ the last §0.3 row ⇒ **G1 goes red**; if the header note moved into the body (not the first non-empty line after the frontmatter) ⇒ **the R2.5 position criterion goes red**. **A11 symmetry spot-check / A12 independent corroboration: no record obtained ⇒ historically missing**. The two places known to have **no mechanical criterion** (the stopped-state report header marker, the `[PILOT]` pilot-metrics slot) = the gap pointed to by **OB-10**.

## 6. Cost and metering (R7.1 / §11.3 fields)

| role | model | stage | calls | note |
|---|---|---|---|---|
| master | `DeepSeek V4.1 Flash (deepseek)` | ⑧ gates / integration / report assembly | throughout | bootstrap-authored this slice's artifacts |
| documenter (⑧a) | `DeepSeek V4.1 Flash` | §1.6 draft / v1.9 mirror / remediation rows / OB-11 + §12.5 / closeout row | **5** | **not a §7.5 billable category** |
| documenter (⑧a) | `DeepSeek V4.1 Flash` | this report + its `_EN` mirror | **2** | as above |
| independent scan (⑥) | `MAI-Code-1.1-Flash (copilot)` | ⑥ independent scan | **1** | `INDEPENDENT SCAN: PASS` (within the budget of 3) |
| independent final review (⑦) | `GPT-5.3-Codex (copilot)` | ⑦ final review | **1** | `FINDINGS` 3 High + 4 Medium |
| independent final review (⑦) | `GPT-5.3-Codex (copilot)` | ⑦ re-review | **1** | `FINDINGS` 1 Medium |

- **Total paid calls = ⑥ ×1 (MAI) ＋ ⑦ ×2 (Codex: 1 final review + 1 re-review) = 3**, all within the §7.5 budget (⑥ ≤ 3; ⑦ 1 + 1); ①⑤ **exempt** ⇒ **V4 Pro 0 calls**; ⑦ **invalid rounds = 0**; probe budget used: **0** (no real external call).
- Difference from the budget declared at startup: the original declaration was ⑥×1 + ⑦×1, the actual was ⑥×1 + ⑦×2 (the re-review is a §5.3 duty), recorded truthfully in the Appendix C.0e "cost" row; the `startedAt` / `retryOf` fields are **not recorded item by item ⇒ historically missing** (`role` / `model` / `stage` / call counts are listed item by item per §7.1).

## 7. Exceptions and observation-item registration

| # | Content | Disposition |
|---|---|---|
| **OB-9** | the independent-verification gap for v1.7 / v1.8 / T3 / T4 | this round's ⑦ **folded in together** these four input classes per OB-9's requirement ⇒ the gap is **closed** |
| **OB-10** | the new §1.6 clauses lack G1–G5 mechanical criteria | **rejected + registered** (beyond the SW-5 boundary); includes the **empirical note** that 4 more places were missed before the commit |
| **OB-11** | the **exception decision standard** for skipping the ⑦ re-run after a ledger-metadata-class remediation | registered and folded into the next revision's ⑦ input; related lesson: a ledger **must not carry expiring literals** |
| **§12.5 gap row** | "the decision standard for the §5.3 vs §7.5 conflict" | can be drafted as soon as this slice lands; to be written formally into §12 at the next revision |
| **SW-4 CI record** | `SW-4`'s CI run has **no record in either** the §0.3 v1.8 row or Appendix C.0d | recorded as "**historically missing**" |

## 8. Explicitly not executed

1. **Scripting the G1–G5 checks in `tools/gates/` still awaits an authorized slice** (this slice's gates were run by the master through an ad-hoc script, deleted once used).
2. **`DEV_PLAN` write-back (scope ruling)**: this slice **does not write back to** `DEV_PLAN{,_EN}.md`. Reason: the `DEV_PLAN` format is **a T027 prerequisite status update / prfrail project implementation plan** and **does not cover "authoring the delivery directive" itself** (the user's explicit ruling of 2026-09-22); the authoritative ledger for directive-governance slices is the **directive's Appendix C.0–C.0e**, and the two ledgers are separate domains and are not mixed. ⇒ `SW-1` / `SW-2` / `SW-3` / `SW-4` **likewise did not write `DEV_PLAN`**; this is not an omission — they were never supposed to be written.
   **Measured basis (which overturned the original assumption)**: `DEV_PLAN{,_EN}.md` contains no `SW-*` section at all (only T027's A0–A7 / B3a–B3c); `SW-1`'s write-back at the time was the **§1 pointer to the new directive** (still present, `DEV_PLAN.md:17`, re-checkable with `git grep`), **not a slice section**.
   **To be discussed**: whether to add a sentence to §11.2 ("the write-back scope of directive-governance slices: Appendix C rather than `DEV_PLAN`"), left for evaluation at the next revision (registered as **OB-12**); **§11.2 is not touched in this slice** (a change to the rule body must not be slipped in during closeout).
3. **The DR fix slice (pilot) was not started**; **gitee was not pushed** (per the rules: gitee is a mirror backup only).
4. After the ⑦ re-review, the **ledger-metadata-class remediation did not re-run ⑦** (for the exception decision standard see **OB-11**).

## 9. Next-step recommendations

- Proceed to the **DR fix slice (pilot)**: per the user's characterisation "**a systemic timing problem within the same package**", **fix DR-2 / DR-3 / DR-4 / DR-5 together**; its input **must** contain ① a **cumulative occurrence-count table** and ② the **attribution judgement** of the interleaved `AUDIT TRAIL BROKEN` evidence from **3 different tests**.
- Fold in the to-do inputs of **OB-10 / OB-11** (whether to add §6.3 mechanical criteria, and the exception decision standard).
- The five §12.1 metrics use **B3c** as the baseline; recommended role effort tiers: the DR slice is a **behaviour-change class** slice (Medium / Complex) ⇒ ① (V4 Pro) and ⑤ are re-enabled, ⑥ is judged per §7.2 (a behaviour change is an enabling condition), and ⑦ keeps 1 final review + 1 re-review.
