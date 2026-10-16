# GATES-TIGHTEN slice validation report (criterion tightening + merged text items)

> Slice: `GATES-TIGHTEN` (v1.33). Status: implementation ② complete, the ⑤ / ⑥ / ⑦ review chain ran, all gates green; stopped at ⑩, and a commit requires authorization in the same turn.

## 1. Change summary

- **Criterion tightening**: G3(a.2) gains a sixth shape (backtick count); the G6-5 endpoint self-check; G6-6 selects its word table from the language declared in the document header; new G7-e (verdict-line closed set) and G5-c (appendix-B call strings against the §2.2 qualified-name column); the G6-3 angle-bracket rule is narrowed; G5-a(1) splits its exit codes.
- **Rule-input change (the user asked for an explicit registration on 2026-10-16)**: this slice introduces the header language declaration fields `lang:zh` / `lang:en` for the G6-6 criterion; the scope is limited to files that carry a §12.4 ledger interval with rows to judge; other bilingual documents are not required to carry them today; if this ever becomes a repository-wide documentation discipline, it charters a separate slice.
- **Eight merged text items**: the §5.1 freedom-list timing note, the §7.2.2 ⑥ mandatory-surface layering, the §7.3 ⑦ tool-surface calibration, a new §9.3 failure-mode row, the §11.2 write-back domains, the §12.5 gap-row timing, the §6.3 bilingual exclusion-marker obligation, and the §6.3 criterion scope and post-freeze appends. **Shape note**: the new §9.3 row follows that section's existing blank-line-separated chunk shape (the pre-existing breakage is not repaired here and is left to a later governance slice).
- **Version package**: header v1.33, a new v1.33 row in §0.3, a new C.0u block in appendix C, and the §12.4 disposition column.
- **OB-73 shape alignment, registered place by place** (§11.1; backticks added / removed only, no content rewritten): CN lines 887 / 933 / 955 / 1370 and EN lines 576 / 1314 / 1597; **three of them in appendix C** = 1314 / 1370 / 1597. **Correction on the record: the EN place 1314 additionally swaps the quote style (「竖」 → `竖`) beyond adding backticks, which is semantically lossless; the other six were verified by paired diff to be pure add/remove.**
- **OB-73 logged-count correction**: the entry recorded 5 places, while recomputation on the `git show HEAD` baseline yields the same **7 places with delta +2** ⇒ all pre-existing drift (not introduced by this slice), so the logged count was **an under-registration**.
- **The disposition column actually changed in 12 places** (11 first-status-word flips + the OB-41 backfill note); the "five" in the task definition refers to the five text-class OBs, while the flips of the seven tightening OBs were implied by necessity. Registered here for accuracy.
- **OB-73 known cost**: the six-shape rule raises the bilingual discipline cost - an EN-only sentence needs a matching marker added on the CN side; ⑦ is asked to rule whether that is acceptable.
- **Residual limit (candidate OB, awaiting the user's ruling)**: all six G3(a.2) shapes are whole-file count / position comparisons ⇒ **a cross-line compensating drift is not covered**; this is stated in the §6.3 approximate-criterion statement.

## 2. Unit implementation and closing readings

| Unit | Content | Closing reading |
|---|---|---|
| 1 | OB-20 G5-a(1) exit-code split | selftest 339/339 |
| 2 | OB-17 G6-3 angle-bracket narrowing | selftest 344/344 |
| 3 | OB-46 G6-5 endpoint self-check (fail-closed) | selftest 347/347 |
| 4 | OB-47 G6-6 word table by declaration + name gate | selftest 356/356 |
| 5 | OB-73 G3(a.2) sixth shape + 7 aligned places | selftest 362/362 |
| 6 | OB-42 G7-e verdict-line closed set | selftest 373/373 |
| 7 | OB-29 G5-c appendix-B call-string set | selftest 381/381 |
| 8 | eight merged text items | selftest 381/381 |
| 9 | version package + §12.4 disposition column | selftest 381/381 |

- **Last full run (`--scope=index`, i.e. the whole staged change set; the snapshot this report pair was finalised against)**: `SELFTEST PASS 385/385`; `TOTAL_FAIL=0`; `G3(a)` directive pair `+59/-24` with all six shapes `true`; `G4a` every inspected file OK; `G4b newCjk=667 unknown=[]`; both `G6-5` and `G7-e` PASS.
- **⑨ native validation (Windows host)**: `go build ./...` OK, `go vet ./...` OK; `go test ./...` reports **14 packages ok and 0 FAIL** (plus 5 "no test files" entries).
- **⑨ tree-wide encoding check**: `tracked=3611 text=3522 binary=89` with **zero violations in the non-frozen main tree**; **17 further ones** sit inside frozen evidence packages (the standing exemption, the same value as the v1.32 baseline).
- **⑨ reference resolvability on changed lines**: **286** section references on added lines and **9** unresolved - composed of 2 example numbers (the **deliberate FAIL-shape example** inside the G6-5 rule text, already present in HEAD) plus 7 legacy numbers (historical section numbering, inside the region that G6-5 declares out of scope and in the inherited class). **No dangling reference was introduced by this slice.**
- The re-checked readings that include this report pair are recorded in the ⑩ stop-point receipt (adding the pair changes the change-set count and what G7 covers).

## 3. In-unit self-corrections (not counted as deviations)

A gate caught red on the spot ⇒ fixed in place ⇒ not a deviation (user ruling of 2026-10-16), item by item:

- G3(a.1) went red on the spot (the first alignment pass had 6 CN / 1 EN edits) ⇒ the directions were rebalanced to three each.
- G3(a.2) went red twice on the bold shape (one missing marker in the OB-29 clause; a spurious asterisk pair on the version-package §0.3 row and a missing marker in the OB-41 note) ⇒ localised mechanically and fixed.
- `create_file` wrote `g5c.js` with CRLF ⇒ G4a went red ⇒ normalised to LF.
- G6-4 matched new wording (the OB-46 `§3.10` row wording) ⇒ the wording was rewritten, with no whitelist entry added.
- A test read the wrong record (no `G3(a.2)` record exists when nothing fails) ⇒ it now reads the shape JSON from the aggregate record.
- The report generator emitted a real backspace (U+0008) for a single-escaped `\b` ⇒ the backslash was doubled and a control-character scan added, now zero on both sides.
- G4b caught one new CJK character (the last character of the former wording is outside the baseline corpus) ⇒ reworded to "shape mismatch", with no whitelist entry added.
- The anchor-correction sentence had been over-replaced into "corrected from `g3a.js:34` to `g3a.js:34`" ⇒ restored to `g3a.js:35` to `g3a.js:34`.
- **Method lesson (kept on the record)**: an EN table cell or paragraph anchor **must not be written from memory** - the OB-29 and OB-73 units each produced a non-matching anchor that way ⇒ always export the verbatim text with a script first, then edit.

## 4. Open items and deviation ledger

- Same-type deviations: **2/3** (A-3 / A-4 of the DIRECTIVE-OB-SWEEP era); this slice adds **no new deviation**.
- **§12.1 deferral registration**: v1.30 / v1.31 / v1.32 and this slice all emitted the §12.1 deferred-upgrade item in their receipts without registering a separate OB in §12.4; this slice follows the same practice; a formal registration is the user's call.
- ⑩ stop point: a commit and a push both require same-turn authorization; this slice has no commit yet.
- Freedom list: `docs/validation/evidence/GATES-TIGHTEN-freedom-list.md` (F-1..F-9, all in-bounds).
- **Candidate OB registrations (user ruling of 2026-10-16: registered in the report and §12.5 only, never as new §12.4 rows)**: (1) the whole-file count unit of G3(a.2) does not cover a per-line compensating drift (the approximate-criterion statement is on the record in §6.3; extending it to a per-line comparison is left to a later governance slice); (2) CJK inter-character spaces in bilingual files have no criterion coverage (measured by the ⑤ pre-review: 10 places on the CN side, none on the `_EN` side, none of them on a line changed by this slice).
- **⑧c wrap-up cleanup (this slice)**: (1) the **tree-wide encoding check**: `tracked=3611 text=3522 binary=89` with zero violations in the non-frozen main tree and 17 inside the frozen evidence packages (the standing exemption); (2) **process artifacts removed**: the whole `tmp/gates-tighten/` directory (including `DESIGN.md`, the probe scripts, the report generator, and both review packages `review/` and `review-r3/`) is gone, leaving only `.gitkeep` under `tmp/`. **Verdict line**: `WRAP-UP CLEANUP: PASS` (zero main-tree violations, zero process-artifact residue).
- **Review chain on the record**: ⑤ `GLM-5.3` r1 `PRE-REVIEW: PASS WITH FIXES` (1 Medium plus several Low / Info) ⇒ fixed ⇒ r2 `PRE-REVIEW: PASS WITH FIXES` (0 Medium+); ⑥ `MAI-Code-1.1-Flash` `INDEPENDENT SCAN: PASS` (0 findings); ⑦ `GPT-5.3-Codex` (Extra High) r1 `FINAL REVIEW: FINDINGS` (2 High + 1 Medium + 1 Low = 3 Medium+) ⇒ fixed ⇒ r2 `FINAL REVIEW: FINDINGS` (1 Medium + 1 Low) ⇒ the user ruled an in-place fix plus a narrow read ⇒ r3 `RE-REVIEW: FINDINGS` (1 Medium).
- **⑦ input packages (rebuilt each round)**: the full package in `tmp/gates-tighten/review/` and the narrow one in `review-r3/`; both were rebuilt before dispatch and checked for "index = disk" (drift=0).
- **Pending authorization**: neither `git commit` nor `git push` **has been run** (the ⑩ stop point requires authorization in the same turn).

## 5. Known gaps and registrations (the section the ruling calls "§8")

- **OB-90 (a pre-existing gap, caught independently by ⑦ r1)**: G3(a) pairs only from the CN side and requires both sides in the change set ⇒ changing only `_EN.md` produces no G3(a) record at all (the bilingual-symmetry criterion fails silently). Anchors: `tools/gates/lib/criteria/g3a.js:34` and `:69`. Registered per the user ruling of 2026-10-16 as §12.4 OB-90 (first status word `To be discussed`) and as a §12.5 gap row; **not fixed here** (fixing it changes the pairing semantics itself, beyond this slice).
- **Why ADR-024 is in the change set**: ADR-024 is a **step-0 by-product** (the kickoff message states that step 0 must register ADR-024 first) and **does not consume an in-bounds slot**; it is a registration addition, not a defect.
- **Medium+ closed inside this slice**: ⑦ r1 High-2 (a `lang:zh-cn` prefix match silently treated as zh) is fixed in place to **whole-token exact comparison**, with `lang:zh-cn` / `lang:en-us` counter-example assertions added (self-test 381 → 383); Medium-3 (stale report readings) is backfilled here; Low-4 (package metadata) is rebuilt with the new input package.
- **⑦ r2 reported a Medium (unclosed, quota exhausted)**: the G6-6 "legal + illegal tag coexist" bypass (a header carrying both `lang:zh` and `lang:zh-cn` lets the latter through); registered per §7.8.3 as §12.4 OB-91 (first status word `To do`, noted as "fold into the next ⑦ input"); **not fixed here** (the ⑦ quota of 1 final + 1 re-review is exhausted, and a third round requires user authorization).
- **Factual correction (⑦ r2 Low)**: OB-90's first anchor is corrected from `g3a.js:35` to `g3a.js:34` (three places updated); a fact-recording correction that changes no rule semantics.
- **⑦ r3 quota exception (ADR-025)**: the user ruled directly on 2026-10-16 to authorise a **narrowly scoped third read** (it only checks whether the G6-6 coexistence bypass is closed) and required the registration to state truthfully that it is "a rule-semantics boundary change, beyond the OB-11 exception literal standard, not usable as a precedent"; **the same-kind count has reached 7 ⇒ the §1.7.4 ⑤ threshold is met, so a slice must be chartered to revise §7.5 formally**.
- **Form slip (not a deviation, recorded as fact)**: the final staging used `git add -- .`, which does not match the §10.3 wording of an explicit per-file list; the outcome is equivalent (`tmp/` is gitignored, porcelain shows exactly the expected 21 paths, no untracked residue) and no OB is opened for it.
- **⑦ r3 (the narrow read) verdict = `RE-REVIEW: FINDINGS` (1 Medium, unclosed)**: the fix closed the r2 path (`lang:zh-cn` / `lang:en-us`, including the illegal tag first / last, multiple lines and the `_EN` side), but a residual of the same kind remains - the token rule `\blang:[A-Za-z][A-Za-z.-]*` truncates `lang:zh_cn` to the legal `lang:zh` and never recognises a bare `lang:` ⇒ it is still let through when it coexists with a legal tag. Under the user's ruling of 2026-10-16 **no fourth read is opened**, so OB-91 closes as "fold into the next ⑦ input" (quota exception registered in ADR-025).
- **Still open inside this slice scope**: the `OB-91` residual (malformed tag forms, see the bullet above; folded into the next ⑦ input under the user's ruling and not fixed further here).

## 6. Next-step advice (candidate slices, awaiting the user's ruling)

- **Candidate slice A (criterion-boundary convergence)**: merge the §4-registered **OB-90** (the G3(a) one-sided pairing bypass) with **OB-91** (G6-6 `declaredLang` switched to whole-token strict matching + coexistence fail-closed); both share the domain, the kind and the review chain, so **splitting them into two slices is not advised**; working name `GATES-EDGE-COVERAGE` (fixed at start-up). Both touch the criterion body in §6.3 ⇒ run it as a `[SLICE]` and do not downgrade ⑦.
- **Candidate slice B (rule-semantics revision)**: revise the §7.5 quota cap and the §1.7.4 ⑤ threshold wording - the basis is the **7th same-kind quota exception on record in ADR-025 ⇒ the ⑤ threshold is met**; working name `DIRECTIVE-TEMP-AUTH-LIMIT` (fixed at start-up). **Chartered independently**, never merged into candidate A.
- **Registration manner (per the §11.2 registration advice and its trail)**: the above are **advice**; this report only keeps the trail, and **nothing is registered unless the user answers explicitly**.

## 7. Commit and CI (⑧b backfill, OB-72)

- **Main commit**: `c6ada79` - `feat: tighten the G3(a.2)/G5/G6/G7 gate suite and land the v1.33 text merges` (changed files 21, `+1275/-83`); pushed to `origin/main` (SSH, range `68afde6..c6ada79`); **gitee not pushed**.
- **CI observation (the run for that commit)**: run `38965047253` (workflow `CI`, event push, attempt 1) **first run green on both legs** - `Go ubuntu-latest` 2m20s, `Go windows-latest` 3m52s; **the first-run-red procedure was not triggered**.
- **Gates-step hard evidence (both legs)**: `SELFTEST: PASS 385/385`; `GATE REPORT scope=ci changed=21`; `G3(a)` all four pairs OK; the report pair has all six shapes `true`; `G4b newCjk=667 unknown=[]`; `TOTAL_FAIL=0`.
- **Package-level evidence**: `ok` lines = ubuntu `Test` 14 + ubuntu `Race` 2 + windows `Test` 14; `FAIL` matches inside the `Test` / `Race` steps: 0; on the windows leg the `Race` and `Contract fixtures` steps are skipped by the platform condition.
- **Reading-scope note**: the `newCjk=667` above belongs to **that commit\'s change set (21 files)**; this backfill commit\'s own index range holds only the report pair (2 files, measured `newCjk=63`), so the two readings do not conflict.

## 8. Artifacts

```artifacts
docs/validation/directive-gates-tighten.md	repo
docs/validation/directive-gates-tighten_EN.md	repo
docs/validation/evidence/GATES-TIGHTEN-freedom-list.md	repo
tools/gates/lib/criteria/g5c.js	repo
```
