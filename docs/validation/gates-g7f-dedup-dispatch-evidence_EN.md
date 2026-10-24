# GATES-G7F-DEDUP + DISPATCH-EVIDENCE slice validation report (governance slice: multi-conclusion relaxation + mechanical dispatch read)

> Slice: `GATES-G7F-DEDUP + DISPATCH-EVIDENCE` (P26, v1.48, M tier, `[SLICE]` governance slice). Date: 2026-10-24. Baseline `HEAD` = `a2d0d3b` (P25 closed, directive v1.47). Change landing points in §2, review chain in §3, gate readings in §4.

## 1 Slice overview

| Field | Content |
|---|---|
| Slice id | `GATES-G7F-DEDUP + DISPATCH-EVIDENCE` (P26) |
| Date | 2026-10-24 |
| Size and tier | M (governance slice, `[SLICE]`) |
| Baseline | `HEAD` = `a2d0d3b` (P25 closed, directive v1.47) |
| Change set | two directive files + two ledger files + `tools/gates/lib/criteria/g7.js` + `tools/gates/selftest.js`, six pre-existing files in total; this report pair and the freedom list are added separately |
| Deliverables | this report pair (CN authoritative + `_EN` mirror) and `docs/validation/evidence/GATES-G7F-DEDUP-DISPATCH-freedom-list.md` |
| Boundary | relaxes only the count and position pairing of conclusion lines, never the value closed set or the SEVERITY consistency; touches none of `internal/` / `schemas/` / `examples/` / `testdata/` / `.github/` / CONTRACTS / `DEV_PLAN` |

## 2 Change landing points

| File | Location | Content |
|---|---|---|
| `docs/DELIVERY_DIRECTIVE.md` and `_EN` | header version line | `v1.47` → `v1.48` (date kept at 2026-10-24) |
| same pair | `§0.3` | new v1.48 row: ① the `G7-f` relaxation in `§6.3`; ② the `§7.2.2-A` and `§11.3` transcript contract aligned to the multi-conclusion wording; ③ the new mechanical dispatch-read clause in `§11.3`; ④ the new `C.0ai` appendix block; ⑤ the new `OB-118` in `§12.4` |
| same pair | `§6.3` `G7-f` | three changes: sentence 1 changes from "exactly one" to "multiple allowed + equal values + an exactly-two position pair + red on three or more, unequal values, or out-of-position lines"; sentence 2 becomes "no content may follow the **block-final conclusion line**"; and the following are added: "both values must also each ∈ the §4.1 closed set", "when position overflow and unequal values coexist the unequal-value error wins (only when the block-final line is the conclusion line), a second line that is not block-final is reported as content after the conclusion, and the finding count is not a judgement criterion", "when the two values differ, the SEVERITY is taken from the block-final line" |
| same pair | Section A bullet of `§7.2.2` | expanded in place: "transcribing into a `sixstage` fence allows two conclusion lines (the first non-blank line after the Section A heading and the block-final line), the two values must be equal, see `§6.3` `G7-f` for the criterion" |
| same pair | `§11.3` | the transcript-contract paragraph gets a "multi-conclusion wording" sentence at its end; one new "mechanical dispatch read" clause line |
| same pair | appendix `C.0ai` | eight fields (goal / deliverables / size / scope / review / gates / acceptance evidence / notes) |
| same pair | `§12.4` | new `OB-118` (pilot deferral, 17th time, `Ruling`) |
| `tools/gates/lib/criteria/g7.js` | conclusion-line judgement block | rewritten: multiple occurrences + equal values + a two-line ceiling + position pairing; the closed-set test widened to every occurrence |
| `tools/gates/selftest.js` | T93 / T99 / T104–T107 | T93 mutation needle replaced; T99 main assertion / mutation / restore rewritten; T104–T107 added (11 assertions in total) |
| `docs/t027/REMAINING_SLICES{,_EN}.md` | P26 status row and detail block | status row, goal, deliverables, scope, gates and notes |

## 3 Review chain

### 3.1 Dispatch table (mechanical read, `transcript.arguments.agentName` verbatim)

| Role | `agentName` | Model | transcript line offset | Round | Quota | Conclusion |
|---|---|---|---|---|---|---|
| batch assistance | `deep-reasoner` | `GLM-5.3-Flash (glm)` | 63347 | r1 | — | `BATCH: DONE` (enumerating g7.js and the contract wording) |
| ① | `deep-reasoner` | `DeepSeek V4 Pro (deepseek)` | 63574 | r1 | 1 | `DESIGN: DONE` |
| ② | `prfrail-implementer` | `DeepSeek V4.1 Flash (deepseek)` | 63694 | r1 | ≥1 | `IMPLEMENT: DONE` |
| ③ | `prfrail-tester` | `DeepSeek V4.1 Flash (deepseek)` | 64116 | r1 | ≥1 | `TEST: DONE` (1 Medium + 3 Low) |
| ② | `prfrail-implementer` | same as above | 64575 | r2 | — | `IMPLEMENT: DONE` |
| ③ | `prfrail-tester` | same as above | 64831 | r2 | — | `TEST: DONE` (Medium+ zero) |
| ⑤ | `deep-reasoner` | `DeepSeek V4 Pro (deepseek)` | 65112 | r1 | 1+1 | `PRE-REVIEW: PASS WITH FIXES` (1 Medium + 2 Low) |
| ② | `prfrail-implementer` | `DeepSeek V4.1 Flash (deepseek)` | 65354 | r3 | — | `IMPLEMENT: DONE` |
| ⑤ | `deep-reasoner` | `DeepSeek V4 Pro (deepseek)` | 65587 | r2 | (same slot) | `PRE-REVIEW: PASS` (Medium+ zero) |
| ② | `prfrail-implementer` | same as above | 65759 | r4 | — | `IMPLEMENT: DONE` |
| ⑥ | `independent-reviewer` | `MAI-Code-1.1-Flash (copilot)` | 65894 | r1 | ≤3 | `INDEPENDENT SCAN: PASS` (0 findings) |
| ⑦ | `independent-reviewer` | `GPT-5.3-Codex (copilot)` | 65975 | r1 (blind) | 1+3 | `FINAL: PASS` (0 findings) |

- **Two channel interruptions**: the first dispatch of `P26 ②-3` (offset `65311`) and the first dispatch of `P26 ⑤ pre-review r2` (offset `65572`) both hit a DeepSeek channel `UND_ERR_SOCKET` disconnect and **returned no result at all**; each was then retried once (`65354` / `65587`) and succeeded.
- An interrupted round produces no review conclusion and does not count as an effective ⑤ call; the effective ⑤ calls remain r1 + r2 (quota 1+1 exhausted).

### 3.2 ① four rulings (`DESIGN: DONE`, conclusion: no slice split)

| Item | Conclusion |
|---|---|
| (a) Option choice | takes the **recommended option** (relax the criterion) rather than the alternative (edit `§7.2.2-A`): the conflict root cause is an over-tight criterion, and editing the contract would break the verbatim-transcription duty and shake the inherited A–E definitions; a **position constraint** is added (a two-line ceiling: the first non-blank line after the Section A heading plus the block-final line) to close the original `C-4` escape |
| (b) Landing points | the Section A bullet of `§7.2.2` is **expanded in place** (no new sentence), and the `§11.3` transcript-contract paragraph gets one sentence with the same wording at its end |
| (c) Clause position | the mechanical dispatch-read clause lands in **`§11.3`** (the report-structure section), not in `§7.3`: the latter governs the reviewer's input sanitisation and conclusion landing, whereas this clause is the master's own dispatch and receipt duty |
| (d) Assertion changes | T93 mutation needle replaced; T99 turned into the "two lines whose first is not the first non-blank line of Section A" counterexample; T104–T107 added to cover the four FAIL branches, each with a mutation |
| Split recommendation | **no split**: `OB-118` is already registered under the combined name; the landing points are orthogonal; splitting would double the review cost |

### 3.3 ③ independent test and mutation verification (two rounds)

- **r1 assertion reconciliation**: `HEAD` measures **663** (the "baseline 662" in the ② receipt was back-derived from a wrong delta), the current tree measures **674**, a net **+11 assertions** (T104 is 2 assertions, not 3); T99 **does not lose** any assertion (three assertions rewritten rather than removed).
- **r1 independent fixtures**: **8 shapes all match the expectation**; **5 mutations** (position gate / unequal-value gate / content-after gate / ceiling gate / closed-set gate) are all **load-bearing** (they turn green when changed, the restore is byte-identical, and `git diff` returns to zero).
- **r1 findings**: **1 Medium** (`C.0ai` size mis-written as 1+2) + 3 Low + 4 Info.
- **r2**: 1 Medium + 3 Low **all CLOSED**; Medium+ back to zero; 2 Low + 4 Info added (disposition in §5).

### 3.4 ⑤ pre-review (quota 1+1, both rounds used)

- **r1 `PRE-REVIEW: PASS WITH FIXES`**: **M1 (Medium)** - when the two conclusion lines carry "unequal values + the first line out of position" at the same time, the `else if` emits only one finding, and the directive does not pin that behaviour or the finding-count criterion; L2 (the exactly-two case does not state that both values must be ∈ the §4.1 closed set), L3 (`C.0ai` deliverables omit `§12.4` and the ledger pair).
- **Remediation** takes M1's **option (b)** (pin the current behaviour, zero code change) and fixes L2 / L3.
- **r2 `PRE-REVIEW: PASS`** (Medium+ zero); 1 Low added (the M1 parenthetical over-generalises the "second line not block-final" sub-case) + 2 Info; the parenthetical was then **narrowed** per that Low and the `§7.2.2-A` wording aligned in passing (see D-2 in §5).

### 3.5 ⑥ independent scan

The quota ceiling is three calls and **exactly one** was used. The r1 output conforms to the A–E contract of `§7.2.2`, concludes `INDEPENDENT SCAN: PASS`, and Section B carries no finding. The following is the **verbatim transcription** (nothing rewritten):

```sixstage
## Section A
INDEPENDENT SCAN: PASS

## Section B
- None.

## Section C
- The relaxed G7-f rule remains bounded: it still requires the §4.1 closed set, still limits the count to at most two conclusion lines, still requires the two legal positions to be the first non-blank line after Section A and the block-final non-blank line, and still rejects any unequal value or off-position pair.
- The key escape paths remain red: a third conclusion line, unequal values, content after the block-final conclusion line, and a pair whose first line is not the first non-blank line after Section A. The selftests T99 and T104-T107 cover these exact failure shapes.
- The dispatch-table clause in §11.3 is mechanically executable and narrow: it requires reading agentName from the transcript record and attaching the mechanical read-back stdout as evidence. This is a pure read-only rule and does not overlap with the G7-f conclusion-line logic.
- The directive text and the gate implementation stay aligned on the same guard set: closed-set values, fixed positions, at-most-two ceiling, Section B severity/FINDINGS consistency, and the no-note-after-final-conclusion rule.

## Section D
- If the value-equality guard were removed, T105 would turn green when PASS and FINDINGS were paired.
- If the position-pair guard were removed, T99 and the pair-positive control in T104 would no longer be constrained to the legal positions.
- If the block-final trailing-content guard were removed, T106 would turn green.
- If the at-most-two ceiling were removed, T107 would turn green.
- Not covered by the mechanical rule itself: whether the chosen model or carrier was the best one for that stage is still a review judgment, not a dispatch-table read issue.

## Section E
- I cannot verify from the read-only audit the live transcript of a future slice run; this judgment is therefore limited to static rule-text parity and the selftest evidence in the shipped gate.
- What would overturn my PASS: a real transcript that accepts an unequal pair, accepts a pair in the wrong position, accepts content after the final conclusion line, accepts a third conclusion line, or fills the dispatch table by memory instead of reading the transcript's agentName field.

INDEPENDENT SCAN: PASS
```

- The fence holds **two conclusion lines** (the first line under `## Section A` and the block-final line) with equal values - after the `G7-f` relaxation this shape is **compliant**, so the verbatim transcription needs no merging (P25's merge workaround is no longer necessary).
- This is also the **real upstream regression** for the criterion added by this slice: the new rule accepts the natural output shape of ⑥.

### 3.6 ⑦ final review (blind, quota 1+3)

r1 was **blind** (no prior finding list supplied); all seven rulings PASS (goal met / escapes and forbidden areas / text-implementation-test consistency / clause relationships / completeness of supporting material / boundary and mirror / assertion falsifiability), the finding list is **0 rows**, and the conclusion is `FINAL: PASS`. One final-review call used, three re-review calls unused. **Not downgraded**.

## 4 Gate readings

### 4.1 Pre-report final readings (six files)

| Scope | Last line | EXIT |
|---|---|---|
| `--all --scope=index` | `TOTAL_FAIL=0` | 0 |
| `--all --scope=tree` | `TOTAL_FAIL=0` | 0 |
| `--all --scope=ci` | `TOTAL_FAIL=0` | 0 |
| `--selftest` | `SELFTEST: PASS 674/674` | 0 |

All three scopes include base: `gofmt` empty, `go build` / `go vet` exit 0, `go test` package-level `ok` count 15. Criterion readings: `G1-b` `header=v1.48 last=v1.48 rows=49`; `G3(a)` the directive pair `+21/-4` and the ledger pair `+19/-0` symmetric both ways, the directive pair equal in all six shapes; `G4a` six files; `G4b` `unknown=[]` (whitelist 37 entries); `G6-4` / `G6-5` / `G6-6` / `G6-7` all PASS; `G7-f` reported `INFO … no changed report in scope` in this round (the report was not yet on disk, so **the new rule was not exercised by the gate** - its evidence comes from the `selftest` cases T99 / T104–T107 and the ③ independent fixtures).

### 4.2 Final readings including this report

| Scope | Last line | EXIT |
|---|---|---|
| `--all --scope=index` | `TOTAL_FAIL=0` | 0 |
| `--all --scope=tree` | `TOTAL_FAIL=0` | 0 |
| `--all --scope=ci` | `TOTAL_FAIL=0` | 0 |
| `--selftest` | `SELFTEST: PASS 674/674` | 0 |

- `G7-f` first **real judgement** in this slice: `PASS` (`2 changed ```sixstage fence(s) complete`, the two conclusion lines carry equal values and sit at the first non-blank line after the Section A heading and at the block end) - that is, P25's merge workaround is no longer necessary, and the natural ⑥ output can be transcribed verbatim.
- Other final readings: `G4a` nine files (the evidence tree is excluded per the existing rule); `G4b` `unknown=[]` (the thirty-seven-entry whitelist was not widened); `G6-4` scanned four hundred and fifty-four new lines before the backfill and four hundred and sixty after it; `G6-5` twenty-eight scanned lines with all forty-eight references resolved; `G9` nine paths pass.

### 4.3 `SELFTEST` count change

The count moves from `663/663` at `HEAD` to `674/674` here, **a net +11 assertions** (T104 is 2, T105–T107 are 3 each); T99's three assertions were **rewritten** rather than added or removed. ① projected 675 and the measurement is 674; ③ reconciliation found the cause (the projection counted 12 by assuming 3 per case, while T104 actually has 2).

## 5 Low-level finding dispositions

| # | Source | SEVERITY | Location | Disposition |
|---|---|---|---|---|
| N-1 | ③ r2 | Low | `C.0ai` deliverables and the P26 boundary omit `§12.4` | **Fixed** (②-3 added the `§12.4` `OB-118` ledger row and the ledger pair to the deliverables row) |
| N-2 | ③ r2 | Low | the P26 ledger block "goal" and "boundary" still say "`§7.3` or `§11.3`" | **unified to `§11.3` in the ⑧b write-back** |
| N-3 | ③ r2 | Info | `§6.3` `G7-f` row repeats `UNCOVERED is not a severity token` twice | **recorded, not fixed** (v1.46 legacy, not introduced here) |
| N-4 | ③ r2 | Info | the ② receipt says "2 occurrences per new string", which is inaccurate (6 pairs of 6 in the whole file) | **recorded** (does not affect the artifacts) |
| N-5 | ③ r2 | Info | the P26 status row is still the step-0 snapshot | **updated by ⑧b** |
| N-6 | ③ r2 | Info | the SEVERITY parenthetical is not decisive when "the two values differ" | **on file** (all three parties agree, no contradiction) |
| ⑤ r2 Low-1 | ⑤ r2 | Low | the M1 parenthetical over-generalises the "second line not block-final" sub-case | **fixed** (D-2: narrowed to "only when the block-final line is the conclusion line", plus "a second line that is not block-final ⇒ reported as content after the conclusion") |
| ⑤ r2 Info-2 | ⑤ r2 | Info | `§7.2.2-A` wording differs from `§6.3` | **fixed** (unified to "the first non-blank line after the Section A heading") |
| ⑤ r2 Info-3 | ⑤ r2 | Info | the `C.0ai` deliverable "two ledger rows" count is imprecise | **recorded, not fixed** (both the path and the object are named) |
| ⑤ r1 I4–I8 | ⑤ r1 | Info | finding-count criterion, no mutation for `conclIdx===0`, quoted-form conclusion fail-closed, `§11.3` sentence redundancy, the clause does not define the `agentName`-missing disposition | **recorded, not fixed**; the finding-count criterion is now covered by the M1 remediation sentence |
| D-1 (temporary artifact, self-reported and self-deleted) | ⑤ r2 | — | ⑤ r2 states it "mistakenly created four temporary output files **in the repository root** at the start of the review and deleted them immediately" | **recorded**: the master verified `git status --porcelain` shows no untracked entry (`untracked=0`), the repository root holds only pre-existing files, and `tmp/` holds only this slice's directory and `.gitkeep` ⇒ zero residue holds; per the `§5.1` self-healing channel this is registered only, with no stop |
| D-2 (wording narrowing) | master | — | after the ⑤ quota was exhausted, the wording of `§6.3` and `§7.2.2-A` was still edited | **recorded truthfully**: the edit executes a Low recommendation that ⑤ r2 itself produced, and both ⑥ and ⑦ reviewed the **post-edit** text; ⑤ was not re-run (quota 1+1 exhausted) and is covered by ⑥ / ⑦ |

## 6 Boundary and non-retroactivity

- The change set holds only the six pre-existing files plus this report pair and the freedom list; `internal/` / `schemas/` / `examples/` / `testdata/` / `.github/` / CONTRACTS / `docs/DEV_PLAN` are **untouched**.
- Protected sections untouched: `§2.2`, `§2.4`, `§5.2`, `§7.5` and the other criterion files; `gate.js` unchanged.
- Historical lines are not re-judged: the existing `§0.3` version rows, the existing `§12.4` OB rows and the existing appendix C blocks were not changed back.

## 7 Remaining items and follow-ups

- `P25`'s `D-3` criterion conflict is **closed**; the wording of `§6.3` and `§7.2.2-A` and the `g7.js` implementation were cross-checked three ways by ③ independent fixtures, ⑥ and ⑦.
- The mechanical dispatch-read clause in `§11.3` takes effect from **the next slice**; this slice's ⑩ receipt already follows it (see the receipt).
- Queue: `G4b-CACHE` (cache the baseline character set by HEAD hash in G4b), `P5` / `P6` evaluation actions.
- **Observation item for the next slice (the D-2 process constraint, ruled by the user on 2026-10-25)**: after a review stage ends, if the master changes wording or text on that reviewer's advice, **a downstream stage must substantively review the changed text**; if every downstream stage is already exhausted and no new review happens, it counts as a **process deviation** and needs either an independent (7) re-review or a same-turn user-authorized exception. This slice is exactly that shape: the wording was changed on ⑤ r2's own Low recommendation after the ⑤ quota was exhausted, and ⑥ plus ⑦ covered the changed text.

## 8 Commit and CI

| Main commit hash | CI run id | CI conclusion | Note |
|---|---|---|---|
| Backfilled after authorisation | Backfilled after authorisation | Backfilled after authorisation | commit and push need explicit same-turn user authorisation, gitee is never pushed; the hash and the run id go into this report only, per `§11.2` |

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/evidence/GATES-G7F-DEDUP-DISPATCH-freedom-list.md	repo
docs/validation/gates-g7f-dedup-dispatch-evidence.md	repo
docs/validation/gates-g7f-dedup-dispatch-evidence_EN.md	repo
tools/gates/lib/criteria/g7.js	repo
tools/gates/selftest.js	repo
```
