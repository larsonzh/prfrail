# Task Execution Checklist (cross-session snapshot v42 · P52 stop-point edition)

Generated 2026-11-03, HEAD = origin/main = `c636dd8` (P52 main commit), directive v1.60. P52 committed and pushed to origin (run `41642139932` green on the first attempt on both legs).
Supersedes: fully replaces v41 and any earlier snapshot. Backfill commits (docs-only) only advance HEAD and do not change the anchor semantics of this section.

This edition (vs v41): ① P52 `ACCEPTANCE-SECURITY-COUNTEREXAMPLES` closed into the baseline (`c636dd8`, seven counterexample test files, zero production changes; **GAP-77 breach proven** (23 out-of-bound bytes read in a real Linux run) + GAP-78/79/80 candidates registered); ② §2 fact pack switches to P52, P51 moves into §3; ③ §4 layer 2 gains top-priority P53 `FIX-GAP-77-SYMLINK-ESCAPE` (layer 3 onwards deferred); ④ §11 updated to "P52 closed, next slice = P53".

## 1. Current anchors (check these five numbers first)

| Item | Value | Check command |
|---|---|---|
| HEAD | `c636dd8` | `git rev-parse --short HEAD` |
| origin/main | `c636dd8` | `git rev-parse --short origin/main` |
| Directive version | **v1.60** | head of `docs/DELIVERY_DIRECTIVE.md` |
| SELFTEST | **697/697** | `node tools/gates/gate.js --selftest` |
| Three scopes | **all TOTAL_FAIL=0** | `gate.js --all --scope=index`, `tree`, `ci` |

If ≥2 items mismatch ⇒ report the repository reality first, then continue.

`tmp` fixed assertion: top level = `.gitkeep` + `gates-cache`; `gates-cache` holds only `g4b`; `tmp/ACCEPTANCE-SECURITY-COUNTEREXAMPLES` was cleared at 8c.

## 2. P52 fact pack

| Field | Value |
|---|---|
| Slice / size | `ACCEPTANCE-SECURITY-COUNTEREXAMPLES` (security-validation `[SLICE]`) / M |
| Goal | Land the GAP-51/52/55/60/64/70 security counterexamples against the acceptance-chain gap of breakpoint A, proving fail-closed under illegal input, path forgery, handle overreach and boundary bypass |
| Changes | 7 new `*_test.go` (T51-1 skip / T52-1/2 / T55-1..7 / T60-1..3 / T64-1..6 / T70-1..4, 23 counterexamples in total); zero production changes |
| Core outputs | **GAP-77 (V-1, HIGH) breach proven**: `openat(O_NOFOLLOW)` call at bounddir_unix.go:64 guards only the final component (:66 is the error-return line); 23 out-of-bound bytes read on Linux (go1.22.12); T51-1 keeps the evidence chain with `t.Skip("Known vulnerability GAP-77: tracked for P53")`. **GAP-78 (V-2)** empty ContextKey, **GAP-79 (V-3)** asymmetric normalisation, **GAP-80 (V-4)** isUnderDirectory prefix residue (all candidates, registered in the P52 ledger and report §2b) |
| Review chain | 1 `DESIGN: DONE`; 3 `PROBE: ALL-GREEN`; 2a-2d all green; 5 r1 `PRE-REVIEW: PASS` (2 Low disposed); 6 `INDEPENDENT SCAN: PASS`; 7 `FINAL REVIEW: PASS` (0 Medium+) |
| Commit / CI | `c636dd8` pushed to origin (never gitee); run `41642139932` green on the first attempt on both legs (both legs `SELFTEST 697/697`, `GATE REPORT scope=ci changed=30`, `TOTAL_FAIL=0`) |
| Metering | R = 0, D = 1 (D1: the 2b cross-compile artefact was written to the system temp directory, violating the v41 §8 temp-location closed set; the controller verified no residue remains) |
| Docs sync | `ARCHITECTURE{,_EN}` +8 each (§5.1 counterexample matrix); `DEV_PLAN{,_EN}` +1 each; ledger P52 ✅ |
| Unclosed | GAP-77 fix deferred to P53; the two junction tests skip where mklink is unavailable; T64-6 simplifies the stage-1 design (backfilled) |

## 3. Completed slices (recent baseline)

| Version | Slice | Main commit + backfill | Conclusion |
|---|---|---|---|
| v1.60 | P45 | `4273d6e` + `3a7541d` | closed |
| v1.60 | P46 | `26ca5a5` + `5e4202e` | closed (R = 1, D = 2) |
| v1.60 | P47 | `bbd1745` + backfill | closed (R = 3, D = 0) |
| v1.60 | P48 | `82e3ac8` + `81dd49f` | closed (R = 9, D = 0; run `41379998477` green on the first attempt) |
| v1.60 | P49 | `75990a8` + backfill | closed (R = 1, D = 2; run `41445657595` green on the first attempt) |
| v1.60 | P50 | `81f9dc9` + backfill | closed (R = 1, D = 1; run `41507722672` green on the first attempt on both legs) |
| v1.60 | P51 | `2b0af2d` + backfill | closed (R = 1, D = 2; run `41569353460` green on the first attempt on both legs) |
| v1.60 | P52 | `c636dd8` + backfill | closed (R = 0, D = 1; run `41642139932` green on the first attempt on both legs; GAP-77 breach proven) |

## 4. Execution queue

General principle: layer 1 clears the record chain → layer 2 strategic mainline → layer 3 pre-release gates → layer 4 exit gates; in-layer order fixed, no jumping between layers.

### Layer 1 (cleared after P49 dequeued)

No remainder. R-03 closed (`RACE-LEG-EXPAND` stood up and closed independently).

### Layer 2 (strategic mainline)

| # | Slice | Identity and contents | Size | 7 budget | Depends on | Boundary |
|---|---|---|---|---|---|---|
| 5 | P53 FIX-GAP-77-SYMLINK-ESCAPE | Fix the Unix middle-symlink escape (openat(O_NOFOLLOW) at bounddir_unix.go:64 guards only the final component; component-by-component openat resolution or O_PATH rework; enable T51-1 by removing the skip) | M | Codex 1 final review | P52 closed | high-risk concurrency/handle class; controller V4 Pro; deliverables include spec updates (below) |

P52 dequeued (`c636dd8`); the only layer-2 item left is item 5, P53.

**P53 deliverables (spec updates, mandatory at launch)**: P53 must update `docs/SECURITY.md` (add the middle-symlink-escape real attack payload and its verification method to the SEC-02 threat matrix) and `docs/TEST_STRATEGY.md` (list the middle-symlink escape as a mandatory counterexample in AT-02 or the relevant security test baseline).

### Layers 3 / 4 (deferred until P53 dequeues)

| # | Slice | Identity and contents | Size | 7 budget | Depends on | Boundary |
|---|---|---|---|---|---|---|
| 6 | CONTRACTS-S7-DURABILITY | contract-side durability semantics and Windows C6 landing | L | Codex 1 final review | R-02 ruling | CONTRACTS first, then code |
| 7 | EXIT-SEED-BUILD-EVIDENCE | AT-13 seed-driven build evidence | M | Codex 1 final review | item 6 | evidence chain only, no CI structure change |
| 8 | AT23-RESTRICTED-CANDIDATE | AT-23 restricted real candidate and sign-off | M | Codex 1 final review | R-08 / R-09 rulings | no sign-off without a real candidate authorisation |

Suspended (not scheduled): `G4B-CORPUS-DECOUPLE` (trigger-gated), the five low-priority items (no trigger, no launch), the rest of round-2 gaps (candidate pool).

## 5. Observations

- Observing (13): `OB-1 / 7 / 10 / 33 / 87 / 99 / 100 / 104 / 105 / 107 / 111 / 113 / 135`.
- Open (11): `OB-19 / 21 / 52 / 53 / 55 / 57 / 58 / 59 / 61 / 75 / 120`.
- To-do (1): `OB-69`. Non-final total **25**.
- GAP registry: **GAP-77 breached and proven (P53 fix)**; GAP-78/79/80 candidates (registered in the P52 report §2b and the ledger).
- Mechanical readings: max OB in §12.4 = **OB-136**; max deferral = the 25th; ADR max = **28**; DR max = **13** (DR-13 fixed).

## 6. Awaiting user rulings (R-01 through R-11)

| Id | Item | Status |
|---|---|---|
| R-01 | P46 strategic review stand-up | completed |
| R-02 | `CONTRACTS-S7-DURABILITY` pacing and ownership | pending (blocks layer 3) |
| R-03 | P45 follow-up slice ownership | completed (`RACE-LEG-EXPAND` stood up and closed independently, layer 1 cleared) |
| R-04 | DR-13 "fixed + observing" wording | pending |
| R-05 | whether DR13-FIX-D3 stands up independently | completed (stood up and closed independently) |
| R-06 | OB-55 ownership and T027-B pacing | pending |
| R-07 | MVA, S1 moved-out items and S2 promise narrowing | pending |
| R-08 | candidate, Windows precondition P and probe authorisation | pending (blocks layer 4) |
| R-09 | AT-24 `unknown` acceptance boundary | pending (blocks P3 closure) |
| R-10 | total budget and per-slice model caps | pending (blocks all scheduling) |
| R-11 | historical-correction scope | completed (tier B + itemised differentiation) |

## 7. Launch and step-0 general rules

1. Check anchors → step 0: register this slice in the ledger (`R/D = 0` baseline, temp artefacts at `tmp/<slice-id>/`), land `input-list.md` and `freedom-list.md`.
2. Every fix slice must provide **regression tests** and explicitly register "tests/lint not executed".
3. Stage 9 three scopes all green + `SELFTEST 697/697` + `go build/vet/test ./...` all green before committing.
4. Model binding (hard): ①⑤ `deep-reasoner` + `DeepSeek V4 Pro (deepseek)`; ②③⑧ `DeepSeek V4.1 Flash (deepseek)`; ⑥ `independent-reviewer` + `MAI-Code-1.1-Flash (copilot)`; ⑦ `independent-reviewer` + `GPT-5.3-Codex (copilot)`; strategic-review rounds 1/3 `independent-reviewer` + `GPT-6.1 Sol (copilot)`; batch assists `deep-reasoner` + `GLM-5.3-Flash (glm)`. **Every `runSubagent` must pass both `agentName` and `model`.**
5. **Pacing discipline**: no routine stops — stage 4 green ⇒ dispatch stage 5 directly; stages 5/6/7 report Medium+ ⇒ fix and re-run that stage in place; **stage 10 is the only routine stop point**.
6. Cross-version/cross-platform triage: for CI-vs-local divergence check toolchain versions and build-tag coverage first; after cross-building re-verify `go env GOOS GOARCH`; before deleting a "locally unused" symbol run at least one opposite-`GOOS` `go vet`.
7. Batch assists: each slice dispatches `deep-reasoner` + `GLM-5.3-Flash (glm)` **≥ 2 times** (strategy and output caps in §8); **strictly serial**, no parallel subagents (R1.1 / R5.1).
8. Controller model policy: default `DeepSeek V4.1 Flash`; switch to Pro for bottom-level concurrency/state-machine/handle-class slices, single designs over 10,000 characters, deep root-cause work over large CI logs, or two consecutive Flash behaviour drifts; **any downgrade must be independently registered per the OB-37 meta-rule**. P53 is a high-risk concurrency/handle class: the controller starts as `DeepSeek V4 Pro`.
9. **Carriers are not substitutable**: stage 2 must dispatch `prfrail-implementer`, stage 3 `prfrail-tester`; the controller only orchestrates, lands files, runs gates and reports — **it must not write product code or tests**.
10. **The launch message is an artefact**: step 0 must land the user launch message verbatim as `docs/validation/evidence/<slice-id>/00-briefing.md`, list it in `MANIFEST.md` and `SHA256SUMS.txt`, and declare it in the `artifacts` fence; without it stage 1 must not start.
11. **Linearisation-point discipline**: for persistent state machines "persist succeeded" must be the linearisation point; the in-memory CAS is only a same-process fast-reject cache; write "credit taken ⟺ marker on disk" with a positive red test.
12. **Evidence-pack and fence sync**: any new evidence payload must add a `MANIFEST.md` entry, recompute `SHA256SUMS.txt` and add an `artifacts` fence line in the same turn; missing any one gets caught by `G8-b` or `G5(d)`.
13. **Reviewer archival discipline**: stages ⑤⑥⑦ output must be archived verbatim; the controller must not rebuild from bullet points; if a summary is truly needed, put it in a separate file disclosed outside the fence.
14. **Assembly-validation slice discipline**: assertions must land on falsifiable mechanical criteria; constant identities must be declared non-sentinel; the report must state whether production code changed at all; an unsupported capability of a layer must not be read as supported.
15. **Evaluation-class slices stop at 10 after stage 9**: with no implementation object, stop at 10 for same-turn authorisation after stage 9 is green; no self-`commit`/`push`.
16. **No single output over 10,000 characters**: batch before the limit; no self-restart; read on demand.
17. **Compaction ends the slice**: at 100K context or a compaction notice ⇒ end the slice immediately; one session runs one L-tier slice.
18. **The launch message must include a "docs sync" section**: when a slice changes the project phase state or contract capability, stage 8 must update `DEV_PLAN.md`, `ARCHITECTURE.md` and peers; a missing section is a structural defect of the launch message.
19. **Environment preflight and reuse iron rule (v39)**: if a later slice needs a local Linux environment, remote SSH, network proxy or other host-dependent environment, the controller must declare it explicitly at step 0. When an environment gap appears (e.g. a command not found), **probe existing environments first (PATH, config, find)**; only install minimally after confirming absence; never blindly reinstall. **Remote SSH sessions must explicitly `source ~/.bashrc` or use `bash -lc` (non-interactive shells skip user PATH and `.bashrc` often carries a non-interactive early-return guard), and must probe user-level directories (`~/.local/bin`, `~/.local/opt`)**; install only after confirming no existing tool, pinned to a CI-aligned version with checksum verification. P53's Linux validation can reuse the remote VM `10.0.0.199` (go1.22.12 at `~/go1.22` already in place).

## 8. Cross-slice collaboration discipline (the launch message must implement)

| Element | Content |
|---|---|
| Carrier table | ①⑤⑥⑦ use `deep-reasoner` / `independent-reviewer` (read-only); ②③⑧ use `prfrail-implementer` / `prfrail-tester` / `prfrail-documenter` (writable); `prfrail-tester` must not run review-layer stages; batch assists must be `deep-reasoner`, never `quick-researcher` |
| Dispatch-table mechanical read | every row's `agentName` in the stage-10 dispatch table must be read verbatim from `transcript.arguments.agentName`, with the raw script stdout attached |
| Serial dispatch | `runSubagent` is serial/blocking; no parallel subagent calls (R1.1 / R5.1); always pass `agentName` and `model` |
| R2.1 add agentName | every `runSubagent` `arguments` must contain both `agentName` and `model` |
| Commit messages | English only, `<type>: <summary>`; with any non-ASCII character use `git commit -F <UTF-8 file>` |
| Independently checkable vs not | independently checkable outputs (enumerations, fact comparisons, reference lists) are not voided by model errors; non-checkable ones (judgements, conclusions, review opinions) must be redone |
| ⑥ before ⑦ | stage 6 must finish before stage 7; if stage 6 is re-dispatched and docs changed, the stage-7 conclusion must be re-evaluated |
| Temp-artefact location and cleanup | always `tmp/<slice-id>/`; only 8c cleans up, and only this slice's directory — never recursively delete a parent; **character counting or intermediate recording must use inline scripts (e.g. `node -e "console.log('...'.length)"`); never create temp files on any disk path (especially `$env:TEMP` or the system temp dir)** |
| Long-command redirection | never PowerShell `>`; use `\| Tee-Object -FilePath <path>` or `cmd /c "..."` |
| Gate execution | intermediate stages use `--all --skip-base`; all green is required before committing |
| Discipline E · index-scope order | precise `git add` first → then `--scope=index`; the receipt must report `G6-1 scanned N` with N > 0 |
| Discipline F · spec-over-example | follow the spec when a launch example conflicts, but register the conflict in the receipt |
| Discipline G · fence = landing snapshot | `artifacts` fence sorted by path, TAB separated, CN/EN byte-identical; new files must add matching lines |
| Discipline H · opposite-GOOS recheck | before deleting/changing constants or functions run at least one opposite-`GOOS` `go vet` |
| Disciplines I–N | verdict closed set / review semantics / persistence first / three-way self-consistency and fence sync / archive verbatim / mechanical criteria and non-sentinel declarations (see §7 items 11-14) |
| **Discipline O · output-length defence** | no single output over 10,000 characters; batch before the limit per §7; no self-restart; read on demand; on a stop-loss trigger stop and report immediately |
| **Discipline P · controller model switch trigger** | see §7 item 8: bottom-level concurrency/state-machine/handle, over-long design docs, deep CI-log root-cause, two consecutive Flash drifts ⇒ switch to Pro; downgrades registered per OB-37 |
| **Discipline Q · batch-assist output cap** | dispatches must state a character cap; the receipt reports the actual characters; over-cap but complete is a metering fact (not a §7 stop-loss); a platform-measured ≥ 10,190 characters may serve as cap calibration |
| **Discipline R · stale-tense batch reconciliation** | ledger-vs-fact mismatches (GAP-12 / GAP-19 / GAP-24 class) go through a one-shot reconciliation slice with code and commits as the fact source; no line-by-line patching |
| **Discipline S · session length and compaction** | at 100K context or a compaction notice ⇒ end the slice immediately; one session runs one L-tier slice; compaction drops early hard stop points and authorisation discipline and induces drift (D-02 was this class) |
| **Discipline T · quantifier shapes** | the ledger and reports must not write "identifier immediately followed by a quantifier": the quantifier is read by `G6-4` as a line-count literal; write "that entry" or "that record" |
| **Discipline U · never quote trigger strings verbatim** | reports must not quote verbatim the strings that trigger criteria: restating the same literal when explaining a red would re-trigger; use descriptive wording instead |
| **Discipline V · split large tasks and guard half-written artefacts** | large tasks must be split into small serial steps with minimal payloads; on timeout check half-written artefacts with `git status`/`git diff`; if half-written and red, capture evidence and `git checkout` precisely — never blindly retry (P50 proof: one-shot stage-2 dispatch timed out five times with two half-written remnants; splitting into 2a-d went green throughout) |

**GLM batch-assist policy**: `GLM-5.3-Flash (glm)` is plentiful; "≥ 2 times" is only the floor — mechanical checks, long-document excerpts, bilingual mirror checks, dispatch drafts and log extraction should be offloaded; no complex judgement or picking; every dispatch counts into the stage-10 cost and metering (output cap per discipline Q).

**P46 lessons (v36 onward)**: ① **evaluation-class slices must stop at 10 after stage 9** — the D-02 root cause was treating one as an implementation pipeline and committing directly; ② **§7 readings are self-referential**: editing the body moves the six-shape readings, so "edit the body first, backfill §7 last, then recompute until stable"; ③ **gate criteria will catch wording**: `G2` catches placeholders, `G6-4` catches count-class literals (placeholders must not remain); the disposition is rewording, never expanding the whitelist.

## 9. Checklist meta-rules

**Quote by first-sentence, not "item N".**

| Rule (first-sentence quote) | Content |
|---|---|
| "Regenerate on every 2 slices / pending-ruling change / directive version bump / queue change / user request / ≥2 anchor mismatches" | when to regenerate |
| "Pending-ruling items must be listed by name" | §6 must name every item |
| "Queue items must carry the six elements" | identity / contents / size tier / stage-7 budget / depends on / boundary |
| "Rule obligations ≠ gate obligations" | never read rule text as a mechanical criterion |
| "No physical-clock assertions" | observation clauses trigger on events, not "N months later" |
| "Gates green before commit" | hard gate |
| "Dispatch-table agentName mechanical read" | read verbatim from `transcript.arguments.agentName` |
| "Governance slices do not backfill DEV_PLAN" | the authoritative ledger for directive-governance slices = appendix C.0-C.0x |
| "A slice name never changes once in the ledger" | name stability |
| "Policy wins over authorisation counts on conflict" | never use this to judge overreach or refuse execution |
| "One fact described in one place" | body or pointer, pick one |
| "Controller-registered text may be worded to completion; reviewer originals must be verbatim" | different standards on the two sides |
| "Factual errors are independent of design decisions" | ledger-report mutual exclusion is a factual error, not a design ruling |
| "The stage-10 receipt must attach a 'fact pack for the next checklist'" | fields: HEAD, CI run, ledger additions, new observations, new discipline candidates, max OB / deferral, ADR/DR max, tmp assertion, differences from the previous edition, unclosed items |
| "[ROUTINE] small fixes take no new slice number but lose no discipline" | same stage-9 gates + CI + fence lines; no semantic change smuggled in |
| "The launch message is an artefact" | `00-briefing.md` landed + MANIFEST and SHA256SUMS + fence declaration |
| "Carrier deviation = highest independence downgrade" | register a D item with an independence-impact note |
| "Mechanical criteria beat witness statements" | zero-model/zero-side-effect assertions must be falsifiable; constant identities marked non-sentinel |
| "Concurrency assertions must cover survival and release" | besides the failure face, assert the winner survives and the slot is reusable after release |
| "Layer boundaries must not widen the reading" | an unsupported capability must be declared explicitly |
| "Evaluation-class slices stop at 10 after stage 9" | no implementation object; stop at 10 for same-turn authorisation after stage 9 green; no self-commit/push |
| "No single output over 10000 characters" | batch before the limit; no self-restart; read on demand |
| "Compaction ends the slice" | 100K context or a compaction notice ⇒ end the slice; one L-tier slice per session |

### Gate pitfalls (P48 six + P49 seven, kept for reuse)

### Since P48

1. **`.go` copies under `tmp/` get compiled by `go build ./...`**: source copies in review input packs must use `*.go.txt`, otherwise the `base go build/vet/test` legs go red.
2. **`G7-b` fence disposition closed set**: unrepo'd temp artefacts write `local-only`; `local` reports unknown disposition.
3. **A report carrying the `INDEPENDENT SCAN` marker must pair a `sixstage` fence**: inside the fence use `## Section A` through `## Section E` once each, keep only the last non-empty line as the conclusion, and disclose merged differences outside the fence.
4. **Writing `.md` must add BOM explicitly**: a node rewrite missing the BOM parameter immediately triggers `G4a` (P48 hit it twice).
5. **`G8-b` wants MANIFEST to list payloads one by one**: a directory row is judged "payload not listed".
6. **Readings must be backfilled last**: archival changes the `changed=` count and the `scanned=` readings; re-measure before finalising the report.

### Since P49

7. **Changing `.github/workflows/*.yml` must sync the `internal/release` workflow freeze contract**: `workflow_test.go`'s `expectedScripts` locks the SHA-256 of every script step body verbatim; stage-1 impact searches must cover that contract (P49's stage 1 missed it, reddening the contract family after r1; closed after r2 sync).
8. **Remote SSH preflight**: non-interactive sessions skip user PATH (`source ~/.bashrc` is often swallowed by an early-return guard); use `bash -lc` and probe user-level directories like `~/.local/bin` (P49 thereby missed the VM's existing go1.27.1, registered as a deviation).
9. **PowerShell → ssh nested quoting**: multiple quoting layers and `$(...)` get parsed by PowerShell first; complex cross-shell commands should upload a remote shell script and execute it.
10. **Stage-6 output must be A-E shaped**: `Section A` (scope) → `Section B` (findings or `(no findings)`) → `Section C` (evidence) → `Section D` (falsifiability audit) → `Section E` (unverifiable parts + reversal conditions) → the block-final conclusion line; a format deviation allows one resend per §7.2, counted into the budget.
11. **Fence and disk enumeration must agree**: after adding any payload, add the report fence line + MANIFEST entry + recompute SHA256SUMS in the same turn (P49 missed the 6-r2 fence line and triggered G5(d)).
12. **`gofmt -l .` scans into hidden directories**: one-off `.go` under `.vscode/scripts/` is listed by `base gofmt` (the `./...` wildcard ignores dot-directories but a bare gofmt does not); run `gofmt -w` before parking temp Go tools.
13. **Recompute hashes with a read-only script**: write-style normalisation fails with `EBUSY` while the editor holds a payload; a read-only recompute that only writes `SHA256SUMS.txt` works around it.

### Since P52 (G1(d)/G8-b record)

14. **A subsection heading's status symbol must match its status row**: `G1(d)` compares the ledger subsection heading (e.g. `· 🔄 in progress`) against the status row (`✅`) closed set; change both at closure — missing one goes red.
15. **Recompute SHA256SUMS after any line-number unification**: content-class corrections (e.g. `:66`→`:64`) change payload hashes and `G8-b` reports a mismatch immediately; recompute in the same turn after any payload content change.

## 10. One-line opener for a new session

Here is my project state snapshot (v42): first check the five anchors in section 1; if they match, continue with the section-4 queue; otherwise report the repository reality.

【paste this v42 checklist in full】

【supplementary docs】`REMAINING_SLICES.md` (P52 completed), `DELIVERY_DIRECTIVE.md` (v1.60), `docs/validation/acceptance-security-counterexamples{,_EN}.md` (P52 report, `c636dd8`, run in §2).

Current slice: P52 closed; next slice = layer-2 item 5 `P53 FIX-GAP-77-SYMLINK-ESCAPE` (controller V4 Pro, deep Codex review).

## 11. Current state in one sentence

P52 `ACCEPTANCE-SECURITY-COUNTEREXAMPLES` closed (`c636dd8` pushed to origin, run `41642139932` green on the first attempt on both legs; 23 counterexamples landed, **GAP-77 breach proven + GAP-78/79/80 candidates**, zero production changes, R=0, D=1); next step = layer-2 item 5 P53 `FIX-GAP-77-SYMLINK-ESCAPE` (component-by-component openat resolution or O_PATH rework; deliverables must update `docs/SECURITY.md` SEC-02 and `docs/TEST_STRATEGY.md` AT-02).
