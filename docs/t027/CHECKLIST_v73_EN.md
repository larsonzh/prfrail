# Cross-session task checklist (snapshot v73, after P89 closes, next item = GATES-CI-PERF)

Generated: 2026-11-25. HEAD = `9a784dc`; directive = v1.73; P89's four CI runs: `45584959722` (main, both legs green), `45587252859` (backfill, `Gates` red), `45589128432` (fix, both legs green), `45590253878` (closing, both legs green).

**Supersession**: this version supersedes v72 (`docs/t027/CHECKLIST_v72{,_EN}.md` is frozen and is not retroactively edited).

**Structure**: same as v72 (section 1 anchors, section 2 fact pack, section 3 completed slices, section 4 execution queue, section 5 observation items, section 6 pending rulings, sections 7 / 8 / 9 pointer sections, section 10 opening prompt, section 11 one-line status); the rule bodies of sections 7 to 9 live in directive section 13 and are not copied here.

**Revision notes (against v72)**:

- (1) Section 1 anchors roll forward: HEAD `9a784dc`; directive **v1.73** (P89 did not bump it); SELFTEST **744/744**; the rolling window drops v70 and adds v72 / v73 (v71 / v72 / v73), with the v71 / v72 counts re-verified by this version and v73's own count left to v74.
- (2) Section 2 becomes the P89 fact pack (`CLOSED`, four commits and four CI runs, R = 0, **D = 1**, directive unchanged at v1.73).
- (3) Section 3 appends one P89 row (carrying v72's twenty-five rows forward = twenty-six rows in total).
- (4) Section 4 queue: P89 leaves; the next item = **`GATES-CI-PERF` (the first of the four tier-3 performance candidates; the user must name it; merge its assessment with "ci frozen cost" when it starts)**; "after that" stays the pre-S2 re-evaluation (S-G1 / S-G8); tier 3 keeps the other three candidates plus the five existing suspended items.
- (5) Section 5: removes "the three tmp leftovers are carried forward" (closed by P89); `find-quotes.js` keeps "awaiting the tmp-lifecycle slice"; P89 adds no new observation item.
- (6) Section 6 pending rulings: the next ruling item = **naming `GATES-CI-PERF` to start**.
- (7) Annex A: obligation transfer = **none**; annex B is carried over unchanged; sections 10 / 11 roll forward, and the `DOCUMENTATION_PLAN{,_EN}` pointer moves to v73.

---

## 1. Current anchors (check these seven numbers first)

**Rolling-window caliber = the most recent three versions (v71 / v72 / v73)**; v70 has been evicted (the same eviction rule as v52).

| Item | Value | Verification command |
|---|---|---|
| HEAD | `9a784dc` | `git rev-parse --short HEAD` |
| origin/main | `9a784dc` | `git rev-parse --short origin/main` |
| Directive version | **v1.73** | head of `docs/DELIVERY_DIRECTIVE.md` |
| SELFTEST | **744/744** | `node tools/gates/gate.js --selftest` |
| Three scopes | **all TOTAL_FAIL=0** | `gate.js --all --scope=index` / `tree` / `ci` |
| v71 whole-file count | **201** (LF count, CN / EN symmetric; **re-verified here**) | count LF in `docs/t027/CHECKLIST_v71{,_EN}.md` |
| v72 whole-file count | **214** (LF count, CN / EN symmetric; **re-verified here**) | count LF in `docs/t027/CHECKLIST_v72{,_EN}.md` |

> This snapshot's HEAD = `9a784dc` is the generation-time anchor after P89's closing commit; once v73 itself lands the HEAD moves forward, so `git rev-parse --short HEAD` is authoritative at that point. v73's own count is left to v74.

> Mechanical readings: ADR max = **48**; OB max = **144**; DR max = **18**; trial-deferral max = the 25th. P89 checked item by item: no new ADR, no new OB, no new DR.

> Rolling window = the most recent three versions (v71 / v72 / v73), v70 evicted; if two or more items disagree, the repository state wins - report the final readings first and continue.

`tmp` fixed assertion: top level = `.gitkeep` + `gates-cache` (holding only `g4b`) + `find-quotes.js` (no reference residue, awaiting the tmp-lifecycle slice). **`P84-CI-FIX` and `P85-SECURITY-PROOF-SUBSYSTEM` were archived into their evidence packs and removed by P89**, so this assertion no longer contains them.

---

## 2. P89 fact pack

- **Slice / size**: `TMP-ARCHIVE-P84-P85` (archive the two tmp directories and rewrite their references) / **S**; main controller `DeepSeek V4.1 Flash` High; the user named it on 2026-11-25.
- **Goal**: archive both former temp directories under the `tmp` root into their own evidence packs, rewrite every live reference, remove the original directories, and leave mechanically verifiable pack self-consistency evidence plus a zero-hit criterion.
- **Verdict**: **`CLOSED`** (2026-11-25); the directive stays at **v1.73** ((1) found no touch on CONTRACTS / schema / contract fixtures / workflows / the CI workflow, so no MINOR bump).
- **Landing**: **110 items** (44 under `P84-CI-FIX`, 66 under `P85-SECURITY-PROOF-SUBSYSTEM`) archived into `.../P84-BOUNDDIR-RESIDUAL/archive/P84-CI-FIX/` and `.../P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/`; byte normalization on 88 items (`.gitattributes` is `* text=auto eol=lf`, so hashes must be taken after normalization) and path-literal normalization on 11 items / 39 occurrences; all 110 items carry an "original / archived sha256" pair and replay from the untouched backup (110/110 pass); six live references rewritten (two ledger pairs, two P84 report-pair cells, one code comment), while the frozen surface keeps its original bytes per section 11.1.
- **Criteria and acceptance**: for all three packs `MANIFEST / SHA256SUMS / disk` are three-way equal (P84 58 / P85 95 / P89 14) with zero hash drift; the residual `git grep -F` hits for the two former paths are **16 = the 16-entry whitelist**, with 0 unexpected and 0 unused; `SELFTEST` 744/744; all three scopes `TOTAL_FAIL=0`; `go build/vet/test` green.
- **Review chain**: (1) x 1 (`DeepSeek V4 Pro`), (2)/(3) by the main controller directly, (5) 2/2 (r1 1 HIGH + 4 MEDIUM + 4 LOW, then r2 three description/count items), (6) 1/3 `PASS`, (7) 2/2 (r1 `FINDINGS`, then r2 `PASS WITH FIXES` whose only fix is closed), (8)a / (9) / (8)b / (8)c, then (10).
- **Metering**: **R = 0, D = 1, probes 0**. D = 1 is the fix round leaving one gate capture in the `tmp` root, deleted as soon as it was noticed (self-healing tier, recorded); (2)/(3) were carried by the main controller directly, the GLM batch assist was replaced by mechanical verification, and F-09 / F-13 were endorsed by (7) - accepted by the user on 2026-11-25 and **not to be taken as precedent**.
- **Four commits and CI**: main commit `5b85f7a` (135 paths) -> run `45584959722` **both legs green** (each leg `SELFTEST` 744/744, `TOTAL_FAIL=0`, `scope=ci changed=135`, 36 package-level ok records) -> backfill `a8eccc9` -> run `45587252859` with `Gates` red on both legs (`G8-b`: three `02-verification/` captures hashed with CRLF against LF blobs) -> fix `e15ad3d` -> run `45589128432` **both legs green** -> closing `9a784dc` -> run `45590253878` **both legs green**.
- **Documents**: report pair `docs/validation/tmp-archive-p84-p85{,_EN}.md` (with the `artifacts` fence and the `sixstage` (6) transcript); evidence pack `docs/validation/evidence/P89-TMP-ARCHIVE-P84-P85/` (14 payloads plus `MANIFEST.md` and `SHA256SUMS.txt`, including the (1) design, two (5) rounds, one (6) round and two (7) rounds); the ledger's P89 completion record.
- **tmp and leftovers**: P89's own temp directory was removed at (8)c; the `tmp` top level is back to three items; `find-quotes.js` still awaits the tmp-lifecycle slice.

### Annex A: obligation-transfer section (this version: **none**)

- **Historical transfer** (first instance traced to P75 O-7): the `s1-exit.md:52` in-line inconsistency of the core task chain was discharged by P76 change 1 plus the S-G2 correction (`0f80c1e`).
- **This version's transfer**: **none** - P89 leaves no obligation; `GATES-CI-PERF` and the other performance candidates remain **tier-3 record-only** and constitute neither review nor delivery obligations.
- **Convention note**: adopted since v59; a landing-slice entry must carry an "admission conditions" face.

### Annex B: AT-23 stage 2 status (unchanged)

- **Status**: ruling A1 says "not applicable" (user ruling 2026-11-09; not a permanent ban).
- **Executable objects before S2**: none (R-3 = (A)).
- **Restart path**: if the user later rules to restart stage 2, a real enforcement hookup (S2) or an alternative path must exist first.
- **Natural restart moment**: S2 starts -> enforcement really hooked up -> S-G8 unlocks -> R-3 re-evaluated.

---

## 3. Completed slices (recent baseline)

**v63 section 3 held P64-P80 as sixteen rows; v64 appended P81 (seventeen); v65 appended P82 (eighteen); v66 carried v65's eighteen and appended P79 (nineteen); v67 carried v66's nineteen and appended P83 (twenty); v68 carried v67's twenty and appended P84 (twenty-one); v69 carried v68's twenty-one and appended P86 (twenty-two); v70 carried v69's twenty-two (readings carried as-is, no retroactive edits) and appended P85 (twenty-three); v71 carried v70's twenty-three (readings as-is) and appended P87 (twenty-four); v72 carried v71's twenty-four (readings as-is) and appended P88 (twenty-five); v73 carries v72's twenty-five (readings as-is) and appends one P89 row. Twenty-six rows in total.**

| Version | Slice | Main commit | Verdict |
|---|---|---|---|
| v1.66 | P64 AT24-SUPERVISED-BLACK-BOX | `2dbc54a` | closed (AT-24 conditions closed; R = 2; run `42948727363` green on both legs first try) |
| v1.66 | P65 GATE-PIT-26-LANDING | `0eda074` | closed (appendix B.7.1 slot plus A-part ledger todo; R = 1; run `42969802819`) |
| v1.66 | P66 CHECKLIST-RULES-HOIST | `1650d5f` | closed (section 13 five subsections plus the new snapshot structure; R = 3, D = 2; run `43022796716`) |
| v1.67 | P67 GATES-DEBT-CLEANUP | `18b6ef9` + `365a930` | closed (g4b LRU, gate pits 27-30, section 12.4 clearing block, B.6 hardening; R = 3, D = 1; run `43124686702`) |
| v1.67 | P68 REMAINING-RULINGS | `e0d745c` + `665a31f` | closed (assessment class; R = 1, D = 1; (7) PASS; five R items ruled; ADR-042; run `43177025096`) |
| v1.67 | P69 T027B-E2E-WIRING | `759c934` | closed (implementation class; R = 0, D = 1; (7) PASS; E-1 recorded; run `43264324421`) |
| v1.68 | P70 S1-STREAM-PROGRESS | `ca88fa3` | closed (implementation class; R = 1, D = 2; (7) PASS; DR-14; runs `43400311736` -> `43403891205` -> `43408438643`) |
| v1.68 | P71 S1-SCOPE-ANNOTATE | `6ee2c61` + `16c98f7` | closed (documentation class; R = 0, D = 2; (7) blind PASS plus final PASS; runs `43463089586` + `43468806314`) |
| v1.68 | P72 REQ020-S2-RESTRICTION | `fb0ab57` + `406f4a0` | closed (documentation class; R = 0, D = 0; (6) skip confirmed by (7); (7) final PASS WITH FIXES; runs `43566904684` + `43568186355`) |
| v1.69 | P73 T027B-ADMISSION | `8c59712` + `f834ce5` | closed (assessment class; R = 0, D = 0; (6) skip confirmed by (7); (7) final PASS; **directive rises to v1.69**; runs `43622770750` + `43624661784`; OB-55 = option b) |
| v1.69 | P74 SANDBOX-INTERFACE | `3300364` + `150423d` | closed (implementation class; R = 3, D = 0; mandatory (6) PASS; (7) final PASS; **R=3 triggers the section 9.2 high-risk mark for the first time**; runs `43669748483` + `43672994996`) |
| v1.69 | P75 R3-REAL-CLI-ASSESSMENT | `ed2e26e` + `70a04cd` + `60146c0` | closed (assessment class; R = 1, D = 1; (7) one blind plus one final; **R-3 = (A)**; ADR-045; runs `43732384194` + `43734829063` + `43739409108`) |
| v1.69 | P76 S1-EXIT-CLOSEOUT-ASSESSMENT | `77f23d2` + `aca19af` + `0f80c1e` | closed (assessment class; R = 1, D = 1; (7) final PASS WITH FIXES; Q2=(a), Q10=(C); ADR-046 / ADR-047; runs `43785603244` + `43787083404` + `43793273520`) |
| v1.69 | P77 ADR010-APPROVAL | `33272be` + `fc00c4b` | closed (documentation class; R = 1, D = 0; ADR-010 = `APPROVED`; S-G4 passes; runs `43986942087` + `43997058549`) |
| v1.69 | P78 RELEASE-ZIP-BINDING | `c935713` + `0741fae` + `308b235` | closed (implementation class; R = 0, D = 0; (7) FINDINGS then re-review PASS; GitHub Release `v0.1.0` bound; S-G3 passes; runs `44043705691` + `44060224254`) |
| v1.69 | P80 SG6-GATELEVEL-CONSOLIDATE | `e00d79b` + `58c2794` | closed (documentation class; R = 2, D = 1; (7) narrow read `RE-REVIEW: PASS`; S-G6 holds on the gate level; runs `44211368890` + `44213070006`) |
| v1.70 | P81 DIRECTIVE-MODEL-UPGRADE | `45590e8` + `fd4d8fb` + `9ca6576` | closed (governance slice; (6) `MAI Code 2 Flash`, (7) `GPT-6.2 Sol`; SELFTEST 737 -> 740; R = 1, D = 2; runs `44281812336` + `44284200392` + `44392113028`) |
| v1.71 | P82 TMP-DISCIPLINE-OUT-OF-REPO | `0264bcd` + `adb8cb9` + `1015430` + `9d83089` + `9fd8cd6` | closed (governance slice; section 5.1 gains the out-of-repo self-healing clause; OB-139 flipped to `handled`; **directive rises to v1.71**; R = 2, D = 3; all five CI runs green) |
| v1.71 | P79 S2-PRE-INDEPENDENT-SECURITY-REVIEW | `a464a44` + `1b7f9c6` | closed (assessment class; all three faces `FINDINGS`, F1-F8 handed to the fix slices; ADR-047 flipped to `COMPENSATED`; R = 2, D = 1; runs `44585713851` + `44589887023`) |
| v1.71 | P83 SECURITY-FINDINGS-REMEDIATION | `ec532a4` + `c820743` + `4323503` | **partially closed** (F2 / F5 / F7 plus the narrowed T9; F1 / F4 partly implemented; F3 / F6 / F8 deferred to P84 / P85; R = 2, D = 0; DR-15; first run red `44795140136` -> `44801690679` green -> observation `44805856028`) |
| v1.71 | P84 BOUNDDIR-RESIDUAL | `d597f1e` + `32bf234` + `509b01b` | **partially closed** (T5#4 fully closed, T4 mostly closed; the residual N13-UNC plus the evidence-attribution Low **were closed by P86**; R = 3, D = 0; DR-16 / DR-17; first run red `44877224368` -> fix `44926215282` green -> observation `44928363670`) |
| v1.71 | P86 BOUNDDIR-UNC-EQUIVALENCE | `f1ae821` + `a5ef122` + `1eb63d9` + `565fc9f` | **closed** (N13-UNC and W10 closed, F1 fully closed; R = 0, D = 0; DR-18; first run `44962802936` green -> guard round red `44965434402` -> fix round `44967047189` green -> observation `44969239294`) |
| v1.71 | P85 SECURITY-PROOF-SUBSYSTEM | `cc7d90f` + `6480bb1` + `18b2613` + `ef1e1f2` | **closed** (T2 / T3 / T8a closed, F3 / F6 / F8 settled, **F1-F8 all closed**; R = 8, D = 0; OB-140/141/142/143/144; first run `45201921016` had two environmental reds on the windows leg -> R=8 fix `6480bb1` -> `45208001981` both legs green) |
| v1.72 | P87 DIRECTIVE-REVIEW-MODEL-POLICY | `8b64ce3` + `2782a09` | **closed** (directive rises to **v1.72**; (7) four-tier selection matrix plus the grading-threshold consensus plus the new section 5.4; G5-b closed set of 11; ADR-048; R = 0, D = 0; run `45294006169` both legs green plus observation `45356014951`; Hydrogen precondition gap registered for a later slice) |
| v1.73 | P88 GATES-PERF-ASSESS | `03c45f7` + `9b44dbf` | **closed** (tooling slice; criterion bodies and decision logic untouched, five same-input A/B pairs byte-equal; SELFTEST 744/744; three scopes 0; CI critical path 380.3 -> 141.8 s; R = 0, D = 0; run `45523666470` green first try plus observation `45524509376`) |
| v1.73 | P89 TMP-ARCHIVE-P84-P85 | `5b85f7a` + `a8eccc9` + `e15ad3d` + `9a784dc` | **closed** (tooling / evidence slice; 110 items archived whole, three packs three-way equal at P84 58 / P85 95 / P89 14, residual hits 16 = the 16-entry whitelist; R = 0, **D = 1**; main run `45584959722` both legs green -> backfill run `45587252859` `Gates` red (`G8-b` CRLF hashing) -> fix run `45589128432` both legs green -> closing run `45590253878` both legs green) |

> **D caliber of section 3** (same as v61): "D = n" here abbreviates the temp-artifact discipline violation count (section 5.1); P89 totals **D = 1** (recorded in the report pair section 6 and the (10) receipt). **The D count and the section 9.2 R count are two independent series**.

---

## 4. Execution queue

### Tier 0 pre-work (✅ discharged; nothing new here)

- P5 / P6 were closed by P57; tier 0 is empty with no residue.

### Tier 1 main-line hard gates

- Items 8-12: all ✅ discharged (see v57 / v58 section 4; trace back through that reference chain if needed).
- **`T027-B` (defensive-refactor phase) = done (partly not applicable)**: C1-C7 / N1-N4 / R-1 / R-2 / R-3 are all handled (R-3 = (A)). **The main line is fully discharged**.
- **`S1 Exit closeout` (gate level done; security remediation F1-F8 fully settled)**: P76 produced the closeout roadmap; P77 closed S-G4; P78 closed S-G3; P80 judged S-G6 to hold at the gate level; the independent security review was completed by P79 and P83 / P84 / P86 / P85 delivered the F1-F8 fixes. Gate table below (P87 / P88 / P89 left it unchanged):

| Gate | Criterion | Status | Basis |
|---|---|---|---|
| S-G1 | the CLI can really start an admitted AI candidate and bill it | **not passing** (Q2=(a) literal criterion) | P76 section 4; ADR-046 |
| S-G2 | `supervised-black-box` implemented and AT-24 holds | **passing** (T028 `COMPLETED` plus AT-24 conditions closed) | P76 section 4; `0f80c1e`; ADR-046 |
| S-G3 | the final ZIP is bound to the release record and re-verified | **passing** (closed by P78) | P78 report sections 5 / 7 |
| S-G4 | ADR-010 is approved by the product owner | **passing** (user approval 2026-11-16; closed by P77) | P77 report section 5 |
| S-G5 | the product owner grants approval | **release part approved (step one)** | ADR-046 Q10=(C); P76 section 12 |
| S-G6 | the dependency fields hold and the four remaining gates close | **holds at the gate level** (P80 combined verdict) | P80 report sections 5 / 6 |
| S-G7 | every AT passes, T018 is done and section 14 fully exits | unreachable (composite gate) | `s1-exit.md:52` / `:56`-`:58` / `:60` |
| S-G8 | the supervised path runs end to end against a real CLI and a real candidate | blocked until the S2 enforcement hookup | P75 A2 / ADR-045; P76 section 6 |

**Next item = `GATES-CI-PERF`**

- First of the four tier-3 performance candidates (lowest risk first); **the user must name it**; the start message records the four elements; its ledger entry is not yet registered (to be added when it starts).
- Scope: CI workflow side - assessing the Gates step as `--all --skip-base --scope=ci`, plus persisting the `go build` / `go test` caches (the latter is an extension beyond P88 report section 7 and **touches no criterion**).
- Boundaries: **touches the CI workflow (user authorization required)**; it shares its origin with the existing "ci frozen cost" item, so **merge the assessment when it starts** to avoid duplication; no criterion-semantics change.
- After that = the pre-S2 re-evaluation (S-G1 / S-G8) (the next T027 main-line stop, pending user scheduling):
  - Prerequisite = the S2 enforcement hookup; the independent security review prerequisite was discharged by P79 and P83 / P84 / P86 / P85 delivered the F1-F8 fixes, so the ADR-047 rollback clause no longer blocks;
  - `S-G8` stays blocked until the S2 enforcement hookup (P73 A2 / P75 A2 / ADR-044 / ADR-045); the `s1-exit{,_EN}.md` gate table is untouched.
- **Post-S2 re-evaluation items**: **S-G5 step two (overall S1 Exit approval) plus S-G7**, in the order of P80 report section 7.

### Tier 2 routine governance

- P72 / P73 / P74 / P75 / P76 / P77 / P78 / P80 / P81 / P82 / P79 / P83 / P84 / P86 / P85 / P87 / P88 / **P89** all ✅ discharged (tier 2 is **empty**).
- Each slice's full section 1.2 field set lives in `docs/t027/REMAINING_SLICES{,_EN}.md`.

### Tier 3 low-priority candidates

| Field | Content |
|---|---|
| Identifiers | `OB-50` / `spawnGit` / `ci frozen cost` / `%TEMP% historical leftovers, 22 items` / `GATES-EXEMPT` plus the four performance candidates (`GATES-CI-PERF` / `GATES-SELFTEST-SPLIT` / `GATES-CRITERIA-ASYNC` / `GATES-BASE-CACHE`) |
| Current state (verified item by item) | OB-50 = **handled** (section 12.4, v1.41); spawnGit = the git-call helper of `tools/gates/selftest.js` (spawn wrapper in `lib/ctx.js`, stderr captured into the result object and not printed), a diagnostic OB candidate from P27 report section 5 item 5 with no section 12.4 entry; ci frozen cost = no section 12.4 entry, from P27 report section 7 "candidate OBs (performance)", now assessed under `GATES-CI-PERF`; `%TEMP%` leftovers = OB-98 handled, OB-99 / OB-100 under observation; GATES-EXEMPT = **not started** (no section 12.4 entry; split out in P33 (1)-3(iv), touches four criterion files); `GATES-CI-PERF` has been promoted to "next item" while the other three remain record-only |
| Contents | See the matching section 12.4 rows; `GATES-PERF-ASSESS` closed as P88 (v1.73) and is no longer listed here; the three remaining performance candidates are in the table below |
| Dependencies | The long suspension covers only the unhandled items (spawnGit / the `%TEMP%` temporal shape OB-99, OB-100 / GATES-EXEMPT), none with a defined trigger |
| Size | Varies (S through M-L) |
| (7) quota | none |
| Boundaries | This version advances no suspended item and keeps the registrations; the other three performance candidates are record-only and constitute no slice |
| Start condition | the user must name it explicitly |

**Performance candidates (risk ascending; `GATES-CI-PERF` has left the tier as the "next item", the table keeps the other three)**:

| Identifier | One-line scope | Relation to existing items / preconditions |
|---|---|---|
| `GATES-SELFTEST-SPLIT` | per-suite timing plus memoized fixture reads for `--selftest`, a G4b inner-loop increment (benefit must be argued first), and a built-in `--timings` | touches `selftest.js` / `g4b.js` (criterion implementations), so (6)/(7) are mandatory; `--timings` is an output-contract surface |
| `GATES-CRITERIA-ASYNC` | make G6 / G7 / G8 asynchronous (argue benefit and the concurrency / ordering contract first) | reuses the existing `runCriteriaConcurrent` scheduling contract; involves concurrent reads of shared caches |
| `GATES-BASE-CACHE` | cache the base verdicts (reuse the `gofmt` / `build` / `vet` / `test` conclusions across scopes) | same family as section 6.3 G4b "never cache verdicts", so **a directive clause must exist first**; highest risk, listed last |

**Queue discipline**: tier boundaries and numbering follow v60; this version does not bump the directive.

---

## 5. Observation items

- Observing (**13**): `OB-1 / 7 / 10 / 33 / 87 / 99 / 100 / 104 / 105 / 107 / 111 / 113 / 135` (carried from v61).
- In discussion (9): `OB-19 / 21 / 52 / 53 / 57 / 58 / 59 / 61 / 75`.
- To do (1): `OB-69`. Non-terminal total **23** (13 + 9 + 1), matching this version's verification.
- **DR registrations**: DR max = **18** (P85 / P86 / P87 / P88 / P89 added none).
- **P87 Hydrogen precondition gap (record-only)**: the light-tier clause (`Gemini 4 Hydrogen`) has landed while its comparison-review measurement is not done; it points at a later business slice (section 1.7.4 temporary authorization plus an OB record) and **must not be described as validated**; backfill a note into `docs/validation/directive-review-model-policy{,_EN}.md` once complete. Registration sites: report section 7 and the ledger's P87 section.
- **Tier-3 performance registrations (rolled forward here)**: `GATES-SELFTEST-SPLIT` / `GATES-CRITERIA-ASYNC` / `GATES-BASE-CACHE` are record-only (the user must name them); `GATES-CI-PERF` has been promoted to the section 4 "next item" (merge its assessment with "ci frozen cost" when it starts).
- **`find-quotes.js` disposition site**: no reference residue (P89 re-checked and deliberately excluded it from the archive and the rewrites); still awaiting the tmp-lifecycle slice.
- Mechanical readings: section 12.4 OB max = **144**; ADR max = **48**; DR max = **18**; trial-deferral max = the 25th.
- **Removed here**: "the three tmp leftovers are carried forward" (closed by P89). **Added here**: `GATES-CI-PERF` promoted out of tier 3, and the `find-quotes.js` disposition site.

---

## 6. Pending user rulings

- **This table is currently empty** (P89's completion record, the three ledger corrections and this v73 snapshot were all authorized together on 2026-11-25).
- **Next ruling item (expected)**: **naming `GATES-CI-PERF` to start** (first tier-3 performance candidate; touches the CI workflow, so user authorization is required; merge its assessment with "ci frozen cost").

---

## 7 / 8 / 9 Pointer sections

Same structure as v60: section 7 start-up and step 0 general rules -> directive section 13.1; section 8 cross-slice collaboration discipline -> section 13.2; section 9 checklist meta-rules and gate pits -> section 13.3 / 13.4. **The rule bodies live in directive section 13 and are not copied here**.

---

## 10. One-line opening prompt for a new session

> Here is my project state snapshot (v73); check the seven anchors in section 1 first, continue with the section 4 queue if they match, and let the repository state win if they do not.
>
> [paste this whole v73 checklist]
>
> [Supporting documents] `REMAINING_SLICES{,_EN}.md` (P89 is closed with a completion record; `GATES-CI-PERF` is the next item, awaiting the start-up naming), `DELIVERY_DIRECTIVE.md` (**v1.73**), `tmp-archive-p84-p85{,_EN}.md` (P89 report pair), `docs/validation/evidence/P89-TMP-ARCHIVE-P84-P85/` (P89 evidence pack), `ADR_REGISTER{,_EN}.md` (ADR-048 stands), `DOCUMENTATION_PLAN{,_EN}.md` (the handover pointer now targets v73), `CHECKLIST_v72{,_EN}.md` (frozen, superseded by v73).
>
> Current state: P89 closed (110 items archived whole, three packs three-way equal, residual hits 16 = the 16-entry whitelist; R = 0, D = 1; directive stays v1.73); next = `GATES-CI-PERF` (first tier-3 performance candidate, touches the CI workflow, the user must name it); the pre-S2 re-evaluation (S-G1 / S-G8) is the next T027 main-line stop, pending user scheduling (section 4 "after that").

## 11. Current status in one line

P89 **closed** (four commits `5b85f7a` / `a8eccc9` / `e15ad3d` / `9a784dc` pushed to origin; all four CI runs end green on both legs, the only red being the backfill run `45587252859` from `G8-b` CRLF hashing, closed by the fix round; R = 0, D = 1; directive unchanged at v1.73; the `tmp` top level is back to three items); next = **`GATES-CI-PERF`** (touches the CI workflow, the user must name it); the pre-S2 re-evaluation (S-G1 / S-G8) is the next T027 main-line stop, pending user scheduling (section 4 "after that").
