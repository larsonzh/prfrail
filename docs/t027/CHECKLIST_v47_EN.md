# Task execution checklist (cross-session snapshot v47 · P57 closure)

Generated: 2026-11-05, HEAD = origin/main = `3884351` (the P57 commit plus the O-1/O-3 cleanup), directive v1.61. Both commits are pushed to origin (`7a43dfa` run `41991695581`, `3884351` run `41992956993`, each green on both legs; never pushed to gitee).
Supersedes v46 and every earlier snapshot. A docs-only backfill commit only moves HEAD; it does not change the anchor semantics of this section.

Changes in this version (vs v46): (1) **P57 `P5-P6-CLOSEOUT-ASSESSMENT` closed** (`7a43dfa`; run `41991695581` green on both legs; P5 ruled that `T027-B` needs no new slice and P6 produced the S1 Exit eight-gate inventory); (2) **O-1/O-3 document cleanup** (`3884351`; run `41992956993` green on both legs; the AT-24/T028 staleness in `s1-exit` and the line drift in `s1-reach` corrected); (3) section 2 now carries the P57 fact pack plus the cleanup and P56 moved into section 3; (4) section 3 gained two rows; (5) **section 4 marks layer 0 ✅ complete and dequeued**, layer 1 stays ⛔ and layers 2 / 3 stay parked until named; (6) section 5 switches to the P57 observations and records that O-1/O-3 are cleared; (7) section 6 gains a P57 ruling note (R-1 / R-2 untouched, R-3 follows R-08 / R-09); (8) **the point-in-time staleness in section 4.3 of the P57 report was corrected by the O-1/O-3 cleanup (`3884351`)**, while the report itself stays unchanged per the ruling; (9) section 9 has no new gate pitfall in v47.

## 1. Current anchors (check these 5 numbers first)

| Item | Value | Command |
|---|---|---|
| HEAD | `3884351` | `git rev-parse --short HEAD` |
| origin/main | `3884351` | `git rev-parse --short origin/main` |
| Directive version | **v1.61** | header of `docs/DELIVERY_DIRECTIVE.md` |
| SELFTEST | **703/703** | `node tools/gates/gate.js --selftest` |
| Three scopes | **all TOTAL_FAIL=0** | `gate.js --all --scope=index`, `tree`, `ci` |

If two or more items disagree, trust the repository and report the final readings before continuing.

Fixed `tmp` assertion: the top level is `.gitkeep` plus `gates-cache`; `gates-cache` contains only `g4b`; P57's `tmp/P57-CLOSEOUT-ASSESSMENT/` and the cleanup's `tmp/ROUTINE-O1-O3-CLEANUP/` were both removed at 8c.

## 2. P57 fact pack

| Field | Value |
|---|---|
| Slice / size | `P57-CLOSEOUT-ASSESSMENT` (assessment class `[SLICE]`) / S; plus a `[ROUTINE]` cleanup (O-1 / O-3) |
| Goal | (1) decide whether `T027-B` needs a new slice; (2) produce the S1 Exit remaining-gate close-out inventory; (3) also correct the staleness and line drift in `s1-exit` and `s1-reach` |
| Change surface | new `docs/validation/p57-closeout-assessment{,_EN}.md` (with its evidence pack); write-back to `DEV_PLAN{,_EN}`, `REMAINING_SLICES{,_EN}`, `S1_SUPPORT_MATRIX{,_EN}`, `S1_RELEASE_NOTES{,_EN}`; cleanup edits to `s1-exit{,_EN}` and `s1-reach-reachability{,_EN}` |
| Core output | (1) **P5: `T027-B` gets no new slice** (C1-C7 and N1-N4 already assigned to P39-P44), with the three residue prongs R-1 / R-2 / R-3; (2) **P6: the S1 Exit eight-gate inventory** (S-G1 to S-G8) plus three registered staleness items; (3) all three were corrected by the cleanup commit |
| Review chain | 5 twice (r1 `FINDINGS` -> r2 `FINDINGS`); the single remaining r2 Medium was closed by **deterministic self-proof** (24/24 citation assertions plus a three-way equal inventory) and the user accepted it without a third round; 6 and 7 skipped per the kickoff message |
| Commit / CI | `7a43dfa` (17 entries) run `41991695581` green on both legs; `3884351` (4 entries) run `41992956993` green on both legs; the three conditional jobs skipped by design |
| Metering | `R = 0`, `D = 0`; `runSubagent` twice (5 r1, 5 r2, both `deep-reasoner` with `DeepSeek V4 Pro`); probes / paid external calls 0 |
| Rulings | the P5 / P6 verdicts and the deterministic closure accepted; R-1 not enqueued, R-2 not ruled, R-3 follows R-08 / R-09; O-1 / O-3 authorised as a `[ROUTINE]` cleanup outside P57 |
| Doc sync | new report pair; `DEV_PLAN` gained the "T027-B closeout assessment" paragraph and the "T018 closeout inventory" pointer; the ledger P5 row now reads 🔄 tracking and P6 reads ✅ assessed; the S1 matrix and release notes carry "not applicable + reason" only; the cleanup added two `artifacts` fences |
| Capability boundary | no product code / CI / schema / fixtures / directive change (v1.61 untouched); no new OB number (observations stay slice-level); no gate verdict is claimed to have changed |
| Not closed | section 4.3 of the P57 report is a `7a43dfa` point-in-time snapshot (its staleness was corrected in `3884351`, the report itself stays unchanged per the ruling); the 5 budget of 1+1 is tight for assessment-class slices |

## 3. Completed slices (recent baseline)

| Version | Slice | Main commit + backfill | Verdict |
|---|---|---|---|
| v1.60 | P49 | `75990a8` + backfill | closed (R = 1, D = 2; run `41445657595` green on the first attempt) |
| v1.60 | P50 | `81f9dc9` + backfill | closed (R = 1, D = 1; run `41507722672`) |
| v1.60 | P51 | `2b0af2d` + backfill | closed (R = 1, D = 2; run `41569353460`) |
| v1.60 | P52 | `c636dd8` + backfill | closed (R = 0, D = 1; run `41642139932`; GAP-77 breach recorded) |
| v1.60 | P53 | `7802819` + backfill | closed (R = 0, D = 0; run `41707160563`; GAP-77 fixed) |
| v1.60 | P54 | `7c0763a` + `f38e0ce` | closed (R = 1, D = 1; run `41840557379`; C6 judgment path plus the downgrade ruling) |
| v1.60 | P55 | `63e2e24` + `c76c6ce` | closed (R = 0, D = 1; runs `41925623071` and `41928327517` both green; AT-13 two-sided evidence chain) |
| v1.61 | P56 | `504caff` | closed (R = 2, D = 1; run `41978355255` green on both legs; the G4b corpus enumeration decoupled from the index, SELFTEST 697 to 703) |
| v1.61 | P57 | `7a43dfa` | closed (assessment class; R = 0, D = 0; run `41991695581` green on both legs; P5 no-slice verdict plus the P6 eight-gate inventory) |
| v1.61 | O-1/O-3 cleanup | `3884351` | closed (`[ROUTINE]`; 4 entries; run `41992956993` green on both legs; the staleness and line drift corrected) |

## 4. Execution queue

**Basis for this ordering**: the layer-0 front-end stocktaking is complete and **dequeued** after P57; the only main line left is the layer-1 hard gate (which needs the R-08 / R-09 rulings to clear its block); layer 2 waits for the user to name it; layer 3 stays parked. The order inside a layer is fixed and layers are not skipped.

### Layer 0 Front-end stocktaking (✅ complete, dequeued)

| Order | Item | Content | Size | Depends on | Boundary |
|---|---|---|---|---|---|
| 0-1 | P5 `T027-B` close-out list | ✅ done: **no new slice**; the three residue prongs are R-1 (P44 end-to-end wiring, not enqueued), R-2 (OB-55 ownership) and R-3 (real CLI end-to-end) | evaluation action | none | dequeued; see section 3 of the P57 report |
| 0-2 | P6 `S1-EXIT-DOC` close-out list | ✅ done: produced the **S1 Exit eight-gate inventory** (S-G1 to S-G8) as the single reference for the later main hard gate | evaluation action | none | dequeued; see section 4 of the P57 report |

> Rule: both were closed and dequeued by P57; their conclusions are the **single reference** for layer 1 and S1 Exit, and no equivalent evaluation action may be started again.

### Layer 1 Main hard gate (advance as soon as the blocker clears)

| Order | Slice | Identifier and content | Size | 7 budget | Depends on | Boundary |
|---|---|---|---|---|---|---|
| 8 | AT23-RESTRICTED-CANDIDATE | AT-23 restricted real candidate and sign-off | M | Codex one final review | **R-08 / R-09 rulings (unruled ⇒ ⛔ blocked)** | no sign-off without a real candidate authorisation |

### Layer 2 Routine governance (backlog)

| Order | Item | Content | Start condition |
|---|---|---|---|
| G-1 | P7 governance-side split-out OB roll-up | Start an evaluation action per section 1.8.1 and roll up the OBs split out in the later-batch section | the user names it (preferred after layer 0) |

### Layer 3 Low-priority alternatives

| Order | Item | Content | Start condition |
|---|---|---|---|
| L-1 | The five low-priority items | nothing starts without a trigger | the user must name them explicitly |

> Layer boundary: layer 0 is complete and dequeued and no longer produces evaluation actions; layer 1's blocker is not lifted by finishing layer 0; neither layer 2 nor layer 3 may be read as "the main line resumed".

## 5. Observations

- Observing (13): `OB-1 / 7 / 10 / 33 / 87 / 99 / 100 / 104 / 105 / 107 / 111 / 113 / 135`.
- To be discussed (10): `OB-19 / 21 / 52 / 53 / 55 / 57 / 58 / 59 / 61 / 75`. **OB-120 was resolved by P56** (`Resolved`), which is why the count dropped from 11 in v45 to 10.
- To do (1): `OB-69`. Non-terminal total **24**.
- GAP register: **GAP-77 fixed (P53 closed)**; GAP-78/79/80 are candidates (registered in the P52 report section 2b and the ledger).
- New P57 observations (slice-level, no OB number): (1) the 5 budget of 1+1 is tight for assessment-class slices (the remaining r2 Medium could only be closed deterministically); (2) no dedicated `ci`-scope assertion (carried over from P56); (3) the dispatch-table mechanical-reading gap (the debug log still has no transcript). **O-1 / O-3 (the staleness and line drift in `s1-exit` and `s1-reach`) were cleared in `3884351`**.
- Mechanical counters: max OB in section 12.4 = **OB-136**; max deferral = the 25th; max ADR = **30**; max DR = **13**.

## 6. Awaiting user rulings (R-01 to R-11)

| Id | Item | Status |
|---|---|---|
| R-01 | P46 strategic review chartering | done |
| R-02 | `CONTRACTS-S7-DURABILITY` cadence and ownership | done (P54 closed) |
| R-03 | P45 follow-up slice ownership | done (`RACE-LEG-EXPAND` closed) |
| R-04 | DR-13 "fixed + observed" calibre | pending |
| R-05 | Charter DR13-FIX-D3 as its own slice | done (chartered and closed) |
| R-06 | OB-55 ownership and the T027-B cadence | pending |
| R-07 | MVA, the S1 removals and the S2 scope narrowing | pending |
| R-08 | Candidate, Windows precondition P and probe authorisation | **pending (blocks layer-1 item 8)** |
| R-09 | AT-24 `unknown` acceptance boundary | **pending (blocks the layer-1 item 8 sign-off calibre)** |
| R-10 | Overall budget and per-slice model caps | pending (blocks scheduling) |
| R-11 | Scope of the historical corrections | done (tier B plus per-item differentiation) |

> P57 ruling note (2026-11-05): R-1 (P44 end-to-end wiring) is **not enqueued**; R-2 (OB-55 ownership) is **not ruled** and stays under R-06; R-3 follows the R-08 / R-09 rulings. R-04 / R-08 / R-09 / R-10 stay pending.

## 7. Start-up and step-0 general rules

1. Check the anchors → step 0: register the slice in the ledger (`R/D = 0` baseline, temp location `tmp/<slice-id>/`) and write `input-list.md` and `freedom-list.md`.
2. Every fixing slice must supply **regression tests** and must declare "tests / lint not executed" explicitly.
3. Submitting requires 9 green on all three scopes, `SELFTEST 703/703` and green `go build/vet/test ./...`.
4. Model binding (hard): 1 and 5 use `deep-reasoner` + `DeepSeek V4 Pro (deepseek)`; 2, 3 and 8 use `DeepSeek V4.1 Flash (deepseek)`; 6 uses `independent-reviewer` + `MAI-Code-1.1-Flash (copilot)`; 7 uses `independent-reviewer` + `GPT-5.3-Codex (copilot)`; strategic-review rounds 1 and 3 use `independent-reviewer` + `GPT-6.1 Sol (copilot)`; batch auxiliary uses `deep-reasoner` + `GLM-5.3-Flash (glm)`. **Every `runSubagent` must pass both `agentName` and `model`.**
5. **Cadence discipline**: normally do not stop - green at 4 ⇒ dispatch 5 directly; a Medium+ from 5/6/7 ⇒ fix and re-run that ring in place; **10 is the only routine stop point**.
6. Cross-version and cross-platform triage: when CI and the host disagree, check toolchain versions and build-tag coverage first; re-check `go env GOOS GOARCH` after a cross build; before deleting a symbol "unused on this host", run `go vet` for the opposite `GOOS` at least once.
7. Batch auxiliary: dispatch `deep-reasoner` + `GLM-5.3-Flash (glm)` **at least twice** per slice (strategy and output caps in section 8); **strictly serial**, never parallel subagents.
8. Controller model policy: default `DeepSeek V4.1 Flash`; switch to Pro for slices touching low-level concurrency, state machines or handle passing, for a single design document above 10,000 characters, for deep root-cause work on large CI logs, or after two consecutive behavioural drifts; **any downgrade must be registered independently under the OB-37 meta-rule**.
9. **Role carriers are not substitutable**: 2 must be dispatched to `prfrail-implementer` and 3 to `prfrail-tester`; the controller only orchestrates, persists, gates and writes the receipt. **Exception**: pure evidence-supplement or log-capture slices may be self-executed by the controller, but only when the user says so in the same turn.
10. **The start message is an artefact**: step 0 must persist the user's start message verbatim as `docs/validation/evidence/<slice-id>/00-briefing.md`, list it in `MANIFEST.md` and `SHA256SUMS.txt`, and declare it in the `artifacts` fence; without it, stage 1 may not start.
11. **Linearisation-point discipline**: for persistent state machines, "persisted successfully" must be the linearisation point; an in-memory CAS is only a same-process fast-reject cache.
12. **Evidence package and fence sync**: a new evidence payload requires, in the same round, a `MANIFEST.md` row, a recomputed `SHA256SUMS.txt` and a line in the `artifacts` fence.
13. **Reviewer archiving discipline**: the outputs of 5/6/7 must be archived verbatim (`04-reviews/`); the controller may not rebuild them from bullet points.
14. **Assembly-verification slice discipline**: assertions must rest on falsifiable mechanical checks; constant identities must be declared non-sentinels; the report must state whether production code changed.
15. **Evaluation slices must stop at 10 after 9**: an evaluation slice has no implementation object, so after green at 9 it stops at 10 and waits for the same-turn authorisation.
16. **A single output must not exceed 10,000 characters**: batch before the limit; never self-restart; read on demand.
17. **Compaction ends the slice**: at a 100K context or a compaction prompt, end the current slice immediately; one session runs at most one L-tier slice.
18. **The start message must contain a "document sync" section**: if the slice changes the project stage or a contract capability, stage 8 must update the authoritative documents; omitting it is a structural defect.
19. **Environment pre-check and reuse rules**: declare external environments at step 0; probe for missing tools (PATH, configuration, `find`) before a minimal install; **remote SSH must use `bash -lc`** to load the login environment and must probe user-level directories; installs must pin the version that matches CI and verify checksums.
20. **Evidence capture must be replayable and faithful**: a captured `COMMAND:` line must match the real argv; if only an abbreviated display form is possible, file an erratum in `anomalies.md` and write the faithful command plus a `commandErrata` pointer into the derived `index.json`; drift or constructed inputs need replayable rebuild steps in `replay.md`.

## 8. Cross-slice collaboration discipline (must be embodied in the start message)

| Element | Content |
|---|---|
| Carrier table | 1/5/6/7 use `deep-reasoner` / `independent-reviewer` (read-only); 2/3/8 use `prfrail-implementer` / `prfrail-tester` / `prfrail-documenter` (writable); `prfrail-tester` may never run a review ring; batch auxiliary must be `deep-reasoner`, never `quick-researcher` |
| Dispatch-table mechanical reading | every `agentName` in the 10 receipt must be read verbatim from `transcript.arguments.agentName` with the script stdout attached; **if the session's debug log carries no transcript (as in P56), the receipt must state "not mechanically verified" and register it as a known fact** |
| Dispatch serialisation | `runSubagent` is serial and blocking; never call several subagents in parallel; always pass `agentName` and `model` |
| Commit messages | always English, `<type>: <summary>`; any non-ASCII content requires `git commit -F <UTF-8 file>` |
| Independently checkable vs not | independently checkable outputs (enumerations, fact comparisons, reference lists) do not invalidate the slice when the model errs; non-checkable ones (judgements, conclusions, review opinions) must be redone |
| 6 before 7 | 6 must finish before 7; when 6 is re-dispatched and documents changed, 7's conclusion must be re-evaluated |
| Temp location and cleanup | always `tmp/<slice-id>/`; only 8c may clean it, and only that slice's directory - never a parent; measure character counts or record intermediate state with inline scripts |
| Long-command redirection | do not use PowerShell `>`; use `\| Tee-Object -FilePath <path>` or `cmd /c "..."` |
| Gate execution | use `--all --skip-base` for intermediate steps; everything must be green before a commit |
| Discipline E · index scope order | stage precisely with `git add` first, then run `--scope=index`; the receipt must report `G6-1 scanned N` with N greater than 0 |
| Discipline F · specification wins | when a start example conflicts with the specification, follow the specification and register the conflict in the receipt |
| Discipline G · fence is the landing snapshot | the `artifacts` fence is path-sorted, TAB separated and byte-identical between CN and EN; touching a new file requires a new row |
| Discipline H · opposite-GOOS check | before deleting or changing a constant or function, run `go vet` for the opposite `GOOS` at least once |
| Discipline O · output-length defence | a single output must not exceed 10,000 characters; batch before the limit; never self-restart |
| Discipline P · controller model switch trigger | see section 7 item 8; downgrades are registered under OB-37 |
| Discipline Q · batch auxiliary output cap | the dispatch must state a character cap; the receipt must report the actual count |
| Discipline R · batch clearing of stale ledger text | when the ledger contradicts the facts, run a one-off clearing slice with code and commits as the source of truth |
| Discipline S · session length and compaction | at a 100K context or a compaction prompt, end the current slice immediately |
| Discipline T · quantifier shape | the ledger and reports must not write "identifier followed by a quantifier": G6-4 reads it as a line-count literal; write "that entry" or "that record" |
| Discipline U · no verbatim trigger strings | reports must not quote the strings that trip a criterion |
| Discipline V · split large tasks, avoid half-written pollution | split large tasks into small serial dispatches; after a timeout always inspect `git status` / `git diff` for partial output and roll back precisely with `git checkout` when needed |
| Discipline W · bilingual report structure | a report pair must share order and shape: the same number of `##` headings, byte-identical `artifacts` fences, and `sixstage` fences with all five sections and the same final conclusion |
| **Discipline X · report before self-delete (candidate, v46)** | a role that discovers a stray temp file must **not delete it immediately**; it reports the path first and waits until the controller has taken the pre-snapshot readings (`tmp` recursive count and `git status --porcelain -uall`), then runs `Remove-Item`; violations are handled under the **non-self-heal tier** (section 5.1) |

**GLM batch-auxiliary strategy**: `GLM-5.3-Flash (glm)` has ample capacity, so "at least twice" is only a floor; mechanical checks, long-document excerpts, bilingual mirror checks, dispatch drafts and log extraction should be delegated whenever possible; complex judgement and ranking are forbidden.

## 9. Checklist meta-rules

**Cite by the first sentence, not by "item N".**

| Rule (first sentence) | Content |
|---|---|
| "every 2 slices / a ruling change / a directive bump / a queue change / a user request / 2+ anchor mismatches ⇒ regenerate" | when to regenerate |
| "pending rulings must be listed one by one" | section 6 must name each one |
| "a queue entry must give six elements" | identifier / content / size tier / 7 budget / dependency / boundary |
| "a rule obligation is not a gate obligation" | rule text is not a mechanical criterion |
| "no physical-clock assertions" | observations trigger on events |
| "gates green before commit" | hard gate |
| "read agentName mechanically from the dispatch table" | read verbatim from `transcript.arguments.agentName`; state it honestly when unavailable |
| "governance slices never write back DEV_PLAN" | the authoritative ledger is appendix C.0 to C.0x |
| "a slice name never changes once it is in the ledger" | stable names |
| "when the authorised quantity conflicts with the policy, the policy wins" | neither overreach nor refusal follows from the quantity |
| "one fact, one place" | body text or a pointer, not both |
| "the controller may reword its own registrations; a reviewer's text is verbatim" | two different standards |
| "a factual error is independent of a design decision" | ledger and report disagreeing is a factual error |
| "the 10 receipt must carry the fact pack for the next checklist" | fields: HEAD, CI run, new ledger rows, new observations, new discipline candidates, max OB / deferral, max ADR / DR, tmp assertion, diff vs the previous version, open items |
| "[ROUTINE] small fixes take no new slice number but lose no discipline" | same 9 gates, CI and fence rows |
| "the start message is an artefact" | persist `00-briefing.md`, list it in MANIFEST and SHA256SUMS, declare it in the fence |
| "a carrier deviation is the top-grade independence downgrade" | register a D item with an independence impact assessment |
| "mechanical criteria beat witness statements" | zero-model / zero-side-effect claims must be falsifiable |
| "concurrency assertions must cover survival and release" | beyond the failure face, assert the winner survives and the slot is reusable after release |
| "layer boundaries must not be widened" | state explicitly what a layer does not support |
| "evaluation slices stop at 10 after 9" | never commit or push on their own |
| "a single output must not exceed 10000 characters" | batch before the limit; never self-restart |
| "compaction ends the slice" | at a 100K context or a compaction prompt, end the current slice |
| "the batch-auxiliary floor cannot be replaced by a controller self-check" | too few dispatches ⇒ top up in the same turn; "the script checked it" is not an exemption |

### Gate pitfalls (P48 six + P49 seven + P52 two + P55 four + P56 one + P57 zero)

> The items from P48 through P55 are carried over from v45 and the P56 entry from v46; **v47 adds no gate pitfall**.

20. **A subagent must report before it deletes its own stray files** (candidate discipline, evidenced by P56): a role created a temp directory outside `tmp/<slice-id>/` and deleted it immediately, so the "reading (2) baseline" of the section 5.1 self-heal tier was unavailable. Correct order: **self-report → the controller takes the pre-snapshot (`tmp` recursive count plus `git status --porcelain -uall`) → self-delete → zero-residue readback**; violations go to the non-self-heal tier.

## 10. One-line opening for a new session

Here is my project snapshot (v47); please check the 5 anchors in section 1 first and then continue with the queue in section 4, or trust the repository when they disagree.

[paste the full v47 checklist here]

[companion documents] `REMAINING_SLICES.md` (P57 assessed), `DELIVERY_DIRECTIVE.md` (v1.61), `docs/validation/p57-closeout-assessment{,_EN}.md` (P57 report, `7a43dfa`, runs in section 2).

Current slice: P57 is closed and layer 0 is dequeued; next is layer-1 item 8, still ⛔ blocked on the R-08 / R-09 rulings, or layer 2 `P7` once the user names it.

## 11. Current state in one line

P57 and the O-1/O-3 cleanup are closed (`7a43dfa` plus `3884351` pushed to origin, runs `41991695581` / `41992956993` green on both legs; P5 returned no new slice and P6 produced the eight-gate inventory, R=0, D=0); next is layer-1 item 8, pending the R-08 / R-09 rulings.
