# DIRECTIVE-TEMP-AUTH-LIMIT slice validation report (temporary-authorization threshold and quota cap: §1.7.4 ⑤ + §7.5 + OB-92 + OB-93)

> Slice: `DIRECTIVE-TEMP-AUTH-LIMIT` (v1.35). Status: implementation ② complete, all gates green; ⑦ r1 `PASS`; stopped at ⑩, and a commit requires authorization in the same turn.

## 1. Change summary

- **§1.7.4 ⑤ gains four judging elements**: (i) "the same class" means the same exceeded clause plus the same exception category (a closed set of five); (ii) counting takes the §12.4 ledger and the ADR register as its single source of truth (an instance registered in both counts once), the window starts at that clause's most recent formal revision and restarts from zero once it lands; (iii) the third instance must not be authorized and escalates at step 0; (iv) a same-template kickoff does not extend the previous slice's authorization and must be re-ruled in the same turn and counted toward that family's running total (the OB-92 landing sentence).
- **The §7.5 ⑦ quota is tiered**: S 1+1, M 1+2, L / governance / high-risk 1+3, with blind review counted separately as +1; a "tier-decision rule" is added (the table's "tier" is the §1.2 size tier, `M-L` counts as M, several matching rows take the most re-reviews, governance and documentation slices are defined, and in-flight slices keep their kickoff quota) plus a "revision basis" (ADR-025's seven instances; the newest, `GATES-TIGHTEN`'s ⑦ r1-r3).
- **§5.1 gains two timing rules**: the review input builder must assert every declared input exists and attach an input manifest (OB-93); the ⑧c cleanup must run after the record lands, item by item, with mechanical read-back, and temp artifacts always go under `tmp/<slice-id>/`.
- **Ledger and appendix**: the §12.4 dispositions of OB-92 / OB-93 flip to `Resolved` and OB-94 is added (the residual once the ⑤ quota was spent, flipped to `Resolved` after the ⑦ re-check); the three §12.5 gap rows are registered at slice start and flip to `Landed` when the slice lands; appendix C.0w; the ADR-024 disposition gains an OB-92 revision note (**without rewriting** the historical ruling text).
- **Boundary**: the §2.2 / §2.4 closed sets, the §5.2 matrix, the §6.3 criterion bodies, CONTRACTS / schema / fixtures / workflows and `internal/**` are untouched; **no ADR is added** (the cap stays 25 - step 0 held that a formal revision path needs no ADR).

## 2. Unit implementation and closing readings

| Unit | Content | Closing reading |
|---|---|---|
| 1 | §1.7.4 ⑤ four elements plus the OB-92 landing sentence | gates green |
| 2 | §7.5 tiered quotas plus the tier-decision rule and revision basis | gates green |
| 3 | §5.1 two timing rules plus OB-93 | gates green |
| 4 | ledger (OB-92/93/94, three §12.5 rows, C.0w, ADR note, P10) | gates green |

- **Last full run (`--scope=index`, the snapshot this report pair is finalised against)**: `SELFTEST PASS 417/417`; `TOTAL_FAIL=0`; `G1-b` reports `header=v1.35 last=v1.35` for both files; all three `G3(a)` pairs have symmetric numstat and all six shapes `true`; `G4b unknown=[]`.
- **New self-tests**: none (documentation-only slice; the denominator stays 417/417).
- **⑨ final readings**: full-tree encoding `ENC-TREE tracked=3617 text=3528 binary=89 violations=0` (the 17 frozen fixtures under `docs/validation/evidence/` are exempt by precedent); changed-line `§X.Y` references `total=277 unresolved=0`; `go build ./...` and `go vet ./...` pass, and `go test ./...` reports `ok=14 FAIL=0` (run by the base gate).

## 3. In-unit self-corrections (not counted as deviations)

A gate caught red on the spot ⇒ fixed in place ⇒ not a deviation, item by item:

- `G1-b`: the first v1.35 log row was placed at the top of the §0.3 table ⇒ header and last-row versions disagreed ⇒ moved to the table tail (last row = current version).
- `G4b`: three characters outside the baseline corpus were hit ⇒ rewritten as "leave ... only to", "three constraints" and "ADR table" respectively, with **no whitelist extension**.
- `G6-4` false positive: the OB-94 row wrote the version number followed by the word for a ledger line, read as an expiring line count ⇒ reworded to "entry".
- `G3(a.1)`: a stray blank line at the tail of the `P10` block made the CN/EN `numstat` 17/18 ⇒ aligned to 18/18.
- Two assertions in the local `checks.js` mis-read a pass as a failure (A4 sliced a fixed length, A5 counted unrelated lines) ⇒ fixed; both pass correctly now.
- Two extra bold markers on the EN side of the ADR-024 row are **pre-existing drift** (HEAD already has CN 20 / EN 22; this slice adds +2 symmetrically on both sides and does not widen it) ⇒ **not changed**, only recorded.

## 4. Open items and deviation ledger

- Same-type deviations: **2/3** (none added by this slice).
- **⑤ r1**: `PRE-REVIEW: PASS WITH FIXES` (4 Medium + 5 Low) ⇒ each fixed in place.
- **⑤ r2**: `PRE-REVIEW: PASS WITH FIXES` (2 Medium residuals + 4 Low) ⇒ all fixed in place; the 1+1 ⑤ quota was spent ⇒ per §7.8.3 the residual was folded into the ⑦ input and registered as OB-94, with **no third round opened**.
- **⑥**: `INDEPENDENT SCAN: PASS` (blind input, Section B empty).
- **⑦ r1**: `FINAL REVIEW: PASS` (0 findings; F-1 / F-2 closed; ⑦ not downgraded, `GPT-5.3-Codex` Extra High).
- **Self-reference registration (a property of this slice, stated explicitly)**: while revising the §1.7.4 ⑤ threshold this slice **hit that very threshold itself**; per the "quota exhausted" row of §7.8.3 it took the "fix plus fold into the ⑦ input" path, handed the closure paths of F-1 / F-2 to ⑦, and registered **OB-94** as the carrier; **nothing is nullified and no extra ⑤ round is granted**.
- **§12.1 pilot deferral (the 8th)**: this slice is a governance / rule-revision slice and does not meet the §12.1 pilot criteria (matching companion ②); per the §12.4 single source of truth it is recorded as the **8th** such instance (the five slices v1.30 through v1.34 are the same case) and this slice **adds no deferral OB**; a formal registration is the user's call.
- ⑩ stop point: a commit and a push both require same-turn authorization; this slice has no commit yet.
- **⑧c cleanup (executed)**: every helper under `tmp/DIRECTIVE-TEMP-AUTH-LIMIT/` (scripts, review packages, raw gate readings) was removed item by item; the read-back shows `Test-Path` as `False` and `tmp/` holding `.gitkeep` only.

## 5. Known gaps and registrations (the section the ruling calls "§8")

- **OB-92 / OB-93 / OB-94 all flip to `Resolved`**: OB-92 closes with the "⑤(iv) landing sentence plus ADR revision note plus no rewrite of the historical ruling" triple; OB-93 closes with the two §5.1 timing rules; OB-94 closes after the ⑦ r1 re-check of F-1 / F-2.
- **F-3..F-5 (Low) remediation record**: F-3 "closing sentence of this section" → "table-following sentence of this section"; F-4 the §12.5 ⑦ budget row gains "from v1.35 the round cap is taken from the slice's tier"; F-5 the step-0 anchor gains "the first item of the §1.2 steps field; a slice without that item must perform it in its kickoff receipt".
- **Still open inside this slice scope**: none.
- **Candidate (recorded in this report only, never as a §12.4 row)**: the tiered §7.5 quotas have no real-slice evidence yet (n=0) ⇒ the next slice should run under the new table and record actual usage; if an M-tier slice returns `PASS` without using the extra round, a later governance slice may return M to 1+1.

## 6. Next-step advice (awaiting the user's ruling)

- **Run the next slice under the new quota table** (M tier, or an L-tier / governance slice) to test whether "1+2 / 1+3 with narrowed input" suffices, and record the actual usage in that slice's "cost and metering" section (R7.1 caliber).
- **The §12.1 pilot** is still the user's call (this slice is a governance slice and does not meet the pilot criteria).

## 7. Commit and CI

- This slice **does not enter the commit stage** before the ⑩ stop point; commit and push need same-turn authorization, so this section carries no CI readings yet.
- After the commit, both legs are observed per §6.4 and the readings are written into this section by a later docs-only commit (which also needs same-turn authorization).

## 8. Artifacts

```artifacts
docs/ADR_REGISTER.md	repo
docs/ADR_REGISTER_EN.md	repo
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/directive-temp-auth-limit.md	repo
docs/validation/directive-temp-auth-limit_EN.md	repo
docs/validation/evidence/DIRECTIVE-TEMP-AUTH-LIMIT-freedom-list.md	repo
```
