# Task execution checklist (cross-session snapshot v37 · P47 governance-baseline closeout edition)

Generated: 2026-11-02, HEAD = origin/main = `bbd1745` (this slice's main commit), directive v1.60.
Supersedes: replaces v36 and every earlier snapshot outright. The docs-only backfill commit only moves HEAD forward and does not change the anchor semantics in this section.

Changes versus v36: (1) P47 is closed and on the ledger, with the seven breakpoint-B verdicts established; (2) the first layer-1 item is dequeued, leaving "layer 1 with two items / layers 2 / 3 / 4"; (3) OB-135 joins the observation list (25 non-terminal items in total); (4) two discipline candidates added (both from G6-4 false positives); (5) R-11 moves to completed.

## 1. Anchors (check these five first)

| Item | Value | Check command |
|---|---|---|
| HEAD | `bbd1745` | `git rev-parse --short HEAD` |
| origin/main | `bbd1745` | `git rev-parse --short origin/main` |
| Directive version | **v1.60** | head of `docs/DELIVERY_DIRECTIVE.md` |
| SELFTEST | **697/697** | `node tools/gates/gate.js --selftest` |
| Three scopes | **TOTAL_FAIL=0 everywhere** | `gate.js --all --scope=index`, `tree`, `ci` |

If two or more items mismatch, treat the repository as authoritative and report the terminal readings before continuing.

Fixed `tmp` assertion: top level holds `.gitkeep` plus `gates-cache`; `gates-cache` holds only `g4b`; `tmp/P47` has been removed.

## 2. P47 fact pack (`GOVERNANCE-BASELINE-VERIFY`)

| Field | Value |
|---|---|
| Slice / size | `GOVERNANCE-BASELINE-VERIFY` (assessment/governance `[SLICE]`) / M |
| Ruling | R-11 = option B plus per-item differentiation (G4 registers only, G5 adds an OB row, G7 registers that the exception did not hold) |
| Verdict spread | closed after backfill: 4 (G4 / G5 / G7 / M6); still open: 3 (P3 / P5 / S7) |
| Backfill | OB-133 to OB-136 added to `§12.4`; pointers appended to three evidence cells in `REMAINING_SLICES` |
| Commit / CI | `bbd1745` / run `41358624169`, both legs green (SELFTEST 697/697, TOTAL_FAIL=0) |
| Accounting | two read-only GLM dispatches (2,972 / 3,480 characters, strictly serial); R = 2, D = 0 |

### 2.1 Seven verdicts

| Item | Verdict | Gap while open |
|---|---|---|
| P3 | Still open | the `unknown` carve-out in REQ-023 versus AT-24 is not cross-referenced; authority awaits R-09 |
| P5 | Still open | CI builds the candidate with the runner and only verifies it with the seed, while RFC R19 requires the seed to build it |
| G4 | Closed after backfill | — |
| G5 | Closed after backfill | — |
| G7 | Closed after backfill | whether to relax the exclusion is carried by a separate slice |
| M6 | Closed after backfill | — |
| S7 | Still open | the three elements now live in OB-136 while the OB-21 row keeps its status |

## 3. Completed slices (recent baseline)

| Version | Slice | Main commit + backfill | Outcome |
|---|---|---|---|
| v1.60 | P44 | `9061b67` + `16abe2b` | closed |
| v1.60 | DR13-FIX | `b123db1` + `eb36992` | closed |
| v1.60 | P45 | `4273d6e` + `3a7541d` | closed |
| v1.60 | P46 | `26ca5a5` + `5e4202e` | closed (R = 1, D = 2) |
| v1.60 | P47 | `bbd1745` + backfill commit | closed (R = 2, D = 0) |

## 4. Execution queue

Principle: layer 1 clears the record chain, then layer 2 the strategic mainline, then layer 3 the pre-release gate, then layer 4 the exit gate; order inside a layer is fixed and layers are never skipped.

### Layer 1 (P47 dequeued, two items left)

| Order | Slice | Identity and content | Size | Step-7 budget | Dependency | Boundary |
|---|---|---|---|---|---|---|
| 1 | DR13-FIX-D3 | cross-process hardening of remove-then-unlock (GAP-71) | S | one Codex final review | R-05 ruling | no DR-13 conclusion change, no CONTRACTS change |
| 2 | RACE-LEG-EXPAND | extend the race leg to console and gates (GAP-42) | M | one Codex final review | R-03 ruling | test coverage only, no production semantics change |

### Layer 2 (strategic mainline)

| Order | Slice | Identity and content | Size | Step-7 budget | Dependency | Boundary |
|---|---|---|---|---|---|---|
| 3 | CONSOLE-MULTITASK-ASSEMBLY | console multi-task production assembly (closes breakpoint A) | L | one Codex final review | layer 1 done plus protocol first | the single-task limit may not be lifted before the protocol changes |
| 4 | ACCEPTANCE-SECURITY-COUNTEREXAMPLES | acceptance-chain security counter-examples (GAP-51 and peers) | M | one Codex final review | item 3 | counter-examples and criteria only, no acceptance semantics change |

### Layers 3 and 4

| Order | Slice | Identity and content | Size | Step-7 budget | Dependency | Boundary |
|---|---|---|---|---|---|---|
| 5 | CONTRACTS-S7-DURABILITY | contract-side durability semantics and Windows C6 landing | L | one Codex final review | R-02 ruling | CONTRACTS first, code second |
| 6 | EXIT-SEED-BUILD-EVIDENCE | AT-13 seed-driven build evidence | M | one Codex final review | item 5 | evidence chain only, no CI structure change |
| 7 | AT23-RESTRICTED-CANDIDATE | AT-23 restricted real candidate and sign-off | M | one Codex final review | R-08 / R-09 rulings | no sign-off without a real-candidate authorisation |

Suspended (not scheduled): `G4B-CORPUS-DECOUPLE` (trigger-based), the five low-priority items (start only on a trigger), and the remaining round-2 gaps (candidate pool).

## 5. Observation items

- Observing (13): `OB-1 / 7 / 10 / 33 / 87 / 99 / 100 / 104 / 105 / 107 / 111 / 113 / 135`.
- To be discussed (11): `OB-19 / 21 / 52 / 53 / 55 / 57 / 58 / 59 / 61 / 75 / 120`.
- To do (1): `OB-69`. Non-terminal total: **25**.
- Mechanical readings: highest OB in `§12.4` = **OB-136**; highest deferral = 25th; ADR max = **28**; DR max = **13** (DR-13 repaired).

## 6. Awaiting user rulings (R-01 to R-11)

| Id | Item | Status |
|---|---|---|
| R-01 | P46 strategic review commissioned | completed |
| R-02 | `CONTRACTS-S7-DURABILITY` cadence and ownership | open (blocks layer 3) |
| R-03 | ownership of the post-P45 slices | open (blocks layer 1 item 2 and layer 2) |
| R-04 | DR-13 "repaired plus observed" reading | open |
| R-05 | whether DR13-FIX D3 is a standalone slice | open (blocks layer 1 item 1) |
| R-06 | OB-55 ownership and T027-B cadence | open |
| R-07 | MVA, the S1 removals and the narrowed S2 promise | open |
| R-08 | candidate, Windows precondition P and probe authorisation | open (blocks layer 4) |
| R-09 | AT-24 `unknown` acceptance boundary | open (blocks P3 closure) |
| R-10 | total budget and per-slice model ceilings | open (blocks all scheduling) |
| R-11 | historical-correction scope | **completed** (option B plus per-item differentiation) |

## 7. Kickoff and step-0 rules

1. Check the anchors, then step 0 registers the ledger (R / D = 0 baseline, temp location `tmp/<sliceId>/`) and lands `input-list.md` and `freedom-list.md`.
2. Every fix slice must ship regression tests; tests or lint not executed must be registered explicitly.
3. Commit only when the three scopes are green, `SELFTEST 697/697` passes and `go build/vet/test ./...` is clean.
4. Model bindings: 1 and 5 use `deep-reasoner` with `DeepSeek V4 Pro (deepseek)`; 2, 3 and 8 use `DeepSeek V4.1 Flash (deepseek)`; 6 uses `independent-reviewer` with `MAI-Code-1.1-Flash (copilot)`; 7 uses `independent-reviewer` with `GPT-5.3-Codex (copilot)`; batch assistance uses `deep-reasoner` with `GLM-5.3-Flash (glm)`. Every `runSubagent` call must pass both `agentName` and `model`.
5. Cadence discipline: never stop in the routine flow; step 10 is the only routine stop point.
6. Batch assistance is dispatched at least twice per slice and **strictly serially**; parallel subagents are forbidden (R1.1 / R5.1).
7. The control seat defaults to `DeepSeek V4.1 Flash`; switch to Pro for low-level concurrency or state-machine work, very long design documents, deep CI-log root-cause analysis, or two consecutive drifts; register any downgrade under OB-37.
8. Role carriers are not substitutable: 2 uses `prfrail-implementer` and 3 uses `prfrail-tester`; the control seat only orchestrates and runs gates.
9. The kickoff message is an artefact: `00-briefing.md` goes into MANIFEST and SHA256SUMS and is declared in the `artifacts` fence.
10. Evidence package and fence stay in step: a new payload needs its MANIFEST entry, recomputed hashes and a fence row in the same round.
11. Reviewer output stays verbatim: 5, 6 and 7 outputs are archived word for word; the control seat must not rebuild them from bullet points.
12. Assessment slices stop at step 10 once step 9 is green; self-issued `commit` or `push` is forbidden.
13. A single output never exceeds 10000 characters; split before the limit and never self-restart.
14. Compaction ends the slice: at 100K of context or on a compaction notice, close the current slice.

## 8. Cross-slice collaboration rules (essentials)

| Element | Content |
|---|---|
| Carrier table | 1, 5, 6 and 7 are read-only; 2, 3 and 8 may write; `prfrail-tester` never runs review stages; batch assistance must use `deep-reasoner` |
| Serial dispatch | `runSubagent` blocks and runs serially; parallelism is forbidden; always pass `agentName` and `model` |
| Commit message language | always English in the form `<type>: <summary>`; use `git commit -F` when non-ASCII appears |
| Temp artefacts | always `tmp/<sliceId>/`, cleaned only by step 8c and only for this slice's directory |
| Gate execution | `--all --skip-base` in intermediate stages and fully green before commit; discipline E: stage precisely, then run `--scope=index` |
| Fence | sorted by path, TAB separated, byte-identical between CN and `_EN` |
| Facts versus judgement | errors in independently checkable output do not void it; judgement items must be redone |

New discipline candidates this version (both from P47's G6-4 false positives):
1. Never write the "number plus counter word" shape in the ledger or reports: a quantifier right after an identifier is read as a line-count literal, so write "the entry" instead.
2. Never quote verbatim the string that tripped a criterion: explaining a red by reproducing the literal trips it again, so use descriptive wording.

## 9. Checklist meta-rules

| Rule (first sentence quoted) | Content |
|---|---|
| "Regenerate every two slices / when pending rulings change / on a directive version bump / when the queue changes / on user request / when two or more anchors mismatch" | regeneration triggers |
| "Pending rulings must each be named" | section 6 names them all |
| "Every queue item must carry six elements" | identity / content / size / step-7 budget / dependency / boundary |
| "A rule obligation is not a gate obligation" | never treat rule text as a mechanical criterion |
| "Gates green before commit" | hard gate |
| "Assessment slices stop at step 10 once step 9 is green" | stop and wait for authorisation when there is no implementation object |
| "Compaction ends the slice" | at 100K or on a compaction notice, close immediately |
| "State each fact in exactly one place" | body text or pointer, never both |
| "Governance slices do not write back to DEV_PLAN" | the authority is appendix C.0 through C.0x |

## 10. One-line opener for a new session

Here is my project state snapshot (v37); check the five anchors in section 1 first, then continue with the queue in section 4, and treat the repository as authoritative on any mismatch.

[paste the full v37 checklist here]

Supporting documents: `REMAINING_SLICES.md` (P47 done), `DELIVERY_DIRECTIVE.md` (v1.60), `docs/validation/governance-baseline-verify{,_EN}.md` (P47 report and the seven verdicts).

Current slice: layer 1 item 1 = `DR13-FIX-D3` (awaiting the R-05 ruling).

## 11. Status in one sentence

P47 is closed (`bbd1745`, both-green run `41358624169`); of the seven breakpoint-B items four are closed after backfill and three stay open (P3 awaits R-09, P5 awaits layer 4, S7 is registered); the next step is layer 1 item 1 `DR13-FIX-D3`, which awaits the R-05 ruling.
