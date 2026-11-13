# Task execution checklist (cross-session snapshot v57 · after P71 closed · P72/P73 awaiting the user)

Generated: 2026-11-14; HEAD = `147008a` (sub-action A: directive bumped to v1.68); directive = **v1.68**; CI run `43483740247` green on both legs (windows 4m48s / ubuntu 3m12s).

Supersedes: this version supersedes v56 (`docs/t027/CHECKLIST_v56{,_EN}.md` is frozen, not retro-edited).

**Structure**: §1 anchors -> §2 fact pack -> §3 completed slices -> §4 execution queue -> §5 observation items -> §6 pending rulings -> §7 / §8 / §9 pointer sections -> §10 opener -> §11 one-line status; the rule entities of §7-§9 have moved into directive §13, so this snapshot does not duplicate them.

Changes in this version (relative to v56): (1) §1 anchors rotated (HEAD / directive v1.68 / the rolling window drops v53 and takes in v56); (2) §2 becomes the P71 fact pack; (3) §3 gains the P70 / P71 rows; (4) §4 layer 2 gains the P72 / P73 pending items; (5) §5 observing goes 13 -> 14 (OB-139 added), the satisfied "TUI carve-out condition 1" item is removed, and a slice-level recommendation "`--progress` coverage" is added.

## 1. Current anchors (check these 8 numbers first)

| Item | Value | Check command |
|---|---|---|
| HEAD | `147008a` | `git rev-parse --short HEAD` |
| origin/main | `147008a` | `git rev-parse --short origin/main` |
| Directive version | **v1.68** | head of `docs/DELIVERY_DIRECTIVE.md` |
| SELFTEST | **737/737** | `node tools/gates/gate.js --selftest` |
| Three scopes | **all TOTAL_FAIL=0** | `gate.js --all --scope=index`, `tree`, `ci` |
| v54 whole-file lines | **180 (LF count, CN/EN symmetric)** | count `docs/t027/CHECKLIST_v54{,_EN}.md` by LF |
| v55 whole-file lines | **182 (LF count, CN/EN symmetric)** | count `docs/t027/CHECKLIST_v55{,_EN}.md` by LF |
| v56 whole-file lines | **185 (LF count, CN/EN symmetric)** | count `docs/t027/CHECKLIST_v56{,_EN}.md` by LF |

> This snapshot's HEAD = `147008a` is the generation-time anchor after commit A; once commit B (v57 itself) lands, HEAD moves forward, and the repository truth is whatever `git rev-parse --short HEAD` reports then.

> The rolling window is the latest 3 versions (v54 / v55 / v56); v53 has rolled out. This snapshot's own line count is carried by v58 (same roll-out rule as v52). If 2 or more anchors disagree => take the repository as truth, report the terminal readings first, then continue.

`tmp` fixed assertion: top level = `.gitkeep` + `gates-cache`; `gates-cache` contains only `g4b`; `tmp/P71-S1-SCOPE-ANNOTATE/` was cleaned at P71's ⑧c, and `tmp/V57-CHECKLIST/` is cleaned at this action's ⑧c (after commit B).

## 2. P71 fact pack

- Slice / size: `S1-SCOPE-ANNOTATE` (P71; layer-2 governance side; documentation slice `[SLICE]`) / **M**; controller `DeepSeek V4.1 Flash (deepseek)`; explicit same-turn user authorisation on 2026-11-13.
- Goal: 5 documentation write-backs (REQ stage ownership + `S1_SUPPORT_MATRIX` + `TEST_STRATEGY` + `s1-exit` + `OPERATIONS`) + 3 additive pointers.
- Conclusion: **closed** (X scheme: no entry deleted, no numbering changed, the MVA composite untouched); main commit `6ee2c61` + back-fill `16c98f7`.
- Step 7: blind `BLIND REVIEW: PASS` + final `FINAL REVIEW: PASS` (0 findings, 0 re-reviews).
- CI: run `43463089586` (windows 5m6s / ubuntu 3m14s) + run `43468806314` (windows 4m27s / ubuntu 2m58s), both green on both legs; the package-level `ok internal/console` lines are present and `FAIL` is 0.
- Review chain: (1) `DESIGN: DONE` (including F-1 stopping to report the facts -> user ruled option A -> design-r2) -> (2) steps a-g -> (3) consistency (controller self-executed, carrier deviation registered) -> (4) gates -> (5) `PRE-REVIEW: PASS WITH FIXES` (5 Low: 2 fixed + 3 turned into content assertions) -> (6) **skipped** (the three §7.2.1 conditions hold, confirmed by step 7) -> (7) blind 1 + final 1.
- Metering: step 5 = 1; step 6 = 0 (skipped); step 7 = blind 1 + final 1 + re-review 0; GLM-5.3-Flash batch 2 times; **R = 0, D = 2** (inside P71, not inheriting P70's R=1 / D=2).
- D=2: step 2a execution carrier `tmp/p71-2a-annotate.js` + step 2d `tmp/p71d-wt` worktree; both left zero residue; the user ruled to accept them; **the recurrence-prevention clause is in report §8** (execution roles under `tmp/` must use `tmp/<sliceID>/`, which P72's kick-off discipline must cite).
- Low-3/4/5: landed in report §5 (turned into the four content assertions X-1 / X-2 / X-3 / X-4 in §4).
- **F-1 registration**: item 4 of section [V] of the kick-off message conflicted with the P68 ruling (REQ-019 explicitly keeps self-hosting out of the removal) => step 1 stopped and reported the facts => the user ruled option A (REQ-009 changes, REQ-007 changes at sub-item level, REQ-021 stays, REQ-019 stays).
- 5 write-back landing points: `PRODUCT_REQUIREMENTS{,_EN}.md` (REQ-007 / REQ-009 stage columns), `S1_SUPPORT_MATRIX{,_EN}.md` (Windows row constraint + tail note), `TEST_STRATEGY{,_EN}.md` (AT-13 / AT-24 owner columns), `s1-exit{,_EN}.md` (AT-13 / AT-24 evidence columns + applicability note), `OPERATIONS{,_EN}.md` (visible debugging observation + three boundaries).
- 3 additive pointers: `remaining-rulings` §5, `p57-closeout-assessment` §4.2, `s1-stream-progress` §8 (original lines unchanged, §11.1).
- Boundaries: the directive body / criterion bodies / `CONTRACTS` / schema / fixtures / CI / v50-v56 / the REQ table structure / the MVA composite / the S1 Exit eight-gate verdicts themselves are unchanged.
- Machine readings: ADR max = **43**; OB max = **139** (added by sub-action A); DR max = **14** (DR-14 added by P70); deferral = **#25**.
- Evidence pack and artifacts: `docs/validation/evidence/P71-S1-SCOPE-ANNOTATE/` (9 payloads + `MANIFEST.md` + `SHA256SUMS.txt`, `G8-b` self-consistent); the report's artifacts fence lists 31 `repo` paths, aligned three ways with step 9's `changed=31` and with "31 paths staged precisely, no `-A`".
- Counterexample readings (report §4): `REQ_count=28` (REQ-001-028 contiguous); the four stage values measured one by one; the `PC-01` / `PC-06` rows identical; the gate table identical (2 BLOCKED / 4 PASS / 3 PARTIAL); `S-G3` / `S-G4` / `S-G5` hits = 0; visible "release support" hits = 0; every CN / EN pair symmetric in numstat (nine pairs).
- Known limits (report §7): the `s1-exit` gate table has nine rows with no `S-Gx` row identifiers (defined in `p57-closeout-assessment` §4.2; this slice builds no mapping); the `S-G2` wording is "the gate-level verdict awaits the v54 snapshot"; the REQ-020 "S2 commitment narrowing registration" is outside P71's scope (moved to P72); `S1_SUPPORT_MATRIX` has no TUI row (expressed through the Windows row constraint plus the tail note).
- Ledger rule (§11.2): a completion record row carries only the report path plus R/D, **never a hash / run**; the CI fact chain goes into report §6 only.

## 3. Completed slices (recent baseline)

| Version | Slice | Main commit | Conclusion |
|---|---|---|---|
| v1.66 | P64 AT24-SUPERVISED-BLACK-BOX | `2dbc54a` | closed (AT-24 conditionally closed; R = 2; run `42948727363` green on first push) |
| v1.66 | P65 GATE-PIT-26-LANDING | `0eda074` | closed (Appendix B.7.1 slot + A-part ledger todo; R = 1; run `42969802819`) |
| v1.66 | P66 CHECKLIST-RULES-HOIST | `1650d5f` | closed (§13 five subsections + new snapshot structure; R = 3, D = 2; run `43022796716`) |
| v1.67 | P67 GATES-DEBT-CLEANUP | `18b6ef9` + `365a930` | closed (g4b LRU + gate pitfalls 27-30 + §12.4 settlement + B.6 hardening; R = 3, D = 1; run `43124686702`) |
| v1.67 | P68 REMAINING-RULINGS | `e0d745c` + `665a31f` | closed (assessment type; R = 1, D = 1; step 7 PASS; five R items ruled; ADR-042; run `43177025096`) |
| v1.67 | P69 T027B-E2E-WIRING | `759c934` | closed (implementation type; R = 0, D = 1; step 7 PASS; E-1 registered; run `43264324421`) |
| v1.68 | P70 S1-STREAM-PROGRESS | `ca88fa3` | closed (implementation type; R = 1, D = 2; step 7 PASS; DR-14; run `43400311736` -> `43403891205` -> `43408438643`) |
| v1.68 | P71 S1-SCOPE-ANNOTATE | `6ee2c61` + `16c98f7` | closed (documentation type; R = 0, D = 2; step 7 blind PASS + final PASS; run `43463089586` + `43468806314`) |

> **The D convention in §3**: "D = n" in this table is shorthand for the temp-file discipline violation count (§5.1) and means the same as "temp-file D" in §2 / §11; the D-1..D-n in a report's "exception handling" chapter is an independent sequence (§5.1 "naming distinction"), and the two do not convert into each other.

## 4. Execution queue

### Layer 0 pre-work (✅ all dequeued; nothing new in this version)

- P5 / P6 were closed by P57; layer 0 is fully dequeued with nothing left.

### Layer 1 main-line hard gates (items 8-12 dequeued; only T027-B left to start)

- Items 8-12: all ✅ dequeued (see v56 §4 for detail).
- **`T027-B` (defensive-refactor phase, main-line milestone)** = ⬜ pending start, depending on **P73 being closed**.
- Start condition = the user names it plus the P73 admission-assessment conclusion; boundary = advance along the P37/P38 design baseline and the completed P39-P44 shards + the P69 wiring already landed; it must not start before P73 is closed.

### Layer 2 routine governance (P69 / P70 / P71 dequeued; P72 / P73 pending start)

- P69 `T027B-E2E-WIRING` ✅ dequeued (`759c934`; R = 0, D = 1; step 7 PASS; ADR-043).
- P70 `S1-STREAM-PROGRESS` ✅ dequeued (`ca88fa3`; implementation type; R = 1, D = 2; step 7 PASS; DR-14 registered).
- P71 `S1-SCOPE-ANNOTATE` ✅ dequeued (`6ee2c61` + `16c98f7`; documentation type; R = 0, D = 2; step 7 blind + final both PASS).
- **P72 `REQ020-S2-RESTRICTION`** ⬜ pending start (tier S `[SLICE]`; the 13 fields are in the ledger §1.2 table):
  - Content = register the S2 commitment narrowing for REQ-020: the stage-ownership column becomes "not required by S1 (C+Go contract only); S2 commitment narrowed (Python/Node official harness + three-language product tasks + the release-support commitment frozen; existing tool scripts / open language IDs / generic capability are not frozen)".
  - Sync surface = AT-01 in `TEST_STRATEGY` + the language rows of `S1_SUPPORT_MATRIX`.
  - Boundary = only the REQ-020 stage column plus the related documents; no entry deleted, no numbering changed, the "acceptable requirements" column and the MVA composite untouched.
  - Dependency / start condition = P68 and P71 closed; started when the user names it.
- **P73 `T027B-ADMISSION`** ⬜ pending start (tier M `[SLICE]`; the 13 fields are in the ledger §1.2 table):
  - Content = a three-way assessment of where OB-55 belongs (a standalone slice / landing the `Sandbox` interface early / keeping it in S2) + the part-B implementation boundary + candidate versions + the resource and probe authorisation needs.
  - Stop point = **after step 1 produces the three-way comparison plus a recommendation, stop and wait for the user to rule on OB-55, then push steps 2-10**.
  - Boundary = no production code, only the assessment plus the ruling registration; no change to the T027-B definition beyond `DEV_PLAN`.
  - Dependency / start condition = P72 closed; the user names it.

### Layer 3 low-priority candidates

#### L-1 - low-priority candidates (parked long-term) · ⬜ never started without a trigger

| Field | Content |
|---|---|
| Identifier | `OB-50` / `spawnGit` / `ci freeze cost` / `%TEMP% historical 22 leftovers` / `GATES-EXEMPT` |
| Current state (checked item by item) | OB-50 = **disposed** (§12.4, v1.41); spawnGit = the git-call helper of `tools/gates/selftest.js` (the spawn wrapper in `lib/ctx.js`) that captures stderr into the result object instead of printing it; the P27 `G4b-CACHE` report §5 item 5 registered it as a "candidate OB (diagnostics)" (Low / observing) and it was never formally registered as a §12.4 OB number; under the §1.8.1 rule for to-be-assessed items it is an "unregistered candidate OB"; candidate fix = attach the last N bytes of stderr to `detail` when spawn fails (no verdict change, diagnostics only); trigger = the next time any git-related assertion fails for an unknown reason, let step 1 assess it; ci freeze cost = **no matching §12.4 entry**; the real source is the P27 report §7 "candidate OB (performance)" (the origin of the 24-37 s ci-scope freeze cost is undetermined) and it is **not formally registered**; %TEMP% historical 22 leftovers = OB-98 (**disposed**, G9-c static form) / OB-99 (**observing**) / OB-100 (**observing**, the window sub-range already disposed); GATES-EXEMPT = **not chartered** (no §12.4 entry; the P33 item 1-3(iv) split, touching four criterion files) |
| Contents | Five items, each in its corresponding §12.4 row; `GATES-EXEMPT` = the P33 item 1-3(iv) split, not chartered |
| Dependency | The long-term park covers only undisposed items (spawnGit / ci freeze cost / the %TEMP% timing shape (OB-99 / OB-100) / GATES-EXEMPT); no definite trigger |
| Size | Varies (S through M-L) |
| Step-7 budget | none |
| Boundary | Each item per §12.4; this version adds no layer-3 entry and advances no parked item |
| Start condition | The user must name it explicitly |

**Queue discipline**:

- Layer boundary: P72 / P73 must be named explicitly by the user and never start automatically.
- Numbering discipline: same as v56 (P numbers contiguous, no gaps; after dequeuing the status word stays in sync with the ledger).
- Sub-actions A and B are not batched with the snapshot generation: the directive bump (A) and the snapshot generation (B) are two separate commits, and each one's CI must be green on both legs before the next step.

## 5. Observation items

- Observing (**14**): `OB-1 / 7 / 10 / 33 / 87 / 99 / 100 / 104 / 105 / 107 / 111 / 113 / 135` (13 items carried from v56) + **OB-139** (new in this version: temp-file location violations at the `tmp/` root, four instances across P70/P71, first state word `observing`).
- Under discussion (10): `OB-19 / 21 / 52 / 53 / 55 / 57 / 58 / 59 / 61 / 75`; OB-55 moved into the T027-B boundary with the R-06 = b ruling (first state word unchanged).
- To do (1): `OB-69`. Non-terminal total **25** (14 + 10 + 1).
- **Slice-level recommendations (2 items, no OB registered)**:
  1. **E-1: the temp directory re-materialized on disk through the editor buffer after ⑧c** (carried from v56) - keep executing the discipline "re-read and clean once more between ⑧c and the ⑩ receipt"; if it recurs across slices (2nd time) => it becomes a formal OB.
  2. **`--progress` covers state transitions only, not log-type output** (new in this version, source P70) - `--progress` prints the chain / task / step state transitions line by line (`progress: <occurredAt> <entity> <fromState>-><toState> <reasonCode>`), but it does **not** forward the live stdout / stderr of gate commands, does **not** consume AgentRunner's `events.jsonl` (tool calls / output summaries), and has **no heartbeat**; P70's original kick-off goal "gate command start/end plus exit codes + AI tool-call summary lines" was narrowed at step 1 because AgentRunner events were not dependable. **Recommendation**: when P72 / P73 start, if long-run observability needs stand out, assess adding a `--progress verbose` level (consuming the `tool` / `output` / `process` events of `agent-runner-replay/runs/<requestId>/events.jsonl`) or a separate `--tail-logs` switch; no OB is registered now, so it is tracked as a slice-level recommendation.
- **Removed in this version**: "the TUI carve-out condition 1 not yet landed = long-run usability degradation" (B1: P70 landed `--progress`, so the condition is satisfied).
- **Flipped to `disposed` in this version**: none.
- Machine readings: §12.4 OB max = **139**; ADR max = **43**; DR max = **14**; deferral max = **#25**.

## 6. Pending user rulings

- This table is currently empty (R-04 / R-06 / R-07 / R-10 were ruled when P68 closed on 2026-11-12; P70 / P71 added no ruling items).
- Next expected ruling: the OB-55 three-way ownership choice at P73's step 1 (it enters this table then).

## 7. Startup and step 0 general rules (pointer section)

> These rules were migrated into the directive as **§13.1 startup and step-0 general rules** (20 items, effective since v1.66); this snapshot does not duplicate the rule entities.
> Historical full text: `docs/t027/CHECKLIST_v53.md` §7 (frozen, not retro-edited).

## 8. Cross-slice collaboration discipline (pointer section)

> These rules were migrated into the directive as **§13.2 cross-slice collaboration discipline** (26 items, effective since v1.66, including disciplines E-X / Z and the discipline-Y pointer); this snapshot does not duplicate the rule entities.
> Historical full text: `docs/t027/CHECKLIST_v53.md` §8 (frozen, not retro-edited).

## 9. Checklist meta-rules and gate pitfalls (pointer section)

> The checklist meta-rules were migrated into the directive as **§13.3 checklist meta-rules** (25 items, effective since v1.66; the first item says to cite the opening sentence instead of "item N").
> Gate pitfalls 1-30 were migrated into the directive as **§13.4 gate pitfalls** (three columns: number / symptom / disposition; symptoms and dispositions kept verbatim); item 26 was formalized by P66 and items 27-30 by P67.
> Historical full text: `docs/t027/CHECKLIST_v53.md` §9 (frozen, not retro-edited).

## 10. One-line opener for a new session

This is my project state snapshot (v57); please check the 8 anchors in section 1 first, then continue with the section 4 queue if they match, otherwise take the repository as truth.

[paste this v57 checklist in full]

[Supporting documents] `REMAINING_SLICES{,_EN}.md` (P71 closed), `DELIVERY_DIRECTIVE.md` (v1.68), `s1-stream-progress{,_EN}.md` (P70 report), `s1-scope-annotate{,_EN}.md` (P71 report), `ADR_REGISTER{,_EN}.md` (including ADR-042, ADR-043).

Current slice: P71 closed; next = **P72 `REQ020-S2-RESTRICTION`**, started when the user names it.

## 11. Current status in one line

P71 closed (`6ee2c61` + `16c98f7` pushed to origin, runs `43463089586` / `43468806314` green on both legs; R = 0, D = 2; step 7 blind + final both PASS); next = P72 `REQ020-S2-RESTRICTION` when the user names it -> P73 `T027B-ADMISSION` (after P72 closes, step 1 stops for the user to rule on OB-55) -> the T027-B main line unfolds slice by slice.
