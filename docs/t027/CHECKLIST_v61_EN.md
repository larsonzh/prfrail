# Task Execution Checklist (cross-session snapshot v61, after P77 closes, S1 Exit closeout in progress)

Generated: 2026-11-17. HEAD = `fc00c4b`; directive = v1.69; two CI runs: the P77 landing commit `43986942087` (ubuntu 3m22s / windows 4m32s, both legs green) and the P77 backfill `43997058549` (**an observation record only, per §11.2**).

**Supersession**: this version supersedes v60 (`docs/t027/CHECKLIST_v60{,_EN}.md` is frozen and is not retro-edited).

**Structure**: §1 anchors -> §2 fact package -> §3 completed slices -> §4 execution queue -> §5 observation items -> §6 pending rulings -> §7 / §8 / §9 pointer sections -> §10 opening -> §11 one-line status; the rule bodies behind §7-§9 have moved into directive §13, and this snapshot does not copy them.

**Revisions against v60**:
- (1) §1 anchors move on: HEAD `fc00c4b`; the rolling window drops v58 and takes in v60 (v59 / v60 / v61).
- (2) §2 becomes the P77 fact package; appendix A traces the hand-off history (nothing new is handed off) and appendix B keeps the AT-23 stage-2 status unchanged.
- (3) §3 adds one P77 row (the thirteen v60 rows plus one, fourteen in total).
- (4) the §4 S1 Exit eight-gate states are updated: S-G4 moves from not passed to passed (ADR-010 approved, P77); the `ADR010-APPROVAL` queue item is dequeued; `RELEASE-ZIP-BINDING` (P78) is the next item with its dependency now met.
- (5) §5 keeps OB-139 under observation (the P77 main slice recorded D = 0, so the sixth occurrence did not trigger); the three slice-level suggestions stay.
- (6) §6 stays empty.

---

## §1 Current anchors (check these eight numbers first)

**Rolling window = the last three versions (v59 / v60 / v61)**; v58 has been dropped (same rule that dropped v52).

| Item | Value | Check command |
|---|---|---|
| HEAD | `fc00c4b` | `git rev-parse --short HEAD` |
| origin/main | `fc00c4b` | `git rev-parse --short origin/main` |
| Directive version | **v1.69** | the header of `docs/DELIVERY_DIRECTIVE.md` |
| SELFTEST | **737/737** | `node tools/gates/gate.js --selftest` |
| Three scopes | **all TOTAL_FAIL=0** | `gate.js --all --scope=index` / `tree` / `ci` |
| v59 whole-file line count | **180** (counted by LF, CN and EN equal) | `docs/t027/CHECKLIST_v59{,_EN}.md` counted by LF |
| v60 whole-file line count | **242** (counted by LF, CN and EN equal) | `docs/t027/CHECKLIST_v60{,_EN}.md` counted by LF |
| v61 whole-file line count | carried forward by v62 (this snapshot's own count) | `docs/t027/CHECKLIST_v61{,_EN}.md` counted by LF |

> The HEAD recorded here, `fc00c4b`, is the anchor at generation time, after the P77 backfill commit; once v61 itself lands the HEAD moves on, and `git rev-parse --short HEAD` then governs.

> The rolling window holds the last three versions (v59 / v60 / v61) and v58 has been dropped (same rule that dropped v52); the line count of v61 itself is carried forward by v62. If two or more items disagree, the repository wins: report the final readings first and only then continue.

`tmp` standing assertion: the top level holds only `.gitkeep` and `gates-cache`; `gates-cache` holds only `g4b`; `tmp/P77-ADR010-APPROVAL/` was cleaned at the P77 ⑧c stage and `tmp/P77-backfill/` was cleaned at the backfill step.

---

## §2 The P77 fact package

- **Slice / size**: `ADR010-APPROVAL` (P77; the EOL notification documentation sync after the product owner approves ADR-010; documentation type `[SLICE]`) / **S**; controller `DeepSeek V4.1 Flash (deepseek)`; the user approved ADR-010 in the same turn on 2026-11-16.
- **Goal**: flip ADR-010 to `APPROVED` (with the four approval elements) and sync the now-effective EOL notification policy to every related document.
- **Conclusion**: **closed** (landing commit `33272be`, backfill `fc00c4b`); **S-G4 moves from not passed to passed**; the p57 and `s1-exit` gate tables are untouched (per the P76 ruling).
- **Four approval elements**: ruling party = the user; date = 2026-11-16; reason = freezing the EOL notification channel before v0.1.0 enters the release flow; basis = the user's same-turn P77 kickoff ruling of 2026-11-16.
- **Landing spots (six file pairs, all in-cell replacement or additive)**: the ADR-010 status cell at `ADR_REGISTER:20` (`` `APPROVED` (four elements); the historical completion sentence kept byte-for-byte plus a closing note ``) plus the support-policy paragraph tail at `:62` plus the decision-table platform/support cell at `:102`; `S1_SUPPORT_MATRIX:21/:23`; `INSTALLATION_PLAN:32`; the additive pointer at `PRODUCT_REQUIREMENTS:107`; residual-risk item 1 at `s1-exit:66` (the gate table is untouched).
- **Kept unchanged and disclosed**: the dated historical conclusion paragraph at `s1-exit:7` (§11.1); the old state quote in the gate-table row `p57:91` (the table is untouched); three dated records in `DEV_PLAN`; the frozen statements in `CHECKLIST_v59/v60` and `s1-reach`.
- **Review chain**: (1) `DESIGN: DONE` (deep-reasoner plus DeepSeek V4 Pro) -> batch assistance `GLM-5.3-Flash` twice -> (2) document sync -> (3) consistency -> (4) gates -> (5) r1 `PRE-REVIEW: PASS WITH FIXES` (one Medium plus three Low) -> (5) r2 `PASS WITH FIXES` (none at medium or above, three Low; the budget is exhausted at 2/2 and the three Low are confirmed by (7)) -> (6) skipped with (7) `SKIP CONFIRMED` plus the tier-(i) coverage/omission statement -> (7) final review `FINAL REVIEW: PASS WITH FIXES` (none at medium or above, one Low) -> (8a) / (9) / (8b) / (8c).
- **Defects fixed (all closed)**: the three Low findings of (5) r2 (the §3.4 difference disclosure, the §3.3 union caliber and the §7.2 gate-reading landing spot) and the single Low of (7) r1 (the §8 frozen-snapshot residue note) were all fixed to the minimal repair and re-checked.
- **CI**: the P77 landing run `43986942087` (ubuntu 3m22s / windows 4m32s, both legs green) and the backfill run `43997058549` (an observation record only).
- **Metering**: **R = 1** (the single Medium from (5) r1, fixed and re-run); **D = 0**.
- **Pre-existing clue (not introduced here)**: `docs/validation/t025-operator-interaction.md` carries a double BOM at HEAD; it is a pre-existing encoding defect for a later slice (this slice registers no OB, per its boundary).
- **Mechanical readings**: **max ADR = 47**; max OB = 139 (none added); max DR = 14; pilot deferral = round 25.
- **Evidence pack**: `docs/validation/evidence/P77-ADR010-APPROVAL/` (`MANIFEST.md` plus 13 `SHA256SUMS.txt` entries over 13 payloads: `00-briefing.md` / `input-list.md` / `freedom-list.md` / `01-arch/design-r1.md` / `02-batch/batch-01-adr010-register-format.md` / `batch-02-eol-survey.md` / `04-reviews/prereview-r1.md` / `prereview-r2.md` / `final-r1.md` / `06-gates/gate-index.txt` / `gate-tree.txt` / `gate-ci.txt` / `gate-selftest.txt`).
- **Ledger caliber** (§11.2): the P77 completion row records only the report paths plus R/D (no hash or run id), and the CI fact chain lives only in §10 of the report.

### Appendix A: the "obligation hand-off" section (historical trace; v61 hands nothing new)

- **Historical hand-off** (first traced back to P75 O-7): the in-cell caliber mismatch in the core task chain row at `s1-exit.md:52` was discharged by P76 correction 1 (the AT-23 fragment) plus the S-G2 correction (`0f80c1e`).
- **Status in this version**: **no new hand-off** (every P77 item closed inside P77, and nothing was handed to a future slice).
- **Convention note**: this treatment has been used since v59; whenever a slice hands an obligation to a future slice, that slice's fact package must carry a hand-off section. **This version only traces the convention; nothing new is handed off.**

### Appendix B: the "AT-23 stage-2 status" section (unchanged)

- **Status**: ruled "not applicable" by A1 (the user ruling of 2026-11-09; this is **not a permanent prohibition**, unless the user later rules explicitly to restart it).
- **What can be executed before S2**: nothing (R-3 = (A), and the real CLI end-to-end run needs real enforcement, which is an S2 premise).
- **Restart path**: if the user later rules to restart stage 2, real enforcement must land first (S2) or an alternative path must exist (P-prime, a candidate swap, or a new intermediate path); an alternative path means "reach an equivalent goal before S2 by another route", not "implement stage 2 before S2".
- **Natural restart point**: S2 starts, enforcement lands, S-G8 unlocks, R-3 is re-judged (A possibly becoming B), and if B is feasible that is the natural restart point for stage 2.

---

## §3 Completed slices (recent baseline)

**Section 3 of v60 carried thirteen rows, P64 through P76; v61 reuses those thirteen (readings copied from v60 as-is, not retro-edited) and appends one P77 row. Fourteen rows in total.**

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

> **The D caliber in §3** (same as v60): "D = n" in this table is shorthand for the temp-file discipline violation count (§5.1); P77 as a whole records D = 0. **The D count and the OB-139 ruling are two separate things**: the P77 main slice recorded D = 0, so the sixth OB-139 occurrence did not trigger.

---

## §4 Execution queue

### Layer 0 preliminary triage (✅ all dequeued, nothing new in this version)

- P5 and P6 were closed by P57; layer 0 is fully dequeued with nothing left.

### Layer 1 main-line hard gates

- Items 8-12: ✅ all dequeued (see §4 of v57 / v58; **the original dequeue records have been pushed into earlier versions by the rolling window, so trace them back layer by layer through the citation chain in §4 of v57**).
- **`T027-B` (the defensive-refactor main-line milestone) = complete (partly not applicable)**: C1-C7 / N1-N4 / R-1 / R-2 / R-3 are all dispositioned (R-3 = (A), ruled by the user on 2026-11-15; the A1 declaration caliber is "complete, partly not applicable"). **The main line is fully dequeued.**
- **`S1 Exit closeout` (in progress)**: **P76 produced the closeout roadmap and re-evaluated the gate-level verdicts, and P77 closed S-G4**. The eight gates stand as follows after P77 closes:

| Gate | Criterion | State after P77 closes | Basis |
|---|---|---|---|
| S-G1 | a real CLI launch of an admitted AI candidate with billing | **not passed** (Q2=(a) criterion-literal, agreeing with G-A being unreachable in s1-reach) | P76 §4; s1-reach G-A; ADR-046 |
| S-G2 | the `supervised-black-box` mode is implemented and the AT-24 verdict holds | **passed** (Q2=(a) criterion-literal; T028 `COMPLETED` plus AT-24 conditional closure) | P76 §4; the S-G2 sentence at `s1-exit.md:52` (S-G2 correction `0f80c1e`); ADR-046 |
| S-G3 | the final ZIP is bound to the release record and rerun | not passed | `INSTALLATION.md` states there is no formal GitHub Release yet; P76 §7 |
| S-G4 | ADR-010 is approved by the product owner | **passed** (ADR-010 is `APPROVED`; user approval on 2026-11-16; P77 closed) | P77 report §5; `ADR_REGISTER.md:20` |
| S-G5 | the product owner approves | unknown (a human gate) | `s1-exit.md:60`; P76 §12 |
| S-G6 | the dependency fields are satisfied and the four remaining gates close | unreachable (**P76 Q8=(c), stated in parallel**: closing S-G3 / S-G5 closes the **gate level**, while T018 at the **task level stays constrained by T027 / T028**) | P76 §10; `DEV_PLAN.md:94` / `:96` |
| S-G7 | every AT passes, T018 is complete and every §14 exit holds | unreachable (a composite gate) | `s1-exit.md:52` / `:56`-`:58` / `:60` |
| S-G8 | the supervised path runs end to end against a real CLI and a real candidate | blocked until S2 enforcement lands | P75 A2 / ADR-045; P76 §6 |

**The independent security review (not an S-G5 criterion, but a precondition of the second approval step)**: the S0-to-S1 gate in `PRODUCT_REQUIREMENTS` §7, which has **never been executed**; **ADR-047** waives it and defers it to S2, with the review to complete before S2 and a due date at the earlier of the S2 start date and 2027-03-31. **The Q10=(C) two-step S-G5 approval distinguishes them**: **step one**, inside S1 and covering the release portion, does not require the review to be complete (ADR-047 defers it to S2), whereas **step two**, after S2 and covering the whole S1 Exit, requires the review to have completed before S2 starts. **Hard boundary: no approval may be read as S1 Exit passing.**

- **S1 Exit cannot pass as a whole inside S1** (the premise of Q2=(a) and Q10=(C)). The remaining actionable items are **S-G3 / S-G5 / S-G6** (S-G4 has passed; S-G7 is composite, S-G8 needs S2, and S-G1 is not passed).
- **No admission residual blocks the closeout follow-up slices**: the only P75 §11.2 residual (the core task chain row at `s1-exit.md:52`) was fully corrected by P76 correction 1 (the AT-23 fragment) plus the S-G2 correction (`0f80c1e`).
- **`ADR010-APPROVAL` (S-G4 closeout)**: **P77 is closed** (landing commit `33272be` plus backfill `fc00c4b`; R = 1, D = 0; S-G4 passed; the p57 / `s1-exit` gate tables untouched). Dequeued; the fact package is in §2 and the report is `docs/validation/adr010-approval{,_EN}.md`.
- **`RELEASE-ZIP-BINDING` (S-G3 closeout · `next item` · P78)**

| Field | Content |
|---|---|
| Identifier | `RELEASE-ZIP-BINDING` (ledger section = **P78**, registered when the kick-off message is drafted) |
| Contents | (1) create the v0.1.0 Release on GitHub; (2) build the ZIP artifact (the Windows amd64 portable ZIP); (3) bind the release record (commit SHA / CI / SHA256SUMS / SPDX SBOM / license manifest); (4) rerun the install lifecycle (first run / upgrade / rollback / uninstall keeping data) |
| Size tier | **M-L** (`[SLICE]`, implementation type; produces a publishable artifact) |
| (7) budget | tier M-L: one final review plus three re-reviews (§7.5) |
| Dependency | **already satisfied** (S-G4 has passed: the ADR-010 approval was landed by P77) |
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

- P72 / P73 / P74 / P75 / P76 / P77 ✅ are all dequeued, so **layer 2 is empty** (see the row summaries in §3).
- The complete §1.2 field definitions for P72 / P73 / P74 / P75 / P76 / P77 live in the matching sections of `docs/t027/REMAINING_SLICES{,_EN}.md`.

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

- Observing (**14**): `OB-1 / 7 / 10 / 33 / 87 / 99 / 100 / 104 / 105 / 107 / 111 / 113 / 135 / 139` (carried from v60).
  - **OB-139** (temp files placed at the `tmp/` root instead of a slice directory; four times in P70 / P71 plus once in P75 makes the fifth occurrence; the user ruled on 2026-11-16 to **keep it under observation** and that **a sixth occurrence becomes pending-plus-a-dedicated-slice**): the S-G2 correction [ROUTINE] D=1 was ruled by the user as not counting, and **the P77 main slice recorded D = 0, so the sixth occurrence did not trigger**; OB-139 stays under observation and is counted again on the next [SLICE] with the same pattern.
  - **The reinforced D=2 recurrence-prevention sentence (standard kick-off content since v59, carried into v61)**: "the executing roles (`prfrail-implementer` / `prfrail-tester` / `prfrail-documenter`) **and the controller** must always work under a `tmp/` slice directory (`tmp/<slice-id>/`); **the same holds for CI observation captures, which must never land at the `tmp/` root**; three D occurrences inside one slice force the slice to stop."
- Pending discussion (9): `OB-19 / 21 / 52 / 53 / 57 / 58 / 59 / 61 / 75`.
- To do (1): `OB-69`. Non-terminal total **24** (14 plus 9 plus 1).
- **Slice-level suggestions (three, not registered as OBs)**:
  1. **`R-ACCUMULATION-MECHANISM` (carried from v58)**: P74 was the first case to trigger the §9.2 rule (R>=2 means high risk) together with the §7.5 high-risk tier (one plus three); later implementation slices should pre-place the dynamic budget sentence in their kick-off message; P75 / P76 / P77 recorded R = 1, below the threshold, so it did not trigger. **No OB is registered now**; a recurrence in a later slice turns it into a formal OB.
  2. **`--progress` covers state transitions only, not log output (carried from v58)**: if long-running observability becomes a real need (the S1 Exit closeout slices are the most likely trigger), evaluate adding `--progress verbose` or `--tail-logs`.
  3. **The "obligation hand-off" convention (in use since v59, first applied to P75 O-7)**: whenever a slice hands an obligation to a future slice, that slice's fact package must carry a hand-off section, the landing slice entry in §4 must carry an "admission condition" section, and §5 may track it as an optional "pending execution list" (without an OB). **No OB is registered now**; if a hand-off is ever left dangling it becomes a formal OB.
- **Removed in this version**: none.
- **Flipped to `settled` in this version**: none.
- Mechanical readings: max OB in §12.4 = **139**; **max ADR = 47**; max DR = 14; max deferral = round 25.

---

## §6 Pending user rulings

- **This table is currently empty** (P77 raised no new ruling item; the OB-139 caliber was ruled on 2026-11-16 and P77 did not trigger the sixth occurrence).
- **Anticipated next ruling**: the ordering and priority of the remaining S1 Exit closeout slices (the Release plus ZIP binding and rerun slice P78 and the pre-S2 independent-security-review slice P79).

---

## §7 / §8 / §9 Pointer sections

Same structure as v60: §7 (kick-off and step-0 general rules) points at directive §13.1; §8 (cross-slice collaboration discipline) at §13.2; §9 (checklist meta-rules and gate pitfalls) at §13.3 / §13.4. **The rule bodies have moved into directive §13 and are not copied here.**

---

## §10 One-line opening for a new session

> Below is my project state snapshot (v61). Check the eight anchors in section 1 first; if they match, continue with the queue in section 4, and if not, the repository state governs.
>
> [paste this checklist v61 in full]
>
> [Supporting documents] `REMAINING_SLICES{,_EN}.md` (P77 closed), `DELIVERY_DIRECTIVE.md` (v1.69), `adr010-approval{,_EN}.md` (the P77 report), `s1-exit-closeout-assessment{,_EN}.md` (the P76 report), `ADR_REGISTER{,_EN}.md` (including ADR-045 / ADR-046 / ADR-047), `DOCUMENTATION_PLAN{,_EN}.md` (including the cross-session hand-off entry pointer) and `CHECKLIST_v60{,_EN}.md` (frozen and superseded by v61).
>
> Current slice: P77 is closed (including the backfill `fc00c4b`). Next: **P78 `RELEASE-ZIP-BINDING` (S-G3 closeout, dependency satisfied)**; after that, P79 `S2-PRE-INDEPENDENT-SECURITY-REVIEW`.

## §11 One-line current state

P77 is closed (`33272be` + `fc00c4b` are on origin and runs `43986942087` + `43997058549` were green on both legs; R = 1, D = 0; **ADR-010 = `APPROVED` (user approval 2026-11-16); S-G4 moves from not passed to passed**; the p57 / `s1-exit` gate tables are untouched; the sixth OB-139 occurrence did not trigger; the t025 double BOM exists at HEAD and waits for a later slice); the directive is v1.69; layer 1 keeps S-G3 / S-G5 / S-G6; next comes **P78 `RELEASE-ZIP-BINDING`** (S-G4 has passed and its dependency is satisfied).