# Cross-session task checklist (snapshot v76, after P92 closes, next item = the G6-4 mutual-exclusion governance slice)

Generated: 2026-11-26. HEAD = `10a1f95`; directive = v1.75; P92 CI run: `45726904763` (main, both legs green on attempt 1); the backfill run is an observation record only per section 11.2.

**Supersession**: this version supersedes v75 (`docs/t027/CHECKLIST_v75{,_EN}.md` is frozen and is not retroactively edited).

**Structure**: same as v75 (section 1 anchors, section 2 fact pack, section 3 completed slices, section 4 execution queue, section 5 observation items, section 6 pending rulings, sections 7 / 8 / 9 pointer sections, section 10 opening prompt, section 11 one-line status); the rule bodies of sections 7 to 9 live in directive section 13 and are not copied here.

**Revision notes (against v75)**:

- (1) Section 1 anchors roll forward: HEAD `10a1f95`; directive **v1.75** (P92 bumped it as a MINOR); SELFTEST **744/744**; the rolling window drops v73 and adds v74 / v75 / v76, with the v74 / v75 counts re-verified by this version (216 / 215) and v76 staying for v77.
- (2) Section 2 becomes the P92 fact pack (`CLOSED`, two commits and the CI run, **R = 3**, D = 1, the directive bumped to v1.75, the benefit-gate conclusion, (7) 1+2 plus one narrow confirmation read, GLM x2).
- (3) Section 3 appends one P92 row (carrying v75 twenty-eight rows forward = twenty-nine rows in total).
- (4) Section 4 queue: P92 leaves; the next item = **the `G6-4 mutual-exclusion governance slice` (tier S; the user must name it; it touches the mutual-exclusion clauses of directive sections 6.3 / 11.3 together with the T23 assertion and the g6 criteria)**; "after that" stays the pre-S2 re-evaluation (S-G1 / S-G8); tier 3 drops `GATES-SELFTEST-SPLIT` (closed), keeping the other two performance candidates and the suspended items.
- (5) Section 5: the "`GATES-SELFTEST-SPLIT` is next" phrasing is removed; the G6-4 mutual-exclusion governance candidate is registered; P92 adds no new observation item.
- (6) Section 6 pending rulings: the next ruling item = **naming the `G6-4 mutual-exclusion governance slice` to start**.
- (7) Annex A: obligation transfer = **none**; annex B is carried over unchanged; sections 10 / 11 roll forward, and the `DOCUMENTATION_PLAN{,_EN}` pointer moves to v76.
- (8) Mechanical readings re-checked: ADR 48 / **OB 145** / DR 18 / trial deferral 25 (P92 adds none, verified item by item).
---

## 1. Current anchors (check these seven numbers first)

**Rolling-window caliber = the most recent three versions (v74 / v75 / v76)**; v73 has been evicted (the same eviction rule as v52).

| Item | Value | Verification command |
|---|---|---|
| HEAD | `10a1f95` | `git rev-parse --short HEAD` |
| origin/main | `10a1f95` | `git rev-parse --short origin/main` |
| Directive version | **v1.75** | head of `docs/DELIVERY_DIRECTIVE.md` |
| SELFTEST | **744/744** | `node tools/gates/gate.js --selftest` |
| Three scopes | **all TOTAL_FAIL=0** | `gate.js --all --scope=index` / `tree` / `ci` |
| v74 whole-file lines | **216** (LF count, CN / EN symmetric; **re-verified here**) | `docs/t027/CHECKLIST_v74{,_EN}.md` counted by LF |
| v75 whole-file lines | **215** (LF count, CN / EN symmetric; **re-verified here**) | `docs/t027/CHECKLIST_v75{,_EN}.md` counted by LF |

> The HEAD = `10a1f95` here is the generation-time anchor after the P92 backfill commit; once v76 lands the HEAD moves on, so trust `git rev-parse --short HEAD`. v76 own line count is left to v77.

> Mechanical readings: ADR max = **48**; OB max = **145**; DR max = **18**; trial deferral max = the 25th. P92 verified item by item: no new ADR, no new OB, no new DR and no deferral change.

> Rolling window = the most recent three versions (v74 / v75 / v76) and v73 is evicted; if two or more items disagree, trust the repository and report the final readings first.

Fixed `tmp` assertion: top level = `.gitkeep` + `gates-cache` (only `g4b` below it) + `find-quotes.js` (no dangling reference; still awaiting the tmp-lifecycle slice). **The P92 temp directory was removed in (8c)** and is not part of this assertion.

---

## 2. P92 fact pack

- **Slice / tier**: `GATES-SELFTEST-SPLIT` (`--selftest` partition timing and fixture-read memoisation plus the benefit gate) / **M**; controller `DeepSeek V4.1 Flash` High; named by the user on 2026-11-26 with same-turn authorisation for the commit and backfill.
- **Goal**: `--selftest` partition timing (composition data on disk) plus fixture-read memoisation, a G4b inner-loop increment and a built-in `--timings` evaluation, with **byte-identical criterion verdicts**; the chain starts with (1) impact ruling plus the **benefit gate**.
- **Verdict**: **`CLOSED`** (2026-11-26); **directive bumped to v1.75 (MINOR)**; zero change to the criterion bodies and their decision logic.
- **Benefit gate (measured)**: **98.5% of the `--selftest` wall clock sits in child processes** (483 spawns, of which 413 are `git rev-parse --show-toplevel` = 75.8 s), in-process disk reads are only 0.27% and fixture reads total < 60 ms, so (b) fixture memoisation, (d) the G4b inner loop and (e) the built-in `--timings` are **moved out of scope**; the added candidate (c) "repo-root discovery memo" clears the bar and lands.
- **Landing**: `tools/gates/selftest.js` collapses the 23 `discoverRepoRoot(__dirname)` call sites into one process-level lazy `selftestRepoRoot()` (+39/-23); `--selftest` goes **150.9 s -> 70.1 s (2.15x)** with **byte-identical stdout** (A/B/B2 captures share one sha256) and `SELFTEST` 744/744.
- **Criteria and acceptance**: (3) A/B byte-identical with a one-byte flip turning `DIFFERENT`; (9) all three scopes `TOTAL_FAIL=0`, `SELFTEST` 744/744 and the base gates green; the report pair mirrors strictly; every one of the 21 pack payloads in `SHA256SUMS` recomputes to 0 mismatches.
- **Review chain**: (1) x1 (`DeepSeek V4 Pro`) -> (2)(3) `prfrail-implementer` / `prfrail-tester` (Flash) -> (5) 1/2 (`PASS WITH FIXES`, 2 Low, `no Medium+`) -> (6) 3/3 (r1/r2 `FINDINGS` on missing downstream artifacts or a mid-convergence pack hash, r3 `INDEPENDENT SCAN: PASS`) -> (7) 1+2 (r1 `FINDINGS` = three non-verbatim fence lines, restored; r2 `PASS WITH FIXES` = tree capture mid-state; r3 `FINDINGS` = stale section 5 reading) plus one narrow confirmation read (`CONFIRM: FINDINGS`, point A closed by a user ruling as a mechanical fact correction) -> (8a) / (9) / (8b) / (8c) -> (10).
- **Accounting**: **R = 3, D = 1, zero probes**. R = 3 = (i) registering `g6-whitelist.txt` turned the T23 assertion red (`SELFTEST` 743/744) and was reverted; (ii) one round trip between fence normalisation and verbatim restoration; (iii) two refills of the section 5 readings. D = 1 = the one-off whitelist registration that was reverted (self-healing tier). **No per-round stop as the P85 precedent suggests; summarised at (10)** - that deviation is recorded and was accepted by the user. Batch helper **GLM x2** (timing-data tidy-up / bilingual mirror structure check).
- **Two commits and CI**: main commit `785115f` (30 paths) -> run `45726904763` **both legs green on attempt 1** (`SELFTEST` 744/744, `scope=ci changed=30`, `TOTAL_FAIL=0`; the `Gates` step took 33 s on windows and 4 s on ubuntu) -> backfill `10a1f95` (2 paths, `+4/-80`), whose run is an observation record only per section 11.2.
- **Docs**: report pair `docs/validation/gates-selftest-split{,_EN}.md` (strictly mirrored CN/EN; the section 5 `G4b` reading is now a **reference-style** statement so the self-referential drift cannot recur; it carries the `artifacts` fence and section 9 commit and CI); evidence pack `docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/` (21 payloads plus `MANIFEST.md` and `SHA256SUMS.txt`); the ledger P92 completion record.
- **tmp and leftovers**: the P92 temp directory was removed in (8c) (the `tmp` top level is back to three items); `find-quotes.js` still awaits the tmp-lifecycle slice. **Conflict record (accepted by the user, not a precedent)**: the `G6-4` counting literals and the section 11.3 verbatim transcript are mutually exclusive, while `tools/gates/g6-whitelist.txt` is pinned empty by the T23 assertion (registering one entry yields 743/744), so this slice can only normalise minimally and disclose it outside the fence; registered as the `G6-4 mutual-exclusion governance slice` candidate (report section 7 item 5).

### Annex A: obligation-transfer section (this version: **none**)

- **Historical transfer** (first case traced to P75 O-7): the in-line caliber inconsistency at `s1-exit.md:52` was discharged by P76 change 1 plus S-G2 (`0f80c1e`).
- **This version**: **none** - P92 leaves no obligation; the `G6-4 mutual-exclusion governance slice` is the next item (pending naming) and `GATES-CRITERIA-ASYNC` / `GATES-BASE-CACHE` remain record-only, so neither constitutes a review or delivery obligation.
- **Convention note**: adopted since v59; the landing slice entry must carry an "admission condition" face.

### Annex B: AT-23 stage 2 status (unchanged)

- **Status**: A1 ruled "not applicable" (user ruling on 2026-11-09; not a permanent ban).
- **Executable object before S2**: none (R-3 = (A)).
- **Restart path**: if the user later rules to restart stage 2, a real enforcement hookup (S2) or an alternative path must exist first.
- **Natural restart point**: S2 starts -> enforcement really hookup -> S-G8 unlocks -> R-3 re-evaluated.

---

## 3. Completed slices (recent baseline)

**v63 section 3 held sixteen rows (P64-P80); v64 added P81 (seventeen); v65 added P82 (eighteen); v66 carried eighteen and added P79 (nineteen); v67 carried nineteen and added P83 (twenty); v68 carried twenty and added P84 (twenty-one); v69 carried twenty-one and added P86 (twenty-two); v70 carried the twenty-two rows verbatim and added P85 (twenty-three); v71 carried twenty-three verbatim and added P87 (twenty-four); v72 carried twenty-four verbatim and added P88 (twenty-five); v73 carried twenty-five verbatim and added P89 (twenty-six); v74 carried twenty-six verbatim and added P90 (twenty-seven); v75 carried twenty-seven verbatim and added P91 (twenty-eight); v76 carries v75 twenty-eight verbatim and adds P92. Twenty-nine rows in total.**

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
| v1.74 | P91 DIRECTIVE-CI-INVOCATION-SYNC | `f963dd7` + `3232d57` | **closed** (governance slice; the directive is bumped to **v1.74**, the section 6.3 CI-invocation sentence now reads `--all --skip-base --scope=ci` plus the base-coverage clause; OB-145 handled; criterion bodies untouched; SELFTEST 744/744; three scopes 0; R = 1, D = 0; main run `45688549781` green on both legs at attempt 1 plus observation `45689325140`) |
| v1.75 | P92 GATES-SELFTEST-SPLIT | `785115f` + `10a1f95` | **Closed** (tool slice; repo-root discovery memo in `selftest` 150.9 -> 70.1 s (2.15x) with byte-identical stdout; zero criterion-body change; directive bumped to **v1.75**; SELFTEST 744/744; three scopes 0; R = 3, D = 1; main run `45726904763` both legs green on attempt 1 plus observation `45727490167`) |

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

**Next item = the `G6-4 mutual-exclusion governance slice`**

- Source = P92 report section 7 item 5; **tier S** (directive / criterion clause plus one assertion policy); **the user must name it**; the four start-up elements must be on record.
- Scope: give a **formal channel** between the `G6-4` counting literals and the section 11.3 verbatim six-stage transcript (a transcript-only registration slot, or an explicit `G6-4` exemption for `sixstage` fences), syncing the `T23` assertion and the `g6-whitelist.txt` policy wording.
- Boundaries: it touches criterion implementations and the directive body, so **both (6) and (7) are mandatory**; the existing verdicts must stay byte-identical; no gitee.

- After that = the pre-S2 re-evaluation (S-G1 / S-G8) (the next T027 main-line stop, pending user scheduling):
  - Prerequisite = the S2 enforcement hookup; the independent security review prerequisite was discharged by P79 and P83 / P84 / P86 / P85 delivered the F1-F8 fixes, so the ADR-047 rollback clause no longer blocks;
  - `S-G8` stays blocked until the S2 enforcement hookup (P73 A2 / P75 A2 / ADR-044 / ADR-045); the `s1-exit{,_EN}.md` gate table is untouched.
- **Post-S2 re-evaluation items**: **S-G5 step two (overall S1 Exit approval) plus S-G7**, in the order of P80 report section 7.

### Tier 2 routine governance

- P72 / P73 / P74 / P75 / P76 / P77 / P78 / P80 / P81 / P82 / P79 / P83 / P84 / P86 / P85 / P87 / P88 / P89 / P90 / P91 / **P92** all ✅ discharged (tier 2 is **empty**).
- Each slice's full section 1.2 field set lives in `docs/t027/REMAINING_SLICES{,_EN}.md`.

### Tier 3 low-priority candidates

| Field | Content |
|---|---|
| Identifiers | `OB-50` / `spawnGit` / `ci frozen cost` / `%TEMP% historical leftovers, 22 items` / `GATES-EXEMPT` plus the two performance candidates (`GATES-CRITERIA-ASYNC` / `GATES-BASE-CACHE`) |
| Current state (verified item by item) | OB-50 = **handled** (section 12.4, v1.41); spawnGit = the git-call helper of `tools/gates/selftest.js` (spawn wrapper in `lib/ctx.js`, stderr captured into the result object and not printed), a diagnostic OB candidate from P27 report section 5 item 5 with no section 12.4 entry; ci frozen cost = handled by P90 (the `--skip-base` Gates step landed and the Go cache was judged already covered); %TEMP% leftovers = OB-98 handled, OB-99 / OB-100 under observation; GATES-EXEMPT = **not started** (no section 12.4 entry; split out in P33 (1)-3(iv), touches four criterion files); `OB-145` was handled by P91 (v1.74); `GATES-SELFTEST-SPLIT` left the tier when P92 closed it (v1.75) |
| Contents | See the matching section 12.4 rows; `GATES-PERF-ASSESS` (P88, v1.73) and `GATES-CI-PERF` (P90) both closed and left the tier; the two remaining performance candidates are in the table below |
| Dependencies | The long suspension covers only the unhandled items (spawnGit / the `%TEMP%` temporal shape OB-99, OB-100 / GATES-EXEMPT), none with a defined trigger |
| Size | Varies (S through M-L) |
| (7) quota | none |
| Boundaries | This version advances no suspended item and keeps the registrations; the two remaining performance candidates are record-only and constitute no slice |
| Start condition | the user must name it explicitly |

**Performance candidates (risk ascending; both rows below are record-only and need the user to name them)**:

| Identifier | One-line scope | Relation to existing items / preconditions |
|---|---|---|
| `GATES-CRITERIA-ASYNC` | make G6 / G7 / G8 asynchronous (argue benefit and the concurrency / ordering contract first) | reuses the existing `runCriteriaConcurrent` scheduling contract; involves concurrent reads of shared caches |
| `GATES-BASE-CACHE` | cache the base verdicts (reuse the `gofmt` / `build` / `vet` / `test` conclusions across scopes) | same family as section 6.3 G4b "never cache verdicts", so **a directive clause must exist first**; highest risk, listed last |

**Queue discipline**: tier boundaries and numbering follow v60; this version does not bump the directive.

---

## 5. Observation items

- Observing (**13**): `OB-1 / 7 / 10 / 33 / 87 / 99 / 100 / 104 / 105 / 107 / 111 / 113 / 135` (carried from v61).
- In discussion (9): `OB-19 / 21 / 52 / 53 / 57 / 58 / 59 / 61 / 75`.
- To do (1): `OB-69`. Non-terminal total **23** (13 + 9 + 1), matching this version verification.
- **DR registrations**: DR max = **18** (P85 / P86 / P87 / P88 / P89 / P90 / P91 / P92 added none).
- **P87 Hydrogen precondition gap (record-only)**: the light-tier clause (`Gemini 4 Hydrogen`) has landed while its comparison-review measurement is not done; it points at a later business slice (section 1.7.4 temporary authorization plus an OB record) and **must not be described as validated**; backfill a note into `docs/validation/directive-review-model-policy{,_EN}.md` once complete. Registration sites: report section 7 and the ledger P87 section.
- **G6-4 mutual-exclusion governance candidate (new here)**: the `G6-4` counting literals and the section 11.3 verbatim transcript are mutually exclusive, and `g6-whitelist.txt` is pinned empty by `T23`, so the conflict **has no compliant channel today**; P92 passed with "minimal normalisation plus an outside-fence disclosure", accepted by the user (not a precedent). Candidate = the `G6-4 mutual-exclusion governance slice` (tier S, needs naming).
- **Tier-3 performance registrations (rolled here)**: `GATES-CRITERIA-ASYNC` / `GATES-BASE-CACHE` stay record-only (the user must name them); `GATES-SELFTEST-SPLIT` left the tier when P92 closed it; `GATES-CI-PERF` left when P90 closed it.
- **`find-quotes.js` disposition**: no dangling reference (P89 re-checked it and excluded it from the archive and rewrite work); still awaiting the tmp-lifecycle slice.
- Mechanical readings: section 12.4 OB max = **145**; ADR max = **48**; DR max = **18**; trial deferral max = the 25th.
- **Removed here**: the "`GATES-SELFTEST-SPLIT` is next" phrasing (now closed). **Added here**: the G6-4 mutual-exclusion governance candidate (P92 adds no other observation item).

---

## 6. Pending user rulings

- **This table is currently empty** (the P92 authorisation sequence, the R/D caliber, the narrow (7) confirmation read and this v76 snapshot were all ruled on 2026-11-26).
- **Next ruling item (predicted)**: **naming the `G6-4 mutual-exclusion governance slice` to start** (tier S; touches directive sections 6.3 / 11.3 plus the `T23` assertion and the `g6-whitelist.txt` policy; (6) and (7) mandatory).

---

## 7 / 8 / 9 Pointer sections

Same structure as v60: section 7 start-up and step 0 general rules -> directive section 13.1; section 8 cross-slice collaboration discipline -> section 13.2; section 9 checklist meta-rules and gate pits -> section 13.3 / 13.4. **The rule bodies live in directive section 13 and are not copied here**.

---

## 10. One-line opening prompt for a new session

> Below is my project status snapshot (v76); check the seven anchors in section 1 first, continue with the section 4 queue if they match, and otherwise trust the repository.

> [paste this whole v76 checklist]

> [Supporting docs] `REMAINING_SLICES{,_EN}.md` (P92 closed with a completion record; the `G6-4 mutual-exclusion governance slice` is next, pending naming), `DELIVERY_DIRECTIVE.md` (**v1.75**), `gates-selftest-split{,_EN}.md` (the P92 report pair), `docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/` (the P92 evidence pack), `ADR_REGISTER{,_EN}.md` (ADR-048 unchanged), `DOCUMENTATION_PLAN{,_EN}.md` (hand-off pointer already at v76), `CHECKLIST_v75{,_EN}.md` (frozen, superseded by v76).

> Now: P92 closed (`--selftest` memo makes 150.9 -> 70.1 s plus the v1.75 package; both legs SELFTEST 744/744 and TOTAL_FAIL=0; R = 3, D = 1); next = the `G6-4 mutual-exclusion governance slice` (tier S, needs the user to name it); the pre-S2 re-evaluation (S-G1 / S-G8) is the next T027 main-line stop, pending user scheduling (section 4 "after that").

## 11. Current status in one line

P92 **closed** (main commit `785115f` plus backfill `10a1f95` pushed to origin; the main run `45726904763` was green on both legs at attempt 1 and the backfill run is observation-only per section 11.2; both legs SELFTEST 744/744, `scope=ci changed=30`, TOTAL_FAIL=0; R = 3, D = 1; **directive bumped to v1.75**; the `tmp` top level is back to three items); next = **the `G6-4 mutual-exclusion governance slice`** (tier S, touches directive sections 6.3 / 11.3 and the `T23` assertion, (6) and (7) mandatory, needs the user to name it); the pre-S2 re-evaluation (S-G1 / S-G8) is the next T027 main-line stop, pending user scheduling (section 4 "after that").