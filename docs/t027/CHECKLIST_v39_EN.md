# Task Execution Checklist (cross-session snapshot v39 · P49 stop-point edition)

Generated: 2026-11-02, HEAD = origin/main = `75990a8` (the P49 main commit), directive v1.60. P49 is committed and green on the first attempt on both legs (run `41445657595`).
Supersession: completely supersedes v38 and any earlier snapshot. Backfill (docs-only) commits only advance HEAD and do not change the anchor semantics of this section.

Revisions in this version (vs v38): (1) P49 `RACE-LEG-EXPAND` closed into the baseline (`75990a8`, run `41445657595` green on the first attempt on both legs); (2) layer 1 cleared, the first layer-2 item (`CONSOLE-MULTITASK-ASSEMBLY`) becomes the next dequeue; (3) a new "environment pre-check and reuse" iron rule added to section 7; (4) **sections 8 and 9 restored to full v36 detail** (the complete discipline A-S table plus all twenty-three checklist meta-rules, with the new disciplines T and U folded in), keeping the seven P49 gate pitfalls as a subsection of section 9 (for reuse in P50); (5) observation items and pending rulings updated (R-03 completed).

## 1. Current anchors (verify these 5 numbers first)

| Item | Value | Check command |
|---|---|---|
| HEAD | `75990a8` | `git rev-parse --short HEAD` |
| origin/main | `75990a8` | `git rev-parse --short origin/main` |
| Directive version | **v1.60** | head of `docs/DELIVERY_DIRECTIVE.md` |
| SELFTEST | **697/697** | `node tools/gates/gate.js --selftest` |
| Three scopes | **all TOTAL_FAIL=0** | `gate.js --all --scope=index`, `tree`, `ci` |

If 2 or more items mismatch ⇒ treat the repository reality as authoritative and report the final readings first.

`tmp` fixed assertion: top level = `.gitkeep` + `gates-cache`; only `g4b` under `gates-cache`; `tmp/RACE-LEG-EXPAND` was removed at 8c.

## 2. P49 fact pack (`RACE-LEG-EXPAND`)

| Field | Value |
|---|---|
| Slice / size | `RACE-LEG-EXPAND` (infrastructure `[SLICE]`) / M |
| Goal | Expand the CI Race leg (Ubuntu-only, `CGO_ENABLED=1 go test -race -count=1`) from `internal/adapters`, `internal/chain` to six packages (plus `console`, `gates`, `guard`, `tickets`) |
| Changes | One line in `.github/workflows/ci.yml` Race `run:` plus one line in `internal/release/workflow_test.go` (frozen-contract digest sync for `test/Race`) |
| Red/green | 3 static red baseline (pre-change ci.yml, exit 1) → green flip after 2 (exit 0 GREEN) → VM real `-race` all four packages green (console 20.410s / gates 3.384s / guard 3.043s / tickets 7.420s, rc=0) |
| Review chain | 5 `PRE-REVIEW: PASS WITH FIXES` (F1 Medium closed); 6 r1 `PASS WITH FINDINGS` (format deviation) → r2 `INDEPENDENT SCAN: PASS` (A-E compliant re-issue); 7 r1 `PASS WITH FIXES` → r2 `FINAL REVIEW: PASS` |
| Commit / CI | `75990a8` pushed to origin (never gitee); run `41445657595` green on the first attempt on both legs (Race leg six packages all ok: console 8.835s / gates 3.265s / guard 2.958s / tickets 5.759s); the first-run obligation is closed |
| Metering | R = 1, D = 2 (2 r2 contract sync; the 1 retrieval miss plus the main-controller environment pre-check blind spot) |
| Doc sync | `ARCHITECTURE{,_EN}` section 1 tail each +2 lines on test facilities; `DEV_PLAN{,_EN}` checked, no outdated description, unchanged |
| Open items | CI first-run observation; flake observation for the four new packages entering the CI Race leg for the first time (GLM-2 B1/B2/B3) |

## 3. Completed slices (recent baseline)

| Version | Slice | Main commit + backfill | Conclusion |
|---|---|---|---|
| v1.60 | P45 | `4273d6e` + `3a7541d` | Closed |
| v1.60 | P46 | `26ca5a5` + `5e4202e` | Closed (R = 1, D = 2) |
| v1.60 | P47 | `bbd1745` + backfill commit | Closed (R = 3, D = 0) |
| v1.60 | P48 | `82e3ac8` + `81dd49f` | Closed (R = 9, D = 0; run `41379998477` green on the first attempt) |
| v1.60 | P49 | `75990a8` + backfill commit | Closed (R = 1, D = 2; run `41445657595` green on the first attempt) |

## 4. Execution queue

General principle: layer 1 clears the record chain → layer 2 strategic main line → layer 3 pre-release gates → layer 4 exit gates; order is fixed within a layer, no skipping between layers.

### Layer 1 (cleared after P49 dequeued)

No remaining items. R-03 is completed (the standalone `RACE-LEG-EXPAND` slice was opened and closed).

### Layer 2 (strategic main line)

| Order | Slice | Identity and contents | Size | 7 budget | Depends on | Boundaries |
|---|---|---|---|---|---|---|
| 3 | CONSOLE-MULTITASK-ASSEMBLY | console multi-task production assembly extension (breakpoint A wrap-up) | L | Codex 1 final review | protocol first | no lifting of the existing single-task restriction before the protocol changes |
| 4 | ACCEPTANCE-SECURITY-COUNTEREXAMPLES | acceptance-chain security counterexamples (GAP-51 etc.) | M | Codex 1 final review | item 3 | counterexamples and criteria only, no acceptance-semantics change |

### Layers 3 / 4

| Order | Slice | Identity and contents | Size | 7 budget | Depends on | Boundaries |
|---|---|---|---|---|---|---|
| 5 | CONTRACTS-S7-DURABILITY | contract-side durability semantics and Windows C6 landing | L | Codex 1 final review | R-02 ruling | CONTRACTS first, then code |
| 6 | EXIT-SEED-BUILD-EVIDENCE | AT-13 seed-driven build evidence | M | Codex 1 final review | item 5 | evidence chain only, no CI structure change |
| 7 | AT23-RESTRICTED-CANDIDATE | AT-23 restricted real candidate and sign-off | M | Codex 1 final review | R-08 / R-09 rulings | no sign-off without a real-candidate authorisation |

Suspended (unscheduled): `G4B-CORPUS-DECOUPLE` (trigger-based), five low-priority items (no trigger, no start), the remaining round-2 gaps (candidate pool).

## 5. Observation items

- Observing (13): `OB-1 / 7 / 10 / 33 / 87 / 99 / 100 / 104 / 105 / 107 / 111 / 113 / 135`.
- To discuss (11): `OB-19 / 21 / 52 / 53 / 55 / 57 / 58 / 59 / 61 / 75 / 120`.
- To do (1): `OB-69`. Non-terminal total **25**.
- Mechanical readings: max OB in section 12.4 = **OB-136**; max carry-forward = the 25th; ADR max = **28**; DR max = **13** (DR-13 fixed).

## 6. Pending user rulings (R-01 to R-11)

| Id | Item | Status |
|---|---|---|
| R-01 | P46 strategic retrospective project | Completed |
| R-02 | `CONTRACTS-S7-DURABILITY` cadence and ownership | Pending (blocks layer 3) |
| R-03 | Ownership of the post-P45 slice | Completed (`RACE-LEG-EXPAND` opened and closed standalone; layer 1 cleared) |
| R-04 | DR-13 "fixed plus observe" wording | Pending |
| R-05 | Whether DR13-FIX-D3 is a standalone slice | Completed (opened and closed standalone) |
| R-06 | OB-55 ownership and T027-B cadence | Pending |
| R-07 | MVA, S1 removals and S2 promise narrowing | Pending |
| R-08 | Candidates, Windows precondition P and probe authorisation | Pending (blocks layer 4) |
| R-09 | AT-24 `unknown` acceptance boundary | Pending (blocks P3 wrap-up) |
| R-10 | Total budget and per-slice model caps | Pending (blocks all scheduling) |
| R-11 | Historical correction scope | Completed (tier B plus per-item differentiation) |

## 7. Kickoff and step-0 general rules

1. Verify the anchors, then run step 0: register the slice ledger (`R/D = 0` baseline, temp area `tmp/<slice-id>/`) and write `input-list.md` and `freedom-list.md`.
2. Every fix slice must ship **regression tests** and register any test or lint not run.
3. Commit only when stage 9 shows three green scopes, `SELFTEST 697/697` and green `go build/vet/test ./...`.
4. Model bindings (hard): stages 1 and 5 `deep-reasoner` + `DeepSeek V4 Pro (deepseek)`; stages 2, 3, 8 `DeepSeek V4.1 Flash (deepseek)`; stage 6 `independent-reviewer` + `MAI-Code-1.1-Flash (copilot)`; stage 7 `independent-reviewer` + `GPT-5.3-Codex (copilot)`; strategic-review rounds 1 and 3 `independent-reviewer` + `GPT-6.1 Sol (copilot)`; batch helper `deep-reasoner` + `GLM-5.3-Flash (glm)`. **Every `runSubagent` must pass both `agentName` and `model`.**
5. **Cadence**: no routine stops; green stage 4 dispatches stage 5 at once; a Medium+ from stages 5/6/7 means fix and re-run that stage in place; **stage 10 is the only routine stop point**.
6. Cross-version and cross-platform triage: on a CI/local split check the toolchain version and build-tag coverage first; after cross builds verify with `go env GOOS GOARCH`; before deleting a symbol believed unused locally, run `go vet` under the opposite `GOOS` at least once.
7. Batch helper: at least **2** dispatches per slice with `deep-reasoner` + `GLM-5.3-Flash (glm)` (policy and output cap in section 8); **strictly serial**, no parallel subagents (R1.1 / R5.1).
8. Orchestrator model policy: `DeepSeek V4.1 Flash` by default; switch to Pro for deep concurrency, state-machine or handle work, for design docs above 10,000 characters, for deep root-cause work over large CI logs, or after two consecutive behaviour drifts; **any downgrade is registered separately under meta-rule OB-37**.
9. **Carriers are not substitutable**: stage 2 must be `prfrail-implementer` and stage 3 must be `prfrail-tester`; the orchestrator only orchestrates, lands files, runs gates and writes the receipt, and **must not write production code or tests**.
10. **The start message is an artifact**: step 0 must land the user start message verbatim as `docs/validation/evidence/<slice-id>/00-briefing.md`, add it to `MANIFEST.md` and `SHA256SUMS.txt`, and declare it in the `artifacts` fence; without it stage 1 must not start.
11. **Linearisation point**: for a persistent state machine, a successful persist must be the linearisation point and in-memory CAS only a same-process fast-reject cache; state that a successful reservation implies the marker is on disk and back it with a positive red test.
12. **Evidence pack and fence sync**: a new evidence payload needs a `MANIFEST.md` entry, a recomputed `SHA256SUMS.txt` and a new `artifacts` fence line in the same turn; missing any one trips `G8-b` or `G5(d)`.
13. **Reviewer archiving**: output from stages 5, 6 and 7 is archived verbatim; the orchestrator must not rebuild it from bullet points; a summary, if needed, is a separate file declared outside the fence.
14. **Assembly-validation slices**: assertions must rest on falsifiable mechanical criteria; constant identities must be declared non-sentinels; the report must state whether production code is unchanged; a capability unsupported in one layer must not be read as supported.
15. **An evaluation slice stops at stage 10 after stage 9**: with no implementation object it stops at stage 10 for a same-turn authorisation and never commits or pushes on its own.
16. **A single output stays within 10,000 characters**: split before overflow, no self-restart, read on demand.
17. **Compaction ends the slice**: at 100K context or on a compaction notice, end the slice; one session carries one L-grade slice.
18. **The kickoff message must contain a "documentation sync update" section**: if the slice changes the project stage or contract capabilities, stage 8 must sync `DEV_PLAN.md`, `ARCHITECTURE.md` and other authoritative documents; a missing section counts as a structural defect of the kickoff message.
19. **Environment pre-check and reuse iron rule (new in v39)**: if a later slice needs a local Linux environment, remote SSH, a network proxy or any other host-dependent external environment, the main controller must declare it explicitly in step 0. If an environment gap is found (e.g. a command is missing), **the existing environment must be explored first (check PATH, configuration, find searches)**; only after confirming the tool is really absent may a minimal installation proceed, and blind re-installation is strictly forbidden. **For remote SSH environments, explicitly `source ~/.bashrc` or use `bash -lc` to load the login environment (a non-interactive shell does not load the user PATH, and `.bashrc` often has a non-interactive early-return guard), and user-level directories (such as `~/.local/bin`, `~/.local/opt`) must be explored**; only after confirming no existing toolchain may installation proceed, and the installation must pin a CI-aligned version with checksum verification.

## 8. Cross-slice collaboration disciplines (mandatory at start-up)

| Element | Content |
|---|---|
| Carriers | stages 1, 5, 6, 7 use `deep-reasoner` / `independent-reviewer` (read-only); stages 2, 3, 8 use `prfrail-implementer` / `prfrail-tester` / `prfrail-documenter` (writable); `prfrail-tester` never runs review stages; the batch helper must be `deep-reasoner`, never `quick-researcher` |
| Mechanical dispatch table | each `agentName` in the stage-10 table is read verbatim from `transcript.arguments.agentName`, with the script stdout attached |
| Serial dispatch | `runSubagent` is serial and blocking, no parallel subagents (R1.1 / R5.1); pass `agentName` and `model` every time |
| R2.1 agentName | every `runSubagent` call carries `arguments` with both `agentName` and `model` |
| Commit messages | always English, format `<type>: <summary>`; use `git commit -F <utf8-file>` whenever the message contains non-ASCII characters |
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
| **Discipline T, quantity form** | the ledger and reports must not write the "identifier immediately followed by a quantity" form: `G6-4` reads such a quantity as a line-number literal; write "the entry" or "the record" instead |
| **Discipline U, no verbatim trigger strings** | reports must not quote a trigger string verbatim: reproducing the same literal while describing a red cause re-triggers the red; use descriptive wording instead |

**GLM batch-helper policy**: `GLM-5.3-Flash (glm)` is plentiful and the floor of 2 dispatches is only a floor; outsource mechanical checks, long-document excerpting, bilingual mirror checks, dispatch drafts and log extraction wherever possible; no complex judgement or ranking; every dispatch is metered in the stage-10 receipt (output cap in discipline Q).

**P46 lessons (from v36)**: (1) **an evaluation slice must stop at stage 10 after stage 9 goes green** and wait for authorisation; D-02 arose from treating it as an implementation pipeline; (2) **the section-7 readings are self-referential**, so edit the body first, backfill section 7 last and recompute until stable; (3) **gate criteria block wording**: `G2` catches placeholder words and `G6-4` catches count-style literals (placeholder words must not remain), and the remedy is rewording, never widening the whitelist.

## 9. Checklist meta rules

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
| An evaluation slice stops at stage 10 after stage 9 | with no implementation object, it stops at stage 10 for a same-turn authorisation and never commits or pushes on its own |
| A single output stays within 10000 characters | split before overflow, no self-restart, read on demand |
| Compaction ends the slice | at 100K context or on a compaction notice, end the slice; one session carries one L-grade slice |

### Gate pitfalls (six from P48 + seven from P49, for reuse)

### From P48 (kept)

1. **`.go` copies under `tmp/` get compiled by `go build ./...`**: source copies for review input must use `*.go.txt`, otherwise the base go build/vet/test legs go red.
2. **`G7-b` fence disposition closed set**: use `local-only` for uncommitted temp artefacts; `local` reports an unknown disposition.
3. **A report mentioning `INDEPENDENT SCAN` must carry a `sixstage` fence**: inside the fence use `## Section A` through `## Section E` once each, keep the conclusion line as the last non-empty line of the block, and disclose merge differences outside the fence.
4. **Writing `.md` must add the BOM explicitly**: omitting the BOM parameter when writing back via node immediately triggers `G4a` (hit twice in P48).
5. **`G8-b` requires MANIFEST to list payloads item by item**: a directory line is judged "payload not listed".
6. **Readings must be backfilled last**: closing the archive changes the `changed=` count and the `scanned=` readings; re-measure before finalising the report.

### From P49 (new)

7. **Changing `.github/workflows/*.yml` must sync the `internal/release` workflow freeze contract**: `expectedScripts` in `workflow_test.go` pins the SHA-256 of every script step body; stage-1 impact searches must cover this contract (the P49 stage-1 miss made the contract family red after r1, closed by the r2 sync).
8. **Remote SSH pre-check**: a non-interactive session does not load the user PATH (`source ~/.bashrc` is often swallowed by an early-return guard); use `bash -lc` and explore user-level directories such as `~/.local/bin` (P49 missed the VM's existing go1.27.1 this way, registered as a deviation).
9. **PowerShell → ssh nested quoting**: nested quotes and `$(...)` get parsed by PowerShell first; for any cross-shell complex command, write a remote shell script, upload it and execute it.
10. **Stage 6 output must be A-E shaped**: `Section A` (scope) → `Section B` (findings or `(no findings)`) → `Section C` (evidence) → `Section D` (falsifiability audit) → `Section E` (parts I cannot verify + what would overturn the PASS) → block-end conclusion line; a format deviation allows one re-issue per section 7.2, counted against the budget.
11. **The fence must match the on-disk enumeration**: after adding any payload, the report fence line plus the MANIFEST entry plus recomputed SHA256SUMS must be added in the same round (P49 missed the 6 r2 fence line and triggered G5(d)).
12. **`gofmt -l .` walks hidden directories**: one-off `.go` files under `.vscode/scripts/` are still listed by `base gofmt` (the `./...` wildcard ignores dot directories, but gofmt itself does not); run `gofmt -w` on temporary Go tooling first.
13. **Recompute hashes with a read-only script**: a write-style normalisation gets `EBUSY` when an editor holds a payload file; a read-only recompute that only writes `SHA256SUMS.txt` avoids it.

## 10. One-line opening for a new session

Here is my project state snapshot (v39); please verify the 5 anchors of section 1 first, continue with the section 4 queue if they match, and treat the repository reality as authoritative otherwise.

[Paste this checklist v39 in full]

[Supplementary documents] `REMAINING_SLICES.md` (P49 completed), `DELIVERY_DIRECTIVE.md` (v1.60), `docs/validation/race-leg-expand{,_EN}.md` (the P49 report; `75990a8`, run `41445657595` green on the first attempt).

Current slice: P49 fully closed; next slice = layer 2 item 3 `CONSOLE-MULTITASK-ASSEMBLY` (protocol first).

## 11. Current status in one sentence

P49 `RACE-LEG-EXPAND` fully closed (`75990a8`, run `41445657595` green on the first attempt on both legs, Race leg six packages all ok); layer 1 cleared; next step = layer 2 item 3 `CONSOLE-MULTITASK-ASSEMBLY` (protocol first).
