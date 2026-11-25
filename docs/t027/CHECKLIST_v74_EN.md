# Cross-session task checklist (snapshot v74, after P90 closes, next item = DIRECTIVE-CI-INVOCATION-SYNC)

Generated: 2026-11-26. HEAD = `9ac8f6b`; directive = v1.73; P90's two CI runs: `45659843522` (main, both legs green on attempt 1) and `45660930671` (backfill, both legs green, observation record only).

**Supersession**: this version supersedes v73 (`docs/t027/CHECKLIST_v73{,_EN}.md` is frozen and is not retroactively edited).

**Structure**: same as v73 (section 1 anchors, section 2 fact pack, section 3 completed slices, section 4 execution queue, section 5 observation items, section 6 pending rulings, sections 7 / 8 / 9 pointer sections, section 10 opening prompt, section 11 one-line status); the rule bodies of sections 7 to 9 live in directive section 13 and are not copied here.

**Revision notes (against v73)**:

- (1) Section 1 anchors roll forward: HEAD `9ac8f6b`; directive **v1.73** (P90 did not bump it); SELFTEST **744/744**; the rolling window drops v71 and adds v72 / v73 / v74, with the v72 / v73 counts re-verified by this version and v74's own count left to v75.
- (2) Section 2 becomes the P90 fact pack (`CLOSED`, two commits and two CI runs, **R = 1**, D = 0, directive unchanged at v1.73, boundary extension ratified by the user and not to be taken as precedent).
- (3) Section 3 appends one P90 row (carrying v73's twenty-six rows forward = twenty-seven rows in total).
- (4) Section 4 queue: P90 leaves; the next item = **`DIRECTIVE-CI-INVOCATION-SYNC` (a small tier-S governance slice; the user must name it; the proposed wording is in P90 report section 8 K4)**; "after that" stays the pre-S2 re-evaluation (S-G1 / S-G8); tier 3 drops `GATES-CI-PERF` and records "ci frozen cost" as handled by P90, keeping the other three candidates and the existing suspended items.
- (5) Section 5: the performance-candidate wording goes from four to three; P90 adds no new observation item; `OB-145` is registered as a candidate (pointing at P91).
- (6) Section 6 pending rulings: the next ruling item = **naming `DIRECTIVE-CI-INVOCATION-SYNC` to start**.
- (7) Annex A: obligation transfer = **none**; annex B is carried over unchanged; sections 10 / 11 roll forward, and the `DOCUMENTATION_PLAN{,_EN}` pointer moves to v74.
- (8) Mechanical readings re-checked: ADR 48 / OB 144 / DR 18 / trial deferral 25 all unchanged; `OB-145` is a candidate, not registered in section 12.4, and does not count toward the OB max.
---

## 1. Current anchors (check these seven numbers first)

**Rolling-window caliber = the most recent three versions (v72 / v73 / v74)**; v71 has been evicted (the same eviction rule as v52).

| Item | Value | Verification command |
|---|---|---|
| HEAD | `9ac8f6b` | `git rev-parse --short HEAD` |
| origin/main | `9ac8f6b` | `git rev-parse --short origin/main` |
| Directive version | **v1.73** | head of `docs/DELIVERY_DIRECTIVE.md` |
| SELFTEST | **744/744** | `node tools/gates/gate.js --selftest` |
| Three scopes | **all TOTAL_FAIL=0** | `gate.js --all --scope=index` / `tree` / `ci` |
| v72 whole-file count | **214** (LF count, CN / EN symmetric; **re-verified here**) | count LF in `docs/t027/CHECKLIST_v72{,_EN}.md` |
| v73 whole-file count | **212** (LF count, CN / EN symmetric; **re-verified here**) | count LF in `docs/t027/CHECKLIST_v73{,_EN}.md` |

> This snapshot's HEAD = `9ac8f6b` is the generation-time anchor after P90's backfill commit; once v74 itself lands the HEAD moves forward, so `git rev-parse --short HEAD` is authoritative at that point. v74's own count is left to v75.

> Mechanical readings: ADR max = **48**; OB max = **144**; DR max = **18**; trial-deferral max = the 25th. P90 checked item by item: no new ADR, no new OB, no new DR; `OB-145` is a candidate and is not registered in section 12.4.

> Rolling window = the most recent three versions (v72 / v73 / v74), v71 evicted; if two or more items disagree, the repository state wins - report the final readings first and continue.

`tmp` fixed assertion: top level = `.gitkeep` + `gates-cache` (holding only `g4b`) + `find-quotes.js` (no reference residue, awaiting the tmp-lifecycle slice). **`P84-CI-FIX` and `P85-SECURITY-PROOF-SUBSYSTEM` were archived into their evidence packs and removed by P89**, so this assertion no longer contains them.

---

## 2. P90 fact pack

- **Slice / size**: `GATES-CI-PERF` (deduplicate the base suite in the CI gates step and assess the Go cache) / **S**; main controller `DeepSeek V4.1 Flash` High; the user named it on 2026-11-25 and ratified the boundary extension plus the commit on 2026-11-26.
- **Goal**: change the CI `Gates` step to `--all --skip-base --scope=ci` (keeping `--selftest`) so the base quartet is not re-run on top of `Build` / `Vet` / `Test`; restore the missing `gofmt` defence; assess the Go build cache.
- **Verdict**: **`CLOSED`** (2026-11-26); the directive stays at **v1.73** ((1) found no touch on CONTRACTS / schema / contract fixtures / section 7.2.1 (3) semantics).
- **Landing**: the `Gates` step runs `--skip-base`; a new `Format` step (`shell: bash`, `set -e` plus a non-empty-output red) restores `gofmt`; two comment lines above `cache: true` (the Go cache is judged "already covered"); the frozen contract was synced in three places (the `test` step-name table, the `test/Gates` digest, the `test/Format` digest with `shell`) plus the necessary `different-shell` change from `bash` to `sh`; `tools/gates/**` untouched.
- **Criteria and acceptance**: in (3), stripping the four base records leaves both sides hashing to `ee00e354…`, and flipping one byte makes the comparator report DIFFERENT; in (9) all three scopes report `TOTAL_FAIL=0` (tree includes all four base records and `go test` is ok=15 FAIL=0); `SELFTEST` 744/744; the evidence pack holds 23 payloads with three-way equality.
- **Review chain**: (1) x 1 (`DeepSeek V4 Pro`), (2)/(3) by the main controller directly, (5) 2/2 (r1 1 Medium + 2 Low, then r2 all CLOSED), (6) 2/3 (r1 `FINDINGS` with 2 Medium, then r2 `PASS`), (7) 2/2 (r1 `FINDINGS` with 1 Medium + 1 Low, then r2 `RE-REVIEW: PASS`), (8)a / (9) / (8)b / (8)c, then (10).
- **Metering**: **R = 1, D = 0, probes 0**. R = 1 is `create_file` writing `.md` without BOM and with CRLF, since normalized (established tool behaviour, not a content defect); no GLM batch assist (mechanical verification, user ruling).
- **Two commits and CI**: main commit `ec221f5` (31 paths, +3132/-4) -> run `45659843522` **both legs green on attempt 1** (each leg `SELFTEST` 744/744, `scope=ci changed=31`, `TOTAL_FAIL=0`; the `Format` step succeeded on both legs, so `shell: bash` works on the windows leg) -> backfill `9ac8f6b` (7 files) -> run `45660930671` both legs green (`changed=7`, observation only per section 11.2).
- **Timing**: the `Gates` step went from 103 to 45 s on windows and from 44 to 6 s on ubuntu (an order-of-magnitude comparison within the same runner family; the controlled equivalence evidence is the (3) A/B sha256, since the two runs cover different change-set sizes).
- **Documents**: report pair `docs/validation/gates-ci-perf{,_EN}.md` (154 lines, with the `artifacts` fence, the `sixstage` (6) transcript and section 14 CI fact chain); evidence pack `docs/validation/evidence/P90-GATES-CI-PERF/` (23 payloads plus `MANIFEST.md` and `SHA256SUMS.txt`); the ledger's P90 completion record.
- **tmp and leftovers**: P90's own temp directory was removed at (8)c; the `tmp` top level is back to three items; `find-quotes.js` still awaits the tmp-lifecycle slice; candidate **OB-145** (the v1.40 sentence at directive `:675` still says `--all --scope=ci`) is handed to `DIRECTIVE-CI-INVOCATION-SYNC`, and leaving the directive untouched was the right call for this slice.

### Annex A: obligation-transfer section (this version: **none**)

- **Historical transfer** (first instance traced to P75 O-7): the `s1-exit.md:52` in-line inconsistency of the core task chain was discharged by P76 change 1 plus the S-G2 correction (`0f80c1e`).
- **This version's transfer**: **none** - P90 leaves no obligation; `DIRECTIVE-CI-INVOCATION-SYNC` and the other three performance candidates are record-only and constitute neither review nor delivery obligations.
- **Convention note**: adopted since v59; a landing-slice entry must carry an "admission conditions" face.

### Annex B: AT-23 stage 2 status (unchanged)

- **Status**: ruling A1 says "not applicable" (user ruling 2026-11-09; not a permanent ban).
- **Executable objects before S2**: none (R-3 = (A)).
- **Restart path**: if the user later rules to restart stage 2, a real enforcement hookup (S2) or an alternative path must exist first.
- **Natural restart moment**: S2 starts -> enforcement really hooked up -> S-G8 unlocks -> R-3 re-evaluated.

---

## 3. Completed slices (recent baseline)

**v63 section 3 held P64-P80 as sixteen rows; v64 appended P81 (seventeen); v65 appended P82 (eighteen); v66 carried v65's eighteen and appended P79 (nineteen); v67 carried v66's nineteen and appended P83 (twenty); v68 carried v67's twenty and appended P84 (twenty-one); v69 carried v68's twenty-one and appended P86 (twenty-two); v70 carried v69's twenty-two (readings carried as-is, no retroactive edits) and appended P85 (twenty-three); v71 carried v70's twenty-three (readings as-is) and appended P87 (twenty-four); v72 carried v71's twenty-four (readings as-is) and appended P88 (twenty-five); v73 carried v72's twenty-five (readings as-is) and appended one P89 row (twenty-six); v74 carries v73's twenty-six (readings as-is) and appends one P90 row. Twenty-seven rows in total.**

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

| v1.73 | P90 GATES-CI-PERF | `ec221f5` + `9ac8f6b` | **closed** (tooling slice; `Gates` runs `--skip-base` plus a new `Format` step, criterion bodies untouched, (3) A/B byte-equal after stripping base; SELFTEST 744/744; three scopes 0; `Gates` 103->45 s on windows and 44->6 s on ubuntu; R = 1, D = 0; main run `45659843522` green on both legs at attempt 1 plus observation `45660930671`) |

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

**Next item = `DIRECTIVE-CI-INVOCATION-SYNC`**

- A small tier-S governance slice (sync the CI invocation described in directive section 6.3); **the user must name it**; the start message records the four elements; its ledger entry is not yet registered (to be added when it starts).
- Scope: change "`--selftest` then `--all --scope=ci`" in `docs/DELIVERY_DIRECTIVE{,_EN}.md` to `--all --skip-base --scope=ci`, and add one sentence stating that base is covered by `Build` / `Vet` / `Test` / `Format`; the proposed wording is in P90 report section 8 K4.
- Boundaries: **touches the directive body, so a version package is required** (header version / section 0.3 row / appendix C); no criterion change and no `ci.yml` change; no gitee.

- After that = the pre-S2 re-evaluation (S-G1 / S-G8) (the next T027 main-line stop, pending user scheduling):
  - Prerequisite = the S2 enforcement hookup; the independent security review prerequisite was discharged by P79 and P83 / P84 / P86 / P85 delivered the F1-F8 fixes, so the ADR-047 rollback clause no longer blocks;
  - `S-G8` stays blocked until the S2 enforcement hookup (P73 A2 / P75 A2 / ADR-044 / ADR-045); the `s1-exit{,_EN}.md` gate table is untouched.
- **Post-S2 re-evaluation items**: **S-G5 step two (overall S1 Exit approval) plus S-G7**, in the order of P80 report section 7.

### Tier 2 routine governance

- P72 / P73 / P74 / P75 / P76 / P77 / P78 / P80 / P81 / P82 / P79 / P83 / P84 / P86 / P85 / P87 / P88 / P89 / **P90** all ✅ discharged (tier 2 is **empty**).
- Each slice's full section 1.2 field set lives in `docs/t027/REMAINING_SLICES{,_EN}.md`.

### Tier 3 low-priority candidates

| Field | Content |
|---|---|
| Identifiers | `OB-50` / `spawnGit` / `ci frozen cost` / `%TEMP% historical leftovers, 22 items` / `GATES-EXEMPT` plus the three performance candidates (`GATES-SELFTEST-SPLIT` / `GATES-CRITERIA-ASYNC` / `GATES-BASE-CACHE`) |
| Current state (verified item by item) | OB-50 = **handled** (section 12.4, v1.41); spawnGit = the git-call helper of `tools/gates/selftest.js` (spawn wrapper in `lib/ctx.js`, stderr captured into the result object and not printed), a diagnostic OB candidate from P27 report section 5 item 5 with no section 12.4 entry; ci frozen cost = no section 12.4 entry, from P27 report section 7 "candidate OBs (performance)", now assessed under `GATES-CI-PERF`; `%TEMP%` leftovers = OB-98 handled, OB-99 / OB-100 under observation; GATES-EXEMPT = **not started** (no section 12.4 entry; split out in P33 (1)-3(iv), touches four criterion files); "ci frozen cost" was handled by P90 (the `--skip-base` Gates step landed and the Go cache was judged already covered), while the other three remain record-only |
| Contents | See the matching section 12.4 rows; `GATES-PERF-ASSESS` closed as P88 (v1.73) and is no longer listed here; the three remaining performance candidates are in the table below |
| Dependencies | The long suspension covers only the unhandled items (spawnGit / the `%TEMP%` temporal shape OB-99, OB-100 / GATES-EXEMPT), none with a defined trigger |
| Size | Varies (S through M-L) |
| (7) quota | none |
| Boundaries | This version advances no suspended item and keeps the registrations; the other three performance candidates are record-only and constitute no slice |
| Start condition | the user must name it explicitly |

**Performance candidates (risk ascending; `GATES-CI-PERF` left the tier when P90 closed it, the table keeps the other three)**:

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
- **DR registrations**: DR max = **18** (P85 / P86 / P87 / P88 / P89 / P90 added none).
- **P87 Hydrogen precondition gap (record-only)**: the light-tier clause (`Gemini 4 Hydrogen`) has landed while its comparison-review measurement is not done; it points at a later business slice (section 1.7.4 temporary authorization plus an OB record) and **must not be described as validated**; backfill a note into `docs/validation/directive-review-model-policy{,_EN}.md` once complete. Registration sites: report section 7 and the ledger's P87 section.
- **Tier-3 performance registrations (rolled forward here)**: `GATES-SELFTEST-SPLIT` / `GATES-CRITERIA-ASYNC` / `GATES-BASE-CACHE` are record-only (the user must name them); `GATES-CI-PERF` left the tier when P90 closed it, and its "ci frozen cost" was handled at the same time.
- **`find-quotes.js` disposition site**: no reference residue (P89 re-checked and deliberately excluded it from the archive and the rewrites); still awaiting the tmp-lifecycle slice.
- Mechanical readings: section 12.4 OB max = **144**; ADR max = **48**; DR max = **18**; trial-deferral max = the 25th.
- **Removed here**: the "next item" wording for `GATES-CI-PERF` (now closed). **Added here**: the `OB-145` candidate slot (the stale sentence at directive `:675`, handed to `DIRECTIVE-CI-INVOCATION-SYNC`; a candidate only, so it does not count toward the section 12.4 OB max).

---

## 6. Pending user rulings

- **This table is currently empty** (the P90 boundary-extension ratification, the F-10 disclosure, the OB-145 candidate registration, the commit and backfill authorization and this v74 snapshot were all ruled on 2026-11-26).
- **Next ruling item (expected)**: **naming `DIRECTIVE-CI-INVOCATION-SYNC` to start** (a small tier-S governance slice; touches the directive body, so a version package is required; the proposed wording is in P90 report section 8 K4).

---

## 7 / 8 / 9 Pointer sections

Same structure as v60: section 7 start-up and step 0 general rules -> directive section 13.1; section 8 cross-slice collaboration discipline -> section 13.2; section 9 checklist meta-rules and gate pits -> section 13.3 / 13.4. **The rule bodies live in directive section 13 and are not copied here**.

---

## 10. One-line opening prompt for a new session

> Here is my project state snapshot (v74); check the seven anchors in section 1 first, continue with the section 4 queue if they match, and let the repository state win if they do not.
>
> [paste this whole v74 checklist]
>
> [Supporting documents] `REMAINING_SLICES{,_EN}.md` (P90 is closed with a completion record; `DIRECTIVE-CI-INVOCATION-SYNC` is the next item, awaiting the start-up naming), `DELIVERY_DIRECTIVE.md` (**v1.73**), `gates-ci-perf{,_EN}.md` (P90 report pair), `docs/validation/evidence/P90-GATES-CI-PERF/` (P90 evidence pack), `ADR_REGISTER{,_EN}.md` (ADR-048 stands), `DOCUMENTATION_PLAN{,_EN}.md` (the handover pointer now targets v74), `CHECKLIST_v73{,_EN}.md` (frozen, superseded by v74).
>
> Current state: P90 closed (CI `Gates` deduplicated plus a `Format` step plus the Go cache judged already covered; both legs SELFTEST 744/744, TOTAL_FAIL=0; R = 1, D = 0; directive stays v1.73); next = `DIRECTIVE-CI-INVOCATION-SYNC` (small tier-S governance slice, the user must name it); the pre-S2 re-evaluation (S-G1 / S-G8) is the next T027 main-line stop, pending user scheduling (section 4 "after that").

## 11. Current status in one line

P90 **closed** (main commit `ec221f5` plus backfill `9ac8f6b` pushed to origin; main run `45659843522` green on both legs at attempt 1 and backfill run `45660930671` green on both legs, observation only; both legs SELFTEST 744/744, `scope=ci changed=31`, TOTAL_FAIL=0; `Gates` 103->45 s on windows and 44->6 s on ubuntu; R = 1, D = 0; directive unchanged at v1.73; the `tmp` top level is back to three items); next = **`DIRECTIVE-CI-INVOCATION-SYNC`** (touches the directive body, the user must name it); the pre-S2 re-evaluation (S-G1 / S-G8) is the next T027 main-line stop, pending user scheduling (section 4 "after that").
