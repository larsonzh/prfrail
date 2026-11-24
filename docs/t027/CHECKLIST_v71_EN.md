# Task execution checklist (cross-session snapshot v71 · after the P87 closeout · next item = P88)

Generated: 2026-11-25. HEAD = `2782a09`; directive = v1.72; the two P87 CI runs: `45294006169` (main round, both legs green, backfilled into the report) and `45356014951` (backfill round, observation only).

**Supersession**: this version supersedes v70 (`docs/t027/CHECKLIST_v70{,_EN}.md` is frozen and not retro-edited).

**Structure**: section 1 anchors → section 2 fact pack → section 3 completed slices → section 4 execution queue → section 5 observation items → section 6 pending rulings → sections 7 / 8 / 9 pointer sections → section 10 opening → section 11 one-line status; the rule bodies of sections 7–9 have moved into directive section 13 and are not copied here.

**Changes in this version (relative to v70)**:

- (1) Section 1 anchors rotate: HEAD `2782a09`; directive **v1.72**; SELFTEST **744/744**; the rolling window drops v68 and takes in v71 (v69 / v70 / v71), with the v69 / v70 line counts re-checked by this version.
- (2) Section 2 becomes the P87 fact pack (`CLOSED`, two commits and two CI runs, R = 0, D = 0, ADR-048).
- (3) Section 3 appends one P87 row (carrying v70's twenty-three rows = twenty-four rows in total).
- (4) Section 4 queue update: P87 leaves the queue; the next item = **P88 `GATES-PERF-ASSESS` (tooling slice; must be named by the user)**; the pre-S2 re-evaluation (S-G1 / S-G8) stays "later" in layer 1 (pending S2 enforcement wiring); the layer-3 `GATES-PERF-ASSESS` slot has been promoted into the P88 ledger section (removed from layer 3).
- (5) Section 5 observation items: DR max = **18**; mechanical readings OB max = **144**, ADR max = **48**; two new entries, "P87 Hydrogen precondition gap" and "P88 sequencing registration".
- (6) Section 6 pending rulings: the next ruling item = **naming the P88 start**.
- (7) Appendix A: "obligation hand-off" becomes **no new hand-off (P87 leaves no obligation behind; the Hydrogen gap is a record-only to-do)**; sections 10 / 11 rotate accordingly.

---

## 1. Current anchors (check these 7 numbers first)

**Rolling window = the latest 3 versions (v69 / v70 / v71)**; v68 has been dropped (same rule as the v52 drop).

| Item | Value | Check command |
|---|---|---|
| HEAD | `2782a09` | `git rev-parse --short HEAD` |
| origin/main | `2782a09` | `git rev-parse --short origin/main` |
| Directive version | **v1.72** | header of `docs/DELIVERY_DIRECTIVE.md` |
| SELFTEST | **744/744** | `node tools/gates/gate.js --selftest` |
| Three scopes | **all TOTAL_FAIL=0** | `gate.js --all --scope=index` / `tree` / `ci` |
| v69 whole-file line count | **209** (LF count, CN / EN symmetric; **re-checked by this version**) | count LF in `docs/t027/CHECKLIST_v69{,_EN}.md` |
| v70 whole-file line count | **203** (LF count, CN / EN symmetric; **re-checked by this version**) | count LF in `docs/t027/CHECKLIST_v70{,_EN}.md` |

> The HEAD = `2782a09` of this snapshot is the generation-time anchor after the P87 backfill commit; once v71 itself lands the HEAD moves forward, and `git rev-parse --short HEAD` governs then. The v71 line count is re-checked by v72.

> Mechanical readings: ADR max = **48**; OB max = **144**; DR max = **18**; pilot deferral max = the 25th.

> Rolling window = the latest 3 versions (v69 / v70 / v71); v68 has been dropped; if 2 or more items disagree ⇒ the repository governs, report the final readings first and then continue.

`tmp` fixed assertion: top level = `.gitkeep` + `gates-cache` + `P84-CI-FIX` (**referenced by the P84 report section 12, retained for archiving per section 1.5.1 (1)**) + `find-quotes.js` (no referencing residue, awaiting the tmp-lifecycle slice) + `P85-SECURITY-PROOF-SUBSYSTEM` (**referenced by the P85 evidence-pack `report.md`, retained for archiving per section 1.5.1 (1)**); `gates-cache` contains only `g4b`. P87's own temp directory was removed at (8)c (the backfill-round helper scripts were removed the same way).

---

## 2. P87 fact pack

- **Slice / size**: `DIRECTIVE-REVIEW-MODEL-POLICY` ((7) review-model selection matrix plus judging-threshold consensus plus the (2) hardening-round trigger) / **M**; master `DeepSeek V4.1 Flash` High; named by the user for a 2026-11-24 start.
- **Goal**: turn the (7) review-model selection rules into a standing four-tier selection matrix and fix the judging-threshold consensus; revise "no mixing within one ring" into a combination-trigger clause; add a new section 5.4 "(2) hardening-round trigger and fallback".
- **Verdict**: **`CLOSED`** (closed 2026-11-25); the directive was raised to **v1.72**; ADR-048.
- **Contract impact**: **no touch** to CONTRACTS / schema / contract fixtures / workflows / `internal/**` ((1) verdict); the directive body and the section 6.3 criterion implementation are touched.
- **Landing**: section 2.2 gains three Gemini final-review rows (`gemini-4-neon` / `carbon` / `hydrogen`) plus the Sol alternative sentence; section 5.2 gains the four-tier selection matrix (including the parallel-alternative tier), the judging-threshold consensus, the combination-trigger clause, the applicability-exclusion sentence and the effect mechanism; a new section 5.4 (S1–S4); section 7.5 gains three price rows and three notes; section 8.2 gains three effort rows and the effort authorization; section 6.3 G5-b closed set of **11**; a v1.72 change-log row in section 0.3; appendix C.0ay.
- **Criteria**: `g5b.js` `CLOSED_SET` of 11 (same order as section 6.3) plus `selftest.js` T132a–c plus four fixtures (three legal / one illegal).
- **Review chain**: (5) 1+1 (r1 `PASS WITH FIXES` → r2 `FINDINGS`; after the quota ran out the Medium was landed verbatim from the authoritative input and handed to (6)/(7) for verification) → (6) 2/3 (`FINDINGS` → `PASS`) → (7) blind 1 (`FINDINGS`) + final 1 (`PASS`, every finding closed) → (8)a / (9) / (8)b / (8)c → (10).
- **Metering**: **R = 0, D = 0, zero probes**; batch assist GLM x 2; the (7) 1+3 quota was not exhausted.
- **Two commits and CI**: main commit `8b64ce3` (+2134/−46, file count per numstat) → run `45294006169` **green on both legs** (windows-latest / ubuntu-latest) → backfill commit `2782a09` (+19/−13) → run `45356014951` (observation record only).
- **Documents**: report pair `docs/validation/directive-review-model-policy{,_EN}.md` (with an `artifacts` fence entry set and the `sixstage` (6) transcription); evidence pack `docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/` (19 payloads + MANIFEST + SHA256SUMS, including six verbatim review transcripts); the P87 ledger completion record.
- **Known precondition gap**: the light-tier clause (`Gemini 4 Hydrogen`) is landed while its comparison-review measurement is not done; it points at a later business slice (section 1.7.4 temporary authorization plus an OB record) and **must not be described as validated**.

### Appendix A: "obligation hand-off" (this version: **no new hand-off**)

- **Historical hand-off** (first case traced to P75 O-7): the in-row caliber mismatch of the core task chain at `s1-exit.md:52` was discharged by P76 change 1 plus the S-G2 correction (`0f80c1e`).
- **This version's hand-off**: **none** — P87 leaves no obligation behind; the Hydrogen gap is a **record-only to-do** (a later business slice fills it and backfills a note into the P87 report), and it is not a review or delivery obligation hand-off.
- **Convention note**: in use since v59; landing-slice entries must carry the "admission condition" orientation.

### Appendix B: "AT-23 stage 2 status" (unchanged)

- **Status**: A1 ruled "not applicable" (user ruling of 2026-11-09; not a permanent ban).
- **Executable object before S2**: none (R-3 = (A)).
- **Restart path**: if the user later rules to restart stage 2, real enforcement wiring (S2) or an alternative path must exist first.
- **Natural restart moment**: S2 starts → enforcement wired → S-G8 unlocks → R-3 re-evaluated.

---

## 3. Completed slices (recent baseline)

**v63 section 3 already carried P64–P80 (sixteen rows); v64 appended P81 (seventeen); v65 appended P82 (eighteen); v66 carried v65's eighteen and appended P79 (nineteen); v67 carried v66's nineteen and appended P83 (twenty); v68 carried v67's twenty and appended P84 (twenty-one); v69 carried v68's twenty-one and appended P86 (twenty-two); v70 carried v69's twenty-two (readings copied, no retro-edit) and appended P85 (twenty-three); v71 carries v70's twenty-three (readings copied, no retro-edit) and appends one P87 row. Twenty-four rows in total.**

| Version | Slice | Main commit | Verdict |
|---|---|---|---|
| v1.66 | P64 AT24-SUPERVISED-BLACK-BOX | `2dbc54a` | closed (AT-24 conditionally closed; R = 2; run `42948727363` green on first try) |
| v1.66 | P65 GATE-PIT-26-LANDING | `0eda074` | closed (appendix B.7.1 slot + A-part ledger to-do; R = 1; run `42969802819`) |
| v1.66 | P66 CHECKLIST-RULES-HOIST | `1650d5f` | closed (section 13 five sub-sections + new snapshot structure; R = 3, D = 2; run `43022796716`) |
| v1.67 | P67 GATES-DEBT-CLEANUP | `18b6ef9` + `365a930` | closed (g4b LRU + gate pitfalls 27–30 + section 12.4 clearing block + B.6 hardening; R = 3, D = 1; run `43124686702`) |
| v1.67 | P68 REMAINING-RULINGS | `e0d745c` + `665a31f` | closed (assessment; R = 1, D = 1; (7) PASS; five R items ruled; ADR-042; run `43177025096`) |
| v1.67 | P69 T027B-E2E-WIRING | `759c934` | closed (implementation; R = 0, D = 1; (7) PASS; E-1 registered; run `43264324421`) |
| v1.68 | P70 S1-STREAM-PROGRESS | `ca88fa3` | closed (implementation; R = 1, D = 2; (7) PASS; DR-14; run `43400311736` → `43403891205` → `43408438643`) |
| v1.68 | P71 S1-SCOPE-ANNOTATE | `6ee2c61` + `16c98f7` | closed (documentation; R = 0, D = 2; (7) blind PASS + final PASS; run `43463089586` + `43468806314`) |
| v1.68 | P72 REQ020-S2-RESTRICTION | `fb0ab57` + `406f4a0` | closed (documentation; R = 0, D = 0; (6) skip confirmed by (7); (7) final PASS WITH FIXES; run `43566904684` + `43568186355`) |
| v1.69 | P73 T027B-ADMISSION | `8c59712` + `f834ce5` | closed (assessment; R = 0, D = 0; (6) skip confirmed by (7); (7) final PASS; **directive raised to v1.69**; run `43622770750` + `43624661784`; OB-55 = option b) |
| v1.69 | P74 SANDBOX-INTERFACE | `3300364` + `150423d` | closed (implementation; R = 3, D = 0; (6) mandatory PASS; (7) final PASS; **R=3 first case triggering section 9.2 high risk**; run `43669748483` + `43672994996`) |
| v1.69 | P75 R3-REAL-CLI-ASSESSMENT | `ed2e26e` + `70a04cd` + `60146c0` | closed (assessment; R = 1, D = 1; (7) blind 1 + final 1; **R-3 = (A)**; ADR-045; run `43732384194` + `43734829063` + `43739409108`) |
| v1.69 | P76 S1-EXIT-CLOSEOUT-ASSESSMENT | `77f23d2` + `aca19af` + `0f80c1e` | closed (assessment; R = 1, D = 1; (7) final PASS WITH FIXES; Q2=(a), Q10=(C); ADR-046 / ADR-047; run `43785603244` + `43787083404` + `43793273520`) |
| v1.69 | P77 ADR010-APPROVAL | `33272be` + `fc00c4b` | closed (documentation; R = 1, D = 0; ADR-010 = `APPROVED`; S-G4 passed; run `43986942087` + `43997058549`) |
| v1.69 | P78 RELEASE-ZIP-BINDING | `c935713` + `0741fae` + `308b235` | closed (implementation; R = 0, D = 0; (7) final FINDINGS → re-review PASS; GitHub Release `v0.1.0` bound; S-G3 passed; run `44043705691` + `44060224254`) |
| v1.69 | P80 SG6-GATELEVEL-CONSOLIDATE | `e00d79b` + `58c2794` | closed (documentation; R = 2, D = 1; (7) narrowed read `RE-REVIEW: PASS`; S-G6 gate-level side holds; run `44211368890` + `44213070006`) |
| v1.70 | P81 DIRECTIVE-MODEL-UPGRADE | `45590e8` + `fd4d8fb` + `9ca6576` | closed (governance; (6) `MAI Code 2 Flash`, (7) `GPT-6.2 Sol`; SELFTEST 737 → 740; R = 1, D = 2; run `44281812336` + `44284200392` + `44392113028`) |
| v1.71 | P82 TMP-DISCIPLINE-OUT-OF-REPO | `0264bcd` + `adb8cb9` + `1015430` + `9d83089` + `9fd8cd6` | closed (governance; section 5.1 gains the out-of-repo self-heal sub-clause; OB-139 flipped to `Resolved`; **directive raised to v1.71**; R = 2, D = 3; all five commits CI-green) |
| v1.71 | P79 S2-PRE-INDEPENDENT-SECURITY-REVIEW | `a464a44` + `1b7f9c6` | closed (assessment; all three faces `FINDINGS`, F1–F8 handed to remediation slices; ADR-047 flipped to `COMPENSATED`; R = 2, D = 1; run `44585713851` + `44589887023`) |
| v1.71 | P83 SECURITY-FINDINGS-REMEDIATION | `ec532a4` + `c820743` + `4323503` | **partially closed** (F2 / F5 / F7 + T9 narrowing closed; F1 / F4 partly implemented; F3 / F6 / F8 deferred to P84 / P85; R = 2, D = 0; DR-15; first run red `44795140136` → `44801690679` green → observation `44805856028`) |
| v1.71 | P84 BOUNDDIR-RESIDUAL | `d597f1e` + `32bf234` + `509b01b` | **partially closed** (T5#4 fully closed, T4 mostly closed; the N13-UNC and evidence-attribution Low residuals **were closed by P86**; R = 3, D = 0; DR-16 / DR-17; first run red `44877224368` → fix `44926215282` green → observation `44928363670`) |
| v1.71 | P86 BOUNDDIR-UNC-EQUIVALENCE | `f1ae821` + `a5ef122` + `1eb63d9` + `565fc9f` | **closed** (N13-UNC + W10 closed, F1 fully closed; R = 0, D = 0; DR-18; first run `44962802936` green → guard round red `44965434402` → fix round `44967047189` green → observation `44969239294`) |
| v1.71 | P85 SECURITY-PROOF-SUBSYSTEM | `cc7d90f` + `6480bb1` + `18b2613` + `ef1e1f2` | **closed** (T2 / T3 / T8a closed, F3 / F6 / F8 settled, **F1–F8 all closed**; R = 8, D = 0; OB-140/141/142/143/144; first run `45201921016` had two environmental windows-leg failures → R=8 fix `6480bb1` → `45208001981` green on both legs) |
| v1.72 | P87 DIRECTIVE-REVIEW-MODEL-POLICY | `8b64ce3` + `2782a09` | **closed** (directive raised to **v1.72**; (7) four-tier selection matrix + judging-threshold consensus + new section 5.4; G5-b closed set of 11; ADR-048; R = 0, D = 0; run `45294006169` green on both legs + observation `45356014951`; Hydrogen precondition gap registered for a later slice) |

> **Section 3 D caliber** (same as v61): the "D = n" in this table is shorthand for the temp-file discipline violation count (section 5.1); P87 is **D = 0** throughout. **The D count and the section 9.2 R count are two independent sequences**.

---

## 4. Execution queue

### Layer 0 upfront triage (✅ already out of the queue; no additions in this version)

- P5 / P6 were settled by P57; layer 0 is fully drained with no residue.

### Layer 1 mainline hard gates

- Items 8–12: all ✅ out of the queue (see v57 / v58 section 4; to trace back, follow the reference chain of v57 section 4 layer by layer).
- **`T027-B` (defensive-refactor phase) = complete (partly not applicable)**: C1–C7 / N1–N4 / R-1 / R-2 / R-3 all handled (R-3 = (A)). **The mainline is fully drained**.
- **`S1 Exit closeout` (gate-level side done; security remediation F1–F8 all settled)**: P76 produced the closeout roadmap; P77 closed S-G4; P78 closed S-G3; P80 judged the S-G6 gate-level side to hold; the independent security review was completed by P79 and P83 / P84 / P86 / P85 finished the F1–F8 fixes. The eight gates are as follows (P87 did not change the gate table):

| Gate | Criterion | Status | Basis |
|---|---|---|---|
| S-G1 | the CLI can really start an admitted AI candidate and bill for it | **fails** (Q2=(a) literal criterion) | P76 section 4; ADR-046 |
| S-G2 | `supervised-black-box` implemented and the AT-24 verdict holds | **passes** (T028 `COMPLETED` + AT-24 conditionally closed) | P76 section 4; `0f80c1e`; ADR-046 |
| S-G3 | the final ZIP is bound to the release record and re-verified | **passes** (P78 closed) | P78 report sections 5 / 7 |
| S-G4 | ADR-010 approved by the product owner | **passes** (user approval 2026-11-16; P77 closed) | P77 report section 5 |
| S-G5 | the product owner grants approval | **release part approved (first step)** | ADR-046 Q10=(C); P76 section 12 |
| S-G6 | dependency fields satisfied and the four remaining gates closed | **gate-level side holds** (P80 combined verdict) | P80 report sections 5 / 6 |
| S-G7 | all ATs pass, T018 complete, section 14 fully exited | unreachable (synthetic gate) | `s1-exit.md:52` / `:56`–`:58` / `:60` |
| S-G8 | the supervised path runs end to end with the real CLI and a real candidate | blocked until S2 enforcement wiring | P75 A2 / ADR-045; P76 section 6 |

**Next item = P88 `GATES-PERF-ASSESS`**

- Tooling slice; **must be named by the user**; the start message records the four elements; the ledger section is registered (⬜ not started, `docs/t027/REMAINING_SLICES{,_EN}.md`).
- Draft essentials: measure the gate timing composition (SELFTEST / three scopes / base-trio execution counts and shares) plus lightweight tuning (de-duplication / cache / parallelism), with **the criteria's verdicts byte-identical**; S tier, (7) 1+1, no blind review.
- Later = the pre-S2 re-evaluation (S-G1 / S-G8) (the next station of the T027 mainline, awaiting the user's sequencing):
  - prerequisite = S2 enforcement wiring; the independent security review prerequisite was met by P79, and P83 / P84 / P86 / P85 finished the F1–F8 fixes ⇒ the ADR-047 fallback clause no longer blocks;
  - `S-G8` stays blocked until S2 enforcement wiring (P73 A2 / P75 A2 / ADR-044 / ADR-045); the `s1-exit{,_EN}.md` gate table is unchanged.
- **Post-S2 re-evaluation items**: **S-G5 step two (overall S1 Exit approval) + S-G7**, in the order of the P80 report section 7.

### Layer 2 routine governance

- P72 / P73 / P74 / P75 / P76 / P77 / P78 / P80 / P81 / P82 / P79 / P83 / P84 / P86 / P85 / **P87** are all ✅ out of the queue (layer 2 is **fully drained**).
- Full section 1.2 field definitions for each slice are in the matching section of `docs/t027/REMAINING_SLICES{,_EN}.md`.

### Layer 3 low-priority candidates

| Field | Content |
|---|---|
| Identifier | `OB-50` / `spawnGit` / `ci frozen cost` / the 22 legacy `%TEMP%` items / `GATES-EXEMPT` |
| Current status (verified item by item) | OB-50 = **resolved** (section 12.4, v1.41); spawnGit = the git-call helper of `tools/gates/selftest.js` (the spawn wrapper of `lib/ctx.js`, stderr collected into the result object and not printed); a candidate OB (diagnostic), source = P27 report section 5 item 5, no section 12.4 entry; ci frozen cost = no section 12.4 entry, source = P27 report section 7 "candidate OB (performance)", not formally registered; the `%TEMP%` legacy items = OB-98 resolved / OB-99 and OB-100 under observation; GATES-EXEMPT = **not started** (no section 12.4 entry; a P33 (1)-3(iv) split, touching four criterion files) |
| Contents | item by item in the matching section 12.4 rows; **`GATES-PERF-ASSESS` was promoted into the P88 ledger section in v71 (⬜ not started), so this layer no longer lists it** |
| Dependency | the long park applies only to the open items (spawnGit / ci frozen cost / the \%TEMP\% timing-shape OB-99 and OB-100 / GATES-EXEMPT), none with a concrete trigger |
| Size | varies (S through M-L) |
| (7) budget | none |
| Boundaries | this version advances no parked item; the registration slots stay |
| Start condition | the user must name it explicitly |

**Queue discipline**: layer boundaries and numbering discipline are the same as v60; this version raises no version.

---

## 5. Observation items

- Observing (**13**): `OB-1 / 7 / 10 / 33 / 87 / 99 / 100 / 104 / 105 / 107 / 111 / 113 / 135` (inherited from v61).
- To be discussed (9): `OB-19 / 21 / 52 / 53 / 57 / 58 / 59 / 61 / 75`.
- To do (1): `OB-69`. Non-terminal total **23** (13 + 9 + 1), matching this version's verified state.
- **DR register**: DR max = **18** (unchanged after v69; P85 / P86 / P87 add no new DR).
- **P87 Hydrogen precondition gap (record-only)**: the light-tier clause (`Gemini 4 Hydrogen`) is landed while its comparison-review measurement is not done; it points at a later business slice (section 1.7.4 temporary authorization plus an OB record) and **must not be described as validated**; after completion a note is backfilled into `docs/validation/directive-review-model-policy{,_EN}.md`. Registered in the report section 7 and the ledger P87 section.
- **P88 sequencing registration**: `GATES-PERF-ASSESS` was promoted from layer 3 to P88 (⬜ not started); its full section 1.2 fields are already in the ledger; the start requires the user to name it.
- Mechanical readings: section 12.4 max OB = **144**; ADR max = **48**; DR max = **18**; pilot deferral max = the 25th.
- **Removed in this version**: none. **Added in this version**: the two entries "P87 Hydrogen precondition gap" and "P88 sequencing registration".

---

## 6. Pending user rulings

- **This table is currently empty** (the P87 completion record, the P88 registration and the v71 snapshot generation were authorized together on 2026-11-25).
- **Next ruling item (expected)**: **naming the P88 start** (`GATES-PERF-ASSESS`; the user must name it; the start message records the four elements; the draft is already in the ledger).

---

## 7 / 8 / 9. Pointer sections

Same structure as v60: section 7 start and step-0 general rules → directive section 13.1; section 8 cross-slice collaboration discipline → section 13.2; section 9 checklist meta-rules and gate pitfalls → sections 13.3 / 13.4. **The rule bodies have moved into directive section 13 and are not copied here**.

---

## 10. One-paragraph opening for a new session

> Below is my project state snapshot (v71). Please check the 7 anchors of section 1 first; if they match, continue with the section 4 queue, otherwise the repository governs.
>
> [paste the full v71 checklist]
>
> [supporting documents] `REMAINING_SLICES{,_EN}.md` (P87 closed with a completion record; P88 registered as the next item, awaiting the start naming), `DELIVERY_DIRECTIVE.md` (**v1.72**), `directive-review-model-policy{,_EN}.md` (the P87 report pair), `docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/` (the P87 evidence pack), `ADR_REGISTER{,_EN}.md` (ADR-048), `DOCUMENTATION_PLAN{,_EN}.md` (the hand-off entry pointer now points at v71), `CHECKLIST_v70{,_EN}.md` (frozen, superseded by v71).
>
> Current: P87 closed (directive v1.72; the (7) four-tier selection matrix + judging-threshold consensus + section 5.4 landed; G5-b closed set of 11); next step = the P88 tooling slice (the user must name it); the pre-S2 re-evaluation (S-G1 / S-G8) = the next station of the T027 mainline, awaiting the user's sequencing (section 4 "later").

## 11. One-line status

P87 **closed** (two commits `8b64ce3` / `2782a09` pushed to origin; the main CI run `45294006169` is green on both legs and the backfill run `45356014951` is observation only; R = 0, D = 0; directive v1.72; ADR-048; the Hydrogen precondition gap is a record-only to-do for a later slice); next step = **the P88 tooling slice `GATES-PERF-ASSESS`** (the user must name it); the pre-S2 re-evaluation (S-G1 / S-G8) = the next station of the T027 mainline, awaiting the user's sequencing (section 4 "later").
