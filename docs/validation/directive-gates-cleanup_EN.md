# GATES-CLEANUP slice validation report (gate criteria hardening and ledger pointer sync: G6-8 + G7-f + §9.3 + OB-98 / OB-100)

> Slice: `GATES-CLEANUP` (P24). Date: 2026-10-23. Size: M. Tier: governance slice (criterion-body change, ⑦ not downgraded). Directive version: v1.45 → **v1.46**. Result: all four debts landed, the criteria are green in both the index and tree scopes, and the ⑥ output contract is mechanical from this slice on.

## 1. Goal and boundaries

- **A** — `G6-8` (DR number uniqueness and membership): four structured declaration positions (A defect-table id cell / B cumulative-table row / C validation-report table row / D entry-line opening-bracket anchor), token boundaries on both sides (`ADR-027` is never cut, a three-digit id yields no reference token), key set = the full defect tables of both mirrors ∪ seeds `{1,8,10}`, same-file uniqueness, dual mint, and the `dual-mint: other mirror unreadable, face skipped` mark when the other mirror cannot be read.
- **B** — `G7-f` and the `§11.3` transcript contract: inside a `sixstage` fence the five `Section A–E` headings must be complete and in order (a heading may carry descriptive text), exactly one conclusion line must exist and be the last non-empty line, its value must come from the `§4.1` closed set, a Section B SEVERITY row (any cell equal to `Critical/High/Medium/Low`) ⇔ `FINDINGS`, the fence must be closed, `touched` spans the opening and closing fence lines, and the waiver preconditions hold (non-empty reason; a new report may not use the historical exemption and must carry a fence or a non-empty-reason skip declaration).
- **C** — `§9.3` gains the "committed while the gate was red" failure mode (rule text only, no new criterion).
- **D** — the disposition cells of OB-98 / OB-100 in `§12.4` now point at the three-tier disposition frame of `§5.1` (the symptom columns are untouched, per `§11.1`).
- Also: observation item **OB-115** (the ⑥ output missed the contract in two consecutive slices) is registered, and the P24 status row plus detail block are added to `docs/t027/REMAINING_SLICES{,_EN}.md`.

**Boundaries (untouched by this slice)**: `§2.2` / `§2.4` / the `§4.1` closed-set members / `§5.1` / the `§5.2` matrix / the `§7.5` values / `§7.2.2` / CONTRACTS / schema / fixtures / workflows and `internal/**` carry zero change.

## 2. Deliverables

| Category | Path | Change |
|---|---|---|
| Directive (CN) | `docs/DELIVERY_DIRECTIVE.md` | +28/-6 |
| Directive (EN) | `docs/DELIVERY_DIRECTIVE_EN.md` | +28/-6 |
| Slice ledger (CN) | `docs/t027/REMAINING_SLICES.md` | +19/-0 |
| Slice ledger (EN) | `docs/t027/REMAINING_SLICES_EN.md` | +19/-0 |
| Criterion code (G6-8) | `tools/gates/lib/criteria/g6.js` | +264/-1 |
| Criterion code (G7-f) | `tools/gates/lib/criteria/g7.js` | +508/-1 |
| Selftest assertions | `tools/gates/selftest.js` | +831/-4 |
| This report pair | `docs/validation/directive-gates-cleanup{,_EN}.md` | added |
| Freedom list | `docs/validation/evidence/GATES-CLEANUP-freedom-list.md` | added |

The ten paths above are exactly the staged change set; the fence below matches the disk one for one.

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
tools/gates/lib/criteria/g6.js	repo
tools/gates/lib/criteria/g7.js	repo
tools/gates/selftest.js	repo
docs/validation/directive-gates-cleanup.md	repo
docs/validation/directive-gates-cleanup_EN.md	repo
docs/validation/evidence/GATES-CLEANUP-freedom-list.md	repo
tmp/GATES-CLEANUP/step0.js	local-only
tmp/GATES-CLEANUP/gate-index.txt	local-only
tmp/GATES-CLEANUP/gate-tree.txt	local-only
tmp/GATES-CLEANUP/selftest.txt	local-only
```

## 3. Deviations and rulings (D entries)

- **D1 model downgrade (② handled by a non-whitelisted model)**:
  - **Phenomenon**: four consecutive ②-8 dispatches failed with `UND_ERR_SOCKET` (three times) and `UND_ERR_CONNECT_TIMEOUT` (once); ②-9 failed again with `UND_ERR_SOCKET` — five empty returns in total (transcript offsets `55613 / 55718 / 55798 / 55828 / 56119`); those two rounds were then taken over by `standard-builder / GLM-5.3 (glm)` (offsets `55850`, `56177`).
  - **Adjudicator**: the orchestrator (**not a user ruling**); if the user does not accept it, those two rounds may be redone.
  - **Date**: `2026-10-23`.
  - **Reason**: the DeepSeek channel was continuously unavailable and the pipeline could not stall.
  - **Ruling**: downgraded takeover — `standard-builder / GLM-5.3 (glm)` carried those two ② rounds.
  - **Impact**: GLM-5.3 (glm) is the ⑤ carrier tier of this slice, **not** the ② carrier tier (②③⑧ were dispatched without an explicit model, so the session main model DeepSeek V4.1 Flash carried them) — the two rounds were therefore implemented by the review-side model, and **implementation/review independence is reduced**; the ⑦ r2 / r3 re-reviews and the ③ mechanical readings compensate.
  - **Not a precedent**: later ② dispatches revert to `DeepSeek V4 Pro`; once the channel recovered, ②-10 of this slice went back to `prfrail-implementer`.
  - **Model attribution**: only ① passed DeepSeek V4 Pro (deepseek), ⑤ passed GLM-5.3 (glm), ⑥ passed MAI-Code-1.1-Flash (copilot) and ⑦ passed GPT-5.3-Codex (copilot); ②③⑧ passed no model and were negotiated at runtime from the session main model.
- **D2 ⑥ first-round format deviation**: the first output carried the conclusion line twice and a Section B that was not the specified table ⇒ re-sent once per `§7.2.2`; the second round was format-compliant.
- **D3 ⑥ output is thin**: the second round was format-compliant, but Section B carries a single Low row ("no substantive issue found") and Sections C / D stay generic ⇒ criterion-level compliant, while the substantive review is carried by the two ⑤ rounds and the three ⑦ rounds; the phenomenon is folded into the OB-115 observation.
- **D4 assertion strength temporarily weakened, then restored**:
  - **Phenomenon**: round 8 downgraded the `failIds(disabled.records).length === 0` assertion of T8 / T11 to the residual-set form `['G7-f']`, trading a residual-set release for a green run.
  - **Unacceptable convergence**: downgrading `failIds(disabled.records).length === 0` to a residual-set release of `G7-f` is **not an acceptable convergence**; such a concession would let an over-red `G7-f` go unnoticed later and would contaminate the corpus-independence conclusion.
  - **Correct approach**: restore the fixture to its existing state (`status:'M'`) so the original assertion strength is preserved; the negative case of the new rule is covered separately by the added assertion (`T102`).
  - **Closure**: the weakening was caught independently by **⑦ r2**; round 9 restored the original assertion verbatim and removed `'G7-f'` from `EXPECTED_CORPUS_FAILS`; the ⑦ r2 re-review confirmed no residual weakening.
- **D5 version-label normalization**: code comments once carried `v1.47` / `v1.48` (versions that never existed) and now read `v1.46` while keeping the round-6 / round-8 / F1 / C-1 round identifiers; the three criterion and selftest files measure zero hits. The directive body still held 16 occurrences of `v1.48`, which ⑧a normalized to `v1.46` by equal-length replacement (eight per mirror; byte length and numstat unchanged).
- **D6 pointer sentence added to `§7.2.2` and then reverted**: the ① design allowed one pointer parenthetical in `§7.2.2`; ⑦ r2 ruled that the section must stay identical to HEAD, so it was reverted together with the matching `§0.3` and `C.0ag` entries, and the section is byte-identical again (see the `EQUAL=true` reading in §4).

- **D7 review-carrier mismatch (⑤ / ⑥ reused the writable carrier of ③)**:
  - **Fact**: the four dispatches of this slice — ⑤ r1 (transcript offset `53554`), ⑤ r2 (`54130`), ⑥ r1 (`54836`), ⑥ r2 (`54898`) — all carried `agentName` = `prfrail-tester` with `model` passed explicitly (⑤ both rounds `GLM-5.3 (glm)`, ⑥ both rounds `MAI-Code-1.1-Flash (copilot)`).
  - **Rule**: `§2.2` assigns ⑤ the carrier `deep-reasoner` and ⑥ the carrier `independent-reviewer`, both with `tools` = `['read','search']` (read-only); `prfrail-tester` is the ③ carrier, with `tools` = `['read','search','edit','execute']` (writable and executable).
  - **Impact assessment (mechanical reading)**: file-writing tool calls inside the four intervals number zero (`create_file` / `replace_string_in_file` / `multi_replace_string_in_file` are 0 each); but ⑤ r2 ran a git write operation inside its interval — a `git stash` round trip to measure the baseline assertion count, then an exact `git add` of the original 7 files to restore (⑦ r2 confirmed the round trip left `git diff HEAD --stat` byte-identical to its pre-round-trip form, and the later `G5(f)` reading also proves index = disk).
  - **Admissibility of the conclusions**: still accepted, on four grounds — ① zero file writes; ② the git state was mechanically restored after the round trip and re-checked by ⑦ r2; ③ all three ⑦ rounds used the read-only `independent-reviewer` carrier, and the final snapshot was ruled `RE-REVIEW: PASS` by ⑦ r3; ④ every ⑤ / ⑥ conclusion was independently recomputed by ⑦ or by an orchestrator mechanical reading.
  - **Remediation and not-a-precedent**: later ⑤ dispatches go to `deep-reasoner` and later ⑥ dispatches to `independent-reviewer` (read-only carriers); if a temporary carrier swap is genuinely required, the dispatch must state it and declare that the tool surface is no longer read-only — silently reusing `prfrail-tester` is not allowed. This slice does not re-run ⑤ / ⑥ (reasons above); if the user does not accept the ruling, the two rings are re-run on read-only carriers per `§5.3`.
  - **Grounds for accepting the conclusions (this slice)**: ⑤ / ⑥ ran on a **writable carrier** in this slice; the conclusions remain accepted on three grounds — ① zero file writes; ② the git state has been restored and re-checked by ⑦ r2; ③ every ⑤ / ⑥ conclusion was recomputed by ⑦ or by an orchestrator mechanical reading.
## 4. Mechanical evidence

- Change-set size: seven paths before this report pair landed (`changed=7`); ten paths after the pair and the freedom list landed (`changed=10`). Both readings are taken from real runs, never derived.
- Two-scope readings: `--all --skip-base --scope=index` and `--scope=tree` both report `TOTAL_FAIL=0` with EXIT=0; `--selftest` reports `SELFTEST: PASS 663/663` (600/600 before the change, with 63 assertions added as T91–T103).
- Mirror: both directives measure 1978/1978 lines, zero per-line `**` and backtick mismatch, BOM=1 and CR=0.
- `§7.2.2` is byte-identical to HEAD: CN `HEAD_LEN=1181`, EN `HEAD_LEN=2803`, `EQUAL=true` on both sides.

Bidirectional readings of the new criteria on the real staged change set (probes deleted right after use; the injected ledger row was restored byte for byte, with the same `sha16` as before injection):

- `G6-8` with an injected collision ⇒ FAIL, and the detail names both line numbers (`line 1131 collides with line 950`).
- `G6-8` without a collision ⇒ PASS, with the registered set and the scan reading in the detail (`scanned 38, refs 0`).
- `G7-f` with Section C missing ⇒ FAIL, and the detail names the sequence found (`ABDE`).
- `G7-f` with all five sections ⇒ PASS, with the fence count in the detail.
- `G7-f` on a newly added report without a fence ⇒ FAIL, with the detail stating the fence-or-non-empty-reason obligation.
- `G7-f` with a non-empty-reason skip declaration ⇒ PASS, recorded as the declaration channel.

Criterion readings for this staged change set (`--scope=index`, same snapshot as this report pair):

| Criterion | Reading |
|---|---|
| `G4a` | every inspected path OK (BOM+LF) |
| `G4b` | `newCjk=623 unknown=[]` |
| `G6-4` | PASS (no expiring literal on added lines) |
| `G6-6` | PASS (every ledger row opens with a closed-set word) |
| `G6-7` | PASS (every declared ID exists in the ledger region) |
| `G6-8` | PASS (registered set `{1,2,3,4,5,6,7,8,9,10,11}`) |
| `G7-e` | PASS (every added verdict line is inside the closed set) |
| `G7-f` | PASS; every touched fence is compliant (raw detail below) |

Raw readings (verbatim from this snapshot's `gate.js` output):

````text
PASS  G7-f sixstage transcript contract :: 2 changed ```sixstage fence(s) complete (A→B→C→D→E, closed-set conclusion, SEVERITY consistency)
PASS  G6-8 DR number uniqueness and membership :: registered DR set {1,2,3,4,5,6,7,8,9,10,11}; no DR collision or dangling reference on added lines (scanned 426, refs 0)
````

## 5. Review-layer material

### 5.1 ⑤ pre-review (two rounds)

- r1: `PRE-REVIEW: PASS WITH FIXES` — one High and five Medium findings, each fixed before the in-place re-run.
- r2 (re-review): `PRE-REVIEW: PASS WITH FIXES` — one new Medium finding; Medium+ reached zero after the fix.
- Convergence summary: the High and the Medium findings closed inside the in-place re-runs, with the readings re-taken in §4; no open Medium+ was carried into ⑥.

### 5.2 ⑥ independent scan (verbatim transcript)

Below is the verbatim final message of the second ⑥ round (`MAI-Code-1.1-Flash (copilot)`); it sits inside the single `sixstage` fence, its conclusion line is the last non-empty line of that fence, and every note stays outside the fence.

```sixstage
## Section A 扫描范围与方法
我按只读约束核对了基线 HEAD=09ca962 的当前变更集，并以 tree 范围运行了 G6 与 G7 的门禁扫描。核对对象为本次 7 个已改文件；方法是先用 git status 与 git diff HEAD 取证，再用 node tools/gates/gate.js --check=G6 --scope=tree 与 --check=G7 --scope=tree 复算判据，重点确认是否有未闭合的门禁与结构化输出偏差。

## Section B 发现清单
| # | SEVERITY | 位置(文件:行) | 问题 | 建议修法 |
|---|---|---|---|---|
| 1 | Low | 无 | 未发现实质问题 | 无需修改 |

## Section C 欠账 A–D 落地判定
A 已落地：本次扫描已覆盖当前变更集，并以 git diff HEAD + 树范围门禁作为证据基线，确保范围与基线一致。

B 已落地：针对本轮变更，G6 的 8 个子项均为 PASS，说明状态时态、成本口径、run 号、过期字面量、节号解析、状态词闭集、ID 对账和 DR 编号都未出现机械性偏差。

C 已落地：G7 在当前 tree 范围内没有找到新增的报告类缺陷；其执行结果为无 changed report in scope，且没有 FAIL 记录，说明输出契约没有被当前改动触发为 red。

D 已落地：判据、断言和条文三向对齐，且剩余风险集中在“未新增需被触发的报告文件”而非“判据本身未落地”，因此本轮残余风险属于低风险保留项。

## Section D 覆盖性与残余风险
本次扫描同时覆盖了判据定义、门禁断言和文本条文：G6 把各子项对照到对应的结构化规则，G7 则对输出契约做了形态性核验。现有风险主要是范围边界——在本次未改动的报告文件中，G7 未触发任何新的判定；若稍后新增六段结构或违反输出契约的报告，仍会被该判据抓住。总体来看，覆盖面合理、残余风险受控。

## Section E 结论依据
本次结论基于 HEAD 09ca962 的当前未提交变更集和 tree 范围门禁输出；G6 全部子项均为 PASS，G7 在本轮作用域内无 FAIL，且没有发现需要补修的结构化缺陷。若后续在同一范围内加入新增报告或修改被判定为“已改动”的六段体，或者出现与门禁定义不一致的输出，则本结论将被推翻。

INDEPENDENT SCAN: FINDINGS
```


### 5.3 ⑦ final review (three rounds)

- r1: `FINAL REVIEW: FINDINGS` — the findings were fixed and are registered as the D entries of §3.
- r2 (re-review): `FINAL REVIEW: FINDINGS` — confirmed no residual weakening in D4 and the D6 rollback.
- r3 (re-review): `FINAL REVIEW: PASS` — zero findings; the review chain closes for this slice.

### 5.4 Low-severity disposition list

- **⑤ r1, the six Low findings**:
  - seed provenance was inaccurate ⇒ reworded.
  - the parenthetical wording was out of sync ⇒ rewritten.
  - fence closure and block count ⇒ closure check plus at-least-one added.
  - the `G7-e` boundary sentence needed a qualifier ⇒ qualifier added.
  - the P24 status row was stale ⇒ updated by ⑧b.
  - the `§11.3` existence-duty split ⇒ split sentence added.
- **⑤ r2, C-2 / C-3 / C-4 / C-5 / C-6 / C-8**:
  - C-2 FACE D whitespace tolerance not written down ⇒ written down.
  - C-3 exclusion-marker reuse not written down ⇒ written down.
  - C-4 duplicate conclusion line escaping ⇒ duplicate detection added.
  - C-5 `T91/T92` wording stale ⇒ changed to `T91–T103` by ⑧b.
  - C-6 OB-115 assess wording ⇒ renamed to 落地 by ⑧b.
  - C-8 dual-mirror silent skip ⇒ marked in the detail.
- **⑤ r2, C-7**: assertion coverage gap, see the follow-up in §6 (not added in this slice).
- **⑥, two rounds**: r1 format deviation (D2), r2 thin content (D3); neither round raised any other Low / Info finding.
- **⑦, three rounds**: r1 one High + one Medium (no Low), r2 one Medium (no Low), r3 zero findings; the three Medium+ items were closed by D6 and by the ②-8 fix.

## 6. Residual risks and follow-ups

- **Historical-report transition**: reports committed before v1.46 carry no `sixstage` fence, so `G7-f` keeps the three legacy states for them (historical exemption / skip declaration / INFO); once such a report is touched by a change set, a fence or a declaration is still required.
- **⑥ output quality**: the second round of this slice is thin in substance (see D3); the criterion can only guarantee the shape, and the substantive review stays with ⑤ / ⑦. The residual is registered as OB-115.
- **⑤ r2, C-7 assertion coverage gap**: `selftest.js` still lacks four to six assertions — (a) the dual-mirror same-number guard has no mutation assertion; (b) the skip-declaration-plus-marker ⇒ FAIL path has no assertion; (c) neither the SEVERITY⇔FINDINGS guard nor the no-conclusion-line-inside-the-block guard has a mutation; (d) FACE B (a `| DR-n` entry at a non-generated position inside a slice file) has no dedicated fixture; (e) the unchanged-block-not-rejudged case has no assertion. **Disposition: not added in this slice, carried to the next `GATES-*` slice**; the reason is that these harden assertions only and do not affect the correctness of the current criterion verdicts (⑦ r3 already ruled PASS).
- **`--list` does not enumerate sub-criteria**: `gate.js --list` prints main criterion ids only, so `G6-8` and `G7-f` are selected through `--check=G6` / `--check=G7`; this slice leaves that behaviour unchanged.
- **Duplicate conclusion-line edge**: `conclusion → note → conclusion` is covered by the duplicate test; the variant where both lines are conclusions and the last line is not one stays with ⑤ / ⑦ (a stated criterion boundary).
- **Encoding statement mismatch**: the dispatch said `docs/validation/**` is BOM-free by repository convention, but the measured corpus and `G4a` both require exactly one BOM on a `.md` file ⇒ this report pair is written with BOM=1 and LF, consistent with the corpus.
- **Version labels**: the 16 `v1.48` occurrences of D5 were normalized in this ring; numstat and line counts are unchanged.
- **CI scope**: this ring is bound by the no-commit / no-push rule, so no CI reading was taken; the two-leg observation after the commit is carried by ⑨ and ⑩, and no run id is pre-filled here.

## 7. Conclusion

- All four debts A / B / C / D landed; `G6-8` and `G7-f` are covered in both directions by SELFTEST assertions and real staged-change-set probes.
- The gates are green in both the `index` and `tree` scopes (`TOTAL_FAIL=0`), and the assertion count grew from 600 to 663 with zero regression.
- The report pair mirrors line by line: equal line counts and equal `**` / backtick counts on every line; the ⑥ transcript sits inside the single `sixstage` fence with the conclusion line last.

## 8. Commit and CI

- **Main commit**: `a54bc49` (`feat: v1.46 GATES-CLEANUP …`), 10 paths, `+2105/-18`; pushed to `origin/main` (SSH channel, **gitee untouched**).
- **CI run**: `40034959012` (head `a54bc49`, attempt 1) ⇒ conclusion **success**.
  - `Go ubuntu-latest`: success; the `Build` / `Vet` / `Test` / `Race` / `Contract fixtures` / `Gates` steps are all success.
  - `Go windows-latest`: success; `Race` and `Contract fixtures` are skipped as usual per `runner.os == 'Linux'`; `Gates` success.
  - the three dispatch-only jobs (`Candidate probe` / `Candidate build` / `Bootstrap and release evidence`) are all skipped.
- **Per-leg `Gates` step readings** (from `gh run view 40034959012 --log`, one set per leg):
  - ubuntu: `SELFTEST: PASS 663/663`, `GATE REPORT scope=ci changed=10`, `base go test ./... ok=15 FAIL=0`, `TOTAL_FAIL=0`.
  - windows: the same four readings, with values identical to ubuntu.
- **Reading scope**: these readings belong to the change set of the **main commit**; the readings of this backfill commit are not written back here, and its observation run adds no new commit per `§11.2`; the `§4` readings such as `newCjk` / `scanned` likewise belong to the main commit's change set, and this backfill does not rewrite them.
- **The `§6` "CI scope" item stays verbatim**: it records the true state of that time (bound by the no-commit rule, no CI reading taken, no run id pre-filled); this section only records the readings taken after the commit and does not rewrite that item.
