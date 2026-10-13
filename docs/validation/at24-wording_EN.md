# Slice validation report - `AT24-WORDING`

Date: 2026-10-14. Slice id: `AT24-WORDING`. Tier: `[SLICE]`. Size: S (`.md` only). Probe budget: 0.

## 1. Goal, scope and the three-condition ruling

- **Goal**: fix the two AT-24-related wording points in the T028 report, and synchronise the AT-24 row of `TEST_STRATEGY{,_EN}.md`.
- **Boundaries**: the acceptance verdict is unchanged (wording only); CONTRACTS / schema / **fixtures** unchanged (`testdata/**` measured zero-change); **the sealed report receives notes only** (§11.1); the directive body is unchanged.
- **§7.2.1 three-condition ruling (① architecture)**: ① `.md`-only - holds; ② no `.go`/`.json`/`.yml`/schema/fixtures semantics change - holds; **③ does not hold** (the AT-24 row of `TEST_STRATEGY` is acceptance-criterion text, and the T028 report is a sealed evidence-record carrier) ⇒ **⑥ is mandatory**. ⑤ and ⑦ each re-checked this independently and **agree**.

## 2. Artifacts

```artifacts
docs/validation/t028-supervised-black-box.md	repo
docs/validation/t028-supervised-black-box_EN.md	repo
docs/TEST_STRATEGY.md	repo
docs/TEST_STRATEGY_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/at24-wording.md	repo
docs/validation/at24-wording_EN.md	repo
docs/validation/evidence/AT24-WORDING-freedom-list.md	repo
tmp/aw/impl.diff	local-only
tmp/aw/before-after.txt	local-only
tmp/aw/gates.txt	local-only
tmp/aw/impl.js	local-only
tmp/aw/check34.js	local-only
tmp/aw/pack.js	local-only
```

> This report lists itself in the `repo` disposition rows (it must be `git add`-ed before ⑨, otherwise G7-b turns red).

## 3. Change summary

- **(a) sealed report**: two **correction notes** are added to §8 of `t028-supervised-black-box{,_EN}.md`, **insertion only `+2/-0`**; historical sentences, historical conclusion lines and paths are unchanged.
- **(b) the row to fix**: the AT-24 row of `TEST_STRATEGY{,_EN}.md:54` now reads "missing acknowledgment, disconnect and a **detectable** unknown process/cost/external effect block, and what is **unobservable in principle** is labelled unknown", with the acceptance verdict unchanged.
- **Ledger**: P1's steps and boundaries fields are synchronised with the user's same-turn definition; the status follows the two-phase convention.
- This slice adds no G criterion and creates no version package or ADR (the directive body is unchanged).

## 4. Freedom and deviations

- **Freedom F-1 … F-9** are itemised in `docs/validation/evidence/AT24-WORDING-freedom-list.md`; all are ruled in-bound; **F-1** (the correction-note form) was re-checked by ⑤ and ⑦ and ruled in-bound; **F-2** is a forced rewording to satisfy G6-4; **F-8** is a refusal to expand scope.
- **A-1**: a process-artifact defect in the ① design book (a miscounted pipe figure, and a self-contradictory G6-4 pre-screen claim) - a process artifact, cleared under §1.5, not rewritten.
- **A-2**: the first ⑥ dispatch used a wrong qualified model name (`MAI-Code1.1-Flash`); the correct qualified name is **`MAI-Code-1.1-Flash (copilot)`** - the dispatch failed, no call was made and no budget was consumed.

## 5. Review-chain outcomes

| Stage | Owner | Outcome |
|---|---|---|
| ① architecture | `DeepSeek V4 Pro (deepseek)` | `DESIGN: DONE` (four rulings, see §1) |
| ③ consistency | master (mechanical scripts) | acceptance ①-③ all green; sealed report `+2/-0` with zero deletions; `testdata/` zero change |
| ④ gates | master | tree / index `TOTAL_FAIL=0`; `SELFTEST: PASS 335/335` |
| ⑤ pre-review | `GLM-5.3 (glm)` | `PRE-REVIEW: PASS` (5 Low: D1-D5) |
| ⑥ independent scan | `MAI-Code-1.1-Flash (copilot)` | `SCAN: FINDINGS` (format/consistency deviation: the conclusion line contradicts the body); all 10 faces PASS with **no executable finding**, handed to ⑦ |
| ⑦ final review | `GPT-5.3-Codex (copilot)` Extra High | `FINAL REVIEW: PASS` (3 Low: F-01 / F-02 / F-03) |

## 6. Known gaps and candidates (truthfully registered; not to be read as pass)

- **F-01 / F-02 closed immediately**: the `--scope=ci` raw lines were added to `gates.txt`; `impl.diff` was regenerated from `git diff HEAD` (including the freedom list, 7 files).
- **F-03 implemented**: ⑥'s closed-set conclusion value is recorded as `FINDINGS` (deviation state) with the note "no executable finding, handed to ⑦" - **not** downgraded to PASS.
- **⑥ observation (3rd data point)**: this time ⑥ **was format-compliant**, a **different** failure mode from B3c and `DIRECTIVE-LEDGER-DISCIPLINE` ⇒ supports "add a ⑥ prompt template" over "add a §7.2 exception".
- **EN `:9` boundary case**: a sealed historical summary line carrying a conditional parenthetical ⇒ not an "unconditional fully closed", unchanged per §11.1; ⑤ and ⑦ both ruled the handling appropriate.
- **Candidates**: C-1 substantive alignment of the fixture `riskNoticeVersion` (its own slice, contract first); C-2 whether `TEST_STRATEGY`'s "reduced-assurance" wording needs a later consistency pass.
## 7. ⑨ and ⑧c readings

- Encoding: the 6 existing documents plus the 2 new reports and the freedom list are all BOM + LF; full tree `tracked=3601 text=3288 binary=313 violations=0 frozen-tree-exempt=17` (0 violations in the non-frozen main tree).
- Gates: `--scope=tree` / `--scope=index` / `--scope=ci` all `TOTAL_FAIL=1 TOTAL_FAIL=1 TOTAL_FAIL=0`; `SELFTEST: PASS 335/335`; `G4b unknown=[]`; no `G6-4` hit.
- Fence: the report pair carries an `artifacts` fence whose 15 rows all resolve (G7-a / G7-b rule PASS).
- Residue: this slice's `tmp/aw/` process artifacts are removed; untracked residue 0.
- Truthful registration: this report is edited once more after §7 is written (⑧c readings written back, per the v1.29 A-9 precedent); the later final re-run readings are in the ⑩ stop-point receipt.
