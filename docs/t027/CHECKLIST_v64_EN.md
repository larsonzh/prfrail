# Task Execution Checklist (cross-session snapshot v64 · after P81 closes · stage-6/stage-7 model replacement done · next = the P82 temp-file discipline fix)

Generated: 2026-11-19. HEAD = `9ca6576`; directive = v1.70; three CI runs: the P81 main commit `44281812336` (commit `45590e8`, both legs success), the P81 §11.2 backfill `44284200392` (commit `fd4d8fb`, observation run) and the probe / OB-139 backfill `44392113028` (commit `9ca6576`, both legs success).

**Supersession**: this version supersedes v63 (`docs/t027/CHECKLIST_v63{,_EN}.md` is frozen and is not retro-edited).

**Structure**: §1 anchors -> §2 fact package -> §3 completed slices -> §4 execution queue -> §5 observation items -> §6 pending rulings -> §7 / §8 / §9 pointer sections -> §10 opening -> §11 one-line status; the rule bodies behind §7-§9 have moved into directive §13, and this snapshot does not copy them.

**Changes in this version (relative to v63)**:
- (1) §1 anchor rollover: HEAD `9ca6576`; directive **v1.70**; SELFTEST **740/740**; the rolling window drops v61 and admits v64 (v62 / v63 / v64).
- (2) §2 becomes the P81 fact pack (the stage-6/stage-7 replacement and the three-way closed set with 8 members, the self-corrected probe-target deviation, the (6) format-deviation re-send, R = 1 / D = 2, the ten v1.70 landing points, and the two wording-side G6-4 / G6-6 self-corrections); appendix A records the P82 hand-off of the OB-139 correction channel, appendix B is unchanged.
- (3) §3 appends one P81 row (the sixteen rows of v63 plus one = seventeen rows).
- (4) §4 queue update: P81 leaves the queue; next = P82 (`TMP-DISCIPLINE-OUT-OF-REPO`, fixes §5.1, bumps to v1.71, start requires the user to name it); then P79; then the post-S2 re-evaluations.
- (5) §5: OB-139 has moved to pending discussion (13 observing, 10 pending, 24 non-terminal), with the P82 registration and the 2026-11-19 ruling recorded (no waiver, not a precedent).
- (6) §6 is cleared; the next ruling item is expected to be the start order of P82 versus P79 (the current advice is P82 first).
- (7) §10 / §11 roll over; the cross-session handoff pointer in `DOCUMENTATION_PLAN{,_EN}.md` now points at v64.

---

## §1 Current anchors (check these eight numbers first)

**Rolling window = the last three versions (v62 / v63 / v64)**; v61 has been dropped (same rule that dropped v52).

| Item | Value | Check command |
|---|---|---|
| HEAD | `9ca6576` | `git rev-parse --short HEAD` |
| origin/main | `9ca6576` | `git rev-parse --short origin/main` |
| Directive version | **v1.70** | the header of `docs/DELIVERY_DIRECTIVE.md` |
| SELFTEST | **740/740** | `node tools/gates/gate.js --selftest` |
| Three scopes | **all TOTAL_FAIL=0** | `gate.js --all --scope=index` / `tree` / `ci` |
| v62 total lines | **217** (LF count, CN / EN symmetric) | `docs/t027/CHECKLIST_v62{,_EN}.md` by LF |
| v63 total lines | **222** (LF count, CN / EN symmetric; re-checked by v64) | `docs/t027/CHECKLIST_v63{,_EN}.md` by LF |
| v64 total lines | **224** (LF count, CN / EN symmetric; re-checked by v65) | `docs/t027/CHECKLIST_v64{,_EN}.md` by LF |

> This snapshot's HEAD = `9ca6576` is the checkpoint after the P81 probe / OB-139 backfill commit; once v64 lands, HEAD advances, so use `git rev-parse --short HEAD` as the live value.

> Rolling window = the most recent 3 versions (v62 / v63 / v64); v61 has dropped out (the same rule that dropped v52); this snapshot's own line count is 224 (LF count), re-checked by v65. If two or more items disagree, trust the repository and report the terminal reading before continuing.

`tmp` fixed assertion: top level = `.gitkeep` + `gates-cache`; `gates-cache` holds only `g4b`; `tmp/P81-DIRECTIVE-MODEL-UPGRADE/` left no residue after the P81 ⑧c cleanup and the two cleanups of the 2026-11-19 back-fill.

---

## §2 P81 fact pack

- **Slice / size**: `DIRECTIVE-MODEL-UPGRADE` (P81; direct replacement of the stage-6 carrier and the stage-7 deep-tier model plus landing the P80 ruling (1) trace; governance `[SLICE]`) / **M**; controller `DeepSeek V4.1 Flash (deepseek)`; the user's 2026-11-18 kickoff message authorised the slice in the same turn.
- **Goal**: replace the stage-6 carrier and the stage-7 deep-tier carrier **directly** per the user ruling of 2026-11-18 (`MAI-Code-1.1-Flash (copilot)` -> `MAI Code 2 Flash (copilot)`; `GPT-6.1 Sol (copilot)` -> `GPT-6.2 Sol (copilot)`), **without adding Astra**; bump the directive to **v1.70**; and land the §7.8.3 trace obligation of P80 ruling (1).
- **Verdict**: **closed**; all **ten landing points** of v1.70 are replaced (the two §2.2 triple rows and the "⑦ three-tier selection" sentence / the §5.2 deep-tier row / the §6.3 G5-b closed set / the §7.5 Sol cost row / the §8.2 tier mapping and tier-authorisation sentence / the Appendix B.6 carrier call string / the v1.70 log row in §0.3 / the P80 authorisation-exception block at the end of §12.4 / Appendix C.0aw); the G5-b closed set is **identical three ways** (the §2.2 modelId column / the §6.3 G5-b closed set / `CLOSED_SET` in `g5b.js`) = **8 members**; SELFTEST rises from 737 to **740**.
- **Criterion boundaries**: the Sol prices **stay at 200 / 1000** (basis = the operator-UI read confirming both generations are level); the (7) effort ruling is **(a)** (§8.2 keeps "Extra High" and allows **Max**); the G5-b **wording and comparison rules stay verbatim**, only two members change; the p57 / `s1-exit` gate tables and the v63 snapshot are untouched.
- **Probes (target deviation self-corrected)**: the first three calls actually targeted three **already-enrolled** models (`DeepSeek V4 Pro` / `MAI-Code-1.1-Flash` / `GPT-5.3-Codex`), which did not match the intended probes of the new stage 6 / new stage 7 and was not disclosed in the (10) receipt, so the user found it on review and ruled on 2026-11-19 to add **two probes** (counted separately under §7.6): one each for the new stage 6 and the new stage 7, **both `PROBE-OK`**, with call strings matching the §2.2 qualified names verbatim; payloads `03-probes/probe-new-mai-code-2-flash.md` and `03-probes/probe-new-gpt-6.2-sol.md`; five probes in total.
- **Review chain**: (1) design `DESIGN: DONE` (deep-reasoner + DeepSeek V4 Pro) -> batch assistance `GLM-5.3-Flash` x 1 -> (2) implementation -> (3) consistency (C-1 to C-10 plus four mutations) -> (4) gates -> (5) r1 `PRE-REVIEW: PASS WITH FIXES` (4 Low, remediated) -> (6) **required**: the r1 format deviated (no Section A-E, no verdict line), so it was re-sent once per §7.2.2 and r2 returned `INDEPENDENT SCAN: PASS` (budget 3, 2 used) -> (7) final review `FINAL REVIEW: PASS` (0 findings) -> ⑧a / ⑨ / ⑧b / ⑧c -> ⑩ -> post-(10) back-fill (extra probes + R correction + OB-139 ruling).
- **Budget and blind review**: the (7) governance budget is **1 + 3** and only **1 was used** (the final review found nothing, so no re-review round was triggered); **no blind review** (none of the three mandatory semantics is touched and this is the first slice of the proportional spot-check window).
- **Metering**: **R = 1** (the probe-target deviation counts as corrective rework after discovery, recorded as-is and not counted as D), **D = 2** (event 1 = the `tmp/`-root misplacement, self-heal grade holds; event 2 = "`%TEMP%` out-of-repository location plus same-command create-delete", which has no self-heal channel under §5.1).
- **Gates and wording-side self-corrections**: during the back-fill, `G6-4` read a phrase whose version-number tail is immediately followed by the CJK character meaning "line" as a numeric line literal, and the `G6-6` EN status closed set **accepts only `To be discussed`** (not `Pending discussion`), so both were fixed on the **wording side** with the criterion bodies and closed-set members untouched; after the back-fill the three scopes report `TOTAL_FAIL=0` and SELFTEST **740/740**.
- **Evidence pack**: `docs/validation/evidence/P81-DIRECTIVE-MODEL-UPGRADE/` (`MANIFEST.md` plus `SHA256SUMS.txt` plus **16 payloads**, including the two `03-probes/` payloads and the four `06-gates/` gate transcripts; the three entry sets agree and re-computing every sha256 yields zero mismatches).
- **CI**: main commit `45590e8` -> run `44281812336` (both legs success; both legs print `SELFTEST 740/740`, `GATE REPORT scope=ci changed=28`, `TOTAL_FAIL=0`, with package-level ok ubuntu 21 / windows 15 and `FAIL` 0); the §11.2 back-fill `fd4d8fb` -> run `44284200392` (observation run, both legs success); the probe / OB-139 back-fill `9ca6576` -> run `44392113028` (both legs success; `changed=12`, `TOTAL_FAIL=0`, 36 package-level ok lines in total and `FAIL` 0).
- **Mechanical readings**: `SELFTEST PASS 740/740`; max ADR = 47; max OB = **139** (no new OB); max DR = 14; deferral = round 25.
- **Ledger convention** (§11.2): the P81 completion record carries only the report path plus R/D (no hash / run); the CI fact chain lives in section 2 of this snapshot and in the (10) receipt.
- **User rulings (2026-11-19; made directly by the user, not usable as a precedent)**:
  1. **Extra probes and R metering**: the two extra probes (new stage 6 / new stage 7) are corrective rework after a discovered deviation, so **R = 1** and no extra D; had either failed, the slice would have stopped and escalated per R2.2 without guessing another call string.
  2. **D acceptance**: both temp-file events of P81 are **accepted as D = 2** and kept as recorded.
  3. **OB-139**: the **sixth trigger is confirmed**, so the entry moves from "under observation" to "to be discussed" and slice **P82** is registered (fix the §5.1 temp-file discipline, bump to v1.71, schedule pending, start requires the user to name it); the ruling grants no waiver and may not be cited as a precedent.

### Appendix A: "obligation transfer" (this version: P82 takes over the OB-139 correction channel)

- **Historical hand-off** (first traced back to P75 O-7): the in-cell caliber mismatch in the core task chain row at `s1-exit.md:52` was discharged by P76 correction 1 (the AT-23 fragment) plus the S-G2 correction (`0f80c1e`).
- **Hand-off in this version**: **there is a new transfer** - the OB-139 correction channel is taken over by **P82** (`TMP-DISCIPLINE-OUT-OF-REPO`), on the user ruling of 2026-11-19; its admission condition is that the start requires the user to name it and that it runs before P79 unless the user rules otherwise; the landing spot is §4 and its ledger segment.
- **Convention note**: this treatment has been used since v59; whenever a slice hands an obligation to a future slice, that slice's fact package must carry a hand-off section.

### Appendix B: the "AT-23 stage-2 status" section (unchanged)

- **Status**: ruled "not applicable" by A1 (the user ruling of 2026-11-09; this is **not a permanent prohibition**, unless the user later rules explicitly to restart it).
- **What can be executed before S2**: nothing (R-3 = (A), and the real CLI end-to-end run needs real enforcement, which is an S2 premise).
- **Restart path**: if the user later rules to restart stage 2, real enforcement must land first (S2) or an alternative path must exist (P-prime, a candidate swap, or a new intermediate path); an alternative path means "reach an equivalent goal before S2 by another route", not "implement stage 2 before S2".
- **Natural restart point**: S2 starts, enforcement lands, S-G8 unlocks, R-3 is re-judged (A possibly becoming B), and if B is feasible that is the natural restart point for stage 2.

---

## §3 Completed slices (recent baseline)

**v63 §3 already carried the sixteen rows P64-P80; v64 §3 reuses those sixteen rows (readings copied from the v63 text, not retro-edited) and appends one P81 row. Seventeen rows in total.**

| Version | Slice | Landing commit | Conclusion |
|---|---|---|---|
| v1.66 | P64 AT24-SUPERVISED-BLACK-BOX | `2dbc54a` | closed (AT-24 conditionally closed; R = 2; run `42948727363` green on the first attempt) |
| v1.66 | P65 GATE-PIT-26-LANDING | `0eda074` | closed (appendix B.7.1 slot plus the A-part ledger todo; R = 1; run `42969802819`) |
| v1.66 | P66 CHECKLIST-RULES-HOIST | `1650d5f` | closed (the five §13 sections plus the new snapshot structure; R = 3, D = 2; run `43022796716`) |
| v1.67 | P67 GATES-DEBT-CLEANUP | `18b6ef9` + `365a930` | closed (g4b LRU, gate pitfalls 27-30, the §12.4 debt-settlement block and the B.6 hardening; R = 3, D = 1; run `43124686702`) |
| v1.67 | P68 REMAINING-RULINGS | `e0d745c` + `665a31f` | closed (assessment type; R = 1, D = 1; (7) PASS; five R items ruled; ADR-042; run `43177025096`) |
| v1.67 | P69 T027B-E2E-WIRING | `759c934` | closed (implementation type; R = 0, D = 1; (7) PASS; E-1 registered; run `43264324421`) |
| v1.68 | P70 S1-STREAM-PROGRESS | `ca88fa3` | closed (implementation type; R = 1, D = 2; (7) PASS; DR-14; runs `43400311736` -> `43403891205` -> `43408438643`) |
| v1.68 | P71 S1-SCOPE-ANNOTATE | `6ee2c61` + `16c98f7` | closed (documentation type; R = 0, D = 2; (7) blind review PASS plus final review PASS; runs `43463089586` + `43468806314`) |
| v1.68 | P72 REQ020-S2-RESTRICTION | `fb0ab57` + `406f4a0` | closed (documentation type; R = 0, D = 0; (6) skipped and confirmed by (7); (7) final review PASS WITH FIXES; runs `43566904684` + `43568186355`) |
| v1.69 | P73 T027B-ADMISSION | `8c59712` + `f834ce5` | closed (assessment type; R = 0, D = 0; (6) skipped and confirmed by (7); (7) final review PASS; **the directive was bumped to v1.69**; runs `43622770750` + `43624661784`; OB-55 = option b) |
| v1.69 | P74 SANDBOX-INTERFACE | `3300364` + `150423d` | closed (implementation type; R = 3, D = 0; (6) required and PASS; (7) final review PASS; **R=3 was the first case to trigger the §9.2 high-risk tier**; runs `43669748483` + `43672994996`; the `Sandbox` check layer landed) |
| v1.69 | P75 R3-REAL-CLI-ASSESSMENT | `ed2e26e` + `70a04cd` + `60146c0` | closed (assessment type; R = 1, D = 1; (6) skipped and confirmed by (7); (7) one blind review plus one final review, blind first as required; **R-3 = (A)**; T027-B = complete (partly not applicable); ADR-045; runs `43732384194` + `43734829063` + `43739409108`) |
| v1.69 | **P76 S1-EXIT-CLOSEOUT-ASSESSMENT** | **`77f23d2` + `aca19af` + `0f80c1e`** | **closed (assessment type; R = 1, D = 1 for P76 as a whole: D = 0 in the main slice and D = 1 in the S-G2 correction [ROUTINE], with the self-healing confirmed by the user; OB-139 ruling: the user decided on 2026-11-16 that this [ROUTINE] D=1 does not count as the sixth occurrence; (6) skipped and confirmed by (7); (7) final review PASS WITH FIXES (none at medium or above; its single Low became the proposed slice name `S1-EXIT-SG2-SYNC`, which never started because the fix landed as the [ROUTINE] correction `0f80c1e`); Q2=(a) leaves S-G1 not passed and S-G2 passed; Q10=(C) two-step; ADR-046 / ADR-047; runs `43785603244` + `43787083404` + `43793273520`)** |
| v1.69 | **P77 ADR010-APPROVAL** | **`33272be` + `fc00c4b`** | **closed (documentation type; R = 1, D = 0; (5) ran twice (r2 with none at medium or above, budget exhausted, the three Low confirmed by (7)) plus (6) skipped and confirmed by (7) plus (7) final review PASS WITH FIXES (none at medium or above); ADR-010 = `APPROVED` (user approval 2026-11-16); S-G4 passed; the p57 / s1-exit gate tables untouched; runs `43986942087` + `43997058549`)** |
| v1.69 | **P78 RELEASE-ZIP-BINDING** | **`c935713` + `0741fae` + `308b235`** | **Closed (implementation type; R = 0, D = 0 (implementation-difference note not counted as D); (5) two rounds (r2 0 Medium+, budget exhausted) + (6) required PASS (0 findings) + (7) final FINDINGS (1 High + 1 Medium) then re-review PASS; **GitHub Release `v0.1.0` created with the tuple bound**; S-G3 passed; the p57 and s1-exit gate tables untouched; runs `44043705691` + `44060224254`)** |
| v1.69 | **P80 SG6-GATELEVEL-CONSOLIDATE** | **`e00d79b` + `58c2794`** | **Closed (documentation type; R = 2 (§9.2 high-risk marking; the tier-S budget was not exceeded, and after stopping and escalating the user ruled option A, one narrowed read), D = 1 (self-imposed); (5) r1 with 0 Medium+; (6) skipped and confirmed by (7); (7) blind review FINDINGS -> final review FINDINGS (1 Medium) -> re-review FINDINGS (1 Medium) -> narrowed read `RE-REVIEW: PASS`; **the S-G6 gate-level side holds**; the p57 and `s1-exit` gate tables untouched; runs `44211368890` + `44213070006`)** |
| v1.70 | **P81 DIRECTIVE-MODEL-UPGRADE** | **`45590e8` + `fd4d8fb` + `9ca6576`** | **Closed (governance `[SLICE]`, M; the stage-6 `MAI-Code-1.1-Flash` becomes `MAI Code 2 Flash` and the stage-7 deep tier `GPT-6.1 Sol` becomes `GPT-6.2 Sol`, a direct replacement without adding Astra; the G5-b closed set agrees three ways with 8 members; SELFTEST 737 -> 740; **R = 1, D = 2**; (6) required (r1 format deviation, r2 PASS, budget 3 with 2 used); (7) one final review PASS (budget 1+3, 1 used); no blind review; the sixth OB-139 trigger moves the entry to "to be discussed" and registers P82; runs `44281812336` + `44284200392` + `44392113028`)** |

> **The D convention in §3** (as in v61): "D = n" in this table is shorthand for the temp-file discipline violation count (§5.1); P81 totals D = 2 (event 1 is the self-healed `tmp/`-root misplacement, event 2 is the "out-of-repository location plus same-command create-delete" case with no self-heal channel). **The D count and the OB-139 ruling are two separate concepts**: the user ruling of 2026-11-19 accepted both P81 events as D = 2 and, on that basis, confirmed the sixth OB-139 trigger (no waiver, not a precedent).

---

## §4 Execution queue

### Layer 0 preliminary triage (✅ all dequeued, nothing new in this version)

- P5 and P6 were closed by P57; layer 0 is fully dequeued with nothing left.

### Layer 1 main-line hard gates

- Items 8-12: ✅ all dequeued (see §4 of v57 / v58; **the original dequeue records have been pushed into earlier versions by the rolling window, so trace them back layer by layer through the citation chain in §4 of v57**).
- **`T027-B` (the defensive-refactor main-line milestone) = complete (partly not applicable)**: C1-C7 / N1-N4 / R-1 / R-2 / R-3 are all dispositioned (R-3 = (A), ruled by the user on 2026-11-15; the A1 declaration caliber is "complete, partly not applicable"). **The main line is fully dequeued.**
- **`S1 Exit close-out` (gate level complete)**: **P76 produced the roadmap; P77 closed S-G4; P78 closed S-G3; P80 decided that the S-G6 gate-level side holds**. The eight gates are listed below (after P80 closes):

| Gate | Criterion | Status (after P80 closes) | Basis |
|---|---|---|---|
| S-G1 | a real CLI launch of an admitted AI candidate with billing | **not passed** (Q2=(a) criterion-literal, agreeing with G-A being unreachable in s1-reach) | P76 §4; s1-reach G-A; ADR-046 |
| S-G2 | the `supervised-black-box` mode is implemented and the AT-24 verdict holds | **passed** (Q2=(a) criterion-literal; T028 `COMPLETED` plus AT-24 conditional closure) | P76 §4; the S-G2 sentence at `s1-exit.md:52` (S-G2 correction `0f80c1e`); ADR-046 |
| S-G3 | Final ZIP bound to a release record and re-verified | **Passed** (the final ZIP carries the bound tuple and passed all four re-verification items; closed in P78) | P78 report sections 5 / 7; `s1-exit:66`; the bound-record section of `S1_SUPPORT_MATRIX` |
| S-G4 | ADR-010 is approved by the product owner | **passed** (ADR-010 is `APPROVED`; user approval on 2026-11-16; P77 closed) | P77 report §5; `ADR_REGISTER.md:20` |
| S-G5 | Product owner issues the approval | **Release part approved (step one)** (ADR-046 Q10=(C); whole S1 Exit approval waits until after S2) | the user's same-turn ruling in the 2026-11-17 kickoff message; P76 section 12; ADR-046 |
| S-G6 | Dependency fields satisfied and the four remaining gates closed | **Gate level holds** (P80 consolidation verdict; S-G3 / S-G4 / S-G5 closed; the T018 task level stays unchecked side by side, still bound by three causes including T027) | P80 report sections 5 / 6; P76 section 10; `DEV_PLAN.md:94` / `:96` |
| S-G7 | every AT passes, T018 is complete and every §14 exit holds | unreachable (a composite gate) | `s1-exit.md:52` / `:56`-`:58` / `:60` |
| S-G8 | the supervised path runs end to end against a real CLI and a real candidate | blocked until S2 enforcement lands | P75 A2 / ADR-045; P76 §6 |

**The independent security review (not an S-G5 criterion, but a precondition of the second approval step)**: the S0-to-S1 gate in `PRODUCT_REQUIREMENTS` §7, which has **never been executed**; **ADR-047** waives it and defers it to S2, with the review to complete before S2 and a due date at the earlier of the S2 start date and 2027-03-31. **The Q10=(C) two-step S-G5 approval distinguishes them**: **step one**, inside S1 and covering the release portion, does not require the review to be complete (ADR-047 defers it to S2), whereas **step two**, after S2 and covering the whole S1 Exit, requires the review to have completed before S2 starts. **Hard boundary: no approval may be read as S1 Exit passing.**

- **S1 Exit cannot pass as a whole inside S1** (the premise of Q2=(a) plus Q10=(C)). **Every advanceable gate-level item is closed** (S-G2 / S-G3 / S-G4 / S-G5 and the S-G6 gate level; S-G7 is a composite gate, S-G8 needs S2, S-G1 does not pass).
- **No admission residual blocks the closeout follow-up slices**: the only P75 §11.2 residual (the core task chain row at `s1-exit.md:52`) was fully corrected by P76 correction 1 (the AT-23 fragment) plus the S-G2 correction (`0f80c1e`).
- **`RELEASE-ZIP-BINDING` (S-G3 close-out)**: **P78 is closed** (release source commit `c935713`; main commit `308b235`; R = 0, D = 0; S-G3 passed; the p57 and `s1-exit` gate tables are untouched). Removed from the queue; the fact pack is in §2; the report is `docs/validation/release-zip-binding{,_EN}.md`.
- **`DIRECTIVE-MODEL-UPGRADE` (stage-6 carrier and stage-7 deep-tier replacement - P81)**: **P81 is closed** (main commit `45590e8`; §11.2 back-fill `fd4d8fb`; probe / OB-139 back-fill `9ca6576`; R = 1, D = 2; **directive v1.70**; the G5-b closed set agrees three ways with 8 members). Removed from the queue; the fact pack is in §2; the report is `docs/validation/directive-model-upgrade{,_EN}.md`.

**Next item = P82 `TMP-DISCIPLINE-OUT-OF-REPO` (governance slice; fixes the §5.1 temp-file discipline, bumps to v1.71; the start requires the user to name it)**

| Field | Content |
|---|---|
| Identifier | `TMP-DISCIPLINE-OUT-OF-REPO` (governance slice; source = the user ruling of 2026-11-19 on the sixth OB-139 trigger: register a fix to §5.1, the directive rises to v1.71; ledger section = **P82**) |
| Contents | add a disposition channel (self-heal or downgrade) for "an out-of-repository location (`%TEMP%` and the like) combined with a same-command create-delete", and state the `tmp/`-root self-heal caliber explicitly; the three readings and the stop-the-slice ceiling of the existing self-heal grade are not changed |
| Size tier | **M** (governance slice; touches the §5.1 body plus §0.3 and appendix C; no criterion code expected to change) |
| (7) budget | M: 1 final review + 2 re-reviews; the (7) model follows the §5.2 tier table (governance defaults to Codex Extra High) |
| Dependency | P81 closed (the sixth trigger was ruled by the user on 2026-11-19); **the start requires the user to name it**; the order relative to P79 is the user's call (the current advice is P82 first) |
| Boundary | no change to the three self-heal readings or the D ceiling; no change to the criteria proper or the gate tables; no change to CONTRACTS / schema / fixtures / product code / CI; no gitee action |

- **`S2-PRE-INDEPENDENT-SECURITY-REVIEW` (the ADR-047 compensation obligation · P79)**

| Field | Content |
|---|---|
| Identifier | `S2-PRE-INDEPENDENT-SECURITY-REVIEW` (ledger section = **P79**, registered when the kick-off message is drafted) |
| Contents | complete the independent security review (the S0-to-S1 gate in `PRODUCT_REQUIREMENTS` §7); archive the review conclusion under `docs/validation/` and update the note at `SECURITY.md:5` to "completed, with the review date and reviewer"; a later slice may flip the ADR-047 risk-acceptance record to "compensated" |
| Size tier | **M-L** (`[SLICE]`, assessment type; the review is work in its own right) |
| (7) budget | tier M-L: 1 final review + 2 re-reviews (§7.5) |
| Dependency | before S2 starts (not an immediate S1 need); the ADR-047 due date is the earlier of the S2 start date and 2027-03-31; **the start order is the user's call (the current advice is P82 first, P79 next)** |
| Boundary | the review scope is the S0-to-S1 security gate; the S1 Exit eight-gate criteria stay; the `SECURITY.md` threat-model baseline stays (only the review conclusion is archived); no directive change |

- **`S-G8`** stays blocked until S2 enforcement lands (P73 A2 / P75 A2 / ADR-044 / ADR-045); the gate table in `s1-exit{,_EN}.md` is untouched (P76 only edited the fragment inside row `:52`).
- **Post-S2 re-evaluations**: S-G5 step two (whole S1 Exit approval), S-G1, S-G8 and S-G7, in the order given by P80 report section 7.

### Layer 2 routine governance

- P72 / P73 / P74 / P75 / P76 / P77 / P78 / P80 / P81 ✅ removed from the queue (tier 2 is **empty**; P82 is a tier-1 queue item awaiting the user naming it; per-row summaries are in §3).
- The full §1.2 field definitions for P72 / P73 / P74 / P75 / P76 / P77 / P78 / P80 / P81 are in the corresponding segments of `docs/t027/REMAINING_SLICES{,_EN}.md`; the P82 fields are in the matching segment of that ledger.

### Layer 3 low-priority options

| Field | Content |
|---|---|
| Identifier | `OB-50` / `spawnGit` / `ci freeze cost` / 22 historical `%TEMP%` items / `GATES-EXEMPT` |
| Current state (verified item by item) | OB-50 = **settled** (§12.4, v1.41); spawnGit = the git-call helper of `tools/gates/selftest.js` (the spawn wrapper in `lib/ctx.js`) that folds stderr into the result object instead of printing it; the P27 `G4b-CACHE` report §5 item 5 registered it as a "candidate OB (diagnostic)" (Low / observing) and it was never formally registered as a §12.4 OB id, so under the §1.8.1 pending-item rule it is an "unregistered candidate OB"; candidate fix = attach the last N bytes of stderr to `detail` when the spawn fails (diagnostics only, the verdict is unchanged); trigger = the next unexplained git-related assertion failure is assessed by (1); ci freeze cost = **no §12.4 entry**, its real origin is item 7 of the P27 report as a "candidate OB (performance)", **not formally registered**; the 22 historical %TEMP% items map to OB-98 (**settled**) / OB-99 (**observing**) / OB-100 (**observing**); GATES-EXEMPT = **not initiated** (no §12.4 entry; the P33 item 1-3(iv) split slice, touching four criterion files) |
| Contents | see the matching §12.4 rows; `GATES-EXEMPT` is the P33 item 1-3(iv) split slice, not initiated |
| Dependency | the long suspension covers only unsettled items (spawnGit / ci freeze cost / the %TEMP% timing shapes (OB-99 / OB-100) / GATES-EXEMPT); there is no concrete trigger |
| Size | varies (S through M-L) |
| (7) budget | none |
| Boundary | per item in §12.4; this version adds no layer-3 item and advances no suspended item |
| Start condition | the user must name the item explicitly |

**Queue discipline**: as in v60 (layer boundaries and numbering discipline; this [ROUTINE] slice runs in a single batch and does not bump the version).

---

## §5 Observation items

- Under observation (**13**): `OB-1 / 7 / 10 / 33 / 87 / 99 / 100 / 104 / 105 / 107 / 111 / 113 / 135` (carried from v61; OB-139 moved to pending discussion on 2026-11-19).
  - **OB-139 moved to pending discussion (user ruling 2026-11-19)**: both temp-file events of P81 are accepted as D = 2, and the "out-of-repository location (`%TEMP%`) plus same-command create-delete" combination has no self-heal channel under section 5.1, so the **sixth trigger is confirmed**; the entry moves from "under observation" to "pending discussion" and slice **P82** is registered (`TMP-DISCIPLINE-OUT-OF-REPO`, directive-governance slice, bumps to v1.71, schedule pending, start requires the user to name it). The ruling grants no waiver and may not be cited as a precedent.
  - **The reinforced D=2 recurrence-prevention sentence (standard kick-off content since v59, carried into v61)**: "the executing roles (`prfrail-implementer` / `prfrail-tester` / `prfrail-documenter`) **and the controller** must always work under a `tmp/` slice directory (`tmp/<slice-id>/`); **the same holds for CI observation captures, which must never land at the `tmp/` root**; three D occurrences inside one slice force the slice to stop."
- Pending discussion (10): `OB-19 / 21 / 52 / 53 / 57 / 58 / 59 / 61 / 75 / 139`.
- To do (1): `OB-69`. Non-terminal total **24** (13 plus 10 plus 1).
- **Slice-level suggestions (three, not registered as OBs)**:
  1. **`R-ACCUMULATION-MECHANISM` (carried from v58)**: P74 was the first case to trigger the §9.2 rule (R>=2 means high risk) together with the §7.5 high-risk tier (one plus three); later implementation slices should pre-place the dynamic budget sentence in their kick-off message; P75 / P76 / P77 recorded R = 1, below the threshold, so it did not trigger. **No OB is registered now**; a recurrence in a later slice turns it into a formal OB.
  2. **`--progress` covers state transitions only, not log output (carried from v58)**: if long-running observability becomes a real need (the S1 Exit closeout slices are the most likely trigger), evaluate adding `--progress verbose` or `--tail-logs`.
  3. **The "obligation hand-off" convention (in use since v59, first applied to P75 O-7)**: whenever a slice hands an obligation to a future slice, that slice's fact package must carry a hand-off section, the landing slice entry in §4 must carry an "admission condition" section, and §5 may track it as an optional "pending execution list" (without an OB). **No OB is registered now**; if a hand-off is ever left dangling it becomes a formal OB.
- **Removed in this version**: none.
- **Flipped to `settled` in this version**: none.
- Mechanical readings: max OB in §12.4 = **139**; **max ADR = 47**; max DR = 14; max deferral = round 25.

---

## §6 Pending user rulings

- **The table is currently empty** (the three P81 rulings of 2026-11-19 have been made, see the last bullet of §2 and §5).
- **Next ruling item (expected)**: **the start order of P82 versus P79** (the current advice is **P82 first**: the OB-139 correction channel already has two instances and is a recurring discipline gap, while the P79 due date is the earlier of the S2 start date and 2027-03-31 and still has room).
- **Three P81 rulings**: the two extra probes with `R = 1`, the acceptance of both temp-file events as `D = 2`, and the sixth OB-139 trigger moving the entry to "to be discussed" with P82 registered; all three are marked "made directly by the user, not usable as a precedent".

---

## §7 / §8 / §9 Pointer sections

Same structure as v60: §7 (kick-off and step-0 general rules) points at directive §13.1; §8 (cross-slice collaboration discipline) at §13.2; §9 (checklist meta-rules and gate pitfalls) at §13.3 / §13.4. **The rule bodies have moved into directive §13 and are not copied here.**

---

## §10 One-line opening for a new session

> This is my project state snapshot (v64); check the eight anchors in section 1 first, then continue per the section 4 queue, and if anything disagrees trust the repository.
>
> [paste this checklist in full v64]
>
> [supporting documents] `REMAINING_SLICES{,_EN}.md` (P81 closed, P82 registered), `DELIVERY_DIRECTIVE.md` (**v1.70**), `directive-model-upgrade{,_EN}.md` (P81 report), `sg6-gatelevel-consolidate{,_EN}.md` (P80 report), `release-zip-binding{,_EN}.md` (P78 report), `ADR_REGISTER{,_EN}.md` (with ADR-045 / ADR-046 / ADR-047), `DOCUMENTATION_PLAN{,_EN}.md` (cross-session handoff entry pointer), `CHECKLIST_v63{,_EN}.md` (frozen, superseded by v64).
>
> Current slice: P81 is closed (main commit `45590e8` plus back-fills `fd4d8fb` and `9ca6576`). Next = **P82 (fixes the §5.1 temp-file discipline, bumps to v1.71, requires the user to name it)**; then P79 (the pre-S2 independent security review).
## §11 One-line current state
P81 closed (main commit `45590e8` plus back-fills `fd4d8fb` and `9ca6576` pushed to origin, runs `44281812336` + `44284200392` + `44392113028` all success; **R = 1 (probe-target deviation, corrective rework), D = 2 (two temp-file events)**; directive **v1.70**, the stage-6 `MAI Code 2 Flash` / stage-7 deep tier `GPT-6.2 Sol` replacement is in place, the G5-b closed set agrees three ways with 8 members, SELFTEST 740/740; the sixth OB-139 trigger moves the entry to "to be discussed" and registers P82); every advanceable tier-1 gate-level item stays closed; next = **P82 (bumps to v1.71, requires the user to name it)**, then P79.
