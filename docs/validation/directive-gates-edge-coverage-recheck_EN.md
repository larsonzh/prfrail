# GATES-EDGE-COVERAGE-RECHECK slice validation report (verification only: both root causes re-checked)

> Slice: `GATES-EDGE-COVERAGE-RECHECK` (v1.37, M tier, verification / criterion slice). Baseline HEAD = `aa07eb4`; this slice is not committed (9).

## 1. Change summary

- **The only change set**: `tools/gates/selftest.js` (479 added lines, 0 deleted; the criterion implementation under `tools/gates/lib/criteria/**` and the criterion text are **untouched**).
- **New counterexample sets**: `T78` = 16 cases (OB-90), `T79` = 50 cases (OB-91, table-driven), `T79 canary` = 10 cases (frozen constants + black-box tokenization + 7 external terminator characters).
- **Selftest reading**: `460 → 536`, `SELFTEST: PASS 536/536`; self-consistency `460 + 16 + 50 + 10 = 536`.
- **Gate readings**: `TOTAL_FAIL=0` in the `index` / `tree` / `ci` scopes.
- **Version package**: `v1.36 → v1.37`, ledger metadata only (recheck confirmations appended to the OB-90 / OB-91 disposition cells, the new OB-97 row, the 0.3 log row, appendix C.0y); **the criterion text is untouched**.

## 2. Counterexample sets and readings

| Unit | Cases | What it can falsify |
|---|---|---|
| `T78` | 16 | the old OB-90 semantics (pairing only from the CN side, both sides required in the change set) would go silently quiet on an `_EN`-only change; covers the one-sided, no-mirror, unreadable-side and directive-pair fail-closed paths |
| `T79` | 50 | the old OB-91 token rule (legal plus illegal coexistence passes, `lang:zh_cn` truncated, bare `lang:` invisible) turns red on regression; table-driven coverage of coexistence order, cross-line and the `_EN` side |
| `T79 canary` | 10 | coverage inflation and self-lock: frozen constants prove the rule input is not defined by the test; black-box tokenization proves the `lang:zh` / `lang:zh_cn` / `lang:zh/en` behaviour; 7 external terminator characters are asserted one by one |

- The three blocks add 76 assertions in total, all **pure additions** with no existing assertion changed; the total and the summary line corroborate each other.

## 3. Real-change-set spot checks

| # | Scenario | Expectation | Actual | Verdict |
|---|---|---|---|---|
| S1 | `_EN` side one line longer | `G3(a.1)` FAIL, exit 1 | `+3/-0 vs _EN +4/-0 MISMATCH`, `TOTAL_FAIL=2` | conforms |
| S2 | non-directive pair, equal numstat, `_EN` bold drift | aggregate `G3(a)` PASS and **no** `G3(a.2)` record | `PASS`, detail `"bold":false`, no `G3(a.2)`, exit 0 | conforms |
| S3 | equal both sides, equal shapes | PASS | all six shapes `true` | conforms |
| S4 | single-language document without a mirror | no record | only `GATE REPORT` + `TOTAL_FAIL=0` | conforms |

- **Why the S2 expectation was corrected**: the previous round wrote `G3(a.2) FAIL` as the expectation, which was **an expectation error** - `G3(a.2)` only records for the `docs/DELIVERY_DIRECTIVE.md` pair; a non-directive pair's bold drift appears only in the aggregate detail. This round judged against the corrected expectation and the measurement matched, without editing the criterion to fit.
- **Index restoration proof**: `git write-tree` is `21b695f…` both before and after the spot checks; the final `git status --porcelain -uall` lists only `M  tools/gates/selftest.js`.

## 4. Review chain

- **⑤ r1**: `PRE-REVIEW: PASS` (0 Medium+, 5 Low).
- **⑥ r1**: `INDEPENDENT SCAN: FINDINGS` (2 must-fix + 3 boundary) ⇒ remediation ⇒ r2 `INDEPENDENT SCAN: PASS WITH FINDINGS` (1 Medium: the canary self-lock depended on the same runtime object) ⇒ remediation ⇒ r3 `INDEPENDENT SCAN: PASS`.
- **⑦ r1**: `FINAL REVIEW: FINDINGS` (2 Medium: missing dispatch receipt / count drift) ⇒ remediation (dispatch audit + final ③ rerun) ⇒ r2 `RE-REVIEW: PASS WITH FIXES` (only 1 Low: audit snapshot timing).
- **Closing the ⑦ r2 Low**: ⑦ r2 noted that the dispatch-audit snapshot lagged (Low); after ⑧a, `audit-dispatch.js` was rerun to overwrite the audit file, so the table shows the final snapshot.
- **The ⑦ quota 1+2 is exhausted, so no third read was opened**; the r2 Low is closed by the report record, without a quota exception.

**OB-95 / OB-96 re-check (1/2)**: this slice is the **1/2 data point for both the OB-95 and OB-96 re-checks**: ① ran once on `DeepSeek V4 Pro (deepseek)` (OB-95), and ②③⑧ were dispatched to `prfrail-implementer` / `prfrail-tester` / `prfrail-documenter` (OB-96). **Neither recurred**, so both stay **under observation without escalation**; the 2/2 data point is the next slice `G9-c` - no recurrence moves them to `Resolved`, a recurrence escalates to `To be discussed`.

## 5. Dispatch self-attestation

| Position | `agentName` | Model | Transcript line offset | Calls |
|---|---|---|---|---|
| ① | `deep-reasoner` | `DeepSeek V4 Pro (deepseek)` | `12211` | 1 |
| ② | `prfrail-implementer` | same as the master | `12377` / `13163` / `13367` | 3 |
| ③ | `prfrail-tester` | same as the master | `12633` / `13811` | 2 |
| ⑤ | `independent-reviewer` | `GLM-5.3 (glm)` | `12922` | 1 |
| ⑥ | `independent-reviewer` | `MAI-Code-1.1-Flash (copilot)` | `13074` / `13312` / `13510` | 3 |
| ⑦ | `independent-reviewer` | `GPT-5.3-Codex (copilot)` | `13546` / `13954` | 2 |
| ⑧ | `prfrail-documenter` | same as the master | `14242` | 1 |

- Audit snapshot: `tmp/GATES-EDGE-COVERAGE-RECHECK/dispatch-audit.txt` (rerun and overwritten; the snapshot holds 13 dispatch records); `14242` is the last snapshot entry and maps to the ⑧a report-and-ledger write; the ⑧b finalization is dispatched this round, after the snapshot, and adds no snapshot row (the snapshot counts `prfrail-documenter` as 1).
- Per-position counts: ① 1 / ② 3 / ③ 2 / ⑤ 1 / ⑥ 3 / ⑦ 2 / ⑧ 1; **no position has zero dispatches**.

## 6. Empirical evidence that the temporary-artifact discipline failed

- **Fact**: the start message already put "temporary artifacts must live under `tmp/<slice-id>/`, a violation stops and escalates" up front, yet `$env:TEMP` accumulated **5** same-type violations (1 in this slice + 4 in `GATES-G1G5`).
- **Conclusion**: **the rule text has no effect on Flash** - the same hard rule was violated repeatedly even though it was written out, placed up front and backed by a stop clause; the chartering condition for `G9-c` (a mechanical temporary-artifact criterion) is **far beyond its threshold**, so it is chartered in the next slice.
- **Same-type nature**: this slice's theme is "is a root cause claimed closed actually closed"; this slice itself repeated exactly the same-type failure on temporary-artifact discipline - **it is its own counterexample and thereby validates its own theme**.
- Every in-repository process artifact lives under `tmp/GATES-EDGE-COVERAGE-RECHECK/` and ⑧c removes them one by one; residue outside the repository has **no mechanical evidence** (recorded truthfully).

## 7. Known gaps and boundaries

- **The five ⑤ Low items**: wording and coverage advice only, no verdict change; recorded against the r1 conclusion.
- **The ⑥ canary boundary**: the r2 Medium showed the first canary depended on the same runtime object (self-lock); r3 closed it with frozen constants plus black-box tokenization, but the structural coupling "constants and the rule live in the same file" remains a **recorded boundary**.
- **① report-level candidates (recorded here only, never as 12.4 rows)**: the `G3(a)` scope asymmetry (the empty-set test on `changedMarkdown` versus the loop over `changedPaths`, so an evidence-tree `_EN` is judged inside a mixed change set); heading **level** drift is uncovered (only presence is compared); a skip caused by dedup order; an illegal tag beyond `LANG_HEADER_WINDOW` is invisible.
- **OB-97** (`Observing`): `lang:zh:en` / `lang:zh:` truncate at the second colon, then declare `zh` without a `UsageError`; anchor `tools/gates/lib/criteria/g6.js`, id about 735; the 6.3 "every expression" sentence over-claims literally ⇒ **registered only, not fixed here**.
- **Pilot deferral**: the **10th** (this is a verification / criterion slice, meeting 12.1 companion ②); this slice adds **no deferral OB**.
- **How the single ⑦ r2 Low was closed**: the single Low from ⑦ r2 (the dispatch-audit snapshot was taken before the final round) was **closed by the master's own mechanical re-run** of the audit script, **not by an independent re-review**; the ⑦ quota (1 + 2) was exhausted, so the loop stopped there and **no third read was opened**. It is a "Low may be deferred" case and is **registered honestly rather than written off**.

## 8. Cost and metering

| Position | Metering | Note |
|---|---|---|
| ① | 1 | design re-check (narrowed input) |
| ② | 3 | first counterexample sets + canary remediation + self-lock independence remediation |
| ③ | 2 | counterexample execution with spot checks + final snapshot rerun |
| ⑤ | 1 | r1 (quota 1+1, only 1 used) |
| ⑥ | 3 | r1 / r2 / r3 (mandatory, 7.2.1 ②) |
| ⑦ | 2 | r1 / r2 (quota 1+2, exhausted) |
| batch assistance | 0 | no applicable task among the four uses |
| probes | 0 | probe quota unused |

**The second data point for the v1.35 tiering**: this slice used ⑦ 1 + 2 in full (r1 `FINDINGS` → r2 `PASS WITH FIXES`, one Low still open) and ⑥ 3 of 3 (r1 `FINDINGS` → r2 `PASS WITH FINDINGS` → r3 `PASS`). It points the same way as `GATES-G1G5` (L tier, ⑦ 1 + 3 used in full) ⇒ **the tiering may be tight for criterion / verification slices**, since the last round keeps producing Low / Medium+ findings. **Nothing is escalated and §7.5 is not changed**; this is a tuning observation, and if the `G9-c` slice reproduces the pattern the tiering can be evaluated separately.

## 9. Commit and CI

- This slice is **not committed** (⑩ stop point); the commit and CI readings are filled in by ⑩ after same-turn authorization.

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/directive-gates-edge-coverage-recheck.md	repo
docs/validation/directive-gates-edge-coverage-recheck_EN.md	repo
docs/validation/evidence/GATES-EDGE-COVERAGE-RECHECK-freedom-list.md	repo
tools/gates/selftest.js	repo
```
