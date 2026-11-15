# Task execution checklist (cross-session snapshot v58 · after P74 closed · the R-3 assessment awaits start)

Generated: 2026-11-15; HEAD = `150423d` (the P74 back-fill commit); directive = **v1.69**; CI main-commit run `43669748483` (P74 main commit) + back-fill run `43672994996` (an observation record only, per §11.2).

Supersedes: this version supersedes v57 (`docs/t027/CHECKLIST_v57{,_EN}.md` is frozen, not retro-edited).

**Structure**: §1 anchors -> §2 fact pack -> §3 completed slices -> §4 execution queue -> §5 observation items -> §6 pending rulings -> §7 / §8 / §9 pointer sections -> §10 opener -> §11 one-line status; the rule entities of §7-§9 have moved into directive §13, so this snapshot does not duplicate them.

Changes in this version (relative to v57): (1) §1 anchors rotated (HEAD `150423d` / directive v1.69 / the rolling window drops v54 and takes in v57); (2) §2 becomes the P74 fact pack (the P72 / P73 fact packs are carried as one-line §3 records plus report paths); (3) §3 gains the P72 / P73 / P74 rows (keeping the eight P64-P71 rows of v57, eleven in total); (4) §4 layer 2 is emptied (P72 / P73 dequeued) and layer 1 T027-B is updated to "only R-3 remains"; (5) §5 slice-level recommendations updated (the `R-ACCUMULATION-MECHANISM` item is added; the `--progress` coverage item is kept; E-1 is now a fixed kick-off 【〇·零】"memory recall" red line - a drafting convention, not formally registered in the directive / an OB / an ADR - and is removed from the recommendations).

## 1. Current anchors (check these 8 numbers first)

| Item | Value | Check command |
|---|---|---|
| HEAD | `150423d` | `git rev-parse --short HEAD` |
| origin/main | `150423d` | `git rev-parse --short origin/main` |
| Directive version | **v1.69** | head of `docs/DELIVERY_DIRECTIVE.md` |
| SELFTEST | **737/737** | `node tools/gates/gate.js --selftest` |
| Three scopes | **all TOTAL_FAIL=0** | `gate.js --all --scope=index`, `tree`, `ci` |
| v55 whole-file lines | **182 (LF count, CN/EN symmetric)** | count `docs/t027/CHECKLIST_v55{,_EN}.md` by LF |
| v56 whole-file lines | **185 (LF count, CN/EN symmetric)** | count `docs/t027/CHECKLIST_v56{,_EN}.md` by LF |
| v57 whole-file lines | **160 (LF count, CN/EN symmetric)** | count `docs/t027/CHECKLIST_v57{,_EN}.md` by LF |

> This snapshot's HEAD = `150423d` is the generation-time anchor after the P74 back-fill commit; once v58 itself lands (commit B), HEAD moves forward, and the repository truth is whatever `git rev-parse --short HEAD` reports then.

> The rolling window is the latest 3 versions (v55 / v56 / v57); v54 has rolled out (same roll-out rule as v52). This snapshot's own line count is carried by v59. If 2 or more anchors disagree => take the repository as truth, report the terminal readings first, then continue.

`tmp` fixed assertion: top level = `.gitkeep` + `gates-cache`; `gates-cache` contains only `g4b`; `tmp/P74-SANDBOX-INTERFACE/` was cleaned at P74's ⑧c.

## 2. P74 fact pack

- Slice / size: `SANDBOX-INTERFACE` (P74; layer-1 part-B introduction slice; implementation `[SLICE]`) / **M**; controller `DeepSeek V4.1 Flash (deepseek)`; same-turn user authorisation on 2026-11-14.
- Goal: land the `Sandbox` check-layer abstraction plus the assembly check points (same layer as the P69 permit, ahead of it); no real enforcement boundary (left to S2).
- Conclusion: **closed** (main commit `3300364` + back-fill `150423d`); step 7 final `FINAL REVIEW: PASS` (0 findings).
- Delivered: `CONTRACTS{,_EN}` §7.2 clauses (a)-(h) + `internal/guard/sandbox{,_windows,_unsupported}.go` + `internal/console/sandbox_preflight.go` + two assembly check points (`runtime_agent.go` / `runtime_supervised.go`, both ahead of the permit) + the trigger `PRFRAIL_SANDBOX_CHECK` (S1 default `off`, zero behaviour change).
- CI: main-commit run `43669748483` (windows 4.9min / ubuntu 3.35min, both legs success); back-fill run `43672994996` (observation record only, §11.2); the package-level `ok internal/guard` / `ok internal/console` lines are present, and the green ubuntu `Test` leg means the `//go:build !windows` test passed on a real Linux host.
- Review chain: (1) `DESIGN: DONE` (including the r2 falsifiability correction) -> (2a) contract -> (3-0) skeleton -> (3) red tests failing -> (2b) implementation -> red-to-green plus mutation check -> (4) gates -> (5) r1 `PASS` (3 Low + 3 Info) -> r2 `PASS` (1 Low + 2 Info, the 1+1 budget exhausted) -> (6) `INDEPENDENT SCAN: PASS` (mandatory, Section A-E complete, Section B no findings, budget 1/3) -> (7) final `PASS` (0 findings; blind review not drawn: none of the three mandatory semantics is touched and the proportional window P72/P73/P74 holds fewer than four slices).
- Metering: **R = 3** (session 1 contributed two: the r1->r2 falsifiability correction and the asymmetric ②a rewrite; session 2 contributed one: the ③0 skeleton shape disagreed with the continuation state file; three distinct root causes, not merged) + **D = 0** (no temp-file violation; OB-139 not repeated); the two batch-assist re-dispatches are process execution (completing session 1's unpersisted work) and **are not counted into R** (registered in report §1.1-c).
- **R=3 trips the §9.2 high-risk marker**: every column of the §5.2 high-risk row was executed (① required / ② required / ③ required / ⑤ required / ⑥ required / ⑦ required / ⑧ required / ⑨ required / ⑧c required); **the R>3 pause threshold was not reached**. **First occurrence** (P72 / P73 had R=0 and did not trip it).
- **The ⑦ budget rises to 1+3 dynamically** (§7.5 tiering rule): tier M set 1+2 at kick-off -> the high-risk row sets 1+3 after R=3; actual use = final 1 + re-review 0 + blind 0, below both ceilings => no overrun; **this is not the §7.5 "no mid-flight tier upgrade" case** (the ⑦ model stayed on the regular `gpt-5.3-codex` tier throughout, never swapped).
- Boundaries: the real enforcement stays in S2; `S-G1` / `S-G8` are not unlocked (P73 condition c landed).
- Segmented execution: across two sessions (session 1: step 0 + ① r2 + ②a; session 2: ②a archival + ③0-⑩), **P74 is one closed slice, not two slices**; §13.1.17 uses the **lenient reading** (context exhaustion = end the session, continue the slice), **acknowledged by the user on 2026-11-15, first use**.
- Deviation numbering: P74-1 segmented execution / P74-2 the ②a archive completed later / P74-3 batch re-dispatch / P74-4 the ⑧ carrier deviation (controller self-execution, all three points present: id + independent impact assessment + must not be a precedent) / P74-5 the ⑦ cross-compile evidence completion - numbered uniformly, an over-and-above registration that does not affect the substantive conclusion.
- Machine readings: ADR max = 44 (P73 ADR-044); OB max = 139 (no new OB); DR max = 14; deferral = #25.
- Evidence pack: `docs/validation/evidence/P74-SANDBOX-INTERFACE/` (24 payloads + `MANIFEST.md` + `SHA256SUMS.txt`, the entry sets equal three ways; `SHA256SUMS` self-check shows 0 mismatches).
- Known limits: the real enforcement stays in S2; the §13.1.17 lenient reading is a first occurrence (report §1.1 + §11); the three uncovered invariants flagged jointly by steps 6 and 7 (no whole-CLI I/O assertion for an illegal environment value / no standalone regression sample for the supervised re-entry path / the non-Windows runtime reading carried by the CI ubuntu leg) are registered in report §11.
- Ledger rule (§11.2): a completion record row carries only the report path plus R/D, **never a hash / run**; the CI fact chain goes into report §14 only.

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
| v1.68 | P72 REQ020-S2-RESTRICTION | `fb0ab57` + `406f4a0` | closed (documentation type; R = 0, D = 0; step 6 skipped with step 7 confirmation; step 7 final PASS WITH FIXES; run `43566904684` + `43568186355`) |
| v1.69 | P73 T027B-ADMISSION | `8c59712` + `f834ce5` | closed (assessment type; R = 0, D = 0; step 6 skipped with step 7 confirmation; step 7 final PASS; **the directive rises to v1.69**; run `43622770750` + `43624661784`; OB-55 = option b) |
| v1.69 | P74 SANDBOX-INTERFACE | `3300364` + `150423d` | closed (implementation type; R = 3, D = 0; step 6 mandatory PASS; step 7 final PASS; **R=3 is the first §9.2 high-risk trip**; run `43669748483` + `43672994996`; the `Sandbox` check layer landed) |

> **The D convention in §3**: "D = n" in this table is shorthand for the temp-file discipline violation count (§5.1) and means the same as "temp-file D" in §2 / §11; the P74-1..P74-5 items in a report's "exception handling" chapter are an independent sequence (§5.1 "naming distinction"), and the two do not convert into each other.

## 4. Execution queue

### Layer 0 pre-work (✅ all dequeued; nothing new in this version)

- P5 / P6 were closed by P57; layer 0 is fully dequeued with nothing left.

### Layer 1 main-line hard gates

- Items 8-12: all ✅ dequeued (see v56 / v57 §4 for detail).
- **`T027-B` (defensive-refactor phase, main-line milestone)**: C1-C7 / N1-N4 were dequeued by P39-P44; R-1 was closed by P69; R-2 was ruled by P73 as option b, and **P74 (the part-B introduction slice, `3300364` + `150423d`) has landed the `Sandbox` abstraction** (the real enforcement stays in S2); **the only remaining item is R-3, the real-CLI end-to-end assessment**.
  - **R-3 start condition**: R-08 / R-09 already ruled (2026-11-06) ✅; once the v58 snapshot closes ✅ the user may name it.
  - **R-3 expected shape**: tier M, assessment type `[SLICE]`, same structure as P73 (step 1 produces the assessment plan and mandates a stop for the user's ruling, then steps 2-10); the slice name is fixed in the kick-off message.
  - **Downstream impact of the R-3 conclusion**: it decides whether the T027-B main-line milestone is declared complete, whether the S-G8 gate advances, and whether the AT-23 stage-1 conditional closure holds.

### Layer 2 routine governance

- P72 / P73 ✅ dequeued (layer 2 is **fully emptied**; see the §3 row summary).
  - **The full §1.2 field definitions of P72 / P73** (slice identifier / goal / dependency / size / steps / current gap / deliverables / acceptance evidence / boundary / tier / probe budget / freedom-list landing path) are in the P72 / P73 sections of `docs/t027/REMAINING_SLICES{,_EN}.md`.

### Layer 3 low-priority candidates

- L-1 low-priority candidates (parked long-term) · ⬜ never started without a trigger

| Field | Content |
|---|---|
| Identifier | `OB-50` / `spawnGit` / `ci freeze cost` / `%TEMP% historical 22 leftovers` / `GATES-EXEMPT` |
| Current status (verified item-by-item) | OB-50 = **disposed** (§12.4, v1.41); spawnGit = the git-call helper of `tools/gates/selftest.js` (the spawn wrapper in `lib/ctx.js`) that captures stderr into the result object instead of printing it; the P27 `G4b-CACHE` report §5 item 5 registered it as a "candidate OB (diagnostics)" (Low / observing) and it was never formally registered as a §12.4 OB number; under the §1.8.1 rule for to-be-assessed items it is an "unregistered candidate OB"; candidate fix = attach the last N bytes of stderr to `detail` when spawn fails (no verdict change, diagnostics only); trigger = the next time any git-related assertion fails for an unknown reason, let step 1 assess it; ci freeze cost = **no matching §12.4 entry**; the real source is the P27 report §7 "candidate OB (performance)" (the origin of the 24-37 s ci-scope freeze cost is undetermined) and it is **not formally registered**; %TEMP% historical 22 leftovers = OB-98 (**disposed**, G9-c static form) / OB-99 (**observing**) / OB-100 (**observing**, the window sub-range already disposed); GATES-EXEMPT = **not chartered** (no §12.4 entry; the P33 item 1-3(iv) split, touching four criterion files) |
| Contains | Five items, each in its corresponding §12.4 row; `GATES-EXEMPT` = the P33 item 1-3(iv) split, not chartered |
| Dependencies | The long-term park covers only undisposed items (spawnGit / ci freeze cost / the %TEMP% timing shape (OB-99 / OB-100) / GATES-EXEMPT); no definite trigger |
| Size | Varies (S through M-L) |
| ⑦ quota | none |
| Boundary | Each item per §12.4; this version adds no layer-3 entry and advances no parked item |
| Start condition | The user must name it explicitly |

**Queue discipline**: same as v57 (layer boundary, numbering discipline, sub-actions A/B are not batched with the snapshot generation: this slice is [ROUTINE] and a single batch; the A / B convention applies only to snapshot generation that bumps the directive, and this slice bumps nothing).

## 5. Observation items

- Observing (**14**): `OB-1 / 7 / 10 / 33 / 87 / 99 / 100 / 104 / 105 / 107 / 111 / 113 / 135` (13 items carried from v57) + **OB-139** (temp-file location violations at the `tmp/` root; four instances across P70 / P71; **P72 / P73 / P74 did not repeat it, so the 5th did not happen**; the state stays `observing`).
- Under discussion (9): `OB-19 / 21 / 52 / 53 / 57 / 58 / 59 / 61 / 75` (**OB-55 was flipped to "disposed" by P73 and left this list**).
- To do (1): `OB-69`. Non-terminal total **24** (14 + 9 + 1).
- **Slice-level recommendations (2 items, no OB registered)**:
  1. **`R-ACCUMULATION-MECHANISM` (recommended tracking, new in this version)**: P74 is the first slice to trip §9.2 "R>=2 => high risk" plus §7.5 "high risk 1+3"; **later implementation kick-off messages should pre-place a dynamic-budget clause** ("if R accumulates to >=2 during execution, apply the §9.2 high-risk marker and the §7.5 high-risk 1+3 budget"), and the dispatch table should pre-place a reminder to upgrade when R reaches 2. **No OB is registered now**; if P75+ repeats the same defect ("no dynamic-budget clause in the kick-off message") => it becomes a formal OB.
  2. **`--progress` covers state transitions only, not log-type output (carried from v57)**: `--progress` prints the chain / task / step state transitions line by line (`progress: <occurredAt> <entity> <fromState>-><toState> <reasonCode>`), but it does **not** forward the live stdout / stderr of gate commands, does **not** consume AgentRunner's `events.jsonl` (tool calls / output summaries), and has **no heartbeat**; P70's original kick-off goal "gate command start/end plus exit codes + AI tool-call summary lines" was narrowed at step 1 because AgentRunner events were not dependable. **Recommendation**: if long-run observability needs later stand out (the R-3 real-CLI end-to-end assessment is the most likely trigger point), assess adding a `--progress verbose` level or a separate `--tail-logs` switch; no OB is registered now, so it is tracked as a slice-level recommendation.
- **Removed in this version**: "E-1: the temp directory re-materialized on disk through the editor buffer after ⑧c" (carried from v57) - **it is now a fixed kick-off 【〇·零】"memory recall" red line** (since P72; P72 / P73 / P74 all carried it; it is a drafting convention, **not yet formally registered in the directive / an OB / an ADR**). None of P72 / P73 / P74 triggered it (each slice's ⑧c second read-back was clean). **E-1 leaves the slice-level recommendations and becomes a kick-off convention**; if P75+ repeats the same shape (the editor buffer re-materializing on disk) => **register it as a new OB**.
- **Flipped to `disposed` in this version**: none (OB-55 was flipped to `disposed` by P73; this version adds no flip).
- Machine readings: §12.4 OB max = **139**; ADR max = **44**; DR max = **14**; deferral max = **#25**.

## 6. Pending user rulings

- This table is currently empty (R-04 / R-06 / R-07 / R-10 were ruled when P68 closed; OB-55 was ruled by P73 as option b; P74 added no ruling items).
- Next expected ruling: after the step-1 deliverable of the R-3 real-CLI end-to-end assessment (the assessment plan) lands, if it touches the maintenance of the AT-23 stage-1 conditional closure or adjusts the S-G8 advance criteria, it may enter this table; to be determined once R-3 starts.

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

This is my project state snapshot (v58); please check the 8 anchors in section 1 first, then continue with the section 4 queue if they match, otherwise take the repository as truth.

[paste this v58 checklist in full]

[Supporting documents] `REMAINING_SLICES{,_EN}.md` (P74 closed), `DELIVERY_DIRECTIVE.md` (v1.69), `sandbox-interface{,_EN}.md` (P74 report), `t027b-admission{,_EN}.md` (P73 report), `req020-s2-restriction{,_EN}.md` (P72 report), `ADR_REGISTER{,_EN}.md` (including ADR-044), `DOCUMENTATION_PLAN{,_EN}.md` (including the cross-session hand-off entry pointer), `CHECKLIST_v57{,_EN}.md` (frozen, superseded by v58).

Current slice: P74 closed; next = **the R-3 real-CLI end-to-end assessment** (tier M assessment type, the slice name to be fixed in the kick-off message; start condition = the v58 snapshot closed + the user names it).

## 11. Current status in one line

P74 closed (`3300364` + `150423d` pushed to origin, runs `43669748483` + `43672994996` green on both legs; R = 3, D = 0; R=3 is the first §9.2 high-risk trip; step 6 mandatory PASS, step 7 final PASS, blind not drawn); directive v1.69 (bumped by P73); the `Sandbox` check layer has landed, the real enforcement stays in S2, and `S-G1` / `S-G8` are not unlocked; next = **the R-3 real-CLI end-to-end assessment** (tier M assessment type, the only remaining T027-B main-line item; starts after the v58 snapshot closes and the user names it).
