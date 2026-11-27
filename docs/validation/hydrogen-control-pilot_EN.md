# P97 HYDROGEN-CONTROL-PILOT validation report (light-tier comparative review pilot, docs-only)

- Slice: `HYDROGEN-CONTROL-PILOT` (P97, `[SLICE]`, S tier). Lead: `DeepSeek V4.1 Flash (deepseek)` High.
- Date: 2026-11-28. Directive version: v1.82 (MINOR). R = 0, D = 1.
- Evidence pack: `docs/validation/evidence/P97-HYDROGEN-CONTROL-PILOT/`.

## 1. Conclusion

The light tier `Gemini 4 Hydrogen (copilot)` completed one comparative review round on the P87 deliverables, closing the P87 report section 7 pre-gap; results are recorded across four faces and P87 is never re-adjudicated. Conclusion line: `SLICE: READY_FOR_REVIEW`.

## 2. Comparison target and input pack

- Target = the P87 evidence pack plus the P87 effective surface of directive v1.81 (read via `git show 8f0c42b:`).
- Input shape = 14 payloads, 1,256,841 bytes (about 1.2 MB); all P87 `04-reviews/`, the P87 report pair, `CHECKLIST_v80` and every P97 artifact were withheld, so the verdict was blind to Hydrogen.
- Inequality note: the original P87 (7) read the change-set diff while this read uses the landed form plus the evidence pack, so non-comparable faces are marked "incomparable".

## 3. Comparison result (four faces)

| Face | P87 original (7) | Hydrogen | Reading |
|---|---|---|---|
| Verdict | `FINAL REVIEW: PASS` | `FINAL REVIEW: FINDINGS` | deviates |
| Finding count | new Medium+ = 0, Low = 0 | Medium+ = 3, Low = 2 | deviates |
| Severity | all bands zero | C0 / H1 / M2 / L2 | deviates |
| Misjudgement face | one blind-review Medium, closed | item 4 shares the same target, graded Low | overlap (grade differs) |

## 4. Hydrogen findings (registered as observations, never re-adjudicated)

- Item 1 High: the section 5.2 deep-tier tier column says `Max` while section 8.2 maps `gemini-4-neon` to `Extra High`, a direct conflict; the original P87 (7) `PASS` missed it and the light tier caught it.
- Item 2 Medium: the section 5.2 matrix lacks literal decidability (undefined terms, missing tier exclusivity and priority).
- Item 3 Medium: section 7.5 lists two context bands for `gemini-4-neon` but omits the price of the other band.
- Items 4 and 5 Low: section 7.5 alternative-quota wording (same target as the P87 blind-review Medium) and the stale self-count inside the `numstat` evidence file.

## 5. Review chain

- (5) Pre-review (`DeepSeek V4 Pro`): `PRE-REVIEW: PASS`, 4 Low items disposed.
- ⑥ skipped: no criterion implementation file touched (`tools/gates/**` unchanged), confirmed by the (7) final review (P72 precedent).
- (7) Pilot round (`Gemini 4 Hydrogen (copilot)`): one round, `FINAL REVIEW: FINDINGS`, outside the slice quota (section 1.7.4 plus OB-146).
- (7) Final review (`DeepSeek V4 Pro`): `FINAL REVIEW: PASS`, 3 Low items cleared.

## 6. Cost comparison

- Hydrogen round 1 delta = 184.3 credits (user-side single source of truth; input pack 1.2 MB).
- Reference: the P85 four-read matrix on one shared input measured Sol 67.9 / Codex 129.9 / Carbon 266.8 / Neon 493.4 credits; this reading uses a different input and is not directly comparable.
- Adoption value of the light tier is ruled by the user after the report is produced.

## 7. Version pack and registrations

- The directive moved to v1.82 (MINOR): header, a new section 0.3 row, and OB-146 in section 12.4; the section 7.5 body is unchanged.
- OB-146 carries the five elements plus the OB-37 four elements; the same-class window starts at v1.72 and this slice is instance 1.
- The ledger P97 section and the P87 "gap filled" entry are both backfilled bilingually.

## 8. Mechanical readings

- Gates tree and index both report `TOTAL_FAIL=0`; `SELFTEST` stays 806/806, equal to the start baseline (no criterion code touched).
- The evidence pack MANIFEST and SHA256SUMS entry sets are equal (8 entries); the bilingual ledger rows are symmetric.

## 9. Observations

- O-1: the P96 out-of-quota narrow re-review lacked a section 12.4 record against the section 1.7.4 item 2 wording (user ruling; recorded only, not retro-fixed).
- O-2 / O-3 / O-4: the Hydrogen High item and its two Medium items; disposition belongs to the user. **Disposition direction (user ruling, 2026-11-28): handle soon - fold into the next governance slice, or a standalone small revision slice named by the user; revision targets = align the section 5.2 deep-tier tier column and tier-authorization sentence with section 8.2, the Neon price band attribution in section 7.5, and the matrix decidability wording.**
- O-5: the two Hydrogen Low items; disposition belongs to the user.

## 10. Anomalies and deviations

- D-1 (D = 1): the input-pack path declaration and the actual layout disagreed (found by the reviewer, no effect on the verdict; see `03-deviation-log.md`).
- Input size: the five target sections were read in full in both CN and EN; over-long historical lines were truncated at the tool layer (outside the target sections).

## 11. Open items and handover

- The three new Hydrogen findings stay undisposed (no re-adjudication of P87); the user rules whether to open a later slice.
- Both CI legs and the commit are evidenced and backfilled through the (9) and (8b) steps.

## 12. Artifact fence

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/directive-review-model-policy.md	repo
docs/validation/directive-review-model-policy_EN.md	repo
docs/validation/hydrogen-control-pilot.md	repo
docs/validation/hydrogen-control-pilot_EN.md	repo
docs/validation/evidence/P97-HYDROGEN-CONTROL-PILOT/00-briefing.md	repo
docs/validation/evidence/P97-HYDROGEN-CONTROL-PILOT/input-list.md	repo
docs/validation/evidence/P97-HYDROGEN-CONTROL-PILOT/freedom-list.md	repo
docs/validation/evidence/P97-HYDROGEN-CONTROL-PILOT/01-arch/01-arch-design.md	repo
docs/validation/evidence/P97-HYDROGEN-CONTROL-PILOT/02-ruling-stop1.md	repo
docs/validation/evidence/P97-HYDROGEN-CONTROL-PILOT/02-comparison-record.md	repo
docs/validation/evidence/P97-HYDROGEN-CONTROL-PILOT/03-deviation-log.md	repo
docs/validation/evidence/P97-HYDROGEN-CONTROL-PILOT/04-reviews/07-pilot-hydrogen.txt	repo
docs/validation/evidence/P97-HYDROGEN-CONTROL-PILOT/MANIFEST.md	repo
docs/validation/evidence/P97-HYDROGEN-CONTROL-PILOT/SHA256SUMS.txt	repo
```
