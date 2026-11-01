# ProofRail P47 governance trusted-baseline verification report: seven breakpoint-B checks and historical backfill (directive v1.60)

Slice `GOVERNANCE-BASELINE-VERIFY` (assessment/governance `[SLICE]`; pipeline trimmed to 0 → 1 → 2 → 8 → 9 → 8c → 10, steps 3–7 skipped); kickoff anchor `HEAD` = `0e46bd0`; directive v1.60; report date 2026-11-02; R / D baseline = 0 / 0; measured R = 2 for this slice (two reruns after gate reds, with multiple criteria in one run merged into one) and D = 0.

## 1 Executive summary

Seven binary verdicts: **three closed after backfill** (G4 / G5 / G7), **three still open** (P3 / P5 / S7), and **M6 follows the G4 reading** and is closed after backfill.

Backfill landed: four rows OB-133 to OB-136 added to `§12.4`; backfill pointers appended to the evidence cells of the P32 / P33 / P34 rows in `REMAINING_SLICES`; this slice's own ledger row and detail block added. **Zero edits to historical report bodies or to already-completed rows.**

A boundary that must be read together with the verdicts: the three "closed after backfill" verdicts mean **the record-chain gap is now registered**; they **do not revoke** the closure of P32 / P33 / P34 and **do not reopen** their acceptance.

## 2 R-11 ruling and correction scope

The user ruled **option B plus per-item differentiation**: G4 registers the current value and the corrected first threshold-breach time without recomputing historical scores; G5 adds an OB row stating that the independent registration does not hold, without rewriting the P33 body; G7 registers "the exception did not hold", while whether the exclusion should be relaxed belongs to a separate slice. Four red lines: add only, never edit history; keep CN and `_EN` in step; keep `G6-1` tense consistent; stop at step 10 once step 9 is green.

## 3 Seven breakpoint-B verdicts

| Item | Verdict | Mechanical basis | Gap (open items) |
|---|---|---|---|
| P3 | Still open | `REQ-023` in `PRODUCT_REQUIREMENTS` says an unknown external side effect blocks the chain, with no carve-out; `AT-24` in `TEST_STRATEGY` introduces a "mark `unknown` when unobservable in principle" carve-out | The two texts do not cross-reference; the authority for the carve-out is undecided (R-09) |
| P5 | Still open | In CI the candidate is built directly by the runner (`go build`, step `Build candidate without executing it`); the seed binary only takes part in `selfhost` verification and the candidate comes from `download-artifact` | RFC R19 requires the seed to **build** candidate N+1; that step has no evidence, and the backfill belongs to layer 4 |
| G4 | Closed after backfill | Per-event rework count is enumerable and the recount table is in the evidence package; the P32 cumulative formula dropped the first two rings | — |
| G5 | Closed after backfill | `§1.7.4` ② requires "OB or ADR" registration and forbids substituting a `§0.3` entry or an appendix-C note; P33 wrote "the record is this section" | — |
| G7 | Closed after backfill | The `§5.1` exception expressly excludes the "legal location, order violation only" case, and the P34 D event falls on that excluded side | — |
| M6 | Closed after backfill | Same source as G4: under cumulative counting the `§9.2` threshold is already met after the second native-verification anomaly | — |
| S7 | Still open | The OB-21 row yields no hit for owner, stop condition or review trigger | Three elements absent; carried by the new OB-136 |

Search scope and per-item criteria are in `01-baseline/verification.md`; the G4 / G5 / G7 correction text is in `01-baseline/correction.md`.

## 4 Backfill landed

| Carrier | Content | Preserved |
|---|---|---|
| `§12.4` OB-133 | P32 rework-count current value and first breach time corrected (G4) | P32 report not rewritten |
| `§12.4` OB-134 | P33 independent-registration gap registered (G5) | P33 closure not revoked |
| `§12.4` OB-135 | P34 exception applicability corrected (G7, observing) | `§5.1` body unchanged |
| `§12.4` `OB-136` | OB-21 three elements supplied (S7) | OB-21 row neither edited nor re-statused |
| `REMAINING_SLICES` P32 / P33 / P34 | Backfill pointers appended to the evidence cells | Status, dependency, size, completion date, evidence path |

Pointers were appended **inside** the evidence cell rather than as a new cell, so the status table keeps its column count. Details are in `02-backfill/registration.md`.

## 5 Pipeline trim and stage statement

- Step 0: landed `00-briefing.md` (verbatim kickoff message), `input-list.md`, `freedom-list.md`, `MANIFEST.md`, `SHA256SUMS.txt`.
- Step 1 baseline check: read-only retrieval plus mechanical recount, producing the two payloads under `01-baseline/`.
- Step 2 backfill: new rows and pointers only, producing `02-backfill/registration.md`.
- Step 8 report assembly: this file and its `_EN` mirror.
- Step 9 three scopes and SELFTEST; step 8c cleanup; step 10 stop point awaiting user authorisation.
- **Steps 3–7 skipped**: this slice has no implementation object; step 4 acts on a code diff, step 6 scans an implementation diff, and step 7 presupposes an implementation review object. **Independence boundary stated honestly**: the verdicts rest on control-seat checking plus read-only batch assistance and **were not reviewed by steps 6 / 7**.
- ⑥ skipped: this slice has no implementation diff, so step 6 has no scan object; the verdicts rest on control-seat checking and read-only batch assistance, and the independence boundary is stated above.

## 6 Models and accounting

| Dispatch | Carrier | Model | Output characters | Dispatch cap |
|---|---|---|---|---|
| 1 ledger enumeration | `deep-reasoner` | `GLM-5.3-Flash (glm)` | 2,972 | 3,000 |
| 2 reference extraction | `deep-reasoner` | `GLM-5.3-Flash (glm)` | 3,480 | 3,500 |

The two dispatches were **strictly serial** (the second was sent only after the first returned), honouring the `R1.1` / `R5.1` rule against parallel subagents; both were **read-only** with no writes. The control-seat checks of G4 / G5 / G7 and of the P5 workflow facts were performed by the control seat itself, read-only, with no writing subagent.

## 7 Gate readings

| Item | Reading |
|---|---|
| `gate.js --all --scope=index / tree / ci` | `TOTAL_FAIL=0` on all three scopes |
| `node tools/gates/gate.js --selftest` | `SELFTEST: PASS 697/697` |
| `G3(a)` bilingual symmetry | Directive pair +4 / −0 equals mirror +4 / −0; ledger pair +25 / −3 equals mirror +25 / −3 |
| `G6-1` status tense | No pending / completion conflict on added lines |
| `G6-4` expiring literal | One red on the first run (a quantifier after the identifier inside a ledger entry was read as a line-count numeric), green after rewording; a process fact of this slice |
| `G6-5` / `G6-6` | Section references resolve; every added ledger row opens with a closed-set status word |
| `G7` report contract | This file and its mirror are inside the G7 scan domain with complete fences |
| `G8` evidence tree | Package self-consistent (`MANIFEST.md` entries and `SHA256SUMS.txt` hash rows in step) |
| Six shapes (CN / `_EN`) | Lines 113 / 113; headings 9 / 9; bold 32 / 32; bars 150 / 150; blank 28 / 28; backticks 184 / 184 |

## 8 Commit and CI

This slice parked at step 10 per the "assessment slices stop at step 10 once step 9 is green" rule, then committed and pushed under a same-turn user authorisation (only `origin/main`, no gitee push).

| Revision | Commit | run | Result |
|---|---|---|---|
| Main commit (all payloads of this slice) | `bbd1745` | `41358624169` | Both green: `ubuntu-latest` and `windows-latest` both success; both legs report `SELFTEST: PASS 697/697` and `TOTAL_FAIL=0` |
| This section backfill (docs-only append) | registered at the append commit | registered at the append commit | registered at the append commit |

Package-level readings (final CI evidence): on `ubuntu-latest`, `internal/adapters` 11.959s, `internal/applier` 0.492s, `internal/chain` 0.565s, `internal/console` 1.912s, `internal/evidence` 0.009s, `internal/gates` 0.126s; on `windows-latest`, `internal/adapters` 29.364s, `internal/applier` 5.253s, `internal/chain` 3.912s, `internal/console` 15.398s, `internal/evidence` 0.057s, `internal/gates` 0.526s, `internal/guard` 2.655s, `internal/release` 3.818s, `internal/repair` 0.028s, `internal/snapshot` 1.320s, `internal/taskdef` 0.062s, `internal/tickets` 4.359s. Both legs report `SELFTEST: PASS 697/697` and `TOTAL_FAIL=0`; this backfill is a docs-only append and changes neither code nor criteria.

## 9 Known limitations

- The P5 criterion contains one interpretive fork: RFC R19 says the seed builds the candidate, while CI currently has the seed verify the candidate and the runner build it. This slice reads the text literally and calls the evidence missing; if the user holds that "seed verification" already satisfies R19, P5 can be re-judged. **The fork awaits a ruling and this slice does not relax it unilaterally.**
- The three "closed after backfill" verdicts close record-chain gaps only; they are not a review of historical scores or technical conclusions.
- Steps 6 / 7 were not dispatched, so the verdicts were not reviewed by an independent review layer.
- This slice does not change the `§5.1` exclusion; "the exception did not hold" is a factual registration only.
- Added ledger rows carry no numerals or line anchors; all recount detail lives in the evidence tree so that the ledger never holds an expiring literal.
- Of the six shapes all six match exactly (lines / headings / bold / bars / blank lines / backticks), so the structural symmetry criterion passes.

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/t027/CHECKLIST_v37.md	repo
docs/t027/CHECKLIST_v37_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/evidence/GOVERNANCE-BASELINE-VERIFY/00-briefing.md	repo
docs/validation/evidence/GOVERNANCE-BASELINE-VERIFY/01-baseline/correction.md	repo
docs/validation/evidence/GOVERNANCE-BASELINE-VERIFY/01-baseline/verification.md	repo
docs/validation/evidence/GOVERNANCE-BASELINE-VERIFY/02-backfill/registration.md	repo
docs/validation/evidence/GOVERNANCE-BASELINE-VERIFY/MANIFEST.md	repo
docs/validation/evidence/GOVERNANCE-BASELINE-VERIFY/SHA256SUMS.txt	repo
docs/validation/evidence/GOVERNANCE-BASELINE-VERIFY/freedom-list.md	repo
docs/validation/evidence/GOVERNANCE-BASELINE-VERIFY/input-list.md	repo
docs/validation/governance-baseline-verify.md	repo
docs/validation/governance-baseline-verify_EN.md	repo
```
