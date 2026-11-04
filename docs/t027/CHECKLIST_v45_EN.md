# Task execution checklist (cross-session snapshot v45 · P55 closed)

Generated: 2026-11-05, HEAD = origin/main = `c76c6ce` (P55 backfill commit), directive v1.60. P55 is closed and pushed to origin (`63e2e24` main commit + `c76c6ce` backfill, neither pushed to gitee).
Supersedes: fully replaces v44 and any earlier snapshot. Backfill commits (docs-only) only move HEAD forward and do not change the anchor semantics in this section.

Changes in this version (vs v44): (1) P55 `EXIT-SEED-BUILD-EVIDENCE` closed (`63e2e24` main commit + `c76c6ce` backfill; runs `41925623071` and `41928327517` both green; cross-platform evidence package with 53 payloads plus a three-role two-generation build); (2) the §2 fact pack is now P55, with P54 moved to §3; (3) P55 appended to §3; (4) item 7 dequeued in §4, **item 8 `AT23-RESTRICTED-CANDIDATE` marked ⛔ blocked (depends on R-08 / R-09)**, and suspended plus pending-assessment items surfaced as alternatives; (5) §5 mechanical readings re-verified (max OB 136, deferral 25, ADR 30, DR 13); (6) R-02 / R-03 / R-05 / R-11 recorded as done in §6, while R-08 / R-09 / R-10 stay pending; (7) four new gate pitfalls for P55 in §9 (isolated module cache polluting `gofmt`, stopped-state wording tripping `G7-d`, captured `COMMAND` fidelity, cross-platform three-role structure).

## 1. Current anchors (check these five numbers first)

| Item | Value | Check command |
|---|---|---|
| HEAD | `c76c6ce` | `git rev-parse --short HEAD` |
| origin/main | `c76c6ce` | `git rev-parse --short origin/main` |
| Directive version | **v1.60** | `docs/DELIVERY_DIRECTIVE.md` header |
| SELFTEST | **697/697** | `node tools/gates/gate.js --selftest` |
| Three scopes | **all TOTAL_FAIL=0** | `gate.js --all --scope=index`, `tree`, `ci` |

If two or more items disagree, trust the repository and report the terminal readings first.

Fixed `tmp` assertion: top level = `.gitkeep` + `gates-cache`; `gates-cache` holds only `g4b`; P55's `tmp/EXIT-SEED-BUILD-EVIDENCE/` was cleared at 8c.

## 2. P55 fact pack

| Field | Value |
|---|---|
| Slice / size | `EXIT-SEED-BUILD-EVIDENCE` (evidence-chain supplement `[SLICE]`) / M |
| Goal | Complete the AT-13 evidence chain (seed and candidate kept separate; external oracle comparison; clean-room success and failure both leave the seed intact) with **no CI structure change and no product code change** |
| Change surface | `.md` and the evidence package only: package `docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/` (53 payloads + 2 metadata files) + report pair `docs/validation/exit-seed-build-evidence{,_EN}.md` + six paired authoritative documents + the P55 ledger entry; **zero `.go`, zero workflow** |
| Core output | (1) dual-platform capture (Windows `go1.27.0` / Linux VM `go1.22.12`); (2) **three-role two-generation structure**: seed `fc99f55…` builds `prfrail`, verifier `18c9340…` builds `prfrail-release`, candidate `baa144d…` builds `prfrail`; (3) `OracleHash=sha256:388f8884…` identical on four sides (both platforms + fixture + the `s1-selfhost.md` CI input table); (4) four counterexamples fail closed on both platforms; (5) seed tree hash stays `equal=true` after failure; (6) Linux builds ran with `GOPROXY=off` (mirror preload then offline) |
| Review chain | 3 two rounds (`prfrail-tester`): `TEST: DONE` (0C/0H/6M/6L) → six Mediums closed via `claim-scope.md` + `anomalies.md` → r2 `TEST: DONE` (6/6 closed, R2-1..R2-6 fixed) → 5 two rounds (`deep-reasoner` + `DeepSeek V4 Pro`): `PRE-REVIEW: PASS WITH FIXES` (1 Medium = abbreviated captured `COMMAND` → erratum E-1 plus faithful derived index command) → 6 (`independent-reviewer` + `MAI-Code-1.1-Flash`): `INDEPENDENT SCAN: PASS` (Section B `(no findings)`) → 7 (`independent-reviewer` + `GPT-5.3-Codex`): `FINAL REVIEW: PASS WITH FIXES` (1 Low wording, fixed) → 8a document sync (`prfrail-documenter`) → 9 three scopes + SELFTEST → 8c cleanup → 10 stop point → commit and push after user authorisation |
| Commit / CI | `63e2e24` (`test: add AT-13 seed-driven build evidence across Windows and Linux (P55)`, 69 entries) + `c76c6ce` (`docs: record P55 slice completion and CI evidence`, 4 entries); run `41925623071` (63e2e24) green, run `41928327517` (c76c6ce) green; the ubuntu leg includes Race and Contract fixtures; the three conditional jobs are skipped by design |
| Accounting | `R = 0`, `D = 1` (early start without a trigger word, corrected by a forced user interruption and converted in place); 9 `runSubagent` calls (3x2, 5x2, 6, 7, 8a, batch x2); 0 probes / paid calls |
| Rulings | The user approved A-4 (three-role structure replacing "run selfhost from the seed", because `fc99f55` has no release tool) and considered it stronger; approved A-5 (`goproxy.cn` preload + `GOPROXY=off` builds); ruled A-3 (controller running step 2) compliant; required the batch assistance to be dispatched (delivered 2/2, `CONSISTENT` + `PARITY: OK`) |
| Document sync | `TEST_STRATEGY` (AT-13 supplement registration + boundary pointer), `DEV_PLAN` (T017 progress), `OPERATIONS` / `S1_SUPPORT_MATRIX` / `S1_RELEASE_NOTES` (three "not applicable + reason" entries), `REMAINING_SLICES` (ledger write-back); CN/EN numstat symmetric pair by pair |
| Capability boundary | `claim-scope.md`: **S-1..S-7 hold / N-1..N-7 are not asserted** (the Windows verifier shares the candidate tree, so independence is not proven there; `SeedHash` hashes an assembled directory; the oracle comparison is structural; behavioural agreement is not a malware-free proof) |
| Open items | No hard open item; observations: Windows toolchain `go1.27.0` versus CI `1.22.12` (N-6), Windows binaries not hash-bound (F-08), counterexample commands carry placeholders (N-4/E-1), E-2 wording erratum |

## 3. Completed slices (recent baseline)

| Version | Slice | Main commit + backfill | Outcome |
|---|---|---|---|
| v1.60 | P49 | `75990a8` + backfill | Closed (R = 1, D = 2; run `41445657595`) |
| v1.60 | P50 | `81f9dc9` + backfill | Closed (R = 1, D = 1; run `41507722672`) |
| v1.60 | P51 | `2b0af2d` + backfill | Closed (R = 1, D = 2; run `41569353460`) |
| v1.60 | P52 | `c636dd8` + backfill | Closed (R = 0, D = 1; run `41642139932`; GAP-77 demonstrated) |
| v1.60 | P53 | `7802819` + backfill | Closed (R = 0, D = 0; run `41707160563`; GAP-77 fixed) |
| v1.60 | P54 | `7c0763a` + `f38e0ce` | Closed (R = 1, D = 1; run `41840557379`; C6 decision path plus downgrade ruling) |
| v1.60 | P55 | `63e2e24` + `c76c6ce` | Closed (R = 0, D = 1; runs `41925623071` and `41928327517` both green; AT-13 dual-platform evidence chain plus three-role two-generation build) |

## 4. Execution queue

General rule: layer 1 clears the record chain, layer 2 is the strategic mainline, layer 3 is the pre-release gate, layer 4 is the exit gate; the order inside a layer is fixed and layers are not skipped.

### Layer 1 (cleared)

Nothing left.

### Layer 2 (strategic mainline)

Nothing left (P53 dequeued; layer 2 is fully clear).

### Layer 3 / layer 4 (current mainline, currently blocked)

Item 7 (P55) is dequeued (`63e2e24` closed). **Next slice = item 8 `AT23-RESTRICTED-CANDIDATE`, currently ⛔ blocked (depends on the R-08 / R-09 ruling).**

| # | Slice | Scope | Size | 7 budget | Depends on | Boundary |
|---|---|---|---|---|---|---|
| 8 | AT23-RESTRICTED-CANDIDATE | AT-23 restricted real candidate and sign-off | M | Codex one final review | **R-08 / R-09 ruling (pending ⇒ ⛔ blocked)** | No sign-off without a real candidate authorisation |

### Alternatives (the only two categories that may move while blocked; both need a user call-out)

| Category | Item | Note | Start condition |
|---|---|---|---|
| Suspended (trigger-based) | `G4B-CORPUS-DECOUPLE` | Corpus versus decision decoupling; suspended pending a trigger | User call-out or trigger met |
| Pending assessment (governance) | P7 governance OB extraction summary | Pending-assessment item in the later-batch section; includes extracted OBs such as OB-77 | Start an assessment action per §1.8.1 first |
| Pending assessment (product) | P5 `T027-B` / P6 `S1-EXIT-DOC` close-out list | Product-side close-out; only three fields exist, so the full §1.2 table must be filled first | Also an §1.8.1 assessment action |
| Low-priority governance | Five low-priority tasks | Not started without a trigger; require an explicit user call-out | User call-out |

> Calibre: none of the four categories above is the layer-3 mainline and none may be read as "layer 3 resumed"; starting them does not lift the `AT23-RESTRICTED-CANDIDATE` block.

## 5. Observations

- Observing (13): `OB-1 / 7 / 10 / 33 / 87 / 99 / 100 / 104 / 105 / 107 / 111 / 113 / 135`.
- Pending discussion (11): `OB-19 / 21 / 52 / 53 / 55 / 57 / 58 / 59 / 61 / 75 / 120`.
- To do (1): `OB-69`. Non-terminal total **25**.
- GAP register: **GAP-77 fixed (P53 closed)**; GAP-78/79/80 candidates (registered in the P52 report §2b and the ledger).
- New P55 observations (no OB number, slice-level only): (1) the Windows toolchain is not the CI one (N-6); (2) captured `COMMAND` fidelity issues (E-1/E-2, corrected by errata); (3) `index.json` once carried stale derived values (rebuilt to the v2 calibre).
- Mechanical readings: max OB in §12.4 = **OB-136**; max deferral = 25th; ADR max = **30**; DR max = **13**.

## 6. Awaiting user rulings (R-01 to R-11)

| ID | Item | Status |
|---|---|---|
| R-01 | P46 strategic review立项 | Done |
| R-02 | `CONTRACTS-S7-DURABILITY` pace and ownership | Done (P54 closed) |
| R-03 | P45 successor ownership | Done (`RACE-LEG-EXPAND` closed) |
| R-04 | DR-13 "fixed + observed" calibre | Pending |
| R-05 | Whether DR13-FIX-D3 is a separate slice | Done (separate slice, closed) |
| R-06 | OB-55 ownership and T027-B pace | Pending |
| R-07 | MVA, S1 removals and S2 commitment narrowing | Pending |
| R-08 | Candidate, Windows premise P and probe authorisation | **Pending (blocks item 8 `AT23-RESTRICTED-CANDIDATE`)** |
| R-09 | AT-24 `unknown` acceptance boundary | **Pending (blocks the item-8 sign-off calibre)** |
| R-10 | Total budget and per-slice model ceiling | Pending (blocks all scheduling) |
| R-11 | Historical correction scope | Done (B tier plus per-item differentiation) |

## 7. Start-up and step-0 rules

1. Verify anchors, then step 0: register the slice ledger (`R/D = 0` baseline, temp location `tmp/<sliceID>/`), write `input-list.md` and `freedom-list.md`.
2. Every fix slice must ship a regression test and must record unexecuted tests/lints explicitly.
3. The three gates at 9 plus `SELFTEST 697/697` plus green `go build/vet/test ./...` are required before committing.
4. Model binding (hard): 1 and 5 use `deep-reasoner` + `DeepSeek V4 Pro (deepseek)`; 2, 3 and 8 use `DeepSeek V4.1 Flash (deepseek)`; 6 uses `independent-reviewer` + `MAI-Code-1.1-Flash (copilot)`; 7 uses `independent-reviewer` + `GPT-5.3-Codex (copilot)`; strategic review round 1/3 uses `independent-reviewer` + `GPT-6.1 Sol (copilot)`; batch assistance uses `deep-reasoner` + `GLM-5.3-Flash (glm)`. **Every `runSubagent` must pass both `agentName` and `model`.**
5. **Pacing discipline**: do not stop routinely — green step 4 goes straight to 5; a Medium+ finding at 5/6/7 means fix then re-run that ring in place; **10 is the only routine stop point**.
6. Cross-version/cross-platform triage: check toolchain version and build-tag coverage first; after cross builds re-check `go env GOOS GOARCH`; before deleting "unused here" symbols run at least one `go vet` for the opposite `GOOS`.
7. Batch assistance: dispatch `deep-reasoner` + `GLM-5.3-Flash (glm)` **at least twice** per slice (policy and output ceiling in §8); **strictly serial**, no parallel subagents.
8. Controller model policy: default `DeepSeek V4.1 Flash`; switch to Pro for low-level concurrency/state-machine/handle slices, design documents over 10,000 characters, heavy CI log root-cause work, or two consecutive behavioural drifts; **any downgrade must be registered per the OB-37 meta-rule**.
9. **Carriers are not interchangeable**: step 2 must go to `prfrail-implementer` and step 3 to `prfrail-tester`; the controller only orchestrates, writes artefacts, runs gates and reports. **Exception**: pure evidence-supplement / log-capture slices may be run by the controller when the user says so in the same round (P55 precedent: start message [6]).
10. **The start message is an artefact**: step 0 must persist the verbatim start message as `docs/validation/evidence/<sliceID>/00-briefing.md`, register it in `MANIFEST.md` and `SHA256SUMS.txt`, and declare it in the `artifacts` fence; without it step 1 must not start.
11. **Linearisation discipline**: for persisted state machines, "persisted" must be the linearisation point; an in-memory CAS is only a same-process fast reject.
12. **Evidence package and fence sync**: every new payload requires a `MANIFEST.md` row, a recomputed `SHA256SUMS.txt`, and an `artifacts` fence line in the same round.
13. **Reviewer archival discipline**: 5/6/7 outputs must be archived verbatim (P55 precedent: `04-reviews/`); the controller must not reconstruct them from bullet points.
14. **Assembly-verification slice discipline**: assertions must land on falsifiable mechanical criteria; constant-identity claims must declare themselves non-sentinel; the report must state whether production code changed.
15. **Assessment-type slices must stop at 10 after 9**: an assessment slice has no implementation object; stop and wait for the same-round authorisation.
16. **A single output must not exceed 10,000 characters**: split before the limit; never self-restart; read on demand.
17. **Compaction ends the slice**: at 100K of context or a compaction prompt, end the current slice immediately; one L-size slice per session.
18. **The start message must contain a "document sync" section**: if the slice changes stage state or contract capability, step 8 must sync the authoritative documents; omitting it is a structural defect.
19. **Environment pre-check and reuse**: declare external environments at step 0; before installing, probe PATH, config and `find`; **remote SSH must use `bash -lc` and probe user-level directories**; pin versions to match CI and verify checksums. P55 reuse note: on VM `10.0.0.199`, `~/go1.22/bin/go` is go1.22.12 (the login shell defaults to a newer 1.27.1, so pin it explicitly).
20. **Capture replayability and fidelity (new in v45)**: a capture's `COMMAND:` line must match the real argv; if only a display abbreviation is written, add an erratum in `anomalies.md` and write the faithful command plus a `commandErrata` pointer into the derived `index.json`; drift or synthetic inputs must have replayable rebuild steps in `replay.md`.

## 8. Cross-slice collaboration discipline (must be in the start message)

| Element | Content |
|---|---|
| Carrier table | 1, 5, 6, 7 use `deep-reasoner` / `independent-reviewer` (read-only); 2, 3, 8 use `prfrail-implementer` / `prfrail-tester` / `prfrail-documenter` (writable); `prfrail-tester` must not run review-layer steps; batch assistance must be `deep-reasoner`, never `quick-researcher` |
| Dispatch table mechanical read | every `agentName` in the 10 receipt must be read verbatim from `transcript.arguments.agentName`, with the script stdout attached |
| Dispatch serialisation | `runSubagent` is serial and blocking; never call several subagents in parallel; always pass `agentName` and `model` |
| Commit messages | always English, `<type>: <summary>`; use `git commit -F <UTF-8 file>` when any non-ASCII character is present |
| Independently checkable vs not | enumerations, fact comparisons and citation lists do not become void on model error; judgments, conclusions and review opinions must be redone |
| 6 before 7 | 6 must complete before 7; if 6 is re-dispatched and documents changed, 7 must be re-evaluated |
| Temp location and cleanup | always `tmp/<sliceID>/`; only 8c may clean, only this slice's directory, never a parent recursively; **character counting and intermediate records must use inline scripts, never files on disk (especially `$env:TEMP`)** |
| Long-command redirection | never PowerShell `>`; use `\| Tee-Object -FilePath <path>` or `cmd /c "..."` |
| Gate execution | use `--all --skip-base` mid-slice; everything must be green before committing |
| Discipline E · index scope order | exact `git add` first, then `--scope=index`; the receipt must report `G6-1 scanned N` with N > 0 |
| Discipline F · spec over example | follow the spec when the start example conflicts, and register the conflict |
| Discipline G · fence equals landed snapshot | `artifacts` fence sorted by path, TAB separated, byte-identical CN/EN; new files need new lines |
| Discipline H · opposite GOOS check | run at least one opposite-`GOOS` `go vet` before removing or changing constants/functions |
| **Discipline O · output length defence** | no single output above 10,000 characters; split before the limit; never self-restart |
| **Discipline P · controller model switch trigger** | see §7 item 8; register downgrades per OB-37 |
| **Discipline Q · batch assistance output ceiling** | state the character ceiling in the dispatch; report the measured length |
| **Discipline R · stale-ledger one-shot clearing** | a ledger that disagrees with reality gets one clearing slice with code and commits as the source of truth |
| **Discipline S · session length and compaction** | at 100K of context or a compaction prompt, end the slice immediately |
| **Discipline T · quantifier form** | never write "identifier immediately followed by a quantifier"; `G6-4` reads it as a line-count literal — write "that entry" or "that record" |
| **Discipline U · no verbatim trigger strings** | reports must not quote the triggering criterion string |
| **Discipline V · split big tasks, avoid half-done pollution** | dispatch big tasks in small serial steps; after a timeout check `git status`/`git diff` for partial artefacts and revert precisely when needed |
| **Discipline W · bilingual report structural parity (new in v45)** | the report pair must match: same `##` heading count, byte-identical `artifacts` fence, complete five-section `sixstage` fence with the same conclusion line; check it with a mechanical script every slice (P55 precedent: batch assistance `PARITY: OK`) |

**GLM batch assistance policy**: `GLM-5.3-Flash (glm)` is plentiful and "at least twice" is only a floor — mechanical cross-checks, long-document extraction, bilingual mirror checks, dispatch drafts and log extraction should be outsourced; no complex judgment or ranking. **P55 lesson: a controller-side script does not replace the dispatch floor** — an under-dispatch must be topped up in the same round and evidenced.

## 9. Checklist element rules

**Cite rules by their first sentence, never by "item N".**

| Rule (first sentence) | Content |
|---|---|
| "Every two slices / ruling change / directive bump / queue change / user request / two or more anchor mismatches ⇒ regenerate" | regeneration triggers |
| "Pending rulings must be named one by one" | §6 must name each item |
| "Queue items must give six elements" | identifier / scope / size tier / 7 budget / dependency / boundary |
| "Rule obligation is not gate obligation" | rule text is not a mechanical criterion |
| "No physical-clock assertions" | observations trigger on events, not on "N months later" |
| "Green gates before commit" | hard gate |
| "Dispatch table agentName mechanical read" | read verbatim from `transcript.arguments.agentName` |
| "Governance slices do not write back DEV_PLAN" | authoritative ledger = appendix C.0–C.0x |
| "A slice name never changes once it enters the ledger" | name stability |
| "Policy wins over authorisation count" | do not treat this as overreach or refuse to execute |
| "One fact is described in exactly one place" | body or pointer, pick one |
| "The controller may reword its own records; reviewer text must stay verbatim" | different standards for the two sides |
| "Factual errors are independent of design decisions" | ledger/report contradiction is a factual error |
| "The 10 receipt must attach the fact pack for the next checklist" | HEAD, CI runs, new ledger rows, new observations, new discipline candidates, max OB / deferral, max ADR/DR, tmp assertion, diff versus previous version, open items |
| "[ROUTINE] small fixes take no new slice number but lose no discipline" | same 9 gate + CI + fence line |
| "The start message is an artefact" | `00-briefing.md` persisted + MANIFEST/SHA256SUMS + fence declaration |
| "Carrier deviation is the highest-grade independence downgrade" | register a D item with an independent impact assessment |
| "Mechanical criteria beat testimonial statements" | zero-model/zero-side-effect claims must be falsifiable |
| "Concurrency assertions must cover survival and release" | besides the failure face, assert the winner survives and the slot can be re-used after release |
| "Layer boundaries may not be read broadly" | an unsupported capability must be declared explicitly |
| "Assessment slices stop at 10 after 9" | no self-authorised commit/push |
| "A single output must not exceed 10,000 characters" | split before the limit; never self-restart |
| "Compaction ends the slice" | at 100K or a compaction prompt, end the slice immediately |
| "The batch-assistance floor cannot be replaced by controller self-checks (v45)" | an under-dispatch must be topped up in the same round; "the script checked it" is not an exemption |

### Gate pitfalls (P48 six + P49 seven + P52 two + P55 four)

### From P48 (carried over)

1. **`.go` copies under `tmp/` are compiled by `go build ./...`**: review-input source copies must be `*.go.txt`.
2. **`G7-b` disposition closed set**: uncommitted temp files are `local-only`; writing `local` reports an unknown disposition.
3. **A report containing an `INDEPENDENT SCAN` marker must carry a `sixstage` fence**: `## Section A` to `## Section E` exactly once each, the conclusion line being the last non-empty line in the block, with merge differences disclosed outside the fence.
4. **`.md` writes must add the BOM explicitly**: a node rewrite without the BOM trips `G4a` immediately.
5. **`G8-b` requires MANIFEST to list payloads one by one**: a directory row is rejected as "payload not listed".
6. **Backfill readings last**: closing a slice changes `changed=` and `scanned=` readings.

### From P49 (carried over)

7. **Changing `.github/workflows/*.yml` requires syncing the `internal/release` workflow freeze contract**.
8. **Remote SSH pre-check**: non-interactive sessions do not load the user PATH; use `bash -lc` and probe user-level directories such as `~/.local/bin`.
9. **PowerShell to ssh quoting**: complex cross-shell commands must be uploaded as a script.
10. **6 output must be A–E shaped**: `Section A` (scope), `B` (findings or `(no findings)`), `C` (evidence), `D` (falsifiability audit), `E` (unverified parts plus overturn conditions), then the conclusion line.
11. **Fence and disk enumeration must agree**: new payloads need a fence line, a MANIFEST row and a recomputed SHA256SUMS in the same round.
12. **`gofmt -l .` walks hidden directories and `tmp/`**: see P55 item 16.
13. **Recompute hashes with a read-only script**: a rewriting normaliser fails with `EBUSY` when an editor holds the payload.

### From P52 (carried over)

14. **Section heading status glyph must match the status row**: `G1(d)` compares both closed sets.
15. **Recompute SHA256SUMS after any line-number correction**: content fixes change payload hashes.

### From P55 (new)

16. **Isolated `GOCACHE`/`GOMODCACHE` under `tmp/<sliceID>/` pollutes `gofmt -l .`**: the module cache contains third-party `.go` sources (e.g. `gopkg.in/yaml.v3`) and `base gofmt` goes red. Fix: delete the whole directory after capture **and clear the session-level `$env:GOCACHE`/`$env:GOMODCACHE`**, otherwise a later `go test` recreates it.
17. **Writing "stopped" in a report trips the `G7-d` stopped-state criterion**: use "the 10 stop point is reached / Reached the 10 stop point"; only a genuinely stopped slice writes "stopped — incomplete".
18. **A capture's `COMMAND:` line must match the real argv**: a display abbreviation (missing `--noop-chain` and similar flags) makes the capture unreplayable as recorded — a fidelity defect that 5 will raise as Medium. Fix: add an E-N erratum in `anomalies.md` and write the faithful command plus `commandErrata` into the derived `index.json`.
19. **Remote Go module acquisition**: the VM may blackhole `proxy.golang.org` / `sum.golang.org` over IPv6 and `GOPROXY=direct` may fail on vanity resolution ⇒ preload from a reachable mirror then set `GOPROXY=off` for builds (fixes download and proves offline buildability); the module cache is read-only, so `chmod -R u+w` before deleting.

## 10. One-line opening for a new session

Here is my project state snapshot (v45); check the five anchors in §1 first, then continue with the §4 queue if they match, otherwise trust the repository.

[Paste the full v45 checklist]

[Supporting documents] `REMAINING_SLICES.md` (P55 done), `DELIVERY_DIRECTIVE.md` (v1.60), `docs/validation/exit-seed-build-evidence{,_EN}.md` (P55 report, `63e2e24`, runs in §2).

Current slice: P55 closed; next slice = layer-3 item 8 `AT23-RESTRICTED-CANDIDATE`, **⛔ blocked (waiting for the R-08 / R-09 ruling)**.

## 11. Current state in one sentence

P55 `EXIT-SEED-BUILD-EVIDENCE` is closed (`63e2e24` + `c76c6ce` pushed to origin, runs `41925623071` and `41928327517` both green; **AT-13 dual-platform evidence chain plus a three-role two-generation build delivered** — seed `fc99f55` / verifier `18c9340` / candidate `baa144d`, `OracleHash` identical on four sides, four counterexamples fail closed on both platforms, failure leaves the seed intact, R=0, D=1); next step = layer-3 item 8 `AT23-RESTRICTED-CANDIDATE`, **currently ⛔ blocked pending the R-08 / R-09 ruling**.
