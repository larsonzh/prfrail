# Task Execution Checklist (cross-session snapshot v66 · after P79 closes · the independent security review is done · next = the P83 security fix slice)

Generated: 2026-11-20. HEAD = `1b7f9c6`; directive = v1.71; two CI runs, both success: the main commit `a464a44` -> `44585713851` (both legs success), and the correction commit `1b7f9c6` -> `44589887023` (observation run).

**Supersession**: this version supersedes v65 (`docs/t027/CHECKLIST_v65{,_EN}.md` is frozen and is not retro-edited).

**Structure**: §1 anchors -> §2 fact package -> §3 completed slices -> §4 execution queue -> §5 observation items -> §6 pending rulings -> §7 / §8 / §9 pointer sections -> §10 opening -> §11 one-line status; the rule bodies behind §7-§9 have moved into directive §13, and this snapshot does not copy them.

**Changes in this version (relative to v65)**:
- (1) §1 anchor rollover: HEAD `1b7f9c6`; directive **v1.71** (unchanged); SELFTEST **740/740**; the rolling window drops v63 and admits v66 (v64 / v65 / v66), and the line count of v65 is re-checked by this version.
- (2) §2 becomes the P79 fact pack (all three faces FINDINGS with F1-F8, the (7) 1+3 user ruling versus 1+2 actually used, R = 2 and D = 1, the six user rulings, and the CI of both commits).
- (3) §3 appends one P79 row (the eighteen rows of v65 plus one = nineteen rows).
- (4) §4 queue update: P79 leaves the queue; next = **P83 `SECURITY-FINDINGS-REMEDIATION`**; after that = **pre-S2 re-evaluation (S-G1 / S-G8)** (moved up from the "post-S2 re-evaluations" line); that line narrows to S-G5 step two plus S-G7; the layer-3 `GATES-PERF-ASSESS` slot stays.
- (5) §5: no new OB or ADR in this version; 13 observing, 9 pending discussion, 1 to do (23 non-terminal, agreeing with the repository re-check).
- (6) §6 is cleared; the next ruling item is the start timing of P83.
- (7) Appendix A rolls over: "obligation transfer" becomes **a new transfer** (F1-F8 of P79 are taken over by P83); §10 / §11 roll over; the cross-session handoff pointer in `DOCUMENTATION_PLAN{,_EN}.md` now points at v66.

---

## §1 Current anchors (check these eight numbers first)

**Rolling window = the last three versions (v64 / v65 / v66)**; v63 has been dropped (same rule that dropped v52).

| Item | Value | Check command |
|---|---|---|
| HEAD | `1b7f9c6` | `git rev-parse --short HEAD` |
| origin/main | `1b7f9c6` | `git rev-parse --short origin/main` |
| Directive version | **v1.71** | the header of `docs/DELIVERY_DIRECTIVE.md` |
| SELFTEST | **740/740** | `node tools/gates/gate.js --selftest` |
| Three scopes | **all TOTAL_FAIL=0** | `gate.js --all --scope=index` / `tree` / `ci` |
| v64 total lines | **224** (LF count, CN / EN symmetric; re-checked by v65) | `docs/t027/CHECKLIST_v64{,_EN}.md` by LF |
| v65 total lines | **220** (LF count, CN / EN symmetric; **re-checked by this version**) | `docs/t027/CHECKLIST_v65{,_EN}.md` by LF |
| v66 total lines | **re-checked by v67** (HEAD advances once this version lands) | `docs/t027/CHECKLIST_v66{,_EN}.md` by LF |

> This snapshot's HEAD = `1b7f9c6` is the checkpoint after the P79 correction commit (the D-value correction); once v66 lands, HEAD advances, so use `git rev-parse --short HEAD` as the live value.

> Rolling window = the most recent three versions (v64 / v65 / v66); v63 has dropped out (the same rule that dropped v52); this snapshot's own line count is re-checked by v67. If two or more items disagree, trust the repository and report the terminal reading before continuing.

`tmp` fixed assertion: top level = `.gitkeep` + `gates-cache`; `gates-cache` holds only `g4b`; the P79 process artifacts were removed item by item during (8)c and left no residue.

---

## §2 P79 fact pack

- **Slice / size**: `S2-PRE-INDEPENDENT-SECURITY-REVIEW` (P79; the ADR-047 compensation duty; assessment `[SLICE]`) / **M-L**; controller `DeepSeek V4.1 Flash (deepseek)`; the user named the start on 2026-11-19 and authorised the later turns in the same rounds.
- **Goal**: complete one independent security review before S2 starts and register the conclusion; three facets = threat model / hard security gate / secrets and credentials (carried word for word from the ADR-047 "uncovered scope", no facet added or removed); scope = the S0-to-S1 security gate in `PRODUCT_REQUIREMENTS` §7 (historical state = never executed).
- **Verdict**: **closed**; all three facets report **`FINDINGS`**, the (7) first round returned `FINAL REVIEW: FINDINGS`, **F1-F8 (3 High + 5 Medium)** in total, and the finding list is **handed over to a fix slice** (user ruling R3 of 2026-11-20); the **ADR-047 status cell now reads `COMPENSATED`** (ruling R1); this slice **changed no product code** and fixed none of the findings.
- **Per-facet findings** (verbatim in report sections 7 / 9): facet 1 threat model, 1 Medium (artifact reads are not bound to a directory handle); facet 2 hard security gate, 2 High (no re-read of revocation / reservation before launch; policy approval not from an independent source) plus 2 Medium (lease evidence derived from the run id only; hook side-effect classification missing); facet 3 secrets and credentials, 1 High (credential / recipient approval binding is not checked first) plus 2 Medium (the Windows native environment block is not zeroed; child-process logs and summaries are not redacted).
- **Review chain**: (1) design (`DeepSeek V4 Pro` Max) -> batch assistance `GLM-5.3-Flash (glm)` **twice** (strictly serial) -> (5) r1/r2 (budget 1+1 **exhausted**) -> (6) r1/r2 (`MAI Code 2 Flash (copilot)`, budget **2 of 3**, r2 `INDEPENDENT SCAN: PASS`) -> (7) final review plus **2 re-reviews** (`GPT-6.2 Sol (copilot)` Max; the third was not used) -> ⑧a / ⑨ / ⑧b / ⑧c -> ⑩.
- **Budget**: the (7) budget follows the user ruling of 2026-11-19 for the **L row, 1 final review + 3 re-reviews** (**it exceeds the §7.5 literal standard "M-L counts as M", is a direct user ruling and may not be cited as a precedent**); actual **1 + 2 (one left)**; probes **0** (no real external call); **no blind review** (none of the three mandatory semantics is touched and the proportional window does not trigger).
- **Metering**: **R = 2** (one stage-5 fix re-run plus one stage-6 fix re-run; the stage-7 re-review round is not counted under §9.2); **the `R>=2` high-risk marker of §9.2 was not set explicitly**, recorded faithfully as a deviation and ratified by the user on 2026-11-20 (no rework, no further R/D counting, a user ruling that may not be cited as a precedent); **D = 1** (after (8)c the controller mistakenly created `tmp/P79-fix/`, which is not a slice directory, and removed the whole directory at once on discovery with zero residue; the self-heal grade holds and the user confirmed it).
- **User rulings (six, 2026-11-20, each marked "a user ruling, not usable as a precedent")**: A ratification of the R = 2 deviation (no rework, no extra counting); B confirmation of the D = 1 self-heal; C authorisation of the ledger D-correction commit; R1 flipping ADR-047 to `COMPENSATED`; R2 the additive note in `PRODUCT_REQUIREMENTS{,_EN}` §7; R3 handing F1-F8 to a fix slice.
- **CI (both commits success)**: main commit `a464a44` -> run **`44585713851`** (28 files total; carrying the report pair whose line count is 393); correction commit `1b7f9c6` -> run **`44589887023`** (4 files: the ledger pair `+1/-1` for the D-value correction plus the report pair `+3/-0` for the ratification notes; **an observation run, and per §11.2 no further commit is made for it**).
- **Documents**: the two-note block in `SECURITY{,_EN}` (review completed + 2026-11-19 + reviewer = the (7) deep tier + conclusion `FINDINGS`; line 7 gains the "before P79" historical qualifier); the additive note in `PRODUCT_REQUIREMENTS{,_EN}` §7 (the original criterion and its blocking force are unchanged); `ADR_REGISTER{,_EN}` ADR-047 = `COMPENSATED` plus the P79 completion note; the ledger P79 segment plus its completion record (including the D-value correction).
- **Report pair terminal state**: the line count is **393** (LF count, **CN / EN strictly co-located**, all six shapes equal), including the (7) verbatim transcription in section 7 and the R/D plus two ratification notes in section 21.
- **Evidence pack**: `docs/validation/evidence/P79-S2-PRE-INDEPENDENT-SECURITY-REVIEW/` (**18 files**, including `MANIFEST.md` plus `SHA256SUMS.txt`; the three entry sets agree and recomputing every sha256 yields zero mismatches).
- **Mechanical readings**: `SELFTEST PASS 740/740`; max ADR = **47**; max OB = **139** (this slice added no OB or ADR); max DR = 14; pilot deferral = round 25 (this slice added no deferral registration).
- **Gates**: the `index` / `tree` / `ci` scopes all report `TOTAL_FAIL=0`; the changed-line counts of the report pair and the ledger pair are symmetric; the G1-b version metadata `v1.71 / 2026-11-19` agrees across both languages (this slice did not change the directive).

### Appendix A: "obligation transfer" (this version: **a new transfer**)

- **Historical hand-off** (first traced back to P75 O-7): the in-cell caliber mismatch in the core task chain row at `s1-exit.md:52` was discharged by P76 correction 1 (the AT-23 fragment) plus the S-G2 correction (`0f80c1e`).
- **Hand-off in this version**: **a new transfer** - the **F1-F8 (3 High + 5 Medium) fix obligation of P79 is taken over by P83 `SECURITY-FINDINGS-REMEDIATION`** (basis = user ruling R3 of 2026-11-20; admission condition = the user must name the start, and **the two High items of facet 2 come first**: re-reading revocation / reservation before launch, and an independent source for policy approval).
- **Convention note**: this treatment has been used since v59; whenever a slice hands an obligation to a future slice, that slice's fact package must carry a hand-off section and the landing entry in §4 must carry an "admission condition" face.

### Appendix B: the "AT-23 stage-2 status" section (unchanged)

- **Status**: ruled "not applicable" by A1 (the user ruling of 2026-11-09; this is **not a permanent prohibition**, unless the user later rules explicitly to restart it).
- **What can be executed before S2**: nothing (R-3 = (A), and the real CLI end-to-end run needs real enforcement, which is an S2 premise).
- **Restart path**: if the user later rules to restart stage 2, real enforcement must land first (S2) or an alternative path must exist (P-prime, a candidate swap, or a new intermediate path); an alternative path means "reach an equivalent goal before S2 by another route", not "implement stage 2 before S2".
- **Natural restart point**: S2 starts, enforcement lands, S-G8 unlocks, R-3 is re-judged (A possibly becoming B), and if B is feasible that is the natural restart point for stage 2.

---

## §3 Completed slices (recent baseline)

**v63 §3 already carried the sixteen rows P64-P80; v64 appended P81 (seventeen rows); v65 appended P82 (eighteen rows); v66 reuses those eighteen rows (readings copied from the v65 text, not retro-edited) and appends one P79 row. Nineteen rows in total.**

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
| v1.71 | **P79 S2-PRE-INDEPENDENT-SECURITY-REVIEW** | **`a464a44` + `1b7f9c6`** | **closed (assessment `[SLICE]`, M-L; all three facets report `FINDINGS`, F1-F8 (3 High + 5 Medium) handed to the fix slice; ADR-047 flips to `COMPENSATED`; **R = 2, D = 1** (the missing high-risk marker at R=2 was ratified by the user, and the D=1 self-heal was confirmed by the user); stage 5 used 2 of 2, stage 6 used 2 of 3, stage 7 used 1 final review + 2 re-reviews (one left); the directive is untouched (v1.71); runs `44585713851` + `44589887023`)** |

> **The D convention in §3** (as in v61): "D = n" in this table is shorthand for the temp-file discipline violation count (§5.1); P79 totals **D = 1** (after (8)c a `tmp/P79-fix/` directory was created by mistake, deleted at once with zero residue, the self-heal grade holds and was confirmed by the user on 2026-11-20). **The D count and the §9.2 R count are two independent series**: P79's R = 2 comes from two stage fix re-runs, and the missing high-risk marker was ratified by the user rather than re-counted.

---

## §4 Execution queue

### Layer 0 preliminary triage (✅ all dequeued, nothing new in this version)

- P5 and P6 were closed by P57; layer 0 is fully dequeued with nothing left.

### Layer 1 main-line hard gates

- Items 8-12: ✅ all dequeued (see §4 of v57 / v58; **the original dequeue records have been pushed into earlier versions by the rolling window, so trace them back layer by layer through the citation chain in §4 of v57**).
- **`T027-B` (the defensive-refactor main-line milestone) = complete (partly not applicable)**: C1-C7 / N1-N4 / R-1 / R-2 / R-3 are all dispositioned (R-3 = (A), ruled by the user on 2026-11-15; the A1 declaration caliber is "complete, partly not applicable"). **The main line is fully dequeued.**
- **`S1 Exit close-out` (gate level complete)**: **P76 produced the roadmap; P77 closed S-G4; P78 closed S-G3; P80 decided that the S-G6 gate-level side holds**. The eight gates are listed below (after P80 closes; P79 did not change the gate table):

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

**The independent security review (not an S-G5 criterion, but a precondition of the second approval step)**: the S0-to-S1 gate in `PRODUCT_REQUIREMENTS` §7, whose historical state was "never executed"; **P79 completed it on 2026-11-19** (conclusion = `FINDINGS`, F1-F8 handed to the fix slice; ADR-047 now reads `COMPENSATED`). **The Q10=(C) two-step S-G5 approval distinguishes them**: **step one**, inside S1 and covering the release portion, does not require the review to be complete, whereas **step two**, after S2 and covering the whole S1 Exit, requires the review to have completed before S2 starts, and **that precondition has been discharged**. **Hard boundary: no approval may be read as S1 Exit passing, and a FINDINGS conclusion may not be read as the gate passing.**

- **S1 Exit cannot pass as a whole inside S1** (the premise of Q2=(a) plus Q10=(C)). **Every advanceable gate-level item is closed** (S-G2 / S-G3 / S-G4 / S-G5 and the S-G6 gate level; S-G7 is a composite gate, S-G8 needs S2, S-G1 does not pass).
- **`TMP-DISCIPLINE-OUT-OF-REPO` (the §5.1 temp-file discipline · P82)**: **P82 is closed** (main commit `0264bcd`; follow-ups `adb8cb9` / `1015430` / `9d83089` / `9fd8cd6`; R = 2, D = 3 (which fired the stop-and-escalate); **directive v1.71**). Removed from the queue; the fact pack is in v65 §2; the report is `docs/validation/tmp-discipline-out-of-repo{,_EN}.md`.

- **`S2-PRE-INDEPENDENT-SECURITY-REVIEW` (the ADR-047 compensation duty · P79)**: **P79 is closed** (main commit `a464a44`; correction commit `1b7f9c6`; R = 2, D = 1; all three facets `FINDINGS`, F1-F8 handed to the fix slice; ADR-047 flips to `COMPENSATED`). Removed from the queue; the fact pack is in §2; the report is `docs/validation/s2-pre-independent-security-review{,_EN}.md`.

**Next item = P83 `SECURITY-FINDINGS-REMEDIATION` (fixes F1-F8 of P79; implementation type, touching product code plus security semantics; the start requires the user to name it)**

| Field | Content |
|---|---|
| Identifier | `SECURITY-FINDINGS-REMEDIATION` (ledger number = **P83**; basis = user ruling R3 of 2026-11-20; source = F1-F8 in P79 report section 9) |
| Contents | fix the eight P79 findings (3 High + 5 Medium); **priority = the two High items of facet 2 first**: `internal/console/runtime_agent.go:1245` (re-read revocation / reservation before launch) and `internal/chain/review.go:116` (independent source for policy approval); then the facet-2 Medium items, then facet 3, then facet 1 |
| Size tier | **M-L** (implementation type; touches product code plus security semantics plus contract wording, so the protocol-first test must be made) |
| (7) budget | per §7.5: when several rows apply, take the one with the most re-reviews (governance / high-risk row = 1 + 3); **decide and register against the repository when the slice starts** |
| Dependency | P79 is closed (F1-F8 registered and archived verbatim); **the start requires the user to name it**; the admission condition is in appendix A |
| Boundary | fix only F1-F8 and the negative tests that serve them; do not change the eight S1 Exit criteria; do not change the `SECURITY{,_EN}` threat-model baseline; do not widen to security facets P79 did not report |
| Start condition | **the user must name it explicitly** (the name is fixed and the priority is written; the start timing is the user's call) |

**After that = pre-S2 re-evaluation (S-G1 / S-G8)** (premise = S2 enforcement lands; the independent-review precondition has been discharged by P79, so the ADR-047 rollback clause no longer blocks it; **this item moves up from the v65 "post-S2 re-evaluations" line into this queue**).

- **`S-G8`** stays blocked until S2 enforcement lands (P73 A2 / P75 A2 / ADR-044 / ADR-045); the gate table in `s1-exit{,_EN}.md` is untouched (P76 only edited the fragment inside row `:52`, and P79 changed nothing there).
- **Post-S2 re-evaluations (narrowed in this version)**: **S-G5 step two (whole S1 Exit approval) plus S-G7**, in the order given by P80 report section 7 (S-G1 / S-G8 have moved up to the pre-S2 line).

### Layer 2 routine governance

- P72 / P73 / P74 / P75 / P76 / P77 / P78 / P80 / P81 / P82 / P79 ✅ removed from the queue (tier 2 is **empty**; per-row summaries are in §3).
- The full §1.2 field definitions for P72 / P73 / P74 / P75 / P76 / P77 / P78 / P80 / P81 / P82 / P79 are in the corresponding segments of `docs/t027/REMAINING_SLICES{,_EN}.md`.

### Layer 3 low-priority options

| Field | Content |
|---|---|
| Identifier | `OB-50` / `spawnGit` / `ci freeze cost` / 22 historical `%TEMP%` items / `GATES-EXEMPT` / `GATES-PERF-ASSESS` |
| Current state (verified item by item) | OB-50 = **settled** (§12.4, v1.41); spawnGit = the git-call helper of `tools/gates/selftest.js` (the spawn wrapper in `lib/ctx.js`) that folds stderr into the result object instead of printing it; the P27 `G4b-CACHE` report §5 item 5 registered it as a "candidate OB (diagnostic)" (Low / observing) and it was never formally registered as a §12.4 OB id, so under the §1.8.1 pending-item rule it is an "unregistered candidate OB"; candidate fix = attach the last N bytes of stderr to `detail` when the spawn fails (diagnostics only, the verdict is unchanged); trigger = the next unexplained git-related assertion failure is assessed by (1); ci freeze cost = **no §12.4 entry**, its real origin is item 7 of the P27 report as a "candidate OB (performance)", **not formally registered**; the 22 historical %TEMP% items map to OB-98 (**settled**) / OB-99 (**observing**) / OB-100 (**observing**); GATES-EXEMPT = **not initiated** (no §12.4 entry; the P33 item 1-3(iv) split slice, touching four criterion files); GATES-PERF-ASSESS = **a new candidate** (registered by the user ruling of 2026-11-19, **not initiated as an implementation**; a gate-cost composition assessment that measures how `SELFTEST` and the three scopes spend their time and how often the `base` trio repeats; an assessment-grade candidate that needs the user to name it at start) |
| Contents | see the matching §12.4 rows for the six items (`GATES-PERF-ASSESS` has no §12.4 entry and is registered only in this table and the ledger); `GATES-EXEMPT` is the P33 item 1-3(iv) split slice, not initiated |
| Dependency | the long suspension covers only unsettled items (spawnGit / ci freeze cost / the %TEMP% timing shapes (OB-99 / OB-100) / GATES-EXEMPT); `GATES-PERF-ASSESS` has no prerequisite and ranks after P83; none of them has a concrete trigger |
| Size | varies (S through M-L) |
| (7) budget | none |
| Boundary | per item in §12.4; this version advances no suspended item and keeps the slot registered |
| Start condition | the user must name the item explicitly |

**Queue discipline**: as in v60 (layer boundaries and numbering discipline; this slice runs in a single batch and does not bump the version).

---

## §5 Observation items

- Under observation (**13**): `OB-1 / 7 / 10 / 33 / 87 / 99 / 100 / 104 / 105 / 107 / 111 / 113 / 135` (carried from v61).
  - **OB-139 is resolved and discharged (P82 closed, v1.71)**: §5.1 adds the "out-of-repo same-command create-and-delete self-heal sub-clause" and fixes the per-element caliber of `tmp/`-root self-heal; **A-6 forms its seventh literal trigger** (fact only, no waiver, not a precedent).
  - **The reinforced D=2 recurrence-prevention sentence (standard kick-off content since v59, carried into v61)**: "the executing roles (`prfrail-implementer` / `prfrail-tester` / `prfrail-documenter`) **and the controller** must always work under a `tmp/` slice directory (`tmp/<slice-id>/`); **the same holds for CI observation captures, which must never land at the `tmp/` root**; three D occurrences inside one slice force the slice to stop."
- Pending discussion (9): `OB-19 / 21 / 52 / 53 / 57 / 58 / 59 / 61 / 75`.
- To do (1): `OB-69`. Non-terminal total **23** (13 plus 9 plus 1), agreeing with the repository re-check of this version.
- **Slice-level suggestions (three, not registered as OBs)**:
  1. **`R-ACCUMULATION-MECHANISM` (carried from v58)**: P74 was the first case to trigger the §9.2 rule (R>=2 means high risk) together with the §7.5 high-risk tier (one plus three); **P79's R = 2 hit that threshold again while the marker was not set** (ratified by the user), which raises the urgency of this suggestion; later implementation slices should pre-place the dynamic budget sentence in their kick-off message **and require the marker to be set on the spot**. **No OB is registered now**.
  2. **`--progress` covers state transitions only, not log output (carried from v58)**: if long-running observability becomes a real need, evaluate adding `--progress verbose` or `--tail-logs`.
  3. **The "obligation hand-off" convention (in use since v59, first applied to P75 O-7)**: this version **actually uses it** (the P79 to P83 transfer in appendix A). **No OB is registered now**.
- **Removed in this version**: none.
- **Flipped to `settled` in this version**: none (P79 added or flipped no OB).
- Mechanical readings: max OB in §12.4 = **139**; **max ADR = 47**; max DR = 14; max pilot deferral = round 25.

---

## §6 Pending user rulings

- **The table is currently empty** (the six P79 rulings were made on 2026-11-20, see the last bullet of §2).
- **Next ruling item (expected)**: **the start timing of P83** (the name is fixed and the priority is written; the start requires the user to name it).
- **Six P79 rulings**: A ratification of the R = 2 deviation; B confirmation of the D = 1 self-heal; C authorisation of the ledger D-correction commit; R1 flipping ADR-047 to `COMPENSATED`; R2 the additive note in §7; R3 handing F1-F8 to a fix slice; all six are marked "made directly by the user, not usable as a precedent".

---

## §7 / §8 / §9 Pointer sections

Same structure as v60: §7 (kick-off and step-0 general rules) points at directive §13.1; §8 (cross-slice collaboration discipline) at §13.2; §9 (checklist meta-rules and gate pitfalls) at §13.3 / §13.4. **The rule bodies have moved into directive §13 and are not copied here.**

---

## §10 One-line opening for a new session

> This is my project state snapshot (v66); check the eight anchors in section 1 first, then continue per the section 4 queue, and if anything disagrees trust the repository.
>
> [paste this checklist in full v66]
>
> [supporting documents] `REMAINING_SLICES{,_EN}.md` (P79 closed, P83 awaiting the start), `DELIVERY_DIRECTIVE.md` (**v1.71**), `s2-pre-independent-security-review{,_EN}.md` (the P79 report, with F1-F8 and the R/D in section 21), `tmp-discipline-out-of-repo{,_EN}.md` (P82 report), `directive-model-upgrade{,_EN}.md` (P81 report), `sg6-gatelevel-consolidate{,_EN}.md` (P80 report), `release-zip-binding{,_EN}.md` (P78 report), `ADR_REGISTER{,_EN}.md` (with ADR-045 / ADR-046 / ADR-047 = `COMPENSATED`), `DOCUMENTATION_PLAN{,_EN}.md` (cross-session handoff entry pointer), `CHECKLIST_v65{,_EN}.md` (frozen, superseded by v66).
>
> Current: P79 is closed (all three facets FINDINGS, F1-F8 handed to the fix slice; R = 2, D = 1; ADR-047 reads `COMPENSATED`); next = P83 `SECURITY-FINDINGS-REMEDIATION` (the start requires the user to name it, with the two facet-2 High items first).

## §11 One-line current state
P79 closed (main commit `a464a44` plus the correction commit `1b7f9c6` pushed to origin, both runs success; all three facets `FINDINGS`, F1-F8 (3 High + 5 Medium) handed to P83; **R = 2, D = 1**; stage 7 used 1 final review + 2 re-reviews (one left); ADR-047 = `COMPENSATED`; the directive stays at v1.71 and SELFTEST 740/740); every advanceable tier-1 gate-level item stays closed; next = **P83 (the start requires the user to name it)**, and after that = **pre-S2 re-evaluation (S-G1 / S-G8)**.
