# GATES-EDGE-COVERAGE slice validation report (criterion-boundary convergence: OB-90 + OB-91)

> Slice: `GATES-EDGE-COVERAGE` (v1.34). Status: implementation ② complete, all gates green; the main commit `3318fef` is pushed and CI is green on both legs (see §7).

## 1. Change summary

- **OB-90 (G3(a) pairing semantics)**: a pair is now formed **from EITHER side** - an `_EN.md`-only change is judged under its CN baseline path (an untouched baseline scores `+0/-0` on its side ⇒ `G3(a.1)` FAIL instead of silence); a baseline `.md` change pairs only when an `_EN` mirror exists in the repository (so a single-language document produces no record); an unreadable side no longer skips the pair, and an incomparable directive pair makes `G3(a.2)` fail closed.
- **OB-91 (G6-6 language tag)**: the token rule becomes **whole-token strict matching plus fail-closed coexistence** - from every `lang:` occurrence the scan consumes the **continuation charset** (letters / digits / `_` / `.` / `-`, a closed enumeration; the terminator set is its complement) and then compares that token whole against the closed set; this closes the v1.33 ⑦ r3 residual (a `lang:zh_cn` truncated to the legal prefix, an invisible bare `lang:`).
- **Version package**: header v1.34, a new v1.34 row in §0.3, a new C.0v block in appendix C, and the §12.4 disposition flips plus the new OB-92.
- **Boundaries**: §2.2 / §2.4 / the §5.2 matrix / the §7.5 numbers are untouched; so are CONTRACTS / schemas / fixtures / workflows and `internal/**`; **no ADR is added** (batch assistance uses the §2.2-whitelisted `GLM-5.3-Flash (glm)`, so §1.7.4 ⑤ is not triggered).

## 2. Unit implementation and closing readings

| Unit | Content | Closing reading |
|---|---|---|
| 1 | OB-90 G3(a) pairs from either side | selftest 385/385 |
| 2 | OB-91 G6-6 whole-token strict matching | selftest 414/414 |
| 3 | directive text v1.34 + §12.4 dispositions and OB-92 | selftest 417/417 |

- **Last full run (`--scope=index`, i.e. the whole staged change set; the snapshot this report pair was finalised against)**: `SELFTEST PASS 417/417`; `TOTAL_FAIL=0`; `G3(a)` directive pair `+22/-6` with all six shapes `true`; `G4a` every inspected file OK; both `G6-5` and `G6-6` PASS.
- New self-tests: T74 (16 assertions, two mutation proofs) and T75 (16 assertions, one mutation proof) ⇒ the suite grows from 385 to 417.
- **⑨ final readings**: full tree encoding `ENC-TREE tracked=3614 text=3525 binary=89 violations=0` (the 17 frozen fixtures under `docs/validation/evidence/` are exempt by precedent); changed-line `§X.Y` references `total=114 unresolved=0`; `go build ./...` OK, `go vet ./...` OK, `go test ./...` `ok=14 FAIL=0`.

## 3. In-unit self-corrections (not counted as deviations)

A gate caught red on the spot ⇒ fixed in place ⇒ not a deviation (user ruling of 2026-10-16), item by item:

- G4b caught one new CJK character (the first character of the former wording is outside the baseline corpus) ⇒ reworded, with no whitelist entry added.
- Marker balancing: after the clauses grew, `bold` / `backticks` were momentarily unequal ⇒ localised per line and topped up on one side, ending with all six shapes `true` and zero per-line drift.
- OB-90's first mutation target was wrong (a single-line patch cannot reproduce the old semantics) ⇒ replaced by the whole guard block, which reverts faithfully to the old one-sided pairing.
- **Method lesson (kept on the record)**: a long English sentence must be balanced against its Chinese counterpart by MARKER COUNT **(bold pairs and backticks equal per line)**, and the result must be re-measured mechanically rather than eyeballed.
- **The ⑦ r1 input package was missing two declared inputs** (`slice-definition.md` / `design-r1.md`): fact = ⑦ returned `PASS` on incomplete input; root cause = the input builder **skips absent sources silently** (not a ⑤ / ⑦ fault); impact = the credibility of that `PASS` rests on the code and the mechanical readings alone; disposition = **no re-run and no nullification** in this slice, the builder now asserts every input (completed before r2), and the gap is registered as OB-93 (`To be discussed`) in §12.4.
- **A post-⑧c record script was placed under `$env:TEMP`**: after removing `tmp/gec/` a further record still had to be written, so the script lived outside the repository; the facts are `$env:TEMP\gec*.js` and `$env:TEMP\gec-*.txt`; cleanup used a wildcard plus `-ErrorAction SilentlyContinue` (**silent when the pattern matches nothing, so the cleanup cannot be proven mechanically**); `tmp/` inside the repository is empty, while **residue outside the repository has no mechanical evidence**; this is a form deviation - not counted, no new OB.

## 4. Open items and deviation ledger

- Same-type deviations: **2/3** (A-3 / A-4 of the DIRECTIVE-OB-SWEEP era); this slice adds **no new deviation**.
- **⑤ pre-review**: r1 `PRE-REVIEW: FINDINGS` (1 Medium - the chartered ledger status row was missing - plus 3 Low) ⇒ each fixed and the stage **re-run in place** ⇒ r2 `PRE-REVIEW: PASS` (0 Medium+; the remediation added three assertions, so the suite went 414 → 417).
- **⑥ independent scan**: `INDEPENDENT SCAN: PASS` (Section B empty; the blind input carried no ⑤ conclusions).
- **⑦ final review r1**: `FINAL REVIEW: PASS` (0 findings; the first round carries **no ⑤ / ⑥ conclusions** by anchoring isolation - only the change set, the design, the slice definition, the mechanical readings and the two-state baseline).
- **§12.1 deferral registration**: five slices (v1.30 through v1.34) emitted the §12.1 deferred-upgrade item in their receipts **without registering a separate OB in §12.4**; a formal registration is the user's call (the counting rule itself is OB-66 in §12.4).
- ⑩ stop point: a commit and a push both require same-turn authorization; the main commit `3318fef` was authorised, pushed, and CI is green (see §7).
- **⑧b finalisation**: the ledger completion row carries the review chain and the closing readings (the section heading marker stays `⬜` per precedent, which records completion in the row rather than in the heading).
- Freedom list: `docs/validation/evidence/GATES-EDGE-COVERAGE-freedom-list.md` (F-1..F-6, all in-bounds).
- **tmp/ process artifacts**: ⑧c has removed the probes, the review packages and the generators under `tmp/gec/`; `tmp/` now holds `.gitkeep` only (this report was finalised after the cleanup and was not regenerated).

## 5. Known gaps and registrations

- **OB-90 is closed** (`Resolved`): the pairing semantics were rewritten along with their self-test coverage - see the §6.3 G3 sub-rule and T74.
- **OB-91 is closed** (`Resolved`): whole-token strict matching plus the coverage invariant - see the §6.3 G6-6 sub-rule and T75.
- **New OB-92** (first status word `To be discussed`): ADR-024's scope wording conflicts with §1.7.4 ④ (found independently at step 0); **registered only, not fixed here**, and proposed as an input to `DIRECTIVE-TEMP-AUTH-LIMIT`.
- **Still open inside this slice scope**: none (OB-92 / OB-93 are registrations, not slice goals).
- **Entered as OB-93** (first status word `To be discussed`): the review input package has no mechanical completeness criterion - the builder skips absent inputs silently (see §3); this slice **registers without fixing**, and the §12.4 row records the candidate fix and where it may land (a standalone slice or `DIRECTIVE-TEMP-AUTH-LIMIT`).
- **Candidate (recorded in this report only, never as a §12.4 row)**: **the ⑧c record-timing conflict** - ⑧c must clear `tmp/`, yet ⑧c's own record has to land on disk ⇒ the two requirements collide; a later governance slice should evaluate "the ⑧c cleanup must run **after** the record lands, and its result must be mechanically provable (never silent)", as an input to `DIRECTIVE-TEMP-AUTH-LIMIT` or a standalone slice.

## 6. Next-step advice (candidate slice, awaiting the user's ruling)

- **Candidate slice `DIRECTIVE-TEMP-AUTH-LIMIT`**: revise the §7.5 quota cap and the §1.7.4 ⑤ threshold wording (ADR-025 records the 7th same-kind quota exception ⇒ the threshold is met) and take OB-92 as one of its inputs; **chartered independently**, never merged into a criterion-fix slice.
- **Registration manner**: the above is **advice**; this report only keeps the trail, and **nothing is registered unless the user answers explicitly**.

## 7. Commit and CI (⑧b backfill, OB-72)

- **Main commit**: `3318fef` - `fix: make G3(a) pairing and G6-6 language tags edge-safe (OB-90 + OB-91)` (10 files, `+639/-50`); pushed to `origin/main` (SSH, range `3393043..3318fef`); **gitee untouched**.
- **CI observation (this commit's own run)**: run `39019183887` (workflow `CI`, event push, attempt 1) **green on both legs on the first attempt** - `Go ubuntu-latest` 2m23s and `Go windows-latest` 3m39s; **no first-run-red handling was needed**.
- **Gates-step hard evidence (measured on both legs)**: `SELFTEST: PASS 417/417`; `GATE REPORT scope=ci changed=10`; the `G3(a)` directive pair has all six shapes `true`; `G4b newCjk=506 unknown=[]`; `TOTAL_FAIL=0`.
- **Package-level evidence**: the `Test` step on both legs covers 14 packages (`ok` 14/14, `FAIL` 0); every line containing `FAIL` in the log is a passing assertion name (for example `must FAIL`), never a failure line.
- **Backfill rule**: this section is written by the docs-only backfill; the backfill commit's own CI is observed only, and no further commit is added for it.

## 8. Artifacts

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/directive-gates-edge-coverage.md	repo
docs/validation/directive-gates-edge-coverage_EN.md	repo
docs/validation/evidence/GATES-EDGE-COVERAGE-freedom-list.md	repo
tools/gates/lib/criteria/g3a.js	repo
tools/gates/lib/criteria/g6.js	repo
tools/gates/selftest.js	repo
```
