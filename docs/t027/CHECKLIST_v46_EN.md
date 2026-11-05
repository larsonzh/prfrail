# Task execution checklist (cross-session snapshot v46 · P56 closure)

Generated: 2026-11-05, HEAD = origin/main = `504caff` (the P56 main commit), directive v1.61. P56 is closed and pushed to origin (`504caff`; run `41978355255` green on both legs; never pushed to gitee).
Supersedes v45 and every earlier snapshot. A docs-only backfill commit only moves HEAD; it does not change the anchor semantics of this section.

Changes in this version (vs v45): (1) P56 `G4B-CORPUS-DECOUPLE` closed (`504caff`; run `41978355255` green on both legs; the G4b corpus enumeration decoupled from the index, cache bumped to `v2`, new T114-T118); (2) section 2 now carries the P56 fact pack and P55 moved into section 3; (3) section 3 gained the P56 row; (4) **section 4 was re-ordered by the user's ruling of 2026-11-05** - the layer order changed from "layer 1 / layer 2 / layer 3+4" to four segments: front-end stocktaking -> main hard gate -> routine governance backlog -> low-priority alternatives, with the rationale in the first paragraph of section 4; (5) section 5 re-reads the mechanical counters (OB-120 is resolved, so the pending counts changed); (6) section 6 keeps R-01 to R-11 unchanged; (7) section 9 gains one **candidate discipline** (a subagent must report before it deletes its own stray files).

## 1. Current anchors (check these 5 numbers first)

| Item | Value | Command |
|---|---|---|
| HEAD | `504caff` | `git rev-parse --short HEAD` |
| origin/main | `504caff` | `git rev-parse --short origin/main` |
| Directive version | **v1.61** | header of `docs/DELIVERY_DIRECTIVE.md` |
| SELFTEST | **703/703** | `node tools/gates/gate.js --selftest` |
| Three scopes | **all TOTAL_FAIL=0** | `gate.js --all --scope=index`, `tree`, `ci` |

If two or more items disagree, trust the repository and report the final readings before continuing.

Fixed `tmp` assertion: the top level is `.gitkeep` plus `gates-cache`; `gates-cache` contains only `g4b`; P56's `tmp/G4B-CORPUS-DECOUPLE/` was removed at 8c.

## 2. P56 fact pack

| Field | Value |
|---|---|
| Slice / size | `G4B-CORPUS-DECOUPLE` (criterion-optimisation class `[SLICE]`) / M |
| Goal | Close OB-120: decouple the G4b corpus enumeration from the index so a cache hit and a cold run agree |
| Change surface | `tools/gates/lib/ctx.js`, `tools/gates/lib/criteria/g4b.js`, `tools/gates/selftest.js`, `tools/gates/cjk-newwords.txt` (one character registered); directive v1.61; ledger; report pair; twelve evidence-package entries |
| Core output | (1) both branches unified on `git ls-tree -r -z --name-only <ref>` (tree/index use `HEAD`); (2) `CACHE_VERSION` and the stored `v` field moved from `v1` to `v2`; (3) new T114-T118 (enumeration source in tree / range / index, cold-hot equivalence under a staged deletion, version consistency); (4) A/B evidence: the old enumeration reported `corpus_old=1428` with one unknown, the new one reports `corpus_new=1429` with an empty difference set |
| Review chain | 5 twice (`PRE-REVIEW: PASS WITH FIXES`) -> 6 `INDEPENDENT SCAN: PASS` -> 7 three times (`FINAL REVIEW: FINDINGS` -> `RE-REVIEW: FINDINGS` -> `RE-REVIEW: PASS`) |
| Commit / CI | `504caff` (22 entries); run `41978355255` green on both legs (Go ubuntu-latest / Go windows-latest both success; the three conditional jobs skipped by design) |
| Metering | `R = 2`, `D = 1`; `runSubagent` 15 times (2 twice, 3 three times, 5 twice, 6 once, 7 three times, 8 twice, batch auxiliary twice); probes / paid external calls 0 |
| Rulings | the user accepted `D = 1` under the **self-heal tier** (the subagent's stray files lay outside `tmp/<slice-id>/` and the controller verified zero residue afterwards); the missing mechanical dispatch-table reading (the debug log has no transcript) is **accepted as a known fact** |
| Doc sync | directive v1.61 (a section 0.3 row, the section 12.4 OB-120 cell flipped to "Resolved", appendix C.0ar; the **section 6.3 G4b criterion text is unchanged** - the implementation caught up with it); ledger flipped to ✅; new report pair; `DEV_PLAN` and `ARCHITECTURE` declared "not applicable, with reason" |
| Capability boundary | zero change to the criterion body; no dedicated `ci`-scope assertion (it shares the implementation line with `range`, covered by T117); "a broken whitelist fails closed before the cache" remains a pre-existing coverage gap |
| Not closed | no hard open item; three observations: the `ci` assertion gap, the broken-whitelist assertion gap, the dispatch-table mechanical-reading gap |

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

## 4. Execution queue

**Basis for this ordering (user ruling, 2026-11-05)**: first start the **stocktaking evaluation actions** (P5 / P6) so that section 1.8.1 establishes the full inventory and finish line of S1 Exit and the remaining workload; then **concentrate resources on the main hard gate** (the real candidate and the `unknown` boundary); afterwards clear technical debt that is not on the critical path; low-priority items start only when the user names them explicitly. The order inside a layer is fixed and layers are not skipped.

### Layer 0 Front-end stocktaking (highest priority; evaluation actions first)

| Order | Item | Content | Size | Depends on | Boundary |
|---|---|---|---|---|---|
| 0-1 | P5 `T027-B` close-out list | Start an **evaluation action** per section 1.8.1: produce "does it hold / a standalone slice or merge / if a slice, fill every section 1.2 field" | pending-evaluation item | none (evaluable immediately) | evaluation only, no product code change; **not started in this version** |
| 0-2 | P6 `S1-EXIT-DOC` close-out list | Same; the goal is the full inventory and the finish line of S1 Exit | pending-evaluation item | none (may run serially after 0-1) | same; **not started in this version** |

> Rule: both are "pending-evaluation items", which section 1.8.1 forbids starting directly; the conclusion may legitimately be "no slice", which is itself a valid output.

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

> Layer boundary: layer 0 is an **evaluation action** and produces no product change; layer 1's blocker is not lifted by starting layer 0; neither layer 2 nor layer 3 may be read as "the main line resumed".

## 5. Observations

- Observing (13): `OB-1 / 7 / 10 / 33 / 87 / 99 / 100 / 104 / 105 / 107 / 111 / 113 / 135`.
- To be discussed (10): `OB-19 / 21 / 52 / 53 / 55 / 57 / 58 / 59 / 61 / 75`. **OB-120 was resolved by P56** (`Resolved`), which is why the count dropped from 11 in v45 to 10.
- To do (1): `OB-69`. Non-terminal total **24**.
- GAP register: **GAP-77 fixed (P53 closed)**; GAP-78/79/80 are candidates (registered in the P52 report section 2b and the ledger).
- New P56 observations (slice-level, no OB number): (1) the `ci`-scope enumeration assertion gap; (2) the "broken whitelist fails closed before the cache" assertion gap; (3) the dispatch-table mechanical-reading gap (this session's debug log holds only `session_start`).
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

### Gate pitfalls (P48 six + P49 seven + P52 two + P55 four + P56 one)

> The items from P48 through P55 are carried over unchanged from v45; only the following entry is new in v46.

20. **A subagent must report before it deletes its own stray files** (candidate discipline, evidenced by P56): a role created a temp directory outside `tmp/<slice-id>/` and deleted it immediately, so the "reading (2) baseline" of the section 5.1 self-heal tier was unavailable. Correct order: **self-report → the controller takes the pre-snapshot (`tmp` recursive count plus `git status --porcelain -uall`) → self-delete → zero-residue readback**; violations go to the non-self-heal tier.

## 10. One-line opening for a new session

Here is my project snapshot (v46); please check the 5 anchors in section 1 first and then continue with the queue in section 4, or trust the repository when they disagree.

[paste the full v46 checklist here]

[companion documents] `REMAINING_SLICES.md` (P56 done), `DELIVERY_DIRECTIVE.md` (v1.61), `docs/validation/g4b-corpus-decouple{,_EN}.md` (P56 report, `504caff`, run in section 2).

Current slice: P56 is closed; next is layer 0 front-end stocktaking (the **evaluation actions** for P5 / P6, started only when the user names them); layer-1 item 8 stays ⛔ blocked on the R-08 / R-09 rulings.

## 11. Current state in one line

P56 `G4B-CORPUS-DECOUPLE` is closed (`504caff` pushed to origin, run `41978355255` green on both legs; **the G4b corpus enumeration moved from the index to the object store**, the cache bumped to `v2`, `SELFTEST` 697 to 703, R=2, D=1); next is layer 0 front-end stocktaking (the P5 / P6 evaluation actions, started after the user names them).
