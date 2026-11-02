# Task execution checklist (cross-session snapshot v38 · P48 closure edition)

Generated: 2026-11-02, HEAD = origin/main = `82e3ac8` (P48 main commit), directive v1.60.
Supersedes: fully replaces v37 and any earlier snapshot. Docs-only backfill commits advance HEAD only; they do not change the anchors of this section.

Changes in this edition (relative to v37): (1) P48 `DR13-FIX-D3` closed into the baseline (CI green on first run); (2) layer-1 first item removed from the queue, layer 1 has one item left (`RACE-LEG-EXPAND`); (3) six P48 gate pitfalls added for P49; (4) observation items and pending rulings unchanged.

## 1. Current anchors (verify these five numbers first)

| Item | Value | Command |
|---|---|---|
| HEAD | `82e3ac8` | `git rev-parse --short HEAD` |
| origin/main | `82e3ac8` | `git rev-parse --short origin/main` |
| Directive version | **v1.60** | header of `docs/DELIVERY_DIRECTIVE.md` |
| SELFTEST | **697/697** | `node tools/gates/gate.js --selftest` |
| Three scopes | **all TOTAL_FAIL=0** | `gate.js --all --scope=index`, `tree`, `ci` |

If two or more items mismatch, follow the repository state: report the terminal readings first and then continue.

`tmp` fixed assertion: top level = `.gitkeep` + `gates-cache`; `gates-cache` holds only `g4b`; `tmp/DR13-FIX-D3` is removed.

## 2. P48 fact pack (`DR13-FIX-D3`)

| Field | Value |
|---|---|
| Slice / size | `DR13-FIX-D3` (implementation `[SLICE]`) / S |
| Goal | `release` removes the path with `os.Remove` before unlocking and closing, closing W1 together with the D1 ownership gate and the D2 stat re-check; full Windows remove-while-open regression |
| Red and green | deterministic AST order gate red on the first run (unlock before remove) plus stress violations 4/35; 6/6 green after the fix |
| Review chain | 5 `PRE-REVIEW: PASS` (2 Low + 3 Info); 6 `INDEPENDENT SCAN: PASS`; 7 `FINAL REVIEW: PASS` (one non-blocking observation) |
| Commit / CI | `82e3ac8` / run `41379998477` green on both legs on the first attempt |
| Metering | R = 9, D = 0; GLM batch assist twice; zero probes |
| Doc sync | `ARCHITECTURE{,_EN}` section 4.2 one lock-order entry each; `DEV_PLAN{,_EN}` has no related milestone and stays unchanged |
| Open items | 5 Low-2 (unix runtime leg covered by the dual-platform CI); 7 observation (theoretical false-green shape of the earliest-call-position order gate); R-N1 fact (an `os.OpenFile` handle blocks removal, errno 32, a designed convergence path) |

## 3. Completed slices (recent baseline)

| Version | Slice | Main commit + backfill | Verdict |
|---|---|---|---|
| v1.60 | P45 | `4273d6e` + `3a7541d` | closed |
| v1.60 | P46 | `26ca5a5` + `5e4202e` | closed (R = 1, D = 2) |
| v1.60 | P47 | `bbd1745` + backfill commit | closed (R = 3, D = 0) |
| v1.60 | P48 | `82e3ac8` + backfill commit | closed (R = 9, D = 0; run `41379998477` green on first run) |

## 4. Execution queue

Principle: layer 1 clears the record chain, layer 2 the strategic mainline, layer 3 the pre-release gates, layer 4 the exit gate; order is fixed inside a layer and layers are not skipped.

### Layer 1 (P48 removed, one item left)

| Order | Slice | Identity and content | Size | 7 budget | Depends on | Boundaries |
|---|---|---|---|---|---|---|
| 2 | RACE-LEG-EXPAND | extend the Race leg to console / gates (GAP-42) | M | Codex one final review | R-03 ruling | test coverage only, no production semantics change |

### Layer 2 (strategic mainline)

| Order | Slice | Identity and content | Size | 7 budget | Depends on | Boundaries |
|---|---|---|---|---|---|---|
| 3 | CONSOLE-MULTITASK-ASSEMBLY | console multi-task production assembly (breakpoint A closeout) | L | Codex one final review | layer 1 done + protocol first | the single-task limit must not be lifted before the protocol changes |
| 4 | ACCEPTANCE-SECURITY-COUNTEREXAMPLES | acceptance-chain security counter-examples (GAP-51 and others) | M | Codex one final review | item 3 | counter-examples and criteria only, no acceptance semantics change |

### Layers 3 and 4

| Order | Slice | Identity and content | Size | 7 budget | Depends on | Boundaries |
|---|---|---|---|---|---|---|
| 5 | CONTRACTS-S7-DURABILITY | contract-side durability semantics and Windows C6 landing | L | Codex one final review | R-02 ruling | change CONTRACTS before code |
| 6 | EXIT-SEED-BUILD-EVIDENCE | AT-13 seed-driven build evidence | M | Codex one final review | item 5 | evidence chain only, no CI structure change |
| 7 | AT23-RESTRICTED-CANDIDATE | AT-23 restricted real candidate and sign-off | M | Codex one final review | R-08 / R-09 rulings | no sign-off without a real-candidate authorisation |

Parked (not scheduled): `G4B-CORPUS-DECOUPLE` (trigger-based), five low-priority items (start only on trigger), remaining round-2 gaps (candidate pool).

## 5. Observation items

- Observing (13): `OB-1 / 7 / 10 / 33 / 87 / 99 / 100 / 104 / 105 / 107 / 111 / 113 / 135`.
- Under discussion (11): `OB-19 / 21 / 52 / 53 / 55 / 57 / 58 / 59 / 61 / 75 / 120`.
- To do (1): `OB-69`. Non-terminal total **25**.
- Mechanical readings: `section 12.4` max OB = **OB-136**; max deferral = 25th; ADR max = **28**; DR max = **13** (DR-13 fixed).

## 6. Pending user rulings (R-01 to R-11)

| ID | Item | Status |
|---|---|---|
| R-01 | P46 strategic review project | done |
| R-02 | `CONTRACTS-S7-DURABILITY` cadence and ownership | pending (blocks layer 3) |
| R-03 | ownership of the follow-up slice after P45 | pending (P48 completed under R-05 standalone; layer-1 item 2 still needs this ruling's scope confirmed) |
| R-04 | DR-13 "fixed plus observation" stance | pending |
| R-05 | whether DR13-FIX-D3 is a standalone slice | done (standalone and closed) |
| R-06 | OB-55 ownership and T027-B cadence | pending |
| R-07 | MVA, S1 removals and S2 scope narrowing | pending |
| R-08 | candidate, Windows premise P and probe authorisation | pending (blocks layer 4) |
| R-09 | AT-24 `unknown` acceptance boundary | pending (blocks P3 closeout) |
| R-10 | total budget and per-slice model caps | pending (blocks all scheduling) |
| R-11 | historical correction scope | done (option B plus per-item differentiation) |

## 7. Kickoff and step-0 rules

1. Verify anchors, then step 0 registers the ledger (R / D = 0 baseline, temp location `tmp/<sliceID>/`), lands `input-list.md` and `freedom-list.md`.
2. Every fix slice must ship regression tests; unexecuted tests or lint must be registered explicitly.
3. Commit only after the three scopes are green, `SELFTEST 697/697` and `go build/vet/test ./...` are all green.
4. Model binding: 1 and 5 `deep-reasoner` + `DeepSeek V4 Pro (deepseek)`; 2/3/8 `DeepSeek V4.1 Flash (deepseek)`; 6 `independent-reviewer` + `MAI-Code-1.1-Flash (copilot)`; 7 `independent-reviewer` + `GPT-5.3-Codex (copilot)`; batch assist `deep-reasoner` + `GLM-5.3-Flash (glm)`. Every `runSubagent` passes both `agentName` and `model`.
5. Rhythm discipline: no routine stops; 10 is the only regular stop point.
6. At least two batch-assist dispatches per slice, **strictly serial**, no parallel subagents (R1.1 / R5.1).
7. Main controller defaults to `DeepSeek V4.1 Flash`; switch to Pro for low-level concurrency or state machines, oversized design documents, large-scale CI log root-causing, or two consecutive drifts; register downgrades under OB-37.
8. Role carriers are not interchangeable: use `prfrail-implementer` for 2 and `prfrail-tester` for 3; the main controller only orchestrates and runs gates.
9. The kickoff message is an artefact: `00-briefing.md` enters MANIFEST and SHA256SUMS and is declared in the `artifacts` fence.
10. Evidence pack and fence stay in sync: every new payload gets a MANIFEST row, a hash recomputation and a fence row in the same round.
11. Reviewer archive discipline: outputs of 5, 6 and 7 are archived verbatim; the main controller must not rebuild them from bullet points.
12. Assessment slices stop at 10 after 9; never `commit` or `push` on your own.
13. Single output must not exceed 10000 characters; batch before the limit and never self-restart.
14. Compaction is slice end: at 100K context or a compaction prompt, end the slice.
15. The kickoff message must include a "document sync" section: if the slice changes the project stage or contract capability, stage 8 must update `DEV_PLAN.md`, `ARCHITECTURE.md` and other authoritative documents; a missing section is a structural defect in the kickoff message.

## 8. Cross-slice collaboration discipline (key points)

| Element | Content |
|---|---|
| Carrier table | 1, 5, 6, 7 are read-only; 2, 3, 8 may write; `prfrail-tester` must never run a review stage; batch assist must be `deep-reasoner` |
| Serial dispatch | `runSubagent` blocks serially; no parallel calls; pass `agentName` and `model` every time |
| Commit message language | always English, `<type>: <summary>`; use `git commit -F` when non-ASCII is present |
| Temp artefacts | always `tmp/<sliceID>/`; only 8c cleans the slice directory |
| Gate execution | intermediate rounds use `--all --skip-base`; all green before commit; discipline E: precise `git add` first, then `--scope=index` |
| Fence | sorted by path, TAB separated, CN and `_EN` byte-identical |
| Fact vs judgement | independently verifiable output errors do not void the slice; judgement errors must be redone |
| Discipline T · quantifier shape | ledger and reports must not put a quantifier right after an identifier: that quantifier is read as a line-count literal by `G6-4`; write "that entry" or "that record" instead |
| Discipline U · no verbatim trigger strings | reports must not quote the string that trips a criterion: repeating the literal makes the report red again; use a descriptive phrasing |

## 9. P48 gate pitfalls (for P49 reuse)

1. **`.go` copies under `tmp/` are compiled by `go build ./...`**: source copies in review input packages must use `*.go.txt`, otherwise the `base go build/vet/test` three legs go red.
2. **`G7-b` fence disposition closed set**: uncommitted temp artefacts use `local-only`; writing `local` reports an unknown disposition.
3. **A report carrying the `INDEPENDENT SCAN` marker needs a `sixstage` fence**: use `## Section A` through `## Section E` exactly once each and keep the verdict line as the last non-empty line in the block, disclosing the merge outside the fence.
4. **Writing `.md` requires an explicit BOM**: omitting the BOM argument when rewriting with node trips `G4a` immediately (hit twice in this slice).
5. **`G8-b` requires MANIFEST to list payloads individually**: a directory row is reported as "payload not listed".
6. **Readings must be backfilled last**: archiving changes the `changed=` count and the `scanned=` reading, so re-measure before writing the report.

## 10. Checklist meta rules

| Rule (opening quotation) | Content |
|---|---|
| "every two slices / ruling change / directive bump / queue change / user request / two or more anchor mismatches means regenerate" | regeneration timing |
| "pending rulings must be named one by one" | section 6 names each |
| "queue items must carry six elements" | identity / content / size / 7 budget / dependency / boundaries |
| "rule duty is not gate duty" | do not turn rule text into a mechanical criterion |
| "gates green before commit" | hard gate |
| "assessment slices stop at 10 after 9" | stop at 10 for authorisation when there is no implementation target |
| "compaction is slice end" | 100K or a compaction prompt means end the slice now |
| "one fact is described in one place" | body or pointer, not both |
| "governance slices do not write back to DEV_PLAN" | authoritative ledger = appendix C.0 to C.0x |

## 11. One-line opening for a new session

Here is my project state snapshot (v38). Verify the five anchors in section 1 first; if they match, continue with the queue in section 4; otherwise follow the repository state.

[paste the full v38 checklist]

[Attachments] `REMAINING_SLICES.md` (P48 completed), `DELIVERY_DIRECTIVE.md` (v1.60), `docs/validation/dr13-fix-d3{,_EN}.md` (P48 report and CI evidence).

Current slice: layer-1 item 2 = `RACE-LEG-EXPAND` (start after confirming the R-03 ruling scope).

## 12. Current state in one line

P48 `DR13-FIX-D3` fully closed (`82e3ac8`, run `41379998477` green on both legs on the first run, `internal/adapters` proven ok on both legs); layer 1 has one item left; next = `RACE-LEG-EXPAND` (confirm the R-03 scope first).
