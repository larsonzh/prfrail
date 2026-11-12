# 05-prereview-r2 — P67 GATES-DEBT-CLEANUP（⑤ 预审第 2 轮，逐字归档）

> 归档纪律（§13.1.13 / §11.3）：⑤⑥⑦ 产出逐字原文归档；主控不改写。以下为 ⑤ r2 回执原文。

# P67 GATES-DEBT-CLEANUP · ⑤ r2 预审复核(额度 2/2)

数据来源:读取主控落盘的 `tmp/P67-GATES-DEBT-CLEANUP/review-input/prereview-r2.diff` 与整改落点文件。本角色只读,未执行任何命令。

## 1. F1–F6 逐条核验

- **[F1] CLOSED** — `忌` 的 α 登记已落地：
  - `tools/gates/cjk-newwords.txt:103`:`忌	docs/validation/evidence/P66-CHECKLIST-RULES-HOIST/03-test/step-A-check.md	「表述禁忌」（P66 证据包 step-A-check.md R2-2 实证未登记字；P67 #2 方案 α 按起跑消息 3.1.3 登记）` — TAB 三字段齐全,与其余登记行同构。
  - 报告围栏双侧补齐:`docs/validation/gates-debt-cleanup.md:55` 与 `docs/validation/gates-debt-cleanup_EN.md:55` 均含 `tools/gates/cjk-newwords.txt	repo`。
  - 解析兼容性已核:`tools/gates/lib/ctx.js` 的 `parseWhitelist` 对 `parts[1]` / `parts[2]` 均 `.trim()`、match 取 `parts[0]`,新行格式有效。
- **[F2] CLOSED** — 坑 7 落盘:`docs/validation/gates-debt-cleanup.md:31`:`- 坑 7（§一 滚动窗口歧义）：不新增准则条目，转 v55 事实包登记（v54 快照冻结不改；v55 生成时继承 §13.5 结构模板的滚动窗口口径）。`;`_EN.md:31` 对应英文 bullet 同句。
- **[F3] CLOSED** — §1 新增第 5 条偏差登记:`docs/validation/gates-debt-cleanup.md:21`:`5. **#3 落点偏差登记（⑤ F3 整改）**：design-r1 §4 的②「补录可追溯字段」实现为「各片 ⑦ 盲审事实以各片报告『审查与扫描记录』为唯一真值源、本块不转抄」……弃 C 部分已在 design-r1 §4 如实标注，本区补登记……`;`_EN.md:21` 对应。
- **[F4] CLOSED** — §13.4 来源分组注双侧同步:`docs/DELIVERY_DIRECTIVE.md:1307` 尾「……+ F2-FIX 一条（第 26 条）+ **P67 四条（第 27–30 条）**。」;`docs/DELIVERY_DIRECTIVE_EN.md:1307` 尾「……+ F2-FIX one (item #26) + **P67 four (items #27–30)**.」。
- **[F5] CLOSED** — B.6 标题行双侧整改:`docs/DELIVERY_DIRECTIVE.md:1531` 与 `docs/DELIVERY_DIRECTIVE_EN.md:1531` — 与 G7-f 张力已消解。
- **[F6] CLOSED** — 时点读数口径 bullet:`docs/validation/gates-debt-cleanup.md:32` 与 `_EN.md:32` 对应。

## 2. 新发现

- **LOW | `tools/gates/cjk-newwords.txt:102` | 范围外格式扰动**:本次 diff 同 hunk 将「军」行第二字段改写为 TAB 后带前导空格。该行不在 F1 修法范围,亦无功能必要;与白名单格式约定产生细微偏离。**影响评估**:`parseWhitelist` 对文件字段 `.trim()`,「军」豁免仍生效;白名单文本参与 G4b 缓存键,文本变化仅使旧缓存条目键失配(磁盘缓存 fail-open),不影响正确性。**最小修复**:恢复该行为无前导空格形态,或在报告登记此扰动的理由。

## 3. 结论

**PRE-REVIEW: PASS WITH FIXES**

## 附：主控整改注记（归档时追加，非 ⑤ 原文）

- 新 LOW（「军」行前导空格）已字节级恢复（`git diff --cached` 核验仅剩 +忌 纯插入）；⑤ 额度 2/2 用尽，该字面修正交 ⑦ 确认（S1-REACH 先例口径）。
