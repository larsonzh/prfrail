# Slice validation report · P66 CHECKLIST-RULES-HOIST

- Directive version: v1.65 (step-0 point in time; step G raised it to v1.66)
- Date: 2026-11-11
- Class: governance `[SLICE]`; tier L; probe budget 0 (documentation-only, local machine)
- Status: **✅ Done** (closed 2026-11-11; stopped at ⑩ awaiting authorization)
- Controller model: `DeepSeek V4.1 Flash` (not switched to Pro; no ADR per briefing §0.0 item 5)
- Evidence pack: `docs/validation/evidence/P66-CHECKLIST-RULES-HOIST/`
- Landing: five new directive subsections (§13.1–§13.5); from v54 the snapshot keeps point-in-time data only

## 1. Execution summary

One-line conclusion: the execution discipline long frozen in v53 §7 / §8 / §9 has been structurally lifted into directive §13 under the "migrate everything" option; v53 stays frozen and, from v54, the snapshot keeps point-in-time data only.

Key readings:

- §13 subsection item counts: 13.1 = 20, 13.2 = 26, 13.3 = 25, 13.4 = 26, 13.5 in three parts 8 / 4 / 5.
- Gate pitfalls 1–26 are numbered contiguously; symptoms and dispositions are kept verbatim (only 3 registered rewrites, see §5 and §9).
- CN/EN symmetry: `git diff --numstat` = `184/3` vs `184/3`.
- Gates: three scopes `TOTAL_FAIL=0`; `SELFTEST: PASS 722/722`.
- R count: R = 3 (§9.2 definition; the per-round breakdown is in §11).

## 2. Three-option comparison and final landing

| Option | Claim | ① verdict | Rationale |
|---|---|---|---|
| A migrate everything | all 98 source lines analysed; the 82 merged rules each get a disposition (migrate / repoint / do not migrate) | **adopted** | complete and auditable snapshot→directive cross-reference; a snapshot reader can find every clause destination in §13 |
| B migrate §7 §8 only | lift only the startup general rules and the collaboration discipline; checklist meta-rules and gate pitfalls stay in the snapshot | not adopted | gate pitfalls and checklist meta-rules would still sit in the snapshot, so the inflation is unsolved |
| C migrate gate pitfalls only | move only the 1–26 pitfall table into the directive | not adopted | the other discipline stays a soft constraint, which misses the "promote to hard constraint" goal |

① verdict: **agree with A**. Final landing = five new directive subsections (§13.1–§13.5) + snapshot slimming from v54 (§7 / §8 / §9 become one pointer line plus a section-name cross-reference).

## 3. Dedup analysis

Rule-level denominator (quotation from `01-arch/dedup-analysis.md`): 82 rules in total (source file line count 98, merged rule count 82). The 14 merges cover multi-line descriptions of the same rule; the 2 non-rule structural rows (pitfall-table header and version note) stay outside the three-way classification.

| Class | Count | Handling rule |
|---|---|---|
| Exact duplicates | 19 | repoint (§13 keeps only the pointer line, no entity copy) |
| Partial overlap | 43 | migrate (standalone entry + pointer) |
| Snapshot-only | 20 | migrate (full body) |

The per-rule table (82 items with rule id / source position / first sentence / target section / relation / remark) is in `01-arch/dedup-analysis.md` §1; this report keeps only the summary table and does not copy the rows.

## 4. Migration mapping

| Landing state | Count | Composition |
|---|---|---|
| Migrated | 70 | partial overlap 43 + snapshot-only 20 + pitfall exact duplicates 7 (verbatim into the pitfall table, with pointer annotations) |
| Repointed | 11 | exact duplicates outside the pitfall table, discipline Y excepted |
| Not migrated | 1 | the discipline-Y note (promoted to §1.7.2 by P59; the note is itself a pointer) |
| Non-rule rows | 2 | pitfall-table header, version note (outside the three-way classification) |

One extra registration `G-0`: the quotation general-rule sentence before the §9 table (not one of the 24 table rows) → `§13.3.1`, migrated verbatim. The per-rule mapping and the section-retarget record are in `01-arch/migration-mapping.md`.

## 5. Step-by-step migration

| Step | Content | `git diff --numstat` (CN = EN) | Landing |
|---|---|---|---|
| A | §13 chapter frame (chapter lead + five empty subsections, zero rule entities) | `24/0` | §13 chapter head and the §13.1–§13.5 headings |
| B | v53 §7, 20 items | `21/0` | §13.1, twenty items |
| C | v53 §8 table rows + discipline-Y note + GLM policy paragraph | `27/0` | §13.2, twenty-six items |
| D | v53 §9 checklist meta-rules, 24 table rows + the quotation general rule | `26/0` | §13.3, twenty-five items |
| E | v53 §9 gate pitfalls 1–26 | `33/0` | §13.4 three-column pitfall table + the source-grouping and version notes |
| F | snapshot structure template | `34/0` | §13.5 three parts (8 kept sections / 4 pointer sections / 5 point-in-time classes) |
| G | metadata, 17 items: header 1 + §0.3 1 + C.0at 13 + B.7.1 1 + B.7.4 1 | metadata step | directive header, §0.3 changelog row, Appendix C.0at, Appendix B.7.1 / B.7.4 |

Each step numstat is the reading recorded in that step check file (`03-test/step-*-check.md`); step G touches only metadata and slot pointers and adds no rule entity, so its items fold into the final cumulative reading `184/3` vs `184/3`.

## 6. New §13 chapter summary

| Section | Responsibility | Item count |
|---|---|---|
| §13.1 Startup and step-0 general rules | anchor check, step-0 registration and startup declarations | 20 (13.1.1–13.1.20) |
| §13.2 Cross-slice collaboration discipline | role/carrier/model binding, mechanical dispatch-table reads, commit and temp-file discipline | 26 (13.2.1–13.2.26) |
| §13.3 Checklist meta-rules | citation form, regeneration timing, wording boundary | 25 (13.3.1–13.3.25) |
| §13.4 Gate pitfalls | incident register: id / symptom / disposition, 1–26 contiguous | 26 (including No. 26 from F2-FIX) |
| §13.5 Snapshot structure template | kept sections, pointer sections and point-in-time scope from v54 | three parts 8 / 4 / 5 |

## 7. §13.5 snapshot structure template summary

- 8 kept sections: §1 current anchors, §2 fact pack, §3 completed slices, §4 execution queue, §5 observations, §6 awaiting user ruling, §10 opening sentence, §11 one-line status.
- 4 pointer sections: §7 → §13.1, §8 → §13.2, §9 checklist meta-rules → §13.3, §9 gate pitfalls → §13.4 (one pointer line plus a section-name cross-reference, no entity copy).
- 5 point-in-time classes: anchor numbers / current fact pack / completed-slice list / queue and observation state / rule text (rule text never enters the snapshot and lives only in §13).

## 8. Known boundaries

- Semantic change after promotion — the rules migrated into §13 are hard constraints from the next slice on, not snapshot experience only.
- §13.4 is an incident register: symptoms and dispositions follow the incident wording; where it points at the same fact as the body or another section (No. 11 with §13.1.12, No. 21 with §9.3, No. 8 with §13.1.19) the body or the existing clause is the normative source, and the table sets up no second normative source.
- A-part: the v53 §9 item-26 candidate is now formalized in §13.4, which closes the P65 ledger A-part carry-over row.
- v53 stays frozen and is not rewritten (§11.1): v50–v53 are untouched and this slice changes no snapshot body.
- ⑥ formal defect downgraded: the round-3 headings did not start with `Section A`, which is a format deviation under §4.1 ⇒ the conclusion is downgraded to low-weight corroboration (see §11).
- EN-side formatting compensation: pitfall entries 16 / 19 / 20 carry bold compensation to satisfy the G3(a.2) six shapes and are not verbatim against the EN snapshot (the CN side is unaffected; recorded in `03-test/step-E-check.md`).
- Pointer-precision residue: the narrow/wide pointer wording in 13.3.5 → §9.3 and 13.3.16 → §1.6 was judged Low and accepted by ⑤, carried to a later slice.
- Rule text lives in one place only (§11.2): from v54 the snapshot replaces §7 / §8 / §9 bodies with pointers.

## 9. Explicitly not executed

- No change to the v50–v53 historical snapshots (including the v53 §7 / §8 / §9 text); v54 is not generated in this slice.
- No change to the criterion body `tools/gates/**`, the CJK whitelist, CONTRACTS / schema / fixtures, product code or CI.
- No change to directive §0–§12 body text; only the §0.3 new changelog row, the header version line, and the two convergence sentences inside §13 (the 13.1.9 document-type degradation pointer and the §13.4 register-scope note).
- ⑥ round 4 was not re-run (allowance exhausted; restoring validity needs a later slice to run a contract-conforming independent scan).
- No commit / push; no VM / proxy / Sol double-run; no new probes.
- `tmp/P66-CHECKLIST-RULES-HOIST/` was not cleaned (⑧c does that); no new temp artefacts.

## 10. Stage transcript (sixstage fence A–E)

Note outside the fence: the five A–E blocks below are the round-3 ⑥ transcript, with headings normalized per the §4.1 contract (the original headings were `## A. Scope` and so on); the weight of ⑥ is downgraded to corroboration per §11.

```sixstage
## Section A

I re-read the amended text in [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md) and [docs/DELIVERY_DIRECTIVE_EN.md](docs/DELIVERY_DIRECTIVE_EN.md), together with §5.2, §9.3 and §11.2, and checked whether the earlier two Medium findings were actually resolved in substance and whether the two Low findings remained acceptable.

## Section B

- None.
- The two earlier Medium findings are materially closed.
- The two Low findings remain within an acceptable range.

## Section C

- In [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md), 13.1.9 states that documentation-type slices degrade to prfrail-documenter per §5.2, and that this clause "sets up no further exception." This is the substantive fix: the governance rule remains in §5.2, and no new exception is created.
- In [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md), 13.4 explicitly says the pitfall table is an incident register and that the body text or existing section remains the normative source under §11.2, with "no second normative source." That directly removes the earlier authority conflict.
- In [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md), 13.3.5 explicitly says "A rule obligation ≠ a gate obligation," and points to §6.3 and §9.3. That is an acceptable boundary and matches the low-item intent.
- In [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md), 13.3.16 says a [ROUTINE] small fix takes no new slice number but loses no discipline, and points to §1.6 for the ring count. That is a pointer fix, not a semantic rewrite, and it is consistent with the original low-item concern being addressed by precision rather than by creating new rule semantics.
- The English mirror in [docs/DELIVERY_DIRECTIVE_EN.md](docs/DELIVERY_DIRECTIVE_EN.md) says the same thing, preserving the authority chain and the no-second-source rule.
- The direct review confirms the 13.1–13.5 sections are present and that the pitfalls list remains a 26-item register, not a second authority layer.

## Section D

- The critical authority question is not whether a note says "no second source," but whether the authority chain is now actually fault-free. It is: §11.2 remains the source-of-truth rule, §5.2 remains the controlling role matrix, and 13.4 explicitly says the register is not a second normative source.
- The earlier Medium 1 is therefore substantively closed: the table is a register, not an authority stream.
- The earlier Medium 2 is also substantively closed: 13.1.9 does not create a new exception; it explicitly says the document-type role degradation remains under §5.2.
- The Low items remain acceptable because they are clarifying boundaries and pointers, not new gate meaning or new rule authority.

## Section E

- No files were modified and no commit or push was made.
- This was a read-only audit of the governing text, and the read evidence supports closure of the earlier Mediums without any remaining High/Medium findings.

INDEPENDENT SCAN: PASS
```

## 11. Review and scan record

- ⑤ pre-review: round 1 raised 2 items to fix (F-1 missing section-retarget record, F-2 quotation general rule outside the denominator) ⇒ after the fix, round 2 returned `RE-REVIEW: PASS` (F-1 / F-2 closed; F-3 / F-4 / F-5 kept as accepted Low items).
- ⑥ independent scan, 3 rounds dispatched: r1 used custom section names (`### Scope` / `### Validation` / `### Findings` / `### Final result`) ⇒ format deviation, void; r2 had a valid A–E shape but concluded `FINDINGS` (2 Medium: the 13.1.9 literal conflict with §5.2 and the §13.4 normative-source ambiguity), later confirmed by the ③-r2 convergence sentences to be stale readings; r3 concluded PASS but headed the blocks `## A. Scope` instead of `## Section A`, still a format deviation under §4.1 ⇒ allowance exhausted plus a formal defect, so the conclusion is downgraded to low-weight corroboration.
- ⑦ final review r1: 6 verdicts (5 PASS + 1 ISSUE) with finding [High] F-01 (the ⑥ PASS cannot serve as a valid pass conclusion) and the conclusion line `FINAL REVIEW: FINDINGS`; the remediation item is to register "⑥ allowance exhausted + formal defect" in this report and the ledger and to downgrade the ⑥ conclusion.
- Closure basis: the ⑦ verdict plus the ④ gates, with ⑥ as low-weight corroboration only (the F-01 remediation item is exactly this section and §12).
- R count (§9.2): the ③ r1 → r2 fix, the ⑥ r2 FINDINGS → r3 fix and the ⑦ r1 F-01 remediation, one each ⇒ R = 3; the void r1 re-dispatch is not counted in R.
- Post-assembly re-check: the changes touch only the reports and the ledger; the rule and structure criteria (G1 / G3 / G4b / G6-4 / G7-e / G7-f) are green, and only the G5(d) fence rows and the G8-b MANIFEST / SHA256SUMS recomputation remain, to be resynced by the controller before ⑨.

### 11.1 Observation candidate (empirical from P66, outside this slice scope) - goes into the v54 fact pack

- **⑥ independent-scan output compliance rate**: all three ⑥ rounds in this slice carried a formal defect (r1 did not use the A-E shape, so it was voided; r2 used A-E but two Medium findings turned out to be stale readings; r3 used ## A. Scope instead of ## Section A), so after the quota was exhausted it was downgraded to low-weight corroboration.
- **Impact**: the ⑥ layer effectiveness failed for this slice; ⑦ caught it twice (final FINDINGS plus re-review PASS), but the ⑥ independent signal never materialised.
- **Candidate disposition**: at v54 or a later governance slice, evaluate whether the ⑥ prompt template needs hardening, or whether a ⑥ format deviation should trigger an automatic re-dispatch rather than consume quota.
- **This slice stance**: not handled here; registered for the v54 fact pack only (user ⑩ ruling 1).

## 12. Boundary and freedom register

| Id | Item | Triggered | Basis / disposition |
|---|---|---|---|
| D-1 | step A run by the controller (no `prfrail-documenter` dispatch) | yes | pure skeleton with zero rule entities; independence impact = low; not continued in B–G (`03-test/step-A-check.md`) |
| D-2 | the two wording-level ③-r2 fixes run by the controller | yes | one sentence and one note each, no new rule entity; independence impact = low; independently re-checked by ⑤ (`03-test/consistency-r2.md`) |
| F-1 | no self-commit / self-push | no | stopped at ⑩ for same-turn authorization |
| F-2 | no new slice, no ruling on the user behalf | no | ① recommends, the user ruled for A |
| F-3 | no VM / proxy / Sol double-run | no | documentation-only, local machine |
| F-4 | no criterion-body or contract change | no | G4b hit fixed by wording; whitelist untouched |
| F-5 | no change to the §0–§12 body | no | only the §0.3 changelog row and the header version line |
| F-6 | no change to v50–v53 | no | v53 stays frozen |
| F-7 | pitfall entries: structure and position only | no | 3 registered rewrites; symptoms and dispositions verbatim |
| F-8 | model binding and serial dispatch | followed | every dispatch carried `agentName` and `model`; strictly serial |
| F-9 | allowance ⑤ 1+1 / ⑥ 3 / ⑦ 1+3 | used up | ⑤ two rounds, ⑥ three rounds, ⑦ one round |
| F-10 | single-output ceiling | batched as instructed | batch-assist receipts report the actual character count |
| F-11 | live ledger read-only | exception under the ⑧ authorization | only the P66 block edited; all other blocks read-only |
| F-12 | step-0 same-round fence declaration | done | report skeleton and evidence pack fenced in the same round |
| F-13 | the skeleton-era fence-waiver declaration replaced by the real fence | done | the waiver line was deleted and the `sixstage` fence landed in this round |
| Hit 1 | G4b new CJK character | handled | the §13.3 lead was reworded (wording boundary); whitelist untouched |
| Hit 2 | G6-4 aggregate-count wording | handled | the §13.5 footnote was reworded to the cumulative item wording |

Two hits in total (one G4b, one G6-4), both closed by wording changes with the criterion body untouched.

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/checklist-rules-hoist.md	repo
docs/validation/checklist-rules-hoist_EN.md	repo
docs/validation/evidence/P66-CHECKLIST-RULES-HOIST/00-briefing.md	repo
docs/validation/evidence/P66-CHECKLIST-RULES-HOIST/01-arch/decision-list.md	repo
docs/validation/evidence/P66-CHECKLIST-RULES-HOIST/01-arch/dedup-analysis.md	repo
docs/validation/evidence/P66-CHECKLIST-RULES-HOIST/01-arch/design-r1.md	repo
docs/validation/evidence/P66-CHECKLIST-RULES-HOIST/01-arch/migration-mapping.md	repo
docs/validation/evidence/P66-CHECKLIST-RULES-HOIST/02-batch/glm-dedup-r1.txt	repo
docs/validation/evidence/P66-CHECKLIST-RULES-HOIST/03-test/consistency-r1.md	repo
docs/validation/evidence/P66-CHECKLIST-RULES-HOIST/03-test/consistency-r2.md	repo
docs/validation/evidence/P66-CHECKLIST-RULES-HOIST/03-test/step-A-check.md	repo
docs/validation/evidence/P66-CHECKLIST-RULES-HOIST/03-test/step-B-check.md	repo
docs/validation/evidence/P66-CHECKLIST-RULES-HOIST/03-test/step-C-check.md	repo
docs/validation/evidence/P66-CHECKLIST-RULES-HOIST/03-test/step-D-check.md	repo
docs/validation/evidence/P66-CHECKLIST-RULES-HOIST/03-test/step-E-check.md	repo
docs/validation/evidence/P66-CHECKLIST-RULES-HOIST/03-test/step-F-check.md	repo
docs/validation/evidence/P66-CHECKLIST-RULES-HOIST/04-reviews/05-prereview-r1.md	repo
docs/validation/evidence/P66-CHECKLIST-RULES-HOIST/04-reviews/05-prereview-r2.md	repo
docs/validation/evidence/P66-CHECKLIST-RULES-HOIST/04-reviews/06-scan-r1.md	repo
docs/validation/evidence/P66-CHECKLIST-RULES-HOIST/04-reviews/06-scan-r2.md	repo
docs/validation/evidence/P66-CHECKLIST-RULES-HOIST/04-reviews/06-scan-r3.md	repo
docs/validation/evidence/P66-CHECKLIST-RULES-HOIST/04-reviews/07-final-r1.md	repo
docs/validation/evidence/P66-CHECKLIST-RULES-HOIST/04-reviews/07-rereview-r1.md	repo
docs/validation/evidence/P66-CHECKLIST-RULES-HOIST/06-gates/gates-r1.txt	repo
docs/validation/evidence/P66-CHECKLIST-RULES-HOIST/06-gates/gates-r2.txt	repo
docs/validation/evidence/P66-CHECKLIST-RULES-HOIST/MANIFEST.md	repo
docs/validation/evidence/P66-CHECKLIST-RULES-HOIST/SHA256SUMS.txt	repo
docs/validation/evidence/P66-CHECKLIST-RULES-HOIST/freedom-list.md	repo
docs/validation/evidence/P66-CHECKLIST-RULES-HOIST/input-list.md	repo
```
