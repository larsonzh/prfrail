# Task execution checklist (cross-session snapshot v67, after the partial P83 closure; next item = P84 start)

Generated: 2026-11-21. HEAD = `4323503`; directive = v1.71; three CI runs: first red `44795140136` (`ec532a4`), fix green `44801690679` (`c820743`), observation green `44805856028` (`4323503`).

**Supersedes**: this version supersedes v66 (`docs/t027/CHECKLIST_v66{,_EN}.md` is frozen and not edited retroactively).

**Structure**: section 1 anchors, section 2 fact pack, section 3 completed slices, section 4 queue, section 5 observations, section 6 pending rulings, sections 7 to 9 pointer sections, section 10 opening line, section 11 one-line status; the rule bodies of sections 7 to 9 moved into directive section 13 and are not copied here.

**Changes against v66**:

- (1) Section 1 anchors rotated: HEAD `4323503`; directive **v1.71** (unchanged); SELFTEST **740/740**; the rolling window drops v64 and takes in v67 (v65 / v66 / v67), and v66's line count is re-verified by this version.
- (2) Section 2 becomes the P83 fact pack (`PARTIALLY CLOSED`, three commits and three CI runs, six (7) rounds, R = 2, D = 0, DR-15).
- (3) Section 3 appends one P83 row (nineteen rows from v66 become twenty).
- (4) Section 4 queue update: P83 leaves the queue; next = **P84 `BOUNDDIR-RESIDUAL`**; then = **P85 `SECURITY-PROOF-SUBSYSTEM`**, then the pre-S2 re-evaluation (S-G1 / S-G8), then the post-S2 re-evaluation items; the tier-3 `GATES-PERF-ASSESS` slot stays.
- (5) Section 5 observations: no new OB or ADR this version; 13 observing, 9 under discussion, 1 pending (23 non-terminal); mechanical reading DR max = 15.
- (6) Section 6 pending rulings cleared; the next ruling item is the P84 start time.
- (7) Appendix A rotated: the duty hand-over is now **P83 residuals to P84 / P85**; sections 10 and 11 follow.

---

## Section 1 current anchors (check these 7 numbers first)

**Rolling window = the latest 3 versions (v65 / v66 / v67)**; v64 has been squeezed out under the same rule as v52.

| Item | Value | Check command |
|---|---|---|
| HEAD | `4323503` | `git rev-parse --short HEAD` |
| origin/main | `4323503` | `git rev-parse --short origin/main` |
| Directive version | **v1.71** | header of `docs/DELIVERY_DIRECTIVE.md` |
| SELFTEST | **740/740** | `node tools/gates/gate.js --selftest` |
| Three scopes | **all TOTAL_FAIL=0** | `gate.js --all --scope=index` / `tree` / `ci` |
| v65 whole-file lines | **220** (counted by LF, CN and EN symmetric; verified by v66) | LF count of `docs/t027/CHECKLIST_v65{,_EN}.md` |
| v66 whole-file lines | **216** (counted by LF, CN and EN symmetric; **re-verified by this version**) | LF count of `docs/t027/CHECKLIST_v66{,_EN}.md` |

> The HEAD of this snapshot, `4323503`, is the anchor after the P83 backfill commit; v67 itself moves HEAD forward once it lands, so `git rev-parse --short HEAD` is authoritative then. v67's own line count is verified by v68.

> Rolling window = the latest 3 versions (v65 / v66 / v67); v64 has been squeezed out; if 2 or more items disagree, the repository wins: report the terminal readings first and continue.

Fixed `tmp` assertion: the top level is `.gitkeep` plus `gates-cache`; `gates-cache` holds only `g4b`; P83's interim files were deleted item by item under (8)c with no residue.

---

## Section 2 P83 fact pack

- **Slice / size**: `SECURITY-FINDINGS-REMEDIATION` (P83; remediates F1 to F8 from P79) / **M-L**; controller `DeepSeek V4.1 Flash (deepseek)`; the user named the start on 2026-11-20 and ruled twice mid-slice on triage and the extra re-review.
- **Goal**: remediate F1 to F8 (3 High plus 5 Medium) from the P79 report section 9 verbatim, with one falsifiable closure assertion plus at least one counterexample regression test per finding.
- **Verdict**: **`PARTIALLY CLOSED`** — **F2 / F5 / F7 plus the T9 narrowing are closed**; **F1 / F4 are partially implemented** (residuals deferred to P84); **F3 / F6 / F8 are deferred to P84 / P85**. **It must not be read as "closed" or "security-closed".**
- **Deferral sources**: T2 STOP (AgentRunner has no persistable, independently re-verifiable frozen policy artifact) to P85; T3 residuals (the Windows name-to-image binding window, Linux same-inode not meaning byte-immutable) to P85; T8a (refuse BYOK when on-disk redaction cannot be proven) to P85; T4 residuals (bounddir case folding and intermediate reparse, root relative-path semantics) to P84; T5#4 residuals (batch predecessor validation, durable proof recovery) to P84.
- **Contract impact**: all eight and the T9 narrowing were judged **DOES-NOT-TOUCH**; `CONTRACTS{,_EN}`, schemas and fixtures are unchanged.
- **Review chain**: (1) plan plus contract triage (`DeepSeek V4 Pro` Max, 2 calls), then (2) eight implementation rounds and (3) seven test rounds (`DeepSeek V4.1 Flash`), (5) 2/2 exhausted (r1 FINDINGS then r2 PASS), (6) 2/3 (r1 FINDINGS then r2 PASS), (7) six rounds of `GPT-6.2 Sol (copilot)` Max (blind 1, final 1, re-review 3, additional 1; the additional round is a direct user ruling and may not be cited as a precedent), then (8)a, (9), (8)b, (8)c and (10).
- **Accounting**: **R = 2** (one (5)-ring remediation re-run plus one (6)-ring remediation re-run); **D = 0**; probes **0**; anomaly **E-1** (an unstaged change in `CHECKLIST_v66.md` whose origin could not be attributed) was reverted per the user ruling and is not counted as D.
- **(5) / (6) caveat**: both `PASS` verdicts bind the **pre-remediation snapshot**; the final snapshot is covered only by (7), whose final review and every re-review returned `FINDINGS`.
- **Defect register**: **DR-15** (`TestP83Fix6T5cCorruptRetainedGenerationBlocksTheRenewal` flaked once locally; not reproduced by re-runs or CI; registered separately under DR boundary 1).
- **Three commits and CI**: main `ec532a4` (`+12245/-283`) produced run **`44795140136`**, **red on both legs** (ubuntu `Test` 5 cases; windows `Gates` reported `base go test` ok=14 FAIL=3); fix commit `c820743` (test files only, `+40/-14`) produced run **`44801690679`**, **success on both legs** (`PASS base go test ./... :: ok=15 FAIL=0`, `SELFTEST 740/740`, `TOTAL_FAIL=0`); backfill commit `4323503` (`+26/-2`) produced observation run **`44805856028`**, success.
- **VM evidence**: Ubuntu `10.0.0.199` (go1.22.12): the three targeted packages `ok`, the full 15 packages `ok`, the five cases `PASS` in verbose mode (no SKIP), and the targeted `-race` run `ok`.
- **Docs**: report pair `docs/validation/security-findings-remediation{,_EN}.md` (strictly mirrored, counted by LF); the `SECURITY{,_EN}` line-5 remediation note; the ledger P83 section and completion record; ADR-047 stays `COMPENSATED`; the evidence pack holds 8 payloads plus MANIFEST and SHA256SUMS (re-hash mismatch 0).

### Appendix A: duty hand-over (this version: **new hand-over**)

- **Historical hand-over** (the first instance traces to P75 O-7): the inconsistent inline wording at `s1-exit.md:52` was discharged by P76 change 1 plus the S-G2 correction (`0f80c1e`).
- **This version's hand-over**: **new hand-over** — P83's five residuals are carried by **P84 (T4 / T5#4)** and **P85 (T2 / T3 / T8a)**; both need a user-named start, and P85 must change the protocol first.
- **Convention note**: in use since v59; the landing slice's entry must carry an admission-condition facet (the "Status" and "Boundaries" fields of the P84 / P85 sections already state it).

### Appendix B: AT-23 stage 2 status (unchanged)

- **Status**: the A1 ruling is "not applicable" (user ruling of 2026-11-09; not a permanent ban).
- **Executable object before S2**: none (R-3 = (A)).
- **Restart path**: if the user later rules to restart stage 2, real enforcement must land first (S2) or an alternative path must exist.
- **Natural restart point**: S2 starts, enforcement lands, S-G8 unlocks, R-3 is re-evaluated.

---

## Section 3 completed slices (recent baseline)

**v63's section 3 already held sixteen rows (P64 to P80); v64 appended P81 (seventeen); v65 appended P82 (eighteen); v66 kept v65's eighteen and appended P79 (nineteen); v67 keeps v66's nineteen (readings copied, not rewritten) and appends one P83 row. Twenty rows in total.**

| Version | Slice | Main commits | Verdict |
|---|---|---|---|
| v1.66 | P64 AT24-SUPERVISED-BLACK-BOX | `2dbc54a` | closed (AT-24 conditions closed; R = 2; run `42948727363` green first time) |
| v1.66 | P65 GATE-PIT-26-LANDING | `0eda074` | closed (appendix B.7.1 slot plus A-part ledger to-do; R = 1; run `42969802819`) |
| v1.66 | P66 CHECKLIST-RULES-HOIST | `1650d5f` | closed (section 13 five parts plus the new snapshot structure; R = 3, D = 2; run `43022796716`) |
| v1.67 | P67 GATES-DEBT-CLEANUP | `18b6ef9` + `365a930` | closed (g4b LRU, gate pitfalls 27 to 30, the 12.4 settlement block, B.6 hardening; R = 3, D = 1; run `43124686702`) |
| v1.67 | P68 REMAINING-RULINGS | `e0d745c` + `665a31f` | closed (assessment type; R = 1, D = 1; (7) PASS; five rulings made; ADR-042; run `43177025096`) |
| v1.67 | P69 T027B-E2E-WIRING | `759c934` | closed (implementation type; R = 0, D = 1; (7) PASS; E-1 recorded; run `43264324421`) |
| v1.68 | P70 S1-STREAM-PROGRESS | `ca88fa3` | closed (implementation type; R = 1, D = 2; (7) PASS; DR-14; runs `43400311736` then `43403891205` then `43408438643`) |
| v1.68 | P71 S1-SCOPE-ANNOTATE | `6ee2c61` + `16c98f7` | closed (documentation type; R = 0, D = 2; (7) blind PASS plus final PASS; runs `43463089586` and `43468806314`) |
| v1.68 | P72 REQ020-S2-RESTRICTION | `fb0ab57` + `406f4a0` | closed (documentation type; R = 0, D = 0; (6) skipped with (7) confirmation; (7) final PASS WITH FIXES; runs `43566904684` and `43568186355`) |
| v1.69 | P73 T027B-ADMISSION | `8c59712` + `f834ce5` | closed (assessment type; R = 0, D = 0; (6) skipped with (7) confirmation; (7) final PASS; **directive raised to v1.69**; runs `43622770750` and `43624661784`; OB-55 = option b) |
| v1.69 | P74 SANDBOX-INTERFACE | `3300364` + `150423d` | closed (implementation type; R = 3, D = 0; (6) mandatory PASS; (7) final PASS; **first R=3 case to trigger the 9.2 high-risk rule**; runs `43669748483` and `43672994996`) |
| v1.69 | P75 R3-REAL-CLI-ASSESSMENT | `ed2e26e` + `70a04cd` + `60146c0` | closed (assessment type; R = 1, D = 1; (7) blind 1 plus final 1; **R-3 = (A)**; ADR-045; runs `43732384194`, `43734829063`, `43739409108`) |
| v1.69 | P76 S1-EXIT-CLOSEOUT-ASSESSMENT | `77f23d2` + `aca19af` + `0f80c1e` | closed (assessment type; R = 1, D = 1; (7) final PASS WITH FIXES; Q2=(a), Q10=(C); ADR-046 and ADR-047; runs `43785603244`, `43787083404`, `43793273520`) |
| v1.69 | P77 ADR010-APPROVAL | `33272be` + `fc00c4b` | closed (documentation type; R = 1, D = 0; ADR-010 = `APPROVED`; S-G4 passed; runs `43986942087` and `43997058549`) |
| v1.69 | P78 RELEASE-ZIP-BINDING | `c935713` + `0741fae` + `308b235` | closed (implementation type; R = 0, D = 0; (7) final FINDINGS then re-review PASS; GitHub Release `v0.1.0` bound; S-G3 passed; runs `44043705691` and `44060224254`) |
| v1.69 | P80 SG6-GATELEVEL-CONSOLIDATE | `e00d79b` + `58c2794` | closed (documentation type; R = 2, D = 1; (7) narrowed read `RE-REVIEW: PASS`; S-G6 gate-level side established; runs `44211368890` and `44213070006`) |
| v1.70 | P81 DIRECTIVE-MODEL-UPGRADE | `45590e8` + `fd4d8fb` + `9ca6576` | closed (governance slice; (6) `MAI Code 2 Flash`, (7) `GPT-6.2 Sol`; SELFTEST 737 to 740; R = 1, D = 2; runs `44281812336`, `44284200392`, `44392113028`) |
| v1.71 | P82 TMP-DISCIPLINE-OUT-OF-REPO | `0264bcd` + `adb8cb9` + `1015430` + `9d83089` + `9fd8cd6` | closed (governance slice; section 5.1 gains the out-of-repo self-heal clause; OB-139 flipped to disposed; **directive raised to v1.71**; R = 2, D = 3; five commits, CI all success) |
| v1.71 | P79 S2-PRE-INDEPENDENT-SECURITY-REVIEW | `a464a44` + `1b7f9c6` | closed (assessment type; all three faces `FINDINGS`, F1 to F8 handed to the fix slice; ADR-047 flipped to `COMPENSATED`; R = 2, D = 1; runs `44585713851` and `44589887023`) |
| v1.71 | P83 SECURITY-FINDINGS-REMEDIATION | `ec532a4` + `c820743` + `4323503` | **partially closed** (F2 / F5 / F7 plus the T9 narrowing closed; F1 / F4 partial; F3 / F6 / F8 deferred to P84 / P85; R = 2, D = 0; DR-15; first run red `44795140136`, then `44801690679` both green, observation `44805856028`) |

> **The D column of section 3** (same as v61): "D = n" abbreviates the temporary-file discipline violation count (section 5.1); P83 totals **D = 0**. **D and the 9.2 R count are two independent sequences.**

---

## Section 4 execution queue

### Tier 0 pre-work (✅ done, nothing new this version)

- P5 and P6 were closed by P57; tier 0 is empty with no residue.

### Tier 1 main-line hard gates

- Items 8 to 12 are all ✅ done (see v57 / v58 section 4; trace back through the v57 section 4 reference chain if needed).
- **`T027-B` (defensive refactor phase) = complete (partly not applicable)**: C1 to C7, N1 to N4, R-1, R-2 and R-3 are all handled (R-3 = (A)). **The main line is fully dequeued.**
- **`S1 Exit close-out` (gate-level side complete)**: P76 produced the close-out roadmap; P77 closed S-G4; P78 closed S-G3; P80 judged the S-G6 gate-level side established; the independent security review was supplied by P79 and partially remediated by P83. The eight gates read as follows (P83 did not change the table):

| Gate | Criterion | Status | Evidence |
|---|---|---|---|
| S-G1 | the CLI can really start an admitted AI candidate and meter it | **not passed** (Q2=(a) literal) | P76 section 4; ADR-046 |
| S-G2 | `supervised-black-box` implemented and AT-24 conditions hold | **passed** (T028 `COMPLETED` plus closed AT-24 conditions) | P76 section 4; `0f80c1e`; ADR-046 |
| S-G3 | the final ZIP is bound to a release record and re-verified | **passed** (P78 closed) | P78 report sections 5 and 7 |
| S-G4 | ADR-010 approved by the product owner | **passed** (user approval 2026-11-16; P77 closed) | P77 report section 5 |
| S-G5 | the product owner approves | **release part approved (first step)** | ADR-046 Q10=(C); P76 section 12 |
| S-G6 | dependency fields satisfied and the four remaining gates closed | **gate-level side established** (P80 combined verdict) | P80 report sections 5 and 6 |
| S-G7 | all ATs pass, T018 complete and section 14 fully exited | unreachable (composite gate) | `s1-exit.md:52` / `:56` to `:58` / `:60` |
| S-G8 | the supervised path runs end to end with a real CLI and a real candidate | blocked until S2 enforcement lands | P75 A2 / ADR-045; P76 section 6 |

**Next = P84 `BOUNDDIR-RESIDUAL` (carries P83's T4 / T5#4 residuals; the user must name the start)**

| Field | Content |
|---|---|
| Identifier | `BOUNDDIR-RESIDUAL` (ledger number **P84**; source = P83 report section 8's P84-d / P84-e plus the user's scheduling ruling of 2026-11-21) |
| Contents | T4 = bounddir case folding, intermediate reparse, root relative-path semantics; T5#4 = batch predecessor validation and a recoverable generation-to-proof transaction |
| Size | **M** (implementation type; platform binding plus transaction semantics) |
| (7) budget | deep tier `GPT-6.2 Sol (copilot)`; 1 final review plus 3 re-reviews; one mandatory blind review (separately counted as +1); **the user's scheduling ruling of 2026-11-21, not a precedent** |
| Dependencies | P83 partially closed (`4323503`); the user must name the start |
| Boundaries | fix the F1 / F4 residuals only; do not touch F3 / F6 / F8; do not change the gate table, directive or CI |
| Start condition | **the user must name it explicitly**, and the start message records the four elements |

**Then = P85 `SECURITY-PROOF-SUBSYSTEM` (carries P83's T2 / T3 / T8a; the user must name the start; protocol first)**

| Field | Content |
|---|---|
| Identifier | `SECURITY-PROOF-SUBSYSTEM` (ledger number **P85**; source = P83 report section 8's P84-a / P84-b / P84-c plus the user's scheduling ruling of 2026-11-21) |
| Contents | T2 = the policy artifact (closes F3); T3 = the immutable install source (closes F6); T8a = the BYOK proof carrier (closes F8) |
| Size | **L** (implementation type; expected to touch `CONTRACTS{,_EN}` plus schemas and fixtures) |
| (7) budget | deep tier `GPT-6.2 Sol (copilot)`; 1 final review plus 3 re-reviews; one mandatory blind review (separately counted as +1); **the user's scheduling ruling of 2026-11-21, not a precedent** |
| Dependencies | P83 partially closed (`4323503`); may run in parallel with P84; the user must name the start |
| Boundaries | close only the deferred parts of F3 / F6 / F8; do not touch F1 / F4; do not change the gate table, directive or CI |
| Start condition | **the user must name it explicitly**; the first step is the contract-impact verdict and protocol-first change |

**Then = the pre-S2 re-evaluation (S-G1 / S-G8)**

- Prerequisite = S2 enforcement lands; the independent security review prerequisite was met by P79 and partially remediated by P83, so ADR-047's rollback clause no longer blocks.
- **`S-G8`** stays blocked until S2 enforcement lands (P73 A2 / P75 A2 / ADR-044 / ADR-045); `s1-exit{,_EN}.md` keeps its gate table unchanged.
- **Post-S2 re-evaluation items**: **S-G5 second step (S1 Exit overall approval) plus S-G7**, following the order in P80 report section 7.

### Tier 2 routine governance

- P72, P73, P74, P75, P76, P77, P78, P80, P81, P82, P79 and P83 are all ✅ dequeued (tier 2 is **empty**).
- Each slice's full 1.2 field definition lives in the matching section of `docs/t027/REMAINING_SLICES{,_EN}.md`.

### Tier 3 low-priority candidates

| Field | Content |
|---|---|
| Identifier | `OB-50`, `spawnGit`, `ci frozen cost`, `%TEMP% 22 legacy items`, `GATES-EXEMPT`, `GATES-PERF-ASSESS` |
| Current status (verified item by item) | OB-50 = **disposed** (section 12.4, v1.41); spawnGit = candidate OB (diagnostic, never formally registered); ci frozen cost = not formally registered (P27 report section 7 candidate); %TEMP% legacy = OB-98 disposed with OB-99 and OB-100 observing; GATES-EXEMPT = **not established**; GATES-PERF-ASSESS = **new candidate** (registered by the user on 2026-11-19, deliberately not implemented) |
| Contents | see the matching section 12.4 rows; `GATES-EXEMPT` is the 1-3(iv) split of P33 |
| Dependencies | the long suspension covers only undisposed items; `GATES-PERF-ASSESS` has no prerequisite and ranks after P84 / P85 |
| Size | varies (S through M-L) |
| (7) budget | none |
| Boundaries | this version advances no suspended item; the registration slots stay |
| Start condition | the user must name it explicitly |

**Queue discipline**: tier boundaries and numbering discipline as in v60; this slice does not bump the version.

---

## Section 5 observations

- Observing (**13**): `OB-1 / 7 / 10 / 33 / 87 / 99 / 100 / 104 / 105 / 107 / 111 / 113 / 135` (carried from v61).
- Under discussion (9): `OB-19 / 21 / 52 / 53 / 57 / 58 / 59 / 61 / 75`.
- Pending (1): `OB-69`. Non-terminal total **23** (13 plus 9 plus 1), matching this version's verification.
- **DR register**: DR max = **15** (v67 adds DR-15: P83's `internal/guard` case `TestP83Fix6T5cCorruptRetainedGenerationBlocksTheRenewal` flaked once locally; registered separately under DR boundary 1 and not blocking the partial closure).
- **(7) budget caveat (P83 carry-over observation)**: of the six (7) rounds (blind 1, final 1, re-review 3, additional 1), the additional round and the 1+3 budget are direct user rulings and may not be cited as precedents; the (5) 2/2 and (6) 2/3 PASS verdicts bind the pre-remediation snapshot.
- Mechanical readings: section 12.4 max OB = **139**; ADR max = **47**; DR max = 15; pilot deferral max = the 25th instance.
- **Removed this version**: none. **Flipped to disposed this version**: none.

---

## Section 6 pending user rulings

- **The table is currently empty** (P83's close-out and scheduling rulings were made on 2026-11-21).
- **Next ruling item (anticipated)**: **the P84 start time** (the slice name and priority are fixed; the user must name the start).

---

## Sections 7 / 8 / 9 pointer sections

Same structure as v60: section 7 start and step-0 general rules lead to directive section 13.1; section 8 cross-slice collaboration discipline leads to 13.2; section 9 checklist meta-rules and gate pitfalls lead to 13.3 / 13.4. **The rule bodies moved into directive section 13 and are not copied here.**

---

## Section 10 one-line opening for a new session

> Here is my project state snapshot (v67); check the seven anchors in section 1 first, continue through the section 4 queue if they match, and trust the repository if they do not.
>
> [paste this whole v67 checklist]
>
> [supporting documents] `REMAINING_SLICES{,_EN}.md` (P84 / P85 established, both awaiting a named start), `DELIVERY_DIRECTIVE.md` (**v1.71**), `security-findings-remediation{,_EN}.md` (the P83 report with the CI record and the deferral list), `SECURITY{,_EN}.md` (line-5 remediation note), `ADR_REGISTER{,_EN}.md` (ADR-047 = `COMPENSATED`), `DOCUMENTATION_PLAN{,_EN}.md` (hand-over entry pointer), `CHECKLIST_v66{,_EN}.md` (frozen, superseded by v67).
>
> Now: P83 partially closed (F2 / F5 / F7 plus the T9 narrowing closed; F1 / F4 partial; F3 / F6 / F8 deferred); next = P84 `BOUNDDIR-RESIDUAL` (the user must name the start).

## Section 11 one-line current status

P83 **partially closed** (three commits `ec532a4` / `c820743` / `4323503` pushed to origin; the first red run `44795140136` is preserved as-is, the fix run `44801690679` is green on both legs, the backfill observation `44805856028` is green; both legs `base go test ok=15 FAIL=0`, `SELFTEST 740/740`, three scopes 0; R = 2, D = 0; DR-15); the tier-1 gate-level items all stay closed; next = **P84 (the user must name the start)**, then **P85**, then the **pre-S2 re-evaluation (S-G1 / S-G8)**.
