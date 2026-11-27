# Task checklist (cross-session snapshot v80 - after P96 closes - queue empty)

Generated: 2026-11-28. HEAD = `90932ee`; directive = v1.81 (P96 bumped it); P96 CI main round: `46053254410` (both legs success); the backfill round `46056837161` is an observation record only per §11.2 (API-level verdict).

**Replacement**: this version replaces v79 (`docs/t027/CHECKLIST_v79{,_EN}.md` frozen, not rewritten).

**Structure**: same as v79 (§1 anchors -> §2 fact pack -> §3 completed slices -> §4 execution queue -> §5 watch items -> §6 pending rulings -> §7 / §8 / §9 pointer sections -> §10 opening -> §11 one-line status); the rule bodies of §7-§9 live in directive §13, so this snapshot does not copy them.

**Changes versus v79**:

- (1) §1 anchors rolled over: HEAD `90932ee`; directive **v1.81** (four amendments from P96); SELFTEST **806/806**; the rolling window drops v77 and takes v78 / v79 / v80, with v78 / v79 re-checked here (211 / 209) and v80 itself left to v81.
- (2) §2 becomes the P96 fact pack (`CLOSED`, tooling class; main commit `81a6794` plus backfill `90932ee`, both-leg CI readings, **R = 12**, D = 0, zero probes; the directive moved v1.78 -> v1.81 in four amendments within one slice; the cache covers build / vet / test only and `gofmt` always runs for real).
- (3) §3 appends one P96 row (the thirty-two rows of v79 plus one = thirty-three).
- (4) §4 queue update: P96 leaves; **the tier-3 performance candidates are cleared** (`GATES-BASE-CACHE` leaves the tier); **next item = none (queued work paused)**; "after that" stays the S2 pre-re-evaluation (S-G1 / S-G8, pending S2 enforcement).
- (5) §5 watch items: the `find-quotes.js` disposition slot is removed (precisely deleted during P96 (8)c and registered); the P87 Hydrogen precondition gap stays; P96 adds no new watch item.
- (6) §6 pending rulings: **none** (waiting for the user to schedule S2 / enforcement).
- (7) Appendix A: the obligation hand-over is **none**; appendix B is copied unchanged; §10 / §11 rolled over, and the `DOCUMENTATION_PLAN{,_EN}` pointer now targets v80.
- (8) Mechanical readings re-checked: ADR 48 / OB 145 / DR 18 / trial deferral 25 (P96 added no ADR / OB / DR, each checked).
---

## §1 Current anchors (check these 7 numbers first)

**Rolling-window caliber = the latest 3 versions (v78 / v79 / v80)**; v77 has been squeezed out (same rule as v52).

| Item | Value | Check command |
|---|---|---|
| HEAD | `90932ee` | `git rev-parse --short HEAD` |
| origin/main | `90932ee` | `git rev-parse --short origin/main` |
| Directive version | **v1.81** | the header of `docs/DELIVERY_DIRECTIVE.md` |
| SELFTEST | **806/806** | `node tools/gates/gate.js --selftest` |
| Three scopes | **all TOTAL_FAIL=0** | `gate.js --all --scope=index` / `tree` / `ci` |
| v78 whole-file lines | **211** (LF count, CN / EN symmetric; **re-checked here**) | `docs/t027/CHECKLIST_v78{,_EN}.md` by LF count |
| v79 whole-file lines | **209** (LF count, CN / EN symmetric; **re-checked here**) | `docs/t027/CHECKLIST_v79{,_EN}.md` by LF count |

> The HEAD = `90932ee` here is the generation-time anchor after the P96 backfill; once v80 lands HEAD moves on, so `git rev-parse --short HEAD` is authoritative then. v80 itself is left to v81 for its line count.

> Mechanical readings: ADR max = **48**; OB max = **145**; DR max = **18**; trial deferral max = the 25th time. P96 checked each: no new ADR, no new OB, no new DR, the deferral count unchanged.

> Rolling window = the latest 3 versions (v78 / v79 / v80); v77 has been squeezed out. If 2 or more items disagree, the repository wins: report the final readings first and then continue.

`tmp` fixed assertion: top level = `.gitkeep` + `gates-cache` (holding `base` and `g4b`). `find-quotes.js` was precisely deleted during P96 (8)c and registered in the ledger; **the P96 temp directory was cleared in (8c)**, so this assertion does not include it.

## §2 P96 fact pack

- **Slice / size**: `GATES-BASE-CACHE` (cross-scope caching of base verdicts; tooling `[SLICE]`) / **M**; controller `DeepSeek V4.1 Flash`; the user named it on 2026-11-27, and stop 1 was ruled as do it + option 1 + deep tier 1+3 (a direct user ruling, not a precedent).
- **Goal**: (1) a benefit gate - measure the four isolated base timings and the repeated payment across the three scopes, and produce a key design and a clause draft; (2) clause first - land the base verdict cache clause in §6.1 (directive MINOR bump) with the criteria code untouched and code only after CI is green; (3) implement the cross-scope base verdict cache (key / three counterexamples / fail-open); (4) zero change to the criteria decision logic, base verdict semantics unchanged, and never cache any G1-G9 criteria verdict.
- **Verdict**: **`CLOSED`** (2026-11-28, tooling class); the cache covers only the three base verdicts `build` / `vet` / `test`, while `gofmt` always runs for real and is never cached (its verdict depends on the PATH / execvp / libuv resolution semantics, which cannot be fully keyed); the key = format version + HEAD + porcelain `-z` digest + untracked map + changed-tracked map + ignored `.go` map + `go version` + the whole `go env` digest (excluding `GOGCCFLAGS`); a non-empty `GOWORK` other than `off` makes that pass abandon the cache.
- **Registration**: the directive went **v1.78 -> v1.79 -> v1.80 -> v1.81** (four amendments in one slice: landing / (6) fixes / `gofmt` exclusion / whole `go env` digest plus conservative `GOWORK` disable); three counterexamples, PASS-only, fail-open, the LRU cap and key-mismatch rejection are all in the clause.
- **Deliverables**: report pair `docs/validation/gates-base-cache{,_EN}.md` (strictly mirrored, with the `artifacts` fence and the (6) `sixstage` verbatim transcript); evidence pack `docs/validation/evidence/P96-GATES-BASE-CACHE/` (11 payloads plus MANIFEST/SHA256SUMS); the P96 ledger completion record; the `tools/gates/lib/baseCache.js` plus `gate.js` / `selftest.js` changes.
- **Review chain**: (1) one (`DeepSeek V4 Pro`; `DESIGN: DONE`) -> (2) one / (3) three (Flash; r1 F-1, r2 F-2, r3 `TEST: PASS`) -> (5) one (`DeepSeek V4 Pro`, `PRE-REVIEW: PASS` plus 4 Lows) -> (6) three (`MAI Code 2 Flash`, r1 / r2 `FINDINGS` -> r3 `INDEPENDENT SCAN: PASS`) -> (7) deep tier 1+3 (`GPT-6.2 Sol` Max, four rounds of `FINDINGS`) plus the Neon cross-check (4 HIGH) plus two narrow re-reviews (out of quota, OB-37) -> (8a) / (9) / (8b) / (8c) -> (10).
- **Metering**: **R = 12, D = 0, zero probes** ((3) 2 / (6) 2 / (7) 3 / convergence and closeout 4 / (9) first-round `G5(d)` 1; the High found in place before (4) is not counted separately per the P88 precedent).
- **Commits and CI**: main commit `81a6794` (17 files, +1506/-19) -> run **`46053254410`** first try, both legs success (both legs `SELFTEST: PASS 806/806`, `GATE REPORT scope=ci changed=17`, `TOTAL_FAIL=0`) -> backfill `90932ee` (8 files, adding the CI readings, the GLM archives and the temp-file disposal) -> run `46056837161` both legs success (`changed=8`), an observation record per §11.2.
- **tmp and leftovers**: the P96 temp directory was cleared in (8)c; `find-quotes.js` was precisely deleted after a whole-repo reference check and registered; `tmp` is back to two items. **Deviations**: the controller ran cold/hot and the counterexamples while (3) reviewed afterwards (endorsed by the user); independence is covered by (6)/(7) and (9).

### Appendix A: obligation hand-over (this version: **none**)

- **Historical hand-over** (first traced to P75 O-7): the inconsistent wording inside the `s1-exit.md:52` core-chain row was honoured by P76 change 1 plus the S-G2 correction (`0f80c1e`).
- **Hand-over in this version**: **none** - P96 leaves no obligation; the queue is empty (next item = none) and the tier-3 performance candidates are cleared, so no review or delivery obligation is handed over.
- **Convention**: in use since v59; a landing slice entry must carry the admission-condition face.

### Annex B: AT-23 stage 2 status (unchanged)

- **Status**: A1 ruled "not applicable" (user ruling on 2026-11-09; not a permanent ban).
- **Executable object before S2**: none (R-3 = (A)).
- **Restart path**: if the user later rules to restart stage 2, a real enforcement hookup (S2) or an alternative path must exist first.
- **Natural restart point**: S2 starts -> enforcement really hookup -> S-G8 unlocks -> R-3 re-evaluated.

---

## 3. Completed slices (recent baseline)

**v63 §3 already held sixteen P64-P80 rows; v64 appended P81 (seventeen); v65 appended P82 (eighteen); v66 kept the eighteen of v65 and appended P79 (nineteen); v67 kept the nineteen of v66 and appended P83 (twenty); v68 kept the twenty of v67 and appended P84 (twenty-one); v69 kept the twenty-one of v68 and appended P86 (twenty-two); v70 kept the twenty-two of v69 (readings copied, not rewritten) and appended P85 (twenty-three); v71 kept the twenty-three of v70 (readings copied, not rewritten) and appended P87 (twenty-four); v72 kept the twenty-four of v71 (readings copied, not rewritten) and appended P88 (twenty-five); v73 kept the twenty-five of v72 (readings copied, not rewritten) and appended P89 (twenty-six); v74 kept the twenty-six of v73 (readings copied, not rewritten) and appended P90 (twenty-seven); v75 kept the twenty-seven of v74 (readings copied, not rewritten) and appended P91 (twenty-eight); v77 kept the twenty-nine of v76 (readings copied, not rewritten) and appended P93; v78 keeps the thirty of v77 (readings copied, not rewritten) and appends P94; v79 keeps the thirty-one of v78 (readings copied, not rewritten) and appends P95; v80 keeps the thirty-two of v79 (readings copied, not rewritten) and appends P96. Thirty-three rows in total.**

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
| v1.76 | P93 DIRECTIVE-G6-4-TRANSCRIPT-EXEMPT | `23afa2d` + `04f2c0d` + `d4ae682` | **Closed** (governance slice; channel = b (a G6-4 scope sub-rule: a CLOSED `sixstage` fence whose five sections are in order is exempt), zero change to the decision logic, single-source reuse of the G7 internals, three fail-closed faces; directive bumped to **v1.76**; SELFTEST 744 -> 757 (T133-T143); R = 4, D = 0; main run `45782329377` green on both legs plus correction `04f2c0d` and backfill `d4ae682`) |
| v1.77 | P94 GATES-EXEMPT | `49f1609` + `afd9fb9` | **Closed** (governance; the three merged verdicts = the exemption extends to G2 + G4b / no generic registry slot / whitelist threshold text only; zero change to the decision logic; directive **v1.77**; SELFTEST 757 -> 770 (T144-T156); R = 6, D = 0; main run `45910223000` green first try + backfill `45910797198` observation) |
| v1.77 | P95 GATES-CRITERIA-ASYNC | `238e45c` + `bea554b` | **Closed** (assessment class; the (1) benefit gate = do not do it (ctx fully synchronous => benefit ≈0); P88 recommendation 4 registered as "assessed, benefit does not hold, not implemented"; zero code change, no directive bump; R = 1, D = 0; main run `45939437180` both legs green + backfill `45940170692` observation) |
| v1.81 | P96 GATES-BASE-CACHE | `81a6794` + `90932ee` | **Closed** (tooling slice; the section 6.1 base verdict cache clause; `gofmt` always runs for real and the cache covers build / vet / test only; directive v1.78 -> v1.81 in four amendments; SELFTEST 770 -> 806 (T157-T165); R = 12, D = 0; main run `46053254410` changed=17 + backfill run `46056837161` changed=8, both legs green) |

> **D caliber of section 3** (same as v61): "D = n" here abbreviates the temp-artifact discipline violation count (section 5.1); P89 totals **D = 1** (recorded in the report pair section 6 and the (10) receipt). **The D count and the section 9.2 R count are two independent series**.

---

## 4. Execution queue

### Tier 0 pre-work (✅ discharged; nothing new here)

- P5 / P6 were closed by P57; tier 0 is empty with no residue.

### Tier 1 main-line hard gates

- Items 8-12: all ✅ discharged (see v57 / v58 section 4; trace back through that reference chain if needed).
- **`T027-B` (defensive-refactor phase) = done (partly not applicable)**: C1-C7 / N1-N4 / R-1 / R-2 / R-3 are all handled (R-3 = (A)). **The main line is fully discharged**.
- **`S1 Exit closeout` (gate level done; security remediation F1-F8 fully settled)**: P76 produced the closeout roadmap; P77 closed S-G4; P78 closed S-G3; P80 judged S-G6 to hold at the gate level; the independent security review was completed by P79 and P83 / P84 / P86 / P85 delivered the F1-F8 fixes. Gate table below (P87 / P88 / P89 / P96 left it unchanged):

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

- Tiers 1 and 2 are both cleared and the tier-3 performance candidates are cleared by this version; this version advances nothing.
- Ways back into the queue: the only entry is real S2 enforcement landing, which unlocks the S2 pre-re-evaluation (S-G1 / S-G8); or the user names a new slice.

- After that = the pre-S2 re-evaluation (S-G1 / S-G8) (the next T027 main-line stop, pending user scheduling):
  - Prerequisite = the S2 enforcement hookup; the independent security review prerequisite was discharged by P79 and P83 / P84 / P86 / P85 delivered the F1-F8 fixes, so the ADR-047 rollback clause no longer blocks;
  - `S-G8` stays blocked until the S2 enforcement hookup (P73 A2 / P75 A2 / ADR-044 / ADR-045); the `s1-exit{,_EN}.md` gate table is untouched.
- **Post-S2 re-evaluation items**: **S-G5 step two (overall S1 Exit approval) plus S-G7**, in the order of P80 report section 7.

### Tier 2 routine governance

- P72 / P73 / P74 / P75 / P76 / P77 / P78 / P80 / P81 / P82 / P79 / P83 / P84 / P86 / P85 / P87 / P88 / P89 / P90 / P91 / P92 / P93 / P94 / P95 / **P96** ✅ left the queue (tier 2 is **fully cleared**).
- Each slice's full section 1.2 field set lives in `docs/t027/REMAINING_SLICES{,_EN}.md`.

### Tier 3 low-priority candidates

| Field | Content |
|---|---|
| Identifier | `OB-50` / `spawnGit` / `ci frozen cost` / `22 historical %TEMP% items` (the performance candidates are cleared) |
| Current state (each verified) | OB-50 = **disposed** (§12.4, v1.41); spawnGit = the git helper of `tools/gates/selftest.js` (the `lib/ctx.js` spawn wrapper keeps stderr in the result object and never prints it), a candidate OB (diagnostic) from P27 report §5 item 5 with no §12.4 row; the ci frozen cost = disposed by P90 (the `--skip-base` Gates step plus go caching judged sufficient); %TEMP% leftovers = OB-98 disposed / OB-99 and OB-100 watching; `OB-145` disposed by P91 (v1.74); `GATES-SELFTEST-SPLIT` left the tier when P92 closed (v1.75); `GATES-EXEMPT` left the tier when P94 closed (v1.77); `GATES-CRITERIA-ASYNC` left the tier when P95 closed (assessed and not implemented); `GATES-BASE-CACHE` left the tier when P96 closed (v1.81) |
| Contents | See the matching section 12.4 rows; `GATES-PERF-ASSESS` (P88, v1.73) and `GATES-CI-PERF` (P90) both closed and left the tier; the performance candidates are in the table below |
| Dependency | the long suspension covers only the undisposed items (spawnGit / the `%TEMP%` timing forms OB-99 and OB-100), none with a defined trigger |
| Size | Varies (S through M-L) |
| (7) quota | none |
| Boundaries | This version advances no suspended item and keeps the registrations; the performance candidates are cleared, so no slice can be opened from them |
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
- **DR registry**: DR max = **18** (none of P85 / P86 / P87 / P88 / P89 / P90 / P91 / P92 / P93 / P94 / P95 / P96 added a DR).
- **P87 Hydrogen prerequisite gap (record-only)**: the lightweight tier (`Gemini 4 Hydrogen`) clause has landed but the comparative review has not been measured; it points at a later business slice (§1.7.4 temporary authorisation + OB record), and it **must not be described as verified**; once covered, annotate `docs/validation/directive-review-model-policy{,_EN}.md`. Registry slot = report §7 and the P87 ledger section.
- **Tier-3 performance candidates (rolled over here)**: `GATES-BASE-CACHE` left the tier when P96 closed (v1.81); `GATES-CRITERIA-ASYNC` left the tier when P95 closed (assessed and not implemented); `GATES-SELFTEST-SPLIT` left the tier when P92 closed; `GATES-CI-PERF` left it when P90 closed; `GATES-EXEMPT` left it when P94 closed.
- **`find-quotes.js` disposition slot**: precisely deleted during P96 (8)c (whole-repo references survive only as historical snapshot traces, with no code dependency) and registered in the P96 ledger completion record; this version no longer keeps the slot.
- Mechanical readings: the largest OB in §12.4 = **145**; ADR max = **48**; DR max = **18**; trial deferral max = the 25th time.
- **Removed here**: the `GATES-BASE-CACHE` tier-3 slot (P96 closed it) and the `find-quotes.js` disposition slot (cleared). **Added here**: none (P96 added no watch item).

## §6 Pending user rulings

- **This table is empty** (the P96 stop-1 ruling, the (7) deep tier 1+3, the two convergence rounds and the backfill authorisation were all ruled on 2026-11-27 / 28).
- **Expected next ruling**: **none** - waiting for the user to schedule S2 / enforcement or to name a new slice.

## 7 / 8 / 9 Pointer sections

Same structure as v60: section 7 start-up and step 0 general rules -> directive section 13.1; section 8 cross-slice collaboration discipline -> section 13.2; section 9 checklist meta-rules and gate pits -> section 13.3 / 13.4. **The rule bodies live in directive section 13 and are not copied here**.

---

## §10 One-line opening for a new session

> Below is my project state snapshot (v80). Check the 7 anchors in section 1 first; if they match, continue with the section 4 queue, otherwise the repository wins.

> [paste this checklist v80 in full]

> [supporting documents] `REMAINING_SLICES{,_EN}.md` (P96 closed with a completion record; the queue is empty), `DELIVERY_DIRECTIVE.md` (**v1.81**), `gates-base-cache{,_EN}.md` (the P96 report pair), `docs/validation/evidence/P96-GATES-BASE-CACHE/` (the P96 evidence pack), `ADR_REGISTER{,_EN}.md` (ADR-048 stands), `DOCUMENTATION_PLAN{,_EN}.md` (the hand-off pointer now targets v80), `CHECKLIST_v79{,_EN}.md` (frozen, replaced by v80).

> Now: P96 closed (tooling class: the base verdict cache covers build / vet / test only and `gofmt` always runs for real, directive v1.78 -> v1.81; both legs SELFTEST 806/806, `changed=17` / `changed=8`, `TOTAL_FAIL=0`; R = 12, D = 0); next = none (queue empty); the S2 pre-re-evaluation (S-G1 / S-G8) unlocks once real S2 enforcement lands.

## §11 One-line current status

P96 **closed** (main commit `81a6794` plus backfill `90932ee` pushed to origin; the main run `46053254410` went green first try and the backfill run `46056837161` is observation-only per §11.2; both legs `SELFTEST 806/806`, `changed=17` / `changed=8`, `TOTAL_FAIL=0`; R = 12, D = 0; the cache covers build / vet / test only and `gofmt` always runs for real, **directive v1.81**; `tmp` back to two items); next = **none (queue empty)**; the S2 pre-re-evaluation (S-G1 / S-G8) unlocks once real S2 enforcement lands.
