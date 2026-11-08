DESIGN: DONE

# P62 R08-R09-RULING-REVISION — ① design-r1（判据设计）

> 性质：① 判据设计产出（只读，本稿不落盘）；主控逐字落盘后交 ② 按此实现。所有目标 `.md` 落盘保持 UTF-8 with BOM + LF（仓库硬规则）。

## 段 A：R-08.5 增补块「继承漂移豁免的适用条件」

插入锚点（CN/EN 同位）：`r08-r09-ruling{,_EN}.md` R-08.5 小节，第三条 bullet 之后、`### R-08.6` 之前。不动既有三行 bullet 与标题；不加新小节号；新增块 CN 与 EN 各 8 行、逐行同位。

### CN 块（8 行，逐字）

```markdown
- **继承漂移豁免的适用条件（五要素，须同时满足）**：
- ① **同基线**：对照基线 = B4 基线存档（`docs/validation/t027-at23-e2e.md` + `docs/validation/evidence/directive-history/b4/archive/`）；只允许以该存档为对照，不得另选或混用其它基线。
- ② **同断言**：断言 id、断言含义、断言所在 harness（`tools/b4-e2e/`）与 B4 基线均不变；harness 有改动或断言被重定义即不满足。
- ③ **漂移方向更严格**：断言值变化方向更保守/更严格（如 `process-tree-taken-away` 从 `true` 变 `false` = 未检测到进程树被取走 = 判定更保守）；反向变化（如 `false` 变 `true` = 更宽松）不满足。
- ④ **四态不变**：`verdictClass` / `chainState` / `taskState` / `stepState` 与 B4 基线逐字一致；任一字段漂移即不满足。
- ⑤ **成因可追溯**：漂移成因可追溯到 B4 之后的特定产品改动（如 P61 的 `02d5797`，A5 pid reports 绑定为第六枚冻结事实）；无法归因到特定改动即不满足。
- **豁免后果**：五要素同时满足 ⇒ 豁免「前提①不满足」的判定，前提①视为满足（继承漂移豁免成立）；任一不满足 ⇒ 不豁免，按本小节既有判据处理。
- **豁免边界**：豁免成立不等于 AT-23 条件闭合——真实候选腿的前提②仍须满足；豁免须由 **⑦（独立终审）** 判定，主控不得自行判定。
```

### EN 块（8 行，逐行同位镜像）

```markdown
- **Applicability conditions of the inherited-drift exemption (five elements, all must hold simultaneously)**:
- ① **Same baseline**: the comparison baseline = the B4 baseline archive (`docs/validation/t027-at23-e2e.md` + `docs/validation/evidence/directive-history/b4/archive/`); only this archive may serve as the comparison target — no other baseline may be chosen or mixed in.
- ② **Same assertion**: the assertion id, the assertion meaning and the assertion's harness (`tools/b4-e2e/`) are all unchanged versus the B4 baseline; any harness change or any redefinition of the assertion fails.
- ③ **Drift direction stricter**: the assertion value changes in the more conservative / stricter direction (e.g. `process-tree-taken-away` from `true` to `false` = no process-tree takeover detected = a more conservative verdict); the reverse direction (e.g. `false` to `true` = more lenient) fails.
- ④ **Four states unchanged**: `verdictClass` / `chainState` / `taskState` / `stepState` are verbatim-identical to the B4 baseline; any drift in any of these fields fails.
- ⑤ **Cause traceable**: the drift cause is traceable to a specific product change after B4 (e.g. P61's `02d5797`, the A5 pid-reports binding as the sixth frozen fact); a cause that cannot be attributed to a specific change fails.
- **Consequence of the exemption**: all five elements hold simultaneously ⇒ the "premise ① unmet" verdict is exempted and premise ① is treated as satisfied (the inherited-drift exemption is established); any element fails ⇒ no exemption, and the existing criteria of this subsection apply.
- **Exemption boundary**: the exemption does not equal AT-23 conditional closure — premise ② of the real candidate leg must still hold; the exemption must be adjudicated by **⑦ (independent final review)**, and the controller must not adjudicate it on its own.
```

## 段 B：R-08.1 事实修正（§8 表行 + 时点注记）

### B-1 表行替换（CN/EN 各 1 行）

CN 旧行：`| R-08.1 封装函数 | 中 | 已落地代码，要过审查链 |`
CN 新行：`| R-08.1 封装函数 | 中 | **P61 ② 新建**（Invoke-CopilotCli / buildCopilotArgs），要过审查链 |`

EN 旧行：`| R-08.1 wrapper function | Medium | Landed code; must pass the review chain |`
EN 新行：`| R-08.1 wrapper function | Medium | **newly created by P61 ②** (Invoke-CopilotCli / buildCopilotArgs); must pass the review chain |`

### B-2 时点注记（CN/EN 各 1 行）

确切放置位置：在 §8 表末行（`| R-09 从放宽改严格（AT-24 已签收后） | 高 | 重跑 AT-24 |` / `| R-09 loosen → strict (after AT-24 signed) | High | Re-run AT-24 |`）之后的空行之后插入注记行，注记行后再保留一个空行，再接 `**原则**` / `**Principle**` 段。即结构为：表末行 → 空行 → 注记行 → 空行 → 原则段；CN/EN 完全同位。不改表其他行。

```markdown
> **时点注记（P62，2026-11-08）**：本行原写「已落地代码」，实况是 P61 ② 才新建（commit 6bb23a1 之前仓库无此函数）。按 §11.1 历史行不改写，本行加注记。事实错误独立于设计决策（R-08.1 裁定本身不变）。
```

```markdown
> **Point-in-time note (P62, 2026-11-08)**: this row originally read "landed code"; in reality the function was only newly created by P61 ② (before commit 6bb23a1 the repository had no such function). Per §11.1 historical rows are not rewritten, so this row gains a note. The factual error is independent of the design decision (the R-08.1 ruling itself is unchanged).
```

## 段 C：判定主体句 + 附注

判定主体句（段 A 末行尾部，最终措辞）：

> 豁免须由 **⑦（独立终审）** 判定，主控不得自行判定。

报告引用用一句话（⑦ 判定时点）：

> ⑦ 判定豁免的时点：在后续 AT-23 复核切片的 ⑦ 终审环节——独立终审员对照五要素逐项核验后，于 `FINAL REVIEW` 结论行作出「豁免成立 / 不成立」的判定；②~⑥ 各前序环节一律不得作出豁免判定，主控在 ⑧ 文档装配只能引用 ⑦ 的判定结论。

## 额外（供 ② 使用）

### §3.3 P61 报告注记（CN + EN 同位，新节「后续指针」）

放置位置：`at23-restricted-candidate{,_EN}.md` 文末 artifacts 围栏闭合 ``` 之后追加（报告无「下一步建议」节，按起跑消息 §3.3 改为末尾新加一节；不改报告原有内容）。结构：空行 → 节标题 → 空行 → 注记行，CN/EN 同位。

CN 节标题：`## 9. 后续指针（P62 追加）`
EN 节标题：`## 9. Follow-up pointer (appended by P62)`

```markdown
> **时点注记（P62，2026-11-08）**：R-08.5 判据已于 P62 修订，增加「继承漂移豁免」五要素；P61 的前提①判定不变（当时 R-08.5 未修订），但后续 AT-23 复核可按 P62 修订后的判据评估。
```

```markdown
> **Point-in-time note (P62, 2026-11-08)**: the R-08.5 criterion was revised in P62 to add the five elements of the "inherited-drift exemption"; P61's premise-① verdict is unchanged (R-08.5 was not yet revised at that time), but a later AT-23 re-check may be evaluated under the P62-revised criterion.
```

### §3.4 s1-exit 括注（CN + EN，AT-23 行状态单元格内）

触发判定：AT-23 行状态单元格含「前提①因继承漂移不满足」字样 ⇒ 按起跑消息 §3.4 执行（非「不适用」）。括注加在该单元格「仍 BLOCKED」/「still BLOCKED」之后、单元格末尾（即 `[P61 报告]` 链接列之前）。

CN 新单元格（整格替换）：

```markdown
| AT-23 | P61 评估完成（PARTIAL）：前提②满足、前提①因继承漂移不满足，待 R-08.5 修订片裁定豁免；仍 BLOCKED（R-08.5 判据已于 P62 修订，含继承漂移豁免五要素；AT-23 现状待复核） | [P61 报告](at23-restricted-candidate.md)：机制腿与 B4 基线 0 漂移（仅 F2 一处断言继承漂移）、真实候选腿实际跑过启动 / 执行 / 停机并记录 unsupported，不构成真实候选兼容性证明；[T026 报告](t026-agent-runner-contract.md) 的历史结论不变 |
```

EN 新单元格（整格替换）：

```markdown
| AT-23 | P61 assessment done (PARTIAL): premise ② met, premise ① unmet because of the inherited drift, awaiting the R-08.5 revision slice's waiver ruling; still BLOCKED (the R-08.5 criterion was revised in P62, including the five elements of the inherited-drift exemption; the current AT-23 status awaits re-check) | [P61 report](at23-restricted-candidate_EN.md): the mechanism legs show 0 drift against the B4 baseline (only one inherited F2 assertion drift), and the real candidate legs actually ran start / execution / stop and recorded unsupported; this does not constitute real-candidate compatibility proof; the historical [T026 report](t026-agent-runner-contract_EN.md) conclusion is unchanged |
```

## 起跑消息五要素 vs 判据语义冲突核查

结论：**无实质冲突**。两处消解说明（不需改判据文本）：
1. 起跑消息 §3.1 后果句写「前提①视为满足（或视为『继承漂移豁免成立』）」的「或」为二选一表述；本稿合并为「前提①视为满足（继承漂移豁免成立）」，语义不变。
2. 起跑消息仅禁止「主控自行判定」；本稿判据文本保持最小表述「⑦（独立终审）判定，主控不得自行判定」，「②~⑥ 各前序环节不得判定」仅写入段 C 报告附注，不入判据文本。
3. 五要素逐项对照起跑消息 §3.1：同基线 / 同断言 / 漂移方向更严格 / 四态不变 / 成因可追溯，名称与锚点一致，无冲突。
