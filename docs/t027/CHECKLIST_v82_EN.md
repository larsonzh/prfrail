# Task checklist (cross-session snapshot v82 - after P98 closes - queue empty)

Generated: 2026-11-28. HEAD = `edeb4e0`; directive = v1.83 (P98 bumped it); P98 CI main-round chain: `46093203794` / `46093550773` / `46095600247` / `46096087094` / `46097218281` (all both legs success); the backfill and correction rounds are observation records only per §11.2 (API-level verdict).

**Replacement**: this version replaces v81 (`docs/t027/CHECKLIST_v81{,_EN}.md` frozen, not rewritten).

**Structure**: same as v79 (§1 anchors -> §2 fact pack -> §3 completed slices -> §4 execution queue -> §5 watch items -> §6 pending rulings -> §7 / §8 / §9 pointer sections -> §10 opening -> §11 one-line status); the rule bodies of §7-§9 live in directive §13, so this snapshot does not copy them.

**Changes versus v81**:

- (1) §1 anchors rolled over: HEAD `edeb4e0`; directive **v1.83** (the P98 revision version); SELFTEST **806/806**; the rolling window drops v79 and takes v80 / v81 / v82, with v80 / v81 re-checked here (210 / 213) and v82 itself left to v83.
- (2) §2 becomes the P98 fact pack (`CLOSED`, governance slice; the seven-step commit chain, both-leg CI readings, **R = 0, temp-artifact D = 1**, zero probes; the directive moved v1.82 -> v1.83; O-2..O-4 / O-5 flipped to handled (P98); **major deviation registered**: seven unauthorised commits retroactively ratified by the user).
- (3) §3 appends one P98 row (the thirty-four rows of v81 plus one = thirty-five).
- (4) §4 queue update: P98 leaves; **next item = none (queued work paused)**; "after that" stays the S2 pre-re-evaluation (S-G1 / S-G8, pending S2 enforcement).
- (5) §5 watch items: O-2..O-4 / O-5 are removed (all handled); O-1 stays (record only); the 23 non-terminal OB items are copied; this version adds no watch item.
- (6) §6 pending rulings: **none**; the expected next ruling = scheduling the OB cleanup slice / the spawnGit assessment.
- (7) Appendix A: the obligation hand-over is **none**; appendix B is copied unchanged; §10 / §11 rolled over, and the `DOCUMENTATION_PLAN{,_EN}` pointer now targets v82.
- (8) Mechanical readings re-checked: ADR 48 / OB 146 / DR 18 / trial deferral 25 (P98 added no ADR / OB / DR, each checked).
---

## §1 Current anchors (check these 7 numbers first)

**Rolling-window caliber = the latest 3 versions (v80 / v81 / v82)**; v79 has been squeezed out (same rule as v52).

| Item | Value | Check command |
|---|---|---|
| HEAD | `edeb4e0` | `git rev-parse --short HEAD` |
| origin/main | `edeb4e0` | `git rev-parse --short origin/main` |
| Directive version | **v1.83** | the header of `docs/DELIVERY_DIRECTIVE.md` |
| SELFTEST | **806/806** | `node tools/gates/gate.js --selftest` |
| Three scopes | **all TOTAL_FAIL=0** | `gate.js --all --scope=index` / `tree` / `ci` |
| v80 whole-file lines | **210** (LF count, CN / EN symmetric; **re-checked here**) | `docs/t027/CHECKLIST_v80{,_EN}.md` by LF count |
| v81 whole-file lines | **213** (LF count, CN / EN symmetric; **re-checked here**) | `docs/t027/CHECKLIST_v81{,_EN}.md` by LF count |

> The HEAD = `edeb4e0` here is the generation-time anchor after the P98 correction commit; once v82 lands HEAD moves on, so `git rev-parse --short HEAD` is authoritative then. v82 itself is left to v83 for its line count.

> Mechanical readings: ADR max = **48**; OB max = **146**; DR max = **18**; trial deferral max = the 25th time. P98 checked each: no new ADR, no new OB, no new DR, the deferral count unchanged.

> Rolling window = the latest 3 versions (v80 / v81 / v82); v79 has been squeezed out. If 2 or more items disagree, the repository wins: report the final readings first and then continue.

`tmp` fixed assertion: top level = `.gitkeep` + `gates-cache` (holding `base` and `g4b`). **The P98 temp directory was cleared in (8c)**, so this assertion does not include it.

## §2 P98 fact pack

- **Slice / size**: `DIRECTIVE-REVIEW-TIER-ALIGN` (the (7) tier-alignment revision; governance `[SLICE]`) / **M**; controller `DeepSeek V4 Pro` Max (a user ruling, not a precedent); the user named it on 2026-11-28, and stop 1 adopted all five items.
- **Goal**: (1) align the section 5.2 deep-tier tier column and tier-authorization sentence with section 8.2 (O-2 High) plus the matrix decidability wording (O-3 Medium); (2) fix the section 7.5 Neon / Carbon price caliber from an operator-UI read (O-4 Medium) plus the alternative-quota wording (O-5 #4); (3) register O-5 #5 without fixing it (frozen evidence-pack file); (4) flip the watch items to handled (P98).
- **Verdict**: **`CLOSED`** (2026-11-28, governance slice); the section 5.2 Neon tier = `Extra High`, byte-identical with section 8.2 and the tier-authorization sentence aligned too; the section 7.5 prices follow the operator-UI read and no longer split by band (Neon 300 / 1500 and Carbon 200 / 1000 are the full-tier list prices); zero change to the criteria implementation and `SELFTEST` stays 806/806.
- **Registration**: the directive went **v1.82 -> v1.83** (MINOR: header + §0.3 row + §5.2 / §7.5 / §8.2 revisions + appendix C.0bf); ADR-048 carries the P98 revision note; the O-2..O-4 / O-5 ledger flips are mirrored in both languages.
- **Deliverables**: report pair `docs/validation/directive-review-tier-align{,_EN}.md` (strictly mirrored, with the `artifacts` fence and the (6) verdict); evidence pack `docs/validation/evidence/P98-DIRECTIVE-REVIEW-TIER-ALIGN/`; the P98 ledger section and completion record plus the P97 watch-item flips.
- **Review chain**: (5) one (`PRE-REVIEW: PASS` plus 3 Info) -> (6) one repeat run ((7) r1 judged §7.2.1 3 not satisfied, so the §7.2 repeat-1-of-3 applies, `MAI Code 2 Flash`, `INDEPENDENT SCAN: FINDINGS`) -> (7) standard tier `GPT-5.3-Codex` one final plus three re-reviews (`FINDINGS` -> `FINDINGS` -> `RE-REVIEW: PASS`) -> (8a) / (9) / (8b) / (8c) -> (10).
- **Metering**: **R = 0, temp-artifact D = 1, zero probes** (D = after (8c) `build-pack-step0.js` was re-flushed by the editor buffer; the user ruled a precise cleanup and registration). **Major deviation registered**: seven commits during (2)-(10) were made without same-turn authorisation; the user retroactively ratified them on 2026-11-28 and ruled they be registered, **not a precedent**; the corrected caliber = "same-turn authorisation" is henceforth read as "a single authorisation after the (10) stop; any commit / push inside a slice needs its own same-turn authorisation".
- **Commits and CI**: `a3fa412` (clause) -> `0c6820b` (ledger) -> `91a9821` -> `f972214` -> `6bb3e90` (report pair and closeout) -> backfill `cf54aa2` -> correction `edeb4e0` (deviation and residue registration); CI runs `46093203794` / `46093550773` / `46095600247` / `46096087094` / `46097218281` all both legs success (both legs `SELFTEST: PASS 806/806`, `TOTAL_FAIL=0`); the backfill and correction rounds are observation records per §11.2.
- **tmp and leftovers**: the P98 temp directory was cleared in (8)c (the `build-pack-step0.js` residue was precisely deleted and read back green); `tmp` is back to two items; O-2..O-4 / O-5 are all handled.

### Appendix A: obligation hand-over (this version: **none**)

- **Historical hand-over** (first traced to P75 O-7): the inconsistent wording inside the `s1-exit.md:52` core-chain row was honoured by P76 change 1 plus the S-G2 correction (`0f80c1e`).
- **Hand-over in this version**: **none** - P98 leaves no obligation; the queue is empty (next item = none), so no review or delivery obligation is handed over.
- **Convention**: in use since v59; a landing slice entry must carry the admission-condition face.

### Annex B: AT-23 stage 2 status (unchanged)

- **Status**: A1 ruled "not applicable" (user ruling on 2026-11-09; not a permanent ban).
- **Executable object before S2**: none (R-3 = (A)).
- **Restart path**: if the user later rules to restart stage 2, a real enforcement hookup (S2) or an alternative path must exist first.
- **Natural restart point**: S2 starts -> enforcement really hookup -> S-G8 unlocks -> R-3 re-evaluated.

---

## 3. Completed slices (recent baseline)

**v63 §3 already held sixteen P64-P80 rows; v64 appended P81 (seventeen); v65 appended P82 (eighteen); v66 kept the eighteen of v65 and appended P79 (nineteen); v67 kept the nineteen of v66 and appended P83 (twenty); v68 kept the twenty of v67 and appended P84 (twenty-one); v69 kept the twenty-one of v68 and appended P86 (twenty-two); v70 kept the twenty-two of v69 (readings copied, not rewritten) and appended P85 (twenty-three); v71 kept the twenty-three of v70 (readings copied, not rewritten) and appended P87 (twenty-four); v72 kept the twenty-four of v71 (readings copied, not rewritten) and appended P88 (twenty-five); v73 kept the twenty-five of v72 (readings copied, not rewritten) and appended P89 (twenty-six); v74 kept the twenty-six of v73 (readings copied, not rewritten) and appended P90 (twenty-seven); v75 kept the twenty-seven of v74 (readings copied, not rewritten) and appended P91 (twenty-eight); v77 kept the twenty-nine of v76 (readings copied, not rewritten) and appended P93; v78 keeps the thirty of v77 (readings copied, not rewritten) and appends P94; v79 keeps the thirty-one of v78 (readings copied, not rewritten) and appends P95; v80 keeps the thirty-two of v79 (readings copied, not rewritten) and appends P96; v81 keeps the thirty-three of v80 (readings copied, not rewritten) and appends P97; v82 keeps the thirty-four of v81 (readings copied, not rewritten) and appends P98. Thirty-five rows in total.**

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
| v1.72 | P87 DIRECTIVE-REVIEW-MODEL-POLICY | `8b64ce3` + `2782a09` | **closed** (directive rises to **v1.72**; (7) four-tier selection matrix plus the grading-threshold consensus plus the new section 5.4; G5-b closed set of 11; ADR-048; R = 0, D = 0; run `45294006169` both legs green plus observation `45356014951`; the Hydrogen pre-gap was filled by P97 and the tier conflict was revised by P98) |
| v1.73 | P88 GATES-PERF-ASSESS | `03c45f7` + `9b44dbf` | **closed** (tooling slice; criterion bodies and decision logic untouched, five same-input A/B pairs byte-equal; SELFTEST 744/744; three scopes 0; CI critical path 380.3 -> 141.8 s; R = 0, D = 0; run `45523666470` green first try plus observation `45524509376`) |
| v1.73 | P89 TMP-ARCHIVE-P84-P85 | `5b85f7a` + `a8eccc9` + `e15ad3d` + `9a784dc` | **closed** (tooling / evidence slice; 110 items archived whole, three packs three-way equal at P84 58 / P85 95 / P89 14, residual hits 16 = the 16-entry whitelist; R = 0, **D = 1**; main run `45584959722` both legs green -> backfill run `45587252859` `Gates` red (`G8-b` CRLF hashing) -> fix run `45589128432` both legs green -> closing run `45590253878` both legs green) |
| v1.73 | P90 GATES-CI-PERF | `ec221f5` + `9ac8f6b` | **closed** (tooling slice; `Gates` runs `--skip-base` plus a new `Format` step, criterion bodies untouched, (3) A/B byte-equal after stripping base; SELFTEST 744/744; three scopes 0; `Gates` 103->45 s on windows and 44->6 s on ubuntu; R = 1, D = 0; main run `45659843522` green on both legs at attempt 1 plus observation `45660930671`) |
| v1.74 | P91 DIRECTIVE-CI-INVOCATION-SYNC | `f963dd7` + `3232d57` | **closed** (governance slice; the directive is bumped to **v1.74**, the section 6.3 CI-invocation sentence now reads `--all --skip-base --scope=ci` plus the base-coverage clause; OB-145 handled; criterion bodies untouched; SELFTEST 744/744; three scopes 0; R = 1, D = 0; main run `45688549781` green on both legs at attempt 1 plus observation `45689325140`) |
| v1.75 | P92 GATES-SELFTEST-SPLIT | `785115f` + `10a1f95` | **Closed** (tool slice; repo-root discovery memo in `selftest` 150.9 -> 70.1 s (2.15x) with byte-identical stdout; zero criterion-body change; directive bumped to **v1.75**; SELFTEST 744/744; three scopes 0; R = 3, D = 1; main run `45726904763` both legs green on attempt 1 plus observation `45727490167`) |
| v1.76 | P93 DIRECTIVE-G6-4-TRANSCRIPT-EXEMPT | `23afa2d` + `04f2c0d` + `d4ae682` | **Closed** (governance slice; channel = b (a G6-4 scope sub-rule: a CLOSED `sixstage` fence whose five sections are in order is exempt), zero change to the decision logic, single-source reuse of the G7 internals, three fail-closed faces; directive bumped to **v1.76**; SELFTEST 744 -> 757 (T133-T143); R = 4, D = 0; main run `45782329377` green on both legs plus correction `04f2c0d` and backfill `d4ae682`) |
| v1.77 | P94 GATES-EXEMPT | `49f1609` + `afd9fb9` | **Closed** (governance; the three merged verdicts = the exemption extends to G2 + G4b / no generic registry slot / whitelist threshold text only; zero change to the decision logic; directive **v1.77**; SELFTEST 757 -> 770 (T144-T156); R = 6, D = 0; main run `45910223000` green first try + backfill `45910797198` observation) |
| v1.77 | P95 GATES-CRITERIA-ASYNC | `238e45c` + `bea554b` | **Closed** (assessment class; the (1) benefit gate = do not do it (ctx fully synchronous => benefit ≈0); P88 recommendation 4 registered as "assessed, benefit does not hold, not implemented"; zero code change, no directive bump; R = 1, D = 0; main run `45939437180` both legs green + backfill `45940170692` observation) |
| v1.81 | P96 GATES-BASE-CACHE | `81a6794` + `90932ee` | **Closed** (tooling slice; the section 6.1 base verdict cache clause; `gofmt` always runs for real and the cache covers build / vet / test only; directive v1.78 -> v1.81 in four amendments; SELFTEST 770 -> 806 (T157-T165); R = 12, D = 0; main run `46053254410` changed=17 + backfill run `46056837161` changed=8, both legs green) |
| v1.82 | P97 HYDROGEN-CONTROL-PILOT | `4f3af5f` + `c30b2de` | **Closed** (docs-only measurement slice; directive bumped to **v1.82** plus OB-146; one Hydrogen comparative review round returned `FINDINGS` (1 High + 2 Medium + 2 Low, 184.3 credits); P87 is never re-adjudicated and its pre-gap is filled; R = 0, D = 0; main run `46077231907` changed=18 + backfill run `46078089184` changed=4, both legs green) |
| v1.83 | P98 DIRECTIVE-REVIEW-TIER-ALIGN | `a3fa412` + `6bb3e90` + `cf54aa2` + `edeb4e0` | **Closed** (governance slice; O-2..O-4 / O-5 all handled: the section 5.2 tier column aligned with section 8.2 and the section 7.5 prices read from the operator UI; directive bumped to **v1.83**; SELFTEST 806/806 unchanged; R = 0, **temp-artifact D = 1**; **major deviation registered**: seven unauthorised commits retroactively ratified by the user, not a precedent; runs `46093203794` / `46093550773` / `46095600247` / `46096087094` / `46097218281` all both legs green) |

> **D caliber of section 3** (same as v61): "D = n" here abbreviates the temp-artifact discipline violation count (section 5.1); the P98 **D = 1** is `build-pack-step0.js` being re-flushed by the editor buffer after (8c) (precisely cleaned per the user ruling); the P97 D-1 belongs to the report deviation-table series and is not written into this table. **The D count and the section 9.2 R count are two independent series**.

---

## 4. Execution queue

### Tier 0 pre-work (✅ discharged; nothing new here)

- P5 / P6 were closed by P57; tier 0 is empty with no residue.

### Tier 1 main-line hard gates

- Items 8-12: all ✅ discharged (see v57 / v58 section 4; trace back through that reference chain if needed).
- **`T027-B` (defensive-refactor phase) = done (partly not applicable)**: C1-C7 / N1-N4 / R-1 / R-2 / R-3 are all handled (R-3 = (A)). **The main line is fully discharged**.
- **`S1 Exit closeout` (gate level done; security remediation F1-F8 fully settled)**: P76 produced the closeout roadmap; P77 closed S-G4; P78 closed S-G3; P80 judged S-G6 to hold at the gate level; the independent security review was completed by P79 and P83 / P84 / P86 / P85 delivered the F1-F8 fixes. Gate table below (P87 / P88 / P89 / P96 / P97 / P98 left it unchanged):

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

**Next item = none (queue empty, queued work paused)**

- Tiers 1 and 2 are both cleared and the tier-3 performance candidates are empty in this version; this version advances nothing.
- Ways back into the queue: the only entry is real S2 enforcement landing, which unlocks the S2 pre-re-evaluation (S-G1 / S-G8); or the user names a new slice (for example the OB cleanup slice / the spawnGit assessment).

- After that = the pre-S2 re-evaluation (S-G1 / S-G8) (the next T027 main-line stop, pending user scheduling):
  - Prerequisite = the S2 enforcement hookup; the independent security review prerequisite was discharged by P79 and P83 / P84 / P86 / P85 delivered the F1-F8 fixes, so the ADR-047 rollback clause no longer blocks;
  - `S-G8` stays blocked until the S2 enforcement hookup (P73 A2 / P75 A2 / ADR-044 / ADR-045); the `s1-exit{,_EN}.md` gate table is untouched.
- **Post-S2 re-evaluation items**: **S-G5 step two (overall S1 Exit approval) plus S-G7**, in the order of P80 report section 7.

### Tier 2 routine governance

- P72 / P73 / P74 / P75 / P76 / P77 / P78 / P80 / P81 / P82 / P79 / P83 / P84 / P86 / P85 / P87 / P88 / P89 / P90 / P91 / P92 / P93 / P94 / P95 / P96 / P97 / **P98** ✅ left the queue (tier 2 is **fully cleared**).
- Each slice's full section 1.2 field set lives in `docs/t027/REMAINING_SLICES{,_EN}.md`.

### Tier 3 low-priority candidates

| Field | Content |
|---|---|
| Identifier | `OB-50` / `spawnGit` / `ci frozen cost` / `22 historical %TEMP% items` (the performance set is empty) |
| Current state (each verified) | OB-50 = **disposed** (§12.4, v1.41); spawnGit = the git helper of `tools/gates/selftest.js` (the `lib/ctx.js` spawn wrapper keeps stderr in the result object and never prints it), a candidate OB (diagnostic) from P27 report §5 item 5 with no §12.4 row; the ci frozen cost = disposed by P90; %TEMP% leftovers = OB-98 disposed / OB-99 and OB-100 watching; `OB-145` disposed by P91 (v1.74); `GATES-SELFTEST-SPLIT` left the tier when P92 closed (v1.75); `GATES-EXEMPT` left the tier when P94 closed (v1.77); `GATES-CRITERIA-ASYNC` left the tier when P95 closed (assessed and not implemented); `GATES-BASE-CACHE` left the tier when P96 closed (v1.81) |
| Contents | See the matching section 12.4 rows; `GATES-PERF-ASSESS` (P88, v1.73) and `GATES-CI-PERF` (P90) both closed and left the tier; the performance candidates are in the table below |
| Dependency | the long suspension covers only the undisposed items (spawnGit / the `%TEMP%` timing forms OB-99 and OB-100), none with a defined trigger |
| Size | Varies (S through M-L) |
| (7) quota | none |
| Boundaries | This version advances no suspended item and keeps the registrations; the performance set is empty, so no slice can be opened from it |
| Start condition | the user must name it explicitly |

**Performance candidates (risk ascending; this version has the empty set)**:

| Identifier | One-line scope | Relation to existing items / preconditions |
|---|---|---|
| none | `GATES-BASE-CACHE` left the tier when P96 closed (v1.81); this version holds no record-only performance candidate | any later addition still requires a directive clause first |

**Queue discipline**: tier boundaries and numbering follow v60; this version does not bump the directive.

---

## §5 Watch items

- Watching (**13**): `OB-1 / 7 / 10 / 33 / 87 / 99 / 100 / 104 / 105 / 107 / 111 / 113 / 135` (carried from v61).
- Under discussion (9): `OB-19 / 21 / 52 / 53 / 57 / 58 / 59 / 61 / 75`.
- To do (1): `OB-69`. Non-terminal total **23** (13 + 9 + 1), matching this version re-check.
- **DR registry**: DR max = **18** (none of P96 / P97 / P98 added a DR).
- **O-1 (the P96 record slot, record only)**: the two P96 out-of-quota narrow re-reviews were not registered in §12.4 as an independent record against §1.7.4 item 2 (present in `freedom-list` / the ledger, absent from §12.4); the user ruled on 2026-11-28 to **record only and not retro-fix**.
- **O-2..O-4 / O-5 (handled)**: O-2 (the section 5.2 deep-tier column conflicting with section 8.2), O-3 (the matrix lacking literal decidability) and O-4 (the section 7.5 Neon price caliber) were revised by P98 `DIRECTIVE-REVIEW-TIER-ALIGN` (directive v1.83); O-5 #4 (the alternative-quota wording) was trimmed by P98 and #5 (the stale `numstat` self-count) is registered as "known, not retro-fixed". This version no longer keeps those four slots.
- **Tier-3 performance candidates (rolled over here)**: `GATES-BASE-CACHE` left the tier when P96 closed (v1.81); `GATES-CRITERIA-ASYNC` left the tier when P95 closed (assessed and not implemented); `GATES-SELFTEST-SPLIT` left the tier when P92 closed; `GATES-CI-PERF` left it when P90 closed; `GATES-EXEMPT` left it when P94 closed.
- **`find-quotes.js` disposition slot**: precisely deleted during P96 (8)c and registered in the P96 ledger completion record; this version no longer keeps the slot.
- Mechanical readings: the largest OB in §12.4 = **146**; ADR max = **48**; DR max = **18**; trial deferral max = the 25th time.
- **Removed here**: the four O-2..O-4 / O-5 watch slots (handled (P98)). **Added here**: none (P98 added no watch item).

## §6 Pending user rulings

- **This table is empty** (the P98 stop-1 five items, the clause and ledger commit authorisation, and the backfill and correction authorisations were all ruled on 2026-11-28).
- **Expected next ruling**: **scheduling the OB cleanup slice / the spawnGit assessment**, pending user naming.

## 7 / 8 / 9 Pointer sections

Same structure as v60: section 7 start-up and step 0 general rules -> directive section 13.1; section 8 cross-slice collaboration discipline -> section 13.2; section 9 checklist meta-rules and gate pits -> section 13.3 / 13.4. **The rule bodies live in directive section 13 and are not copied here**.

---

## §10 One-line opening for a new session

> Below is my project state snapshot (v82). Check the 7 anchors in section 1 first; if they match, continue with the section 4 queue, otherwise the repository wins.

> [paste this checklist v82 in full]

> [supporting documents] `REMAINING_SLICES{,_EN}.md` (P98 closed with a completion record; the queue is empty), `DELIVERY_DIRECTIVE.md` (**v1.83**), `directive-review-tier-align{,_EN}.md` (the P98 report pair), `docs/validation/evidence/P98-DIRECTIVE-REVIEW-TIER-ALIGN/` (the P98 evidence pack), `ADR_REGISTER{,_EN}.md` (ADR-048 stands with the P98 revision note), `DOCUMENTATION_PLAN{,_EN}.md` (the hand-off pointer now targets v82), `CHECKLIST_v81{,_EN}.md` (frozen, replaced by v82).

> Now: P98 closed (governance slice: the section 5.2 tier column aligned with section 8.2, the section 7.5 prices read from the operator UI, O-2..O-5 all handled, directive v1.82 -> v1.83; both legs SELFTEST 806/806, `TOTAL_FAIL=0`; R = 0, temp-artifact D = 1); next = none (queue empty); the OB cleanup slice / spawnGit assessment awaits naming; the S2 pre-re-evaluation (S-G1 / S-G8) unlocks once real S2 enforcement lands.

## §11 One-line current status

P98 **closed** (the seven-step chain `a3fa412` -> `0c6820b` -> `91a9821` -> `f972214` -> `6bb3e90` -> backfill `cf54aa2` -> correction `edeb4e0` pushed to origin; the five main-round runs are both-legs green and the backfill and correction rounds are observation-only per §11.2; both legs `SELFTEST 806/806`, `TOTAL_FAIL=0`; R = 0, temp-artifact D = 1; **directive v1.83**; **the major deviation is registered**: seven unauthorised commits retroactively ratified by the user, not a precedent, and inside a slice every commit / push now needs its own same-turn authorisation; `tmp` back to two items); next = **none (queue empty)**; the OB cleanup slice / spawnGit assessment awaits naming; the S2 pre-re-evaluation (S-G1 / S-G8) unlocks once real S2 enforcement lands.
