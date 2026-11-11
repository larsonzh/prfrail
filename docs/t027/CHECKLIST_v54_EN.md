# Task Execution Checklist (cross-session snapshot v54 · §13 migration completed)

Generated: 2026-11-11, HEAD = origin/main = `1650d5f` (the P66 commit), directive v1.66, CI run `43022796716` (P66) green on both legs. This snapshot and its accompanying commit are the commit right after `1650d5f`; that hash cannot be known at generation time, and once it lands HEAD moves forward one step — take the live value from `git rev-parse --short HEAD` at that point.
Supersedes: fully replaces v53 and every earlier snapshot (historical snapshot files are kept, not deleted). A backfill commit (docs-only) only moves HEAD forward and does not change the semantics of this section's anchors.
**Structure note**: this is the **first snapshot under the new structure** after the §13 migration — §7 / §8 / §9 become pointer sections (pointing at directive §13.1 / §13.2 / §13.3 / §13.4) and rule entities are no longer copied into this snapshot (directive §11.2, "describe the same fact in exactly one place"). Historical snapshots v50–v53 are never retro-edited (directive §11.1).

Changes in this version (relative to v53): (1) **backfill of P64 / P65 / P66**, the three slices closed after v53 but missing from its queue; (2) **§2 fact pack replaced with P66**; (3) **§4 queue gains P67 `GATES-DEBT-CLEANUP`**, and layer-1 follow-up candidate item 12 is ✅ dequeued; (4) **§5 flips two observation items to `已处置`** (the G5(d) item-26 pitfall candidate and the in-snapshot content inflation) and **adds five debt observations**; (5) **§6 narrows to four pending rulings**; (6) **§7 / §8 / §9 become pointer sections**.

## 1. Current anchors (verify these 8 numbers first)

| Item | Value | Check command |
|---|---|---|
| HEAD | `1650d5f` | `git rev-parse --short HEAD` |
| origin/main | `1650d5f` | `git rev-parse --short origin/main` |
| Directive version | **v1.66** | header of `docs/DELIVERY_DIRECTIVE.md` |
| SELFTEST | **722/722** | `node tools/gates/gate.js --selftest` |
| Three scopes | **all TOTAL_FAIL=0** | `gate.js --all --scope=index` / `tree` / `ci` |
| v51 whole-file size | **270 (LF count, CN/EN symmetric)** | `docs/t027/CHECKLIST_v51{,_EN}.md` counted by LF |
| v52 whole-file size | **277 (LF count, CN/EN symmetric)** | `docs/t027/CHECKLIST_v52{,_EN}.md` counted by LF |
| v53 whole-file size | **292 (LF count, CN/EN symmetric)** | `docs/t027/CHECKLIST_v53{,_EN}.md` counted by LF |

> v52 size calibre note (effective from v54): v53 §1 recorded the v52 whole-file size as 278, a reading error at generation time; from v54 on the count is uniformly by LF, and v52 is actually 277 (CN/EN symmetric). v53 itself is never retro-edited (historical snapshots are frozen).

If 2 or more items disagree, take the repository as ground truth and report the final readings before continuing.

`tmp` fixed assertion: top level = `.gitkeep` + `gates-cache`; `gates-cache` contains only `g4b`; `tmp/P66-CHECKLIST-RULES-HOIST/` was cleaned at 8c.

## 2. P66 fact pack

| Field | Value |
|---|---|
| Slice / size | `CHECKLIST-RULES-HOIST` (governance `[SLICE]`, layer-2 governance side) / L |
| Goal | Structurally migrate the execution discipline already frozen in v53 §7 / §8 / §9 into a new directive §13 "Cross-slice collaboration and startup general rules", promoting soft constraints to hard constraints |
| Verdict | **closed**; the five §13 subsections are on disk (§13.1 startup / §13.2 cross-slice discipline / §13.3 checklist meta-rules / §13.4 gate pitfalls / §13.5 snapshot structure template); ⑦ final `FINAL REVIEW: FINDINGS` (F-01 = the ⑥ formal defect) then `RE-REVIEW: PASS` after remediation |
| Migration verdict | rule count **82** (source rows 98, with 14 same-rule multi-row merges and 2 non-rule rows): fully duplicated 19 / partially overlapping 43 / snapshot-only 20 ⇒ **migrated 70 / turned into pointers 11 / not migrated 1** (discipline Y, already promoted to §1.7.2); plus one backfilled citation meta-rule from just before the §9 table, landing in §13.3.1 |
| §13 item readings | §13.1 = 20 / §13.2 = 26 / §13.3 = 25 / §13.4 = 26 (pitfalls 1–26, three-column table) / §13.5 = 8 retained sections + 4 pointer sections + 5 point-in-time data classes |
| Pitfall item 26 | Formalized by this slice (it was a v53 §5 candidate) ⇒ this also closes P65's "A-part carry-over" |
| Change surface | directive `DELIVERY_DIRECTIVE{,_EN}.md` (§13 five subsections, the v1.66 row in §0.3, Appendix C.0at, the Appendix B.7.1 slot, Appendix B.7.4 reading set, header version line); report pair `docs/validation/checklist-rules-hoist{,_EN}.md`; evidence pack `evidence/P66-CHECKLIST-RULES-HOIST/` (25 payloads, MANIFEST and SHA256SUMS self-consistent); ledger P66 block |
| Review chain | ① `DESIGN: DONE` (four artefacts complete) → ② steps A–G in strict series (each independently verified) → ③ consistency check (2 CONFIRMED DEFECTs → r2 fix) → ④ gates → ⑤ `PASS WITH FIXES` → `RE-REVIEW: PASS` → ⑥ three rounds all formally defective (downgraded to corroboration) → ⑦ final `FINDINGS` → `RE-REVIEW: PASS` |
| Commit / CI | `1650d5f`, run `43022796716` green on the first attempt on both legs (windows / ubuntu 15 package-level `ok` each) |
| Metering | ⑦ = 1 final + 1 re-review + 0 blind (within the 1+3 allowance); ⑤ = 2 runs (allowance exhausted); ⑥ = 3 runs (allowance exhausted); GLM 1 run; R = 3, D = 2 |
| Rulings | controller Flash (no Pro switch, no ADR); the user's ⑩ accepted D-1 / D-2 / R2-1 deviations and every gate-hit correction; the ⑥ downgrade was endorsed with no re-run |
| Capability boundary | no change to the criterion body / CONTRACTS / schema / fixtures / product code / CI; no change to directive §0–§12 body text; no change to v50–v53; pitfall 26 symptom and disposition kept verbatim (only 3 registered rewrites) |
| Not closed | g4b cache governance, the §12.4 sampling-record debt and ⑥ prompt hardening — all carried into P67 in §4 |

## 3. Completed slices (recent baseline)

| Version | Slice | Main commit | Verdict |
|---|---|---|---|
| v1.60 | P49 | `75990a8` + backfill commit | closed (R = 1, D = 2; run `41445657595` green on the first attempt) |
| v1.60 | P50 | `81f9dc9` + backfill commit | closed (R = 1, D = 1; run `41507722672`) |
| v1.60 | P51 | `2b0af2d` + backfill commit | closed (R = 1, D = 2; run `41569353460`) |
| v1.60 | P52 | `c636dd8` + backfill commit | closed (R = 0, D = 1; run `41642139932`; GAP-77 breach reproduced) |
| v1.60 | P53 | `7802819` + backfill commit | closed (R = 0, D = 0; run `41707160563`; GAP-77 fixed) |
| v1.60 | P54 | `7c0763a` + `f38e0ce` | closed (R = 1, D = 1; run `41840557379`; C6 ruling path + downgrade ruling) |
| v1.60 | P55 | `63e2e24` + `c76c6ce` | closed (R = 0, D = 1; runs `41925623071` + `41928327517` green; AT-13 two-sided evidence chain) |
| v1.61 | P56 | `504caff` | closed (R = 2, D = 1; run `41978355255` green; G4b corpus enumeration decoupled from the index, SELFTEST 697 → 703) |
| v1.61 | P57 | `7a43dfa` | closed (evaluation; R = 0, D = 0; run `41991695581` green; P5 not a slice + P6 eight-gate list) |
| v1.62 | P58 | `f798cea` + `8d30e66` | closed (evaluation; R = 0, D = 0; runs `42062391497` + `42069357349` green; P7 not a slice, row kept) |
| v1.63 | P59 | `75b90a8` | closed (governance; R = 0, D = 0; run `42121119502` green on the first attempt; R-12 ruling A2×B1×Y1 landed) |
| v1.64 | P60 | `e381898` + `dafe004` | closed (governance; R = 0, D = 2; runs `42323986343` + `42328605588` green; J2×J3 landed) |
| v1.64 | P61 | `6bb23a1` + `05be20b` | closed (acceptance, PARTIAL close-out; R = 0, D = 1; run `42452377967` green; ADR-036) |
| v1.64 | P62 | `d0971e5` | closed (governance; R = 0, D = 0; run `42505811938` green on the first attempt; R-08.5 five elements landed; ADR-037) |
| v1.64 | P63 | `de37c58` | closed (adjudication; R = 0, D = 0; run `42531794673` both legs success; ⑦ ruled the F2 exemption invalid; ADR-038) |
| v1.64 | F2-FIX | `57631a3` | closed (fix; R = 0, D = 0; run `42568763867` both legs green; AT-23 stage 1 conditionally closed; ADR-039) |
| v1.66 | P64 AT24-SUPERVISED-BLACK-BOX | `2dbc54a` | closed (AT-24 conditionally closed; R = 2; ⑦ final PARTIAL → re-review `RE-REVIEW: PASS`; run `42948727363` green on the first attempt) |
| v1.66 | P65 GATE-PIT-26-LANDING | `0eda074` | closed (Appendix B.7.1 slot + A-part ledger carry-over; directive v1.65; R = 1; run `42969802819` green on the first attempt) |
| v1.66 | P66 CHECKLIST-RULES-HOIST | `1650d5f` | closed (the five §13 subsections + the new snapshot structure; directive v1.66; R = 3, D = 2; ⑦ final `FINDINGS` → re-review `RE-REVIEW: PASS`; run `43022796716` green on the first attempt) |

> **D calibre in §3**: the "D = n" in each row abbreviates the count of temporary-artefact discipline violations (§5.1) and means the same as "temp-file D" in §2 / §11; the D-1..D-n entries in a report's "exception handling" deviation table form an independent series (§5.1 "naming distinction") and the two are never converted into each other.

## 4. Execution queue

**Ordering basis for this version**: layer 0 / layer 1 main hard gates / layer-1 follow-up candidate item 12 (AT-24) are **all dequeued**; layer 2 gains P67 `GATES-DEBT-CLEANUP` (four debt classes) awaiting an explicit user call; layer 3 stays not started.

### Layer 0 pre-work (✅ done, dequeued)

| No. | Item | Content | Size | Dependency | Boundary |
|---|---|---|---|---|---|
| 0-1 | P5 `T027-B` close-out list | ✅ done: ruled **no new slice**; remaining items split three ways R-1 / R-2 / R-3 | evaluation action | — | dequeued; deliverables in P57 report §3 |
| 0-2 | P6 `S1-EXIT-DOC` close-out list | ✅ done: S1 Exit eight-gate list (S-G1 to S-G8) | evaluation action | — | dequeued; deliverables in P57 report §4 |

> Calibre: both items were closed out and dequeued by P57; their conclusions are the **only source of truth** for layer 1 and S1 Exit, and no same-shape evaluation action may be raised again.

### Layer 1 main hard gates (all ✅ dequeued)

| No. | Slice | Identity and content | Size | ⑦ allowance | Dependency | Boundary |
|---|---|---|---|---|---|---|
| 8 | AT23-RESTRICTED-CANDIDATE | ✅ dequeued (stage 1 conditionally closed; stage 2 not applicable, user ruling A1 of 2026-11-09) | M | — | P61 assessment done; P62 revised R-08.5; P63 adjudication; F2-FIX closed | dequeued |
| 9 | R-08.5 criterion revision slice | ✅ done and dequeued (P62 closed, `d0971e5`) | S | — | — | dequeued |
| 10 | F2 fix slice | ✅ done and dequeued (F2-FIX closed, `57631a3`) | S | Codex 1 final (+1 blind) | P63 ruled the exemption invalid | dequeued; only `tools/agent-stub` changed |
| 11 | AT-23 re-check under the P62-revised criterion (P63) | ✅ done and dequeued (P63 closed, `de37c58`) | S | Codex 1 final (+1 blind) | P62 closed | dequeued; adjudication + ledger/report updates only |

### Layer 1 follow-up candidates (all ✅ dequeued)

| No. | Slice | Content | Size | ⑦ allowance | Dependency | Boundary |
|---|---|---|---|---|---|---|
| 12 | AT24-SUPERVISED-BLACK-BOX (AT-24) | ✅ dequeued (P64 closed, `2dbc54a`; AT-24 conditionally closed) | M | Codex 1 final + blind judged per §7.7 | user call to start; AT-23 stage 1 conditionally closed | dequeued; product code untouched |

### Layer 2 routine governance (debt cleanup)

| No. | Item | Content | Start condition |
|---|---|---|---|
| G-1 | P7 governance-side OB roll-up | ✅ done and dequeued: closed by P58, verdict = no slice, P7 row kept at ⏸ | dequeued |
| G-3 | `DIRECTIVE-SCOPE-GUARD` (P59) | ✅ done and dequeued: closed by P59 (2026-11-06), R-12 ruling landed in v1.63 | dequeued |
| G-4 | `DIRECTIVE-JUDGE-BLINDSPOTS` (P60) | ✅ done and dequeued: closed by P60 (2026-11-07), J2×J3 landed in v1.64 | dequeued |
| G-5 | P65 `GATE-PIT-26-LANDING` | ✅ done and dequeued: B-part landed in v1.65; A-part completed by P66 (pitfall 26 formalized) | dequeued |

### Layer 2 addition (P67 · awaiting an explicit user call)

| No. | Slice | Content | Size | ⑦ allowance | Dependency | Boundary |
|---|---|---|---|---|---|---|
| 13 | GATES-DEBT-CLEANUP | One-shot cleanup of four debts: #1 g4b cache LRU; #2 formalize the pitfall-27 candidate (four pitfalls merged); #3 clear the §12.4 sampling-record debt; #4 harden the ⑥ prompt | L | 1+3 | P66 closed (already satisfied) | no change to existing directive §1–§12 clauses; no change to v50–v53; for #1 the criterion-body change is limited to tightening the LRU and does not alter existing G4b semantics; for #4, if it touches the §7.2.2 failure-disposition semantics, stop and escalate |

### Layer 3 low-priority options

| No. | Item | Content | Start condition |
|---|---|---|---|
| L-1 | five low-priority tasks | no trigger, no start | explicit user call only |

> Layer boundary: dequeued items in layer 0 / layer 1 (including follow-up candidates) / layer 2 produce no further evaluation actions; P67 requires an explicit user call and never auto-starts; layer 3 must not be read as "the main line resuming".

## 5. Observation items

- Observing (13): `OB-1 / 7 / 10 / 33 / 87 / 99 / 100 / 104 / 105 / 107 / 111 / 113 / 135`.
- Pending discussion (10): `OB-19 / 21 / 52 / 53 / 55 / 57 / 58 / 59 / 61 / 75`.
- Pending action (1): `OB-69`. Non-terminal total **24**.
- GAP registry: GAP-77 fixed (P53 closed); GAP-78/79/80 are candidates (registered in P52 report §2b and the ledger).
- **Flipped to `已处置` in this version (2)**:
  1. **G5(d) pitfall item-26 candidate** — formalized as item 26 in directive §13.4 by P66 (symptom and disposition kept verbatim); P65's "A-part carry-over" is closed at the same time.
  2. **In-snapshot frozen-content inflation** — resolved by the P66 §13 migration: the rule entities of §7 / §8 / §9 moved into the directive, and this snapshot keeps pointers only.
- **Added in this version (5, all found at P66 close-out, all → P67)**:
  1. **g4b cache has no LRU or cap**: P64 90 → P65 94 → P66 100, a 3–6 `.json` increment per slice, with no eviction. → P67 #1
  2. **The §12.4 proportional-sampling record block has not advanced since v1.32**: a debt that degrades v54's point-in-time accuracy. → P67 #3
  3. **⑥ independent-scan formal defects (all three P66 rounds invalid)**: r1 did not use the A–E shape (voided), r2 used A–E but two Medium findings were stale readings, r3 used `## A. Scope` instead of `## Section A` (downgraded to low-weight corroboration after the allowance was exhausted). → P67 #4
  4. **Three wording traps**: `忌` (a G4b unregistered character), the aggregate-count wording (the G6-4 aggregate-count shape), `候选` (read as a point-in-time word after being migrated into a formal table). → P67 #2
  5. **Line-tail truncation trap (P66 empirical)**: a CN insertion pushed the tail of an existing line into a standalone line (`⑨ 前重同步。`), producing CN 220 / EN 219 asymmetry; same family as the `\n\n` swallowing trap. → P67 #2
- Mechanical readings: largest OB in `§12.4` = **OB-138**; largest deferral = the 25th; ADR max = **39**; DR max = **13**.

## 6. Pending user rulings (R-04 / R-06 / R-07 / R-10)

| Id | Item | Status |
|---|---|---|
| R-04 | DR-13 "fixed + observed" calibre | pending |
| R-06 | OB-55 ownership and the T027-B cadence | pending |
| R-07 | MVA, S1 exit items and narrowing the S2 commitment | pending |
| R-10 | total budget and per-slice model ceilings | pending (blocks scheduling) |

> Decided / completed items no longer occupy this table: R-01 / R-02 / R-03 / R-05 / R-11 are complete; R-08 / R-09 were decided on 2026-11-06 (see `docs/validation/r08-r09-ruling{,_EN}.md`); R-12 was decided on 2026-11-06 (A2×B1×Y1, landed in v1.63).
> P66 ruling note (2026-11-11): the three ⑥ formal defects are confirmed, the layer is downgraded to low-weight corroboration and is **not re-run**; the D-1 / D-2 / R2-1 deviations and all gate-hit corrections were accepted; R-04 / R-06 / R-07 / R-10 remain pending.

## 7. Startup and step 0 general rules (pointer section)

> These rules were migrated by P66 into directive **§13.1 Startup and step-0 general rules** (effective from v1.66, 20 items); this snapshot no longer copies rule entities.
> The historical full text is in `docs/t027/CHECKLIST_v53.md` §7 (frozen, never retro-edited).

## 8. Cross-slice collaboration discipline (pointer section)

> These rules were migrated by P66 into directive **§13.2 Cross-slice collaboration discipline** (effective from v1.66, 26 items, including disciplines E–X / Z and the discipline-Y pointer); this snapshot no longer copies rule entities.
> The historical full text is in `docs/t027/CHECKLIST_v53.md` §8 (frozen, never retro-edited).

## 9. Checklist meta-rules and gate pitfalls (pointer section)

> The checklist meta-rules were migrated by P66 into directive **§13.3 Checklist meta-rules** (effective from v1.66, 25 items, whose first item is "cite by the opening quotation, not by 'item N'").
> Gate pitfalls 1–26 were migrated by P66 into directive **§13.4 Gate pitfalls** (three-column table: number / symptom / disposition; symptoms and dispositions kept verbatim); item 26 was formalized by P66.
> The historical full text is in `docs/t027/CHECKLIST_v53.md` §9 (frozen, never retro-edited).

## 10. One-line opener for a new session

This is my project state snapshot (v54); verify the eight anchors in section 1 first, then continue with the queue in section 4 if they match, and otherwise take the repository as ground truth.

【paste this checklist in full, v54】

【supporting documents】`REMAINING_SLICES.md` (P66 closed), `DELIVERY_DIRECTIVE.md` (v1.66, with the new §13), `docs/validation/checklist-rules-hoist{,_EN}.md` (P66 report, §11 carries the ⑥ downgrade record), `docs/validation/r08-r09-ruling{,_EN}.md` (R-08/R-09 rulings).

Current slice: P66 is closed; every dequeued item in layers 0 / 1 / 2 is complete; next = P67 `GATES-DEBT-CLEANUP` on an explicit user call.

## 11. Current status in one line

P66 closed (`1650d5f` pushed to origin, run `43022796716` green on both legs; the five §13 subsections plus the new snapshot structure; directive v1.66; R = 3, D = 2; the ⑥ formal defect downgraded); next = P67 `GATES-DEBT-CLEANUP` awaiting an explicit user call.
