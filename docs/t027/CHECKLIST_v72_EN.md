# Task execution checklist (cross-session snapshot v72 · after the P88 closeout · next item = TMP-ARCHIVE-P84-P85)

Generated: 2026-11-25. HEAD = `9b44dbf`; directive = v1.73; the two P88 CI runs: `45523666470` (main round, both legs green, backfilled into the report) and `45524509376` (backfill round, observation only).

**Supersession**: this version supersedes v71 (`docs/t027/CHECKLIST_v71{,_EN}.md` is frozen and not retro-edited).

**Structure**: section 1 anchors → section 2 fact pack → section 3 completed slices → section 4 execution queue → section 5 observation items → section 6 pending rulings → sections 7 / 8 / 9 pointer sections → section 10 opening → section 11 one-line status; the rule bodies of sections 7–9 have moved into directive section 13 and are not copied here.

**Changes in this version (relative to v71)**:

- (1) Section 1 anchors rotate: HEAD `9b44dbf`; directive **v1.73**; SELFTEST **744/744**; the rolling window drops v69 and takes in v72 (v70 / v71 / v72), with the v70 / v71 line counts re-checked by this version.
- (2) Section 2 becomes the P88 fact pack (`CLOSED`, two commits and two CI runs, R = 0, D = 0, directive v1.73).
- (3) Section 3 appends one P88 row (carrying v71's twenty-four rows = twenty-five rows in total).
- (4) Section 4 queue update: P88 leaves the queue; the next item = **`TMP-ARCHIVE-P84-P85` (tier-S tooling slice; must be named by the user)**; the pre-S2 re-evaluation (S-G1 / S-G8) stays "later" in layer 1; layer 3 gains four performance candidates in ascending risk order: `GATES-CI-PERF` / `GATES-SELFTEST-SPLIT` / `GATES-CRITERIA-ASYNC` / `GATES-BASE-CACHE`.
- (5) Section 5 observation items: "P88 sequencing registration" is removed; "four performance candidates registered in layer 3" and "the three tmp leftovers are carried by `TMP-ARCHIVE-P84-P85`" are added; DR max = **18**, OB max = **144**, ADR max = **48** are carried over.
- (6) Section 6 pending rulings: the next ruling item = **naming the `TMP-ARCHIVE-P84-P85` start**.
- (7) Appendix A: "obligation hand-off" becomes **no new hand-off (P88 leaves no obligation behind; the performance candidates are layer-3 records)**; sections 10 / 11 rotate accordingly.

---

## 1. Current anchors (check these 7 numbers first)

**Rolling window = the latest 3 versions (v70 / v71 / v72)**; v69 has been dropped (same rule as the v52 drop).

| Item | Value | Check command |
|---|---|---|
| HEAD | `9b44dbf` | `git rev-parse --short HEAD` |
| origin/main | `9b44dbf` | `git rev-parse --short origin/main` |
| Directive version | **v1.73** | header of `docs/DELIVERY_DIRECTIVE.md` |
| SELFTEST | **744/744** | `node tools/gates/gate.js --selftest` |
| Three scopes | **all TOTAL_FAIL=0** | `gate.js --all --scope=index` / `tree` / `ci` |
| v70 whole-file line count | **203** (LF count, CN / EN symmetric; **re-checked by this version**) | count LF in `docs/t027/CHECKLIST_v70{,_EN}.md` |
| v71 whole-file line count | **201** (LF count, CN / EN symmetric; **re-checked by this version**) | count LF in `docs/t027/CHECKLIST_v71{,_EN}.md` |

> The HEAD = `9b44dbf` of this snapshot is the generation-time anchor after the P88 backfill commit; once v72 itself lands the HEAD moves forward, and `git rev-parse --short HEAD` governs then. The v72 line count is re-checked by v73.

> Mechanical readings: ADR max = **48**; OB max = **144**; DR max = **18**; pilot deferral max = the 25th.

> Rolling window = the latest 3 versions (v70 / v71 / v72); v69 has been dropped; if 2 or more items disagree ⇒ the repository governs, report the final readings first and then continue.

`tmp` fixed assertion: top level = `.gitkeep` + `gates-cache` + `P84-CI-FIX` (**referenced by the P84 report section 12, retained for archiving per section 1.5.1 (1), carried by `TMP-ARCHIVE-P84-P85`**) + `find-quotes.js` (no referencing residue, awaiting the tmp-lifecycle slice) + `P85-SECURITY-PROOF-SUBSYSTEM` (**referenced by the P85 evidence-pack `report.md`, retained for archiving per section 1.5.1 (1), carried by `TMP-ARCHIVE-P84-P85`**); `gates-cache` contains only `g4b`. P88's own temp directory was removed at (8)c (the backfill-round helper scripts were removed the same way).

---

## 2. P88 fact pack

- **Slice / size**: `GATES-PERF-ASSESS` (gate timing-composition measurement plus lightweight efficiency tuning) / **S**; master `DeepSeek V4.1 Flash` High; named by the user for a 2026-11-25 start.
- **Goal**: measure and land the gate timing composition with evidence (SELFTEST / the three scopes / the execution counts and shares of the four base checks); land lightweight tuning on that composition (batched blob prefetch / raw-reading memos) with **byte-identical criterion verdicts**.
- **Verdict**: **`CLOSED`** (closed 2026-11-25); the directive was raised to **v1.73**; appendix C.0az.
- **Contract impact**: **no touch** to CONTRACTS / schema / contract fixtures / workflows / `internal/**` / the CI workflow ((1) verdict); the criterion **implementation** file `tools/gates/lib/ctx.js` is touched while the criterion bodies and decision logic stay unchanged.
- **Landing**: `ctx.js` batched blob prefetch (ONE `git cat-file --batch` with 12 fail-open paths that discard the whole batch and fall back to per-file reads) plus raw-reading memos (`readBytes` / `exists` / `trackedAt` / `numstat`); `g4b.js` / `gate.js` / `selftest.js` are byte-identical to HEAD and the cache key `key12` does not move; a v1.73 changelog entry in section 0.3; appendix C.0az.
- **Criteria and acceptance**: five same-input A/B pairs (tree full / index / range / ci cold corpus / selftest) are byte-identical (including the `key12` and `corpus=` of `INFO G4b-cache`); `SELFTEST` stays 744/744; the three scopes report `TOTAL_FAIL=0`; the report pair mirrors strictly with equal line counts on both sides.
- **Performance result**: the ci criteria leg on a cold corpus drops from 170277 ms to 12929 ms (13.17x) and the CI critical path (`--selftest && --all --scope=ci`) from 380.3 s to 141.8 s (2.68x); parallelising the four base checks was **not** done (recorded as a larger-refactor suggestion) and the breakdown of the 265 s `--selftest` is handed to a later slice.
- **Review chain**: (1) x 1 (`DeepSeek V4 Pro`) → (2) x 1 (`prfrail-implementer`) → (3) x 1 (`prfrail-tester`) → (5) 1/2 `PASS` → (6) 2/3 (r1 `FINDINGS`, r2 `PASS` after the paired evidence was supplied) → (7) one final review of 2 `FINAL REVIEW: PASS` (the re-review quota was not used) → (8)a / (9) / (8)b / (8)c → (10).
- **Metering**: **R = 0, D = 0, zero probes**; one GLM batch-assist call (report-pair mirror check, `MIRROR: PASS`); no blind review.
- **Two commits and CI**: main commit `03c45f7` (31 artifacts) → run `45523666470` **green on both legs** (ubuntu 3m09s / windows 2m58s, green on the first attempt) → backfill commit `9b44dbf` (4 files +18/-4) → run `45524509376` (observation record only).
- **Documents**: report pair `docs/validation/gates-perf-assess{,_EN}.md` (with the `artifacts` fence and the `sixstage` (6) transcription); evidence pack `docs/validation/evidence/P88-GATES-PERF-ASSESS/` (22 payloads + MANIFEST + SHA256SUMS, including four verbatim review transcripts and one batch-assist transcript); the P88 ledger completion record.
- **tmp and leftovers**: P88's own temp directory and backfill scripts are cleared (P88 residue under `tmp` = 0); the three leftovers from other slices (`P84-CI-FIX` / `P85-SECURITY-PROOF-SUBSYSTEM` / `find-quotes.js`) are recorded as they are, and the first two are carried by `TMP-ARCHIVE-P84-P85`.

### Appendix A: "obligation hand-off" (this version: **no new hand-off**)

- **Historical hand-off** (first case traced to P75 O-7): the in-row caliber mismatch of the core task chain at `s1-exit.md:52` was discharged by P76 change 1 plus the S-G2 correction (`0f80c1e`).
- **Obligation hand-off in this version**: **no new hand-off** — P88 leaves no obligation behind; the four performance candidates are **layer-3 records** and create no review or delivery obligation hand-off.
- **Convention note**: in use since v59; landing-slice entries must carry the "admission condition" orientation.

### Appendix B: "AT-23 stage 2 status" (unchanged)

- **Status**: A1 ruled "not applicable" (user ruling of 2026-11-09; not a permanent ban).
- **Executable object before S2**: none (R-3 = (A)).
- **Restart path**: if the user later rules to restart stage 2, real enforcement wiring (S2) or an alternative path must exist first.
- **Natural restart moment**: S2 starts → enforcement wired → S-G8 unlocks → R-3 re-evaluated.

---

## 3. Completed slices (recent baseline)

**v63 section 3 already carried P64–P80 (sixteen rows); v64 appended P81 (seventeen); v65 appended P82 (eighteen); v66 carried v65's eighteen and appended P79 (nineteen); v67 carried v66's nineteen and appended P83 (twenty); v68 carried v67's twenty and appended P84 (twenty-one); v69 carried v68's twenty-one and appended P86 (twenty-two); v70 carried v69's twenty-two (readings copied, no retro-edit) and appended P85 (twenty-three); v71 carries v70's twenty-three (readings copied, no retro-edit) and appends one P87 row. Twenty-four rows in total; v72 carries v71's twenty-four rows (readings left as recorded, not retro-edited) and appends one P88 row, for twenty-five rows in total.**

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
| v1.73 | P88 GATES-PERF-ASSESS | `03c45f7` + `9b44dbf` | **Closed** (tooling slice; criterion bodies and decision logic untouched, five same-input A/B pairs byte-identical; SELFTEST 744/744; three scopes 0; CI critical path 380.3 to 141.8 s; R = 0, D = 0; run `45523666470` green on both legs on the first attempt plus observation `45524509376`) |

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

**Next item = `TMP-ARCHIVE-P84-P85`**

- Tooling slice; **must be named by the user**; the start message records the four elements; the ledger section is registered (⬜ not started, `docs/t027/REMAINING_SLICES{,_EN}.md`).
- Scope: archive `tmp/P84-CI-FIX` and `tmp/P85-SECURITY-PROOF-SUBSYSTEM` into their own evidence packs, rewrite the references and remove both directories; the reference points are whatever grep measures (the measured list is in that ledger section).
- Boundary: **leave `find-quotes.js` alone** (it belongs to the tmp-lifecycle slice), do not wipe tmp wholesale and do not touch the directive body; tier S; dependency = the v72 snapshot closeout.
- Later = the pre-S2 re-evaluation (S-G1 / S-G8) (the next T027 mainline station, pending user scheduling):
  - Prerequisite = S2 enforcement wiring; the independent security-review prerequisite was fulfilled by P79 and F1-F8 were fixed by P83 / P84 / P86 / P85 => the ADR-047 rollback clause no longer blocks;
  - `S-G8` stays blocked until S2 enforcement wiring (P73 A2 / P75 A2 / ADR-044 / ADR-045); the `s1-exit{,_EN}.md` gate table is unchanged.
- **Post-S2 re-evaluation items**: **S-G5 step two (the overall S1 Exit approval) + S-G7**, in the order of the P80 report section 7.

### Layer 2 routine governance

- P72 / P73 / P74 / P75 / P76 / P77 / P78 / P80 / P81 / P82 / P79 / P83 / P84 / P86 / P85 / **P87** / **P88** are all ✅ out of the queue (layer 2 is **fully drained**).
- Full section 1.2 field definitions for each slice are in the matching section of `docs/t027/REMAINING_SLICES{,_EN}.md`.

### Layer 3 low-priority alternatives

| Field | Content |
|---|---|
| Identifier | `OB-50` / `spawnGit` / `ci freeze cost` / `%TEMP% legacy 22 items` / `GATES-EXEMPT` |
| Current state (item by item) | OB-50 = **disposed** (section 12.4, v1.41); spawnGit = the git-call helper of `tools/gates/selftest.js` (the spawn wrapper in `lib/ctx.js`, stderr folded into the result object and never printed); a candidate OB (diagnostics), source = the P27 report section 5 item 5, no section 12.4 entry; ci freeze cost = no section 12.4 entry, source = the P27 report section 7 "candidate OB (performance)", never formally registered, and its evaluation is carried by `GATES-CI-PERF` from this version on; %TEMP% legacy = OB-98 disposed / OB-99 and OB-100 under observation; GATES-EXEMPT = **not started** (no section 12.4 entry; split off by P33 (1)-3(iv), touches four criterion files) |
| Contents | item by item in the matching section 12.4 rows; `GATES-PERF-ASSESS` was promoted into the P88 ledger section at v71 and closed before v72 (2026-11-25), so this layer no longer lists it; this version adds four performance candidates (table below) |
| Dependency | the standing hold covers only the undisposed items (spawnGit / the %TEMP% timing shapes OB-99 and OB-100 / GATES-EXEMPT), none with a defined trigger |
| Size | varies (S through M-L) |
| (7) budget | none |
| Boundary | this version advances no held item and keeps the registration slot; the four performance candidates are records and constitute no project start |
| Start condition | must be named explicitly by the user |

**Performance candidates (new in this version, ascending risk; all records, all requiring the user to name them)**:

| Identifier | One-line scope | Relation to existing items / prerequisite |
|---|---|---|
| `GATES-CI-PERF` | CI-workflow side: evaluate running the Gates step as `--all --skip-base --scope=ci`, plus persisting the `go build` / `go test` caches (the latter is an extension beyond this slice's report section 7 and touches no criterion) | same origin as the existing "ci freeze cost" => **evaluate together at project start** to avoid duplication; touches the CI workflow (needs user authorization) |
| `GATES-SELFTEST-SPLIT` | `--selftest` per-suite timing and fixture-read memoisation plus a G4b inner-loop increment (gain must be argued first) plus a built-in `--timings` | touches `selftest.js` / `g4b.js` (criterion implementation) => (6)/(7) mandatory; `--timings` is an output-protocol surface |
| `GATES-CRITERIA-ASYNC` | make G6 / G7 / G8 asynchronous (argue the gain and the concurrency / ordering contract first) | reuses the existing `runCriteriaConcurrent` scheduling contract; touches concurrent reads of shared caches |
| `GATES-BASE-CACHE` | cache the base-check verdicts (reuse the `gofmt` / `build` / `vet` / `test` results across scopes) | same family as the section 6.3 G4b rule "never cache a verdict" => **a directive clause must land first**; highest risk, placed last |

**Queue discipline**: layer boundaries and numbering discipline are the same as v60; this version raises no version.

---

## 5. Observation items

- Under observation (**13**): `OB-1 / 7 / 10 / 33 / 87 / 99 / 100 / 104 / 105 / 107 / 111 / 113 / 135` (carried from v61).
- Pending discussion (9): `OB-19 / 21 / 52 / 53 / 57 / 58 / 59 / 61 / 75`.
- To do (1): `OB-69`. Non-terminal total **23** (13 + 9 + 1), consistent with this version's re-check.
- **DR register**: DR max = **18** (P85 / P86 / P87 / P88 added none).
- **P87 Hydrogen precondition gap (record state)**: the light-tier clause (`Gemini 4 Hydrogen`) is landed while its comparison-review measurement is not done; it points at a later business slice (section 1.7.4 temporary authorization plus an OB record) and **must not be described as validated**; the note is backfilled into `docs/validation/directive-review-model-policy{,_EN}.md` once done. Registration slot = that report's section 7 and the P87 ledger section.
- **Four performance candidates registered in layer 3**: `GATES-CI-PERF` / `GATES-SELFTEST-SPLIT` / `GATES-CRITERIA-ASYNC` / `GATES-BASE-CACHE` are registered in ascending risk order in section 4 layer 3 (records requiring the user to name them; `GATES-CI-PERF` is evaluated together with "ci freeze cost" at project start).
- **Three tmp leftovers carried over**: `P84-CI-FIX` and `P85-SECURITY-PROOF-SUBSYSTEM` are carried by `TMP-ARCHIVE-P84-P85`; `find-quotes.js` has no referencing residue and still awaits the tmp-lifecycle slice.
- Mechanical readings: section 12.4 OB max = **144**; ADR max = **48**; DR max = **18**; pilot deferral max = the 25th.
- **Removed in this version**: the P88 sequencing registration (closed). **Added in this version**: the four performance candidates registered in layer 3, and the three tmp leftovers carried over.

---

## 6. Pending user rulings

- **This table is currently empty** (the P88 completion record and the v72 snapshot generation were authorized together on 2026-11-25).
- **Next ruling item (expected)**: **naming the `TMP-ARCHIVE-P84-P85` start** (tier-S tooling slice; must be named by the user; the start message records the four elements; the draft is already in the ledger).

---

## 7 / 8 / 9. Pointer sections

Same structure as v60: section 7 start and step-0 general rules → directive section 13.1; section 8 cross-slice collaboration discipline → section 13.2; section 9 checklist meta-rules and gate pitfalls → sections 13.3 / 13.4. **The rule bodies have moved into directive section 13 and are not copied here**.

---

## 10. One-paragraph opening for a new session

> Below is my project state snapshot (v72). Check the seven anchors in section 1 first; if they agree, continue with the section 4 queue, and if not, the repository governs.
>
> [paste this checklist v72 in full]
>
> [companion documents] `REMAINING_SLICES{,_EN}.md` (P88 is closed with a completion record; `TMP-ARCHIVE-P84-P85` is registered as the next item and awaits naming), `DELIVERY_DIRECTIVE.md` (**v1.73**), `gates-perf-assess{,_EN}.md` (the P88 report pair), `docs/validation/evidence/P88-GATES-PERF-ASSESS/` (the P88 evidence pack), `ADR_REGISTER{,_EN}.md` (ADR-048 stands), `DOCUMENTATION_PLAN{,_EN}.md` (the hand-over entry now points at v72), `CHECKLIST_v71{,_EN}.md` (frozen, superseded by v72).
>
> Current: P88 closed (directive v1.73; batched blob prefetch plus read memos, criterion logic untouched, five A/B pairs byte-identical; CI critical path 380.3 to 141.8 s); next = the `TMP-ARCHIVE-P84-P85` tooling slice (must be named by the user); the pre-S2 re-evaluation (S-G1 / S-G8) = the next T027 mainline station, pending user scheduling (section 4, later).

## 11. One-line status

P88 **closed** (two commits `03c45f7` / `9b44dbf` pushed to origin; main CI `45523666470` green on both legs and backfill `45524509376` observation only; R = 0, D = 0; directive v1.73; criterion logic untouched; CI critical path 380.3 to 141.8 s); next = the **`TMP-ARCHIVE-P84-P85` tooling slice** (must be named by the user); the pre-S2 re-evaluation (S-G1 / S-G8) = the next T027 mainline station, pending user scheduling (section 4, later).
