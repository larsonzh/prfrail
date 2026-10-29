# ProofRail P35 strategic review report: five independent rounds (directive v1.56)

Slice `STRATEGIC-REVIEW-2026-10-30` (an evaluation `[SLICE]`); all five rounds ran read-only on the `independent-reviewer` carrier and the master transcribed them verbatim; baseline `HEAD` = `ebec7d1`, batch 1 commit `4ca38b6`, batch 2 commit `d5a76f7`; directive v1.56; report date 2026-10-29; R = 1, D = 0.

## 1 Executive summary

The five rounds converge on a **double break**: the real-AI production loop is not closed yet, and the record chain used to judge progress and compliance has breaks of its own; the governance and management rounds, blind to each other, both recomputed P32's rework tally as under-counted (7, where the report said 5), and the strategy round ruled that AT-23 must be split into three independently provable acceptance units instead of being **downgraded**.

## 2 Round-by-round findings

### 2.1 Round 1 · product layer

- **High｜P3**: the unknown semantics of REQ-023 and AT-24 never reached the product and business authorities, so the acceptance criteria diverge (`docs/PRODUCT_REQUIREMENTS.md`:90 vs `docs/TEST_STRATEGY.md`:54).
- **High｜P5**: AT-13 is recorded as done, yet the workflow builds the candidate directly and lacks the seed-driven build evidence the RFC requires.
- **Medium**: stale status passages, US11 outside the success criteria, and a state enumeration that omits `TERMINAL_PENDING`.

### 2.2 Round 2 · governance layer

- **High｜G4**: recomputed from its own items, P32's R is 7 rather than the 5 the report wrote, which delayed the first threshold crossing.
- **High｜G5**: P33 replaced the OB or ADR registration required by §1.7.4 with "this paragraph is the independent record".
- **High｜G7**: P34 passed "position valid, ordering only", which the exception explicitly excludes.

### 2.3 Round 3 · management layer

- **High｜M6**: an independent recomputation reached the same conclusion, R is at least 7, and the tally should have paused the slice right after the second ⑨ item.
- **Medium**: governance slices are four fifths of completions, four self-breeding chains, stale OB write-backs, and 12 missing completion rows.

### 2.4 Round 4 · strategy layer

- **High｜S1**: AT-23 must stay the S1 exit but be split into mechanism, assembly and real-candidate units; that redefines granularity, it does not lower the bar.
- **High｜S6**: T027 has an internal way out: rule on T027-B first, then close the production entry, then run a bounded real-candidate experiment.
- **High｜S7**: for the long-open item take "redefine and split", and give OB-21 an owner, stop conditions and a reconsideration trigger.

### 2.5 Round 5 · synthesis

- **High｜Z2**: the four rounds assemble two main chains, production breakage and record credibility, and neither stands alone.
- **High｜Z3**: three shared blind spots: the trust root of the evidence chain, cross-run authorisation windows, and the full lifecycle of sensitive data.
- **High｜Z4**: do only three things: close production, run security counter-examples on the same acceptance chain, and establish a credible baseline.

## 3 Priority ruling

The quarterly trade-off and the three-way disposition of round 5 are recorded here in full.

| Rank | Item | Proposed slice |
|---|---|---|
| 1 | Close one trustworthy real-AI production loop | `T027-PRODUCTION-CLOSEOUT` |
| 2 | Run independent security counter-examples on the same acceptance chain | `ACCEPTANCE-TRUST-NEGATIVES` |
| 3 | Establish a credible baseline and test governance with a business slice | `GOV-EFFICACY-BUSINESS-PILOT` |

**Adopt directly**: status and mapping corrections, count and recovery cursor, narrowed compliance claims, ledger gaps, and no over-claims.

**Needs a ruling**: the AT-23 split, the unknown blocking semantics, the black-box versus formal-execution boundary, commitments' fate, governance spend, and strategic closure.

**Needs a slice**: `T027-PRODUCTION-CLOSEOUT`, `T028-CLI-TARGET-CLOSEOUT`, `AT13-SEED-HOSTED-PROOF`, `AT12-TUI-CLOSEOUT`, `AT22-SESSION-RESUME`, `GOV-EFFICACY-BUSINESS-PILOT`, `CROSS-PROJECT-VALUE-PROOF`.

## 4 Items needing a user ruling

1. Whether AT-23 is split into three units; if so, who signs the old-to-new threshold mapping and its counter-conditions.
2. Whether T027-B opens now, and how the Windows first-release scope and the resource cap are frozen.
3. In which mode unknown blocks, and how black-box low assurance stays distinct from formal execution.
4. Once P32's R is 7, how far the correction of historical reports and the ledger goes.
5. Whether non-blocking governance is frozen, and what happens to pilots that miss their target.
6. Whether the old switches (baseline-mode=vcs, window-mode=legacy) are kept, deferred or retired.

## 5 Pipeline trimming declaration

- ⑥ skipped: this is an evaluation slice with no implementation and no diff, so ⑥ has no object (per the v24 plan and the same-turn authorisation).
- ⑦ does not apply: ⑦ presupposes an implementation object; there is none here, and this is **not** a case of external review replacing it.
- ⑧a is replaced by master transcription: all five originals were transcribed and compared with `sha256`.
- ⑨ kept: encoding gate, three scopes and the selftest.
- ⑧c kept: one-shot temp cleanup with a terminal read-back.
- ⑩ stop: this slice stops at ⑩ for commit and push authorisation.

```artifacts
docs/validation/strategic-review-2026-10-30.md	repo
docs/validation/strategic-review-2026-10-30_EN.md	repo
docs/validation/evidence/STRATEGIC-REVIEW-2026-10-30/MANIFEST.md	repo
docs/validation/evidence/STRATEGIC-REVIEW-2026-10-30/00-doc-inventory.md	repo
docs/validation/evidence/STRATEGIC-REVIEW-2026-10-30/01-product/review.md	repo
docs/validation/evidence/STRATEGIC-REVIEW-2026-10-30/02-governance/review.md	repo
docs/validation/evidence/STRATEGIC-REVIEW-2026-10-30/03-management/review.md	repo
docs/validation/evidence/STRATEGIC-REVIEW-2026-10-30/04-strategy/review.md	repo
docs/validation/evidence/STRATEGIC-REVIEW-2026-10-30/05-synthesis/review.md	repo
```
