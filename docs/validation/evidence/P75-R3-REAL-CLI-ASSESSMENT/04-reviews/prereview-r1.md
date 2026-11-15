# prereview-r1 — ⑤ 预审第 1 轮原文（逐字归档）

载体：`deep-reasoner` + `DeepSeek V4 Pro (deepseek)`。切片：P75 R3-REAL-CLI-ASSESSMENT。实得字符数：约 2170。原文如下。

---

PRE-REVIEW: PASS WITH FIXES

**核验结论（必核 1–7）**
- 必核 1：对比表与正文无矛盾；(A)/(B)/(C) 成本/风险四类下游/可回滚性均有，但缺「不做什么」（见 #5）。
- 必核 2：§4 四要素齐备（裁决人=用户/2026-11-15/理由/依据）。
- 必核 3：九处引用全部命中——p57 §3.3/§4.2、`s1-exit.md:43`、at23 `:19-20/:423/:499`、t027b 裁定 5(`:43`)/7(`:45`)、sandbox `:133`、B1 `:4/:111`、B2 `:12/:18-19/:162-164`、r08-r09 `:28-57`。
- 必核 4：`DEV_PLAN.md:203`、`p57:76/:99`、`ADR_REGISTER.md:56`(ADR-045) 属实；p57 仅 +4 行加法式指针，S-G1 行零改动。
- 必核 5：§11.2 逐项三项+§11.3 准入条件与 p57 §4.3 一致（但内容过期，见 #1–#3）。
- 必核 6：变更集 = 围栏 22 文件，无 CONTRACTS/schema/fixtures/CI/产品代码改动；`s1-exit`、`s1-reach-reachability` 零改动。
- 必核 7：围栏存在，22 行均为「路径+TAB+处置」；无过期「数字+行」字面量、无占位符残留。

**发现**

1. `Medium | docs/validation/r3-real-cli-assessment.md:155 | 「\`docs/validation/s1-exit.md:44\` 的 AT-24 所在行写「T028 尚未实现」，与 \`docs/DEV_PLAN.md:119\` 的 \`COMPLETED\` 互斥」 | 磁盘现况：s1-exit.md:44 已是「条件闭合（2026-11-11，P64 验收…T028 已完成）」（提交 2dbc54a/3884351），滞后已消除。最小修复：重核后改写为「已更正，本项义务作废」，同步 §11.3。`
2. `Medium | 同文件:157 | 「\`docs/validation/s1-reach-reachability.md\` §2 的 G-D 行把 T028 记为 \`PLANNED / NOT IMPLEMENTED\`，且其 \`DEV_PLAN\` 引用已漂移（\`:91\` / \`:93\` → 实际 \`:94\` / \`:96\`）」 | 磁盘现况：该文件 G-D 行(:25)已写「DEV_PLAN.md:119（T028 COMPLETED…）；s1-exit.md:44（已更正为「条件闭合」）」，无漂移。最小修复：改为「已更正，义务作废」。`
3. `Medium | 同文件:156 | 「\`docs/validation/s1-exit.md:52\` 的核心任务链行同型时态滞后（T028 部分已滞后）」 | 磁盘现况：:52 已写「T028 已完成、AT-24 已验收 = 条件闭合」，T028 不再滞后；残余滞后为「AT-23 未通过」「T027 未完成」（与 :43 及 A1 口径矛盾）。最小修复：改写为当前实际滞后点。`
   （注：三项与 p57 §4.3 逐字一致，但 p57 §4.3 系 2026-11-05 时点登记、已过期；P75 原样搬运未重核。）
4. `Medium | 同文件:106、:149 | 「（落点见 §13）」「新立 ADR-045，见 §13」 | 报告无 §12/§13（止于 §11.3），悬空引用。最小修复：改为「见 \`docs/ADR_REGISTER.md:56\`」。`
5. `Medium | 同文件:53-55、:57-64 | 三选项正文与对比表均无「不做什么」，表缺「做什么/不做什么」轴，违反 freedom-list.md F-3「做什么/不做什么/成本/风险/可回滚性」五轴承诺 | 最小修复：按 design-r1 §2.1 补「不做」列/行。`
6. `Medium | evidence/P75-R3-REAL-CLI-ASSESSMENT/06-gates/gate-tree.txt:20、:41、:52 | 「FAIL G5(d) staged is a subset of declared … stagedNotDeclared gate-index/tree」「FAIL G8-b … not listed in 归档条目清单/SHA256SUMS」「TOTAL_FAIL=2」 | 与 00-briefing「三作用域 TOTAL_FAIL=0」矛盾；两 FAIL 为声明前时点快照（现 fence/MANIFEST/SHA256SUMS 已含四文件），但归档门禁输出未刷新、包内无重跑证据。最小修复：以定稿后重跑结果替换 gate-tree/gate-index 并重算 SHA256SUMS，或显式登记两 FAIL 及解决证据。`
7. `Medium | 同报告:75 | 「依据 | 用户 2026-11-15 同轮裁定原文（含五项附带裁定 A1–A5 与 O-7 = ② 的裁定）」 | 证据包未归档裁定原文（00-briefing 仅起跑消息），「可回溯到用户原话」磁盘上不可验。最小修复：落盘裁定消息原文（如 02-batch/）并在依据行引用。`
8. `Low | 同报告:3 | 「状态：② 阶段成文（① 已停点裁定；②a–②c 分步落地）」 | ④ 门禁已归档、现处 ⑤ 预审，状态行过期。最小修复：更新状态行。`
9. `Low | docs/t027/REMAINING_SLICES.md:2074 | 「现状缺口 | R-3 未裁定：…」 | R-3 已裁定 (A)，字段未同步（完成记录按模板待 ⑧）。最小修复：⑧ 完成记录时同步。`

**无法验证的部分**
- 用户 R-3/A1–A5/O-7 裁定原文与会话原话的逐字一致性（磁盘无原文归档）。
- 最后一次 tree 门禁重跑的真实结果（终端筛选输出不可见；磁盘 gate-tree.txt 仍为 FAIL 快照）。
- ⑥ 是否可跳过（属 ⑦ 终审判定）；`_EN` 译文语义等价（门禁仅结构对称校验）。

实得字符数：约 2170。
