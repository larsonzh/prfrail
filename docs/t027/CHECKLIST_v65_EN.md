# Task Execution Checklist (cross-session snapshot v65 · after P82 closes · the section 5.1 temp-file discipline fix is done · next = the P79 security review)

Generated: 2026-11-19. HEAD = `9fd8cd6`; directive = v1.71; five CI runs, all success: the main commit `0264bcd` -> `44447198220` (both legs success), then `adb8cb9` -> `44449369232`, `1015430` -> `44450683923`, `9d83089` -> `44451589047` and `9fd8cd6` -> `44454737938` (the last four are observation runs).

**Supersession**: this version supersedes v64 (`docs/t027/CHECKLIST_v64{,_EN}.md` is frozen and is not retro-edited).

**Structure**: §1 anchors -> §2 fact package -> §3 completed slices -> §4 execution queue -> §5 observation items -> §6 pending rulings -> §7 / §8 / §9 pointer sections -> §10 opening -> §11 one-line status; the rule bodies behind §7-§9 have moved into directive §13, and this snapshot does not copy them.

**Changes in this version (relative to v64)**:
- (1) §1 anchor rollover: HEAD `9fd8cd6`; directive **v1.71**; SELFTEST **740/740**; the rolling window drops v62 and admits v65 (v63 / v64 / v65).
- (2) §2 becomes the P82 fact pack (the §5.1 sub-clause and the `tmp/` caliber, D = 3 with the stop-and-escalate, OB-139 resolved plus its seventh literal trigger, R = 2, five CI runs all success, and the four user rulings).
- (3) §3 appends one P82 row (the seventeen rows of v64 plus one = eighteen rows).
- (4) §4 queue update: P82 leaves the queue; next = P79 (`S2-PRE-INDEPENDENT-SECURITY-REVIEW`, the start requires the user to name it); the layer-3 `GATES-PERF-ASSESS` candidate stays (registered by ruling, not initiated).
- (5) §5: OB-139 resolved and discharged (P82 closed), and A-6 is its seventh literal trigger; 13 observing, 9 pending, 23 non-terminal (re-checked against the repository, unchanged).
- (6) §6 is cleared; the next ruling item is expected to be the start timing of P79.
- (7) §10 / §11 roll over; the cross-session handoff pointer in `DOCUMENTATION_PLAN{,_EN}.md` now points at v65.

---

## §1 Current anchors (check these eight numbers first)

**Rolling window = the last three versions (v63 / v64 / v65)**; v62 has been dropped (same rule that dropped v52).

| Item | Value | Check command |
|---|---|---|
| HEAD | `9fd8cd6` | `git rev-parse --short HEAD` |
| origin/main | `9fd8cd6` | `git rev-parse --short origin/main` |
| Directive version | **v1.71** | the header of `docs/DELIVERY_DIRECTIVE.md` |
| SELFTEST | **740/740** | `node tools/gates/gate.js --selftest` |
| Three scopes | **all TOTAL_FAIL=0** | `gate.js --all --scope=index` / `tree` / `ci` |
| v63 total lines | **222** (LF count, CN / EN symmetric; re-checked by v64) | `docs/t027/CHECKLIST_v63{,_EN}.md` by LF |
| v64 total lines | **224** (LF count, CN / EN symmetric; **re-checked by v65**) | `docs/t027/CHECKLIST_v64{,_EN}.md` by LF |
| v65 total lines | **re-checked by v66** (HEAD advances once this version lands) | `docs/t027/CHECKLIST_v65{,_EN}.md` by LF |

> This snapshot's HEAD = `9fd8cd6` is the checkpoint after the P82 user-closure-ruling commit; once v65 lands, HEAD advances, so use `git rev-parse --short HEAD` as the live value.

> Rolling window = the most recent 3 versions (v63 / v64 / v65); v62 has dropped out (the same rule that dropped v52); this snapshot's own line count is re-checked by v66. If two or more items disagree, trust the repository and report the terminal reading before continuing.

`tmp` fixed assertion: top level = `.gitkeep` + `gates-cache`; `gates-cache` holds only `g4b`; `tmp/P82-TMP-DISCIPLINE-OUT-OF-REPO/` left no residue after the P82 ⑧c cleanup and the later window cleanups.

---

## §2 P82 fact pack

- **Slice / size**: `TMP-DISCIPLINE-OUT-OF-REPO` (P82; the §5.1 temp-file discipline gains a disposition channel plus the `tmp/`-root self-heal caliber; governance `[SLICE]`) / **M**; controller `DeepSeek V4.1 Flash (deepseek)`; the user named the start on 2026-11-19 and authorised the slice in the same turn.
- **Goal**: per the user's 2026-11-19 registration ruling on the sixth OB-139 trigger, fix the §5.1 temp-file discipline by adding a disposition channel for "an out-of-repo location (`$env:TEMP` / `%TEMP%` and the like) plus a same-command create-and-delete", stating the per-element caliber of `tmp/`-root self-heal and revising the "out-of-repo locations are strict" sentence; flip the OB-139 disposition cell to `Resolved`; bump the directive to **v1.71**.
- **Verdict**: **closed**; the v1.71 landing points are the new §5.1 "out-of-repo same-command create-and-delete self-heal sub-clause (v1.71)" and the new "per-element caliber for `tmp/`-root self-heal", the revised "out-of-repo locations are strict" sentence, the v1.71 log row in §0.3, the OB-139 flip to `Resolved` in §12.4, and appendix `C.0ax`; CN/EN are strictly co-located with all six shapes equal; no change to the three self-heal readings, the D ceiling, the five gates of the v1.56 exception, the §6.3 criteria or `tools/gates/**` (**zero code change**).
- **Channel design**: the **self-heal tier gains a sub-clause** (not a downgrade tier) - the sub-clause states only the applicability conditions, the applicability precondition (the command must yield at least one creation action, and creation paths must correspond one-to-one with deletion actions) and the differentiating elements a-d (the out-of-repo caliber of gate 3, gate 5, the no-self-attestation constraint and the pre-snapshot completeness waiver); gates 1, 2 and 4 plus the waiver, the completeness substitute, the review-layer confirmation, the counting and the fallback sentence **follow the same-numbered wording of v1.56** (one caliber written once, to prevent drift).
- **Review chain**: (1) design `DESIGN: DONE` (deep-reasoner + `DeepSeek V4 Pro`) -> batch assistance `GLM-5.3-Flash` **twice** (instance enumeration and bilingual-mirror check) -> (2) implementation three times (initial, stage-5 r1 fix, stage-6 r1 fix) -> (3) consistency twice (C-1 to C-8 plus three mutations, all three caught) -> (4) gates -> (5) **three times** (r1 `PASS WITH FIXES` 3M+1L -> r2 `PASS WITH FIXES` 1M+2L -> the user-authorised narrowed read r3 `PASS`) -> (6) **`MAI Code 2 Flash` first engagement**: r1 `FINDINGS` (3M) -> r2 `INDEPENDENT SCAN: PASS` (0 findings, budget 3 with 2 used) -> (7) r1 `FINDINGS` (1 High + 1 Medium, both missing-artifact type) -> r2 `RE-REVIEW: PASS` (both CLOSED) -> ⑧a / ⑨ / ⑧b / ⑧c -> ⑩.
- **Budget and blind review**: the (7) governance budget is **1 + 3**: one final review plus three in-budget re-reviews are exhausted; in addition there were **two user-authorised narrowed reads** (B-prime and A-9, off-budget and adding no R). The (5) budget of 1 + 1 was exhausted and the user authorised one narrowed read. **No blind review** (none of the three mandatory semantics is touched and this is not the trigger slot of the proportional window).
- **Metering**: **R = 2** (one is the re-run after the three stage-5 r1 Medium fixes; the other is the second fix re-run after the fix script rebuilt the snapshot from HEAD and overwrote the tier-3 registration), so **R = 2 has raised the §9.2 high-risk mark**; **D = 3** (A-6 = a `tmp/`-root misplacement; A-7 = a `$env:TEMP` artifact; A-9 = a CI-watch capture written to `$env:TEMP`), all three handled under the §5.1 self-heal tier with evidence and zero residue.
- **D = 3 stop-and-escalate**: **D reaching 3 hits the §5.1 cumulative ceiling, so the slice stopped and escalated**; the §7.8.5 five-field trace is recorded (object = post-fix re-review / stop-and-escalate; basis = the §5.1 ceiling plus §7.8.3; verdict lines verbatim = `RE-REVIEW: PASS` and `RE-REVIEW: FINDINGS`; budget count stage 5 = 3 / stage 6 = 2 / stage 7 = 1+4); **the closing steps were authorised by the user ruling of 2026-11-19 (a §7.8.3 user-authorised exception) and may not be cited as a precedent**.
- **OB-139 disposition**: the cell flips from `Pending` to **`Resolved`** (the phenomenon column stays verbatim); **A-6 is registered as OB-139's seventh literal trigger** (recurring while the entry sat in pending discussion), recording the fact only, adding no rule and granting no waiver.
- **Probes**: **0** (this slice changes no model and needs no real external call).
- **Gates**: `SELFTEST: PASS 740/740`; the `index` / `tree` / `ci` scopes all report `TOTAL_FAIL=0`; the G1-b version metadata `v1.71 / 2026-11-19` agrees across both languages.
- **Evidence pack**: `docs/validation/evidence/P82-TMP-DISCIPLINE-OUT-OF-REPO/` (`MANIFEST.md` plus `SHA256SUMS.txt` plus its payload set: the kickoff message, the design, two batch-assistance payloads, seven review transcripts, four gate transcripts, the freedom list and the input list); the three entry sets agree and re-computing every sha256 yields **zero mismatches**.
- **CI (five commits, all success)**: main commit `0264bcd` -> run **`44447198220`** (`Go ubuntu-latest` plus `Go windows-latest` both green); `adb8cb9` -> `44449369232`; `1015430` -> `44450683923`; `9d83089` -> `44451589047`; `9fd8cd6` -> `44454737938` (the last four are observation runs).
- **Mechanical readings**: `SELFTEST PASS 740/740`; max ADR = 47; max OB = **139** (no new OB); max DR = 14; deferral = round 25.
- **Ledger convention** (§11.2): the P82 completion record carries only the report path plus R/D (no hash / run); the CI fact chain lives in section 2 of this snapshot and in the (10) receipt.
- **User rulings (2026-11-19; made directly by the user, not usable as a precedent)**:
  1. **D acceptance and OB-139**: the three P82 temp-file events are accepted as **D = 3**, and A-6 forms OB-139's seventh literal trigger (no waiver, not a precedent).
  2. **Closing the D=3 stop-and-escalate**: the registration is accepted as faithful and the user authorised the §11.2 back-fill plus the (10) receipt (a §7.8.3 user-authorised exception, this slice only, not a precedent).
  3. **Closing the A-9 narrowed-read FINDINGS**: the user directly ruled that the FINDINGS is closed by `9d83089` (basis = an independent verbatim check that the A-9 readings are complete, CN/EN are symmetric and the gates are green; a user ruling, not a precedent).
  4. **The stage-5 narrowed read**: after the stage-5 budget ran out the user authorised one narrowed read (this slice only, not a precedent).

### Appendix A: "obligation transfer" (this version: no new transfer, the P82 obligation is discharged)

- **Historical hand-off** (first traced back to P75 O-7): the in-cell caliber mismatch in the core task chain row at `s1-exit.md:52` was discharged by P76 correction 1 (the AT-23 fragment) plus the S-G2 correction (`0f80c1e`).
- **Hand-off in this version**: **no new transfer** - the v64 hand-off "OB-139 correction channel taken over by P82" has been discharged by P82 closing (evidence = the §3 row and the P82 report).
- **Convention note**: this treatment has been used since v59; whenever a slice hands an obligation to a future slice, that slice's fact package must carry a hand-off section.

### Appendix B: the "AT-23 stage-2 status" section (unchanged)

- **Status**: ruled "not applicable" by A1 (the user ruling of 2026-11-09; this is **not a permanent prohibition**, unless the user later rules explicitly to restart it).
- **What can be executed before S2**: nothing (R-3 = (A), and the real CLI end-to-end run needs real enforcement, which is an S2 premise).
- **Restart path**: if the user later rules to restart stage 2, real enforcement must land first (S2) or an alternative path must exist (P-prime, a candidate swap, or a new intermediate path); an alternative path means "reach an equivalent goal before S2 by another route", not "implement stage 2 before S2".
- **Natural restart point**: S2 starts, enforcement lands, S-G8 unlocks, R-3 is re-judged (A possibly becoming B), and if B is feasible that is the natural restart point for stage 2.

---

## §3 Completed slices (recent baseline)

**v63 §3 already carried the sixteen rows P64-P80; v64 §3 reused them and appended one P81 row (seventeen rows); v65 §3 reuses those seventeen rows (readings copied from the v64 text, not retro-edited) and appends one P82 row. Eighteen rows in total.**

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
| v1.69 | **P78 RELEASE-ZIP-BINDING** | **`c935713` + `0741fae` + `308b235`** | **Closed (implementation type; R = 0, D = 0 (implementation-difference note not counted as D); (5) two rounds (r2 0 Medium+, budget exhausted) + (6) required PASS (0 findings) + (7) final FINDINGS (1 High + 1 Medium) then re-review PASS; **GitHub Release `v0.1.0` created with the tuple bound**; S-G3 passed; the p57 and `s1-exit` gate tables untouched; runs `44043705691` + `44060224254`)** |
| v1.69 | **P80 SG6-GATELEVEL-CONSOLIDATE** | **`e00d79b` + `58c2794`** | **Closed (documentation type; R = 2 (§9.2 high-risk marking; the tier-S budget was not exceeded, and after stopping and escalating the user ruled option A, one narrowed read), D = 1 (self-imposed); (5) r1 with 0 Medium+; (6) skipped and confirmed by (7); (7) blind review FINDINGS -> final review FINDINGS (1 Medium) -> re-review FINDINGS (1 Medium) -> narrowed read `RE-REVIEW: PASS`; **the S-G6 gate-level side holds**; the p57 and `s1-exit` gate tables untouched; runs `44211368890` + `44213070006`)** |
| v1.70 | **P81 DIRECTIVE-MODEL-UPGRADE** | **`45590e8` + `fd4d8fb` + `9ca6576`** | **Closed (governance `[SLICE]`, M; the stage-6 `MAI-Code-1.1-Flash` becomes `MAI Code 2 Flash` and the stage-7 deep tier `GPT-6.1 Sol` becomes `GPT-6.2 Sol`, a direct replacement without adding Astra; the G5-b closed set agrees three ways with 8 members; SELFTEST 737 -> 740; **R = 1, D = 2**; (6) required (r1 format deviation, r2 PASS, budget 3 with 2 used); (7) one final review PASS (budget 1+3, 1 used); no blind review; the sixth OB-139 trigger moves the entry to "to be discussed" and registers P82; runs `44281812336` + `44284200392` + `44392113028`)** |
| v1.71 | **P82 TMP-DISCIPLINE-OUT-OF-REPO** | **`0264bcd` + `adb8cb9` + `1015430` + `9d83089` + `9fd8cd6`** | **closed (governance `[SLICE]`, M; §5.1 gains the "out-of-repo same-command create-and-delete self-heal sub-clause" and the "per-element caliber for `tmp/`-root self-heal", the "out-of-repo locations are strict" sentence is revised, OB-139 flips to `Resolved`, and the directive rises to v1.71; **R = 2, D = 3**; D=3 fired the §5.1 stop-and-escalate with the §7.8.5 trace; stage 5 ran 3 times (including a user-authorised narrowed read), stage 6 `MAI Code 2 Flash` first engagement r1 FINDINGS -> r2 PASS (budget 2 of 3), stage 7 used 1 final review + 3 re-reviews plus 2 user-authorised narrowed reads; no blind review; all five CI runs success)** |

> **The D convention in §3** (as in v61): "D = n" in this table is shorthand for the temp-file discipline violation count (§5.1); P82 totals **D = 3** (A-6 the `tmp/`-root misplacement, A-7 the `$env:TEMP` artifact, A-9 the CI-watch capture written to `$env:TEMP`), all three handled under the self-heal tier with evidence and zero residue, and **D reaching 3 fired the stop-and-escalate** (§7.8.5 five-field trace). **The D count and the OB-139 ruling are two separate concepts**: A-6 additionally forms OB-139's seventh literal trigger (fact only, no waiver, not a precedent).

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
- **`DIRECTIVE-MODEL-UPGRADE` (stage-6 carrier and stage-7 deep-tier replacement - P81)**: **P81 is closed** (main commit `45590e8`; §11.2 back-fill `fd4d8fb`; probe / OB-139 back-fill `9ca6576`; R = 1, D = 2; **directive v1.70**; the G5-b closed set agrees three ways with 8 members). Removed from the queue; the fact pack is in v64 §2; the report is `docs/validation/directive-model-upgrade{,_EN}.md`.

- **`TMP-DISCIPLINE-OUT-OF-REPO` (the §5.1 temp-file discipline · P82)**: **P82 is closed** (main commit `0264bcd`; follow-ups `adb8cb9` / `1015430` / `9d83089` / `9fd8cd6`; R = 2, D = 3 (which fired the stop-and-escalate); **directive v1.71**). Removed from the queue; the fact pack is in §2; the report is `docs/validation/tmp-discipline-out-of-repo{,_EN}.md`.

**Next item = P79 `S2-PRE-INDEPENDENT-SECURITY-REVIEW` (the ADR-047 compensation obligation; `[SLICE]`, assessment type; the start requires the user to name it)**

| Field | Content |
|---|---|
| Identifier | `S2-PRE-INDEPENDENT-SECURITY-REVIEW` (ledger section = **P79**, registered when the kick-off message is drafted) |
| Contents | complete the independent security review (the S0-to-S1 gate in `PRODUCT_REQUIREMENTS` §7); archive the review conclusion under `docs/validation/` and update the note at `SECURITY.md:5` to "completed, with the review date and reviewer"; a later slice may flip the ADR-047 risk-acceptance record to "compensated" |
| Size tier | **M-L** (`[SLICE]`, assessment type; the review is work in its own right) |
| (7) budget | tier M-L: 1 final review + 2 re-reviews (§7.5) |
| Dependency | before S2 starts (not an immediate S1 need); the ADR-047 due date is the earlier of the S2 start date and 2027-03-31; **the start order is the user's call (P82 closed, P79 awaiting the start)** |
| Boundary | the review scope is the S0-to-S1 security gate; the S1 Exit eight-gate criteria stay; the `SECURITY.md` threat-model baseline stays (only the review conclusion is archived); no directive change |

- **`S-G8`** stays blocked until S2 enforcement lands (P73 A2 / P75 A2 / ADR-044 / ADR-045); the gate table in `s1-exit{,_EN}.md` is untouched (P76 only edited the fragment inside row `:52`).
- **Post-S2 re-evaluations**: S-G5 step two (whole S1 Exit approval), S-G1, S-G8 and S-G7, in the order given by P80 report section 7.

### Layer 2 routine governance

- P72 / P73 / P74 / P75 / P76 / P77 / P78 / P80 / P81 / P82 ✅ removed from the queue (tier 2 is **empty**; per-row summaries are in §3).
- The full §1.2 field definitions for P72 / P73 / P74 / P75 / P76 / P77 / P78 / P80 / P81 are in the corresponding segments of `docs/t027/REMAINING_SLICES{,_EN}.md`; the P82 fields are in the matching segment of that ledger.

### Layer 3 low-priority options

| Field | Content |
|---|---|
| Identifier | `OB-50` / `spawnGit` / `ci freeze cost` / 22 historical `%TEMP%` items / `GATES-EXEMPT` / `GATES-PERF-ASSESS` |
| Current state (verified item by item) | OB-50 = **settled** (§12.4, v1.41); spawnGit = the git-call helper of `tools/gates/selftest.js` (the spawn wrapper in `lib/ctx.js`) that folds stderr into the result object instead of printing it; the P27 `G4b-CACHE` report §5 item 5 registered it as a "candidate OB (diagnostic)" (Low / observing) and it was never formally registered as a §12.4 OB id, so under the §1.8.1 pending-item rule it is an "unregistered candidate OB"; candidate fix = attach the last N bytes of stderr to `detail` when the spawn fails (diagnostics only, the verdict is unchanged); trigger = the next unexplained git-related assertion failure is assessed by (1); ci freeze cost = **no §12.4 entry**, its real origin is item 7 of the P27 report as a "candidate OB (performance)", **not formally registered**; the 22 historical %TEMP% items map to OB-98 (**settled**) / OB-99 (**observing**) / OB-100 (**observing**); GATES-EXEMPT = **not initiated** (no §12.4 entry; the P33 item 1-3(iv) split slice, touching four criterion files); GATES-PERF-ASSESS = **a new candidate** (registered by the user ruling of 2026-11-19, **not initiated as an implementation**; a gate-cost composition assessment that measures how `SELFTEST` and the three scopes spend their time and how often the `base` trio repeats; an assessment-grade candidate that needs the user to name it at start) |
| Contents | see the matching §12.4 rows for the six items (`GATES-PERF-ASSESS` has no §12.4 entry and is registered only in this table and the ledger); `GATES-EXEMPT` is the P33 item 1-3(iv) split slice, not initiated |
| Dependency | the long suspension covers only unsettled items (spawnGit / ci freeze cost / the %TEMP% timing shapes (OB-99 / OB-100) / GATES-EXEMPT); `GATES-PERF-ASSESS` has no prerequisite and ranks after P79; none of them has a concrete trigger |
| Size | varies (S through M-L) |
| (7) budget | none |
| Boundary | per item in §12.4; this version adds one layer-3 candidate (`GATES-PERF-ASSESS`, registered by the user ruling of 2026-11-19 and not initiated as an implementation) and advances no suspended item |
| Start condition | the user must name the item explicitly |

**Queue discipline**: as in v60 (layer boundaries and numbering discipline; this [ROUTINE] slice runs in a single batch and does not bump the version).

---

## §5 Observation items

- Under observation (**13**): `OB-1 / 7 / 10 / 33 / 87 / 99 / 100 / 104 / 105 / 107 / 111 / 113 / 135` (carried from v61).
  - **OB-139 is resolved and discharged (P82 closed, v1.71)**: §5.1 adds the "out-of-repo same-command create-and-delete self-heal sub-clause" and fixes the per-element caliber of `tmp/`-root self-heal; **A-6 forms its seventh literal trigger** (fact only, no waiver, not a precedent); this slice records **D = 3** and fired the stop-and-escalate.
  - **The reinforced D=2 recurrence-prevention sentence (standard kick-off content since v59, carried into v61)**: "the executing roles (`prfrail-implementer` / `prfrail-tester` / `prfrail-documenter`) **and the controller** must always work under a `tmp/` slice directory (`tmp/<slice-id>/`); **the same holds for CI observation captures, which must never land at the `tmp/` root**; three D occurrences inside one slice force the slice to stop."
- Pending discussion (9): `OB-19 / 21 / 52 / 53 / 57 / 58 / 59 / 61 / 75`.
- To do (1): `OB-69`. Non-terminal total **23** (13 plus 9 plus 1).
- **Slice-level suggestions (three, not registered as OBs)**:
  1. **`R-ACCUMULATION-MECHANISM` (carried from v58)**: P74 was the first case to trigger the §9.2 rule (R>=2 means high risk) together with the §7.5 high-risk tier (one plus three); later implementation slices should pre-place the dynamic budget sentence in their kick-off message; P75 / P76 / P77 recorded R = 1, below the threshold, so it did not trigger. **No OB is registered now**; a recurrence in a later slice turns it into a formal OB.
  2. **`--progress` covers state transitions only, not log output (carried from v58)**: if long-running observability becomes a real need (the S1 Exit closeout slices are the most likely trigger), evaluate adding `--progress verbose` or `--tail-logs`.
  3. **The "obligation hand-off" convention (in use since v59, first applied to P75 O-7)**: whenever a slice hands an obligation to a future slice, that slice's fact package must carry a hand-off section, the landing slice entry in §4 must carry an "admission condition" section, and §5 may track it as an optional "pending execution list" (without an OB). **No OB is registered now**; if a hand-off is ever left dangling it becomes a formal OB.
- **Removed in this version**: none.
- **Flipped to `settled` in this version**: `OB-139` (the disposition cell actually changed in the P82 main commit `0264bcd`; the v64 snapshot carried the preview wording).
- Mechanical readings: max OB in §12.4 = **139**; **max ADR = 47**; max DR = 14; max deferral = round 25.

---

## §6 Pending user rulings

- **The table is currently empty** (the four P82 rulings of 2026-11-19 have been made, see the last bullet of §2).
- **Next ruling item (expected)**: **the start timing of P79** (P82 started and closed on 2026-11-19; the P79 due date is the earlier of the S2 start date and 2027-03-31).
- **Four P82 rulings**: the three temp-file events accepted as `D = 3` with A-6 as OB-139's seventh literal trigger; authorising the closing steps of the D=3 stop-and-escalate; the user's closure of the A-9 narrowed-read FINDINGS; and the narrowed read after the stage-5 budget ran out; all four are marked "made directly by the user, not usable as a precedent".

---

## §7 / §8 / §9 Pointer sections

Same structure as v60: §7 (kick-off and step-0 general rules) points at directive §13.1; §8 (cross-slice collaboration discipline) at §13.2; §9 (checklist meta-rules and gate pitfalls) at §13.3 / §13.4. **The rule bodies have moved into directive §13 and are not copied here.**

---

## §10 One-line opening for a new session

> This is my project state snapshot (v65); check the eight anchors in section 1 first, then continue per the section 4 queue, and if anything disagrees trust the repository.
>
> [paste this checklist in full v65]
>
> [supporting documents] `REMAINING_SLICES{,_EN}.md` (P82 closed, P79 awaiting the start), `DELIVERY_DIRECTIVE.md` (**v1.71**), `tmp-discipline-out-of-repo{,_EN}.md` (P82 report), `directive-model-upgrade{,_EN}.md` (P81 report), `sg6-gatelevel-consolidate{,_EN}.md` (P80 report), `release-zip-binding{,_EN}.md` (P78 report), `ADR_REGISTER{,_EN}.md` (with ADR-045 / ADR-046 / ADR-047), `DOCUMENTATION_PLAN{,_EN}.md` (cross-session handoff entry pointer), `CHECKLIST_v64{,_EN}.md` (frozen, superseded by v65).
>
> Current: P82 is closed (R = 2, D = 3, with the D=3 stop-and-escalate closed by a user ruling); next = P79 (the pre-S2 independent security review, the start requires the user to name it).

## §11 One-line current state
P82 closed (main commit `0264bcd` plus `adb8cb9` + `1015430` + `9d83089` + `9fd8cd6` pushed to origin, all five runs success; **R = 2, D = 3**, D=3 fired the stop-and-escalate which a user ruling closed and traced; directive **v1.71**: §5.1 gains the "out-of-repo same-command create-and-delete self-heal sub-clause" and the `tmp/`-root per-element caliber, OB-139 flips to `Resolved`, and appendix C.0ax lands; stage-6 `MAI Code 2 Flash` first engagement r1 FINDINGS -> r2 PASS; SELFTEST 740/740); every advanceable tier-1 gate-level item stays closed; next = **P79 (the start requires the user to name it)**.
