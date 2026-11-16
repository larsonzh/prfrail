# P76 起跑消息原文（§13.1.10 起跑消息即工件）

切片：P76 `S1-EXIT-CLOSEOUT-ASSESSMENT`（`[SLICE]` · 评估类 M 档）
主控：DeepSeek V4.1 Flash
① 架构：DeepSeek V4 Pro（`deep-reasoner`）
⑤ 预审：GLM-5.3（`glm`）
⑦ 终审：GPT-5.3-Codex（常规档）
批量辅助：GLM-5.3-Flash × 3

---

## 【〇·零】记忆唤醒（开场必读；本段为起草惯例红线，未登入准则 / OB / ADR）

你是 VS Code 主控（DeepSeek V4.1 Flash）。本片是 P76，M 档评估类 `[SLICE]`，片名 = `S1-EXIT-CLOSEOUT-ASSESSMENT`。① 阶段出 S1 Exit 收口路线图后强制停下等用户裁定，再由裁定驱动 ②–⑩ 执行。

**关键锚点（步骤 0 须实测核对）**：
- HEAD / origin = `b90de21`（v59 补正提交）；CI run = `43758070811`（双腿 success）
- 准则 = v1.69；SELFTEST = 737/737；三作用域 = `TOTAL_FAIL=0`
- 滚动窗口 = v57/v58/v59 = 160/148/180
- 当前快照 = **v59**（含本次补正；v59 §二 ⑤ 编号体系已修正为 P75 §11.2 编号 vs P57 §4.3 原编号）

**当前状态（v59 锚点）**：
- P75 已闭环：R-3 = (A) 暂不可起跑；T027-B = 完成（部分不适用）；S-G8 = blocked 至 S2；AT-23 阶段 1 条件闭合 = 永久事实；ADR-045 已登记
- 层 1 T027-B 主线全部出队；层 2 清空
- S1 Exit 八门剩余可推进项 = S-G3 / S-G4 / S-G5 / S-G6（S-G7 合成门；S-G8 blocked）

**P76 的 ②b 阶段将修正 `s1-exit.md:52` 核心任务链行内口径不一致——这是「宣告 S1 Exit 通过」的准入条件，不是 P76 起跑的准入条件。** 未修正前不得宣告 S1 Exit 通过。本片 ②b 负责修正；① 阶段方案书负责给出修正建议（正确口径判定见 Q3）。

**P76 ① 阶段核心产出**：S1 Exit 收口路线图——Q1–Q11 逐项（见【三】）。

**临时件纪律（v59 起强化 D=2 防再犯句）**：「执行角色（`prfrail-implementer` / `prfrail-tester` / `prfrail-documenter`）**与主控**在 `tmp/` 下操作时，一律使用 `tmp/<切片ID>/`；**CI 观察捕获件同理（不得落 `tmp/` 根）**；片内 D 累计达 3 次 ⇒ 强制停片」。本片临时件一律落 `tmp/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/`。

**OB-139 观察中**（第 5 次已发生；用户 2026-11-16 裁定 ① 维持观察中）：本片若再犯（第 6 次）⇒ **转待议 + 立项**（不再停下升级）。

**动态额度预置（`R-ACCUMULATION-MECHANISM`，v59 §五 切片级建议 1）**：若执行中 R 累积到 ≥2 ⇒ **⑦** 按 §9.2 标记高风险 + §7.5 高风险 1+3 额度执行（**⑤ 额度不受 R 累积影响**）。P76 预期 R = 0–1，未触发时按 M 档标准。

## 【〇】级别判定与环境

`[SLICE]`（评估类 M 档；同 P73 / P75 结构）。主控 = `DeepSeek V4.1 Flash (deepseek)`。① 架构 = `DeepSeek V4 Pro`（`deep-reasoner`）；⑤ 预审 = `GLM-5.3 (glm)`（首选，按 §5.2）；⑥ 独立扫描 = 按 §7.2.1 三条件判断（评估类通常可跳过，由 ⑦ 确认）；⑦ 终审 = `GPT-5.3-Codex`（常规档）。**不派 Sol 双打**（非不可测语义 / 断电耐久 / 安全边界重构）。

**不做真实外部调用**：探针额度 = 0；Ubuntu VM / 网络代理 / 真实 CLI 端到端均不启用（R-08.4 不触发）。

**批量辅助**：GLM-5.3-Flash，三次（详见【六】）。

## 【一】目标与背景

**P76 目标**：在 P75 闭环后的 S1 Exit 八门现状基础上，出具 **S1 Exit 收口路线图**——回答如何从「剩余可推进项 = S-G3 / S-G4 / S-G5 / S-G6」到达「S1 Exit 宣布」。

**必须回答的核心问题（Q1–Q11）**：

| # | 问题 | 核心张力 |
|---|---|---|
| **Q1** | S1 Exit 八门（S-G1–S-G8，来自 P57 §4.2）与 s1-reach 十门（G-A–G-I + OB-21，来自 S1-REACH 2026-09-29）的完整映射 | 两套判定闭集、时点、依据不同；重评前必须先建映射 |
| **Q2** | S-G1 / S-G2 的「门级判定待重评」口径——看 AT-23 / AT-24 层成就，还是看 CLI 真实启动 / visible 端到端 | v59 §四 S-G1 判据「CLI 可真实启动」与状态「AT-23 阶段 1 条件闭合」张力 |
| **Q3** | P75 §11.2 唯一残余 `s1-exit.md:52` 核心任务链行内口径不一致的正确口径判定 | 行内「BLOCKED」vs 「PARTIAL」（或其它）哪个正确；决定 S-G7 判定是否变化 |
| **Q4** | S-G7 依据栏在 Q3 修正后的同步更新 | S-G7 是合成门，随其它门自动收口；但 Q3 修正后其判定可能变化 |
| **Q5** | S-G3 解锁路径（最终 ZIP 绑定发行记录 + 复验） | INSTALLATION.md 明写「当前尚无正式 GitHub Release」；需先有 Release 才能 ZIP 绑定 |
| **Q6** | S-G4 解锁路径（ADR-010 批准） | S1_SUPPORT_MATRIX 明写「该提议尚未获批」 |
| **Q7** | S-G5「产品所有者批准」与 `PRODUCT_REQUIREMENTS` §7「独立安全评审结论」的关系 | §7 原文为「**进入 S1 前**」的阻断门禁（S0→S1 门）；可能已通过而并非当前门禁；① 阶段须先核实历史状态再讨论关系 |
| **Q8** | S-G6 门级与任务级不同源（`DEV_PLAN:94` 的 T018 任务级依赖字段含 T027/T028） | 四门清空后 T018 仍不可勾选 |
| **Q9** | `S1_SUPPORT_MATRIX`「CLI Agent / AgentRunner」行口径滞后（仍写「T026/T027/AT-23 阻断 S1」） | P75 后 R-3 = (A)，口径已滞后 |
| **Q10** | 产品所有者批准时机与形式 | 所有可推进门收口后 / S-G8 保持 blocked 前提下的「附条件批准」？ |
| **Q11** | REQ 阶段归属不得越界（REQ-007 黑箱子项推 S2 / REQ-009 推 S2 / REQ-020 S1 仅 C+Go / AT-13 仍在 S1） | 任何「某门未通过」的判定不能把这些推 S2 的项算进 S1 判定 |

**P76 边界**：
- 真实 enforcement 留 S2（S-G8 保持 blocked）
- 不改 `s1-exit{,_EN}` 的其它部分（除 Q3 修正 `:52` 外）
- 不改 `s1-reach-reachability{,_EN}`（P75 A2 明确零改动）
- 不重新执行 AT-23 / AT-24（已条件闭合 / 闭合）
- 不启动真实 CLI 端到端探针（R-3 = A 前提未变）

## 【二】步骤 0（锚点核对 + 建包 + 落盘清单）

**步骤 0a：锚点核对（8 项）**

1. `git rev-parse --short HEAD` 与 `origin/main` —— 期望均 = `b90de21`
2. 准则 `docs/DELIVERY_DIRECTIVE.md` 头部 = `版本：v1.69`
3. `node tools/gates/gate.js --selftest` = `SELFTEST: PASS 737/737`
4. 三作用域 `--scope=index`、`--scope=tree`、`--scope=ci` = 均 `TOTAL_FAIL=0`
5. `tmp` 顶层 = `gates-cache` + `.gitkeep`；排除 `gates-cache` 后递归文件数 = 1
6. `git status --porcelain -uall` 行数 = 0
7. 滚动窗口实测（按 LF 计数，CN / EN 对称）：**v57 = 160、v58 = 148、v59 = 180**
8. 若 ≥2 项不符 ⇒ **停下报实况**，不得静默修正

**步骤 0b：证据包建包**（`docs/validation/evidence/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/`）

三件载荷（步骤 0 落盘；后续按轮追加）：
- `00-briefing.md`：本起跑消息原文（§13.1.10 起跑消息即工件）
- `input-list.md`：本片输入清单（§13.1.1）
- `freedom-list.md`：执行自由度清单（步骤 0 落盘；⑥ 之后、⑦ 之前定稿，§1.7.3）

**步骤 0c：输入清单落盘**（`tmp/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/input-list.md`）

按 §13.1.1 列出本片全部输入（清单见【三】所列文件族）。

**步骤 0 结束条件**：8 锚点相符 + 证据包三件落盘 + 输入清单落盘 + `tmp` 零残留后进入 ①。

## 【三】① 阶段：S1 Exit 收口路线图（核心产出 + 停点）

**① 架构师**（`deep-reasoner` + `DeepSeek V4 Pro`）产出**评估方案书**（落 `tmp/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/design-r1.md`），必须逐项回答 Q1–Q11，每项含：现状核对（精确 `file:line`）、分析（核心张力 + 可选项）、推荐（带理由）、需用户裁定的选项（A/B/C 及后果）。

**① 方案书语言**：单语（中文）。双语报告在 ②c 阶段落 `docs/validation/s1-exit-closeout-assessment{,_EN}.md`。

**① 阶段结束条件**：评估方案书成文（`DESIGN: DONE`）⇒ **强制停下**，出 ① 阶段回执，等用户裁定收口路径（Q1–Q11 的选项）。

**停点纪律**：
- ① 阶段**不 commit**；方案书在 `tmp/` 临时件
- ① 阶段**不改任何文档**（除 `tmp/` 临时件外）
- ① 阶段回执须含：Q1–Q11 逐项现状 + 推荐 + 需用户裁定选项 + 主控对收口顺序的整体建议

## 【四】用户裁定后的 ②–⑩ 流程（预告）

本起跑消息只覆盖步骤 0 + ① 阶段 + 停点；②–⑩ 由续跑消息覆盖（②a 裁定登记 / ②b 文档同步 / ②c 报告成文 / ③ 一致性检查 / ④ 门禁 / ⑤ 预审 / ⑥ 独立扫描 / ⑦ 终审（盲审先于终审）/ ⑧a–⑨–⑧b–⑧c–⑩）。

## 【五】审查层配置

- ⑤ 预审 `GLM-5.3 (glm)`：r1 + 若 r1 报 Medium+ ⇒ 整改后 r2（retryOf = r1）；硬上限 2 次；额度独立于 ⑦ 额度。
- ⑥ 独立扫描：按 §7.2.1 三条件判断（① 只落文档台账 ② 不含 .go/.json/schema/fixtures/脚本语义变更 ③ 不涉状态机/门禁判据/证据模型/角色/授权契约语义）。
- ⑦ 终审 `GPT-5.3-Codex`：M 档 1+2（不含盲审）；盲审另计 +1 且先于终审；新窗口自 P76 起算（P76/P77/P78 不触发按比例抽查，P79 触发）；必抽三语义 = 所有权 / 停机 / 身份，P76 预期不触。⑦ 总额度上限 4（R≥2 升 §9.2 高风险 1+3 ⇒ 上限 5）。

## 【六】批量辅助配置（GLM-5.3-Flash）

| # | 任务 | 产出 |
|---|---|---|
| **批量 1** | `s1-exit{,_EN}` + `p57-closeout-assessment{,_EN}` §4.2 的八门现状枚举 | `batch-01-s1-exit-p57.md` |
| **批量 2** | `s1-reach-reachability{,_EN}` §2 十门清单枚举 | `batch-02-s1-reach.md` |
| **批量 3** | `PRODUCT_REQUIREMENTS{,_EN}` §7 + `S1_SUPPORT_MATRIX{,_EN}` + `INSTALLATION{,_EN}` 的滞后口径枚举 + `docs/validation/` 独立安全评审检索 | `batch-03-lagging-docs.md` |

批量辅助只做事实枚举，不做判定 / 推荐 / 排序；产出交 ① 架构师作为输入；不计入 ⑦ 额度。

## 【七】门禁与边界

门禁：步骤 0 前后、① 完成后、②b 后、commit 前三作用域 + SELFTEST。

边界（不做清单）：不改准则（维持 v1.69、不升版）；不改 CONTRACTS / ARCHITECTURE / BUSINESS_WORKFLOWS / CODING_CONVENTIONS / DOCUMENTATION_PLAN / OPERATIONS / RFC-*；不改 `CHECKLIST_v59{,_EN}`（v59 冻结）与 v50–v58 快照；不改 `s1-reach-reachability{,_EN}`；不重跑 AT-23 / AT-24；不启真实 CLI 端到端探针；不改 `.go`/`.json`/`.yml`/schema/fixtures/CI/`go.mod`；不新增 OB（除非用户裁定）；不 commit / push 前先停在待授权。

## 【八】文档同步更新清单（分层）

第一层（确定会触）：`s1-exit{,_EN}` 修 `:52`；`REMAINING_SLICES{,_EN}` P76 完成记录（⑧b）；证据包追加；`s1-exit-closeout-assessment{,_EN}` 报告新建。

第二层（视裁定）：`p57` / `DEV_PLAN` / `S1_SUPPORT_MATRIX` / `PRODUCT_REQUIREMENTS` / `INSTALLATION*` / `ADR_REGISTER` / `SECURITY` / `S1_RELEASE_NOTES` / `TEST_STRATEGY`。

第三层（确定不触）：准则 / CONTRACTS / ARCHITECTURE / BUSINESS_WORKFLOWS / CODING_CONVENTIONS / DOCUMENTATION_PLAN / OPERATIONS / RFC-* / CHECKLIST_v59 / v50–v58 快照 / s1-reach / at23 / at24 / r08-r09 / 代码与 CI。

## 【九】授权（分两批）

本起跑消息只授权第一批：步骤 0 + ① 阶段 + 停点。① 阶段结束时**不 commit**（方案书在 `tmp/`）。**无 commit / push 授权**。

第二批（续跑消息）覆盖 ②–⑩ + commit + push（仅 origin，不推 gitee）。

本片不触发 §1.7.4 临时授权。探针额度 = 0。

## 【十】① 阶段回执要求

1. 步骤 0 的 8 项锚点实测（含滚动窗口 160 / 148 / 180）
2. 证据包三件落盘确认
3. 批量辅助三次的产出摘要（每次字符数 + 覆盖范围 + 落盘路径；批量 3 须含独立安全评审检索结论）
4. ① 阶段方案书成文（`DESIGN: DONE`）：Q1–Q11 逐项现状 + 分析 + 推荐 + 需裁定选项；Q7 历史状态核实结论；收口顺序整体建议；⑦ 盲审是否必抽的判定；⑥ 是否可跳过的判定
5. P75 EN §11.2 术语一致性 + v59 EN 内部术语一致性的核实结论
6. `tmp` 零残留 + porcelain 0（① 阶段不 commit，工作区应干净）
7. 给续跑消息的锚点：HEAD / CI run / 临时件路径

## 【十一】执行节奏与红线

1. 全程在 ① 阶段结束后强制停下。
2. 触发 §1.7.2 停升级边界 ⇒ 停下报实况，不得自裁。
3. 临时件一律 `tmp/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/`；只增不删；① 阶段结束不清理。
4. 不改准则版本（维持 v1.69）。
5. ① 阶段方案书单语（中文）；最终报告 ②c 双语。
6. v59 冻结不改写。
7. ① 阶段不改 `s1-exit` / `s1-reach`（只给修正建议）；②b 改 `s1-exit:52`。
8. 动态额度预置：R≥2 ⇒ ⑦ 按 §9.2 高风险 1+3；⑤ 额度不受 R 累积影响。
9. G6-4 自查：避免 `\d{2,}\s*行` 类字面。
10. E-1 流程纪律：① 阶段结束前二次回读 `tmp/` 五条读数。
11. D=2 防再犯句适用：`tmp/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/`；片内 D 累计达 3 次强制停片。
12. OB-139 观察中：本片再犯（第 6 次）⇒ 转待议 + 立项。
