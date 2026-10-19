# Slice DIRECTIVE-EVIDENCE-WINDOW validation report (directive v1.39, MINOR)

> Slice: `DIRECTIVE-EVIDENCE-WINDOW` (working name `G9-d`); directive version v1.39; tier M; class `[SLICE]`; this file is the **⑧b final pass** (§9 is filled in after ⑨ executes).
> Baseline HEAD = `caf3b0f`; the ⑧c cleaning record, main commit hash and CI run id are filled in by the stage authorised at ⑩ in the same turn.
> Bilingual: this file is the CN authority; `docs/validation/directive-evidence-window_EN.md` is a section-by-section mirror.

## 1 Change summary

**Current change set** (`git diff --cached --numstat`, 9 paths, +664/−12 in total; the first 5 rows come from ⑧a and the last 4 were added at ⑧b (this report pair plus the ledger pair)):

| Path | + | − | Note |
|---|---|---|---|
| `docs/ADR_REGISTER.md` | 1 | 0 | ADR-026 registry row |
| `docs/ADR_REGISTER_EN.md` | 1 | 0 | ADR-026 `_EN` mirror row |
| `docs/DELIVERY_DIRECTIVE.md` | 23 | 6 | v1.39 directive body |
| `docs/DELIVERY_DIRECTIVE_EN.md` | 23 | 6 | strict `_EN` mirror |
| `docs/validation/evidence/DIRECTIVE-EVIDENCE-WINDOW-freedom-list.md` | 38 | 0 | new; execution-freedom list |
| `docs/validation/directive-evidence-window.md` | 270 | 0 | new; CN authority of this report |
| `docs/validation/directive-evidence-window_EN.md` | 270 | 0 | new; section-by-section `_EN` mirror |
| `docs/t027/REMAINING_SLICES.md` | 19 | 0 | P14 definition block plus the completion row |
| `docs/t027/REMAINING_SLICES_EN.md` | 19 | 0 | `_EN` mirror increment |

**Directive-side landing points** (`docs/DELIVERY_DIRECTIVE.md` and its `_EN` mirror, 8 hunks each):

- Header version line: v1.38 → v1.39.
- §0.3 new v1.39 row.
- §1.5 cleanup-timing parenthetical: temp artifacts inside the window are cleaned before ⑩ per the §5.1 sub-clause.
- §5.1 temp-artifact exception sentence + new "⑧c post-evidence window" sub-clause.
- §5.3 half-sentence in the ⑧c row.
- §12.4 OB-99 / OB-100 disposition columns.
- Appendix B.7.3 ⑩ receipt gains a "final-state readback" slot.
- Appendix C.0aa new ledger block.

**Change nature**: MINOR revision; it only adds the §5.1 window disposition rule and the evidence slot, without touching criterion bodies or pipeline structure.

**Execution freedom**: the substantive choices exercised by this slice are itemised in §13 and match `docs/validation/evidence/DIRECTIVE-EVIDENCE-WINDOW-freedom-list.md`.

## 2 Execution pipeline

| Step | Role / carrier | Model | State | Rework count |
|---|---|---|---|---|
| 0 dual-identity registration + ADR-026 | orchestrator | orchestrator model | landed | 0 |
| ① design | `deep-reasoner` | `DeepSeek V4 Pro (deepseek)` | `DESIGN: DONE` | 0 |
| ② implementation | `prfrail-implementer` | same model as orchestrator | `IMPLEMENT: DONE` (5 dispatches, 2 remediation rounds) | 2 |
| ③ testing | `prfrail-tester` | same model as orchestrator | r1 + r2 complete | 1 |
| ④ gates and integration | orchestrator | orchestrator model | complete | 0 |
| ⑤ pre-review | `independent-reviewer` | `GLM-5.3 (glm)` | r1 + r2 complete (1+1 exhausted) | 1 |
| ⑥ independent scan | `independent-reviewer` | `MAI-Code-1.1-Flash (copilot)` | `CODE-REVIEW: PASS` (blind) | 0 |
| ⑦ terminal review | `independent-reviewer` | `GPT-5.3-Codex (copilot)` Extra High | r1 + r2 + r3 complete (1+2 exhausted) | 2 |
| ⑧a draft | `prfrail-documenter` | same model as orchestrator | this report landed | 0 |
| ⑨ native verification | orchestrator | orchestrator model | complete (see §9.1) | 0 |
| ⑧b final pass | `prfrail-documenter` | same model as orchestrator | this report pair finalised | 0 |
| ⑧c close-out cleaning | orchestrator | orchestrator model | runs after the ⑩ same-turn authorisation | — |
| ⑩ stop-point | orchestrator | orchestrator model | stop-point reached; not committed / not pushed | — |

- The ADR-026 row of step 0 was landed by ② (see §3 deviation #2).
- The 5 dispatches of ② cover: ADR-026 for step 0, the six v1.39 clauses, the ③ defect remediation, the ⑤ remediation and the ⑦ remediation (the last dispatch carries the three ⑦ items).
- ③ ran r1 + r2; ⑤ ran r1 + r2; ⑦ ran r1 + r2 + r3.
- Elapsed time was not captured by this slice; it is marked "historically missing" and not invented.

## 3 Exceptions and fallbacks

| # | Event | Disposition |
|---|---|---|
| 1 | The freedom list was produced after ⑦ r1 / r2 (§1.7 requires after ⑥ and before ⑦) | Closed by the narrowed ⑦ r3 review (uses the third ⑦ slot); item §4 deviation #1 of the list records the timestamp evidence. |
| 2 | ADR-026 was landed by ② (not ⑧) | Landing-role deviation only; the authorisation source is unaffected; item §4 deviation #2 of the list records it. |
| 3 | ③ r1 finding F3 (the OB-99 anchor paragraph name "known boundary" should be "boundary 2") was not fixed | The anchor sits in the §12.4 phenomenon column; §11.1 forbids retroactive rewriting; ⑦ r2 accepted it. |
| 4 | ③ finding D12 (`_EN` spelling of `tmp/<slice-id>` in 3 places) was not fixed | `G3(a.1)` requires equal numstat on both sides, so a one-sided fix is necessarily red; changing both sides would need filler edits, which is rejected; left to a later slice that genuinely touches the CN rows. |
| 5 | Batch helper #1 found the declared six-shape `**` count stale: `06-six-shapes.md` in the package hardcodes 4284, measured 4298 | The substantive conclusion (CN / `_EN` equality) is unchanged and `G3(a.2)` `bold:true` still holds; §8 registers the measured value. |
| 6 | Batch helper #2's excerpt has one terminology error (it wrote "clean before ⑧b" for "clean before ⑩") | That sentence was not adopted; this report follows the directive text. |

- Items 1 and 3 are **timing / historical-row discipline** deviations, not implementation defects; item 4 is a trade-off forced by the criterion.
- Item 5 only corrects a reading and does not change any gate conclusion; item 6 only affects an excerpt, not the clause text.

## 4 Environment and tier

| Item | Value |
|---|---|
| OS | Windows (local `d:\LZProjects\prfrail`) |
| Node.js | `v24.17.0` |
| Go | `go1.27.0 windows/amd64` |
| Baseline HEAD | `caf3b0f` |
| Tier | M |
| Class | `[SLICE]` |
| ⑥ | mandatory (rule-semantics change) |
| ⑦ budget | 1+2 (exhausted) |
| ⑤ budget | 1+1 (exhausted) |
| Batch helpers | 2 (`GLM-5.3-Flash (glm)` number check + `GLM-5.3-FlashX (glm)` long-document excerpt) |
| Probes | 0 |

## 5 Review conclusions (⑤⑥⑦)

| Round | Id | Severity | Conclusion / disposition |
|---|---|---|---|
| ⑤ r1 | F1 | Medium | The ⑤ confirmation timing for prefetch omissions was unreachable; remediated to "the user decides at ⑩". |
| ⑤ r1 | F6 | Medium | Reading ① did not cover the first column of the index; remediated by adding a first-column set comparison. |
| ⑤ r1 | F2 | Low | §0.3 changed to a pointer sentence. |
| ⑤ r1 | F3 / F4 / F5 | Low | Remediated. |
| ⑤ r2 | N1 | Low | "The declared change set" was not anchored to a single source of truth; escalated as F-02 at ⑦ r1 and fixed. |
| ⑤ r2 | N2 | Low | Bidirectionality of "identical one by one". |
| ⑤ r2 | N3 | Low | The failure-shape enumeration was not synchronised with F6. |
| ⑥ | — | — | `CODE-REVIEW: PASS`, zero defects (blind; no ⑤ conclusions supplied). |
| ⑦ r1 | F-01 | Medium | OB-100 status word `已处置` conflicted with the "window only" scope; reverted to `观察中` and added "window sub-range handled". |
| ⑦ r1 | F-02 | Medium | "The change set declared by this slice" was not anchored to a single source of truth; reading ① now anchors to §5.2 G5 `declaredFiles`. |
| ⑦ r1 | F-03 | Low | "Created then self-cleaned inside the window" was not auditable; the sub-clause gained a "judge final-state observability only" boundary sentence. |
| ⑦ r2 | — | — | F-01 / F-02 / F-03 all closed; `TERMINAL-REVIEW: PASS`. |
| ⑦ r3 | F9 + 2 related + timestamp | Low | `FREEDOM-LIST: PASS WITH FIXES`; all three comments landed. |

- **⑤ r1** `GLM-5.3 (glm)`: `PRE-REVIEW: PASS WITH FIXES`; after remediating F1, F6 and F2, **r2** returned `RE-REVIEW: PASS WITH FIXES` with only 3 Low findings left (N1 / N2 / N3).
- **The ⑤ budget 1+1 is exhausted**; no third round was opened; N1 was escalated as F-02 and fixed at ⑦ r1.
- **⑥** `MAI-Code-1.1-Flash (copilot)` blind review: `CODE-REVIEW: PASS`, zero defects.
- **⑦ r1** `GPT-5.3-Codex (copilot)` Extra High blind review: `TERMINAL-REVIEW: FINDINGS` (F-01 / F-02 Medium, F-03 Low); after remediation the narrowed **r2** returned `TERMINAL-REVIEW: PASS`.
- **⑦ r3** narrowed to the freedom list only: `FREEDOM-LIST: PASS WITH FIXES`; the F9 wording, the two related synchronisations and the timestamp on deviation #1 all landed.

## 6 Review input packages (including blind-isolation proof)

| Package | Location | Content |
|---|---|---|
| ⑤ r1 | `tmp/G9-d/review-5/` | 14 items + `MANIFEST.md` (with SHA256 / byte sizes); all existence assertions passed. |
| ⑤ r2 | `tmp/G9-d/review-5b/` | narrowed package, 5 items. |
| ⑥ | `tmp/G9-d/review-6/` | 14 items; **no ⑤ conclusion of any kind attached**. |
| ⑦ r1 | `tmp/G9-d/review-7/` | 13 items; **no ⑤ / ⑥ conclusion of any kind attached**. |
| ⑦ r2 | `tmp/G9-d/review-7b/` | narrowed package, 5 items. |
| ⑦ r3 | — | narrowed input judging the freedom list only. |

- **Blind-isolation proof (⑥)**: the package contains no ⑤ conclusion file and no earlier review notes; the reviewer's scratch directory is separate from the input package.
- **Blind-isolation proof (⑦ r1)**: the package contains no ⑤ / ⑥ conclusion; the brief states "no earlier findings on purpose".
- **Whether the first round carried the list**: ⑥ and ⑦ r1 both **did not carry** the freedom list (it was produced later, see §3 item 1); ⑦ r3 is a narrowed input that judges the list only.
- **Input sanitisation**: the package contains only the slice definition, design, test report, change diff, numstat, gate readings and directive excerpts; no orchestrator conclusion is injected.
- This slice satisfies the ⑦ r1 statement that the first round carries no ⑤ / ⑥ list.

## 7 Falsifiability

- **Observable test of the rule revision (① (e), two criteria)**:
- Criterion one: delete the sub-clause ⇒ a later slice report again records the same shape "temp artifact recreated inside the window / location + timing violation / root cause is the §5.1 rule gap".
- Criterion two: delete the "pre-⑩ final-state readback" sentence ⇒ the ⑩ receipt no longer has to carry the four readings verbatim, so the in-window shape "location violation created but never staged" has no mechanical interception point before commit.
- **Failure shapes stated by the clause itself**: a `tmp/` top-level entry other than `.gitkeep`; `??` or a non-empty second column greater than 0; a ⑩ receipt without the verbatim readings.
- **Boundary**: only final-state observability is judged; the in-window process shape "created then self-cleaned" cannot be observed mechanically (OB-99 record state).
- **Mutation testing**: this slice does not change any criterion body and has zero `tools/**` changes, so no new `selftest.js` mutation case is added; regression is covered by the full `selftest` run.

## 8 Gate results

| Command | Result |
|---|---|
| `node tools/gates/selftest.js` | `SELFTEST: PASS 552/552` |
| `gate.js --all --scope=index` | `TOTAL_FAIL=0` |
| `gate.js --all --scope=tree` | `TOTAL_FAIL=0` |
| `gate.js --all --scope=ci` | `TOTAL_FAIL=0` |
| `gate.js --check=G3-a --scope=index` | PASS |
| `gate.js --check=G4b --scope=index` | PASS (`newCjk=558 unknown=[]`) |
| `gate.js --check=G6-4 --scope=index` | PASS |
| `gate.js --check=G6-5 --scope=index` | PASS (`scanned 28 refs 41`) |
| Six shapes (CN / `_EN`, measured equal) | see the table below |

**CN / `_EN` six shapes** (`docs/DELIVERY_DIRECTIVE.md` against its `_EN` mirror; measured values):

| Shape | CN | EN |
|---|---|---|
| lines (split) | 1823 | 1823 |
| lines (physical) | 1822 | 1822 |
| `**` | 4298 | 4298 |
| `\|` | 2897 | 2897 |
| blank | 442 | 442 |
| backticks | 4298 | 4298 |
| headings | 122 (same positions) | 122 (same positions) |

- Note: the declared `**` value 4284 is a stale reading; the measured value is 4298 (see §3 item 5); both sides are equal.
- `G4b` reading `newCjk=558 unknown=[]`; `G6-4` reading `scanned 626`; `G6-5` reading `scanned 28 refs 41`.
- This report pair is staged at ⑧b; every `--scope=index` reading above is a post-staging re-run that already includes the pair.

**This report pair, six shapes** (`docs/validation/directive-evidence-window.md` against `_EN`; measured after the ⑧b final pass, all six equal; the ⑧a baseline in parentheses):

| Shape | CN | EN |
|---|---|---|
| lines | 270 (242) | 270 (242) |
| headings | 17 | 17 |
| blank | 54 (49) | 54 (49) |
| `**` | 71 (62) | 71 (62) |
| `\|` | 541 (485) | 541 (485) |
| backticks | 306 (256) | 306 (256) |

## 9 ⑨ native verification and commit state (⑨ filled in; ⑧c / ⑩ see §12)

### 9.1 ⑨ native verification

**Method**: a pristine worktree (`git worktree add --detach tmp/G9-d/wt HEAD`, HEAD = `caf3b0f`) plus `git apply --index` replaying the staged change set (apply exit 0).

**Readings**:

- `--all` in all three scopes returns `TOTAL_FAIL=0` (exit 0).
- `SELFTEST: PASS 552/552`.
- `--check=G3-a|G4b|G6-4|G6-5` (`--scope=index`): all four `PASS`, `TOTAL_FAIL=0`, exit 0.
- The `v1.39` version header hits; the sub-clause heading "⑧c post-evidence window (v1.39)" hits; the `ADR-026` row hits.
- `wt_status_lines=5` (the change set is self-contained and does not depend on untracked review packs).

**Close-out**: `git worktree remove --force` exit 0, `git worktree prune`; the main repository is unaffected.

### 9.2 ⑧c close-out cleaning record

⑧c runs after the ⑩ same-turn authorisation; every reading below is an executed readback and no number is rewritten.

1. **Repository-wide encoding audit (not only the change set)**: `3633` version-controlled files, of which `3544` text files were checked file by file for BOM/LF ⇒ **0 violations**; `17` hits match the rule but sit inside frozen evidence packs ⇒ per OB-27 re-verified for self-consistency: the four packs' `SHA256SUMS` entries `10 / 29 / 74 / 77`, totalling **190/190 passing, 0 missing, 0 mismatched**.
2. **Temporary-file cleanup (file by file)**: the `75` files under `tmp/G9-d/` were deleted one by one with `Remove-Item -LiteralPath`, after which the empty directories were removed bottom-up and then the root directory. Readback proof: `Test-Path tmp/G9-d` = `False`; the `tmp` top level holds only `.gitkeep`; the recursive file count of `tmp` = `1`; the `??` entries of `git status --porcelain -uall` = `0`.
3. **⑨'s temporary worktree**: `git worktree remove --force` + `git worktree prune` already run; `git worktree list` shows only the main worktree.
4. **Editor-side diagnostics**: workspace-wide `get_errors` ⇒ **zero errors**.
5. **Re-run after cleanup**: `SELFTEST: PASS 552/552`; `gate.js --all --scope=index|tree|ci` reports `TOTAL_FAIL=0` in all three runs (exit 0).
6. **Sequencing compliance statement**: records were persisted first (this report's §1–§13, the freedom list, ledger P14) → then cleanup → then the gates were re-run, satisfying §5.1 (v1.35)'s requirement that cleanup must follow record persistence.
7. **F4 registration**: this slice's temporary directory is `tmp/G9-d/` (**provisional-named directory**; the registered name is `DIRECTIVE-EVIDENCE-WINDOW`), zeroed before ⑩ per the previous item; that transitional state is mechanically covered by readings ②③ of the §5.1 sub-clause.
8. **Four readings before ⑩ (verbatim into the ⑩ receipt)**: ① `??` entries of `git status --porcelain -uall` = 0, non-empty second-column entries = 0, the set of non-empty first-column entries = this slice's declared change set (`9` paths); ② `(Get-ChildItem -LiteralPath tmp -Recurse -File -Force | Measure-Object).Count` = `1`; ③ `(Get-ChildItem -LiteralPath tmp -Force).Name -join ','` = `.gitkeep`; ④ `Test-Path tmp/DIRECTIVE-EVIDENCE-WINDOW` = `False`.

### 9.3 Commit and CI

- **Main commit**: `00fbd2d` (`docs: add the post-⑧c evidence-window rule (v1.39)`), `9` files, `686` insertions / `16` deletions; pushed to `origin/main`; **gitee untouched**.
- **CI run**: `39235252911` (workflow `CI`, created `2026-10-19T10:08:16Z`), conclusion **success**.
- **Both legs and timestamps**: `Go ubuntu-latest` **success** (`10:08:19Z` → `10:10:55Z`); `Go windows-latest` **success** (`10:08:20Z` → `10:12:23Z`).
- **`Gates` step readings (identical on both legs)**: `SELFTEST: PASS 552/552`, `GATE REPORT scope=ci changed=9`, `TOTAL_FAIL=0`.
- **Package-level evidence**: `30` `ok   github.com/larsonzh/prfrail/...` lines across both legs (including the duplicated counts of the Race step), `--- FAIL:` = 0, `FAIL <pkg>` = 0.
- Green on the first run: no §6.4 disposition, no re-run.
- The backfill commit (this sub-clause) is docs-only; its own CI result is reported in the ⑩ follow-up.

## 10 Cost and accounting

| Slot | Count | Model / note |
|---|---|---|
| ① | 1 | `DeepSeek V4 Pro (deepseek)` |
| ② | 5 | same model as orchestrator (step 0 ADR-026 + six v1.39 clauses + ③ defect remediation + ⑤ remediation + ⑦ remediation; the last dispatch carries the three ⑦ items) |
| ③ | 2 | same model as orchestrator (r1 + r2) |
| ④ | 1 | orchestrator (gates and integration) |
| ⑤ | 2 | `GLM-5.3 (glm)` (1+1 exhausted) |
| ⑥ | 1 | `MAI-Code-1.1-Flash (copilot)` (budget 3) |
| ⑦ | 3 | `GPT-5.3-Codex (copilot)` (1+2 exhausted) |
| ⑧ | 2 | same model as orchestrator (⑧a draft + ⑧b final pass) |
| ⑨ | 1 | orchestrator |
| batch helper `GLM-5.3-Flash (glm)` | 1 | number check (found the stale six-shape `**` declaration) |
| batch helper `GLM-5.3-FlashX (glm)` | 1 | long-document excerpt (8 hunks / 23 insertions / 6 deletions, matching numstat; one terminology error not adopted) |
| probes | 0 | no probe budget for this slice |

## 11 Explicitly not executed

- No `git commit`, no `git push`, no gitee contact.
- No change to `tools/**`, §6.3, the §5.2 matrix, §7.5, §11.2, §11.3, or the pipeline structure.
- No permanent FlashX whitelist entry (ADR-026 states it ends when the clearance completes).
- This slice uses no probe budget (0).

## 12 Next steps

1. ⑨ native verification: done (encoding gate + G1–G9 + reference resolvability + `selftest`, readings in §9.1).
2. ⑧b final pass: done (§9 and ledger P14 filled in, the report pair and ledger pair staged, gates re-run).
3. ⑧c close-out: clean `tmp/G9-d/` per the §5.1 sub-clause and prove it by readback.
4. ⑩ stop-point: wait for same-turn user authorisation before commit and push, then backfill the commit and CI facts.

## 13 Execution freedom

Item-by-item reference to `docs/validation/evidence/DIRECTIVE-EVIDENCE-WINDOW-freedom-list.md` (the list is authoritative; this table does not change its meaning).

| # | Choice | Landing point |
|---|---|---|
| F1 | The three candidate paths form a fixed order: default one-shot diskless output; alternative one prefetch before ⑧c (mandatory when enumerable); alternative two recreate after ⑧c (exception only + forced order) | §5.1 sub-clause |
| F2 | No new CI criterion (no `G10`) and no §6.3 change | ① verdict; zero §6.3 change |
| F3 | Slice name changed to `DIRECTIVE-EVIDENCE-WINDOW` (no G number consumed) | ledger P14, report filenames, freedom-list filename |
| F4 | The readback uses four fixed commands; the count uses `Get-ChildItem -Recurse -File -Force` (git semantics explicitly disabled) | §5.1 sub-clause reading ② |
| F5 | Reading ① gains "first-column row set = declared change set", anchored to §5.2 G5 `declaredFiles` | §5.1 sub-clause reading ① |
| F6 | The confirmer for prefetch omissions is the user (decided at the ⑩ stop-point), not the ⑤ pre-review | §5.1 sub-clause alternative one |
| F7 | The appendix ledger block number is `C.0aa` (`C.0z` already used) | appendix C |
| F8 | The ⑩ receipt template (appendix B.7.3) gains a "final-state readback" slot | appendix B.7.3 |
| F9 | The OB-100 status word stays `观察中`; the disposition cell adds the boundaries "window sub-range handled" and "other periods still not mechanised" | §12.4 |
| F10 | The OB-99 "known boundary" anchor in the phenomenon column is unchanged (§11.1 forbids retroactive rewriting of the phenomenon column) | ⑦ r2 accepted |
| F11 | The 3 `tmp/<slice-id>` spellings in `_EN` are unchanged (`G3(a.1)` requires equal numstat; a one-sided fix is red) | ③ r2 reconfirmed |
| F12 | This slice uses `tmp/G9-d/` (working name) and did not migrate it after the rename | ⑧c record |

**Freedom list §2 related mechanical sync (not freedom)**: the ⑤ r1 F1 / F6 / F2 and ⑦ r1 F-01 / F-02 / F-03 remediations are defect fixes; the §0.3 reading enumeration became a pointer sentence; §1.5 cleanup timing gained the in-window temp-artifact disposition; the §12.4 OB-99 disposition cell gained the boundary limit.

**Freedom list §4 known deviations**: ① the freedom list came after ⑦ r1 / r2 and was closed by the narrowed ⑦ r3 review; ② ADR-026 was landed by ②, a landing-role deviation; ③ the ③ r1 F3 (OB-99 anchor paragraph name) was not fixed because §11.1 forbids retroactive rewriting.

```artifacts
docs/validation/directive-evidence-window.md	repo
docs/validation/directive-evidence-window_EN.md	repo
docs/validation/evidence/DIRECTIVE-EVIDENCE-WINDOW-freedom-list.md	repo
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/ADR_REGISTER.md	repo
docs/ADR_REGISTER_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
```
