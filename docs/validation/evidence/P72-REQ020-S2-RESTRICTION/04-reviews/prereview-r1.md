# P72 `REQ020-S2-RESTRICTION` — ⑤ 预审 r1（deep-reasoner / DeepSeek V4 Pro）

日期 2026-11-14。输入：`diff-r1.txt`（已暂存改动，证据包除外）、`01-arch/design-r1.md`、`r07-refined.md:49`、`req020-s2-restriction{,_EN}.md`、磁盘实文件。结论：**PRE-REVIEW: PASS WITH FIXES**（2 Low + 4 Info；无 Critical/High/Medium）。

## 发现表

| # | 严重度 | 位置 | 问题 | 建议 |
|---|---|---|---|---|
| 1 | Low | `req020-s2-restriction.md` §3；`design-r1.md` §6 判定 5 | 引用节号错误：`s1-scope-annotate` 待办 ③ 实位于 §7 已知边界（`:90`），design 与报告 §3 写作「§5」 | 统一改为「§7 ③」 |
| 2 | Low | `req020-s2-restriction{,_EN}.md` §1 / §5 | 预写未发生事实：§5 已写 ⑤ PASS 并预写 ⑥ 跳过 / ⑦ 终审 PASS | ⑧ 前按实际结果改写 |
| 3 | Info | `PRODUCT_REQUIREMENTS.md:82` / `S1_SUPPORT_MATRIX.md:12` vs `r07-refined.md:49` | 「不冻」枚举较冻结原文少「黄金配置」一词 | 报告 §7 说明差异来源 |
| 4 | Info | `S1_SUPPORT_MATRIX.md:29/31/33` | 该文件有「切片同步」尾注惯例，P72 未追加 | ⑧ 前决定（非硬要求） |
| 5 | Info | `design-r1.md` §2 末行 | 「不动」理由措辞笼统（`DEV_PLAN:32` / `OPERATIONS:121` / `s1-exit:15` 属引用但均为已闭环事实/边界句） | ⑧ 收口时可精确化 |
| 6 | Info | `req020-s2-restriction.md` §2 第 5 行 | 台账 🔄 已在步骤 0 落地 | 表述改为「步骤 0 已落 🔄；⑧ 落 ✅」 |

## 逐项结论

- A: PASS — 全部文本改动与 `design-r1` §3.1 / §3.2 / §3.4 / §3.6 目标逐字一致；无未授权额外改动。
- B: PASS — 6 对文档逐句对称；EN 术语与 `REMAINING_SLICES_EN.md:1993` 既有译本一致；指针 4 行 CN/EN 对称。
- C: PASS — 6 条反例全部不成立（col4 CN 40 / EN 103 原值；REQ 28/28 连续；表内 7 行、`^\| (Python|Node)`=0；状态列不变；保留句在行内）。
- D: PASS — 暂存集合不含准则 / `CONTRACTS` / `SECURITY` / `ARCHITECTURE` / `DEV_PLAN` / `ADR_REGISTER` / `CHECKLIST_v5x` / `s1-exit` / `p57`；MVA 与八门零改动；无 Python/Node 独立行。
- E: PASS — AT-01 保留 C+Go 契约可执行性，仅把「三语言实跑」归 S2；AT-11 判不动合理。
- F: PASS — 报告 §2 / §4 / §7 与磁盘实测一致；§1 / §5 见 #2，§3 节号见 #1。
- G: PASS — BOM/LF 合规；围栏成对；新增文本对 `共 N` / `\d{2,} 行` / run id 形态零命中。

## 判定

PRE-REVIEW: PASS WITH FIXES

## 主控整改记录

- #1：报告 §3 节号改为「§7 ③」（CN/EN）；`design-r1` 作为 ① 证据不改写，改在报告 §7 登记该笔误。
- #2：报告 §1 结论改为「改动面已落地（审查链结论于 ⑧ 定稿）」；§5 改为记录 ⑤ 实际判定、⑥ 待 ⑦ 确认、⑦ 于 ⑧ 定稿。
- #3 / #5 / #4：写入报告 §7（含「不补尾注」的决定与理由）。
- #6：报告 §2 台账行改为「步骤 0 已落 🔄；⑧ 落 ✅ + 完成记录」。
