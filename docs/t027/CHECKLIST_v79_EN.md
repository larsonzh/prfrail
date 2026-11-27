# Task checklist (cross-session snapshot v79 - after P95 closes - queue empty)

Generated: 2026-11-27. HEAD = `bea554b`; directive = v1.77 (P95 did not bump it); P95 CI main round: `45939437180` (both legs success); the backfill round `45940170692` is an observation record only per §11.2 (API-level verdict).

**Replacement**: this version replaces v78 (`docs/t027/CHECKLIST_v78{,_EN}.md` frozen, not rewritten).

**Structure**: same as v78 (§1 anchors -> §2 fact pack -> §3 completed slices -> §4 execution queue -> §5 watch items -> §6 pending rulings -> §7 / §8 / §9 pointer sections -> §10 opening -> §11 one-line status); the rule bodies of §7-§9 live in directive §13, so this snapshot does not copy them.

**Changes versus v78**:

- (1) §1 anchors rolled over: HEAD `bea554b`; directive **v1.77** (P95 did not bump it); SELFTEST **770/770**; the rolling window drops v76 and takes v77 / v78 / v79, with v77 / v78 re-checked here (217 / 211) and v79 itself left to v80.
- (2) §2 becomes the P95 fact pack (`CLOSED`, assessment class; main commit and backfill, both-leg CI readings, **R = 1**, D = 0, zero probes; the (1) benefit gate = do not implement; P88 report section 7 recommendation 4 registered as "assessed, benefit does not hold, not implemented"; zero criteria code change, directive stays v1.77).
- (3) §3 appends one P95 row (the thirty-one rows of v78 plus one = thirty-two).
- (4) §4 queue update: P95 leaves; **next item = none (queue empty, queued work paused)**; "after that" stays the S2 pre-re-evaluation (S-G1 / S-G8, pending S2 enforcement); tier 3 drops `GATES-CRITERIA-ASYNC` (assessed and not implemented), keeping `GATES-BASE-CACHE` as the only performance candidate plus the other suspended items.
- (5) §5 watch items: the `GATES-CRITERIA-ASYNC` candidate is removed; P95 adds no new watch item.
- (6) §6 pending rulings: **none** (waiting for the user to schedule S2 / enforcement or to name `GATES-BASE-CACHE`).
- (7) Appendix A: the obligation hand-over is **none**; appendix B is copied unchanged; §10 / §11 rolled over, and the `DOCUMENTATION_PLAN{,_EN}` pointer now targets v79.
- (8) Mechanical readings re-checked: ADR 48 / OB 145 / DR 18 / trial deferral 25 (P95 added none, each checked).
---

## §1 Current anchors (check these 7 numbers first)

**Rolling-window caliber = the latest 3 versions (v77 / v78 / v79)**; v76 has been squeezed out (same rule as v52).

| Item | Value | Check command |
|---|---|---|
| HEAD | `bea554b` | `git rev-parse --short HEAD` |
| origin/main | `bea554b` | `git rev-parse --short origin/main` |
| Directive version | **v1.77** | the header of `docs/DELIVERY_DIRECTIVE.md` |
| SELFTEST | **770/770** | `node tools/gates/gate.js --selftest` |
| Three scopes | **all TOTAL_FAIL=0** | `gate.js --all --scope=index` / `tree` / `ci` |
| v77 whole-file lines | **217** (LF count, CN / EN symmetric; **re-checked here**) | `docs/t027/CHECKLIST_v77{,_EN}.md` by LF count |
| v78 whole-file lines | **211** (LF count, CN / EN symmetric; **re-checked here**) | `docs/t027/CHECKLIST_v78{,_EN}.md` by LF count |

> The HEAD = `bea554b` here is the generation-time anchor after the P95 backfill; once v79 lands HEAD moves on, so `git rev-parse --short HEAD` is authoritative then. v79 itself is left to v80 for its line count.

> Mechanical readings: ADR max = **48**; OB max = **145**; DR max = **18**; trial deferral max = the 25th time. P95 checked each: no new ADR, no new OB, no new DR, the deferral count unchanged.

> Rolling window = the latest 3 versions (v77 / v78 / v79); v76 has been squeezed out. If 2 or more items disagree, the repository wins: report the final readings first and then continue.

`tmp` fixed assertion: top level = `.gitkeep` + `gates-cache` (only `g4b` beneath) + `find-quotes.js` (no reference residue; awaiting the tmp lifecycle slice). **The P95 temp directory was cleared in (8c)**, so this assertion does not include it.

## §2 P95 fact pack

- **Slice / size**: `GATES-CRITERIA-ASYNC` (assessing the G6 / G7 / G8 asynchronisation; tooling `[SLICE]`) / **M**; controller `DeepSeek V4 Pro` Max (user ruling 2026-11-27, not a precedent); the user named it on 2026-11-27 and in the same turn ruled for candidate A with commit-and-push authorisation.
- **Goal**: (1) a benefit gate - measure the serial time composition of G6 / G7 / G8 inside a single-process `--all` run and the concurrency benefit ceiling, and assess the cost of shared-cache concurrent reads and of keeping the report-order contract; (2) if approved, switch G6 / G7 / G8 to true `runCriteriaConcurrent` concurrency; (3) keep the criteria verdicts byte-identical and the report order unchanged.
- **Verdict**: **`CLOSED`** (2026-11-27, assessment class); the (1) benefit gate = **do not do it** - the ctx read APIs are fully synchronous (`readFileSync` / `spawnSync`) and a single Node process cannot overlap synchronous I/O through asynchronisation, so the theoretical benefit ceiling is ≈0 and the practical bound is ≤250 ms out of a 5.8 s wall clock; the user ruled for candidate A at the (1) stop (assessment closure, **not a precedent**).
- **Registration**: P88 report section 7 recommendation 4 registered as "**assessed, benefit does not hold, not implemented**"; zero criteria code change, no directive bump (v1.77 unchanged); the `CHECKLIST_v78` pair is untouched (absorbed by this v79).
- **Deliverables**: report pair `docs/validation/gates-criteria-async{,_EN}.md` (strictly mirrored, with the `artifacts` fence and the (6) `sixstage` verbatim transcript); evidence pack `docs/validation/evidence/P95-GATES-CRITERIA-ASYNC/` (10 payloads plus MANIFEST/SHA256SUMS); the P95 ledger completion record.
- **Review chain**: (1) one (`DeepSeek V4 Pro` Max) -> (2) one / (3) one (Flash; (3) `CONSISTENCY: PASS`) -> (5) one (`PRE-REVIEW: PASS`, 2 Lows fixed + 3 Info) -> (6) one (`MAI Code 2 Flash`, `INDEPENDENT SCAN: PASS`) -> (7) 1+1 (r1 `FINDINGS` 1 Medium + 1 Low -> fixed -> r2 `RE-REVIEW: PASS`; the S-caliber quota is a user ruling and is fully used) -> (8a) / (9) / (8b) / (8c) -> (10); batch helper `GLM-5.3-Flash` twice (`MIRROR: PASS` / `NUMBERS: PASS`).
- **Metering**: **R = 1, D = 0, zero probes** (R = the (7) r1 fix plus the r2 re-review round; the two (5) Lows inside a PASS verdict and the (4) gate in-place fixes are not counted per the P88 precedent, each registered in report section 7).
- **Commits and CI**: main commit `238e45c` (16 files, +884/-0) -> run **`45939437180`** first try, both legs success (both legs `SELFTEST: PASS 770/770`, `GATE REPORT scope=ci changed=16`, `TOTAL_FAIL=0`; package-level ok ubuntu 21 / windows 15, FAIL 0) -> backfill `bea554b` (4 files, +18/-4, adding the CI readings) -> run `45940170692` both legs success, an observation record per §11.2 (API-level verdict).
- **tmp and leftovers**: the P95 temp directory was cleared in (8c) (`tmp` back to three items); `find-quotes.js` still awaits the tmp lifecycle slice. **Deviations**: none (the (10) dispatch table is annotated "not mechanically verified" because the debug log carries no transcript record; registered as-is and accepted by the user).

### Appendix A: obligation hand-over (this version: **none**)

- **Historical hand-over** (first traced to P75 O-7): the inconsistent wording inside the `s1-exit.md:52` core-chain row was honoured by P76 change 1 plus the S-G2 correction (`0f80c1e`).
- **Hand-over in this version**: **none** - P95 leaves no obligation; the queue is empty (next item = none) and `GATES-BASE-CACHE` is record-only, so it creates no review or delivery obligation.
- **Convention**: in use since v59; a landing slice entry must carry the admission-condition face.

### Annex B: AT-23 stage 2 status (unchanged)

- **Status**: A1 ruled "not applicable" (user ruling on 2026-11-09; not a permanent ban).
- **Executable object before S2**: none (R-3 = (A)).
- **Restart path**: if the user later rules to restart stage 2, a real enforcement hookup (S2) or an alternative path must exist first.
- **Natural restart point**: S2 starts -> enforcement really hookup -> S-G8 unlocks -> R-3 re-evaluated.

---

## 3. Completed slices (recent baseline)

**v63 §3 already held sixteen P64-P80 rows; v64 appended P81 (seventeen); v65 appended P82 (eighteen); v66 kept the eighteen of v65 and appended P79 (nineteen); v67 kept the nineteen of v66 and appended P83 (twenty); v68 kept the twenty of v67 and appended P84 (twenty-one); v69 kept the twenty-one of v68 and appended P86 (twenty-two); v70 kept the twenty-two of v69 (readings copied, not rewritten) and appended P85 (twenty-three); v71 kept the twenty-three of v70 (readings copied, not rewritten) and appended P87 (twenty-four); v72 kept the twenty-four of v71 (readings copied, not rewritten) and appended P88 (twenty-five); v73 kept the twenty-five of v72 (readings copied, not rewritten) and appended P89 (twenty-six); v74 kept the twenty-six of v73 (readings copied, not rewritten) and appended P90 (twenty-seven); v75 kept the twenty-seven of v74 (readings copied, not rewritten) and appended P91 (twenty-eight); v77 kept the twenty-nine of v76 (readings copied, not rewritten) and appended P93; v78 keeps the thirty of v77 (readings copied, not rewritten) and appends P94; v79 keeps the thirty-one of v78 (readings copied, not rewritten) and appends P95. Thirty-two rows in total.**

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
| v1.76 | P93 DIRECTIVE-G6-4-TRANSCRIPT-EXEMPT | `23afa2d` + `04f2c0d` + `d4ae682` | **Closed** (governance slice; channel = b (a G6-4 scope sub-rule: a CLOSED fence whose five sections are in order is exempt), zero change to the decision logic, single-source reuse of the G7 internals, three fail-closed faces; directive bumped to **v1.76**; SELFTEST 744 -> 757 (T133-T143); R = 4, D = 0; main run `45782329377` green on both legs plus correction `04f2c0d` and backfill `d4ae682`) |
| v1.77 | P94 GATES-EXEMPT | `49f1609` + `afd9fb9` | **Closed** (governance; the three merged verdicts = the exemption extends to G2 + G4b / no generic registry slot / whitelist threshold text only; zero change to the decision logic; directive **v1.77**; SELFTEST 757 -> 770 (T144-T156); R = 6, D = 0; main run `45910223000` green first try + backfill `45910797198` observation) |
| v1.77 | P95 GATES-CRITERIA-ASYNC | `238e45c` + `bea554b` | **Closed** (assessment class; the (1) benefit gate = do not do it (ctx fully synchronous => benefit ≈0); P88 recommendation 4 registered as "assessed, benefit does not hold, not implemented"; zero code change, no directive bump; R = 1, D = 0; main run `45939437180` both legs green + backfill `45940170692` observation) |

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

**Next item = none (queue empty, queued work paused)**

- Tiers 1 and 2 are both cleared and tier 3 holds record-only candidates; this version advances nothing.
- Two ways back into the queue: (1) real S2 enforcement lands, which unlocks the S2 pre-re-evaluation (S-G1 / S-G8); (2) the user names a tier-3 performance candidate (now only `GATES-BASE-CACHE`).

- After that = the pre-S2 re-evaluation (S-G1 / S-G8) (the next T027 main-line stop, pending user scheduling):
  - Prerequisite = the S2 enforcement hookup; the independent security review prerequisite was discharged by P79 and P83 / P84 / P86 / P85 delivered the F1-F8 fixes, so the ADR-047 rollback clause no longer blocks;
  - `S-G8` stays blocked until the S2 enforcement hookup (P73 A2 / P75 A2 / ADR-044 / ADR-045); the `s1-exit{,_EN}.md` gate table is untouched.
- **Post-S2 re-evaluation items**: **S-G5 step two (overall S1 Exit approval) plus S-G7**, in the order of P80 report section 7.

### Tier 2 routine governance

- P72 / P73 / P74 / P75 / P76 / P77 / P78 / P80 / P81 / P82 / P79 / P83 / P84 / P86 / P85 / P87 / P88 / P89 / P90 / P91 / P92 / P93 / P94 / **P95** ✅ left the queue (tier 2 is **fully cleared**).
- Each slice's full section 1.2 field set lives in `docs/t027/REMAINING_SLICES{,_EN}.md`.

### Tier 3 low-priority candidates

| Field | Content |
|---|---|
| Identifier | `OB-50` / `spawnGit` / `ci frozen cost` / `22 historical %TEMP% items` + the performance candidate (`GATES-BASE-CACHE`) |
| Current state (each verified) | OB-50 = **disposed** (§12.4, v1.41); spawnGit = the git helper of `tools/gates/selftest.js` (the `lib/ctx.js` spawn wrapper keeps stderr in the result object and never prints it), a candidate OB (diagnostic) from P27 report §5 item 5 with no §12.4 row; the ci frozen cost = disposed by P90 (the `--skip-base` `Gates` step plus go caching judged sufficient); %TEMP% leftovers = OB-98 disposed / OB-99 and OB-100 watching; `OB-145` disposed by P91 (v1.74); `GATES-SELFTEST-SPLIT` left the tier when P92 closed (v1.75); `GATES-EXEMPT` left the tier when P94 closed (v1.77); `GATES-CRITERIA-ASYNC` left the tier when P95 closed (assessed and not implemented) |
| Contents | See the matching section 12.4 rows; `GATES-PERF-ASSESS` (P88, v1.73) and `GATES-CI-PERF` (P90) both closed and left the tier; the performance candidate is in the table below |
| Dependency | the long suspension covers only the undisposed items (spawnGit / the `%TEMP%` timing forms OB-99 and OB-100), none with a defined trigger |
| Size | Varies (S through M-L) |
| (7) quota | none |
| Boundaries | This version advances no suspended item and keeps the registrations; the single remaining performance candidate is record-only and constitutes no slice |
| Start condition | the user must name it explicitly |

**Performance candidates (risk ascending; the single row below is still record-only and needs the user to name it)**:

| Identifier | One-line scope | Relation to existing items / preconditions |
|---|---|---|
| `GATES-BASE-CACHE` | cache the base verdicts (reuse the `gofmt` / `build` / `vet` / `test` conclusions across scopes) | same family as section 6.3 G4b "never cache verdicts", so **a directive clause must exist first**; highest risk, listed last |

**Queue discipline**: tier boundaries and numbering follow v60; this version does not bump the directive.

---

## §5 Watch items

- Watching (**13**): `OB-1 / 7 / 10 / 33 / 87 / 99 / 100 / 104 / 105 / 107 / 111 / 113 / 135` (carried from v61).
- Under discussion (9): `OB-19 / 21 / 52 / 53 / 57 / 58 / 59 / 61 / 75`.
- To do (1): `OB-69`. Non-terminal total **23** (13 + 9 + 1), matching this version re-check.
- **DR registry**: DR max = **18** (none of P85 / P86 / P87 / P88 / P89 / P90 / P91 / P92 / P93 / P94 / P95 added a DR).
- **P87 Hydrogen prerequisite gap (record-only)**: the lightweight tier (`Gemini 4 Hydrogen`) clause has landed but the comparative review has not been measured; it points at a later business slice (§1.7.4 temporary authorisation + OB record), and it **must not be described as verified**; once covered, annotate `docs/validation/directive-review-model-policy{,_EN}.md`. Registry slot = report §7 and the P87 ledger section.
- **Tier-3 performance candidates (rolled over here)**: `GATES-BASE-CACHE` stays record-only (the user must name it); `GATES-CRITERIA-ASYNC` left the tier when P95 closed (assessed and not implemented); `GATES-SELFTEST-SPLIT` left the tier when P92 closed; `GATES-CI-PERF` left it when P90 closed; `GATES-EXEMPT` left it when P94 closed.
- **`find-quotes.js` disposition slot**: no reference residue (P89 re-checked and excluded it from archiving and rewriting); it still awaits the tmp lifecycle slice.
- Mechanical readings: the largest OB in §12.4 = **145**; ADR max = **48**; DR max = **18**; trial deferral max = the 25th time.
- **Removed here**: the `GATES-CRITERIA-ASYNC` tier-3 slot (P95 closed it by assessment). **Added here**: none (P95 added no watch item).

## §6 Pending user rulings

- **This table is empty** (the P94 authorisation sequence, the R/D caliber, (7) 1+2 and the v78 snapshot generation were all ruled on 2026-11-27).
- **Expected next ruling**: **none** - waiting for the user to schedule S2 / enforcement or to name the tier-3 performance candidate (`GATES-BASE-CACHE`).

## 7 / 8 / 9 Pointer sections

Same structure as v60: section 7 start-up and step 0 general rules -> directive section 13.1; section 8 cross-slice collaboration discipline -> section 13.2; section 9 checklist meta-rules and gate pits -> section 13.3 / 13.4. **The rule bodies live in directive section 13 and are not copied here**.

---

## §10 One-line opening for a new session

> Below is my project state snapshot (v79). Check the 7 anchors in section 1 first; if they match, continue with the section 4 queue, otherwise the repository wins.

> [paste this checklist v79 in full]

> [supporting documents] `REMAINING_SLICES{,_EN}.md` (P95 closed with a completion record; the queue is empty), `DELIVERY_DIRECTIVE.md` (**v1.77**), `gates-criteria-async{,_EN}.md` (the P95 report pair), `docs/validation/evidence/P95-GATES-CRITERIA-ASYNC/` (the P95 evidence pack), `ADR_REGISTER{,_EN}.md` (ADR-048 stands), `DOCUMENTATION_PLAN{,_EN}.md` (the hand-off pointer now targets v79), `CHECKLIST_v78{,_EN}.md` (frozen, replaced by v79).

> Now: P95 closed (assessment class: G6 / G7 / G8 asynchronisation not done, P88 section 7 recommendation 4 registered as not implemented, directive stays v1.77; both legs SELFTEST 770/770, `changed=16`, `TOTAL_FAIL=0`; R = 1, D = 0); next = none (queue empty); the S2 pre-re-evaluation (S-G1 / S-G8) unlocks once real S2 enforcement lands.

## §11 One-line current status

P95 **closed** (main commit `238e45c` plus backfill `bea554b` pushed to origin; the main run `45939437180` went green first try and the backfill run `45940170692` is observation-only per §11.2; both legs `SELFTEST 770/770`, `changed=16`, `TOTAL_FAIL=0`; R = 1, D = 0; the (1) benefit gate = do not do it, so zero criteria code change and **the directive stays at v1.77**; P88 section 7 recommendation 4 registered as "assessed, benefit does not hold, not implemented"; `tmp` back to three items); next = **none (queue empty)**; the S2 pre-re-evaluation (S-G1 / S-G8) unlocks once real S2 enforcement lands.
