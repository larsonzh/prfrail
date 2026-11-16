# Task Execution Checklist (cross-session snapshot v60, after P76 closes, S1 Exit closeout in progress)

Generated: 2026-11-16. HEAD = `0f80c1e`; directive = v1.69; three CI runs: the P76 landing commit `43785603244` (ubuntu 3m24s / windows 5m12s, both legs green), the P76 backfill `43787083404` (**an observation record only, per §11.2**) and the S-G2 correction [ROUTINE] slice main run `43793273520` (ubuntu 3m16s / windows 4m28s, both legs green).

**Supersession**: this version supersedes v59 (`docs/t027/CHECKLIST_v59{,_EN}.md` is frozen and is not retro-edited).

**Structure**: §1 anchors -> §2 fact package -> §3 completed slices -> §4 execution queue -> §5 observation items -> §6 pending rulings -> §7 / §8 / §9 pointer sections -> §10 opening -> §11 one-line status; the rule bodies behind §7-§9 have moved into directive §13, and this snapshot does not copy them.

**Revisions against v59**: (1) §1 anchors move on (HEAD `0f80c1e`, directive v1.69, the rolling window drops v56 and takes in v59); (2) §2 becomes the P76 fact package, its appendix A "obligation hand-off" becomes a **historical trace** (P75 O-7 residual settled by P76, no new hand-off) and appendix B "AT-23 stage-2 status" **stays unchanged**; (3) §3 adds one P76 row (the twelve v59 rows plus one, thirteen in total); (4) the §4 S1 Exit eight-gate states are updated to the P76 Q2=(a) re-evaluation (S-G1 not passed / S-G2 passed; the S-G5 and independent-security-review relation is clarified by ADR-047; S-G6 states gate level and task level side by side per Q8=(c)); the S-G3 / S-G4 / S-G6 closeout paths are made explicit; (5) the §5 OB-139 caliber follows the user ruling of 2026-11-16 that "the S-G2 correction [ROUTINE] D=1 does not count as the sixth occurrence", and the three slice-level suggestions stay; (6) §6 stays empty (Q2-Q10 are already ruled).

---

## §1 Current anchors (check these eight numbers first)

**Rolling window = the last three versions (v57 / v58 / v59)**; v56 has been dropped (same rule that dropped v52).

| Item | Value | Check command |
|---|---|---|
| HEAD | `0f80c1e` | `git rev-parse --short HEAD` |
| origin/main | `0f80c1e` | `git rev-parse --short origin/main` |
| Directive version | **v1.69** | the header of `docs/DELIVERY_DIRECTIVE.md` |
| SELFTEST | **737/737** | `node tools/gates/gate.js --selftest` |
| Three scopes | **all TOTAL_FAIL=0** | `gate.js --all --scope=index` / `tree` / `ci` |
| v57 whole-file line count | **160** (counted by LF, CN and EN equal) | `docs/t027/CHECKLIST_v57{,_EN}.md` counted by LF |
| v58 whole-file line count | **148** (counted by LF, CN and EN equal) | `docs/t027/CHECKLIST_v58{,_EN}.md` counted by LF |
| v59 whole-file line count | **180** (counted by LF, CN and EN equal) | `docs/t027/CHECKLIST_v59{,_EN}.md` counted by LF |

> The HEAD recorded here, `0f80c1e`, is the anchor at generation time, after the S-G2 correction commit; once v60 itself lands the HEAD moves on, and `git rev-parse --short HEAD` then governs.

> The rolling window holds the last three versions (v57 / v58 / v59) and v56 has been dropped (same rule that dropped v52); the line count of v60 itself is carried forward by v61. If two or more items disagree, the repository wins: report the final readings first and only then continue.

`tmp` standing assertion: the top level holds only `.gitkeep` and `gates-cache`; `gates-cache` holds only `g4b`; `tmp/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/` was cleaned at the P76 ⑧c stage and `tmp/P76-SG2-SYNC/` at the S-G2 correction ⑧c stage.

---

## §2 The P76 fact package

- **Slice / size**: `S1-EXIT-CLOSEOUT-ASSESSMENT` (P76; the S1 Exit closeout assessment; assessment type `[SLICE]`) / **M**; controller `DeepSeek V4.1 Flash (deepseek)`; authorized by the user in the same turn on 2026-11-16.
- **Goal**: the S1 Exit closeout roadmap, with Q1-Q11 analysed item by item (the eight-gate and ten-gate mapping, the gate-level re-evaluation caliber, the S-G3 to S-G6 unlock paths, and the product-owner approval timing).
- **Conclusion**: **closed** (landing commit `77f23d2`, backfill `aca19af`, S-G2 correction `0f80c1e`); the core ruling is that S1 Exit **cannot pass as a whole inside S1** (the S-G1 criterion is unmet, S-G8 stays blocked until S2, and S-G7 is unreachable).

**The eight rulings** (all taken by the user in the same turn on 2026-11-16):

- **Q2=(a) criterion-literal**: the S-G1 criterion ("a real CLI launch of an admitted AI candidate with billing") is unmet ⇒ **not passed**; the S-G2 criterion ("the `supervised-black-box` mode is implemented and the AT-24 verdict holds") is met (T028 `COMPLETED` plus AT-24 conditional closure) ⇒ **passed**. Basis: it agrees with G-A being "unreachable" in s1-reach, and the core-chain block is not weakened to PARTIAL.
- **Q3=A additive (the AT-23 fragment)**: inside row `:52`, "T027 is incomplete and AT-23 is not passed" becomes "T027 is incomplete and **AT-23 as a whole is not passed**" (CN / EN `as a whole`); the status cell `BLOCKED` stays, and the nine-row gate table keeps its structure.
- **Q7=A waived and deferred to S2**: the independent security review has **never been executed** (basis: `SECURITY.md:5`, no review report anywhere in the repository, and `DEV_PLAN:27`); **ADR-047** records the waiver reason, the obligation to complete the review before S2, the due date (the earlier of the S2 start date and 2027-03-31) and the uncovered scope; SECURITY carries a note and PRODUCT_REQUIREMENTS §7 an additive pointer.
- **Q9 included**: the tail clauses of `S1_SUPPORT_MATRIX{,_EN}.md:13` / `:14` are corrected (`:13` now reads "R-3 = (A), cannot start for now; S-G8 stays blocked until S2 enforcement lands (ADR-045)" and `:14` reads "AT-24 conditionally closed; the real CLI end-to-end path depends on S2 enforcement (R-3 = (A), ADR-045)"); `:5` is untouched.
- **Q10=(C) two steps**: step one, inside S1, approves ADR-010 and the v0.1.0 Release so that S-G3 / S-G4 close, with S-G5 booked as "release portion approved, whole-stage approval pending S2"; step two, after S2, approves S1 Exit as a whole. **Hard boundary: no approval may be read as S1 Exit passing.**
- **Closeout order** (the user accepted the P76 §12 sequence): the `:52` admission correction, then the S-G2 gate-level settlement, then S-G4 (ADR-010 approval), then S-G3 (Release plus ZIP binding plus rerun), then the S-G6 gate-level side stated in parallel, then the stepwise S-G5 approval, and finally S-G1 / S-G8 / S-G7 re-evaluated after S2.
- **Additional (the p57 gate-table rows stay unchanged)**: the S-G1 / S-G2 / S-G7 verdict words in p57 §4.2 are untouched and an additive pointer is placed after §4.2. **⑥ may therefore be skipped.**
- **X authorization (Q3 flexibility)**: Q3 is a **directional ruling** (the direction is "A additive"); the concrete landing is chosen by the controller between option A and option B, with **A preferred and B used only if A proves insufficient on read-back**; if the row turns out to be unambiguous (rare), stop and escalate under §1.7.2 rather than deciding unilaterally. On read-back the P76 controller confirmed that `:52` already states "stage 2 not applicable" twice (once in CN, once in EN), so **option A was chosen and the qualifier "整体" added** (`AT-23 未通过` becoming `AT-23 整体未通过`, EN `as a whole`) without touching the status cell.

**The S-G2 correction (`0f80c1e`)**: the P76 ⑩ receipt registered one residual, namely that the S-G2 sentence inside row `:52` still read "pending the v54 snapshot, must not be cited as that gate passing" (an earlier point-in-time caliber) and therefore conflicted with Q2=(a). The correction replaces the S-G2 gate-level verdict sentence with "**has been confirmed as passed by P76 under the criterion-literal approach** (see `s1-exit-closeout-assessment{,_EN}.md` and ADR-046); the previous `pending the v54 snapshot` was an earlier point-in-time statement and has been superseded by P76", keeps the leading "user ruling ③" half, and follows the disk text where it differs from the brief (CN has no space inside 「用户裁决③」 and EN actually says caliber rather than `criterion`); the **A1 escalation branch did not trigger**, because the AT-24 snapshot marking and the S-G2 gate-level verdict are independent.

**ADRs**: **ADR-046** (the S1 Exit closeout caliber and approval timing; a decision record of the same class as ADR-042 / 043 / 044 / 045 and therefore not bound by "must not be used as a precedent") and **ADR-047** (the independent security review waived and deferred to S2; a risk-acceptance decision record). Max ADR = **47** (ADR-045 belongs to P75).

**CI**: the P76 landing run `43785603244` (ubuntu 3m24s / windows 5m12s, both legs green), the backfill run `43787083404` (an observation record only) and the S-G2 correction run `43793273520` (ubuntu 3m16s / windows 4m28s, both legs green).

**Review chain (the P76 main slice)**: (1) `DESIGN: DONE` (deep-reasoner plus DeepSeek V4 Pro) -> (2a) ruling registration -> (2b) document sync -> (2c) report write-up -> (3) consistency -> (4) gates -> (5) r1 `PRE-REVIEW: FINDINGS` (1 High plus 5 Low) -> (5) r2 `PASS WITH FIXES` (none at medium or above, 2 Low) -> (6) skipped with (7) `SKIP CONFIRMED` -> (7) final review `PASS WITH FIXES` (none at medium or above; its single Low was turned into the **proposed slice name `S1-EXIT-SG2-SYNC`**, which **never started**, because the fix landed instead as the [ROUTINE] correction `0f80c1e`) -> (8a) / (9) / (8b) / (8c).

**The S-G2 correction review**: the [ROUTINE] tier (no (5), (6) or (7); the controller executed it directly); all eight step-0 anchors matched, all four gate groups were green, and both CI legs succeeded.

**Metering**: the P76 main slice recorded **R = 1** (the single High from (5) r1, fixed and re-run) and **D = 0** (no temp-file violation in the main slice). The S-G2 correction [ROUTINE] recorded **D = 1** (`git show --output=` dropped two files at the `tmp/` root; it self-corrected on the spot with zero residue) and the **user confirmed the self-healing on 2026-11-16**; **in the same ruling the user decided that this [ROUTINE] D=1 does not count as the sixth OB-139 occurrence** (a [ROUTINE] correction is not a formal [SLICE]), so OB-139 stays under observation and is counted again only on the next [SLICE] with the same pattern.

**Numbered deviations**: P76 has none (the (5) r1 High and the (5) r2 and (7) r1 Lows all closed through the normal loop).

**Boundaries**: real enforcement stays in S2; S-G8 stays blocked; S-G1 is not passed (unreachable inside S1); the probe budget is zero; nothing under `s1-reach` / `CONTRACTS` / schema / fixtures / CI / product code / v50-v59 was touched.

**Mechanical readings**: **max ADR = 47**; max OB = 139 (none added); max DR = 14; pilot deferral = round 25.

**Evidence pack**: `docs/validation/evidence/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/` (`MANIFEST.md` plus 15 `SHA256SUMS.txt` entries over 15 payloads: `00-briefing.md` / `input-list.md` / `freedom-list.md` / `01-arch/design-r1.md` / `02-batch/batch-01-s1-exit-p57.md` / `batch-02-s1-reach.md` / `batch-03-lagging-docs.md` / `02-ruling/ruling-2026-11-16.md` / `04-reviews/prereview-r1.md` / `prereview-r2.md` / `final-r1.md` / `06-gates/gate-index.txt` / `gate-tree.txt` / `gate-ci.txt` / `gate-selftest.txt`).

**Known boundaries**: **the AT-23 stage-1 condition closure is a permanent fact (ruling A3), and R-08.5 states that the conditional closure is not a real-candidate compatibility proof**; the risk surfaces listed in §17 of the P76 report plus the uncovered scope in ADR-047; S1 Exit cannot pass as a whole inside S1 (the premise of the Q10=(C) two-step approval); the S-G5 and independent-security-review relation is clarified by ADR-047, with a mandatory review before S2.

**Ledger caliber** (§11.2): the P76 completion row records only the report paths plus R/D (no hash or run id), and the CI fact chain lives only in §16 of the report.

### Appendix A: the "obligation hand-off" section (historical trace; v60 hands nothing new)

- **Historical hand-off** (first traced back to P75 O-7): of the three stale statements in P75 §11.2, items one and two (originally items one and three in P57 §4.3) were already closed by the P64 and s1-reach corrections; the only residual was the in-cell caliber mismatch in the core task chain row at `s1-exit.md:52` (originally item two, later item three in P75 §11.2), whose landing point was the S1 Exit closeout slice, i.e. P76, where **correction 1 (the AT-23 fragment) and the S-G2 correction (`0f80c1e`) have now discharged it**.
- **Status in this version**: **no new hand-off** (every P76 ruling closed inside P76, and nothing was handed to a future slice).
- **Convention note**: this "obligation hand-off" treatment has been used since v59 (first applied to P75 O-7): whenever a slice hands an obligation to a future slice, that slice's fact package must carry a hand-off section, the landing slice entry in §4 must carry an "admission condition" section, and §5 may optionally track it as a "pending execution list" (without registering an OB). **This version only traces the convention; nothing new is handed off.**

### Appendix B: the "AT-23 stage-2 status" section (unchanged)

- **Status**: ruled "not applicable" by A1 (the user ruling of 2026-11-09; this is **not a permanent prohibition**, unless the user later rules explicitly to restart it).
- **What can be executed before S2**: nothing (R-3 = (A), and the real CLI end-to-end run needs real enforcement, which is an S2 premise).
- **Restart path**: if the user later rules to restart stage 2, real enforcement must land first (S2) or an alternative path must exist (P-prime, a candidate swap, or a new intermediate path); an alternative path means "reach an equivalent goal before S2 by another route", not "implement stage 2 before S2".
- **Natural restart point**: S2 starts, enforcement lands, S-G8 unlocks, R-3 is re-judged (A possibly becoming B), and if B is feasible that is the natural restart point for stage 2.

---

## §3 Completed slices (recent baseline)

**Section 3 of v59 carried twelve rows, P64 through P75; v60 reuses those twelve (readings copied from v59 as-is, not retro-edited) and appends one P76 row. Thirteen rows in total.**

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

> **The D caliber in §3** (same as v59): "D = n" in this table is shorthand for the temp-file discipline violation count (§5.1); P76 as a whole records D = 1 (D = 0 in the main slice plus D = 1 in the S-G2 correction [ROUTINE]). **The D count and the OB-139 ruling are two separate things**: D = 1 is an objective count (already self-corrected and confirmed by the user), while the sixth-occurrence ruling was made by the user on 2026-11-16 and excludes it (a [ROUTINE] correction is not a formal [SLICE]).

---

## §4 Execution queue

### Layer 0 preliminary triage (✅ all dequeued, nothing new in this version)

- P5 and P6 were closed by P57; layer 0 is fully dequeued with nothing left.

### Layer 1 main-line hard gates

- Items 8-12: ✅ all dequeued (see §4 of v57 / v58; **the original dequeue records have been pushed into earlier versions by the rolling window, so trace them back layer by layer through the citation chain in §4 of v57**).
- **`T027-B` (the defensive-refactor main-line milestone) = complete (partly not applicable)**: C1-C7 / N1-N4 / R-1 / R-2 / R-3 are all dispositioned (R-3 = (A), ruled by the user on 2026-11-15; the A1 declaration caliber is "complete, partly not applicable"). **The main line is fully dequeued.**
- **`S1 Exit closeout` (in progress)**: **P76 produced the closeout roadmap and re-evaluated the gate-level verdicts**. The eight gates stand as follows after the P76 Q2=(a) re-evaluation:

| Gate | Criterion | State after the P76 re-evaluation | Basis |
|---|---|---|---|
| S-G1 | a real CLI launch of an admitted AI candidate with billing | **not passed** (Q2=(a) criterion-literal, agreeing with G-A being unreachable in s1-reach) | P76 §4; s1-reach G-A; ADR-046 |
| S-G2 | the `supervised-black-box` mode is implemented and the AT-24 verdict holds | **passed** (Q2=(a) criterion-literal; T028 `COMPLETED` plus AT-24 conditional closure) | P76 §4; the S-G2 sentence at `s1-exit.md:52` (S-G2 correction `0f80c1e`); ADR-046 |
| S-G3 | the final ZIP is bound to the release record and rerun | not passed | `INSTALLATION.md` states there is no formal GitHub Release yet; P76 §7 |
| S-G4 | ADR-010 is approved by the product owner | not passed | ADR-010 is `PROPOSED / NOT APPROVED`; P76 §8 |
| S-G5 | the product owner approves | unknown (a human gate) | `s1-exit.md:60`; P76 §12 |
| S-G6 | the dependency fields are satisfied and the four remaining gates close | unreachable (**P76 Q8=(c), stated in parallel**: closing S-G3 / S-G4 / S-G5 closes the **gate level**, while T018 at the **task level stays constrained by T027 / T028**) | P76 §10; `DEV_PLAN.md:94` / `:96` |
| S-G7 | every AT passes, T018 is complete and every §14 exit holds | unreachable (a composite gate) | `s1-exit.md:52` / `:56`-`:58` / `:60` |
| S-G8 | the supervised path runs end to end against a real CLI and a real candidate | blocked until S2 enforcement lands | P75 A2 / ADR-045; P76 §6 |

**The independent security review (not an S-G5 criterion, but a precondition of the second approval step)**: the S0-to-S1 gate in `PRODUCT_REQUIREMENTS` §7, which has **never been executed**; **ADR-047** waives it and defers it to S2, with the review to complete before S2 and a due date at the earlier of the S2 start date and 2027-03-31. **The Q10=(C) two-step S-G5 approval distinguishes them**: **step one**, inside S1 and covering the release portion, does not require the review to be complete (ADR-047 defers it to S2), whereas **step two**, after S2 and covering the whole S1 Exit, requires the review to have completed before S2 starts. **Hard boundary: no approval may be read as S1 Exit passing.**

- **S1 Exit cannot pass as a whole inside S1** (the premise of Q2=(a) and Q10=(C)). The remaining actionable items are **S-G3 / S-G4 / S-G5 / S-G6** (S-G7 is composite, S-G8 needs S2, and S-G1 is not passed).
- **No admission residual blocks the closeout follow-up slices**: the only P75 §11.2 residual (the core task chain row at `s1-exit.md:52`) was fully corrected by P76 correction 1 (the AT-23 fragment) plus the S-G2 correction (`0f80c1e`).
- **`ADR010-APPROVAL` (S-G4 closeout · `next item` · starting on the user's 2026-11-16 call · P77)**

| Field | Content |
|---|---|
| Identifier | `ADR010-APPROVAL` (ledger section = **P77**, registered when the kick-off message is drafted) |
| Contents | (1) flip the ADR-010 status in `ADR_REGISTER{,_EN}.md` to `APPROVED` (currently `PROPOSED / NOT APPROVED`); (2) update the "this proposal is not yet approved" sentence in the tail of `S1_SUPPORT_MATRIX{,_EN}.md` to approved plus the 90-day EOL window in force; (3) update the EOL item at `INSTALLATION_PLAN{,_EN}.md:32`; (4) update any matching sentence in `S1_RELEASE_NOTES{,_EN}.md`; the controller performs the document updates, and the product owner's approval is stated by the user in the same turn as this slice's kick-off message |
| Size tier | **S** (`[SLICE]`, documentation type) |
| (7) budget | tier S: one final review plus one re-review (§7.5) |
| Dependency | none (the product owner approves personally); may run in parallel with `RELEASE-ZIP-BINDING` |
| Boundary | no other ADR is touched; the S1_SUPPORT_MATRIX platform list and support levels stay; the INSTALLATION install model stays; no directive change; no product code |

- **`RELEASE-ZIP-BINDING` (S-G3 closeout · P78)**

| Field | Content |
|---|---|
| Identifier | `RELEASE-ZIP-BINDING` (ledger section = **P78**, registered when the kick-off message is drafted) |
| Contents | (1) create the v0.1.0 Release on GitHub; (2) build the ZIP artifact (the Windows amd64 portable ZIP); (3) bind the release record (commit SHA / CI / SHA256SUMS / SPDX SBOM / license manifest); (4) rerun the install lifecycle (first run / upgrade / rollback / uninstall keeping data) |
| Size tier | **M-L** (`[SLICE]`, implementation type; produces a publishable artifact) |
| (7) budget | tier M-L: one final review plus three re-reviews (§7.5) |
| Dependency | **`ADR010-APPROVAL` closed** (S-G4 passed and the version comes from ADR-010) |
| Boundary | the ZIP structure stays (ADR-009 froze it: manual extraction, a dedicated directory, explicit-path invocation, no PATH / registry / system-directory changes); no product code; no formal Linux release (S1 is Windows formal support plus Linux core CI only); **no code signing** (ADR-006 simplified release trust: fixed commit SHA, GitHub CI, SHA256SUMS, SBOM and a license manifest are the required evidence); no gitee push |

- **`S2-PRE-INDEPENDENT-SECURITY-REVIEW` (the ADR-047 compensation obligation · P79)**

| Field | Content |
|---|---|
| Identifier | `S2-PRE-INDEPENDENT-SECURITY-REVIEW` (ledger section = **P79**, registered when the kick-off message is drafted) |
| Contents | complete the independent security review (the S0-to-S1 gate in `PRODUCT_REQUIREMENTS` §7); archive the review conclusion under `docs/validation/` and update the note at `SECURITY.md:5` to "completed, with the review date and reviewer"; a later slice may flip the ADR-047 risk-acceptance record to "compensated" |
| Size tier | **M-L** (`[SLICE]`, assessment type; the review is work in its own right) |
| (7) budget | tier M-L: one final review plus three re-reviews (§7.5) |
| Dependency | before S2 starts (not an immediate S1 need); the ADR-047 due date is the earlier of the S2 start date and 2027-03-31 |
| Boundary | the review scope is the S0-to-S1 security gate; the S1 Exit eight-gate criteria stay; the `SECURITY.md` threat-model baseline stays (only the review conclusion is archived); no directive change |

- **`S-G8`** stays blocked until S2 enforcement lands (P73 A2 / P75 A2 / ADR-044 / ADR-045); the gate table in `s1-exit{,_EN}.md` is untouched (P76 only edited the fragment inside row `:52`).

### Layer 2 routine governance

- P72 / P73 / P74 / P75 / P76 ✅ are all dequeued, so **layer 2 is empty** (see the row summaries in §3).
- The complete §1.2 field definitions for P72 / P73 / P74 / P75 / P76 live in the matching sections of `docs/t027/REMAINING_SLICES{,_EN}.md`.

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

**Queue discipline**: as in v59 (layer boundaries and numbering discipline; this [ROUTINE] slice runs in a single batch and does not bump the version).

---

## §5 Observation items

- Observing (**14**): `OB-1 / 7 / 10 / 33 / 87 / 99 / 100 / 104 / 105 / 107 / 111 / 113 / 135 / 139` (carried from v59).
  - **OB-139** (temp files placed at the `tmp/` root instead of a slice directory; four times in P70 / P71 plus once in P75 makes the fifth occurrence; the user ruled on 2026-11-16 to **keep it under observation** and that **a sixth occurrence becomes pending-plus-a-dedicated-slice**): **in the same ruling the user decided that the S-G2 correction [ROUTINE] D=1 does not count as the sixth occurrence** (a [ROUTINE] correction is not a formal [SLICE]); OB-139 stays under observation and is counted again on the next [SLICE] with the same pattern. The P76 main slice recorded D = 0 and the S-G2 correction [ROUTINE] D = 1 (self-corrected on the spot, zero residue, self-healing confirmed by the user).
  - **The reinforced D=2 recurrence-prevention sentence (standard kick-off content since v59, carried into v60)**: "the executing roles (`prfrail-implementer` / `prfrail-tester` / `prfrail-documenter`) **and the controller** must always work under a `tmp/` slice directory (`tmp/<slice-id>/`); **the same holds for CI observation captures, which must never land at the `tmp/` root**; three D occurrences inside one slice force the slice to stop."
- Pending discussion (9): `OB-19 / 21 / 52 / 53 / 57 / 58 / 59 / 61 / 75`.
- To do (1): `OB-69`. Non-terminal total **24** (14 plus 9 plus 1).
- **Slice-level suggestions (three, not registered as OBs)**:
  1. **`R-ACCUMULATION-MECHANISM` (carried from v58)**: P74 was the first case to trigger the §9.2 rule (R>=2 means high risk) together with the §7.5 high-risk tier (one plus three); later implementation slices should pre-place the dynamic budget sentence in their kick-off message; P75 and P76 recorded R = 1, below the threshold, so it did not trigger. **No OB is registered now**; a recurrence on P77 or later turns it into a formal OB.
  2. **`--progress` covers state transitions only, not log output (carried from v58)**: if long-running observability becomes a real need (the S1 Exit closeout slices are the most likely trigger), evaluate adding `--progress verbose` or `--tail-logs`.
  3. **The "obligation hand-off" convention (in use since v59, first applied to P75 O-7)**: whenever a slice hands an obligation to a future slice, that slice's fact package must carry a hand-off section, the landing slice entry in §4 must carry an "admission condition" section, and §5 may track it as an optional "pending execution list" (without an OB). **No OB is registered now**; if a hand-off is ever left dangling it becomes a formal OB.
- **Removed in this version**: none.
- **Flipped to `settled` in this version**: none.
- Mechanical readings: max OB in §12.4 = **139**; **max ADR = 47**; max DR = 14; max deferral = round 25.

---

## §6 Pending user rulings

- **This table is currently empty** (Q2 / Q3 / Q7 / Q9 / Q10, the closeout order, the additional item and X were all ruled in P76; the P76 main slice raised no new ruling item; the user confirmed the S-G2 correction D=1; and the user ruled that it does not count as the sixth OB-139 occurrence).
- **Anticipated next ruling**: the ordering and priority of the S1 Exit closeout follow-up slices (the ADR-010 approval slice, the Release plus ZIP binding and rerun slice, and the pre-S2 independent-security-review slice).

---

## §7 / §8 / §9 Pointer sections

Same structure as v59: §7 (kick-off and step-0 general rules) points at directive §13.1; §8 (cross-slice collaboration discipline) at §13.2; §9 (checklist meta-rules and gate pitfalls) at §13.3 / §13.4. **The rule bodies have moved into directive §13 and are not copied here.**

---

## §10 One-line opening for a new session

> Below is my project state snapshot (v60). Check the eight anchors in section 1 first; if they match, continue with the queue in section 4, and if not, the repository state governs.
>
> [paste this checklist v60 in full]
>
> [Supporting documents] `REMAINING_SLICES{,_EN}.md` (P76 closed), `DELIVERY_DIRECTIVE.md` (v1.69), `s1-exit-closeout-assessment{,_EN}.md` (the P76 report), `r3-real-cli-assessment{,_EN}.md` (the P75 report), `sandbox-interface{,_EN}.md` (the P74 report), `t027b-admission{,_EN}.md` (the P73 report), `req020-s2-restriction{,_EN}.md` (the P72 report), `ADR_REGISTER{,_EN}.md` (including ADR-045 / ADR-046 / ADR-047), `DOCUMENTATION_PLAN{,_EN}.md` (including the cross-session hand-off entry pointer) and `CHECKLIST_v59{,_EN}.md` (frozen and superseded by v60).
>
> Current slice: P76 is closed (including the S-G2 correction). Next: the **S1 Exit closeout follow-up slices** (the ADR-010 approval slice, the Release plus ZIP binding and rerun slice, and the pre-S2 independent-security-review slice), pending the user's scheduling.

## §11 One-line current state

P76 is closed (`77f23d2` + `aca19af` + `0f80c1e` are on origin and runs `43785603244` + `43787083404` + `43793273520` were green on both legs; R = 1 and D = 1 for P76 as a whole: D = 0 in the main slice and D = 1 in the S-G2 correction [ROUTINE], with the self-healing confirmed by the user; on OB-139 the user ruled on 2026-11-16 that this [ROUTINE] D=1 does not count as the sixth occurrence; **Q2=(a) leaves S-G1 not passed and S-G2 passed**; **Q3=A leaves AT-23 as a whole not passed**; **Q7=A waives the independent security review under ADR-047 and defers it to S2**; **Q10=(C) approves in two steps**; ADR-046 / ADR-047 are registered; the S-G2 correction is closed); the directive is v1.69; layer 1 T027-B is fully dequeued and layer 2 is empty; next comes the **S1 Exit closeout follow-up slices** (the ADR-010 approval slice, the Release plus ZIP binding and rerun slice, and the pre-S2 independent-security-review slice), pending the user's scheduling.
