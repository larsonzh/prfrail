# Slice Validation Report — `DIRECTIVE-OB-SWEEP`

Date: 2026-10-16. Slice id: `DIRECTIVE-OB-SWEEP`. Tier: `[SLICE]`. Size: S (`.md` only). Probe budget: none.

## 1. Goal, scope and the three-condition ruling

- **Goal**: batch-dispose three `.md`-only OB fixes - OB-60 (appendix B.1's architecture design document gains the "cross-package shared contracts" and "failure-path enumeration" sections), OB-82 (§7.7's proportional spot-check window gains a counting anchor and a record location), OB-84 (the §12.5 gap row gains a budget caliber for the ⑦ cost).
- **Boundaries**: unchanged §2.2 / §2.4 / §5.2 matrix / §6.3 criterion bodies / §7.5 figures; unchanged CONTRACTS / schema / fixtures / workflows; no change under `internal/**`; role set and carriers untouched; no gitee push (this slice stops at ⑩, uncommitted and unpushed).
- **§7.2.1 three-condition ruling (① architecture)**: (1) `.md` only holds; (2) no `.go` / `.json` / `.yml` / schema / fixtures semantic change holds; **condition ③ does not hold** (this slice changes the §7.7 review-system rule text, the §12.5 cost caliber and the appendix B.1 role-template output structure) ⇒ **⑥ is mandatory**. ⑤ and ⑥ each confirmed this independently.

## 2. Artifact list

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/directive-ob-sweep.md	repo
docs/validation/directive-ob-sweep_EN.md	repo
docs/validation/evidence/DIRECTIVE-OB-SWEEP-freedom-list.md	repo
tmp/obsweep/gates-index.txt	local-only
tmp/obsweep/gates-tree.txt	local-only
tmp/obsweep/impl.diff	local-only
tmp/obsweep/norm.js	local-only
tmp/obsweep/numstat.txt	local-only
tmp/obsweep/proof.js	local-only
tmp/obsweep/proof.txt	local-only
tmp/obsweep/selftest.txt	local-only
tmp/obsweep/shape.js	local-only
```

> This report pair lists itself in the `repo` disposition rows (it must be `git add`-ed before ⑨, otherwise G7-b turns red).

## 3. Change summary

- **OB-60**: appendix B.1's "Output" list grows from 5 items to 7 (new item 6 "cross-package shared contracts" and item 7 "failure-path enumeration"), with one added applicability-trim line (not required at tier S, mandatory at tier M or above and for any cross-package change); the role / model / carrier / input / forbidden rows are verbatim untouched.
- **OB-82**: a sub-item is added under the §7.7 proportional spot-check bullet, giving the window counting anchor (starting at v1.32; the registration order of the ledger's "later batch slices" section) and the record location (the "proportional spot-check record" block at the end of §12.4); the original judgement-rule sentence is verbatim untouched (mechanical proof in §6).
- **OB-84**: in §12.5 the "no budget model for ⑦'s three readings" gap row now says "Landed" in its timing column and gives the budget caliber (= the §7.5 call cap × the measured per-reading cost); the §7.5 region has zero change.
- **Version package**: header v1.31 → v1.32; a v1.32 log entry is added to §0.3; appendix C gains the C.0t ledger block.
- **Ledger**: P4 gains a Status field row and a time-scoped parenthetical on its "Current gap" row; both the heading status symbol and the status row are `🔄` (in progress; flipped to `✅` by ⑧b).
- **OB disposition cells**: the first status words of OB-60 / OB-82 / OB-84 flip from `待议` / `待议` / `观察中` to `已处置`.

## 4. Item-by-item file:line assertions (acceptance evidence ③)

- **OB-60**: the appendix B.1 fence carries "cross-package shared contracts" and "failure-path enumeration" (items 6 / 7, `docs/DELIVERY_DIRECTIVE.md:1066-1067`) and the applicability-trim line (`:1068`), the latter including "items 6 / 7 are **not required at tier S**"; the `_EN` side carries items 6 / 7 and the applicability trim at the matching position.
- **OB-82**: §7.7 carries the "window counting anchor and record location" sub-item (`docs/DELIVERY_DIRECTIVE.md:674-675`, the original sentence at `:674` and the new sub-item at `:675`) including "the window starts at v1.32"; §12.4's end carries the "proportional spot-check record" block (`:976-977`); the `_EN` side matches at the same position.
- **OB-84**: the §12.5 gap row (`docs/DELIVERY_DIRECTIVE.md:991`) carries "**Landed**" and "the §7.5 figures are unchanged"; the `_EN` side matches at the same position.
- **Disposition cells**: the last cell of the OB-60 / OB-82 / OB-84 rows (`docs/DELIVERY_DIRECTIVE.md:938` / `:960` / `:962`) starts with `**Resolved**` (`**已处置**` in CN).
- **Version**: the header (`docs/DELIVERY_DIRECTIVE.md:3`) `Version: v1.32` and the §0.3 last row (`:78`) `v1.32` agree (G1-b passed in the real run).

## 5. Freedom and anomalies

- Freedom items F-1 … F-8 are itemised in `docs/validation/evidence/DIRECTIVE-OB-SWEEP-freedom-list.md` and are all judged in-bounds.
- **A-1 (tool form)**: creating the freedom list landed it as BOM-free + CRLF, against the repository convention ⇒ a pure encoding normalization was applied (BOM + LF, zero content change).
- **A-2 (caliber disclosure)**: the new §7.7 sub-item is an **additive** completion; ⑨'s reference resolvability is judged over the **changed-line scope**, and inherited citations elsewhere in the file are not defects.
- **A-3 (process deviation - the §5.3 ordering, a master orchestration mistake)**: after ⑤ r1 reported 3 Medium + 5 Low, ⑤ was **not re-run in place** as §5.3 requires; the flow went straight to ⑥, and ⑤ r2 was only **run later** after ⑦ r1 reported `FINDINGS` and its remediation was done. ⑤'s 1+1 cap was **not consumed** at that time (so this is not a cap exception), and r1 carried Medium+ ⇒ the **OB-11 exception does not apply**; recorded truthfully, not silenced.
- **A-4 (the late run's input package and the item-by-item closure map)**: ⑤ r2's input package was the **same index snapshot** taken after both the ⑤ r1 remediation and the ⑦ r1 remediation (`git status --porcelain` proves index = disk); ⑤ r1's 3 Medium + 5 Low were **each judged CLOSED** in ⑤ r2 (M-1 / M-2 / M-3 / L-1 / L-2 / L-3 / L-4 / L-5); the 4 Lows ⑤ r2 reported (N-L1 … N-L4) were remediated and then **confirmed on behalf of ⑤** by ⑦ r2 on the final snapshot. (A-3 / A-4 are factual records added under the user's direction in the same turn; they do not trigger a ⑦ re-run.)
- **A-5 (the boundary of ⑥'s coverage, a factual record)**: ⑥'s `PASS` is bound to the **intermediate snapshot taken after the ⑤ r1 remediation and before the ⑦ r1 remediation**; **the final snapshot is covered by ⑦ r2**, and **⑥ was not re-run on the final snapshot** (⑥ is a supplementary scan layer and its conclusion does not prove review sufficiency). User ruling of 2026-10-16: "run-later coverage" is accepted and the **deviation is not dissolved**.

## 6. Mechanical proof that §7.5 is unchanged, and verbatim retention of the rule sentence

- **Zero change in the §7.5 region**: every `git diff -U0` hunk lands on the version line, §0.3, appendix B.1, §7.7, §12.4 and §12.5, and no hunk covers §7.5.
- **Verbatim retention of the rule sentence**: the §7.7 proportional spot-check sentence has no deleted line in the diff; after the change it is byte-identical to the `HEAD` version.

## 7. Review-chain verdicts

| Stage | Owner | Verdict |
|---|---|---|
| ① architecture | `DeepSeek V4 Pro (deepseek)` | `DESIGN: DONE` (the three-condition and three-OB boundary rulings all hold) |
| ③ consistency | master (mechanical scripts) | references and figures measured item by item; `G3(a)` five shapes equal |
| ④ gates | master | tree / index `TOTAL_FAIL=0`; `SELFTEST: PASS 335/335` |
| ⑤ pre-review | `GLM-5.3 (glm)` | see the review-chain log in §11 |
| ⑥ independent scan | `MAI-Code-1.1-Flash (copilot)` | see the review-chain log in §11 |
| ⑦ final review | `GPT-5.3-Codex (copilot)` Extra High standard tier | see the review-chain log in §11 |

## 8. Known gaps and candidates (recorded truthfully; not a pass)

- **The line numbers in P4's source citation were already stale before this slice**: P4's "Sources" paragraph cites `docs/DELIVERY_DIRECTIVE.md:915` (OB-60) and others, yet the measured line **before** this slice's change was already 936 (938 after it) ⇒ that staleness predates this slice; per §1.7.2 this slice does not widen scope, and the item is recorded as a candidate only.
- **The §7.7 window start has no historical caliber**: this slice sets the window start to v1.32 and does not retroactively rule on earlier slices; if a later review demands retroactive coverage, that is a rule-semantic change ⇒ it needs a separate slice.
- **No dedicated mechanical criterion**: all three items are wording and structure changes, with no dedicated criterion beyond the generic G1–G8 set.
- **This slice is the window's first slice, so the §7.7 self-reference is out of order in time**: the rule-creating slice is also the window's first, and its decision and reason are registered at wrap-up rather than before start, which is literally unsatisfiable for the very slice that creates the rule; the substantive duty is met (the §12.4 end block registers this slice's decision).
- **Candidate (no scope widening)**: this slice's ⑤ ordering deviation (a class-C orchestration mistake) is suggested to a later directive-governance slice for assessment of whether to add it to the §9.3 fault-mode library; **this slice changes no §12.4 entry**.

## 9. ⑨ and ⑧c readings

- Encoding: every `.md` changed or created by this slice is BOM + LF; the ⑨ whole-tree per-file check reports `tracked=3607 text=3518 binary=89 violations=17`, with **zero violations outside the frozen main tree** (all 17 sit inside the frozen archive packs under `docs/validation/evidence/**`, exempt per §5.1).
- Reference resolvability (changed-line scope): 212 `§X.Y` references, 0 unresolved; G6-5 separately reports `scanned 40, refs 56` with none unresolved.
- Gates: `--scope=tree` and `--scope=index` both `TOTAL_FAIL=0`; `--selftest` `SELFTEST: PASS 335/335`; `G3(a)` gives `+32/-5` vs `+32/-5` for the directive pair and `+3/-2` vs `+3/-2` for the ledger pair; `G4b unknown=[]`; `G6-4` and `G6-5` have no hits.
- Parity: the directive pair's five shapes (line count / heading line-number sequence / bold / bars / blank lines) are all equal.
- Residue: the `tmp/obsweep/` and `tmp/cn.diff` process artifacts are removed in ⑧c; no untracked residue remains.

## 10. Commit and CI

- Commit: `2d6a8e4` (`docs: DIRECTIVE-OB-SWEEP 三项 OB 批量处置落地（v1.32）`; 7 files `+319/-14`).
- Push: `origin/main`, range `cb64d6b..2d6a8e4` (**no gitee push**).
- CI: run `38810986151` (head `2d6a8e4`) **first run green on both legs** (windows 3m45s / ubuntu 2m38s); both legs `SELFTEST: PASS 335/335`, `GATE REPORT scope=ci changed=7`, `TOTAL_FAIL=0`; package level `ok` 30 / `FAIL` 0 (the ubuntu leg additionally runs Race: `internal/adapters` and `internal/chain` both ok); **§6.4 was not triggered**.
- Commit-message language **deviation** registration: this slice's commit summary was written in Chinese, **deviating** from the "English `<type>: <summary>`" discipline; **disposition = no history rewrite**, with the discipline effective from the next slice; basis = the user's ruling text of 2026-10-16.
- This section is backfilled by the docs-only commit after CI turned green (§11.2); that backfill commit's own CI run is an observation record only; the N-L3 caliber = the P4 definition row's "three scopes" is historical text that this slice does not change, and this slice's self-measured caliber is the C.0t wording "the tree and index scopes (measured)".

## 11. Review-chain log (final before ⑩)

- The verdict lines of ⑤ pre-review, ⑥ independent scan and ⑦ final review are appended at the end of this section by ⑧b (before ⑩); their caps follow §7.5 (⑤ 1+1, ⑥ 3, ⑦ 1+1).
- This slice's ⑦ = `GPT-5.3-Codex (copilot)` at the Extra High standard tier; **not downgraded, OB-37 not triggered**; probes 0.
- Factual-record note: this section is written after ⑦ and records facts (verdict lines and measured values); per the OB-11 exception caliber it does not trigger a ⑦ re-run.
- **⑤ pre-review**: r1 = `PRE-REVIEW: PASS WITH FIXES` (3 Medium + 5 Low, remediated); **no in-place re-review followed ⑤ r1** (§5.3 ordering deviation, see A-3 in §5), and r2 was **run later** after the ⑦ r1 remediation = `PRE-REVIEW: PASS WITH FIXES` (**0 Medium+**, 4 Low).
- **For the user to rule on**: whether "⑤ r2 run later, after ⑦ r1" is accepted as **covering** the §5.3 ordering duty (rather than an in-place closure); this slice does not dissolve the deviation on its own.
- **⑤'s cap is exhausted (1+1)** ⇒ the literal remediation of its 4 r2 Lows is handed to the ⑦ re-review for confirmation (recorded truthfully in this report).
- **⑥ independent scan**: `INDEPENDENT SCAN: PASS` (Section A–E complete, Section B empty, verdict line consistent with Section B).
- **⑦ final review**: r1 = `FINAL REVIEW: FINDINGS` (1 High + 2 Medium, remediated with a rebuilt input package) ⇒ r2 = `RE-REVIEW: PASS` (all 7 items CLOSED, zero new Medium+).
- **The root cause of ⑦ r1's High** (input package inconsistent with the disk) was removed under the "input-package freshness discipline": `impl.diff` / `numstat.txt` / `proof.txt` and the three gate readings all come from the same index snapshot taken after the fixes.
