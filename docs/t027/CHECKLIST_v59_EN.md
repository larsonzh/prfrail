# Task execution checklist (cross-session snapshot v59 · after P75 closes · S1 Exit closure assessment pending)

Generated: 2026-11-16; HEAD = `60146c0`; directive = **v1.69**; CI main run `43732384194` (P75 main commit) + backfill run `43734829063` (P75 backfill) + supplement run `43739409108` (P75 supplement; all three double green; the last two are observation records only per §11.2).

Supersession: this version supersedes v58 (`docs/t027/CHECKLIST_v58{,_EN}.md` frozen, not rewritten).

**Structure note**: §1 anchors -> §2 fact pack -> §3 completed slices -> §4 execution queue -> §5 observations -> §6 pending rulings -> §7 / §8 / §9 pointer sections -> §10 one-line opening -> §11 one-line status; the rule bodies of §7-§9 moved into directive §13, so this snapshot does not copy them.

This version's changes (versus v58): (1) §1 anchors refreshed (HEAD `60146c0` / directive v1.69 / rolling window drops v55 and admits v58); (2) §2 becomes the P75 fact pack plus **Annex A "obligation hand-off"** and **Annex B "AT-23 stage-2 state"**; (3) §3 appends the P75 row (the eleven v58 rows plus one = twelve rows); (4) §4 layer 1 updates T027-B to "complete (partly not applicable)" plus full dequeue plus a new "S1 Exit closure" item (admission condition = the single P75 §11.2 residual is fixed first); (5) §5 registers the 5th OB-139 occurrence plus the strengthened D=2 recurrence-prevention clause plus slice-level advice item 3 "obligation hand-off convention".

## 1. Current anchors (check these 8 numbers first)

| Item | Value | Check command |
|---|---|---|
| HEAD | `60146c0` | `git rev-parse --short HEAD` |
| origin/main | `60146c0` | `git rev-parse --short origin/main` |
| Directive version | **v1.69** | `docs/DELIVERY_DIRECTIVE.md` header |
| SELFTEST | **737/737** | `node tools/gates/gate.js --selftest` |
| Three scopes | **all TOTAL_FAIL=0** | `gate.js --all --scope=index` / `tree` / `ci` |
| v56 whole-file line count | **185** (LF count, CN/EN symmetric) | `docs/t027/CHECKLIST_v56{,_EN}.md` by LF |
| v57 whole-file line count | **160** (LF count, CN/EN symmetric) | `docs/t027/CHECKLIST_v57{,_EN}.md` by LF |
| v58 whole-file line count | **148** (LF count, CN/EN symmetric) | `docs/t027/CHECKLIST_v58{,_EN}.md` by LF |

> This snapshot's HEAD = `60146c0` is the generation anchor after the P75 supplement commit; once v59 lands, HEAD moves forward and `git rev-parse --short HEAD` is authoritative.

> Rolling window = the latest 3 versions (v56 / v57 / v58); v55 has been dropped (same rule as v52); this snapshot v59's own line count continues in v60. If 2 or more items mismatch, take the repository as authoritative and report the final readings first.

`tmp` fixed assertion: the top level = `.gitkeep` + `gates-cache`; `gates-cache` holds only `g4b`; `tmp/P75-R3-REAL-CLI-ASSESSMENT/` was cleaned at P75 ⑧c.

## 2. P75 fact pack

- Slice / size: `R3-REAL-CLI-ASSESSMENT` (P75; layer-1 T027-B main-line milestone residual assessment; assessment `[SLICE]`) / **M**; controller `DeepSeek V4.1 Flash (deepseek)`; user authorised in the same round on 2026-11-15.
- Goal: answer the three R-3 residual questions of the T027-B main-line milestone (can it be declared complete / can S-G8 advance / is the closed AT-23 stage-1 state maintained); no production code, no AT-23 stage-1 re-run.
- Conclusion: **closed** (main commit `ed2e26e` + backfill `70a04cd` + supplement `60146c0`); core ruling = **R-3 = (A) cannot start for now** (same-round user ruling of 2026-11-15).
- **A1**: T027-B declaration = "**complete (partly not applicable)**" - the main line C1-C7 / N1-N4 / R-1 / R-2 has fully landed; R-3 is partly not applicable (real CLI end-to-end needs the S2 enforcement premise and is not done inside S1).
- **A2**: S-G8 = stays blocked until S2 enforcement lands; the `s1-exit{,_EN}.md` file is unchanged (mechanically checked); the S-G8 row of `p57-closeout-assessment{,_EN}.md` §4.2 gets one additive pointer (report: `docs/validation/r3-real-cli-assessment{,_EN}.md`).
- **A3**: the closed AT-23 stage-1 state = **permanent fact registration** (a historical fact, not revoked by the R-3 conclusion); stage 2 is not applicable (user ruling A1 of 2026-11-09); the P61 point record (2026-11-09 F2-FIX) is in report §7; the S-G1 row of `p57` §4.2 is unchanged (the gate-level verdict and the AT-23 stage-1 closure are two different dimensions).
- **A4**: three R-3 landing points - the T027-B section of `DEV_PLAN{,_EN}` (`docs/DEV_PLAN.md:203`), the P75 completion record in `REMAINING_SLICES{,_EN}` (R = 1, D = 1), and the additive pointer on the R-3 row of `p57` §3.3 (`docs/validation/p57-closeout-assessment.md:76`).
- **A5**: **ADR-045** is added (`docs/ADR_REGISTER.md:56` / `_EN:57`) - a **decision record** (same class as ADR-042 / ADR-043 / ADR-044, **not bound by "must not be used as a precedent"**): R-3 = (A) / T027-B = complete (partly not applicable) / S-G8 blocked until S2 / AT-23 stage-1 closure = permanent fact; effective 2026-11-15; changes neither ADR-044, nor R-08 / R-09, nor the eight-gate decisions.
- **O-1 to O-8**: 7 in-slice rulings (O-1 / O-2 / O-3 / O-4 / O-5 / O-6 / O-8) plus 1 pending (O-7); the user ruled on 2026-11-16 that **O-7 = option 2 (no separate slice)**. **The hand-off is in Annex A**.
- **Step-5 proactive verification precedent**: the step-5 r1 reported "the three O-7 items are stale" => the controller re-checked the disk => **P75 §11.2 items 1 / 2** (corresponding to P57 §4.3 original items 1 / 3) were already closed by P64 / s1-reach corrections, and **P75 §11.2 item 3** (corresponding to P57 §4.3 original item 2 = `s1-exit.md:52`) is the sole residue; the entries were corrected per the facts (correct application of §13.3.14).
- CI: main run `43732384194` (ubuntu 3m19s / windows 4m41s, both green); backfill run `43734829063` (observation only); supplement run `43739409108` (ubuntu 3m16s / windows 4m28s, both green).
- Review chain: step 1 `DESIGN: DONE` (three-way + eight items + recommendation (B); the user ruled (A)) -> 2a ruling registration -> 2b `DEV_PLAN` + two `p57` pointers + ADR-045 -> 2c report §5-§11 -> step 3 consistency (carrier deviation P75-3) -> step 4 gates -> step 5 r1 `PASS WITH FIXES` (7 Medium + 2 Low) -> step 5 r2 `PASS WITH FIXES` (3 Low fixed) -> step 6 skipped (step 7 `SKIP CONFIRMED`) -> **step 7 blind review 1 before final review 1** (`BLIND REVIEW: PASS`, 0 findings; the §7.7 proportional draw hit it; widest span = 3 decision faces x 3 downstream items; the three mandatory semantics not touched) -> step 7 final `FINAL REVIEW: PASS` (0 findings) -> step 8 assembly -> step 9 verification -> 8c cleanup.
- Metering: **R = 1** (P75-1, step-5 r1 rework) + **D = 1** (P75-6 temp-file violation); **step 7 = blind 1 + final 1 = 2 calls** (cap 4, not exceeded); R = 1 < 2 => the §9.2 dynamic high-risk upgrade did not trigger.
- Numbered deviations: P75-1 (step-5 r1 rework, R = 1) / P75-2 (authoring-time gate wording fixes, not counted as R) / P75-3 (step-3 carrier deviation) / P75-4 (freedom-list finalisation timing deviation) / P75-5 (subagent temp file, zero residue) / **P75-6 (temp-file violation, D = 1, the 5th OB-139 occurrence)**.
- Boundary: real enforcement stays in S2; `S-G1` / `S-G8` not unlocked; no real CLI end-to-end probe started; `s1-exit` / `s1-reach` / `CONTRACTS` / schema / fixtures / CI / product code unchanged.
- Mechanical readings: **ADR max = 45**; OB max = 139 (none added); DR max = 14; pilot deferral = the 25th time.
- Evidence pack: `docs/validation/evidence/P75-R3-REAL-CLI-ASSESSMENT/` (15 payloads + `MANIFEST.md` + `SHA256SUMS.txt`).
- Known boundaries: real enforcement stays in S2; the O items and the O-7 hand-off are recorded in report §11; step 6 was skipped (confirmed by step 7), so there is no "uncovered invariants flagged jointly by steps 6 and 7".
- Ledger convention (§11.2): the completion record carries only the report path plus R/D (no hash / run); the CI fact chain lives only in report §16.

### Annex A: "obligation hand-off" (new from v59; first case = P75's O-7 = option 2)

- Hand-off content: the **single residual** among the three stale statements of P75 report §11.2 (the inconsistent wording inside the core task-chain row at `s1-exit.md:52`, corresponding to P57 §4.3 original item 2).
- Landing: it must be fixed before the S1 Exit closure slice starts; S1 Exit must not be declared passed before that fix.
- Precedent reference: P75 §11.2 lists all three items plus fix timing, the admission condition, and a three-row comparison table (P75 §11.2 items against the P57 §4.3 originals, CN `:162` / `:164-166`).
- Convention note: **this "obligation hand-off" section is adopted from v59 onward (first case traced to P75 O-7)** - whenever a slice hands an obligation to a future slice, this snapshot's §2 fact pack must carry an "obligation hand-off" section; §4's landing-slice item must carry an "admission condition" section; §5 may optionally track it as "pending actions" (no OB registration).

### Annex B: "AT-23 stage-2 state" (new from v59)

- State: ruling A1 says "not applicable" (user ruling of 2026-11-09; **not a permanent ban** - unless the user later explicitly rules a restart).
- Executable object before S2: none (R-3 = (A); real CLI end-to-end needs real enforcement, an S2 premise).
- Restart path: if the user later rules a stage-2 restart, real enforcement must land first (S2) or an alternative path (P-prime / candidate swap / a new intermediate path) must exist; an alternative path is not "implementing stage 2 before S2" but "reaching an equivalent goal before S2".
- Natural restart point: S2 starts -> enforcement lands -> S-G8 unlocks -> R-3 is re-assessed (A -> possibly B) -> if B is feasible, that is the natural stage-2 restart point.

## 3. Completed slices (recent baseline)

| Version | Slice | Main commit | Conclusion |
|---|---|---|---|
| v1.66 | P64 AT24-SUPERVISED-BLACK-BOX | `2dbc54a` | closed (AT-24 conditional closure; R = 2; run `42948727363` double green on first try) |
| v1.66 | P65 GATE-PIT-26-LANDING | `0eda074` | closed (Annex B.7.1 slot + A-part ledger todo; R = 1; run `42969802819`) |
| v1.66 | P66 CHECKLIST-RULES-HOIST | `1650d5f` | closed (§13 five sections + new snapshot structure; R = 3, D = 2; run `43022796716`) |
| v1.67 | P67 GATES-DEBT-CLEANUP | `18b6ef9` + `365a930` | closed (g4b LRU + gate pits 27-30 + §12.4 clearing block + B.6 hardening; R = 3, D = 1; run `43124686702`) |
| v1.67 | P68 REMAINING-RULINGS | `e0d745c` + `665a31f` | closed (assessment; R = 1, D = 1; step 7 PASS; five R items ruled; ADR-042; run `43177025096`) |
| v1.67 | P69 T027B-E2E-WIRING | `759c934` | closed (implementation; R = 0, D = 1; step 7 PASS; E-1 registered; run `43264324421`) |
| v1.68 | P70 S1-STREAM-PROGRESS | `ca88fa3` | closed (implementation; R = 1, D = 2; step 7 PASS; DR-14; run `43400311736` -> `43403891205` -> `43408438643`) |
| v1.68 | P71 S1-SCOPE-ANNOTATE | `6ee2c61` + `16c98f7` | closed (documents; R = 0, D = 2; step 7 blind PASS + final PASS; run `43463089586` + `43468806314`) |
| v1.68 | P72 REQ020-S2-RESTRICTION | `fb0ab57` + `406f4a0` | closed (documents; R = 0, D = 0; step 6 skipped, confirmed by step 7; step 7 final PASS WITH FIXES; run `43566904684` + `43568186355`) |
| v1.69 | P73 T027B-ADMISSION | `8c59712` + `f834ce5` | closed (assessment; R = 0, D = 0; step 6 skipped, confirmed by step 7; step 7 final PASS; **directive bumped to v1.69**; run `43622770750` + `43624661784`; OB-55 = option b) |
| v1.69 | P74 SANDBOX-INTERFACE | `3300364` + `150423d` | closed (implementation; R = 3, D = 0; step 6 mandatory PASS; step 7 final PASS; **first R=3 case triggering §9.2 high risk**; run `43669748483` + `43672994996`; `Sandbox` check layer landed) |
| v1.69 | P75 R3-REAL-CLI-ASSESSMENT | `ed2e26e` + `70a04cd` + `60146c0` | closed (assessment; R = 1, D = 1; step 6 skipped, confirmed by step 7; step 7 blind 1 + final 1, blind-before-final order compliant; **R-3 = (A)**; T027-B = complete (partly not applicable); ADR-045; run `43732384194` + `43734829063` + `43739409108`) |

> **§3 D convention**: "D = n" here is shorthand for the temp-file discipline violation count (§5.1); among the P75-N entries in the report's exception section, only the temp-file deviation (P75-6) counts as D; the other deviations (P75-1 to P75-5) do not.

## 4. Execution queue

### Layer 0 prerequisite sweep (✅ done; nothing new in this version)

- P5 / P6 were closed by P57; layer 0 is empty with no residue.

### Layer 1 main-line hard gates

- Items 8-12: all ✅ done (see v56 / v57 §4).
- **`T027-B` (defensive refactor period, main-line milestone) = complete (partly not applicable)**: C1-C7 / N1-N4 / R-1 / R-2 / R-3 are all handled (R-3 = (A), user ruling of 2026-11-15; the A1 declaration = "complete (partly not applicable)"). **The main line is fully dequeued**.
- **`S1 Exit closure` (next item)**: the eight gates as in the table below (structure + state + basis, three elements each).

| Gate | Criterion | State | Basis |
|---|---|---|---|
| S-G1 | the CLI can genuinely launch an admitted AI candidate and bill it | stage-1 conditional closure reached; gate-level verdict pending re-assessment | `s1-exit.md:43`; F2-FIX (2026-11-09); P75 §7 |
| S-G2 | `supervised-black-box` implemented and the AT-24 verdict holds | conditionally closed; gate-level verdict pending re-assessment | `s1-exit.md:44` already corrected (P64); P57 §4.3 original item 1 closed |
| S-G3 | the final ZIP is bound to the release record and re-verified | not passed | install-package guide PARTIAL + ADR-009 |
| S-G4 | ADR-010 approved by the product owner | not passed | ADR-010 `PROPOSED / NOT APPROVED` |
| S-G5 | the product owner approves | unknown (human gate; its relationship with `PRODUCT_REQUIREMENTS` §7 "independent security review conclusion" awaits clarification in the S1 Exit closure assessment step-1 stage) | `s1-exit.md:60` |
| S-G6 | dependency fields satisfied and the four remaining gates closed | unreachable | `DEV_PLAN.md:94` / `:96` |
| S-G7 | all ATs pass, T018 complete, §14 fully exited | unreachable (composite gate) | `s1-exit.md:52` / `:56`-`:58` / `:60` |
| S-G8 | the supervised path runs end-to-end on a real CLI and a real candidate | blocked until S2 enforcement lands | P75 A2; ADR-045 |

- Remaining actionable = S-G3 / S-G4 / S-G5 / S-G6 (S-G7 is a composite gate, S-G8 needs S2).
- **Admission condition**: the single residual of P75 report §11.2 (the inconsistent wording inside the core task-chain row at `s1-exit.md:52`, corresponding to P57 §4.3 original item 2) must be fixed first; S1 Exit must not be declared passed before that fix.
- Suggested shape (to be confirmed when the user names it): an S1 Exit closure assessment (M-tier assessment, same shape as P73 / P75), step 1 produces the closure roadmap (the gate-level re-assessment wording for S-G1 / S-G2 plus the unlock conditions for S-G3 / S-G4 / S-G5 / S-G6 plus the O-7 single-residual fix plus the product-owner approval timing) -> stop for the user's ruling; the assessment then drives the execution slices.
- **`S-G8`** = stays blocked until S2 enforcement lands (P73 A2 / P75 A2 / ADR-044 / ADR-045); the `s1-exit{,_EN}.md` file is unchanged (mechanically checked).

### Layer 2 routine governance

- P72 / P73 / P74 / P75 ✅ done (layer 2 is **fully empty**; see the §3 row summaries).
- The full §1.2 field definitions for P72 / P73 / P74 / P75 live in the corresponding sections of `docs/t027/REMAINING_SLICES{,_EN}.md`.

### Layer 3 low-priority backlog

| Field | Content |
|---|---|
| Identifier | `OB-50` / `spawnGit` / `ci frozen cost` / `22 historical %TEMP% leftovers` / `GATES-EXEMPT` |
| Current state (item by item) | OB-50 = **disposed** (§12.4, v1.41); spawnGit = the git helper in `tools/gates/selftest.js` (the spawn wrapper in `lib/ctx.js`) that captures stderr into the result object instead of printing it; the P27 `G4b-CACHE` report §5 item 5 registered it as a "candidate OB (diagnosis)" (Low / observing) and it was never formally registered as a §12.4 OB id; under the §1.8.1 open-item rule it is an "unregistered candidate OB"; candidate fix = attach the last N bytes of stderr to `detail` when spawn fails (no verdict change, diagnostics only); trigger = the next unexplained git-related assertion failure, assessed at step 1; ci frozen cost = **no §12.4 entry**, its real origin is the P27 report §7 "candidate OB (performance)", **unregistered**; the 22 historical %TEMP% leftovers = OB-98 (**disposed**) / OB-99 (**observing**) / OB-100 (**observing**); GATES-EXEMPT = **not raised** (no §12.4 entry; the P33 step-1-3(iv) split, touching 4 criterion files) |
| Contents | the five items are in the §12.4 rows; `GATES-EXEMPT` = the P33 step-1-3(iv) split, not raised |
| Dependency | the long suspension covers only undisposed items (spawnGit / ci frozen cost / %TEMP% timing forms (OB-99 / OB-100) / GATES-EXEMPT); no clear trigger |
| Size | varies (S to M-L) |
| Step-7 quota | none |
| Boundary | item by item in §12.4; this version adds no layer-3 item and advances no suspended item |
| Start condition | must be named explicitly by the user |

**Queue discipline**: same as v58 (layer boundaries, numbering discipline; this slice is `[ROUTINE]` single batch; no directive bump).

## 5. Observations

- Observing (**14**): `OB-1 / 7 / 10 / 33 / 87 / 99 / 100 / 104 / 105 / 107 / 111 / 113 / 135` (13 items, carried from v58) + **OB-139** (temp-file location violation at the `tmp/` root; four times in P70 / P71 + **once in P75 = the 5th occurrence has happened**; the user ruled on 2026-11-16 = **option 1, keep observing**; **a 6th occurrence moves it to pending discussion plus a governance slice**).
  - **Strengthened D=2 recurrence-prevention clause (standard briefing content from v59)**: "execution roles (`prfrail-implementer` / `prfrail-tester` / `prfrail-documenter`) **and the controller** always use `tmp/<slice-id>/` under `tmp/`; **the same applies to CI observation captures (never at the `tmp/` root)**; three D occurrences in one slice force a stop".
- Pending discussion (9): `OB-19 / 21 / 52 / 53 / 57 / 58 / 59 / 61 / 75`.
- To do (1): `OB-69`. Non-terminal total **24** (14 + 9 + 1).
- **Slice-level advice (3 items, no OB registration)**:
  1. **`R-ACCUMULATION-MECHANISM` (carried from v58)**: P74 was the first case triggering §9.2 "R>=2 => high risk" plus §7.5 "high risk 1+3"; **later implementation slices should pre-place the dynamic-quota sentence in their briefings**; P75's R = 1 < 2 did not trigger it. **No OB registration now**; a P76+ repeat of the same defect turns it into a formal OB.
  2. **`--progress` covers state transitions only, not log-like output (carried from v58)**: it prints state transitions line by line but does not forward the live stdout / stderr of gate commands, does not consume AgentRunner `events.jsonl`, and has no heartbeat; if long-task observability becomes prominent later (the S1 Exit closure execution slice is the most likely trigger), evaluate adding `--progress verbose` or `--tail-logs`.
  3. **"Obligation hand-off" convention (adopted from v59, first case P75 O-7)**: whenever a slice hands an obligation to a future slice, the §2 fact pack must carry an "obligation hand-off" section + the §4 landing-slice item must carry an "admission condition" section + §5 may optionally track a "pending actions" list (no OB registration). **No OB registration now**; a later "dangling hand-off" turns it into a formal OB.
- **Removed in this version**: none (E-1 was removed in v58 and is now fixed content of the briefing 'memory wake-up' red line).
- **Flipped to disposed in this version**: none.
- Mechanical readings: §12.4 max OB = **139**; **ADR max = 45**; DR max = 14; max deferral = the 25th time.

## 6. Pending user rulings

- **This table is currently empty** (R-04 / R-06 / R-07 / R-10 were ruled at P68; OB-55 was ruled at P73 = option b; R-3 was ruled at P75 = (A) plus A1-A5; O-7 was ruled on 2026-11-16 = option 2; P75 adds no ruling item).
- **Next expected ruling**: once the S1 Exit closure assessment's step-1 output (the closure roadmap) lands, the user must rule on the closure order of S-G3 / S-G4 / S-G5 / S-G6 and on the product-owner approval timing.

## 7. Kick-off and step-0 general rules (pointer section)

> These rules moved into directive **§13.1 kick-off and step-0 general rules** (effective from v1.66, 20 items); this snapshot does not copy them.
> The historical full text is in `docs/t027/CHECKLIST_v53.md` §7 (frozen, not rewritten).

## 8. Cross-slice collaboration discipline (pointer section)

> These rules moved into directive **§13.2 cross-slice collaboration discipline** (effective from v1.66, 26 items, including disciplines E-X / Z and the discipline Y pointer); this snapshot does not copy them.
> The historical full text is in `docs/t027/CHECKLIST_v53.md` §8 (frozen, not rewritten).

## 9. Checklist meta-rules and gate pits (pointer section)

> The meta-rules moved into directive **§13.3 checklist meta-rules** (effective from v1.66, 25 items, the first being "cite the opening sentence, never 'item N'").
> Gate pits 1-30 moved into directive **§13.4 gate pits** (three columns: id / symptom / disposition; symptom and disposition kept verbatim), with pit 26 formalised by P66 and 27-30 by P67.
> The historical full text is in `docs/t027/CHECKLIST_v53.md` §9 (frozen, not rewritten).

## 10. One-line opening for a new session

> Here is my project state snapshot (v59); check the 8 anchors in section 1 first, continue with the queue in section 4 if they match, and take the repository as authoritative otherwise.
>
> [paste this checklist v59 in full]
>
> [supporting documents] `REMAINING_SLICES{,_EN}.md` (P75 closed), `DELIVERY_DIRECTIVE.md` (v1.69), `r3-real-cli-assessment{,_EN}.md` (P75 report), `sandbox-interface{,_EN}.md` (P74 report), `t027b-admission{,_EN}.md` (P73 report), `req020-s2-restriction{,_EN}.md` (P72 report), `ADR_REGISTER{,_EN}.md` (with ADR-045), `DOCUMENTATION_PLAN{,_EN}.md` (with the cross-session hand-off entry), `CHECKLIST_v58{,_EN}.md` (frozen, superseded by v59).
>
> Current slice: P75 is closed; next = **S1 Exit closure assessment** (M-tier assessment; the slice name is set by its briefing; start condition = v59 snapshot closed plus the user naming it).

## 11. One-line current status

P75 closed (`ed2e26e` + `70a04cd` + `60146c0` pushed to origin, runs `43732384194` + `43734829063` + `43739409108` all double green; R = 1, D = 1; **R-3 = (A)**; **T027-B = complete (partly not applicable)**; **S-G8 = blocked until S2 enforcement lands**; **AT-23 stage-1 closure = permanent fact**; **ADR-045** registered; O-7 = option 2, its single residual handed to the S1 Exit closure slice as an admission condition); directive v1.69; the `Sandbox` check layer has landed and real enforcement stays in S2; next = **S1 Exit closure assessment** (M-tier assessment, roadmap for S-G3 / S-G4 / S-G5 / S-G6; started after the v59 snapshot closes plus the user naming it).
