# `DIRECTIVE-V21-OBS-SWEEP` Slice Validation Report (Directive v1.55: the five v21 observation items consolidated)

> Slice `DIRECTIVE-V21-OBS-SWEEP` (P33, directive-governance slice) | Tier `[SLICE]` | Level High (M tier) | Date 2026-10-28.
> Status: ✅ **Complete** (directive v1.55; ⑤ 1 round, ⑥ 1 round, ⑦ 1 final review; three scopes `TOTAL_FAIL=0`; `SELFTEST: PASS 697/697`; commit / push / CI run only after same-turn explicit authorization at ⑩).
> [中文](directive-v21-obs-sweep.md)

## 1 Change summary

**What this slice does**: closes out in one slice the five v21 observation items produced by P32, raising the directive from v1.54 to v1.55 — (1) Appendix B.7.1's kickoff-message template gains a "⑥⑦ input-package scope" slot (slot level only; the B.7 opening paragraph is untouched and `§7.3` is not restated); (2) `§11.3` gains an "8a self-check" subsection (four recorded-state requirements, no mechanical criterion); (3) the end of the `§6.3` recorded-state supplement gains a recorded-state exemption sentence for "text inside the report's §1 exception-registration area trips criteria", with the marker mechanism split out as `GATES-EXEMPT`; (4) `§5.1` gains the "boundary between remediation and evidence gathering" paragraph; (5) `§11.3` gains an "A-NN numbering integrity" subsection (binding new reports only). The header moves to v1.55, `§0.3` gets the change row, Appendix C gets `C.0ap`, `_EN` mirrors all of it in place, and the ledger writes back the P33 boundary row.

**The five rulings** (user-ruled group by group; ① produced the comparison table only and took no position):

| Group | Ruling | Landing | Nature |
|---|---|---|---|
| ①-1 | (i) | Appendix B.7.1 template gains the "⑥⑦ input-package scope" slot | slot level |
| ①-2 | (i) | `§11.3` gains the "8a self-check" subsection | recorded-state duty |
| ①-3 | (ii) + (iv) | `§6.3` recorded-state exemption sentence + the split-out slice `GATES-EXEMPT` | recorded-state duty + split |
| ①-4 | (i) | `§5.1` "boundary between remediation and evidence gathering" paragraph | explicit limitation |
| ①-5 | (i) | `§11.3` "A-NN numbering integrity" subsection | recorded-state duty |
| Boundary | extended to Appendix B.7.1 (excluding the B.7 opening paragraph) | ledger P33 boundary row | ledger sync |

**Version package**: the header `docs/DELIVERY_DIRECTIVE.md:3` reads v1.55; the `§0.3` change row is `:101`; Appendix C adds `C.0ap` at `:2064`; the `_EN` side mirrors all of it in place.

**What this slice does not do**: does not change `§2.2` / `§2.4` / the `§5.2` matrix / the `§7.5` figures; adds no mechanical criterion of any kind; does not touch `tools/gates/**`, `internal/**`, CONTRACTS / schema / fixtures / workflows, `DEV_PLAN{,_EN}`; does not push to gitee.

**One-line conclusion**: all five v21 observation items are disposed of at slot level or as recorded-state duties, with no mechanical criterion added; the marker mechanism is split out as `GATES-EXEMPT` under ruling ①-3(iv) and is not implemented here.

**Artifact list**: line-for-line identical to the `artifacts` fence at the end (directive pair, ledger pair, this report pair, freedom list).

## 2 Execution pipeline

| Step | Role | Carrier | Model (qualified name) | Status | Duration | Rework count |
|---|---|---|---|---|---|---|
| Step 0 | controller | — | — | done (three P33 ledger spots + the P32 stale-field correction) | history missing | 0 |
| ① architecture | architect | `deep-reasoner` | `DeepSeek V4 Pro (deepseek)` | done (five-group fork comparison table, conclusion line `ARCHITECT: OPTIONS READY`) | history missing | 1 (R1) |
| ruling | user | — | — | ruled (①-1(i) / ①-2(i) / ①-3(ii)+(iv) / ①-4(i) / ①-5(i) + boundary extension) | history missing | 0 |
| ② implementation | implementer | `prfrail-implementer` | `DeepSeek V4.1 Flash (deepseek)` | done (ten landing points) | history missing | 0 |
| ③ testing | tester | `prfrail-tester` | `DeepSeek V4.1 Flash (deepseek)` | done (consistency checks; 2 observations, 1 rejected, see A-3 in §3; case B is in the R = 4 registration at the end of §3) | history missing | 1 (R4) |
| ④ gates | controller | — | — | 5 runs (final readings in §8) | history missing | 1 (R2) |
| ⑤ pre-review | pre-reviewer | `deep-reasoner` | `DeepSeek V4 Pro (deepseek)` | r1 `PRE-REVIEW: PASS` (1 Low + 4 Info, the Low remediated) | history missing | 1 (R3) |
| ⑥ scan | independent scanner | `independent-reviewer` | `MAI-Code-1.1-Flash (copilot)` | done (`INDEPENDENT SCAN: PASS`, zero findings; the internal budget of 3 exhausted) | history missing | 2 (⑥ re-sends, not counted in R) |
| ⑦ final review | independent final reviewer | `independent-reviewer` | `GPT-5.3-Codex (copilot)` | done (`FINAL REVIEW: PASS`, no re-review triggered) | history missing | 0 |
| ⑧a documentation | documenter | `prfrail-documenter` | `DeepSeek V4.1 Flash (deepseek)` | done (this report's draft) | history missing | 0 |
| ⑨ native validation | controller | — | — | done (three scopes `TOTAL_FAIL=0`; `SELFTEST: PASS 697/697`; readings in §8) | history missing | 0 |
| ⑧b finalize documentation | documenter | `prfrail-documenter` | `DeepSeek V4.1 Flash (deepseek)` | this step | — | 0 |
| ⑧c closeout cleanup | controller | — | — | not done (follows this step) | history missing | 0 |
| ⑩ stop point | user | — | — | not done (commit / push / CI not authorized) | history missing | 0 |

- **Rework count R = 4**: R1 = ①'s first call misread the `§5.2` role matrix (it read "required" as a nonexistent "technical review" column), then the controller checked it mechanically and re-dispatched under `§9.1`; R2 = ④'s first tree run made `G3(b)` report `MISMATCH` (the `_EN` side used a plural form that the singular pattern does not match), then both sides were rewritten into the equal `CI run <id>` form and the re-run passed; R3 = ⑤ r1 raised 1 Low, and the gates re-ran green after remediation; R4 = ③ case B (the transcript holds two different ③ dispatch prompts, `sha ede7fe40` vs `437a0dee`, each ending with `TEST: DONE`, which is a substantive rework). Item-by-item grounds, the OB-37 pre-authorization behind R = 4, and why ⑥'s re-sends are excluded from R are in A-1 of §3 and in the R = 4 registration at the end of §3.
- **R=4 triggers the `§5.2` "R≥2 high-risk" row**: ①⑤⑥⑦⑨ are all mandatory in this slice anyway, so the actual flow is unchanged (see A-2 of §3).
- **Duration column**: this slice has no per-step timing readings, so under Discipline B every cell says "history missing"; nothing is estimated or invented.

**Item (2) landing list (ten landing points, for ⑦'s point-by-point check)**:

| # | Landing | CN | EN |
|---|---|---|---|
| 1 | header v1.55 | `:3` | `:3` |
| 2 | `§0.3` v1.55 change row | `:101` | `:101` |
| 3 | `§5.1` boundary sentence | `:511` | `:511` |
| 4 | `§6.3` recorded-state exemption | `:670` | `:670` |
| 5 | `§11.3` structure-chain two items | `:952` | `:952` |
| 6 | `§11.3` "8a self-check" subsection | `:962` | `:962` |
| 7 | `§11.3` "A-NN numbering integrity" subsection | `:964` | `:964` |
| 8 | Appendix B.7.1 slot | `:1360` | `:1360` |
| 9 | Appendix C.0ap | `:2064` | `:2064` |
| 10 | ledger P33 boundary row | `REMAINING_SLICES.md:1129` | `_EN.md:1116` |

The line numbers above are measured values after landing; the directive pair is `+24/-3` and the ledger pair `+25/-1`, equal on both sides. This report lists itself in a `repo` disposition row, so it must be `git add`-ed before ⑨ or G7-b judges it red.

## 3 Exceptions and fallback record

**A-1 (item-by-item grounds for R)**: the item-by-item grounds are in the table below; the ⑥ re-sends are registered truthfully, but per `§7.2.2` they consume ⑥'s internal budget and do not enter the global R of `§9.2`; the make-up of R = 4 and the OB-37 pre-authorization are in the R = 4 registration at the end of §3.

| Event | Counts in R | Ground |
|---|---|---|
| R1 ①'s first call misread the `§5.2` role matrix (it read "required" as a nonexistent "technical review" column), then the controller verified mechanically and re-dispatched under `§9.1` | counted | `§9.2` "controller fallback intervention" |
| R2 ④'s first tree run made `G3(b)` report `MISMATCH` (the `_EN` side used a plural form that the singular pattern does not match), then both sides became the equal `CI run <id>` form and the re-run passed | counted | `§9.2` "fix-and-rerun after discovery" |
| R3 ⑤ r1 raised 1 Low (the ledger's "observation A landing declaration" collided literally with the new `§5.1` boundary sentence), remediated, then the gates re-ran green | counted | `§9.2` "fix-and-rerun after discovery" |
| ⑥ re-send  first output lacked the Section A–E structure, so it was re-dispatched; the first re-dispatch failed on the context window and succeeded after trimming | **not counted** | `§7.2.2` "⑥ cost cap … a re-send is limited to 'truncation / exception' and counts against the cap" — it consumes ⑥'s internal budget (3 allowed, 3 used) |
| R4 ③ case B: the transcript holds two different ③ dispatch prompts (`sha ede7fe40` vs `437a0dee`), each ending with `TEST: DONE`, which is a substantive rework | counted | `§9.2` "fix-and-rerun after discovery" (same class as R2 / R3) |

**A-2 (R=4 trips the high-risk row)**: R=4 hits the "R≥2 high-risk" row of `§5.2`; ①⑤⑥⑦⑨ are all mandatory in this slice anyway, so the actual flow is unchanged.

**A-3 (③ tester's observation 2, rejected)**: the tester reported a half-width space at `docs/DELIVERY_DIRECTIVE.md:511` between the closing punctuation of "修复动作）" and "本款不排除"; the controller enumerated every half-width space on that line, found eleven at positions 22, 49, 82, 103, 105, 132, 134, 141, 143, 148 and 150, and found character 71 (the full-width semicolon) adjacent to character 72 (`本`) with no space at all — all eleven are ordinary separators (after `⑧c 后`, after the equals sign, at `⑩ `, at ` → ` and the like) ⇒ judged a **false positive**, no file changed. The rejection **does not bar ⑦ from raising the same class of point again**; a fresh finding by ⑦ requires this rejection to be re-assessed.

**A-4 (the controller's own state failure, registered)**: after ② landed, one round of output carried nothing but placeholders and did not advance the pipeline (self-corrected inside the same round; `git status` proves ② had in fact landed). Nature = executor output degradation; disposition = self-correction inside the same round, no D violation produced; registered as a v22 observation candidate ("the controller's output degrades — when to halt and when to retry must be defined"); a recurrence in the next two or three slices escalates it to dedicated treatment.

**A-5 (⑥ format deviation, registered across slices)**: P32's ⑥ deviated once and this slice's ⑥ deviated once (no Section A–E on the first output) ⇒ judged systemic; this slice **registers only and does not dispose of it**, leaving v22 to evaluate whether a "Section A–E precheck" belongs in the task-package template.

**A-6 (circular-dependency handling = option (a), the most conservative wording)**: this slice chose (a): every piece of text in report §1 that could trip G2 / G6-4 / G4b is rewritten as a code-point expression or as a synonym that keeps the original meaning; quoted forbidden strings are written in code-point form, and **no meaning is dropped**.

**Temporary-artifact count D (reading before 8b)**: it read zero before 8b — every temporary artifact of this slice stayed under `tmp/DIRECTIVE-V21-OBS-SWEEP/`, nothing was left outside the repo or at the repo root, and nothing was cleaned early; at the 8b point in time one new occurrence was registered, and the D = 1 registration at the end of §3 governs.

**Scope of rewriting**: ⑥'s original text is kept verbatim inside this report's `sixstage` fence with zero rewriting; the ⑧a prose outside the fence follows the A-6 wording rule.

**(1) ② review conclusion and the rejected correction suggestions**

- ② was re-dispatched to `prfrail-implementer` + `DeepSeek V4.1 Flash (deepseek)` as a review pass (in this slice the ② writing was done by the controller first; the dispatched review satisfies the kickoff message's hard requirement of "at least 1 ② call").
- Review result: all 10 landing points are **faithful** to the five rulings, with no drift and no omission; the CN / `_EN` ranges of the 9 hunks match one by one; both directive files have the same line count of 2142, and all 10 landing points sit on identical line numbers.
- ② raised 2 further correction suggestions, and the controller **rejected each one**, with the grounds below:

| ② suggestion | Controller ruling | Grounds for rejection |
|---|---|---|
| `DELIVERY_DIRECTIVE.md:101` (the v1.55 entry in §0.3) gains "and registers **OB-127** in §12.4" | **rejected** | measured: `:100` (the v1.54 entry) has `hasOB126=false` — the previous slice's §0.3 change row omits the OB number too; ② cited the older v1.52 convention, while the nearer v1.54 and C.0ao follow the "§0.3 records the directive change only, the OB number goes to §12.4" caliber. Adopting it would break consistency with the nearest precedents |
| `DELIVERY_DIRECTIVE.md:2068` (C.0ap) gains "and registers OB-127 in §12.4 in sync" | **rejected** | same as above; C.0ap must follow the same caliber as the §0.3 change row, and C.0ao lists no OB number either |

- Because the corrections were rejected, the ② review **does not amount to "correction adopted + gates re-run"** ⇒ no extra +1 on R.

**(2) R = 4 registered (OB-37 pre-authorization)**

- R = 4 = R1 (①'s first call misread the `§5.2` role matrix, reading "required" as a nonexistent "technical review" column → 10 mechanical checks by the controller plus a `§9.1` re-dispatch) + R2 (④'s first tree run reported a `G3(b)` run reading of `1/0` MISMATCH because the EN side used the plural `runs`, which the `\brun \d+\b` pattern does not match → both sides became the equal `CI run <id>` form and the re-run passed) + R3 (⑤ r1 Low: the ledger's "observation A landing declaration" collided literally with the new `§5.1` boundary sentence → gates re-run after remediation) + R4 (③ case B: the transcript holds two different ③ dispatch prompts, `sha ede7fe40` vs `437a0dee`, each ending with `TEST: DONE`, which is a substantive rework).
- Routed through the **OB-37 meta-rule pre-authorization** (the user ruled in the same round of 2026-10-28: "if the final R ≥ 4, always take the OB-37 pre-authorization").
- Four elements: ruling party = user; date = 2026-10-28; reason = the slice has already spent substantial cost, ⑨ is fully green and 8b has not run, so halting is expensive; basis = the user's reply in this round, verbatim.
- **This item must not serve as a precedent**.

**(3) D = 1 registered (the user ruled to continue directly; the registration carries one level more)**

- Facts: during the review ② created `_p33_en_diff.txt` at the **`tmp/` root** (not under `tmp/<slice id>/`) and deleted it with `Remove-Item` inside the same command chain; the residue check reads zero (the only non-cache item at the `tmp` top level is the `DIRECTIVE-V21-OBS-SWEEP` directory).
- Tier: the self-deletion preceded the controller's reading ⇒ under `§5.1` this is the **non-self-healing tier**.
- **The one-level-stronger wording (verbatim)**: "this slice's D=1 breaks not only the default tier of `§5.1` but **also the hard constraint of this slice's kickoff message §七 二 item 8: 'this slice enforces the temporary-artifact discipline strictly: any violation is handled at the `§5.1` non-self-healing tier (halt and escalate), and the `§5.1` self-healing tier does not apply to this slice'** — a direct ruling by the user in the same round. This slice's exception **must not serve as a precedent**; a same-class violation in the next slice is still handled by that slice's own kickoff temporary-artifact discipline (this slice's exception does not generalise)."
- Four elements: ruling party = user; date = 2026-10-28; reason = ②'s violation is a first occurrence with zero residue and no intent; the slice is already at R=4, ⑨ is fully green and 8b has not run, so halting is expensive; and being a governance slice, halting would mean "the governance slice fails its own goal"; basis = the user's reply in this round, verbatim ("authorize continuing 8b").
- Independent trace = this paragraph (no separate OB / ADR is raised).

## 4 Environment and tier

- **Tier**: High (the fixed controller tier of `§8.2`); probe budget 0 (no real external call was made).
- **Toolchain**: Node (`gate.js`), Go 1.22+ (`gofmt` / `build` / `vet` / `test`).
- **Operating system**: Windows (local).
- **Baseline**: `HEAD` = `dbf5298` (directive v1.54); P32 is closed.
- **Change surface**: only the directive pair and the ledger pair, plus this report pair and the freedom list; `internal/**` and `tools/gates/**` are untouched.
- **Review isolation**: ⑤ full package / ⑥ blind scan package / ⑦ blind final-review package, three separate sets (see §6).

## 5 Review conclusion summary (⑤⑥⑦)

- **⑤ r1**: carried by `deep-reasoner`, verdict `PRE-REVIEW: PASS` — nothing at Medium or above; 1 Low (the ledger's "observation A landing declaration" collided literally with the new `§5.1` boundary sentence → remediated: both sides now read "a remediation action must not rebuild; only when ⑩ evidence gathering needs a rebuild, follow the four actions of backup path two") + 4 Info (first, the reading for CI run `40550644240` needs a `gh run view` re-check; second, the P33 ledger's status and fix direction still stand at ①'s point in time and are advanced at close-out; third, the deliverables row says "Appendix B" while the boundary row says "Appendix B.7.1"; fourth, the P32 status-row update passes factual verification).
- **⑥**: carried by `independent-reviewer`, verdict `INDEPENDENT SCAN: PASS` (zero findings; the full Section A to Section E text is in the fence below; its internal budget of 3 is exhausted, see A-1 of §3 and §10).
- **⑦ r1**: carried by `independent-reviewer`, verdict `FINAL REVIEW: PASS` — all eight dimensions pass, with no findings at all (High / Medium / Low / Info all zero); no re-review triggered.

⑥'s final message is transcribed verbatim under the `§11.3` transcription contract (the fence is the authoritative transcription; only notes may sit outside it):

```sixstage
## Section A
INDEPENDENT SCAN: PASS — [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md#L511)

## Section B
无发现

## Section C
- v1.55 进入点与本轮五项处置均已在版本日志里落地： [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md#L101-L101) 与 [docs/DELIVERY_DIRECTIVE_EN.md](docs/DELIVERY_DIRECTIVE_EN.md#L101-L101)
- §5.1 的“整改与取证的划界”已写入，并且明确了“不重建临时目录”和“备选路径二”的顺序： [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md#L511-L511)
- §6.3 的记录态豁免句已落地，并说明该类触发不设机械判据、由人工核验： [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md#L670-L670)
- §11.3 的“8a 自查”与“A-NN 编号完整性”已并列落位，且都写明为记录态、手工核验： [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md#L951-L964)
- B.7.1 的起跑消息槽位“⑥⑦ 输入包范围”已落地，且仅为槽位层、不重述 §7.3： [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md#L1357-L1360)
- 这五处在中英双语中保持同位一致： [docs/DELIVERY_DIRECTIVE_EN.md](docs/DELIVERY_DIRECTIVE_EN.md#L511-L511), [docs/DELIVERY_DIRECTIVE_EN.md](docs/DELIVERY_DIRECTIVE_EN.md#L962-L964), [docs/DELIVERY_DIRECTIVE_EN.md](docs/DELIVERY_DIRECTIVE_EN.md#L2064-L2074)

## Section D
- 删除“⑥⑦ 输入包范围”槽位后：不会触发当前自动门禁的红灯；它是模板槽位，不是机械判据，缺失只能在人工审查里发现。该不变量已经被设计为“slot-only + 记录态”，因此不属于自动判红项。
- 删除“8a 自查”小节后：不会因 G7-f 自动红；它被正文写明为“报告结构要求（记录态）”、“不设机械判据”。因此它既不属于 G6-5，也不属于 G7-f 的机械匹配对象。
- 删除“A-NN 编号完整性”小节后：不会因 G6-5 自动红；正文已明确“本条只约束新增报告，且不并入 G6-5”。该条属于人工校验的完整性义务，而非可脚本证明的结构化判据。
- 未覆盖的不变量：记录态要求、slot-only 模板扩口、以及历史报告不追溯改写，都是一类“有义务但无机械判据”的字段；这个不变量是故意保留的，不是漏项。若要把它们转成机械门禁，就必须另开规则修订，不在本片内做“看起来像判据”的伪机械化。
- 这意味着本轮的可证伪性是：删掉任一条上述句子，人工复查会红；自动门禁仍然不会因这个句子消失而红。该结论与本片写法一致，也正是该片对 v21 观察项的处理方向。

## Section E
- 我无法验证的部分：我没有对实际 runtime gate 进行二次执行，也没有对第三方自动门禁实跑；本次只能在只读约束下核对 diff、正文与中英镜像内容。我的结论仅基于这份材料的文本证据与位置证据。
- 什么会推翻我的 PASS：如果任一关键落点在工作树中缺失、或 CN/EN 发生实质不一致、或如果后续证据显示“记录态豁免”并未真正覆盖 G2/G6-4/G4b 的触发场景，那么这份 PASS 会被推翻。反过来说，若本轮的文案与对应条文一致，且该五条观察项的落点都在规定范围内，那么当前判定应保持 PASS。

INDEPENDENT SCAN: PASS — [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md#L511)
```

**Transcription normalization note**: ⑥'s five Section lines arrive as plain text without a `#` prefix; per the `§11.3` requirement that they appear as level 2–4 headings starting with `Section X`, they became `##` headings and their line endings became LF. Markdown links, U+201C / U+201D quotation marks and full-width punctuation are preserved verbatim; not one character or word inside the fence was changed.

### ⑧a Self-check (first application of v1.55)

| Item | Conclusion | Evidence path |
|---|---|---|
| ① mechanical audit of the cost table | present | every row of the §10 cost table corresponds to a `runSubagent` `agentName` / `model` record; the ⑩ receipt attaches the raw stdout of the mechanical-reading script |
| ② "resolved" items carry file:line | present | every remediation / finding in §3 and §5 carries a file:line |
| ③ old-wording grep before delivery | present | a grep for the three old wordings v1.54 / R=4 / "technical review"; every hit is a legitimate version reference, detailed in the note line at the end of this subsection |
| ④ every Info / Low of each ⑥ / ⑦ round lands in the report | present | ⑥'s Section B says "no findings"; ⑦ reports none; ⑤'s 1 Low and 4 Info are in the first item of §5 |

- ③'s measured reading: `v1.54` has 5 lines (the version-migration sentence in §1, the baseline sentence in §4, item 1 of §7, plus the ③ row and this note line), all legitimate references to a version before or after the upgrade; `R=4` has 2 lines (the ③ row and this note line), which is the search term referring to itself; "technical review" has 4 lines on the CN side (§2's R1 line, §3's A-1 line, the ③ row and this note line), all describing the column that was misread rather than reusing the old wording, while the EN side renders it as technical review and therefore has no literal hit.

### A-NN numbering integrity (first application of v1.55)

- The entity definitions of this report live in §3 as A-1 through A-6, each carrying its own heading definition (the `**A-N（…）**` form).
- No dangling anchor: every `A-NN` reference in the text is A-1 through A-6, each defined in §3; this report cites no anchor from the P32 report.

## 6 Review input package summary (with the blind-review isolation proof)

- **⑤**: full package (slice definition + the user's five-group ruling record + full diff + contract excerpts + ④ gate readings); per `§2.2` R2.6 anti-anchoring it excludes ①'s reasoning.
- **⑥**: blind scan package, holding only the slice fields, the full diff, the contract excerpts (before / after) and the read-only constraint; the ruling record and ⑤'s findings are removed.
- **⑦**: full material package, holding the raw gate output (`gate-tree.txt` / `gate-index.txt` / `gate-ci.txt` / `gate-selftest.txt`) and the freedom list; it excludes the ⑤ / ⑥ lists and conclusions.
- **Checklist attached in the first round**: ⑤'s first round carried ①'s comparison table (without the reasoning); ⑥'s and ⑦'s first rounds both carried **no** other reviewer's list, so blindness holds.
- **Input sanitization**: all three sets contain only in-repository paths and slice fields — no IP, user name, host name or credential; every reviewer is read-only.
- **Isolation conclusion**: ⑥ and ⑦ owe their blindness to the absence of any other reviewer's list; ⑤ and ⑦ come from different model families (`DeepSeek V4 Pro` and `GPT-5.3-Codex`), so independence holds.

## 7 Falsifiability (mechanical checks + mutation test)

1. Revert the `§0.3` v1.55 entry to v1.54 ⇒ G1-b goes red (header version disagrees with the last changelog row).
2. Delete any new `_EN` line ⇒ G3(a.1) symmetry goes red.
3. Delete any Section heading inside the `sixstage` fence ⇒ G7-f goes red (the A→B→C→D→E sequence).
4. Append any line after the conclusion line inside the fence ⇒ G7-f goes red (nothing may follow the conclusion line).
5. A new report without an `artifacts` fence ⇒ G7-a goes red; a fence row declared `repo` whose path does not resolve ⇒ G7-b goes red.
6. Does not go red (registered honestly): deleting the `§11.3` "8a self-check" subsection on both sides at once keeps the numstat symmetric and trips no criterion, leaving only ⑤ / ⑦ human guarding; this is a direct consequence of the slice's "recorded state (no mechanical criterion)" positioning, not an omission.
7. Also does not go red (registered honestly): deleting the `§6.3` recorded-state exemption sentence, the `§5.1` boundary sentence or the Appendix B.7.1 slot row is verified by humans only; ⑥'s Section D confirmed this independently.
8. Mutation test: this is a documentation-governance slice with no code mutation surface ⇒ recorded as a non-equivalent substitute per `§7.2.1` (iii).

## 8 Gate results

- `gofmt -l .` / `go build ./...` / `go vet ./...` / `go test ./...`: this slice changes `.md` files only, so the Go side is **not applicable** (registered reason: the change type is documentation-only).
- **④ point-in-time readings**: `gate.js --all --scope=tree` `TOTAL_FAIL=0`; `--scope=index` `TOTAL_FAIL=0`, `G6-1 scanned 98`, `G3(a)` shows the directive pair `+24/-3` against `_EN +24/-3` and the ledger pair `+25/-1` against `+25/-1`, `G3(b)` reads `run` at `2/2 OK` and `hash` at `3/3 OK`; `--scope=ci` `TOTAL_FAIL=0`, `G6-1 scanned 4`; `--selftest` `SELFTEST: PASS 697/697`.
- **④ execution count**: 5 times (precheck after step 0 / precheck after ② / ④'s first red run / ④'s green re-run / 1 after the ⑤ Low remediation).
- **⑧a re-test (tree scope; this report pair and the freedom list are this step's new artifacts; 2 extra runs: the first test and the re-test after the fix)**: `TOTAL_FAIL=1`, and the only failure is G7-b — the three new files are not yet in the index, which the controller closes with a precise `git add` before ⑨, so it is not a documentation defect; everything else is green: `G2` no hit, `G1-b` v1.55 on both sides, `G3(a)` this report pair `+244/-0` against `+244/-0`, `G3(b)` all three pairs OK, `G4b` reports an empty `unknown`, `G7-a` sees a fence on both reports, `G7-e` and `G7-f` pass.
- **Six shapes**: this report pair's line count / heading line numbers / bold markers / pipes / blank lines / backticks are equal on both sides (measured in the same re-test).
- **⑨ point-in-time readings (after `git add`)**: `--scope=index` `TOTAL_FAIL=0`, `G6-1 scanned 586`, `G7-b` 14 artifact rows OK, `G7-f` PASS; `--scope=tree` `TOTAL_FAIL=0`, `G6-1 scanned 586`; `--scope=ci` `TOTAL_FAIL=0`, `G6-1 scanned 4`; `--selftest` `SELFTEST: PASS 697/697`. `G3(a)` shows the report pair `+244/-0` against `_EN +244/-0` OK, and all six shapes lines/head/bold/bars/blank/backticks read `true`.

## 9 ⑧c closeout cleanup record

- `tmp/DIRECTIVE-V21-OBS-SWEEP/` is this slice's only temporary directory; after this report is written the controller removes it item by item and reads the result back.
- Temporary-artifact count D: registered as 1 at the 8b point in time (see the D = 1 registration at the end of §3); the read-back after the 8c cleanup is zero residue.
- Encoding check of this slice's new files: UTF-8 with BOM + pure LF; 1 BOM each, 0 CR bytes.
- The repository-wide encoding and line-ending sweep runs after this report is written; the existing frozen-evidence exception set under `docs/validation/evidence/` is unchanged.
- `tmp/gates-cache/g4b/` belongs to the runtime-cache exception set and is left untouched.
- Post-cleanup read-back: `Test-Path tmp\DIRECTIVE-V21-OBS-SWEEP` is `False`; `git status --porcelain -uall` shows no unexpected untracked path.

## 10 Cost and accounting

| Stage | Carrier | Model | Count |
|---|---|---|---|
| ① | `deep-reasoner` | `DeepSeek V4 Pro (deepseek)` | 2 (first misread call terminated + successful `§9.1` re-dispatch) |
| ② | `prfrail-implementer` | `DeepSeek V4.1 Flash (deepseek)` | 1 (review pass) |
| ③ | `prfrail-tester` | `DeepSeek V4.1 Flash (deepseek)` | 2 (per the transcript; case B) |
| ⑤ | `deep-reasoner` | `DeepSeek V4 Pro (deepseek)` | 1 (r1 PASS; budget 1+1, 1 used) |
| ⑥ | `independent-reviewer` | `MAI-Code-1.1-Flash (copilot)` | 3 (one format violation / one context overflow / one success; budget 3/3 exhausted) |
| ⑦ | `independent-reviewer` | `GPT-5.3-Codex (copilot)` | 1 (r1 PASS; budget 1+3, 1 used) |
| ⑧ | `prfrail-documenter` | `DeepSeek V4.1 Flash (deepseek)` | 2 (8a + 8b) |
| batch helper | — | — | 0 |
| probe | — | — | 0 |

- **Mechanical-audit caliber**: the table above corresponds row by row to the `runSubagent` `agentName` / `model` records; the ⑩ receipt attaches the raw stdout of the mechanical-reading script.
- The ④ gates and the ⑨ native validation run on the controller side and count no model call.
- **② role-split deviation, registered**: the ② writing itself was done by the controller first and then reviewed by a dispatched `prfrail-implementer`; that breaches the kickoff message's "at least 1 ② call, zero dispatch = halt and escalate", and it was closed after the user ruled on 2026-10-28 to pick (A), the dispatched review pass.

## 11 Explicitly unexecuted items

- commit, push and CI observation (require same-turn explicit authorization after the ⑩ stop point). **Completed after the ⑩ stop point**: commit `06231e1`; CI run `40608643891`, both the ubuntu and windows jobs succeeded (per-leg `G6-1 scanned 650`, `SELFTEST: PASS 697/697`, `TOTAL_FAIL=0`).
- Real probes (budget 0).
- Batch helper calls (0 times).
- `DEV_PLAN{,_EN}` write-back (a governance slice's authoritative ledger is Appendix C, so per the OB-12 caliber of `§11.2` nothing is written back).
- Implementation of the `GATES-EXEMPT` split slice (this slice only rules on the split; the marker mechanism is not implemented here).
- Mechanical criteria for the `§11.3` 8a self-check and A-NN integrity (①-2(ii) / ①-5(ii) were ruled deferred).

## 12 Precedent-ban cross-check

- This slice **does not rely on** any clause marked "must not serve as a precedent"; the two exceptions this slice raises itself (the OB-37 pre-authorization for R = 4, and D = 1) are marked "must not serve as a precedent" verbatim at the end of §3.
- Citation check: the slice cites neither `§1.7.4` ④, nor relies on the OB-37 meta-rule ③ of `§7.8.4`, nor any OB row of the `§12.4` ledger whose disposition cell carries that marker, nor any same-class row in `ADR_REGISTER`; R = 4 and D = 1 take the meta-rule's **independent registration path** (four elements at the end of §3) and do not invoke its item ③.
- The slice's exceptions number exactly two and both were ruled by the user in the same round: R = 4 continues under the OB-37 pre-authorization and D = 1 continues as a non-self-healing-tier exception (both registered at the end of §3 and both marked as not serving as a precedent); ⑤'s 1+1 execution budget was not exhausted and ⑦ triggered no re-review, so neither was expanded.
- The `GATES-EXEMPT` split is registered as a ruling only and is no source of exception inside this slice.

## 13 Next-step recommendations

- The next slice should start `GATES-EXEMPT` under ruling ①-3(iv): give the report's §1 exception-registration area a marker mechanism (a `gate-exempt` region) so that G2 / G6-4 / G4b skip inside it; that slice must touch `tools/gates/**` and therefore dispatches ⑦ as a governance slice.
- A later slice should evaluate the mechanical-criterion shapes of ①-2(ii) and ①-5(ii) (the four 8a self-check items, A-NN numbering integrity); this slice ruled them deferred.
- v22 should evaluate two cross-slice observations: a Section A–E precheck for ⑥ (the systemic deviation registered as A-5 here), and a halt-and-retry criterion for controller output degradation (registered as A-4 here).
- Role-tier recommendation: ①⑤ stay on `DeepSeek V4 Pro (deepseek)`; ②③⑧ stay on `DeepSeek V4.1 Flash (deepseek)`; ⑥ on `MAI-Code-1.1-Flash (copilot)`; ⑦ on `GPT-5.3-Codex (copilot)` in the regular tier, not downgraded; the OB-37 meta-rule is used in this slice only for the exception registrations of R = 4 and D = 1 (see the end of §3) and for no model downgrade.

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/directive-v21-obs-sweep.md	repo
docs/validation/directive-v21-obs-sweep_EN.md	repo
docs/validation/evidence/DIRECTIVE-V21-OBS-SWEEP-freedom-list.md	repo
```
