# Task execution checklist (cross-session snapshot v55 · post-P68 closeout)

Generated: 2026-11-12, HEAD = `665a31f` (the P68 back-fill commit), directive v1.67, CI run `43177025096` (P68 main commit) green on both legs. This snapshot and its accompanying commit are the commit right after `665a31f`; that hash cannot be known at generation time, and once it lands HEAD moves forward one step - take the live value from `git rev-parse --short HEAD` at that point.
Supersedes: v54 and every older snapshot in full (historical snapshot files are kept, not deleted). A docs-only back-fill commit only advances HEAD and does not change the meaning of this section's anchors.
**Structure note**: this version keeps the post-§13-migration structure - §7 / §8 / §9 are pointer sections (to directive §13.1 / §13.2 / §13.3 / §13.4) and rule entities are not duplicated here (§11.2, "one fact described in one place"). Snapshots v50-v54 are not retro-edited (§11.1).

Changes versus v54: (1) **P68 recorded** (assessment type, closed); (2) **§2 fact pack switched to P68**; (3) **§4 layer 2 gains three items P69 / P70 / P71**, and `T027-B` is stated explicitly as a main-line milestone; (4) **§5 gains one observation** (the TUI carve-out condition not yet landed = long-run usability degradation); (5) **§6 cleared** (R-04 / R-06 / R-07 / R-10 ruled).

## 1. Current anchors (check these 8 numbers first)

| Item | Value | Check command |
|---|---|---|
| HEAD | `665a31f` | `git rev-parse --short HEAD` |
| origin/main | `665a31f` | `git rev-parse --short origin/main` |
| Directive version | **v1.67** | head of `docs/DELIVERY_DIRECTIVE.md` |
| SELFTEST | **737/737** | `node tools/gates/gate.js --selftest` |
| Three scopes | **all TOTAL_FAIL=0** | `gate.js --all --scope=index`, `tree`, `ci` |
| v52 whole-file lines | **277 (LF count, CN/EN symmetric)** | count `docs/t027/CHECKLIST_v52{,_EN}.md` by LF |
| v53 whole-file lines | **292 (LF count, CN/EN symmetric)** | count `docs/t027/CHECKLIST_v53{,_EN}.md` by LF |
| v54 whole-file lines | **180 (LF count, CN/EN symmetric)** | count `docs/t027/CHECKLIST_v54{,_EN}.md` by LF |

> The rolling window is the latest 3 versions (v52 / v53 / v54); v51 has rolled out. This snapshot's own line count is carried by v56 (same roll-out rule as v51). If 2 or more anchors disagree => take the repository as truth, report the terminal readings first, then continue.

`tmp` fixed assertion: top level = `.gitkeep` + `gates-cache`; `gates-cache` contains only `g4b`; `tmp/P68-REMAINING-RULINGS/` was cleaned at ⑧c.

## 2. P68 fact pack

| Field | Value |
|---|---|
| Slice / tier | `REMAINING-RULINGS` (assessment-type `[SLICE]`, layer-2 governance side) / M |
| Goal | (1) confirm the stage-progression seal (verify in-repo where the A1 ruling of 2026-11-09 landed); (2) for each of the four R items (R-04 / R-06 / R-07 / R-10): current-state anchors + candidate options + impact analysis + advice |
| Outcome | **Closed** (assessment type); five rulings landed: R-04 = a / R-06 = b / R-07 = b1 plus conditions / R-10 = a / OB-11 = user exception plus ADR-042 |
| Main / back-fill commits | `e0d745c` (21 files +3591/−0) + `665a31f` (2 files +16/−0) |
| CI | run `43177025096` (main commit, both legs success) + `43178253595` (back-fill, both legs success) |
| Review chain | (1) `DESIGN: DONE` → (3) `CONSISTENCY: PASS` (34 anchors machine-checked, 1 Low line-number fix) → (4) three scopes `TOTAL_FAIL=0` → (5) r1 `PASS` (1 Low remediated) → (6) **skipped** (§7.2.1 three conditions, confirmed by (7)) → (7) r1 `FINAL REVIEW: PASS` (0 findings) |
| Metering | (7) = final 1 + re-review 0 + blind 0; (5) = 1; (6) = 0 (skipped); GLM-5.3-Flash 2 runs (2418 / 1428 chars); **R = 1, D = 1** |
| Rulings (user, 2026-11-12) | R-04 = a (definitional closeout, keeping "reopen on recurrence"); R-06 = b (close the P44 end-to-end wiring R-1 first, then start T027-B); R-07 = b1 plus conditions (X-style deferral; condition 1 = the TUI carve-out presupposes a `run --progress` addition inside S1; condition 2 = the visible debugging guidance goes into `OPERATIONS` with three boundaries); R-10 = a (keep §7.5; no directive edit) |
| Key finding | **The CLI is fully silent**: `run` prints six summary lines only when finished, has no per-task/per-step lines, no `--verbose` / `--progress` and no heartbeat; the AI child stdout/stderr go to `logs/` files only; `events/state-events.jsonl` is the structured single total-order source but has no subscription API. **visible is the only observable channel during long runs** (and per ADR-012 it is not authoritative evidence) |
| Evidence | Report pair `docs/validation/remaining-rulings{,_EN}.md`; evidence pack `docs/validation/evidence/P68-REMAINING-RULINGS/` (16 items, including `01-arch/r07-refined.md`); ADR-042 |
| Boundary | No change to the directive body / criteria / product code / CONTRACTS / schema / fixtures / CI / v50-v54; no ruling made on the user's behalf |
| Machine readings | ADR max = **42**; OB max = **138**; DR max = **13**; pilot deferral = **#25** (P68 is assessment-type, so no new deferral was registered) |

## 3. Completed slices (recent baseline)

| Version | Slice | Main commit | Outcome |
|---|---|---|---|
| v1.66 | P64 AT24-SUPERVISED-BLACK-BOX | `2dbc54a` | closed (AT-24 conditional closure; R = 2; run `42948727363` first-run green) |
| v1.66 | P65 GATE-PIT-26-LANDING | `0eda074` | closed (Appendix B.7.1 slot plus the A-part ledger item; R = 1; run `42969802819`) |
| v1.66 | P66 CHECKLIST-RULES-HOIST | `1650d5f` | closed (§13 five subsections plus the new snapshot structure; R = 3, D = 2; run `43022796716`) |
| v1.67 | P67 GATES-DEBT-CLEANUP | `18b6ef9` + `365a930` | closed (g4b LRU + gate pitfalls 27-30 + §12.4 settlement block + B.6 hardening; R = 3, D = 1; run `43124686702`) |
| v1.67 | P68 REMAINING-RULINGS | `e0d745c` + `665a31f` | closed (assessment type; R = 1, D = 1; (7) PASS; five R rulings landed; ADR-042; run `43177025096`) |

> **D convention for §3**: each "D = n" here is shorthand for the temporary-artefact discipline violation count (§5.1) and means the same as "temp-artefact D" in §2 / §11; the D-1..D-n entries in a report's deviation table are a separate series (§5.1 "naming distinction") and the two are never converted into each other.

## 4. Execution queue

**Ordering basis**: layer 0 / layer 1 hard gates and follow-up candidates are all dequeued; layer 2 gains three items after P68 (P69 / P70 / P71), all derived from the R-07 conditions and the R-06 ruling; `T027-B` is a main-line milestone whose sub-slices are named one by one after P69 closes.

### Layer 0 pre-work (✅ all dequeued; nothing new in this version)

| # | Item | Content | Status |
|---|---|---|---|
| 0-1 | P5 `T027-B` closeout list | verdict: no new slice; the remainder splits into R-1 / R-2 / R-3 | ✅ dequeued by P57 |
| 0-2 | P6 `S1-EXIT-DOC` closeout list | the S1 Exit eight-gate list (S-G1 to S-G8) | ✅ dequeued by P57 |

### Layer 1 main-line hard gates (all ✅ dequeued)

| # | Slice / milestone | Status | Note |
|---|---|---|---|
| 8-11 | AT23-RESTRICTED-CANDIDATE / R-08.5 revision / F2 fix / P63 re-check | ✅ dequeued | stage-1 conditional closure (A1 ruling 2026-11-09) |
| 12 | AT24-SUPERVISED-BLACK-BOX | ✅ dequeued | AT-24 conditional closure (P64) |
| - | **`T027-B` (defensive refactor phase, main-line milestone)** | ⬜ not started | tier M/L; **depends on P69 closing**; starts only when the user names it; boundary = proceed per the P37/P38 design baseline and the completed P39-P44 slices. **Not a single P number**; register it per §1.8.2 when it starts |

### Layer 2 routine governance (three new items after P68)

#### P69 — `T027B-E2E-WIRING` (R-1 closeout) · ⬜ not started

| Field | Content |
|---|---|
| Identifier | `T027B-E2E-WIRING` (layer-2 governance side; source = R-1 in P57 report §3.3) |
| Goal | Wire the P44 readiness / permit verdicts into the production assembly path and the CLI entry, closing R-1 |
| Dependency | P44 (closed); P67 / P68 (closed) |
| Size | M |
| Steps | (1) current-state check (the seam between P44 verdict outputs and the production assembly) → (2) wiring → (3) tests → (4) gates → (5) pre-review → (6) scan → (7) final review → (8) → (9) → (8c) → (10) |
| Current gap | The P44 verdict implementation and tests are done; the end-to-end wiring is not (P57 §3.3 R-1) |
| Deliverables | Wiring implementation + tests + report pair; evidence pack |
| Acceptance evidence | The production assembly path actually consumes the P44 verdicts; related tests and gates all green |
| Boundary | Do not change the P44 verdict body or the `CONTRACTS` semantics; wiring only; no gitee push |
| Tier | `[SLICE]` |
| Probe budget | 0 (local machine, no external calls) |
| Freedom-list path | `docs/validation/evidence/P69-T027B-E2E-WIRING/freedom-list.md` |

#### P70 — `S1-STREAM-PROGRESS` (TUI carve-out condition 1) · ⬜ not started

| Field | Content |
|---|---|
| Identifier | `S1-STREAM-PROGRESS` (layer-2 governance side; source = R-07 b1 condition 1) |
| Goal | Add a lightweight `run --progress` stream switch inside S1: off by default, consuming the existing event stream, printing per-task/step start and end plus gate command start/end and exit codes plus AI tool-call summary lines |
| Dependency | P68 (closed, R-07 condition 1) |
| Size | S |
| Steps | (1) event-surface check → (2) switch and printing → (3) tests → (4) gates → (5) → (6) → (7) → (8) → (9) → (8c) → (10) |
| Current gap | `run` has no live output at all (six summary lines only when finished); no `--verbose` / `--progress`; no heartbeat |
| Deliverables | `internal/console` changes + tests + report pair; evidence pack |
| Acceptance evidence | Default behaviour unchanged; with the switch on, per-task/step lines are visible; `--json` semantics unchanged; gates all green |
| Boundary | `console` + tests; no new event model; no view layer; no change to `CONTRACTS` / `--json` semantics; no gitee push |
| Tier | `[SLICE]` |
| Probe budget | 0 |
| Freedom-list path | `docs/validation/evidence/P70-S1-STREAM-PROGRESS/freedom-list.md` |

#### P71 — `S1-SCOPE-ANNOTATE` (TUI carve-out condition 2 plus the stage-ownership column) · ⬜ not started

| Field | Content |
|---|---|
| Identifier | `S1-SCOPE-ANNOTATE` (layer-2 governance side; source = R-07 b1 condition 2 and the X-style write-back) |
| Goal | (1) write the visible debugging guidance into `OPERATIONS{,_EN}` with three boundaries; (2) mark the stage-ownership column of REQ-009 / REQ-019 / REQ-007 / REQ-021 as "not required by S1, deferred to S2"; (3) align `S1_SUPPORT_MATRIX` with the TUI carve-out and the visible degraded guarantee; (4) align the AT-13 / AT-24 items in `TEST_STRATEGY`; (5) add the gate-applicability note to `s1-exit` |
| Dependency | P68 (closed, R-07 condition 2); order-wise recommended after P69 and P70 |
| Size | M (documentation only; touches product-side release wording) |
| Steps | (1) current-state and difference check → (2) per-document write-back → (3) consistency check → (4) gates → (5) → (6) → (7) → (8) → (9) → (8c) → (10) |
| Current gap | The S1 scope and stage ownership are not yet written per R-07 b1; the visible guidance is missing |
| Deliverables | The five document write-backs above + report pair; evidence pack |
| Acceptance evidence | Stage-ownership marks in place with no entry deleted; the three visible boundaries present; `S1_SUPPORT_MATRIX` and `TEST_STRATEGY` consistent; gates all green |
| Boundary | No change to REQ numbering or item count (X-style); no change to the MVA acceptance composite; no change to directive §7.5 / §6.3 criteria; no gitee push; **starting it needs an explicit same-turn user authorization, it is not `[ROUTINE]`** |
| Tier | `[SLICE]` |
| Probe budget | 0 |
| Freedom-list path | `docs/validation/evidence/P71-S1-SCOPE-ANNOTATE/freedom-list.md` |

### Layer 3 low-priority candidates

| # | Item | Content | Start condition |
|---|---|---|---|
| L-1 | Five low-priority tasks | not started without a trigger | explicit user naming only |

**Layer boundaries**: layer 0 / layer 1 (including dequeued items) / dequeued layer-2 items produce no further assessment actions; the three new layer-2 items start only when the user names them and never automatically; `T027-B` is a main-line milestone rather than a single P number, its sub-slices are named one by one after P69 closes and take P numbers after P71 (P72, P73...); layer 3 must not be read as "the main line resuming".

**P-numbering discipline**: P numbers are contiguous and only grow; a slice name never changes once it is in the ledger (§13.3.10). Milestone numbers (e.g. P39-P44 of T027-B) and slice numbers (from P45) are two separate series; a slice always starts under a new P number, never under a milestone sub-task number.

## 5. Observation items

- Observing (13): `OB-1 / 7 / 10 / 33 / 87 / 99 / 100 / 104 / 105 / 107 / 111 / 113 / 135` (13 items carried from v54, no change in this version).
- Under discussion (10): `OB-19 / 21 / 52 / 53 / 55 / 57 / 58 / 59 / 61 / 75`; **OB-55 moved into the T027-B boundary with the R-06 = b ruling** (first state word unchanged, still under discussion).
- To do (1): `OB-69`. Non-terminal total **24** (13 + 10 + 1).
- **New in this version: 1 slice-level recommendation (no OB registered)**: **the TUI carve-out condition 1 not yet landed = long-run usability degradation** - until `run --progress` lands, the TUI carve-out means users still get no live progress during long runs; when `P70` starts it should become a formal observation item with either "accept the degradation" or "land the condition" concluded.
- **Flipped to `disposed` in this version**: none (P68 closed no OB; ADR-042 registers an OB-11 exception, it does not flip an OB).
- Machine readings: `§12.4` OB max = **138**; deferral max = **#25**; ADR max = **42**; DR max = **13**.

## 6. Pending user rulings

| Id | Item | Status |
|---|---|---|
| - | - | **this table is currently empty** |

> R-01 / R-02 / R-03 / R-05 / R-08 / R-09 / R-11 / R-12 were ruled or completed earlier; **R-04 / R-06 / R-07 / R-10 were ruled on 2026-11-12 when P68 closed**; this table is currently empty.

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

This is my project state snapshot (v55); please check the 8 anchors in section 1 first, then continue with the section 4 queue if they match, otherwise take the repository as truth.

[paste this v55 checklist in full]

[Supporting documents] `REMAINING_SLICES.md` (P68 closed), `DELIVERY_DIRECTIVE.md` (v1.67, including §13), `docs/validation/remaining-rulings{,_EN}.md` (P68 report), `docs/validation/evidence/P68-REMAINING-RULINGS/01-arch/r07-refined.md` (the full R-07 advice), `docs/ADR_REGISTER{,_EN}.md` (including ADR-042).

Current slice: P68 closed; next = **P69 `T027B-E2E-WIRING`**, started when the user names it.

## 11. Current status in one line

P68 closed (`e0d745c` + `665a31f` pushed to origin, runs `43177025096` / `43178253595` green on both legs; five R rulings landed; ADR-042; R = 1, D = 1; (7) PASS); next = **P69 `T027B-E2E-WIRING`** when the user names it → P70 `S1-STREAM-PROGRESS` → P71 `S1-SCOPE-ANNOTATE` → `T027-B` (expanded slice by slice after P69 closes).
