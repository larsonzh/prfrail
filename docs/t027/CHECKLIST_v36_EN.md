# Task execution checklist (cross-session snapshot v36, strategic-review closeout)

Generated 2026-11-02; HEAD = origin/main = `5e4202e`; directive v1.60.
Supersedes v35 and any snapshot older than today.
Changes from v35: (1) anchors moved to the P46 closeout chain (`5e4202e`); (2) section 2 carries the P46 fact pack (breakpoint A partially closed, B not closed); (3) the queue is re-laid out after P46 round 3 as four layers: first three / strategic mainline / pre-release gates / exit gates; (4) three new lessons and five discipline candidates are added, including the **session-length and compaction discipline**; (5) the 24 non-terminal observation items and all P39–P45 candidates are inherited.

---

## 1 Current anchors (check these five first)

| Item | Value | Verification command |
|---|---|---|
| HEAD | `5e4202e` | `git rev-parse --short HEAD` |
| origin/main | `5e4202e` | `git rev-parse --short origin/main` |
| Directive version | **v1.60** | header of `docs/DELIVERY_DIRECTIVE.md` |
| SELFTEST | **697/697** | `node tools/gates/gate.js --selftest` |
| Three scopes | **all TOTAL_FAIL=0** | `gate.js --all --scope=index\|tree\|ci` |

**If two or more items disagree, the repository wins**: have the orchestrator report terminal readings first.

**[Three-scope green convention]** On a clean tree `index/tree` both report `scanned 0` (vacuous pass) and only `ci` scans N lines, so discipline E requires the actual `scanned N`. At v36 generation on `5e4202e` the `ci` leg measured `G6-1 scanned 0`; before the commit it measured `scanned 410`.

**[Fixed tmp/ assertions]** `tmp` top level = `.gitkeep,gates-cache`; one subdirectory under `gates-cache` (`g4b`); **files under `tmp/gates-cache/` outside `g4b/` = 0**; `Test-Path tmp/STRATEGIC-REVIEW-2026-11-02` is `False`, and the P40–P45 directories plus `tmp/DR13-FIX` are likewise `False`.

> **Two distinct readings**: files under `gates-cache` outside `g4b` = 0, and `tmp/` excluding `gates-cache` = 1 (`.gitkeep`); the two are not interchangeable.

---

## 2 P46 strategic-review fact pack

### 2.1 Conclusions and boundaries

| Item | Value |
|---|---|
| Slice | `STRATEGIC-REVIEW-2026-11-02` (evaluation class; stages 1–7 skipped; three review rounds plus compensation review) |
| Baseline | `f497ff3` (v35 landing); prior P45 closeout `3a7541d` |
| Commit chain | corrected draft `26ca5a5` (+56/−16, two files) plus backfill `5e4202e` (+10/−4, two files) |
| CI | `41344276390` both green (final CI evidence); `41345581562` (docs backfill) both green; intermediate draft `41339104941` both green (**process fact, not final evidence**) |
| **Breakpoint A** | **partially closed**: a zero-model engine chain does not mean a real multi-task AI production chain |
| **Breakpoint B** | **not closed** (governance record chain; seven items P3/P5/G4/G5/G7/M6/S7 remain not closed for insufficient evidence) |
| Round-1 verdicts | not closed 7, partially closed 5 (S1/S6/Z2/Z3/Z4), closed 0 |
| Round-2 gaps | 76 items (`GAP-01`–`GAP-76`), batch A 33 / batch B 43, **carrying no priority** |
| Round-3 conclusion | November mainline = console multi-task assembly extension, then `CONTRACTS-S7-DURABILITY`; finish the first three first |
| Compensation review | `COMPENSATION REVIEW: PASS`, covering **only format, structure and verdict-line consistency**, not judgement quality |

**Key boundary (read alongside the conclusions)**: the P45 console single-task limit (`runtime_agent.go:774`, `runtime_supervised.go:335`) counts against breakpoint A, so breakpoint A must not be upgraded to closed; round 1 reported not closed for insufficient evidence within an evidence search limited to the input pack, which does not mean not remediated.

### 2.2 Models and metering

| Role | agentName | model | Output |
|---|---|---|---|
| Round 1 / round 3 | `independent-reviewer` | `GPT-6.1 Sol (copilot)` | 3,018 / 4,027 characters |
| Round 2 | `deep-reasoner` | `DeepSeek V4 Pro (deepseek)` | 8,501 total (part 1 3,804 / part 2 4,695) |
| Compensation review | `independent-reviewer` | `MAI-Code-1.1-Flash (copilot)` | 650 characters |
| Stage 8a report writing | `prfrail-documenter` | `DeepSeek V4.1 Flash (deepseek)` | 2,614 characters |
| Batch helpers 1–2 | `deep-reasoner` | `GLM-5.3-Flash (glm)` | 10,190 / 9,627 characters (**both above their dispatch caps**) |
| Batch helper 3 | `deep-reasoner` | `GLM-5.3-Flash (glm)` | six-shape comparison, recomputed three times, last recomputation 590 characters; 11 dispatches in total |

**Metering anomaly**: GLM#1 emitted 10,190 characters, above both its 6,000-character dispatch cap and the 10,000 cap in section 7, yet returned intact; the platform demonstrably returns at least 10,190 characters, the basis for recalibrating the batch-helper cap (see section 8 item 5).

### 2.3 Deviations D-01 and D-02

| Item | Content |
|---|---|
| D-01 | The section-7 six-shape readings are self-referential (editing sections 3/6/9 moves them); resolved by editing the body first and backfilling section 7 last, then `SIX-SHAPE: EQUAL` |
| D-02 | After stage 9 the orchestrator ran `commit/push` without waiting for the stage-10 stop point or a same-turn authorisation (violating the section 10.1 authorisation and section 10.4 stop-point disciplines). Stop-loss: no force push, no revert, no amend; the intermediate-draft CI counts as a process fact only |
| R count | +1 (fallback intervention: the orchestrator corrected one discipline deviation) |

Status: [done]

---

## 3 Completed slices (recent baseline)

| Version | Slice | Main commit + backfill | Conclusion |
|---|---|---|---|
| v1.60 | P43 | `b6e0f34` + `609bb19` + `7004282` + `8a46b01` | closed (R=0, D=0; D7/D7b/D8) |
| v1.60 | P44 | `9061b67` + `16abe2b` | closed (R=0, D=0; D1–D6, DR-13 added) |
| v1.60 | DR13-FIX | `b123db1` + `eb36992` | closed (P45 prerequisite fix; DR-13 repaired) |
| v1.60 | P45 | `4273d6e` + `3a7541d` | closed (R=0, D=0; D10/D11) |
| v1.60 | v35 snapshot | `f497ff3` | baseline (P45 closed, before P46) |
| **v1.60** | **P46 strategic review** | **`26ca5a5` + `5e4202e`** | **closed (R=1, D=2; D-02 discipline deviation, finalised after five fixes)** |

---

## 4 Execution queue (four layers, per P46 round 3)

**Principle: clear the governance record chain (layer 1), then the strategic mainline (layer 2), then the pre-release gates (layer 3), then the exit gates (layer 4). Order inside a layer is fixed; layers must not be skipped.**

### Layer 1, the first three (hard sequence from P46 round 3)

| # | Slice | Carries | Size | Depends on |
|---|---|---|---|---|
| 1 | **P47 `GOVERNANCE-BASELINE-VERIFY`** (trustworthy baseline verification and historical re-registration for breakpoint B) | breakpoint B not closed (seven items P3/P5/G4/G5/G7/M6/S7) | M | **ruling R-11** |
| 2 | `DR13-FIX-D3` (cross-process hardening: remove before unlock) | DR13-FIX leftover D3 (GAP-71) | S | ruling R-05 |
| 3 | `RACE-LEG-EXPAND` (extend the Race leg to console/gates) | breakpoint B item G5 and GAP-42 | M | ruling R-03 |

### Layer 2, the strategic mainline (November)

| # | Slice | Carries | Size | Depends on |
|---|---|---|---|---|
| 4 | `CONSOLE-MULTITASK-ASSEMBLY` (console multi-task production assembly) | mandatory position 4; closes breakpoint A (GAP-14/56/57/62/63) | L | layer 1 done + protocol first |
| 5 | `ACCEPTANCE-SECURITY-COUNTEREXAMPLES` (acceptance-chain safety counterexamples) | mandatory position 10 (GAP-51/52/55/60/64/70) | M | item 4 |

### Layer 3, pre-release gates

| # | Slice | Carries | Size | Depends on |
|---|---|---|---|---|
| 6 | `CONTRACTS-S7-DURABILITY` (contract durability semantics and Windows C6) | mandatory positions 8 and 9 (GAP-26/34/58) | L | ruling R-02 |

### Layer 4, exit gates

| # | Slice | Carries | Size | Depends on |
|---|---|---|---|---|
| 7 | `EXIT-SEED-BUILD-EVIDENCE` (AT-13 seed-driven build evidence) | mandatory position 11 (breakpoints A/B item P5) | M | layer 3 done |
| 8 | `AT23-RESTRICTED-CANDIDATE` (AT-23 restricted real candidate and sign-off) | mandatory position 12 (GAP-14/17/65) | M | rulings R-08/R-09 |

### Suspended items (not scheduled in this version)

| Item | Status | Note |
|---|---|---|
| `G4B-CORPUS-DECOUPLE` (OB-120) | suspended, trigger-based | trigger = the G4b corpus breaks again or the cache hit rate drops |
| Low priority (`OB-50` / `spawnGit` / `ci frozen cost` / `%TEMP% leftovers` / `GATES-EXEMPT`) | suspended long-term | no trigger, no start |
| The 76 round-2 gaps minus the ones already carried | candidate pool | absorbed item by item inside layers 1–4, not scheduled separately |

### P47 six elements (layer 1 item 1, ready to start)

| Element | Content |
|---|---|
| ID | `GOVERNANCE-BASELINE-VERIFY` (P47) |
| Content | For each of the seven not-closed-for-insufficient-evidence items of breakpoint B, verify a trustworthy baseline: confirm the evidence search scope, re-register missing records, and conclude either upgraded after evidence or kept open with the gap stated |
| Size | M (governance/evaluation; expected zero or minimal production-code change, declared explicitly if any) |
| Stage-7 budget | default `independent-reviewer` + `GPT-5.3-Codex (copilot)`, one final review; for pure docs that touch no gate criteria, `deep-reasoner` + `DeepSeek V4 Pro (deepseek)` is allowed, registered under OB-37 |
| Dependencies | **ruling R-11 (historical-correction scope: P32 counts / P33 independent registration / P34 exception)**; ruling R-04 (DR-13 convention) as reference |
| Boundary | verification and registration only; no directive or contract body edits; no lifting of the console single-task limit; no reopening of DR-13; repaired-plus-observed must not be read as permanently closed |

---

## 5 Observation items

### Observing (12)

`OB-1 / 7 / 10 / 33 / 87 / 99 / 100 / 104 / 105 / 107 / 111 / 113`

### Under discussion (11)

`OB-19 / 21 / 52 / 53 / 55 / 57 / 58 / 59 / 61 / 75 / 120`

### To do (1)

`OB-69`

**24 non-terminal items in total** (12 observing + 11 under discussion + 1 to do); section 12.4 is unchanged in this version, so the counts match v34/v35.

### Ruled (28) and additions in the P46 window

OB-117 to OB-132 plus the existing additions keep the v35 convention; the P46-window rulings (three-round finalisation, the five-fix instruction, D-02 handling, the two-step commit/push authorisation) are **slice-level rulings** and are not registered as OB.

### Mechanical readings (verifiable convention)

- Section 12.4 maximum OB = **OB-132**; maximum deferral = the 25th; ADR max = **28**; DR max = **13** (repaired).
- Section 12.4 entries: 132; `disposed 75 / ruled 28 / registered 5 / observing 12 / under discussion 11 / to do 1`.

---

## 6 Awaiting the user ruling (P46 numbering R-01 to R-11)

| ID | Item | Status | Blocks |
|---|---|---|---|
| R-01 | Whether to launch the P46 strategic review | **done** | — |
| R-02 | Cadence and placement of `CONTRACTS-S7-DURABILITY` | awaiting | layer 3 (item 6) |
| R-03 | Ownership of the P45 follow-ups (console multi-task assembly / Race-leg expansion R-4 / `registry.len()` assertion) | awaiting | layer 1 item 3 and layer 2 item 4 |
| R-04 | Convention for DR-13 repaired-plus-observed and for reopening on recurrence | awaiting | CI stability convention |
| R-05 | Whether DR13-FIX D3 becomes its own slice | awaiting | layer 1 item 2 |
| R-06 | Ownership of OB-55 and the T027-B cadence | awaiting | later implementation |
| R-07 | Approval of MVA, the S1 removals and the S2 commitment narrowing | awaiting | later slices |
| R-08 | Candidates, Windows prerequisite P and resource/probe authorisation | awaiting | layer 4 item 8 |
| R-09 | Acceptance boundary for AT-24 `unknown` and keeping unknown side effects paused | awaiting | layer 4 item 8 |
| R-10 | Total budget and per-slice model caps | awaiting | all scheduling |
| R-11 | **Historical-correction scope (P32 counts / P33 independent registration / P34 exception)** | awaiting | **layer 1 item 1 (P47)** |

---

## 7 Start-up and step-0 general rules

1. Check the anchors, then run step 0: register the slice ledger (`R/D = 0` baseline, temp area `tmp/<slice-id>/`) and write `input-list.md` and `freedom-list.md`.
2. Every fix slice must ship **regression tests** and register any test or lint not run.
3. Commit only when stage 9 shows three green scopes, `SELFTEST 697/697` and green `go build/vet/test ./...`.
4. Model bindings (hard): stages 1 and 5 `deep-reasoner` + `DeepSeek V4 Pro (deepseek)`; stages 2, 3, 8 `DeepSeek V4.1 Flash (deepseek)`; stage 6 `independent-reviewer` + `MAI-Code-1.1-Flash (copilot)`; stage 7 `independent-reviewer` + `GPT-5.3-Codex (copilot)`; strategic-review rounds 1 and 3 `independent-reviewer` + `GPT-6.1 Sol (copilot)`; batch helper `deep-reasoner` + `GLM-5.3-Flash (glm)`. **Every `runSubagent` must pass both `agentName` and `model`.**
5. **Cadence**: no routine stops; green stage 4 dispatches stage 5 at once; a Medium+ from stages 5/6/7 means fix and re-run that stage in place; **stage 10 is the only routine stop point**.
6. Cross-version and cross-platform triage: on a CI/local split check the toolchain version and build-tag coverage first; after cross builds verify with `go env GOOS GOARCH`; before deleting a symbol believed unused locally, run `go vet` under the opposite `GOOS` at least once.
7. Batch helper: at least **2** dispatches per slice with `deep-reasoner` + `GLM-5.3-Flash (glm)` (policy and output cap in section 8).
8. Orchestrator model policy: `DeepSeek V4.1 Flash` by default; switch to Pro for deep concurrency, state-machine or handle work, for design docs above 10,000 characters, for deep root-cause work over large CI logs, or after two consecutive behaviour drifts; **any downgrade is registered separately under meta-rule OB-37**.
9. **Carriers are not substitutable**: stage 2 must be `prfrail-implementer` and stage 3 must be `prfrail-tester`; the orchestrator only orchestrates, lands files, runs gates and writes the receipt, and **must not write production code or tests**.
10. **The start message is an artifact**: step 0 must land the user start message verbatim as `docs/validation/evidence/<slice-id>/00-briefing.md`, add it to `MANIFEST.md` and `SHA256SUMS.txt`, and declare it in the `artifacts` fence; without it stage 1 must not start.
11. **Linearisation point**: for a persistent state machine, a successful persist must be the linearisation point and in-memory CAS only a same-process fast-reject cache; state that a successful reservation implies the marker is on disk and back it with a positive red test.
12. **Evidence pack and fence sync**: a new evidence payload needs a `MANIFEST.md` entry, a recomputed `SHA256SUMS.txt` and a new `artifacts` fence line in the same turn; missing any one trips `G8-b` or `G5(d)`.
13. **Reviewer archiving**: output from stages 5, 6 and 7 is archived verbatim; the orchestrator must not rebuild it from bullet points; a summary, if needed, is a separate file declared outside the fence.
14. **Assembly-validation slices**: assertions must rest on falsifiable mechanical criteria; constant identities must be declared non-sentinels; the report must state whether production code is unchanged; a capability unsupported in one layer must not be read as supported.

---

## 8 Cross-slice collaboration disciplines (mandatory at start-up)

| Element | Content |
|---|---|
| Carriers | stages 1, 5, 6, 7 use `deep-reasoner` / `independent-reviewer` (read-only); stages 2, 3, 8 use `prfrail-implementer` / `prfrail-tester` / `prfrail-documenter` (writable); `prfrail-tester` never runs review stages; the batch helper must be `deep-reasoner`, never `quick-researcher` |
| Mechanical dispatch table | each `agentName` in the stage-10 table is read verbatim from `transcript.arguments.agentName`, with the script stdout attached |
| R2.1 agentName | every `runSubagent` call carries `arguments` with both `agentName` and `model` |
| Non-ASCII commit messages | use `git commit -F <utf8-file>` whenever the message contains non-ASCII characters |
| Independently verifiable or not | independently verifiable output (enumeration, fact comparison, reference lists) survives a model error; non-verifiable output (judgements, conclusions, review opinions) must be redone |
| Stage 6 before stage 7 | stage 6 completes first; if stage 6 is re-dispatched with document changes, stage 7 must be re-evaluated |
| Temp placement and cleanup | always `tmp/<slice-id>/`; only stage 8c cleans it, deleting just this slice directory and never its parent |
| Long-command redirection | avoid PowerShell `>`; use `\| Tee-Object -FilePath <path>` or `cmd /c "..."` |
| Gate execution | use `--all --skip-base` mid-slice; all green before commit |
| Discipline E, index scope order | precise `git add` first, then `--scope=index`; the receipt reports `G6-1 scanned N` with N above 0 |
| Discipline F, specification first | when a start-up example conflicts with the specification the specification wins, and the conflict is registered |
| Discipline G, the fence is the landed snapshot | the `artifacts` fence is path-sorted, TAB-separated and byte-identical across CN/EN; new files need a new fence line |
| Discipline H, opposite-GOOS check | before deleting or changing a constant or function run `go vet` under the opposite `GOOS` at least once |
| Disciplines I to N | closed verdict set / review semantics / persistence first / three-way pack consistency / verbatim archiving / mechanical criteria and non-sentinel declarations (see section 7 items 11–14) |
| **Discipline O, output-length defence** | a single output stays within 10,000 characters; split per the section-7 plan before overflow; self-restart is forbidden; read on demand rather than loading everything; stop and report as soon as a stop-loss condition triggers |
| **Discipline P, orchestrator model switch** | see section 7 item 8: deep concurrency, state machines or handles, oversized design docs, root-cause work over large CI logs, or two consecutive Flash drifts trigger a switch to Pro; downgrades registered under OB-37 |
| **Discipline Q, batch-helper output cap** | the dispatch states a character cap; the receipt reports the actual count; an over-cap but intact result is registered as a metering fact outside the section-7 stop-loss scope; the measured platform ceiling of at least 10,190 characters calibrates the cap |
| **Discipline R, bulk settlement of lagging state** | a ledger-versus-fact mismatch (the GAP-12 / GAP-19 / GAP-24 class) is cleared in one settlement slice with code and commits as the fact source; line-by-line patching is not acceptable |
| **Discipline S, session length and compaction** | on reaching 100K context or on a compaction notice, end the current slice at once; one session carries one L-grade slice; compaction evicts early hard stop-point and authorisation rules and invites drift (D-02 was that class) |

**GLM batch-helper policy**: `GLM-5.3-Flash (glm)` is plentiful and the floor of 2 dispatches is only a floor; outsource mechanical checks, long-document excerpting, bilingual mirror checks, dispatch drafts and log extraction wherever possible; no complex judgement or ranking; every dispatch is metered in the stage-10 receipt (output cap in discipline Q).

**P46 lessons (from v36)**: (1) **an evaluation slice must stop at stage 10 after stage 9 goes green** and wait for authorisation; D-02 arose from treating it as an implementation pipeline; (2) **the section-7 readings are self-referential**, so edit the body first, backfill section 7 last and recompute until stable; (3) **gate criteria block wording**: `G2` catches placeholder words and `G6-4` catches count-style literals (placeholder words must not remain), and the remedy is rewording, never widening the whitelist.

---

## 9 Checklist meta-rules

**Quote the opening words rather than a clause number.**

| Rule (opening words) | Content |
|---|---|
| Every two slices, an awaiting-ruling change, a directive upgrade, a queue change, a user request, or two or more anchor mismatches trigger regeneration | when to regenerate |
| Awaiting-ruling items are listed one by one | section 6 lists every item |
| Every queue item states six elements | id / content / size / stage-7 budget / dependencies / boundary |
| Rule obligations are not gate obligations | rule text is not a mechanical criterion |
| No wall-clock assertions | observation clauses trigger on events, never on a passage of months |
| Gates green before commit | hard gate |
| Dispatch-table agentName is read mechanically | verbatim from `transcript.arguments.agentName` |
| Governance slices do not write back DEV_PLAN | their authoritative ledger is appendix C.0 to C.0x |
| A slice name never changes once in the ledger | names are stable |
| Policy wins over an authorised quantity | never grounds to declare overreach or refuse execution |
| One description per fact | either the body or a pointer |
| Orchestrator text may be completed; reviewer text stays verbatim | the two standards differ |
| Factual errors are independent of design decisions | ledger and report contradictions are factual, not design questions |
| The stage-10 receipt carries a fact pack for the next checklist | fields: HEAD, CI run, new ledger rows, new observation items, new discipline candidates, maximum OB and deferral, maximum ADR and DR, tmp assertions, differences from the previous version, open items |
| Routine small fixes add no slice number but lose no discipline | same stage-9 gates, CI and fence line, with no incidental semantic change |
| The start message is an artifact | `00-briefing.md` is landed, added to MANIFEST and SHA256SUMS, and declared in the fence |
| A carrier deviation is the highest independence downgrade | it needs a D entry with an independent impact assessment |
| Mechanical criteria beat testimonial statements | zero-model and zero-side-effect claims must be falsifiable; constant identities are marked non-sentinel |
| Concurrency assertions cover survival and release | besides the losing side, assert the winner stays alive and the slot can be re-taken after release |
| Layer boundaries are not widened | unsupported capabilities are declared explicitly |
| An evaluation slice stops at stage 10 after stage 9 | **new in v36**: with no implementation object, it stops at stage 10 for a same-turn authorisation and never commits or pushes on its own |
| A single output stays within 10000 characters | **new in v36**: split before overflow, no self-restart, read on demand |
| Compaction ends the slice | **new in v36**: at 100K context or on a compaction notice, end the slice; one session carries one L-grade slice |

---

## 10 One-paragraph opener for a new session

```
Here is my project snapshot (v36). Check the five anchors in section 1 first; if they match, continue with the four-layer queue in section 4; otherwise the repository wins.

[paste this checklist v36 in full]

[supporting documents]
- REMAINING_SLICES.md (ledger; P46 ✅; DR13-FIX ✅; DR max = 13 with DR-13 repaired)
- DELIVERY_DIRECTIVE.md (v1.60)
- docs/validation/strategic-review-2026-11-02{,_EN}.md (the P46 report with D-01/D-02 and the three rounds)
- docs/validation/p45-e2e-zero-model{,_EN}.md and dr13-fix{,_EN}.md (prior slice reports)
- docs/CONTRACTS.md sections 4.1 and 4.2 (sealing barrier and K, readiness grading and probe permits)

Current slice: layer 1 item 1, P47 `GOVERNANCE-BASELINE-VERIFY` (depends on ruling R-11).
Flash status: [to start / running / awaiting receipt]
```

---

## 11 Current state in one sentence

P35 to P46 are all closed (`5e4202e`, directive v1.60); the three P46 review rounds confirmed **breakpoint A partially closed and breakpoint B not closed** and produced the four-layer November sequence; **the next slice is layer 1 item 1, P47 (trustworthy baseline verification for breakpoint B), which depends on ruling R-11**; D-02 is registered as a discipline deviation and produced the rule that an evaluation slice must stop at stage 10 after stage 9.
